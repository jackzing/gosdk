package rpc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestRPC_GetLatestBlock(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(block)
}

func TestRPC_GetBlocks(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	latestBlock, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	blocks, err := rp.GetBlocks(latestBlock.Number-1, latestBlock.Number, true)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(blocks)
}

func TestRPC_GetBlocksWithLimit(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	latestBlock, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	metadata := &Metadata{
		PageSize: 5,
	}

	pageResult, err := rp.GetBlocksWithLimit(latestBlock.Number-1, latestBlock.Number, true, metadata)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(pageResult)
}

func TestRPC_GetBlockByHash(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	latestBlock, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	block, err := rp.GetBlockByHash(latestBlock.Hash, true)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(block)
}

func TestRPC_GetBatchBlocksByHash(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	latestBlock, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	blocks, err := rp.GetBatchBlocksByHash([]string{latestBlock.Hash}, true)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(blocks)
}

func TestRPC_GetBlockByNumber(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	latestBlock, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	//nolint
	rp.GetBlockByNumber("latest", false)
	block, err := rp.GetBlockByNumber(latestBlock.Number, true)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(block)
}

func TestRPC_GetBatchBlocksByNumber(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	latestBlock, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}

	blocks, err := rp.GetBatchBlocksByNumber([]uint64{latestBlock.Number}, true)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(blocks)
}

func TestRPC_GetAvgGenTimeByBlkNum(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	block, err := rp.GetLatestBlock()
	if err != nil {
		t.Error(err)
		return
	}
	avgTime, err := rp.GetAvgGenTimeByBlockNum(block.Number-2, block.Number)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(avgTime)
}

func TestRPC_GetBlockByTime(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	start := time.Now().Unix() - 1
	end := time.Now().Unix()
	blockInterval, err := rp.GetBlocksByTime(uint64(start), uint64(end))
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(blockInterval)
}

func TestRPC_QueryTPS(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	tpsInfo, err := rp.QueryTPS(1, 1778959217012956575)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(tpsInfo)
}

func TestRPC_GetGenesisBlock(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	blkNum, err := rp.GetGenesisBlock()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, true, strings.HasPrefix(blkNum, "0x"))
}

func TestRPC_GetChainHeight(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	blkNum, err := rp.GetChainHeight()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, true, strings.HasPrefix(blkNum, "0x"))
}
