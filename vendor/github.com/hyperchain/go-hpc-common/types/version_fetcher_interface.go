package types

// VersionFetcher defines a series of method to get relevant version from ledger.
type VersionFetcher interface {
	// GetTxVersion returns the currently running tx version from db or trusted config.
	GetTxVersion() (txVersion []byte, err error)
	// GetBlockVersion returns the currently running block version from db.
	GetBlockVersion() (blkVersion []byte, err error)
	// GetConsensusVersion returns the currently running consensus version from db.
	GetConsensusVersion() (consensusVersion []byte, err error)

	// GetSelfSupportedVersion returns supported version list of this node
	GetSelfSupportedVersion() SupportedVersion
	// IsCompatible checks if this binary supports all running version stored in db or trusted config.
	// If shouldCheckConsensus is false, it will ignore to check consensus version.
	// It returns false if the versions are incompatible or the versions are nil, otherwise, returns true.
	IsCompatible(shouldCheckConsensus bool) bool
}
