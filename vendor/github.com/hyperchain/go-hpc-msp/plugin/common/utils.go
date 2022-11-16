package common

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/go-crypto-standard/asym/secp256k1"
	"github.com/meshplus/crypto"
)

//GetModeName get mod name
func GetModeName(mode int) (m string) {
	switch {
	//Hash
	case mode&0xffffff00 == 0:
		switch mode {
		case crypto.SHA2_224:
			m = "SHA2_224"
		case crypto.SHA2_256:
			m = "SHA2_256"
		case crypto.SHA2_384:
			m = "SHA2_384"
		case crypto.SHA2_512:
			m = "SHA2_512"
		case crypto.SHA3_224:
			m = "SHA3_224"
		case crypto.SHA3_256:
			m = "SHA3_256"
		case crypto.SHA3_384:
			m = "SHA3_384"
		case crypto.SHA3_512:
			m = "SHA3_512"
		case crypto.KECCAK_224:
			m = "KECCAK_224"
		case crypto.KECCAK_256:
			m = "KECCAK_256"
		case crypto.KECCAK_384:
			m = "KECCAK_384"
		case crypto.KECCAK_512:
			m = "KECCAK_512"
		case crypto.SM3:
			m = "SM3"
		case crypto.Sm3WithPublicKey:
			m = "Sm3WithPublicKey"
		}
	//Asymmetric Algo
	case mode&0xffff00ff == 0:
		switch mode {
		case crypto.Sm2p256v1:
			m = "Sm2p256v1"
		case crypto.Secp256k1:
			m = "Secp256k1"
		case crypto.Secp256r1:
			m = "Secp256r1"
		case crypto.Secp384r1:
			m = "Secp384r1"
		case crypto.Secp521r1:
			m = "Secp521r1"
		case crypto.Secp256k1Recover:
			m = "Secp256k1Recover"
		case crypto.Ed25519:
			m = "Ed25519"
		case crypto.Rsa2048:
			m = "Rsa2048"
		case crypto.Rsa3072:
			m = "Rsa3072"
		case crypto.Rsa4096:
			m = "Rsa4096"
		default:
			m = "None"
		}
	//Symmetrical Algo for Encrypt and Decrypt
	case mode&0xff00ffff == 0:
		switch mode {
		case crypto.Sm4 | crypto.CBC:
			m = "SM4_CBC"
		case crypto.Sm4 | crypto.ECB:
			m = "SM4_ECB"
		case crypto.Aes | crypto.CBC:
			m = "AES_CBC"
		case crypto.Aes | crypto.ECB:
			m = "AES_ECB"
		case crypto.Aes | crypto.GCM:
			m = "AES_GCM"
		case crypto.Des3 | crypto.CBC:
			m = "3DES_CBC"
		case crypto.Des3 | crypto.ECB:
			m = "3DES_ECB"
		case crypto.Des3 | crypto.GCM:
			m = "3DES_GCM"
		case crypto.TEE:
			m = "TEE"
		}
	case mode == 0xffffffff:
		m = "DEFAULT ALGO"
	}
	if m == "" {
		return "None"
	}
	return m
}

//GetKey get key
func GetKey(function Function, mode int) uint64 {
	return uint64(function)<<32 | uint64(mode)
}

//GetModeFromKey get mod from key
func GetModeFromKey(k uint64) int {
	return int(k & 0x00000000ffffffff)
}

//GetRandomStr get random string with length
func GetRandomStr(length int) string {
	ret := make([]byte, length)
	_, _ = rand.Read(ret)
	return hex.EncodeToString(ret)
}

//ModeGetCurve get curve form mode
func ModeGetCurve(mode int) (elliptic.Curve, error) {
	switch mode {
	case crypto.Secp521r1:
		return elliptic.P521(), nil
	case crypto.Secp384r1:
		return elliptic.P384(), nil
	case crypto.Secp256r1:
		return elliptic.P256(), nil
	case crypto.Secp256k1, crypto.Secp256k1Recover:
		return secp256k1.S256(), nil
	case crypto.Sm2p256v1:
		return gm.GetSm2Curve(), nil
	default:
		return nil, fmt.Errorf("unknown mode")
	}
}

//ModeFromCurve get mode from curve
func ModeFromCurve(curve elliptic.Curve) int {
	switch curve {
	case elliptic.P521():
		return crypto.Secp521r1
	case elliptic.P384():
		return crypto.Secp384r1
	case elliptic.P256():
		return crypto.Secp256r1
	case secp256k1.S256():
		return crypto.Secp256k1
	case gm.GetSm2Curve():
		return crypto.Sm2p256v1
	default:
		return crypto.None
	}
}
