package types

// Metadata defines the parameters related to pagination in the
// API query request.
type Metadata struct {
	// the number of returned value in a query request
	PageSize int32 `json:"pagesize"`

	// Bookmark defines the starting position of a query request
	Bookmark *Bookmark `json:"bookmark"`

	// Backward indicates the direction of a query request.
	// true means to search backward from the bookmark position,
	// otherwise to search forward from the bookmark position
	Backward bool `json:"backward"`
}

// Bookmark defines the starting position of a query request,
// its unit is the position of specific transaction, so it
// must consist of block number and transaction index.
type Bookmark struct {
	BlockNumber uint64 `json:"blkNum"`
	TxIndex     int64  `json:"txIndex"`
}

// TransactionsPaginationResult defines the pagination result structure.
type TransactionsPaginationResult struct {
	Hasmore  bool                 `json:"hasmore"`
	Metadata *Metadata            `json:"metadata"`
	Txs      []*TransactionResult `json:"txs"`
}

// TransactionResult defines the returned data structure of querying
// transaction sub-document in indexdb.
type TransactionResult struct {
	TxHash      string      `bson:"txHash,omitempty" json:"txHash,omitempty"`
	BlockNumber uint64      `bson:"blkNum" json:"blkNum"`
	TxIndex     int64       `bson:"txIndex" json:"txIndex"`
	From        string      `bson:"from,omitempty" json:"from,omitempty"`
	To          string      `bson:"to,omitempty" json:"to,omitempty"`
	ExtraID     interface{} `bson:"extraId,omitempty" json:"extraId,omitempty"`
}

// BlockStatResult defines the blocks statistics result.
type BlockStatResult struct {
	FirstNumber uint64 `bson:"first" json:"first"`
	LastNumber  uint64 `bson:"last" json:"last"`
	Count       int    `bson:"count" json:"count"`
	TxCount     int    `bson:"txCount" json:"txCount"`
}
