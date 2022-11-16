package tee

//
//import "C"
//import (
//	"crypto"
//	"crypto/rand"
//	"encoding/asn1"
//	"errors"
//	"github.com/denisbrodbeck/machineid"
//	"github.com/hyperchain/go-crypto-standard/asym"
//	flatoDB "github.com/ultramesh/flato-db-interface"
//	//include c code dir
//	"io"
//	"sync"
//)
//

//
////keyStorageEnclaveImpl implement of enclave for key storage
//type keyStorageEnclaveImpl struct {
//	ID               []byte
//	PublicKey        []byte
//	PrivateKeyCipher []byte
//	Algo             string
//}
//
////GetKeyStoreEnclave get key store by id
//func GetKeyStoreEnclave(path string, id []byte, db flatoDB.FlatoDB, algo string) (KeyStore, error) {
//	err := loadEnclave(path)
//	if err != nil {
//		return nil, err
//	}
//	e := &keyStorageEnclaveImpl{}
//	e.Algo = algo
//	get(id, db, e)
//	return e, err
//}
//
////NewKeyStoreEnclave generate a new enclave
//func NewKeyStoreEnclave(pathFromConfig string, generate bool, importKey []byte, db flatoDB.FlatoDB, algo string) (KeyStore, []byte, error) {
//	lerr := loadEnclave(pathFromConfig)
//	if lerr != nil {
//		return nil, nil, lerr
//	}
//
//	pk, vkCipher, err := generateKey(eid, algo)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	tmp := &keyStorageEnclaveImpl{
//		ID:               make([]byte, 32),
//		PublicKey:        pk,
//		PrivateKeyCipher: vkCipher,
//		Algo:             algo,
//	}
//
//	_, _ = rand.Read(tmp.ID)
//	if !put(db, tmp) {
//		return nil, nil, errors.New("store key to db error")
//	}
//	return tmp, tmp.ID, nil
//}
//
////Close close
//func (e *keyStorageEnclaveImpl) Close() {
//	if eid == 0 {
//		return
//	}
//	close(eid)
//}
//
////IsSGX always true
//func (e *keyStorageEnclaveImpl) IsSGX() bool {
//	return true
//}
//
////Public Get public key
//func (e *keyStorageEnclaveImpl) Public() crypto.PublicKey {
//	return e.PublicKey
//}
//
//func (e *keyStorageEnclaveImpl) Sign(_ io.Reader, digest []byte, _ crypto.SignerOpts) (signature []byte, err error) {
//	if eid == 0 {
//		return nil, errors.New("enclave need init")
//	}
//	signature, err = sign(eid, e.PrivateKeyCipher, digest, e.Algo)
//	if err != nil {
//		return nil, err
//	}
//	if e.Algo == Secp256k1 {
//		k := new(asym.ECDSAPublicKey)
//		_ = k.FromBytes(e.PublicKey, asym.AlgoP256K1)
//		ok, _ := k.Verify(nil, signature, digest)
//		if !ok {
//			signature = signature[:64]
//			signature = append(signature, byte(0x01))
//		}
//	}
//	return signature, nil
//}
//
//func (e *keyStorageEnclaveImpl) Verify(sign, hash []byte) (bool, error) {
//	if eid == 0 {
//		return false, errors.New("enclave need init")
//	}
//	return verify(eid, e.PublicKey, sign, hash, e.Algo)
//}
//
//func getHostID() ([]byte, error) {
//	id, err := machineid.ID()
//	if err != nil {
//		return nil, err
//	}
//	return []byte(id), err
//}
//
//func get(id []byte, db flatoDB.FlatoDB, in interface{}) {
//	hostID, err := getHostID()
//	if err != nil {
//		return
//	}
//	id = append(id, hostID...)
//	bs, err := db.Get(id, flatoDB.DBTYPE_CERT)
//	if err != nil {
//		return
//	}
//	_, err = asn1.Unmarshal(bs, in)
//	if err != nil {
//		return
//	}
//}
//
//func put(db flatoDB.FlatoDB, in *keyStorageEnclaveImpl) bool {
//	hostID, err := getHostID()
//	if err != nil {
//		return false
//	}
//	id := append(in.ID, hostID...)
//	bs, err := asn1.Marshal(*in)
//	if err != nil {
//		return false
//	}
//
//	err = db.Put(id, bs, flatoDB.DBTYPE_CERT)
//	if err != nil {
//		return false
//	}
//	return true
//}
