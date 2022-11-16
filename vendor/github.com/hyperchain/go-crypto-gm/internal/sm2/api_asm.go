//go:build (amd64 || arm64) && !gmnoasm
// +build amd64 arm64
// +build !gmnoasm

package sm2

import (
	"crypto/elliptic"
	"io"
	"math/big"

	"github.com/hyperchain/go-crypto-gm/internal/sm2/internal"
	"golang.org/x/sys/cpu"
)

/*
API
func Sm2() elliptic.Curve
func Sign(dgst []byte, reader io.Reader, key []byte) ([]byte, uint8, error)
func verifySignature(sig, dgst []byte, X []byte, Y []byte) (bool, error)
func GetBatchHeap() interface{}
func PutBatchHeap(in interface{})
func BatchVerifyInit(ctxin interface{}, publicKey, signature, msg [][]byte) bool
func BatchVerifyEnd(ctxin interface{}) bool
func BatchVerify(publicKey, signature, msg [][]byte) error
*/

//Sm2 sm2 curve
var Sm2 func() elliptic.Curve

//Sign generate signature
var Sign func(dgst []byte, reader io.Reader, key []byte) ([]byte, uint8, error)

//VerifySignature verify
var VerifySignature func(sig, dgst []byte, X []byte, Y []byte) (bool, error)

//GetBatchHeap get data heap
var GetBatchHeap func() interface{}

//PutBatchHeap put data heap
var PutBatchHeap func(in interface{})

//BatchVerifyInit batch verify init
var BatchVerifyInit func(ctxin interface{}, publicKey, signature, msg [][]byte) error

//BatchVerifyEnd batch verify finish
var BatchVerifyEnd func(ctxin interface{}) error

//BatchVerify batch verify
var BatchVerify func(publicKey, signature, msg [][]byte) error

//ComplementCoordinates get corrdinates
var ComplementCoordinates func(x *big.Int, tildeY byte) error //the result is placed in 'x'

func init() {
	if (cpu.X86.HasSSE2 && cpu.X86.HasBMI2 && cpu.X86.HasSSE42 && cpu.X86.HasSSE41) || cpu.ARM64.HasAES {
		Sm2 = sm264bit
		Sign = sign64bit
		VerifySignature = verifySignature64bit
		GetBatchHeap = getBatchHeap64bit
		PutBatchHeap = putBatchHeap64bit
		BatchVerifyInit = batchVerifyInit64bit
		BatchVerifyEnd = batchVerifyEnd64bit
		BatchVerify = batchVerify64bit
		ComplementCoordinates = complementCoordinates64bit
		return
	}

	Sm2 = internal.Sm2_32bit
	Sign = internal.Sign32bit
	VerifySignature = internal.VerifySignature32bit
	GetBatchHeap = internal.GetBatchHeap32bit
	PutBatchHeap = internal.PutBatchHeap32bit
	BatchVerifyInit = internal.BatchVerifyInit_32bit
	BatchVerifyEnd = internal.BatchVerifyEnd_32bit
	BatchVerify = internal.BatchVerify_32bit
	ComplementCoordinates = internal.ComplementCoordinates32bit
}

//MarshalSig marshal signature
func MarshalSig(x, y []byte) []byte

//Unmarshal unmarshal signature
func Unmarshal(in []byte) (x []byte, y []byte)

//CompressCoordinates compress coordinates
func CompressCoordinates(in []byte, x, y *big.Int) {
	in[0] = byte(y.Bit(0)) + 2
	x.FillBytes(in[1:])
	return
}
