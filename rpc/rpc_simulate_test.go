package rpc

import (
	"fmt"
	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/go-crypto-standard/hash"
	"github.com/jackzing/gosdk/abi"
	"github.com/jackzing/gosdk/common"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRPC_SimulateDeployContract(t *testing.T) {
	t.Skip("not support simulate")
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Simulate(true).Deploy("0x60606040526000805463ffffffff19169055341561001957fe5b5b61012a806100296000396000f300606060405263ffffffff60e060020a6000350416633ad14af38114603e57806348fe842114605c578063569c5f6d14606b578063d09de08a146091575bfe5b3415604557fe5b605a63ffffffff6004358116906024351660a0565b005b3415606357fe5b605a60c2565b005b3415607257fe5b607860d2565b6040805163ffffffff9092168252519081900360200190f35b3415609857fe5b605a60df565b005b6000805463ffffffff808216850184011663ffffffff199091161790555b5050565b6000805463ffffffff191690555b565b60005463ffffffff165b90565b6000805463ffffffff8082166001011663ffffffff199091161790555b5600a165627a7a72305820caa934a33fe993d03f87bdf39706fada68ddde78182e0110fd43e8c323d5984a0029")
	transaction.Sign(guomiKey)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	receipt, _ := rp.DeployContract(transaction)
	fmt.Println("address:", receipt.ContractAddress)
}

func TestRPC_SimulateInvokeContract(t *testing.T) {
	t.Skip("not support simulate")
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy("0x60606040526000805463ffffffff19169055341561001957fe5b5b61012a806100296000396000f300606060405263ffffffff60e060020a6000350416633ad14af38114603e57806348fe842114605c578063569c5f6d14606b578063d09de08a146091575bfe5b3415604557fe5b605a63ffffffff6004358116906024351660a0565b005b3415606357fe5b605a60c2565b005b3415607257fe5b607860d2565b6040805163ffffffff9092168252519081900360200190f35b3415609857fe5b605a60df565b005b6000805463ffffffff808216850184011663ffffffff199091161790555b5050565b6000805463ffffffff191690555b565b60005463ffffffff165b90565b6000805463ffffffff8082166001011663ffffffff199091161790555b5600a165627a7a72305820caa934a33fe993d03f87bdf39706fada68ddde78182e0110fd43e8c323d5984a0029")
	transaction.Sign(guomiKey)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	receipt, err := rp.DeployContract(transaction)
	assert.Nil(t, err, "??????????????????")
	fmt.Println("address:", receipt.ContractAddress)

	contractAddress := receipt.ContractAddress
	ABI, serr := abi.JSON(strings.NewReader(abiContract))
	if serr != nil {
		t.Error(serr)
		return
	}
	packed, serr := ABI.Pack("add", uint32(1), uint32(2))
	if serr != nil {
		t.Error(serr)
		return
	}
	address, privateKey := testPrivateAccount()
	transaction = NewTransaction(address).Invoke(contractAddress, packed).Simulate(true)
	transaction.Sign(privateKey)
	receipt, _ = rp.InvokeContract(transaction)
	fmt.Println("ret:", receipt.Ret)
}

func TestRPC_SimulateInvokeContract_Increment(t *testing.T) {
	t.Skip("not support simulate")
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy("0x60606040526000805463ffffffff19169055341561001957fe5b5b61012a806100296000396000f300606060405263ffffffff60e060020a6000350416633ad14af38114603e57806348fe842114605c578063569c5f6d14606b578063d09de08a146091575bfe5b3415604557fe5b605a63ffffffff6004358116906024351660a0565b005b3415606357fe5b605a60c2565b005b3415607257fe5b607860d2565b6040805163ffffffff9092168252519081900360200190f35b3415609857fe5b605a60df565b005b6000805463ffffffff808216850184011663ffffffff199091161790555b5050565b6000805463ffffffff191690555b565b60005463ffffffff165b90565b6000805463ffffffff8082166001011663ffffffff199091161790555b5600a165627a7a72305820caa934a33fe993d03f87bdf39706fada68ddde78182e0110fd43e8c323d5984a0029")
	transaction.Sign(guomiKey)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	receipt, err := rp.DeployContract(transaction)
	assert.Nil(t, err, "??????????????????")
	fmt.Println("address:", receipt.ContractAddress)

	contractAddress := receipt.ContractAddress
	ABI, serr := abi.JSON(strings.NewReader(abiContract))
	if serr != nil {
		t.Error(serr)
		return
	}
	packed, serr := ABI.Pack("increment")
	if serr != nil {
		t.Error(serr)
		return
	}
	address, privateKey := testPrivateAccount()
	transaction = NewTransaction(address).Invoke(contractAddress, packed).Simulate(true)
	transaction.Sign(privateKey)
	receipt, _ = rp.InvokeContract(transaction)
	fmt.Println("ret:", receipt.Ret)
}

func TestRPC_SimulateInvokeContract_GetSum(t *testing.T) {
	t.Skip("not support simulate")
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy("0x60606040526000805463ffffffff19169055341561001957fe5b5b61012a806100296000396000f300606060405263ffffffff60e060020a6000350416633ad14af38114603e57806348fe842114605c578063569c5f6d14606b578063d09de08a146091575bfe5b3415604557fe5b605a63ffffffff6004358116906024351660a0565b005b3415606357fe5b605a60c2565b005b3415607257fe5b607860d2565b6040805163ffffffff9092168252519081900360200190f35b3415609857fe5b605a60df565b005b6000805463ffffffff808216850184011663ffffffff199091161790555b5050565b6000805463ffffffff191690555b565b60005463ffffffff165b90565b6000805463ffffffff8082166001011663ffffffff199091161790555b5600a165627a7a72305820caa934a33fe993d03f87bdf39706fada68ddde78182e0110fd43e8c323d5984a0029")
	transaction.Sign(guomiKey)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	receipt, err := rp.DeployContract(transaction)
	assert.Nil(t, err, "??????????????????")
	fmt.Println("address:", receipt.ContractAddress)

	contractAddress := receipt.ContractAddress
	ABI, serr := abi.JSON(strings.NewReader(abiContract))
	if serr != nil {
		t.Error(serr)
		return
	}
	packed, serr := ABI.Pack("getSum")
	if serr != nil {
		t.Error(serr)
		return
	}
	address, privateKey := testPrivateAccount()
	transaction = NewTransaction(address).Invoke(contractAddress, packed).Simulate(true)
	transaction.Sign(privateKey)
	receipt, _ = rp.InvokeContract(transaction)
	if receipt != nil {
		fmt.Println("ret:", receipt.Ret)
	} else {
		fmt.Println("nil receipt")
	}
}

func TestRPC_SimulateInvokeContract_GetSum_Range(t *testing.T) {
	t.Skip("not support simulate")
	guomiKey := testGuomiKey()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy("0x60606040526000805463ffffffff19169055341561001957fe5b5b61012a806100296000396000f300606060405263ffffffff60e060020a6000350416633ad14af38114603e57806348fe842114605c578063569c5f6d14606b578063d09de08a146091575bfe5b3415604557fe5b605a63ffffffff6004358116906024351660a0565b005b3415606357fe5b605a60c2565b005b3415607257fe5b607860d2565b6040805163ffffffff9092168252519081900360200190f35b3415609857fe5b605a60df565b005b6000805463ffffffff808216850184011663ffffffff199091161790555b5050565b6000805463ffffffff191690555b565b60005463ffffffff165b90565b6000805463ffffffff8082166001011663ffffffff199091161790555b5600a165627a7a72305820caa934a33fe993d03f87bdf39706fada68ddde78182e0110fd43e8c323d5984a0029")
	transaction.Sign(guomiKey)
	receipt, err := rp.DeployContract(transaction)
	assert.Nil(t, err, "??????????????????")
	fmt.Println("address:", receipt.ContractAddress)

	contractAddress := receipt.ContractAddress

	address, privateKey := testPrivateAccount()
	var wg sync.WaitGroup
	var count int32
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			ABI, serr := abi.JSON(strings.NewReader(abiContract))
			if serr != nil {
				t.Error(serr)
				return
			}
			packed, serr := ABI.Pack("getSum")
			if serr != nil {
				t.Error(serr)
				return
			}
			transaction := NewTransaction(address).Invoke(contractAddress, packed).Simulate(true)
			transaction.Sign(privateKey)
			receipt, iErr := rp.InvokeContract(transaction)
			if receipt != nil {
				fmt.Println("ret:", receipt.Ret)
			} else {
				fmt.Printf("nil receipt, err: %s\n", iErr)
			}
			wg.Done()
			atomic.AddInt32(&count, 1)
		}()
	}
	wg.Wait()
	t.Logf("total success count: %d\n", atomic.LoadInt32(&count))
}

// maintain contract by opcode 1
func TestRPC_SimulateMaintainContract(t *testing.T) {
	t.Skip("not support simulate")
	guomiKey := testGuomiKey()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy("0x60606040526000805463ffffffff19169055341561001957fe5b5b61012a806100296000396000f300606060405263ffffffff60e060020a6000350416633ad14af38114603e57806348fe842114605c578063569c5f6d14606b578063d09de08a146091575bfe5b3415604557fe5b605a63ffffffff6004358116906024351660a0565b005b3415606357fe5b605a60c2565b005b3415607257fe5b607860d2565b6040805163ffffffff9092168252519081900360200190f35b3415609857fe5b605a60df565b005b6000805463ffffffff808216850184011663ffffffff199091161790555b5050565b6000805463ffffffff191690555b565b60005463ffffffff165b90565b6000805463ffffffff8082166001011663ffffffff199091161790555b5600a165627a7a72305820caa934a33fe993d03f87bdf39706fada68ddde78182e0110fd43e8c323d5984a0029")
	transaction.Sign(guomiKey)
	receipt, err := rp.DeployContract(transaction)
	assert.Nil(t, err, "??????????????????")
	fmt.Println("address:", receipt.ContractAddress)

	contractAddress := receipt.ContractAddress
	fmt.Println("-----------------------------")

	transactionUpdate := NewTransaction(common.BytesToAddress(newAddress).Hex()).Maintain(1, contractAddress, binContract).Simulate(true)
	transactionUpdate.Sign(guomiKey)
	receiptUpdate, err := rp.MaintainContract(transactionUpdate)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("address:", receiptUpdate.ContractAddress)
}

// maintain contract by opcode 2 and 3
func TestRPC_SimulateMaintainContract2(t *testing.T) {
	t.Skip("not support simulate")
	guomiKey := testGuomiKey()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy("0x60606040526000805463ffffffff19169055341561001957fe5b5b61012a806100296000396000f300606060405263ffffffff60e060020a6000350416633ad14af38114603e57806348fe842114605c578063569c5f6d14606b578063d09de08a146091575bfe5b3415604557fe5b605a63ffffffff6004358116906024351660a0565b005b3415606357fe5b605a60c2565b005b3415607257fe5b607860d2565b6040805163ffffffff9092168252519081900360200190f35b3415609857fe5b605a60df565b005b6000805463ffffffff808216850184011663ffffffff199091161790555b5050565b6000805463ffffffff191690555b565b60005463ffffffff165b90565b6000805463ffffffff8082166001011663ffffffff199091161790555b5600a165627a7a72305820caa934a33fe993d03f87bdf39706fada68ddde78182e0110fd43e8c323d5984a0029")
	transaction.Sign(guomiKey)
	receipt, err := rp.DeployContract(transaction)
	assert.Nil(t, err, "??????????????????")
	fmt.Println("address:", receipt.ContractAddress)

	contractAddress := receipt.ContractAddress

	// freeze contract
	transactionFreeze := NewTransaction(common.BytesToAddress(newAddress).Hex()).Maintain(2, contractAddress, "")
	transactionFreeze.Sign(guomiKey)
	receiptFreeze, err := rp.MaintainContract(transactionFreeze)
	assert.Nil(t, err)
	fmt.Println(receiptFreeze.TxHash)

	//// unfreeze contract
	transactionUnfreeze := NewTransaction(common.BytesToAddress(newAddress).Hex()).Maintain(3, contractAddress, "").Simulate(true)
	transactionUnfreeze.Sign(guomiKey)
	receiptUnFreeze, err := rp.MaintainContract(transactionUnfreeze)
	assert.Nil(t, err)
	fmt.Println(receiptUnFreeze.TxHash)
}
