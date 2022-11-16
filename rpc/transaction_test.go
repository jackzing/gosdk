package rpc

import (
	_ "embed"
	"github.com/jackzing/gosdk/account"
	"github.com/jackzing/gosdk/bvm"
	"github.com/jackzing/gosdk/common/hexutil"
	"github.com/jackzing/gosdk/kvsql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction(t *testing.T) {
	tax := NewTransaction("0x0000000000000000")
	tax.setTxVersion("1.0")
	tax.SetNonce(int64(1))
	tax.SetExtra("extra")
	tax.SetFrom("0x0000000000000000")
	tax.SetHasExtra(true)
	tax.SetIsDeploy(false)
	tax.SetIsInvoke(false)
	tax.SetIsMaintain(true)
	tax.SetIsPrivateTxm(false)
	tax.Simulate(false)
	tax.SetIsValue(false)
	tax.SetOpcode(1)
	tax.SetParticipants([]string{})
	tax.SetTo("0x0000000000000000")
	tax.SetVmType("HVM")
	tax.Nonce(int64(1))
	tax.Timestamp(int64(1))
	tax.Value(int64(1))
	tax.SetPayload("nothing")
	tax.SetValue(int64(123))
	tax.SetTimestamp(int64(1))
	tax.SetSignature("signature")
	tax.SerializeToString()
}

func TestNeedHashString(t *testing.T) {
	tax := NewTransaction("0x0000000000000000")
	tax.SetNonce(int64(1))
	tax.SetExtra("extra")
	tax.SetFrom("0x0000000000000000")
	tax.SetHasExtra(true)
	tax.SetIsDeploy(false)
	tax.SetIsInvoke(false)
	tax.SetIsMaintain(true)
	tax.SetIsPrivateTxm(false)
	tax.Simulate(false)
	tax.SetIsValue(false)
	tax.SetOpcode(1)
	tax.SetParticipants([]string{})
	tax.SetTo("0x0000000000000000")
	tax.SetVmType("HVM")
	tax.Nonce(int64(1))
	tax.Timestamp(int64(1))
	tax.Value(int64(1))
	tax.SetPayload("nothing")
	tax.SetValue(int64(123))
	tax.SetTimestamp(int64(1))

	tax.setTxVersion("1.8")
	expect18 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM"
	assert.Equal(t, expect18, needHashString(tax))

	tax.setTxVersion("2.0")
	expect20 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.0"
	assert.Equal(t, expect20, needHashString(tax))

	tax.setTxVersion("2.1")
	expect21 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.1&extraid="
	assert.Equal(t, expect21, needHashString(tax))

	tax.setTxVersion("2.2")
	expect22 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.2&extraid=&cname="
	assert.Equal(t, expect22, needHashString(tax))

	tax.setTxVersion("2.3")
	expect23 := "from=0x0000000000000000&to=0x0000000000000000&value=0x7b&payload=0xnothing&timestamp=0x1&nonce=0x1&opcode=1&extra=extra&vmtype=HVM&version=2.3&extraid=&cname="
	assert.Equal(t, expect23, needHashString(tax))
}

func TestAddKVSQLType(t *testing.T) {
	// todo open after kvsql
	t.Skip("latest flato not support kvsql")
	hrpc := NewRPCWithPath("../conf")

	js, err := account.NewAccountSm2("12345678")
	assert.Equal(t, nil, err)

	gmAcc, err := account.GenKeyFromAccountJson(js, "12345678")

	assert.Equal(t, nil, err)
	newAddress := gmAcc.(*account.SM2Key).GetAddress()
	transaction := NewTransaction(newAddress.Hex())
	transaction.VMType(KVSQL)
	transaction.Deploy(hexutil.Encode([]byte("KVSQL")))

	transaction.Sign(gmAcc)

	// 建库
	txReceipt, err := hrpc.SignAndDeployContract(transaction, gmAcc)
	assert.Equal(t, nil, err)

	// 建表
	str := "CREATE TABLE IF NOT EXISTS testTable (id bigint(20) NOT NULL, name varchar(32) NOT NULL, exp bigint(20), money double(16,2) NOT NULL DEFAULT '99', primary key (id), unique key name (name));"
	tranInvoke := NewTransaction(newAddress.Hex()).InvokeSql(txReceipt.ContractAddress, []byte(str))
	tranInvoke.VMType(KVSQL)
	tranInvoke.Sign(gmAcc)
	res, err := hrpc.SignAndInvokeContract(tranInvoke, gmAcc)
	assert.Equal(t, nil, err)
	t.Log(res.Ret)

	// 插入
	in := "insert into testTable (id, name, exp, money) values (1, 'test', 1, 1.1);"
	tranInvoke2 := NewTransaction(newAddress.Hex()).InvokeSql(txReceipt.ContractAddress, []byte(in))
	tranInvoke2.VMType(KVSQL)
	tranInvoke2.Sign(gmAcc)
	_, err = hrpc.SignAndInvokeContract(tranInvoke2, gmAcc)
	assert.Equal(t, nil, err)

	// 获取
	sl := "select * from testTable where id = 1;"
	tranInvoke3 := NewTransaction(newAddress.Hex()).InvokeSql(txReceipt.ContractAddress, []byte(sl))
	tranInvoke3.VMType(KVSQL)
	tranInvoke3.Sign(gmAcc)
	res, err = hrpc.SignAndInvokeContract(tranInvoke3, gmAcc)
	assert.Equal(t, nil, err)
	t.Log(res.Ret)
	b, _ := hexutil.Decode(res.Ret)
	rs := kvsql.DecodeRecordSet(b)
	t.Log(rs)
}

func TestTransaction_Serialize(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	js, err := account.NewAccountSm2("12345678")
	assert.Equal(t, nil, err)

	gmAcc, err := account.GenKeyFromAccountJson(js, "12345678")

	assert.Equal(t, nil, err)
	newAddress := gmAcc.(*account.SM2Key).GetAddress()
	transaction := NewTransaction(newAddress.Hex())
	transaction.VMType(KVSQL)
	transaction.Deploy(hexutil.Encode([]byte("KVSQL")))
	transaction.txVersion = rp.txVersion
	transaction.Sign(gmAcc)
	method := CONTRACT + "deployContract"
	assert.Equal(t, nil, err)
	transaction.isDeploy = true
	//param := transaction.Serialize()
	param := make(map[string]interface{})
	param["from"] = transaction.from

	if !(transaction.isDeploy || transaction.isByName) {
		param["to"] = transaction.to
	}

	param["timestamp"] = transaction.timestamp
	param["nonce"] = transaction.nonce

	if !transaction.isMaintain {
		param["simulate"] = transaction.simulate
	}

	param["vmType"] = transaction.vmType

	if transaction.isValue {
		param["value"] = transaction.value
	} else if transaction.isMaintain && (transaction.opcode == 2 || transaction.opcode == 3) {

	} else {
		param["payload"] = transaction.payload
	}

	param["signature"] = transaction.signature

	if transaction.isMaintain || transaction.isDID {
		param["opcode"] = transaction.opcode
	}

	if transaction.hasExtra {
		param["extra"] = transaction.extra
	}

	if transaction.extraIdInt64 != nil || len(transaction.extraIdInt64) > 0 {
		param["extraIdInt64"] = transaction.extraIdInt64
	}
	if transaction.extraIdString != nil || len(transaction.extraIdString) > 0 {
		param["extraIdString"] = transaction.extraIdString
	}
	if transaction.cName != "" {
		param["cName"] = transaction.cName
	}
	if transaction.optionExtra != "" {
		param["optionExtra"] = transaction.optionExtra
	}

	rp.CallByPolling(method, param, transaction.isPrivateTx)
}

//pki certs
var (
	//go:embed pkiAccount/p256_idcert.pfx
	p256IDCert []byte
	//go:embed pkiAccount/p256_idcert.cert
	p256IDCertCert []byte
	//go:embed pkiAccount/sm2_idcert.pfx
	sm2IDCert []byte
	//go:embed pkiAccount/sm2_idcert.cert
	sm2IDCertCert []byte
	//go:embed pkiAccount/k1_idcert.pfx
	k1IDCert []byte
	//go:embed pkiAccount/k1_idcert.cert
	k1IDCertCert []byte
)

/*
	获取PKI账户的步骤：
	1. certgen或者第三方生成证书，证书中cn为小写的不带'0x'前缀的hex地址，证书类型为idcert, 地址需要有5字节0xff前缀
	2. 使用openssl将私钥和证书转换为pfx格式
	3. 注册账户，将地址和证书注册
	PS：
	0. certgen生成证书命令：
		certgen gc --priv ./p256_idcert.key -c p256 --cn "ffffffffff0102030405060708090a0b0c0d0e0f" --ct 'p256IDCert' --to 2100-01-01 root1.ca ./root1.priv ./p256_idcert.cert
	1. pem转pfx：
		openssl pkcs12 -export -out certificate.pfx -inkey privateKey.key -in certificate.crt
	2. pfx转pem：
		openssl pkcs12 -in p256IDCert.pfx -out p256IDCert.cert -clcerts -nokeys
		openssl pkcs12 -in p256IDCert.pfx -out p256IDCert.key.pem -nocerts -nodes
*/
func TestPKIAccountTx(t *testing.T) {
	//suggest using center MSP
	adminJSON := `{"address":"0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd","version":"4.0","algo":"0x03",
"publicKey":"0400ddbadb932a0d276e257c6df50599a425804a3743f40942d031f806bf14ab0c57aed6977b1ad14646672f9b9ce385f2c98c4581267b611f48f4b7937de386ac",
"privateKey":"16acbf6b4f09a476a35ebd4c01e337238b5dceceb6ff55ff0c4bd83c4f91e11b"}`
	adminInterface, err := account.GenKeyFromAccountJson(adminJSON, "")
	assert.Equal(t, err, nil)
	admin := adminInterface.(*account.ECDSAKey)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("p256 pki", func(t *testing.T) {
		testPKIAccount(t, admin, "ffffffffff0102030405060708090a0b0c0d0e0f", p256IDCert, p256IDCertCert, rp)
	})
	t.Run("sm2 pki", func(t *testing.T) {
		testPKIAccount(t, admin, "ffffffffff0102030405060708090a0b0c0d0e10", sm2IDCert, sm2IDCertCert, rp)
	})
	t.Run("k1 pki", func(t *testing.T) {
		testPKIAccount(t, admin, "ffffffffff0102030405060708090a0b0c0d0e11", k1IDCert, k1IDCertCert, rp)
	})
}

func testPKIAccount(t *testing.T, admin *account.ECDSAKey, addr string, pfx, cert []byte, rp *RPC) {
	key, err := account.NewAccountPKI("123456", pfx)
	if err != nil {
		t.Error(err)
		return
	}
	newAddress := key.GetAddress()
	assert.Equal(t, err, nil)

	t.Run("register", func(t *testing.T) {
		//注册需要admin操作
		//这里传入的addr需要和证书中的一致
		op := bvm.NewAccountRegisterOperation(addr, cert)
		payload1 := bvm.EncodeOperation(op)
		tx1 := NewTransaction(admin.GetAddress().Hex()).Invoke(op.Address(), payload1).VMType(BVM)
		receipt, err := rp.SignAndInvokeContract(tx1, admin)
		assert.Equal(t, err, nil)
		rec, derr := hexutil.Decode(receipt.Ret)
		assert.Equal(t, derr, nil)
		t.Log(string(rec))
	})

	t.Run("send tx", func(t *testing.T) {
		tx2 := NewTransaction(newAddress.Hex()).Transfer(newAddress.Hex(), int64(0))
		receipt, err := rp.SignAndSendTx(tx2, key)
		assert.Equal(t, err, nil)
		t.Log(receipt.Ret)
	})
}
