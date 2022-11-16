package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/hyperchain/go-hpc-common/utils"
)

//HPCAccount wrap address and DIDAccount
type HPCAccount struct {
	Address
	DidAccount *DIDAccount
	IsDID      bool
}

// NewHPCAccount create HPCAccount
func NewHPCAccount() *HPCAccount {
	return &HPCAccount{}
}

// NewAccountFromAddress create HPCAccount by address
func NewAccountFromAddress(address Address) *HPCAccount {
	return &HPCAccount{
		Address: address,
		IsDID:   false,
	}
}

// NewAccountFromDID create HPCAccount by didAccount
func NewAccountFromDID(didAccount *DIDAccount) *HPCAccount {
	return &HPCAccount{
		DidAccount: didAccount,
		IsDID:      true,
	}
}

//MarshalJSON marshal the given HPCAccount to json
func (hpcAccount *HPCAccount) MarshalJSON() ([]byte, error) {
	if hpcAccount.IsDID {
		return hpcAccount.DidAccount.MarshalJSON()
	}
	return hpcAccount.Address.MarshalJSON()
}

//UnmarshalJSON parse HPCAccount from raw json data
func (hpcAccount *HPCAccount) UnmarshalJSON(data []byte) error {
	var err error
	hpcAccount.DidAccount = NewDIDAccount()
	if err = hpcAccount.DidAccount.UnmarshalJSON(data); err == nil {
		hpcAccount.IsDID = true
		return nil
	}
	hpcAccount.DidAccount = nil
	if err1 := hpcAccount.Address.UnmarshalJSON(data); err1 == nil {
		hpcAccount.IsDID = false
		return nil
	}
	return err
}

//GetChainID return the DIDAccount's chainID
func (hpcAccount *HPCAccount) GetChainID() []byte {
	if hpcAccount.IsDID && hpcAccount.DidAccount != nil {
		return hpcAccount.DidAccount.chainID
	}
	return nil
}

//Hex the hex string representation of the underlying address
func (hpcAccount *HPCAccount) Hex() string {
	if !hpcAccount.IsDID {
		return hpcAccount.Address.Hex()
	}
	return hpcAccount.DidAccount.Hex()
}

//Length return the length
func (hpcAccount *HPCAccount) Length() int {
	if !hpcAccount.IsDID {
		return len(hpcAccount.Address)
	}
	return len(hpcAccount.DidAccount.origin)
}

//Bytes return the address's bytes
func (hpcAccount *HPCAccount) Bytes() []byte {
	if !hpcAccount.IsDID {
		return hpcAccount.Address.Bytes()
	}
	return hpcAccount.DidAccount.origin[:]
}

//Str is the string representation of the underlying hpcAccount
func (hpcAccount *HPCAccount) Str() string {
	if !hpcAccount.IsDID {
		return hpcAccount.Address.Str()
	}
	return hpcAccount.DidAccount.Str()
}

// Copy copys the given HPCAccount
func (hpcAccount *HPCAccount) Copy() *HPCAccount {
	res := &HPCAccount{
		Address: hpcAccount.Address,
		IsDID:   hpcAccount.IsDID,
	}
	if hpcAccount.DidAccount != nil {
		res.DidAccount = &DIDAccount{
			origin:  utils.CopyBytes(hpcAccount.DidAccount.origin),
			prefix:  utils.CopyBytes(hpcAccount.DidAccount.prefix),
			chainID: utils.CopyBytes(hpcAccount.DidAccount.chainID),
			did:     utils.CopyBytes(hpcAccount.DidAccount.did),
		}
	}
	return res
}

// ToString return string representation for user
func (hpcAccount *HPCAccount) ToString() string {
	if hpcAccount.IsDID {
		return hpcAccount.DidAccount.Str()
	}
	return hpcAccount.Address.Hex()
}

// GetAddress return the Address of HPCAccount
func (hpcAccount *HPCAccount) GetAddress() (Address, error) {
	if !hpcAccount.IsDID {
		return hpcAccount.Address, nil
	}
	return TransformDIDToAddr(hpcAccount.DidAccount.origin)
}

// TransformDIDToAddr convert did address to Address type
func TransformDIDToAddr(didAddr []byte) (Address, error) {
	hash := sha256.New()
	_, err := hash.Write(didAddr)
	if err != nil {
		return Address{}, fmt.Errorf("transform didAddr `%s` failed:%s", string(didAddr), err.Error())
	}
	ret := hash.Sum(nil)
	addr := BytesToAddress(ret)
	copy(addr[0:], DIDAccountAddrPrefix)
	return addr, nil
}

// IsDIDAccount check whether the address is didAccount by compare prefix
func IsDIDAccount(addr Address) bool {
	if bytes.HasPrefix(addr.Bytes(), DIDAccountAddrPrefix) {
		return true
	}
	return false
}
