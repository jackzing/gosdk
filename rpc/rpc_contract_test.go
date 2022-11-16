package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/go-crypto-standard/asym"
	"github.com/hyperchain/go-crypto-standard/hash"
	"github.com/jackzing/gosdk/abi"
	"github.com/jackzing/gosdk/account"
	"github.com/jackzing/gosdk/bvm"
	"github.com/jackzing/gosdk/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRPC_Contract_Smock(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	_, privateKey := testPrivateAccount()
	adminCount := 6
	source, _ := ioutil.ReadFile("../conf/contract/Accumulator.sol")
	//bin := `6060604052341561000f57600080fd5b5b6104c78061001f6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680635b6beeb914610049578063e15fe02314610120575b600080fd5b341561005457600080fd5b6100a4600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061023a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100e55780820151818401525b6020810190506100c9565b50505050905090810190601f1680156101125780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561012b57600080fd5b6101be600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061034f565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101ff5780820151818401525b6020810190506101e3565b50505050905090810190601f16801561022c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102426103e2565b6000826040518082805190602001908083835b60208310151561027b57805182525b602082019150602081019050602083039250610255565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103425780601f1061031757610100808354040283529160200191610342565b820191906000526020600020905b81548152906001019060200180831161032557829003601f168201915b505050505090505b919050565b6103576103e2565b816000846040518082805190602001908083835b60208310151561039157805182525b60208201915060208101905060208303925061036b565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090805190602001906103d79291906103f6565b508190505b92915050565b602060405190810160405280600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061043757805160ff1916838001178555610465565b82800160010185558215610465579182015b82811115610464578251825591602001919060010190610449565b5b5090506104729190610476565b5090565b61049891905b8082111561049457600081600090555060010161047c565b5090565b905600a165627a7a723058208ac1d22e128cf8381d7ac66b4c438a6a906ccf5ee583c3a9e46d4cdf7b3f94580029`
	//abiRaw := `[{"constant":false,"inputs":[{"name":"key","type":"string"}],"name":"getHash","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"key","type":"string"},{"name":"value","type":"string"}],"name":"setHash","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

	t.Run("set_proposal.contract.vote.enable_false", func(t *testing.T) {
		setProposalContractVoteEnable(t, rp, adminCount, false)
	})

	// proposal.contract.vote.enable = false
	t.Run("deploy_success", func(t *testing.T) {
		// use deployContract deploy contract success
		contractAddr := deploySuccess(t, rp, privateKey)

		// invoke contract
		invokeContractSuccess(t, contractAddr, rp, privateKey)

		// get proposal.contract.vote.enable value
		assertContractVote(t, false, rp)
	})

	t.Run("deploy_by_vote_fail", func(t *testing.T) {
		ope := bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil)
		createProposalForContractFail(t, rp, privateKey, ope)
	})

	t.Run("upgrade_success", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false, rp)
		// deploy
		contractAddr := deploySuccess(t, rp, privateKey)
		// upgrade
		tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(1, contractAddr, binContract).VMType(EVM)
		tx.Sign(privateKey)
		re, err := rp.MaintainContract(tx)
		assert.Nil(t, err)
		assert.NotNil(t, re)

		// contract status is normal
		assertContractStatusSuccess(t, rp, contractAddr, `"normal"`)
	})

	t.Run("upgrade_by_vote_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false, rp)
		// deploy
		contractAddr := deploySuccess(t, rp, privateKey)
		// upgrade fail
		ope := bvm.NewContractUpgradeContractOperation(source, common.Hex2Bytes(binContract), "evm", contractAddr, nil)
		createProposalForContractFail(t, rp, privateKey, ope)
	})

	t.Run("freeze_success", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false, rp)
		// deploy
		contractAddr := deploySuccess(t, rp, privateKey)
		// freeze success
		maintainSuccess(t, rp, privateKey, contractAddr, 2)

		// contract status is frozen
		assertContractStatusSuccess(t, rp, contractAddr, `"frozen"`)
	})

	t.Run("freeze_by_vote_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false, rp)
		// deploy
		contractAddr := deploySuccess(t, rp, privateKey)
		// freeze fail
		ope := bvm.NewContractMaintainContractOperation(contractAddr, "evm", 2)
		createProposalForContractFail(t, rp, privateKey, ope)
	})

	t.Run("unfreeze_success", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false, rp)
		// deploy
		contractAddr := deploySuccess(t, rp, privateKey)
		// freeze
		maintainSuccess(t, rp, privateKey, contractAddr, 2)
		// unfreeze
		maintainSuccess(t, rp, privateKey, contractAddr, 3)

		// contract status is normal
		assertContractStatusSuccess(t, rp, contractAddr, `"normal"`)
	})

	t.Run("unfreeze_by_vote_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false, rp)
		// deploy
		contractAddr := deploySuccess(t, rp, privateKey)
		// freeze
		maintainSuccess(t, rp, privateKey, contractAddr, 2)
		// unfreeze fail
		ope := bvm.NewContractMaintainContractOperation(contractAddr, "evm", 3)
		createProposalForContractFail(t, rp, privateKey, ope)
	})

	t.Run("destroy_success", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false, rp)
		// deploy
		contractAddr := deploySuccess(t, rp, privateKey)
		// destroy_success
		maintainSuccess(t, rp, privateKey, contractAddr, 5)
		// contract status is 'destroy'
		assertContractStatusSuccess(t, rp, contractAddr, `"destroy"`)
	})

	t.Run("destroy_by_vote_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = false
		assertContractVote(t, false, rp)
		// deploy
		contractAddr := deploySuccess(t, rp, privateKey)
		// destroy fail
		ope := bvm.NewContractMaintainContractOperation(contractAddr, "evm", 5)
		createProposalForContractFail(t, rp, privateKey, ope)
	})

	t.Run("set_proposal.contract.vote.enable_true", func(t *testing.T) {
		setProposalContractVoteEnable(t, rp, adminCount, true)
	})

	// proposal for contract
	t.Run("create_by_contractManager", func(t *testing.T) {
		// grant contractManager to a new account
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, rp, adminCount, bvm.NewPermissionGrantOperation("contractManager", newKey.GetAddress().Hex()))

		// use new account create proposal for contract
		createProposalForContractSuccess(t, rp, newKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))

		// cancel proposal
		proposal, _ := rp.GetProposal()
		invokeProposalContractSuccess(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), newKey, t)
	})

	t.Run("create_by_admin", func(t *testing.T) {
		// grant admin to a new account
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, rp, adminCount, bvm.NewPermissionGrantOperation("admin", newKey.GetAddress().Hex()))

		// use new account create proposal for contract
		createProposalForContractSuccess(t, rp, newKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))

		// cancel proposal
		proposal, _ := rp.GetProposal()
		invokeProposalContractSuccess(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), newKey, t)
	})

	t.Run("create_by_normal", func(t *testing.T) {
		newKey := genNewAccountKey(t)

		// use new account create proposal for contract
		createProposalForContractSuccess(t, rp, newKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))

		// cancel proposal
		proposal, _ := rp.GetProposal()
		invokeProposalContractSuccess(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), newKey, t)
	})

	t.Run("vote_by_contractManager", func(t *testing.T) {
		createProposalForContractSuccess(t, rp, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))

		proposal, _ := rp.GetProposal()
		key, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
		invokeProposalContractSuccess(rp, bvm.NewProposalVoteOperation(int(proposal.ID), true), key, t)

		invokeProposalContractSuccess(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
	})

	t.Run("vote_by_normal", func(t *testing.T) {
		createProposalForContractSuccess(t, rp, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))

		newKey := genNewAccountKey(t)
		t.Log(newKey.GetAddress().Hex())
		proposal, _ := rp.GetProposal()
		invokeProposalContractFail(rp, bvm.NewProposalVoteOperation(int(proposal.ID), true), newKey, t)

		invokeProposalContractSuccess(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
	})

	t.Run("execute_by_creator", func(t *testing.T) {
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		assertContractStatusSuccess(t, rp, res[0].Msg, `"normal"`)
	})

	t.Run("execute_by_other", func(t *testing.T) {
		createProposalForContractSuccess(t, rp, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		proposal, _ := rp.GetProposal()
		voteProposalByAdminCount(int(proposal.ID), adminCount, t, 0, rp)
		newKey := genNewAccountKey(t)
		invokeProposalContractFail(rp, bvm.NewProposalExecuteOperation(int(proposal.ID)), newKey, t)

		invokeProposalContractSuccess(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
	})

	t.Run("cancel_when_voting_by_creator", func(t *testing.T) {
		createProposalForContractSuccess(t, rp, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		proposal, _ := rp.GetProposal()
		assert.Equal(t, "VOTING", proposal.Status)
		invokeProposalContractSuccess(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
		proposal, _ = rp.GetProposal()
		assert.Equal(t, "CANCEL", proposal.Status)
	})

	t.Run("cancel_when_waitingExe_by_creator", func(t *testing.T) {
		createProposalForContractSuccess(t, rp, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		proposal, _ := rp.GetProposal()
		voteProposalByAdminCount(int(proposal.ID), adminCount, t, 0, rp)
		proposal, _ = rp.GetProposal()
		assert.Equal(t, "WAITING_EXE", proposal.Status)
		invokeProposalContractSuccess(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
		proposal, _ = rp.GetProposal()
		assert.Equal(t, "CANCEL", proposal.Status)
	})

	t.Run("cancel_when_voting_by_other", func(t *testing.T) {
		createProposalForContractSuccess(t, rp, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		proposal, _ := rp.GetProposal()
		assert.Equal(t, "VOTING", proposal.Status)
		newKey := genNewAccountKey(t)
		invokeProposalContractFail(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), newKey, t)
		proposal, _ = rp.GetProposal()
		assert.Equal(t, "VOTING", proposal.Status)

		invokeProposalContractSuccess(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
	})

	t.Run("cancel_when_waitingExe_by_other", func(t *testing.T) {
		createProposalForContractSuccess(t, rp, privateKey, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		proposal, _ := rp.GetProposal()
		voteProposalByAdminCount(int(proposal.ID), adminCount, t, 0, rp)
		proposal, _ = rp.GetProposal()
		assert.Equal(t, "WAITING_EXE", proposal.Status)
		newKey := genNewAccountKey(t)
		invokeProposalContractFail(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), newKey, t)
		proposal, _ = rp.GetProposal()
		assert.Equal(t, "WAITING_EXE", proposal.Status)

		invokeProposalContractSuccess(rp, bvm.NewProposalCancelOperation(int(proposal.ID)), privateKey, t)
	})

	// proposal for permission
	t.Run("create_contractManager_role", func(t *testing.T) {
		key, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
		invokeProposalContractFail(rp, bvm.NewProposalCreateOperationForPermission(bvm.NewPermissionCreateRoleOperation("contractManager")), key, t)
	})

	t.Run("delete_admin_role", func(t *testing.T) {
		key, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
		invokeProposalContractFail(rp, bvm.NewProposalCreateOperationForPermission(bvm.NewPermissionDeleteRoleOperation("admin")), key, t)
	})

	t.Run("delete_contractManager_role", func(t *testing.T) {
		key, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
		invokeProposalContractFail(rp, bvm.NewProposalCreateOperationForPermission(bvm.NewPermissionDeleteRoleOperation("contractManager")), key, t)
	})

	t.Run("grant_admin_to_normal", func(t *testing.T) {
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, rp, adminCount, bvm.NewPermissionGrantOperation("admin", newKey.GetAddress().Hex()))
		roles, err := rp.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, "admin", roles[0])
	})

	t.Run("grant_contractManager_to_normal", func(t *testing.T) {
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, rp, adminCount, bvm.NewPermissionGrantOperation("contractManager", newKey.GetAddress().Hex()))
		roles, err := rp.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, "contractManager", roles[0])
	})

	t.Run("revoke_contractManager", func(t *testing.T) {
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, rp, adminCount, bvm.NewPermissionGrantOperation("contractManager", newKey.GetAddress().Hex()))
		roles, err := rp.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, "contractManager", roles[0])

		completePermissionProposal(t, rp, adminCount, bvm.NewPermissionRevokeOperation("contractManager", newKey.GetAddress().Hex()))
		roles, err = rp.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 0)
	})

	t.Run("revoke_admin", func(t *testing.T) {
		newKey := genNewAccountKey(t)
		completePermissionProposal(t, rp, adminCount, bvm.NewPermissionGrantOperation("admin", newKey.GetAddress().Hex()))
		roles, err := rp.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, "admin", roles[0])

		completePermissionProposal(t, rp, adminCount, bvm.NewPermissionRevokeOperation("admin", newKey.GetAddress().Hex()))
		roles, err = rp.GetRoles(newKey.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Len(t, roles, 0)
	})

	// proposal.contract.vote.enable = true
	t.Run("deploy_by_vote_success", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true, rp)

		// deploy by vote success
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg
		// invoke
		invokeContractSuccess(t, contractAddr, rp, privateKey)
	})

	t.Run("deploy_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true, rp)

		// deploy fail
		tx := NewTransaction(privateKey.GetAddress().Hex()).Deploy(binContract).VMType(EVM)
		tx.Sign(privateKey)
		_, err := rp.DeployContract(tx)
		assert.Error(t, err)
	})

	t.Run("upgrade_by_vote_success", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true, rp)

		// deploy
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// update by vote success
		res = manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractUpgradeContractOperation(source, common.Hex2Bytes(binContract), "evm", contractAddr, nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
	})

	t.Run("upgrade_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true, rp)

		// deploy
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// update fail
		tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(1, contractAddr, binContract).VMType(EVM)
		tx.Sign(privateKey)
		_, err := rp.MaintainContract(tx)
		assert.Error(t, err)
	})

	t.Run("freeze_by_vote_success", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true, rp)

		// deploy
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// freeze by vote success
		res = manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractMaintainContractOperation(contractAddr, "evm", 2))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)

		assertContractStatusSuccess(t, rp, contractAddr, `"frozen"`)
	})

	t.Run("freeze_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true, rp)

		// deploy
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// freeze fail
		tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(2, contractAddr, "").VMType(EVM)
		tx.Sign(privateKey)
		_, err := rp.MaintainContract(tx)
		assert.Error(t, err)

		assertContractStatusSuccess(t, rp, contractAddr, `"normal"`)

	})

	t.Run("unfreeze_by_vote_success", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true, rp)

		// deploy
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// freeze
		res = manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractMaintainContractOperation(contractAddr, "evm", 2))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)

		assertContractStatusSuccess(t, rp, contractAddr, `"frozen"`)
		// unfreeze by vote success
		res = manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractMaintainContractOperation(contractAddr, "evm", 3))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		assertContractStatusSuccess(t, rp, contractAddr, `"normal"`)

	})

	t.Run("unfreeze_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true, rp)

		// deploy
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		// freeze
		res = manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractMaintainContractOperation(contractAddr, "evm", 2))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)

		assertContractStatusSuccess(t, rp, contractAddr, `"frozen"`)
		// unfreeze fail
		tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(3, contractAddr, "").VMType(EVM)
		tx.Sign(privateKey)
		_, err := rp.MaintainContract(tx)
		assert.Error(t, err)

		assertContractStatusSuccess(t, rp, contractAddr, `"frozen"`)

	})

	t.Run("destroy_by_vote_success", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true, rp)

		// deploy
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		assertContractStatusSuccess(t, rp, contractAddr, `"normal"`)
		// destroy by vote success
		res = manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractMaintainContractOperation(contractAddr, "evm", 5))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		assertContractStatusSuccess(t, rp, contractAddr, `"destroy"`)
	})

	t.Run("destroy_fail", func(t *testing.T) {
		// proposal.contract.vote.enable = true
		assertContractVote(t, true, rp)

		// deploy
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg

		assertContractStatusSuccess(t, rp, contractAddr, `"normal"`)
		// unfreeze fail
		tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(5, contractAddr, "").VMType(EVM)
		tx.Sign(privateKey)
		_, err := rp.MaintainContract(tx)
		assert.Error(t, err)
		assertContractStatusSuccess(t, rp, contractAddr, `"normal"`)

	})

}

func TestRPC_ThresholdAndVP(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	t.Run("init_node", func(t *testing.T) {
		creator, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
		nodes := []string{"node1", "node2", "node3", "node4"}
		ns := "global"
		role := "vp"
		pub := []byte("pub1")
		var ops []bvm.NodeOperation
		for _, n := range nodes {
			ops = append(ops, bvm.NewNodeAddNodeOperation(pub, n, role, ns))
			ops = append(ops, bvm.NewNodeAddVPOperation(n, ns))
		}
		res := completeProposal(t, rp, creator, 6, bvm.NewProposalCreateOperationForNode(ops...))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
	})

	t.Run("remove_vp", func(t *testing.T) {
		creator, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
		res := completeProposal(t, rp, creator, 6, bvm.NewProposalCreateOperationForNode(bvm.NewNodeRemoveVPOperation("node1", "global")))
		assert.NotEqual(t, bvm.SuccessCode, res[0].Code)
		t.Log(res[0].Msg)
	})

	t.Run("SetProposalThreshold", func(t *testing.T) {
		creator, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
		invokeProposalContractFail(rp, bvm.NewProposalCreateOperationByConfigOps(bvm.NewSetProposalThreshold(7)), creator, t)
	})

	t.Run("revoke", func(t *testing.T) {
		creator, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
		ops := []bvm.PermissionOperation{
			bvm.NewPermissionRevokeOperation("admin", creator.GetAddress().Hex()),
			bvm.NewPermissionRevokeOperation("contractManager", creator.GetAddress().Hex()),
		}
		res := completeProposal(t, rp, creator, 6, bvm.NewProposalCreateOperationForPermission(ops...))
		assert.Len(t, res, 2)
		//assert.NotEqual(t, bvm.SuccessCode,res[0].Code)
		//assert.NotEqual(t, bvm.SuccessCode,res[1].Code)
	})

}

func TestRPC_ContractByName(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	_, privateKey := testPrivateAccount()
	adminCount := 6
	source, _ := ioutil.ReadFile("../conf/contract/Accumulator.sol")

	t.Run("set_proposal.contract.vote.enable_false", func(t *testing.T) {
		setProposalContractVoteEnable(t, rp, adminCount, false)
	})

	t.Run("maintain_by_name_success", func(t *testing.T) {
		contractAddr := deploySuccess(t, rp, privateKey)
		ri := rand.Int()
		contractName := "name" + strconv.Itoa(ri)
		setCName(contractAddr, contractName, rp, t)

		// freeze success
		tx := NewTransaction(privateKey.GetAddress().Hex()).MaintainByName(2, contractName, "").VMType(EVM)
		tx.Sign(privateKey)
		re, err := rp.MaintainContract(tx)
		assert.Nil(t, err)
		assert.NotNil(t, re)

		status, stdError := rp.GetContractStatus(contractAddr)
		assert.Nil(t, stdError)
		statu, stdError := rp.GetContractStatusByName(contractName)
		assert.Nil(t, stdError)
		assert.Equal(t, status, statu)
		assert.Equal(t, `"frozen"`, status)
	})

	t.Run("set_proposal.contract.vote.enable_true", func(t *testing.T) {
		setProposalContractVoteEnable(t, rp, adminCount, true)
	})

	t.Run("maintain_by_name_by_vote_success", func(t *testing.T) {
		// deploy
		res := manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(binContract), "evm", nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)
		contractAddr := res[0].Msg
		ri := rand.Int()
		contractName := "name" + strconv.Itoa(ri)
		setCName(contractAddr, contractName, rp, t)

		// update by vote success
		res = manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractUpgradeOperationByName(source, common.Hex2Bytes(binContract), "evm", contractName, nil))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)

		// freeze by vote success
		res = manageContractByVote(t, rp, privateKey, adminCount, bvm.NewContractMaintainOperationByName(contractName, "evm", 2))
		assert.Equal(t, bvm.SuccessCode, res[0].Code)

		status, stdError := rp.GetContractStatus(contractAddr)
		assert.Nil(t, stdError)
		statu, stdError := rp.GetContractStatusByName(contractName)
		assert.Nil(t, stdError)
		assert.Equal(t, status, statu)
		assert.Equal(t, `"frozen"`, status)
	})
}

func completeProposal(t *testing.T, rp *RPC, creatorKey account.Key, adminCount int, opt bvm.BuiltinOperation) []*bvm.OpResult {
	invokeProposalContractSuccess(rp, opt, creatorKey, t)
	proposal, _ := rp.GetProposal()
	voteProposalByAdminCount(int(proposal.ID), adminCount, t, 1, rp)
	ret := invokeProposalContractSuccess(rp, bvm.NewProposalExecuteOperation(int(proposal.ID)), creatorKey, t)
	var res []*bvm.OpResult
	_ = json.Unmarshal(ret, &res)
	return res
}

func manageContractByVote(t *testing.T, rp *RPC, creatorKey account.Key, adminCount int, ops ...bvm.ContractOperation) []*bvm.OpResult {
	createProposalForContractSuccess(t, rp, creatorKey, ops...)
	proposal, _ := rp.GetProposal()
	voteProposalByAdminCount(int(proposal.ID), adminCount, t, 0, rp)
	ret := invokeProposalContractSuccess(rp, bvm.NewProposalExecuteOperation(int(proposal.ID)), creatorKey, t)
	var res []*bvm.OpResult
	_ = json.Unmarshal(ret, &res)
	assert.Len(t, res, len(ops))
	return res
}

func invokeContractSuccess(t *testing.T, contractAddr string, rp *RPC, privateKey *account.ECDSAKey) {
	ABI, er := abi.JSON(strings.NewReader(abiContract))
	assert.Nil(t, er)
	packed, er := ABI.Pack("add", uint32(1), uint32(2))
	assert.Nil(t, er)
	tx := NewTransaction(privateKey.GetAddress().Hex()).Invoke(contractAddr, packed)
	tx.Sign(privateKey)
	_, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
}

func completePermissionProposal(t *testing.T, rp *RPC, adminCount int, ops ...bvm.PermissionOperation) {
	key, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
	invokeProposalContractSuccess(rp, bvm.NewProposalCreateOperationForPermission(ops...), key, t)
	proposal, _ := rp.GetProposal()
	voteProposalByAdminCount(int(proposal.ID), adminCount, t, 1, rp)
	invokeProposalContractSuccess(rp, bvm.NewProposalExecuteOperation(int(proposal.ID)), key, t)
}

func genNewAccountKey(t *testing.T) account.Key {
	pwd := "12347890"
	newAccount, err := account.NewAccount(pwd)
	assert.Nil(t, err)
	newKey, err := account.NewAccountFromAccountJSON(newAccount, pwd)
	assert.Nil(t, err)
	return newKey
}

func setProposalContractVoteEnable(t *testing.T, rp *RPC, adminCount int, enable bool) {
	// create proposal
	key, _ := account.NewAccountFromAccountJSON(accountJsons[0], password)
	invokeProposalContractSuccess(rp, bvm.NewProposalCreateOperationByConfigOps(bvm.NewSetContactVoteEnable(enable)), key, t)

	proposal, _ := rp.GetProposal()
	// vote
	voteProposalByAdminCount(int(proposal.ID), adminCount, t, 1, rp)
	// execute
	invokeProposalContractSuccess(rp, bvm.NewProposalExecuteOperation(int(proposal.ID)), key, t)
}

func voteProposalByAdminCount(pid, adminCount int, t *testing.T, startCount int, rp *RPC) {
	// vote
	for i := startCount; i < adminCount; i++ {
		var (
			key account.Key
		)
		if i < 4 {
			key, _ = account.NewAccountFromAccountJSON(accountJsons[i], pwd)
		} else {
			key, _ = account.NewAccountSm2FromAccountJSON(accountJsons[i], pwd)
		}
		operation := bvm.NewProposalVoteOperation(pid, true)
		invokeProposalContractSuccess(rp, operation, key, t)
	}
}

func invokeBVMSuccess(rp *RPC, operation bvm.BuiltinOperation, key account.Key, t *testing.T) *bvm.Result {
	payload := bvm.EncodeOperation(operation)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(operation.Address(), payload).VMType(BVM)
	tx.Sign(key)
	re, err := rp.InvokeContract(tx)
	assert.Nil(t, err)
	assert.NotNil(t, re)
	result := bvm.Decode(re.Ret)
	return result
}

func invokeProposalContractSuccess(rp *RPC, operation bvm.BuiltinOperation, key account.Key, t *testing.T) []byte {
	result := invokeBVMSuccess(rp, operation, key, t)
	assert.True(t, result.Success)
	t.Log(result)
	return []byte(result.Ret)
}

func invokeProposalContractFail(rp *RPC, operation bvm.BuiltinOperation, key account.Key, t *testing.T) {
	result := invokeBVMSuccess(rp, operation, key, t)
	assert.False(t, result.Success)
	t.Log(result)
}

func assertContractStatusSuccess(t *testing.T, rp *RPC, contractAddr string, expect string) {
	status, err := rp.GetContractStatus(contractAddr)
	assert.Nil(t, err)
	assert.Equal(t, expect, status)
}

func createProposalForContractFail(t *testing.T, rp *RPC, privateKey *account.ECDSAKey, ops ...bvm.ContractOperation) {
	_, err := createProposalForManagerContract(t, rp, privateKey, ops...)
	assert.Error(t, err)
}

func createProposalForContractSuccess(t *testing.T, rp *RPC, key account.Key, ops ...bvm.ContractOperation) {
	re, err := createProposalForManagerContract(t, rp, key, ops...)
	assert.Nil(t, err)
	assert.NotNil(t, re)
	result := bvm.Decode(re.Ret)
	assert.True(t, result.Success)
	t.Log(result)
}

func createProposalForManagerContract(t *testing.T, rp *RPC, key account.Key, ops ...bvm.ContractOperation) (*TxReceipt, StdError) {
	contractOpt := bvm.NewProposalCreateOperationForContract(ops...)
	payload := bvm.EncodeOperation(contractOpt)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(contractOpt.Address(), payload).VMType(BVM)
	tx.Sign(key)
	return rp.ManageContractByVote(tx)
}

func maintainSuccess(t *testing.T, rp *RPC, privateKey *account.ECDSAKey, contractAddr string, op int64) {
	tx := NewTransaction(privateKey.GetAddress().Hex()).Maintain(op, contractAddr, "").VMType(EVM)
	tx.Sign(privateKey)
	re, err := rp.MaintainContract(tx)
	assert.Nil(t, err)
	assert.NotNil(t, re)
}

func deploySuccess(t *testing.T, rp *RPC, privateKey *account.ECDSAKey) string {
	tx := NewTransaction(privateKey.GetAddress().Hex()).Deploy(binContract).VMType(EVM)
	tx.Sign(privateKey)
	re, err := rp.DeployContract(tx)
	assert.Nil(t, err)
	assert.NotNil(t, re)
	t.Log(re.ContractAddress)
	return re.ContractAddress
}

func assertContractVote(t *testing.T, expect bool, rp *RPC) {
	config, err := rp.GetConfig()
	assert.Nil(t, err)
	v := viper.New()
	v.SetConfigType("toml")
	er := v.ReadConfig(strings.NewReader(config))
	assert.Nil(t, er)
	assert.Equal(t, expect, v.GetBool("proposal.contract.vote.enable"))
}

func TestInt(t *testing.T) {
	bigThreshold := new(big.Int)
	_ = json.Unmarshal(nil, bigThreshold)
	t.Log(bigThreshold.Int64())
}

func TestRPC_CompileContract(t *testing.T) {
	//nolint
	compileContract("../conf/contract/Accumulator.sol")
}

func TestRPC_ManageContractByVote(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	source, _ := ioutil.ReadFile("../conf/contract/Accumulator.sol")
	bin := `6060604052341561000f57600080fd5b5b6104c78061001f6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680635b6beeb914610049578063e15fe02314610120575b600080fd5b341561005457600080fd5b6100a4600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061023a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100e55780820151818401525b6020810190506100c9565b50505050905090810190601f1680156101125780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561012b57600080fd5b6101be600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061034f565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101ff5780820151818401525b6020810190506101e3565b50505050905090810190601f16801561022c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102426103e2565b6000826040518082805190602001908083835b60208310151561027b57805182525b602082019150602081019050602083039250610255565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103425780601f1061031757610100808354040283529160200191610342565b820191906000526020600020905b81548152906001019060200180831161032557829003601f168201915b505050505090505b919050565b6103576103e2565b816000846040518082805190602001908083835b60208310151561039157805182525b60208201915060208101905060208303925061036b565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090805190602001906103d79291906103f6565b508190505b92915050565b602060405190810160405280600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061043757805160ff1916838001178555610465565b82800160010185558215610465579182015b82811115610464578251825591602001919060010190610449565b5b5090506104729190610476565b5090565b61049891905b8082111561049457600081600090555060010161047c565b5090565b905600a165627a7a723058208ac1d22e128cf8381d7ac66b4c438a6a906ccf5ee583c3a9e46d4cdf7b3f94580029`
	ope := bvm.NewContractDeployContractOperation(source, common.Hex2Bytes(bin), "evm", nil)
	contractOpt := bvm.NewProposalCreateOperationForContract(ope)
	payload := bvm.EncodeOperation(contractOpt)
	_, privateKey := testPrivateAccount()
	tx := NewTransaction(privateKey.GetAddress().Hex()).Invoke(contractOpt.Address(), payload).VMType(BVM)
	re, err := rp.SignAndManageContractByVote(tx, privateKey)
	assert.NoError(t, err)
	result := bvm.Decode(re.Ret)
	assert.True(t, result.Success)
	var proposal bvm.ProposalData
	_ = proto.Unmarshal(result.Ret, &proposal)
	t.Log(proposal.String())
}

func TestRPC_DeployContract(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0])
	receipt, _ := rp.SignAndDeployContract(transaction, guomiKey)
	fmt.Println("address:", receipt.ContractAddress)
}

func TestRPC_InvokeContract(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	contractAddress, err := deployContract(cr.Bin[0], cr.Abi[0])
	ABI, serr := abi.JSON(strings.NewReader(cr.Abi[0]))
	if err != nil {
		t.Error(serr)
		return
	}
	packed, serr := ABI.Pack("add", uint32(1), uint32(2))
	if serr != nil {
		t.Error(serr)
		return
	}
	address, privateKey := testPrivateAccount()
	rp, serr := NewJsonRPC()
	assert.Nil(t, serr)
	transaction := NewTransaction(address).Invoke(contractAddress, packed)
	receipt, _ := rp.SignAndInvokeContract(transaction, privateKey)
	fmt.Println("ret:", receipt.Ret)
}

func TestRPC_GetCode(t *testing.T) {
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	receipt, err := rp.SignAndDeployContract(transaction, guomiKey)
	assert.Nil(t, err)
	_, err = rp.GetCode(receipt.ContractAddress)
	assert.Nil(t, err)
}

func TestRPC_GetContractCountByAddr(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	count, err := rp.GetContractCountByAddr(common.BytesToAddress(newAddress).Hex())
	if err != nil {
		t.Error(err)
	}
	fmt.Println(count)
}

func TestRPC_EncryptoMessage(t *testing.T) {
	t.Skip("method not found")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	count, err := rp.EncryptoMessage(100, 10, "123456")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(count)
}

func TestRPC_CheckHmValue(t *testing.T) {
	t.Skip("method not found")
	rp := NewRPC()
	count, err := rp.CheckHmValue([]uint64{1, 2}, []string{"123", "456"}, "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(count)
}

func TestRPC_DeployContractWithArgs(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator2.sol")
	var arg [32]byte
	copy(arg[:], "test")
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0]).DeployArgs(cr.Abi[0], uint32(10), arg)
	receipt, _ := rp.SignAndDeployContract(transaction, guomiKey)
	fmt.Println("address:", receipt.ContractAddress)

	fmt.Println("-----------------------------------")

	address, privateKey := testPrivateAccount()
	rp, serr := NewJsonRPC()
	assert.Nil(t, serr)
	ABI, _ := abi.JSON(strings.NewReader(cr.Abi[0]))
	packed, _ := ABI.Pack("getMul")
	transaction1 := NewTransaction(address).Invoke(receipt.ContractAddress, packed)
	receipt1, err := rp.SignAndInvokeContract(transaction1, privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("ret:", receipt1.Ret)

	var p0 []byte
	var p1 int64
	var p2 common.Address
	testV := []interface{}{&p0, &p1, &p2}
	fmt.Println(reflect.TypeOf(testV))
	decode(ABI, &testV, "getMul", receipt1.Ret)
	fmt.Println(string(p0), p1, p2.Hex())
}

func TestRPC_DeployContractWithStringArgs(t *testing.T) {
	t.Skip("solc")
	guomiKey := testGuomiKey()
	cr, _ := compileContract("../conf/contract/Accumulator2.sol")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0]).DeployStringArgs(cr.Abi[0], "10", "test")
	receipt, _ := rp.SignAndDeployContract(transaction, guomiKey)
	fmt.Println("address:", receipt.ContractAddress)
	fmt.Println("-----------------------------------")
	ABI, _ := abi.JSON(strings.NewReader(cr.Abi[0]))
	packed, _ := ABI.Encode("getMul")
	address, privateKey := testPrivateAccount()
	invokeTx := NewTransaction(address).Invoke(receipt.ContractAddress, packed)
	invokeReceipt, err := rp.SignAndInvokeContract(invokeTx, privateKey)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("ret:", invokeReceipt.Ret)
	ret, e := ABI.Decode("getMul", common.FromHex(invokeReceipt.Ret))
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Printf("%v\n", ret)
}

func TestRPC_UnpackLog(t *testing.T) {
	t.Skip("solc")
	cr, _ := compileContract("../conf/contract/Accumulator2.sol")
	var arg [32]byte
	copy(arg[:], "test")
	guomiKey := testGuomiKey()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0]).DeployArgs(cr.Abi[0], uint32(10), arg)
	receipt, _ := rp.SignAndDeployContract(transaction, guomiKey)
	fmt.Println("address:", receipt.ContractAddress)

	fmt.Println("-----------------------------------")

	address, privateKey := testPrivateAccount()
	ABI, _ := abi.JSON(strings.NewReader(cr.Abi[0]))
	packed, _ := ABI.Pack("getHello")
	transaction1 := NewTransaction(address).Invoke(receipt.ContractAddress, packed)
	transaction1.Sign(privateKey)
	receipt1, err := rp.InvokeContract(transaction1)
	if err != nil {
		t.Error(err)
		return
	}
	test := struct {
		Addr int64   `abi:"addr1"`
		Msg1 [8]byte `abi:"msg"`
	}{}

	// testLog
	sysErr := ABI.UnpackLog(&test, "sayHello", receipt1.Log[0].Data, receipt1.Log[0].Topics)
	if sysErr != nil {
		t.Error(sysErr)
		return
	}
	msg, sysErr := abi.ByteArrayToString(test.Msg1)
	if sysErr != nil {
		t.Error(sysErr)
		return
	}
	assert.Equal(t, int64(1), test.Addr, "解码失败")
	assert.Equal(t, "test", msg, "解码失败")
}

func TestRPC_SendTx(t *testing.T) {
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	guomiKey, _ := asym.GenerateKey(asym.AlgoP256R1)
	pubKey := &account.ECDSAKey{ECDSAPrivateKey: guomiKey}
	address, _ := testPrivateAccount()
	newAddress := pubKey.GetAddress()
	for i := 0; i < 10; i++ {
		transaction := NewTransaction(newAddress.Hex()).Transfer(address, int64(0))
		receipt, err := rp.SignAndSendTx(transaction, pubKey)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(receipt.TxHash)
		fmt.Println(transaction.GetTransactionHash(1000000))
	}
}

func TestVC_MQ(t *testing.T) {
	t.Skip("ci")
	//APP.sol
	_ = `// SPDX-License-Identifier: hyperchain

pragma solidity >=0.7.0 <0.9.0;

contract App {
    mapping(int => bytes) public triggerInfo;
    int nonce;
    address proxyContract;

    function callback(bytes32 taskID,string memory result,string memory proof) public pure returns (bool, bool, string memory){
        require(taskID != 0, "taskID");
        bytes memory r = bytes(result);
        bytes memory p = bytes(proof);
        require(r.length > 0 && p.length > 0, "length");
        return (true, false, "app callback success");
    }

    function trigger(string memory input) public {
        //computeAndProve(input, this, "callback")
        bytes memory payload = abi.encodeWithSignature("computeAndProve(string,address,string)", input, this, "callback");
        (bool success1, bytes memory response) = proxyContract.call(payload);
        require(success1,"call back error");
        nonce++;
        triggerInfo[nonce] = response;
    }

    function setAddr(address proxy) public {
        proxyContract = proxy;
    }
}`
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	//输入{"in":"4"}，输出{"out":19}
	t.Run("proxy", func(t *testing.T) {
		abiStr := `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"taskID","type":"bytes32"},{"indexed":false,"internalType":"string","name":"circuitID","type":"string"},{"indexed":false,"internalType":"address","name":"BusinessContractAddr","type":"address"},{"indexed":false,"internalType":"string","name":"BusinessContractMethod","type":"string"},{"indexed":false,"internalType":"string","name":"input","type":"string"}],"name":"EVENT_COMPUTE","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"taskID","type":"bytes32"},{"indexed":false,"internalType":"string","name":"circuitID","type":"string"},{"indexed":false,"internalType":"string","name":"proof","type":"string"},{"indexed":false,"internalType":"string","name":"result","type":"string"},{"indexed":false,"internalType":"string","name":"response","type":"string"}],"name":"EVENT_FINISH","type":"event"},{"inputs":[],"name":"circuitID","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"input","type":"string"},{"internalType":"address","name":"businessContractAddr","type":"address"},{"internalType":"string","name":"businessContractMethod","type":"string"}],"name":"computeAndProve","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"taskID","type":"bytes32"},{"internalType":"string","name":"result","type":"string"},{"internalType":"string","name":"proof","type":"string"}],"name":"finish","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"s","type":"string"}],"name":"fromHex","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"uint8","name":"c","type":"uint8"}],"name":"fromHexChar","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"pkContent","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"s","type":"bytes"}],"name":"stringToBytes32","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"taskStorage","outputs":[{"internalType":"bool","name":"Exist","type":"bool"},{"internalType":"address","name":"BusinessContractAddr","type":"address"},{"internalType":"string","name":"BusinessContractMethod","type":"string"},{"internalType":"string","name":"PublicInput","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"pk","type":"bytes"},{"internalType":"bytes","name":"vk","type":"bytes"},{"internalType":"bytes32","name":"vkTag","type":"bytes32"}],"name":"update","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"input","type":"bytes"},{"internalType":"bytes","name":"proof","type":"bytes"}],"name":"verifyProof","outputs":[{"internalType":"bool","name":"r","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"vkContent","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function"}]`
		bin := `6080604052348015620000125760006000fd5b505b60006040518061072001604052806106ea81526020016200442a6106ea91399050600060405180610ea00160405280610e68815260200162004b14610e689139905060007feb33da442302b2c45c0949a191f7ea6e79898e32e866de2ab288ee213249dc8660001b905060006200009184620000c760201b60201c565b90506000620000a684620000c760201b60201c565b9050620000bb8183856200032560201b60201c565b50505050505b62000baa565b60606000829050600060028251620000e0919062000ab1565b141515620000ee5760006000fd5b600060028251620001009190620008f4565b67ffffffffffffffff81111562000140577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519080825280601f01601f191660200182016040528015620001735781602001600182028036833780820191505090505b5090506000600090505b600283516200018d9190620008f4565b81101562000312576200020b836001836002620001ab91906200092f565b620001b7919062000857565b815181101515620001f1577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b60f81c6200037f60201b60201c565b601062000275858460026200022191906200092f565b8151811015156200025b577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b60f81c6200037f60201b60201c565b62000281919062000991565b6200028d9190620008b5565b60f81b8282815181101515620002cc577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505b80620003089062000a62565b905080506200017d565b508092505050620003205650505b919050565b8160016000506000836000191660001916815260200190815260200160002060005090805190602001906200035c92919062000739565b5082600060005090805190602001906200037892919062000739565b505b505050565b60007f30000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916101580156200046057507f39000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191611155b156200049f577f300000000000000000000000000000000000000000000000000000000000000060f81c82620004979190620009d4565b905062000734565b7f61000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916101580156200057e57507f66000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191611155b15620005cb577f610000000000000000000000000000000000000000000000000000000000000060f81c82600a620005b79190620008b5565b620005c39190620009d4565b905062000734565b7f41000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191610158015620006aa57507f46000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191611155b15620006f7577f410000000000000000000000000000000000000000000000000000000000000060f81c82600a620006e39190620008b5565b620006ef9190620009d4565b905062000734565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200072b9062000822565b60405180910390fd5b919050565b828054620007479062000a29565b90600052602060002090601f0160209004810192826200076b5760008555620007bc565b82601f106200078657805160ff1916838001178555620007bc565b82800160010185558215620007bc579182015b82811115620007bb578251826000509090559160200191906001019062000799565b5b509050620007cb9190620007cf565b5090565b620007d5565b80821115620007f15760008181506000905550600101620007d5565b50905662000ba9565b60006200080960048362000845565b9150620008168262000b7f565b6020820190505b919050565b600060208201905081810360008301526200083d81620007fa565b90505b919050565b60008282526020820190505b92915050565b6000620008648262000a10565b9150620008718362000a10565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115620008a957620008a862000aec565b5b82820190505b92915050565b6000620008c28262000a1b565b9150620008cf8362000a1b565b92508260ff03821115620008e857620008e762000aec565b5b82820190505b92915050565b6000620009018262000a10565b91506200090e8362000a10565b925082151562000923576200092262000b1d565b5b82820490505b92915050565b60006200093c8262000a10565b9150620009498362000a10565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161562000985576200098462000aec565b5b82820290505b92915050565b60006200099e8262000a1b565b9150620009ab8362000a1b565b92508160ff0483118215151615620009c857620009c762000aec565b5b82820290505b92915050565b6000620009e18262000a1b565b9150620009ee8362000a1b565b92508282101562000a045762000a0362000aec565b5b82820390505b92915050565b60008190505b919050565b600060ff821690505b919050565b60006002820490506001821680151562000a4457607f821691505b6020821081141562000a5b5762000a5a62000b4e565b5b505b919050565b600062000a6f8262000a10565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141562000aa55762000aa462000aec565b5b6001820190505b919050565b600062000abe8262000a10565b915062000acb8362000a10565b925082151562000ae05762000adf62000b1d565b5b82820690505b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b565b7f6661696c000000000000000000000000000000000000000000000000000000006000820152505b565b5b6138708062000bba6000396000f3fe60806040523480156100115760006000fd5b50600436106100ae5760003560e01c8063c39ffd4b11610072578063c39ffd4b14610190578063cc5bfa64146101ae578063e84df22c146101ca578063ef2e50d7146101e6578063f9f6694e14610216578063fd639b9514610249576100ae565b806307a72184146100b45780632ecb20d3146100e45780637510656b146101145780638e7e34d714610130578063b8e72af614610160576100ae565b60006000fd5b6100ce60048036038101906100c99190612737565b610267565b6040516100db9190612db0565b60405180910390f35b6100fe60048036038101906100f99190612947565b610444565b60405161010b9190613070565b60405180910390f35b61012e600480360381019061012991906127fb565b6107e9565b005b61014a6004803603810190610145919061287f565b61083f565b6040516101579190612eeb565b60405180910390f35b61017a60048036038101906101759190612780565b610a7e565b6040516101879190612d40565b60405180910390f35b6101986114b3565b6040516101a59190612eeb565b60405180910390f35b6101c860048036038101906101c391906126b3565b611544565b005b6101e460048036038101906101df91906128c3565b611bc1565b005b61020060048036038101906101fb919061265d565b611da9565b60405161020d9190612eeb565b60405180910390f35b610230600480360381019061022b919061265d565b611e4c565b6040516102409493929190612d5c565b60405180910390f35b610251611fc2565b60405161025e9190612f46565b60405180910390f35b6000600083838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509050600060009050604082511415156102fe576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102f590612fcc565b60405180910390fd5b6000600090505b602081101561042e5780601f61031b91906132fc565b60086103279190613265565b61039684600184600261033a9190613265565b61034491906131a2565b81518110151561037d577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b60f81c61044463ffffffff16565b60106103fb868560026103a99190613265565b8151811015156103e2577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b60f81c61044463ffffffff16565b61040591906132c0565b61040f91906131f9565b60ff16901b8217915081505b8061042590613477565b90508050610305565b508060001b9250505061043e5650505b92915050565b60007f30000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161015801561052457507f39000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191611155b1561055f577f300000000000000000000000000000000000000000000000000000000000000060f81c826105589190613331565b90506107e4565b7f61000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161015801561063d57507f66000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191611155b15610684577f610000000000000000000000000000000000000000000000000000000000000060f81c82600a61067391906131f9565b61067d9190613331565b90506107e4565b7f41000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161015801561076257507f46000000000000000000000000000000000000000000000000000000000000007effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168260f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191611155b156107a9577f410000000000000000000000000000000000000000000000000000000000000060f81c82600a61079891906131f9565b6107a29190613331565b90506107e4565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107db90612fab565b60405180910390fd5b919050565b81600160005060008360001916600019168152602001908152602001600020600050908051906020019061081e92919061221f565b50826000600050908051906020019061083892919061221f565b505b505050565b6060600082905060006002825161085691906134c1565b1415156108635760006000fd5b6000600282516108739190613231565b67ffffffffffffffff8111156108b2577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519080825280601f01601f1916602001820160405280156108e45781602001600182028036833780820191505090505b5090506000600090505b600283516108fc9190613231565b811015610a6c576109728360018360026109169190613265565b61092091906131a2565b815181101515610959577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b60f81c61044463ffffffff16565b60106109d7858460026109859190613265565b8151811015156109be577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b60f81c61044463ffffffff16565b6109e191906132c0565b6109eb91906131f9565b60f81b8282815181101515610a29577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505b80610a6390613477565b905080506108ee565b508092505050610a795650505b919050565b600060006000610a948585611fde63ffffffff16565b91509150818585604051602001610aac929190612cf4565b604051602081830303815290604052901515610afe576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610af59190612f46565b60405180910390fd5b50604b81610b0c91906131a2565b85859050148585604051602001610b24929190612caa565b604051602081830303815290604052901515610b76576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b6d9190612f46565b60405180910390fd5b50600060016000506000610bb88888600987610b9291906131a2565b90604988610ba091906131a2565b92610bad9392919061316b565b61026763ffffffff16565b600019166000191681526020019081526020016000206000508054610bdc90613410565b80601f0160208091040260200160405190810160405280929190818152602001828054610c0890613410565b8015610c555780601f10610c2a57610100808354040283529160200191610c55565b820191906000526020600020905b815481529060010190602001808311610c3857829003601f168201915b5050505050905060008151118686600985610c7091906131a2565b90604986610c7e91906131a2565b92610c8b9392919061316b565b604051602001610c9c929190612ccf565b604051602081830303815290604052901515610cee576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ce59190612f46565b60405180910390fd5b506000868690506020610d0191906131a2565b898990506020610d1191906131a2565b83516020610d1f91906131a2565b610d2991906131a2565b610d3391906131a2565b67ffffffffffffffff811115610d72577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519080825280601f01601f191660200182016040528015610da45781602001600182028036833780820191505090505b5090506000610db9835161219963ffffffff16565b905060006000905080505b6020811015610e91578181815181101515610e08577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b8382815181101515610e4e577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505b8080610e8990613477565b915050610dc4565b6000905080505b8351811015610f71578381815181101515610edc577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b83826020610ef591906131a2565b815181101515610f2e577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505b8080610f6990613477565b915050610e98565b610f838b8b905061219963ffffffff16565b915081506000905080505b6020811015611073578181815181101515610fd2577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b838286516020610fed91906131a2565b610ff791906131a2565b815181101515611030577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505b808061106b90613477565b915050610f8e565b6000905080505b8a8a905081101561116c578a8a8281811015156110c0577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9050013560f81c60f81b83826020875160206110dc91906131a2565b6110e691906131a2565b6110f091906131a2565b815181101515611129577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505b808061116490613477565b91505061107a565b61117e8989905061219963ffffffff16565b915081506000905080505b60208110156112885781818151811015156111cd577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b83828d8d905060206111ea91906131a2565b875160206111f891906131a2565b61120291906131a2565b61120c91906131a2565b815181101515611245577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505b808061128090613477565b915050611189565b6000905080505b8888905081101561139b5788888281811015156112d5577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9050013560f81c60f81b8382602060208f8f90506112f391906131a2565b8851602061130191906131a2565b61130b91906131a2565b61131591906131a2565b61131f91906131a2565b815181101515611358577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505b808061139390613477565b91505061128f565b60006113a56122aa565b8a8a905060206113b591906131a2565b8d8d905060206113c591906131a2565b875160206113d391906131a2565b6113dd91906131a2565b6113e791906131a2565b92508250602081846020880160fb6107d05a03fa9150816000811461140b5761140d565bfe5b50811515611450576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161144790612f69565b60405180910390fd5b600081600060018110151561148e577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602002015114985050505050505050506114ab5650505050505050505b949350505050565b600060005080546114c390613410565b80601f01602080910402602001604051908101604052809291908181526020018280546114ef90613410565b801561153c5780601f106115115761010080835404028352916020019161153c565b820191906000526020600020905b81548152906001019060200180831161151f57829003601f168201915b505050505081565b60006002600050600085600019166000191681526020019081526020016000206000506040518060800160405290816000820160009054906101000a900460ff161515151581526020016000820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160005080546115f690613410565b80601f016020809104026020016040519081016040528092919081815260200182805461162290613410565b801561166f5780601f106116445761010080835404028352916020019161166f565b820191906000526020600020905b81548152906001019060200180831161165257829003601f168201915b5050505050815260200160028201600050805461168b90613410565b80601f01602080910402602001604051908101604052809291908181526020018280546116b790613410565b80156117045780601f106116d957610100808354040283529160200191611704565b820191906000526020600020905b8154815290600101906020018083116116e757829003601f168201915b505050505081526020015050905080600001511515611758576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161174f90612fed565b60405180910390fd5b60003073ffffffffffffffffffffffffffffffffffffffff1663b8e72af685856040518363ffffffff1660e01b8152600401611795929190612f0e565b60206040518083038186803b1580156117ae5760006000fd5b505afa1580156117c3573d600060003e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117e791906125c7565b905080846040516020016117fb9190612d19565b60405160208183030381529060405290151561184d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016118449190612f46565b60405180910390fd5b50600082604001516040518060400160405280601781526020017f28627974657333322c737472696e672c737472696e672900000000000000000081526020015060405160200161189f929190612c46565b6040516020818303038152906040528686866040516024016118c393929190612e35565b604051602081830303815290604052906040516118e09190612c6b565b60405180910390207bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905060006000846020015173ffffffffffffffffffffffffffffffffffffffff16836040516119649190612c2e565b600060405180830381855afa9150503d806000811461199f576040519150601f19603f3d011682016040523d82523d6000602084013e6119a4565b606091505b50915091508185604001516040516020016119bf9190612c83565b604051602081830303815290604052901515611a11576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611a089190612f46565b60405180910390fd5b5060006000600083806020019051810190611a2c91906125f2565b925092509250821515611a74576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611a6b90612f8a565b60405180910390fd5b600260005060008c6000191660001916815260200190815260200160002060006000820160006101000a81549060ff02191690556000820160016101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600182016000611adf91906122cc565b600282016000611aef91906122cc565b50508115611b5b577fe1426943f65f9b8b1ae886dcbaedf7a0a2cda97712c2e4d5b7d150f61da365118b6040518060600160405280604081526020016137fb604091398a602001518b6040015185604051611b4e959493929190612dcc565b60405180910390a1611bb3565b7f761fc27afeadfc822438b6a79bd12a170f47c2f965dead06da4128a659c207568b6040518060600160405280604081526020016137fb604091398b8d85604051611baa959493929190612e7b565b60405180910390a15b50505050505050505b505050565b600060024333868686604051602001611bde95949392919061300e565b604051602081830303815290604052604051611bfa9190612c2e565b602060405180830381855afa158015611c18573d600060003e3d6000fd5b5050506040513d601f19601f82011682018060405250810190611c3b9190612688565b9050600060405180608001604052806001151581526020018573ffffffffffffffffffffffffffffffffffffffff1681526020018481526020018681526020015090508060026000506000846000191660001916815260200190815260200160002060005060008201518160000160006101000a81548160ff02191690831515021790555060208201518160000160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816001016000509080519060200190611d2692919061230c565b506060820151816002016000509080519060200190611d4692919061230c565b509050507fe1426943f65f9b8b1ae886dcbaedf7a0a2cda97712c2e4d5b7d150f61da36511826040518060600160405280604081526020016137fb60409139868689604051611d99959493929190612dcc565b60405180910390a150505b505050565b60016000506020528060005260406000206000915090508054611dcb90613410565b80601f0160208091040260200160405190810160405280929190818152602001828054611df790613410565b8015611e445780601f10611e1957610100808354040283529160200191611e44565b820191906000526020600020905b815481529060010190602001808311611e2757829003601f168201915b505050505081565b60026000506020528060005260406000206000915090508060000160009054906101000a900460ff16908060000160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001016000508054611eae90613410565b80601f0160208091040260200160405190810160405280929190818152602001828054611eda90613410565b8015611f275780601f10611efc57610100808354040283529160200191611f27565b820191906000526020600020905b815481529060010190602001808311611f0a57829003601f168201915b505050505090806002016000508054611f3f90613410565b80601f0160208091040260200160405190810160405280929190818152602001828054611f6b90613410565b8015611fb85780601f10611f8d57610100808354040283529160200191611fb8565b820191906000526020600020905b815481529060010190602001808311611f9b57829003601f168201915b5050505050905084565b6040518060600160405280604081526020016137fb6040913981565b6000600060006040518060400160405280600981526020017f22566b546167223a22000000000000000000000000000000000000000000000081526020015090506000600090506000600090505b8686905081101561217d578282815181101515612072577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602001015160f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191687878381811015156120da577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916141561214f57818061211790613477565b925050825182141561214a576001600184518361213491906132fc565b61213e91906131a2565b94509450505050612192565b612169565b6007821461215e576000612161565b60015b60ff16915081505b5b808061217590613477565b91505061202c565b50600060008090509350935050506121925650505b9250929050565b6060602067ffffffffffffffff8111156121dc577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519080825280601f01601f19166020018201604052801561220e5781602001600182028036833780820191505090505b50905080508160208201525b919050565b82805461222b90613410565b90600052602060002090601f01602090048101928261224d5760008555612299565b82601f1061226657805160ff1916838001178555612299565b82800160010185558215612299579182015b828111156122985782518260005090905591602001919060010190612278565b5b5090506122a69190612397565b5090565b6040518060200160405280600190602082028036833780820191505090505090565b5080546122d890613410565b6000825580601f106122ea5750612309565b601f0160209004906000526020600020908101906123089190612397565b5b50565b82805461231890613410565b90600052602060002090601f01602090048101928261233a5760008555612386565b82601f1061235357805160ff1916838001178555612386565b82800160010185558215612386579182015b828111156123855782518260005090905591602001919060010190612365565b5b5090506123939190612397565b5090565b61239c565b808211156123b6576000818150600090555060010161239c565b5090566137f9565b60006123d16123cc846130b3565b61308c565b9050828152602081018484840111156123ea5760006000fd5b6123f58482856133cb565b505b9392505050565b600061241161240c846130e5565b61308c565b90508281526020810184848401111561242a5760006000fd5b6124358482856133cb565b505b9392505050565b600061245161244c846130e5565b61308c565b90508281526020810184848401111561246a5760006000fd5b6124758482856133db565b505b9392505050565b60008135905061248d8161378d565b5b92915050565b6000815190506124a3816137a8565b5b92915050565b6000813590506124b9816137c3565b5b92915050565b6000815190506124cf816137c3565b5b92915050565b6000600083601f84011215156124ec5760006000fd5b8235905067ffffffffffffffff8111156125065760006000fd5b60208301915083600182028301111561251f5760006000fd5b5b9250929050565b600082601f830112151561253b5760006000fd5b813561254b8482602086016123be565b9150505b92915050565b600082601f83011215156125695760006000fd5b81356125798482602086016123fe565b9150505b92915050565b600082601f83011215156125975760006000fd5b81516125a784826020860161243e565b9150505b92915050565b6000813590506125c0816137de565b5b92915050565b6000602082840312156125da5760006000fd5b60006125e884828501612494565b9150505b92915050565b600060006000606084860312156126095760006000fd5b600061261786828701612494565b935050602061262886828701612494565b925050604084015167ffffffffffffffff8111156126465760006000fd5b61265286828701612583565b9150505b9250925092565b6000602082840312156126705760006000fd5b600061267e848285016124aa565b9150505b92915050565b60006020828403121561269b5760006000fd5b60006126a9848285016124c0565b9150505b92915050565b600060006000606084860312156126ca5760006000fd5b60006126d8868287016124aa565b935050602084013567ffffffffffffffff8111156126f65760006000fd5b61270286828701612555565b925050604084013567ffffffffffffffff8111156127205760006000fd5b61272c86828701612555565b9150505b9250925092565b600060006020838503121561274c5760006000fd5b600083013567ffffffffffffffff8111156127675760006000fd5b612773858286016124d6565b92509250505b9250929050565b6000600060006000604085870312156127995760006000fd5b600085013567ffffffffffffffff8111156127b45760006000fd5b6127c0878288016124d6565b9450945050602085013567ffffffffffffffff8111156127e05760006000fd5b6127ec878288016124d6565b92509250505b92959194509250565b600060006000606084860312156128125760006000fd5b600084013567ffffffffffffffff81111561282d5760006000fd5b61283986828701612527565b935050602084013567ffffffffffffffff8111156128575760006000fd5b61286386828701612527565b9250506040612874868287016124aa565b9150505b9250925092565b6000602082840312156128925760006000fd5b600082013567ffffffffffffffff8111156128ad5760006000fd5b6128b984828501612555565b9150505b92915050565b600060006000606084860312156128da5760006000fd5b600084013567ffffffffffffffff8111156128f55760006000fd5b61290186828701612555565b93505060206129128682870161247e565b925050604084013567ffffffffffffffff8111156129305760006000fd5b61293c86828701612555565b9150505b9250925092565b60006020828403121561295a5760006000fd5b6000612968848285016125b1565b9150505b92915050565b61297b81613366565b825250505b565b61298b81613379565b825250505b565b61299b81613386565b825250505b565b60006129ae8385613141565b93506129bb8385846133cb565b82840190505b9392505050565b60006129d382613117565b6129dd818561312f565b93506129ed8185602086016133db565b6129f6816135b9565b84019150505b92915050565b6000612a0d82613117565b612a178185613141565b9350612a278185602086016133db565b8084019150505b92915050565b6000612a3f82613123565b612a49818561314d565b9350612a598185602086016133db565b612a62816135b9565b84019150505b92915050565b6000612a7982613123565b612a83818561315f565b9350612a938185602086016133db565b8084019150505b92915050565b6000612aad60148361314d565b9150612ab8826135cb565b6020820190505b919050565b6000612ad160208361314d565b9150612adc826135f5565b6020820190505b919050565b7f63616c6c206261636b206572726f723a000000000000000000000000000000008152505b565b6000612b1c60048361314d565b9150612b278261361f565b6020820190505b919050565b6000612b40602483613141565b9150612b4b82613649565b6024820190505b919050565b6000612b64602983613141565b9150612b6f82613699565b6029820190505b919050565b6000612b8860198361314d565b9150612b93826136e9565b6020820190505b919050565b6000612bac602583613141565b9150612bb782613713565b6025820190505b919050565b7f7665726966792070726f6f66206572726f723a000000000000000000000000008152505b565b6000612bf760158361314d565b9150612c0282613763565b6020820190505b919050565b612c17816133b2565b825250505b565b612c27816133bd565b825250505b565b6000612c3a8284612a02565b91508190505b92915050565b6000612c528285612a02565b9150612c5e8284612a02565b91508190505b9392505050565b6000612c778284612a6e565b91508190505b92915050565b6000612c8e82612ae8565b601082019150612c9e8284612a02565b91508190505b92915050565b6000612cb582612b33565b9150612cc28284866129a2565b91508190505b9392505050565b6000612cda82612b57565b9150612ce78284866129a2565b91508190505b9392505050565b6000612cff82612b9f565b9150612d0c8284866129a2565b91508190505b9392505050565b6000612d2482612bc3565b601382019150612d348284612a02565b91508190505b92915050565b6000602082019050612d556000830184612982565b5b92915050565b6000608082019050612d716000830187612982565b612d7e6020830186612972565b8181036040830152612d908185612a34565b90508181036060830152612da48184612a34565b90505b95945050505050565b6000602082019050612dc56000830184612992565b5b92915050565b600060a082019050612de16000830188612992565b8181036020830152612df38187612a34565b9050612e026040830186612972565b8181036060830152612e148185612a34565b90508181036080830152612e288184612a34565b90505b9695505050505050565b6000606082019050612e4a6000830186612992565b8181036020830152612e5c8185612a34565b90508181036040830152612e708184612a34565b90505b949350505050565b600060a082019050612e906000830188612992565b8181036020830152612ea28187612a34565b90508181036040830152612eb68186612a34565b90508181036060830152612eca8185612a34565b90508181036080830152612ede8184612a34565b90505b9695505050505050565b60006020820190508181036000830152612f0581846129c8565b90505b92915050565b60006040820190508181036000830152612f2881856129c8565b90508181036020830152612f3c81846129c8565b90505b9392505050565b60006020820190508181036000830152612f608184612a34565b90505b92915050565b60006020820190508181036000830152612f8281612aa0565b90505b919050565b60006020820190508181036000830152612fa381612ac4565b90505b919050565b60006020820190508181036000830152612fc481612b0f565b90505b919050565b60006020820190508181036000830152612fe581612b7b565b90505b919050565b6000602082019050818103600083015261300681612bea565b90505b919050565b600060a0820190506130236000830188612c0e565b6130306020830187612972565b81810360408301526130428186612a34565b90506130516060830185612972565b81810360808301526130638184612a34565b90505b9695505050505050565b60006020820190506130856000830184612c1e565b5b92915050565b60006130966130a8565b90506130a28282613445565b5b919050565b600060405190505b90565b600067ffffffffffffffff8211156130ce576130cd613588565b5b6130d7826135b9565b90506020810190505b919050565b600067ffffffffffffffff821115613100576130ff613588565b5b613109826135b9565b90506020810190505b919050565b6000815190505b919050565b6000815190505b919050565b60008282526020820190505b92915050565b60008190505b92915050565b60008282526020820190505b92915050565b60008190505b92915050565b600060008585111561317d5760006000fd5b8386111561318b5760006000fd5b600185028301915084860390505b94509492505050565b60006131ad826133b2565b91506131b8836133b2565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156131ed576131ec6134f5565b5b82820190505b92915050565b6000613204826133bd565b915061320f836133bd565b92508260ff03821115613225576132246134f5565b5b82820190505b92915050565b600061323c826133b2565b9150613247836133b2565b925082151561325957613258613526565b5b82820490505b92915050565b6000613270826133b2565b915061327b836133b2565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156132b4576132b36134f5565b5b82820290505b92915050565b60006132cb826133bd565b91506132d6836133bd565b92508160ff04831182151516156132f0576132ef6134f5565b5b82820290505b92915050565b6000613307826133b2565b9150613312836133b2565b925082821015613325576133246134f5565b5b82820390505b92915050565b600061333c826133bd565b9150613347836133bd565b92508282101561335a576133596134f5565b5b82820390505b92915050565b600061337182613391565b90505b919050565b600081151590505b919050565b60008190505b919050565b600073ffffffffffffffffffffffffffffffffffffffff821690505b919050565b60008190505b919050565b600060ff821690505b919050565b828183376000838301525050505b565b60005b838110156133fa5780820151818401525b6020810190506133de565b83811115613409576000848401525b505050505b565b60006002820490506001821680151561342a57607f821691505b6020821081141561343e5761343d613557565b5b505b919050565b61344e826135b9565b810181811067ffffffffffffffff8211171561346d5761346c613588565b5b806040525050505b565b6000613482826133b2565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156134b5576134b46134f5565b5b6001820190505b919050565b60006134cc826133b2565b91506134d7836133b2565b92508215156134e9576134e8613526565b5b82820690505b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b565b6000601f19601f83011690505b919050565b7f7665726966792d6f70636f64652d6661696c65640000000000000000000000006000820152505b565b7f63616c6c206261636b2066756e6374696f6e2072657475726e73206572726f726000820152505b565b7f6661696c000000000000000000000000000000000000000000000000000000006000820152505b565b7f7665726966792070726f6f66206661696c3a2067657420766b5461672065727260008201527f6f723a20000000000000000000000000000000000000000000000000000000006020820152505b565b7f7665726966792070726f6f66206661696c3a2067657420766b20627920766b5460008201527f6167206572726f723a00000000000000000000000000000000000000000000006020820152505b565b7f766b54616720737472696e67206c656e677468206572726f72000000000000006000820152505b565b7f7665726966792070726f6f66206661696c3a2066696e6420766b54616720657260008201527f726f723a200000000000000000000000000000000000000000000000000000006020820152505b565b7f63616e27742066696e64207461736b20627920494400000000000000000000006000820152505b565b61379681613366565b811415156137a45760006000fd5b505b565b6137b181613379565b811415156137bf5760006000fd5b505b565b6137cc81613386565b811415156137da5760006000fd5b505b565b6137e7816133bd565b811415156137f55760006000fd5b505b565bfe38663030333135636633326439323832656564316662646165646432336138663238356261613038626131633834313734366161613532323939376635323838a2646970667358221220387ce93f200e40eef64e0fef693e2202327bfb75c39f117e6da9d9337375114464736f6c6343000804003330313030303030303637373236663734363833313336303030303030303030353632366533323335333430303030303030343030303030303033303030303031383030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303130303030303038303139386539333933393230643438336137323630626662373331666235643235663161613439333333356139653731323937653438356237616566333132633231383030646565663132316631653736343236613030363635653563343437393637343332326434663735656461646434366465626435636439393266366564303930363839643035383566663037356563396539396164363930633333393562633462333133333730623338656633353561636461646364313232393735623132633835656135646238633664656234616162373138303864636234303866653364316537363930633433643337623463653663633031363666613764616130303030303038303139386539333933393230643438336137323630626662373331666235643235663161613439333333356139653731323937653438356237616566333132633231383030646565663132316631653736343236613030363635653563343437393637343332326434663735656461646434366465626435636439393266366564303930363839643035383566663037356563396539396164363930633333393562633462333133333730623338656633353561636461646364313232393735623132633835656135646238633664656234616162373138303864636234303866653364316537363930633433643337623463653663633031363666613764616130303030303034303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030343030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303430303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030343030303030303033303030303030303430303030303030303030303030303034303030303030303130303030303034303635363233333333363436313334333433323333333033323632333236333334333536333330333933343339363133313339333136363337363536313336363533373339333833393338363533333332363533383336333636343635333236313632333233383338363536353332333133333332333433393634363333383336303030303030303536323665333233353334303030303031386633303832303138623031303166663034303536323665333233353334303230313031303432303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303130343230303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030313034323030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303031303432303330363434653732653133316130323962383530343562363831383135383564323833336538343837396239373039313433653166353933663030303030303030343230333036343465373265313331613032396238353034356236383138313538356432383333653834383739623937303931343365316635393366303030303030303330323230343230303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030313330323230343230303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030313330343430343230303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030313034323033303634346537326531333161303239623835303435623638313831353835643238333365383438373962393730393134336531663539336630303030303030333034343034323030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303031303432303330363434653732653133316130323962383530343562363831383135383564323833336538343837396239373039313433653166353933663030303030303030303030303034303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030343030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303430303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030313030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303230303030303038303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303038303139386539333933393230643438336137323630626662373331666235643235663161613439333333356139653731323937653438356237616566333132633231383030646565663132316631653736343236613030363635653563343437393637343332326434663735656461646434366465626435636439393266366564303930363839643035383566663037356563396539396164363930633333393562633462333133333730623338656633353561636461646364313232393735623132633835656135646238633664656234616162373138303864636234303866653364316537363930633433643337623463653663633031363666613764616130303030303034303264393662313231343836616239646137626635343965353764326638613663633139383361333336393033353234666230356463643530373435376636336331323939303866666337623764356633643837316663313230613965653462626535623762353633323961376137393235396137343637646237613235353634303030303030343030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303031303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030323030303030303830313938653933393339323064343833613732363062666237333166623564323566316161343933333335613965373132393765343835623761656633313263323138303064656566313231663165373634323661303036363565356334343739363734333232643466373565646164643436646562643563643939326636656430393036383964303538356666303735656339653939616436393063333339356263346233313333373062333865663335356163646164636431323239373562313263383565613564623863366465623461616237313830386463623430386665336431653736393063343364333762346365366363303136366661376461613030303030303430303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030313030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303230303030303034303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030383030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030343030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303031333036343465373265313331613032396238353034356236383138313538356439373831366139313638373163613864336332303863313664383763666434353030303030303430303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303038303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303034303036613762363461663866343134626362656566343535623164613532303863396235393262383365653635393938323463616136643265653931343161373630386537346534333863656533316163313034636535396239346534356665393861393764386638613665373536363463653838656635613431653732666263`
		contractAddress, err := deployContract(bin, abiStr)
		t.Log(contractAddress)
		assert.Nil(t, err) //0x244f3c929900517497f935160a72a6656a1332cc
	})

	abiStr := `[{"inputs":[{"internalType":"bytes32","name":"taskID","type":"bytes32"},{"internalType":"string","name":"result","type":"string"},{"internalType":"string","name":"proof","type":"string"}],"name":"callback","outputs":[{"internalType":"bool","name":"","type":"bool"},{"internalType":"bool","name":"","type":"bool"},{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"proxy","type":"address"}],"name":"setAddr","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"input","type":"string"}],"name":"trigger","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"int256","name":"","type":"int256"}],"name":"triggerInfo","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function"}]`
	bin := `60806040523480156100115760006000fd5b50610017565b610d7c806100266000396000f3fe60806040523480156100115760006000fd5b50600436106100515760003560e01c80638c3884621461005757806399eaa6a114610089578063c603761c146100b9578063d1d80fdf146100d557610051565b60006000fd5b610071600480360381019061006c9190610629565b6100f1565b6040516100809392919061088a565b60405180910390f35b6100a3600480360381019061009e91906106ad565b6101f6565b6040516100b091906108c9565b60405180910390f35b6100d360048036038101906100ce91906106d8565b610299565b005b6100ef60048036038101906100ea91906105fe565b610457565b005b600060006060600060001b866000191614151515610144576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161013b90610930565b60405180910390fd5b6000859050600085905060008251118015610160575060008151115b15156101a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161019890610951565b60405180910390fd5b600160006040518060400160405280601481526020017f6170702063616c6c6261636b207375636365737300000000000000000000000081526020015094509450945050506101ed5650505b93509350939050565b6000600050602052806000526040600020600091509050805461021890610af6565b80601f016020809104026020016040519081016040528092919081815260200182805461024490610af6565b80156102915780601f1061026657610100808354040283529160200191610291565b820191906000526020600020905b81548152906001019060200180831161027457829003601f168201915b505050505081565b600081306040516024016102ae9291906108ec565b6040516020818303038152906040527fe84df22c000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905060006000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16836040516103789190610872565b6000604051808303816000865af19150503d80600081146103b5576040519150601f19603f3d011682016040523d82523d6000602084013e6103ba565b606091505b5091509150811515610401576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103f890610972565b60405180910390fd5b60016000818150548092919061041690610b5d565b9190509090555080600060005060006001600050548152602001908152602001600020600050908051906020019061044f92919061049c565b505050505b50565b80600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b50565b8280546104a890610af6565b90600052602060002090601f0160209004810192826104ca5760008555610516565b82601f106104e357805160ff1916838001178555610516565b82800160010185558215610516579182015b8281111561051557825182600050909055916020019190600101906104f5565b5b5090506105239190610527565b5090565b61052c565b80821115610546576000818150600090555060010161052c565b509056610d45565b600061056161055c846109ba565b610993565b90508281526020810184848401111561057a5760006000fd5b610585848285610ab1565b505b9392505050565b60008135905061059d81610cf4565b5b92915050565b6000813590506105b381610d0f565b5b92915050565b6000813590506105c981610d2a565b5b92915050565b600082601f83011215156105e45760006000fd5b81356105f484826020860161054e565b9150505b92915050565b6000602082840312156106115760006000fd5b600061061f8482850161058e565b9150505b92915050565b600060006000606084860312156106405760006000fd5b600061064e868287016105a4565b935050602084013567ffffffffffffffff81111561066c5760006000fd5b610678868287016105d0565b925050604084013567ffffffffffffffff8111156106965760006000fd5b6106a2868287016105d0565b9150505b9250925092565b6000602082840312156106c05760006000fd5b60006106ce848285016105ba565b9150505b92915050565b6000602082840312156106eb5760006000fd5b600082013567ffffffffffffffff8111156107065760006000fd5b610712848285016105d0565b9150505b92915050565b61072581610a47565b825250505b565b6000610737826109ec565b6107418185610a04565b9350610751818560208601610ac1565b61075a81610c3a565b84019150505b92915050565b6000610771826109ec565b61077b8185610a16565b935061078b818560208601610ac1565b8084019150505b92915050565b6107a181610a8b565b825250505b565b60006107b3826109f8565b6107bd8185610a22565b93506107cd818560208601610ac1565b6107d681610c3a565b84019150505b92915050565b60006107ef600683610a22565b91506107fa82610c4c565b6020820190505b919050565b6000610813600683610a22565b915061081e82610c76565b6020820190505b919050565b6000610837600f83610a22565b915061084282610ca0565b6020820190505b919050565b600061085b600883610a22565b915061086682610cca565b6020820190505b919050565b600061087e8284610766565b91508190505b92915050565b600060608201905061089f600083018661071c565b6108ac602083018561071c565b81810360408301526108be81846107a8565b90505b949350505050565b600060208201905081810360008301526108e3818461072c565b90505b92915050565b6000606082019050818103600083015261090681856107a8565b90506109156020830184610798565b81810360408301526109268161084e565b90505b9392505050565b60006020820190508181036000830152610949816107e2565b90505b919050565b6000602082019050818103600083015261096a81610806565b90505b919050565b6000602082019050818103600083015261098b8161082a565b90505b919050565b600061099d6109af565b90506109a98282610b2b565b5b919050565b600060405190505b90565b600067ffffffffffffffff8211156109d5576109d4610c09565b5b6109de82610c3a565b90506020810190505b919050565b6000815190505b919050565b6000815190505b919050565b60008282526020820190505b92915050565b60008190505b92915050565b60008282526020820190505b92915050565b6000610a3f82610a6a565b90505b919050565b600081151590505b919050565b60008190505b919050565b60008190505b919050565b600073ffffffffffffffffffffffffffffffffffffffff821690505b919050565b6000610a9682610a9e565b90505b919050565b6000610aa982610a6a565b90505b919050565b828183376000838301525050505b565b60005b83811015610ae05780820151818401525b602081019050610ac4565b83811115610aef576000848401525b505050505b565b600060028204905060018216801515610b1057607f821691505b60208210811415610b2457610b23610bd8565b5b505b919050565b610b3482610c3a565b810181811067ffffffffffffffff82111715610b5357610b52610c09565b5b806040525050505b565b6000610b6882610a5f565b91507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610b9b57610b9a610ba7565b5b6001820190505b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b565b6000601f19601f83011690505b919050565b7f7461736b494400000000000000000000000000000000000000000000000000006000820152505b565b7f6c656e67746800000000000000000000000000000000000000000000000000006000820152505b565b7f63616c6c206261636b206572726f7200000000000000000000000000000000006000820152505b565b7f63616c6c6261636b0000000000000000000000000000000000000000000000006000820152505b565b610cfd81610a34565b81141515610d0b5760006000fd5b505b565b610d1881610a54565b81141515610d265760006000fd5b505b565b610d3381610a5f565b81141515610d415760006000fd5b505b565bfea2646970667358221220b5831e1969d9494378cda42ae92cfa78802022bddf007cf933479f00cd35e96f64736f6c63430008040033`
	t.Run("app", func(t *testing.T) {
		appAddress, err := deployContract(bin, abiStr)
		t.Log(appAddress) //0x31cf62472b1856d94553d2fe78f3bb067afb0714
		assert.Nil(t, err)

		ABIBefore, serr := abi.JSON(strings.NewReader(abiStr))
		assert.Nil(t, serr)
		addr, _ := hex.DecodeString("244f3c929900517497f935160a72a6656a1332cc")
		packed, serr := ABIBefore.Pack("setAddr", common.BytesToAddress(addr))
		if serr != nil {
			t.Error(serr)
			return
		}

		r1Key, _ := asym.GenerateKey(asym.AlgoP256R1)
		pubKey := &account.ECDSAKey{ECDSAPrivateKey: r1Key}
		newAddress := pubKey.GetAddress()
		transaction := NewTransaction(newAddress.Hex()).Invoke(appAddress, packed)
		ret, err := rp.SignAndInvokeContract(transaction, pubKey)
		assert.Nil(t, err)
		t.Log(ret.Ret)
	})

	t.Run("trigger", func(t *testing.T) {
		ABIBefore, serr := abi.JSON(strings.NewReader(abiStr))
		assert.Nil(t, serr)
		packed, serr := ABIBefore.Pack("trigger", `{"in":4}`)
		if serr != nil {
			t.Error(serr)
			return
		}

		r1Key, _ := asym.GenerateKey(asym.AlgoP256R1)
		pubKey := &account.ECDSAKey{ECDSAPrivateKey: r1Key}
		newAddress := pubKey.GetAddress()
		transaction := NewTransaction(newAddress.Hex()).Invoke( /*app contract*/ "0x31cf62472b1856d94553d2fe78f3bb067afb0714", packed)
		ret, err := rp.SignAndInvokeContract(transaction, pubKey)
		assert.Nil(t, err)
		t.Log(ret.Ret)
	})

}

func TestRPC_SendED25519Tx(t *testing.T) {
	address, _ := testPrivateAccount()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	accountJSON, _ := account.NewAccountED25519("12345678")
	t.Log("account", accountJSON)
	ekey, err := account.GenKeyFromAccountJson(accountJSON, "12345678")
	assert.Nil(t, err)
	newAddress := ekey.(*account.ED25519Key).GetAddress()

	transaction := NewTransaction(newAddress.Hex()).Transfer(address, int64(0))
	receipt, err := rp.SignAndSendTx(transaction, ekey)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, err)
	fmt.Println(receipt.Ret)
}

func TestRPC_SendTx_With_ExtraID(t *testing.T) {
	guomiKey := testGuomiKey()
	address, _ := testPrivateAccount()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	addr := guomiKey.GetAddress()
	transaction := NewTransaction(hex.EncodeToString(addr[:])).Transfer(address, int64(0))
	transaction.SetExtraIDInt64(123, 456)
	transaction.SetExtraIDString("abc")
	receipt, err := rp.SignAndSendTx(transaction, guomiKey)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(receipt.Ret)
}

func TestRPC_SM2Account(t *testing.T) {
	//accountJson, _ := account.NewAccountSm2("")
	strAcc := `{"address":"0x8485147cbf02dec93ee84f81824a3b60e355f5cd","publicKey":"04a1b4c82a2a13e15a11e3ee9316504de0c3b54d46f5c189ae42603c9cd07a50fdca2ac35d0ceef4a8466ccb182f52403d9a58b573e1bf6fd4f52c31493bf7241b","privateKey":"f67136bf3caa4197a1cfaf38a5392ff94dae91bda700f8898b11cf49891a47bb","privateKeyEncrypted":false}`
	key, _ := account.NewAccountSm2FromAccountJSON(strAcc, "")
	//fmt.Println(accountJson)
	//key, _ := account.NewAccountSm2FromAccountJSON(accountJson, "")
	pubKey, _ := key.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	address, _ := testPrivateAccount()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Transfer(address, int64(0))
	receipt, err := rp.SignAndSendTx(transaction, key)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, 66, len(receipt.TxHash))

	accountJSON, _ := account.NewAccountSm2("12345678")
	aKey, syserr := account.NewAccountSm2FromAccountJSON(accountJSON, "12345678")
	if syserr != nil {
		t.Error(syserr)
	}
	newAddress2 := aKey.GetAddress()

	transaction1 := NewTransaction(newAddress2.Hex()).Transfer(address, int64(0))
	receipt1, err := rp.SignAndSendTx(transaction1, aKey)
	if err != nil {
		t.Error(err)
		return
	}
	//assert.EqualValues(t, transaction1.GetTransactionHash(DefaultTxGasLimit), receipt1.TxHash)
	assert.EqualValues(t, 66, len(receipt1.TxHash))
}

// maintain contract by opcode 1
func TestRPC_MaintainContract(t *testing.T) {
	t.Skip("solc")
	contractOriginFile := "../conf/contract/Accumulator.sol"
	contractUpdateFile := "../conf/contract/AccumulatorUpdate.sol"
	compileOrigin, _ := compileContract(contractOriginFile)
	compileUpdate, _ := compileContract(contractUpdateFile)
	contractAddress, err := deployContract(compileOrigin.Bin[0], compileOrigin.Abi[0])
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("contractAddress:", contractAddress)

	// test invoke before update
	ABIBefore, serr := abi.JSON(strings.NewReader(compileOrigin.Abi[0]))
	assert.Nil(t, serr)
	packed, serr := ABIBefore.Pack("add", uint32(11), uint32(1))
	if serr != nil {
		t.Error(serr)
		return
	}
	guomiKey := testGuomiKey()
	rp, serr := NewJsonRPC()
	assert.Nil(t, serr)
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	transactionInvokeBe := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
	receiptBe, err := rp.SignAndInvokeContract(transactionInvokeBe, guomiKey)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(receiptBe.Ret)

	var result1 uint32
	decode(ABIBefore, &result1, "add", receiptBe.Ret)
	fmt.Println(result1)

	fmt.Println("-----------------------------")

	transactionUpdate := NewTransaction(common.BytesToAddress(newAddress).Hex()).Maintain(1, contractAddress, compileUpdate.Bin[0])
	//transactionUpdate, err := NewMaintainTransaction(guomiKey.GetAddress(), contractAddress, compileUpdate.Bin[0], 1, EVM)
	transactionUpdate.Sign(guomiKey)
	receiptUpdate, err := rp.MaintainContract(transactionUpdate)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(receiptUpdate.ContractAddress)

	// test invoke after update
	ABI, serr := abi.JSON(strings.NewReader(compileUpdate.Abi[0]))
	if serr != nil {
		t.Error(err)
		return
	}
	packed2, serr := ABI.Pack("addUpdate", uint32(1), uint32(2))
	assert.Nil(t, serr)
	transactionInvoke := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed2)
	receiptInvoke, err := rp.SignAndInvokeContract(transactionInvoke, guomiKey)
	assert.Nil(t, err)
	t.Log(receiptInvoke.Ret)
	var result2 uint32
	decode(ABI, &result2, "addUpdate", receiptInvoke.Ret)
	fmt.Println(result2)
}

// maintain contract by opcode 2 and 3
func TestRPC_MaintainContract2(t *testing.T) {
	contractAddress, _ := deployContract(binContract, abiContract)
	ABI, _ := abi.JSON(strings.NewReader(abiContract))
	// invoke first
	guomiKey := testGuomiKey()
	packed, _ := ABI.Pack("getSum")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	transaction1 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
	receipt1, err := rp.SignAndInvokeContract(transaction1, guomiKey)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("invoke first:", receipt1.Ret)

	// freeze contract
	transactionFreeze := NewTransaction(common.BytesToAddress(newAddress).Hex()).Maintain(2, contractAddress, "")
	//transactionFreeze, _ := NewMaintainTransaction(guomiKey.GetAddress(), contractAddress, "", 2, EVM)
	transactionFreeze.Sign(guomiKey)
	receiptFreeze, err := rp.MaintainContract(transactionFreeze)
	assert.Nil(t, err)
	fmt.Println(receiptFreeze.TxHash)
	status, err := rp.GetContractStatus(contractAddress)
	assert.Nil(t, err)
	fmt.Println("contract status >>", status)

	// invoke after freeze
	transaction2 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
	receipt2, err := rp.SignAndInvokeContract(transaction2, guomiKey)
	if err != nil {
		fmt.Println("invoke second receipt2 is null ", receipt2 == nil)
		fmt.Println(err)
	}

	// unfreeze contract
	transactionUnfreeze := NewTransaction(common.BytesToAddress(newAddress).Hex()).Maintain(3, contractAddress, "")
	//transactionUnfreeze, _ := NewMaintainTransaction(guomiKey.GetAddress(), contractAddress, "", 3, EVM)
	transactionUnfreeze.Sign(guomiKey)
	receiptUnFreeze, err := rp.MaintainContract(transactionUnfreeze)
	assert.Nil(t, err)
	fmt.Println(receiptUnFreeze.TxHash)
	status, _ = rp.GetContractStatus(contractAddress)
	fmt.Println("contract status >>", status)

	// invoke after unfreeze
	transaction3 := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
	receipt3, err := rp.SignAndInvokeContract(transaction3, guomiKey)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("invoke third:", receipt3.Ret)
}

func TestRPC_GetContractStatus(t *testing.T) {
	t.Skip("the node can get the account")
	contractAddress, _ := deployContract(binContract, abiContract)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	statu, err := rp.GetContractStatus(contractAddress)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(statu)
}

func TestRPC_GetContractStatusByName(t *testing.T) {
	t.Skip("set contract name `HashContract`")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	dateTime, err := rp.GetContractStatusByName("HashContract")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(dateTime)
}

func TestRPC_GetCreator(t *testing.T) {
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	receipt, _ := rp.SignAndDeployContract(transaction, guomiKey)

	accountAddress, err := rp.GetCreator(receipt.ContractAddress)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(accountAddress)
}

func TestRPC_GetCreatorByName(t *testing.T) {
	t.Skip("set contract name `HashContract`")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	dateTime, err := rp.GetCreatorByName("HashContract")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(dateTime)
}

func TestRPC_GetCreateTime(t *testing.T) {
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract)
	receipt, _ := rp.SignAndDeployContract(transaction, guomiKey)

	dateTime, err := rp.GetCreateTime(receipt.ContractAddress)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(dateTime)
}

func TestRPC_GetCreateTimeByName(t *testing.T) {
	t.Skip("set contract name `HashContract`")
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	dateTime, err := rp.GetCreateTimeByName("HashContract")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(dateTime)
}

func TestRPC_GetDeployedList(t *testing.T) {
	guomiKey := testGuomiKey()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	list, err := rp.GetDeployedList(common.BytesToAddress(newAddress).Hex())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(list))
}

func TestRPC_InvokeContractReturnHash(t *testing.T) {
	t.Skip("pressure test, do not put this test in CI")
	cr, _ := compileContract("../conf/contract/Accumulator.sol")
	contractAddress, err := deployContract(cr.Bin[0], cr.Abi[0])
	ABI, serr := abi.JSON(strings.NewReader(cr.Abi[0]))
	if err != nil {
		t.Error(serr)
		return
	}
	packed, serr := ABI.Pack("add", uint32(1), uint32(2))
	if serr != nil {
		t.Error(serr)
		return
	}
	address, privateKey := testPrivateAccount()
	rp, serr := NewJsonRPC()
	assert.Nil(t, serr)
	transaction := NewTransaction(address).Invoke(contractAddress, packed)
	transaction.Sign(privateKey)
	tt := time.After(1 * time.Minute)
	counter := 0
	for {
		_, err = rp.InvokeContractReturnHash(transaction)
		if err != nil {
			t.Error(err)
			return
		}
		select {
		case <-tt:
			fmt.Println(counter)
			return
		default:
			counter++
		}
	}
}

func TestRPC_InvokeContract2(t *testing.T) {
	to := `0x12345678901234567890123456789012345678901234567890`
	rawABI := `[{"constant":false,"inputs":[{"name":"num1","type":"uint32"},{"name":"num2","type":"uint32"}],"name":"add","outputs":[{"name":"","type":"uint32"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"getSum","outputs":[{"name":"","type":"uint32"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getHello","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"increment","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"addr1","type":"address"},{"indexed":false,"name":"msg","type":"bytes32"}],"name":"sayHello","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"msg","type":"bytes32"},{"indexed":false,"name":"sum","type":"uint32"}],"name":"saySum","type":"event"}]`

	address, _ := testPrivateAccount()
	t.Run("normal input", func(t *testing.T) {
		tx1 := NewTransaction(address).InvokeContract(to, rawABI, "add", "111", "111")

		ABI, err := abi.JSON(strings.NewReader(rawABI))
		assert.NoError(t, err)
		payload, err := ABI.Pack("add", uint32(111), uint32(111))
		tx2 := NewTransaction(address).Invoke(to, payload)
		assert.NoError(t, err)
		assert.Equal(t, tx2.payload, tx1.payload)
	})

	t.Run("error input", func(t *testing.T) {
		errTx := NewTransaction(address).InvokeContract(to, strings.Replace(rawABI, "[", "{", -1), "add", "111", "111")
		assert.Nil(t, errTx)

		errTx = NewTransaction(address).InvokeContract(to, rawABI, "123", "111", "111")
		assert.Nil(t, errTx)
	})
}
