package software

import (
	"bytes"
	"crypto/rand"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"strings"

	"math/big"
	"net"
	"net/url"
	"time"

	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/hyperchain/go-hpc-msp/plugin/pkix"
	"github.com/meshplus/crypto"
)

// A Certificate represents an X.509 certificate.
type Certificate struct {
	Raw                     []byte // Complete ASN.1 DER derEncode (certificate, signature algorithm and signature).
	RawTBSCertificate       []byte // Certificate part of raw ASN.1 DER derEncode.
	RawSubjectPublicKeyInfo []byte // DER encoded SubjectPublicKeyInfo.
	RawSubject              []byte // DER encoded Subject
	RawIssuer               []byte // DER encoded Issuer

	Signature          []byte
	SignatureAlgorithm SignatureAlgorithm

	PublicKeyAlgorithm PublicKeyAlgorithm
	PublicKey          crypto.VerifyKey

	Version             int
	SerialNumber        *big.Int
	Issuer              pkix.Name
	Subject             pkix.Name
	NotBefore, NotAfter time.Time // Validity bounds.
	KeyUsage            KeyUsage

	// Extensions contains raw X.509 extensions. When parsing certificates,
	// this can be used to extract non-critical extensions that are not
	// parsed by this package. When marshaling certificates, the Extensions
	// field is ignored, see ExtraExtensions.
	Extensions []pkix.Extension

	// ExtraExtensions contains extensions to be copied, raw, into any
	// marshaled certificates. Values override any extensions that would
	// otherwise be produced based on the other fields. The ExtraExtensions
	// field is not populated when parsing certificates, see Extensions.
	ExtraExtensions []pkix.Extension

	// UnhandledCriticalExtensions contains a list of extension IDs that
	// were not (fully) processed when parsing. Verify will fail if this
	// slice is non-empty, unless verification is delegated to an OS
	// library which understands all the critical extensions.
	//
	// Users can access these extensions using Extensions and can remove
	// elements from this slice if they believe that they have been
	// handled.
	UnhandledCriticalExtensions []asn1.ObjectIdentifier

	ExtKeyUsage        []ExtKeyUsage           // Sequence of extended key usages.
	UnknownExtKeyUsage []asn1.ObjectIdentifier // Encountered extended key usages unknown to this package.

	// BasicConstraintsValid indicates whether IsCA, MaxPathLen,
	// and MaxPathLenZero are valid.
	BasicConstraintsValid bool
	IsCA                  bool

	// MaxPathLen and MaxPathLenZero indicate the presence and
	// value of the BasicConstraints' "pathLenConstraint".
	//
	// When parsing a certificate, a positive non-zero MaxPathLen
	// means that the field was specified, -1 means it was unset,
	// and MaxPathLenZero being true mean that the field was
	// explicitly set to zero. The case of MaxPathLen==0 with MaxPathLenZero==false
	// should be treated equivalent to -1 (unset).
	//
	// When generating a certificate, an unset pathLenConstraint
	// can be requested with either MaxPathLen == -1 or using the
	// zero value for both MaxPathLen and MaxPathLenZero.
	MaxPathLen int
	// MaxPathLenZero indicates that BasicConstraintsValid==true
	// and MaxPathLen==0 should be interpreted as an actual
	// maximum path length of zero. Otherwise, that combination is
	// interpreted as MaxPathLen not being set.
	MaxPathLenZero bool

	//使用公钥计算的sha1结果，对于sm2来说，是65字节长度裸公钥的sha1值
	SubjectKeyID   []byte
	AuthorityKeyID []byte

	// RFC 5280, 4.2.2.1 (Authority Information Access)
	OCSPServer            []string
	IssuingCertificateURL []string

	// Subject Alternate Name values. (Note that these values may not be valid
	// if invalid values were contained within a parsed certificate. For
	// example, an element of DNSNames may not be a valid DNS domain name.)
	DNSNames       []string
	EmailAddresses []string
	IPAddresses    []net.IP
	URIs           []*url.URL

	// Name constraints
	PermittedDNSDomainsCritical bool // if true then the name constraints are marked critical.
	PermittedDNSDomains         []string
	ExcludedDNSDomains          []string
	PermittedIPRanges           []*net.IPNet
	ExcludedIPRanges            []*net.IPNet
	PermittedEmailAddresses     []string
	ExcludedEmailAddresses      []string
	PermittedURIDomains         []string
	ExcludedURIDomains          []string

	// CRL Distribution Points
	CRLDistributionPoints []string

	PolicyIdentifiers []asn1.ObjectIdentifier
}

//ParseCertificate already support ra
// input is pem
func ParseCertificate(pemBytes string) (*Certificate, error) {
	//if input is pem format, try to parse
	der := []byte(pemBytes)
	block, _ := pem.Decode(der)

	if block != nil {
		der = block.Bytes
	}

	var cert certificate
	rest, err := asn1.Unmarshal(der, &cert)
	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, asn1.SyntaxError{Msg: "trailing data"}
	}

	x509Cert, err := parseCertificate(&cert)

	if err != nil {
		return nil, err
	}

	return x509Cert, nil
}

//MarshalCertificate rev to parseCertificate
func MarshalCertificate(template *Certificate) (cert []byte, err error) {
	if template == nil {
		return nil, errors.New("param is nil")
	}
	template.Version = template.Version - 1

	if template.SerialNumber == nil {
		return nil, errors.New("x509: no SerialNumber given")
	}

	var signatureAlgorithm, publicKeyAlgorithm pkix.AlgorithmIdentifier
	_, signatureAlgorithm, err = signParamsForPublicKey(template.SignatureAlgorithm)
	if err != nil {
		return nil, err
	}
	publicKeyAlgorithm, err = GetPublicKeyAlgorithmFromMode(template.PublicKey.GetKeyInfo())
	if err != nil {
		return nil, err
	}
	publicKeyBytes := template.PublicKey.Bytes()
	encodedPublicKey := asn1.BitString{
		Bytes:     publicKeyBytes,
		BitLength: 8 * len(publicKeyBytes),
	}

	var asn1Issuer, asn1Subject []byte
	asn1Issuer, err = issuerBytesWhenMarshal(template)
	if err != nil {
		return
	}
	asn1Subject, err = subjectBytesWhenMarshal(template)
	if err != nil {
		return
	}

	template.ExtraExtensions = template.Extensions
	extensions, err := buildExtensions(template, bytes.Equal(asn1Subject, emptyASN1Subject),
		template.AuthorityKeyID, template.SubjectKeyID)
	if err != nil {
		return
	}

	c := tbsCertificate{
		Version:            template.Version,
		SerialNumber:       template.SerialNumber,
		SignatureAlgorithm: signatureAlgorithm,
		Issuer:             asn1.RawValue{FullBytes: asn1Issuer},
		Validity:           validity{template.NotBefore.UTC(), template.NotAfter.UTC()},
		Subject:            asn1.RawValue{FullBytes: asn1Subject},
		PublicKey:          publicKeyInfo{nil, publicKeyAlgorithm, encodedPublicKey},
		Extensions:         extensions,
	}

	tbsCertContents, err := asn1.Marshal(c)
	if err != nil {
		return
	}

	c.Raw = tbsCertContents
	signature := template.Signature

	return asn1.Marshal(certificate{
		nil,
		c,
		signatureAlgorithm,
		asn1.BitString{Bytes: signature, BitLength: len(signature) * 8},
	})
}

//VerifyCert already support ra
func VerifyCert(cert *Certificate, ca *Certificate) (bool, error) {
	if cert.NotBefore.After(time.Now()) || cert.NotAfter.Before(time.Now()) {
		return false, errors.New("this cert is expired")
	}

	err := cert.CheckSignatureFrom(ca)
	if err != nil {
		return false, err
	}

	return true, nil
}

//GenCert generate cert from ca
func GenCert(ca *Certificate, privatekey crypto.SignKey, publicKey crypto.VerifyKey,
	o, cn, gn string, isCA bool, from, to time.Time, webAddr ...string) ([]byte, error) {

	if !bytes.Equal(ca.PublicKey.Bytes(), privatekey.Bytes()) {
		return nil, errors.New("public key in ca does not match private key")
	}

	return createCertByCaAndPublicKey(ca, privatekey, publicKey, isCA, o, cn, gn, from, to, webAddr...)
}

//NewSelfSignedCert generate self-signature certificate
func NewSelfSignedCert(o, cn, gn string, ct CurveType, from, to time.Time, webAddr ...string) (
	[]byte, string, error) {
	var (
		err                error
		mode               int
		signatureAlgorithm SignatureAlgorithm
		privKey            crypto.SignKey
		privIndex          string
	)
	switch ct {
	case CurveTypeSm2:
		signatureAlgorithm = SM3WithSM2
		mode = crypto.Sm2p256v1
	case CurveTypeP256:
		signatureAlgorithm = ECDSAWithSHA256
		mode = crypto.Secp256r1
	case CurveTypeK1:
		signatureAlgorithm = ECDSAWithSHA256
		mode = crypto.Secp256k1
	}
	privIndex, privKey, err = createPrivateKey(mode)
	if err != nil {
		return nil, "", err
	}

	t, err := generateTemplate(o, cn, gn, from, to, signatureAlgorithm, webAddr...)
	if err != nil {
		return nil, "", err
	}
	t.IsCA = true

	cert, err := CreateCertificate(rand.Reader, t, t, privKey, privKey)
	if err != nil {
		return nil, "", err
	}

	return cert, privIndex, nil
}

//SelfSignedCert generate self-signature certificate by privKey and pubKey
func SelfSignedCert(o, cn, gn string, webAddr []string, privKey crypto.SignKey, from, to time.Time) ([]byte, error) {
	var (
		err                error
		signatureAlgorithm SignatureAlgorithm
	)

	t, err := generateTemplate(o, cn, gn, from, to, signatureAlgorithm, webAddr...)
	if err != nil {
		return nil, err
	}
	t.IsCA = true

	cert, err := CreateCertificate(rand.Reader, t, t, privKey, privKey)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

//CertTypeOID oid fo certType
var CertTypeOID asn1.ObjectIdentifier = []int{1, 2, 86, 1}

//AssertCertType assert cert type with specified type，return boolean
func AssertCertType(expect crypto.CertType, certificate *Certificate) bool {
	for _, v := range certificate.Extensions {
		if CertTypeOID.Equal(v.ID) {
			value := common.ParseCertType(string(v.Value))
			if value == expect {
				return true
			}
			if value == crypto.ERCert && (expect == crypto.ECert || expect == crypto.RCert) {
				return true
			}
		}
	}
	return false
}

func createCertByCaAndPublicKey(ca *Certificate, caPrivate crypto.SignKey, subPublic crypto.VerifyKey, isCa bool,
	o, cn, gn string, from, to time.Time, webAddr ...string) (certDER []byte, err error) {
	var signatureAlgorithm SignatureAlgorithm
	//If the private key of ca is sm2,
	// the generated private key of the cert is also sm2.
	switch caPrivate.GetKeyInfo() {
	case crypto.Sm2p256v1:
		signatureAlgorithm = SM3WithSM2
	case crypto.Secp256k1, crypto.Secp256r1:
		signatureAlgorithm = ECDSAWithSHA256
	default:
		return nil, errors.New("private curve neither k1 nor r1")
	}
	template, gerr := generateTemplate(o, cn, gn, from, to, signatureAlgorithm, webAddr...)
	if gerr != nil {
		return nil, gerr
	}
	template.IsCA = isCa
	cert, cerr := CreateCertificate(rand.Reader, template, ca, subPublic, caPrivate)
	if cerr != nil {
		return nil, cerr
	}
	return cert, nil
}

func generateTemplate(o, cn, gn string, from, to time.Time, signatureAlgorithm SignatureAlgorithm, webAddr ...string) (*Certificate, error) {
	gn = strings.ToLower(gn)
	if gn != "ecert" && gn != "rcert" && gn != "sdkcert" && gn != "" && gn != "idcert" {
		return nil, errors.New("gn should be one of ecert, rcert, sdkcert, idcert or empty")
	}

	//parse SAN
	var IP []net.IP
	var DNSName []string
	if len(webAddr) > 0 {
		var URL []*url.URL
		IP, URL = parseSAN(webAddr)
		for _, u := range URL {
			DNSName = append(DNSName, u.String())
		}
	}

	Subject := pkix.Name{
		CommonName:         cn,
		Organization:       []string{o},
		OrganizationalUnit: []string{gn},
		Country:            []string{"CN"},
		ExtraNames:         []pkix.AttributeTypeAndValue{{Type: []int{2, 5, 4, 42}, Value: gn}},
	}
	random, _ := rand.Int(rand.Reader, big.NewInt(1<<63-1))
	template := &Certificate{
		SerialNumber: random,
		Subject:      Subject,

		NotBefore: from,
		NotAfter:  to,

		SignatureAlgorithm: signatureAlgorithm,
		KeyUsage: KeyUsageCertSign | KeyUsageDigitalSignature | KeyUsageCRLSign |
			KeyUsageContentCommitment | KeyUsageKeyEncipherment | KeyUsageKeyAgreement,
		ExtKeyUsage: []ExtKeyUsage{ExtKeyUsageClientAuth, ExtKeyUsageServerAuth,
			ExtKeyUsageCodeSigning, ExtKeyUsageEmailProtection},
		BasicConstraintsValid: true,
		IPAddresses:           IP,
		DNSNames:              DNSName,
	}

	t := common.ParseCertType(gn)
	if t != crypto.UnknownCertType {
		template.ExtraExtensions = append(template.ExtraExtensions,
			pkix.Extension{
				ID:    CertTypeOID,
				Value: []byte(t.String()),
			})
	}
	return template, nil
}

func parseSAN(in []string) (IPAddresses []net.IP, URIs []*url.URL) {
	IPAddresses = make([]net.IP, 0)
	URIs = make([]*url.URL, 0)
	for _, v := range in {
		if ip := net.ParseIP(v); ip != nil {
			IPAddresses = append(IPAddresses, ip)
			continue
		}
		if u, err := url.Parse(v); err == nil {
			URIs = append(URIs, u)
			continue
		}
	}
	return
}

func parseCertificate(in *certificate) (out *Certificate, err error) {
	out = &Certificate{
		Raw:                     in.Raw,
		SerialNumber:            in.TBSCertificate.SerialNumber,
		Version:                 in.TBSCertificate.Version + 1,
		RawTBSCertificate:       in.TBSCertificate.Raw,
		RawSubjectPublicKeyInfo: in.TBSCertificate.PublicKey.Raw,
		RawSubject:              in.TBSCertificate.Subject.FullBytes,
		RawIssuer:               in.TBSCertificate.Issuer.FullBytes,
		Signature:               in.SignatureValue.RightAlign(),
	}

	out.SignatureAlgorithm = getSignatureAlgorithmFromAI(in.TBSCertificate.SignatureAlgorithm)
	out.PublicKeyAlgorithm = GetPublicKeyAlgorithmFromAlgorithmIdentifier(in.TBSCertificate.PublicKey.Algorithm)
	rawPub, mode, perr := ParsePKIXPublicKey(in.TBSCertificate.PublicKey.Raw)
	if perr != nil {
		return nil, perr
	}
	engine := GetSoftwareEngine()
	out.PublicKey, err = engine.GetVerifyKey(rawPub, mode)
	if err != nil {
		return nil, err
	}

	var issuer, subject pkix.RDNSequence
	var rest []byte
	if rest, err = asn1.Unmarshal(in.TBSCertificate.Subject.FullBytes, &subject); err != nil {
		return nil, err
	} else if len(rest) != 0 {
		return nil, errors.New("x509: trailing data after X.509 subject")
	}

	if rest, err = asn1.Unmarshal(in.TBSCertificate.Issuer.FullBytes, &issuer); err != nil {
		return nil, err
	} else if len(rest) != 0 {
		return nil, errors.New("x509: trailing data after X.509 subject")
	}

	out.Issuer.FillFromRDNSequence(&issuer)
	out.Subject.FillFromRDNSequence(&subject)

	out.NotBefore = in.TBSCertificate.Validity.NotBefore
	out.NotAfter = in.TBSCertificate.Validity.NotAfter

	for _, e := range in.TBSCertificate.Extensions {
		out.Extensions = append(out.Extensions, e)
		unhandled := false

		if len(e.ID) == 4 && e.ID[0] == 2 && e.ID[1] == 5 && e.ID[2] == 29 {
			switch e.ID[3] {
			case 15:
				// RFC 5280, 4.2.1.3
				var usageBits asn1.BitString
				if rest, err = asn1.Unmarshal(e.Value, &usageBits); err != nil {
					return nil, err
				} else if len(rest) != 0 {
					return nil, errors.New("x509: trailing data after X.509 KeyUsage")
				}

				var usage int
				for i := 0; i < 9; i++ {
					if usageBits.At(i) != 0 {
						usage |= 1 << uint(i)
					}
				}
				out.KeyUsage = KeyUsage(usage)

			case 19:
				// RFC 5280, 4.2.1.9
				var constraints basicConstraints
				if rest, err = asn1.Unmarshal(e.Value, &constraints); err != nil {
					return nil, err
				} else if len(rest) != 0 {
					return nil, errors.New("x509: trailing data after X.509 BasicConstraints")
				}

				out.BasicConstraintsValid = true
				out.IsCA = constraints.IsCA
				out.MaxPathLen = constraints.MaxPathLen
				out.MaxPathLenZero = out.MaxPathLen == 0
				// TODO: map out.MaxPathLen to 0 if it has the -1 default value? (Issue 19285)
			case 17:
				out.DNSNames, out.EmailAddresses, out.IPAddresses, out.URIs, err = parseSANExtension(e.Value)
				if err != nil {
					return nil, err
				}

				if len(out.DNSNames) == 0 && len(out.EmailAddresses) == 0 && len(out.IPAddresses) == 0 && len(out.URIs) == 0 {
					// If we didn't parse anything then we do the critical check, below.
					unhandled = true
				}

			case 30:
				unhandled, err = parseNameConstraintsExtension(out, e)
				if err != nil {
					return nil, err
				}

			case 31:
				// RFC 5280, 4.2.1.13

				// CRLDistributionPoints ::= SEQUENCE SIZE (1..MAX) OF DistributionPoint
				//
				// DistributionPoint ::= SEQUENCE {
				//     distributionPoint       [0]     DistributionPointName OPTIONAL,
				//     reasons                 [1]     ReasonFlags OPTIONAL,
				//     cRLIssuer               [2]     GeneralNames OPTIONAL }
				//
				// DistributionPointName ::= CHOICE {
				//     fullName                [0]     GeneralNames,
				//     nameRelativeToCRLIssuer [1]     RelativeDistinguishedName }

				var cdp []distributionPoint
				if rest, err := asn1.Unmarshal(e.Value, &cdp); err != nil {
					return nil, err
				} else if len(rest) != 0 {
					return nil, errors.New("x509: trailing data after X.509 CRL distribution point")
				}

				for _, dp := range cdp {
					// Per RFC 5280, 4.2.1.13, one of distributionPoint or cRLIssuer may be empty.
					if len(dp.DistributionPoint.FullName) == 0 {
						continue
					}

					for _, fullName := range dp.DistributionPoint.FullName {
						if fullName.Tag == 6 {
							out.CRLDistributionPoints = append(out.CRLDistributionPoints, string(fullName.Bytes))
						}
					}
				}

			case 35:
				// RFC 5280, 4.2.1.1
				var a authKeyID
				if rest, err := asn1.Unmarshal(e.Value, &a); err != nil {
					return nil, err
				} else if len(rest) != 0 {
					return nil, errors.New("x509: trailing data after X.509 authority key-id")
				}
				out.AuthorityKeyID = a.ID

			case 37:
				// RFC 5280, 4.2.1.12.  Extended Key Usage

				// id-ce-extKeyUsage OBJECT IDENTIFIER ::= { id-ce 37 }
				//
				// ExtKeyUsageSyntax ::= SEQUENCE SIZE (1..MAX) OF KeyPurposeId
				//
				// KeyPurposeId ::= OBJECT IDENTIFIER

				var keyUsage []asn1.ObjectIdentifier
				if rest, err := asn1.Unmarshal(e.Value, &keyUsage); err != nil {
					return nil, err
				} else if len(rest) != 0 {
					return nil, errors.New("x509: trailing data after X.509 ExtendedKeyUsage")
				}

				for _, u := range keyUsage {
					if extKeyUsage, ok := extKeyUsageFromOID(u); ok {
						out.ExtKeyUsage = append(out.ExtKeyUsage, extKeyUsage)
					} else {
						out.UnknownExtKeyUsage = append(out.UnknownExtKeyUsage, u)
					}
				}

			case 14:
				// RFC 5280, 4.2.1.2
				var keyid []byte
				if rest, err := asn1.Unmarshal(e.Value, &keyid); err != nil {
					return nil, err
				} else if len(rest) != 0 {
					return nil, errors.New("x509: trailing data after X.509 key-id")
				}
				out.SubjectKeyID = keyid

			case 32:
				// RFC 5280 4.2.1.4: Certificate Policies
				var policies []policyInformation
				if rest, err := asn1.Unmarshal(e.Value, &policies); err != nil {
					return nil, err
				} else if len(rest) != 0 {
					return nil, errors.New("x509: trailing data after X.509 certificate policies")
				}
				out.PolicyIdentifiers = make([]asn1.ObjectIdentifier, len(policies))
				for i, policy := range policies {
					out.PolicyIdentifiers[i] = policy.Policy
				}

			default:
				// Unknown extensions are recorded if critical.
				unhandled = true
			}
		} else if e.ID.Equal(oidExtensionAuthorityInfoAccess) {
			// RFC 5280 4.2.2.1: Authority Information Access
			var aia []authorityInfoAccess
			if rest, err := asn1.Unmarshal(e.Value, &aia); err != nil {
				return nil, err
			} else if len(rest) != 0 {
				return nil, errors.New("x509: trailing data after X.509 authority information")
			}

			for _, v := range aia {
				// GeneralName: uniformResourceIdentifier [6] IA5String
				if v.Location.Tag != 6 {
					continue
				}
				if v.Method.Equal(oidAuthorityInfoAccessOcsp) {
					out.OCSPServer = append(out.OCSPServer, string(v.Location.Bytes))
				} else if v.Method.Equal(oidAuthorityInfoAccessIssuers) {
					out.IssuingCertificateURL = append(out.IssuingCertificateURL, string(v.Location.Bytes))
				}
			}
		} else {
			// Unknown extensions are recorded if critical.
			unhandled = true
		}

		if e.Critical && unhandled {
			out.UnhandledCriticalExtensions = append(out.UnhandledCriticalExtensions, e.ID)
		}
	}

	return out, nil
}

func subjectBytesWhenMarshal(cert *Certificate) ([]byte, error) {
	if len(cert.RawSubject) > 0 {
		return cert.RawSubject, nil
	}
	var ret pkix.RDNSequence
	for _, atv := range cert.Subject.Names {
		ret = append(ret, []pkix.AttributeTypeAndValue{atv})
	}
	return asn1.Marshal(ret)
}

func issuerBytesWhenMarshal(cert *Certificate) ([]byte, error) {
	if len(cert.RawIssuer) > 0 {
		return cert.RawIssuer, nil
	}
	var ret pkix.RDNSequence
	for _, atv := range cert.Issuer.Names {
		ret = append(ret, []pkix.AttributeTypeAndValue{atv})
	}

	return asn1.Marshal(ret)
}

// signParamsForPublicKey returns the parameters to use SignatureAlgorithm
// If requestedSigAlgo is not zero then it overrides the default
// signature algorithm.
func signParamsForPublicKey(requestedSigAlgo SignatureAlgorithm) (hashFunc Hash, sigAlgo pkix.AlgorithmIdentifier, err error) {
	var pubType PublicKeyAlgorithm

	switch requestedSigAlgo {
	case
		MD2WithRSA,
		MD5WithRSA,
		SHA1WithRSA,
		SHA256WithRSA,
		SHA384WithRSA,
		SHA512WithRSA:

		pubType = RSA
		hashFunc = SHA256
		sigAlgo.Algorithm = oidSignatureSHA256WithRSA
		sigAlgo.Parameters = asn1.NullRawValue

	case
		ECDSAWithSHA1,
		ECDSAWithSHA256,
		ECDSAWithSHA384,
		ECDSAWithSHA512:

		pubType = ECDSA

		switch requestedSigAlgo {
		case ECDSAWithSHA256:
			hashFunc = SHA256
			sigAlgo.Algorithm = oidSignatureECDSAWithSHA256
		case ECDSAWithSHA384:
			hashFunc = SHA384
			sigAlgo.Algorithm = oidSignatureECDSAWithSHA384
		case ECDSAWithSHA512:
			hashFunc = SHA512
			sigAlgo.Algorithm = oidSignatureECDSAWithSHA512
		default:
			err = errors.New("x509: unknown ecdsa sign algo")
		}
	case
		SM3WithSM2,
		SHA1WithSM2,
		SHA256WithSM2,
		SHA512WithSM2,
		SHA224WithSM2,
		SHA384WithSM2,
		RMD160WithSM2:

		pubType = SM2

		switch requestedSigAlgo {
		case SM3WithSM2:
			hashFunc = SM3WithPublicKey
			sigAlgo.Algorithm = oidSignatureSM3WithSM2
		case SHA1WithSM2:
			hashFunc = SHA1
			sigAlgo.Algorithm = oidSignatureSHA1WithSM2
		case SHA256WithSM2:
			hashFunc = SHA256
			sigAlgo.Algorithm = oidSignatureSHA256WithRSA
		case SHA512WithSM2:
			hashFunc = SHA512
			sigAlgo.Algorithm = oidSignatureSHA512WithSM2
		case SHA224WithSM2:
			hashFunc = SHA224
			sigAlgo.Algorithm = oidSignatureSHA224WithSM2
		case SHA384WithSM2:
			hashFunc = SHA384
			sigAlgo.Algorithm = oidSignatureSHA384WithSM2
		case RMD160WithSM2:
			hashFunc = RIPEMD160
			sigAlgo.Algorithm = oidSignatureRMD160WithSM2
		default:
			err = errors.New("x509: unknown SM2 sign algo")
		}
	default:
		err = errors.New("x509: only RSA, ECDSA and SM2 keys supported")
		return
	}

	if err != nil {
		return
	}

	found := false
	for _, details := range signatureAlgorithmDetails {
		if details.algo == requestedSigAlgo {
			if details.pubKeyAlgo != pubType {
				err = errors.New("x509: requested SignatureAlgorithm does not match private key type")
				return
			}
			sigAlgo.Algorithm, hashFunc = details.oid, details.hash
			if hashFunc == 0 {
				err = errors.New("x509: cannot sign with hash function requested")
				return
			}
			//if requestedSigAlgo.isRSAPSS() {
			//	sigAlgo.Parameters = rsaPSSParameters(std.Hash(hashFunc))
			//}
			found = true
			break
		}
	}

	if !found {
		err = errors.New("x509: unknown SignatureAlgorithm")
	}

	return
}

// These structures reflect the ASN.1 structure of X.509 certificates.:
type certificate struct {
	Raw                asn1.RawContent
	TBSCertificate     tbsCertificate
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

type tbsCertificate struct {
	Raw                asn1.RawContent
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Issuer             asn1.RawValue
	Validity           validity
	Subject            asn1.RawValue
	PublicKey          publicKeyInfo
	UniqueID           asn1.BitString   `asn1:"optional,tag:1"`
	SubjectUniqueID    asn1.BitString   `asn1:"optional,tag:2"`
	Extensions         []pkix.Extension `asn1:"optional,explicit,tag:3"`
}
