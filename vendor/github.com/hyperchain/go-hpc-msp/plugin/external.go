package plugin

import (
	"fmt"
	"path/filepath"
	ld "plugin"
	"strings"

	"github.com/hyperchain/go-hpc-common/bvmcom"
	hc "github.com/hyperchain/go-hpc-common/cryptocom"
	msp "github.com/hyperchain/go-hpc-msp"
	"github.com/hyperchain/go-hpc-msp/config"
	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/meshplus/crypto"
)

const getFunctionName = "GetLevel" //func() []crypto.Level

type externalMux struct {
	logger    msp.Logger
	config    config.EncryptionConfigInterface
	soft      hc.Engine
	allPlugin map[string]*externalPlugin //Key: plugin name 用于保证每个plugin只加载一次
	mapping   map[uint64]*functionInfo   //Key: function | mode

	random common.RandomFuncSlot
	//FakeHash(0x00)
	//SHA1(0x10)
	//SHA2_224(0x21),SHA2_256(0x20),SHA2_384(0x22),SHA2_512(0x23)
	//SHA3_224(0x31),SHA3_256(0x30),SHA3_384(0x32),SHA3_512(0x33)
	//KECCAK_224(0x41),KECCAK_256(0x40),KECCAK_384(0x42),KECCAK_512(0x43)
	//SM3(0x50),Sm3WithPublicKey(0x60)
	//SelfDefined(0x70)
	hash [8]common.HashFuncSlot

	//AES_CBC(0x120000),AES_ECB(0x220000),AES_GCM(0x320000)
	//3DES_CBC(0x130000),3DES_ECB(0x230000),3DES_GCM(0x330000)
	//SM4_CBC(0x010000)
	//TEE(0x040000)
	//SelfDefined(0x000000)
	crypt [8]common.CryptFuncSlot

	//Sm2p256v1(0x0100), Secp256k1(0x0200), Secp256r1(0x0300), Secp256k1Recover(0x0600), ED25519(0x2000)
	//SelfDefined(0x0700)
	verify     [6]common.VerifyFuncSlot
	sign       common.SignGetFuncSlot
	signCreate common.SignCreateFuncSlot

	certificate common.ParseCertificateSlot
	ca          common.ParseCASlot
	issue       common.IssueSlot
	generate    common.GenerateLocalCASlot

	skInit  common.KeyAgreementInitSlot
	skFinal common.KeyAgreementFinalSlot
	acc     [6]common.VerifyAccSlot
}

func (mux *externalMux) SetRandomFunc(in common.RandomFuncSlot) {
	mux.random = in
}

func (mux *externalMux) SetHashFunc(algo int, in common.HashFuncSlot) error {
	index, err := hashAlgoMap2Index(algo)
	if err != nil {
		return fmt.Errorf("set hash func error: %v", err)
	}
	mux.hash[index] = in
	return nil
}

func (mux *externalMux) SetCryptFunc(algo int, in common.CryptFuncSlot) error {
	index, err := cryptoAlgoMap2Index(algo)
	if err != nil {
		return fmt.Errorf("set crypto func error: %v", err)
	}
	mux.crypt[index] = in
	return nil
}

func (mux *externalMux) SetVerifyFunc(algo int, in common.VerifyFuncSlot) error {
	index, err := signAlgoMap2Index(algo)
	if err != nil {
		return fmt.Errorf("set verify func error: %v", err)
	}
	mux.verify[index] = in
	return nil
}

func (mux *externalMux) SetACCVerifyFunc(algo int, in common.VerifyAccSlot) error {
	index, err := signAlgoMap2Index(algo)
	if err != nil {
		return fmt.Errorf("set verify func error: %v", err)
	}
	mux.acc[index] = in
	return nil
}

func (mux *externalMux) SetCertFunc(in crypto.PluginSignFuncL2, caMode bvmcom.CAMode) {
	mux.certificate = in.ParseCertificate
	mux.ca = in.ParseAllCA
	mux.signCreate = in.CreateSignKey
}

func (mux *externalMux) SetDistributedCertFunc(in crypto.PluginSignFuncL3) {
	mux.issue = in.Issue
	mux.generate = in.GenerateLocalCA
}

func (mux *externalMux) SetSessionKeyFunc(in1 common.KeyAgreementInitSlot, in3 common.KeyAgreementFinalSlot) {
	mux.skInit, mux.skFinal = in1, in3
}

//func (mux *externalMux) SetEncFunc(mode uint64, in common.EncKeyFuncSlot) {
//	mux.enc[mode] = in
//}
//
//func (mux *externalMux) SetCreateDecKeyFunc(in common.DecCreateFuncSlot) {
//	mux.decCreate = in
//}
//
//func (mux *externalMux) SetDecFunc(in common.DecGetFuncSlot) {
//	mux.dec = in
//}

func (mux *externalMux) loadPlugin(path string) (_ *externalPlugin, err error) {
	defer func() {
		r := recover()
		if r != nil {
			mux.logger.Errorf("load %v error: %v", path, r)
			err = fmt.Errorf("load plugin '%v' error: %v", path, r)
		}
	}()

	p, err := ld.Open(path)
	if err != nil {
		return nil, err
	}
	ext := &externalPlugin{
		Plugin: *p,
		path:   path,
		name:   filepath.Base(path),
	}

	//获取符号getLevelArrayFuncName
	sym1, err := ext.Lookup(getFunctionName)
	if err != nil {
		return nil, err
	}
	getFunction, ok := sym1.(func() []crypto.Level)
	if !ok {
		return nil, fmt.Errorf("type assert 'func() []crypto.Level' error")
	}

	ext.function = getFunction()
	mux.allPlugin[ext.path] = ext
	mux.logger.Noticef("load plugin '%v' success: %v", path, ext.String())
	return ext, nil
}

func (mux *externalMux) load(pluginPath string, slotPosition common.Function, logger crypto.Logger) (err error) {
	const errNotImpl = "plugin %v don't implement %v function interface, please check ns_static.config"
	if pluginPath == msp.Default || len(pluginPath) == 0 {
		return nil
	}

	//检查so是否加载过，获取到externalPlugin
	plugin, ok := mux.allPlugin[pluginPath]
	if !ok {
		plugin, err = mux.loadPlugin(pluginPath)
		if err != nil {
			return fmt.Errorf("load plugin '%v' err: %v", pluginPath, err.Error())
		}
	}

	//检查是否有需要的方法
	var implemented bool
	switch slotPosition {
	case common.Random:
		var r crypto.PluginRandomFunc
		for _, i := range plugin.function {
			r, implemented = i.(crypto.PluginRandomFunc)
			if implemented {
				break
			}
		}
		if r == nil {
			return fmt.Errorf(errNotImpl, plugin.GetName(), slotPosition)
		}
		if err = testRandom(pluginPath, r); err != nil {
			return err
		}
		mux.SetRandomFunc(r.Rander)
	case common.Hash:
		var r crypto.PluginHashFunc
		for _, i := range plugin.function {
			r, implemented = i.(crypto.PluginHashFunc)
			if implemented {
				break
			}
		}
		if r == nil {
			return fmt.Errorf(errNotImpl, plugin.GetName(), slotPosition)
		}

		ls, _ := r.GetLevel()
		for _, algo := range ls {
			if err = testHash(pluginPath, r, algo, mux.soft); err != nil {
				return err
			}
			if err = mux.SetHashFunc(algo, r.GetHash); err != nil {
				return fmt.Errorf("load %v error: %v", plugin.GetName(), err)
			}
		}
	case common.Crypt:
		var r crypto.PluginCryptFunc
		for _, i := range plugin.function {
			r, implemented = i.(crypto.PluginCryptFunc)
			if implemented {
				break
			}
		}
		if r == nil {
			return fmt.Errorf(errNotImpl, plugin.GetName(), slotPosition)
		}

		ls, _ := r.GetLevel()
		for _, algo := range ls {
			if err = testCrypt(pluginPath, r, algo, mux.soft); err != nil {
				return err
			}
			if err = mux.SetCryptFunc(algo, r.GetSecretKey); err != nil {
				return fmt.Errorf("load %v error: %v", plugin.GetName(), err)
			}
		}
	case common.Signature:
		var r crypto.PluginSignFuncL0
		var level int
		for _, i := range plugin.function {
			r, implemented = i.(crypto.PluginSignFuncL3)
			if implemented {
				level = 3
				break
			}
			r, implemented = i.(crypto.PluginSignFuncL2)
			if implemented {
				level = 2
				break
			}
			r, implemented = i.(crypto.PluginSignFuncL1)
			if implemented {
				level = 1
				break
			}
			r, implemented = i.(crypto.PluginSignFuncL0)
			if implemented {
				level = 0
				break
			}
		}
		if r == nil {
			return fmt.Errorf(errNotImpl, plugin.GetName(), slotPosition)
		}

		if mux.config.CAMode() == bvmcom.Distributed && level < 3 {
			logger.Errorf("the certificate plugin under distributed CA must implement crypto.PluginSignFuncL3")
			return fmt.Errorf(errNotImpl, plugin.GetName(), slotPosition)
		}

		switch level {
		case 0:
			if err = mux.LoadL0(r, pluginPath); err != nil {
				return err
			}
		case 1:
			if err = mux.LoadL1(r.(crypto.PluginSignFuncL1), pluginPath); err != nil {
				return err
			}
		case 2:
			if err = mux.LoadL2(r.(crypto.PluginSignFuncL2), pluginPath); err != nil {
				return err
			}
		default:
			if err = mux.LoadL3(r.(crypto.PluginSignFuncL3), pluginPath); err != nil {
				return err
			}
		}

	case common.SessionKey:
		var r crypto.PluginGenerateSessionKeyFunc
		for _, i := range plugin.function {
			r, implemented = i.(crypto.PluginGenerateSessionKeyFunc)
			if implemented {
				break
			}
		}
		if r == nil {
			return fmt.Errorf(errNotImpl, plugin.GetName(), slotPosition)
		}
		if err = testSessionKey(pluginPath, r); err != nil {
			return err
		}
		mux.SetSessionKeyFunc(r.KeyAgreementInit, r.KeyAgreementFinal)
	case common.Accelerate:
		var r hc.PluginAccelerateFunc
		for _, i := range plugin.function {
			r, implemented = i.(hc.PluginAccelerateFunc)
			if implemented {
				break
			}
		}
		if r == nil {
			return fmt.Errorf(errNotImpl, plugin.GetName(), slotPosition)
		}
		ls, _ := r.GetLevel()
		for _, algo := range ls {
			if err := testACCVerify(algo, pluginPath, r); err != nil {
				return fmt.Errorf("load %v error: %v", pluginPath, err)
			}

			if err := mux.SetACCVerifyFunc(algo, r.GetVerifyACCKey); err != nil {
				return fmt.Errorf("load %v error: %v", pluginPath, err)
			}
		}
		return nil
	}
	logger.Noticef("load %v from %v success", slotPosition, pluginPath)
	return
}

func (mux *externalMux) LoadL0(r crypto.PluginSignFuncL0, pluginPath string) error {
	ls, _ := r.GetLevel()
	for _, algo := range ls {
		if err := testVerify(algo, pluginPath, r); err != nil {
			return fmt.Errorf("load %v error: %v", pluginPath, err)
		}

		if err := mux.SetVerifyFunc(algo, r.GetVerifyKey); err != nil {
			return fmt.Errorf("load %v error: %v", pluginPath, err)
		}
	}
	return nil
}

func (mux *externalMux) LoadL1(r crypto.PluginSignFuncL1, pluginPath string) error {
	if err := testSign(pluginPath, r); err != nil {
		return err
	}
	mux.sign = r.GetSignKey
	return mux.LoadL0(r, pluginPath)
}

func (mux *externalMux) LoadL2(r crypto.PluginSignFuncL2, pluginPath string) error {
	if err := testCert(pluginPath, r); err != nil {
		return err
	}
	mux.SetCertFunc(r, mux.config.CAMode())
	return mux.LoadL1(r, pluginPath)
}

func (mux *externalMux) LoadL3(r crypto.PluginSignFuncL3, pluginPath string) error {
	if err := testDistributedCert(pluginPath, r); err != nil {
		return err
	}
	mux.SetDistributedCertFunc(r)
	return mux.LoadL2(r, pluginPath)
}

//externalPlugin externalMux plugin
type externalPlugin struct {
	ld.Plugin
	function   []crypto.Level
	path, name string
}

func (e *externalPlugin) GetName() string {
	return e.name
}

func (e *externalPlugin) GetPath() string {
	return e.name
}

func (e *externalPlugin) GetLevelArray() []crypto.Level {
	return e.function
}

type functionInfo struct {
	from     *externalPlugin
	index    int
	software bool
}

func (e *externalPlugin) String() string {
	var builder = make([]string, 0, 6)
	for _, i := range e.function {
		if _, ok := i.(crypto.PluginRandomFunc); ok {
			builder = append(builder, "[random]")
		}
		if _, ok := i.(crypto.PluginHashFunc); ok {
			builder = append(builder, "[hash]")
		}
		if _, ok := i.(crypto.PluginCryptFunc); ok {
			builder = append(builder, "[crypto]")
		}

		if _, ok := i.(crypto.PluginSignFuncL3); ok {
			builder = append(builder, "[signature(L3)]")
		} else if _, ok := i.(crypto.PluginSignFuncL2); ok {
			builder = append(builder, "[signature(L2)]")
		} else if _, ok := i.(crypto.PluginSignFuncL1); ok {
			builder = append(builder, "[signature(L1)]")
		} else if _, ok := i.(crypto.PluginSignFuncL0); ok {
			builder = append(builder, "[signature(L0)]")
		}

		if _, ok := i.(crypto.PluginGenerateSessionKeyFunc); ok {
			builder = append(builder, "[session key]")
		}

		if _, ok := i.(hc.PluginAccelerateFunc); ok {
			builder = append(builder, "[accelerate]")
		}
	}

	return strings.Join(builder, ", ")
}
