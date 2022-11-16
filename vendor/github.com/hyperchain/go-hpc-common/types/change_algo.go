package types

import (
	"bytes"
	"encoding/binary"

	"github.com/meshplus/crypto"
)

//GetLatestAlgo return latest block number algo set
func (algoMap BlockAlgoSets) GetLatestAlgo() AlgoSet {
	if len(algoMap) == 0 {
		return AlgoSet{}
	}
	return algoMap[len(algoMap)-1].AlgoSet
}

// CheckAlgo check algo if valid
func CheckAlgo(algo string, method string, version TXVersion) bool {
	hashAlgo, encryptAlgo := GetSupportAlgo(version)
	switch method {
	case HashMethod:
		for i := range hashAlgo {
			if hashAlgo[i] == algo {
				return true
			}
		}
	case EncryptMethod:
		for i := range encryptAlgo {
			if encryptAlgo[i] == algo {
				return true
			}
		}
	}
	return false
}

// GetAlgoInt return algo int
func GetAlgoInt(name string) int {
	return algoToInt[name]
}

// GetSupportAlgo return concurrent support algo
func GetSupportAlgo(version TXVersion) ([]string, []string) {
	var hashAlgo []string
	var encryptAlgo []string

	for _, v := range hashVersionSet {
		if version.Compare(v.version) > -1 {
			hashAlgo = append(hashAlgo, v.algo)
		}
	}

	for _, v := range encryptVersionSet {
		if version.Compare(v.version) > -1 {
			encryptAlgo = append(encryptAlgo, v.algo)
		}
	}

	return hashAlgo, encryptAlgo
}

//MultiEncode The encoding method after changing the algorithm
func MultiEncode(algo int, data []byte) []byte {
	var a [encVerLen]byte
	binary.BigEndian.PutUint32(a[:], uint32(algo))
	switch {
	//hash
	case algo&0xffffff00 == 0:
		return bytes.Join([][]byte{{EncodeVersion}, a[:], data}, nil)
		//encrypt
	case algo&0xff00ffff == 0:
		return bytes.Join([][]byte{{EncodeVersion}, a[:], encMagic, data}, nil)
	}
	return nil
}

//MultiDecode decode multi-encode data
func MultiDecode(data []byte, method string) (int, []byte) {
	if len(data) < 5 {
		return 0, nil
	}

	switch method {
	case HashMethod:
		if len(data) == 32 || len(data) == 24 {
			return crypto.KECCAK_256, data
		}
		algo := binary.BigEndian.Uint32(data[1:5])
		return int(algo), data[5:]
	case MultiCacheEnc, FileLogEnc:
		indexStart := 1 + encVerLen
		indexEnd := 1 + encVerLen + len(encMagic)

		if (len(data) < indexEnd || !bytes.Equal(data[indexStart:indexEnd], encMagic)) && method == MultiCacheEnc {
			return crypto.Aes | crypto.CBC, data
		}
		if (len(data) < indexEnd || !bytes.Equal(data[indexStart:indexEnd], encMagic)) && method == FileLogEnc {
			return crypto.TEE, data
		}

		algo := binary.BigEndian.Uint32(data[1:indexStart])
		return int(algo), data[indexEnd:]
	}

	return 0, nil
}
