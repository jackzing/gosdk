package rpc

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/hyperchain/go-crypto-standard/asym"
	"github.com/jackzing/gosdk/account"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	txGrpc, err := g.NewTransactionGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)

	queueName := "queue_name2"
	var stop uint32
	var start = make(chan bool)
	t.Run("mq", func(t *testing.T) {
		t.Parallel()
		que, err := g.NewGrpcMQ()
		if err != nil {
			t.Error(err)
		}
		//注册队列
		_, err = que.Register(&RegisterMeta{
			RoutingKeys: []routingKey{MQBlock, MQLog},
			QueueName:   queueName,
			From:        "",
			Signature:   "",
			IsVerbose:   true,
			FromBlock:   "",
			ToBlock:     "",
			Addresses:   nil,
			Topics:      nil,
			Delay:       false,
		})
		assert.Nil(t, err)

		<-start
		stream, err := que.Consume(&ConsumeParams{
			QueueName: queueName,
		})

		assert.Nil(t, err)

		go func() {
			for atomic.LoadUint32(&stop) == 0 {
				res, err := stream.Recv()
				assert.Nil(t, err)
				t.Logf("[MQ]: %v", res.String())
			}

			_, err = que.UnRegister(&UnRegisterMeta{
				QueueName: queueName,
				From:      "",
				Signature: "",
			})
			t.Log("unregister")
			assert.Nil(t, err)
		}()
	})

	t.Run("tx", func(t *testing.T) {
		t.Parallel()
		address, _ := testPrivateAccount()
		for i := 0; i < 13; i++ {
			if i == 7 {
				start <- true
			}
			guomiKey, _ := asym.GenerateKey(asym.AlgoP256R1)
			pubKey := &account.ECDSAKey{ECDSAPrivateKey: guomiKey}
			newAddress := pubKey.GetAddress()
			transaction := NewTransaction(newAddress.Hex()).Transfer(address, int64(0))
			transaction.Sign(pubKey)

			receipt, err := txGrpc.SendTransactionReturnReceipt(transaction)
			assert.Nil(t, err)
			t.Log(receipt.TxHash)
			t.Log(transaction.GetTransactionHash(1000000))
		}
		atomic.StoreUint32(&stop, 1)
	})
}

func TestGrpcMQ_Register(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	que, err := g.NewGrpcMQ()
	if err != nil {
		t.Error(err)
	}
	ans, err := que.Register(&RegisterMeta{
		RoutingKeys: []routingKey{"MQBlock", "MQLog"},
		QueueName:   "test2",
		From:        "",
		Signature:   "",
		IsVerbose:   true,
		FromBlock:   "",
		ToBlock:     "",
		Addresses:   nil,
		Topics:      nil,
		Delay:       false,
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "test2", ans.QueueName)
}

func TestGrpcMQ_GetAllQueueNames(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	que, err := g.NewGrpcMQ()
	if err != nil {
		t.Error(err)
	}
	ans, err := que.GetAllQueueNames()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "test2", ans[0])
}

func TestGrpcMQ_UnRegister(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	que, err := g.NewGrpcMQ()
	if err != nil {
		t.Error(err)
	}
	ans, err := que.UnRegister(&UnRegisterMeta{
		QueueName: "test2",
		From:      "",
		Signature: "",
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, ans.Success)
}

func TestGrpcMQ_Consume(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	que, err := g.NewGrpcMQ()
	if err != nil {
		t.Error(err)
	}
	stream, err := que.Consume(&ConsumeParams{
		QueueName: "test2",
	})
	if err != nil {
		t.Error(err)
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Println(res)
	}
}

func TestGrpcMQ_StopConsume(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	que, err := g.NewGrpcMQ()
	if err != nil {
		t.Error(err)
	}
	ans, err := que.StopConsume(&StopConsumeParams{
		QueueName: "test2",
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, ans)
}
