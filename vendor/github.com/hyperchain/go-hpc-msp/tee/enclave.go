package tee

import "sync"

var eid uint64
var once sync.Once

const (
	//Secp256k1 is sign type
	Secp256k1 = "secp256k1"
	//aesBlockSize
	aesBlockSize = 16
	//bodyMaxSize https header
	bodyMaxSize = 1024 * 1024
	//errStatusCode https response status
	errStatusCode = 303
	//cipherLength encrypt cipher length
	cipherLength = 560
)

func loadEnclave(path string) error {
	var err error
	once.Do(func() {
		err = getEnclaveLibSignSo(path)
		if err != nil {
			return
		}
		absPath := toAbs(path)
		//global eid
		eid, err = start(absPath)
	})
	if err != nil {
		return err
	}
	return nil
}
