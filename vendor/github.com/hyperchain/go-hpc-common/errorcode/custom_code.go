package code

// JSON-RPC custom error code
// Range -32001 to -32099
var (
	ErrDBNotFound           = New(-32001, errDBNotFoundMsg)
	ErrOutOfBalance         = New(-32002, errOutOfBalanceMsg)
	ErrSignatureInvalid     = New(-32003, errSignatureInvalidMsg)
	ErrContractDeploy       = New(-32004, errContractDeployMsg)
	ErrContractInvoke       = New(-32005, errContractInvokeMsg)
	ErrSystemTooBusy        = New(-32006, errSystemTooBusyMsg)
	ErrRepeatedTx           = New(-32007, errRepeatedTxMsg)
	ErrContractPermission   = New(-32008, errContractPermissionMsg)
	ErrAccountNotExist      = New(-32009, errAccountNotExistMsg)
	ErrNamespaceNotFound    = New(-32010, errNamespaceNotFoundMsg)
	ErrNoBlockGenerated     = New(-32011, errNoBlockGeneratedMsg)
	ErrSubNotExist          = New(-32012, errSubNotExistMsg)
	ErrSnapshot             = New(-32013, errSnapshotMsg)
	ErrUnlockAccount        = New(-32014, errUnlockAccountMsg)
	ErrInvokeNativeContract = New(-32015, errInvokeNativeContractMsg)
	ErrInvalidNodeHash      = New(-32016, errInvalidNodeHashMsg)
	ErrUnReachablePeer      = New(-32017, errUnReachablePeerMsg)
	ErrInconsistent         = New(-32018, errInconsistentMsg)
	ErrInvalidCollection    = New(-32019, errInvalidCollectionMsg)
	ErrSelfGovService       = New(-32020, errSelfGovServiceMsg)
	ErrDeprecatedAPI        = New(-32021, errDeprecatedAPIMsg)
	ErrDIDInvocation        = New(-32022, errDIDInvocationMsg)
	ErrCvpNotSupport        = New(-32023, errCvpNotSupport)
	ErrConsensusAbnormal    = New(-32024, errConsensusAbnormal)
	ErrDispatcherFull       = New(-32025, errDispatcherFull)
	ErrQPSLimit             = New(-32026, errQPSLimit)
	ErrSimulateLimit        = New(-32027, errSimulateLimit)
	ErrKVSQLEXEC            = New(-32028, errKVSQLEXECMsg)
	ErrGetProof             = New(-32029, errGetProofMsg)
	ErrValidateProof        = New(-32030, errValidateProofMsg)
	ErrLpNotSupport         = New(-32031, errLpNotSupport)

	ErrInvalidToken = New(-32097, errInvalidTokenMsg)
	ErrUnauthorized = New(-32098, errUnauthorizedMsg)
	ErrCert         = New(-32099, errCertErrorMsg)

	ErrCVPBefore = New(-32605, errCVPBefore)
	ErrCVPAfter  = New(-32606, errCVPAfter)
)

// JSON-RPC custom error message
const (
	errDBNotFoundMsg           = "DB not found"
	errOutOfBalanceMsg         = "Out of balance"
	errSignatureInvalidMsg     = "Invalid signature"
	errContractDeployMsg       = "Deploy contract failed"
	errContractInvokeMsg       = "Invoke contract failed"
	errSystemTooBusyMsg        = "System is too busy"
	errRepeatedTxMsg           = "Repeated transaction"
	errContractPermissionMsg   = "Contract invocation permission not enough"
	errAccountNotExistMsg      = "Account dose not exist or account balance is 0"
	errNamespaceNotFoundMsg    = "The namespace does not exist"
	errNoBlockGeneratedMsg     = "There is no block generated"
	errSubNotExistMsg          = "Required subscription does not existed or has expired"
	errSnapshotMsg             = "The process of snapshot or archive happened error"
	errUnlockAccountMsg        = "Failed to unlock local node account"
	errInvokeNativeContractMsg = "Failed to invoke native contract"
	errInvalidNodeHashMsg      = "Invalid participant node hashes"
	errUnReachablePeerMsg      = "Unreachable peers"
	errInconsistentMsg         = "Inconsistent peer"
	errInvalidCollectionMsg    = "Invalid private transaction with no collection info"
	errSelfGovServiceMsg       = "ACO service is not available"
	errDeprecatedAPIMsg        = "Deprecated API"
	errDIDInvocationMsg        = "DID invocation failed"
	errCvpNotSupport           = "cvp not support api"
	errLpNotSupport            = "lp not support api"
	errConsensusAbnormal       = "Consensus status abnormal"
	errDispatcherFull          = "Dispatcher full"
	errQPSLimit                = "QPS limit"
	errSimulateLimit           = "Simulate limit"
	errKVSQLEXECMsg            = "sql exec failed"
	errGetProofMsg             = "get proof failed"
	errValidateProofMsg        = "validate proof failed"

	errInvalidTokenMsg = "Invalid token"
	errUnauthorizedMsg = "Unauthorized, Please check your cert"
	errCertErrorMsg    = "Cert error"

	errCVPBefore = "cvp intercept before error"
	errCVPAfter  = "cvp intercept after error"
)
