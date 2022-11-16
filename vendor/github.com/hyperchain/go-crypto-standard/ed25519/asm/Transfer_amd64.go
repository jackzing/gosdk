//go:build amd64
// +build amd64

package asm

//Get2DArray transfer 2D slice [][]byte to C style unsigned char**
func Get2DArray(out []uint64, in [][]byte)
