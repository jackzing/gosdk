package config

import (
	"github.com/hyperchain/go-hpc-common/bvmcom"
	"github.com/hyperchain/go-hpc-common/p2pcom"
)

//GetMSPClassicConfig Get MSP ClassicConfig
func GetMSPClassicConfig() *MspConfigForTest {
	MSPConfigMock := &MspConfigForTest{}
	MSPConfigMock.EngineRand = "default"
	MSPConfigMock.EngineHash = "default"
	MSPConfigMock.EngineCrypt = "default"
	MSPConfigMock.EngineSign = "default"
	MSPConfigMock.EngineCert = "default"
	MSPConfigMock.EngineSessionKey = "default"
	MSPConfigMock.EngineAccelerate = "default"
	return MSPConfigMock
}

//MspConfigForTest config for test
type MspConfigForTest struct {
	CaMode           bvmcom.CAMode
	EngineRand       string
	EngineHash       string
	EngineCrypt      string
	EngineSign       string
	EngineCert       string
	EngineSessionKey string
	EngineAccelerate string
}

//EncryptionCAPath config for test
func (m *MspConfigForTest) EncryptionCAPath() string {
	//TODO implement me
	panic("implement me")
}

//EncryptionCertPath config for test
func (m *MspConfigForTest) EncryptionCertPath() string {
	//TODO implement me
	panic("implement me")
}

//EncryptionRPCAuthEnable config for test
func (m *MspConfigForTest) EncryptionRPCAuthEnable() bool {
	//TODO implement me
	panic("implement me")
}

//EncryptionSecurityAlgo config for test
func (m *MspConfigForTest) EncryptionSecurityAlgo() string {
	//TODO implement me
	panic("implement me")
}

//SelfIsVP config for test
func (m *MspConfigForTest) SelfIsVP() bool {
	//TODO implement me
	panic("implement me")
}

//NsName config for test
func (m *MspConfigForTest) NsName() string {
	//TODO implement me
	panic("implement me")
}

//GetNewFlag config for test
func (m *MspConfigForTest) GetNewFlag() bool {
	//TODO implement me
	panic("implement me")
}

//GetHostname config for test
func (m *MspConfigForTest) GetHostname() string {
	//TODO implement me
	panic("implement me")
}

//ConsensusAlgo config for test
func (m *MspConfigForTest) ConsensusAlgo() string {
	//TODO implement me
	panic("implement me")
}

//EncryptionGPULibPath config for test
func (m *MspConfigForTest) EncryptionGPULibPath() string {
	//TODO implement me
	panic("implement me")
}

//EncryptionEngineRand rand config
func (m *MspConfigForTest) EncryptionEngineRand() string {
	return m.EngineRand
}

//EncryptionEngineHash hash config
func (m *MspConfigForTest) EncryptionEngineHash() string {
	return m.EngineHash
}

//EncryptionEngineCrypt crypt config
func (m *MspConfigForTest) EncryptionEngineCrypt() string {
	return m.EngineCrypt
}

//EncryptionEngineSign signature config
func (m *MspConfigForTest) EncryptionEngineSign() string {
	return m.EngineSign
}

//EncryptionEngineCert certificate config
func (m *MspConfigForTest) EncryptionEngineCert() string {
	return m.EngineCert
}

//EncryptionEngineSessionKey session key config
func (m *MspConfigForTest) EncryptionEngineSessionKey() string {
	return m.EngineSessionKey
}

//EncryptionEngineACC accelerate config
func (m *MspConfigForTest) EncryptionEngineACC() string {
	return m.EngineAccelerate
}

//GetPeerType config for test
func (m *MspConfigForTest) GetPeerType() p2pcom.PeerType {
	//TODO implement me
	panic("implement me")
}

// CAMode return ca mode.
func (m *MspConfigForTest) CAMode() bvmcom.CAMode {
	return m.CaMode
}

//RootCAs config for test
func (m *MspConfigForTest) RootCAs() []bvmcom.FileInfo {
	//TODO implement me
	panic("implement me")
}
