package common

import (
	"github.com/meshplus/crypto"
)

//ParseCertType unmarshal cert Type
func ParseCertType(certType string) crypto.CertType {
	for i := crypto.ECert; i < crypto.UnknownCertType; i++ {
		if i.String() == certType {
			return i
		}
	}
	return crypto.UnknownCertType
}
