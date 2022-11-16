package common

import (
	"github.com/meshplus/crypto"
)

//ModeIsSignatureAlgo is it a signature algorithm
func ModeIsSignatureAlgo(mode int) bool {
	return mode != 0 && mode&0xffff00ff == 0
}

//ModeIsHashAlgo is it a Hash algorithm
func ModeIsHashAlgo(mode int) bool {
	return mode != 0 && mode&0xffffff00 == 0
}

//ModeIsEncryptAlgo is it a encrypt algorithm
func ModeIsEncryptAlgo(mode int) bool {
	return mode != 0 && mode&0xff00ffff == 0
}

//ModeIsRSAAlgo is it a RSA signature algorithm
func ModeIsRSAAlgo(mode int) bool {
	return ModeIsSignatureAlgo(mode) && mode&0xf000 == crypto.Rsa2048
}

//ModeIsECDSAAlgo is it a ECDSA signature algorithm
func ModeIsECDSAAlgo(mode int) bool {
	return mode != 0 && mode&0xfffff0ff == 0 && mode != crypto.Sm2p256v1
}
