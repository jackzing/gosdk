package account

import (
	"encoding/hex"
	"errors"
	"fmt"
	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/go-crypto-standard/hash"
	"github.com/hyperchain/gosdk/common"
)

type AggSigner struct {
	account        interface{}
	publicKeys     []gm.PubKey
	accountAddress string
	publicKey      gm.PubKey
	commitment     gm.Commitment
	index          int
	aggContext     *gm.AggContext
	hasDID         bool
}

func NewAggSignerFromAccount(key Key, publicKeys []gm.PubKey) (*AggSigner, error) {
	switch key.(type) {
	case *SM2Key:
		publicKey, err := key.PublicBytes()
		if err != nil {
			return nil, err
		}
		var index = -1
		for i, v := range publicKeys {
			if hex.EncodeToString(publicKey) == hex.EncodeToString(v) {
				index = i
				break
			}
		}
		if index == -1 {
			return nil, fmt.Errorf("account %s is not in publicKeys", key.GetAddress().Hex())
		}
		var aggContext = new(gm.AggContext)
		commitment, p, err := aggContext.Init(index, publicKeys...)
		if err != nil {
			return nil, err
		}
		address := key.GetAddress()
		return &AggSigner{
			account:        key,
			publicKeys:     publicKeys,
			accountAddress: address.Hex(),
			publicKey:      p,
			index:          index,
			aggContext:     aggContext,
			hasDID:         false,
			commitment:     commitment,
		}, nil
	default:
		return nil, errors.New("not support key type")
	}
}

func NewAggSignerFromDIDAccount(did *DIDKey, publicKeys []gm.PubKey) (*AggSigner, error) {
	switch did.GetNormalKey().(type) {
	case *SM2Key:
		publicKey, err := did.PublicBytes()
		if err != nil {
			return nil, err
		}
		var index = -1
		for i, v := range publicKeys {
			if hex.EncodeToString(publicKey) == hex.EncodeToString(v) {
				index = i
				break
			}
		}
		if index == -1 {
			return nil, fmt.Errorf("account %s is not in publicKeys", did.GetAddress())
		}
		var aggContext = new(gm.AggContext)
		commitment, p, err := aggContext.Init(index, publicKeys...)
		if err != nil {
			return nil, err
		}
		didAddress := did.GetAddress()
		return &AggSigner{
			account:        did,
			publicKeys:     publicKeys,
			accountAddress: didAddress,
			publicKey:      p,
			index:          index,
			aggContext:     aggContext,
			hasDID:         true,
			commitment:     commitment,
		}, nil
	default:
		return nil, errors.New("not support key type")
	}
}

func (a *AggSigner) AggCommitment(in ...gm.Commitment) (gm.Commitment, error) {
	return a.aggContext.AggCommitment(in...)
}

func (a *AggSigner) PartSign(msg []byte, aggCommitment gm.Commitment) (gm.Signature, error) {
	var account *SM2Key
	if a.hasDID {
		account = a.account.(*DIDKey).GetNormalKey().(*SM2Key)
	} else {
		account = a.account.(*SM2Key)
	}
	keyBytes, err := account.Bytes()
	if err != nil {
		return nil, err
	}
	hashMsg, err := gm.NewSM3Hasher().Hash(msg)
	if err != nil {
		return nil, err
	}
	return a.aggContext.PartSign(keyBytes, hashMsg, aggCommitment)
}

func (a *AggSigner) AggSign(signs ...gm.Signature) (gm.Signature, error) {
	aggFlag := byte(7)
	didAggFlag := byte(7 | 128)
	if a.hasDID {
		aggFlag = didAggFlag
	}
	ans, err := a.aggContext.AggSign(signs...)
	if err != nil {
		return nil, err
	}
	var res gm.Signature
	res = append(res, aggFlag)
	res = append(res, a.publicKey...)
	res = append(res, ans...)
	return res, nil
}

func (a *AggSigner) Verify(msg []byte, sign gm.Signature) error {
	realSign := make([]byte, 65)
	copy(realSign, sign[66:])
	hashMsg, err := gm.NewSM3Hasher().Hash(msg)
	if err != nil {
		return err
	}
	return a.aggContext.Verify(a.publicKey, hashMsg, realSign)
}

func (a *AggSigner) Commitment() gm.Commitment {
	return a.commitment
}

func (a *AggSigner) PubKey() gm.PubKey {
	return a.publicKey
}

func (a *AggSigner) Address() common.Address {
	bs := a.PubKey()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(bs)
	return common.BytesToAddress(h[12:])
}

func (a *AggSigner) OriginAccountAddress() string {
	return a.accountAddress
}

func (a *AggSigner) SetHasDID(flag bool) {
	a.hasDID = flag
}
