package bvm

import (
	"github.com/gogo/protobuf/proto"
	"github.com/jackzing/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecode(t *testing.T) {
	logStr := "08011284030000000a0000000f53657446696c746572456e61626c65000000010000000566616c73650000000e53657446696c74657252756c6573000000010000001f5b7b22616c6c6f775f616e796f6e65223a66616c73652c226964223a307d5d00000015536574436f6e73656e737573426174636853697a65000000010000000331303000000014536574436f6e73656e737573506f6f6c53697a65000000010000000332303000000013536574436f6e73656e73757353657453697a65000000010000000235300000001453657450726f706f73616c5468726573686f6c640000000100000001340000001253657450726f706f73616c54696d656f7574000000010000000c3438303030303030303030300000001253657450726f706f73616c54696d656f7574000000010000000c34383030303030303030303000000018536574436f6e7472616374566f74655468726573686f6c6400000001000000013300000015536574436f6e7472616374566f7465456e61626c6500000001000000047472756518d0bbb4bbebcec4d11620d0abc786c9d7c4d116280132720a2a3078303030663161376130386363633438653564333066383038353063663163663238336161336162641242307831366133313237363930393565323832633130323861376138376634336563623864623937653765643166386234663466626438393435356363353266386666180140064801522a3078303030663161376130386363633438653564333066383038353063663163663238336161336162645a05302e312e30"
	a := common.Hex2Bytes(logStr)
	var pro ProposalData
	proto.Unmarshal(a, &pro)
	assert.Equal(t, uint64(1), pro.Id)
	code, err := DecodeProposalCode(pro.Code)
	assert.Nil(t, err)
	assert.Equal(t, "[{\"MethodName\":\"SetFilterEnable\",\"Params\":[\"false\"]},{\"MethodName\":\"SetFilterRules\",\"Params\":[\"[{\\\"allow_anyone\\\":false,\\\"id\\\":0}]\"]},{\"MethodName\":\"SetConsensusBatchSize\",\"Params\":[\"100\"]},{\"MethodName\":\"SetConsensusPoolSize\",\"Params\":[\"200\"]},{\"MethodName\":\"SetConsensusSetSize\",\"Params\":[\"50\"]},{\"MethodName\":\"SetProposalThreshold\",\"Params\":[\"4\"]},{\"MethodName\":\"SetProposalTimeout\",\"Params\":[\"480000000000\"]},{\"MethodName\":\"SetProposalTimeout\",\"Params\":[\"480000000000\"]},{\"MethodName\":\"SetContractVoteThreshold\",\"Params\":[\"3\"]},{\"MethodName\":\"SetContractVoteEnable\",\"Params\":[\"true\"]}]", code)
}

func TestPayload(t *testing.T) {
	payload := "0x0000000644697265637400000002000002d400000002000000074164644e6f64650000000400000284308202803082022ca00302010202083207cade0b3bc127300a06082a8648ce3d0403023043310b300906035504061302434e310e300c060355040a1305666c61746f31093007060355040b1300310e300c060355040313056e6f64653131093007060355042a13003020170d3230313233313030303030305a180f32313230313233313030303030305a30819a310b300906035504061302434e313d303b060355040a133465794a77624746305a6d397962534936496d5a735958527649697769646d567963326c7662694936496a41754d433478496e303d310e300c060355040b13056563657274310e300c060355040313056e6f646536311c301a0603550405131333363035303733303831373533333837333033310e300c060355042a130565636572743056301006072a8648ce3d020106052b8104000a03420004b349fac31631546b2d4ce7ee5ebf32fe0f7ebd3df6ebc966c471a73dd1b456e6409ea756432ac4091346a5ad609aadc5d432862e0860e03e3dcc3e444da4b840a381b23081af300e0603551d0f0101ff0404030201ee30310603551d25042a302806082b0601050507030206082b0601050507030106082b0601050507030306082b06010505070304300f0603551d130101ff040530030101ff301d0603551d0e04160414c4b54b65ee7c3038bcabfd41f9f8da94ec407841301f0603551d23041830168014c8778e1445f7273299b3152f2902025b08b64258300b0603551d11040430028200300c06032a560104056563657274300a06082a8648ce3d0403020342008f96822984ce72212f1c2ce5afff68279e3d85564f21810706c7ddd78be871c061f554ebb70225259b91a6feb031fb3cd352af182e45223f495dc15088a3877a01000000056e6f64653600000002767000000006676c6f62616c00000005416464565000000002000000056e6f64653600000006676c6f62616c000000044e4f4445"
	res, err := DecodePayload(common.Hex2Bytes(payload))
	assert.Nil(t, err)
	assert.Equal(t, ContractMethod("Direct"), res.Method())
}
