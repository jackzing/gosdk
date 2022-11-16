package tee

import (
	"crypto"
)
import fc "github.com/meshplus/crypto"

//Enclave is an interface of sgx
type Enclave interface {
	Close()
	IsSGX() bool //should always true
}

//SecretStore secret store
type SecretStore interface {
	Enclave
	fc.Cryptor
}

//KeyStore key store
type KeyStore interface {
	Enclave
	crypto.Signer
	Verify(sign, hash []byte) (bool, error)
}

//Oracle oracle
type Oracle interface {
	Enclave
	Fetch(url string, method string, body string, head map[string]string) (Response, error)
}
