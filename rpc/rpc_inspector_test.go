package rpc

import (
	"fmt"
	"github.com/jackzing/gosdk/account"
	"github.com/jackzing/gosdk/bvm"
	"github.com/jackzing/gosdk/common"
	"github.com/jackzing/gosdk/hvm"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRPC_InspectorRole(t *testing.T) {
	t.Skip()
	//newAccount, err := account.NewAccount("")
	//assert.Nil(t, err)
	//t.Log(newAccount)
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	defaultPwd := ""
	newAccount := `{"address":"0x430d637bb98033e1a2406ab8c24949e142fa4e75","algo":"0x03","version":"4.0","publicKey":"0x0454bbfe47a39685ac16d8263b0f7e3e0406996dc218543399411dfff12e304768945d0afe07c6b932fb203d1fb670fc46839adce54efed8e153fd21779c29d02b","privateKey":"e549225831723fdf9a7d4530b17b746e2c9da37fc1a982f4927f94d2b2972ac2"}`
	ac, err := account.NewAccountFromAccountJSON(newAccount, defaultPwd)
	assert.Nil(t, err)

	role := "user"
	role1 := "account"
	t.Run("add_role_success", func(t *testing.T) {
		stdError := rp.AddRoleForNode(ac.GetAddress().Hex(), role, role1)
		assert.Nil(t, stdError)
	})

	t.Run("add_role_already_exist", func(t *testing.T) {
		stdError := rp.AddRoleForNode(ac.GetAddress().Hex(), role)
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address:%s already has roles:[%s]", ac.GetAddress().Hex(), role)))
	})

	t.Run("delete_role_success", func(t *testing.T) {
		stdError := rp.DeleteRoleFromNode(ac.GetAddress().Hex(), role1)
		assert.Nil(t, stdError)
	})

	t.Run("delete_role_success", func(t *testing.T) {
		stdError := rp.DeleteRoleFromNode(ac.GetAddress().Hex(), role1)
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address:%s has not roles:[%s]", ac.GetAddress().Hex(), role1)))
	})

	t.Run("get_role", func(t *testing.T) {
		roles, err := rp.GetRoleFromNode(ac.GetAddress().Hex())
		assert.Nil(t, err)
		assert.Equal(t, roles, []string{role})
	})

	t.Run("get_all_role", func(t *testing.T) {
		roles, err := rp.GetAllRolesFromNode()
		assert.Nil(t, err)
		assert.NotNil(t, roles)
	})

	// inspector enable
	t.Run("auth_forbidden_and_authorized_and_methods", func(t *testing.T) {
		// set rules
		user := "user"
		admin := "admin"
		rule := &InspectorRule{
			AllowAnyone:     false,
			AuthorizedRoles: []string{admin},
			ForbiddenRoles:  []string{user},
			ID:              0,
			Method:          []string{"node_*"},
			Name:            "node",
		}
		stdError := rp.SetRulesInNode([]*InspectorRule{rule})
		assert.Nil(t, stdError)

		// add user to account A
		accountA, _ := account.NewAccount(defaultPwd)
		acA, _ := account.NewAccountFromAccountJSON(accountA, defaultPwd)
		stdError = rp.AddRoleForNode(acA.GetAddress().Hex(), user)
		assert.Nil(t, stdError)
		// use account A call node_getNodeStates has not permission
		rp.SetAccount(acA)
		_, stdError = rp.GetNodeStates()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodeStates", acA.GetAddress().Hex())))
		// use account A call node_getNodes has not permission
		rp.SetAccount(acA)
		_, stdError = rp.GetNodes()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodes", acA.GetAddress().Hex())))

		// add admin to account B
		accountB, _ := account.NewAccount(defaultPwd)
		acB, _ := account.NewAccountFromAccountJSON(accountB, defaultPwd)
		stdError = rp.AddRoleForNode(acB.GetAddress().Hex(), admin)
		assert.Nil(t, stdError)
		// use account B call node_getNodeStates has permission
		rp.SetAccount(acB)
		nodeState, stdError := rp.GetNodeStates()
		assert.Nil(t, stdError)
		assert.NotNil(t, nodeState)
		t.Log(nodeState)
	})

	t.Run("auth_allow_anyone_false", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = [""]
		//id = 0
		//methods = ["node_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rule := &InspectorRule{
			AllowAnyone:     false,
			AuthorizedRoles: []string{admin},
			ForbiddenRoles:  []string{},
			ID:              0,
			Method:          []string{"node_*"},
			Name:            "node",
		}
		stdError := rp.SetRulesInNode([]*InspectorRule{rule})
		assert.Nil(t, stdError)

		// add user to account B
		accountB, _ := account.NewAccount(defaultPwd)
		acB, _ := account.NewAccountFromAccountJSON(accountB, defaultPwd)
		stdError = rp.AddRoleForNode(acB.GetAddress().Hex(), user)
		assert.Nil(t, stdError)
		// use account B call node_getNodeStates has not permission
		rp.SetAccount(acB)
		_, stdError = rp.GetNodeStates()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodeStates", acB.GetAddress().Hex())))
	})

	t.Run("auth_rules_id", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = [""]
		//id = 0
		//methods = ["node_*"]
		//name = "node"
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = ["user"]
		//id = 1
		//methods = ["node_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ID:              0,
				Method:          []string{"node_*"},
				Name:            "node",
			},
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"node_*"},
				Name:            "node",
			},
		}
		stdError := rp.SetRulesInNode(rules)
		assert.Nil(t, stdError)
		inspectorRules, stdError := rp.GetRulesFromNode()
		assert.Nil(t, stdError)
		assert.Len(t, inspectorRules, len(rules))

		// add user to account C
		accountC, _ := account.NewAccount(defaultPwd)
		acC, _ := account.NewAccountFromAccountJSON(accountC, defaultPwd)
		stdError = rp.AddRoleForNode(acC.GetAddress().Hex(), user)
		assert.Nil(t, stdError)
		// use account C call node_getNodeStates has not permission
		rp.SetAccount(acC)
		_, stdError = rp.GetNodeStates()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodeStates", acC.GetAddress().Hex())))

	})

	t.Run("auth_priority", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = true
		//authorized_roles = ["admin"]
		//forbidden_roles = ["user"]
		//id = 1
		//methods = ["node_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     true,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"node_*"},
				Name:            "node",
			},
		}
		stdError := rp.SetRulesInNode(rules)
		assert.Nil(t, stdError)

		// add user to account D
		accountD, _ := account.NewAccount(defaultPwd)
		acD, _ := account.NewAccountFromAccountJSON(accountD, defaultPwd)
		stdError = rp.AddRoleForNode(acD.GetAddress().Hex(), user, admin)
		assert.Nil(t, stdError)
		// use account D call node_getNodeStates has not permission
		rp.SetAccount(acD)
		_, stdError = rp.GetNodeStates()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodeStates", acD.GetAddress().Hex())))
	})

	t.Run("auth_transaction", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = ["","user"]
		//id = 1
		//methods = ["tx_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"contract_*"},
				Name:            "node",
			},
		}
		stdError := rp.SetRulesInNode(rules)
		assert.Nil(t, stdError)

		// add user,admin to account H
		accountH, _ := account.NewAccount(defaultPwd)
		acH, _ := account.NewAccountFromAccountJSON(accountH, defaultPwd)
		stdError = rp.AddRoleForNode(acH.GetAddress().Hex(), user, admin)
		assert.Nil(t, stdError)

		// send tx
		transaction := NewTransaction(acH.GetAddress().Hex()).Transfer("bfa5bd992e3eb123c8b86ebe892099d4e9efb783", int64(0))
		transaction.Sign(acH)
		receipt, stdError := rp.SendTx(transaction)
		assert.Nil(t, stdError)
		assert.NotNil(t, receipt)

		// deploy contract
		deployJar, err := DecompressFromJar("../hvmtestfile/fibonacci/fibonacci-1.0-fibonacci.jar")
		assert.Nil(t, err)
		transaction = NewTransaction(acH.GetAddress().Hex()).Deploy(common.Bytes2Hex(deployJar)).VMType(HVM)
		transaction.Sign(acH)
		receipt, err = rp.DeployContract(transaction)
		assert.Nil(t, err)
		assert.NotNil(t, receipt)

		// invoke contract
		abiPath := "../hvmtestfile/fibonacci/hvm.abi"
		abiJson, rerr := common.ReadFileAsString(abiPath)
		assert.Nil(t, rerr)
		abi, gerr := hvm.GenAbi(abiJson)
		assert.Nil(t, gerr)
		easyBean := "invoke.InvokeFibonacci"
		beanAbi, err := abi.GetBeanAbi(easyBean)
		assert.Nil(t, err)
		payload, err := hvm.GenPayload(beanAbi)
		assert.Nil(t, err)
		transaction1 := NewTransaction(acH.GetAddress().Hex()).Invoke(receipt.ContractAddress, payload).VMType(HVM)
		transaction1.Sign(acH)
		_, err = rp.InvokeContract(transaction1)
		assert.Nil(t, err)

		// manager contract by vote
		ope := bvm.NewContractDeployContractOperation([]byte("source"), deployJar, "hvm", nil)
		contractOpt := bvm.NewProposalCreateOperationForContract(ope)
		payload = bvm.EncodeOperation(contractOpt)
		tx := NewTransaction(acH.GetAddress().Hex()).Invoke(contractOpt.Address(), payload).VMType(BVM)
		tx.Sign(acH)
		_, err = rp.ManageContractByVote(tx)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "ManagerContractByVote] is not enable"))
	})

	t.Run("auth_more_methods_and_revoke", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = ["","user"]
		//id = 1
		//methods = ["contract_*", "block_*","archive_*", "sub_*", "node_*", "cert_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"contract_*", "block_*", "archive_*", "sub_*", "node_*", "cert_*"},
				Name:            "node",
			},
		}
		stdError := rp.SetRulesInNode(rules)
		assert.Nil(t, stdError)

		// add user,admin to account D
		accountD, _ := account.NewAccount(defaultPwd)
		acD, _ := account.NewAccountFromAccountJSON(accountD, defaultPwd)
		stdError = rp.AddRoleForNode(acD.GetAddress().Hex(), user, admin)
		assert.Nil(t, stdError)

		rp.SetAccount(acD)
		// get contract_*
		_, err := rp.GetContractCountByAddr(acD.GetAddress().Hex())
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", acD.GetAddress().Hex())))

		// get block_*
		_, err = rp.GetChainHeight()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", acD.GetAddress().Hex())))

		// get archive_*
		_, err = rp.ListSnapshot()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", acD.GetAddress().Hex())))

		// get node_*
		_, err = rp.GetNodes()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", acD.GetAddress().Hex())))

		// sub_*, cert_* has not get api

		// delete user from account
		// add user,admin to account D
		stdError = rp.DeleteRoleFromNode(acD.GetAddress().Hex(), user)
		assert.Nil(t, stdError)

		rp.SetAccount(acD)
		// get contract_*
		_, err = rp.GetContractCountByAddr(acD.GetAddress().Hex())
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "Account dose not exist or account balance is 0"))

		// get block_*
		_, err = rp.GetChainHeight()
		assert.Nil(t, err)

		// get archive_*
		_, err = rp.ListSnapshot()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "The process of snapshot or archive happened error"))

		// get node_*
		_, err = rp.GetNodes()
		assert.Nil(t, err)
	})

}

func TestRPC_InspectorRole_did(t *testing.T) {
	t.Skip()
	// 注意此单测只允许配置一个节点
	defaultPwd := ""

	rp, err := NewJsonRPC()
	assert.Nil(t, err)

	newAc := `{"address":"0x430d637bb98033e1a2406ab8c24949e142fa4e75","algo":"0x03","version":"4.0","publicKey":"0x0454bbfe47a39685ac16d8263b0f7e3e0406996dc218543399411dfff12e304768945d0afe07c6b932fb203d1fb670fc46839adce54efed8e153fd21779c29d02b","privateKey":"e549225831723fdf9a7d4530b17b746e2c9da37fc1a982f4927f94d2b2972ac2"}`
	acKey, err := account.NewAccountFromAccountJSON(newAc, defaultPwd)
	assert.Nil(t, err)
	rp.SetAccount(acKey)
	rp.SetLocalChainID()
	newAccount := `{"address":"0x430d637bb98033e1a2406ab8c24949e142fa4e75","algo":"0x03","version":"4.0","publicKey":"0x0454bbfe47a39685ac16d8263b0f7e3e0406996dc218543399411dfff12e304768945d0afe07c6b932fb203d1fb670fc46839adce54efed8e153fd21779c29d02b","privateKey":"e549225831723fdf9a7d4530b17b746e2c9da37fc1a982f4927f94d2b2972ac2"}`
	ac, err := account.GenDIDKeyFromAccountJson(newAccount, defaultPwd)
	assert.Nil(t, err)

	didKey := account.NewDIDAccount(ac.(account.Key), rp.chainID, common.RandomString(10))
	puKey, _ := GenDIDPublicKeyFromDIDKey(didKey)
	document := NewDIDDocument(didKey.GetAddress(), puKey, nil)
	tx := NewTransaction(didKey.GetAddress()).Register(document)
	_, derr := rp.SendDIDTransaction(tx, didKey)
	assert.Nil(t, derr)

	role := "user"
	role1 := "account"
	t.Run("add_role_success", func(t *testing.T) {
		stdError := rp.AddRoleForNode(didKey.GetAddress(), role, role1)
		assert.Nil(t, stdError)
	})

	t.Run("add_role_already_exist", func(t *testing.T) {
		stdError := rp.AddRoleForNode(didKey.GetAddress(), role)
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address:%s already has roles:[%s]", didKey.GetAddress(), role)))
	})

	t.Run("delete_role_success", func(t *testing.T) {
		stdError := rp.DeleteRoleFromNode(didKey.GetAddress(), role1)
		assert.Nil(t, stdError)
	})

	t.Run("delete_role_fail", func(t *testing.T) {
		stdError := rp.DeleteRoleFromNode(didKey.GetAddress(), role1)
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address:%s has not roles:[%s]", didKey.GetAddress(), role1)))
	})

	t.Run("get_role", func(t *testing.T) {
		roles, err := rp.GetRoleFromNode(didKey.GetAddress())
		assert.Nil(t, err)
		assert.Equal(t, roles, []string{role})
	})

	t.Run("get_all_role", func(t *testing.T) {
		roles, err := rp.GetAllRolesFromNode()
		assert.Nil(t, err)
		assert.NotNil(t, roles)
	})

	// inspector enable
	t.Run("auth_forbidden_and_authorized_and_methods", func(t *testing.T) {
		// set rules
		user := "user"
		admin := "admin"
		rule := &InspectorRule{
			AllowAnyone:     false,
			AuthorizedRoles: []string{admin},
			ForbiddenRoles:  []string{user},
			ID:              0,
			Method:          []string{"node_*"},
			Name:            "node",
		}
		stdError := rp.SetRulesInNode([]*InspectorRule{rule})
		assert.Nil(t, stdError)

		// add user to account A
		accountA, _ := account.NewAccount(defaultPwd)
		acA, _ := account.GenDIDKeyFromAccountJson(accountA, defaultPwd)
		didAcA := account.NewDIDAccount(acA.(account.Key), rp.chainID, common.RandomString(10))
		puKey, _ := GenDIDPublicKeyFromDIDKey(didAcA)
		document := NewDIDDocument(didAcA.GetAddress(), puKey, nil)
		tx := NewTransaction(didAcA.GetAddress()).Register(document)
		_, derr = rp.SendDIDTransaction(tx, didAcA)
		assert.Nil(t, derr)

		stdError = rp.AddRoleForNode(didAcA.GetAddress(), user)
		assert.Nil(t, stdError)
		// use account A call node_getNodeStates has not permission
		// 设置节点账户
		rp.SetAccount(didAcA)
		_, stdError = rp.GetNodeStates()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodeStates", didAcA.GetAddress())))
		// use account A call node_getNodes has not permission
		_, stdError = rp.GetNodes()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodes", didAcA.GetAddress())))

		// add admin to account B
		accountB, _ := account.NewAccount(defaultPwd)
		acB, _ := account.GenDIDKeyFromAccountJson(accountB, defaultPwd)
		didAcB := account.NewDIDAccount(acB.(account.Key), rp.chainID, common.RandomString(11))
		puKey, _ = GenDIDPublicKeyFromDIDKey(didAcB)
		document = NewDIDDocument(didAcB.GetAddress(), puKey, nil)
		tx = NewTransaction(didAcB.GetAddress()).Register(document)
		rp.SendDIDTransaction(tx, didAcB)

		stdError = rp.AddRoleForNode(didAcB.GetAddress(), admin)
		assert.Nil(t, stdError)
		// use account B call node_getNodeStates has permission
		rp.SetAccount(didAcB)
		nodeState, stdError := rp.GetNodeStates()
		assert.Nil(t, stdError)
		assert.NotNil(t, nodeState)
		t.Log(nodeState)
	})

	t.Run("auth_rules_id", func(t *testing.T) {
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ID:              0,
				Method:          []string{"node_*"},
				Name:            "node",
			},
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"node_*"},
				Name:            "node",
			},
		}
		stdError := rp.SetRulesInNode(rules)
		assert.Nil(t, stdError)
		inspectorRules, stdError := rp.GetRulesFromNode()
		assert.Nil(t, stdError)
		assert.Len(t, inspectorRules, len(rules))
	})

	t.Run("auth_priority", func(t *testing.T) {
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     true,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"node_*"},
				Name:            "node",
			},
		}
		stdError := rp.SetRulesInNode(rules)
		assert.Nil(t, stdError)

		// add user to account D
		accountD, _ := account.NewAccount(defaultPwd)
		acD, _ := account.GenDIDKeyFromAccountJson(accountD, defaultPwd)
		didAcD := account.NewDIDAccount(acD.(account.Key), rp.chainID, common.RandomString(5))
		puKey, _ := GenDIDPublicKeyFromDIDKey(didAcD)
		document := NewDIDDocument(didAcD.GetAddress(), puKey, nil)
		tx := NewTransaction(didAcD.GetAddress()).Register(document)
		rp.SendDIDTransaction(tx, didAcD)

		stdError = rp.AddRoleForNode(didAcD.GetAddress(), user, admin)
		assert.Nil(t, stdError)
		// use account D call node_getNodeStates has not permission
		rp.SetAccount(didAcD)
		_, stdError = rp.GetNodeStates()
		assert.Error(t, stdError)
		assert.True(t, strings.Contains(stdError.Error(), fmt.Sprintf("address %s has not permission to access node_getNodeStates", didAcD.GetAddress())))
	})

	t.Run("auth_more_methods_and_revoke", func(t *testing.T) {
		// set rules
		//[[inspector.rules]]
		//allow_anyone = false
		//authorized_roles = ["admin"]
		//forbidden_roles = ["","user"]
		//id = 1
		//methods = ["contract_*", "block_*","archive_*", "sub_*", "node_*", "cert_*"]
		//name = "node"
		user := "user"
		admin := "admin"
		rules := []*InspectorRule{
			{
				AllowAnyone:     false,
				AuthorizedRoles: []string{admin},
				ForbiddenRoles:  []string{user},
				ID:              1,
				Method:          []string{"contract_*", "block_*", "archive_*", "sub_*", "node_*", "cert_*"},
				Name:            "node",
			},
		}
		stdError := rp.SetRulesInNode(rules)
		assert.Nil(t, stdError)

		// add user,admin to account D
		accountD, _ := account.NewAccount(defaultPwd)
		acD, _ := account.GenDIDKeyFromAccountJson(accountD, defaultPwd)
		didAcD := account.NewDIDAccount(acD.(account.Key), rp.chainID, common.RandomString(8))
		puKey, _ := GenDIDPublicKeyFromDIDKey(didAcD)
		document := NewDIDDocument(didAcD.GetAddress(), puKey, nil)
		tx := NewTransaction(didAcD.GetAddress()).Register(document)
		rp.SendDIDTransaction(tx, didAcD)

		stdError = rp.AddRoleForNode(didAcD.GetAddress(), user, admin)
		assert.Nil(t, stdError)

		rp.SetAccount(didAcD)
		// get contract_*
		_, err := rp.GetContractCountByAddr(didAcD.GetAddress())
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", didAcD.GetAddress())))

		// get block_*
		_, err = rp.GetChainHeight()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", didAcD.GetAddress())))

		// get archive_*
		_, err = rp.ListSnapshot()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", didAcD.GetAddress())))

		// get node_*
		_, err = rp.GetNodes()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("address %s has not permission to access", didAcD.GetAddress())))

		// sub_*, cert_* has not get api

		// delete user from account
		// add user,admin to account D
		stdError = rp.DeleteRoleFromNode(didAcD.GetAddress(), user)
		assert.Nil(t, stdError)

		rp.SetAccount(didAcD)
		// get contract_*
		_, err = rp.GetContractCountByAddr(didAcD.GetAddress())
		assert.Nil(t, err)

		// get block_*
		_, err = rp.GetChainHeight()
		assert.Nil(t, err)

		// get archive_*
		_, err = rp.ListSnapshot()
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "The process of snapshot or archive happened error"))

		// get node_*
		_, err = rp.GetNodes()
		assert.Nil(t, err)
	})

}

func TestRPC_InspectorNotEnable(t *testing.T) {
	t.Skip()
	rp, err := NewJsonRPC()
	assert.Nil(t, err)
	user := "user"
	admin := "admin"
	rules := []*InspectorRule{
		{
			AllowAnyone:     false,
			AuthorizedRoles: []string{admin},
			ForbiddenRoles:  []string{user},
			ID:              0,
			Method:          []string{"node_*"},
			Name:            "node",
		},
	}
	stdError := rp.SetRulesInNode(rules)
	assert.Nil(t, stdError)

	// add user, admin to account
	accountD, _ := account.NewAccount("")
	acD, _ := account.NewAccountFromAccountJSON(accountD, "")
	stdError = rp.AddRoleForNode(acD.GetAddress().Hex(), user, admin)
	assert.Nil(t, stdError)

	rp.SetAccount(acD)
	states, err := rp.GetNodeStates()
	assert.Nil(t, err)
	assert.NotNil(t, states)
}
