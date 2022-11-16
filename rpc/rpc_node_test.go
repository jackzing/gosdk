package rpc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRPC_GetNodes(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	nodeInfo, err := rp.GetNodes()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(nodeInfo))
	logger.Info(nodeInfo)
}

func TestRPC_GetNodeHash(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	nodeHash, err := rp.GetNodeHash()
	assert.Nil(t, err)
	assert.NotNil(t, nodeHash)
}

func TestRPC_GetNodeHashById(t *testing.T) {
	id := 1
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	hashByID, err := rp.GetNodeHashByID(id)
	assert.Nil(t, err)
	assert.NotNil(t, hashByID)
}

func TestRPC_DeleteNodeVP(t *testing.T) {
	t.Skip("do not delete VP in CI")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	hash1, _ := rp.GetNodeHashByID(1)
	success, _ := rp.DeleteNodeVP(hash1)
	assert.Equal(t, true, success)

	hash11, _ := rp.GetNodeHashByID(1)
	fmt.Println(hash11)
}

func TestRPC_DeleteNodeNVP(t *testing.T) {
	t.Skip("do not delete NVP in CI")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	hash1, _ := rp.GetNodeHashByID(6)
	rp, _ = rp.BindNodes(1)
	success, _ := rp.DeleteNodeNVP(hash1)
	assert.Equal(t, true, success)
}

func TestRPC_DisconnectNodeVP(t *testing.T) {
	t.Skip("do not delete NVP in CI")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	hash1, _ := rp.GetNodeHashByID(1)
	rp, _ = rp.BindNodes(6)
	success, _ := rp.DisconnectNodeVP(hash1)
	assert.Equal(t, true, success)
}

func TestRPC_ReplaceNodeCerts(t *testing.T) {
	t.Skip("do not delete NVP in CI")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	rp, _ = rp.BindNodes(1)
	hash1, err := rp.ReplaceNodeCerts("node1")
	assert.Nil(t, err)
	fmt.Println(hash1)
}

func TestRPC_GetNodeStates(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	states, err := rp.GetNodeStates()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 4, len(states))
}
