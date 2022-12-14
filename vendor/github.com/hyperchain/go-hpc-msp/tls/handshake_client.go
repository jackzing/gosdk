// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

import (
	"bytes"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"fmt"

	"io"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/hyperchain/go-hpc-msp/plugin/software"
	"github.com/meshplus/crypto"
)

type clientHandshakeState struct {
	c            *Conn
	serverHello  *serverHelloMsg
	hello        *clientHelloMsg
	suite        *cipherSuite
	finishedHash finishedHash
	masterSecret []byte
	session      *ClientSessionState
}

func makeClientHello(config *Config) (*clientHelloMsg, error) {
	if len(config.ServerName) == 0 && !config.InsecureSkipVerify {
		return nil, errors.New("tls: either ServerName or InsecureSkipVerify must be specified in the tls.Config")
	}

	nextProtosLength := 0
	for _, proto := range config.NextProtos {
		l := len(proto)
		if l == 0 || l > 255 {
			return nil, errors.New("tls: invalid NextProtos value")
		}
		nextProtosLength += 1 + l

	}

	if nextProtosLength > 0xffff {
		return nil, errors.New("tls: NextProtos values too large")
	}

	hello := &clientHelloMsg{
		vers:                         config.maxVersion(),
		compressionMethods:           []uint8{compressionNone},
		random:                       make([]byte, 32),
		ocspStapling:                 true,
		scts:                         true,
		serverName:                   hostnameInSNI(config.ServerName),
		supportedCurves:              config.curvePreferences(),
		supportedPoints:              []uint8{pointFormatUncompressed},
		nextProtoNeg:                 len(config.NextProtos) > 0,
		secureRenegotiationSupported: true,
		alpnProtocols:                config.NextProtos,
	}
	possibleCipherSuites := config.cipherSuites()
	hello.cipherSuites = make([]uint16, 0, len(possibleCipherSuites))

NextCipherSuite:
	for _, suiteID := range possibleCipherSuites {
		for _, suite := range cipherSuites {
			if suite.id != suiteID {
				continue
			}
			// Don't advertise TLS 1.2-only cipher suites unless
			// we're attempting TLS 1.2.
			if hello.vers < VersionTLS12 && suite.flags&suiteTLS12 != 0 {
				continue
			}
			hello.cipherSuites = append(hello.cipherSuites, suiteID)
			continue NextCipherSuite
		}
	}

	_, err := io.ReadFull(config.rand(), hello.random)
	if err != nil {
		return nil, errors.New("tls: short read from Rand: " + err.Error())
	}

	if len(config.Certificates) > 0 {
		smOk := config.Certificates[0].PrivateKey.GetKeyInfo() == crypto.Sm2p256v1
		if smOk {
			hello.vers = VersionGMTLS
		}
	}
	if hello.vers >= VersionTLS12 {
		hello.supportedSignatureAlgorithms = supportedSignatureAlgorithms
	}

	return hello, nil
}

func (c *Conn) clientHandshake() error {
	if c.config == nil {
		c.config = defaultConfig()
	}

	// This may be a renegotiation handshake, in which case some fields
	// need to be reset.
	c.didResume = false

	hello, err := makeClientHello(c.config)
	if err != nil {
		return err
	}

	if c.handshakes > 0 {
		hello.secureRenegotiation = c.clientFinished[:]
	}

	var session *ClientSessionState
	var cacheKey string
	sessionCache := c.config.ClientSessionCache
	if c.config.SessionTicketsDisabled {
		sessionCache = nil
	}

	if sessionCache != nil {
		hello.ticketSupported = true
	}

	// Session resumption is not allowed if renegotiating because
	// renegotiation is primarily used to allow a client to send a client
	// certificate, which would be skipped if session resumption occurred.
	if sessionCache != nil && c.handshakes == 0 {
		// Try to resume a previously negotiated TLS session, if
		// available.
		cacheKey = clientSessionCacheKey(c.conn.RemoteAddr(), c.config)
		candidateSession, ok := sessionCache.Get(cacheKey)
		if ok {
			// Check that the ciphersuite/version used for the
			// previous session are still valid.
			cipherSuiteOk := false
			for _, id := range hello.cipherSuites {
				if id == candidateSession.cipherSuite {
					cipherSuiteOk = true
					break
				}
			}

			versOk := candidateSession.vers >= c.config.minVersion() &&
				candidateSession.vers <= c.config.maxVersion()
			if versOk && cipherSuiteOk {
				session = candidateSession
			}
		}
	}

	if session != nil {
		hello.sessionTicket = session.sessionTicket
		// A random session ID is used to detect when the
		// server accepted the ticket and is resuming a session
		// (see RFC 5077).
		hello.sessionID = make([]byte, 16)
		if _, rerr := io.ReadFull(c.config.rand(), hello.sessionID); rerr != nil {
			return errors.New("tls: short read from Rand: " + rerr.Error())
		}
	}

	hs := &clientHandshakeState{
		c:       c,
		hello:   hello,
		session: session,
	}

	if err = hs.handshake(); err != nil {
		return err
	}

	// If we had a successful handshake and hs.session is different from
	// the one already cached - cache a new one
	if sessionCache != nil && hs.session != nil && session != hs.session {
		sessionCache.Put(cacheKey, hs.session)
	}

	return nil
}

// Does the handshake, either a full one or resumes old session.
// Requires hs.c, hs.hello, and, optionally, hs.session to be set.
func (hs *clientHandshakeState) handshake() error {
	c := hs.c

	// send ClientHello
	if _, err := c.writeRecord(recordTypeHandshake, hs.hello.marshal()); err != nil {
		return err
	}

	msg, err := c.readHandshake()
	if err != nil {
		return err
	}

	var ok bool
	if hs.serverHello, ok = msg.(*serverHelloMsg); !ok {
		_ = c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(hs.serverHello, msg)
	}

	if err = hs.pickTLSVersion(); err != nil {
		return err
	}

	if err = hs.pickCipherSuite(); err != nil {
		return err
	}

	isResume, err := hs.processServerHello()
	if err != nil {
		return err
	}

	hs.finishedHash = newFinishedHash(c.vers, hs.suite)

	// No signatures of the handshake are needed in a resumption.
	// Otherwise, in a full handshake, if we don't have any certificates
	// configured then we will never send a CertificateVerify message and
	// thus no signatures are needed in that case either.
	if isResume || (len(c.config.Certificates) == 0 && c.config.GetClientCertificate == nil) {
		hs.finishedHash.discardHandshakeBuffer()
	}

	_, _ = hs.finishedHash.Write(hs.hello.marshal())
	_, _ = hs.finishedHash.Write(hs.serverHello.marshal())

	c.buffering = true
	if isResume {
		if err := hs.establishKeys(); err != nil {
			return err
		}
		if err := hs.readSessionTicket(); err != nil {
			return err
		}
		if err := hs.readFinished(c.serverFinished[:]); err != nil {
			return err
		}
		c.clientFinishedIsFirst = false
		if err := hs.sendFinished(c.clientFinished[:]); err != nil {
			return err
		}
		if _, err := c.flush(); err != nil {
			return err
		}
	} else {
		if err := hs.doFullHandshake(); err != nil {
			return err
		}
		if err := hs.establishKeys(); err != nil {
			return err
		}
		if err := hs.sendFinished(c.clientFinished[:]); err != nil {
			return err
		}
		if _, err := c.flush(); err != nil {
			return err
		}
		c.clientFinishedIsFirst = true
		if err := hs.readSessionTicket(); err != nil {
			return err
		}
		if err := hs.readFinished(c.serverFinished[:]); err != nil {
			return err
		}
	}

	c.ekm = ekmFromMasterSecret(c.vers, hs.suite, hs.masterSecret, hs.hello.random, hs.serverHello.random)
	c.didResume = isResume
	atomic.StoreUint32(&c.handshakeStatus, 1)

	return nil
}

func (hs *clientHandshakeState) pickTLSVersion() error {
	vers, ok := hs.c.config.mutualVersion(hs.serverHello.vers)
	if (!ok || vers < VersionTLS10) && hs.serverHello.vers != VersionGMTLS {
		// TLS 1.0 is the minimum version supported as a client.
		_ = hs.c.sendAlert(alertProtocolVersion)
		return fmt.Errorf("tls: server selected unsupported protocol version %x", hs.serverHello.vers)
	}

	hs.c.vers = vers
	hs.c.haveVers = true

	return nil
}

func (hs *clientHandshakeState) pickCipherSuite() error {
	if hs.suite = mutualCipherSuite(hs.hello.cipherSuites, hs.serverHello.cipherSuite); hs.suite == nil {
		_ = hs.c.sendAlert(alertHandshakeFailure)
		return errors.New("tls: server chose an unconfigured cipher suite")
	}

	hs.c.cipherSuite = hs.suite.id
	return nil
}

func (hs *clientHandshakeState) doFullHandshake() error {
	c := hs.c

	msg, err := c.readHandshake()
	if err != nil {
		return err
	}
	certMsg, ok := msg.(*certificateMsg)
	if !ok || len(certMsg.certificates) == 0 {
		_ = c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(certMsg, msg)
	}
	_, _ = hs.finishedHash.Write(certMsg.marshal())

	if c.handshakes == 0 {
		// If this is the first handshake on a connection, process and
		// (optionally) verify the server's certificates.
		certs := make([]*software.Certificate, len(certMsg.certificates))
		for i, asn1Data := range certMsg.certificates {
			cert, perr := software.ParseCertificate(string(asn1Data))
			if perr != nil {
				_ = c.sendAlert(alertBadCertificate)
				Errorf("tls: failed to parse certificate from server: " + perr.Error())
				return errors.New("tls: failed to parse certificate from server: " + perr.Error())
			}
			certs[i] = cert
		}

		if !c.config.InsecureSkipVerify {
			opts := software.VerifyOptions{
				Roots:         c.config.RootCAs,
				CurrentTime:   c.config.time(),
				DNSName:       c.config.ServerName,
				Intermediates: software.NewCertPool(),
			}

			for i, cert := range certs {
				if i == 0 { //0 is entry cert
					continue
				}
				opts.Intermediates.AddCert(cert)
			}
			c.verifiedChains, err = certs[0].Verify(opts)
			if err != nil {
				Errorf("alert bad certificate: %v", err.Error())
				_ = c.sendAlert(alertBadCertificate)
				return err
			}
		}

		if c.config.VerifyPeerCertificate != nil {
			if verr := c.config.VerifyPeerCertificate(certMsg.certificates, c.verifiedChains); verr != nil {
				Errorf("alert bad certificate: %v", verr.Error())
				_ = c.sendAlert(alertBadCertificate)
				return verr
			}
		}

		if certs[0].PublicKey.GetKeyInfo() < crypto.Sm2p256v1 || certs[0].PublicKey.GetKeyInfo() > crypto.Rsa4096 {
			_ = c.sendAlert(alertUnsupportedCertificate)
			return fmt.Errorf("tls: server's certificate contains an unsupported type of public key: %T", certs[0].PublicKey)
		}

		c.peerCertificates = certs
	} else {
		// This is a renegotiation handshake. We require that the
		// server's identity (i.e. leaf certificate) is unchanged and
		// thus any previous trust decision is still valid.
		//
		// See https://mitls.org/pages/attacks/3SHAKE for the
		// motivation behind this requirement.
		if !bytes.Equal(c.peerCertificates[0].Raw, certMsg.certificates[0]) {
			_ = c.sendAlert(alertBadCertificate)
			Errorf("tls: server's identity changed during renegotiation")
			return errors.New("tls: server's identity changed during renegotiation")
		}
	}

	msg, err = c.readHandshake()
	if err != nil {
		return err
	}

	cs, ok := msg.(*certificateStatusMsg)
	if ok {
		// RFC4366 on Certificate Status Request:
		// The server MAY return a "certificate_status" message.

		if !hs.serverHello.ocspStapling {
			// If a server returns a "CertificateStatus" message, then the
			// server MUST have included an extension of type "status_request"
			// with empty "extension_data" in the extended server hello.

			_ = c.sendAlert(alertUnexpectedMessage)
			return errors.New("tls: received unexpected CertificateStatus message")
		}
		_, _ = hs.finishedHash.Write(cs.marshal())

		if cs.statusType == statusTypeOCSP {
			c.ocspResponse = cs.response
		}

		msg, err = c.readHandshake()
		if err != nil {
			return err
		}
	}

	keyAgreement := hs.suite.ka(c.vers)
	ecckeyAgreement, eccOk := keyAgreement.(*eccAgreement)
	if eccOk {
		if len(c.peerCertificates) > 1 {
			ecckeyAgreement.EncCert = c.peerCertificates[1].Raw
		}
	}

	skx, ok := msg.(*serverKeyExchangeMsg)
	if ok {
		_, _ = hs.finishedHash.Write(skx.marshal())
		err = keyAgreement.processServerKeyExchange(c.config, hs.hello, hs.serverHello, c.peerCertificates[0], skx)
		if err != nil {
			_ = c.sendAlert(alertUnexpectedMessage)
			return err
		}

		msg, err = c.readHandshake()
		if err != nil {
			return err
		}
	}

	var chainToSend *Certificate
	var certRequested bool
	certReq, ok := msg.(*certificateRequestMsg)
	if ok {
		certRequested = true
		_, _ = hs.finishedHash.Write(certReq.marshal())

		if chainToSend, err = hs.getCertificate(certReq); err != nil {
			_ = c.sendAlert(alertInternalError)
			return err
		}

		msg, err = c.readHandshake()
		if err != nil {
			return err
		}
	}

	shd, ok := msg.(*serverHelloDoneMsg)
	if !ok {
		_ = c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(shd, msg)
	}
	_, _ = hs.finishedHash.Write(shd.marshal())

	// If the server requested a certificate then we have to send a
	// Certificate message, even if it's empty because we don't have a
	// certificate to send.
	if certRequested {
		certMsg = new(certificateMsg)
		certMsg.certificates = chainToSend.Certificate
		_, _ = hs.finishedHash.Write(certMsg.marshal())
		if _, werr := c.writeRecord(recordTypeHandshake, certMsg.marshal()); werr != nil {
			return werr
		}
	}

	preMasterSecret, ckx, err := keyAgreement.generateClientKeyExchange(c.config, hs.hello, c.peerCertificates[0])
	if err != nil {
		_ = c.sendAlert(alertInternalError)
		return err
	}
	if ckx != nil {
		_, _ = hs.finishedHash.Write(ckx.marshal())
		if _, err := c.writeRecord(recordTypeHandshake, ckx.marshal()); err != nil {
			return err
		}
	}
	engine := software.GetSoftwareEngine()
	if chainToSend != nil && len(chainToSend.Certificate) > 0 {
		certVerify := &certificateVerifyMsg{
			hasSignatureAndHash: c.vers >= VersionTLS12,
		}

		key := chainToSend.PrivateKey
		if reflect.DeepEqual(chainToSend.PrivateKey, nil) {
			_ = c.sendAlert(alertInternalError)
			return fmt.Errorf("tls: client certificate private key of type %T does not implement crypto.SignKey", chainToSend.PrivateKey)
		}

		signatureAlgorithm, sigType, hashFunc, err := pickSignatureAlgorithm(chainToSend.PrivateKey, certReq.supportedSignatureAlgorithms, hs.hello.supportedSignatureAlgorithms, c.vers)
		if err != nil {
			_ = c.sendAlert(alertInternalError)
			return err
		}
		// SignatureAndHashAlgorithm was introduced in TLS 1.2.
		if certVerify.hasSignatureAndHash {
			certVerify.signatureAlgorithm = signatureAlgorithm
		}

		digest, err := hs.finishedHash.hashForClientCertificate(sigType, hashFunc, hs.masterSecret)
		if sigType == signaturePKCS1v15 {
			digest = append(digest, make([]byte, 4)...)
			binary.BigEndian.PutUint32(digest[len(digest)-4:], uint32(hashFunc))
		}

		if err != nil {
			_ = c.sendAlert(alertInternalError)
			return err
		}

		//signOpts := std.SignerOpts(hashFunc)
		//if sigType == signatureRSAPSS {
		//	signOpts = &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: std.Hash(hashFunc)}
		//}
		if signatureAlgorithm == SM2WithSM3 {
			if key.GetKeyInfo() != crypto.Sm2p256v1 {
				_ = c.sendAlert(alertInternalError)
				return err
			}
			h, _ := engine.GetHash(crypto.Sm3WithPublicKey)
			digest, _ = h.BatchHash([][]byte{key.Bytes(), digest})
		}
		fake, _ := engine.GetHash(crypto.FakeHash)
		certVerify.signature, err = key.Sign(digest, fake, c.config.rand())
		if err != nil {
			_ = c.sendAlert(alertInternalError)
			return err
		}

		_, _ = hs.finishedHash.Write(certVerify.marshal())
		if _, err := c.writeRecord(recordTypeHandshake, certVerify.marshal()); err != nil {
			return err
		}
	}

	hs.masterSecret = masterFromPreMasterSecret(c.vers, hs.suite, preMasterSecret, hs.hello.random, hs.serverHello.random)
	if err := c.config.writeKeyLog(hs.hello.random, hs.masterSecret); err != nil {
		_ = c.sendAlert(alertInternalError)
		return errors.New("tls: failed to write to key log: " + err.Error())
	}

	hs.finishedHash.discardHandshakeBuffer()

	return nil
}

func (hs *clientHandshakeState) establishKeys() error {
	c := hs.c

	clientMAC, serverMAC, clientKey, serverKey, clientIV, serverIV :=
		keysFromMasterSecret(c.vers, hs.suite, hs.masterSecret, hs.hello.random, hs.serverHello.random, hs.suite.macLen, hs.suite.keyLen, hs.suite.ivLen)
	var clientCipher, serverCipher interface{}
	var clientHash, serverHash macFunction
	if hs.suite.cipher != nil {
		clientCipher = hs.suite.cipher(clientKey, clientIV, false /* enc */)
		clientHash = hs.suite.mac(c.vers, clientMAC)
		serverCipher = hs.suite.cipher(serverKey, serverIV, true /* dec */)
		serverHash = hs.suite.mac(c.vers, serverMAC)
	} else {
		clientCipher = hs.suite.aead(clientKey, clientIV)
		serverCipher = hs.suite.aead(serverKey, serverIV)
	}

	c.in.prepareCipherSpec(c.vers, serverCipher, serverHash)
	c.out.prepareCipherSpec(c.vers, clientCipher, clientHash)
	return nil
}

func (hs *clientHandshakeState) serverResumedSession() bool {
	// If the server responded with the same sessionID then it means the
	// sessionTicket is being used to resume a TLS session.
	return hs.session != nil && hs.hello.sessionID != nil &&
		bytes.Equal(hs.serverHello.sessionID, hs.hello.sessionID)
}

func (hs *clientHandshakeState) processServerHello() (bool, error) {
	c := hs.c

	if hs.serverHello.compressionMethod != compressionNone {
		_ = c.sendAlert(alertUnexpectedMessage)
		return false, errors.New("tls: server selected unsupported compression format")
	}

	if c.handshakes == 0 && hs.serverHello.secureRenegotiationSupported {
		c.secureRenegotiation = true
		if len(hs.serverHello.secureRenegotiation) != 0 {
			_ = c.sendAlert(alertHandshakeFailure)
			return false, errors.New("tls: initial handshake had non-empty renegotiation extension")
		}
	}

	if c.handshakes > 0 && c.secureRenegotiation {
		var expectedSecureRenegotiation [24]byte
		copy(expectedSecureRenegotiation[:], c.clientFinished[:])
		copy(expectedSecureRenegotiation[12:], c.serverFinished[:])
		if !bytes.Equal(hs.serverHello.secureRenegotiation, expectedSecureRenegotiation[:]) {
			_ = c.sendAlert(alertHandshakeFailure)
			return false, errors.New("tls: incorrect renegotiation extension contents")
		}
	}

	clientDidNPN := hs.hello.nextProtoNeg
	clientDidALPN := len(hs.hello.alpnProtocols) > 0
	serverHasNPN := hs.serverHello.nextProtoNeg
	serverHasALPN := len(hs.serverHello.alpnProtocol) > 0

	if !clientDidNPN && serverHasNPN {
		_ = c.sendAlert(alertHandshakeFailure)
		return false, errors.New("tls: server advertised unrequested NPN extension")
	}

	if !clientDidALPN && serverHasALPN {
		_ = c.sendAlert(alertHandshakeFailure)
		return false, errors.New("tls: server advertised unrequested ALPN extension")
	}

	if serverHasNPN && serverHasALPN {
		_ = c.sendAlert(alertHandshakeFailure)
		return false, errors.New("tls: server advertised both NPN and ALPN extensions")
	}

	if serverHasALPN {
		c.clientProtocol = hs.serverHello.alpnProtocol
		c.clientProtocolFallback = false
	}
	c.scts = hs.serverHello.scts

	if !hs.serverResumedSession() {
		return false, nil
	}

	if hs.session.vers != c.vers {
		_ = c.sendAlert(alertHandshakeFailure)
		return false, errors.New("tls: server resumed a session with a different version")
	}

	if hs.session.cipherSuite != hs.suite.id {
		_ = c.sendAlert(alertHandshakeFailure)
		return false, errors.New("tls: server resumed a session with a different cipher suite")
	}

	// Restore masterSecret and peerCerts from previous state
	hs.masterSecret = hs.session.masterSecret
	c.peerCertificates = hs.session.serverCertificates
	c.verifiedChains = hs.session.verifiedChains
	return true, nil
}

func (hs *clientHandshakeState) readFinished(out []byte) error {
	c := hs.c

	_ = c.readRecord(recordTypeChangeCipherSpec)
	if c.in.err != nil {
		return c.in.err
	}

	msg, err := c.readHandshake()
	if err != nil {
		return err
	}
	serverFinished, ok := msg.(*finishedMsg)
	if !ok {
		_ = c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(serverFinished, msg)
	}

	verify := hs.finishedHash.serverSum(hs.masterSecret)
	if len(verify) != len(serverFinished.verifyData) ||
		subtle.ConstantTimeCompare(verify, serverFinished.verifyData) != 1 {
		_ = c.sendAlert(alertHandshakeFailure)
		return errors.New("tls: server's Finished message was incorrect")
	}
	_, _ = hs.finishedHash.Write(serverFinished.marshal())
	copy(out, verify)
	return nil
}

func (hs *clientHandshakeState) readSessionTicket() error {
	if !hs.serverHello.ticketSupported {
		return nil
	}

	c := hs.c
	msg, err := c.readHandshake()
	if err != nil {
		return err
	}
	sessionTicketMsg, ok := msg.(*newSessionTicketMsg)
	if !ok {
		_ = c.sendAlert(alertUnexpectedMessage)
		return unexpectedMessageError(sessionTicketMsg, msg)
	}
	_, _ = hs.finishedHash.Write(sessionTicketMsg.marshal())

	hs.session = &ClientSessionState{
		sessionTicket:      sessionTicketMsg.ticket,
		vers:               c.vers,
		cipherSuite:        hs.suite.id,
		masterSecret:       hs.masterSecret,
		serverCertificates: c.peerCertificates,
		verifiedChains:     c.verifiedChains,
	}

	return nil
}

func (hs *clientHandshakeState) sendFinished(out []byte) error {
	c := hs.c

	if _, err := c.writeRecord(recordTypeChangeCipherSpec, []byte{1}); err != nil {
		return err
	}
	if hs.serverHello.nextProtoNeg {
		nextProto := new(nextProtoMsg)
		proto, fallback := mutualProtocol(c.config.NextProtos, hs.serverHello.nextProtos)
		nextProto.proto = proto
		c.clientProtocol = proto
		c.clientProtocolFallback = fallback

		_, _ = hs.finishedHash.Write(nextProto.marshal())
		if _, err := c.writeRecord(recordTypeHandshake, nextProto.marshal()); err != nil {
			return err
		}
	}

	finished := new(finishedMsg)
	finished.verifyData = hs.finishedHash.clientSum(hs.masterSecret)
	_, _ = hs.finishedHash.Write(finished.marshal())
	if _, err := c.writeRecord(recordTypeHandshake, finished.marshal()); err != nil {
		return err
	}
	copy(out, finished.verifyData)
	return nil
}

// tls11SignatureSchemes contains the signature schemes that we synthesise for
// a TLS <= 1.1 connection, based on the supported certificate types.
var tls11SignatureSchemes = []SignatureScheme{ECDSAWithP256AndSHA256, ECDSAWithP384AndSHA384, ECDSAWithP521AndSHA512, PKCS1WithSHA256, PKCS1WithSHA384, PKCS1WithSHA512, PKCS1WithSHA1}

const (
	// tls11SignatureSchemesNumECDSA is the number of initial elements of
	// tls11SignatureSchemes that use ECDSA.
	tls11SignatureSchemesNumECDSA = 3
	// tls11SignatureSchemesNumRSA is the number of trailing elements of
	// tls11SignatureSchemes that use RSA.
	tls11SignatureSchemesNumRSA = 4
)

func (hs *clientHandshakeState) getCertificate(certReq *certificateRequestMsg) (*Certificate, error) {
	c := hs.c

	var rsaAvail, ecdsaAvail bool
	for _, certType := range certReq.certificateTypes {
		switch certType {
		case certTypeRSASign:
			rsaAvail = true
		case certTypeECDSASign:
			ecdsaAvail = true
		}
	}

	if c.config.GetClientCertificate != nil {
		var signatureSchemes []SignatureScheme

		if !certReq.hasSignatureAndHash {
			// Prior to TLS 1.2, the signature schemes were not
			// included in the certificate request message. In this
			// case we use a plausible list based on the acceptable
			// certificate types.
			signatureSchemes = tls11SignatureSchemes
			if !ecdsaAvail {
				signatureSchemes = signatureSchemes[tls11SignatureSchemesNumECDSA:]
			}
			if !rsaAvail {
				signatureSchemes = signatureSchemes[:len(signatureSchemes)-tls11SignatureSchemesNumRSA]
			}
		} else {
			signatureSchemes = certReq.supportedSignatureAlgorithms
		}

		return c.config.GetClientCertificate(&CertificateRequestInfo{
			AcceptableCAs:    certReq.certificateAuthorities,
			SignatureSchemes: signatureSchemes,
		})
	}

	// RFC 4346 on the certificateAuthorities field: A list of the
	// distinguished names of acceptable certificate authorities.
	// These distinguished names may specify a desired
	// distinguished name for a root CA or for a subordinate CA;
	// thus, this message can be used to describe both known roots
	// and a desired authorization space. If the
	// certificate_authorities list is empty then the client MAY
	// send any certificate of the appropriate
	// ClientCertificateType, unless there is some external
	// arrangement to the contrary.

	// We need to search our list of client certs for one
	// where SignatureAlgorithm is acceptable to the server and the
	// Issuer is in certReq.certificateAuthorities
findCert:
	for i, chain := range c.config.Certificates {
		if !rsaAvail && !ecdsaAvail {
			continue
		}

		for j, cert := range chain.Certificate {
			x509Cert := chain.Leaf
			// parse the certificate if this isn't the leaf
			// node, or if chain.Leaf was nil
			if j != 0 || x509Cert == nil {
				var err error
				if x509Cert, err = software.ParseCertificate(string(cert)); err != nil {
					_ = c.sendAlert(alertInternalError)
					return nil, errors.New("tls: failed to parse client certificate #" + strconv.Itoa(i) + ": " + err.Error())
				}
			}

			switch {
			case rsaAvail && x509Cert.PublicKeyAlgorithm == software.RSA:
			case ecdsaAvail && x509Cert.PublicKeyAlgorithm == software.ECDSA:
			case ecdsaAvail && x509Cert.PublicKeyAlgorithm == software.SM2:
			default:
				continue findCert
			}

			if len(certReq.certificateAuthorities) == 0 {
				// they gave us an empty list, so just take the
				// first cert from c.config.Certificates
				return &chain, nil
			}

			for _, ca := range certReq.certificateAuthorities {
				if bytes.Equal(x509Cert.RawIssuer, ca) {
					return &chain, nil
				}
			}
		}
	}

	// No acceptable certificate found. Don't send a certificate.
	return new(Certificate), nil
}

// clientSessionCacheKey returns a key used to cache sessionTickets that could
// be used to resume previously negotiated TLS sessions with a server.
func clientSessionCacheKey(serverAddr net.Addr, config *Config) string {
	if len(config.ServerName) > 0 {
		return config.ServerName
	}
	return serverAddr.String()
}

// mutualProtocol finds the mutual Next Protocol Negotiation or ALPN protocol
// given list of possible protocols and a list of the preference order. The
// first list must not be empty. It returns the resulting protocol and flag
// indicating if the fallback case was reached.
func mutualProtocol(protos, preferenceProtos []string) (string, bool) {
	for _, s := range preferenceProtos {
		for _, c := range protos {
			if s == c {
				return s, false
			}
		}
	}

	return protos[0], true
}

// hostnameInSNI converts name into an approriate hostname for SNI.
// Literal IP addresses and absolute FQDNs are not permitted as SNI values.
// See https://tools.ietf.org/html/rfc6066#section-3.
func hostnameInSNI(name string) string {
	host := name
	if len(host) > 0 && host[0] == '[' && host[len(host)-1] == ']' {
		host = host[1 : len(host)-1]
	}
	if i := strings.LastIndex(host, "%"); i > 0 {
		host = host[:i]
	}
	if net.ParseIP(host) != nil {
		return ""
	}
	for len(name) > 0 && name[len(name)-1] == '.' {
		name = name[:len(name)-1]
	}
	return name
}
