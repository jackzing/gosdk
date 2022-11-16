package types

import (
	"encoding/json"
	"fmt"
	"github.com/hyperchain/go-hpc-common/utils"
)

// InnerOracleRequest oracle inner request
type InnerOracleRequest struct {
	//web request
	URL    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`

	//specified
	BizID string `json:"bizID"`

	//callback
	CallBackAddress string `json:"callBackAddress"`
	CallBackMethod  string `json:"callBackMethod"`
}

// ToOracleRequest convert InnerOracleRequest to OracleRequest
func (req *InnerOracleRequest) ToOracleRequest() OracleRequest {
	return OracleRequest{
		URL:             req.URL,
		Header:          req.Header,
		Method:          req.Method,
		Body:            req.Body,
		BizID:           []byte(req.BizID),
		CallBackAddress: utils.DecodeString(req.CallBackAddress),
		CallBackMethod:  req.CallBackMethod,
	}
}

// OracleRequest oracle request
type OracleRequest struct {
	//web request
	URL    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`

	//specified
	TxHash []byte `json:"txHash"`
	BizID  []byte `json:"bizId"`

	//callback
	CallBackAddress []byte `json:"callBackAddress"`
	CallBackMethod  string `json:"callBackMethod"`
}

// OracleEvents oracle events
type OracleEvents []*OracleEvent

// OracleEvent oracle event
type OracleEvent struct {
	OracleRequest
	UUID            []byte `json:"uuid"`
	ContractAddress []byte `json:"contractAddress"`
}

// NewOracleEvent create OracleEvent
func NewOracleEvent(req OracleRequest, uuid, contractAddr []byte) *OracleEvent {
	return &OracleEvent{req, uuid, contractAddr}
}

// EncodeOracles encode Logs to bytes
func (os OracleEvents) EncodeOracles() ([]byte, error) {
	return json.Marshal(os)
}

// DecodeOracles decode Oracles from bytes
func DecodeOracles(buf []byte) (OracleEvents, error) {
	var tmp OracleEvents
	err := json.Unmarshal(buf, &tmp)
	return tmp, err
}

// ToOracleTrans oracle event to oracle trans
func (os OracleEvents) ToOracleTrans() []OracleTrans {
	var ret = make([]OracleTrans, len(os))
	for idx, o := range os {
		ret[idx] = OracleTrans{
			URL:             o.URL,
			Method:          o.Method,
			Header:          o.Header,
			Body:            o.Body,
			TxHash:          utils.BytesToHex(o.TxHash),
			BizID:           utils.BytesToHex(o.BizID),
			CallBackAddress: utils.BytesToHex(o.CallBackAddress),
			CallBackMethod:  o.CallBackMethod,
			UUID:            utils.BytesToHex(o.UUID),
			ContractAddress: utils.BytesToHex(o.ContractAddress),
		}
	}
	return ret
}

// String toString
func (oracle OracleEvent) String() string {
	return fmt.Sprintf("oracle request:%v, uuid:%v, contract addr:%v", oracle.OracleRequest, oracle.UUID, oracle.ContractAddress)
}

// OracleTrans oracle trans
type OracleTrans struct {
	//web request
	URL    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`

	//specified
	TxHash string `json:"txHash"`
	BizID  string `json:"bizId"`

	//callback
	CallBackAddress string `json:"callBackAddress"`
	CallBackMethod  string `json:"callBackMethod"`

	UUID            string `json:"uuid"`
	ContractAddress string `json:"contractAddress"`
}

// OracleResponse oracle response
type OracleResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Header  map[string]string `json:"header"`
	Body    string            `json:"body"`

	UUID  []byte `json:"uuid"`
	BizID []byte `json:"bizId"`

	CallerContract []byte `json:"callerContract"`
	CallbackMethod string `json:"callbackMethod"`
}
