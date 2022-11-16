package modifycache

// ModifyCache is used to cache changes of builtin key in stateDB with different SeqNo. It is concurrent-safe.
//
// Changes are generated from stateDB (checkout between applyTransaction and commitLeger)
//
// It will be used as a component in Executor/Validator/Committer in flato to cache changes.
// It can also be passed to flato-cm to confirm changes
//
// ModifyCache will be reset while DBCommit is called.
type ModifyCache interface {
	// Set sets the ModifySet checkout from stateDB into cache
	Set(ModifySet)
	// Get returns the ModifySet with given SeqNo
	Get(seqNo uint64) (ModifySet, bool)
	// Del deletes the ModifySet with given SeqNo
	Del(seqNo uint64)
	// TrimBefore remove modify set before SeqNo (include SeqNo)
	TrimBefore(seqNo uint64)
	// Reset clear all modify sets
	Reset()
	// Len returns the length of modify cache
	Len() int
}

// ModifySet is used to cache changes of builtin key in stateDB. It isn't concurrent-safe.
//
// It's a component of stateDB (checkout and reset between applyTransaction and commitLeger)
type ModifySet interface {
	// Set set the modify items with given key, old value and new value
	Set(key string, oldVal, newVal []byte)
	// Get returns the modify items with given key
	Get(key string) (item ModifyItem, exist bool)
	// GetAll returns the all modify items
	GetAll() map[string]ModifyItem
	// GetSeqNo returns the sequence number of ModifySet
	GetSeqNo() uint64
	// SetSeqNo sets the sequence number of ModifySet
	SetSeqNo(uint64)
	// Len returns the length of ModifySet
	Len() int
	// Prune removes all modify items with the same old value and new value
	Prune()
	// Del deletes certain key of modify set
	Del(key string)
	// PersistToFile persist ModifySet into file
	PersistToFile(path string) error
	// RecoverFromFile recover ModifySet from file
	RecoverFromFile(path string) error
}

// ModifyItem store change details
// Each change of builtin storage in stateDB will generate a item log the Name/OldVal/NewVal
type ModifyItem interface {
	// SetNewVal sets new value
	SetNewVal([]byte)
	// SetOldVal sets old value
	SetOldVal([]byte)

	// GetNewVal returns new value
	GetNewVal() []byte
	// GetNewVal returns old value
	GetOldVal() []byte
	// GetName returns Name
	GetName() string
}
