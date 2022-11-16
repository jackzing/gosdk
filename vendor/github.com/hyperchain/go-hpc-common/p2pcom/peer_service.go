package p2pcom

import (
	"strings"

	"github.com/hyperchain/go-hpc-common/types/protos"
	"github.com/hyperchain/go-hpc-common/utils"

	"golang.org/x/crypto/sha3"
)

// PeerType represents the type of peer
type PeerType int

// the definition of peer type
const (
	VP  PeerType = iota // validate peer
	NVP                 // non-validate peer
	CVP                 // candidate peer
	SPV
	LP // light peer

	ALL
	UNKNOWN
)

func (pt PeerType) String() string {
	switch pt {
	case VP:
		return "VP"
	case NVP:
		return "NVP"
	case CVP:
		return "CVP"
	case SPV:
		return "SPV"
	case LP:
		return "LP"
	case ALL:
		return "ALL"
	default:
		return "not known peer type"
	}
}

// ConvertStringToPeerType convert string type to PeerType.
func ConvertStringToPeerType(t string) (tp PeerType) {
	switch strings.ToUpper(t) {
	case "VP":
		tp = VP
	case "NVP":
		tp = NVP
	case "CVP":
		tp = CVP
	case "SPV":
		tp = SPV
	case "ALL":
		tp = ALL
	case "LP":
		tp = LP
	default:
		tp = UNKNOWN
	}
	return
}

// ModuleName represents the id of module
type ModuleName string

const (
	// CONSENSUS is concensus module
	CONSENSUS ModuleName = "CONSENSUS"
	// MEMPOOL is mempool module
	MEMPOOL ModuleName = "MEMPOOL"
	// NONVP is nvp module
	NONVP ModuleName = "NONVP"
	// SYNCMGR is sync chain module
	SYNCMGR ModuleName = "SYNCMGR"
	// FILEMGR is filemgr sync file module
	FILEMGR ModuleName = "FILEMGR"
	// VPMGR is vp mgr module
	VPMGR ModuleName = "VPMGR"
	// EPOCHMGR is epoch mgr module
	EPOCHMGR ModuleName = "EPOCHMGR"
	// SESSION is session module，it may be refactored later
	SESSION ModuleName = "SESSION"
)

// String returns string type.
func (mn ModuleName) String() string {
	return string(mn)
}

// MsgWriter provides writing of messages. Implementation
// should ensure that WriteMsg can be called simultaneously
// from multiple goroutines.
type MsgWriter interface {
	WriteMsg(msgcode uint64, payload []byte) error
	WriteMsgSync(msgcode uint64, payload []byte) error
	WriteMsgByGossip(msgcode uint64, dataType protos.GossipDataType, data *protos.GossipNsData) error
}

// P2PService provides p2p service
type P2PService interface {
	// Reconnect is used to reconnect peer
	Reconnect()
}

// PeerService represents peer service in hyperchain protocol. If a upper module
// registered peer service, it would receive a peer which has established connection
// through AddPeer() callback. When this peer disconnected, it would receive a remove
// request through RemovePeer() callback.
type PeerService struct {
	// ModuleName module name
	ModuleName ModuleName

	// Type the peer type that module wants to receive.
	Type PeerType

	// Versions represents the supported version list for this module.
	// This version list will be sent to the remote peer by p2p so that to
	// negotiate the shared version of this module.
	Versions []string

	// VersionMatch is a callback implemented by upper module but called by
	// p2p after finishing logical handshaking. It is used to get the version
	// number of the module.
	//
	// The parameters in VersionMatch:
	// 1. selfVersions is the local node supported version.
	// 2. remoteVersions is the remote peer supported version.
	// The returned value indicates the final selected version number for this module.
	VersionMatch func(selfVersions, remoteVersions []string) (string, error)

	// AddPeer will be call when establish logic connect with a peer.
	//
	// the parameters in AddPeer:
	// 1. peerBasicInfo is the remote peer basic info
	// 2. writer is used to write message to other peer
	// if the returned error is not nil, it will print an error log.
	AddPeer func(peerBasicInfo *PeerBasicInfo, writer MsgWriter) error

	// RemovePeer will be call when disconnect with this peer hash.
	//
	// the parameters in RemovePeer:
	// 1. hash is the peer hash which will been disconnected.
	// if the returned error is not nil, it will print an error log.
	RemovePeer func(hash string) error

	// RecvMsg will be call when receive msg to self module from other peer every time.
	//
	// RecvMsg function's implementation should make sure safe in concurrent call
	// and if a not nil error return, will disconnect physical connect and retry
	// so make sure return a not nil error when module suffer serious trouble and need peer reconnect.
	//
	// the parameters in RecvMsg:
	// 1. payload is the msg content
	// 2. hash is the peer hash which is the msg from
	// 3. hostname is the peer hostname which is the msg from
	// if the returned error is not nil, it will disconnect logical connection and
	// then automatically reconnect after a period of time.
	RecvMsg func(payload []byte, hash string, hostname string) error
}

// PeerBasicInfo defines the basic infos of the p2p peer.
type PeerBasicInfo struct {
	Hostname string
	Hash     string
	PeerType PeerType
	Versions map[ModuleName]string // the value is the largest version negotiated by the two nodes.
	P2PService
}

// GetHash return the peer's hash.
func (pbi *PeerBasicInfo) GetHash() string {
	return pbi.Hash
}

// GetHostname return the peer's hostname.
func (pbi *PeerBasicInfo) GetHostname() string {
	return pbi.Hostname
}

// IsVP returns the peer's role, if the peer is vp it returns
// true, otherwise returns false. false means the peer may
// be a nvp, cvp or spv.
func (pbi *PeerBasicInfo) IsVP() bool {
	return pbi.PeerType == VP
}

// GetPeerType returns the peer type.
func (pbi *PeerBasicInfo) GetPeerType() PeerType {
	return pbi.PeerType
}

// GetVersion returns versions negotiated by the two nodes.
func (pbi *PeerBasicInfo) GetVersion() map[ModuleName]string {
	return pbi.Versions
}

// ReconnectPeer is used to reconnect all peers actively.
func (pbi *PeerBasicInfo) ReconnectPeer() {
	pbi.P2PService.Reconnect()
}

// GetPeerHash calculate the (hostname, namespace) hash, return a hex string
func GetPeerHash(hostname string, namespace string) string {
	hasher := sha3.NewLegacyKeccak256()
	//nolint
	hasher.Write([]byte(hostname + namespace))
	dst := hasher.Sum(nil)
	return utils.BytesToHex(dst)
}
