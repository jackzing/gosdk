package code

// JSON-RPC specified error code
var (
	ErrInvalidMessage = New(-32700, errInvalidMessageMsg)
	ErrCallback       = New(-32000, errCallbackMsg)
)

// JSON-RPC specified error code
var (
	ErrInvalidRequest  = New(-32600, errInvalidRequestMsg)
	ErrMethodNotFound  = New(-32601, errMethodNotFoundMsg)
	ErrInvalidParams   = New(-32602, errInvalidParamsMsg)
	ErrInterceptBefore = New(-32603, errInterceptBeforeMsg)
	ErrInterceptAfter  = New(-32604, errInterceptAfterMsg)
)

// JSON-RPC specified error message
const (
	errInvalidMessageMsg  = "Invalid JSON message"
	errCallbackMsg        = "Internal server error"
	errInvalidRequestMsg  = "Invalid request"
	errMethodNotFoundMsg  = "Method not found"
	errInvalidParamsMsg   = "Invalid params"
	errInterceptBeforeMsg = "Intercepted in before"
	errInterceptAfterMsg  = "Intercepted in after"
)
