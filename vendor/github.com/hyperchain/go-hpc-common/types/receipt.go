// Package types Copyright 2016-2017 Flato Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package types

import (
	"github.com/hyperchain/go-hpc-common/types/protos"
	"github.com/hyperchain/go-hpc-common/utils"

	"github.com/willf/bloom"
)

// ReceiptTrans are used to show in web.
type ReceiptTrans struct {
	Version           string                `json:"version"`
	Bloom             string                `json:"bloom"`
	CumulativeGasUsed int64                 `json:"cumulativeGasUsed"`
	TxHash            string                `json:"txHash"`
	ContractAddress   string                `json:"contractAddress"`
	GasUsed           int64                 `json:"gasUsed"`
	Ret               string                `json:"ret"`
	Status            protos.Receipt_STATUS `json:"status"`
	Message           string                `json:"message"`
	Logs              []LogTrans            `json:"logs"`
	VMType            string                `json:"vmType"`
	Oracles           []OracleTrans         `json:"oracles"`
}

// ToReceiptTrans receipt to receiptTrans
func ToReceiptTrans(r protos.Receipt) (receiptTrans *ReceiptTrans) {
	logs, err := RetrieveLogs(&r)
	var logsValue []LogTrans
	if err != nil {
		logsValue = nil
	} else {
		logsValue = logs.ToLogsTrans()
	}
	oracles, err := RetrieveOracles(&r)
	var oraclesValue []OracleTrans
	if err != nil {
		oraclesValue = nil
	} else {
		oraclesValue = oracles.ToOracleTrans()
	}
	return &ReceiptTrans{
		Version:           string(r.Version),
		GasUsed:           r.GasUsed,
		Bloom:             "",
		ContractAddress:   BytesToAddress(r.ContractAddress).Hex(),
		CumulativeGasUsed: r.CumulativeGasUsed,
		Ret:               utils.ToHex(r.Ret),
		TxHash:            BytesToHash(r.TxHash).String(),
		Status:            r.Status,
		Message:           string(r.Message),
		Logs:              logsValue,
		VMType:            r.VmType.String(),
		Oracles:           oraclesValue,
	}
}

// NewReceipt creates a barebone transaction receipt, copying the init fields.
func NewReceipt(receiptVersion []byte, cumulativeGasUsed uint64, vmType int32) *protos.Receipt {
	return &protos.Receipt{Version: receiptVersion, CumulativeGasUsed: int64(cumulativeGasUsed), VmType: protos.Receipt_VmType(vmType)}
}

// RetrieveLogs retrieve logs
func RetrieveLogs(r *protos.Receipt) (Logs, error) {
	return DecodeLogs((*r).Logs)
}

// SetLogs set logs
func SetLogs(r *protos.Receipt, logs Logs) error {
	buf, err := (&logs).EncodeLogs()
	if err != nil {
		return err
	}
	r.Logs = buf
	return nil
}

//RetrieveOracles retrieve oracles
func RetrieveOracles(r *protos.Receipt) (OracleEvents, error) {
	return DecodeOracles((*r).Oracles)
}

// CreateBloom create bloom
func CreateBloom(receipts []*protos.Receipt) ([]byte, error) {
	blom := bloom.New(256*8, 3)

	for _, r := range receipts {
		logs, err := RetrieveLogs(r)
		if err != nil {
			return nil, err
		}
		for _, log := range logs {
			blom.Add(log.Address.Bytes())
			for _, topic := range log.Topics {
				blom.Add(topic.Bytes())
			}
		}
	}
	return blom.GobEncode()
}

// BloomFilter return bloom filter
func BloomFilter(r *protos.Receipt) (*bloom.BloomFilter, error) {
	blom := bloom.New(256*8, 3)
	if err := blom.GobDecode(r.GetBloom()); err != nil {
		return nil, err
	}
	return blom, nil
}

// BloomLookup bloom lookup
func BloomLookup(bloom *bloom.BloomFilter, content []byte) bool {
	return bloom.Test(content)
}

// Receipts is a wrapper around a Receipt array to implement types.DerivableList.
type Receipts []*protos.Receipt

// ReceiptForStorage is Receipt for storage
type ReceiptForStorage protos.Receipt

// Len returns the number of receipts in this list.
func (r Receipts) Len() int { return len(r) }
