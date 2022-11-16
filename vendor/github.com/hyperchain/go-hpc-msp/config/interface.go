package config

import (
	"github.com/hyperchain/go-hpc-common/bvmcom"
	"github.com/hyperchain/go-hpc-common/p2pcom"
)

//EncryptionConfigInterface Encryption Config Interface
type EncryptionConfigInterface interface {
	EncryptionCAPath() string
	EncryptionCertPath() string
	EncryptionRPCAuthEnable() bool
	EncryptionSecurityAlgo() string

	SelfIsVP() bool
	NsName() string
	GetNewFlag() bool
	GetHostname() string
	ConsensusAlgo() string
	EncryptionGPULibPath() string

	EncryptionEngineRand() string
	EncryptionEngineHash() string
	EncryptionEngineCrypt() string
	EncryptionEngineSign() string
	EncryptionEngineCert() string
	EncryptionEngineSessionKey() string
	EncryptionEngineACC() string

	GetPeerType() p2pcom.PeerType
	CAMode() bvmcom.CAMode
	RootCAs() []bvmcom.FileInfo
}
