package bvmcom

import (
	"reflect"
	"strconv"

	"github.com/hyperchain/go-hpc-common/types/protos"
)

// NewResult create a new Result and return
func NewResult(success bool, ret []byte, err error) *Result {

	result := &Result{
		Success: success,
		Ret:     ret,
	}
	if err != nil {
		result.Err = err.Error()
	}
	return result
}

// NewInvalidResult create a new invalid Result and return
func NewInvalidResult(err error) *Result {

	result := &Result{
		Success:   false,
		isInvalid: true,
	}
	if err != nil {
		result.Err = err.Error()
	}
	return result
}

// Result contract execute result
type Result struct {
	Success bool
	Ret     []byte
	Err     string

	// a flag for if need be defined as invalid tx used in bvm/exec.go
	isInvalid bool
	// a label to identify whether to execute contract management proposal, etc.
	Label string `json:"label,omitempty"`
}

// IsInvalid if need be defined as invalid tx
func (result *Result) IsInvalid() bool {
	return result.isInvalid
}

// SetLabel set label of result.
func (result *Result) SetLabel(label string) {
	result.Label = label
}

// GetExecContractProposalLabel return the label of execute contract proposal.
func GetExecContractProposalLabel() string {
	return GetExecProposalLabel(protos.ProposalData_CONTRACT)
}

// GetExecProposalLabel return the label of execute proposal with given pty.
func GetExecProposalLabel(pty protos.ProposalData_PType) string {
	return ProposalPrefix + split + pty.String()
}

// ParseArgs parse in to m.parameters
func ParseArgs(m reflect.Value, in []string) ([]reflect.Value, error) {
	argsNum := m.Type().NumIn()
	args := make([]reflect.Value, argsNum)
	for i := 0; i < argsNum; i++ {
		t := m.Type().In(i).Kind()
		inArg := string(in[i])
		switch t {
		case reflect.Bool:
			ret, err := strconv.ParseBool(inArg)
			if err != nil {
				return nil, err
			}
			args[i] = reflect.ValueOf(ret)
		case reflect.Int:
			ret, err := strconv.Atoi(inArg)
			if err != nil {
				return nil, err
			}
			args[i] = reflect.ValueOf(ret)
		case reflect.Uint64:
			ret, err := strconv.ParseUint(inArg, 10, 0)
			if err != nil {
				return nil, err
			}
			args[i] = reflect.ValueOf(ret)
		case reflect.Int64:
			ret, err := strconv.ParseInt(inArg, 10, 0)
			if err != nil {
				return nil, err
			}
			args[i] = reflect.ValueOf(ret)
		case reflect.String:
			args[i] = reflect.ValueOf(inArg)
		case reflect.Slice:
			args[i] = reflect.ValueOf([]byte(inArg))
		default:
			args[i] = reflect.ValueOf(inArg)
		}
	}
	return args, nil
}
