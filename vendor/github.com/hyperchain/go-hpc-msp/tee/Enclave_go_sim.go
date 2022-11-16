//go:build !sgx
// +build !sgx

package tee

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	inter "github.com/hyperchain/go-crypto-standard"
	"github.com/hyperchain/go-crypto-standard/hash"
)

//NewKeyStorageEnclave is an instance of key storage
func start(_ string) (uint64, error) {
	return 1, nil //没有实际意义,大于0即可
}

func close(eid uint64) {
}

func encrypt(_ uint64, vk []byte) ([]byte, error) {
	return teeEncrypt(vk)
}

func decrypt(_ uint64, in []byte) ([]byte, error) {
	return teeDecrypt(in)
}

//func sign(eid uint64, vkCipher []byte, digest []byte, algo string) ([]byte, error) {
//	kb, err := decrypt(eid, vkCipher)
//	if err != nil {
//		return nil, err
//	}
//	k := new(asym.ECDSAPrivateKey)
//	if algo == Secp256r1 {
//		k = new(asym.ECDSAPrivateKey)
//		err = k.FromBytes(kb, asym.AlgoP256R1)
//	} else {
//		k = new(asym.ECDSAPrivateKey)
//		err = k.FromBytes(kb, asym.AlgoP256K1)
//	}
//	if err != nil {
//		return nil, err
//	}
//
//	return k.Sign(nil, digest, rand.Reader)
//}

//func verify(eid uint64, pk, sign, digest []byte, algo string) (bool, error) {
//	k := new(asym.ECDSAPublicKey)
//	var err error
//	if algo == Secp256r1 {
//		k = new(asym.ECDSAPublicKey)
//		err = k.FromBytes(pk, asym.AlgoP256R1)
//	} else {
//		k = new(asym.ECDSAPublicKey)
//		err = k.FromBytes(pk, asym.AlgoP256K1)
//	}
//	if err != nil {
//		return false, err
//	}
//	return k.Verify(nil, sign, digest)
//}
//
//func generateKey(eid uint64, algo string) ([]byte, []byte, error) {
//	k := new(asym.ECDSAPrivateKey)
//	if algo == Secp256r1 {
//		k, _ = asym.GenerateKey(asym.AlgoP256R1)
//	} else {
//		k, _ = asym.GenerateKey(asym.AlgoP256K1)
//	}
//
//	kb, _ := k.Bytes()
//	pk, _ := k.ECDSAPublicKey.Bytes()
//	vkCipher, err := encrypt(eid, kb)
//	return pk /*65*/, vkCipher, err
//}

func getEnclaveLibSignSo(path string) error {
	return nil
}

type keyRequest struct {
	KeyName           [8]byte
	CPUSvnAndSoOn     [16 + 16 + 32 + 4 + 2]byte //AttributeMask KeyID  MiscMask  ConfigSvn
	Reserved2         []byte
	TextOffSetAndSoOn [4 + 12 + 4 + 12 + 16]byte //Reserved3 PayloadSize Reserved4 PayloadTag
	Payload           []byte
}

var keyRequestPool = sync.Pool{
	New: func() interface{} {
		req := new(keyRequest)
		req.Reserved2 = make([]byte, 434)
		return req
	},
}

func teeEncrypt(secret []byte) ([]byte, error) {
	req := keyRequestPool.Get().(*keyRequest)
	defer func() {
		req.Payload = nil
		keyRequestPool.Put(req)
	}()

	var err error
	req.Payload, err = new(inter.AES).Encrypt(inter.AESKey(getKey()), secret, rand.Reader)
	if err != nil {
		return nil, err
	}

	// 17 <= l <= 32
	l := aesBlockSize - len(secret)%aesBlockSize + aesBlockSize
	binary.LittleEndian.PutUint64(req.KeyName[:], uint64(cipherLength-l))

	con := [][]byte{req.KeyName[:],
		req.CPUSvnAndSoOn[:], req.Reserved2[:434-l], req.TextOffSetAndSoOn[:], req.Payload}

	return bytes.Join(con, nil), nil
}

func teeDecrypt(cipher []byte) ([]byte, error) {
	if len(cipher) < 8 {
		return cipher, nil
	}
	num := binary.LittleEndian.Uint64(cipher[:8])
	c := cipher[num:]
	return new(inter.AES).Decrypt(inter.AESKey(getKey()), c)
}

func getKey() []byte {
	var msg = []byte{0x51, 0x75, 0x6c, 0x69, 0x61, 0x6e, 0x20, 0x54, 0x65, 0x63, 0x68, 0x6e, 0x6f, 0x6c, 0x6f, 0x67,
		0x79, 0x20, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x74,
		0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x20, 0x00, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x20, 0x62, 0x6c,
		0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x20, 0x74, 0x65, 0x61, 0x6d, 0x20, 0x77, 0x69,
		0x20, 0x63, 0x68, 0x61, 0x69, 0x72, 0x6d, 0x61, 0x6e, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65,
		0x20, 0x62, 0x6f, 0x61, 0x72, 0x00, 0x2e, 0x20, 0x54, 0x68, 0x65, 0x20, 0x63, 0x6f, 0x6d, 0x70,
		0x61, 0x6e, 0x79, 0x20, 0x68, 0x61, 0x73, 0x20, 0x61, 0x20, 0x74, 0x65, 0x61, 0x6d, 0x20, 0x6f,
		0x66, 0x20, 0x6e, 0x65, 0x61, 0x72, 0x6c, 0x79, 0x20, 0x32, 0x30, 0x30, 0x20, 0x70, 0x65, 0x6f,
		0x70, 0x00, 0x65, 0x2c, 0x20, 0x39, 0x00, 0x00, 0x00, 0x6f, 0x66, 0x20, 0x77, 0x68, 0x6f, 0x6d,
		0x20, 0x61, 0x72, 0x65, 0x20, 0x74, 0x65, 0x63, 0x68, 0x6e, 0x69, 0x63, 0x69, 0x61, 0x6e, 0x73,
		0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x20, 0x54, 0x68, 0x69, 0x73, 0x20, 0x70, 0x00, 0x61, 0x74, 0x66,
		0x6f, 0x72, 0x6d, 0x20, 0x72, 0x61, 0x6e, 0x6b, 0x73, 0x20, 0x66, 0x69, 0x72, 0x73, 0x74, 0x20,
		0x69, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x74, 0x65, 0x63, 0x68, 0x6e, 0x69, 0x63, 0x61, 0x6c}
	hasher := hash.NewHasher(hash.KECCAK_256)
	key, _ := hasher.Hash(msg)
	return key
}

func getHTTPSResponse(eid uint64, method, url string, body string, head map[string]string) (Response, error) {
	res := Response{
		Status: errStatusCode,
		Body:   "",
		Header: nil,
	}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return res, err
	}
	for k, v := range head {
		req.Header.Add(k, v)
	}

	trustCA := x509.NewCertPool()
	resp, err := client(trustCA).Do(req)
	if err != nil {
		return res, err
	}

	var html string
	if resp.StatusCode == 200 {
		robots, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return res, err
		}
		//nolint
		resp.Body.Close()
		if len(robots) > bodyMaxSize {
			return res, errors.New("the size of the request body exceeds 1mb")
		}
		html = string(robots)
	} else {
		html = ""
	}
	reponseHeader := make(map[string]string)
	for k, v := range resp.Header {
		var tmpValue string
		for _, value := range v {
			tmpValue += value
		}
		reponseHeader[k] = tmpValue
	}
	//nolint
	defer resp.Body.Close()
	res.Header = reponseHeader
	res.Status = resp.StatusCode
	res.Body = html

	return res, nil
}

func client(certPool *x509.CertPool) *http.Client {
	var tr *http.Transport
	tr = &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, time.Second*10) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			err = conn.SetDeadline(time.Now().Add(time.Second * 10)) //设置发送接受数据超时
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
		ResponseHeaderTimeout: time.Second * 10,
		TLSClientConfig: &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{Transport: tr}
}
