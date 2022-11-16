package types

// NodeStateRequest is the request to get all nodes' state
type NodeStateRequest struct {
	Hash string
	Key  int64
}

// NodeStateResponse is the response to get all nodes' state
type NodeStateResponse struct {
	Key  int64
	Info NodeStateInfo
}

// NodeStateInfo records the node status(including consensus status)
type NodeStateInfo struct {
	ID          uint64 `json:"id"`
	Hostname    string `json:"hostname"` // hostname of node
	Hash        string `json:"hash"`
	Status      string `json:"status"` // TIMEOUT, NORMAL, VIEWCHANGE...
	View        uint64 `json:"view"`
	Epoch       uint64 `json:"epoch"`
	Checkpoint  uint64 `json:"checkpoint"`  // latest checkpoint height of node
	BlockHeight uint64 `json:"blockHeight"` // latest block height of node
	BlockHash   string `json:"blockHash"`   // latest block hash of node
	Version     string `json:"version"`     // commercial version of node
}

// NodeInfo records the node status about consensus
type NodeInfo struct {
	ID        uint64
	Namespace string
	Hash      string
	Hostname  string
	IsPrimary bool
	IsVP      bool
}

// StateList records the list of all nodes' state info
type StateList []NodeStateInfo

// NodesList records the node info
type NodesList []NodeInfo

// SystemStatus records the system status(including consensus state)
type SystemStatus struct {
	Status     uint64
	View       uint64
	Epoch      uint64
	ReplicaID  uint64
	Checkpoint uint64
}

// constants for node status
const (
	TIMEOUT = iota
	SystemNormal
	ConsensusViewChange
	ConsensusRecovery
	ConsensusTxPoolFull
	ConsensusTxSetFull
	ConsensusConfigChange
	ConsensusStateUpdate
	ConsensusPending
)

// String returns a long description of SystemStatus
func (status SystemStatus) String() string {
	switch status.Status {
	case TIMEOUT:
		return "TIMEOUT"
	case SystemNormal:
		return "NORMAL"
	case ConsensusViewChange:
		return "system is in view change"
	case ConsensusRecovery:
		return "system is in recovery"
	case ConsensusConfigChange:
		return "system is in config change"
	case ConsensusTxPoolFull:
		return "system is too busy"
	case ConsensusTxSetFull:
		return "system is too busy"
	case ConsensusStateUpdate:
		return "system is in state update"
	case ConsensusPending:
		return "system is in pending state"
	default:
		return "Unknown status"
	}
}

// Description returns a short description of SystemStatus
func (status SystemStatus) Description() string {
	switch status.Status {
	case TIMEOUT:
		return "TIMEOUT"
	case SystemNormal:
		return "NORMAL"
	case ConsensusViewChange:
		return "VIEW_CHANGE"
	case ConsensusRecovery:
		return "RECOVERY"
	case ConsensusConfigChange:
		return "CONFIG_CHANGE"
	case ConsensusTxPoolFull:
		return "TX_POOL_FULL"
	case ConsensusTxSetFull:
		return "TX_SET_FULL"
	case ConsensusStateUpdate:
		return "STATE_UPDATE"
	case ConsensusPending:
		return "ConsensusPending"
	default:
		return "UNKNOWN"
	}
}

// GetView returns the current view
func (status SystemStatus) GetView() uint64 {
	return status.View
}

// GetEpoch returns the current view
func (status SystemStatus) GetEpoch() uint64 {
	return status.Epoch
}

// GetID returns the node's id
func (status SystemStatus) GetID() uint64 {
	return status.ReplicaID
}

// GetCheckpoint returns the node's checkpoint
func (status SystemStatus) GetCheckpoint() uint64 {
	return status.Checkpoint
}

func (sl StateList) Len() int {
	return len(sl)
}

func (sl StateList) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

func (sl StateList) Less(i, j int) bool {
	return sl[i].ID < sl[j].ID
}
