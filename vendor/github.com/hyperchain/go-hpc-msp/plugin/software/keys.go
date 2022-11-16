package software

import (
	"fmt"

	"github.com/meshplus/crypto"
)

// UnmarshalPrivateKey unmarshals a pkcs8 der to private key
func UnmarshalPrivateKey(index string) (key crypto.SignKey, err error) {
	engine := GetSoftwareEngine()
	k, err := engine.GetSignKey(index)
	if err == nil {
		return k, nil
	}
	return
}

// MarshalPublicKey marshal a public key to the pem forma
func MarshalPublicKey(publicKey crypto.VerifyKey) ([]byte, error) {
	return MarshalPKIXPublicKey(publicKey.Bytes(), publicKey.GetKeyInfo())
}

// UnmarshalPublicKey unmarshal a der to public key
func UnmarshalPublicKey(derBytes []byte) (pub crypto.VerifyKey, err error) {
	//parse der
	rawpub, mode, err := ParsePKIXPublicKey(derBytes)
	if err != nil {
		return nil, fmt.Errorf("parse pkix pub error: %v", err.Error())
	}
	engine := GetSoftwareEngine()
	return engine.GetVerifyKey(rawpub, mode)
}
