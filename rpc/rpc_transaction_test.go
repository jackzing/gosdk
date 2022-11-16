package rpc

import (
	"fmt"
	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/go-crypto-standard/hash"
	"github.com/hyperchain/gosdk/abi"
	"github.com/hyperchain/gosdk/common"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRPC_GetTransactions(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	txs, err := rp.GetTransactionsByBlkNum(block.Number-1, block.Number)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(txs[0].Invalid)
}

func TestRPC_GetTransactionsWithLimit(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	metadata := &Metadata{
		PageSize: 1,
		Bookmark: &Bookmark{
			BlockNumber: block.Number,
			TxIndex:     0,
		},
		Backward: false,
	}

	pageResult, err := rp.GetTransactionsByBlkNumWithLimit(block.Number-1, block.Number, metadata)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pageResult)
}

func TestRPC_GetInvalidTransactionsWithLimit(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	metadata := &Metadata{
		PageSize: 1,
		Bookmark: &Bookmark{
			BlockNumber: 1,
			TxIndex:     0,
		},
		Backward: false,
	}

	pageResult, err := rp.GetInvalidTransactionsByBlkNumWithLimit(block.Number-1, block.Number, metadata)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pageResult)
}

func TestRPC_GetInvalidTxByBlockNumber(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	txInfos, err := rp.GetInvalidTransactionsByBlkNum(block.Number)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(txInfos)
}

func TestRPC_GetInvalidTxByBlockHash(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	txInfos, err := rp.GetInvalidTransactionsByBlkHash(block.Hash)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(txInfos)
}

func TestRPC_GetInvalidTxsCount(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	count, err := rp.GetInvalidTxCount()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(count)
}

func TestRPC_GetDiscardTx(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	txs, err := rp.GetDiscardTx()
	if err != nil {
		//t.Error(err)
		return
	}
	fmt.Println(len(txs))
	if len(txs) > 0 {
		fmt.Println(txs[len(txs)-1].Hash)
	}
}

func TestRPC_GetTransactionByHash(t *testing.T) {
	//t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	transaction.Sign(guomiKey)
	receipt, _ := rp.DeployContract(transaction)
	fmt.Println("txhash:", receipt.TxHash)

	hash := receipt.TxHash
	tx, err := rp.GetTransactionByHash(hash)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tx.Hash)
	assert.Equal(t, receipt.TxHash, tx.Hash)
}

func TestRPC_GetBatchTxByHash(t *testing.T) {
	//t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction1 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	transaction1.Sign(guomiKey)
	receipt1, _ := rp.DeployContract(transaction1)
	fmt.Println("txhash1:", receipt1.TxHash)

	// 模拟一个无法查询到交易的hash
	txHash2 := "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	txhashes := make([]string, 0)
	txhashes = append(txhashes, receipt1.TxHash, txHash2)

	txs, err := rp.GetBatchTxByHash(txhashes)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(txs))
	fmt.Println(txs[0].Hash, txs[1].Hash)
	assert.Equal(t, receipt1.TxHash, txs[0].Hash)
	assert.Equal(t, txHash2, txs[1].Hash)
}

func TestRPC_GetTxByBlkHashAndIdx(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, _ := rp.GetLatestBlock()
	info, err := rp.GetTxByBlkHashAndIdx(block.Hash, 0)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(info)
	assert.EqualValues(t, 66, len(info.Hash))
}

func TestRPC_GetTxByBlkNumAndIdx(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, _ := rp.GetLatestBlock()
	info, err := rp.GetTxByBlkNumAndIdx(block.Number, 0)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(info)
	assert.EqualValues(t, 66, len(info.Hash))
}

func TestRPC_GetTxAvgTimeByBlockNumber(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, _ := rp.GetLatestBlock()
	avgTime, err := rp.GetTxAvgTimeByBlockNumber(block.Number-2, block.Number)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(avgTime)
}

func TestRPC_GetBatchReceipt(t *testing.T) {
	//t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, _ := rp.GetLatestBlock()
	trans, _ := rp.GetTransactionsByBlkNum(block.Number-2, block.Number)
	// 模拟一个无法查询到回执的hash
	txHash2 := "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	hashes := []string{trans[0].Hash, txHash2}
	txs, err := rp.GetBatchReceipt(hashes)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, len(txs))
}

func TestRPC_GetTransactionsCountByTime(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	count, err := rp.GetTransactionsCountByTime(1, uint64(time.Now().UnixNano()))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(count)
}

func TestRPC_GetTxCountByContractAddr(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	guomiKey := testGuomiKey()
	cAddress, _ := deployContract(binContract, abiContract)
	ABI, _ := abi.JSON(strings.NewReader(abiContract))
	packed, _ := ABI.Pack("getSum")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction.Sign(guomiKey)
	//nolint
	rp.InvokeContract(transaction)
	transaction2 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction2.Sign(guomiKey)
	//nolint
	rp.InvokeContract(transaction2)

	block, _ := rp.GetLatestBlock()
	count, err := rp.GetTxCountByContractAddr(block.Number-1, block.Number, cAddress, false)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 2, count.Count)
}

func TestRPC_GetTransactionsCountByMethodID(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	cAddress, _ := deployContract(binContract, abiContract)
	abiStr, _ := abi.JSON(strings.NewReader(abiContract))
	methodID := string(abiStr.Constructor.Id())

	block, _ := rp.GetLatestBlock()
	count, err := rp.GetTransactionsCountByMethodID(block.Number-1, block.Number, cAddress, methodID)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(count)
}

func TestRPC_GetTxByTime(t *testing.T) {
	//t.Skip("the length of result is too long")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	infos, err := rp.GetTxByTime(1, uint64(time.Now().UnixNano()))
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(infos)
	assert.EqualValues(t, true, len(infos) > 0)
}

func TestRPC_GetTxByTimeWithLimit(t *testing.T) {
	t.Skip()
	metadata := &Metadata{
		PageSize: 1,
		Bookmark: &Bookmark{
			BlockNumber: 1,
			TxIndex:     0,
		},
		Backward: false,
	}
	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	pageResult, err := rp.GetTxByTimeWithLimit(1, uint64(time.Now().UnixNano()), metadata)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pageResult)
}

func TestRPC_GetDiscardTransactionsByTime(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	infos, err := rp.GetDiscardTransactionsByTime(1, uint64(time.Now().UnixNano()))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(infos)
}

func TestRPC_GetNextPageTxs(t *testing.T) {
	//t.Skip("hyperchain snapshot will case error")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	guomiKey := testGuomiKey()
	cAddress, _ := deployContract(binContract, abiContract)
	ABI, _ := abi.JSON(strings.NewReader(abiContract))
	packed, _ := ABI.Pack("getSum")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction.Sign(guomiKey)
	//nolint
	rp.InvokeContract(transaction)
	transaction2 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction2.Sign(guomiKey)
	//nolint
	rp.InvokeContract(transaction2)

	block, _ := rp.GetLatestBlock()

	infos, err := rp.GetNextPageTxs(block.Number-10, 0, 1, block.Number, 0, 10, false, cAddress)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 3, len(infos))

	t.Skip()
	txs, err := rp.GetNextPageTxs(block.Number-10, 0, 1, block.Number, 0, 10, false, "")
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 10, len(txs))
}

func TestRPC_GetPrevPageTxs(t *testing.T) {
	//t.Skip("hyperchain snapshot will case error")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	guomiKey := testGuomiKey()
	cAddress, _ := deployContract(binContract, abiContract)
	ABI, _ := abi.JSON(strings.NewReader(abiContract))
	packed, _ := ABI.Pack("getSum")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction.Sign(guomiKey)
	//nolint
	rp.InvokeContract(transaction)
	transaction2 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(cAddress, packed)
	transaction2.Sign(guomiKey)
	//nolint
	rp.InvokeContract(transaction2)

	block, _ := rp.GetLatestBlock()

	infos, err := rp.GetPrevPageTxs(block.Number, 0, 1, block.Number, 0, 10, false, cAddress)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 2, len(infos))

	t.Skip()
	txs, err := rp.GetPrevPageTxs(block.Number-10, 0, 1, block.Number, 0, 10, false, "")
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 10, len(txs))
}

func TestRPC_GetBlkTxCountByHash(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	count, err := rp.GetBlkTxCountByHash(block.Hash)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(count)
}

func TestRPC_GetBlkTxCountByNumber(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, err := rp.GetLatestBlock()
	hex := "0x" + strconv.FormatUint(block.Number, 16)
	fmt.Println("=====", block, hex)
	if err != nil {
		t.Error(err)
		return
	}
	count, err := rp.GetBlkTxCountByNumber(hex)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(count)
}

func TestRPC_GetTxCount(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	txCount, err := rp.GetTxCount()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(txCount.Count)
}
