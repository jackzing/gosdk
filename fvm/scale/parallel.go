package scale

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hyperchain/gosdk/common"
	"github.com/hyperchain/gosdk/common/parallel"
)

type MutexElement struct {
	StoreId       int
	ParaIndexList []int
}

func Parallel(abiMap *AbiMap, addr string, method string, params ...interface{}) ([]byte, error) {
	curAbi, err := abiMap.GetAbi(addr)
	if err != nil {
		return nil, err
	}
	// payload
	payload, err := curAbi.Encode(method, params...)
	if err != nil {
		return nil, err
	}

	curMethod, err := curAbi.GetMethod(method)
	if err != nil {
		return nil, err
	}
	visit := make(map[string]bool)
	elementMap, addrs, err := parallelMethod(abiMap, addr, curMethod, visit, params...)
	if err != nil {
		return nil, err
	}
	if len(elementMap) == 0 && curMethod.ParallelLevel == parallel.GlobalMutex {
		return payload, nil
	}

	var buf bytes.Buffer
	buf.Write(common.EncodeInt32(ParallelVersion))
	buf.Write(common.EncodeInt32(int32(curMethod.ParallelLevel)))
	buf.Write(common.EncodeInt32(int32(len(elementMap))))

	for _, ad := range addrs {
		if v, ok := elementMap[ad]; ok {
			addrEncode := []byte(ad)
			buf.Write(addrEncode[len(addrEncode)-4:])
			buf.Write(common.EncodeInt32(int32(len(v))))
			for _, k := range v {
				buf.Write(common.EncodeInt32(int32(k.StoreId)))
				buf.Write(common.EncodeInt32(int32(len(k.ParaIndexList))))
				for _, i := range k.ParaIndexList {
					tempAbi, err := abiMap.GetAbi(ad)
					if err != nil {
						return nil, err
					}
					buf.Write(common.EncodeInt32(int32(tempAbi.curParaIndex[i-1].start)))
					buf.Write(common.EncodeInt32(int32(tempAbi.curParaIndex[i-1].end)))
				}
			}
		}
	}

	var res bytes.Buffer
	// total parallel byte len
	res.Write([]byte("para"))
	res.Write(common.EncodeInt32(int32(buf.Len())))
	res.Write(buf.Bytes())
	res.Write(payload)

	return res.Bytes(), nil
}

func parallelMethod(abiMap *AbiMap, addr string, curMethod *Method, curVisit map[string]bool, params ...interface{}) (map[string][]MutexElement, []string, error) {
	var addrs []string
	addrs = append(addrs, addr)
	mutexElementMap := make(map[string][]MutexElement)
	// bfs search to combine mutex variable
	// first to count current function mutex storeField
	storeFields, err := parallelMutexFields(curMethod, curVisit, params...)
	if err != nil {
		return nil, nil, err
	}
	storeFieldTotal := len(storeFields)
	if storeFieldTotal > 128 {
		return nil, nil, errors.New("too many store field")
	}
	if len(storeFields) != 0 {
		mutexElementMap[addr] = append(mutexElementMap[addr], storeFields...)
	}

	// second to count call method mutex storeField
	var methodNames []string
	for _, v := range curMethod.MethodCalls {
		methodNames = append(methodNames, v)
	}
	var crossMethods []MutexCall
	for _, v := range curMethod.MutexCalls {
		crossMethods = append(crossMethods, v)
	}

	methodNum := 0
	methodVisit := make(map[string]bool)
	for len(methodNames) != 0 {
		if methodNum > 128 {
			return nil, nil, errors.New("too many method call")
		}
		methodNum += 1
		tempMethodName := methodNames[0]
		if methodVisit[tempMethodName] {
			methodNames = methodNames[1:]
			continue
		}
		methodVisit[tempMethodName] = true
		tempAbi, err := abiMap.GetAbi(addr)
		if err != nil {
			return nil, nil, err
		}
		tempCurMethod, err := tempAbi.GetMethod(tempMethodName)
		if err != nil {
			return nil, nil, err
		}
		t, err := parallelMutexFields(tempCurMethod, curVisit, params...)
		if err != nil {
			return nil, nil, err
		}
		storeFieldTotal += len(t)
		if storeFieldTotal > 128 {
			return nil, nil, errors.New("too many store field")
		}
		if len(t) != 0 {
			mutexElementMap[addr] = append(mutexElementMap[addr], t...)
		}
		methodNames = methodNames[1:]
		for _, v := range tempCurMethod.MethodCalls {
			methodNames = append(methodNames, v)
		}
		for _, v := range tempCurMethod.MutexCalls {
			crossMethods = append(crossMethods, v)
		}
	}

	// cross call
	if len(crossMethods) > 128 {
		return nil, nil, errors.New("too many cross call method")
	}
	for len(crossMethods) != 0 {
		tempCross := crossMethods[0]
		crossAddr := tempCross.Address
		if tempCross.Address != "" {
			addrs = append(addrs, crossAddr)
		} else if tempCross.AddressIndex > 0 {
			add, err := convertPrimitive("string", params[tempCross.AddressIndex])
			if err != nil {
				return nil, nil, err
			}
			crossAddr = add.GetVal().(string)
			addrs = append(addrs, crossAddr)
		} else {
			return nil, nil, errors.New("error abi, check it")
		}

		tempVisit := make(map[string]bool)
		methodNum += len(tempCross.Methods)
		if methodNum > 128 {
			return nil, nil, errors.New("too many method call")
		}
		// cross call method
		for _, v := range tempCross.Methods {
			tempAbi, err := abiMap.GetAbi(crossAddr)
			if err != nil {
				return nil, nil, err
			}
			tempCurMethod, err := tempAbi.GetMethod(v)
			if err != nil {
				return nil, nil, err
			}
			elements, err := parallelMutexFields(tempCurMethod, tempVisit, nil)
			if err != nil {
				return nil, nil, err
			}
			storeFieldTotal += len(elements)
			if storeFieldTotal > 128 {
				return nil, nil, errors.New("too many store field")
			}
			if len(elements) != 0 {
				mutexElementMap[crossAddr] = append(mutexElementMap[crossAddr], elements...)
			}

			var tempMethodCall = tempCurMethod.MethodCalls
			crossMethodVisit := make(map[string]bool)
			for len(tempMethodCall) != 0 {
				tempMethodName := tempMethodCall[0]
				if crossMethodVisit[tempMethodName] {
					tempMethodCall = tempMethodCall[1:]
					continue
				}
				crossMethodVisit[tempMethodName] = true
				tempCurMethod2, err := tempAbi.GetMethod(tempMethodName)
				if err != nil {
					return nil, nil, err
				}
				t, err := parallelMutexFields(tempCurMethod2, tempVisit, nil)
				if err != nil {
					return nil, nil, err
				}
				storeFieldTotal += len(t)
				if storeFieldTotal > 128 {
					return nil, nil, errors.New("too many store field")
				}
				if len(t) != 0 {
					mutexElementMap[crossAddr] = append(mutexElementMap[crossAddr], t...)
				}
				tempMethodCall = tempMethodCall[1:]
				for _, j := range tempCurMethod2.MethodCalls {
					tempMethodCall = append(tempMethodCall, j)
				}
			}
		}
		crossMethods = crossMethods[1:]
	}

	return mutexElementMap, addrs, nil
}

func parallelMutexFields(curMethod *Method, curVisit map[string]bool, params ...interface{}) ([]MutexElement, error) {
	if len(curMethod.MutexFields) > 128 {
		return nil, errors.New("too many store field")
	}
	var mutexElement []MutexElement
	if curMethod.ParallelLevel == parallel.ContractMutex {
		mutexElement = append(mutexElement, MutexElement{StoreId: 0, ParaIndexList: nil})
	}
	for _, v := range curMethod.MutexFields {
		for _, i := range v.ParallelIndex {
			if i > len(curMethod.Input) {
				return nil, fmt.Errorf("wrong param index in para_index, filed_name is %s", v.FieldName)
			}
		}
		if len(v.ParallelIndex) != 0 && len(params) != 0 {
			mutexElement = append(mutexElement, MutexElement{
				StoreId:       v.FieldId,
				ParaIndexList: v.ParallelIndex,
			})
		} else {
			if !curVisit[v.FieldName] {
				curVisit[v.FieldName] = true
				mutexElement = append(mutexElement, MutexElement{
					StoreId:       v.FieldId,
					ParaIndexList: nil,
				})
			}
		}
	}
	return mutexElement, nil
}
