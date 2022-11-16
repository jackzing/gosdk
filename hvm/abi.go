package hvm

import (
	"errors"
	"strings"

	"github.com/hyperchain/gosdk/common"
	"github.com/hyperchain/gosdk/common/parallel"
)

// Version abi version.
type Version string

const (
	// Version1 version v1.
	Version1 Version = "v1"
)

// Type abi type.
type Type string

// BeanType bean type.
type BeanType string

// nolint
const (
	Void   Type = "Void"
	Bool   Type = "Bool"
	Char   Type = "Char"
	Byte   Type = "Byte"
	Short  Type = "Short"
	Int    Type = "Int"
	Long   Type = "Long"
	Float  Type = "Float"
	Double Type = "Double"
	String Type = "String"
	Array  Type = "Array"
	List   Type = "List"
	Map    Type = "Map"
	Struct Type = "Struct"
)

const (
	// InvokeBean invoke bean type.
	InvokeBean BeanType = "InvokeBean"
	// MethodBean invoke  directly type.
	MethodBean BeanType = "MethodBean"
)

// Entry field/parameter entry.
type Entry struct {
	Name       string  `json:"name"`
	EntryType  Type    `json:"type"`
	Properties []Entry `json:"properties,omitempty"`
	StructName string  `json:"structName,omitempty"`
}

// BeanAbi invoke abi message.
type BeanAbi struct {
	BeanVersion Version  `json:"version"`
	BeanName    string   `json:"beanName"`
	Inputs      []Entry  `json:"inputs"`
	Output      Entry    `json:"output"`
	ClassBytes  string   `json:"classBytes"`
	Structs     []Entry  `json:"structs"`
	BeanType    BeanType `json:"beanType"`

	ParallelLevel parallel.ParallelLevel `json:"parallelLevel"`
	MutexFields   []MutexField           `json:"mutexFields"`
	MutexCalls    []MutexCall            `json:"mutexCalls"`
}

// MutexField mutex field message from abi.
type MutexField struct {
	FieldID   int   `json:"fieldID"`
	ParaIndex []int `json:"paraIndex"`
}

// MutexCall mutex cross call.
type MutexCall struct {
	AddressIndex int      `json:"addressIndex"`
	Address      string   `json:"Address"`
	Methods      []string `json:"methods"`
}

// MutexElement mutex element message of tx.
type MutexElement struct {
	// Address means contract address, if address is empty, represent current contract.
	Address string `json:"address"`
	// FieldID means mutex StoreField id, if is 0, mean contract mutex.
	FieldID   int      `json:"fieldID"`
	ParaIndex []uint64 `json:"paraIndex"`
}

// Abi hvm contract abi.
type Abi []BeanAbi

// GetBeanAbi get invoke bean abi.
func (abi Abi) GetBeanAbi(beanName string) (*BeanAbi, error) {
	for _, beanAbi := range abi {
		if beanAbi.BeanName == beanName {
			return &beanAbi, nil
		}
	}
	return nil, errors.New("can not find invoke bean " + beanName)
}

// GetMethodAbi get invoke directly bean abi.
func (abi Abi) GetMethodAbi(methodName string) (*BeanAbi, error) {
	index := strings.IndexByte(methodName, '(')
	if index == -1 {
		for _, beanAbi := range abi {
			if beanAbi.BeanName == methodName && beanAbi.BeanType == MethodBean {
				return &beanAbi, nil
			}
		}
	} else {
		if !strings.HasSuffix(methodName, ")") {
			return nil, errors.New("methodName is not legal")
		}
		name := methodName[0:index]
		paramStr := methodName[index+1 : len(methodName)-1]
		params := make([]string, 0)
		if len(paramStr) != 0 {
			params = strings.Split(paramStr, ",")
		}
		for _, beanAbi := range abi {
			if beanAbi.BeanName == name && beanAbi.BeanType == MethodBean && beanAbi.checkMethodInputs(params) {
				return &beanAbi, nil
			}
		}
	}
	return nil, errors.New("can not find method bean " + methodName)
}

func (beanAbi *BeanAbi) resolveStruct(input Entry, param interface{}) string {
	result := "{"
	s, err := beanAbi.getStruct(input.StructName)
	if err != nil {
		return ""
	}
	properties := s.Properties
	realParams := param.([]interface{})
	var (
		entryRes string
	)

	for i, prop := range properties {
		p := realParams[i]
		entryRes, _, _ = beanAbi.resolveEntry(prop, p)
		result += entryRes
		if i < len(properties)-1 {
			result += ","
		}
	}

	result += "}"
	return result
}

func (beanAbi *BeanAbi) resolveList(input Entry, param interface{}) string {
	result := "["
	entry := input.Properties[0]
	//entryType := entry.EntryType
	realParams := param.([]interface{})
	for i, p := range realParams {
		result += beanAbi.resolveNestListOrMap(entry, p)

		if i < len(realParams)-1 {
			result += ","
		}
	}

	result += "]"
	return result
}

func (beanAbi *BeanAbi) resolveMap(input Entry, param interface{}) string {
	result := "{"
	keyEntry := input.Properties[0]
	valEntry := input.Properties[1]

	realParams := param.([]interface{})
	for i, p := range realParams {

		kv := p.([]interface{})
		result += beanAbi.resolveEntrySingle(keyEntry, kv[0])
		result += ":"
		result += beanAbi.resolveNestListOrMap(valEntry, kv[1])

		if i < len(realParams)-1 {
			result += ","
		}
	}

	result += "}"
	return result
}

func (beanAbi *BeanAbi) resolveNestListOrMap(valEntry Entry, param interface{}) (result string) {
	switch valEntry.EntryType {
	case Bool, Byte, Short, Int, Long, Float, Double:
		result += beanAbi.resolveEntrySingle(valEntry, param)
	case Char, String:
		result += beanAbi.resolveEntrySingle(valEntry, param)
	case Struct:
		result += beanAbi.resolveStruct(valEntry, param)
	case List:
		result += beanAbi.resolveList(valEntry.Properties[0], param)
	case Map:
		result += beanAbi.resolveMap(valEntry.Properties[0], param)
	case Array:
		result += beanAbi.resolveArray(valEntry.Properties[0], param)

	}
	return
}

func (beanAbi *BeanAbi) resolveArray(input Entry, param interface{}) string {
	result := "["
	entry := input.Properties[0]
	//entryType := entry.EntryType
	realParams := param.([]interface{})
	for i, p := range realParams {
		result += beanAbi.resolveEntrySingle(entry, p)

		if i < len(realParams)-1 {
			result += ","
		}
	}

	result += "]"
	return result
}

func (beanAbi *BeanAbi) resolveEntrySingle(entry Entry, param interface{}) (result string) {
	switch entry.EntryType {
	case Bool, Byte, Short, Int, Long, Float, Double:
		if str, ok := param.(string); ok {
			result = str
		}
	case Char, String:
		if str, ok := param.(string); ok {
			result = "\"" + str + "\""
		}
		// TODO list<map<k, v>>, map<k, list<T>>
	case Struct:
		result = beanAbi.resolveStruct(entry, param)
	}
	return
}

func (beanAbi *BeanAbi) resolveEntry(input Entry, param interface{}) (entryResult string, start, end int) {
	start = len("\"") + len(input.Name) + len("\":")
	switch input.EntryType {
	case Bool, Byte, Short, Int, Long, Float, Double:
		if str, ok := param.(string); ok {
			entryResult = "\"" + input.Name + "\":" + str
		}
	case Char, String:
		if str, ok := param.(string); ok {
			entryResult = "\"" + input.Name + "\":\"" + str + "\""
			start++
			end = -1
		}
	case List:
		entryResult = "\"" + input.Name + "\":" + beanAbi.resolveList(input, param)
	case Map:
		entryResult = "\"" + input.Name + "\":" + beanAbi.resolveMap(input, param)
	case Array:
		entryResult = "\"" + input.Name + "\":" + beanAbi.resolveArray(input, param)
	case Struct:
		entryResult = "\"" + input.Name + "\":" + beanAbi.resolveStruct(input, param)
	}
	end += len(entryResult)
	return
}

func (beanAbi *BeanAbi) resolveJSONEntry(input Entry, param string) (entryResult string, start, end int) {
	start = len("\"") + len(input.Name) + len("\":")
	switch input.EntryType {
	case Char, String:
		entryResult = "\"" + input.Name + "\":\"" + strings.ReplaceAll(param, `"`, `\"`) + "\""
		start++
		end = -1
	default:
		entryResult = "\"" + input.Name + "\":" + param
	}
	end += len(entryResult)
	return
}

func (beanAbi BeanAbi) encode(startPos int, params ...interface{}) (result string, paraIndex []int, _ error) {
	if len(beanAbi.Inputs) != len(params) {
		return "", nil, errors.New("param count is not enough")
	}
	var (
		entryResult string
		start       int
		end         int
	)
	paraIndex = make([]int, len(params)*2)
	result = "{"
	for i, input := range beanAbi.Inputs {
		param := Convert(params[i])
		entryResult, start, end = beanAbi.resolveEntry(input, param)
		paraIndex[2*i] = startPos + len(result) + start
		paraIndex[2*i+1] = startPos + len(result) + end
		result += entryResult
		if i < len(params)-1 {
			result += ","
		}
	}
	result += "}"

	return result, paraIndex, nil
}

func (beanAbi BeanAbi) encodeJSON(startPos int, params ...interface{}) (result string, paraIndex []int, _ error) {
	if len(beanAbi.Inputs) != len(params) {
		return "", nil, errors.New("param count is not equal with inputs count")
	}
	var (
		str         string
		entryResult string
		ok          bool
		start       int
		end         int
	)

	paraIndex = make([]int, len(params)*2)
	result = "{"
	for i, input := range beanAbi.Inputs {
		param := params[i]

		if str, ok = param.(string); !ok {
			return "", nil, errors.New("param should be string")
		}

		entryResult, start, end = beanAbi.resolveJSONEntry(input, str)
		paraIndex[2*i] = startPos + len(result) + start
		paraIndex[2*i+1] = startPos + len(result) + end
		result += entryResult
		if i < len(params)-1 {
			result += ","
		}
	}
	result += "}"
	return result, paraIndex, nil
}

func (beanAbi BeanAbi) classBytes() []byte {
	return common.Hex2Bytes(beanAbi.ClassBytes)
}

func (beanAbi BeanAbi) getStruct(name string) (*Entry, error) {
	for _, s := range beanAbi.Structs {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, errors.New("can not find struct " + name)
}

func (beanAbi BeanAbi) checkMethodInputs(params []string) bool {
	if len(beanAbi.Inputs) != len(params) {
		return false
	}
	for i, input := range beanAbi.Inputs {
		if input.StructName != params[i] {
			return false
		}
	}
	return true
}
