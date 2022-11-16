package software

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	gm "github.com/hyperchain/go-crypto-gm"
	inter "github.com/hyperchain/go-crypto-standard"
	"github.com/hyperchain/go-crypto-standard/asym"
	"github.com/hyperchain/go-crypto-standard/ed25519"
	"github.com/hyperchain/go-crypto-standard/hash"
	hc "github.com/hyperchain/go-hpc-common/cryptocom"
	"github.com/hyperchain/go-hpc-msp/pemencode"
	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/hyperchain/go-hpc-msp/plugin/pkix"
	"github.com/hyperchain/go-hpc-msp/tee"
	"github.com/meshplus/crypto"
)

//Engine software Engine
type Engine struct {
	Outer hc.Engine
}

//GetLevel get level
func (s *Engine) GetLevel() ([]int, uint8) {
	panic("do not use")
}

//GetVerifyACCKey get verify acc key
func (s *Engine) GetVerifyACCKey(mode int) (hc.VerifyACCKey, error) {
	return nil, crypto.ErrNotSupport
}

//KeyAgreementInit step 1 of key argument
//data1 is random secret, should be stored locally
//data2 will be sent to remote peer
func (s *Engine) KeyAgreementInit() (data1, data2 []byte, err error) {
	key, _ := gm.GenerateSM2Key()
	data1 = key.K[:]
	data2, _ = key.PublicKey.Bytes()
	return
}

//KeyAgreementFinal step 3 of key argument
func (s *Engine) KeyAgreementFinal(algo string, data1, data2FromPeer []byte) (crypto.SecretKey, error) {
	pubkey := new(gm.SM2PublicKey)
	err := pubkey.FromBytes(data2FromPeer, crypto.Sm2p256v1)
	if err != nil {
		return nil, fmt.Errorf("parse date2FromPeer error: %v", hex.EncodeToString(data2FromPeer))
	}
	if len(data1) != 32 {
		return nil, fmt.Errorf("parse data1 error: length is not 32")
	}
	curve := gm.GetSm2Curve()
	x, y := curve.ScalarMult(new(big.Int).SetBytes(pubkey.X[:]), new(big.Int).SetBytes(pubkey.Y[:]), data1)
	z := bytes.Join([][]byte{x.Bytes(), y.Bytes()}, nil)

	var engine hc.Engine = s
	if !reflect.ValueOf(s.Outer).IsNil() {
		engine = s.Outer
	}

	//generate key
	var key crypto.SecretKey
	switch strings.ToLower(strings.TrimSpace(algo)) {
	case "sm4":
		//hasher if sm4, use sm3
		iterator := GetKeyByKDF(z, gm.NewSM3Hasher(), 16)
		mode := crypto.Sm4 | crypto.CBC
		key, err = engine.GetSecretKey(mode, nil, iterator())
	case "3des":
		iterator := GetKeyByKDF(z, hash.NewHasher(hash.SHA2_256), 32)
		mode := crypto.Des3 | crypto.CBC
		key, err = engine.GetSecretKey(mode, nil, iterator())
	case "aes":
		iterator := GetKeyByKDF(z, hash.NewHasher(hash.SHA2_256), 32)
		mode := crypto.Aes | crypto.CBC
		key, err = engine.GetSecretKey(mode, nil, iterator())
	case "selfdefined":
		iterator := GetKeyByKDF(z, hash.NewHasher(hash.SHA2_256), 32)
		key, err = engine.GetSecretKey(crypto.SelfDefinedCrypt, nil, iterator())
	default:
		return nil, fmt.Errorf("unknown symmetric encryption algorithm, please modify config and restart node")
	}
	if err != nil {
		return nil, fmt.Errorf("GetSecretKey error: %v", err)
	}
	return key, nil
}

//GetKeyByKDF returns a hash function
func GetKeyByKDF(z []byte, hasher crypto.Hasher, length int) func() []byte {
	hash, _ := hasher.Hash(z)
	return func() []byte {
		hash, _ = hasher.BatchHash([][]byte{z, hash})
		if len(hash) > length {
			hash = hash[:length]
		}
		return hash
	}
}

//Issue issue cert
func (s *Engine) Issue(ca crypto.CA, hostname string, ct crypto.CertType, ext map[string]string, vk crypto.VerifyKey) ([]byte, error) {
	x509ca, ok := ca.(*X509CA)
	if !ok {
		return nil, fmt.Errorf("[default Engine] ca is not x509ca, but %T", ca)
	}
	if x509ca.privateKey == nil {
		return nil, fmt.Errorf("[default Engine] ca does not have private key (it's remote ca)")
	}
	ret, err := GenCert(x509ca.ca, x509ca.privateKey, vk, pkix.GetOrganization(ext), hostname, ct.String(), false,
		x509ca.ca.NotBefore, x509ca.ca.NotAfter)
	if err != nil {
		return nil, fmt.Errorf("issue cert from [default Engine] error: %v", err)
	}
	return ret, nil
}

//GenerateLocalCA generate local ca
func (s *Engine) GenerateLocalCA(hostName string) (skIndex string, ca crypto.CA, err error) {
	now := time.Now()
	certByte, vkDer, _ := NewSelfSignedCert("{}", hostName, "", CurveTypeK1, now, now.Add(time.Hour*876000))
	pem, _ := pemencode.DER2PEM(certByte, pemencode.PEMCertificate)
	cert, _ := ParseCertificate(string(pem))

	sk, err := s.GetSignKey(string(vkDer))
	if err != nil {
		panic("generate local ca error" + err.Error())
	}
	result := &X509CA{
		privateKey: sk,
		ski:        hex.EncodeToString(cert.SubjectKeyID),
		hostname:   hostName,
		ca:         cert,
	}
	return vkDer, result, nil
}

//GetSoftwareEngine get SoftwareEngine
func GetSoftwareEngine() hc.Engine {
	ret := &Engine{}
	return ret
}

//Rander get rander
func (s *Engine) Rander() (io.Reader, error) {
	return rand.Reader, nil
}

//GetHash get hash
func (s *Engine) GetHash(mode int) (crypto.Hasher, error) {
	switch mode {
	case crypto.SM3:
		return gm.NewSM3Hasher(), nil
	case crypto.Sm3WithPublicKey:
		return gm.NewSM3IDHasher().(crypto.Hasher), nil
	case crypto.FakeHash:
		return hash.GetFakeHasher(), nil
	case crypto.SelfDefinedHash:
		return hash.NewHasher(crypto.KECCAK_256), nil
	default:
		ret := hash.NewHasher(hash.HashType(mode))
		if ret == nil {
			return nil, fmt.Errorf("software Engine: GetHash error: %v", mode)
		}
		return ret, nil
	}
}

//GetSecretKey get secret key
func (s *Engine) GetSecretKey(mode int, pwd, key []byte) (crypto.SecretKey, error) {
	r := &SecretKey{
		key: key,
	}
	switch mode {
	case crypto.Sm4 | crypto.CBC:
		r.c = &gm.SM4{}
	case crypto.Aes | crypto.CBC:
		r.c = &inter.AES{}
	case crypto.Des3 | crypto.CBC:
		r.c = &inter.TripleDES{}
	case crypto.TEE:
		r.c, _ = tee.NewSecretStoreEnclave()
	}
	return r, nil
}

//GetVerifyKey key is raw or rich-text with crypto.None
func (s *Engine) GetVerifyKey(key []byte, mode int) (crypto.VerifyKey, error) {
	if mode == crypto.None {
		key, _ = pemencode.PEM2DER(key)
		raw, mod, err := ParsePKIXPublicKey(key)
		if err != nil {
			return nil, fmt.Errorf("defailt software vk rich-format is PKIX, got %v", hex.EncodeToString(raw))
		}
		key, mode = raw, mod
	}
	return parseDERPublicKey(key, mode)
}

//GetSignKey for software Engine, key index is DER or raw
// for ecdsa\sm2\ed25519, if keyIndex is raw, you mast support a non-None mode
// for rsa, we only accept PKCS1 DER
// OpenSSL 0.9.8 generates PKCS#1 private keys by default
// OpenSSL 1.0.0 generates PKCS#8 keys.
// OpenSSL ecparam generates SEC1 EC private keys for ECDSA
func (s *Engine) GetSignKey(keyIndex string) (crypto.SignKey, error) {
	return getPrivateKey(keyIndex)
}

//CreateSignKey return printable keyIndex
// for software Engine, always use sm2
func (s *Engine) CreateSignKey() (index string, k crypto.SignKey, err error) {
	return createPrivateKey(crypto.Sm2p256v1)
}

//GetEncKey parse enc key
func (s *Engine) GetEncKey(key []byte, mode int) (crypto.EncKey, error) {
	var err error
	if mode == crypto.None {
		key, mode, err = ParsePKIXPublicKey(key)
		if err != nil {
			return nil, fmt.Errorf("parse pkix public Key error")
		}
	}
	return parseDERPublicKey(key, mode)
}

//GetDecKey parse dec key
func (s *Engine) GetDecKey(keyIndex string) (crypto.DecKey, error) {
	return getPrivateKey(keyIndex)
}

func getPrivateKey(keyIndex string) (crypto.PrivateKey, error) {
	der, t := pemencode.PEM2DER([]byte(keyIndex))
	if t == pemencode.PLUGINEncodeKey {
		debug.PrintStack()
		return nil, fmt.Errorf("parse self-defined format key error, software not support")
	}
	if t == pemencode.PEMECCPrivateKey || t == pemencode.PEMAnyPrivateKey || t == pemencode.PEMRSAPrivateKey {
		return parseDERPrivateKey(der)
	}
	return nil, fmt.Errorf("unknown keyIndex type: %v", keyIndex)
}

//CreateDecKey create dec key
func (s *Engine) CreateDecKey() (index string, k crypto.DecKey, err error) {
	return createPrivateKey(crypto.Sm2p256v1)
}

//ParseCertificate in is pem or der
func (s *Engine) ParseCertificate(inStr string) (crypto.Cert, error) {
	c, err := ParseCertificate(inStr)
	if err != nil {
		return nil, fmt.Errorf("parse certificate error: %v", err)
	}

	//判断算法类型
	if common.ModeIsRSAAlgo(c.PublicKey.GetKeyInfo()) {
		return nil, fmt.Errorf("unsupport RSA certificate")
	}

	//判断是否过期
	n := time.Now()
	if n.Before(c.NotBefore) || n.After(c.NotAfter) {
		return nil, fmt.Errorf("cert(CN=%v) is expired, from: %v to: %v, ignore this cert", c.Subject.CommonName, c.NotBefore, c.NotAfter)
	}

	return &X509Cert{cert: c}, nil
}

//ParseAllCA parse all ca
func (s *Engine) ParseAllCA(idStr []string) ([]crypto.CA, error) {
	var caSet sync.Map                                 //caName,string => ca, crypto.CA
	var tempPrivateKeySet, tempCertificateSet sync.Map //publicStr,string => file, fileInfo

	for _, info := range idStr {
		rerr := s.bytes2CA(info, &caSet, &tempPrivateKeySet, &tempCertificateSet)
		if rerr != nil {
			//d.GetLogger().Warningf("parse file %v failed, ignore it, reason: %v", info.Path, rerr)
			continue
		}
	}

	ret := make([]crypto.CA, 0, 2)
	caSet.Range(func(caName, caInfo interface{}) bool {
		ret = append(ret, caInfo.(crypto.CA))
		return true
	})

	tempCertificateSet.Range(func(caName, caInfo interface{}) bool {
		info := caInfo.(*fileInfo).parsedCA
		ret = append(ret, &X509CA{
			ca:       info,
			hostname: info.Subject.CommonName,
			ski:      hex.EncodeToString(info.SubjectKeyID),
		})
		return true
	})

	return ret, nil
}

//SetAlgo set algo
func (s *Engine) SetAlgo(data []byte) error {
	panic("do not use")
}

//GetDefaultAlgo get default algo
func (s *Engine) GetDefaultAlgo() (int, int) {
	panic("do not use")
}

func parseDERPublicKey(pkixInner []byte, mode int) (crypto.PublicKey, error) {
	var inner crypto.Verifier
	switch {
	case common.ModeIsECDSAAlgo(mode):
		tmp := new(asym.ECDSAPublicKey)
		err := tmp.FromBytes(pkixInner, mode)
		if err != nil {
			return nil, err
		}
		inner = tmp
	case common.ModeIsRSAAlgo(mode):
		tmp := new(asym.RSAPublicKey)
		err := tmp.FromBytes(pkixInner, mode) //asn1(N, E)
		if err != nil {
			return nil, err
		}
		if modeFromTmp, err := ModeFromRSAMod(tmp.N.BitLen()); err != nil || modeFromTmp != mode {
			return nil, fmt.Errorf("mode and Key size do not match, %v", modeFromTmp)
		}
		inner = tmp
	case mode == crypto.Sm2p256v1:
		tmp := new(gm.SM2PublicKey)
		if err := tmp.FromBytes(pkixInner, crypto.Sm2p256v1); err != nil {
			return nil, err
		}
		inner = tmp
	case mode == crypto.Ed25519:
		tmp := new(ed25519.EDDSAPublicKey)
		if err := tmp.FromBytes(pkixInner, crypto.Ed25519); err != nil {
			return nil, err
		}
		inner = tmp
	default:
		return nil, fmt.Errorf("unknown mode")
	}
	return &PublicKey{
		Mode: mode,
		Key:  inner,
	}, nil
}

//CreateSoftwarePrivateKey create software private key
func CreateSoftwarePrivateKey(mode int) (string, crypto.PrivateKey, error) {
	return createPrivateKey(mode)
}

func createPrivateKey(mode int) (string, crypto.PrivateKey, error) {
	var sk crypto.Signer
	var vk crypto.Verifier
	switch {
	case common.ModeIsECDSAAlgo(mode):
		tmp, _ := asym.GenerateKey(mode)
		sk, vk = tmp, &tmp.ECDSAPublicKey
	case common.ModeIsRSAAlgo(mode):
		rsaMod, _ := ModeGetRSAMod(mode)
		tmp, _ := rsa.GenerateKey(rand.Reader, rsaMod)
		sk, vk = (*asym.RSAPrivateKey)(tmp), (*asym.RSAPublicKey)(&tmp.PublicKey)
	case mode == crypto.Sm2p256v1:
		tmp, _ := gm.GenerateSM2Key()
		sk, vk = tmp, &tmp.PublicKey
	case mode == crypto.Ed25519:
		tmp, _ := ed25519.GenerateKey(rand.Reader)
		sk, vk = tmp, tmp.Public().(*ed25519.EDDSAPublicKey)
	default:
		return "", nil, fmt.Errorf("createPrivateKey unknown mode")
	}
	k := &PrivateKey{
		PublicKey: PublicKey{
			Mode: mode,
			Key:  vk,
		},
		PrivKey: sk,
	}
	der, err := MarshalPKCS8PrivateKey(k)
	if err != nil {
		return "", nil, err
	}

	//to pem
	buf := bytes.NewBuffer(nil)
	err = pem.Encode(buf, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	})
	if err != nil {
		return "", nil, err
	}
	KeyIndex := buf.Bytes()

	//note: 持久化私钥需要另外调用persistPriv，有两处地方使用。
	// 1.分布式ca新节点生成defaultID使用
	// 2.DummyCA生成defaultID使用
	//另外分布式CA给自己生成ca（persistCA）中也调用persistPriv
	return string(KeyIndex), k, err
}

//MarshalPrivateKeyPKCS8 marshal private key with PKCS8
func MarshalPrivateKeyPKCS8(key []byte, mode int) (index []byte, err error) {
	var sk crypto.Signer
	var vk crypto.Verifier
	switch {
	case common.ModeIsECDSAAlgo(mode):
		tmp := new(asym.ECDSAPrivateKey)
		if err := tmp.FromBytes(key, mode); err != nil {
			return nil, err
		}
		tmp.CalculatePublicKey()
		sk, vk = tmp, &tmp.ECDSAPublicKey
	case common.ModeIsRSAAlgo(mode):
		tmp := new(asym.RSAPrivateKey)
		if err := tmp.FromBytes(key, mode); err != nil {
			return nil, err
		}
		sk, vk = tmp, (*asym.RSAPublicKey)(&tmp.PublicKey)
	case mode == crypto.Sm2p256v1:
		tmp := new(gm.SM2PrivateKey)
		if err := tmp.FromBytes(key, mode); err != nil {
			return nil, err
		}
		tmp.CalculatePublicKey()
		sk, vk = tmp, &tmp.PublicKey
	default:
		return nil, fmt.Errorf("unknown mode")
	}
	tmp := &PrivateKey{
		PublicKey: PublicKey{
			Mode: mode,
			Key:  vk,
		},
		PrivKey: sk,
	}
	return MarshalPKCS8PrivateKey(tmp)
}

//parseDERPrivateKey parse private Key in pkcs8, sec1, pkcs1, sm2
func parseDERPrivateKey(der []byte) (crypto.PrivateKey, error) {
	if k, err := ParsePKCS8PrivateKey(der); err == nil {
		return k, nil
	}
	if k, err := ParseSMPrivateKey(der); err == nil {
		return &PrivateKey{
			PublicKey: PublicKey{
				Mode: crypto.Sm2p256v1,
				Key:  &k.PublicKey,
			},
			PrivKey: k,
		}, nil
	}
	if k, err := ParseECPrivateKey(der); err == nil {
		return &PrivateKey{
			PublicKey: PublicKey{
				Mode: common.ModeFromCurve(k.Curve),
				Key:  &k.ECDSAPublicKey,
			},
			PrivKey: k,
		}, nil
	}
	if k, err := ParsePKCS1PrivateKey(der); err == nil {
		mode, _ := ModeFromRSAMod(k.N.BitLen())
		return &PrivateKey{
			PublicKey: PublicKey{
				Mode: mode,
				Key:  (*asym.RSAPublicKey)(&k.PublicKey),
			},
			PrivKey: (*asym.RSAPrivateKey)(k),
		}, nil
	}
	return nil, crypto.ErrNotSupport
}

func checkCertRoughly(certificate *Certificate, isCA bool) error {
	if isCA && !certificate.IsCA {
		return fmt.Errorf("need a ca")
	}
	if time.Now().Before(certificate.NotBefore) ||
		time.Now().After(certificate.NotAfter) {
		return fmt.Errorf("time is non-compliant")
	}
	return nil
}
