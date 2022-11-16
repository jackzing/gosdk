package java

import (
	"encoding/base64"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperchain/go-hpc-common/types/protos"
	"github.com/hyperchain/gosdk/common"
	"github.com/mholt/archiver/v3"
	"github.com/opentracing/opentracing-go/log"
	"io/ioutil"
	"os"
	"strings"
)

var logger = common.GetLogger("java")

// ReadJavaContract read compiled java contract from the given
// path and return the payload used to deploy
// params indicates the constructor params
func ReadJavaContract(path string, params ...string) (string, error) {
	err := archiver.Archive([]string{path}, "tmp.tar.gz")
	if err != nil {
		logger.Error(err)
		return "", nil
	}

	tar, err := ioutil.ReadFile("tmp.tar.gz")
	if err != nil {
		logger.Error(err)
		return "", nil
	}

	err = os.Remove("tmp.tar.gz")
	if err != nil {
		log.Error(err)
		return "", err
	}

	invokeArgs := protos.InvokeArgs{
		Code: tar,
	}

	for _, p := range params {
		invokeArgs.Args = append(invokeArgs.Args, []byte(p))
	}

	res, err := proto.Marshal(&invokeArgs)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return common.Bytes2Hex(res), nil
}

// EncodeJavaFunc encodes method and params to invoke contract
func EncodeJavaFunc(methodName string, params ...string) []byte {
	invokeArgs := protos.InvokeArgs{
		MethodName: "invoke",
	}

	invokeArgs.Args = append(invokeArgs.Args, []byte(methodName))

	for _, p := range params {
		invokeArgs.Args = append(invokeArgs.Args, []byte(p))
	}

	res, err := proto.Marshal(&invokeArgs)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return res
}

// DecodeJavaResult decodes the return value of a java contract to string
func DecodeJavaResult(ret string) string {
	ret = strings.TrimPrefix(ret, "0x")
	return string(common.Hex2Bytes(ret))
}

// DecodeJavaLog decode the log value of a contract to string
func DecodeJavaLog(data string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(string(common.Hex2Bytes(data)))
	if err != nil {
		logger.Errorf("decode log failed: %v", err)
		return "", err
	}

	return string(res), nil
}
