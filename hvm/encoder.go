package hvm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jackzing/gosdk/common/parallel"

	"github.com/hyperchain/go-hpc-common/types"
	"github.com/jackzing/gosdk/common"
)

const (
	parallelPrefixSize  = 4
	parallelVersionSize = 1
	addressLastUseLen   = 4

	invokeBeanParameterStartPos = 4 + 2 // | class length(4B) | name length(2B)
)

func GenPayload(beanAbi *BeanAbi, params ...interface{}) (payload []byte, err error) {
	switch beanAbi.BeanType {
	case MethodBean:
		payload, _, err = methodBeanPayload(beanAbi, params...)
	case InvokeBean:
		fallthrough
	default:
		payload, _, err = invokeBeanPayload(beanAbi, params...)
	}
	return
}

func GenParallelPayload(abiGetter AbiGetter, addr types.Address, beanName string, invokeDirectly bool, params ...interface{}) (payload []byte, err error) {
	abi, err := abiGetter.GetAbi(addr.Hex())
	if err != nil {
		return nil, err
	}
	var (
		paraIndex []int
		beanAbi   *BeanAbi
	)
	if invokeDirectly {
		beanAbi, err = abi.GetMethodAbi(beanName)
	} else {
		beanAbi, err = abi.GetBeanAbi(beanName)
	}
	if err != nil {
		return nil, err
	}

	switch beanAbi.BeanType {
	case MethodBean:
		payload, paraIndex, err = methodBeanPayload(beanAbi, params...)
	case InvokeBean:
		fallthrough
	default:
		payload, paraIndex, err = invokeBeanPayload(beanAbi, params...)
	}
	if err != nil {
		return nil, err
	}
	mutexData, err := buildMutexElement(beanAbi, abiGetter, paraIndex, params...)
	if err != nil {
		return nil, err
	}
	if mutexData != nil {
		payload = append(mutexData, payload...)
	}
	return payload, err
}

// | class length(4B) | name length(2B) | class | class name | bin |
func invokeBeanPayload(beanAbi *BeanAbi, params ...interface{}) ([]byte, []int, error) {
	classBytes := beanAbi.classBytes()

	if len(classBytes) > 0xffff {
		return nil, nil, errors.New("the bean class is too large") // 64k
	}

	beanName := []byte(beanAbi.BeanName)
	isJson := true
	for _, str := range params {
		if _, ok := str.(string); !ok {
			isJson = false
			break
		}
	}
	var bin string
	var err error
	var paraIndex []int
	startPos := invokeBeanParameterStartPos + len(classBytes) + len(beanName)

	// todo invoke directly, parameter (Struct, Array, List, Map) encode result are different with invoke bean. Current we don't support (Struct, Array, List, Map) as parallel parameter.
	if isJson {
		bin, paraIndex, err = beanAbi.encodeJSON(startPos, params...)
	} else {
		bin, paraIndex, err = beanAbi.encode(startPos, params...)
	}

	if err != nil {
		return nil, nil, err
	}
	binBytes := []byte(bin)

	result := make([]byte, 0)
	classLenByte := common.IntToBytes4(len(classBytes))
	nameLenByte := common.IntToBytes2(len(beanName))
	result = append(result, classLenByte[:]...)
	result = append(result, nameLenByte[:]...)
	result = append(result, classBytes...)
	result = append(result, beanName...)
	result = append(result, binBytes...)

	return result, paraIndex, nil
}

func methodBeanPayload(methodAbi *BeanAbi, params ...interface{}) ([]byte, []int, error) {
	if len(params) != len(methodAbi.Inputs) {
		return nil, nil, errors.New("param count is not enough")
	}

	paramBuilder := NewParamBuilder(methodAbi.BeanName)

	for i, input := range methodAbi.Inputs {
		param := Convert(params[i])
		if input.Name == "" {
			return nil, nil, errors.New("input StructName is empty")
		}
		switch input.EntryType {
		case Bool, Byte, Short, Int, Long, Float, Double, Char, String:
			if str, ok := param.(string); ok {
				paramBuilder.appendPayload([]byte(input.StructName), []byte(str))
			}
		case Struct, Array, List, Map:
			paramBuilder.AddObject(input.StructName, param)
		}
	}
	return paramBuilder.Build(), paramBuilder.ParamIndex(), nil
}

func buildMutexElement(beanAbi *BeanAbi, abiGetter AbiGetter, paraIndex []int, params ...interface{}) (_ []byte, err error) {
	var elements []*MutexElement
	// build current contract invoke message
	if beanAbi.ParallelLevel == parallel.GlobalMutex {
		return nil, nil
	} else {
		elements, err = parseMutexElement(beanAbi, paraIndex, "", false)
		if err != nil {
			return nil, err
		}
	}

	// build cross call message
	for _, mutexCall := range beanAbi.MutexCalls {
		// get cross call contract address
		var (
			callAddr string
			ok       bool
		)
		if mutexCall.Address != "" {
			callAddr = mutexCall.Address
		} else {
			if mutexCall.AddressIndex > len(params) || mutexCall.AddressIndex < 1 {
				return nil, fmt.Errorf("illegal mutexCall address index %d", mutexCall.AddressIndex)
			}
			callAddr, ok = params[mutexCall.AddressIndex-1].(string)
			if !ok {
				return nil, fmt.Errorf("the %d param %v should be string", mutexCall.AddressIndex, params[mutexCall.AddressIndex-1])
			}
		}
		callAbi, err := abiGetter.GetAbi(callAddr)
		if err != nil {
			return nil, err
		}

		for _, callMethod := range mutexCall.Methods {
			beanAbi, err := callAbi.GetMethodAbi(callMethod)
			if err != nil {
				return nil, err
			}
			if beanAbi.ParallelLevel == parallel.GlobalMutex {
				return nil, nil
			} else {
				crossEle, err := parseMutexElement(beanAbi, nil, callAddr, true)
				if err != nil {
					return nil, err
				}
				elements = append(elements, crossEle...)
			}
		}
	}

	return encodeMutexElement(elements)
}

func parseMutexElement(beanAbi *BeanAbi, paraIndex []int, address string, crossCall bool) ([]*MutexElement, error) {
	elements := make([]*MutexElement, 0)
	if beanAbi.ParallelLevel == parallel.ContractMutex {
		elements = append(elements, &MutexElement{Address: "", FieldID: 0, ParaIndex: nil})
	} else {
		for _, field := range beanAbi.MutexFields {
			var element *MutexElement
			if crossCall {
				element = &MutexElement{
					Address:   address[len(address)-addressLastUseLen:],
					FieldID:   field.FieldID,
					ParaIndex: make([]uint64, 0),
				}
			} else {
				element = &MutexElement{
					FieldID:   field.FieldID,
					ParaIndex: make([]uint64, len(field.ParaIndex)*2),
				}
				for i, pi := range field.ParaIndex {
					pi--
					if pi < 0 || pi > len(paraIndex)/2 {
						return nil, fmt.Errorf("illegal paraIndex")
					}
					element.ParaIndex[2*i], element.ParaIndex[2*i+1] = uint64(paraIndex[2*pi]), uint64(paraIndex[2*pi+1])
				}
			}
			elements = append(elements, element)
		}
	}
	return elements, nil
}

const (
	parallelPrefix = "para"

	parallelVersion = 1

	mutexDataAddressMaxNum = 127
	mutexDataNestedMaxNum  = 127
	mutexDataElementMaxNum = 127
)

type mutexElementEncoder struct {
	addressData []byte
	addressMap  map[string]int
	buf         bytes.Buffer
	elementNum  uint8
}

func (encoder *mutexElementEncoder) encode(element *MutexElement) error {
	var (
		addrIndex int
		ok        bool
		num       uint8
	)

	// address is empty: use default address(to)
	if element.Address != "" {
		address := element.Address
		if addrIndex, ok = encoder.addressMap[address]; !ok {
			if len(encoder.addressMap)+1 > mutexDataAddressMaxNum {
				return fmt.Errorf("cross call too much contracts, should not more than %d", mutexDataAddressMaxNum)
			}
			encoder.addressData = append(encoder.addressData, address[:]...)
			addrIndex = len(encoder.addressMap) + 1
			encoder.addressMap[address] = addrIndex
		}
	}

	encoder.buf.Write(common.EncodeUint64(uint64(addrIndex)))
	encoder.buf.Write(common.EncodeUint64(uint64(element.FieldID)))

	if element.FieldID != 0 {
		if len(element.ParaIndex) != 0 {
			num = uint8(len(element.ParaIndex) / 2)
			if num > mutexDataNestedMaxNum {
				return fmt.Errorf("too much nested parameter, should not more than %d", mutexDataNestedMaxNum)
			}
		}

		encoder.buf.Write(common.EncodeUint64(uint64(num)))
		for i := 0; i < int(num); i++ {
			encoder.buf.Write(common.EncodeUint64(element.ParaIndex[2*i]))
			encoder.buf.Write(common.EncodeUint64(element.ParaIndex[2*i+1]))
		}
	} else {
		encoder.buf.WriteByte(0)
	}
	encoder.elementNum++
	return nil
}

// mutex data: paraPrefix:version:payloadPos:addressNum:addressData:elementNum:elementData
func (encoder *mutexElementEncoder) build() []byte {
	var payloadPos int
	prefix := []byte(parallelPrefix)

	version := []byte{parallelVersion}

	addressNum := common.EncodeUint64(uint64(len(encoder.addressMap)))
	addressData := encoder.addressData

	elementNum := common.EncodeUint64(uint64(encoder.elementNum))
	elementData := encoder.buf.Bytes()

	payloadPos = getPayloadPos(parallelPrefixSize + parallelVersionSize + len(addressNum) + len(addressData) + len(elementNum) + len(elementData))
	payloadPosData := common.EncodeUint64(uint64(payloadPos))

	return mergeBytes(prefix, version, payloadPosData, addressNum, addressData, elementNum, elementData)
}

func mergeBytes(data ...[]byte) []byte {
	dataLen := 0
	for _, d := range data {
		dataLen += len(d)
	}
	res := make([]byte, dataLen)
	pos := 0
	for _, d := range data {
		copy(res[pos:], d)
		pos += len(d)
	}
	return res
}

func encodeMutexElement(elements []*MutexElement) (_ []byte, err error) {
	if len(elements) > mutexDataElementMaxNum {
		return nil, fmt.Errorf("too much elements, should not more than %d", mutexDataElementMaxNum)
	}

	encoder := &mutexElementEncoder{
		addressData: make([]byte, 0),
		addressMap:  make(map[string]int, 0),
	}

	for _, element := range elements {
		if err = encoder.encode(element); err != nil {
			return
		}
	}

	return encoder.build(), nil
}

func getPayloadPos(dataLen int) int {
	posLen := 1
	for {
		maxLen := 1<<(posLen*7+1) - 1
		if maxLen >= posLen+dataLen {
			break
		}
		posLen += 1
	}
	return posLen + dataLen
}
