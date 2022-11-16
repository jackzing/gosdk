package scale

import (
	"github.com/hyperchain/gosdk/common"
	"github.com/stretchr/testify/assert"
	_ "net/http/pprof"
	"strings"
	"testing"
)

const (
	abi_1 = `{
  "contract": {
    "name": "SetHash",
    "constructor": {
      "input": []
    }
  },
  "methods": [
    {
      "name": "set_hash",
      "input": [
        {
          "type_id": 0
        },
        {
          "type_id": 0
        }
      ],
      "output": [],
      "parallel_level": 2,
      "mutex_fields": [
        {
          "field_name": "key",
          "field_id": 1,
          "parallel_index": [
            1
          ]
        }
      ],
      "mutex_calls": [],
      "method_calls": [
        "get_hash"
      ]
    },
    {
      "name": "get_hash",
      "input": [
        {
          "type_id": 0
        }
      ],
      "output": [
        {
          "type_id": 1
        }
      ],
      "parallel_level": 2,
      "mutex_fields": [],
      "mutex_calls": [
        {
          "address_index": 1,
          "address": "0x01",
          "methods": [
            "m1",
            "m2"
          ]
        }
      ],
      "method_calls": []
    }
  ],
  "types": [
    {
      "id": 0,
      "type": "primitive",
      "primitive": "str"
    },
    {
      "id": 1,
      "type": "primitive",
      "primitive": "str"
    }
  ]
}`

	abi_2 = `{
  "contract": {
    "name": "Fib",
    "constructor": {
      "input": []
    }
  },
  "methods": [
    {
      "name": "fib",
      "input": [
        {
          "type_id": 0
        }
      ],
      "output": [
        {
          "type_id": 0
        }
      ],
      "parallel_level": 0,
      "mutex_fields": [],
      "mutex_calls": [],
      "method_calls": []
    },
    {
      "name": "m1",
      "input": [],
      "output": [],
      "parallel_level": 2,
      "mutex_fields": [
        {
          "field_name": "key2",
          "field_id": 1,
          "parallel_index": []
        },
        {
          "field_name": "key2",
          "field_id": 2,
          "parallel_index": []
        },
        {
          "field_name": "key3",
          "field_id": 3,
          "parallel_index": []
        }
      ],
      "mutex_calls": [],
      "method_calls": []
    },
    {
      "name": "m2",
      "input": [],
      "output": [],
      "parallel_level": 2,
      "mutex_fields": [
        {
          "field_name": "map1",
          "field_id": 4,
          "parallel_index": []
        }
      ],
      "mutex_calls": [],
      "method_calls": [
        "m1"
      ]
    }
  ],
  "types": [
    {
      "id": 0,
      "type": "primitive",
      "primitive": "u64"
    }
  ]
}`
)

func TestName(t *testing.T) {

	addr1 := "0x02"
	addr2 := "0x01"

	abi1, err := JSON(strings.NewReader(abi_1))
	assert.Nil(t, err)

	abi2, err := JSON(strings.NewReader(abi_2))
	assert.Nil(t, err)

	abiMap := NewAbiMap()
	abiMap.SetAbi(addr1, &abi1)
	abiMap.SetAbi(addr2, &abi2)

	parallel, err := Parallel(abiMap, addr1, "set_hash", "hello", "world")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "0x706172611701020230783032010101090e3078303103010003000400207365745f686173681468656c6c6f14776f726c64", common.ToHex(parallel))
}

func BenchmarkName(b *testing.B) {
	addr1 := "0x02"
	addr2 := "0x01"

	abi1, err := JSON(strings.NewReader(abi_1))
	assert.Nil(b, err)

	abi2, err := JSON(strings.NewReader(abi_2))
	assert.Nil(b, err)
	abiMap := NewAbiMap()
	abiMap.SetAbi(addr1, &abi1)
	abiMap.SetAbi(addr2, &abi2)

	for i := 0; i < b.N; i++ {
		_, err := Parallel(abiMap, addr1, "set_hash", "hello", "world")
		if err != nil {
			b.Error(err)
		}
	}
}
