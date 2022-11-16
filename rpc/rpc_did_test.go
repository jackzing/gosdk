package rpc

import (
	"fmt"
	"github.com/jackzing/gosdk/account"
	"github.com/jackzing/gosdk/bvm"
	"github.com/jackzing/gosdk/common"
	"github.com/jackzing/gosdk/common/hexutil"
	"github.com/jackzing/gosdk/hvm"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const genesisAccountJson = `{"address":"0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd","version":"4.0", "algo":"0x03","publicKey":"0400ddbadb932a0d276e257c6df50599a425804a3743f40942d031f806bf14ab0c57aed6977b1ad14646672f9b9ce385f2c98c4581267b611f48f4b7937de386ac","privateKey":"16acbf6b4f09a476a35ebd4c01e337238b5dceceb6ff55ff0c4bd83c4f91e11b"}`

func TestDID(t *testing.T) {

	t.Run("setchainID", func(t *testing.T) {
		testDID_SetChainID(t)
	})

	t.Run("getchainID", func(t *testing.T) {
		testDID_GetChainID(t)
	})

	t.Run("DID_account", func(t *testing.T) {
		testDID_Account(t)
	})

	t.Run("DID_Credential", func(t *testing.T) {
		testDID_Credential(t)
	})

	t.Run("DID_Set_GetExtra", func(t *testing.T) {
		testDID_Set_GetExtra(t)
	})

	t.Run("DID_EXTransaction", func(t *testing.T) {
		testDID_EXTransaction(t)
	})
}

func testDID_SetChainID(t *testing.T) {
	//t.Skip("skip this test")
	rpc := NewRPC()
	key, _ := account.GenKeyFromAccountJson(genesisAccountJson, "")
	opt := bvm.NewDIDSetChainIDOperation("chainID_01")
	payload := bvm.EncodeOperation(opt)
	tx := NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	_, err := rpc.SignAndInvokeContract(tx, key)
	assert.Nil(t, err)
	//result := bvm.Decode(re.Ret)

}

func testDID_GetChainID(t *testing.T) {
	//t.Skip("skip this test")
	rpc := NewRPC()
	res, _ := rpc.GetNodeChainID()
	fmt.Println(res)

}

func registerDIDAccount(rpc *RPC, admins []string) *account.DIDKey {
	rpc.SetLocalChainID()
	randNum := common.RandInt(100)
	var accountJson string
	password := "hyper"
	switch randNum % 3 {
	case 0:
		accountJson, _ = account.NewAccountED25519(password)
	case 1:
		accountJson, _ = account.NewAccountSm2(password)
	case 2:
		accountJson, _ = account.NewAccountDID(password)
	}
	key, _ := account.GenDIDKeyFromAccountJson(accountJson, password)
	suffix := common.RandomString(10)
	didKey := account.NewDIDAccount(key.(account.Key), rpc.chainID, suffix)

	puKey, _ := GenDIDPublicKeyFromDIDKey(didKey)
	document := NewDIDDocument(didKey.GetAddress(), puKey, admins)
	tx := NewTransaction(didKey.GetAddress()).Register(document)
	_, err := rpc.SendDIDTransaction(tx, didKey)
	if err == nil {
		return didKey
	}
	return nil
}

func testDID_Account(t *testing.T) {
	//t.Skip("skip this test")
	rpc := NewRPC()
	adminAccount := registerDIDAccount(rpc, nil)
	ac := registerDIDAccount(rpc, []string{adminAccount.GetAddress()})

	tx := NewTransaction(adminAccount.GetAddress()).MaintainDID(ac.GetAddress(), DID_FREEZE)
	_, err := rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)
	tx = NewTransaction(adminAccount.GetAddress()).MaintainDID(ac.GetAddress(), DID_UNFREEZE)
	_, err = rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)

	tempAc := registerDIDAccount(rpc, nil)
	puKey, _ := GenDIDPublicKeyFromDIDKey(tempAc)
	tx = NewTransaction(adminAccount.GetAddress()).UpdatePublicKey(ac.GetAddress(), puKey)
	_, err = rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)

	tx = NewTransaction(adminAccount.GetAddress()).UpdateAdmins(ac.GetAddress(), []string{tempAc.GetAddress()})
	_, err = rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)

	adminAccount = tempAc
	tx = NewTransaction(adminAccount.GetAddress()).MaintainDID(ac.GetAddress(), DID_UNFREEZE)
	_, err = rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)
	tx = NewTransaction(adminAccount.GetAddress()).MaintainDID(ac.GetAddress(), DID_ABANDON)
	_, err = rpc.SendDIDTransaction(tx, adminAccount)
	assert.Nil(t, err)
}

func testDID_Credential(t *testing.T) {
	//t.Skip("skip this test")
	rpc := NewRPC()
	holder := registerDIDAccount(rpc, nil)
	issuer := registerDIDAccount(rpc, nil)

	cred := NewDIDCredential("type", issuer.GetAddress(), holder.GetAddress(), "", time.Now().UnixNano(), time.Now().UnixNano()+1e11)
	err := cred.Sign(issuer)
	assert.Nil(t, err)
	tx := NewTransaction(holder.GetAddress()).UploadCredential(cred)
	_, err = rpc.SendDIDTransaction(tx, holder)
	assert.Nil(t, err)

	tx = NewTransaction(issuer.GetAddress()).DownloadCredential(cred.ID)
	_, err = rpc.SendDIDTransaction(tx, issuer)
	assert.Nil(t, err)

	tx = NewTransaction(holder.GetAddress()).DestroyCredential("ddas")
	_, err = rpc.SendDIDTransaction(tx, holder)
	assert.NotNil(t, err)
	//res, _ := hexutil.Decode(receipt.Ret)
	//assert.True(t, strings.Contains(string(res), "the credential does not exist"))

}

func testDID_Set_GetExtra(t *testing.T) {
	t.Skip("skip this test")
	rpc := NewRPC()
	from := registerDIDAccount(rpc, nil)
	to := registerDIDAccount(rpc, []string{from.GetAddress()})
	key := "key_test"
	value := "value_test"

	tx := NewTransaction(from.GetAddress()).DIDSetExtra(to.GetAddress(), key, value)
	_, err := rpc.SendDIDTransaction(tx, from)
	assert.Nil(t, err)

	tx1 := NewTransaction(from.GetAddress()).DIDGetExtra(to.GetAddress(), key)
	receipt, err := rpc.SendDIDTransaction(tx1, from)
	assert.Nil(t, err)
	res, _ := hexutil.Decode(receipt.Ret)
	assert.Equal(t, value, string(res))
}

func testDID_EXTransaction(t *testing.T) {
	//t.Skip("skip this test")
	rpc := NewRPC()
	ac := registerDIDAccount(rpc, nil)
	_, err := rpc.GetDIDDocument(ac.GetAddress())
	assert.Nil(t, err)

	_, err = rpc.GetNodeChainID()
	assert.Nil(t, err)

	holder := registerDIDAccount(rpc, nil)
	issuer := registerDIDAccount(rpc, nil)
	cred := NewDIDCredential("type", issuer.GetAddress(), holder.GetAddress(), "", time.Now().UnixNano(), time.Now().UnixNano()+1e11)
	cred.Sign(issuer)
	tx := NewTransaction(holder.GetAddress()).UploadCredential(cred)
	_, err = rpc.SendDIDTransaction(tx, holder)
	assert.Nil(t, err)

	_, err = rpc.GetCredentialPrimaryMessage(cred.ID)
	assert.Nil(t, err)

	valid, _ := rpc.CheckCredentialValid(cred.ID)
	assert.True(t, valid)

	abandoned, _ := rpc.CheckCredentialAbandoned(cred.ID)
	assert.True(t, !abandoned)
}

func TestDID_contract(t *testing.T) {
	t.Skip("skip this test")
	rpc := NewRPC()
	ac := registerDIDAccount(rpc, nil)

	deployJar, err := DecompressFromJar("../hvmtestfile/fibonacci/fibonacci-1.0-fibonacci.jar")
	if err != nil {
		t.Error(err)
	}

	transaction := NewTransaction(ac.GetAddress()).Deploy(common.Bytes2Hex(deployJar)).VMType(HVM)
	receipt, err := rpc.SignAndDeployContract(transaction, ac)
	assert.Nil(t, err)

	abiPath := "../hvmtestfile/fibonacci/hvm.abi"
	abiJson, rerr := common.ReadFileAsString(abiPath)
	assert.Nil(t, rerr)
	abi, gerr := hvm.GenAbi(abiJson)
	if gerr != nil {
		logger.Error(gerr)
	}

	easyBean := "invoke.InvokeFibonacci"
	beanAbi, err := abi.GetBeanAbi(easyBean)
	if err != nil {
		logger.Error(err)
	}

	payload, err := hvm.GenPayload(beanAbi)
	if err != nil {
		logger.Error(err)
	}
	transaction1 := NewTransaction(ac.GetAddress()).Invoke(receipt.ContractAddress, payload).VMType(HVM)
	invokeContract, err := rpc.SignAndInvokeContract(transaction1, ac)
	if err != nil {
		t.Error(err)
	}
	t.Log(invokeContract)
}

func TestDID_R1Algo(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	rp.SetLocalChainID()
	var accountJson string
	password := "hyper"
	accountJson, _ = account.NewAccountR1(password)
	key, _ := account.GenDIDKeyFromAccountJson(accountJson, password)
	suffix := common.RandomString(10)
	didKey := account.NewDIDAccount(key.(account.Key), rp.chainID, suffix)

	puKey, _ := GenDIDPublicKeyFromDIDKey(didKey)
	document := NewDIDDocument(didKey.GetAddress(), puKey, nil)
	tx := NewTransaction(didKey.GetAddress()).Register(document)
	_, err = rp.SendDIDTransaction(tx, didKey)
	assert.Nil(t, err)
	tx = NewTransaction(didKey.GetAddress()).Deploy(binContract).VMType(EVM)
	re, err := rp.SignAndDeployContract(tx, didKey)
	assert.Nil(t, err)
	assert.NotNil(t, re)
	t.Log(re.ContractAddress)
}
