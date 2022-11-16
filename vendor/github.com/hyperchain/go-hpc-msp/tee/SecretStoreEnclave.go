package tee

import "C"
import (
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

//secretStorageEnclaveImpl implement of enclave for key storage
type secretStorageEnclaveImpl struct {
}

//NewSecretStoreEnclave new SecretStoreEnclave
func NewSecretStoreEnclave() (SecretStore, error) {
	return &secretStorageEnclaveImpl{}, nil
}

func toAbs(pathFromConfig string) string {
	if strings.HasPrefix(pathFromConfig, "http") {
		return pathFromConfig
	}
	if path.IsAbs(pathFromConfig) {
		return pathFromConfig
	}
	pwd, _ := os.Getwd()
	return path.Join(pwd, pathFromConfig)
}

//Close close
func (e *secretStorageEnclaveImpl) Close() {
	if eid == 0 {
		return
	}
	close(eid)
}

//IsSGX always true
func (e *secretStorageEnclaveImpl) IsSGX() bool {
	return true
}

//Encrypt encrypt
func (e *secretStorageEnclaveImpl) Encrypt(_, plaintext []byte, _ io.Reader) (cipherText []byte, err error) {
	if len(plaintext) == 0 {
		return plaintext, nil
	}
	if eid == 0 {
		return nil, errors.New("enclave need init")
	}
	return encrypt(eid, plaintext)
}

//Decrypt remote attest
func (e *secretStorageEnclaveImpl) Decrypt(_, cipherText []byte) (plaintext []byte, err error) {
	if len(cipherText) == 0 {
		return cipherText, nil
	}
	if eid == 0 {
		return nil, errors.New("enclave need init")
	}
	return decrypt(eid, cipherText)
}
