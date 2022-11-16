package software

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/hyperchain/go-hpc-msp/pemencode"
	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/hyperchain/go-hpc-msp/plugin/pkix"
	"github.com/meshplus/crypto"
)

type fileInfo struct {
	derEncode     []byte //cert: der; privateKey: pub der
	alreadyPaired bool
	parsedSK      crypto.SignKey //for privateKey
	parsedCA      *Certificate   //for ca
}

func (s *Engine) bytes2CA(caInfo string, caPairing *sync.Map, tempStoreSk *sync.Map, tempStoreCA *sync.Map) error {
	ca, privateKey, caName, err := s.pairingCA(caInfo, tempStoreSk, tempStoreCA)
	if err != nil {
		return err
	}

	if ca == nil || privateKey == nil {
		return nil
	}

	//key: ca.String()
	caPairing.Store(caName, &X509CA{
		ca:         ca.parsedCA,
		hostname:   pkix.GetIdentityNameFromString(caName).CN,
		ski:        hex.EncodeToString(ca.parsedCA.SubjectKeyID),
		privateKey: ca.parsedSK,
	})
	return nil
}

//caInfo: input
//caPairing, tempStoreCA: output
func (s *Engine) pairingCA(caInfo string, tempStoreVk *sync.Map, tempStoreCA *sync.Map) (
	casOut *fileInfo, privateKeyOut *fileInfo, caNameOut string, err error) {

	b, t := pemencode.PEM2DER([]byte(caInfo))
	switch t {
	case pemencode.PLUGINEncodeKey:
		if err != nil {
			return nil, nil, "", fmt.Errorf("parse private key for ca in plugin encodeing error: %v: %v", err, caInfo)
		}
		sk, err := s.GetSignKey(caInfo)
		if err != nil {
			return nil, nil, "", fmt.Errorf("get signKey for ca in plugin encodeing error: %v: %v", err, caInfo)
		}
		//判断算法类型
		if common.ModeIsRSAAlgo(sk.GetKeyInfo()) {
			return nil, nil, "", fmt.Errorf("unsupport RSA private key, ignore this key: %v", caInfo)
		}

		publicKey := hex.EncodeToString(sk.Bytes())
		newFileInfo := &fileInfo{
			derEncode: b,
			parsedSK:  sk,
		}
		tempStoreVk.Store(publicKey, newFileInfo)

		//try to pair
		val, ok := tempStoreCA.Load(publicKey)
		if !ok { //配对失败
			return nil, nil, "", nil
		}
		//不允许存在一个私钥多个CA的情况，因此删除
		tempStoreCA.Delete(publicKey)

		//配对成功，dears为配对的证书, 证书的alreadyPaired也置true
		dears := val.(*fileInfo)
		newFileInfo.alreadyPaired = true
		newFileInfo.parsedCA = dears.parsedCA
		dears.alreadyPaired = true
		dears.parsedSK = sk
		name := &pkix.IdentityName{
			CN:           dears.parsedCA.Subject.CommonName,
			SerialNumber: hex.EncodeToString(dears.parsedCA.SubjectKeyID),
		}
		caNameOut = name.String()
		casOut, privateKeyOut = dears, newFileInfo
	case pemencode.PEMECCPrivateKey, pemencode.PEMAnyPrivateKey, pemencode.PEMRSAPrivateKey: //maybe it's ecc key
		sk, gerr := s.GetSignKey(caInfo)
		if gerr != nil {
			return nil, nil, "", fmt.Errorf("parse private key error: %v: %v", gerr, caInfo)
		}

		//判断算法类型
		if common.ModeIsRSAAlgo(sk.GetKeyInfo()) {
			return nil, nil, "", fmt.Errorf("unsupport RSA private key, ignore this key: %v", caInfo)
		}

		publicKey := hex.EncodeToString(sk.Bytes())
		newFileInfo := &fileInfo{
			derEncode: b,
			parsedSK:  sk,
		}
		tempStoreVk.Store(publicKey, newFileInfo)

		//try to pair
		val, ok := tempStoreCA.Load(publicKey)
		if !ok { //配对失败
			return nil, nil, "", nil
		}
		//不允许存在一个私钥多个CA的情况, 因此删除掉已经配对的CA
		tempStoreCA.Delete(publicKey)

		//配对成功，dears为配对的证书（可能多个）, 证书的alreadyPaired也置true
		dears := val.(*fileInfo)
		newFileInfo.alreadyPaired = true
		newFileInfo.parsedCA = dears.parsedCA
		dears.alreadyPaired = true
		dears.parsedSK = sk
		name := &pkix.IdentityName{
			CN:           dears.parsedCA.Subject.CommonName,
			SerialNumber: hex.EncodeToString(dears.parsedCA.SubjectKeyID),
		}
		caNameOut = name.String()
		casOut, privateKeyOut = dears, newFileInfo
	default: //maybe it's ca
		c, err := ParseCertificate(caInfo)
		if err != nil {
			return nil, nil, "", fmt.Errorf("parse ca error: %v: %v", err, caInfo)
		}

		err = checkCertRoughly(c, true)
		if err != nil {
			return nil, nil, "", fmt.Errorf("ca check error: %v: %v", err, caInfo)
		}

		//暂存证书
		publicKeyIndex := hex.EncodeToString(c.PublicKey.Bytes())
		newFileInfo := &fileInfo{
			derEncode: b,
			parsedCA:  c,
		}
		tempStoreCA.Store(publicKeyIndex, newFileInfo)

		//try to pair
		val, ok := tempStoreVk.Load(publicKeyIndex)
		if !ok { //配对失败
			return nil, nil, "", nil
		}
		//不允许存在一个私钥多个CA的情况
		tempStoreCA.Delete(publicKeyIndex)

		//配对成功, dear为配对的私钥, 私钥的alreadyPaired也置true
		dear := val.(*fileInfo)
		dear.alreadyPaired = true
		newFileInfo.alreadyPaired = true
		dear.parsedCA = newFileInfo.parsedCA
		newFileInfo.parsedSK = dear.parsedSK
		casOut, privateKeyOut = newFileInfo, dear
		caNameOut = (&pkix.IdentityName{
			CN:           newFileInfo.parsedCA.Subject.CommonName,
			SerialNumber: hex.EncodeToString(dear.parsedCA.SubjectKeyID),
		}).String()
	case pemencode.PEMPublicKey:
		return nil, nil, "", fmt.Errorf("parse failed, this file is nether ecc private key nor certificate: %v", caInfo)
	}
	return casOut, privateKeyOut, caNameOut, nil
}
