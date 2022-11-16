package types

import "fmt"

/**
@collection: blocks
@primary key:  BlockNumber
@index: writeTime, txs.hash, txs.from, txs.to
*/

// BlockDoc is the data structure of a document in indexdb.
type BlockDoc struct {
	BlockNumber    uint64            `bson:"_id" json:"_id"`
	BlockWriteTime int64             `bson:"writeTime,omitempty" json:"writeTime,omitempty"`
	Txs            []*TransactionDoc `bson:"txs,omitempty" json:"txs,omitempty"`
}

// TransactionDoc is the data structure of a sub-document in indexdb.
type TransactionDoc struct {
	TxHash  string      `bson:"hash,omitempty" json:"hash,omitempty"`
	TxIndex int64       `bson:"index" json:"index"`
	From    string      `bson:"from,omitempty" json:"from,omitempty"`
	To      string      `bson:"to,omitempty" json:"to,omitempty"`
	ExtraID interface{} `bson:"eid,omitempty" json:"extraId,omitempty"`
}

// QueryTxArgs is the filter condition used in indexdb.
type QueryTxArgs struct {
	TxHash      *string     `bson:"txs.hash,omitempty" json:"txHash,omitempty"`
	BlockNumber *uint64     `bson:"_id,omitempty" json:"blkNumber,omitempty"`
	TxIndex     *int64      `bson:"txs.index,omitempty" json:"txIndex,omitempty"`
	From        *string     `bson:"txs.from,omitempty" json:"from,omitempty"`
	To          *string     `bson:"txs.to,omitempty" json:"to,omitempty"`
	ExtraID     interface{} `bson:"txs.eid,omitempty" json:"extraId,omitempty"`
}

func (qt *QueryTxArgs) String() string {
	msg := "&QueryTxArgs{"
	if qt.TxHash != nil {
		msg += "TxHash: " + *qt.TxHash
	}
	if qt.BlockNumber != nil {
		msg += fmt.Sprintf(" BlockNumber: %v", *qt.BlockNumber)
	}
	if qt.TxIndex != nil {
		msg += fmt.Sprintf(" TxIndex: %v", *qt.TxIndex)
	}
	if qt.From != nil {
		msg += " From: " + *qt.From
	}
	if qt.To != nil {
		msg += " To: " + *qt.To
	}
	if qt.ExtraID != nil {
		msg += fmt.Sprintf(" ExtraID: %v", qt.ExtraID)
	}
	msg += "}"
	return msg
}

// FieldID defines different field identification on document.
type FieldID int

// defines different field identification on document
const (
	FieldBlkNum FieldID = iota
	FieldWriteTime
	FieldTxFrom
	FieldTxTo
	FieldTxHash
	FieldExtraID
)

// FieldName map field id to corresponding field name.
var FieldName = map[FieldID]string{
	FieldBlkNum:    "_id",
	FieldWriteTime: "writeTime",
	FieldTxFrom:    "txs.from",
	FieldTxTo:      "txs.to",
	FieldTxHash:    "txs.hash",
	FieldExtraID:   "txs.eid",
}

// FieldNameToID map field name to corresponding field id.
var FieldNameToID = map[string]FieldID{
	"_id":       FieldBlkNum,
	"writeTime": FieldWriteTime,
	"txs.from":  FieldTxFrom,
	"txs.to":    FieldTxTo,
	"txs.hash":  FieldTxHash,
	"txs.eid":   FieldExtraID,
}
