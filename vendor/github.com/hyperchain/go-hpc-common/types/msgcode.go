package types

import "github.com/hyperchain/go-hpc-common/p2pcom"

// flato main protocol message code
// Note: please remember to update SortedMsgCode and ModuleNameMap!!
const (
	SessionMsg   = 0x00 //todo this message code may be removed later
	ConsensusMsg = 0x01
	MemPoolMsg   = 0x02
	NonVpMsg     = 0x03
	SyncMgrMsg   = 0x04
	FileMgrMsg   = 0x05
	VPMgrMsg     = 0x06
	EpochMgrMsg  = 0x07
	LENGTH       = 0x10
)

// SortedMsgCode defines the order of module, which is used to read the weight occupied
// by each module from the namespace configuration file in order.
var SortedMsgCode = []uint64{SessionMsg, ConsensusMsg, MemPoolMsg, NonVpMsg, SyncMgrMsg, FileMgrMsg, VPMgrMsg, EpochMgrMsg}

// ModuleNameMap defines the correspondence between message code and module name.
var ModuleNameMap = map[uint64]p2pcom.ModuleName{
	SessionMsg:   p2pcom.SESSION,
	ConsensusMsg: p2pcom.CONSENSUS,
	NonVpMsg:     p2pcom.NONVP,
	SyncMgrMsg:   p2pcom.SYNCMGR,
	FileMgrMsg:   p2pcom.FILEMGR,
	MemPoolMsg:   p2pcom.MEMPOOL,
	VPMgrMsg:     p2pcom.VPMGR,
	EpochMgrMsg:  p2pcom.EPOCHMGR,
}

// In order to prevent forgetting to define, a mandatory check is added to the init function.
func init() {
	if len(SortedMsgCode) != len(ModuleNameMap) {
		panic("please check code about the definition of flato main protocol message code")
	}
	for _, msgCode := range SortedMsgCode {
		if ModuleNameMap[msgCode] == "" {
			panic("not found module name by msg code, please check code")
		}
	}
}

// GetModuleName returns the module name by msg code, returns
// the msg code if not found module name.
func GetModuleName(msgcode uint64) interface{} {
	n, ok := ModuleNameMap[msgcode]
	if ok {
		return n
	}
	return msgcode
}
