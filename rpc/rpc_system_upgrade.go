package rpc

import (
	"encoding/json"
	"errors"
	"github.com/hyperchain/go-hpc-common/types"
	"github.com/hyperchain/gosdk/bvm"
)

func (rpc *RPC) GetTxVersionByID(id int) (string, StdError) {
	if id < 1 {
		return "", NewSystemError(errors.New("id must equal to or more than 1"))
	}
	url := rpc.hrm.nodes[id-1].url
	method := TRANSACTION + "getTransactionsVersion"
	req := rpc.jsonRPC(method)
	data, err := rpc.callWithReqURL(req, url)
	if err != nil {
		return "", err
	}
	var txVersion string
	if sysErr := json.Unmarshal(data, &txVersion); sysErr != nil {
		return "", NewSystemError(sysErr)
	}
	return txVersion, nil
}

func (rpc *RPC) SetSupportedVersionByID(id int) (*TxReceipt, StdError) {
	if id < 1 {
		return nil, NewSystemError(errors.New("id must equal to or more than 1"))
	}
	url := rpc.hrm.nodes[id-1].url
	receipt, err := rpc.CallByPollingWithURL(VERSION+"setSupportedVersion", url, nil)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (rpc *RPC) GetSupportedVersionByHostname(hostname string) (types.SupportedVersion, StdError) {
	method := VERSION + "getSupportedVersionByHostname"
	data, err := rpc.call(method, hostname)
	if err != nil {
		return nil, err
	}

	var vr types.SupportedVersion
	if sysErr := json.Unmarshal(data, &vr); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}
	return vr, nil
}

type VersionResult struct {
	AvailableHyperchainVersion map[string]types.RunningVersion `json:"availableHyperchainVersions"`
	RunningHyperchainVersion   map[string]types.RunningVersion `json:"runningHyperchainVersions"`
}

func (rpc *RPC) GetVersions() (*VersionResult, StdError) {
	method := VERSION + "getVersions"
	data, err := rpc.call(method)
	if err != nil {
		return nil, err
	}

	var vr VersionResult
	if sysErr := json.Unmarshal(data, &vr); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}
	return &vr, nil
}

func (rpc *RPC) NewProposalCreateOperationForSystemUpgrade(targetVersion string) (bvm.BuiltinOperation, error) {
	rversions, stdErr := rpc.GetHyperchainVersionFromBin(targetVersion)
	if stdErr != nil {
		return nil, stdErr
	}

	suo, nerr := bvm.NewSystemUpgradeOperation(rversions)
	if nerr != nil {
		return nil, NewSystemError(nerr)
	}

	proposalOP := bvm.NewProposalCreateOperationForSystemUpgrade(suo)
	return proposalOP, nil
}

func (rpc *RPC) GetHyperchainVersionFromBin(hyperchainVersion string) (types.RunningVersion, StdError) {
	method := VERSION + "getHyperchainVersionFromBin"
	data, err := rpc.call(method, hyperchainVersion)
	if err != nil {
		return nil, err
	}

	var vr types.RunningVersion
	if sysErr := json.Unmarshal(data, &vr); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}
	return vr, nil
}

func (rpc *RPC) callWithReqURL(req *JSONRequest, randomURL string) (json.RawMessage, StdError) {
	body, sysErr := json.Marshal(req)
	if sysErr != nil {
		return nil, NewSystemError(sysErr)
	}

	data, err := rpc.hrm.SyncRequestSpecificURL(body, randomURL, GENERAL, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp *JSONResponse
	if sysErr = json.Unmarshal(data, &resp); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}
	if resp.Code != SuccessCode {
		return nil, NewServerError(resp.Code, resp.Message)
	}
	return resp.Result, nil
}

func (rpc *RPC) CallWithURL(method string, param interface{}, randomURL string) (*TxReceipt, StdError) {
	req := rpc.jsonRPC(method, param)
	result, err := rpc.callWithReqURL(req, randomURL)
	if err != nil {
		return nil, err
	}

	var receipt TxReceipt
	if sysErr := json.Unmarshal(result, &receipt); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}
	return &receipt, nil
}

// CallByPollingWithURL call and get tx receipt by polling
func (rpc *RPC) CallByPollingWithURL(method string, url string, params interface{}) (*TxReceipt, StdError) {
	var hash string

	data, cerr := rpc.callWithSpecificURL(method, url, params)
	if cerr != nil {
		return nil, cerr
	}

	if sysErr := json.Unmarshal(data, &hash); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}

	for i := int64(0); i < rpc.resTime; i++ {
		txReceipt, innErr, success := rpc.GetTxReceiptByPolling(hash, false)
		if success {
			return txReceipt, innErr
		}
		continue
	}
	return nil, NewRequestTimeoutError(errors.New("request time out"))
}
