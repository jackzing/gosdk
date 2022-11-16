package types

// version_util.go used to define version information of node layer or utility methods.

import (
	"encoding/json"
	"fmt"
	"sort"
)

// VersionKey is the key for Version's storage
var VersionKey = []byte("version-key-")

const (
	// DBVersion represents the version of DB
	DBVersion = "1.1"
	// TxMetaVersion represents tx meta
	TxMetaVersion = TxMetaVersion10
)

const (
	// DBVersion00 use flato-db 0.2.7
	DBVersion00 = "0.0"

	// DBVersion10 use flato-db 0.2.8, add metaDB
	DBVersion10 = "1.0"

	// DBVersion11 use flato-db 0.2.10, add didDB and credentialDB
	DBVersion11 = "1.1"
)

const (
	// TxMetaVersion10 use flato-1.2.0
	TxMetaVersion10 uint64 = 1
)

// Version represents the version of flato, including DB...
type Version struct {
	// DBVersion if the version of flatodb
	DBVersion string `json:"dbVersion"`
}

// MarshalVersion marshals Version to bytes
func MarshalVersion(v *Version) ([]byte, error) {
	return json.Marshal(v)
}

// UnmarshalVersion unmarshals bytes to Version
func UnmarshalVersion(data []byte, v *Version) error {
	return json.Unmarshal(data, v)
}

// ========================================== TxVersion ===================================

// TXVersion type of txVersion
type TXVersion string

// Compare the value of TXVersion , return 1 if a > b ,0 if a = b, -1 if a < b
func (a TXVersion) Compare(b TXVersion) int {
	return CompareVersion(string(a), string(b))
}

var (
	// ErrInvalidTxVersion invalid tx version error
	ErrInvalidTxVersion = func(version interface{}) error {
		return fmt.Errorf("invalid tx version:%v", version)
	}
)

// GetTXVersion return the TXVersion value of tx version
func GetTXVersion(version string) (TXVersion, error) {
	_, ok := ChainVersionList[TXVersionTag][version]
	if !ok {
		return "", ErrInvalidTxVersion(version)
	}
	return TXVersion(version), nil
}

// ========================================== BlkVersion ===================================

// BlkVersion type of blockVersion
type BlkVersion string

// Compare the value of BlkVersion , return 1 if a > b ,0 if a = b, -1 if a < b
func (a BlkVersion) Compare(b BlkVersion) int {
	return CompareVersion(string(a), string(b))
}

var (
	// ErrInvalidBlkVersion invalid block version error
	ErrInvalidBlkVersion = func(version interface{}) error {
		return fmt.Errorf("invalid block version:%v", version)
	}
)

// GetBlkVersion return the BlkVersion value of tx version
func GetBlkVersion(version string) (BlkVersion, error) {
	if version == "" {
		return "0.0", nil
	}
	_, ok := ChainVersionList[BlockVersionTag][version]
	if !ok {
		return "", ErrInvalidBlkVersion(version)
	}
	return BlkVersion(version), nil
}

// ========================================== util function ===================================

// CompareVersion compares the value of a and b , return 1 if a > b ,0 if a = b, -1 if a < b
func CompareVersion(a, b string) int {
	indexa, indexb := 0, 0
	parta, partb := 0, 0
	lengtha, lengthb := len(a), len(b)
	f := func(str string, index, length int) (int, int) {
		part := 0
		for ; index < length; index++ {
			if str[index] != '.' {
				part = part << 4
				part += int(str[index] - '0')
			} else {
				index++
				break
			}
		}
		return part, index
	}

	for indexa < lengtha || indexb < lengthb {
		parta, indexa = f(a, indexa, lengtha)
		partb, indexb = f(b, indexb, lengthb)
		if parta > partb {
			return 1
		} else if parta < partb {
			return -1
		}
	}
	return 0
}

// UseNewHashFunc returns whether this tx used new hash function or not
func UseNewHashFunc(v []byte) bool {
	version := TXVersion(v)
	if v == nil || len(v) == 0 || version.Compare(TxVersion24) > 0 {
		return true
	}
	return false
}

// PersistBlkMeta returns whether this blk need to persist block meta
// Deprecated
func PersistBlkMeta(v []byte) bool {
	version := BlkVersion(v)
	if v == nil || len(v) == 0 || version.Compare(BlkVersion24) > 0 {
		return true
	}
	return false
}

type stringVersionArr []string

func (sa stringVersionArr) Len() int { return len(sa) }

func (sa stringVersionArr) Swap(i, j int) { sa[i], sa[j] = sa[j], sa[i] }

func (sa stringVersionArr) Less(i, j int) bool {
	result := CompareVersion(sa[i], sa[j])
	if result >= 0 {
		return false
	}
	return true
}

// GetSortedVersionList get version map according to the given tag and
// then convert it to sorted string array.
func GetSortedVersionList(tag VersionTag) []string {
	mp := ChainVersionList[tag]
	list := make(stringVersionArr, len(mp))
	i := 0
	for k := range mp {
		list[i] = k
		i++
	}
	sort.Sort(list)
	return list
}

// GetHyperchainVersionByRunningVersion returns hyperchain version according to
// the given running version.
func GetHyperchainVersionByRunningVersion(rv RunningVersion) string {
	if rv == nil || len(rv) == 0 {
		return ""
	}

	var version string

loop:
	for hv, hrvs := range HyperchainReleaseVersion {
		if len(hrvs) != len(rv) {
			continue loop
		}
		for tag, hrv := range hrvs {
			rvForTag, aOK := rv[tag]
			if !aOK {
				continue loop
			}
			if rvForTag != hrv {
				continue loop
			}
		}
		version = hv
	}

	return version
}
