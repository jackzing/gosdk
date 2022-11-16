package cryptocom

import (
	"github.com/meshplus/crypto"
)

//Engine crypto engine
type Engine interface {
	crypto.PluginRandomFunc
	crypto.PluginHashFunc
	crypto.PluginCryptFunc
	crypto.PluginSignFuncL3
	crypto.PluginGenerateSessionKeyFunc
	PluginAccelerateFunc
	SetAlgo([]byte) error
	GetDefaultAlgo() (int, int)
}

//PluginAccelerateFunc plugin interface
type PluginAccelerateFunc interface {
	crypto.Level
	GetVerifyACCKey(mode int) (VerifyACCKey, error)
}

//VerifyACCKey plugin interface
type VerifyACCKey interface {
	VerifyACC(key, sign, hashRet [][]byte) error
	Destroy() error
}
