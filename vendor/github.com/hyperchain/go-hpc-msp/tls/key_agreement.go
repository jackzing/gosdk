// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

import (
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"errors"

	"github.com/hyperchain/go-hpc-msp/plugin/software"
	"github.com/hyperchain/go-hpc-msp/tls/internal"

	"io"
	"math/big"

	"golang.org/x/crypto/curve25519"

	"bytes"

	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/meshplus/crypto"
)

var errClientKeyExchange = errors.New("tls: invalid ClientKeyExchange message")
var errServerKeyExchange = errors.New("tls: invalid ServerKeyExchange message")

// rsaKeyAgreement implements the standard TLS key agreement where the client
// encrypts the pre-master secret to the server's public key.
type rsaKeyAgreement struct{}

func (ka rsaKeyAgreement) generateServerKeyExchange(config *Config, cert *Certificate, clientHello *clientHelloMsg, hello *serverHelloMsg) (*serverKeyExchangeMsg, error) {
	return nil, nil
}

func (ka rsaKeyAgreement) processClientKeyExchange(config *Config, cert *Certificate, ckx *clientKeyExchangeMsg, version uint16) ([]byte, error) {
	if len(ckx.ciphertext) < 2 {
		return nil, errClientKeyExchange
	}

	ciphertext := ckx.ciphertext
	if version != VersionSSL30 {
		ciphertextLen := int(ckx.ciphertext[0])<<8 | int(ckx.ciphertext[1])
		if ciphertextLen != len(ckx.ciphertext)-2 {
			return nil, errClientKeyExchange
		}
		ciphertext = ckx.ciphertext[2:]
	}
	priv, ok := cert.PrivateKey.(crypto.DecKey)
	if !ok {
		return nil, errors.New("tls: certificate private key does not implement crypto.Decrypter")
	}
	// Perform constant time RSA PKCS#1 v1.5 decryption
	preMasterSecret, err := priv.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}
	// We don't check the version number in the premaster secret. For one,
	// by checking it, we would leak information about the validity of the
	// encrypted pre-master secret. Secondly, it provides only a small
	// benefit against a downgrade attack and some implementations send the
	// wrong version anyway. See the discussion at the end of section
	// 7.4.7.1 of RFC 4346.
	return preMasterSecret, nil
}

func (ka rsaKeyAgreement) processServerKeyExchange(config *Config, clientHello *clientHelloMsg, serverHello *serverHelloMsg, cert *software.Certificate, skx *serverKeyExchangeMsg) error {
	return errors.New("tls: unexpected ServerKeyExchange")
}

func (ka rsaKeyAgreement) generateClientKeyExchange(config *Config, clientHello *clientHelloMsg, cert *software.Certificate) ([]byte, *clientKeyExchangeMsg, error) {
	preMasterSecret := make([]byte, 48)
	preMasterSecret[0] = byte(clientHello.vers >> 8)
	preMasterSecret[1] = byte(clientHello.vers)
	_, err := io.ReadFull(config.rand(), preMasterSecret[2:])
	if err != nil {
		return nil, nil, err
	}

	encrypter, ok := cert.PublicKey.(crypto.EncKey)
	if !ok {
		return nil, nil, errors.New("not a rsa publickey of enc")
	}

	encrypted, err := encrypter.Encrypt(preMasterSecret, config.rand())
	if err != nil {
		return nil, nil, err
	}
	ckx := new(clientKeyExchangeMsg)
	ckx.ciphertext = make([]byte, len(encrypted)+2)
	ckx.ciphertext[0] = byte(len(encrypted) >> 8)
	ckx.ciphertext[1] = byte(len(encrypted))
	copy(ckx.ciphertext[2:], encrypted)
	return preMasterSecret, ckx, nil
}

// sha1Hash calculates a SHA1 hash over the given byte slices.
func sha1Hash(slices [][]byte) []byte {
	hsha1 := sha1.New()
	for _, slice := range slices {
		_, _ = hsha1.Write(slice)
	}
	return hsha1.Sum(nil)
}

// md5SHA1Hash implements TLS 1.0's hybrid hash function which consists of the
// concatenation of an MD5 and SHA1 hash.
func md5SHA1Hash(slices [][]byte) []byte {
	md5sha1 := make([]byte, md5.Size+sha1.Size)
	hmd5 := md5.New()
	for _, slice := range slices {
		_, _ = hmd5.Write(slice)
	}
	copy(md5sha1, hmd5.Sum(nil))
	copy(md5sha1[md5.Size:], sha1Hash(slices))
	return md5sha1
}

// hashForServerKeyExchange hashes the given slices and returns their digest
// using the given hash function (for >= TLS 1.2) or using a default based on
// the sigType (for earlier TLS versions).
func hashForServerKeyExchange(sigType uint8, hashFunc software.Hash, version uint16, slices ...[]byte) ([]byte, error) {
	if version >= VersionTLS12 {
		h := hashFunc.New()
		for _, slice := range slices {
			_, _ = h.Write(slice)
		}
		digest := h.Sum(nil)
		return digest, nil
	}
	if sigType == signatureECDSA {
		return sha1Hash(slices), nil
	}
	return md5SHA1Hash(slices), nil
}

func curveForCurveID(id CurveID) (elliptic.Curve, bool) {
	switch id {
	case CurveP256:
		return elliptic.P256(), true
	case CurveP384:
		return elliptic.P384(), true
	case CurveP521:
		return elliptic.P521(), true
	case CurveGMP256:
		return gm.GetSm2Curve(), true
	default:
		return nil, false
	}

}

// ecdheKeyAgreement implements a TLS key agreement where the server
// generates an ephemeral EC public/private key pair and signs it. The
// pre-master secret is then calculated using ECDH. The signature may
// either be ECDSA or RSA.
type ecdheKeyAgreement struct {
	version    uint16
	isRSA      bool
	privateKey []byte
	curveid    CurveID

	// publicKey is used to store the peer's public value when X25519 is
	// being used.
	publicKey []byte
	// x and y are used to store the peer's public value when one of the
	// NIST curves is being used.
	x, y *big.Int
}

func (ka *ecdheKeyAgreement) generateServerKeyExchange(config *Config, cert *Certificate, clientHello *clientHelloMsg, hello *serverHelloMsg) (*serverKeyExchangeMsg, error) {
	preferredCurves := config.curvePreferences()

NextCandidate:
	for _, candidate := range preferredCurves {
		for _, c := range clientHello.supportedCurves {
			if candidate == c {
				ka.curveid = c
				break NextCandidate
			}
		}
	}

	if ka.curveid == 0 {
		return nil, errors.New("tls: no supported elliptic curves offered")
	}
	if ka.curveid == CurveGMP256 {
		Infof("generateServerKeyExchange, the server selects GMP256 curve")
	} else {
		Infof("generateServerKeyExchange, the server has selected non guomi curve: %v", ka.curveid)
	}
	var ecdhePublic []byte

	if ka.curveid == X25519 {
		var scalar, public [32]byte
		if _, err := io.ReadFull(config.rand(), scalar[:]); err != nil {
			return nil, err
		}

		curve25519.ScalarBaseMult(&public, &scalar)
		ka.privateKey = scalar[:]
		ecdhePublic = public[:]
	} else {
		curve, ok := curveForCurveID(ka.curveid)
		if !ok {
			return nil, errors.New("tls: preferredCurves includes unsupported curve")
		}

		var x, y *big.Int
		var err error
		ka.privateKey, x, y, err = elliptic.GenerateKey(curve, config.rand())
		if err != nil {
			return nil, err
		}
		ecdhePublic = elliptic.Marshal(curve, x, y)
		//fmt.Printf("###generateServerKeyExchange产生了临时的秘钥对%v\n", ka.curveid)
	}

	// https://tools.ietf.org/html/rfc4492#section-5.4
	serverECDHParams := make([]byte, 1+2+1+len(ecdhePublic))
	serverECDHParams[0] = 3 // named curve
	serverECDHParams[1] = byte(ka.curveid >> 8)
	serverECDHParams[2] = byte(ka.curveid)
	serverECDHParams[3] = byte(len(ecdhePublic))
	copy(serverECDHParams[4:], ecdhePublic)

	//priv, ok := cert.PrivateKey.(crypto.Signer)
	//if !ok {
	//	return nil, errors.New("tls: certificate private key does not implement crypto.Signer")
	//}

	priv := cert.PrivateKey
	//first para is self cert's public key
	signatureAlgorithm, sigType, hashFunc, err := pickSignatureAlgorithm(priv, clientHello.supportedSignatureAlgorithms, supportedSignatureAlgorithms, ka.version)
	if err != nil {
		return nil, err
	}
	if (sigType == signaturePKCS1v15 || sigType == signatureRSAPSS) != ka.isRSA {
		//fmt.Printf("###tls: certificate cannot be used with the selected cipher suite\n")
		return nil, errors.New("tls: certificate cannot be used with the selected cipher suite")
	}
	//fmt.Printf("###sig:%v, hash:%v, signAlgo: %v,isRSA:%v", sigType, hashFunc, signatureAlgorithm, ka.isRSA)
	var digest []byte
	if signatureAlgorithm == SM2WithSM3 {
		digest = bytes.Join([][]byte{clientHello.random, hello.random, serverECDHParams}, nil)
	} else {
		digest, err = hashForServerKeyExchange(sigType, hashFunc, ka.version, clientHello.random, hello.random, serverECDHParams)
		if err != nil {
			return nil, err
		}
	}

	var sig []byte
	if signatureAlgorithm == SM2WithSM3 {
		ok := priv.GetKeyInfo() == crypto.Sm2p256v1
		if !ok {
			return nil, errors.New("tls: failed to sign SM2WithSM3 parameters")
		}
		sig, err = priv.Sign(digest, internal.New(), rand.Reader)
		if err != nil {
			return nil, errors.New("tls: failed to sign SM2DHE parameters: " + err.Error())
		}
	} else {
		//signOpts := crypto.SignerOpts(hashFunc)
		//if sigType == signatureRSAPSS {
		//	signOpts = &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.Hash(hashFunc)}
		//}
		//sig, err = priv.Sign(config.rand(), digest, signOpts)
		if sigType == signaturePKCS1v15 {
			digest = append(digest, make([]byte, 4)...)
			binary.BigEndian.PutUint32(digest[len(digest)-4:], uint32(hashFunc))
		}
		sig, err = priv.Sign(digest, internal.Copy(), config.rand())
		if err != nil {
			return nil, errors.New("tls: failed to sign ECDHE parameters: " + err.Error())
		}
	}

	skx := new(serverKeyExchangeMsg)
	sigAndHashLen := 0
	if ka.version >= VersionTLS12 {
		sigAndHashLen = 2
	}
	skx.key = make([]byte, len(serverECDHParams)+sigAndHashLen+2+len(sig))
	copy(skx.key, serverECDHParams)
	k := skx.key[len(serverECDHParams):]
	if ka.version >= VersionTLS12 {
		k[0] = byte(signatureAlgorithm >> 8)
		k[1] = byte(signatureAlgorithm)
		k = k[2:]
	}
	k[0] = byte(len(sig) >> 8)
	k[1] = byte(len(sig))
	copy(k[2:], sig)

	return skx, nil
}

func (ka *ecdheKeyAgreement) processClientKeyExchange(config *Config, cert *Certificate, ckx *clientKeyExchangeMsg, version uint16) ([]byte, error) {
	if len(ckx.ciphertext) == 0 || int(ckx.ciphertext[0]) != len(ckx.ciphertext)-1 {
		return nil, errClientKeyExchange
	}

	if ka.curveid == X25519 {
		if len(ckx.ciphertext) != 1+32 {
			return nil, errClientKeyExchange
		}

		var theirPublic, sharedKey, scalar [32]byte
		copy(theirPublic[:], ckx.ciphertext[1:])
		copy(scalar[:], ka.privateKey)
		curve25519.ScalarMult(&sharedKey, &scalar, &theirPublic)
		return sharedKey[:], nil
	}

	curve, ok := curveForCurveID(ka.curveid)
	if !ok {
		panic("internal error")
	}
	x, y := elliptic.Unmarshal(curve, ckx.ciphertext[1:]) // Unmarshal also checks whether the given point is on the curve
	if x == nil {
		return nil, errClientKeyExchange
	}
	x, _ = curve.ScalarMult(x, y, ka.privateKey)
	preMasterSecret := make([]byte, (curve.Params().BitSize+7)>>3)
	xBytes := x.Bytes()
	copy(preMasterSecret[len(preMasterSecret)-len(xBytes):], xBytes)

	return preMasterSecret, nil
}

func (ka *ecdheKeyAgreement) processServerKeyExchange(config *Config, clientHello *clientHelloMsg, serverHello *serverHelloMsg, cert *software.Certificate, skx *serverKeyExchangeMsg) error {
	if len(skx.key) < 4 {
		return errServerKeyExchange
	}
	if skx.key[0] != 3 { // named curve
		return errors.New("tls: server selected unsupported curve")
	}
	ka.curveid = CurveID(skx.key[1])<<8 | CurveID(skx.key[2])

	publicLen := int(skx.key[3])
	if publicLen+4 > len(skx.key) {
		return errServerKeyExchange
	}
	serverECDHParams := skx.key[:4+publicLen]
	publicKey := serverECDHParams[4:]

	sig := skx.key[4+publicLen:]
	if len(sig) < 2 {
		return errServerKeyExchange
	}

	if ka.curveid == X25519 {
		if len(publicKey) != 32 {
			return errors.New("tls: bad X25519 public value")
		}
		ka.publicKey = publicKey
	} else {
		curve, ok := curveForCurveID(ka.curveid)
		if !ok {
			return errors.New("tls: server selected unsupported curve")
		}
		ka.x, ka.y = elliptic.Unmarshal(curve, publicKey) // Unmarshal also checks whether the given point is on the curve
		if ka.x == nil {
			return errServerKeyExchange
		}
	}

	var signatureAlgorithm SignatureScheme
	if ka.version >= VersionTLS12 {
		// handle SignatureAndHashAlgorithm
		signatureAlgorithm = SignatureScheme(sig[0])<<8 | SignatureScheme(sig[1])
		sig = sig[2:]
		if len(sig) < 2 {
			return errServerKeyExchange
		}
	}
	_, sigType, hashFunc, err := pickSignatureAlgorithm(cert.PublicKey, []SignatureScheme{signatureAlgorithm}, clientHello.supportedSignatureAlgorithms, ka.version)
	if err != nil {
		return err
	}
	if (sigType == signaturePKCS1v15 || sigType == signatureRSAPSS) != ka.isRSA {
		return errServerKeyExchange
	}

	sigLen := int(sig[0])<<8 | int(sig[1])
	if sigLen+2 != len(sig) {
		return errServerKeyExchange
	}
	sig = sig[2:]

	var digest []byte
	if signatureAlgorithm == SM2WithSM3 {
		digest = bytes.Join([][]byte{clientHello.random, serverHello.random, serverECDHParams}, nil)
		if !cert.PublicKey.Verify(digest, internal.New(), sig) {
			return errors.New("sm2 verify err")
		}
		return nil
	}
	digest, err = hashForServerKeyExchange(sigType, hashFunc, ka.version, clientHello.random, serverHello.random, serverECDHParams)
	if err != nil {
		return err
	}
	if sigType == signaturePKCS1v15 {
		digest = append(digest, make([]byte, 4)...)
		binary.BigEndian.PutUint32(digest[len(digest)-4:], uint32(hashFunc))
	}
	return verifyHandshakeSignature(sigType, cert.PublicKey, hashFunc, digest, sig)
}

func (ka *ecdheKeyAgreement) generateClientKeyExchange(config *Config, clientHello *clientHelloMsg, cert *software.Certificate) ([]byte, *clientKeyExchangeMsg, error) {
	if ka.curveid == 0 {
		return nil, nil, errors.New("tls: missing ServerKeyExchange message")
	}

	var serialized, preMasterSecret []byte

	if ka.curveid == X25519 {
		var ourPublic, theirPublic, sharedKey, scalar [32]byte

		if _, err := io.ReadFull(config.rand(), scalar[:]); err != nil {
			return nil, nil, err
		}

		copy(theirPublic[:], ka.publicKey)
		curve25519.ScalarBaseMult(&ourPublic, &scalar)
		curve25519.ScalarMult(&sharedKey, &scalar, &theirPublic)
		serialized = ourPublic[:]
		preMasterSecret = sharedKey[:]
	} else {
		curve, ok := curveForCurveID(ka.curveid)
		if !ok {
			panic("internal error")
		}
		priv, mx, my, err := elliptic.GenerateKey(curve, config.rand())
		if err != nil {
			return nil, nil, err
		}
		x, _ := curve.ScalarMult(ka.x, ka.y, priv)
		preMasterSecret = make([]byte, (curve.Params().BitSize+7)>>3)
		xBytes := x.Bytes()
		copy(preMasterSecret[len(preMasterSecret)-len(xBytes):], xBytes)

		serialized = elliptic.Marshal(curve, mx, my)
	}

	ckx := new(clientKeyExchangeMsg)
	ckx.ciphertext = make([]byte, 1+len(serialized))
	ckx.ciphertext[0] = byte(len(serialized))
	copy(ckx.ciphertext[1:], serialized)

	return preMasterSecret, ckx, nil
}

type eccAgreement struct {
	EncCert    []byte
	privateKey crypto.SignKey
}

func (ka *eccAgreement) generateServerKeyExchange(config *Config, cert *Certificate, clientHello *clientHelloMsg, hello *serverHelloMsg) (*serverKeyExchangeMsg, error) {
	//var encCert *gmx509.Certificate
	//for _, asn1Data := range cert.Certificate[1:] {
	//	cert, err := gmx509.ParseCertificate(asn1Data)
	//	if err != nil {
	//		continue
	//	}
	//	if cert.KeyUsage&gmx509.KeyUsageKeyEncipherment != 0 {
	//		encCert = cert
	//		break
	//	}
	//}
	engine := software.GetSoftwareEngine()
	encCert := cert.Certificate[1]
	if encCert == nil {
		return nil, errors.New("generateServerKeyExchange: no enc certificate from server ")
	}
	ka.EncCert = encCert
	if len(cert.Certificate) > 0 && bytes.Equal(ka.EncCert, cert.Certificate[0]) {
		ka.privateKey = cert.PrivateKey
	}

	skx := new(serverKeyExchangeMsg)

	priv := cert.PrivateKey
	h, _ := engine.GetHash(crypto.Sm3WithPublicKey)
	_, _ = h.Write(priv.Bytes())
	_, _ = h.Write(clientHello.random)
	_, _ = h.Write(hello.random)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(len(ka.EncCert)))
	_, _ = h.Write(b[1:])
	_, _ = h.Write(ka.EncCert)
	digest := h.Sum(nil)
	sig, err := priv.Sign(digest, internal.Copy(), rand.Reader)
	if err != nil {
		return nil, errors.New("generateServerKeyExchange: " + err.Error())
	}
	skx.key = make([]byte, 2+len(sig))
	binary.BigEndian.PutUint16(skx.key[:2], uint16(len(sig)))
	copy(skx.key[2:], sig)
	return skx, nil
}

func (ka *eccAgreement) processServerKeyExchange(config *Config, clientHello *clientHelloMsg, serverHello *serverHelloMsg, cert *software.Certificate, skx *serverKeyExchangeMsg) error {
	if len(skx.key) < 2 {
		return errServerKeyExchange
	}
	siglen := binary.BigEndian.Uint16(skx.key)
	if int(siglen)+2 != len(skx.key) {
		return errServerKeyExchange
	}

	if ka.EncCert == nil {
		return errors.New("processServerKeyExchange: no enc certificate from server ")
	}

	pub := cert.PublicKey

	engine := software.GetSoftwareEngine()
	h, _ := engine.GetHash(crypto.Sm3WithPublicKey)
	_, _ = h.Write(pub.Bytes())
	_, _ = h.Write(clientHello.random)
	_, _ = h.Write(serverHello.random)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(len(ka.EncCert)))
	_, _ = h.Write(b[1:])
	_, _ = h.Write(ka.EncCert)
	suc := pub.Verify(h.Sum(nil), internal.Copy(), skx.key[2:])
	if !suc {
		return errors.New("tls: invalid signature by the server certificate: ")
	}
	return nil
}

func (ka *eccAgreement) generateClientKeyExchange(config *Config, clientHello *clientHelloMsg, cert *software.Certificate) ([]byte, *clientKeyExchangeMsg, error) {
	preMasterSecret := make([]byte, 48)
	binary.BigEndian.PutUint16(preMasterSecret[:2], clientHello.vers)
	_, err := io.ReadFull(rand.Reader, preMasterSecret[2:])
	if err != nil {
		return nil, nil, err
	}

	cert, perr := software.ParseCertificate(string(ka.EncCert))
	if perr != nil {
		return nil, nil, errors.New("generateClientKeyExchange: no sm2 enc cert")
	}

	pub, ok := cert.PublicKey.(crypto.EncKey)
	if !ok {
		return nil, nil, errors.New("generateClientKeyExchange: no sm2 enc publicKey")
	}

	c1, encerr := pub.Encrypt(preMasterSecret, rand.Reader)
	if encerr != nil {
		return nil, nil, encerr
	}
	ciphertext := make([]byte, 3)
	ciphertext[0] = 48
	ciphertext[1] = 129
	bytes := make([]byte, 1)
	bytes[0] = 2
	if c1[1] > 127 {
		bytes = append(bytes, 33)
		bytes = append(bytes, 0)
	} else {
		bytes = append(bytes, 32)
	}
	bytes = append(bytes, c1[1:32+1]...)
	bytes = append(bytes, 2)
	if c1[1+32] > 127 {
		bytes = append(bytes, 33)
		bytes = append(bytes, 0)
	} else {
		bytes = append(bytes, 32)
	}
	bytes = append(bytes, c1[32+1:32+32+1]...)
	bytes = append(bytes, 4)
	bytes = append(bytes, 32)
	bytes = append(bytes, c1[len(c1)-32:]...)
	bytes = append(bytes, 4)
	bytes = append(bytes, byte(len(c1)-32*3-1))
	bytes = append(bytes, c1[1+32+32:len(c1)-32]...)
	ciphertext[2] = byte(len(bytes))
	ciphertext = append(ciphertext, bytes...)

	data := make([]byte, 2+len(ciphertext))
	binary.BigEndian.PutUint16(data[:2], uint16(len(ciphertext)))
	copy(data[2:], ciphertext)
	ckx := &clientKeyExchangeMsg{
		nil,
		data,
	}
	return preMasterSecret, ckx, nil
}

func (ka *eccAgreement) processClientKeyExchange(config *Config, cert *Certificate, ckx *clientKeyExchangeMsg, version uint16) ([]byte, error) {
	if len(ckx.ciphertext) == 0 || int(binary.BigEndian.Uint16(ckx.ciphertext[:2])) != len(ckx.ciphertext)-2 {
		return nil, errClientKeyExchange
	}

	privateKey := ka.privateKey
	if privateKey == nil {
		for _, cert := range config.Certificates {
			if bytes.Equal(cert.Certificate[0], ka.EncCert) {
				privateKey = cert.PrivateKey
			}
		}
		if privateKey == nil {
			return nil, errors.New("processClientKeyExchange: no sm2 enc privateKey")
		}
	}
	priv, ok := privateKey.(crypto.DecKey)
	if !ok {
		return nil, errors.New("processClientKeyExchange: no sm2 enc publicKey")
	}

	raw := ckx.ciphertext[2:]
	//todo
	ciphertext := make([]byte, 1)
	ciphertext[0] = 4

	if len(raw) < 128 || raw[0] != 48 && raw[1] != 129 {
		return nil, errors.New("processClientKeyExchange: err ciphertext")
	}
	off := 5 + raw[4] - 32
	ciphertext = append(ciphertext, raw[off:off+32]...)
	off += 32
	off += raw[off+1] - 32
	off += 2
	ciphertext = append(ciphertext, raw[off:off+32]...)
	off += 32
	ok = raw[off+1] == 32
	if !ok {
		return nil, errors.New("processClientKeyExchange: err ciphertext")
	}
	off += 2
	c3 := make([]byte, 32)
	copy(c3, raw[off:off+32])
	off += 32
	ok = raw[off+1] == 48
	if !ok {
		return nil, errors.New("processClientKeyExchange: err ciphertext")
	}
	off += 2
	ciphertext = append(ciphertext, raw[off:off+48]...)
	ciphertext = append(ciphertext, c3...)

	plaintext, err := priv.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}

	if version != binary.BigEndian.Uint16(plaintext[:2]) {
		return nil, errClientKeyExchange
	}
	return plaintext, nil
}
