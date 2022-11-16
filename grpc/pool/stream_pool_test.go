package pool

import (
	"github.com/hyperchain/go-hpc-common/types/protos"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	var opt []grpc.DialOption
	opt = append(opt, grpc.WithInsecure())
	conn, err := grpc.Dial(":11001", opt...)
	if err != nil {
		log.Fatalf("init client pool err: %v", err)
	}
	defer conn.Close()

	t3 := protos.NewGrpcApiTransactionClient(conn)
	t.Log(t3)
}
