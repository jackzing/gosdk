package hvm

import (
	"fmt"
	"testing"

	"github.com/hyperchain/go-hpc-common/types"
	"github.com/jackzing/gosdk/common"
	"github.com/stretchr/testify/assert"
)

func TestAbi_GetMethodAbi(t *testing.T) {
	methodInvokeAbi := "../hvmtestfile/methodInvoke/hvm.abi"
	abiJSON, err := common.ReadFileAsString(methodInvokeAbi)
	assert.Nil(t, err)

	abi, err := GenAbi(abiJSON)
	assert.Nil(t, err)

	_, err = abi.GetMethodAbi("Hello")
	assert.Nil(t, err)

	methodAbi2, err := abi.GetMethodAbi("Hello()")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(methodAbi2.Inputs))

	methodAbi3, err := abi.GetMethodAbi("Hello(java.lang.String)")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(methodAbi3.Inputs))

	methodAbi4, err := abi.GetMethodAbi("Hello(int,java.lang.String)")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(methodAbi4.Inputs))
}

func TestStaticAnalyze(t *testing.T) {
	abiMap := NewAbiMap()

	t.Run("parallel invoke directly nestedMap", func(t *testing.T) {
		addrStr := "0x81df5d4ef1c94c7629f786fce261d84b4edf73b9"
		addr := types.HexToAddress(addrStr)
		abiPath := "../hvmtestfile/staticanalyze/parallel-contract.abi"
		abiJSON, err := common.ReadFileAsString(abiPath)
		assert.Nil(t, err)
		abi, err := GenAbi(abiJSON)
		assert.Nil(t, err)
		abiMap.SetAbi(addrStr, abi)

		beanName := "transfer"
		var paray = []struct {
			f1 string
			f2 string
		}{
			{"a1", "a2"},
			{"a3", "a4"},
			{"a1", "a3"},
			{"a10", "a20"},
			{"a11", "a20"},
			{"a3", "a5"},
		}
		for _, p := range paray {
			payload, err := GenParallelPayload(abiMap, addr, beanName, true, p.f1, p.f2, 100)
			assert.Nil(t, err)

			res, err := parseMutexData(payload, true)
			assert.Nil(t, err)
			assert.Len(t, res, 2)
			assert.Equal(t, res[0].addrIndex, 0)
			assert.Equal(t, res[0].fieldID, 1)
			assert.Equal(t, res[0].params[0], p.f1)
			assert.Equal(t, res[1].params[0], p.f2)
		}
	})

	t.Run("parallel invoke bean", func(t *testing.T) {
		addrStr := "0x7629f786fce261d84b4edf73b981df5d4ef1c94c"
		addr := types.HexToAddress(addrStr)
		abiPath := "../hvmtestfile/staticanalyze/parallel-contract.abi"
		abiJSON, err := common.ReadFileAsString(abiPath)
		assert.Nil(t, err)
		abi, err := GenAbi(abiJSON)
		assert.Nil(t, err)
		abiMap.SetAbi(addrStr, abi)

		beanName := "demo.invoke.Transfer"
		var paray = []struct {
			f1 string
			f2 string
		}{
			{"a1", "a2"},
			{"a3", "a4"},
			{"a1", "a3"},
			{"a10", "a20"},
			{"a11", "a20"},
			{"a3", "a5"},
		}
		for _, p := range paray {
			payload, err := GenParallelPayload(abiMap, addr, beanName, false, p.f1, p.f2, 100)
			assert.Nil(t, err)

			res, err := parseMutexData(payload, false)
			assert.Nil(t, err)
			assert.Len(t, res, 2)
			assert.Equal(t, res[0].addrIndex, 0)
			assert.Equal(t, res[0].fieldID, 1)
			assert.Equal(t, res[0].params[0], p.f1)
			assert.Equal(t, res[1].params[0], p.f2)

		}
	})

	t.Run("parallel invoke directly hyperTable", func(t *testing.T) {
		addrStr2 := "0xbfa5bd992e3eb123c8b86ebe892099d4e9efb783"
		addr2 := types.HexToAddress(addrStr2)
		abiPath := "../hvmtestfile/staticanalyze/sethash.abi"
		abiJSON, err := common.ReadFileAsString(abiPath)
		assert.Nil(t, err)
		abi, err := GenAbi(abiJSON)
		assert.Nil(t, err)
		abiMap.SetAbi(addrStr2, abi)

		var paray = []struct {
			k1 string
			k2 string
			k3 string
			v  string
		}{
			{"ka1", "ka2", "ka3", "v"},
			{"ka1", "ka2222222222", "", "vv"},
			{"ka1", "", "ka3", ""},
		}
		for _, p := range paray {
			payload, err := GenParallelPayload(abiMap, addr2, "register", true, p.k1, p.k2, p.k3, p.v)
			assert.Nil(t, err)

			res, err := parseMutexData(payload, true)
			assert.Nil(t, err)
			assert.Len(t, res, 1)
			assert.Equal(t, res[0].addrIndex, 0)
			assert.Equal(t, res[0].fieldID, 1)
			assert.Len(t, res[0].params, 3)
			assert.Equal(t, res[0].params[0], p.k1)
			assert.Equal(t, res[0].params[1], p.k2)
			assert.Equal(t, res[0].params[2], p.k3)
		}
	})

	t.Run("parallel invoke directly cross call", func(t *testing.T) {
		abiPath0 := "../hvmtestfile/staticanalyze/parallel-contract.abi"
		abiJSON0, err := common.ReadFileAsString(abiPath0)
		assert.Nil(t, err)
		abi0, err := GenAbi(abiJSON0)
		assert.Nil(t, err)
		addrStr0 := "0x916a99c5bb4a3c39f7c0f579c9df5b9931af49a8"
		addr0 := types.HexToAddress(addrStr0)
		abiMap.SetAbi(addrStr0, abi0)

		abiPath1 := "../hvmtestfile/staticanalyze/parallel-contract-a.abi"
		abiJSON1, err := common.ReadFileAsString(abiPath1)
		assert.Nil(t, err)
		abi1, err := GenAbi(abiJSON1)
		assert.Nil(t, err)
		addrStr1 := "0x916a3c39f7c0f579c9df5b9931af49a899c5bb4a"
		//addr1 := types.HexToAddress(addrStr1)
		abiMap.SetAbi(addrStr1, abi1)

		key, val := "key", "value"
		payload, err := GenParallelPayload(abiMap, addr0, "crossSetAMapKV", true, key, val)
		assert.Nil(t, err)
		res, err := parseMutexData(payload, true)
		assert.Nil(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, res[0].addrIndex, 1)
		assert.Equal(t, res[0].fieldID, 1)
		assert.Len(t, res[0].params, 0)
	})

	t.Run("compare parallel invoke bean and invoke directly", func(t *testing.T) {
		addrStr := "0x7629f786fce261d84b4edf73b981df5d4ef1c94c"
		addr := types.HexToAddress(addrStr)
		abiPath := "../hvmtestfile/staticanalyze/parallel-contract.abi"
		abiJSON, err := common.ReadFileAsString(abiPath)
		assert.Nil(t, err)
		abi, err := GenAbi(abiJSON)
		assert.Nil(t, err)
		abiMap.SetAbi(addrStr, abi)

		type Man struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		testParams := [][]interface{}{
			{
				"true", "true", "13", "23", "15", "26", "17", "28", "119", "1110", // 1~10
				"1.11", "1.12", "2.13", "2.14", "s", "a", "Sp17", // 11~17
				"{\"name\":\"Jack\",\"age\":18}", "[1,19]", "[2,20]", // 18~20
				"[{\"name\":\"QQQ\",\"age\":21},{\"name\":\"PPP\",\"age\":21}]", "[\"1s22\",\"2s22\"]", // 21, 22
				"[{\"name\":\"ZZZ\",\"age\":23},{\"name\":\"XXX\",\"age\":23}]", // 23
				"{\"CCC\":{\"name\":\"CCC\",\"age\":24},\"DDD\":{\"name\":\"DDD\",\"age\":24}}}",
			},
			{
				"true", "true", "13", "23", "15", "26", "17", "28", "119", "1110", // 1~10
				"1.11", "1.12", "2.13", "2.14", "s", "a", "Sp17", // 11~17
				Man{"Jack", 18}, []int{1, 19}, []int{2, 20}, // 18~20
				[]Man{{"QQQ", 21}, {"PPP", 21}}, []string{"1s22", "2s22"}, // 21, 22
				[]Man{{"ZZZ", 23}, {"XXX", 23}}, // 23
				map[string]Man{"CCC": {"CCC", 24}, "DDD": {"DDD", 24}},
			},
			{
				true, true, 13, 23, 15, 26, 17, 28, 119, 1110, // 1~10
				1.11, 1.12, 2.13, 2.14, "s", "a", "Sp17", // 11~17
				Man{"Jack", 18}, []int{1, 19}, []int{2, 20}, // 18~20
				[]Man{{"QQQ", 21}, {"PPP", 21}}, []string{"1s22", "2s22"}, // 21, 22
				[]Man{{"ZZZ", 23}, {"XXX", 23}}, // 23
				map[string]Man{"CCC": {"CCC", 24}, "DDD": {"DDD", 24}},
			},
		}
		parallelParameterLen := 17

		// not support Struct, Map, List, Array as parallel parameter
		for _, params := range testParams {
			payload, err := GenParallelPayload(abiMap, addr, "complexParameter", true, params...)
			assert.Nil(t, err)
			res1, err := parseMutexData(payload, true)
			assert.Nil(t, err)
			//fmt.Println(res1[0])
			assert.Len(t, res1[0].params, parallelParameterLen)

			payload, err = GenParallelPayload(abiMap, addr, "demo.invoke.ComplexInvokeBean", false, params...)
			assert.Nil(t, err)
			res2, err := parseMutexData(payload, false)
			assert.Nil(t, err)
			//fmt.Println(res2[0])
			assert.Len(t, res2[0].params, parallelParameterLen)

			for i := 0; i < parallelParameterLen; i++ {
				assert.Equal(t, res1[0].params[i], res2[0].params[i], fmt.Sprintf("para %d", i))
			}
		}
	})
}

type mutexData struct {
	addrIndex int
	fieldID   int
	params    []string
}

func parseMutexData(payload []byte, methodBean bool) ([]*mutexData, error) {
	var index uint64 = 5
	// read pos
	pos, n, err := common.DecodeUint64(common.NewSliceBytes(payload[index:]))
	if err != nil {
		return nil, err
	}

	parallel := payload[:pos]
	payload = payload[pos:]

	index += n
	// read addr
	addrNum, n, err := common.DecodeUint64(common.NewSliceBytes(parallel[index:]))
	if err != nil {
		return nil, err
	}
	index += n
	index = index + 4*addrNum
	// read eleNum
	eleNum, n, err := common.DecodeUint64(common.NewSliceBytes(parallel[index:]))
	if err != nil {
		return nil, err
	}
	index += n
	res := make([]*mutexData, 0)
	for i := 0; i < int(eleNum); i++ {
		// read addrIndex
		addrIndex, n, err := common.DecodeUint64(common.NewSliceBytes(parallel[index:]))
		if err != nil {
			return nil, err
		}

		index += n
		// read fieldID
		fieldID, n, err := common.DecodeUint64(common.NewSliceBytes(parallel[index:]))
		if err != nil {
			return nil, err
		}
		index += n

		paramNum, n, err := common.DecodeUint64(common.NewSliceBytes(parallel[index:]))
		if err != nil {
			return nil, err
		}
		index += n

		md := &mutexData{
			addrIndex: int(addrIndex),
			fieldID:   int(fieldID),
			params:    make([]string, 0),
		}

		for i := 0; i < int(paramNum); i++ {
			var (
				pos1 uint64
				pos2 uint64
			)

			pos1, n, err = common.DecodeUint64(common.NewSliceBytes(parallel[index:]))
			if err != nil {
				return nil, err
			}
			index += n
			pos2, n, err = common.DecodeUint64(common.NewSliceBytes(parallel[index:]))
			if err != nil {
				return nil, err
			}
			index += n

			if methodBean {
				pos1 -= 4
				pos2 -= 4
			}

			md.params = append(md.params, string(payload[pos1:pos2]))
		}
		res = append(res, md)
	}
	return res, nil
}
