// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"fmt"

	"github.com/hyperchain/go-hpc-common/utils"

	json "github.com/json-iterator/go"
)

// Log is the event log struct of Flato
type Log struct {
	// consensus fields
	Address Address `json:"Address"`
	Topics  []Hash  `json:"Topics"`
	Data    []byte  `json:"Data"`

	// Derived fields (don't reorder!)
	BlockNumber uint64 `json:"BlockNumber"`
	TxHash      Hash   `json:"TxHash"`
	TxIndex     uint   `json:"TxIndex"`
	BlockHash   Hash   `json:"BlockHash"`
	Index       uint   `json:"Index"`
}

// NewLog assign block number as 0 temporarily
// because the blcok number in env is a seqNo actually
// primary's seqNo may not equal to other's
// correctly block number and block hash will be assigned in the commit phase
func NewLog(address Address, topics []Hash, data []byte, number uint64) *Log {
	return &Log{Address: address, Topics: topics, Data: data, BlockNumber: 0}
}

func (l *Log) String() string {
	return fmt.Sprintf(`log: %x %x %x %x %d %x %d`, l.Address, l.Topics, l.Data, l.TxHash, l.TxIndex, l.BlockHash, l.Index)
}

// EncodeLog encode Log to bytes
func (l *Log) EncodeLog() ([]byte, error) {
	return json.Marshal(*l)
}

// DecodeLog decode Log from bytes
func DecodeLog(buf []byte) (Log, error) {
	var tmp Log
	err := json.Unmarshal(buf, &tmp)
	return tmp, err
}

// Logs logs
type Logs []*Log

// EncodeLogs encode Logs to bytes
func (ls *Logs) EncodeLogs() ([]byte, error) {
	return json.Marshal(*ls)
}

// DecodeLogs decode Logs from bytes
func DecodeLogs(buf []byte) (Logs, error) {
	var tmp Logs
	err := json.Unmarshal(buf, &tmp)
	return tmp, err
}

// LogTrans is the real struct store in db
type LogTrans struct {
	Address     string   `json:"address"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
	BlockNumber uint64   `json:"blockNumber"`
	BlockHash   string   `json:"blockHash"`
	TxHash      string   `json:"txHash"`
	TxIndex     uint     `json:"txIndex"`
	Index       uint     `json:"index"`
}

// ToLogsTrans construct the real struct
func (ls Logs) ToLogsTrans() []LogTrans {
	var ret = make([]LogTrans, len(ls))
	for idx, log := range ls {
		var topics = make([]string, len(log.Topics))
		for ti, t := range log.Topics {
			topics[ti] = t.String()
		}
		ret[idx] = LogTrans{
			Address:     log.Address.Hex(),
			Data:        utils.BytesToHex(log.Data),
			BlockNumber: log.BlockNumber,
			BlockHash:   log.BlockHash.String(),
			Topics:      topics,
			TxHash:      log.TxHash.String(),
			Index:       log.Index,
			TxIndex:     log.TxIndex,
		}
	}
	return ret
}
