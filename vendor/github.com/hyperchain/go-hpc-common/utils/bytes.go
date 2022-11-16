// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
)

var letterBytes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
	//EncMagic enc type
	EncMagic = "enc"
	//MagicNumberLen len of magic number
	encTypLen = 4
)

const (
	//MultiHashNumber multi hash code
	MultiHashNumber uint32 = iota + 1
)

// ToHex converts []bytes to hex string
func ToHex(b []byte) string {
	h := BytesToHex(b)
	// Prefer output of "0x0" instead of "0x"
	if len(h) == 0 {
		h = "0"
	}
	return "0x" + h
}

// DecodeString converts hex string to []byte
func DecodeString(s string) []byte {
	if len(s) >= len(EncMagic)+encTypLen && s[0:len(EncMagic)] == EncMagic {
		b := make([]byte, encTypLen)
		//nolint
		base64.RawStdEncoding.Decode(b[1:], []byte(s[len(EncMagic):len(EncMagic)+encTypLen]))
		magicNumber := binary.BigEndian.Uint32(b)
		switch magicNumber {
		case MultiHashNumber:
			if len(s) == len(EncMagic)+encTypLen {
				return nil
			}
			return multiHashDecoding(s[len(EncMagic)+encTypLen:])
		default:
			return nil
		}
	}
	if len(s) > 1 {
		if s[0:2] == "0x" {
			s = s[2:]
		}
		if len(s)%2 == 1 {
			s = "0" + s
		}
		return HexToBytes(s)
	}
	return nil
}

//multiHashDecoding base64 decoding rule
func multiHashDecoding(s string) []byte {
	ret, _ := base64.RawStdEncoding.DecodeString(s[:])
	return ret
}

// NumberToBytes returns the number in bytes with the specified base
func NumberToBytes(num interface{}, bits int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		fmt.Println("NumberToBytes failed:", err)
	}

	return buf.Bytes()[buf.Len()-(bits/8):]
}

// BytesToNumber attempts to cast a byte slice to a unsigned integer
func BytesToNumber(b []byte) uint64 {
	var number uint64

	if len(b) > 8 {
		return 0
	}
	// Make sure the buffer is 64bits
	data := make([]byte, 8)
	data = append(data[:8-len(b)], b...)

	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.BigEndian, &number)
	if err != nil {
		fmt.Println("BytesToNumber failed:", err)
	}

	return number
}

// ReadVarInt reads a variable length number in big endian byte order
func ReadVarInt(buff []byte) (ret uint64) {
	switch l := len(buff); {
	case l > 4:
		d := LeftPadBytes(buff, 8)
		_ = binary.Read(bytes.NewReader(d), binary.BigEndian, &ret)
	case l > 2:
		var num uint32
		d := LeftPadBytes(buff, 4)
		_ = binary.Read(bytes.NewReader(d), binary.BigEndian, &num)
		ret = uint64(num)
	case l > 1:
		var num uint16
		d := LeftPadBytes(buff, 2)
		_ = binary.Read(bytes.NewReader(d), binary.BigEndian, &num)
		ret = uint64(num)
	default:
		var num uint8
		_ = binary.Read(bytes.NewReader(buff), binary.BigEndian, &num)
		ret = uint64(num)
	}

	return
}

// CopyBytes returns an exact copy of the provided bytes
func CopyBytes(b []byte) (copiedBytes []byte) {
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

// IsHex returns if the given string is in hex format
func IsHex(str string) bool {
	l := len(str)
	return l >= 4 && l%2 == 0 && str[0:2] == "0x"
}

// BytesToHex converts []byte to hex string
func BytesToHex(d []byte) string {
	return hex.EncodeToString(d)
}

// HexToBytes converts hex string to []byte
func HexToBytes(str string) []byte {
	if len(str) >= 2 && str[0:2] == "0x" {
		str = str[2:]
	}
	h, _ := hex.DecodeString(str)

	return h
}

// Hex2Bytes converts hex string to []byte
func Hex2Bytes(str string) []byte {
	if len(str) >= 2 && str[0:2] == "0x" {
		str = str[2:]
	}
	h, _ := hex.DecodeString(str)

	return h
}

// Hex2BytesWithError converts hex string to []byte
func Hex2BytesWithError(str string) ([]byte, error) {
	if len(str) >= 2 && str[0:2] == "0x" {
		str = str[2:]
	}
	h, err := hex.DecodeString(str)

	return h, err
}

// FormatData format the given data by remove '\' and '0x'
func FormatData(data string) []byte {
	if len(data) == 0 {
		return nil
	}
	// Simple stupid
	d := new(big.Int)
	if data[0:1] == "\"" && data[len(data)-1:] == "\"" {
		return RightPadBytes([]byte(data[1:len(data)-1]), 32)
	} else if len(data) > 1 && data[:2] == "0x" {
		d.SetBytes(HexToBytes(data[2:]))
	} else {
		d.SetString(data, 0)
	}

	return BigToBytes(d, 256)
}

// ParseData parse the given interface value to []byte
func ParseData(data ...interface{}) (ret []byte) {
	for _, item := range data {
		switch t := item.(type) {
		case string:
			var str []byte
			if IsHex(t) {
				str = HexToBytes(t[2:])
			} else {
				str = []byte(t)
			}

			ret = append(ret, RightPadBytes(str, 32)...)
		case []byte:
			ret = append(ret, LeftPadBytes(t, 32)...)
		}
	}

	return
}

// RightPadBytes right pad the given slice with l length
func RightPadBytes(slice []byte, l int) []byte {
	if l < len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[0:len(slice)], slice)

	return padded
}

// LeftPadBytes left pad the given slice with l length
func LeftPadBytes(slice []byte, l int) []byte {
	if l < len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}

// LeftPadString left pad the given string with l length
func LeftPadString(str string, l int) string {
	if l < len(str) {
		return str
	}

	zeros := BytesToHex(make([]byte, (l-len(str))/2))

	return zeros + str

}

// RightPadString right pad the given string with l length
func RightPadString(str string, l int) string {
	if l < len(str) {
		return str
	}

	zeros := BytesToHex(make([]byte, (l-len(str))/2))

	return str + zeros

}

// PaddedBigBytes encodes a big integer as a big-endian byte slice. The length
// of the slice is at least n bytes.
func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	ReadBits(bigint, ret)
	return ret
}

// RandBytes return a slice with size random int
func RandBytes(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

// ReadBits encodes the absolute value of bigint as big-endian bytes. Callers must ensure
// that buf has enough space. If buf is too short the result will be incomplete.
func ReadBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}

// BytesToInt32 converts bytes of hex string to int32
func BytesToInt32(bytes []byte) int {
	if len(bytes) > 4 {
		return -1
	}
	result := 0
	for _, b := range bytes {
		result <<= 8
		result |= int(b & 0xff)
	}
	return result
}

// IntToBytes2 convert int to [2]byte
// NOTE: i should less than 2^15
func IntToBytes2(i int) [2]byte {
	result := [2]byte{}
	result[0] = (byte)((i >> 8) & 0xff)
	result[1] = (byte)(i & 0xff)
	return result
}

// IntToBytes4 convert int to [4]byte
// NOTE: i should less than 2^31
func IntToBytes4(i int) [4]byte {
	result := [4]byte{}
	result[0] = (byte)((i >> 24) & 0xff)
	result[1] = (byte)((i >> 16) & 0xff)
	result[2] = (byte)((i >> 8) & 0xff)
	result[3] = (byte)(i & 0xff)
	return result
}
