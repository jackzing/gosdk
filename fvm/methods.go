package fvm

import (
	"github.com/jackzing/gosdk/fvm/scale"
	"io"
)

func GenAbi(reader io.Reader) (*scale.Abi, error) {
	a, err := scale.JSON(reader)
	return &a, err
}

func Encode(abi *scale.Abi, method string, params ...interface{}) ([]byte, error) {
	return abi.Encode(method, params...)
}

func DecodeRet(abi *scale.Abi, method string, val []byte) (*scale.InvokeBean, error) {
	return abi.DecodeRet(val, method)
}

func EncodeParallel(abiMap *scale.AbiMap, addr string, method string, params ...interface{}) ([]byte, error) {
	return scale.Parallel(abiMap, addr, method, params...)
}
