package rpc

import (
	"encoding/hex"
	"fmt"
	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/jackzing/gosdk/account"
	"github.com/jackzing/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAggSign(t *testing.T) {
	rp := NewRPC()
	user1, err := account.NewAccountSm2FromPriv("14f02ad2f0fc12f25b9a233fb182fa064679c1077fdbbc77d1a945806ebfee35")
	assert.Nil(t, err)

	user2, err := account.NewAccountSm2FromPriv("5fa355d98ba611d56b05d18a3b665a9e276f21487043c2c88e8a0ce84c0249cb")
	assert.Nil(t, err)
	pk1, err := user1.PublicBytes()
	assert.Nil(t, err)
	pk2, err := user2.PublicBytes()
	assert.Nil(t, err)
	pks := []gm.PubKey{pk1, pk2}

	first, err := account.NewAggSignerFromAccount(user1, pks)
	assert.Nil(t, err)

	second, err := account.NewAggSignerFromAccount(user2, pks)
	assert.Nil(t, err)

	aggCommitment, err := first.AggCommitment(first.Commitment(), second.Commitment())
	assert.Nil(t, err)
	address, _ := testPrivateAccount()

	transaction := NewTransaction(first.Address().Hex()).Participant(&Participant{
		Initiator:   pk1,
		Withholding: [][]byte{pk2},
	}).Transfer(address, int64(0))
	s1, err := first.PartSign([]byte(transaction.NeedHashString()), aggCommitment)
	assert.Nil(t, err)
	s2, err := second.PartSign([]byte(transaction.NeedHashString()), aggCommitment)
	assert.Nil(t, err)

	aggSign, err := first.AggSign(s1, s2)
	assert.Nil(t, err)
	fmt.Println(common.Bytes2Hex(aggSign))
	assert.Nil(t, first.Verify([]byte(transaction.NeedHashString()), aggSign))
	transaction.Signature(common.Bytes2Hex(aggSign))

	resp, err := rp.SignAndInvokeContract(transaction, nil)
	assert.Nil(t, err)
	fmt.Println(resp)

}

func TestDIDAggSign(t *testing.T) {
	rp := NewRPC()
	iniDID(t)
	suffix1 := common.RandomString(15)
	suffix2 := common.RandomString(20)

	// 生成账户并注册
	user1, err := account.NewAccountSm2FromPriv("14f02ad2f0fc12f25b9a233fb182fa064679c1077fdbbc77d1a945806ebfee35")
	assert.Nil(t, err)
	did1 := account.NewDIDAccount(user1, "chainID_01", suffix1)
	puKey1, err := GenDIDPublicKeyFromDIDKey(did1)
	assert.Nil(t, err)
	document1 := NewDIDDocument(did1.GetAddress(), puKey1, nil)
	tx1 := NewTransaction(did1.GetAddress()).Register(document1)
	_, stdErr := rp.SendDIDTransaction(tx1, did1)
	assert.Nil(t, stdErr)

	// 生成账户并注册
	user2, err := account.NewAccountSm2FromPriv("5fa355d98ba611d56b05d18a3b665a9e276f21487043c2c88e8a0ce84c0249cb")
	assert.Nil(t, err)
	did2 := account.NewDIDAccount(user2, "chainID_01", suffix2)
	puKey2, err := GenDIDPublicKeyFromDIDKey(did2)
	assert.Nil(t, err)
	document2 := NewDIDDocument(did2.GetAddress(), puKey2, nil)
	tx2 := NewTransaction(did2.GetAddress()).Register(document2)
	_, stdErr = rp.SendDIDTransaction(tx2, did2)
	assert.Nil(t, stdErr)

	k1, _ := did1.PublicBytes()
	k2, _ := did2.PublicBytes()
	kb := []gm.PubKey{k1, k2}

	signer, err := account.NewAggSignerFromDIDAccount(did1, kb)
	assert.Nil(t, err)
	t.Logf("signer: %v", hex.EncodeToString(signer.Commitment()))

	signer2, err := account.NewAggSignerFromDIDAccount(did2, kb)
	assert.Nil(t, err)
	t.Logf("signer: %v", hex.EncodeToString(signer2.Commitment()))
	address, _ := testPrivateAccount()
	tx3 := NewTransaction(signer.Address().Hex()).Participant(&Participant{
		Initiator:   []byte(did1.GetAddress()),
		Withholding: [][]byte{k2},
	}).Transfer(address, int64(0))

	aggCommitment, err := signer.AggCommitment(signer.Commitment(), signer2.Commitment())
	assert.Nil(t, err)

	s1, err := signer.PartSign([]byte(tx3.NeedHashString()), aggCommitment)
	assert.Nil(t, err)
	s2, err := signer2.PartSign([]byte(tx3.NeedHashString()), aggCommitment)
	assert.Nil(t, err)

	aggSign, err := signer.AggSign(s1, s2)
	assert.Nil(t, err)
	assert.Nil(t, signer.Verify([]byte(tx3.NeedHashString()), aggSign))
	tx3.Signature(common.Bytes2Hex(aggSign))
	resp, err := rp.SignAndInvokeContract(tx3, nil)
	assert.Nil(t, err)
	fmt.Println(resp)

}
