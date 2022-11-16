package fvm_test

import (
	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/gosdk/account"
	"github.com/hyperchain/gosdk/common"
	"github.com/hyperchain/gosdk/fvm"
	"github.com/hyperchain/gosdk/rpc"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestDeploy(t *testing.T) {
	rp := rpc.NewRPCWithPath("../../conf")
	wasmPath1 := "./wasm/SetHash-gc.wasm"
	buf, err := ioutil.ReadFile(wasmPath1)
	if err != nil {
		t.Fatal(err)
	}
	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}

	transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Deploy(common.Bytes2Hex(buf)).VMType(rpc.FVM)
	transaction.Sign(guomiKey)
	rep, err := rp.SignAndDeployContract(transaction, guomiKey)
	if err != nil {
		t.Error(err)
	}

	invokeInput := []byte{32, 115, 101, 116, 95, 104, 97, 115, 104, 24, 107, 101, 121, 48, 48, 49, 100, 116, 104, 105, 115, 32, 105, 115, 32, 116, 104, 101, 32, 118, 97, 108, 117, 101, 32, 111, 102, 32, 48, 48, 48, 49}
	invokeTrans := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(rep.ContractAddress, invokeInput).VMType(rpc.FVM)
	invokeTrans.Sign(guomiKey)
	_, err = rp.SignAndInvokeContract(invokeTrans, guomiKey)
	if err != nil {
		t.Error(err)
	}

	invokeInput2 := []byte{32, 103, 101, 116, 95, 104, 97, 115, 104, 24, 107, 101, 121, 48, 48, 49}
	invokeTrans2 := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(rep.ContractAddress, invokeInput2).VMType(rpc.FVM)
	invokeTrans2.Sign(guomiKey)
	recipt2, err := rp.SignAndInvokeContract(invokeTrans2, guomiKey)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "dthis is the value of 0001", string(common.Hex2Bytes(recipt2.Ret)))
}

func TestDeployWithParams(t *testing.T) {
	t.Skip("wasm bytes have changed")
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
	rp := rpc.NewRPCWithPath("../../conf")
	wasmPath1 := "./wasm/abi_test.wasm"
	buf, err := ioutil.ReadFile(wasmPath1)
	if err != nil {
		t.Fatal(err)
	}
	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}

	a, err := fvm.GenAbi(strings.NewReader(DeployParam))
	if err != nil {
		t.Error(err)
	}

	deployParams, err := fvm.Encode(a, "", 2)
	if err != nil {
		t.Error(err)
	}
	//buf = append(buf, deployParams...)

	transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).DeployWithArgs(buf, deployParams).VMType(rpc.FVM)
	transaction.Sign(guomiKey)
	rep, err := rp.SignAndDeployContract(transaction, guomiKey)
	if err != nil {
		t.Error(err)
	}

	invokeInput, err := fvm.Encode(a, "get_hello")
	if err != nil {
		t.Error(err)
	}

	invokeTrans := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(rep.ContractAddress, invokeInput).VMType(rpc.FVM)
	invokeTrans.Sign(guomiKey)
	returns, err := rp.SignAndInvokeContract(invokeTrans, guomiKey)
	if err != nil {
		t.Error(err)
	}
	res, err := fvm.DecodeRet(a, "get_hello", common.Hex2Bytes(returns.Ret))
	assert.Nil(t, err)
	assert.Equal(t, uint32(2), res.Params[0].GetVal())
}

func TestDemo(t *testing.T) {
	const currentABI = `{
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
	rp := rpc.NewRPCWithPath("../../conf")
	wasmPath1 := "./wasm/SetHash-gc.wasm"
	buf, err := ioutil.ReadFile(wasmPath1)
	if err != nil {
		t.Fatal(err)
	}
	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		&gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}

	transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Deploy(common.Bytes2Hex(buf)).VMType(rpc.FVM)
	transaction.Sign(guomiKey)
	rep, err := rp.SignAndDeployContract(transaction, guomiKey)
	if err != nil {
		t.Error(err)
	}
	a, err := fvm.GenAbi(strings.NewReader(currentABI))
	if err != nil {
		t.Error(err)
	}
	invokeInput, err := fvm.Encode(a, "set_hash", "key", "value")
	if err != nil {
		t.Error(err)
	}
	invokeTrans := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(rep.ContractAddress, invokeInput).VMType(rpc.FVM)
	invokeTrans.Sign(guomiKey)
	_, err = rp.SignAndInvokeContract(invokeTrans, guomiKey)
	if err != nil {
		t.Error(err)
	}

	invokeInput2, err := fvm.Encode(a, "get_hash", "key")
	if err != nil {
		t.Error(err)
	}
	invokeTrans2 := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(rep.ContractAddress, invokeInput2).VMType(rpc.FVM)
	invokeTrans2.Sign(guomiKey)
	recipt2, err := rp.SignAndInvokeContract(invokeTrans2, guomiKey)
	if err != nil {
		t.Error(err)
	}

	res, err := fvm.DecodeRet(a, "get_hash", common.Hex2Bytes(recipt2.Ret))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "value", res.Params[0].GetVal())
}
