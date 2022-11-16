package types

import (
	"github.com/hyperchain/go-hpc-common/subscription"
	"github.com/hyperchain/go-hpc-common/types/protos"
)

// NewAnchorTx contains to-call anchorTx and a resultCh for execution result
type NewAnchorTx struct {
	Tx       *protos.Transaction
	ResultCh chan *protos.AnchorTx
}

// CrossChainAck used for cm to inform cmc finished checkpoint event
type CrossChainAck struct {
	Namespace string
	Prev      uint64
	Now       uint64
}

// CrossChainContractResult is the contract result, stored in receipt result
type CrossChainContractResult struct {
	Ret []byte

	Success         bool   // 执行成功或失败
	TargetNamespace []byte // 目标调用分区
	Calldata        []byte // 目标合约及调用参数 + 本交易回滚namespace+目标合约及调用参数
}

// cmc to cm

// CrossChainEvent represents a block
type CrossChainEvent struct {
	SeqNo uint64
	Txs   []*protos.CrossChainTx
}

// CrossModuleClientRecover used for cmc recover self status
type CrossModuleClientRecover struct {
	StableCheckpoint uint64
}

// NamespaceStart is generated by cmc and will send this event to cm
type NamespaceStart struct {
	Namespace   string
	CmcReceiver *subscription.TypeMux
	Ack         chan bool
	VFetcher    VersionFetcher // versionFetcher used to get chain version under this namespace
}

// NamespaceStop is generated by cmc and will send this event to cm
type NamespaceStop struct {
	Namespace string
	Ack       chan bool
}

// CrossChainCheckpoint represents a block scope within a checkpoint phase
type CrossChainCheckpoint struct {
	Namespace string
	Prev, Now uint64
	Events    []*CrossChainEvent
}

// InformReplaceLedgerFinish used by cmc to inform cm replace ledger finished
type InformReplaceLedgerFinish struct {
	Namespace          string
	CurrentGenesis     uint64
	GlobalRegisteredNs map[string]uint16
	NormalRegStatus    uint16
	NormalHostname     string
	AckCh              chan error
}

// AnchorTxStage is the anchorTx stage value
type AnchorTxStage int

const (
	// ANCHORDEFAULT is a default value
	ANCHORDEFAULT = iota
	// ANCHORSTAGE1 - global to-register
	ANCHORSTAGE1
	// ANCHORSTAGE2 - normal to-register
	ANCHORSTAGE2
	// ANCHORSTAGE3 - global confirm-register
	ANCHORSTAGE3
	// ANCHORSTAGE4 - normal confirm-register
	ANCHORSTAGE4
	// ANCHORSTAGE5 - stage3 failed, following stage is ANCHORSTAGE5
	ANCHORSTAGE5
	// ANCHORSTAGE6 - global to-unregister
	ANCHORSTAGE6
	// ANCHORSTAGE7 - normal to-unregister
	ANCHORSTAGE7
	// ANCHORSTAGE8 - when cross-chain txs cleared, normal confirm-unregister
	ANCHORSTAGE8
	// ANCHORSTAGE9 - global confirm-unregister
	ANCHORSTAGE9
	// ANCHORSTAGE10 - when tx8 failed, following tx is ANCHORSTAGE10
	ANCHORSTAGE10
	// ANCHORSTAGE11 - replace anchor
	ANCHORSTAGE11
	// ANCHORSTAGE12 - anchor replaced
	ANCHORSTAGE12
	// ANCHORSTAGENOGENERATE - limit value
	ANCHORSTAGENOGENERATE
)

// Increase return a +1 stage value
func (s AnchorTxStage) Increase() AnchorTxStage {
	return AnchorTxStage(int(s) + 1)
}

// CrossChainTxStage represents cross chain tx stage
type CrossChainTxStage int

const (
	// STAGEINVALID is a default value
	STAGEINVALID CrossChainTxStage = iota
	// STAGE1 - hvm, tx1, normal ns
	STAGE1
	// STAGE2 - bvm, tx2, global ns
	STAGE2
	// STAGE3 - hvm, tx3, normal ns
	STAGE3
	// STAGE4 - bvm, tx4, global ns
	STAGE4
	// STAGE5 - hvm, tx5(rollback tx1), normal ns
	STAGE5
	// STAGE6 - bvm, confirm tx5 result, global ns
	STAGE6
	// STAGE7 - bvm, confirm tx3 failed result, global
	STAGE7
	// STAGE8 - hvm, rollback tx3, normal ns
	STAGE8
	// STAGE9 - bvm, confirm tx8 result, global ns
	STAGE9
	// STAGENOGENERATE - mark no proceeding tx
	STAGENOGENERATE
)

// Increase return a +1 stage value
func (s CrossChainTxStage) Increase() CrossChainTxStage {
	return CrossChainTxStage(int(s) + 1)
}

// CrossChainIDSeperator used as CrossChainID seperator
const CrossChainIDSeperator = "__@__"

// DefaultCrossChainEventThreshold used for cmc recover logic, one CrossChainCheckpoint contains DefaultCrossChainEventThreshold event
const DefaultCrossChainEventThreshold = 10

// CMRootPath represents cm root dir
const CMRootPath = "cross_chain"

// AnchorInfoEntry represents entry in anchor.meta
type AnchorInfoEntry struct {
	Namespace      string `json:"ns"`
	RegisterStatus int    `json:"reg_status"`
}

// CheckCMCFinishedCK used for archMgr to check whether current seqNo has been finished in cmc logic
type CheckCMCFinishedCK struct {
	SeqNo uint64
	Ack   chan bool
}
