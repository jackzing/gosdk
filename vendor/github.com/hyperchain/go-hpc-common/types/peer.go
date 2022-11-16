package types

// PeerArchive is the event passed from archMgr to nvpMgr.
type PeerArchive struct {
	BlockNumber uint64
	Namespace   string
	Ack         chan bool
}

// PeerArchiveFinish is the event passed from nvpMgr to archMgr.
type PeerArchiveFinish struct {
	Namespace string
	Ack       chan bool
}

// PeerRestore is the event passed from archMgr to nvpMgr.
type PeerRestore struct {
	Namespace string
	Ack       chan bool
}

// PeerRestoreFinish is the event passed from nvpMgr to archMgr.
type PeerRestoreFinish struct {
	Namespace string
	Ack       chan bool
}

// PeerReplaceLedger is the event passed from syncMgr to nvpMgr.
type PeerReplaceLedger struct {
	Namespace string
	Ack       chan struct{}
}

// PeerReplaceLedgerFinish is the event passed from nvpMgr to syncMgr.
type PeerReplaceLedgerFinish struct {
	Namespace string
	Ack       chan struct{}
}

// PeerStartNamespace type.
type PeerStartNamespace struct {
	Namespace string
	Ack       chan error
}

// PeerStopNamespace type.
type PeerStopNamespace struct {
	Namespace string
	Ack       chan error
}

// InformStop type.
type InformStop struct {
	Namespace string
	Hostname  string
}

const (
	// DisconnectInformSuccess indicates inform success.
	DisconnectInformSuccess uint32 = iota
	// DisconnectInformFailed indicates inform failed.
	DisconnectInformFailed
)

// PeerDisconnect is the event from nvpMgr to nvpSingle.
type PeerDisconnect struct {
	Ack chan uint32
}

// StatusQuery is the event from peerMgr to remoteSingle.
type StatusQuery struct {
	Ack chan string
}

// RuleQuery is the event from peerMgr to remoteSingle.
type RuleQuery struct {
	Ack chan string
}

// RuleUpdate is the event from peerMgr to remoteSingle.
type RuleUpdate struct {
	Rule string
	Ack  chan error
}

// RelayUpdate is the event from peerMgr to remoteSingle.
type RelayUpdate struct {
	Relay bool
	Ack   chan error
}
