package code

import (
	"fmt"
)

// Common error message format
const (
	FmtNotFoundMethod      = "the method %s_%s does not exist/is not available"
	FmtNotFoundBlockNumber = "not found block number %#x"
	FmtNotFoundBlockHash   = "not found block hash %s"
	FmtNotFoundTxHash      = "not found transaction %s"
	FmtNotFoundNamespace   = "not found namespace %s"
	FmtNotFoundAccount     = "not found account %s"
	FmtNotFoundNode        = "not found node %v"
	FmtNotFoundReceipt     = "not found receipt %v"
	FmtNotFoundDiscardTx   = "not found discard transaction %v"
	FmtNewDBsFailed        = "new accountDB and stateDB failed"
)

var _codes = make(map[int]string)

// New checks whether the code is unique or not.
func New(c int, message string) int {
	if _, ok := _codes[c]; ok {
		panic(fmt.Sprintf("error code: %d already exist", c))
	}
	_codes[c] = message
	return c
}

// RPCError implements RPC error, is add support for error codec over regular go errors
type RPCError interface {
	// RPC error code
	Code() int
	// Error message
	Error() string
}

// CustomError defines error type, including error message and error code.
type CustomError struct {
	ErrCode int
	Message string
}

// Error returns error message.
func (ce *CustomError) Error() string {
	return ce.Message
}

// Code returns error code.
func (ce *CustomError) Code() int {
	return ce.ErrCode
}

// NewError creates and returns a new instance of CustomError.
func NewError(code int, format string, v ...interface{}) RPCError {
	desc := _codes[code]
	return NewCustomError(code, desc, format, v...)
}

// NewCustomError is an util function used by package errorcode and bvmcom.
// If you want to construct a RPCError for your application, please use NewError()
// not NewCustomError().
func NewCustomError(code int, desc string, format string, v ...interface{}) RPCError {
	err := &CustomError{code, desc}
	if format != "" {
		if err.Message != "" {
			err.Message += ": " + fmt.Sprintf(format, v...)
		} else {
			err.Message = fmt.Sprintf(format, v...)
		}
	}
	return err
}
