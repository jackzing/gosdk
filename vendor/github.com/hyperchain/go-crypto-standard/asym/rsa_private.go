package asym

import (
	std "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"io"
)

//RSAPrivateKey RSA private key.
type RSAPrivateKey rsa.PrivateKey

//Bytes marshal
func (key *RSAPrivateKey) Bytes() ([]byte, error) {
	return x509.MarshalPKCS1PrivateKey((*rsa.PrivateKey)(key)), nil
}

//FromBytes unmarshal
func (key *RSAPrivateKey) FromBytes(k []byte, opt int) error {
	rsaKey, err := x509.ParsePKCS1PrivateKey(k)
	if err != nil {
		return err
	}
	key.PublicKey = rsaKey.PublicKey
	key.D = rsaKey.D
	key.Primes = rsaKey.Primes
	key.Precomputed = rsaKey.Precomputed
	return nil
}

//Sign get signature of specific digest
func (key *RSAPrivateKey) Sign(hashType, digest []byte, reader io.Reader) (signature []byte, err error) {
	var hashAlgo = std.Hash(binary.BigEndian.Uint32(hashType))
	return rsa.SignPKCS1v15(reader, (*rsa.PrivateKey)(key), hashAlgo, digest)
}

//GenerateRSAKey generate a pair of rsa key,input is algorithm type
func GenerateRSAKey(opt int) (*RSAPrivateKey, error) {
	switch opt {
	case AlgoRSA2048:
		ret, err := rsa.GenerateKey(rand.Reader, 2048)
		return (*RSAPrivateKey)(ret), err
	case AlgoRSA3072:
		ret, err := rsa.GenerateKey(rand.Reader, 3072)
		return (*RSAPrivateKey)(ret), err
	case AlgoRSA4096:
		ret, err := rsa.GenerateKey(rand.Reader, 4096)
		return (*RSAPrivateKey)(ret), err
	default:
		return nil, errors.New(errIllegalInputParameter)
	}
}

//Public get public key
func (key *RSAPrivateKey) Public() std.PublicKey {
	return &key.PublicKey
}
