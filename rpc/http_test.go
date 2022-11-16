package rpc

import (
	"fmt"
	"github.com/jackzing/gosdk/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRPC_Reconnect(t *testing.T) {
	rpc := NewRPC()
	rpc.hrm.ReConnectNode(0)

}

func TestRPC_GetNode(t *testing.T) {
	rpc := NewRPC()
	rpc.hrm.nodes[0].status = false
	fmt.Print(rpc.GetNodes())

}

func TestSort(t *testing.T) {
	cf, err := config.NewFromFile("../conf")
	assert.Nil(t, err)
	hrm := newHTTPRequestManager(cf, "../conf")
	_, err1 := hrm.selectNodeURL()
	if err1 != nil {
		t.Fatal(err1)
	}
}
