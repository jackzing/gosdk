package cryptocom

import (
	"github.com/hyperchain/go-hpc-common/types/protos"
)

//VerifyTube a verify service
type VerifyTube interface {
	//GetPromise get promise
	GetPromise([]*protos.Transaction)
	//SyncVerify synchronous interface
	SyncVerify(a uint8, sign, key, from []byte, msg string) error
	//ClearCertCache clear cert cache
	ClearCertCache() error
	//Query it should be called only once for each transaction
	//especially should not be called multiple times in multiple threads
	Query(*protos.Transaction) error
	Start() error
	Stop()
}
