// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tls

import (
	"errors"
	"fmt"

	gmx509 "github.com/hyperchain/go-hpc-msp/plugin/software"
	"github.com/hyperchain/go-hpc-msp/tls/internal"
	"github.com/meshplus/crypto"
)

// pickSignatureAlgorithm selects a signature algorithm that is compatible with
// the given public key and the list of algorithms from the peer and this side.
// The lists of signature algorithms (peerSigAlgs and ourSigAlgs) are ignored
// for tlsVersion < VersionTLS12.
//
// The returned SignatureScheme codepoint is only meaningful for TLS 1.2,
// previous TLS versions have a fixed hash function.
func pickSignatureAlgorithm(pubkey crypto.VerifyKey, peerSigAlgs, ourSigAlgs []SignatureScheme, tlsVersion uint16) (sigAlg SignatureScheme, sigType uint8, hashFunc gmx509.Hash, err error) {
	if tlsVersion < VersionTLS12 || len(peerSigAlgs) == 0 {
		// For TLS 1.1 and before, the signature algorithm could not be
		// negotiated and the hash is fixed based on the signature type.
		// For TLS 1.2, if the client didn't send signature_algorithms
		// extension then we can assume that it supports SHA1. See
		// https://tools.ietf.org/html/rfc5246#section-7.4.1.4.1
		switch pubkey.GetKeyInfo() {
		case crypto.Rsa2048, crypto.Rsa3072, crypto.Rsa4096:
			if tlsVersion < VersionTLS12 {
				return 0, signaturePKCS1v15, gmx509.MD5SHA1, nil
			}
			return PKCS1WithSHA1, signaturePKCS1v15, gmx509.SHA1, nil
		case crypto.Secp256k1, crypto.Secp256r1, crypto.Secp384r1, crypto.Secp521r1:
			return ECDSAWithSHA1, signatureECDSA, gmx509.SHA1, nil
		case crypto.Sm2p256v1:
			return SM2WithSM3, signatureECDSA, gmx509.SM3WithPublicKey, nil
		default:
			return 0, 0, 0, fmt.Errorf("tls: unsupported public key: %T", pubkey)
		}
	}
	for _, sigAlg := range peerSigAlgs {
		if !isSupportedSignatureAlgorithm(sigAlg, ourSigAlgs) {
			continue
		}
		hashAlg, err := lookupTLSHash(sigAlg)
		if err != nil {
			panic("tls: supported signature algorithm has an unknown hash function")
		}
		sigType := signatureFromSignatureScheme(sigAlg)
		switch pubkey.GetKeyInfo() {
		case crypto.Rsa2048, crypto.Rsa3072, crypto.Rsa4096:
			if sigType == signaturePKCS1v15 || sigType == signatureRSAPSS {
				return sigAlg, sigType, hashAlg, nil
			}
		case crypto.Secp256k1, crypto.Secp256r1, crypto.Secp384r1, crypto.Secp521r1:
			if sigType == signatureECDSA && sigAlg != SM2WithSM3 {
				return sigAlg, sigType, hashAlg, nil
			}
		case crypto.Sm2p256v1:
			if sigType == signatureECDSA && sigAlg == SM2WithSM3 {
				return sigAlg, sigType, hashAlg, nil
			}
		default:
			return 0, 0, 0, fmt.Errorf("tls: unsupported public key: %T", pubkey)
		}
	}
	return 0, 0, 0, errors.New("tls: peer doesn't support any common signature algorithms")
}

// verifyHandshakeSignature verifies a signature against pre-hashed handshake
// contents.
func verifyHandshakeSignature(sigType uint8, pubkey crypto.VerifyKey, hashFunc gmx509.Hash, digest, sig []byte) error {
	hasher := internal.Copy()
	switch sigType {
	case signatureECDSA:
		//first, try to turn to gmx509.PublicKry
		switch pubkey.GetKeyInfo() {
		case crypto.Sm2p256v1:
			if len(sig) == 0 {
				return errors.New("tls: SM2 signature contained zero or negative values")
			}
			if !pubkey.Verify(digest, hasher, sig) {
				return errors.New("tls: SM2 signature verification failure")
			}
		case crypto.Secp256k1, crypto.Secp256r1, crypto.Secp384r1, crypto.Secp521r1:
			if !pubkey.Verify(digest, hasher, sig) {
				return errors.New("tls: ECDSA verification failure")
			}
		default:
			return errors.New("tls: ECDSA signing requires a ECDSA public key")
		}
	case signaturePKCS1v15:
		if pubkey.GetKeyInfo() > crypto.Rsa4096 || pubkey.GetKeyInfo() < crypto.Rsa2048 {
			return errors.New("tls: RSA signing requires a RSA public key")
		}
		if !pubkey.Verify(digest, hasher, sig) {
			return fmt.Errorf("tls: rsa signature verify error")
		}
	case signatureRSAPSS:
		//pubKey, ok := pubkey.(*rsa.PublicKey)
		//if !ok {
		//	return errors.New("tls: RSA signing requires a RSA public key")
		//}
		//signOpts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}
		//if err := rsa.VerifyPSS(pubKey, crypto.Hash(hashFunc), digest, sig, signOpts); err != nil {
		//	return err
		//}
		return errors.New("tls: not support RSAPSS")
	default:
		return errors.New("tls: unknown signature algorithm")
	}
	return nil
}
