package types

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperchain/go-hpc-common/cryptocom"

	"github.com/meshplus/crypto"

	"github.com/hyperchain/go-hpc-common/types/protos"
	"github.com/hyperchain/go-hpc-common/utils"

	"github.com/pingcap/failpoint"
)

// const length
const (
	// NumberLength is the number length of block hash
	NumberLength = 8
)

// CalBlockHash return block hash
func CalBlockHash(block *protos.Block, engine cryptocom.Engine) Hash {
	da, _ := engine.GetDefaultAlgo()
	hashAlgo := da
	if len(block.BlockHash) != 0 {
		hashAlgo, _ = MultiDecode(block.BlockHash[NumberLength:], HashMethod)
	}

	hasher, err := engine.GetHash(hashAlgo)
	if err != nil {
		return Hash{}
	}
	var res []byte
	var jerr error

	v := block.GetVersion()
	version, gerr := GetBlkVersion(string(v))
	if gerr != nil || v == nil || len(v) == 0 {
		// this can only happen in sync chain when the provider
		// provided an invalid block, and empty result can not
		// pass the checker's judge function in sync chain.
		return Hash{}
	}
	switch {
	case version.Compare(BlkVersion40) >= 0:
		var qcBytes []byte
		var qcAuthors []string
		if block.QuorumCert != nil {
			var merr error
			qcBytes, merr = dumpQCBytes(block.QuorumCert)
			if merr != nil {
				panic(jerr)
			}
			qcAuthors = block.QuorumCert.SortedSignatures().Authors()
		}
		res, jerr = json.Marshal([]interface{}{
			block.ParentHash,
			block.Number,
			block.Timestamp,
			block.TxRoot,
			block.ReceiptRoot,
			block.MerkleRoot,
			block.InvalidRoot,
			block.Author,
			block.Signature,
			qcBytes,
			qcAuthors,
			block.Version,
		})

	case version.Compare(BlkVersion36) >= 0:
		var qcBytes []byte
		if block.QuorumCert != nil {
			var merr error
			qcBytes, merr = dumpQCBytes(block.QuorumCert)
			if merr != nil {
				panic(jerr)
			}
		}
		res, jerr = json.Marshal([]interface{}{
			block.ParentHash,
			block.Number,
			block.Timestamp,
			block.TxRoot,
			block.ReceiptRoot,
			block.MerkleRoot,
			block.InvalidRoot,
			qcBytes,
			block.Author,
			block.Signature,
		})
	case version.Compare(BlkVersion31) > 0:
		res, jerr = json.Marshal([]interface{}{
			block.ParentHash,
			block.Number,
			block.Timestamp,
			block.TxRoot,
			block.ReceiptRoot,
			block.MerkleRoot,
			block.InvalidRoot,
		})
	default:
		res, jerr = json.Marshal([]interface{}{
			block.ParentHash,
			block.Number,
			block.Timestamp,
			block.TxRoot,
			block.ReceiptRoot,
			block.MerkleRoot,
		})
	}

	failpoint.Inject("fp-CalBlockHash-1", func(_ failpoint.Value) {
		jerr = errors.New("expected error")
	})
	if jerr != nil {
		// copy the history logic
		panic(jerr)
	}
	h, herr := hasher.Hash(res)

	failpoint.Inject("fp-CalBlockHash-2", func(_ failpoint.Value) {
		herr = errors.New("expected error")
	})
	if herr != nil {
		return Hash{}
	}
	hash := BytesToHash(h)

	ret := MultiEncode(hashAlgo, hash)

	switch {
	case version.Compare(BlkVersion24) > 0:
		if hashAlgo == crypto.KECCAK_256 {
			binary.BigEndian.PutUint64(hash[0:NumberLength], block.GetNumber())
			return hash
		}
		blockNumber := make([]byte, NumberLength)
		binary.BigEndian.PutUint64(blockNumber, block.GetNumber())
		ret = append(blockNumber, ret...)
		return ret
	default:
		return hash
	}
}

func dumpQCBytes(qc *protos.QuorumCert) ([]byte, error) {
	return proto.Marshal(&protos.QuorumCert{
		VoteData: qc.VoteData,
		SignedLedgerInfo: &protos.LedgerInfoWithSignatures{
			LedgerInfo: qc.SignedLedgerInfo.LedgerInfo,
		},
	})
}

// ParseNumberFromBlockHash returns block number by parsing block hash
func ParseNumberFromBlockHash(h Hash) uint64 {
	return utils.ReadVarInt(h[0:NumberLength])
}

// BlockJSON is block can Marshal type for json
type BlockJSON struct {
	Version      string             ` json:"version,omitempty"`
	ParentHash   string             `json:"parentHash,omitempty"`
	BlockHash    string             `json:"blockHash,omitempty"`
	Transactions []*TransactionJSON `json:"transactions,omitempty"`
	Timestamp    int64              `json:"timestamp,omitempty"`
	MerkleRoot   string             `json:"merkleRoot,omitempty"`
	TxRoot       string             `json:"txRoot,omitempty"`
	ReceiptRoot  string             `json:"receiptRoot,omitempty"`
	Number       uint64             `json:"number,omitempty"`
	WriteTime    int64              `json:"writeTime,omitempty"`
	CommitTime   int64              `json:"commitTime,omitempty"`
	EvmTime      int64              `json:"evmTime,omitempty"`
	Bloom        string             `json:"bloom,omitempty"`
}

// CovertBlockJSON cover a block message to can Marshal type
func CovertBlockJSON(t *protos.Block) *BlockJSON {
	txs := make([]*TransactionJSON, len(t.Transactions))
	for i, tx := range t.Transactions {
		txs[i] = CovertTransactionJSON(tx)
	}
	return &BlockJSON{
		Version:      string(t.Version),
		ParentHash:   utils.ToHex(t.ParentHash),
		BlockHash:    utils.ToHex(t.BlockHash),
		Transactions: txs,
		Timestamp:    t.Timestamp,
		MerkleRoot:   utils.ToHex(t.MerkleRoot),
		TxRoot:       utils.ToHex(t.TxRoot),
		ReceiptRoot:  utils.ToHex(t.ReceiptRoot),
		Number:       t.Number,
		WriteTime:    t.WriteTime,
		CommitTime:   t.CommitTime,
		EvmTime:      t.EvmTime,
		Bloom:        utils.ToHex(t.Bloom),
	}
}
