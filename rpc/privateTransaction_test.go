package rpc

import (
	"fmt"
	"github.com/hyperchain/go-crypto-standard/asym"
	"github.com/hyperchain/go-crypto-standard/hash"
	"github.com/jackzing/gosdk/common"
	"strings"
	"testing"

	"github.com/jackzing/gosdk/abi"
	"github.com/jackzing/gosdk/account"
)

func generateAccount() *account.ECDSAKey {
	accountJSON, err := account.NewAccount("12345678")
	if err != nil {
		logger.Error(err)
		return nil
	}

	key, err := account.NewAccountFromAccountJSON(accountJSON, "12345678")
	if err != nil {
		logger.Error(err)
		return nil
	}
	return key
}

func deployPrivateContract(key *account.ECDSAKey, cr *CompileResult) (*TxReceipt, []string) {
	rp, err := NewJsonRPC()
	if err != nil {
		logger.Error(err)
		return nil, nil
	}
	rpcAPI, _ := rp.BindNodes(1)
	nodeHashList := make([]string, 0)
	nodeHash, err := rpcAPI.GetNodeHashByID(1)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}
	nodeHashList = append(nodeHashList, nodeHash)
	pubKeyBytes, _ := key.Public().(*asym.ECDSAPublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKeyBytes)
	address := h[12:]

	tranDeploy := NewPrivateTransaction(common.BytesToAddress(address).Hex(), nodeHashList).Deploy(cr.Bin[0])
	tranDeploy.Sign(key)
	txDeploy, stdErr := rpcAPI.DeployContract(tranDeploy)
	if stdErr != nil {
		logger.Error(stdErr)
		return nil, nil
	}
	return txDeploy, nodeHashList
}

func invokePrivateTransaction(key *account.ECDSAKey) *TxReceipt {
	rp, err := NewJsonRPC()
	if err != nil {
		logger.Error(err)
		return nil
	}
	rpcAPI, _ := rp.BindNodes(1)
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	txReceipt, nodeHashList := deployPrivateContract(key, cr)

	// invoke
	ABI, _ := abi.JSON(strings.NewReader(cr.Abi[0]))
	packed, _ := ABI.Pack("add", uint32(1), uint32(2))

	kvExtra := NewKVExtra()
	//nolint
	kvExtra.AddKV("name", "ice")
	//nolint
	kvExtra.AddKV("friend", map[string]string{
		"name": "ice",
		"age":  "24",
	})
	pubKeyBytes, _ := key.Public().(*asym.ECDSAPublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKeyBytes)
	address := h[12:]

	tranInvoke := NewPrivateTransaction(common.BytesToAddress(address).Hex(), nodeHashList).Invoke(txReceipt.ContractAddress, packed).KVExtra(kvExtra)
	tranInvoke.Sign(key)
	txInvoke, stdErr := rpcAPI.InvokeContract(tranInvoke)
	if stdErr != nil {
		logger.Error(stdErr)
		return nil
	}
	return txInvoke
}

func TestRPC_DeployPrivateTransaction(t *testing.T) {
	t.Skip("no private")
	key := generateAccount()
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	txReceipt, _ := deployPrivateContract(key, cr)
	fmt.Println(txReceipt)
	logger.Info("contract address:", txReceipt.ContractAddress)
}

func TestRPC_InvokePrivateTransaction(t *testing.T) {
	t.Skip("no private")
	key := generateAccount()
	txInvoke := invokePrivateTransaction(key)
	logger.Info("invoke transaction hash:", txInvoke.TxHash)
}
