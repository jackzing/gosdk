package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/go-hpc-common/types"
	"github.com/hyperchain/gosdk/account"
	"github.com/hyperchain/gosdk/bvm"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
	"time"
)

var accounts = []string{
	`{"address":"0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd","version":"4.0", "algo":"0x03","publicKey":"0400ddbadb932a0d276e257c6df50599a425804a3743f40942d031f806bf14ab0c57aed6977b1ad14646672f9b9ce385f2c98c4581267b611f48f4b7937de386ac","privateKey":"16acbf6b4f09a476a35ebd4c01e337238b5dceceb6ff55ff0c4bd83c4f91e11b"}`,
	`{"address":"0x856e2b9a5fa82fd1b031d1ff6863864dbac7995d","version":"4.0","algo":"0x13","publicKey":"047ea464762c333762d3be8a04536b22955d97231062442f81a3cff46cb009bbdbb0f30e61ade5705254d4e4e0c0745fb3ba69006d4b377f82ecec05ed094dbe87","privateKey":"71b9acc4ee2b32b3d2c79b5abe9e118e5f73765aee5e7755d6aa31f12945036d"}`,
	`{"address":"0x6201cb0448964ac597faf6fdf1f472edf2a22b89","version":"4.0", "algo":"0x03","publicKey":"04e482f140d70a1b8ec8185cc699db5b391ea5a7b8e93e274b9f706be9efdaec69542eb32a61421ba6219230b9cf87bf849fa01c1d10a8d298cbe3dcfa5954134c","privateKey":"21ff03a654c939f0c9b83e969aaa9050484aa4108028094ee2e927ba7e7d1bbb"}`,
	`{"address":"0xb18c8575e3284e79b92100025a31378feb8100d6","version":"4.0", "algo":"0x03","publicKey":"042169a7260acaff308228579aab2a2c6b3a790922c6a6b58b218cdd7ce0b1db0fbfa6f68737a452010b9d138187b8321288cae98f07fc758bb67bb818292cab9b","privateKey":"aa9c83316f68c17bcc21cf20a4733ae2b2bf76ad1c745f634c0ebf7d5094500e"}`,
	`{"address":"0xe93b92f1da08f925bdee44e91e7768380ae83307","version":"4.0","algo":"0x03","publicKey":"047196daf5d4d1fe339da58e2fe0543bbfec9a464b76546f180facdcc56315b8eeeca50474100f15fb17606695ce24a1f8e5a990600c1c4ea9787ba4dd65c8ce3e","privateKey":"8cdfbe86deb690e331453a84a98c956f0422dd1e783c3a02aed9180b1f4516a9"}`,
	`{"address":"fbca6a7e9e29728773b270d3f00153c75d04e1ad","version":"4.0","algo":"0x13","publicKey":"049c330d0aea3d9c73063db339b4a1a84d1c3197980d1fb9585347ceeb40a5d262166ee1e1cb0c29fd9b2ef0e4f7a7dfb1be6c5e759bf411c520a616863ee046a4","privateKey":"5f0a3ea6c1d3eb7733c3170f2271c10c1206bc49b6b2c7e550c9947cb8f098e3"}`,
}

// AccJSON account json
type AccJSON struct {
	Address    string `json:"address"`
	Algo       string `json:"algo"`
	Version    string `json:"version"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func (aj *AccJSON) isSM() bool {
	return strings.HasPrefix(aj.Algo, "0x1")
}

func TestSendTx(t *testing.T) {
	// 创建交易
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	guomiKey, _ := gm.GenerateSM2Key()
	pubKey := &account.SM2Key{SM2PrivateKey: guomiKey}
	newAddress := pubKey.GetAddress()

	var wg sync.WaitGroup
	for i := 0; i <= 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tx := NewTransaction(newAddress.Hex())
			tx.SetTo("0xcf8dc52bab9775e3df68d7e2f82f52a382bf7706")
			if i/2 == 0 {
				tx.setTxVersion("3.0")
			} else {
				tx.setTxVersion("2.0")
			}

			// 签名
			receipt, err := rp.SignAndSendTx(tx, pubKey)
			if err != nil {
				panic(err)
			}
			if !receipt.Valid {
				panic("invalid tx")
			}
		}(i)
	}
	wg.Wait()
}

func TestGetTxVersionAndSendTx(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	// 创建交易
	guomiKey, _ := gm.GenerateSM2Key()
	pubKey := &account.SM2Key{SM2PrivateKey: guomiKey}
	newAddress := pubKey.GetAddress()
	tx := NewTransaction(newAddress.Hex())
	tx.SetTo("0xcf8dc52bab9775e3df68d7e2f82f52a382bf7706")

	txVersion, gerr := rp.GetTxVersionByID(2)
	if gerr != nil {
		panic(gerr)
	}
	tx.setTxVersion(txVersion)

	// 签名
	tx.Sign(pubKey)

	receipt, err := rp.SendTx(tx)
	if err != nil {
		panic(err)
	}
	if !receipt.Valid {
		panic("invalid tx")
	}
}

func TestRPC_GetVersions(t *testing.T) {
	ast := assert.New(t)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	vr, err := rp.GetVersions()
	ast.Nil(err)
	ast.NotNil(vr)
	t.Logf("%#v", vr.RunningHyperchainVersion)
	t.Logf("%#v", vr.AvailableHyperchainVersion)
}

func TestRPC_SetSupportedVersionByID(t *testing.T) {
	ast := assert.New(t)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	for i := 1; i <= 4; i++ {
		receipt, err := rp.SetSupportedVersionByID(i)
		ast.Nil(err)
		ast.NotNil(receipt)
		time.Sleep(1 * time.Second)
	}
}

func TestRPC_GetLatestVersion_BVM(t *testing.T) {
	key, aerr := account.NewAccountFromAccountJSON(accounts[0], "")
	if aerr != nil {
		t.Logf("new account err: %v", aerr)
	}
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	operation := bvm.NewGetLatestVersionsOperation()
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	result, err := rp.SignAndInvokeContract(tx, key)
	if err != nil {
		t.Fatal(err)
	}
	bvmResult := bvm.Decode(result.Ret)
	var vr struct {
		MaxSupportedVersions map[string]map[types.VersionTag]string `json:"maxSupportedVersions"`
		AvailableVersion     types.AvailableVersion                 `json:"availableVersions"`
		RunningVersion       types.RunningVersion                   `json:"runningVersions"`
	}
	if uerr := json.Unmarshal(bvmResult.Ret, &vr); uerr != nil {
		t.Fatal(uerr)
	}
	t.Log(string(bvmResult.Ret))
	t.Log(vr)
}

func TestRPC_GetSupportedVersionByHostname_BVM(t *testing.T) {
	key, aerr := account.NewAccountFromAccountJSON(accounts[0], "")
	if aerr != nil {
		t.Logf("new account err: %v", aerr)
	}
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	operation := bvm.NewGetSupportedVersionByHostnameOperation("node1")
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	result, err := rp.SignAndInvokeContract(tx, key)
	if err != nil {
		t.Fatal(err)
	}
	bvmResult := bvm.Decode(result.Ret)
	var supportedVersion types.SupportedVersion
	if uerr := json.Unmarshal(bvmResult.Ret, &supportedVersion); uerr != nil {
		t.Fatal(uerr)
	}
	t.Log(string(bvmResult.Ret))
}

func TestRPC_GetSupportedVersionByHostname(t *testing.T) {
	ast := assert.New(t)

	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("success", func(t *testing.T) {
		receipt, err := rp.SetSupportedVersionByID(1)
		ast.Nil(err)
		ast.NotNil(receipt)
		time.Sleep(1 * time.Second)
		sv, gerr := rp.GetSupportedVersionByHostname("node1")
		ast.NotNil(sv)
		ast.Nil(gerr)
		t.Logf("%#v", sv)
	})

	t.Run("not found node's hostname", func(t *testing.T) {
		rvs, err := rp.GetSupportedVersionByHostname("not found")
		ast.Nil(rvs)
		ast.Error(err)
	})
}

func TestRPC_GetHyperchainVersionFromBin(t *testing.T) {
	ast := assert.New(t)

	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("success", func(t *testing.T) {
		rvs, err := rp.GetHyperchainVersionFromBin("2.9.0")
		ast.NotNil(rvs)
		ast.Nil(err)
		ast.Len(rvs, 3)
		t.Logf("%#v", rvs)
		ast.Equal("4.0", rvs[types.TXVersionTag])
		ast.Equal("4.0", rvs[types.BlockVersionTag])
	})

	t.Run("not found hyperchain version", func(t *testing.T) {
		rvs, err := rp.GetHyperchainVersionFromBin("not found")
		ast.Nil(rvs)
		ast.Error(err)

		rvs, err = rp.GetHyperchainVersionFromBin("111")
		ast.Nil(rvs)
		ast.Error(err)
	})
}

func TestRPC_SystemUpgrade(t *testing.T) {
	ast := assert.New(t)

	creatorKey, aerr := account.NewAccountFromAccountJSON(accounts[0], "")
	if aerr != nil {
		t.Logf("new account err: %v", aerr)
	}
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("system upgrade", func(t *testing.T) {

		targetVersion := "2.9.0"
		cancelProposal := func(pId uint64) error {
			cancelOperation := bvm.NewProposalCancelOperation(int(pId))
			payload := bvm.EncodeOperation(cancelOperation)
			tx := NewTransaction(creatorKey.GetAddress().Hex()).Invoke(cancelOperation.Address(), payload).VMType(BVM)
			receipt, err := rp.SignAndInvokeContract(tx, creatorKey)
			if err != nil {
				return err
			}
			bvmResult := bvm.Decode(receipt.Ret)
			if bvmResult == nil {
				return errors.New("bvm result is nil")
			}
			if !bvmResult.Success {
				return fmt.Errorf("failed to cancel proposal %v, reason: %v", pId, string(bvmResult.Ret))
			}
			t.Logf("successfully cancel proposal %v", pId)
			return nil
		}

		// 0. system upgrade failed, because there is no available version used to upgrade
		{
			// create proposal
			proposalOP, err := rp.NewProposalCreateOperationForSystemUpgrade(targetVersion)
			ast.NotNil(proposalOP)
			ast.Nil(err)
			payload := bvm.EncodeOperation(proposalOP)
			tx := NewTransaction(creatorKey.GetAddress().Hex()).Invoke(proposalOP.Address(), payload).VMType(BVM)
			receipt, err := rp.SignAndInvokeContract(tx, creatorKey)
			if err != nil {
				t.Fatal(err)
			}
			bvmResult := bvm.Decode(receipt.Ret)
			if bvmResult == nil {
				t.Fatal("bvm result is nil")
			}
			ast.False(bvmResult.Success)
			t.Log(string(bvmResult.Ret))
		}

		// 1. initial chain through set supportedVersion
		{
			TestRPC_SetSupportedVersionByID(t)

			operation := bvm.NewGetLatestVersionsOperation()
			payload := bvm.EncodeOperation(operation)
			tx := NewTransaction(creatorKey.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
			result, err := rp.SignAndInvokeContract(tx, creatorKey)
			if err != nil {
				t.Fatal(err)
			}
			bvmResult := bvm.Decode(result.Ret)
			var vr struct {
				SupportedVersions types.SupportedVersions `json:"supportedVersions"`
				AvailableVersion  types.AvailableVersion  `json:"availableVersions"`
				RunningVersion    types.RunningVersion    `json:"runningVersions"`
			}
			if uerr := json.Unmarshal(bvmResult.Ret, &vr); uerr != nil {
				t.Fatal(uerr)
			}
			if vr.AvailableVersion == nil {
				t.Log("not found available version!!!!!!")
				return
			}
			t.Log(string(bvmResult.Ret))
		}

		// 2. create proposal
		var pd bvm.ProposalData
		{
			// create proposal
			proposalOP, err := rp.NewProposalCreateOperationForSystemUpgrade(targetVersion)
			if err != nil {
				t.Fatal(err)
			}
			payload := bvm.EncodeOperation(proposalOP)

			tx := NewTransaction(creatorKey.GetAddress().Hex()).Invoke(proposalOP.Address(), payload).VMType(BVM)
			receipt, err := rp.SignAndInvokeContract(tx, creatorKey)
			if err != nil {
				t.Fatal(err)
			}

			bvmResult := bvm.Decode(receipt.Ret)
			if bvmResult == nil {
				t.Fatal("bvm result is nil")
			}
			if !bvmResult.Success {
				t.Log(string(bvmResult.Ret))
				t.Logf(bvmResult.Err)
				return
			}
			// get proposal
			err = proto.Unmarshal(bvmResult.Ret, &pd)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("proposal: %#v", pd)
		}

		// 3. vote proposal
		{
			voteOperation := bvm.NewProposalVoteOperation(int(pd.Id), true)
			payload := bvm.EncodeOperation(voteOperation)
			for i := range accounts {
				// 获取账户并且创建交易
				accJSON := &AccJSON{}
				umerr := json.Unmarshal([]byte(accounts[i]), accJSON)
				if umerr != nil {
					t.Logf("unmarshal account err: %v", umerr)
					continue
				}

				var (
					tx   *Transaction
					akey interface{}
				)
				if accJSON.isSM() {
					sm2key, aerr := account.NewAccountSm2FromAccountJSON(accounts[i], "")
					if aerr != nil {
						_ = cancelProposal(pd.Id)
						t.Fatalf("new account err: %v", aerr)
					}
					tx = NewTransaction(sm2key.GetAddress().Hex()).Invoke(voteOperation.Address(), payload).VMType(BVM)
					akey = sm2key
				} else {
					ecdsakey, aerr := account.NewAccountFromAccountJSON(accounts[i], "")
					if aerr != nil {
						_ = cancelProposal(pd.Id)
						t.Fatalf("new account err: %v", aerr)
					}
					tx = NewTransaction(ecdsakey.GetAddress().Hex()).Invoke(voteOperation.Address(), payload).VMType(BVM)
					akey = ecdsakey
				}

				// 调用合约
				re, err := rp.SignAndInvokeContract(tx, akey)
				if err != nil {
					_ = cancelProposal(pd.Id)
					t.Fatalf("invoke err: %v", err)
				}
				bvmResult := bvm.Decode(re.Ret)
				if bvmResult == nil {
					_ = cancelProposal(pd.Id)
					t.Fatal("failed to decode ret")
				}
				if !bvmResult.Success {
					t.Log(bvmResult.Err, string(bvmResult.Ret))
					continue
				}
			}
		}

		// 4. execute proposal
		{
			// check the status of proposal
			p, err := rp.GetProposal()
			if err != nil {
				t.Fatal(err)
			}
			if p == nil {
				t.Fatal("no proposal")
			}
			t.Logf("the proposal %v status is %v", p.ID, p.Status)
			if p.Status != bvm.ProposalData_WAITING_EXE.String() {
				t.Fatalf("the status of proposal %v is %v, not waiting exe", p.ID, p.Status)
			}

			peo := bvm.NewProposalExecuteOperation(int(pd.Id))
			payload := bvm.EncodeOperation(peo)
			tx := NewTransaction(creatorKey.GetAddress().Hex()).Invoke(peo.Address(), payload).VMType(BVM)
			receipt, err := rp.SignAndInvokeContract(tx, creatorKey)
			if err != nil {
				t.Fatal(err)
			}
			bvmResult := bvm.Decode(receipt.Ret)
			if bvmResult == nil {
				t.Fatal("bvm result is nil")
			}
			if !bvmResult.Success {
				t.Log(string(bvmResult.Ret))
				t.Fatal(bvmResult.Err)
			}

			var ops []*bvm.OpResult
			uerr := json.Unmarshal(bvmResult.Ret, &ops)
			if uerr != nil {
				t.Fatal(uerr)
			}
			for _, or := range ops {
				t.Logf("%#v", *or)
			}
		}
	})
}

// 获取提案并且对提案投票
func TestAddVP_Distrubuted(t *testing.T) {
	// 获取提案
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	p, err := rp.GetProposal()
	if err != nil {
		panic(err)
	}
	if p == nil {
		t.Log("no proposal")
		return
	}
	t.Logf("proposal: %#v", *p)
	t.Logf("the proposal %v status is %v", p.ID, p.Status)

	// 对提案投票
	urls := []string{
		"localhost:8081",
		"localhost:8082",
		"localhost:8083",
		"localhost:8084",
	}

	voteOperation := bvm.NewProposalVoteOperation(int(p.ID), true)
	payload := bvm.EncodeOperation(voteOperation)
	for i := range accounts {
		// 获取账户并且创建交易
		accJSON := &AccJSON{}
		umerr := json.Unmarshal([]byte(accounts[i]), accJSON)
		if umerr != nil {
			t.Logf("unmarshal account err: %v", umerr)
			continue
		}

		var (
			tx   *Transaction
			akey interface{}
		)
		if accJSON.isSM() {
			sm2Key, aerr := account.NewAccountSm2FromAccountJSON(accounts[i], "")
			if aerr != nil {
				t.Logf("new account err: %v", aerr)
				continue
			}
			tx = NewTransaction(sm2Key.GetAddress().Hex()).Invoke(voteOperation.Address(), payload).VMType(BVM)
			akey = sm2Key
		} else {
			ecdsaKey, aerr := account.NewAccountFromAccountJSON(accounts[i], "")
			if aerr != nil {
				t.Logf("new account err: %v", aerr)
				continue
			}
			tx = NewTransaction(ecdsaKey.GetAddress().Hex()).Invoke(voteOperation.Address(), payload).VMType(BVM)
			akey = ecdsaKey
		}

		// 调用合约
		re, err := rp.SignAndInvokeContract(tx, akey)
		if err != nil {
			t.Logf("%v invoke err: %v", urls[0], err)
			continue
		}
		result := bvm.Decode(re.Ret)
		t.Logf("%#v", result)
		//if result.Ret != nil {
		//	var res []*bvm.OpResult
		//	if uerr := json.Unmarshal(result.Ret, &res); uerr != nil {
		//		t.Logf("unmarshal err: %v", uerr)
		//		continue
		//	}
		//	for _, r := range res {
		//		t.Logf("%#v", *r)
		//	}
		//}
	}

	// 再次获取提案，确认提案状态
	time.Sleep(1 * time.Second)
	p, err = rp.GetProposal()
	if err != nil {
		panic(err)
	}
	if p == nil {
		t.Log("no proposal")
		return
	}
	t.Logf("proposal: %#v", *p)
	t.Logf("the proposal %v status is %v", p.ID, p.Status)

	time.Sleep(1 * time.Second)
	p, err = rp.GetProposal()
	if err != nil {
		panic(err)
	}
	if p == nil {
		t.Log("no proposal")
		return
	}
	t.Logf("proposal: %#v", *p)
	t.Logf("the proposal %v status is %v", p.ID, p.Status)

}

func TestDeleteVP(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	hostname := "node5"
	ns := "global"

	// 创建提案
	{
		sm2Key, aerr := account.NewAccountSm2FromAccountJSON(accounts[1], "")
		if aerr != nil {
			panic(aerr)
		}
		delOp := bvm.NewNodeRemoveVPOperation(hostname, ns)
		op := bvm.NewProposalCreateOperationForNode(delOp)
		payload := bvm.EncodeOperation(op)
		tx := NewTransaction(sm2Key.GetAddress().Hex()).Invoke(op.Address(), payload).VMType(BVM)
		re, err := rp.SignAndInvokeContract(tx, sm2Key)
		if err != nil {
			t.Logf("invoke err: %v", err)
			return
		}
		result := bvm.Decode(re.Ret)
		t.Logf("%#v", result)
	}

	// 获取提案
	time.Sleep(1 * time.Second)
	p, err := rp.GetProposal()
	if err != nil {
		panic(err)
	}
	if p == nil {
		t.Log("no proposal")
		return
	}
	t.Logf("proposal: %#v", *p)
	t.Logf("the proposal %v status is %v", p.ID, p.Status)

	// 投票提案
	{
		voteOperation := bvm.NewProposalVoteOperation(int(p.ID), true)
		payload := bvm.EncodeOperation(voteOperation)
		for i := range accounts {
			// 获取账户并且创建交易
			accJSON := &AccJSON{}
			umerr := json.Unmarshal([]byte(accounts[i]), accJSON)
			if umerr != nil {
				t.Logf("unmarshal account err: %v", umerr)
				continue
			}

			var (
				tx   *Transaction
				akey interface{}
			)
			if accJSON.isSM() {
				sm2Key, aerr := account.NewAccountSm2FromAccountJSON(accounts[i], "")
				if aerr != nil {
					t.Logf("new account err: %v", aerr)
					continue
				}
				tx = NewTransaction(sm2Key.GetAddress().Hex()).Invoke(voteOperation.Address(), payload).VMType(BVM)
				akey = sm2Key
			} else {
				ecdsaKey, aerr := account.NewAccountFromAccountJSON(accounts[i], "")
				if aerr != nil {
					t.Logf("new account err: %v", aerr)
					continue
				}
				tx = NewTransaction(ecdsaKey.GetAddress().Hex()).Invoke(voteOperation.Address(), payload).VMType(BVM)
				akey = ecdsaKey
			}

			// 调用合约
			re, err := rp.SignAndInvokeContract(tx, akey)
			if err != nil {
				t.Logf("invoke err: %v", err)
				continue
			}
			result := bvm.Decode(re.Ret)
			t.Logf("%#v", result)
		}

		// 再次获取提案，确认提案状态
		time.Sleep(1 * time.Second)
		p, err = rp.GetProposal()
		if err != nil {
			panic(err)
		}
		if p == nil {
			t.Log("no proposal")
			return
		}
		t.Logf("proposal: %#v", *p)
		t.Logf("the proposal %v status is %v", p.ID, p.Status)
	}

	// 执行提案
	{
		sm2Key, aerr := account.NewAccountSm2FromAccountJSON(accounts[1], "")
		if aerr != nil {
			panic(aerr)
		}
		op := bvm.NewProposalExecuteOperation(int(p.ID))
		payload := bvm.EncodeOperation(op)
		tx := NewTransaction(sm2Key.GetAddress().Hex()).Invoke(op.Address(), payload).VMType(BVM)
		re, err := rp.SignAndInvokeContract(tx, sm2Key)
		if err != nil {
			t.Logf("invoke err: %v", err)
			return
		}
		result := bvm.Decode(re.Ret)
		t.Logf("%#v", result)
	}
}

func TestDecodeBVMResult(t *testing.T) {
	result := bvm.Decode("0x7b2253756363657373223a747275652c22526574223a22573373695932396b5a5349364d6a417766537837496d4e765a4755694f6a49774d48307365794a6a6232526c496a6f744d7a41774d445173496d317a5a794936496d4e68624777675157526b566c41675a584a794f6d466a59323931626e51364d4867785a474d774d324d79597a6c6959546b7a5a6a63355a57497a4d6d4a6c5a6a4e694e6a5a6d4e6a566b596a56695a5751774e7a646949474673636d56685a486b676147467a49484a7662475536626d396b5a55396d566c41696656303d222c22457272223a22227d")
	t.Logf("%#v", *result)
	t.Log(string(result.Ret))
}
