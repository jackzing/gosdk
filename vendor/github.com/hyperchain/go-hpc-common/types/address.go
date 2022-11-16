// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/hyperchain/go-hpc-common/utils"

	json "github.com/json-iterator/go"
)

// const length
const (
	AddressLength = 20
)

var errHashJSONLength = errors.New("common: unmarshalJSON failed: hash must be exactly 32 bytes")

type (
	// Address type used in Hyperchain
	Address [AddressLength]byte
)

// StringToHex converts string to hex format
func StringToHex(s string) string {
	if len(s) >= 2 && s[:2] == "0x" {
		return s
	}
	return "0x" + s
}

// HexToString converts hex string to string without '0x' prefix
func HexToString(s string) string {
	if len(s) >= 2 && s[:2] == "0x" {
		return s[2:]
	}
	return s
}

// BytesToAddress converts []byte to Address
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// StringToAddress converts string to Address
func StringToAddress(s string) Address { return BytesToAddress([]byte(s)) }

// BigToAddress converts big Int to Address
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// HexToAddress converts hex string to Address
func HexToAddress(s string) Address { return BytesToAddress(utils.DecodeString(s)) }

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Ethereum address or not.
func IsHexAddress(s string) bool {
	if len(s) == 2+2*AddressLength && utils.IsHex(s) {
		return true
	}
	if len(s) == 2*AddressLength && utils.IsHex("0x"+s) {
		return true
	}
	return false
}

// Str is the string representation of the underlying Address
func (a Address) Str() string { return string(a[:]) }

// Bytes is the []byte representation of the underlying Address
func (a Address) Bytes() []byte { return a[:] }

// Big is the bigInt representation of the underlying Address
func (a Address) Big() *big.Int { return utils.Bytes2Big(a[:]) }

// Hash is the Hash representation of the underlying Address
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex is the hex string representation of the underlying Address
func (a Address) Hex() string { return "0x" + utils.BytesToHex(a[:]) }

// IsZero returns if Hash is empty
func (a Address) IsZero() bool { return a == Address{} }

// SetBytes sets the address to the value of b. If b is larger than len(a) it will panic
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// SetString set string `s` to a. If s is larger than len(a) it will panic
func (a *Address) SetString(s string) { a.SetBytes([]byte(s)) }

// Set sets a to other
func (a *Address) Set(other Address) {
	for i, v := range other {
		a[i] = v
	}
}

// MarshalJSON marshal the given Address to json
func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Hex())
}

// UnmarshalJSON parse address from raw json data
func (a *Address) UnmarshalJSON(data []byte) error {
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}

	if len(data) > 2 && data[0] == '0' && data[1] == 'x' {
		data = data[2:]
	}

	if len(data) != 2*AddressLength {
		return fmt.Errorf("invalid address length, expected %d got %d bytes", 2*AddressLength, len(data))
	}

	n, err := hex.Decode(a[:], data)
	if err != nil {
		return err
	}

	if n != AddressLength {
		return fmt.Errorf("invalid address")
	}

	a.Set(HexToAddress(string(data)))
	return nil
}

// PP Pretty Prints a byte slice in the following format:
// 	hex(value[:4])...(hex[len(value)-4:])
func PP(value []byte) string {
	if len(value) <= 8 {
		return utils.BytesToHex(value)
	}

	return fmt.Sprintf("%x...%x", value[:4], value[len(value)-4])
}

// AddressArrayEquals checks if the arrays of addresses are the same
func AddressArrayEquals(a []Address, b []Address) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
