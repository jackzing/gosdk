// Copyright 2015, 2018, 2019 Opsmate, Inc. All rights reserved.
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkcs12

import (
	"crypto/elliptic"
	"encoding/asn1"
	"errors"
	"fmt"
	"io"

	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/go-crypto-standard/asym"
	"github.com/hyperchain/go-crypto-standard/asym/secp256k1"
	"github.com/hyperchain/go-hpc-msp/pemencode"
	"github.com/hyperchain/go-hpc-msp/plugin/software"
	"github.com/meshplus/crypto"
)

var (
	// see https://tools.ietf.org/html/rfc7292#appendix-D
	oidCertTypeX509Certificate = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 9, 22, 1})
	oidPKCS8ShroundedKeyBag    = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 12, 10, 1, 2})
	oidCertBag                 = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 12, 10, 1, 3})
)

type certBag struct {
	ID   asn1.ObjectIdentifier
	Data []byte `asn1:"tag:0,explicit"`
}

func decodePkcs8ShroudedKeyBag(asn1Data, password []byte) (privateKey interface{}, err error) {
	pkinfo := new(encryptedPrivateKeyInfo)
	if err = unmarshal(asn1Data, pkinfo); err != nil {
		return nil, errors.New("pkcs12: error decoding PKCS#8 shrouded key bag: " + err.Error())
	}

	skDer, err := pbDecrypt(pkinfo, password)
	if err != nil {
		return nil, errors.New("pkcs12: error decrypting PKCS#8 shrouded key bag: " + err.Error())
	}

	pem, _ := pemencode.DER2PEM(skDer, pemencode.PEMAnyPrivateKey)

	engine := software.GetSoftwareEngine()
	if privateKey, err = engine.GetSignKey(string(pem)); err != nil {
		return nil, errors.New("pkcs12: error parsing private key: " + err.Error())
	}

	return privateKey, nil
}

func encodePkcs8ShroudedKeyBag(rand io.Reader, privateKey crypto.Signer, password []byte) (asn1Data []byte, err error) {
	var pkData []byte
	b, err := privateKey.Bytes()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	var mode int
	switch key := privateKey.(type) {
	case *gm.SM2PrivateKey:
		mode = crypto.Sm2p256v1
	case *asym.ECDSAPrivateKey:
		if key.Curve == secp256k1.S256() {
			mode = crypto.Secp256k1
		} else if key.Curve == elliptic.P256() {
			mode = crypto.Secp256r1
		}
	}

	pkData, err = software.MarshalPrivateKeyPKCS8(b, mode)
	if err != nil {
		return nil, fmt.Errorf("marshal key failed, reason:%v", err)
	}

	randomSalt := make([]byte, 8)
	_, _ = rand.Read(randomSalt)
	var paramBytes []byte
	if paramBytes, err = asn1.Marshal(pbeParams{Salt: randomSalt, Iterations: 2048}); err != nil {
		return nil, errors.New("pkcs12: error encoding params: " + err.Error())
	}

	var pkinfo encryptedPrivateKeyInfo
	pkinfo.AlgorithmIdentifier.Algorithm = oidPBEWithSHAAnd3KeyTripleDESCBC
	pkinfo.AlgorithmIdentifier.Parameters.FullBytes = paramBytes

	if err = pbEncrypt(&pkinfo, pkData, password); err != nil {
		return nil, errors.New("pkcs12: error encrypting PKCS#8 shrouded key bag: " + err.Error())
	}

	if asn1Data, err = asn1.Marshal(pkinfo); err != nil {
		return nil, errors.New("pkcs12: error encoding PKCS#8 shrouded key bag: " + err.Error())
	}

	return asn1Data, nil
}

func decodeCertBag(asn1Data []byte) (x509Certificates []byte, err error) {
	bag := new(certBag)
	if err := unmarshal(asn1Data, bag); err != nil {
		return nil, errors.New("pkcs12: error decoding cert bag: " + err.Error())
	}
	if !bag.ID.Equal(oidCertTypeX509Certificate) {
		return nil, NotImplementedError("only X509 certificates are supported")
	}
	return bag.Data, nil
}

func encodeCertBag(x509Certificates []byte) (asn1Data []byte, err error) {
	var bag certBag
	bag.ID = oidCertTypeX509Certificate
	bag.Data = x509Certificates
	if asn1Data, err = asn1.Marshal(bag); err != nil {
		return nil, errors.New("pkcs12: error encoding cert bag: " + err.Error())
	}
	return asn1Data, nil
}
