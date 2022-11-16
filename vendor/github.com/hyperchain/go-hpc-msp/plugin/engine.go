package plugin

import (
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"github.com/hyperchain/go-hpc-common/bvmcom"
	hc "github.com/hyperchain/go-hpc-common/cryptocom"
	"github.com/hyperchain/go-hpc-common/types"
	msp "github.com/hyperchain/go-hpc-msp"
	"github.com/hyperchain/go-hpc-msp/config"
	"github.com/hyperchain/go-hpc-msp/logger"
	"github.com/hyperchain/go-hpc-msp/pemencode"
	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/hyperchain/go-hpc-msp/plugin/software"
	"github.com/meshplus/crypto"
)

//GetCryptoEngine get crypto engine
func GetCryptoEngine(con config.EncryptionConfigInterface, logger msp.Logger) (hc.Engine, error) {
	// 1.init software engine
	engine := GetSoftwareEngine(logger).(*encryptEngine)
	engine.logger = logger

	//分布式CA中才会允许调用issue和generate
	if con.CAMode() != bvmcom.Distributed {
		pa, pb := getPanicFunc(con.CAMode())
		engine.externalMux.issue = pa
		engine.externalMux.generate = pb
	}
	// 2.init externalMux
	engine.externalMux.allPlugin = make(map[string]*externalPlugin)
	engine.externalMux.mapping = make(map[uint64]*functionInfo)
	//engine.enc = make(map[uint64]common.EncKeyFuncSlot)

	if err := engine.loadExternal(con, logger); err != nil {
		return nil, fmt.Errorf("init externalMux crypto engine error: %v", err)
	}
	return engine, nil
}

//GetSoftwareEngine get software engine
func GetSoftwareEngine(log ...msp.Logger) hc.Engine {
	soft := software.GetSoftwareEngine()
	var logSelector msp.Logger = logger.MSPFakeLogger
	if len(log) > 0 {
		logSelector = log[0]
	}
	ret := &encryptEngine{
		externalMux:        externalMux{soft: soft, logger: logSelector},
		logger:             logSelector,
		defaultHashAlgo:    crypto.KECCAK_256,
		defaultEncryptAlgo: crypto.Aes | crypto.CBC,
	}
	ret.registerSoftware(soft)
	ret.externalMux.soft = soft
	return ret
}

type unLoadError struct {
	code         errLoad
	otherReason  string
	pluginName   string
	slotPosition common.Function
}

type errLoad int

const (
	loadErrConfigReason errLoad = iota
	loadErrAlgoImplement
)

func (e *unLoadError) Error() string {
	switch e.code {
	case loadErrConfigReason:
		//config reason
		return fmt.Sprintf("config reason, %s from %s, detail: %s", e.slotPosition, e.pluginName, e.otherReason)
	case loadErrAlgoImplement:
		//algo implement error
		return fmt.Sprintf("implement err, %s from %s, detail: %s", e.slotPosition, e.pluginName, e.otherReason)
	default:
		//return unexpect error
		return fmt.Sprintf("unexpect error, %s from %s, detail: %v", e.slotPosition, e.pluginName, e.otherReason)
	}
}

//encryptEngine encryption mux
type encryptEngine struct {
	externalMux
	defaultHashAlgo    uint32
	defaultEncryptAlgo uint32
	logger             msp.Logger
}

func (engine *encryptEngine) GetVerifyACCKey(mode int) (hc.VerifyACCKey, error) {
	index, err := signAlgoMap2Index(mode)
	if err != nil {
		return nil, fmt.Errorf("parse mode error: %v", err)
	}
	f := engine.acc[index]
	if f == nil {
		return nil, crypto.ErrNotSupport
	}
	return f(mode)
}

func (engine *encryptEngine) ParseAllCA(s []string) ([]crypto.CA, error) {
	return engine.ca(s)
}

func (engine *encryptEngine) Issue(ca crypto.CA, hostname string, ct crypto.CertType, ext map[string]string, vk crypto.VerifyKey) ([]byte, error) {
	return engine.issue(ca, hostname, ct, ext, vk)
}

func (engine *encryptEngine) GenerateLocalCA(hostName string) (skIndex string, ca crypto.CA, err error) {
	return engine.generate(hostName)
}

func (engine *encryptEngine) ParseCertificate(s string) (crypto.Cert, error) {
	return engine.certificate(s)
}

func (engine *encryptEngine) registerSoftware(s hc.Engine) {
	engine.externalMux.random = s.Rander
	for i := range engine.externalMux.hash {
		engine.externalMux.hash[i] = s.GetHash
	}
	for i := range engine.externalMux.crypt {
		engine.externalMux.crypt[i] = s.GetSecretKey
	}
	for i := range engine.externalMux.verify {
		engine.externalMux.verify[i] = s.GetVerifyKey
	}
	engine.externalMux.sign = s.GetSignKey
	engine.externalMux.signCreate = s.CreateSignKey
	engine.externalMux.certificate = s.ParseCertificate
	engine.externalMux.ca = s.ParseAllCA
	engine.externalMux.issue = s.Issue
	engine.externalMux.generate = s.GenerateLocalCA
	engine.externalMux.skInit = s.KeyAgreementInit
	engine.externalMux.skFinal = s.KeyAgreementFinal

	//将外部的Engine注入到software中，这是为了在软件密钥协商后能够用到插件的对称加密
	se := s.(*software.Engine)
	se.Outer = engine
}

func (engine *encryptEngine) loadExternal(config config.EncryptionConfigInterface, logger msp.Logger) error {
	logger.Noticef("start load externalMux crypto engine")
	engine.externalMux.config = config
	//1.random
	if err := engine.externalMux.load(config.EncryptionEngineRand(), common.Random, engine.logger); err != nil {
		return err
	}
	//2.hash
	if err := engine.externalMux.load(config.EncryptionEngineHash(), common.Hash, engine.logger); err != nil {
		return err
	}
	//3.crypt
	if err := engine.externalMux.load(config.EncryptionEngineCrypt(), common.Crypt, engine.logger); err != nil {
		return err
	}
	//4.signature
	if err := engine.externalMux.load(config.EncryptionEngineSign(), common.Signature, engine.logger); err != nil {
		return err
	}
	//6.session key
	if err := engine.externalMux.load(config.EncryptionEngineSessionKey(), common.SessionKey, engine.logger); err != nil {
		return err
	}
	//7.accelerate
	if err := engine.externalMux.load(config.EncryptionEngineACC(), common.Accelerate, engine.logger); err != nil {
		return err
	}

	return nil
}

//GetLevel get implemented algorithm
func (engine *encryptEngine) GetLevel() ([]int, uint8) {
	panic("do not use")
}

//Rander Random reader
func (engine *encryptEngine) Rander() (io.Reader, error) {
	return engine.random()
}

//GetHash get Hash function
//如果没有指定SelfDefinedHash，则默认自定义算法为KECCAK_256
func (engine *encryptEngine) GetHash(mode int) (crypto.Hasher, error) {

	if mode == types.DefaultAlgo {
		mode = int(atomic.LoadUint32(&engine.defaultHashAlgo))
	}

	index, err := hashAlgoMap2Index(mode)
	if err != nil {
		return nil, fmt.Errorf("parse mode error: %v", err)
	}
	return engine.hash[index](mode)
}

//GetSecretKey get secret Key
func (engine *encryptEngine) GetSecretKey(mode int, pwd, key []byte) (crypto.SecretKey, error) {

	if mode == types.DefaultAlgo {
		mode = int(atomic.LoadUint32(&engine.defaultEncryptAlgo))
	}

	index, err := cryptoAlgoMap2Index(mode)
	if err != nil {
		return nil, fmt.Errorf("parse mode error: %v", err)
	}
	return engine.crypt[index](mode, pwd, key)
}

//GetVerifyKey get Signature Key
func (engine *encryptEngine) GetVerifyKey(key []byte, mode int) (crypto.VerifyKey, error) {
	if mode == crypto.None {
		//1.尝试软件
		vk, err := engine.soft.GetVerifyKey(key, mode)
		if err == nil {
			return vk, nil
		}
		//2.依次尝试
		for _, f := range engine.verify {
			vk, err = f(key, mode)
			if err == nil {
				return vk, err
			}
		}
		return nil, fmt.Errorf("parse vk without algo info error: %v", hex.EncodeToString(key))
	}

	index, err := signAlgoMap2Index(mode)
	if err != nil {
		return nil, fmt.Errorf("parse mode error: %v", err)
	}
	return engine.verify[index](key, mode)
}

//GetSignKey get sign Key
func (engine *encryptEngine) GetSignKey(keyIndex string) (crypto.SignKey, error) {
	_, t := pemencode.PEM2DER([]byte(keyIndex))
	if strings.HasPrefix(keyIndex, "plugin") {
		keyIndex = strings.TrimSpace(keyIndex[len("plugin"):])
	}
	switch t {
	case pemencode.PEMPublicKey, pemencode.PEMCertificate:
		return nil, fmt.Errorf("input is not privateKey index")
	default:
		//pemencode.PEMECCPrivateKey, pemencode.PEMAnyPrivateKey, pemencode.PLUGINEncodeKey, pemencode.PEMRSAPrivateKey
		sk, err := engine.sign(keyIndex)
		if err != nil {
			engine.logger.Debugf("parse sk by plugin error, try to use software: %v", err)
			return engine.soft.GetSignKey(keyIndex)
		}
		return sk, nil
	}
}

//CreateSignKey sign Key
func (engine *encryptEngine) CreateSignKey() (index string, k crypto.SignKey, err error) {
	return engine.signCreate()
}

////GetEncKey get enc Key
//func (engine *encryptEngine) GetEncKey(key []byte, mode int) (crypto.EncKey, error) {
//	if mode == types.DefaultAlgo {
//		mode = engine.da.encryptAlgo
//	}
//	if mode != crypto.None {
//		return engine.EncKey(key, mode)
//	}
//
//	if raw, modeInner, err := software.ParsePKIXPublicKey(key); err == nil {
//		return engine.EncKey(raw, modeInner)
//	}
//
//	return nil, crypto.ErrNotSupport
//}
//
////GetDecKey get dec Key
//func (engine *encryptEngine) GetDecKey(key string) (crypto.DecKey, error) {
//	return engine.DecGet(key)
//}
//
////CreateDecKey create dec Key
//func (engine *encryptEngine) CreateDecKey() (index string, k crypto.DecKey, err error) {
//	return engine.DecCreate()
//}

//SetAlgo engine set default algo
func (engine *encryptEngine) SetAlgo(data []byte) error {
	engine.logger.Debugf("[msp encryptEngine] start to set algo")

	as := types.DefaultAlgoSet
	if len(data) != 0 {
		bs, err := bvmcom.UnmarshalAlgoMap(data)
		if err != nil {
			return err
		}
		engine.logger.Debugf("[msp encryptEngine] input data is not nil [%v]", string(data))
		as = bs.GetLatestAlgo()
	}

	engine.defaultHashAlgo = uint32(types.GetAlgoInt(as.HashAlgo))
	engine.defaultEncryptAlgo = uint32(types.GetAlgoInt(as.EncryptAlgo))

	engine.logger.Debugf("[msp encryptEngine] set algo success, hashAlgo[%v], encryptAlgo[%v]", as.HashAlgo, as.EncryptAlgo)
	return nil
}

//GetDefaultAlgo engine get default algo
func (engine *encryptEngine) GetDefaultAlgo() (int, int) {
	return int(atomic.LoadUint32(&engine.defaultHashAlgo)),
		int(atomic.LoadUint32(&engine.defaultEncryptAlgo))
}

func (engine *encryptEngine) KeyAgreementInit() (data1, data2ToPeer []byte, err error) {
	return engine.skInit()
}

func (engine *encryptEngine) KeyAgreementFinal(algo string, data1, data2FromPeer []byte) (crypto.SecretKey, error) {
	return engine.skFinal(algo, data1, data2FromPeer)
}
