package types

import (
	"bytes"
	"sort"

	"golang.org/x/crypto/sha3"
)

//  state & account proof related

// Inode struct
type Inode struct {
	Key   []byte `json:"key,omitempty"`
	Value []byte `json:"value,omitempty"`
	Hash  []byte `json:"hash,omitempty"`
}

// Inodes struct
type Inodes []*Inode

// ProofNode struct
type ProofNode struct {
	IsData bool   `json:"isData"`
	Key    []byte `json:"key,omitempty"`
	Hash   []byte `json:"hash,omitempty"`
	Inodes Inodes `json:"inodes,omitempty"`
	Index  int    `json:"index"`
}

// ProofPath struct
type ProofPath []*ProofNode

// Validate validates poof with given key
func Validate(key []byte, proof ProofPath) bool {
	if len(proof) == 0 || !proof[len(proof)-1].IsData {
		return false
	}
	var nextHash []byte
	for _, elem := range proof {
		current := elem
		if len(nextHash) != 0 && !bytes.Equal(nextHash, current.Hash) {
			return false
		}
		index := sort.Search(len(current.Inodes), func(i int) bool { return bytes.Compare(current.Inodes[i].Key, key) != -1 })
		exact := (len(current.Inodes) > 0 && index < len(current.Inodes) && bytes.Equal(current.Inodes[index].Key, key))
		if !exact {
			index--
		}
		if index != elem.Index || (current.IsData && !bytes.Equal(current.Inodes[index].Key, key)) {
			return false
		}
		res := CalProofNodeHash(current)
		if !bytes.Equal(res, current.Hash) {
			return false
		}
		nextHash = current.Inodes[index].Hash
	}
	return true
}

// CalProofNodeHash calculate hash for given ProofNode
func CalProofNodeHash(node *ProofNode) []byte {
	buff := make([]byte, 0)
	if node.IsData {
		for _, in := range node.Inodes {
			buff = append(buff, in.Key...)
			buff = append(buff, in.Value...)
		}
	} else {
		for _, in := range node.Inodes {
			buff = append(buff, in.Hash...)
		}
	}

	hasher := sha3.NewLegacyKeccak256()
	//nolint
	hasher.Write(buff)
	return hasher.Sum(nil)
}

// Tx proof related

// MerkleProofNode struct
type MerkleProofNode struct {
	Hash  []byte `json:"hash,omitempty"`
	Index int    `json:"index,omitempty"`
}

// MerkleProofPath struct
type MerkleProofPath []*MerkleProofNode

// ValidateMerkleProof validate merkleProof using targetTxHash and blockTxRootHash
func ValidateMerkleProof(proof []*MerkleProofNode, targetHash []byte, rootHash []byte) bool {
	if len(proof) == 1 {
		return bytes.Equal(proof[0].Hash, rootHash) && bytes.Equal(proof[0].Hash, targetHash)
	}
	hasher := sha3.NewLegacyKeccak256()
	var heads []int
	for i, p := range proof {
		if p.Index == 0 {
			heads = append(heads, i)
		}
	}
	if len(heads) == 0 {
		return false
	}
	// find targetHash in proof list
	findTargetHash := false
	for i := heads[len(heads)-1]; i < len(proof); i++ {
		if bytes.Equal(proof[i].Hash, targetHash) {
			findTargetHash = true
		}
	}
	if !findTargetHash {
		return false
	}
	// cal rootHash
	headlength := len(heads)
	proofLength := len(proof)
	for proofLength != 1 {
		var hashB []byte
		curIndex := heads[headlength-1]
		expectIndex := 0
		breakFlag := false

		for i := curIndex; i < proofLength; i++ {
			// ??????????????????index???????????????hash
			if expectIndex != proof[i].Index {
				//nolint
				hasher.Write(hashB)
				bs := hasher.Sum(nil)
				hasher.Reset()
				proof[curIndex] = &MerkleProofNode{
					Index: proof[i].Index - 1,
					Hash:  bs,
				}
				temp := curIndex + 1
				// ???????????????????????????
				for j := i; j < proofLength; j++ {
					proof[temp] = proof[j]
					temp++
				}
				// ??????curIndex???0????????????????????????index???0????????????
				if curIndex != 0 {
					headlength--
				}
				proofLength -= i - curIndex - 1
				breakFlag = true
				break
			}
			hashB = append(hashB, proof[i].Hash...)
			expectIndex++
		}
		// ??????????????????????????????????????????????????????
		if !breakFlag {
			//nolint
			hasher.Write(hashB)
			bs := hasher.Sum(nil)
			hasher.Reset()
			tempIndex := 0
			if curIndex != 0 {
				// ??????????????????????????????????????????
				tempIndex = proof[curIndex-1].Index + 1
			}
			proof[curIndex] = &MerkleProofNode{
				Index: tempIndex,
				Hash:  bs,
			}
			// ??????curIndex???0????????????????????????index???0??????????????????????????????heads???????????????????????????????????????
			if curIndex != 0 {
				headlength--
			}
			proofLength -= proofLength - curIndex - 1
		}
	}
	root := proof[0].Hash

	return bytes.Equal(root, rootHash)
}
