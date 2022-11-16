package types

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"reflect"
	"strings"

	"github.com/hyperchain/go-hpc-common/utils"
)

//HashLength default hash length
const HashLength = 32

//HashConsensus consensus temporary use
type HashConsensus [HashLength]byte

// Hash type used in Flato
type Hash []byte

// BytesToHash converts []byte to Hash
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// StringToHash converts string to Hash
func StringToHash(s string) Hash { return BytesToHash(utils.DecodeString(s)) }

//Str print string, now use String()
// Deprecated
func (h Hash) Str() string { return h.String() }

// Bytes is the []byte representation of the underlying hash
func (h Hash) Bytes() []byte {
	//Compatibility processing, when the h length is 0, that is, hash{}, should return [32] byte {}
	if len(h) == 0 {
		//return value
		emptyHash := [HashLength]byte{}
		return emptyHash[:]
	}
	//return pointer
	return h
}

//Bytes32 return 32 bytes
func (h Hash) Bytes32() [32]byte {
	var h32 [HashLength]byte
	if len(h) > HashLength {
		//return value
		copy(h32[:], h[len(h)-HashLength:])
		return h32
	}
	//return value
	copy(h32[HashLength-len(h):], h)
	return h32
}

// UnmarshalJSON parses a hash in its hex from to a hash.
func (h *Hash) UnmarshalJSON(input []byte) error {
	length := len(input)
	if length >= 2 && input[0] == '"' && input[length-1] == '"' {
		input = input[1 : length-1]
	}
	// strip "0x" for length check
	if len(input) > 1 && strings.ToLower(string(input[:2])) == "0x" {
		input = input[2:]
	}
	//match old rule
	if len(input) < len(utils.EncMagic) {
		return errHashJSONLength
	}

	if strings.ToLower(string(input[:len(utils.EncMagic)])) != utils.EncMagic {
		if len(input) != HashLength*2 {
			return errHashJSONLength
		}
	}

	h.SetBytes(utils.DecodeString(string(input)))
	return nil
}

// MarshalJSON marshal the Hash to json
func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// SetBytes sets the hash to the value of b. If b is larger than len(h) it will panic
func (h *Hash) SetBytes(b []byte) {
	if len(b) == 0 {
		*h = Hash{}.Bytes()
		return
	}
	if len(b) < HashLength {
		*h = make(Hash, HashLength)
		copy((*h)[HashLength-len(b):], b)
	} else {
		*h = make(Hash, len(b))
		copy(*h, b)
	}
}

// Set sets other to h
func (h *Hash) Set(other Hash) {
	*h = make(Hash, len(other))
	copy(*h, other)
}

//String return hash string
func (h Hash) String() string {
	if len(h) == HashLength || len(h) == 0 {
		return "0x" + utils.BytesToHex(h.Bytes())
	}
	return multiEncodeToString(h)
}

// Generate generate hash reflect
func (h *Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	if len(*h) == 0 {
		*h = make(Hash, HashLength)
	}
	m := rand.Intn(HashLength)
	for i := HashLength - 1; i > m; i-- {
		(*h)[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(*h)
}

// EmptyHash judges if the given hash is empty
func EmptyHash(h Hash) bool {
	//Compatibility processing, for example [32] byte {0}, should also be empty Hash
	emptyHash := [HashLength]byte{}
	if bytes.Equal(h, emptyHash[:]) {
		return true
	}
	return len(h) == 0
}

// FullHash returns a full Hash(all bits are set)
func FullHash() Hash {
	h := make(Hash, HashLength)
	for i := 0; i < HashLength; i++ {
		h[i] = 0xff
	}
	return h
}

//String HashConsensus string
func (h HashConsensus) String() string {
	return hex.EncodeToString(h[:])
}

//multiEncodeToString base64 encoding rule
func multiEncodeToString(b []byte) string {
	num := make([]byte, 4)
	binary.BigEndian.PutUint32(num, utils.MultiHashNumber)
	num = num[1:]

	magic := base64.RawStdEncoding.EncodeToString(num)
	ret := base64.RawStdEncoding.EncodeToString(b)

	return utils.EncMagic + magic + ret
}
