// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/hyperchain/go-hpc-msp/plugin/software"
	"github.com/hyperchain/go-hpc-msp/tls/internal"
	"github.com/meshplus/crypto"
)

// serverHandshakeState contains details of a server handshake in progress.
// It's discarded once the handshake has completed.
type serverHandshakeState struct {
	c                     *Conn
	clientHello           *clientHelloMsg
	hello                 *serverHelloMsg
	suite                 *cipherSuite
	ellipticOk            bool
	ecdsaOk               bool
	sm2OK                 bool
	rsaDecryptOk          bool
	rsaSignOk             bool
	sessionState          *sessionState
	finishedHash          finishedHash
	masterSecret          []byte
	certsFromClient       [][]byte
	cert                  *Certificate
	cachedClientHelloInfo *ClientHelloInfo
}

// serverHandshake performs a TLS handshake as a server.
func (c *Conn) serverHandshake() error {
	// If this is the first server handshake, we generate a random key to
	// encrypt the tickets with.
	c.config.serverInitOnce.Do(func() { c.config.serverInit(nil) })

	hs := serverHandshakeState{
		c: c,
	}
	isResume, err := hs.readClientHello()
	if err != nil {
		return err
	}

	// For an overview of TLS handshaking, see https://tools.ietf.org/html/rfc5246#section-7.3
	c.buffering = true
	if isResume {
		// The client has included a session ticket and so we do an abbreviated handshake.
		if err := hs.doResumeHandshake(); err != nil {
			return err
		}
		if err := hs.establishKeys(); err != nil {
			return err
		}
		// ticketSupported is set in a resumption handshake if the
		// ticket from the client was encrypted with an old session
		// ticket key and thus a refreshed ticket should be sent.
		if hs.hello.ticketSupported {
			if err := hs.sendSessionTicket(); err != nil {
				return err
			}
		}
		if err := hs.sendFinished(c.serverFinished[:]); err != nil {
			return err
		}
		if _, err := c.flush(); err != nil {
			return err
		}
		c.clientFinishedIsFirst = false
		if err := hs.readFinished(nil); err != nil {
			return err
		}
		c.didResume = true
	} else {
		// The client didn't include a session ticket, or it wasn't
		// valid so we do a full handshake.
		if err := hs.doFullHandshake(); err != nil {
			return err
		}
		if err := hs.establishKeys(); err != nil {
			return err
		}
		if err := hs.readFinished(c.clientFinished[:]); err != nil {
			return err
		}
		c.clientFinishedIsFirst = true
		c.buffering = true
		if err := hs.sendSessionTicket(); err != nil {
			return err
		}
		if err := hs.sendFinished(nil); err != nil {
			return err
		}
		if _, err := c.flush(); err != nil {
			return err
		}
	}

	c.ekm = ekmFromMasterSecret(c.vers, hs.suite, hs.masterSecret, hs.clientHello.random, hs.hello.random)
	atomic.StoreUint32(&c.handshakeStatus, 1)

	return nil
}

// readClientHello reads a ClientHello message from the client and decides
// whether we will perform session resumption.
func (hs *serverHandshakeState) readClientHello() (isResume bool, err error) {
	c := hs.c

	msg, err := c.readHandshake()
	if err != nil {
		return false, err
	}
	var ok bool
	hs.clientHello, ok = msg.(*clientHelloMsg)
	if !ok {
		_ = c.sendAlert(alertUnexpectedMessage)
		return false, unexpectedMessageError(hs.clientHello, msg)
	}

	if c.config.GetConfigForClient != nil {
		if newConfig, gerr := c.config.GetConfigForClient(hs.clientHelloInfo()); gerr != nil {
			_ = c.sendAlert(alertInternalError)
			return false, gerr
		} else if newConfig != nil {
			newConfig.serverInitOnce.Do(func() { newConfig.serverInit(c.config) })
			c.config = newConfig
		}
	}

	c.vers, ok = c.config.mutualVersion(hs.clientHello.vers)
	if !ok {
		_ = c.sendAlert(alertProtocolVersion)
		return false, fmt.Errorf("tls: client offered an unsupported, maximum protocol version of %x", hs.clientHello.vers)
	}
	c.haveVers = true

	hs.hello = new(serverHelloMsg)

	supportedCurve := false
	preferredCurves := c.config.curvePreferences()
Curves:
	for _, curve := range hs.clientHello.supportedCurves {
		for _, supported := range preferredCurves {
			if supported == curve {
				supportedCurve = true
				break Curves
			}
		}
	}

	supportedPointFormat := false
	for _, pointFormat := range hs.clientHello.supportedPoints {
		if pointFormat == pointFormatUncompressed {
			supportedPointFormat = true
			break
		}
	}
	hs.ellipticOk = supportedCurve && supportedPointFormat

	foundCompression := false
	// We only support null compression, so check that the client offered it.
	for _, compression := range hs.clientHello.compressionMethods {
		if compression == compressionNone {
			foundCompression = true
			break
		}
	}

	if !foundCompression {
		_ = c.sendAlert(alertHandshakeFailure)
		return false, errors.New("tls: client does not support uncompressed connections")
	}

	hs.hello.vers = c.vers
	hs.hello.random = make([]byte, 32)
	_, err = io.ReadFull(c.config.rand(), hs.hello.random)
	if err != nil {
		_ = c.sendAlert(alertInternalError)
		return false, err
	}

	if len(hs.clientHello.secureRenegotiation) != 0 {
		_ = c.sendAlert(alertHandshakeFailure)
		return false, errors.New("tls: initial handshake had non-empty renegotiation extension")
	}

	hs.hello.secureRenegotiationSupported = hs.clientHello.secureRenegotiationSupported
	hs.hello.compressionMethod = compressionNone
	if len(hs.clientHello.serverName) > 0 {
		c.serverName = hs.clientHello.serverName
	}

	if len(hs.clientHello.alpnProtocols) > 0 {
		if selectedProto, fallback := mutualProtocol(hs.clientHello.alpnProtocols, c.config.NextProtos); !fallback {
			hs.hello.alpnProtocol = selectedProto
			c.clientProtocol = selectedProto
		}
	} else {
		// Although sending an empty NPN extension is reasonable, Firefox has
		// had a bug around this. Best to send nothing at all if
		// c.config.NextProtos is empty. See
		// https://golang.org/issue/5445.
		if hs.clientHello.nextProtoNeg && len(c.config.NextProtos) > 0 {
			hs.hello.nextProtoNeg = true
			hs.hello.nextProtos = c.config.NextProtos
		}
	}

	hs.cert, err = c.config.getCertificate(hs.clientHelloInfo())
	if err != nil {
		_ = c.sendAlert(alertInternalError)
		return false, err
	}
	if hs.clientHello.scts {
		hs.hello.scts = hs.cert.SignedCertificateTimestamps
	}

	hs.sm2OK = hs.cert.PrivateKey.GetKeyInfo() == crypto.Sm2p256v1
	hs.ecdsaOk = hs.cert.PrivateKey.GetKeyInfo() >= crypto.Secp256k1 && hs.cert.PrivateKey.GetKeyInfo() <= crypto.Secp521r1
	hs.rsaSignOk = hs.cert.PrivateKey.GetKeyInfo() >= crypto.Rsa2048 && hs.cert.PrivateKey.GetKeyInfo() <= crypto.Rsa4096
	hs.rsaDecryptOk = hs.rsaSignOk

	if hs.checkForResumption() {
		return true, nil
	}

	var preferenceList, supportedList []uint16
	if c.config.PreferServerCipherSuites {
		preferenceList = c.config.cipherSuites()
		supportedList = hs.clientHello.cipherSuites
	} else {
		preferenceList = hs.clientHello.cipherSuites
		supportedList = c.config.cipherSuites()
	}

	for _, id := range preferenceList {
		if hs.setCipherSuite(id, supportedList, c.vers) {
			//fmt.Printf("###???????????????suit:%x\n", id)
			break
		}
	}

	if hs.suite == nil {
		_ = c.sendAlert(alertHandshakeFailure)
		return false, errors.New("tls: no cipher suite supported by both client and server")
	}

	// See https://tools.ietf.org/html/rfc7507.
	for _, id := range hs.clientHello.cipherSuites {
		if id == TLS_FALLBACK_SCSV {
			// The client is doing a fallback connection.
			if hs.clientHello.vers < c.config.maxVersion() {
				_ = c.sendAlert(alertInappropriateFallback)
				return false, errors.New("tls: client using inappropriate protocol fallback")
			}
			break
		}
	}

	return false, nil
}

// checkForResumption reports whether we should perform resumption on this connection.
func (hs *serverHandshakeState) checkForResumption() bool {
	c := hs.c

	if c.config.SessionTicketsDisabled {
		return false
	}

	var ok bool
	var sessionTicket = append([]uint8{}, hs.clientHello.sessionTicket...)
	if hs.sessionState, ok = c.decryptTicket(sessionTicket); !ok {
		return false
	}

	// Never resume a session for a different TLS version.
	if c.vers != hs.sessionState.vers {
		return false
	}

	cipherSuiteOk := false
	// Check that the client is still offering the ciphersuite in the session.
	for _, id := range hs.clientHello.cipherSuites {
		if id == hs.sessionState.cipherSuite {
			cipherSuiteOk = true
			break
		}
	}
	if !cipherSuiteOk {
		return false
	}

	// Check that we also support the ciphersuite from the session.
	if !hs.setCipherSuite(hs.sessionState.cipherSuite, c.config.cipherSuites(), hs.sessionState.vers) {
		return false
	}

	sessionHasClientCerts := len(hs.sessionState.certificates) != 0
	needClientCerts := c.config.ClientAuth == RequireAnyClientCert || c.config.ClientAuth == RequireAndVerifyClientCert
	if needClientCerts && !sessionHasClientCerts {
		return false
	}
	if sessionHasClientCerts && c.config.ClientAuth == NoClientCert {
		return false
	}

	return true
}

func (hs *serverHandshakeState) doResumeHandshake() error {
	c := hs.c

	hs.hello.cipherSuite = hs.suite.id
	// We echo the client's session ID in the ServerHello to let it know
	// that we're doing a resumption.
	hs.hello.sessionID = hs.clientHello.sessionID
	hs.hello.ticketSupported = hs.sessionState.usedOldKey
	hs.finishedHash = newFinishedHash(c.vers, hs.suite)
	hs.finishedHash.discardHandshakeBuffer()
	_, _ = hs.finishedHash.Write(hs.clientHello.marshal())
	_, _ = hs.finishedHash.Write(hs.hello.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, hs.hello.marshal()); err != nil {
		return err
	}

	if len(hs.sessionState.certificates) > 0 {
		if _, err := hs.processCertsFromClient(hs.sessionState.certificates); err != nil {
			return err
		}
	}

	hs.masterSecret = hs.sessionState.masterSecret

	return nil
}

func (hs *serverHandshakeState) doFullHandshake() error {
	c := hs.c
	if c.handshakeErr != nil {
		Errorf("doFullHandshake: handshakeErr: %v, isClient: %v ,"+
			"hs.cipherSuit.id: %x",
			c.handshakeErr, c.isClient, hs.suite.id,
		)
	}

	if hs.clientHello.ocspStapling && len(hs.cert.OCSPStaple) > 0 {
		hs.hello.ocspStapling = true
	}

	hs.hello.ticketSupported = hs.clientHello.ticketSupported && !c.config.SessionTicketsDisabled
	hs.hello.cipherSuite = hs.suite.id

	hs.finishedHash = newFinishedHash(hs.c.vers, hs.suite)
	if c.config.ClientAuth == NoClientCert {
		// No need to keep a full record of the handshake if client
		// certificates won't be used.
		hs.finishedHash.discardHandshakeBuffer()
	}
	_, _ = hs.finishedHash.Write(hs.clientHello.marshal())
	_, _ = hs.finishedHash.Write(hs.hello.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, hs.hello.marshal()); err != nil {
		return err
	}

	certMsg := new(certificateMsg)
	certMsg.certificates = hs.cert.Certificate
	_, _ = hs.finishedHash.Write(certMsg.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, certMsg.marshal()); err != nil {
		return err
	}

	if hs.hello.ocspStapling {
		certStatus := new(certificateStatusMsg)
		certStatus.statusType = statusTypeOCSP
		certStatus.response = hs.cert.OCSPStaple
		_, _ = hs.finishedHash.Write(certStatus.marshal())
		if _, err := c.writeRecord(recordTypeHandshake, certStatus.marshal()); err != nil {
			return err
		}
	}

	//KeyAgreement
	keyAgreement := hs.suite.ka(c.vers)
	skx, err := keyAgreement.generateServerKeyExchange(c.config, hs.cert, hs.clientHello, hs.hello)
	if err != nil {
		_ = c.sendAlert(alertHandshakeFailure)
		return err
	}
	if skx != nil {
		_, _ = hs.finishedHash.Write(skx.marshal())
		if _, werr := c.writeRecord(recordTypeHandshake, skx.marshal()); werr != nil {
			return werr
		}
	}

	if c.config.ClientAuth >= RequestClientCert {
		// Request a client certificate
		certReq := new(certificateRequestMsg)
		certReq.certificateTypes = []byte{
			byte(certTypeRSASign),
			byte(certTypeECDSASign),
		}
		if c.vers >= VersionTLS12 {
			certReq.hasSignatureAndHash = true
			certReq.supportedSignatureAlgorithms = supportedSignatureAlgorithms
		}

		// An empty list of certificateAuthorities signals to
		// the client that it may send any certificate in response
		// to our request. When we know the CAs we trust, then
		// we can send them down, so that the client can choose
		// an appropriate certificate to give to us.
		if c.config.ClientCAs != nil {
			certReq.certificateAuthorities = c.config.ClientCAs.Subjects()
		}
		_, _ = hs.finishedHash.Write(certReq.marshal())
		if _, werr := c.writeRecord(recordTypeHandshake, certReq.marshal()); werr != nil {
			return werr
		}
	}

	helloDone := new(serverHelloDoneMsg)
	_, _ = hs.finishedHash.Write(helloDone.marshal())
	if _, werr := c.writeRecord(recordTypeHandshake, helloDone.marshal()); werr != nil {
		return werr
	}

	if _, ferr := c.flush(); ferr != nil {
		return ferr
	}

	var pub crypto.VerifyKey // public key for client auth, if any

	msg, err := c.readHandshake()
	if err != nil {
		return err
	}

	var ok bool
	// If we requested a client certificate, then the client must send a
	// certificate message, even if it's empty.
	if c.config.ClientAuth >= RequestClientCert {
		if certMsg, ok = msg.(*certificateMsg); !ok {
			_ = c.sendAlert(alertUnexpectedMessage)
			return unexpectedMessageError(certMsg, msg)
		}
		_, _ = hs.finishedHash.Write(certMsg.marshal())

		if len(certMsg.certificates) == 0 {
			// The client didn't actually send a certificate
			switch c.config.ClientAuth {
			case RequireAnyClientCert, RequireAndVerifyClientCert:
				_ = c.sendAlert(alertBadCertificate)
				Errorf("tls: client didn't provide a certificate")
				return errors.New("tls: client didn't provide a certificate")
			}
		}

		pub, err = hs.processCertsFromClient(certMsg.certificates)
		if err != nil {
			return err
		}

		msg, err = c.readHandshake()
		if err != nil {
			return err
		}
	}

	// Get client key exchange
	ckx, ok := msg.(*clientKeyExchangeMsg)
	if !ok {
		_ = c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(ckx, msg)
	}
	_, _ = hs.finishedHash.Write(ckx.marshal())

	preMasterSecret, err := keyAgreement.processClientKeyExchange(c.config, hs.cert, ckx, c.vers)
	if err != nil {
		_ = c.sendAlert(alertHandshakeFailure)
		return err
	}
	hs.masterSecret = masterFromPreMasterSecret(c.vers, hs.suite, preMasterSecret, hs.clientHello.random, hs.hello.random)
	if werr := c.config.writeKeyLog(hs.clientHello.random, hs.masterSecret); werr != nil {
		_ = c.sendAlert(alertInternalError)
		return werr
	}

	// If we received a client cert in response to our certificate request message,
	// the client will send us a certificateVerifyMsg immediately after the
	// clientKeyExchangeMsg. This message is a digest of all preceding
	// handshake-layer messages that is signed using the private key corresponding
	// to the client's certificate. This allows us to verify that the client is in
	// possession of the private key of the certificate.
	if len(c.peerCertificates) > 0 {
		msg, err = c.readHandshake()
		if err != nil {
			return err
		}
		certVerify, ok := msg.(*certificateVerifyMsg)
		if !ok {
			_ = c.sendAlert(alertUnexpectedMessage)
			return unexpectedMessageError(certVerify, msg)
		}

		// Determine the signature type.
		sigAlgo, sigType, hashFunc, perr := pickSignatureAlgorithm(pub, []SignatureScheme{certVerify.signatureAlgorithm}, supportedSignatureAlgorithms, c.vers)
		if perr != nil {
			Errorf("pickSignatureAlgorithm: pick signature algorithm failed")
			_ = c.sendAlert(alertIllegalParameter)
			return perr
		}

		var digest []byte
		engine := software.GetSoftwareEngine()
		if sigAlgo == SM2WithSM3 && hs.finishedHash.version >= VersionTLS12 {
			pass := pub.Verify(hs.finishedHash.buffer, internal.New(), certVerify.signature)
			if !pass {
				_ = c.sendAlert(alertBadCertificate)
				Errorf("tls: could not validate signature of connection nonces: validate signature failed")
				return errors.New("tls: could not validate signature of connection nonces: validate signature failed")
			}
		} else {
			digest, perr = hs.finishedHash.hashForClientCertificate(sigType, hashFunc, hs.masterSecret)
			if perr == nil {
				if sigAlgo == SM2WithSM3 {
					if pub.GetKeyInfo() != crypto.Sm2p256v1 {
						_ = c.sendAlert(alertInternalError)
						return err
					}
					h, _ := engine.GetHash(crypto.Sm3WithPublicKey)
					digest, _ = h.BatchHash([][]byte{pub.Bytes(), digest})
					pass := pub.Verify(digest, internal.Copy(), certVerify.signature)
					if !pass {
						_ = c.sendAlert(alertBadCertificate)
						Errorf("tls: could not validate signature of connection nonces: validate signature failed")
						return errors.New("tls: could not validate signature of connection nonces: validate signature failed")
					}
				} else {
					if sigType == signaturePKCS1v15 {
						digest = append(digest, make([]byte, 4)...)
						binary.BigEndian.PutUint32(digest[len(digest)-4:], uint32(hashFunc))
					}
					perr = verifyHandshakeSignature(sigType, pub, hashFunc, digest, certVerify.signature)
					if perr != nil {
						_ = c.sendAlert(alertBadCertificate)
						Errorf("tls: could not validate signature of connection nonces: " + perr.Error())
						return errors.New("tls: could not validate signature of connection nonces: " + perr.Error())
					}
				}
			}
		}

		_, _ = hs.finishedHash.Write(certVerify.marshal())
	}

	hs.finishedHash.discardHandshakeBuffer()

	return nil
}

func (hs *serverHandshakeState) establishKeys() error {
	c := hs.c

	clientMAC, serverMAC, clientKey, serverKey, clientIV, serverIV :=
		keysFromMasterSecret(c.vers, hs.suite, hs.masterSecret, hs.clientHello.random, hs.hello.random, hs.suite.macLen, hs.suite.keyLen, hs.suite.ivLen)

	var clientCipher, serverCipher interface{}
	var clientHash, serverHash macFunction

	if hs.suite.aead == nil {
		clientCipher = hs.suite.cipher(clientKey, clientIV, true /* for reading */)
		clientHash = hs.suite.mac(c.vers, clientMAC)
		serverCipher = hs.suite.cipher(serverKey, serverIV, false /* not for reading */)
		serverHash = hs.suite.mac(c.vers, serverMAC)
	} else {
		clientCipher = hs.suite.aead(clientKey, clientIV)
		serverCipher = hs.suite.aead(serverKey, serverIV)
	}

	c.in.prepareCipherSpec(c.vers, clientCipher, clientHash)
	c.out.prepareCipherSpec(c.vers, serverCipher, serverHash)

	return nil
}

func (hs *serverHandshakeState) readFinished(out []byte) error {
	c := hs.c

	_ = c.readRecord(recordTypeChangeCipherSpec)
	if c.in.err != nil {
		return c.in.err
	}

	if hs.hello.nextProtoNeg {
		msg, err := c.readHandshake()
		if err != nil {
			return err
		}
		nextProto, ok := msg.(*nextProtoMsg)
		if !ok {
			_ = c.sendAlert(alertUnexpectedMessage)
			return unexpectedMessageError(nextProto, msg)
		}
		_, _ = hs.finishedHash.Write(nextProto.marshal())
		c.clientProtocol = nextProto.proto
	}

	msg, err := c.readHandshake()
	if err != nil {
		return err
	}
	clientFinished, ok := msg.(*finishedMsg)
	if !ok {
		_ = c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(clientFinished, msg)
	}

	verify := hs.finishedHash.clientSum(hs.masterSecret)
	if len(verify) != len(clientFinished.verifyData) ||
		subtle.ConstantTimeCompare(verify, clientFinished.verifyData) != 1 {
		_ = c.sendAlert(alertHandshakeFailure)
		return errors.New("tls: client's Finished message is incorrect")
	}

	_, _ = hs.finishedHash.Write(clientFinished.marshal())
	copy(out, verify)
	return nil
}

func (hs *serverHandshakeState) sendSessionTicket() error {
	if !hs.hello.ticketSupported {
		return nil
	}

	c := hs.c
	m := new(newSessionTicketMsg)

	var err error
	state := sessionState{
		vers:         c.vers,
		cipherSuite:  hs.suite.id,
		masterSecret: hs.masterSecret,
		certificates: hs.certsFromClient,
	}
	m.ticket, err = c.encryptTicket(&state)
	if err != nil {
		return err
	}

	_, _ = hs.finishedHash.Write(m.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, m.marshal()); err != nil {
		return err
	}

	return nil
}

func (hs *serverHandshakeState) sendFinished(out []byte) error {
	c := hs.c

	if _, err := c.writeRecord(recordTypeChangeCipherSpec, []byte{1}); err != nil {
		return err
	}

	finished := new(finishedMsg)
	finished.verifyData = hs.finishedHash.serverSum(hs.masterSecret)
	_, _ = hs.finishedHash.Write(finished.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, finished.marshal()); err != nil {
		return err
	}

	c.cipherSuite = hs.suite.id
	copy(out, finished.verifyData)

	return nil
}

// processCertsFromClient takes a chain of client certificates either from a
// Certificates message or from a sessionState and verifies them. It returns
// the public key of the leaf certificate.
func (hs *serverHandshakeState) processCertsFromClient(certificates [][]byte) (crypto.VerifyKey, error) {
	c := hs.c

	hs.certsFromClient = certificates
	certs := make([]*software.Certificate, len(certificates))
	var err error
	for i, asn1Data := range certificates {
		if certs[i], err = software.ParseCertificate(string(asn1Data)); err != nil {
			Warningf("server: get certificate(index:%v) from client: %v", i, base64.StdEncoding.EncodeToString(asn1Data))
			_ = c.sendAlert(alertBadCertificate)
			Errorf("tls: failed to parse client certificate: " + err.Error())
			return nil, errors.New("tls: failed to parse client certificate: " + err.Error())
		}
	}

	if c.config.ClientAuth >= VerifyClientCertIfGiven && len(certs) > 0 {
		opts := software.VerifyOptions{
			Roots:         c.config.ClientCAs,
			CurrentTime:   c.config.time(),
			Intermediates: software.NewCertPool(),
			KeyUsages:     []software.ExtKeyUsage{software.ExtKeyUsageClientAuth},
		}

		for _, cert := range certs[1:] {
			opts.Intermediates.AddCert(cert)
		}

		chains, err := certs[0].Verify(opts)
		if err != nil {
			_ = c.sendAlert(alertBadCertificate)
			Errorf("tls: failed to verify client's certificate: %v, client cert: \n%v", err.Error(),
				base64.StdEncoding.EncodeToString(certificates[0]))
			Errorf("server root ca congfig: %v", opts.Roots.PrintDebugInfo())
			Errorf("server intermediates ca congfig: %v", opts.Intermediates.PrintDebugInfo())
			return nil, errors.New("tls: failed to verify client's certificate: " + err.Error())
		}

		c.verifiedChains = chains
	}

	if c.config.VerifyPeerCertificate != nil {
		if err := c.config.VerifyPeerCertificate(certificates, c.verifiedChains); err != nil {
			_ = c.sendAlert(alertBadCertificate)
			Errorf("bad certificate: %v", err.Error())
			return nil, err
		}
	}

	if len(certs) == 0 {
		return nil, nil
	}

	if certs[0].PublicKey.GetKeyInfo() < crypto.Sm2p256v1 || certs[0].PublicKey.GetKeyInfo() > crypto.Rsa4096 {
		_ = c.sendAlert(alertUnsupportedCertificate)
		return nil, fmt.Errorf("tls: client's certificate contains an unsupported public key of type %T", certs[0].PublicKey)
	}

	c.peerCertificates = certs
	return certs[0].PublicKey, nil
}

// setCipherSuite sets a cipherSuite with the given id as the serverHandshakeState
// suite if that cipher suite is acceptable to use.
// It returns a bool indicating if the suite was set.
func (hs *serverHandshakeState) setCipherSuite(id uint16, supportedCipherSuites []uint16, version uint16) bool {
	for _, supported := range supportedCipherSuites {
		if id == supported {
			var candidate *cipherSuite

			for _, s := range cipherSuites {
				if s.id == id {
					candidate = s
					break
				}
			}
			if candidate == nil {
				continue
			}
			// Don't select a ciphersuite which we can't
			// support for this client.
			if candidate.flags&suiteECDHE != 0 { //EC??????
				if !hs.ellipticOk {
					continue
				}
				if candidate.flags&suiteECDSA != 0 {
					if !hs.ecdsaOk {
						continue
					}
				} else if !hs.rsaSignOk {
					continue
				}
			} else if candidate.flags&suiteSMDH != 0 && hs.sm2OK { //????????????
				if !hs.ellipticOk {
					continue
				}
			} else if hs.rsaDecryptOk {
				if candidate.flags&suiteECDSA != 0 {
					continue
				}
				if candidate.flags&suiteSMDH != 0 {
					continue
				}
			} else {
				continue
			}
			if version < VersionTLS12 && candidate.flags&suiteTLS12 != 0 {
				continue
			}
			hs.suite = candidate
			return true
		}
	}
	return false
}

// suppVersArray is the backing array of ClientHelloInfo.SupportedVersions
var suppVersArray = [...]uint16{VersionTLS12, VersionTLS11, VersionTLS10, VersionSSL30}

func (hs *serverHandshakeState) clientHelloInfo() *ClientHelloInfo {
	if hs.cachedClientHelloInfo != nil {
		return hs.cachedClientHelloInfo
	}

	var supportedVersions []uint16
	if hs.clientHello.vers > VersionTLS12 {
		supportedVersions = suppVersArray[:]
	} else if hs.clientHello.vers >= VersionSSL30 {
		supportedVersions = suppVersArray[VersionTLS12-hs.clientHello.vers:]
	}

	hs.cachedClientHelloInfo = &ClientHelloInfo{
		CipherSuites:      hs.clientHello.cipherSuites,
		ServerName:        hs.clientHello.serverName,
		SupportedCurves:   hs.clientHello.supportedCurves,
		SupportedPoints:   hs.clientHello.supportedPoints,
		SignatureSchemes:  hs.clientHello.supportedSignatureAlgorithms,
		SupportedProtos:   hs.clientHello.alpnProtocols,
		SupportedVersions: supportedVersions,
		Conn:              hs.c.conn,
	}

	return hs.cachedClientHelloInfo
}
