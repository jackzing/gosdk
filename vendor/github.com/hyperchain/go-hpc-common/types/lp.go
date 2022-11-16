package types

// LPArchive is the event passed from archMgr to nvpMgr.
type LPArchive struct {
	BlockNumber uint64
	Ack         chan bool
}

// LPArchiveFinish is the event passed from nvpMgr to archMgr.
type LPArchiveFinish struct {
	Ack chan bool
}

// LPReplaceLedger is the event passed from syncMgr to nvpMgr.
type LPReplaceLedger struct {
	Ack chan struct{}
}

// LPReplaceLedgerFinish is the event passed from nvpMgr to syncMgr.
type LPReplaceLedgerFinish struct {
	Ack chan struct{}
}

const (
	// LPInformFailed indicates inform failed.
	LPInformFailed uint32 = iota
	// LPInformSuccess indicates inform success.
	LPInformSuccess
)

// LPDisconnect is the event from nvpMgr to nvpSingle.
type LPDisconnect struct {
	Ack chan uint32
}
