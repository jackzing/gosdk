package rpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersionGrpc_SetSupportedVersion(t *testing.T) {
	g := NewGRPC()
	vg, err := g.NewVersionGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer vg.Close()

	ans, err := vg.SetSupportedVersion()
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}

func TestVersionGrpc_SetSupportedVersionReturnReceipt(t *testing.T) {
	g := NewGRPC()
	tg, err := g.NewVersionGrpc(ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()

	ans, err := tg.SetSupportedVersionReturnReceipt()
	if err != nil {
		t.Error(err)
	}
	t.Log(ans)
}
