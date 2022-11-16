package types

// TxContent represents tx content and execution result, used for tx list in block event
type TxContent struct {
	Extra           string     `json:"extra"`
	ExtraID         string     `json:"extraID"`
	BlockNumber     string     `json:"blockNumber,omitempty"` // block number where this transaction was in
	BlockHash       *Hash      `json:"blockHash,omitempty"`   // hash of the block where this transaction was in
	TxIndex         int64      `json:"txIndex,omitempty"`     // transaction index in the block
	From            string     `json:"from"`                  // the address of sender
	To              string     `json:"to"`                    // the address of receiver
	Amount          string     `json:"amount,omitempty"`      // transfer amount
	Timestamp       string     `json:"timestamp"`             // the unix timestamp for when the transaction was generated
	Nonce           string     `json:"nonce"`
	Payload         string     `json:"payload,omitempty"`
	Signature       string     `json:"signature,omitempty"`
	Version         string     `json:"version"`
	TxHash          string     `json:"txHash"`
	VMType          string     `json:"vmType"`
	ContractAddress string     `json:"contractAddress"`
	GasUsed         int64      `json:"gasUsed"`
	Ret             string     `json:"ret"`
	Log             []LogTrans `json:"log"`
}
