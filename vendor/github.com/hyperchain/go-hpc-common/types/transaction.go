// Copyright 2016-2017 Flato Corp.
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
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	cryptohash "hash"
	"math/big"
	"strconv"
	"strings"

	"github.com/hyperchain/go-hpc-common/types/protos"
	"github.com/hyperchain/go-hpc-common/utils"

	"github.com/gogo/protobuf/proto"
	"github.com/meshplus/crypto"
	"github.com/pingcap/failpoint"
	"golang.org/x/crypto/sha3"
)

// const length
const (
	hashSplit = "@"
	// NumberLength is the number length of block hash
	TimeLength = 8
)

const eccKeyLength, ed25519KeyLength, flagLength = 65, 32, 1

// Hash returns the transaction hash calculated using Keccak256.
func hash(tx *protos.Transaction) Hash {
	var h32 [32]byte
	if CompareVersion(string(tx.Version), TxVersion35) > 0 {
		binary.BigEndian.PutUint64(h32[0:TimeLength], uint64(tx.GetTimestamp()))
		copy(h32[TimeLength:], tx.Signature[len(tx.Signature)-24:])
		return h32[:]
	}
	res, jerr := json.Marshal([]interface{}{
		tx.From,
		tx.To,
		tx.Value,
		tx.Timestamp,
		tx.Nonce,
		tx.Signature,
		tx.ExpirationTimestamp,
		tx.Participant,
	})

	failpoint.Inject("fp-Hash-1", func() {
		jerr = errors.New("Hash : excepted error")
	})

	if jerr != nil {
		// copy the history logic
		panic(jerr)
	}
	hasher := sha3.NewLegacyKeccak256()
	_, herr := hasher.Write(res)
	failpoint.Inject("fp-Hash-2", func() {
		herr = errors.New("Hash : excepted error")
	})
	if herr != nil {
		return Hash{}
	}
	h := hasher.Sum(nil)

	hash := BytesToHash(h)
	binary.BigEndian.PutUint64(hash[0:TimeLength], uint64(tx.GetTimestamp()))
	return hash
}

// GetTimestampFromHash returns tx's timestamp by parsing tx hash
func GetTimestampFromHash(h Hash) int64 {
	return int64(utils.ReadVarInt(h[0:NumberLength]))
}

// GetHash returns the TransactionHash directly if hash has been recorded in tx,
// else calculate and return.
func GetHash(tx *protos.Transaction) Hash {
	if len(tx.TransactionHash) == 0 {
		return hash(tx)
	}
	return BytesToHash(tx.TransactionHash)
}

// SignHash returns the whole hash of this transaction calculated using Keccak256.
func SignHash(tx *protos.Transaction, ch crypto.Hasher) Hash {
	value := new(protos.TransactionValue)
	hashErr := proto.Unmarshal(tx.Value, value)
	if hashErr != nil {
		hh, _ := ch.Hash([]byte("invalid hash"))
		return BytesToHash(hh)
	}
	needHash := SerialString(tx, value)
	ret, _ := ch.Hash([]byte(needHash))
	return BytesToHash(ret)
}

// SerialString serial tx and value to need hash string
// Deprecated
func SerialString(tx *protos.Transaction, value *protos.TransactionValue) string {
	version, err := GetTXVersion(string(tx.Version))
	// version check
	// (1) version not exist, Note! if here has a invalid number format, it will also be parse as 0
	// (2) version > 1.3 generally (1.6)
	// (3) version parse failed, err != nil
	// Note: for very old v1.2 version, the version field was not exist,
	//       so, we may need below logic to avoid version is empty bug
	//       but, actually, most customers are using latest v1.2, so we ignore the old
	//       version firstly.
	//if (version == 0 && err.Error() == "strconv.ParseFloat: parsing \"\": invalid syntax") || (version > 1.3 && err == nil) {
	failpoint.Inject("SerialString-1", func(_ failpoint.Value) {
		if value.Payload == nil {
			value.Payload = []byte("1")
			defer func() {
				value.Payload = nil
			}()
		}
	})
	// if version <= 1.3 hashString = from + to + value + timestamp + nonce
	// if 1.3 < version <= 1.7 hashString = from + to + value + timestamp + nonce + opcode + extra + vmtype
	// if 1.7 < version <= 2.0 hashString = from + to + value + timestamp + nonce + opcode + extra + vmtype + version
	// if 2.0 < version <= 2.1 hashString = from + to + value + timestamp + nonce + opcode + extra + vmtype + version + extraid
	// if 2.1 < version <= 2.2 hashString = from + to + value + timestamp + nonce + opcode + extra + vmtype + version + extraid + cname
	// if 2.2 < version hashString = from + to + value + timestamp + nonce + opcode + extra + vmtype + version + extraid + cname + price + gasLimit + expirationTimestamp + participant

	needHashString := "from=" + utils.ToHex(tx.From) +
		"&to=" + utils.ToHex(tx.To) +
		"&value=0x" + strconv.FormatInt(value.Amount, 16) +
		"&payload=" + utils.ToHex(value.Payload) +
		"&timestamp=0x" + strconv.FormatInt(tx.Timestamp, 16) +
		"&nonce=0x" + strconv.FormatInt(tx.Nonce, 16) +
		"&opcode=" + strconv.FormatInt(int64(value.Op), 16) +
		"&extra=" + string(value.Extra) +
		"&vmtype=" + value.VmType.String() +
		"&version=" + string(tx.Version)

	if err == nil && version.Compare(TxVersion21) >= 0 {
		needHashString += "&extraid=" + string(value.ExtraId)
	}

	if err == nil && version.Compare(TxVersion22) >= 0 {
		needHashString += "&cname=" + string(tx.CName)
	}
	if err == nil && version.Compare(TxVersion36) >= 0 {
		needHashString += "&price=0x" + strconv.FormatInt(value.Price, 16)
		needHashString += "&gasLimit=0x" + strconv.FormatInt(value.GasLimit, 16)
		needHashString += "&expirationTimestamp=0x" + strconv.FormatInt(tx.ExpirationTimestamp, 16)
		if tx.Participant == nil {
			needHashString += "&initiator=" + "&withholding="
		} else {
			needHashString += "&initiator=" + utils.ToHex(tx.Participant.Initiator)
			str := "["
			for i, bs := range tx.Participant.Withholding {
				str += utils.ToHex(bs)
				if i < len(tx.Participant.Withholding)-1 {
					str += ","
				}
			}
			str += "]"
			needHashString += "&withholding=" + str
		}
	}
	return needHashString
}

// SerialBytes serial tx and value to need hash string
//nolint: errcheck
func SerialBytes(tx *protos.Transaction, v *protos.TransactionValue, h cryptohash.Hash) {
	const (
		from                = "from=0x"
		to                  = "&to=0x"
		value               = "&value=0x"
		payload             = "&payload=0x"
		timestamp           = "&timestamp=0x"
		nonce               = "&nonce=0x"
		op                  = "&opcode="
		extra               = "&extra="
		vmtype              = "&vmtype="
		version             = "&version="
		extraID             = "&extraid="
		caname              = "&cname="
		price               = "&price=0x"
		gasLimit            = "&gasLimit=0x"
		expirationTimestamp = "&expirationTimestamp=0x"
		initiator           = "&initiator="
		withholding         = "&withholding="
	)
	txVersion, err := GetTXVersion(string(tx.Version))
	h.Write([]byte(from))
	hexTo(tx.From, h)
	h.Write([]byte(to))
	hexTo(tx.To, h)
	h.Write([]byte(value))
	h.Write([]byte(strconv.FormatInt(v.Amount, 16)))
	h.Write([]byte(payload))
	hexTo(v.Payload, h)
	h.Write([]byte(timestamp))
	h.Write([]byte(strconv.FormatInt(tx.Timestamp, 16)))
	h.Write([]byte(nonce))
	h.Write([]byte(strconv.FormatInt(tx.Nonce, 16)))
	h.Write([]byte(op))
	h.Write([]byte(strconv.FormatInt(int64(v.Op), 16)))
	h.Write([]byte(extra))
	h.Write(v.Extra)
	h.Write([]byte(vmtype))
	h.Write([]byte(v.VmType.String()))
	h.Write([]byte(version))
	h.Write(tx.Version)

	if err == nil && txVersion.Compare(TxVersion21) >= 0 {
		h.Write([]byte(extraID))
		h.Write(v.ExtraId)
	}

	if err == nil && txVersion.Compare(TxVersion22) >= 0 {
		h.Write([]byte(caname))
		h.Write(tx.CName)
	}

	if err == nil && txVersion.Compare(TxVersion36) >= 0 {
		h.Write([]byte(price))
		h.Write([]byte(strconv.FormatInt(v.Price, 16)))
		h.Write([]byte(gasLimit))
		h.Write([]byte(strconv.FormatInt(v.GasLimit, 16)))
		h.Write([]byte(expirationTimestamp))
		h.Write([]byte(strconv.FormatInt(tx.ExpirationTimestamp, 16)))
		if tx.Participant == nil {
			h.Write([]byte(initiator))
			h.Write([]byte(withholding))
		} else {
			h.Write([]byte(initiator))
			h.Write([]byte("0x"))
			hexTo(tx.Participant.Initiator, h)
			h.Write([]byte(withholding))
			h.Write([]byte("["))
			for i, bs := range tx.Participant.Withholding {
				h.Write([]byte("0x"))
				hexTo(bs, h)
				if i < len(tx.Participant.Withholding)-1 {
					h.Write([]byte(","))
				}
			}
			h.Write([]byte("]"))
		}
	}
}

func hexTo(src []byte, hr cryptohash.Hash) {
	var buf [64]byte
	if len(src) == 0 {
		_, _ = hr.Write([]byte{'0'})
		return
	}
	cycNum := len(src) >> 5
	s := 0
	for i := 0; i < cycNum; i++ {
		hex.Encode(buf[:], src[s:s+32])
		_, _ = hr.Write(buf[:])
		s += 32
	}

	if s == len(src) {
		return
	}

	hex.Encode(buf[:(len(src)-s)*2], src[s:])
	_, _ = hr.Write(buf[:(len(src)-s)*2])
}

// NewTransaction returns a new transaction.
func NewTransaction(txVersion []byte, from []byte, to []byte, value []byte, timestamp, nonce, expirationTimestamp int64, participant *protos.Participant) *protos.Transaction {
	transaction := &protos.Transaction{
		Version:             txVersion,
		From:                from,
		To:                  to,
		Value:               value,
		Timestamp:           timestamp,
		Nonce:               nonce,
		ExpirationTimestamp: expirationTimestamp,
		Participant:         participant,
		Other:               &protos.NonHash{},
	}

	return transaction
}

// NewTransactionValue returns a new TransactionValue.
func NewTransactionValue(price, gasLimit, amount int64, payload []byte, opcode int32, extra, extraID []byte, vmType protos.TransactionValue_VmType) *protos.TransactionValue {
	return &protos.TransactionValue{
		Price:    price,
		GasLimit: gasLimit,
		Amount:   amount,
		Payload:  payload,
		Extra:    extra,
		ExtraId:  extraID,
		Op:       protos.TransactionValue_Opcode(opcode),
		VmType:   vmType,
	}
}

// GetTransactionValue returns the TransactionValue.
func GetTransactionValue(tx *protos.Transaction) *protos.TransactionValue {
	transactionValue := &protos.TransactionValue{}
	_ = proto.Unmarshal(tx.Value, transactionValue)
	return transactionValue
}

// RetrievePayload returns the tx payload.
func RetrievePayload(tv *protos.TransactionValue) []byte {
	return utils.CopyBytes(tv.Payload)
}

// RetrieveExtra returns the tx payload.
func RetrieveExtra(tv *protos.TransactionValue) []byte {
	return utils.CopyBytes(tv.Extra)
}

// RetrieveGas returns the tx gas limit, note the type is uint64
func RetrieveGas(tv *protos.TransactionValue) uint64 {
	return uint64(tv.GasLimit)
}

// RetrieveGasPrice returns the tx gas price.
func RetrieveGasPrice(tv *protos.TransactionValue) *big.Int {
	return new(big.Int).Set(big.NewInt(tv.Price))
}

// RetrieveAmount returns the tx amount.
func RetrieveAmount(tv *protos.TransactionValue) *big.Int {
	return new(big.Int).Set(big.NewInt(tv.Amount))
}

// AppendNodeHash is used to append the given node hash to value of `tx.Other.NodeHash`
// with special characters. If `tx.Other.NodeHash` is nil, assigned the given node
// hash to it directly.
func AppendNodeHash(tx *protos.Transaction, hash string) error {
	if strings.HasPrefix(hash, "0x") {
		hash = hash[2:]
	}
	if len(tx.GetOther().GetNodeHash()) != 0 {
		tx.GetOther().NodeHash = append(tx.GetOther().NodeHash, []byte(hashSplit+hash)...)
	} else {
		tx.GetOther().NodeHash = []byte(hash)
	}
	return nil
}

// GetNodeHash returns the hash value of all nodes that have forwarded this transaction.
func GetNodeHash(tx *protos.Transaction) []string {
	nodesHash := string(tx.GetOther().GetNodeHash())
	return strings.Split(nodesHash, hashSplit)
}

// GetPrivateTxHash returns the private tx hash if backward compatible.
func GetPrivateTxHash(tx *protos.Transaction) []byte {
	return tx.GetOther().PrivateTxHash
}

// SetPrivateTxHash sets the private tx hash if backward compatible.
func SetPrivateTxHash(tx *protos.Transaction, hash []byte) {
	tx.GetOther().PrivateTxHash = hash
}

// IsPrivate returns if this tx is corresponding with a private tx.
func IsPrivate(tx *protos.Transaction) bool {
	return tx.GetOther().PrivateTxHash != nil
}

// GetPrivateCollection returns the private tx collection if backward compatible.
func GetPrivateCollection(tx *protos.Transaction) []string {
	return tx.GetOther().Collection
}

// SetPrivateCollection sets the private tx collection if backward compatible.
func SetPrivateCollection(tx *protos.Transaction, collection []string) {
	tx.GetOther().Collection = collection
}

// GetPrivateNonce returns the private tx nonce if backward compatible.
func GetPrivateNonce(tx *protos.Transaction) uint64 {
	return tx.GetOther().Nonce
}

// SetPrivateNonce sets the private tx nonce if backward compatible.
func SetPrivateNonce(tx *protos.Transaction, nonce uint64) {
	tx.GetOther().Nonce = nonce
}

// IsConfigTx returns if this tx is corresponding with a config tx.
func IsConfigTx(tx *protos.Transaction) bool {
	return tx.IsConfigTx()
}

// PrivateMeta records the private transaction hash together with
// its collection info
type PrivateMeta struct {
	BlockNum   uint64
	TxHash     string
	Collection []string
	Nonce      uint64
}

//IsInValidOp judge an opcode is valid
func IsInValidOp(opcode protos.TransactionValue_Opcode) bool {
	if _, ok := protos.TransactionValue_Opcode_name[int32(opcode)]; ok {
		return false
	}
	return true
}

//IsManagementOp judge an opcode is manage opcode
func IsManagementOp(opcode protos.TransactionValue_Opcode) bool {
	return opcode == protos.TransactionValue_FREEZE ||
		opcode == protos.TransactionValue_UNFREEZE ||
		opcode == protos.TransactionValue_DESTROY
}

//TransactionJSON is transaction JSON friendly form
type TransactionJSON struct {
	Version         string       `json:"version,omitempty"`
	From            string       `json:"from,omitempty"`
	To              string       `json:"to,omitempty"`
	Value           string       `json:"value,omitempty"`
	Timestamp       int64        `json:"timestamp,omitempty"`
	Signature       string       `json:"signature,omitempty"`
	ID              uint64       `json:"id,omitempty"`
	TransactionHash string       `json:"transactionHash,omitempty"`
	Nonce           int64        `json:"nonce,omitempty"`
	Other           *NonHashJSON `json:"other,omitempty"`
	TxType          int32        `json:"txType,omitempty"`
}

//NonHashJSON is NonHash JSON friendly form
type NonHashJSON struct {
	NodeHash      string   `json:"nodeHash,omitempty"`
	PrivateTxHash string   `json:"privateTxHash,omitempty"`
	Collection    []string `json:"collection,omitempty"`
	Nonce         uint64   `son:"nonce,omitempty"`
}

//CovertTransactionJSON covert a transaction struct to a transactionJSON struct
func CovertTransactionJSON(t *protos.Transaction) *TransactionJSON {
	tj := &TransactionJSON{}
	tj.Timestamp = t.Timestamp
	tj.Signature = utils.BytesToHex(t.Signature)
	tj.Value = utils.BytesToHex(t.Value)
	tj.Version = string(t.Version)
	tj.From = utils.BytesToHex(t.From)
	tj.To = utils.BytesToHex(t.To)
	tj.ID = t.Id
	tj.Nonce = t.Nonce
	tj.TransactionHash = utils.BytesToHex(t.TransactionHash)
	tj.TxType = int32(t.TxType)
	nonHashJSON := covertNonHashJSON(t.Other)
	tj.Other = nonHashJSON
	return tj
}

func covertNonHashJSON(nonHash *protos.NonHash) *NonHashJSON {
	result := &NonHashJSON{}
	result.Nonce = nonHash.Nonce
	result.Collection = nonHash.Collection
	result.NodeHash = utils.BytesToHex(nonHash.NodeHash)
	result.PrivateTxHash = utils.BytesToHex(nonHash.PrivateTxHash)
	return result
}

//ParseTxSignature parse tx signature
// 0 ecdsa secp256k1 recover
// 1 :sm2 | 2 :ed25519 | 4 :pki | 5: ecdsa r1 | 6: ecdsa(pubKey)
func ParseTxSignature(signature, fromAddress []byte) (algo uint8, sign []byte, key []byte) {
	if len(signature) == 0 || len(fromAddress) == 0 {
		return 255, nil, nil
	}
	algo = signature[0]
	switch algo & 0x7F {
	case 0:
		sign, key = signature[flagLength:], fromAddress[:]
	case 1, 5, 6, 7, 8:
		sign, key = signature[eccKeyLength+flagLength:], signature[1:eccKeyLength+flagLength]
	case 2:
		sign, key = signature[ed25519KeyLength+flagLength:], signature[1:ed25519KeyLength+flagLength]
	case 4:
		sign, key = signature[flagLength:], nil
	default:
		algo = 255
	}
	return
}

// TransactionWrapper transaction wrapper
type TransactionWrapper struct {
	Tx      *protos.Transaction
	TxValue *protos.TransactionValue
	ID      int
}
