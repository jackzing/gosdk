package bvm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	a := &ProposalOperationImpl{}
	assert.Equal(t, true, IsSystemOpt(a))
}
