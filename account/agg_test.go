package account

import (
	"bytes"
	"encoding/hex"
	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/jackzing/gosdk/common"
	"github.com/meshplus/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAggSignerFromAccount(t *testing.T) {
	k1K, _ := hex.DecodeString("14f02ad2f0fc12f25b9a233fb182fa064679c1077fdbbc77d1a945806ebfee35")
	k2K, _ := hex.DecodeString("5fa355d98ba611d56b05d18a3b665a9e276f21487043c2c88e8a0ce84c0249cb")
	k3K, _ := hex.DecodeString("365a2cf192b9dea77771790f218267ba950b663c64f3167c506077f661de3f4e")

	k1, _ := hex.DecodeString("04a642d77b4c8688124e343af9531de801813f1a880b664ce38a78d76ed128d07b1815004fe1de9c903b349afba9c74775ecc1692a6dd230689ea096521402e834")
	k2, _ := hex.DecodeString("042837e98f62ce38e93d41c949425b8e9a5b46107f2bd8725cc7db91b317c6d3199fd87f51883f9cc76a5535ec3b52dac4778597d26e0bc4a0b164cb4b6e87a3d5")
	k3, _ := hex.DecodeString("0445a79c4d689b9f19dbf19d58d1cb1952a8b557399f2c3a94a1e52fa21fa47c1acf7bc909e76449914a61a7ada1a855c3d3142181ce35b65cc172e5388dc494d6")
	ka := []*gm.SM2PrivateKey{new(gm.SM2PrivateKey), new(gm.SM2PrivateKey), new(gm.SM2PrivateKey)}
	assert.Nil(t, ka[0].FromBytes(k1K, crypto.Sm2p256v1))
	assert.Nil(t, ka[0].FromBytes(k2K, crypto.Sm2p256v1))
	assert.Nil(t, ka[0].FromBytes(k3K, crypto.Sm2p256v1))
	kb := []gm.PubKey{k1, k2, k3}

	user1, err := NewAccountSm2FromPriv("14f02ad2f0fc12f25b9a233fb182fa064679c1077fdbbc77d1a945806ebfee35")
	assert.Nil(t, err)
	signer, err := NewAggSignerFromAccount(user1, kb)
	assert.Nil(t, err)
	t.Logf("signer: %v", hex.EncodeToString(signer.commitment))

	t.Run("flow", func(t *testing.T) {
		user2, err := NewAccountSm2FromPriv("5fa355d98ba611d56b05d18a3b665a9e276f21487043c2c88e8a0ce84c0249cb")
		assert.Nil(t, err)
		signer2, err := NewAggSignerFromAccount(user2, kb)
		assert.Nil(t, err)
		t.Logf("signer2: %v", hex.EncodeToString(signer2.commitment))
		user3, err := NewAccountSm2FromPriv("365a2cf192b9dea77771790f218267ba950b663c64f3167c506077f661de3f4e")
		assert.Nil(t, err)
		signer3, err := NewAggSignerFromAccount(user3, kb)
		assert.Nil(t, err)
		t.Logf("signer3: %v", hex.EncodeToString(signer3.commitment))

		// 比较三者的聚合公钥
		assert.True(t, bytes.Equal(signer.publicKey, signer2.publicKey))
		assert.True(t, bytes.Equal(signer2.publicKey, signer3.publicKey))

		// 聚合随机承诺
		aggCom, err := signer.AggCommitment(signer.commitment, signer2.commitment, signer3.commitment)
		assert.Nil(t, err)
		message := []byte("One ping only, please.One ping only, please.One ping only, please.One ping only, please.One ping only, please.")
		// 部分签名
		s0, err := signer.PartSign(message, aggCom)
		assert.Nil(t, err)
		s1, err := signer2.PartSign(message, aggCom)
		assert.Nil(t, err)
		s2, err := signer3.PartSign(message, aggCom)
		assert.Nil(t, err)
		// 聚合签名
		aggSign, err := signer.AggSign(s0, s1, s2)
		assert.Nil(t, err)
		t.Logf("aggSign: %v", hex.EncodeToString(aggSign))
		//5.验证
		assert.Nil(t, signer.Verify(message, aggSign))
		assert.Nil(t, signer2.Verify(message, aggSign))
		assert.Nil(t, signer3.Verify(message, aggSign))
	})
}

func TestNewAggSignerFromDIDAccount(t *testing.T) {
	k1K, _ := hex.DecodeString("14f02ad2f0fc12f25b9a233fb182fa064679c1077fdbbc77d1a945806ebfee35")
	k2K, _ := hex.DecodeString("5fa355d98ba611d56b05d18a3b665a9e276f21487043c2c88e8a0ce84c0249cb")
	k3K, _ := hex.DecodeString("365a2cf192b9dea77771790f218267ba950b663c64f3167c506077f661de3f4e")

	k1, _ := hex.DecodeString("04a642d77b4c8688124e343af9531de801813f1a880b664ce38a78d76ed128d07b1815004fe1de9c903b349afba9c74775ecc1692a6dd230689ea096521402e834")
	k2, _ := hex.DecodeString("042837e98f62ce38e93d41c949425b8e9a5b46107f2bd8725cc7db91b317c6d3199fd87f51883f9cc76a5535ec3b52dac4778597d26e0bc4a0b164cb4b6e87a3d5")
	k3, _ := hex.DecodeString("0445a79c4d689b9f19dbf19d58d1cb1952a8b557399f2c3a94a1e52fa21fa47c1acf7bc909e76449914a61a7ada1a855c3d3142181ce35b65cc172e5388dc494d6")
	ka := []*gm.SM2PrivateKey{new(gm.SM2PrivateKey), new(gm.SM2PrivateKey), new(gm.SM2PrivateKey)}
	assert.Nil(t, ka[0].FromBytes(k1K, crypto.Sm2p256v1))
	assert.Nil(t, ka[0].FromBytes(k2K, crypto.Sm2p256v1))
	assert.Nil(t, ka[0].FromBytes(k3K, crypto.Sm2p256v1))
	kb := []gm.PubKey{k1, k2, k3}
	suffix := common.RandomString(10)

	user1, err := NewAccountSm2FromPriv("14f02ad2f0fc12f25b9a233fb182fa064679c1077fdbbc77d1a945806ebfee35")
	assert.Nil(t, err)
	did1 := NewDIDAccount(user1, "chainID_01", suffix)
	signer, err := NewAggSignerFromDIDAccount(did1, kb)
	assert.Nil(t, err)
	t.Logf("signer: %v", hex.EncodeToString(signer.commitment))

	user2, err := NewAccountSm2FromPriv("5fa355d98ba611d56b05d18a3b665a9e276f21487043c2c88e8a0ce84c0249cb")
	assert.Nil(t, err)
	did2 := NewDIDAccount(user2, "chainID_01", suffix)
	signer2, err := NewAggSignerFromDIDAccount(did2, kb)
	assert.Nil(t, err)
	t.Logf("signer2: %v", hex.EncodeToString(signer2.commitment))

	user3, err := NewAccountSm2FromPriv("365a2cf192b9dea77771790f218267ba950b663c64f3167c506077f661de3f4e")
	assert.Nil(t, err)
	did3 := NewDIDAccount(user3, "chainID_01", suffix)
	signer3, err := NewAggSignerFromDIDAccount(did3, kb)
	assert.Nil(t, err)
	t.Logf("signer3: %v", hex.EncodeToString(signer3.commitment))

	// 比较三者的聚合公钥
	assert.True(t, bytes.Equal(signer.publicKey, signer2.publicKey))
	assert.True(t, bytes.Equal(signer2.publicKey, signer3.publicKey))

	// 聚合随机承诺
	aggCom, err := signer.AggCommitment(signer.commitment, signer2.commitment, signer3.commitment)
	assert.Nil(t, err)
	message := []byte("One ping only, please.One ping only, please.One ping only, please.One ping only, please.One ping only, please.")
	// 部分签名
	s0, err := signer.PartSign(message, aggCom)
	assert.Nil(t, err)
	s1, err := signer2.PartSign(message, aggCom)
	assert.Nil(t, err)
	s2, err := signer3.PartSign(message, aggCom)
	assert.Nil(t, err)
	// 聚合签名
	aggSign, err := signer.AggSign(s0, s1, s2)
	assert.Nil(t, err)
	t.Logf("aggSign: %v", hex.EncodeToString(aggSign))
	//5.验证
	assert.Nil(t, signer.Verify(message, aggSign))
	assert.Nil(t, signer2.Verify(message, aggSign))
	assert.Nil(t, signer3.Verify(message, aggSign))
}
