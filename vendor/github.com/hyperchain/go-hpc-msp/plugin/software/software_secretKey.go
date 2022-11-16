package software

import (
	"crypto/rand"
	"io"

	"github.com/meshplus/crypto"
)

//SecretKey secret Key
type SecretKey struct {
	key []byte
	c   crypto.Cryptor
}

//Encrypt encrypt
func (s *SecretKey) Encrypt(src []byte, reader io.Reader) []byte {
	if len(s.key) == 0 {
		panic("key length is zero")
	}
	if c, err := s.c.Encrypt(s.key, src, reader); err == nil {
		return c
	}
	return nil
}

//Decrypt decrypt
func (s *SecretKey) Decrypt(src []byte) []byte {
	if len(s.key) == 0 {
		panic("key length is zero")
	}
	if p, err := s.c.Decrypt(s.key, src); err == nil {
		return p
	}
	return nil
}

//Destroy destroy
func (s *SecretKey) Destroy() {
	_, _ = rand.Read(s.key)
	s.key = nil
	s.c = nil
}

//MockSecretKey use in pure encrypt
type MockSecretKey struct {
}

//Encrypt encrypt
func (s *MockSecretKey) Encrypt(src []byte, reader io.Reader) []byte {
	return src
}

//Decrypt decrypt
func (s *MockSecretKey) Decrypt(src []byte) []byte {
	return src
}

//Destroy destroy Key
func (s *MockSecretKey) Destroy() {
}
