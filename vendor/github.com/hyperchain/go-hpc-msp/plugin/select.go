package plugin

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/go-crypto-standard/hash"
	"github.com/hyperchain/go-hpc-common/bvmcom"
	hc "github.com/hyperchain/go-hpc-common/cryptocom"
	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/meshplus/crypto"
)

func fail(err error, name string, f common.Function) error {
	return &unLoadError{
		code:         loadErrAlgoImplement,
		pluginName:   name,
		slotPosition: f,
		otherReason:  err.Error(),
	}
}

//Random插件选择和算法无关
//和算法无关是指，要么都采用要么都不采用
//和算法耦合是指部分算法采用该实现，部分算法不采用
//算法耦合意味着需要做软件后补，并且算法必须正确实现
func testRandom(name string, obj crypto.PluginRandomFunc) error {
	var tmp [2048]byte
	r, err := obj.Rander()
	if err != nil {
		return fail(err, name, common.Random)
	}
	_, err = r.Read(tmp[:])
	if err != nil {
		return fail(err, name, common.Random)
	}
	return nil
}

//0 -> FakeHash(0x00)
//1 -> SHA1(0x10)
//2 -> SHA2_256(0x20)
//3 -> SHA3_256(0x30)
//4 -> KECCAK_256(0x40)
//5 -> SM3(0x50)
//6 -> Sm3WithPublicKey(0x60)
//7 -> SelfDefinedHash(0x70)
// 0<i<8
var hashAlgo = [8]int{crypto.FakeHash, crypto.SHA1, crypto.SHA2_256, crypto.SHA3_256, crypto.KECCAK_256, crypto.SM3, crypto.Sm3WithPublicKey, crypto.SelfDefinedHash}

//nolint
func hashIndexMap2Algo(i int) int {
	return hashAlgo[i]
}

func hashAlgoMap2Index(algo int) (int, error) {
	t := algo & 0xff >> 4
	if t > 7 {
		return 0, fmt.Errorf("hash algo error: %v", algo)
	}
	return t, nil
}

//0 -> SM4_CBC(0x110000)
//1 -> AES_CBC(0x120000)
//2 -> 3DES_CBC(0x130000)
//3 -> TEE(0x040000)
//4 -> SelfDefinedCrypt(0x050000)
// 0<i<5
var cryptoAlgo = [5]int{crypto.Sm4 | crypto.CBC, crypto.Aes | crypto.CBC, crypto.Des3 | crypto.CBC, crypto.TEE, crypto.SelfDefinedCrypt}

//nolint
func cryptoIndexMap2Algo(i int) int {
	return cryptoAlgo[i]
}

func cryptoAlgoMap2Index(algo int) (int, error) {
	t := algo&0x0f0000>>crypto.Symmetrical - 1
	if t > 4 {
		return 0, fmt.Errorf("crypto algo error: %v", algo)
	}
	return t, nil
}

//0 -> Sm2p256v1(0x0100)
//1 -> Secp256k1(0x0200)
//2 -> Secp256r1(0x0300)
//3 -> Secp256k1Recover(0x0600)
//4 -> Ed25519(0x2000)
//5 -> SelfDefined(0x0700)
//0<i<6
var signAlgo = [6]int{crypto.Sm2p256v1, crypto.Secp256k1, crypto.Secp256r1, crypto.Secp256k1Recover, crypto.Ed25519, crypto.SelfDefinedSign}

//nolint
func signIndexMap2Algo(i int) int {
	return signAlgo[i]
}

func signAlgoMap2Index(algo int) (int, error) {
	switch algo {
	case crypto.Sm2p256v1:
		return 0, nil
	case crypto.Secp256k1:
		return 1, nil
	case crypto.Secp256r1:
		return 2, nil
	case crypto.Secp256k1Recover:
		return 3, nil
	case crypto.Ed25519:
		return 4, nil
	case crypto.SelfDefinedSign:
		return 5, nil
	default:
		return 0, fmt.Errorf("sign algo error: %v", algo)
	}
}

func isValidCryptoAlgo(algo int) bool {
	return algo == crypto.Sm4|crypto.CBC ||
		algo == crypto.Aes|crypto.CBC ||
		algo == crypto.Des3|crypto.CBC ||
		algo == crypto.SelfDefinedCrypt
}

func getKey(algo int) []byte {
	var ret []byte
	switch algo {
	case crypto.Sm4 | crypto.CBC:
		ret = make([]byte, 16)
	case crypto.Des3 | crypto.CBC:
		ret = make([]byte, 24)
	default: //crypto.SelfDefinedCrypt，crypto.Aes | crypto.CBC
		ret = make([]byte, 32)
	}
	_, _ = rand.Read(ret)
	return ret
}

func testCrypt(name string, obj crypto.PluginCryptFunc, algo int, soft hc.Engine) error {
	//检查algo是合法的algo
	if !isValidCryptoAlgo(algo) {
		return fail(fmt.Errorf("unknown crypto algo: %v, accept %v", algo,
			"crypto.Sm4|crypto.CBC, crypto.Aes|crypto.CBC, crypto.Des3|crypto.CBC, crypto.SelfDefinedCrypt"), name, common.Crypt)
	}

	//测试加解密算法匹配
	seed := getKey(algo)
	sk, err := obj.GetSecretKey(algo, nil, seed)
	if err != nil {
		return fail(fmt.Errorf("get secretKey with %v error: %v", algo, err), name, common.Crypt)
	}
	cipher := sk.Encrypt(common.Message, rand.Reader)
	plan := sk.Decrypt(cipher)
	if !bytes.Equal(plan, common.Message) {
		return fail(fmt.Errorf("encrypted and decrypted result are't match: %v", algo), name, common.Crypt)
	}

	//对非自定义算法需要测试正确性
	if algo != crypto.SelfDefinedCrypt {
		skStd, err := soft.GetSecretKey(algo, nil, seed)
		if err != nil {
			return fmt.Errorf("get secret key from software when self-test error: %v", err)
		}
		msg := bytes.Repeat(common.Message, 33)
		cipher1 := skStd.Encrypt(msg, rand.Reader)
		ret := sk.Decrypt(cipher1)
		if !bytes.Equal(ret, msg) {
			return fail(fmt.Errorf("decrypt correctness error"), name, common.Crypt)
		}

		cipher2 := sk.Encrypt(msg, rand.Reader)
		ret = skStd.Decrypt(cipher2)
		if !bytes.Equal(ret, msg) {
			return fail(fmt.Errorf("encrypt correctness error"), name, common.Crypt)
		}
	}

	return nil
}

func isValidHashAlgo(algo int) bool {
	return algo == crypto.SHA2_256 ||
		algo == crypto.SHA3_256 ||
		algo == crypto.SM3 ||
		algo == crypto.KECCAK_256 ||
		algo == crypto.SelfDefinedHash
}

//Hash插件选择和算法耦合，因为如地址计算等部分场景要求必须提供某种算法，而用户插件很可能不做实现
//办法就是根据算法选择使用哪种插件，然后使用软件实现候补未实现的算法
//hash必须是32字节结果，自定义hash也是
func testHash(name string, obj crypto.PluginHashFunc, algo int, soft hc.Engine) error {
	//检查algo是合法的algo
	if !isValidHashAlgo(algo) {
		return fail(fmt.Errorf("unknown hash algo: %v, accept %v", algo,
			"crypto.SHA2_256, crypto.SHA3_256, crypto.KECCAK_256, crypto.SM3, crypto.SelfDefinedHash"), name, common.Hash)
	}

	//检查基本功能
	hasher, err := obj.GetHash(algo)
	if err != nil {
		return fail(fmt.Errorf("get hasher with %v error: %v", algo, err), name, common.Hash)
	}

	for j := 0; j < 5; j++ {
		n, wErr := hasher.Write(common.Message)
		if wErr != nil {
			return fail(fmt.Errorf("hasher.Write error: %v", wErr), name, common.Hash)
		}
		if n != len(common.Message) {
			return fail(fmt.Errorf("result 'n' error"), name, common.Hash)
		}
	}

	got := hasher.Sum(nil)

	//检查两次hash的一致性
	hasher2, err := obj.GetHash(algo)
	if err != nil {
		return fail(fmt.Errorf("get hasher with %v error: %v", algo, err), name, common.Hash)
	}
	for i := 0; i < 5; i++ {
		n, wErr := hasher2.Write(common.Message)
		if wErr != nil {
			return fail(fmt.Errorf("hasher.Write error: %v", wErr), name, common.Hash)
		}
		if n != len(common.Message) {
			return fail(fmt.Errorf("result 'n' error"), name, common.Hash)
		}
	}

	got2 := hasher2.Sum(nil)

	if !bytes.Equal(got, got2) {
		return fail(fmt.Errorf("the results were different"), name, common.Hash)
	}

	//检查hash结果需要为32字节
	if len(got) != 32 {
		return fail(fmt.Errorf("hash algo result should be 32 bytes"), name, common.Hash)
	}

	//如果不是自定义算法，需要检查正确性
	if algo != crypto.SelfDefinedHash {
		softHasher, _ := soft.GetHash(algo)
		for i := 0; i < 5; i++ {
			_, _ = softHasher.Write(common.Message)
		}
		expect := softHasher.Sum(nil)
		if !bytes.Equal(got, expect) {
			return fail(fmt.Errorf("correctness error"), name, common.Hash)
		}
	}

	return nil
}

func testSign(name string, obj crypto.PluginSignFuncL1) error {
	return nil
}

func testCert(name string, obj crypto.PluginSignFuncL2) error {
	index, key1, err := obj.CreateSignKey()
	if err != nil {
		return fail(fmt.Errorf("create sign key error when self test: %v", err), name, common.Signature)
	}

	key2, err := obj.GetSignKey(index)
	if err != nil {
		return fail(fmt.Errorf("get sign key error when self test: %v", err), name, common.Signature)
	}

	hasher := hash.NewHasher(hash.KECCAK_256)
	sign, err := key1.Sign(common.Message, hasher, rand.Reader)
	if err != nil {
		return fail(fmt.Errorf("sign error when self test: %v", err), name, common.Signature)
	}
	hasher.Reset()
	if !key2.Verify(common.Message, hasher, sign) {
		return fail(fmt.Errorf("verify error when self test, sign: %v", hex.EncodeToString(sign)), name, common.Signature)
	}

	if !bytes.Equal(key2.Bytes(), key1.Bytes()) {
		return fail(fmt.Errorf("key.Bytes error when self test, sign: %v : %v",
			hex.EncodeToString(key1.Bytes()), hex.EncodeToString(key2.Bytes())), name, common.Signature)
	}

	key3, err := obj.GetVerifyKey(key1.Bytes(), key1.GetKeyInfo())
	if err != nil {
		return fail(fmt.Errorf("get verify key error when self test: %v", err), name, common.Signature)
	}
	hasher.Reset()
	if !key3.Verify(common.Message, hasher, sign) {
		return fail(fmt.Errorf("verify error when self test, sign: %v", hex.EncodeToString(sign)), name, common.Signature)
	}
	return nil
}

//defaultPub 预置公钥
//测试数据通过TestGenerateTestCaseForVerify生成
var defaultPub = map[int]string{
	crypto.Sm2p256v1:        "04a42dbf3feb999a6d14820fa073ab7b9fdfd8025a7223b80a4c9a47696e2256b48e654b8d7072d2fd3d2f9d9c392941df773637a9c35ec4212181d37a8b50d5c1",
	crypto.Ed25519:          "5518cd62f07484af3e898a84d3e9eaf33cd44f80ccd789acad0659cd748ae499",
	crypto.Secp256k1:        "04d92cfcff9a2e394511fda0c1b0df7c97e89c934d6390e656dbcf3c67aa7c03feeff4997f06b077584df78d5a85a5b655205b3691c2191b8c343b2013bd96fdd1",
	crypto.Secp256r1:        "045884efee4c9fe4bb05c7b7097056e38523aa59968710f8c7d558599a3f25e1418b092247fc847aa682e1c49cfa0efc83debbac83e2582c52db8f45af1c91cee8",
	crypto.Secp256k1Recover: "05b9fab0eb57cf695def3aa702cc2fbab8c15d69",
}

//defaultSign 对[]byte("hyperchain")的签名
var defaultSign = map[int]string{
	//sm3 with id
	crypto.Sm2p256v1: "3046022100ea7e32780b08701d219117c82289227c69b9d8edb1d45778d04e318a62eb1f86022100c19cc19b2e592be85495cb80dc1e46eb5c9363491d7e225aae20763980b48d91",
	//fake hash
	crypto.Ed25519: "7a385ca9bec3b969e7101b1eb95e45c4200a7b287b19582a936e116deaa3965f94d65c017467f589da89e9f4cf7b997de38c1143dd8509a8e36b9d8f07edca03",
	//keccak256
	crypto.Secp256k1: "304402206b0f0a77b0cb40a6b77f1a4cd1f5f40aba4534ef6c441748289b6a477209a02502203ad13f961b37a5be2dc184d7c5a94d137daf0432249a48b020921c9a034fdc15",
	//keccak256
	crypto.Secp256r1: "3045022100f238546f698a42a9166fd266694f1a7569cfd0bd27fe2ad00fd1a416018ed23302205a49acf7a39918549a3f942facd623e8695a0bf93e58f038c44f4b214447b53b",
	//keccak256
	crypto.Secp256k1Recover: "ddbf8eae9d6fab78e0e73713484cd85c54af2b7acaf8cc55cfdfe9ea5aac271a3ce756af14ebd6b96ddef526ac222c34233497efa69d1e90e44499102d1879b400",
}

func isValidSignatureAlgo(algo int) bool {
	return algo == crypto.Secp256r1 ||
		algo == crypto.Secp256k1 ||
		algo == crypto.Sm2p256v1 ||
		algo == crypto.Secp256k1Recover ||
		algo == crypto.Ed25519 ||
		algo == crypto.SelfDefinedSign
}

func testACCVerify(algo int, name string, obj hc.PluginAccelerateFunc) error {
	accer, err := obj.GetVerifyACCKey(algo)
	if err != nil {
		return fmt.Errorf("test acc interface for %v error: %v", algo, err)
	}
	if algo == crypto.SelfDefinedSign {
		return fmt.Errorf("acceleration of custom algorithms is not supported")
	}
	key, _ := hex.DecodeString(defaultPub[algo])
	signature, _ := hex.DecodeString(defaultSign[algo])

	switch algo {
	case crypto.Sm2p256v1:
		hasher := gm.NewSM3IDHasher()
		_, _ = hasher.Write(key)
		_, _ = hasher.Write(common.Message)
		h := hasher.Sum(nil)
		err = accer.VerifyACC([][]byte{key}, [][]byte{signature}, [][]byte{h})
		if err != nil {
			return fail(fmt.Errorf("sm2 self-test acc error: %v", err), name, common.Signature)
		}

		key2 := make([]byte, len(key))
		signature2 := make([]byte, len(signature))
		h2 := make([]byte, len(h))
		copy(key2, key)
		copy(signature2, signature)
		copy(h2, h)
		h[12] = 0
		err = accer.VerifyACC([][]byte{key, key2}, [][]byte{signature, signature2}, [][]byte{h, h2})
		if err == nil {
			return fail(fmt.Errorf("sm2 self-test case2 acc error: %v", err), name, common.Signature)
		}

	case crypto.Secp256r1, crypto.Secp256k1, crypto.Secp256k1Recover:
		hasher := hash.NewHasher(crypto.KECCAK_256)
		_, _ = hasher.Write(common.Message)
		h := hasher.Sum(nil)
		err = accer.VerifyACC([][]byte{key}, [][]byte{signature}, [][]byte{h})
		if err != nil {
			return fail(fmt.Errorf("ecdsa self-test acc error: %v", err), name, common.Signature)
		}

		key2 := make([]byte, len(key))
		signature2 := make([]byte, len(signature))
		h2 := make([]byte, len(h))
		copy(key2, key)
		copy(signature2, signature)
		copy(h2, h)
		h[12] = 0
		err = accer.VerifyACC([][]byte{key, key2}, [][]byte{signature, signature2}, [][]byte{h, h2})
		if err == nil {
			return fail(fmt.Errorf("ecdsa self-test acc case2 error: %v", err), name, common.Signature)
		}
	case crypto.Ed25519:
		h := bytes.Repeat(common.Message, 7)
		err = accer.VerifyACC([][]byte{key}, [][]byte{signature}, [][]byte{h})
		if err != nil {
			return fail(fmt.Errorf("ed25519 self-test acc error: %v", err), name, common.Signature)
		}

		key2 := make([]byte, len(key))
		signature2 := make([]byte, len(signature))
		h2 := make([]byte, len(h))
		copy(key2, key)
		copy(signature2, signature)
		copy(h2, h)
		h[12] = 0
		err = accer.VerifyACC([][]byte{key, key2}, [][]byte{signature, signature2}, [][]byte{h, h2})
		if err == nil {
			return fail(fmt.Errorf("ed25519 self-test case2 acc error: %v", err), name, common.Signature)
		}
	}
	return nil
}

//testVerify 验证验签算法的正确性
//如果不是自定义算法，则需要完全和原有的系统兼容
func testVerify(algo int, name string, obj crypto.PluginSignFuncL0) error {
	//1.检查algo是合法的algo
	if !isValidSignatureAlgo(algo) {
		return fail(fmt.Errorf("unknown signature algo: %v, accept %v", algo,
			"crypto.Secp256r1, crypto.Secp256k1, crypto.Sm2p256v1, crypto.Secp256k1Recover, crypto.Ed25519, crypto.SelfDefinedSign"), name, common.Hash)
	}

	if algo == crypto.SelfDefinedSign {
		return nil
	}

	//2.对非自定义的签名算法检验正确性
	key, _ := hex.DecodeString(defaultPub[algo])
	signature, _ := hex.DecodeString(defaultSign[algo])

	//2.1测试能够正确解析裸公钥
	k, err := obj.GetVerifyKey(key, algo)
	if err != nil {
		return fail(fmt.Errorf("get verify key error: %v", err), name, common.Signature)
	}
	//2.2测试能够验签
	switch algo {
	case crypto.Sm2p256v1:
		hasher := gm.NewSM3IDHasher()
		_, _ = hasher.Write(k.Bytes())
		_, _ = hasher.Write(common.Message)
		hashRet := hasher.Sum(nil)
		fake := hash.GetFakeHasher()
		if !k.Verify(hashRet, fake, signature) {
			return fail(fmt.Errorf("sm2 self-test error: verify eror"), name, common.Signature)
		}
		//2.3测试能够识别出错误签名
		signature[37], signature[5], signature[55] = 1, 2, 3
		if k.Verify(hashRet, fake, signature) {
			return fail(fmt.Errorf("sm2 self-test error: case 2 expect err"), name, common.Signature)
		}
		fake.Reset()
		//2.4测试能够接受空签名并返回错误
		if k.Verify(hashRet, fake, nil) {
			return fail(fmt.Errorf("sm2 self-test error: case 3 expect err"), name, common.Signature)
		}
	case crypto.Secp256r1, crypto.Secp256k1, crypto.Secp256k1Recover:
		hasher := hash.NewHasher(crypto.KECCAK_256)
		if !k.Verify(common.Message, hasher, signature) {
			return fail(fmt.Errorf("ecdsa self-test error: verify eror"), name, common.Signature)
		}
		//2.3测试能够识别出错误签名
		signature[37], signature[5], signature[55] = 1, 2, 3
		if k.Verify(common.Message, hasher, signature) {
			return fail(fmt.Errorf("ecdsa self-test error: case 2 expect err"), name, common.Signature)
		}
		hasher.Reset()
		//2.4测试能够接受空签名并返回错误
		if k.Verify(common.Message, hasher, nil) {
			return fail(fmt.Errorf("ecdsa self-test error: case 3 expect err"), name, common.Signature)
		}
	case crypto.Ed25519:
		hasher := hash.GetFakeHasher()
		if !k.Verify(common.Message, hasher, signature) {
			return fail(fmt.Errorf("ed25519 self-test error: verify eror"), name, common.Signature)
		}
		//2.3测试能够识别出错误签名
		signature[37], signature[5] = 1, 2
		if k.Verify(common.Message, hasher, signature) {
			return fail(fmt.Errorf("ed25519 self-test error: case 2 expect err"), name, common.Signature)
		}
		hasher.Reset()
		//2.4测试能够接受空签名并返回错误
		if k.Verify(common.Message, hasher, nil) {
			return fail(fmt.Errorf("ed25519 self-test error: case 3 expect err"), name, common.Signature)
		}
	}

	return nil
}

func testDistributedCert(name string, obj crypto.PluginSignFuncL3) error {
	_, ca, err := obj.GenerateLocalCA("test_node1")
	if err != nil {
		return fail(fmt.Errorf("generate local ca error: %v", err), name, common.Signature)
	}
	_, k, err := obj.CreateSignKey()
	if err != nil {
		return fail(fmt.Errorf("generate sign key error: %v", err), name, common.Signature)
	}

	cert, err := obj.Issue(ca, "test_node2", crypto.ECert, make(map[string]string), k)
	if err != nil {
		return fail(fmt.Errorf("generate cert error: %v", err), name, common.Signature)
	}

	c, err := obj.ParseCertificate(string(cert))
	if err != nil {
		return fail(fmt.Errorf("parse self-generaed cert error: %v", err), name, common.Signature)
	}

	if c.GetHostName() != "test_node2" || c.GetCAHostName() != "test_node1" {
		return fail(fmt.Errorf("get hostname or get ca hostname error: got %v and %v", c.GetCAHostName(), c.GetCAHostName()), name, common.Signature)
	}

	_, err = obj.ParseAllCA([]string{ca.String()})
	if err != nil {
		return fail(fmt.Errorf("parse ca '%v' error: %v", ca.String(), err), name, common.Signature)
	}

	err = c.VerifyCert([]string{ca.String()}, nil)
	if err != nil {
		return fail(fmt.Errorf("verify self generated cert '%v' error: %v", string(cert), err), name, common.Signature)
	}

	return nil
}

func testSessionKey(name string, obj crypto.PluginGenerateSessionKeyFunc) error {
	d1, d2, err := obj.KeyAgreementInit()
	if err != nil {
		return fail(fmt.Errorf("KeyAgreementInit error: %v", err), name, common.SessionKey)
	}

	sk2, err := obj.KeyAgreementFinal("sm4", d1, d2)
	if err != nil {
		return fail(fmt.Errorf("KeyAgreementFinal error: %v", err), name, common.SessionKey)
	}

	cipher := sk2.Encrypt(common.Message, rand.Reader)
	m2 := sk2.Decrypt(cipher)
	if !bytes.Equal(m2, common.Message) {
		return fmt.Errorf("key argument error")
	}
	sk2.Destroy()
	return nil
}

func getPanicFunc(mod bvmcom.CAMode) (common.IssueSlot, common.GenerateLocalCASlot) {
	return func(ca crypto.CA, hostname string, ct crypto.CertType, ext map[string]string, vk crypto.VerifyKey) ([]byte, error) {
			panic(fmt.Sprintf("shod not call this at center CA mode: method 'Issue', mode: %v", mod.String()))
		}, func(hostName string) (skIndex string, ca crypto.CA, err error) {
			panic(fmt.Sprintf("shod not call this at center CA mode: method 'GenerateLocalCA', mode: %v", mod.String()))
		}

}
