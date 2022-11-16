package scale

import (
	"github.com/jackzing/gosdk/common"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

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

const VecJson = `{
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
          "type_id": 2
        },
        {
          "type_id": 0
        }
      ],
      "output": []
    },
    {
      "name": "set_hash2",
      "input": [
        {
          "type_id": 2
        },
        {
          "type_id": 1
        }
      ],
      "output": []
    }
  ],
  "types": [
    {
      "id": 0,
      "type": "vec",
      "fields": [
        {
          "type_id": 1
        }
      ]
    },
    {
      "id": 1,
      "type": "vec",
      "fields": [
        {
          "type_id": 2
        }
      ]
    },
    {
      "id": 2,
      "type": "primitive",
      "primitive": "String"
    }
  ]
}`

const DeployParam = `{
  "contract": {
    "name": "AbiTest",
    "constructor": {
      "input": [
        {
          "type_id": 0
        }
      ]
    }
  },
  "methods": [
    {
      "name": "set_hash",
      "input": [
        {
          "type_id": 1
        },
        {
          "type_id": 1
        }
      ],
      "output": []
    },
    {
      "name": "get_hash",
      "input": [
        {
          "type_id": 1
        }
      ],
      "output": [
        {
          "type_id": 2
        }
      ]
    },
    {
      "name": "get_enum",
      "input": [],
      "output": [
        {
          "type_id": 3
        }
      ]
    },
    {
      "name": "set_enum",
      "input": [
        {
          "type_id": 3
        }
      ],
      "output": []
    },
    {
      "name": "get_option",
      "input": [],
      "output": [
        {
          "type_id": 5
        }
      ]
    },
    {
      "name": "set_option",
      "input": [
        {
          "type_id": 5
        }
      ],
      "output": []
    },
    {
      "name": "get_tuple",
      "input": [],
      "output": [
        {
          "type_id": 6
        }
      ]
    },
    {
      "name": "set_tuple",
      "input": [
        {
          "type_id": 4
        }
      ],
      "output": []
    },
    {
      "name": "get_hello",
      "input": [],
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
      "primitive": "u32"
    },
    {
      "id": 1,
      "type": "primitive",
      "primitive": "str"
    },
    {
      "id": 2,
      "type": "primitive",
      "primitive": "str"
    },
    {
      "id": 3,
      "type": "enum",
      "variants": [
        [
          {
            "type_id": 0
          },
          {
            "type_id": 4
          }
        ],
        [
          {
            "type_id": 0
          }
        ]
      ]
    },
    {
      "id": 4,
      "type": "tuple",
      "fields": [
        {
          "type_id": 0
        },
        {
          "type_id": 0
        }
      ]
    },
    {
      "id": 5,
      "type": "enum",
      "variants": [
        [],
        [
          {
            "type_id": 0
          }
        ]
      ]
    },
    {
      "id": 6,
      "type": "tuple",
      "fields": [
        {
          "type_id": 0
        },
        {
          "type_id": 7
        }
      ]
    },
    {
      "id": 7,
      "type": "primitive",
      "primitive": "i32"
    }
  ]
}`

func TestAbi_Encode(t *testing.T) {
	a, err := JSON(strings.NewReader(VecJson))
	if err != nil {
		t.Error(err)
	}
	res, err := a.EncodeCompact("set_hash", &CompactString{Val: "key"}, &CompactVec{Val: []Compact{
		&CompactVec{
			Val: []Compact{
				&CompactString{Val: "hello"},
				&CompactString{Val: "world"},
			},
			NextList: []TypeString{StringName},
		},
		&CompactVec{
			Val: []Compact{
				&CompactString{Val: "hello"},
				&CompactString{Val: "world"},
			},
			NextList: []TypeString{StringName},
		},
	}, NextList: []TypeString{VecName, StringName}})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "207365745f686173680c6b657908081468656c6c6f14776f726c64081468656c6c6f14776f726c64", common.Bytes2Hex(res))
}

func TestAbi_Decode(t *testing.T) {
	a, err := JSON(strings.NewReader(VecJson))
	if err != nil {
		t.Error(err)
	}
	res, err := a.DecodeInput("set_hash", common.Hex2Bytes("207365745f686173680c6b657908081468656c6c6f14776f726c64081468656c6c6f14776f726c64"))
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestAbi_Encode3(t *testing.T) {
	a, err := JSON(strings.NewReader(VecJson))
	if err != nil {
		t.Error(err)
	}
	_, err = a.EncodeCompact("set_hash", &CompactString{Val: "key"},
		&CompactVec{
			Val: []Compact{
				&CompactString{Val: "hello"},
				&CompactString{Val: "world"},
			},
			NextList: []TypeString{StringName},
		})
	assert.NotNil(t, err)
}

func TestAbi_Encode4(t *testing.T) {
	a, err := JSON(strings.NewReader(VecJson))
	if err != nil {
		t.Error(err)
	}

	res, err := a.EncodeCompact("set_hash2", &CompactString{Val: "key"}, &CompactVec{Val: []Compact{
		&CompactString{Val: "hello", Type: StringName},
		&CompactString{Val: "world", Type: StringName},
	}, NextList: []TypeString{StringName}})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "247365745f68617368320c6b6579081468656c6c6f14776f726c64", common.Bytes2Hex(res))
}

func TestAbi_Decode4(t *testing.T) {
	a, err := JSON(strings.NewReader(VecJson))
	if err != nil {
		t.Error(err)
	}
	res, err := a.DecodeInput("set_hash2", common.Hex2Bytes("247365745f68617368320c6b6579081468656c6c6f14776f726c64"))
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestEncodeStruct(t *testing.T) {
	a, err := JSON(strings.NewReader(structBin))
	if err != nil {
		t.Error(err)
	}
	t.Run("encode", func(t *testing.T) {
		res, err := a.EncodeCompact("make_student", &CompactStruct{Val: []Compact{
			&FixU32{
				Val: uint32(1),
			},
			&CompactString{Val: "test", Type: StringName},
			&CompactVec{Val: []Compact{
				&CompactVec{
					Val: []Compact{
						&CompactString{Val: "hello"},
						&CompactString{Val: "world"},
					},
					NextList: []TypeString{StringName},
				},
			}, NextList: []TypeString{VecName, StringName}},
		}}, &CompactStruct{Val: []Compact{
			&CompactVec{Val: []Compact{
				&CompactVec{
					Val: []Compact{
						&CompactStruct{Val: []Compact{}},
					},
					NextList: []TypeString{StructName},
				},
			}, NextList: []TypeString{VecName, StructName}},
		}})
		assert.Nil(t, err)
		assert.Equal(t, "306d616b655f73747564656e7401000000107465737404081468656c6c6f14776f726c640404", common.Bytes2Hex(res))
	})
	t.Run("h", func(t *testing.T) {
		c := &CompactString{}
		c.Decode([]byte{48, 109, 97, 107, 101, 95, 115, 116, 117, 100, 101, 110, 116})
		assert.Equal(t, "make_student", c.GetVal())
	})
}

func TestConstruct(t *testing.T) {
	a, err := JSON(strings.NewReader(DeployParam))
	if err != nil {
		t.Error(err)
	}
	t.Run("encodeConstruct", func(t *testing.T) {
		ans, err := a.encodeConstruct(2)
		assert.Nil(t, err)
		assert.Equal(t, []byte{0, 11, 6, 112, 97, 114, 97, 109, 115, 4, 78, 2, 0, 0, 0}, ans)
	})
	t.Run("encode", func(t *testing.T) {
		ans, err := a.Encode("", 2)
		assert.Nil(t, err)
		assert.Equal(t, []byte{0, 11, 6, 112, 97, 114, 97, 109, 115, 4, 78, 2, 0, 0, 0}, ans)
	})
	t.Run("encodeCompact", func(t *testing.T) {
		ans, err := a.encodeConstructCompact(&FixU32{Val: uint32(2)})
		assert.Nil(t, err)
		assert.Equal(t, []byte{0, 11, 6, 112, 97, 114, 97, 109, 115, 2, 0, 0, 0}, ans)
	})

	t.Run("encodeCompact", func(t *testing.T) {
		ans, err := a.EncodeCompact("", &FixU32{Val: uint32(2)})
		assert.Nil(t, err)
		assert.Equal(t, []byte{0, 11, 6, 112, 97, 114, 97, 109, 115, 2, 0, 0, 0}, ans)
	})
}

func TestTupleAbi(t *testing.T) {
	a, err := JSON(strings.NewReader(DeployParam))
	if err != nil {
		t.Error(err)
	}
	t.Run("encode", func(t *testing.T) {
		ans, err := a.Encode("set_tuple", []interface{}{
			1,
			2,
		})
		assert.Nil(t, err)
		t.Log(ans)
	})
	t.Run("decode", func(t *testing.T) {
		ans, err := a.DecodeRet([]byte{1, 0, 0, 0, 2, 0, 0, 0}, "get_tuple")
		assert.Nil(t, err)
		assert.Equal(t, TupleName, ans.Params[0].GetType())
		for _, v := range ans.Params {
			t.Log(GetCompactValue(v))
		}
	})
}

func TestEnumAbi(t *testing.T) {
	a, err := JSON(strings.NewReader(DeployParam))
	if err != nil {
		t.Error(err)
	}
	t.Run("encode", func(t *testing.T) {
		ans, err := a.Encode("set_enum", []interface{}{
			1,
			2,
		})
		assert.Nil(t, err)
		t.Log(ans)
	})
	t.Run("encode2", func(t *testing.T) {
		ans, err := a.Encode("set_enum", []interface{}{
			0,
			2, []interface{}{3, 4},
		})
		assert.Nil(t, err)
		t.Log(ans)
	})
	t.Run("decode", func(t *testing.T) {
		ans, err := a.DecodeRet([]byte{1, 1, 0, 0, 0}, "get_enum")
		assert.Nil(t, err)
		assert.Equal(t, EnumName, ans.Params[0].GetType())
		for _, v := range ans.Params {
			t.Log(GetCompactValue(v))
		}
	})
}

func TestOptionAbi(t *testing.T) {
	a, err := JSON(strings.NewReader(DeployParam))
	if err != nil {
		t.Error(err)
	}
	t.Run("encode", func(t *testing.T) {
		ans, err := a.Encode("set_option", []interface{}{
			0,
		})
		assert.Nil(t, err)
		t.Log(ans)
	})
	t.Run("decode", func(t *testing.T) {
		ans, err := a.DecodeRet([]byte{1, 1, 0, 0, 0}, "get_option")
		assert.Nil(t, err)
		assert.Equal(t, EnumName, ans.Params[0].GetType())
		for _, v := range ans.Params {
			t.Log(GetCompactValue(v))
		}
	})
	t.Run("decode2", func(t *testing.T) {
		ans, err := a.DecodeRet([]byte{0}, "get_option")
		assert.Nil(t, err)
		assert.Equal(t, EnumName, ans.Params[0].GetType())
		for _, v := range ans.Params {
			t.Log(GetCompactValue(v))
		}
	})
}
