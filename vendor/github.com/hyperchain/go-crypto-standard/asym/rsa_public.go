package asym

import (
	std "crypto"
	"crypto/rsa"
	"encoding/asn1"
	"encoding/binary"
	"errors"
)

//RSAPublicKey RSA public key.
// never new(RSAPublicKey), use NewRSAPublicKey()
type RSAPublicKey rsa.PublicKey

//Bytes marshal
func (key *RSAPublicKey) Bytes() ([]byte, error) {
	return asn1.Marshal(*key)
}

//FromBytes unmarshal
func (key *RSAPublicKey) FromBytes(k []byte, mode int) error {
	rest, err := asn1.Unmarshal(k, key)
	if err != nil {
		return err
	}
	if len(rest) != 0 {
		return errors.New("trailing data after RSA public key")
	}
	if key.N.Sign() <= 0 {
		return errors.New("x509: RSA modulus is not a positive number")
	}
	if key.E <= 0 {
		return errors.New("x509: RSA public exponent is not a positive number")
	}
	return err
}

// Verify verify the signature by ECDSAPublicKey self, so the first parameter will be ignored.
func (key *RSAPublicKey) Verify(hashType, signature, digest []byte) (bool, error) {
	var hashAlgo = std.Hash(binary.BigEndian.Uint32(hashType))
	err := rsa.VerifyPKCS1v15((*rsa.PublicKey)(key), hashAlgo, digest, signature)
	if err != nil {
		return false, err
	}
	return true, nil
}
