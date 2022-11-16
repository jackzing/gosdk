package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"
	"time"

	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/go-crypto-standard/hash"
	"github.com/hyperchain/go-hpc-common/types"
	"github.com/jackzing/gosdk/abi2"
	"github.com/jackzing/gosdk/hvm"

	"github.com/jackzing/gosdk/abi"
	"github.com/jackzing/gosdk/account"
	"github.com/jackzing/gosdk/common"
	"github.com/stretchr/testify/assert"
)

const (
	abiContract = `[{"constant":false,"inputs":[{"name":"num1","type":"uint32"},{"name":"num2","type":"uint32"}],"name":"add","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"archiveSum","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"getSum","outputs":[{"name":"","type":"uint32"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"increment","outputs":[],"payable":false,"type":"function"}]`
	binContract = "0x60606040526000805463ffffffff19169055341561001957fe5b5b61012a806100296000396000f300606060405263ffffffff60e060020a6000350416633ad14af38114603e57806348fe842114605c578063569c5f6d14606b578063d09de08a146091575bfe5b3415604557fe5b605a63ffffffff6004358116906024351660a0565b005b3415606357fe5b605a60c2565b005b3415607257fe5b607860d2565b6040805163ffffffff9092168252519081900360200190f35b3415609857fe5b605a60df565b005b6000805463ffffffff808216850184011663ffffffff199091161790555b5050565b6000805463ffffffff191690555b565b60005463ffffffff165b90565b6000805463ffffffff8082166001011663ffffffff199091161790555b5600a165627a7a72305820caa934a33fe993d03f87bdf39706fada68ddde78182e0110fd43e8c323d5984a0029"
)

func testGuomiKey() *account.SM2Key {
	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	_ = pri.FromBytes(common.FromHex(guomiPri), 0)
	return &account.SM2Key{
		SM2PrivateKey: &gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
}

func testPrivateAccount() (string, *account.ECDSAKey) {
	address := "bfa5bd992e3eb123c8b86ebe892099d4e9efb783"
	privateKey, _ := account.NewAccountFromPriv("a1fd6ed6225e76aac3884b5420c8cdbb4fde1db01e9ef773415b8f2b5a9b77d4")
	return address, privateKey
}

func TestRPC_New(t *testing.T) {
	rpc, err := NewJsonRPC()
	assert.Nil(t, err)
	logger.Info(rpc)
}

func TestRPC_NewWithPath(t *testing.T) {
	rpc, err := NewJsonRPCWithPath("../conf")
	assert.Nil(t, err)
	logger.Info(rpc)
}

func TestRPC_DefaultRPC(t *testing.T) {
	rpc := DefaultRPC(NewNode("localhost", "8081", "11001")).Https("../conf/certs/tls/tlsca.ca", "../conf/certs/tls/tls_peer.cert", "../conf/certs/tls/tls_peer.priv").Tcert(true, "../conf/certs/sdkcert.cert", "../conf/certs/sdkcert.priv", "../conf/certs/unique.pub", "../conf/certs/unique.priv")
	logger.Info(rpc)
}

/*---------------------------------- archive ----------------------------------*/

func TestRPC_Snapshot(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	res, err := rp.Snapshot(1)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_FlatoVersion_Snapshot(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	res, err := rp.MakeSnapshot4Flato(1)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_QuerySnapshotExist(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	res, err := rp.QuerySnapshotExist("0x5d86cce7e537cd0e0346468889801196")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_CheckSnapshot(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	res, err := rp.CheckSnapshot("0x5d86cce7e537cd0e0346468889801196")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_Archive(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	res, err := rp.Archive("0x5d86cce7e537cd0e0346468889801196", false)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_Restore(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	res, err := rp.Restore("0x5d86cce7e537cd0e0346468889801196", false)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}
func TestRPC_RestoreAll(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	res, err := rp.RestoreAll(false)
	if err != nil {
		t.Error(err)

	}

	fmt.Println(res)
}

func TestRPC_QueryArchiveExist(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	res, err := rp.QueryArchiveExist("0x5d86cce7e537cd0e0346468889801196")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_QueryLatestArchive(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	res, err := rp.QueryLatestArchive()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

func TestRPC_Pending(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	res, err := rp.QueryArchiveExist("0x5d86cce7e537cd0e0346468889801196")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(res)
}

/*---------------------------------- config ----------------------------------*/

func TestRPC_GetProposal(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	proposal, err := rp.GetProposal()
	if err != nil {
		t.Error(err)
		return
	}
	bytes, _ := json.Marshal(proposal)
	fmt.Println(string(bytes))
}

func TestRPC_GetHosts(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	data, err := rp.GetHosts("vp")
	if err != nil {
		t.Error(err)
		return
	}
	bytes, _ := json.Marshal(data)
	fmt.Println(string(bytes))
}

func TestRPC_GetVSet(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	data, err := rp.GetVSet()
	if err != nil {
		t.Error(err)
		return
	}
	bytes, _ := json.Marshal(data)
	fmt.Println(string(bytes))
}

func TestRPC_GetConfig(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	config, err := rp.GetConfig()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(config)
}

var accountJsons = []string{
	//ecdsa
	`{"address":"0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd","version":"4.0", "algo":"0x03","publicKey":"0400ddbadb932a0d276e257c6df50599a425804a3743f40942d031f806bf14ab0c57aed6977b1ad14646672f9b9ce385f2c98c4581267b611f48f4b7937de386ac","privateKey":"16acbf6b4f09a476a35ebd4c01e337238b5dceceb6ff55ff0c4bd83c4f91e11b"}`,
	`{"address":"0x6201cb0448964ac597faf6fdf1f472edf2a22b89","version":"4.0", "algo":"0x03","publicKey":"04e482f140d70a1b8ec8185cc699db5b391ea5a7b8e93e274b9f706be9efdaec69542eb32a61421ba6219230b9cf87bf849fa01c1d10a8d298cbe3dcfa5954134c","privateKey":"21ff03a654c939f0c9b83e969aaa9050484aa4108028094ee2e927ba7e7d1bbb"}`,
	`{"address":"0xb18c8575e3284e79b92100025a31378feb8100d6","version":"4.0", "algo":"0x03","publicKey":"042169a7260acaff308228579aab2a2c6b3a790922c6a6b58b218cdd7ce0b1db0fbfa6f68737a452010b9d138187b8321288cae98f07fc758bb67bb818292cab9b","privateKey":"aa9c83316f68c17bcc21cf20a4733ae2b2bf76ad1c745f634c0ebf7d5094500e"}`,
	`{"address":"0xe93b92f1da08f925bdee44e91e7768380ae83307","version":"4.0","algo":"0x03","publicKey":"047196daf5d4d1fe339da58e2fe0543bbfec9a464b76546f180facdcc56315b8eeeca50474100f15fb17606695ce24a1f8e5a990600c1c4ea9787ba4dd65c8ce3e","privateKey":"8cdfbe86deb690e331453a84a98c956f0422dd1e783c3a02aed9180b1f4516a9"}`,
	//sm2
	`{"address":"0xfbca6a7e9e29728773b270d3f00153c75d04e1ad","version":"4.0","algo":"0x13","publicKey":"049c330d0aea3d9c73063db339b4a1a84d1c3197980d1fb9585347ceeb40a5d262166ee1e1cb0c29fd9b2ef0e4f7a7dfb1be6c5e759bf411c520a616863ee046a4","privateKey":"5f0a3ea6c1d3eb7733c3170f2271c10c1206bc49b6b2c7e550c9947cb8f098e3"}`,
	`{"address":"0x856e2b9a5fa82fd1b031d1ff6863864dbac7995d","publicKey":"047ea464762c333762d3be8a04536b22955d97231062442f81a3cff46cb009bbdbb0f30e61ade5705254d4e4e0c0745fb3ba69006d4b377f82ecec05ed094dbe87","privateKey":"71b9acc4ee2b32b3d2c79b5abe9e118e5f73765aee5e7755d6aa31f12945036d"}`,
}

var r1Account = `{"address":"0xd2548c5e47c54d7ae8e1319cefa10d5832f37542","version":"4.0", "algo":"0x031","publicKey":"04183125ce9ddcee9c2b190dd2fff4257b042c3c09a5cb6613a936d4aec6e880bfe8a94f11a4efdbc62c5cbe25d3acc3287b1c63215e113108d0d7e2dd8f0b4f3c","privateKey":"4921395a4105bd52fc7d8d183ecd385f5896e4a3d79bd3b8fe0ee96fb473f68b"}`

func TestRPC_GetAllRoles(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	roles, err := rp.GetAllRoles()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(roles)
}

func TestRPC_IsRoleExist(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	exist, err := rp.IsRoleExist("admin")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(exist)
}

func TestRPC_GetAddressByName(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	addr, err := rp.GetAddressByName("HashContract")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(addr)
}

func TestRPC_GetNameByAddress(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	name, err := rp.GetNameByAddress("0x0000000000000000000000000000000000ffff01")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(name)
}

func TestRPC_GetAllCNS(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	all, err := rp.GetAllCNS()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(all)
}

func TestRPC_GetGenesisInfo(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	all, err := rp.GetGenesisInfo()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(all)
}

func TestRPC_AddRoleForNode(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	_, privateKey := testPrivateAccount()
	rp.SetAccount(privateKey)
	err = rp.AddRoleForNode(privateKey.GetAddress().Hex(), "accountManager")
	assert.Nil(t, err)
}

func TestRPC_GetRoleFromNode(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	roles, err := rp.GetRoleFromNode("0x6201cb0448964ac597faf6fdf1f472edf2a22b89")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(roles)
}

func TestRPC_GetAddressFromNode(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	roles, err := rp.GetAddressFromNode("accountManager")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(roles)
}

func TestRPC_GetAllRolesFromNode(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	roles, err := rp.GetAllRolesFromNode()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(roles)
}

func TestRPC_DeleteRoles(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	err = rp.DeleteRoleFromNode("0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd", "accountManager")
	assert.Nil(t, err)
}

func TestRPC_SetRulesInNode(t *testing.T) {
	t.Skip()
	rule := &InspectorRule{
		Name:            "rule",
		ID:              1,
		AllowAnyone:     false,
		AuthorizedRoles: []string{"accountManager"},
		Method:          []string{"account_*"},
	}
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	err = rp.SetRulesInNode([]*InspectorRule{rule})
	assert.Nil(t, err)
}

func TestRPC_GetRulesFromNode(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	_, privateKey := testPrivateAccount()
	rp.SetAccount(privateKey)
	rules, err := rp.GetRulesFromNode()
	if err != nil {
		t.Error(err)
		return
	}
	marshal, _ := json.Marshal(rules)
	t.Log(string(marshal))
}

func TestRPC_SetAccount(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	_, privateKey := testPrivateAccount()
	rp.SetAccount(key)
	balance, stdError := rp.GetBalance(privateKey.GetAddress().Hex())
	t.Log(stdError)
	t.Log(balance)
}

/**************************** self function ******************************/

func compileContract(path string) (*CompileResult, error) {
	contract, _ := common.ReadFileAsString(path)
	rp, err := NewJsonRPC()
	if err != nil {
		return nil, err
	}
	cr, err := rp.CompileContract(contract)
	if err != nil {
		logger.Error("can not get compile return, ", err.Error())
		return nil, err
	}
	fmt.Println("abi:", cr.Abi[0])
	fmt.Println("bin:", cr.Bin[0])
	fmt.Println("type:", cr.Types[0])

	return cr, err
}

func decode(contractAbi abi.ABI, v interface{}, method string, ret string) (result interface{}) {
	if err := contractAbi.UnpackResult(v, method, ret); err != nil {
		logger.Error(NewSystemError(err).String())
	}
	result = v
	return result
}

func deployContract(bin, abi string, params ...interface{}) (string, StdError) {
	var transaction *Transaction
	var err StdError
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	if len(params) == 0 {

		transaction = NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(bin)
	} else {
		transaction = NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(bin).DeployArgs(abi, params)
	}
	transaction.Sign(guomiKey)
	rp, serr := NewJsonRPC()
	if serr != nil {
		return "", NewSystemError(serr)
	}
	txReceipt, err := rp.DeployContract(transaction)
	if err != nil {
		logger.Error(err)
	}
	return txReceipt.ContractAddress, nil
}

func Test_demo1(t *testing.T) {
	t.Skip("solc")
	contract, _ := common.ReadFileAsString("../conf/contract/contract02/assembly.sol")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	guomiKey := testGuomiKey()
	cr, _ := rp.CompileContract(contract)
	abiStr := cr.Abi[0]
	bin := cr.Bin[0]
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	deployTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(bin)
	deployTx.Sign(guomiKey)
	deployRe, _ := rp.DeployContract(deployTx)
	contractAddress := deployRe.ContractAddress
	ABI, _ := abi2.JSON(strings.NewReader(abiStr))
	{
		packed1, _ := ABI.Pack("addition", big.NewInt(100001), big.NewInt(100001))
		invokeTx1 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed1)
		invokeTx1.Sign(guomiKey)
		invokeRe1, _ := rp.InvokeContract(invokeTx1)

		var p0 *big.Int
		if err := ABI.UnpackResult(&p0, "addition", invokeRe1.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(p0)
	}
}

func Test_demo2(t *testing.T) {
	abiStr := "[ { \"inputs\": [], \"name\": \"g\", \"outputs\": [ { \"internalType\": \"bytes4\", \"name\": \"\", \"type\": \"bytes4\" } ], \"stateMutability\": \"pure\", \"type\": \"function\" } ]"
	bin := "608060405234801561001057600080fd5b5060de8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063e2179b8e14602d575b600080fd5b60336047565b604051603e91906063565b60405180910390f35b600063b3de648b60e01b905090565b605d81607c565b82525050565b6000602082019050607660008301846056565b92915050565b60007fffffffff000000000000000000000000000000000000000000000000000000008216905091905056fea26469706673582212204f79d51af42d5aa5cddf512cc1550b5db9dde2fb7ef962170d1fcab6ea59e4da64736f6c63430008040033"
	guomiKey := testGuomiKey()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	deployTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(bin)
	deployTx.Sign(guomiKey)
	deployRe, _ := rp.DeployContract(deployTx)
	contractAddress := deployRe.ContractAddress
	ABI, _ := abi2.JSON(strings.NewReader(abiStr))
	{
		packed1, _ := ABI.Pack("g")
		invokeTx1 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed1)
		invokeTx1.Sign(guomiKey)
		invokeRe1, _ := rp.InvokeContract(invokeTx1)

		var p0 [4]byte
		if err := ABI.UnpackResult(&p0, "g", invokeRe1.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(string(p0[:]))
	}
}

func Test_TypeCheck(t *testing.T) {
	t.Skip("solc")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	guomiKey := testGuomiKey()
	contract, _ := common.ReadFileAsString("../conf/contract/TypeCheck.sol")
	cr, _ := rp.CompileContract(contract)
	abiStr := cr.Abi[0]
	bin := cr.Bin[0]
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	deployTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(bin)
	deployTx.Sign(guomiKey)
	deployRe, _ := rp.DeployContract(deployTx)
	contractAddress := deployRe.ContractAddress

	ABI, _ := abi.JSON(strings.NewReader(abiStr))

	// invoke fun1
	{
		var data32 [32]byte
		copy(data32[:], "data32")
		var data8 [8]byte
		copy(data8[:], "byte8")
		packed1, _ := ABI.Pack("fun1", []byte("data1"), data32, data8)
		invokeTx1 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed1)
		invokeTx1.Sign(guomiKey)
		invokeRe1, _ := rp.InvokeContract(invokeTx1)

		var p0 []byte
		var p1 [32]byte
		var p2 [8]byte
		testV := []interface{}{&p0, &p1, &p2}
		if err := ABI.UnpackResult(&testV, "fun1", invokeRe1.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(string(p0), string(p1[:]), string(p2[:]))
	}

	// invoke fun2
	{
		bigInt1 := big.NewInt(-100001)
		bigInt2 := big.NewInt(-1000001)
		bigInt3 := big.NewInt(10000001)
		int1 := int64(-10001)
		int2 := int8(101)
		packed, _ := ABI.Pack("fun2", bigInt1, bigInt2, bigInt3, int1, int2)
		invokeTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
		invokeTx.Sign(guomiKey)
		invokeRe, _ := rp.InvokeContract(invokeTx)

		var p0 interface{}
		var p1 *big.Int
		var p2 *big.Int
		var p3 interface{}
		var p4 int8
		testV := []interface{}{&p0, &p1, &p2, &p3, &p4}
		if err := ABI.UnpackResult(&testV, "fun2", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(p0, p1.Int64(), p2, p3, p4)
	}

	// invoke fun3
	{
		bigInt1 := big.NewInt(100001)
		bigInt2 := big.NewInt(1000001)
		bigInt3 := big.NewInt(10000001)
		int1 := uint64(10001)
		int2 := uint8(101)
		packed, _ := ABI.Pack("fun3", bigInt1, bigInt2, bigInt3, int1, int2)
		invokeTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
		invokeTx.Sign(guomiKey)
		invokeRe, _ := rp.InvokeContract(invokeTx)

		var p0 interface{}
		var p1 *big.Int
		var p2 *big.Int
		var p3 interface{}
		var p4 uint8
		testV := []interface{}{&p0, &p1, &p2, &p3, &p4}
		if err := ABI.UnpackResult(&testV, "fun3", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(p0, p1, p2, p3, p4)
	}

	// invoke fun4
	{
		bigInt1 := big.NewInt(-100001)
		a16int := int16(-10001)
		bigInt3 := big.NewInt(10001)
		bigInt4 := big.NewInt(1111111)
		a16uint := uint16(10001)
		bigInt6 := big.NewInt(111111)
		packed, _ := ABI.Pack("fun4", bigInt1, a16int, bigInt3, bigInt4, a16uint, bigInt6)
		invokeTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
		invokeTx.Sign(guomiKey)
		invokeRe, _ := rp.InvokeContract(invokeTx)

		var p0 interface{}
		var p1 int16
		var p2 *big.Int
		var p3 interface{}
		var p4 uint16
		var p5 *big.Int
		testV := []interface{}{&p0, &p1, &p2, &p3, &p4, &p5}
		if err := ABI.UnpackResult(&testV, "fun4", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(p0, p1, p2, p3, p4, p5)
	}

	// invoke fun5
	{
		address := common.Address{}
		address.SetString("2312321312")
		packed, _ := ABI.Pack("fun5", "data1", address)
		invokeTx := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
		invokeTx.Sign(guomiKey)
		invokeRe, _ := rp.InvokeContract(invokeTx)

		var p0 string
		var p1 common.Address
		testV := []interface{}{&p0, &p1}
		if err := ABI.UnpackResult(&testV, "fun5", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		fmt.Println(p0, p1)
	}
}

func TestRPC_ListenContract(t *testing.T) {
	t.Skip("hyperchain needs to start the radar service first")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	guomiKey := testGuomiKey()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0])
	transaction.Sign(guomiKey)
	receipt, _ := rp.DeployContract(transaction)
	fmt.Println("address:", receipt.ContractAddress)
	srcCode, _ := ioutil.ReadFile("../conf/contract/Accumulator.sol")
	result, err := rp.ListenContract(string(srcCode), receipt.ContractAddress)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)
}

func TestRPC_Setter(t *testing.T) {
	rp := NewRPC()
	rp.Namespace("global").ResendTimes(int64(1)).FirstPollTime(int64(1))
	rp.FirstPollInterval(int64(1)).SecondPollInterval(int64(1)).SecondPollTime(int64(1))
	rp.ReConnTime(int64(1))
	rp.AddNode("127.0.0.1", "8081", "10001")
	//nolint
	rp.BindNodes(1)
}

func TestRPC_Simulate(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	guomiKey := testGuomiKey()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0]).Simulate(true)
	transaction.Sign(guomiKey)
	receipt, _ := rp.DeployContract(transaction)
	fmt.Println("address:", receipt.ContractAddress)
}

func TestRPC_SignAndInvokeContractCombineReturns(t *testing.T) {
	accountJSON, _ := account.NewAccountED25519("12345678")
	ekey, err := account.GenKeyFromAccountJson(accountJSON, "12345678")
	assert.Nil(t, err)
	newAddress := ekey.(*account.ED25519Key).GetAddress()
	address, _ := testPrivateAccount()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	transaction := NewTransaction(newAddress.Hex()).Transfer(address, int64(0))
	transaction.Sign(ekey)
	txreceipt, info, err := rp.SignAndInvokeContractCombineReturns(transaction, ekey)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, info)
	assert.NotNil(t, txreceipt)
	assert.Equal(t, txreceipt.TxHash, info.Hash)
	assert.Nil(t, err)
}

func TestGetGasPrice(t *testing.T) {
	//t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	price, err := rp.GetGasPrice()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(0), price)
}

func TestV2Contest(t *testing.T) {

	type T struct {
		X *big.Int
		Y *big.Int
	}

	type S struct {
		A *big.Int
		B []*big.Int
		C []T
	}

	abiStr := `[
	{
		"anonymous": false,
		"inputs": [
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "a",
						"type": "uint256"
					},
					{
						"internalType": "uint256[]",
						"name": "b",
						"type": "uint256[]"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "x",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "y",
								"type": "uint256"
							}
						],
						"internalType": "struct Test.T[]",
						"name": "c",
						"type": "tuple[]"
					}
				],
				"indexed": false,
				"internalType": "struct Test.S",
				"name": "ss",
				"type": "tuple"
			},
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "x",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "y",
						"type": "uint256"
					}
				],
				"indexed": false,
				"internalType": "struct Test.T",
				"name": "tt",
				"type": "tuple"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "uu",
				"type": "uint256"
			}
		],
		"name": "Event",
		"type": "event"
	},
	{
		"inputs": [
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "a",
						"type": "uint256"
					},
					{
						"internalType": "uint256[]",
						"name": "b",
						"type": "uint256[]"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "x",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "y",
								"type": "uint256"
							}
						],
						"internalType": "struct Test.T[]",
						"name": "c",
						"type": "tuple[]"
					}
				],
				"internalType": "struct Test.S",
				"name": "ss",
				"type": "tuple"
			},
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "x",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "y",
						"type": "uint256"
					}
				],
				"internalType": "struct Test.T",
				"name": "tt",
				"type": "tuple"
			},
			{
				"internalType": "uint256",
				"name": "uu",
				"type": "uint256"
			}
		],
		"name": "f",
		"outputs": [
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "a",
						"type": "uint256"
					},
					{
						"internalType": "uint256[]",
						"name": "b",
						"type": "uint256[]"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "x",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "y",
								"type": "uint256"
							}
						],
						"internalType": "struct Test.T[]",
						"name": "c",
						"type": "tuple[]"
					}
				],
				"internalType": "struct Test.S",
				"name": "",
				"type": "tuple"
			},
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "x",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "y",
						"type": "uint256"
					}
				],
				"internalType": "struct Test.T",
				"name": "",
				"type": "tuple"
			},
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "g",
		"outputs": [
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "a",
						"type": "uint256"
					},
					{
						"internalType": "uint256[]",
						"name": "b",
						"type": "uint256[]"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "x",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "y",
								"type": "uint256"
							}
						],
						"internalType": "struct Test.T[]",
						"name": "c",
						"type": "tuple[]"
					}
				],
				"internalType": "struct Test.S",
				"name": "",
				"type": "tuple"
			},
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "x",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "y",
						"type": "uint256"
					}
				],
				"internalType": "struct Test.T",
				"name": "",
				"type": "tuple"
			},
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "pure",
		"type": "function"
	}
]`

	bin := "608060405234801561001057600080fd5b50610776806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80636f2be7281461003b578063e2179b8e1461006d575b600080fd5b6100556004803603810190610050919061035b565b61008d565b60405161006493929190610581565b60405180910390f35b6100756100ec565b60405161008493929190610581565b60405180910390f35b610095610103565b61009d610124565b60007f91661573ac280570ddebaf850e5b8163eb47f1a96616e9121884430fb79e657b8686866040516100d293929190610581565b60405180910390a185858592509250925093509350939050565b6100f4610103565b6100fc610124565b6000909192565b60405180606001604052806000815260200160608152602001606081525090565b604051806040016040528060008152602001600081525090565b600061015161014c846105e4565b6105bf565b9050808382526020820190508285604086028201111561017057600080fd5b60005b858110156101a0578161018688826102fa565b845260208401935060408301925050600181019050610173565b5050509392505050565b60006101bd6101b884610610565b6105bf565b905080838252602082019050828560208602820111156101dc57600080fd5b60005b8581101561020c57816101f28882610346565b8452602084019350602083019250506001810190506101df565b5050509392505050565b600082601f83011261022757600080fd5b813561023784826020860161013e565b91505092915050565b600082601f83011261025157600080fd5b81356102618482602086016101aa565b91505092915050565b60006060828403121561027c57600080fd5b61028660606105bf565b9050600061029684828501610346565b600083015250602082013567ffffffffffffffff8111156102b657600080fd5b6102c284828501610240565b602083015250604082013567ffffffffffffffff8111156102e257600080fd5b6102ee84828501610216565b60408301525092915050565b60006040828403121561030c57600080fd5b61031660406105bf565b9050600061032684828501610346565b600083015250602061033a84828501610346565b60208301525092915050565b60008135905061035581610729565b92915050565b60008060006080848603121561037057600080fd5b600084013567ffffffffffffffff81111561038a57600080fd5b6103968682870161026a565b93505060206103a7868287016102fa565b92505060606103b886828701610346565b9150509250925092565b60006103ce8383610505565b60408301905092915050565b60006103e68383610563565b60208301905092915050565b60006103fd8261065c565b610407818561068c565b93506104128361063c565b8060005b8381101561044357815161042a88826103c2565b975061043583610672565b925050600181019050610416565b5085935050505092915050565b600061045b82610667565b610465818561069d565b93506104708361064c565b8060005b838110156104a157815161048888826103da565b97506104938361067f565b925050600181019050610474565b5085935050505092915050565b60006060830160008301516104c66000860182610563565b50602083015184820360208601526104de8282610450565b915050604083015184820360408601526104f882826103f2565b9150508091505092915050565b60408201600082015161051b6000850182610563565b50602082015161052e6020850182610563565b50505050565b60408201600082015161054a6000850182610563565b50602082015161055d6020850182610563565b50505050565b61056c816106ae565b82525050565b61057b816106ae565b82525050565b6000608082019050818103600083015261059b81866104ae565b90506105aa6020830185610534565b6105b76060830184610572565b949350505050565b60006105c96105da565b90506105d582826106b8565b919050565b6000604051905090565b600067ffffffffffffffff8211156105ff576105fe6106e9565b5b602082029050602081019050919050565b600067ffffffffffffffff82111561062b5761062a6106e9565b5b602082029050602081019050919050565b6000819050602082019050919050565b6000819050602082019050919050565b600081519050919050565b600081519050919050565b6000602082019050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b6000819050919050565b6106c182610718565b810181811067ffffffffffffffff821117156106e0576106df6106e9565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f8301169050919050565b610732816106ae565b811461073d57600080fd5b5056fea2646970667358221220644bf522edc6e6267d9b2fed23a91764fee081f2257abd117e966f91733a60df64736f6c63430008040033"

	contractAddress, err := deployContract(bin, abiStr)
	if err != nil {
		t.Error(err)
	}
	logger.Info(contractAddress)

	ABI, _ := abi2.JSON(strings.NewReader(abiStr))

	accJson, _ := account.NewAccountSm2("123")
	logger.Debug(accJson)
	key, _ := account.GenKeyFromAccountJson(accJson, "123")
	rp, terr := NewJsonRPC()
	assert.Nil(t, terr)
	//调用合约
	{
		//调用方法g
		invokePayload, _ := ABI.Pack("g")

		transaction := NewTransaction(key.(*account.SM2Key).GetAddress().Hex()).Invoke(contractAddress, invokePayload).VMType("EVM")
		transaction.Sign(key)
		invokeRe, _ := rp.InvokeContract(transaction)
		logger.Info(invokeRe.Ret)

		var p0 S
		var p1 T
		var p2 *big.Int
		testV := []interface{}{&p0, &p1, &p2}
		if err := ABI.UnpackResult(&testV, "g", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		t.Log(p0, p1, p2)
		ret, _ := json.Marshal(p0)
		t.Log(string(ret))
		assert.Equal(t, "{\"A\":0,\"B\":[],\"C\":[]}", string(ret))
	}
	{
		//调用方法f
		s1 := new(S)
		t1 := new(T)
		t1.X = big.NewInt(1)
		t1.Y = big.NewInt(1)

		s1.A = big.NewInt(1)
		s1.B = []*big.Int{big.NewInt(1), big.NewInt(1)}
		s1.C = []T{
			*t1,
		}
		invokePayload, _ := ABI.Pack("f", s1, t1, big.NewInt(1))

		transaction := NewTransaction(key.(*account.SM2Key).GetAddress().Hex()).Invoke(contractAddress, invokePayload).VMType("EVM")

		invokeRe, err := rp.SignAndInvokeContract(transaction, key)
		if err != nil {
			t.Error(err)
			return
		}
		logger.Info(invokeRe.Ret)
		logger.Info(invokeRe.Log[0].Data)

		var p0 S
		var p1 T
		var p2 *big.Int
		testV := []interface{}{&p0, &p1, &p2}
		if err := ABI.UnpackResult(&testV, "f", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		logger.Info(p0, p1, p2)
	}
}

func TestDeployParallel(t *testing.T) {
	t.Skip()
	key, _ := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	deployJar, err := DecompressFromJar("../hvmtestfile/staticanalyze/parallel-contract-1.0-bank.jar")
	if err != nil {
		t.Error(err)
	}

	tx := NewTransaction(key.GetAddress().Hex()).Deploy(common.Bytes2Hex(deployJar)).VMType(HVM)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	receipt, err := rp.SignAndDeployContract(tx, key)
	assert.Nil(t, err)
	t.Log("contract address:", receipt.ContractAddress)
}

func TestCallParallel(t *testing.T) {
	t.Skip()
	key, _ := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	contractAddr := "0x1e548137be17e1a11f0642c9e22dfda64e61fe6d"
	addr := types.HexToAddress(contractAddr)
	abiPath := "../hvmtestfile/staticanalyze/parallel-contract.abi"
	abiJson, err := common.ReadFileAsString(abiPath)
	assert.Nil(t, err)

	abiMap := hvm.NewAbiMap()
	abi0, err := hvm.GenAbi(abiJson)
	assert.Nil(t, err)
	abiMap.SetAbi(contractAddr, abi0)
	beanName := "transfer"

	gNum, txNum := 1, 1
	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	ch := make(chan int, gNum)

	for i := 0; i < gNum; i++ {
		go func(id int) {
			for j := 0; j < txNum; j++ {
				params := []interface{}{fmt.Sprintf("a-%d_%d", id, j), fmt.Sprintf("b-%d_%d", id, j), 100}
				// 1. generate parallel tx's payload and parallel
				payload, err := hvm.GenParallelPayload(abiMap, addr, beanName, true, params...)
				assert.Nil(t, err)

				// 2. build parallel tx
				tx := NewTransaction(key.GetAddress().Hex()).Invoke(contractAddr, payload).VMType(HVM)
				tx.Sign(key)

				// 3. invoke
				receipt, err := rp.InvokeContract(tx)
				assert.Nil(t, err)
				fmt.Println(receipt)
			}
			ch <- id
		}(i)
	}

	for i := 0; i < gNum; i++ {
		id := <-ch
		fmt.Printf("%d finish\n", id)
	}
}

func Test_node_http(t *testing.T) {
	t.Skip("需要手动执行")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// 手动停一个节点
	for {
		t.Logf("%+v", rp)
		height, err := rp.GetChainHeight()
		assert.Nil(t, err)
		t.Logf("区块高度: %s", height)
		time.Sleep(500 * time.Millisecond)
	}
}
