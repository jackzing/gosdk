package pemencode

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

//PEMType is pem type
type PEMType int

//pem type enum
const (
	PEMECCPrivateKey PEMType = iota
	PEMRSAPrivateKey
	PEMAnyPrivateKey
	PEMPublicKey
	PEMCertificate
	PLUGINEncodeKey
	PEMInvalidPEMType
)

//pem file header
const (
	PemTypeECPrivateKey  = "EC PRIVATE KEY"
	PemTypeCertificate   = "CERTIFICATE"
	PemTypePublicKey     = "PUBLIC KEY"
	PemTypeRSAPrivateKey = "RSA PRIVATE KEY"
	PemTypeAnyPrivateKey = "PRIVATE KEY"
)

// PEM2DER pem to der
func PEM2DER(raw []byte) ([]byte, PEMType) {
	if len(raw) == 0 {
		return nil, PEMInvalidPEMType
	}

	//see function ParsePluginKey(...) in plugin package
	if strings.HasPrefix(string(raw), "plugin") {
		return raw, PLUGINEncodeKey
	}

	block, _ := pem.Decode(bytes.TrimSpace(raw))
	if block == nil {
		return raw, PEMInvalidPEMType
	}
	return block.Bytes, getPemType(block.Type)
}

//DER2PEM encode der to pem
func DER2PEM(in []byte, t PEMType) ([]byte, error) {
	pb := new(pem.Block)
	if t >= PEMInvalidPEMType || t < 0 {
		return nil, errors.New("unknown pem type")
	}
	switch t {
	case PEMPublicKey:
		pb.Type = PemTypePublicKey
	case PEMECCPrivateKey:
		pb.Type = PemTypeECPrivateKey
	case PEMRSAPrivateKey:
		pb.Type = PemTypeRSAPrivateKey
	case PEMAnyPrivateKey:
		pb.Type = PemTypeAnyPrivateKey
	case PEMCertificate:
		pb.Type = PemTypeCertificate
	}
	pb.Bytes = in
	ret := pem.EncodeToMemory(pb)
	return ret, nil
}

//DER2PEMWithEncryption encode der to pem with encryption
func DER2PEMWithEncryption(in []byte, t PEMType, pwd [32]byte) ([]byte, error) {
	var pemType string
	if t >= PEMInvalidPEMType || t < 0 {
		return nil, errors.New("unknown pem type")
	}
	switch t {
	case PEMPublicKey:
		pemType = PemTypePublicKey
	case PEMECCPrivateKey:
		pemType = PemTypeECPrivateKey
	case PEMRSAPrivateKey:
		pemType = PemTypeRSAPrivateKey
	case PEMAnyPrivateKey:
		pemType = PemTypeAnyPrivateKey
	case PEMCertificate:
		pemType = PemTypeCertificate
	}

	pb, err := x509.EncryptPEMBlock(rand.Reader, pemType, in, pwd[:], x509.PEMCipherAES256)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(pb), nil
}

// PEM2DERWithEncryption decode pem to der with password
// if pem is encrypted pem ,pwd is mast not nil
func PEM2DERWithEncryption(raw []byte, pwd *[32]byte) ([]byte, PEMType) {
	if len(raw) == 0 {
		return nil, PEMInvalidPEMType
	}
	//block := new(pem.Block)
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, PEMInvalidPEMType
	}

	if x509.IsEncryptedPEMBlock(block) {
		if pwd == nil {
			return nil, PEMInvalidPEMType
		}
		var err error
		block.Bytes, err = x509.DecryptPEMBlock(block, pwd[:])
		if err != nil {
			return nil, PEMInvalidPEMType
		}
	}

	return block.Bytes, getPemType(block.Type)
}

func getPemType(blockType string) PEMType {
	r := PEMInvalidPEMType
	switch {
	case strings.Contains(blockType, "CERT"):
		r = PEMCertificate
	case strings.Contains(blockType, "RSA"):
		r = PEMRSAPrivateKey
	case strings.Contains(blockType, "EC"):
		r = PEMECCPrivateKey
	case strings.Contains(blockType, "PRIVATE"):
		r = PEMAnyPrivateKey
	case strings.Contains(blockType, "PUB"):
		r = PEMPublicKey
	}
	return r
}
