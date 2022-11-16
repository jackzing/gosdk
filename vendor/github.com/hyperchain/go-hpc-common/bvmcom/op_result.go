package bvmcom

// OpResult is the result of operation
type OpResult struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

// NewSuccOpResult generates success result of operation
func NewSuccOpResult() *OpResult {
	return &OpResult{
		Code: SuccessCode,
		Msg:  "",
	}
}

// NewSuccOpResultWithMsg generates success result of operation with message
func NewSuccOpResultWithMsg(msg string) *OpResult {
	return &OpResult{
		Code: SuccessCode,
		Msg:  msg,
	}
}

// NewErrOpResult generates error result of operation
func NewErrOpResult(code int32, err error) *OpResult {
	return &OpResult{
		Code: code,
		Msg:  err.Error(),
	}
}
