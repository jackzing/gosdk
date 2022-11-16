package rpc

import (
	"encoding/json"
	"github.com/hyperchain/go-hpc-common/types"
	"github.com/jackzing/gosdk/account"
	"github.com/jackzing/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRPC_GetAccountProof(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	accountStr := "e93b92f1da08f925bdee44e91e7768380ae83307"
	res, err := rp.GetAccountProof(accountStr)
	if err != nil {
		t.Error(err)
		return
	}
	ast := assert.New(t)
	ast.True(ValidateAccountProof(accountStr, res))
}

func TestRPC_GetTxProof(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, _ := rp.GetLatestBlock()
	info, err := rp.GetTxByBlkHashAndIdx(block.Hash, 0)
	if err != nil {
		t.Error(err)
		return
	}
	res, err2 := rp.GetTxProof(info.Hash)
	if err2 != nil {
		t.Error(err2)
		return
	}
	ast := assert.New(t)
	ast.True(ValidateTxProof(info.Hash, block.TxRoot, res))
}

func TestRPC_ArchiveSnapshot(t *testing.T) {
	t.Skip("need flato 1.5.0")
	deployJar, err := DecompressFromJar("../hvmtestfile/fibonacci/fibonacci-1.0-fibonacci.jar")
	if err != nil {
		t.Error(err)
	}
	accountJson, sysErr := account.NewAccountJson(account.SMRAW, "")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	key, sysErr := account.GenKeyFromAccountJson(accountJson, "")
	if sysErr != nil {
		logger.Error(sysErr)
		return
	}
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	newAddress := key.(*account.SM2Key).GetAddress()
	transaction := NewTransaction(newAddress.Hex()).Deploy(common.Bytes2Hex(deployJar)).VMType(HVM)
	transaction.Sign(key)
	receipt, err := rp.DeployContract(transaction)
	assert.Nil(t, err)
	t.Log("contract address:", receipt.ContractAddress)

	b, err := rp.GetLatestBlock()
	assert.Nil(t, err)

	time.Sleep(20 * time.Second)
	id, err := rp.Snapshot(b.Number)
	assert.Nil(t, err)
	t.Log(b.Number)
	t.Log(b.MerkleRoot)
	t.Log(id)
}

func TestRPC_GetStateProof(t *testing.T) {
	// should send to archiveReader!
	t.Skip("need archiveReader")
	id := "0x5b1a5bb7b10d15bc9d47701eed9c9349"
	seq := 2
	contractAddr := "0x6de31be7a30204189d70bd202340c6d9b395523e"
	merkleRoot := "0xaa2fd673656f4bada6ff6d8588498239eeb3202214a24005d6cf0138a9f30a79"
	proofParam := &ProofParam{
		Meta: &LedgerMetaParam{
			SnapshotID: id,
			SeqNo:      uint64(seq),
		},
		Key: &KeyParam{
			Address:   types.HexToAddress(contractAddr),
			FieldName: "hyperMap1",
			Params:    []string{"key1"},
			VMType:    "HVM",
		},
	}
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	proof, err := rp.GetStateProof(proofParam)
	assert.Nil(t, err)
	t.Log(proof)
	b, _ := json.Marshal(proof)
	t.Log(string(b))

	ok, err := rp.ValidateStateProof(proofParam, proof, merkleRoot)
	assert.Nil(t, err)
	assert.True(t, ok)
}
