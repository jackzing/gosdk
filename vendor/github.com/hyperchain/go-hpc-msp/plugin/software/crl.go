package software

import (
	"bytes"
	x509std "crypto/x509"
	pkixStd "crypto/x509/pkix"
	"encoding/asn1"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/hyperchain/go-hpc-msp/plugin/pkix"
	"github.com/meshplus/crypto"
)

//CFCA mode
const (
	CFCAModeNone = "none"
	CFCAModeRA   = "ra"
	CFCAModeCRL  = "crl"
)

//ra response status
const (
	RACertStatusNotDownloaded = 3 + iota
	RACertStatusValid
	RACertStatusFrozen
	RACertStatusRevoked
)

//error derEncode
const (
	ERRUnknownStatus = "unknown certificate status"
	ERRNotDownload   = "certificate not downloaded"
	ERRRevoked       = "certificate has been revoked"
	ERRFrozen        = "the certificate has been frozen"
	ERRNotExit       = "certificate does not exist"
)

//CRL CRL is a Thread-safe certificate revocation list
type CRL struct {
	url   string
	crl   *pkixStd.CertificateList
	mutex *sync.RWMutex
}

//NewCRL create a crl Instance
func NewCRL(url string, quit <-chan bool) (*CRL, error) {
	result := &CRL{
		url:   url,
		mutex: new(sync.RWMutex),
	}
	if err := result.update(); err != nil {
		return nil, err
	}
	go func() {
		for {
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 8, 0, rand.Intn(60), 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			select {
			case <-t.C:
				err := result.update()
				if err != nil {
					continue
				}
			case <-quit:
				return
			}
		}
	}()
	return result, nil
}

//CheckRevocation Verify that the certificate has been revoked
func (c *CRL) CheckRevocation(cert *Certificate) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return CheckRevocationWithCRL(cert, c.crl)
}

func (c *CRL) update() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	newCRL, err := FetchCRL(c.url)
	if err != nil {
		return err
	}
	c.crl = newCRL
	return nil
}

// FetchCRL fetch CRL with a HTTP GET request to a given URL
// returns a reference to a pkix.CertificateList instance
func FetchCRL(url string) (cl *pkixStd.CertificateList, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to fetch CRL,the status code is %v", resp.StatusCode)
	}

	crl, err := ioutil.ReadAll(resp.Body)
	defer func() {
		err = resp.Body.Close()
	}()
	if err != nil {
		return nil, err
	}

	return x509std.ParseDERCRL(crl)
}

//CheckRevocationWithCRL CheckRevocation verifies if a given certificate is revoked in
// reference to current CRL
func CheckRevocationWithCRL(cert *Certificate, crl *pkixStd.CertificateList) (bool, error) {
	if cert == nil {
		return false, errors.New("invalid cert")
	}
	if crl == nil {
		return false, errors.New("invalid crl")
	}

	for _, i := range crl.TBSCertList.RevokedCertificates {
		if cert.SerialNumber.Cmp(i.SerialNumber) == 0 {
			return true, nil
		}
	}
	return false, nil
}

// CheckRevocation verifies if a given certificate is revoked via
// its CRL distribution point
func CheckRevocation(cert *Certificate) (bool, error) {
	for _, url := range cert.CRLDistributionPoints {
		crl, err := FetchCRL(url)
		if err != nil {
			continue
		} else {
			revoked, err := CheckRevocationWithCRL(cert, crl)
			if err != nil {
				return false, err
			}
			return revoked, nil
		}
	}
	return false, errors.New("failed to check revocation state")
}

//CheckRevocationWithURL check revocation from specific url
func CheckRevocationWithURL(cert *Certificate, url string) (bool, error) {
	crl, err := FetchCRL(url)
	if err != nil {
		return false, err
	}
	return CheckRevocationWithCRL(cert, crl)
}

//ra  http://ucrl.cfca.com.cn/OCA1/SM2/allCRL.crl
//crl http://114.55.107.52:8080/raWeb/CSHttpServlet

// CertificateList represents the ASN.1 structure of the same name. See RFC
// 5280, section 5.1. Use Certificate.CheckCRLSignature to verify the
// signature.
type CertificateList struct {
	TBSCertList        TBSCertificateList
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

// HasExpired reports whether certList should have been updated by now.
func (certList *CertificateList) HasExpired(now time.Time) bool {
	return !now.Before(certList.TBSCertList.NextUpdate)
}

// TBSCertificateList represents the ASN.1 structure of the same name. See RFC
// 5280, section 5.1.
type TBSCertificateList struct {
	Raw                 asn1.RawContent
	Version             int `asn1:"optional,default:0"`
	Signature           pkix.AlgorithmIdentifier
	Issuer              pkix.RDNSequence
	ThisUpdate          time.Time
	NextUpdate          time.Time            `asn1:"optional"`
	RevokedCertificates []RevokedCertificate `asn1:"optional"`
	Extensions          []pkix.Extension     `asn1:"tag:0,optional,explicit"`
}

// RevokedCertificate represents the ASN.1 structure of the same name. See RFC
// 5280, section 5.1.
type RevokedCertificate struct {
	SerialNumber   *big.Int
	RevocationTime time.Time
	Extensions     []pkix.Extension `asn1:"optional"`
}

// pemCRLPrefix is the magic string that indicates that we have a PEM encoded
// CRL.
var pemCRLPrefix = []byte("-----BEGIN X509 CRL")

// pemType is the type of a PEM encoded CRL.
var pemType = "X509 CRL"

// ParseCRL parses a CRL from the given bytes. It's often the case that PEM
// encoded CRLs will appear where they should be DER encoded, so this function
// will transparently handle PEM encoding as long as there isn't any leading
// garbage.
func ParseCRL(crlBytes []byte) (*CertificateList, error) {
	if bytes.HasPrefix(crlBytes, pemCRLPrefix) {
		block, _ := pem.Decode(crlBytes)
		if block != nil && block.Type == pemType {
			crlBytes = block.Bytes
		}
	}
	return ParseDERCRL(crlBytes)
}

// ParseDERCRL parses a DER encoded CRL from the given bytes.
func ParseDERCRL(derBytes []byte) (*CertificateList, error) {
	certList := new(CertificateList)
	if rest, err := asn1.Unmarshal(derBytes, certList); err != nil {
		return nil, err
	} else if len(rest) != 0 {
		return nil, errors.New("x509: trailing data after CRL")
	}
	return certList, nil
}

// CreateCRL returns a DER encoded CRL, signed by this Certificate, that
// contains the given list of revoked certificates.
func (c *Certificate) CreateCRL(rand io.Reader, priv crypto.SignKey, revokedCerts []RevokedCertificate, now, expiry time.Time) (crlBytes []byte, err error) {
	if reflect.DeepEqual(priv, nil) {
		return nil, fmt.Errorf("private key is nil")
	}
	hashFunc, signatureAlgorithm, err := signingParamsForPublicKey(priv, 0)
	if err != nil {
		return nil, err
	}

	// Force revocation times to UTC per RFC 5280.
	revokedCertsUTC := make([]RevokedCertificate, len(revokedCerts))
	for i, rc := range revokedCerts {
		rc.RevocationTime = rc.RevocationTime.UTC()
		revokedCertsUTC[i] = rc
	}

	tbsCertList := TBSCertificateList{
		Version:             1,
		Signature:           signatureAlgorithm,
		Issuer:              c.Subject.ToRDNSequence(),
		ThisUpdate:          now.UTC(),
		NextUpdate:          expiry.UTC(),
		RevokedCertificates: revokedCertsUTC,
	}

	// Authority Key ID
	if len(c.SubjectKeyID) > 0 {
		var aki pkix.Extension
		aki.ID = oidExtensionAuthorityKeyID
		aki.Value, err = asn1.Marshal(authKeyID{ID: c.SubjectKeyID})
		if err != nil {
			return
		}
		tbsCertList.Extensions = append(tbsCertList.Extensions, aki)
	}

	tbsCertListContents, err := asn1.Marshal(tbsCertList)
	if err != nil {
		return
	}

	var signature []byte
	//pkcs1v15
	if common.ModeIsRSAAlgo(priv.GetKeyInfo()) {
		var hashTypeUsedInPKCS1v15 [4]byte
		binary.BigEndian.PutUint32(hashTypeUsedInPKCS1v15[:], uint32(hashFunc))
		signature, err = priv.Sign(append(tbsCertListContents, hashTypeUsedInPKCS1v15[:]...), hashFunc.New(), rand)
	} else {
		signature, err = priv.Sign(tbsCertListContents, hashFunc.New(), rand)
	}
	if err != nil {
		return
	}

	return asn1.Marshal(CertificateList{
		TBSCertList:        tbsCertList,
		SignatureAlgorithm: signatureAlgorithm,
		SignatureValue:     asn1.BitString{Bytes: signature, BitLength: len(signature) * 8},
	})
}
