package bvmcom

import (
	"github.com/hyperchain/go-hpc-common/types"
)

var (
	// hashContractAddr the address of builtin hash contract
	hashContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 1})
	// proposalContractAddr the address of builtin proposal contract
	proposalContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 2})
	// accountContractAddr the address of builtin account contract
	accountContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 4})
	// certContractAddr the address of builtin cert contract
	certContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 5})
	// didContractAddr the address of builtin did set chain contract
	didContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 6})
	// crossChainNormalContractAddr the address of builtin non-system namespace cross-chain contract
	crossChainNormalContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 8})
	// mpcContractAddr the address of builtin mpc contract
	mpcContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 9})
	// crossChainSystemContractAddr the address of builtin system namespace cross-chain contract
	crossChainSystemContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 0x0a})
	// rootCAContractAddr the address of builtin rootCA contract
	rootCAContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 0x0b})
	// versionContractAddr the address of builtin system upgrade contract
	versionContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 0x0c})
	// hashChangeContractAddr the address of builtin hash change contract
	hashChangeContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 0x0d})
	// gasManagerContractAddr the address of builtin gas manager contract
	gasManagerContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 0x0e})
	// didCredentialContractAddr the address of did credential contract
	didCredentialContractAddr = types.BytesToAddress([]byte{0xff, 0xff, 0x0f})

	deployAddr = types.BytesToAddress([]byte{200, 179, 7, 153, 211, 68, 254, 61, 184, 194, 95, 217, 122, 216, 124, 245, 24, 33, 35, 109})
)

func getAddrCopy(addr types.Address) (ret types.Address) {
	copy(ret[:], addr[:])
	return
}

// DIDCredentialContractAddr the address of did credential contract
func DIDCredentialContractAddr() types.Address {
	return getAddrCopy(didCredentialContractAddr)
}

// HashContractAddr the address of builtin hash contract
func HashContractAddr() types.Address {
	return getAddrCopy(hashContractAddr)
}

// CertContractAddr the address of builtin cert contract
func CertContractAddr() types.Address {
	return getAddrCopy(certContractAddr)
}

// ProposalContractAddr the address of builtin proposal contract
func ProposalContractAddr() types.Address {
	return getAddrCopy(proposalContractAddr)
}

// AccountContractAddr the address of builtin account contract
func AccountContractAddr() types.Address {
	return getAddrCopy(accountContractAddr)
}

// DIDContractAddr the address of builtin did set chain contract
func DIDContractAddr() types.Address {
	return getAddrCopy(didContractAddr)
}

// MPCContractAddr the address of builtin mpc contract
func MPCContractAddr() types.Address {
	return getAddrCopy(mpcContractAddr)
}

//HashChangeContractAddr the address of hash change contract
func HashChangeContractAddr() types.Address {
	return getAddrCopy(hashChangeContractAddr)
}

// RootCAContractAddr the address of builtin rootCA contract
func RootCAContractAddr() types.Address {
	return getAddrCopy(rootCAContractAddr)
}

// CrossChainSystemContractAddr the address of builtin system namespace cross-chain contract
func CrossChainSystemContractAddr() types.Address {
	return getAddrCopy(crossChainSystemContractAddr)
}

// CrossChainNormalContractAddr the address of builtin non-system namespace cross-chain contract
func CrossChainNormalContractAddr() types.Address {
	return getAddrCopy(crossChainNormalContractAddr)
}

// GasManagerContractAddr the address of builtin namespace gas manager contract
func GasManagerContractAddr() types.Address {
	return getAddrCopy(gasManagerContractAddr)
}

// VersionContractAddr the address of builtin version contract
func VersionContractAddr() types.Address {
	return getAddrCopy(versionContractAddr)
}
