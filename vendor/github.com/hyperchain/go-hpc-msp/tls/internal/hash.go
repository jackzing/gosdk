package internal

import "hash"

//Copy hash Copy
func Copy() hash.Hash {
	return new(bytes)
}

type bytes []byte

func (data *bytes) Reset() {
	*data = nil
}

func (data *bytes) Size() int {
	return len(*data)
}

func (data *bytes) BlockSize() int {
	return 0
}

func (data *bytes) Write(p []byte) (n int, err error) {
	*data = append(*data, p...)
	return n, nil
}

func (data *bytes) Sum(b []byte) []byte {
	return append(*data, b...)
}
