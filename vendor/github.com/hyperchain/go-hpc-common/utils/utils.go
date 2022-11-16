// Hyperchain License
// Copyright (C) 2016 The Hyperchain Authors.

package utils

import (
	"errors"
	"github.com/hyperchain/go-hpc-common/types/protos"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// IntArrayEquals checks if the arrays of ints are the same
func IntArrayEquals(a []int, b []int) bool {
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

// StringArrayEquals checks if the arrays of strings are the same
func StringArrayEquals(a []string, b []string) bool {
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

// Clone clones the passed slice
func Clone(src []byte) []byte {
	clone := make([]byte, len(src))
	copy(clone, src)

	return clone
}

// TransportEncode transfer string to hex
func TransportEncode(str string) string {
	b := []byte(str)
	return ToHex(b)
}

// TransportDecode transfer hex to string
func TransportDecode(str string) string {
	b := HexToBytes(str)
	return string(b)
}

// CalculateSize calculates the length of a string
func CalculateSize(sizeStr string) (int64, error) {
	switch {
	case strings.Contains(sizeStr, "mb"):
		strs := strings.Split(sizeStr, "mb")
		size, err := strconv.ParseInt(strs[0], 10, 64)
		if err != nil {
			return 0, err
		}
		return size * 1024 * 1024, nil
	case strings.Contains(sizeStr, "kb"):
		strs := strings.Split(sizeStr, "kb")
		size, err := strconv.ParseInt(strs[0], 10, 64)
		if err != nil {
			return 0, err
		}
		return size * 1024, nil
	default:
		return 0, errors.New("incorrect configuration parameters " + sizeStr + ".")
	}
}

// CheckNilElems checks if provided struct has nil elements, returns error if provided
// param is not a struct pointer and returns all nil elements' name if has.
func CheckNilElems(i interface{}) (string, []string, error) {
	typ := reflect.TypeOf(i)
	value := reflect.Indirect(reflect.ValueOf(i))

	if typ.Kind() != reflect.Ptr {
		return "", nil, errors.New("got a non-ptr to check if has nil elements")
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		return "", nil, errors.New("got a non-struct to check if has nil elements")
	}

	structName := typ.Name()
	nilElems := make([]string, 0)
	hasNil := false

	for i := 0; i < typ.NumField(); i++ {
		kind := typ.Field(i).Type.Kind()
		if kind == reflect.Chan || kind == reflect.Map {
			elemName := typ.Field(i).Name
			if value.FieldByName(elemName).IsNil() {
				nilElems = append(nilElems, elemName)
				hasNil = true
			}
		}
	}
	if hasNil {
		return structName, nilElems, nil
	}
	return structName, nil, nil
}

// IsValidNamespace judge if namespace name is valid
func IsValidNamespace(namespace string) error {
	if strings.HasPrefix(namespace, "system") {
		return errors.New("can not new namespace with prefix 'system'")
	}
	if len(namespace) > 20 {
		return errors.New("can not use namespace name whose length longer than 20")
	}
	return nil
}

// GetOrderedTxList sort the valid transaction and invalid transaction
func GetOrderedTxList(block *protos.Block) []*protos.Transaction {
	txs := make([]*protos.Transaction, len(block.Transactions)+len(block.InvalidRecords))
	for _, record := range block.InvalidRecords {
		txs[record.Index] = record.Tx
	}
	j := 0
	for i := 0; i < len(txs); i++ {
		if txs[i] == nil {
			txs[i] = block.Transactions[j]
			j++
		}
	}
	return txs
}

// SizeOfFile returns the size of a file.
func SizeOfFile(name string) (int64, error) {
	info, sterr := os.Stat(name)
	if sterr != nil {
		return 0, sterr
	}
	return info.Size(), nil
}
