package types

import (
	"fmt"
	"github.com/hyperchain/go-hpc-common/utils"
	json "github.com/json-iterator/go"
)

// KeyValueMap a map definition
type KeyValueMap map[string]interface{}

// NewKeyValueMap new a key value map instance
func NewKeyValueMap() KeyValueMap {
	return make(map[string]interface{})
}

// Add a key into the map
func (m KeyValueMap) Add(key string, val interface{}) {
	if val == nil {
		val = uint64(0)
	}
	m[key] = val
}

// Exist judge key exist in the map or not
func (m KeyValueMap) Exist(key string) bool {
	_, ok := m[key]
	return ok
}

// Get value from keyValue map
func (m KeyValueMap) Get(key string, val interface{}) error {
	enc, ok := m[key]
	if !ok {
		return fmt.Errorf("key %s not exist", key)
	}
	if val == nil {
		return fmt.Errorf("val is nil")
	}
	valBytes, err := json.Marshal(enc)
	if err != nil {
		return err
	}
	return json.Unmarshal(valBytes, val)
}

// ContractApprover contract approver struct
type ContractApprover struct {
	ApproverList   []Address `json:"approverList"`
	NeededApproval int       `json:"neededApproval"`
}

// PrivateRawData private raw data struct
type PrivateRawData struct {
	Collection      []string `json:"collection"`
	Payload         string   `json:"payload"`
	PublicSignature string   `json:"publicSignature"`
}

// ContractUpgradeExtra contract upgrade extra info struct
type ContractUpgradeExtra struct {
	ContractAddr         Address  `json:"contractAddr"`
	TxHash               [32]byte `json:"txHash"`
	ContractProposalType int64    `json:"contractProposalType"`
}

// SystemUpgradeExtra system upgrade extra info struct
type SystemUpgradeExtra struct {
	Extra         []byte `json:"extra"`
	IsHardUpgrade bool   `json:"isHardUpgrade"`
}

// CAManageExtra ca manager extra info struct
type CAManageExtra struct {
	// CaSubType ca subscribe type
	CaSubType uint8 `json:"caSubType"` // 0: enterACO, 1: leaveACO, 2: cancel proposal
	// MemberAddr member address
	MemberAddr Address `json:"memberAddr"`
	// TargetID target identify
	TargetID int64 `json:"targetId"`
}

const (
	// ExtraVersionKey extra version key
	ExtraVersionKey = "__version__"
	// ExtraVersionValue extra version value
	ExtraVersionValue = "1.0"

	// ContractCreationKey extra key for contract creation
	ContractCreationKey = "__contractCreation__" // extra key for contract creation
	// PrivateDataKey extra key for private raw data
	PrivateDataKey = "__privateData__" // extra key for private raw data
	// PrivateDataKeyPubSig private data key public signature
	PrivateDataKeyPubSig = "publicSignature"
	// PrivateDataKeyCollection private data collection
	PrivateDataKeyCollection = "collection"
	// PrivateDataKeyPayload private data key payload
	PrivateDataKeyPayload = "payload"
	// ContractUpgradeKey contract upgrade key
	ContractUpgradeKey = "__contractUpgrade__"
	// SystemUpgradeKey system upgrade key
	SystemUpgradeKey = "__systemUpgrade__"
	// CAManageKey ca manager key
	CAManageKey = "__caManage__"
)

// GenerateContractApproverExtra generate the contract approver extra
func GenerateContractApproverExtra(approverList []Address, neededApproval int) ([]byte, error) {
	kv := NewKeyValueMap()
	contractApprover := ContractApprover{
		ApproverList:   approverList,
		NeededApproval: neededApproval,
	}
	kv.Add(ContractCreationKey, contractApprover)
	enc, err := json.Marshal(kv)
	if err != nil {
		return nil, err
	}
	return []byte(utils.BytesToHex(enc)), nil
}

// ParseContractApproverExtra parse extra field and returns approverList and neededApproval
func ParseContractApproverExtra(extra []byte) ([]Address, int) {
	contractApprover := &ContractApprover{}
	extra = utils.HexToBytes(string(extra))
	entries := &KeyValueMap{}
	err := json.Unmarshal(extra, entries)
	if err == nil {
		err = entries.Get(ContractCreationKey, contractApprover)
		if err == nil {
			if len(contractApprover.ApproverList) > 0 && contractApprover.NeededApproval > 0 && contractApprover.NeededApproval <= len(contractApprover.ApproverList) {
				return contractApprover.ApproverList, contractApprover.NeededApproval
			}
		}
	}
	return nil, -1
}
