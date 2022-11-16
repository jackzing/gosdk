package rpc

import (
	"fmt"
	"github.com/jackzing/gosdk/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestGetECert(t *testing.T) {
	confRootPath := "../conf"
	vip := viper.New()
	vip.SetConfigFile(filepath.Join(confRootPath, common.DefaultConfRelPath))
	err := vip.ReadInConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("read conf from %s error", filepath.Join(confRootPath, common.DefaultConfRelPath)))
	}
	tcm := NewTCertManager(vip, confRootPath)
	if tcm != nil {
		assert.Equal(t, true, len(tcm.GetECert()) > 0)
	}
}
