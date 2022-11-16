package fvm

import (
	"fmt"
	"github.com/hyperchain/gosdk/common"
	"github.com/hyperchain/gosdk/fvm/scale"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	const structBin = `{
  "contract": {
    "name": "MyContract",
    "constructor": {
      "input": []
    }
  },
  "methods": [
    {
      "name": "add",
      "input": [
        {
          "type_id": 0
        }
      ],
      "output": [
        {
          "type_id": 0
        }
      ]
    },
    {
      "name": "make_student",
      "input": [
        {
          "type_id": 1
        },
        {
          "type_id": 6
        }
      ],
      "output": [
        {
          "type_id": 0
        }
      ]
    }
  ],
  "types": [
    {
      "id": 0,
      "type": "primitive",
      "primitive": "u64"
    },
    {
      "id": 1,
      "type": "struct",
      "fields": [
        {
          "type_id": 2
        },
        {
          "type_id": 3
        },
        {
          "type_id": 4
        }
      ]
    },
    {
      "id": 2,
      "type": "primitive",
      "primitive": "u32"
    },
    {
      "id": 3,
      "type": "primitive",
      "primitive": "str"
    },
    {
      "id": 4,
      "type": "vec",
      "fields": [
        {
          "type_id": 5
        }
      ]
    },
    {
      "id": 5,
      "type": "vec",
      "fields": [
        {
          "type_id": 3
        }
      ]
    },
    {
      "id": 6,
      "type": "struct",
      "fields": [
        {
          "type_id": 7
        }
      ]
    },
    {
      "id": 7,
      "type": "vec",
      "fields": [
        {
          "type_id": 8
        }
      ]
    },
    {
      "id": 8,
      "type": "array",
      "fields": [
        {
          "type_id": 1
        }
      ],
      "array_len": 10
    }
  ]
}`
	a, err := GenAbi(strings.NewReader(structBin))
	if err != nil {
		t.Error(err)
	}
	t.Run("encode", func(t *testing.T) {
		res, err := Encode(a, "make_student", []interface{}{
			uint32(1),
			"test",
			[]interface{}{
				[]interface{}{
					"hello",
					"world",
				},
			},
		}, []interface{}{
			[]interface{}{},
		})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "306d616b655f73747564656e7401000000107465737404081468656c6c6f14776f726c6400", common.Bytes2Hex(res))
	})
	t.Run("decode", func(t *testing.T) {
		l, err := a.DecodeInput("make_student", common.Hex2Bytes("306d616b655f73747564656e7401000000107465737404081468656c6c6f14776f726c6400"))
		assert.Nil(t, err)
		assert.Equal(t, uint32(1), l.Params[0].GetVal().([]scale.Compact)[0].GetVal())
	})
}

func TestVec(t *testing.T) {
	const hello = `{
  "contract": {
    "name": "StructJsonContract",
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
      "output": []
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
      ]
    },
    {
      "name": "setVecStr",
      "input": [
        {
          "type_id": 2
        }
      ],
      "output": [
        {
          "type_id": 2
        }
      ]
    },
    {
      "name": "setVecInt",
      "input": [
        {
          "type_id": 3
        }
      ],
      "output": [
        {
          "type_id": 3
        }
      ]
    },
    {
      "name": "setVecStruct3",
      "input": [
        {
          "type_id": 5
        }
      ],
      "output": [
        {
          "type_id": 5
        }
      ]
    },
    {
      "name": "setVec2Int",
      "input": [
        {
          "type_id": 14
        }
      ],
      "output": [
        {
          "type_id": 14
        }
      ]
    },
    {
      "name": "setVec2Struct3",
      "input": [
        {
          "type_id": 16
        }
      ],
      "output": [
        {
          "type_id": 16
        }
      ]
    },
    {
      "name": "setVec3Int",
      "input": [
        {
          "type_id": 15
        }
      ],
      "output": [
        {
          "type_id": 15
        }
      ]
    },
    {
      "name": "setVec3Struct3",
      "input": [
        {
          "type_id": 17
        }
      ],
      "output": [
        {
          "type_id": 17
        }
      ]
    },
    {
      "name": "setStruct1",
      "input": [
        {
          "type_id": 8
        }
      ],
      "output": [
        {
          "type_id": 8
        }
      ]
    },
    {
      "name": "setStruct2",
      "input": [
        {
          "type_id": 7
        }
      ],
      "output": [
        {
          "type_id": 7
        }
      ]
    },
    {
      "name": "setStruct3",
      "input": [
        {
          "type_id": 6
        }
      ],
      "output": [
        {
          "type_id": 6
        }
      ]
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
    },
    {
      "id": 2,
      "type": "array",
      "fields": [
        {
          "type_id": 0
        }
      ],
      "array_len": 2
    },
    {
      "id": 3,
      "type": "array",
      "fields": [
        {
          "type_id": 4
        }
      ],
      "array_len": 2
    },
    {
      "id": 4,
      "type": "primitive",
      "primitive": "i64"
    },
    {
      "id": 5,
      "type": "array",
      "fields": [
        {
          "type_id": 6
        }
      ],
      "array_len": 2
    },
    {
      "id": 6,
      "type": "struct",
      "fields": [
        {
          "type_id": 7
        },
        {
          "type_id": 9
        },
        {
          "type_id": 10
        },
        {
          "type_id": 11
        },
        {
          "type_id": 4
        },
        {
          "type_id": 12
        },
        {
          "type_id": 13
        },
        {
          "type_id": 3
        },
        {
          "type_id": 14
        },
        {
          "type_id": 15
        },
        {
          "type_id": 8
        },
        {
          "type_id": 8
        },
        {
          "type_id": 7
        }
      ]
    },
    {
      "id": 7,
      "type": "struct",
      "fields": [
        {
          "type_id": 8
        },
        {
          "type_id": 9
        },
        {
          "type_id": 10
        },
        {
          "type_id": 11
        },
        {
          "type_id": 4
        },
        {
          "type_id": 12
        },
        {
          "type_id": 13
        },
        {
          "type_id": 3
        },
        {
          "type_id": 14
        },
        {
          "type_id": 15
        },
        {
          "type_id": 8
        }
      ]
    },
    {
      "id": 8,
      "type": "struct",
      "fields": [
        {
          "type_id": 9
        },
        {
          "type_id": 10
        },
        {
          "type_id": 11
        },
        {
          "type_id": 4
        },
        {
          "type_id": 12
        },
        {
          "type_id": 13
        },
        {
          "type_id": 3
        },
        {
          "type_id": 14
        },
        {
          "type_id": 15
        }
      ]
    },
    {
      "id": 9,
      "type": "primitive",
      "primitive": "i8"
    },
    {
      "id": 10,
      "type": "primitive",
      "primitive": "i16"
    },
    {
      "id": 11,
      "type": "primitive",
      "primitive": "i32"
    },
    {
      "id": 12,
      "type": "primitive",
      "primitive": "u8"
    },
    {
      "id": 13,
      "type": "primitive",
      "primitive": "u16"
    },
    {
      "id": 14,
      "type": "array",
      "fields": [
        {
          "type_id": 3
        }
      ],
      "array_len": 2
    },
    {
      "id": 15,
      "type": "array",
      "fields": [
        {
          "type_id": 14
        }
      ],
      "array_len": 2
    },
    {
      "id": 16,
      "type": "array",
      "fields": [
        {
          "type_id": 5
        }
      ],
      "array_len": 2
    },
    {
      "id": 17,
      "type": "array",
      "fields": [
        {
          "type_id": 16
        }
      ],
      "array_len": 2
    }
  ]
}`
	a, err := GenAbi(strings.NewReader(hello))
	if err != nil {
		t.Error(err)
	}
	t.Run("encode", func(t *testing.T) {
		res, err := Encode(a, "setVecInt", []interface{}{
			1, 1,
		})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "24736574566563496e7401000000000000000100000000000000", common.Bytes2Hex(res))
	})
	t.Run("decode", func(t *testing.T) {
		res, err := DecodeRet(a, "setVecInt", common.Hex2Bytes("0xffffffffffffff7fffffffffffffff7f"))
		if err != nil {
			t.Error(err)
		}
		//assert.Equal(t, )
		for _, v := range res.Params {
			fmt.Println(v.GetVal())
		}
	})
	t.Run("", func(t *testing.T) {
		res, err := DecodeRet(a, "setVecStruct3", common.Hex2Bytes("0x790903c81000008f13000000000000789001981e000000000000ee0f000000000000100c000000000000ab070000000000004423000000000000b81a000000000000c906000000000000e817000000000000e323000000000000671f0000000000003f040000000000006004000000000000810a000000000000d11c0000000000003d7502d9060000f1060000000000005c7201f60b0000000000002c110000000000003d1e0000000000003f260000000000003e02000000000000de0f00000000000004100000000000006726000000000000f306000000000000320d000000000000fb13000000000000be0c000000000000b509000000000000ef0200000000000061340231040000500b0000000000001aa800981f0000000000002e16000000000000b902000000000000410a000000000000661100000000000051200000000000001d18000000000000db14000000000000d101000000000000ba1000000000000057060000000000000e270000000000000a0800000000000011150000000000006ecc03ca0900001b120000000000006829036402000000000000690f0000000000002914000000000000be1600000000000005170000000000002a010000000000008013000000000000d808000000000000d00b00000000000076010000000000001308000000000000991c000000000000391100000000000074130000000000001adc00930b000038050000000000006e9d01a20c00000000000013040000000000000a0b000000000000051b000000000000c71a0000000000001b0d0000000000008f19000000000000bb0400000000000075050000000000005422000000000000670f000000000000d0250000000000003c1d0000000000000f130000000000006c8001201c0000c40c00000000000068c902a615000000000000881f0000000000007a020000000000009c1300000000000031220000000000000e140000000000008804000000000000db190000000000009714000000000000d626000000000000c2260000000000005808000000000000290a000000000000ea1e000000000000631703a8240000f81e00000000000033b002eb140000000000006517000000000000eb24000000000000b709000000000000f705000000000000251e0000000000000c1d0000000000007401000000000000e80e000000000000661a000000000000d80e0000000000002824000000000000a91600000000000056020000000000000fa801ce0b0000751d0000000000001fa5028212000000000000a10500000000000041240000000000006508000000000000fe18000000000000c80b000000000000e606000000000000ed19000000000000c70b000000000000ea16000000000000bd26000000000000c20f000000000000e11f00000000000012230000000000007eca00811b0000de24000000000000375302411b000000000000cc1900000000000054250000000000005126000000000000d1040000000000000d21000000000000710c0000000000005d19000000000000d51c000000000000be0a0000000000004a1900000000000065120000000000005d2200000000000061070000000000006b6600ac0b0000a90d00000000000027b9002b0c000000000000641b000000000000250a00000000000044170000000000004e26000000000000b1180000000000000708000000000000d40e0000000000008c250000000000005b080000000000003c1d000000000000f82300000000000021050000000000004103000000000000490503e91a00007e0200000000000059b0010208000000000000e81b000000000000f821000000000000e315000000000000ff0f000000000000ce10000000000000fd06000000000000b80c0000000000006702000000000000fc16000000000000e6150000000000009502000000000000ba24000000000000d50900000000000067bf004b240000de1a0000000000000cce013d240000000000008d25000000000000bb15000000000000820a000000000000d21b000000000000ba12000000000000c4210000000000009817000000000000550300000000000024200000000000007a03000000000000f80a000000000000d623000000000000ac060000000000005d1f03a0240000f70c0000000000003dfc01d500000000000000fd030000000000009e1e0000000000006d09000000000000d7180000000000005a1f000000000000b60d0000000000004805000000000000050c0000000000002823000000000000a420000000000000891d000000000000c500000000000000031c00000000000011e50243180000621e000000000000395700d321000000000000e726000000000000d8220000000000003b0a000000000000ff130000000000005512000000000000e400000000000000cf240000000000003a04000000000000341b000000000000fa1b0000000000008019000000000000b304000000000000300d0000000000001e0802fe1e00007e0a0000000000003a68017d1f0000000000003e100000000000002a18000000000000b7040000000000001a00000000000000b726000000000000e62300000000000049140000000000008009000000000000b50a000000000000cc0e000000000000a90b0000000000009e1c0000000000003c060000000000007eae03bd04000063010000000000001a22023503000000000000a301000000000000bc11000000000000ef0b0000000000009c0d000000000000a213000000000000bd180000000000004c170000000000006e18000000000000731b000000000000f624000000000000b1110000000000008401000000000000810f0000000000002adf013a1f0000241a00000000000058a4004a0c000000000000bf12000000000000ef0a0000000000000d110000000000003c1c0000000000005a04000000000000231f0000000000009020000000000000bc0a000000000000981900000000000020000000000000003016000000000000671f000000000000e9110000000000007d9d00b6060000081200000000000068d400c71b0000000000003a23000000000000c20b000000000000171a000000000000de160000000000001b07000000000000bb2100000000000033250000000000001e0900000000000065180000000000007a1e000000000000f92300000000000050190000000000006418000000000000"))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, 1, len(res.Params))
		assert.Equal(t, scale.ArrayName, res.Params[0].GetType())
		for _, v := range res.Params {
			t.Log(scale.GetCompactValue(v))
		}
	})
}

func TestArray(t *testing.T) {
	const he = `{
  "contract": {
    "name": "StructJsonContract",
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
      "output": []
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
      ]
    },
    {
      "name": "setVecStr",
      "input": [
        {
          "type_id": 2
        }
      ],
      "output": [
        {
          "type_id": 2
        }
      ]
    },
    {
      "name": "setVecInt",
      "input": [
        {
          "type_id": 3
        }
      ],
      "output": [
        {
          "type_id": 3
        }
      ]
    },
    {
      "name": "setVecStruct3",
      "input": [
        {
          "type_id": 5
        }
      ],
      "output": [
        {
          "type_id": 5
        }
      ]
    },
    {
      "name": "setVec2Int",
      "input": [
        {
          "type_id": 14
        }
      ],
      "output": [
        {
          "type_id": 14
        }
      ]
    },
    {
      "name": "setVec2Struct3",
      "input": [
        {
          "type_id": 17
        }
      ],
      "output": [
        {
          "type_id": 17
        }
      ]
    },
    {
      "name": "setVec3Int",
      "input": [
        {
          "type_id": 16
        }
      ],
      "output": [
        {
          "type_id": 16
        }
      ]
    },
    {
      "name": "setVec3Struct3",
      "input": [
        {
          "type_id": 19
        }
      ],
      "output": [
        {
          "type_id": 19
        }
      ]
    },
    {
      "name": "setStruct1",
      "input": [
        {
          "type_id": 8
        }
      ],
      "output": [
        {
          "type_id": 8
        }
      ]
    },
    {
      "name": "setStruct2",
      "input": [
        {
          "type_id": 7
        }
      ],
      "output": [
        {
          "type_id": 7
        }
      ]
    },
    {
      "name": "setStruct3",
      "input": [
        {
          "type_id": 6
        }
      ],
      "output": [
        {
          "type_id": 6
        }
      ]
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
    },
    {
      "id": 2,
      "type": "vec",
      "fields": [
        {
          "type_id": 0
        }
      ]
    },
    {
      "id": 3,
      "type": "vec",
      "fields": [
        {
          "type_id": 4
        }
      ]
    },
    {
      "id": 4,
      "type": "primitive",
      "primitive": "i64"
    },
    {
      "id": 5,
      "type": "array",
      "fields": [
        {
          "type_id": 6
        }
      ],
      "array_len": 2
    },
    {
      "id": 6,
      "type": "struct",
      "fields": [
        {
          "type_id": 7
        },
        {
          "type_id": 9
        },
        {
          "type_id": 10
        },
        {
          "type_id": 11
        },
        {
          "type_id": 4
        },
        {
          "type_id": 12
        },
        {
          "type_id": 13
        },
        {
          "type_id": 3
        },
        {
          "type_id": 14
        },
        {
          "type_id": 16
        },
        {
          "type_id": 8
        },
        {
          "type_id": 8
        },
        {
          "type_id": 7
        }
      ]
    },
    {
      "id": 7,
      "type": "struct",
      "fields": [
        {
          "type_id": 8
        },
        {
          "type_id": 9
        },
        {
          "type_id": 10
        },
        {
          "type_id": 11
        },
        {
          "type_id": 4
        },
        {
          "type_id": 12
        },
        {
          "type_id": 13
        },
        {
          "type_id": 3
        },
        {
          "type_id": 14
        },
        {
          "type_id": 16
        },
        {
          "type_id": 8
        }
      ]
    },
    {
      "id": 8,
      "type": "struct",
      "fields": [
        {
          "type_id": 9
        },
        {
          "type_id": 10
        },
        {
          "type_id": 11
        },
        {
          "type_id": 4
        },
        {
          "type_id": 12
        },
        {
          "type_id": 13
        },
        {
          "type_id": 3
        },
        {
          "type_id": 14
        },
        {
          "type_id": 16
        }
      ]
    },
    {
      "id": 9,
      "type": "primitive",
      "primitive": "i8"
    },
    {
      "id": 10,
      "type": "primitive",
      "primitive": "i16"
    },
    {
      "id": 11,
      "type": "primitive",
      "primitive": "i32"
    },
    {
      "id": 12,
      "type": "primitive",
      "primitive": "u8"
    },
    {
      "id": 13,
      "type": "primitive",
      "primitive": "u16"
    },
    {
      "id": 14,
      "type": "array",
      "fields": [
        {
          "type_id": 15
        }
      ],
      "array_len": 2
    },
    {
      "id": 15,
      "type": "array",
      "fields": [
        {
          "type_id": 4
        }
      ],
      "array_len": 2
    },
    {
      "id": 16,
      "type": "vec",
      "fields": [
        {
          "type_id": 14
        }
      ]
    },
    {
      "id": 17,
      "type": "array",
      "fields": [
        {
          "type_id": 18
        }
      ],
      "array_len": 2
    },
    {
      "id": 18,
      "type": "vec",
      "fields": [
        {
          "type_id": 6
        }
      ]
    },
    {
      "id": 19,
      "type": "array",
      "fields": [
        {
          "type_id": 17
        }
      ],
      "array_len": 2
    }
  ]
}`
	a, err := GenAbi(strings.NewReader(he))
	if err != nil {
		t.Error(err)
	}
	t.Run("encode", func(t *testing.T) {
		res, err := DecodeRet(a, "setVec2Struct3", common.Hex2Bytes("0x08790903c81000008f1300000000000078900108981e000000000000ee0f000000000000100c000000000000ab070000000000004423000000000000b81a00000000000008c906000000000000e817000000000000e323000000000000671f0000000000003f040000000000006004000000000000810a000000000000d11c0000000000003d7502d9060000f1060000000000005c720108f60b0000000000002c110000000000003d1e0000000000003f260000000000003e02000000000000de0f0000000000000804100000000000006726000000000000f306000000000000320d000000000000fb13000000000000be0c000000000000b509000000000000ef0200000000000061340231040000500b0000000000001aa80008981f0000000000002e16000000000000b902000000000000410a00000000000066110000000000005120000000000000081d18000000000000db14000000000000d101000000000000ba1000000000000057060000000000000e270000000000000a0800000000000011150000000000006ecc03ca0900001b12000000000000682903086402000000000000690f0000000000002914000000000000be1600000000000005170000000000002a01000000000000088013000000000000d808000000000000d00b00000000000076010000000000001308000000000000991c000000000000391100000000000074130000000000001adc00930b000038050000000000006e9d0108a20c00000000000013040000000000000a0b000000000000051b000000000000c71a0000000000001b0d000000000000088f19000000000000bb0400000000000075050000000000005422000000000000670f000000000000d0250000000000003c1d0000000000000f130000000000006c8001201c0000c40c00000000000068c90208a615000000000000881f0000000000007a020000000000009c1300000000000031220000000000000e14000000000000088804000000000000db190000000000009714000000000000d626000000000000c2260000000000005808000000000000290a000000000000ea1e000000000000631703a8240000f81e00000000000033b00208eb140000000000006517000000000000eb24000000000000b709000000000000f705000000000000251e000000000000080c1d0000000000007401000000000000e80e000000000000661a000000000000d80e0000000000002824000000000000a91600000000000056020000000000000fa801ce0b0000751d0000000000001fa502088212000000000000a10500000000000041240000000000006508000000000000fe18000000000000c80b00000000000008e606000000000000ed19000000000000c70b000000000000ea16000000000000bd26000000000000c20f000000000000e11f00000000000012230000000000007eca00811b0000de2400000000000037530208411b000000000000cc1900000000000054250000000000005126000000000000d1040000000000000d2100000000000008710c0000000000005d19000000000000d51c000000000000be0a0000000000004a1900000000000065120000000000005d2200000000000061070000000000006b6600ac0b0000a90d00000000000027b900082b0c000000000000641b000000000000250a00000000000044170000000000004e26000000000000b118000000000000080708000000000000d40e0000000000008c250000000000005b080000000000003c1d000000000000f82300000000000021050000000000004103000000000000490503e91a00007e0200000000000059b001080208000000000000e81b000000000000f821000000000000e315000000000000ff0f000000000000ce1000000000000008fd06000000000000b80c0000000000006702000000000000fc16000000000000e6150000000000009502000000000000ba24000000000000d50900000000000067bf004b240000de1a0000000000000cce01083d240000000000008d25000000000000bb15000000000000820a000000000000d21b000000000000ba1200000000000008c4210000000000009817000000000000550300000000000024200000000000007a03000000000000f80a000000000000d623000000000000ac060000000000005d1f03a0240000f70c0000000000003dfc0108d500000000000000fd030000000000009e1e0000000000006d09000000000000d7180000000000005a1f00000000000008b60d0000000000004805000000000000050c0000000000002823000000000000a420000000000000891d000000000000c500000000000000031c00000000000011e50243180000621e00000000000039570008d321000000000000e726000000000000d8220000000000003b0a000000000000ff13000000000000551200000000000008e400000000000000cf240000000000003a04000000000000341b000000000000fa1b0000000000008019000000000000b304000000000000300d0000000000001e0802fe1e00007e0a0000000000003a6801087d1f0000000000003e100000000000002a18000000000000b7040000000000001a00000000000000b72600000000000008e62300000000000049140000000000008009000000000000b50a000000000000cc0e000000000000a90b0000000000009e1c0000000000003c060000000000007eae03bd04000063010000000000001a2202083503000000000000a301000000000000bc11000000000000ef0b0000000000009c0d000000000000a21300000000000008bd180000000000004c170000000000006e18000000000000731b000000000000f624000000000000b1110000000000008401000000000000810f0000000000002adf013a1f0000241a00000000000058a400084a0c000000000000bf12000000000000ef0a0000000000000d110000000000003c1c0000000000005a0400000000000008231f0000000000009020000000000000bc0a000000000000981900000000000020000000000000003016000000000000671f000000000000e9110000000000007d9d00b6060000081200000000000068d40008c71b0000000000003a23000000000000c20b000000000000171a000000000000de160000000000001b0700000000000008bb2100000000000033250000000000001e0900000000000065180000000000007a1e000000000000f9230000000000005019000000000000641800000000000008790903c81000008f1300000000000078900108981e000000000000ee0f000000000000100c000000000000ab070000000000004423000000000000b81a00000000000008c906000000000000e817000000000000e323000000000000671f0000000000003f040000000000006004000000000000810a000000000000d11c0000000000003d7502d9060000f1060000000000005c720108f60b0000000000002c110000000000003d1e0000000000003f260000000000003e02000000000000de0f0000000000000804100000000000006726000000000000f306000000000000320d000000000000fb13000000000000be0c000000000000b509000000000000ef0200000000000061340231040000500b0000000000001aa80008981f0000000000002e16000000000000b902000000000000410a00000000000066110000000000005120000000000000081d18000000000000db14000000000000d101000000000000ba1000000000000057060000000000000e270000000000000a0800000000000011150000000000006ecc03ca0900001b12000000000000682903086402000000000000690f0000000000002914000000000000be1600000000000005170000000000002a01000000000000088013000000000000d808000000000000d00b00000000000076010000000000001308000000000000991c000000000000391100000000000074130000000000001adc00930b000038050000000000006e9d0108a20c00000000000013040000000000000a0b000000000000051b000000000000c71a0000000000001b0d000000000000088f19000000000000bb0400000000000075050000000000005422000000000000670f000000000000d0250000000000003c1d0000000000000f130000000000006c8001201c0000c40c00000000000068c90208a615000000000000881f0000000000007a020000000000009c1300000000000031220000000000000e14000000000000088804000000000000db190000000000009714000000000000d626000000000000c2260000000000005808000000000000290a000000000000ea1e000000000000631703a8240000f81e00000000000033b00208eb140000000000006517000000000000eb24000000000000b709000000000000f705000000000000251e000000000000080c1d0000000000007401000000000000e80e000000000000661a000000000000d80e0000000000002824000000000000a91600000000000056020000000000000fa801ce0b0000751d0000000000001fa502088212000000000000a10500000000000041240000000000006508000000000000fe18000000000000c80b00000000000008e606000000000000ed19000000000000c70b000000000000ea16000000000000bd26000000000000c20f000000000000e11f00000000000012230000000000007eca00811b0000de2400000000000037530208411b000000000000cc1900000000000054250000000000005126000000000000d1040000000000000d2100000000000008710c0000000000005d19000000000000d51c000000000000be0a0000000000004a1900000000000065120000000000005d2200000000000061070000000000006b6600ac0b0000a90d00000000000027b900082b0c000000000000641b000000000000250a00000000000044170000000000004e26000000000000b118000000000000080708000000000000d40e0000000000008c250000000000005b080000000000003c1d000000000000f82300000000000021050000000000004103000000000000490503e91a00007e0200000000000059b001080208000000000000e81b000000000000f821000000000000e315000000000000ff0f000000000000ce1000000000000008fd06000000000000b80c0000000000006702000000000000fc16000000000000e6150000000000009502000000000000ba24000000000000d50900000000000067bf004b240000de1a0000000000000cce01083d240000000000008d25000000000000bb15000000000000820a000000000000d21b000000000000ba1200000000000008c4210000000000009817000000000000550300000000000024200000000000007a03000000000000f80a000000000000d623000000000000ac060000000000005d1f03a0240000f70c0000000000003dfc0108d500000000000000fd030000000000009e1e0000000000006d09000000000000d7180000000000005a1f00000000000008b60d0000000000004805000000000000050c0000000000002823000000000000a420000000000000891d000000000000c500000000000000031c00000000000011e50243180000621e00000000000039570008d321000000000000e726000000000000d8220000000000003b0a000000000000ff13000000000000551200000000000008e400000000000000cf240000000000003a04000000000000341b000000000000fa1b0000000000008019000000000000b304000000000000300d0000000000001e0802fe1e00007e0a0000000000003a6801087d1f0000000000003e100000000000002a18000000000000b7040000000000001a00000000000000b72600000000000008e62300000000000049140000000000008009000000000000b50a000000000000cc0e000000000000a90b0000000000009e1c0000000000003c060000000000007eae03bd04000063010000000000001a2202083503000000000000a301000000000000bc11000000000000ef0b0000000000009c0d000000000000a21300000000000008bd180000000000004c170000000000006e18000000000000731b000000000000f624000000000000b1110000000000008401000000000000810f0000000000002adf013a1f0000241a00000000000058a400084a0c000000000000bf12000000000000ef0a0000000000000d110000000000003c1c0000000000005a0400000000000008231f0000000000009020000000000000bc0a000000000000981900000000000020000000000000003016000000000000671f000000000000e9110000000000007d9d00b6060000081200000000000068d40008c71b0000000000003a23000000000000c20b000000000000171a000000000000de160000000000001b0700000000000008bb2100000000000033250000000000001e0900000000000065180000000000007a1e000000000000f92300000000000050190000000000006418000000000000"))
		if err != nil {
			t.Error(err)
		}
		for _, v := range res.Params {

			t.Log(scale.GetCompactValue(v))
		}
	})

	t.Run("decode", func(t *testing.T) {
		res, err := DecodeRet(a, "setVec3Int", common.Hex2Bytes("0x00"))
		if err != nil {
			t.Error(err)
		}
		for _, v := range res.Params {
			t.Log(scale.GetCompactValue(v))
		}
	})
}

func TestName(t *testing.T) {
	abiFirst := `{
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

	abiSecond := `{
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

	addr1 := "0x02"
	addr2 := "0x01"

	abi1, err := scale.JSON(strings.NewReader(abiFirst))
	assert.Nil(t, err)

	abi2, err := scale.JSON(strings.NewReader(abiSecond))
	assert.Nil(t, err)

	abiMap := scale.NewAbiMap()
	abiMap.SetAbi(addr1, &abi1)
	abiMap.SetAbi(addr2, &abi2)

	parallel, err := EncodeParallel(abiMap, addr1, "set_hash", "hello", "world")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "0x706172611701020230783032010101090e3078303103010003000400207365745f686173681468656c6c6f14776f726c64", common.ToHex(parallel))
}
