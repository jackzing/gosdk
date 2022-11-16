package scale

import (
	"fmt"
	"strings"
)

// AbiMap abi map.
type AbiMap struct {
	abis map[string]*Abi
}

// NewAbiMap create a new abi Map.
func NewAbiMap() *AbiMap {
	return &AbiMap{
		abis: make(map[string]*Abi),
	}
}

func (m *AbiMap) GetAbi(addr string) (*Abi, error) {
	addr = strings.ToLower(addr)
	if addr[0:2] == "0x" {
		addr = addr[2:]
	}
	abi, ok := m.abis[addr]
	if !ok {
		return nil, fmt.Errorf("can not find abi of %s", addr)
	}
	return abi, nil
}

func (m *AbiMap) SetAbi(addr string, abi *Abi) {
	addr = strings.ToLower(addr)
	if addr[0:2] == "0x" {
		addr = addr[2:]
	}
	m.abis[addr] = abi
}
