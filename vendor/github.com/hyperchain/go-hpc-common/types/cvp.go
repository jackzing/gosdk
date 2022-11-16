package types

import (
	"github.com/hyperchain/go-hpc-common/modifycache"
	"github.com/hyperchain/go-hpc-common/types/protos"
)

// BJR is the memory data,
// parsed from remote fetchResp, and provided by vp
type BJR struct {
	SeqNo        uint64
	BlockBytes   []byte
	JournalBytes []byte
	ReceiptBytes []byte
	NonApply     bool
}

// InitVP struct
type InitVP struct {
	Payload *protos.CvpPayload
	Ack     chan error
}

// ModifyNamespace struct
type ModifyNamespace struct {
	Namespace string
	Enable    bool
	NsPayload *protos.NamespacePayload
	Ack       chan error
}

// ModifyHosts struct
// @Deprecated
type ModifyHosts struct {
	Host, Addr string
	Ack        chan error
}

// ModifyDomain struct
type ModifyDomain struct {
	Delete       bool
	Domain, Addr string
	Ack          chan error
}

// ModifyAuth struct
type ModifyAuth struct {
	Namespace string
	Method    string
	Rules     string
	Ack       chan error
}

// TxModify struct
type TxModify struct {
	Namespace string
	MSet      modifycache.ModifySet
	Ack       chan *TxModifyResp
}

// ConfCmdModify struct
type ConfCmdModify struct {
	List *protos.CvpConfCmdList
	Ack  chan error
}

// TxModifyResp struct
type TxModifyResp struct {
	MSet modifycache.ModifySet
	Err  error
}

// ReplaceConfig struct
type ReplaceConfig struct {
	Ack chan error
}

// IPCArgs struct
type IPCArgs struct {
	Args     []string
	Ret      *[]string
	ExitCode *int
}

// APIArgs struct
type APIArgs struct {
	Args     []string
	Payload  interface{}
	Result   *[]interface{}
	ExitCode *int
}

// MQAPIArgsPayload struct
type MQAPIArgsPayload struct {
	Namespace string
	Meta      *RegisterMeta
}

const (
	// CvpAddHost indicates AddHost operation.
	CvpAddHost = "AddHost"
	// CvpDeleteHost indicates DeleteHost operation.
	CvpDeleteHost = "DeleteHost"
	// CvpAddNvp indicates AddNvp operation.
	CvpAddNvp = "AddNvp"
	// CvpDeleteNvp indicates DeleteNvp operation.
	CvpDeleteNvp = "DeleteNvp"
	// CvpAddNode indicates CvpAddNode operation.
	CvpAddNode = "CvpAddNode"
	// CvpDeleteNode indicates CvpDeleteNode operation.
	CvpDeleteNode = "CvpDeleteNode"
	// CvpAddLp indicates AddNvp operation.
	CvpAddLp = "AddLp"
	// CvpDeleteLp indicates DeleteNvp operation.
	CvpDeleteLp = "DeleteLp"
)

// ConfCmd struct
type ConfCmd struct {
	Method string
	Params []string
	OldVal interface{}
	Key    string
	Ack    chan error
}

// CmdRespWithAck struct
type CmdRespWithAck struct {
	Resp *protos.CmdResp
	Ack  chan struct{}
}

// CVPConfCmdList struct.
type CVPConfCmdList struct {
	Namespace string
	Lists     []*ConfCmd
}

// CVPCertsInfo struct.
type CVPCertsInfo struct {
	Namespace string
}

const (
	// CVP_PROVIDER is the provider of cvp.
	CVP_PROVIDER string = "provider"
	// CVP_RECEIVER is the receiver of cvp.
	CVP_RECEIVER string = "receiver"
)
