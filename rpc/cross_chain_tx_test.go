package rpc

import (
	"errors"
	"github.com/hyperchain/gosdk/hvm"
	"testing"

	"github.com/hyperchain/go-hpc-common/utils"
	"github.com/hyperchain/gosdk/account"
	"github.com/hyperchain/gosdk/bvm"
	"github.com/hyperchain/gosdk/common"

	"github.com/stretchr/testify/assert"
)

// TestRPC_BVMCrossChainAnchorTxContent handle the user-operated register stage1 and stage2,
// make sure the bvm inner parse logic can handle the tx payload
func TestRPC_BVMCrossChainAnchorTxContent(t *testing.T) {
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], "")
	assert.Nil(t, err)
	t.Run("register_tx1_tx2", func(t *testing.T) {
		op := bvm.NewSystemAnchorOperation(bvm.RegisterAnchor, "node1", "ns1")
		payload := bvm.EncodeOperation(op)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(op.Address(), payload).VMType(BVM)
		methodName, params, derr := bvmDecode(common.Hex2Bytes(tx.payload))
		assert.Nil(t, derr)
		assert.Equal(t, methodName, string(bvm.RegisterAnchor))
		assert.Equal(t, len(params), 2)
		assert.Equal(t, params[0], "node1")
		assert.Equal(t, params[1], "ns1")
	})
	t.Run("unregister_tx1", func(t *testing.T) {
		op := bvm.NewSystemAnchorOperation(bvm.UnRegisterAnchor, "node1", "ns1")
		payload := bvm.EncodeOperation(op)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(op.Address(), payload).VMType(BVM)
		methodName, params, derr := bvmDecode(common.Hex2Bytes(tx.payload))
		assert.Nil(t, derr)
		assert.Equal(t, methodName, string(bvm.UnRegisterAnchor))
		assert.Equal(t, len(params), 2)
		assert.Equal(t, params[0], "node1")
		assert.Equal(t, params[1], "ns1")
	})
	t.Run("replace_tx1", func(t *testing.T) {
		operation := bvm.NewSystemAnchorOperation(bvm.ReplaceAnchor, "node2", "ns1", "node1", "ns1")
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		methodName, params, derr := bvmDecode(common.Hex2Bytes(tx.payload))
		assert.Nil(t, derr)
		assert.Equal(t, methodName, string(bvm.ReplaceAnchor))
		assert.Equal(t, len(params), 4)
		assert.Equal(t, params[0], "node2")
		assert.Equal(t, params[1], "ns1")
		assert.Equal(t, params[2], "node1")
		assert.Equal(t, params[3], "ns1")
	})
	t.Run("crossChain_tx1", func(t *testing.T) {
		invokeDirectMagic := []byte{254, 255, 251, 206}
		abiPath := "../hvmtestfile/crosschain/hvm.abi"
		abiJson, rerr := common.ReadFileAsString(abiPath)
		assert.Nil(t, rerr)
		abi, gerr := hvm.GenAbi(abiJson)
		assert.Nil(t, gerr)
		addAbi, abierr := abi.GetMethodAbi("Add")
		assert.Nil(t, abierr)
		payload, perr := hvm.GenPayload(addAbi, "11")
		assert.Nil(t, perr)
		assert.Equal(t, string(payload[0:4]), string(invokeDirectMagic))
		binPayload := payload[4:]
		methodNameLen := (uint16(binPayload[0]) << 8) | uint16(binPayload[1])
		assert.Equal(t, methodNameLen, uint16(3))
	})
	t.Run("crossChain_timeout", func(t *testing.T) {
		crossChainID := "ns1__@__0x16b22229e93fdad80cef0920cc26d2fb6acb32945f8f5fe2248f0adcc3aa2c29__@__ns2"
		operation := bvm.NewSystemAnchorOperation(bvm.Timeout, crossChainID)
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		methodName, params, derr := bvmDecode(common.Hex2Bytes(tx.payload))
		assert.Nil(t, derr)
		assert.Equal(t, methodName, string(bvm.Timeout))
		assert.Equal(t, len(params), 1)
		assert.Equal(t, params[0], crossChainID)
	})
}

// util part for tx content check

// bvmDecode decode bvm payload, copied from flato-common@v0.2.39
func bvmDecode(payload []byte) (string, []string, error) {
	defaultLen := 4
	if len(payload) < defaultLen {
		return "", nil, errors.New("invalid payload")
	}
	methodNameLen := utils.BytesToInt32(payload[:defaultLen])
	if len(payload) < defaultLen+methodNameLen {
		return "", nil, errors.New("invalid payload")
	}
	methodName := string(payload[defaultLen : defaultLen+methodNameLen])
	if len(payload) < defaultLen+methodNameLen+defaultLen {
		return "", nil, errors.New("invalid payload")
	}
	paramCount := utils.BytesToInt32(payload[defaultLen+methodNameLen : defaultLen+methodNameLen+defaultLen])
	paramsInput := payload[defaultLen+methodNameLen+defaultLen:]
	params := make([]string, paramCount)
	index := 0
	for i := 0; i < paramCount; i++ {
		if len(paramsInput) < index+defaultLen {
			return "", nil, errors.New("invalid payload")
		}
		paramLen := utils.BytesToInt32(paramsInput[index : index+defaultLen])
		if len(paramsInput) < index+defaultLen+paramLen {
			return "", nil, errors.New("invalid payload")
		}
		param := string(paramsInput[index+defaultLen : index+defaultLen+paramLen])
		params[i] = param
		index += defaultLen + paramLen
	}
	return methodName, params, nil
}
