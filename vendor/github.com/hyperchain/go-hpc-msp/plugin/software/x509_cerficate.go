package software

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hyperchain/go-hpc-msp/pemencode"
	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/hyperchain/go-hpc-msp/plugin/pkix"
	"github.com/meshplus/crypto"
)

//X509Cert software x509 cert
type X509Cert struct {
	cert *Certificate
}

//String return printable Cert text
func (x *X509Cert) String() string {
	der, _ := MarshalCertificate(x.cert)
	pem, _ := pemencode.DER2PEM(der, pemencode.PEMCertificate)
	return string(pem)
}

//VerifyCert untrustedPubList is raw public
func (x *X509Cert) VerifyCert(caList []string, untrustedList [][]byte) (err error) {
	//1.检查是否有匹配CA
	var caCerts []*Certificate
	for i := range caList {
		if string(caList[i]) == common.DummyCAName {
			return nil
		}
		caCert, pErr := ParseCertificate(caList[i])
		if pErr != nil {
			//忽略不能解析的ca，它们可能来自其他插件因而采用了非x509格式
			continue
		}
		if bytes.Equal(caCert.SubjectKeyID, x.GetAuthorityKeyIdentifier()) {
			caCerts = append(caCerts, caCert)
		}
	}
	if len(caCerts) == 0 {
		return fmt.Errorf("can't find ca for cert '%v'", x.String())
	}
	//2.检查是否被本地或者链上吊销
	vk := x.GetVerifyKey()
	for i := range untrustedList {
		if bytes.Equal(vk.Bytes(), untrustedList[i]) || bytes.Equal(x.cert.Signature, untrustedList[i]) {
			return fmt.Errorf("cert is revoked by ledger: %v", hex.EncodeToString(untrustedList[i]))
		}
	}

	//3. 验证证书
	var pass = false
	var errDetail strings.Builder
	for i := range caCerts {
		if bytes.Equal(x.cert.AuthorityKeyID, caCerts[i].SubjectKeyID) {
			_, err = VerifyCert(x.cert, caCerts[i])
			if err != nil {
				errDetail.WriteString(fmt.Sprintf(`"ca %v":"%v"; `, caCerts[i].Subject.CommonName, err.Error()))
			}
		} else {
			errDetail.WriteString(fmt.Sprintf(`"ca %v":"not match, subcert's issuer is %v'"; `,
				caCerts[i].Subject.CommonName,
				pkix.GetIdentityNameFromPKIXName(x.cert.Issuer)))
		}

		pass = pass || (err == nil)
	}
	if !pass {
		return fmt.Errorf("verify certificate error: %v", errDetail.String())
	}

	//4. 是否CRL吊销
	if len(x.cert.CRLDistributionPoints) != 0 {
		// find in crl
		exist, perr := CheckRevocation(x.cert)
		if perr != nil {
			return fmt.Errorf("isRevoked error: CRL CheckRevocation error: %s", perr.Error())
		}
		if exist {
			return fmt.Errorf("cert revoked by CRL: [%v]", strings.Join(x.cert.CRLDistributionPoints, ", "))
		}
	}
	return nil
}

//GetCertType get cert type
func (x *X509Cert) GetCertType() crypto.CertType {
	name := pkix.GetIdentityNameFromPKIXName(x.cert.Subject)
	return name.GetCertType()
}

//GetHostName hostname is CN
func (x *X509Cert) GetHostName() string {
	return x.cert.Subject.CommonName
}

//GetCAHostName get ca hostname
func (x *X509Cert) GetCAHostName() string {
	return x.cert.Issuer.CommonName
}

//GetExtName get Org
func (x *X509Cert) GetExtName() map[string]string {
	name := pkix.GetIdentityNameFromPKIXName(x.cert.Subject)
	ret, _ := pkix.ParseOrganization(name.O)
	return ret
}

//GetAuthorityKeyIdentifier get AKI
func (x *X509Cert) GetAuthorityKeyIdentifier() []byte {
	return x.cert.AuthorityKeyID
}

//GetVerifyKey get public key
func (x *X509Cert) GetVerifyKey() crypto.VerifyKey {
	return x.cert.PublicKey
}

//X509CA software x509 ca
type X509CA struct {
	ca         *Certificate
	hostname   string
	ski        string
	privateKey crypto.SignKey
}

//GetHostName get hostname
func (x *X509CA) GetHostName() string {
	return x.hostname
}

//String return ca text
func (x *X509CA) String() string {
	der, _ := MarshalCertificate(x.ca)
	pem, _ := pemencode.DER2PEM(der, pemencode.PEMCertificate)
	return string(pem)
}

//GetKeyIdentifier get key identifier，i.e. SKI
func (x *X509CA) GetKeyIdentifier() []byte {
	return x.ca.SubjectKeyID
}

//GetPubKeyForPairing get PubKey for Pairing, it's typically 65Bytes
func (x *X509CA) GetPubKeyForPairing() []byte {
	return x.ca.PublicKey.Bytes()
}
