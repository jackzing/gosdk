package types

// version_chain.go used to define version information of chain layer.

// VersionTag defines version classification.
type VersionTag string

// Chain level
const (
	// TXVersionTag is the label of tx version, you must upgrade tx version
	// if the transaction processing logic changes which influences result
	// of execution.
	TXVersionTag VersionTag = "tx_version"
	// BlockVersionTag is the label of block version, you must upgrade block
	// version if the block processing logic changes, for example, the algorithm
	// of root hash in block head changes.
	BlockVersionTag VersionTag = "block_version"
	// EncodeVersionTag is the label of hash algorithm version.
	EncodeVersionTag VersionTag = "encode_version" // hash algo

	// ConsensusVersionTag is the label of consensus version.
	ConsensusVersionTag VersionTag = "consensus_version"
)

// from release 1.4.* to 1.6.1(not contain), was use 1.3
// from release 1.6.1 to 1.6.17, 1.7.1 to 1.7.8 was use 1.4
// for 1.6.18+, 1.8.0 to 1.8.1(not contain) use 1.5
// for 1.8.2-1.8.4, use 1.6  (this for hvm timestamp and evm no code contract #hpc-1047 #hpc-1060)
// for 1.8.5+, use 1.7 (this for evm update #hpc-1142)
// for 1.8.7+  use 1.8 (this for log index bug fix hpc-1326 2020/03/23)
// flato 0.0.3 to 0.0.4, use 2.0
// flato 0.0.5, use 2.1
// flato 0.0.7, use 2.2
// flato 0.1.0, use 2.3 (this for evm suicide logic change)
// flato 1.0.2, use 2.4 (for hvm compatible and choose sync_chain mode by tx version)
// flato 1.0.4, use 2.5 (if index data stores into meta)
// flato 1.0.7, use 2.6 (add did transaction)
// flato 1.0.8, use 2.7 (for bvm compatible)
// flato 1.0.9, use 2.8 (return gas used for all transaction)
// flato 1.0.10, use 2.9 (calculate gas when json marshal the result of hvm)
// flato 1.0.11, use 3.0 (for hvm optimize)
// flato 1.0.12, use 3.1 (kvsql pipeline execute)
// flato 1.2.0, block version use 3.2, tx meta use 1(add invalid tx to block)
// flato 1.3.0, use tx version 3.3 (prohibited some operations of float in hvm and prohibited float in kvsql)
// flato 1.4.0, use tx version 3.4 (refactor kvsql gas)
// flato 1.5.0, use tx version 3.5 (add tx root merkle calculation)
// hyperchain 2.7.0, use tx version 3.6(supplementary configuration on-chain function in bvm)
// hyperchain 2.9.0, use tx version 4.0 (1.support consensus switching. 2. implement the draft of system upgrade)
// hyperchain 2.10.0 use tx version 4.1 (1.support did refactor. 2. support custom hash algo and custom symmetric encryption algo)
// NOTE: please remember to update map HyperchainReleaseVersion and ChainVersionList !!!!
const (
	HyperchainVersion290  = "2.9.0"
	HyperchainVersion2100 = "2.10.0"
)

const (
	// DefaultSOLOGenesisRunningVersion is used in flato-config
	DefaultSOLOGenesisRunningVersion = HyperchainVersion2100
)

// HyperchainReleaseVersion used to define relevant chain version for
// every hyperchain open version.
// NOTE: this map can only add new value, but not update or delete old value.
var HyperchainReleaseVersion = map[string]RunningVersion{
	HyperchainVersion290: {
		TXVersionTag:     TxVersion40,
		BlockVersionTag:  BlkVersion40,
		EncodeVersionTag: encodeVersionMap[EncodeVersion00],
	},
	HyperchainVersion2100: {
		TXVersionTag:     TxVersion41,
		BlockVersionTag:  BlkVersion40,
		EncodeVersionTag: encodeVersionMap[EncodeVersion00],
	},
}

// ChainVersionList defines all versions supported by this hyperchain binary.
// NOTE: all history version should be preserved for chain data. We allow
// preserve part of history version for chain protocol(for example: consensus version).
var ChainVersionList = map[VersionTag]map[string]string{ //
	TXVersionTag: {
		TxVersion20: TxVersion20,
		TxVersion21: TxVersion21,
		TxVersion22: TxVersion22,
		TxVersion23: TxVersion23,
		TxVersion24: TxVersion24,
		TxVersion25: TxVersion25,
		TxVersion26: TxVersion26,
		TxVersion27: TxVersion27,
		TxVersion28: TxVersion28,
		TxVersion29: TxVersion29,
		TxVersion30: TxVersion30,
		TxVersion31: TxVersion31,
		TxVersion32: TxVersion32,
		TxVersion33: TxVersion33,
		TxVersion34: TxVersion34,
		TxVersion35: TxVersion35,
		TxVersion36: TxVersion36,
		TxVersion40: TxVersion40,
		TxVersion41: TxVersion41,
	},
	BlockVersionTag: {
		BlkVersion15: BlkVersion15,
		BlkVersion16: BlkVersion16,
		BlkVersion17: BlkVersion17,
		BlkVersion18: BlkVersion18,
		BlkVersion20: BlkVersion20,
		BlkVersion21: BlkVersion21,
		BlkVersion22: BlkVersion22,
		BlkVersion23: BlkVersion23,
		BlkVersion24: BlkVersion24,
		BlkVersion25: BlkVersion25,
		BlkVersion26: BlkVersion26,
		BlkVersion27: BlkVersion27,
		BlkVersion28: BlkVersion28,
		BlkVersion29: BlkVersion29,
		BlkVersion30: BlkVersion30,
		BlkVersion31: BlkVersion31,
		BlkVersion32: BlkVersion32,
		BlkVersion33: BlkVersion33,
		BlkVersion34: BlkVersion34,
		BlkVersion35: BlkVersion35,
		BlkVersion36: BlkVersion36,
		BlkVersion40: BlkVersion40,
	},
	EncodeVersionTag: {
		encodeVersionMap[EncodeVersion00]: encodeVersionMap[EncodeVersion00],
	},
}

// tx version in string
const (
	// TxVersion20 tx version 2.0
	TxVersion20 = "2.0"

	// TxVersion21 tx version 2.1
	TxVersion21 = "2.1"

	// TxVersion22 tx version 2.2
	TxVersion22 = "2.2"

	// TxVersion23 tx version 2.3
	TxVersion23 = "2.3"

	// TxVersion24 tx version 2.4
	TxVersion24 = "2.4"

	// TxVersion25 tx version 2.5
	TxVersion25 = "2.5"

	// TxVersion26 tx version 2.6
	TxVersion26 = "2.6"

	// TxVersion27 tx version 2.7
	TxVersion27 = "2.7"

	// TxVersion28 tx version 2.8
	TxVersion28 = "2.8"

	// TxVersion29 tx version 2.9
	TxVersion29 = "2.9"

	// TxVersion30 tx version 3.0
	TxVersion30 = "3.0"

	// TxVersion31 tx version 3.1
	TxVersion31 = "3.1"

	// TxVersion32 tx version 3.2
	TxVersion32 = "3.2"

	// TxVersion33 tx version 3.3
	TxVersion33 = "3.3"

	// TxVersion34 tx version 3.4
	TxVersion34 = "3.4"

	// TxVersion35 tx version 3.5
	TxVersion35 = "3.5"

	// TxVersion36 tx version 3.6
	TxVersion36 = "3.6"

	// TxVersion40 tx version 4.0
	TxVersion40 = "4.0"

	// TxVersion41 tx version 4.1
	TxVersion41 = "4.1"
)

// block version in string
const (
	// BlkVersion15 tx version 1.5
	BlkVersion15 = "1.5"

	// BlkVersion16 tx version 1.6
	BlkVersion16 = "1.6"

	// BlkVersion17 tx version 1.7
	BlkVersion17 = "1.7"

	// BlkVersion18 tx version 1.8
	BlkVersion18 = "1.8"

	// BlkVersion20 tx version 2.0
	BlkVersion20 = "2.0"

	// BlkVersion21 tx version 2.1
	BlkVersion21 = "2.1"

	// BlkVersion22 tx version 2.2
	BlkVersion22 = "2.2"

	// BlkVersion23 tx version 2.3
	BlkVersion23 = "2.3"

	// BlkVersion24 tx version 2.4
	BlkVersion24 = "2.4"

	// BlkVersion25 tx version 2.5
	BlkVersion25 = "2.5"

	// BlkVersion26 tx version 2.6
	BlkVersion26 = "2.6"

	// BlkVersion27 tx version 2.7
	BlkVersion27 = "2.7"

	// BlkVersion28 tx version 2.8
	BlkVersion28 = "2.8"

	// BlkVersion29 tx version 2.9
	BlkVersion29 = "2.9"

	// BlkVersion30 tx version 3.0
	BlkVersion30 = "3.0"

	// BlkVersion31 tx version 3.1
	BlkVersion31 = "3.1"

	// BlkVersion32 add invalidRoot into block struct
	BlkVersion32 = "3.2"

	// BlkVersion33 tx version 3.3
	BlkVersion33 = "3.3"

	// BlkVersion34 tx version 3.4
	BlkVersion34 = "3.4"

	// BlkVersion35 tx version 3.5
	BlkVersion35 = "3.5"

	// BlkVersion36 tx version 3.6
	BlkVersion36 = "3.6"

	// BlkVersion40 tx version 4.0
	BlkVersion40 = "4.0"
)

// receipt version in string
const (
	// ReceiptVersion25 receipt version 2.5
	ReceiptVersion25 = "2.5"
)

const (
	// EncodeVersion encode version
	// Deprecated: please read encode version from ledger
	EncodeVersion = EncodeVersion00
)

const (
	//EncodeVersion00 encode initial version
	EncodeVersion00 byte = 0x00
)

// encode version in string
var (
	encodeVersionMap = map[byte]string{
		EncodeVersion00: "0",
	}
)
