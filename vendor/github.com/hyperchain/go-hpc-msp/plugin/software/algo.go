package software

import (
	"fmt"

	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/meshplus/crypto"
)

//ModeGetRSAMod get RSA mod from mode
func ModeGetRSAMod(mode int) (int, error) {
	if !common.ModeIsRSAAlgo(mode) {
		return 0, fmt.Errorf("it is'n RSA algo mode")
	}
	return (mode&0x0f00 + 0x200) << 2, nil
}

//ModeFromRSAMod get mode from RSA mode
func ModeFromRSAMod(rsaMod int) (int, error) {
	r := (rsaMod + 1023) >> 10
	if r < 3 {
		return crypto.Rsa2048, nil
	}
	mode := crypto.Rsa2048 | ((r - 2) << 8)
	if !common.ModeIsRSAAlgo(mode) {
		return crypto.None, fmt.Errorf("it is'n RSA algo mode")
	}
	return mode, nil
}
