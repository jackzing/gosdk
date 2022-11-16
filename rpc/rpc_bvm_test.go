package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperchain/go-hpc-common/utils"
	"github.com/jackzing/gosdk/account"
	"github.com/jackzing/gosdk/bvm"
	"github.com/jackzing/gosdk/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var notAdminAddress = "0x2a307e1e5b53863242a465bf99ca6e94947da898"
var managerRole = "manager"
var pwd = "123456"
var cf = `
[filter]
    enable = false
    [[filter.rules]]
    allow_anyone = false
    authorized_roles = ["admin"]
    forbidden_roles = ["20"]
    id = 0
    name = "bvm auth"
    to = ["0x0000000000000000000000000000000000ffff02"]
    vm = ["bvm"]
[consensus]
  algo = "RBFT"

[proposal]
	timeout   = "15m"
	threshold = 3
`

func TestUpdateConsensusAlgo(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("has not set", func(t *testing.T) {
		results := completeProposal(t, rp, key, 6, bvm.NewProposalCreateOperationByConfigOps(bvm.NewUpdateConsensusAlgo("NoxBFT")))
		for _, res := range results {
			t.Logf("code:%v, msg:%v", res.Code, res.Msg)
			assert.True(t, strings.Contains(res.Msg, "has not set value"))
		}
	})

	waitingStateNormal()

	t.Run("update success", func(t *testing.T) {
		results := completeProposal(t, rp, key, 6, bvm.NewProposalCreateOperationByConfigOps(bvm.NewSetConsensusAlgo("RBFT"), bvm.NewUpdateConsensusAlgo("NoxBFT")))
		for _, res := range results {
			t.Logf("code:%v, msg:%v", res.Code, res.Msg)
		}
	})

}

func waitingStateNormal() {
	rp, _ := NewJsonRPC()
	for {

		states, _ := rp.GetNodeStates()
		for _, state := range states {
			if strings.ToUpper(state.Status) == "NORMAL" {
				return
			}
		}
		time.Sleep(time.Second)
	}
}

func TestRPC_BVMSet(t *testing.T) {
	t.Skip()
	key, err := account.NewAccountFromAccountJSON(r1Account, pwd)
	assert.Nil(t, err)

	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	operation := bvm.NewHashSetOperation("0x1231", "0x456")
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))
}

func TestRPC_BVMSetGensisiInfo(t *testing.T) {
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)

	operation := bvm.NewSetGenesisInfoForHpcOperation(&bvm.GenesisInfo{
		GenesisAccount: map[string]string{
			"0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd": "10000",
		},
		GenesisNodes: []*bvm.GenesisNode{
			{
				Hostname: "node1",
				CertContent: `-----BEGIN CERTIFICATE-----
MIICSTCCAfWgAwIBAgIBATAKBggqhkjOPQQDAjB0MQkwBwYDVQQIEwAxCTAHBgNV
BAcTADEJMAcGA1UECRMAMQkwBwYDVQQREwAxDjAMBgNVBAoTBWZsYXRvMQkwBwYD
VQQLEwAxDjAMBgNVBAMTBW5vZGUxMQswCQYDVQQGEwJaSDEOMAwGA1UEKhMFZWNl
cnQwIBcNMjAwNTIxMDU1MzU2WhgPMjEyMDA0MjcwNjUzNTZaMHQxCTAHBgNVBAgT
ADEJMAcGA1UEBxMAMQkwBwYDVQQJEwAxCTAHBgNVBBETADEOMAwGA1UEChMFZmxh
dG8xCTAHBgNVBAsTADEOMAwGA1UEAxMFbm9kZTIxCzAJBgNVBAYTAlpIMQ4wDAYD
VQQqEwVlY2VydDBWMBAGByqGSM49AgEGBSuBBAAKA0IABBI3ewNK21vHNOPG6U3X
mKJohSNNz72QKDxUpRt0fCJHwaGYfSvY4cnqkbliclfckUTpCkFSRr4cqN6PURCF
zkWjeTB3MA4GA1UdDwEB/wQEAwIChDAmBgNVHSUEHzAdBggrBgEFBQcDAgYIKwYB
BQUHAwEGAioDBgOBCwEwDwYDVR0TAQH/BAUwAwEB/zANBgNVHQ4EBgQEAQIDBDAP
BgNVHSMECDAGgAQBAgMEMAwGAypWAQQFZWNlcnQwCgYIKoZIzj0EAwIDQgD6Dzzv
31kMiIlYtQRQjjs5m9pSvMZtmq2vOJW6/5J3kBuUkjTJfawqDdoxVh4yw06F/IBQ
gSu97EZSaseY9hweAA==
-----END CERTIFICATE-----`,
			},
			{
				Hostname: "node2",
				CertContent: `-----BEGIN CERTIFICATE-----
MIICSTCCAfWgAwIBAgIBATAKBggqhkjOPQQDAjB0MQkwBwYDVQQIEwAxCTAHBgNV
BAcTADEJMAcGA1UECRMAMQkwBwYDVQQREwAxDjAMBgNVBAoTBWZsYXRvMQkwBwYD
VQQLEwAxDjAMBgNVBAMTBW5vZGUyMQswCQYDVQQGEwJaSDEOMAwGA1UEKhMFZWNl
cnQwIBcNMjAwNTIxMDU1MTE0WhgPMjEyMDA0MjcwNjUxMTRaMHQxCTAHBgNVBAgT
ADEJMAcGA1UEBxMAMQkwBwYDVQQJEwAxCTAHBgNVBBETADEOMAwGA1UEChMFZmxh
dG8xCTAHBgNVBAsTADEOMAwGA1UEAxMFbm9kZTExCzAJBgNVBAYTAlpIMQ4wDAYD
VQQqEwVlY2VydDBWMBAGByqGSM49AgEGBSuBBAAKA0IABBI3ewNK21vHNOPG6U3X
mKJohSNNz72QKDxUpRt0fCJHwaGYfSvY4cnqkbliclfckUTpCkFSRr4cqN6PURCF
zkWjeTB3MA4GA1UdDwEB/wQEAwIChDAmBgNVHSUEHzAdBggrBgEFBQcDAgYIKwYB
BQUHAwEGAioDBgOBCwEwDwYDVR0TAQH/BAUwAwEB/zANBgNVHQ4EBgQEAQIDBDAP
BgNVHSMECDAGgAQBAgMEMAwGAypWAQQFZWNlcnQwCgYIKoZIzj0EAwIDQgB3Cfo8
/Vdzzlz+MW+MIVuYQkcNkACY/yU/IXD1sHDGZQWcGKr4NR7FHJgsbjGpbUiCofw4
4rK6biAEEAOcv1BQAA==
-----END CERTIFICATE-----`,
			},
			{
				Hostname: "node3",
				CertContent: `-----BEGIN CERTIFICATE-----
MIICRjCCAfKgAwIBAgIBATAKBggqhkjOPQQDAjB0MQkwBwYDVQQIEwAxCTAHBgNV
BAcTADEJMAcGA1UECRMAMQkwBwYDVQQREwAxDjAMBgNVBAoTBWZsYXRvMQkwBwYD
VQQLEwAxDjAMBgNVBAMTBW5vZGUzMQswCQYDVQQGEwJaSDEOMAwGA1UEKhMFZWNl
cnQwIBcNMjAwNTIxMDU1MTQ0WhgPMjEyMDA0MjcwNjUxNDRaMHQxCTAHBgNVBAgT
ADEJMAcGA1UEBxMAMQkwBwYDVQQJEwAxCTAHBgNVBBETADEOMAwGA1UEChMFZmxh
dG8xCTAHBgNVBAsTADEOMAwGA1UEAxMFbm9kZTExCzAJBgNVBAYTAlpIMQ4wDAYD
VQQqEwVlY2VydDBWMBAGByqGSM49AgEGBSuBBAAKA0IABBI3ewNK21vHNOPG6U3X
mKJohSNNz72QKDxUpRt0fCJHwaGYfSvY4cnqkbliclfckUTpCkFSRr4cqN6PURCF
zkWjdjB0MA4GA1UdDwEB/wQEAwIChDAmBgNVHSUEHzAdBggrBgEFBQcDAgYIKwYB
BQUHAwEGAioDBgOBCwEwDAYDVR0TAQH/BAIwADANBgNVHQ4EBgQEAQIDBDAPBgNV
HSMECDAGgAQBAgMEMAwGAypWAQQFZWNlcnQwCgYIKoZIzj0EAwIDQgCalJzkOAqk
IU4AMQGeWzFmdtRYJXZiElyqfrCn7Zg08Ssx14ZMO8K2cRsCncO+c/a/0IqwObEO
wL4C2ich1g5bAA==
-----END CERTIFICATE-----`,
			},
			{
				Hostname: "node4",
				CertContent: `-----BEGIN CERTIFICATE-----
MIICSTCCAfWgAwIBAgIBATAKBggqhkjOPQQDAjB0MQkwBwYDVQQIEwAxCTAHBgNV
BAcTADEJMAcGA1UECRMAMQkwBwYDVQQREwAxDjAMBgNVBAoTBWZsYXRvMQkwBwYD
VQQLEwAxDjAMBgNVBAMTBW5vZGU0MQswCQYDVQQGEwJaSDEOMAwGA1UEKhMFZWNl
cnQwIBcNMjAwNTIxMDU1MzI0WhgPMjEyMDA0MjcwNjUzMjRaMHQxCTAHBgNVBAgT
ADEJMAcGA1UEBxMAMQkwBwYDVQQJEwAxCTAHBgNVBBETADEOMAwGA1UEChMFZmxh
dG8xCTAHBgNVBAsTADEOMAwGA1UEAxMFbm9kZTExCzAJBgNVBAYTAlpIMQ4wDAYD
VQQqEwVlY2VydDBWMBAGByqGSM49AgEGBSuBBAAKA0IABBI3ewNK21vHNOPG6U3X
mKJohSNNz72QKDxUpRt0fCJHwaGYfSvY4cnqkbliclfckUTpCkFSRr4cqN6PURCF
zkWjeTB3MA4GA1UdDwEB/wQEAwIChDAmBgNVHSUEHzAdBggrBgEFBQcDAgYIKwYB
BQUHAwEGAioDBgOBCwEwDwYDVR0TAQH/BAUwAwEB/zANBgNVHQ4EBgQEAQIDBDAP
BgNVHSMECDAGgAQBAgMEMAwGAypWAQQFZWNlcnQwCgYIKoZIzj0EAwIDQgATn17d
VxBI7s/4D/KNU0T3wvQeNrj6ZYmiNRjB/JIwPH7MhM5cTeJx0rs0K27rM/pFZtbq
+0W3ll/JfjQ6WxAAAQ==
-----END CERTIFICATE-----`,
			},
		},
	})
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	payload := bvm.EncodeOperation(operation)
	t.Log(utils.BytesToHex(payload))
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))

	t.Log("-------------------------------")
	operation = bvm.NewHashGetOperation("the_key_for_genesis_info")
	payload = bvm.EncodeOperation(operation)
	tx = NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err = rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))
}

func TestRPC_BVMGet(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	operation := bvm.NewHashGetOperation("0x1231")
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))
}

func TestBVMCert(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	ecert := []byte("-----BEGIN CERTIFICATE-----\nMIICHzCCAcSgAwIBAgIIapt5s0h7G4owCgYIKoEcz1UBg3UwPTELMAkGA1UEBhMC\nQ04xEzARBgNVBAoTCkh5cGVyY2hhaW4xDjAMBgNVBAMTBW5vZGUxMQkwBwYDVQQq\nEwAwIBcNMjEwMzI1MDAwMDAwWhgPMjEyMTAzMjUwMDAwMDBaMEQxCzAJBgNVBAYT\nAkNOMRMwEQYDVQQKEwpIeXBlcmNoYWluMQ4wDAYDVQQDEwVub2RlNTEQMA4GA1UE\nKhMHc2RrY2VydDBZMBMGByqGSM49AgEGCCqBHM9VAYItA0IABG401JscKfKj0rT3\nxN8Dwyen8mVCnXC3GBNkaENJEnqOO4jw0wT331CcX47bHMcSMRfpprbbv4cUj8cV\ncXNa9J6jgaQwgaEwDgYDVR0PAQH/BAQDAgHuMDEGA1UdJQQqMCgGCCsGAQUFBwMC\nBggrBgEFBQcDAQYIKwYBBQUHAwMGCCsGAQUFBwMEMAwGA1UdEwEB/wQCMAAwHQYD\nVR0OBBYEFPT6cvqWN9MBuhhlnmPrCQZG2iKoMB8GA1UdIwQYMBaAFJq1kzm0Q76P\nxf84+ZRlfrWBKy27MA4GAypWAQQHc2RrY2VydDAKBggqgRzPVQGDdQNJADBGAiEA\n3vcQvDi91E5GTsvV/IhKqrfuLkrnudN+3+QtocUX2IMCIQC6Ct1CS4c60SaE59tI\n3a/wjXSyWIYGN6Rwt0k0KFbF+w==\n-----END CERTIFICATE-----\n")
	priv := []byte("-----BEGIN EC PRIVATE KEY-----\nMHgCAQECIQClNEoZsGgZLfdMgYyMCWH8I0PLZynFp2U+wnsSzJ6z+6AKBggqgRzP\nVQGCLaFEA0IABG401JscKfKj0rT3xN8Dwyen8mVCnXC3GBNkaENJEnqOO4jw0wT3\n31CcX47bHMcSMRfpprbbv4cUj8cVcXNa9J4=\n-----END EC PRIVATE KEY-----\n")
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	operation, ferr := bvm.NewCertRevokeOperation(ecert, priv)
	assert.Nil(t, ferr)
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	res := bvm.Decode(re.Ret)
	assert.Equal(t, res.Err, "only support to revoke sdkcert, this cert type is ecert")

	operation = bvm.NewCertCheckOperation(ecert)
	payload = bvm.EncodeOperation(operation)
	tx = NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err = rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))
}

func TestBVMFreezeCert(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	cert := []byte("-----BEGIN CERTIFICATE-----\nMIICFDCCAbqgAwIBAgIIbGmp7HEb95UwCgYIKoEcz1UBg3UwPTELMAkGA1UEBhMC\nQ04xEzARBgNVBAoTCkh5cGVyY2hhaW4xDjAMBgNVBAMTBW5vZGUxMQkwBwYDVQQq\nEwAwHhcNMjEwMzEwMDAwMDAwWhcNMjUwMzEwMDAwMDAwWjA/MQswCQYDVQQGEwJD\nTjEOMAwGA1UEChMFZmxhdG8xDjAMBgNVBAMTBW5vZGUxMRAwDgYDVQQqEwdzZGtj\nZXJ0MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAE1hoClj022lTxWSUCw0Ht4PT+dr8/\nn0BQLeuQVBCnZWKNntBg6cMyVSbMVtcyhAyB8s4+tvzS5bIOqYjLqdO18KOBpDCB\noTAOBgNVHQ8BAf8EBAMCAe4wMQYDVR0lBCowKAYIKwYBBQUHAwIGCCsGAQUFBwMB\nBggrBgEFBQcDAwYIKwYBBQUHAwQwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUEo46\neuyltTBBzeqlUhbr7DhPVvowHwYDVR0jBBgwFoAUmrWTObRDvo/F/zj5lGV+tYEr\nLbswDgYDKlYBBAdzZGtjZXJ0MAoGCCqBHM9VAYN1A0gAMEUCIHnScuepuomkq2OT\nprJL44lxsSkc4Zhpq6c+IpX5cbmZAiEA6l2BMWHuDrVudJ2COYWo8E42mvn7lLPD\nmpMkfrWt5ek=\n-----END CERTIFICATE-----\n")
	priv := []byte("-----BEGIN EC PRIVATE KEY-----\nMHQCAQEEICKWeh1X4x1cZI+nfsAw5VXDgLPspN9vixkTlOTSllknoAcGBSuBBAAK\noUQDQgAE1hoClj022lTxWSUCw0Ht4PT+dr8/n0BQLeuQVBCnZWKNntBg6cMyVSbM\nVtcyhAyB8s4+tvzS5bIOqYjLqdO18A==\n-----END EC PRIVATE KEY-----\n")
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)

	operation, ferr := bvm.NewCertFreezeOperation(cert, priv)
	assert.Nil(t, ferr)
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	res := bvm.Decode(re.Ret)
	assert.True(t, res.Success)
	t.Log(res)

	operation = bvm.NewCertCheckOperation(cert)
	payload = bvm.EncodeOperation(operation)
	tx = NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err = rp.InvokeContract(tx)
	assert.Nil(t, err)
	res = bvm.Decode(re.Ret)
	assert.True(t, res.Success)
	t.Log(res)

	operation, err = bvm.NewCertUnfreezeOperation(cert, priv)
	assert.Nil(t, err)
	payload = bvm.EncodeOperation(operation)
	tx = NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err = rp.InvokeContract(tx)
	assert.Nil(t, err)
	res = bvm.Decode(re.Ret)
	assert.True(t, res.Success)
	t.Log(res)

	operation = bvm.NewCertCheckOperation(cert)
	payload = bvm.EncodeOperation(operation)
	tx = NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err = rp.InvokeContract(tx)
	assert.Nil(t, err)
	res = bvm.Decode(re.Ret)
	assert.False(t, res.Success)
	t.Log(res)
}

func TestBVMCheckCert(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	ecert := []byte("-----BEGIN CERTIFICATE-----\nMIICHTCCAcOgAwIBAgIIed1vBe+JODkwCgYIKoEcz1UBg3UwPTELMAkGA1UEBhMC\nQ04xEzARBgNVBAoTCkh5cGVyY2hhaW4xDjAMBgNVBAMTBW5vZGUxMQkwBwYDVQQq\nEwAwIBcNMjEwMzI1MDAwMDAwWhgPMjEyMTAzMjUwMDAwMDBaMEIxCzAJBgNVBAYT\nAkNOMRMwEQYDVQQKEwpIeXBlcmNoYWluMQ4wDAYDVQQDEwVub2RlNjEOMAwGA1UE\nKhMFZWNlcnQwWTATBgcqhkjOPQIBBggqgRzPVQGCLQNCAAQ6uzqDCLapNh7AR8v2\nxSF1CEe7+ZqpBqQrb6i07L0h1AyC77t1ykE03JPPf2BaGyj+WI2jWK3QtCFiulfr\nYjfvo4GlMIGiMA4GA1UdDwEB/wQEAwIB7jAxBgNVHSUEKjAoBggrBgEFBQcDAgYI\nKwYBBQUHAwEGCCsGAQUFBwMDBggrBgEFBQcDBDAPBgNVHRMBAf8EBTADAQH/MB0G\nA1UdDgQWBBT9ZbjckJMem6i2brgxHkzqIZF+OTAfBgNVHSMEGDAWgBSatZM5tEO+\nj8X/OPmUZX61gSstuzAMBgMqVgEEBWVjZXJ0MAoGCCqBHM9VAYN1A0gAMEUCIDOB\nTuFtkbup8iYH3W5iE4bo4cfV7NshFMtkfsh0O3ISAiEAs8+PYufzvjg7crkmL8rs\nYy80FcF/AV1EluqfFWS2iN4=\n-----END CERTIFICATE-----\n")
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	operation := bvm.NewCertCheckOperation(ecert)
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))
}

func TestBVMCreatePermissionProposal(t *testing.T) {
	//t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	var operations []bvm.PermissionOperation
	operations = append(operations, bvm.NewPermissionCreateRoleOperation(managerRole))
	operations = append(operations, bvm.NewPermissionGrantOperation(managerRole, notAdminAddress))
	operations = append(operations, bvm.NewPermissionRevokeOperation(managerRole, notAdminAddress))
	operations = append(operations, bvm.NewPermissionDeleteRoleOperation(managerRole))
	proposalCreateOperation := bvm.NewProposalCreateOperationForPermission(operations...)
	payload := bvm.EncodeOperation(proposalCreateOperation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(proposalCreateOperation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))
}

func TestRPC_BVMVoteProposal(t *testing.T) {
	//t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	for i := 4; i < 5; i++ {
		key, err := account.NewAccountSm2FromAccountJSON(accountJsons[i], pwd)
		assert.Nil(t, err)
		operation := bvm.NewProposalVoteOperation(1, true)
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.Nil(t, err)
		t.Log(bvm.Decode(re.Ret))
	}
}

func TestRPC_BVMExecuteProposal(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	operation := bvm.NewProposalExecuteOperation(1)
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))
}

func TestRPC_BVMCreateConfigProposal(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))
}

func TestRPC_BVMConfigProposal(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
	payload := bvm.EncodeOperation(operation)
	fmt.Println(payload)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	var proposal bvm.ProposalData
	assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
	t.Log(proposal)

	for i := 1; i < 4; i++ {
		key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
		assert.Nil(t, err)
		operation := bvm.NewProposalVoteOperation(int(proposal.Id), false)
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.Nil(t, err)
		var proposal bvm.ProposalData
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
	}

	key, err = account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	operation = bvm.NewProposalExecuteOperation(int(proposal.Id))
	payload = bvm.EncodeOperation(operation)
	tx = NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err = rp.InvokeContract(tx)
	assert.Nil(t, err)
	assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
	t.Log(proposal)
}

func TestRPC_BVMCancelProposal(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	operation := bvm.NewProposalCancelOperation(2)
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))
}

func TestRPC_BVMCreateNodeProposal(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)
	var operations []bvm.NodeOperation
	operations = append(operations, bvm.NewNodeAddNodeOperation([]byte("pub"), "node1", "vp", "global"))
	operations = append(operations, bvm.NewNodeAddVPOperation("node1", "global"))
	operations = append(operations, bvm.NewNodeRemoveVPOperation("node1", "global"))
	proposalCreateOperation := bvm.NewProposalCreateOperationForNode(operations...)
	payload := bvm.EncodeOperation(proposalCreateOperation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(proposalCreateOperation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))

	// cancel
	operation := bvm.NewProposalCancelOperation(3)
	payload = bvm.EncodeOperation(operation)
	tx = NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err = rp.InvokeContract(tx)
	assert.Nil(t, err)
	t.Log(bvm.Decode(re.Ret))
}

func TestRPC_BVMVote(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("branch1_accept", func(t *testing.T) {
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		for i := 0; i < 6; i++ {
			var (
				key account.Key
				err error
			)
			if i < 4 {
				key, err = account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			} else {
				key, err = account.NewAccountSm2FromAccountJSON(accountJsons[i], pwd)
			}

			assert.Nil(t, err)
			operation := bvm.NewProposalVoteOperation(int(p.ID), true)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.Nil(t, err)
			t.Log(bvm.Decode(re.Ret))
		}
	})

	t.Run("branch2_reject", func(t *testing.T) {
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		for i := 0; i < 6; i++ {
			var (
				key account.Key
				err error
			)
			if i < 4 {
				key, err = account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			} else {
				key, err = account.NewAccountSm2FromAccountJSON(accountJsons[i], pwd)
			}

			assert.Nil(t, err)
			operation := bvm.NewProposalVoteOperation(int(p.ID), false)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.Nil(t, err)
			t.Log(bvm.Decode(re.Ret))
		}
	})

	t.Run("branch3_query", func(t *testing.T) {
		p, _ := rp.GetProposal()
		t.Log(p.Status)
	})
}

func TestRPC_BVMAddNode1(t *testing.T) {
	t.Skip()
	var (
		p *ProposalRaw
	)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("step1_add_node_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		hostname := "node5"
		ns := "global"
		op := bvm.NewNodeAddVPOperation(hostname, ns)
		operations := bvm.NewProposalCreateOperationForNode(op)
		payload := bvm.EncodeOperation(operations)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operations.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.True(t, bvm.Decode(re.Ret).Success)
		p, err = rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
	})

	t.Run("step2_vote", func(t *testing.T) {
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		for i := 0; i < 6; i++ {
			var (
				key account.Key
				err error
			)
			if i < 4 {
				key, err = account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			} else {
				key, err = account.NewAccountSm2FromAccountJSON(accountJsons[i], pwd)
			}

			assert.Nil(t, err)
			operation := bvm.NewProposalVoteOperation(int(p.ID), true)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.Nil(t, err)
			t.Log(bvm.Decode(re.Ret))
		}
	})
	t.Run("step3_execute_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		operation := bvm.NewProposalExecuteOperation(int(p.ID))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.True(t, bvm.Decode(re.Ret).Success)
		p, err = rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
	})

	t.Run("step4_check_vp", func(t *testing.T) {
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
		set, err := rp.GetVSet()
		assert.NoError(t, err)
		t.Log(set)
		m, err := rp.GetHosts("vp")
		assert.NoError(t, err)
		var list []string
		for k := range m {
			list = append(list, k)
		}
		t.Log(list)
	})
}

func TestRPC_BVMDelNode1(t *testing.T) {
	t.Skip()
	var (
		p *ProposalRaw
	)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("step1_del_node_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		hostname := "node5"
		ns := "global"
		op := bvm.NewNodeRemoveVPOperation(hostname, ns)
		operations := bvm.NewProposalCreateOperationForNode(op)
		payload := bvm.EncodeOperation(operations)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operations.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.True(t, bvm.Decode(re.Ret).Success)
		p, err = rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
	})

	t.Run("step2_vote", func(t *testing.T) {
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		for i := 0; i < 6; i++ {
			var (
				key account.Key
				err error
			)
			if i < 4 {
				key, err = account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			} else {
				key, err = account.NewAccountSm2FromAccountJSON(accountJsons[i], pwd)
			}

			assert.Nil(t, err)
			operation := bvm.NewProposalVoteOperation(int(p.ID), true)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.Nil(t, err)
			t.Log(bvm.Decode(re.Ret))
		}
	})
	t.Run("step3_execute_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		operation := bvm.NewProposalExecuteOperation(int(p.ID))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.True(t, bvm.Decode(re.Ret).Success)
		p, err = rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
	})

	t.Run("step4_check_vp", func(t *testing.T) {
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
		set, err := rp.GetVSet()
		assert.NoError(t, err)
		t.Log(set)
		m, err := rp.GetHosts("vp")
		assert.NoError(t, err)
		var list []string
		for k := range m {
			list = append(list, k)
		}
		t.Log(list)
	})
}

func TestRPC_BVMStateChange1(t *testing.T) {
	t.Skip()

	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64619
	// 1. create proposal
	// 2. cancel proposal and query
	// 3. create another proposal
	// 4. cancel another proposal
	var proposal bvm.ProposalData
	cf := `
[consensus]
  algo = "RBFT"
`
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
		t.Log(proposal)
	})

	t.Run("step2_cancel_proposal_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalCancelOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		t.Log(bvm.Decode(re.Ret))
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_CANCEL)], p.Status)
	})

	t.Run("step3_create_another_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		lastProposalId := proposal.Id
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		assert.NotEqual(t, lastProposalId, proposal.Id)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
	})

	t.Run("step4_cancel_another_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalCancelOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		t.Log(bvm.Decode(re.Ret))
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_CANCEL)], p.Status)
	})

}

func TestRPC_BVMStateChange2(t *testing.T) {
	t.Skip()

	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b6461d
	// 1. create proposal
	// 2. vote
	// 3. execute and query
	var proposal bvm.ProposalData
	cf := `
[consensus]
  algo = "RBFT"
`
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
	})

	t.Run("step2_vote", func(t *testing.T) {
		for i := 0; i < 4; i++ {
			key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			var proposal bvm.ProposalData
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
			t.Log(proposal)
		}
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
	})

	t.Run("step3_execute_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalExecuteOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.True(t, bvm.Decode(re.Ret).Success)
		t.Log(proposal)
		conf, err := rp.GetConfig()
		assert.NoError(t, err)
		CheckConfig(t, cf, conf)
		t.Log(conf)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
	})
}

func TestRPC_BVMStateChange3(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64622
	// 1. create proposal
	// 2. vote
	// 3. cancel and query
	var proposal bvm.ProposalData
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)

	})

	t.Run("step2_vote", func(t *testing.T) {
		for i := 1; i < 4; i++ {
			key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			var proposal bvm.ProposalData
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		}
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
	})

	t.Run("step3_cancel_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalCancelOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_CANCEL)], p.Status)

	})
}

func TestRPC_BVMStateChange4(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64623
	// 1. create proposal
	// 2. vote(accept)
	// 3. out of time and query
	var proposal bvm.ProposalData
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
	})

	t.Run("step2_vote", func(t *testing.T) {
		for i := 1; i < 4; i++ {
			key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			var proposal bvm.ProposalData
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		}
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
	})
	time.Sleep(11 * time.Minute)
	t.Run("step3_out_of_time_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalExecuteOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.False(t, bvm.Decode(re.Ret).Success)
		t.Log(proposal)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_TIMEOUT)], p.Status)
	})
}

func TestRPC_BVMStateChange5(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64623
	// 1. create proposal
	// 2. vote(reject)
	// 3. create another proposal
	var proposal bvm.ProposalData
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
	})

	t.Run("step2_vote", func(t *testing.T) {
		for i := 1; i < 4; i++ {
			key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalVoteOperation(int(proposal.Id), false)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			var proposal bvm.ProposalData
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		}
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_REJECT)], p.Status)

	})

	t.Run("step3_create_another_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		lastId := proposal.Id
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
		assert.NotEqual(t, lastId, p.ID)
		t.Log(proposal)

	})
	t.Run("step4_cancel_another_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalCancelOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_CANCEL)], p.Status)
		t.Log(proposal)
	})
}

func TestRPC_BVMStateChange6(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b6462c
	// 1. create proposal
	// 2. overtime
	// 3. create another proposal
	var proposal bvm.ProposalData
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
	})

	t.Run("step2_vote", func(t *testing.T) {
		for i := 1; i < 4; i++ {
			key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			var proposal bvm.ProposalData
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		}
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
	})
	time.Sleep(11 * time.Minute)
	t.Run("step3_create_another_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		lastId := proposal.Id
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
		assert.NotEqual(t, lastId, p.ID)
		t.Log(proposal)
	})
}

func TestRPC_BVMCreateProposal1(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64630
	// 1. admin create proposal
	// 2. admin cancel proposal
	var proposal bvm.ProposalData
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
		t.Log(proposal)
	})

	t.Run("step2_cancel_proposal_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalCancelOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		t.Log(bvm.Decode(re.Ret))
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_CANCEL)], p.Status)
	})
}

func TestRPC_BVMCreateProposal2(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64630
	// 1. non-admin create proposal
	pwd := "12345678"
	t.Run("step1_create_proposal", func(t *testing.T) {
		nonAdminAcc, err := account.NewAccount(pwd)
		assert.NoError(t, err)
		key, err := account.GenKeyFromAccountJson(nonAdminAcc, pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.False(t, bvm.Decode(re.Ret).Success)
	})

}

func TestRPC_BVMExecution1(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b6464b
	// 1. create proposal
	// 2. vote
	// 3. execute and query
	var proposal bvm.ProposalData
	cf := `
[proposal]
	timeout   = "15m"
	threshold = 3
`
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
	})

	t.Run("step2_vote", func(t *testing.T) {
		for i := 1; i < 4; i++ {
			key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			var proposal bvm.ProposalData
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
			t.Log(proposal)
		}
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
	})

	t.Run("step3_execute_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalExecuteOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.True(t, bvm.Decode(re.Ret).Success)
		t.Log(proposal)
		conf, err := rp.GetConfig()
		assert.NoError(t, err)
		CheckConfig(t, cf, conf)
		t.Log(conf)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
	})
}

func TestRPC_BVMExecution2(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64650
	// 1. create proposal
	// 2. vote
	// 3. non-creator execute and query
	// 4. cancel
	var proposal bvm.ProposalData
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
	})

	t.Run("step2_vote", func(t *testing.T) {
		for i := 1; i < 4; i++ {
			key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			var proposal bvm.ProposalData
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
			t.Log(proposal)
		}
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
	})

	t.Run("step3_non-creator_execute_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[1], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalExecuteOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.False(t, bvm.Decode(re.Ret).Success)
		t.Log(proposal)
		conf, err := rp.GetConfig()
		assert.NoError(t, err)
		CheckConfig(t, cf, conf)
		t.Log(conf)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
	})

	t.Run("step4_cancel_proposal_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalCancelOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		t.Log(bvm.Decode(re.Ret))
		p, err := rp.GetProposal()
		assert.Equal(t, proposal.Id, p.ID)
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_CANCEL)], p.Status)
	})
}

func TestRPC_BVMCancel1(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b6465c
	// 1. create proposal
	// 2. cancel and query
	var proposal bvm.ProposalData
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
	})

	t.Run("step2_cancel_proposal_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalCancelOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		t.Log(bvm.Decode(re.Ret))
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_CANCEL)], p.Status)
	})
}

func TestRPC_BVMCancel2(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b6465c
	// 1. create proposal
	// 2. vote
	// 3. cancel and query
	var proposal bvm.ProposalData
	t.Run("step1_create_proposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(cf))
		payload := bvm.EncodeOperation(operation)
		fmt.Println(payload)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		t.Log(proposal)
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_VOTING)], p.Status)
	})

	t.Run("step2_vote", func(t *testing.T) {
		for i := 1; i < 4; i++ {
			key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			var proposal bvm.ProposalData
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
			t.Log(proposal)
		}
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
	})

	t.Run("step3_cancel_proposal_and_query", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalCancelOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		t.Log(bvm.Decode(re.Ret))
		p, err := rp.GetProposal()
		assert.NoError(t, err)
		assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_CANCEL)], p.Status)
	})
}

func TestRPC_BVMSubscribe(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("step1_subscribe", func(t *testing.T) {
		wsCli := rp.GetWebSocketClient()
		subID, err := wsCli.SubscribeForProposal(1, &TestEventHandler{})
		if err != nil {
			t.Error(err.String())
			return
		}

		time.Sleep(30 * time.Minute)
		//解订阅
		_ = wsCli.UnSubscribe(subID)
		time.Sleep(1 * time.Second)
		//关闭连接
		_ = wsCli.CloseConn(1)
		time.Sleep(1 * time.Second)
	})
}

func TestRPC_BVMAuth1(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64674
	// 1. create admin role
	t.Run("step1_create_admin", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		createRole := bvm.NewPermissionCreateRoleOperation("admin")
		operation := bvm.NewProposalCreateOperationForPermission(createRole)
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.False(t, bvm.Decode(re.Ret).Success)
	})

}

func TestRPC_BVMAuth2(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64679
	// 1. delete admin role
	t.Run("step1_create_admin", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.NoError(t, err)
		createRole := bvm.NewPermissionDeleteRoleOperation("admin")
		operation := bvm.NewProposalCreateOperationForPermission(createRole)
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		assert.False(t, bvm.Decode(re.Ret).Success)
	})

}

func TestRPC_BVMAuth3_4_5(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64674
	// 1. create other role
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b6467f
	// 2. set role to address
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b64684
	// 3. rm role from address
	var proposal bvm.ProposalData
	role := "auth3"
	t.Run("step1", func(t *testing.T) {
		t.Run("step1_1_create_role_proposal", func(t *testing.T) {
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			createRole := bvm.NewPermissionCreateRoleOperation(role)
			operation := bvm.NewProposalCreateOperationForPermission(createRole)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		})

		t.Run("step1_2_vote", func(t *testing.T) {
			for i := 1; i < 4; i++ {
				key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
				assert.NoError(t, err)
				operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
				payload := bvm.EncodeOperation(operation)
				tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
				tx.Sign(key)
				re, err := rp.InvokeContract(tx)
				assert.NoError(t, err)
				var proposal bvm.ProposalData
				assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
				t.Log(proposal)
			}
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
		})

		t.Run("step1_3_execute_and_query", func(t *testing.T) {
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalExecuteOperation(int(proposal.Id))
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.True(t, bvm.Decode(re.Ret).Success)
			t.Log(proposal)
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
			roles, err := rp.GetAllRoles()
			assert.NoError(t, err)
			_, exist := roles[role]
			assert.True(t, exist)

			addrs, err := rp.GetAccountsByRole(role)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(addrs))

		})
	})

	pwd := "12345678"
	nonAdminAcc, _ := account.NewAccount(pwd)
	key, _ := account.GenKeyFromAccountJson(nonAdminAcc, pwd)
	addr := key.(account.Key).GetAddress().Hex()
	t.Run("step2", func(t *testing.T) {
		t.Run("step2_1_set_role_proposal", func(t *testing.T) {
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			createRole := bvm.NewPermissionGrantOperation(role, addr)
			operation := bvm.NewProposalCreateOperationForPermission(createRole)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		})

		t.Run("step2_2_vote", func(t *testing.T) {
			for i := 1; i < 4; i++ {
				key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
				assert.NoError(t, err)
				operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
				payload := bvm.EncodeOperation(operation)
				tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
				tx.Sign(key)
				re, err := rp.InvokeContract(tx)
				assert.NoError(t, err)
				var proposal bvm.ProposalData
				assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
				t.Log(proposal)
			}
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
		})

		t.Run("step2_3_execute_and_query", func(t *testing.T) {
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalExecuteOperation(int(proposal.Id))
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.True(t, bvm.Decode(re.Ret).Success)
			t.Log(proposal)
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
			roles, err := rp.GetAllRoles()
			assert.NoError(t, err)
			_, exist := roles[role]
			assert.True(t, exist)

			addrs, err := rp.GetAccountsByRole(role)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(addrs))
			assert.Equal(t, addr, addrs[0])

			roleList, err := rp.GetRoles(addr)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(roleList))
			assert.Equal(t, role, roleList[0])

		})

	})

	t.Run("step3", func(t *testing.T) {
		t.Run("step3_1_revoke_role_proposal", func(t *testing.T) {
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			createRole := bvm.NewPermissionRevokeOperation(role, addr)
			operation := bvm.NewProposalCreateOperationForPermission(createRole)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		})

		t.Run("step3_2_vote", func(t *testing.T) {
			for i := 1; i < 4; i++ {
				key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
				assert.NoError(t, err)
				operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
				payload := bvm.EncodeOperation(operation)
				tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
				tx.Sign(key)
				re, err := rp.InvokeContract(tx)
				assert.NoError(t, err)
				var proposal bvm.ProposalData
				assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
				t.Log(proposal)
			}
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
		})

		t.Run("step3_3_execute_and_query", func(t *testing.T) {
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalExecuteOperation(int(proposal.Id))
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.True(t, bvm.Decode(re.Ret).Success)
			t.Log(proposal)
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
			roles, err := rp.GetAllRoles()
			assert.NoError(t, err)
			_, exist := roles[role]
			assert.True(t, exist)

			addrs, err := rp.GetAccountsByRole(role)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(addrs))

			roleList, err := rp.GetRoles(addr)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(roleList))

		})
	})

}

func TestRPC_BVMAuth_did(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	role_admin := "allowRole"
	role_user := "user"
	didKey_admin := registerDIDAccount(rp, nil)
	didKey_user := registerDIDAccount(rp, nil)
	var proposal bvm.ProposalData

	t.Run("set tx rules", func(t *testing.T) {
		//t.Skip()
		fmt.Println(time.Now())
		var ops []bvm.ConfigOperation
		ops = append(ops, bvm.NewSetFilterEnable(true))
		rule := &bvm.NsFilterRule{
			Name:            "did_tx_auth",
			To:              []string{"*"},
			VM:              []string{"hvm"},
			AllowAnyone:     false,
			AuthorizedRoles: []string{role_admin},
			ForbiddenRoles:  []string{role_user},
		}
		ops = append(ops, bvm.NewSetFilterRules([]*bvm.NsFilterRule{rule}))
		operation := bvm.NewProposalCreateOperationByConfigOps(ops...)
		key, err := account.GenKeyFromAccountJson(accountJsons[0], pwd)
		assert.Nil(t, err)
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		re, err := rp.SignAndInvokeContract(tx, key)
		assert.NoError(t, err)
		assert.True(t, bvm.Decode(re.Ret).Success)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		testVoteAndExecute(t)
	})

	t.Run("grant roles", func(t *testing.T) {
		time.Sleep(5 * time.Second) // 预留足够时间，使得节点共识恢复
		fmt.Println(time.Now())
		key, err := account.GenKeyFromAccountJson(accountJsons[0], pwd)
		assert.NoError(t, err)
		createRole := bvm.NewPermissionCreateRoleOperation(role_admin)
		createRole2 := bvm.NewPermissionCreateRoleOperation(role_user)
		grantRole := bvm.NewPermissionGrantOperation(role_admin, utils.ToHex([]byte(didKey_admin.GetAddress())))
		grantRole2 := bvm.NewPermissionGrantOperation(role_user, utils.ToHex([]byte(didKey_user.GetAddress())))
		operation := bvm.NewProposalCreateOperationForPermission(createRole, createRole2, grantRole, grantRole2)
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		re, err := rp.SignAndInvokeContract(tx, key)
		assert.NoError(t, err)
		assert.True(t, bvm.Decode(re.Ret).Success)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		testVoteAndExecute(t)
	})

	t.Run("test_tx_auth", func(t *testing.T) {
		time.Sleep(5 * time.Second)
		fmt.Println(time.Now())
		res, err1 := rp.GetRoles(utils.ToHex([]byte(didKey_admin.GetAddress())))
		fmt.Println(res)
		fmt.Println(err1)
		deployJar, err := DecompressFromJar("../hvmtestfile/fibonacci/fibonacci-1.0-fibonacci.jar")
		if err != nil {
			t.Error(err)
		}

		transaction := NewTransaction(didKey_admin.GetAddress()).Deploy(common.Bytes2Hex(deployJar)).VMType(HVM)
		_, err = rp.SignAndDeployContract(transaction, didKey_admin)
		assert.Nil(t, err)

		transaction = NewTransaction(didKey_user.GetAddress()).Deploy(common.Bytes2Hex(deployJar)).VMType(HVM)
		_, err = rp.SignAndDeployContract(transaction, didKey_user)
		assert.NotNil(t, err)
		fmt.Println(err)
	})
	t.Run("revoke & delete roles", func(t *testing.T) {
		time.Sleep(5 * time.Second)
		fmt.Println(time.Now())
		key, err := account.GenKeyFromAccountJson(accountJsons[0], pwd)
		assert.NoError(t, err)
		revokeRole := bvm.NewPermissionRevokeOperation(role_admin, utils.ToHex([]byte(didKey_admin.GetAddress())))
		revokeRole2 := bvm.NewPermissionRevokeOperation(role_user, utils.ToHex([]byte(didKey_user.GetAddress())))
		deleteRole := bvm.NewPermissionDeleteRoleOperation(role_admin)
		deleteRole2 := bvm.NewPermissionDeleteRoleOperation(role_user)
		operation := bvm.NewProposalCreateOperationForPermission(revokeRole, revokeRole2, deleteRole, deleteRole2)
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		re, err := rp.SignAndInvokeContract(tx, key)
		assert.NoError(t, err)
		assert.True(t, bvm.Decode(re.Ret).Success)
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		testVoteAndExecute(t)
	})

}

func testVoteAndExecute(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	proposal, err := rp.GetProposal()
	assert.Nil(t, err)
	for i := 1; i < 6; i++ {
		key, err := account.GenKeyFromAccountJson(accountJsons[i], pwd)
		assert.NoError(t, err)
		operation := bvm.NewProposalVoteOperation(int(proposal.ID), true)
		payload := bvm.EncodeOperation(operation)
		tx := NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		var proposal bvm.ProposalData
		assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
	}
	p, err := rp.GetProposal()
	assert.NoError(t, err)
	assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)

	key, err1 := account.GenKeyFromAccountJson(accountJsons[0], pwd)
	assert.NoError(t, err1)
	operation := bvm.NewProposalExecuteOperation(int(proposal.ID))
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.(account.Key).GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	re, err := rp.SignAndInvokeContract(tx, key)
	assert.NoError(t, err)
	assert.True(t, bvm.Decode(re.Ret).Success)
	p, err = rp.GetProposal()
	assert.NoError(t, err)
	assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
}

func TestRPC_BVMFilter1_2(t *testing.T) {
	t.Skip()

	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b6468a
	// 1. create role and set role to address
	// 2. set filter rule and try to send transaction
	// http://teambition.hyperchain.cn:8099/project/5cf48431010dd5597861529c/testplan/5e0b10cec8ead008b6b540fd/testcase/5e0d929fc8ead008b6b6468d
	// 3. set filter rule and try to send transaction

	role := "filter1"
	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	var proposal bvm.ProposalData
	pwd := "12345678"
	//nonAdminAcc := `{"address":"0xd93f64dc5d887b4f594d62e60a5b29330c04b8a4","algo":"0x02","version":"1.0","publicKey":"04b6b579ac348009b890c35eddd704bce7175026d56cdf5c2bb42377ba4076d7b126cdda5d81914762b5a704659b39d285ca110d83d462799deabd56f97be4eae8","privateKey":"9dafa6ba3d7624bec05080cc033f8500691d929655c58695eaa37eb0d40bfb624cbe236a80ce82b9"}`
	nonAdminAcc, _ := account.NewAccount(pwd)
	genKey, err := account.GenKeyFromAccountJson(nonAdminAcc, pwd)
	assert.NoError(t, err)
	genAddr := genKey.(account.Key).GetAddress().Hex()

	t.Run("step1", func(t *testing.T) {
		t.Run("step1_1_create_role_and_set_to_addr_proposal", func(t *testing.T) {
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			createRole := bvm.NewPermissionCreateRoleOperation(role)
			setRole := bvm.NewPermissionGrantOperation(role, genAddr)
			operation := bvm.NewProposalCreateOperationForPermission(createRole, setRole)
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		})

		t.Run("step1_2_vote", func(t *testing.T) {
			for i := 1; i < 4; i++ {
				key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
				assert.NoError(t, err)
				operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
				payload := bvm.EncodeOperation(operation)
				tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
				tx.Sign(key)
				re, err := rp.InvokeContract(tx)
				assert.NoError(t, err)
				var proposal bvm.ProposalData
				assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
				t.Log(proposal)
			}
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
		})

		t.Run("step1_3_execute_and_query", func(t *testing.T) {
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalExecuteOperation(int(proposal.Id))
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.True(t, bvm.Decode(re.Ret).Success)
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
			roles, err := rp.GetAllRoles()
			assert.NoError(t, err)
			_, exist := roles[role]
			assert.True(t, exist)

			addrs, err := rp.GetAccountsByRole(role)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(addrs))

		})
	})

	t.Run("step2", func(t *testing.T) {
		t.Run("step2_1_set_filter_proposal", func(t *testing.T) {
			config1 := `
[filter]
    enable = true
    [[filter.rules]]
    allow_anyone = false
    authorized_roles = ["admin"]
    forbidden_roles = ["20"]
    id = 0
    name = "bvm auth"
    to = ["0x0000000000000000000000000000000000ffff02"]
    vm = ["bvm"]

    [[filter.rules]]
    allow_anyone = false
    authorized_roles = ["admin", "filter1"]
    forbidden_roles = []
    id = 1
    name = "evm auth"
    to = ["*"]
    vm = ["evm"]
`
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(config1))
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		})

		t.Run("step2_2_vote", func(t *testing.T) {
			for i := 1; i < 4; i++ {
				key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
				assert.NoError(t, err)
				operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
				payload := bvm.EncodeOperation(operation)
				tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
				tx.Sign(key)
				re, err := rp.InvokeContract(tx)
				assert.NoError(t, err)
				var proposal bvm.ProposalData
				assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
				t.Log(proposal)
			}
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
		})

		t.Run("step2_3_execute_and_check_result", func(t *testing.T) {
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalExecuteOperation(int(proposal.Id))
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.True(t, bvm.Decode(re.Ret).Success)
			t.Log(proposal)
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
			roles, err := rp.GetAllRoles()
			assert.NoError(t, err)
			_, exist := roles[role]
			assert.True(t, exist)

			addrs, err := rp.GetAccountsByRole(role)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(addrs))
			assert.Equal(t, genAddr, addrs[0])

			roleList, err := rp.GetRoles(genAddr)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(roleList))
			assert.Equal(t, role, roleList[0])

		})

		t.Run("step2_4_check_filter_rule", func(t *testing.T) {
			// do not care payload is right or not
			operation := bvm.NewProposalVoteOperation(1, false)
			payload := bvm.EncodeOperation(operation)

			tx := NewTransaction(genAddr).Invoke(genAddr, payload).VMType(EVM)
			tx.Sign(genKey)
			txHash, err := rp.InvokeContractReturnHash(tx)
			assert.NoError(t, err)
			assert.NotEqual(t, "", txHash)

			acc, _ := account.NewAccount(pwd)
			key, _ := account.GenKeyFromAccountJson(acc, pwd)
			addr := key.(account.Key).GetAddress().Hex()
			tx = NewTransaction(addr).Invoke(addr, payload).VMType(EVM)
			tx.Sign(key)
			_, err = rp.InvokeContractReturnHash(tx)
			assert.Error(t, err)
		})
	})

	t.Run("step3", func(t *testing.T) {
		t.Run("step3_1_set_filter_proposal", func(t *testing.T) {
			config2 := `
[filter]
    enable = true
    [[filter.rules]]
    allow_anyone = false
    authorized_roles = ["admin"]
    forbidden_roles = ["20"]
    id = 0
    name = "bvm auth"
    to = ["0x0000000000000000000000000000000000ffff02"]
    vm = ["bvm"]

    [[filter.rules]]
    allow_anyone = true
    authorized_roles = ["admin"]
    forbidden_roles = ["filter1"]
    id = 1
    name = "evm auth"
    to = ["*"]
    vm = ["evm"]
`
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			operation, _ := bvm.NewProposalCreateOperationForConfig([]byte(config2))
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
		})

		t.Run("step3_2_vote", func(t *testing.T) {
			for i := 1; i < 4; i++ {
				key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
				assert.NoError(t, err)
				operation := bvm.NewProposalVoteOperation(int(proposal.Id), true)
				payload := bvm.EncodeOperation(operation)
				tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
				tx.Sign(key)
				re, err := rp.InvokeContract(tx)
				assert.NoError(t, err)
				var proposal bvm.ProposalData
				assert.NoError(t, proto.Unmarshal([]byte(bvm.Decode(re.Ret).Ret), &proposal))
				t.Log(proposal)
			}
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_WAITING_EXE)], p.Status)
		})

		t.Run("step3_3_execute_and_check_result", func(t *testing.T) {
			key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
			assert.NoError(t, err)
			operation := bvm.NewProposalExecuteOperation(int(proposal.Id))
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			assert.True(t, bvm.Decode(re.Ret).Success)
			t.Log(proposal)
			p, err := rp.GetProposal()
			assert.NoError(t, err)
			assert.Equal(t, bvm.ProposalData_Status_name[int32(bvm.ProposalData_COMPLETED)], p.Status)
			roles, err := rp.GetAllRoles()
			assert.NoError(t, err)
			_, exist := roles[role]
			assert.True(t, exist)

			addrs, err := rp.GetAccountsByRole(role)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(addrs))
			assert.Equal(t, genAddr, addrs[0])

			roleList, err := rp.GetRoles(genAddr)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(roleList))
			assert.Equal(t, role, roleList[0])
		})

		t.Run("step3_4_check_filter_rule", func(t *testing.T) {
			tx := NewTransaction(genAddr).Invoke(genAddr, []byte("12345678")).VMType(EVM)
			tx.Sign(genKey)
			_, err := rp.InvokeContractReturnHash(tx)
			assert.Error(t, err)

			acc, _ := account.NewAccount(pwd)
			key, _ := account.GenKeyFromAccountJson(acc, pwd)
			addr := key.(account.Key).GetAddress().Hex()
			tx = NewTransaction(addr).Invoke(addr, []byte("12345678")).VMType(EVM)
			tx.Sign(key)
			txHash, err := rp.InvokeContractReturnHash(tx)
			assert.NoError(t, err)
			assert.NotEqual(t, "", txHash)
		})
	})
}

func CheckConfig(t *testing.T, expected, got string) {
	v1 := viper.New()
	v2 := viper.New()
	assert.NoError(t, v1.ReadConfig(strings.NewReader(expected)))
	assert.NoError(t, v2.ReadConfig(strings.NewReader(got)))
	keys1 := v1.AllKeys()
	for _, k := range keys1 {
		assert.Equal(t, v1.Get(k), v2.Get(k))
	}
}

func TestRPC_SetCName(t *testing.T) {
	t.Skip()
	var proposal bvm.ProposalData
	addr := "0x0000000000000000000000000000000000ffff01"
	name := "HashContract"
	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	t.Run("SetHash", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.Nil(t, err)
		opt := bvm.NewHashSetOperation("0x123", "0x456")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		result := bvm.Decode(re.Ret)
		assert.True(t, result.Success)
	})

	t.Run("CreateProposalForSetCName", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.Nil(t, err)
		setCNameOpt := bvm.NewCNSSetCNameOperation(addr, name)
		cnsOpt := bvm.NewProposalCreateOperationForCNS(setCNameOpt)
		payload := bvm.EncodeOperation(cnsOpt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(cnsOpt.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		result := bvm.Decode(re.Ret)
		assert.True(t, result.Success)

		_ = proto.Unmarshal([]byte(result.Ret), &proposal)
	})

	t.Run("VoteProposal", func(t *testing.T) {
		for i := 1; i < 4; i++ {
			key, err := account.NewAccountFromAccountJSON(accountJsons[i], pwd)
			assert.Nil(t, err)
			opt := bvm.NewProposalVoteOperation(int(proposal.Id), true)
			payload := bvm.EncodeOperation(opt)
			tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
			tx.Sign(key)
			re, err := rp.InvokeContract(tx)
			assert.NoError(t, err)
			result := bvm.Decode(re.Ret)
			assert.True(t, result.Success)
		}
	})

	t.Run("ExecuteProposal", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.Nil(t, err)
		opt := bvm.NewProposalExecuteOperation(int(proposal.Id))
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		result := bvm.Decode(re.Ret)
		assert.True(t, result.Success)
		t.Log(result)
	})

	t.Run("GetProposal", func(t *testing.T) {
		proposal, _ := rp.GetProposal()
		t.Log(proposal)
		assert.Equal(t, bvm.ProposalData_COMPLETED.String(), proposal.Status)
	})

	t.Run("GetAddressByName", func(t *testing.T) {
		addressByName, stdError := rp.GetAddressByName(name)
		assert.Nil(t, stdError)
		assert.Equal(t, addr, addressByName)
	})

	t.Run("GetNameByAddress", func(t *testing.T) {
		nameByAddress, stdError := rp.GetNameByAddress(addr)
		assert.Nil(t, stdError)
		assert.Equal(t, name, nameByAddress)
	})

	t.Run("InvokeContractByName", func(t *testing.T) {
		key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
		assert.Nil(t, err)
		opt := bvm.NewHashGetOperation("0x123")
		payload := bvm.EncodeOperation(opt)
		tx := NewTransaction(key.GetAddress().Hex()).InvokeByName(name, payload).VMType(BVM)
		tx.Sign(key)
		re, err := rp.InvokeContract(tx)
		assert.NoError(t, err)
		result := bvm.Decode(re.Ret)
		assert.True(t, result.Success)
		t.Log(result)
	})

	t.Run("getTransactionByBikNumAndIndex", func(t *testing.T) {
		info, _ := rp.GetTxByBlkNumAndIdx(1, 0)
		marshal, _ := json.Marshal(info)
		t.Log(string(marshal))
		info, _ = rp.GetTxByBlkNumAndIdx(7, 0)
		marshal, _ = json.Marshal(info)
		t.Log(string(marshal))
	})

	t.Run("GetStatus", func(t *testing.T) {
		status, stdError := rp.GetContractStatus(addr)
		assert.Nil(t, stdError)
		status1, stdError := rp.GetContractStatusByName(name)
		assert.Nil(t, stdError)
		t.Log(status1)
		assert.Equal(t, status, status1)
	})

	t.Run("GetCreator", func(t *testing.T) {
		creator, stdError := rp.GetCreator(addr)
		assert.Nil(t, stdError)
		creatorByName, stdError := rp.GetCreatorByName(name)
		assert.Nil(t, stdError)
		t.Log(creatorByName)
		assert.Equal(t, creator, creatorByName)
	})

	t.Run("GetCreateTime", func(t *testing.T) {
		createTime, stdError := rp.GetCreateTime(addr)
		assert.Nil(t, stdError)
		createTimeByName, stdError := rp.GetCreateTimeByName(name)
		assert.Nil(t, stdError)
		t.Log(createTimeByName)
		assert.Equal(t, createTime, createTimeByName)
	})
}

func TestHash_GetHashSupportAlgo(t *testing.T) {
	t.Skip()
	rp := NewRPC()
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)

	opt := bvm.NewHashGetSupportOperation()
	payload := bvm.EncodeOperation(opt)
	fmt.Println("payload :", hex.EncodeToString(payload))
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	res, serr := rp.SignAndInvokeContract(tx, key)
	assert.Nil(t, serr)
	ret := bvm.Decode(res.Ret)
	t.Log(string(ret.Ret))
}

func TestHash_GetHashAlgo(t *testing.T) {
	//t.Skip()
	rp := NewRPC()
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)

	opt := bvm.NewHashGetAlgoOperation()
	payload := bvm.EncodeOperation(opt)
	fmt.Println("payload :", hex.EncodeToString(payload))
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	res, serr := rp.SignAndInvokeContract(tx, key)
	assert.Nil(t, serr)
	ret := bvm.Decode(res.Ret)
	t.Log(string(ret.Ret))
}

func TestHash_ChangeHashAlgo(t *testing.T) {
	//t.Skip()
	rp := NewRPC()
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)

	ha := bvm.AlgoSet{
		HashAlgo:    "SELF_DEFINED_HASH",
		EncryptAlgo: "SELF_DEFINED_CRYPTO",
	}
	data, _ := json.Marshal(ha)
	opt := bvm.NewHashChangeHashAlgo(data)
	payload := bvm.EncodeOperation(opt)
	fmt.Println("payload :", hex.EncodeToString(payload))
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	res, serr := rp.SignAndInvokeContract(tx, key)
	assert.Nil(t, serr)
	ret := bvm.Decode(res.Ret)
	assert.True(t, ret.Success)
	t.Log(ret.Err)
}

func TestConfig_SetConsensusAlgo(t *testing.T) {
	t.Skip()
	rp := NewRPC()
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)

	// invalid consensus algo
	opt := bvm.NewProposalCreateOperationByConfigOps(bvm.NewSetConsensusAlgo("HOTSTUFF"))
	payload := bvm.EncodeOperation(opt)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	res, serr := rp.SignAndInvokeContract(tx, key)
	assert.Nil(t, serr)
	ret := bvm.Decode(res.Ret)
	assert.False(t, ret.Success)
	t.Log(ret.Err)
	t.Log(string(ret.Ret))

	// support noxbft
	opResults := completeProposal(t, rp, key, 6, bvm.NewProposalCreateOperationByConfigOps(bvm.NewSetConsensusAlgo("NoxBFT")))
	for _, rest := range opResults {
		assert.Equal(t, 200, rest.Code)
	}

	// cannot set again
	opt = bvm.NewProposalCreateOperationByConfigOps(bvm.NewSetConsensusAlgo("SOLO"))
	payload = bvm.EncodeOperation(opt)
	tx = NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	res, serr = rp.SignAndInvokeContract(tx, key)
	assert.Nil(t, serr)
	ret = bvm.Decode(res.Ret)
	assert.False(t, ret.Success)
	t.Log(ret.Err)
	t.Log(string(ret.Ret))

}

func TestConfig_DisableFeature(t *testing.T) {
	t.Skip()
	rp := NewRPC()
	key, err := account.NewAccountFromAccountJSON(accountJsons[0], pwd)
	assert.Nil(t, err)

	// disable set
	opt := bvm.NewProposalCreateOperationByConfigOps(bvm.NewSetConsensusPoolSize(10), bvm.NewSetConsensusBatchSize(10), bvm.NewSetConsensusSetSize(10))
	payload := bvm.EncodeOperation(opt)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	res, serr := rp.SignAndInvokeContract(tx, key)
	assert.Nil(t, serr)
	ret := bvm.Decode(res.Ret)
	assert.False(t, ret.Success)
	t.Log(ret.Err)
	t.Log(string(ret.Ret))
}
