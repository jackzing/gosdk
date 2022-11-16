package plugin

import (
	"github.com/meshplus/crypto"
)

//KeyStore key store
type KeyStore interface {
	CreateSignKey(persistent bool, mode int) (index []byte, k crypto.SignKey, err error)
	GetSignKey(keyIndex []byte, mode int) (crypto.SignKey, error)
}
