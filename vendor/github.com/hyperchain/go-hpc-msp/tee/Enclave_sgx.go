//go:build sgx
// +build sgx

package tee

/*
#cgo CFLAGS : -I../include -I../proxy -I/opt/intel/sgxsdk/include -w -DDEBUG
#cgo linux LDFLAGS: -L/opt/intel/sgxsdk/lib64 -L/opt/intel/mbedtls -L/usr/local -lmbedtls_SGX_u -lsgx_urts_sim -lsgx_uae_service_sim -lsgx_tservice_sim
//linux下的仿真模式，只需要sgxsdk
//运行硬件模式的前提是需要安装sgx的相关环境：driver，psw，sdk
//linux下的硬件模式需要重新生成硬件模式下的enlave.sign.so，另外需要将此处链接的动态库改为sgx_urts, sgx_uae_service

//StoreKey   SignSecp256r1    GenerateSecp256r1    VerifySecp256r1    UnsealData
//initialize_enclave    sgx_destroy_enclave    wrap
#include "KeyStore.h"
#include "KeyStore.c"
#include "Enclave_u.h"
#include "Enclave_u.c"
#include "util.h"
#include "util.c"
*/
import "C"
import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

//NewKeyStorageEnclave is an instance of key storage
func start(absPath string) (eid uint64, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("start enclave error: %v", e)
		}
	}()
	cStr := C.CString(absPath)
	defer C.free(unsafe.Pointer(cStr))
	eid = uint64(C.initialize_enclave(cStr)) //Initialize the enclave
	if eid == 0 {
		return 0, fmt.Errorf("init enclave error")
	}
	var in, out C.flato_sgx_parameters
	C.wrap(-1, C.ulong(eid), in, C.NULL, 0, out, C.NULL, 0)
	return eid, nil
}

func getEnclaveLibSignSo(path string) error {
	// from local path
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return nil
}

func close(eid uint64) {
	//C.sgx_destroy_enclave(C.ulong(eid))
}

func encrypt(eid uint64, vk []byte) ([]byte, error) {
	vklen := len(vk)
	var cipher = make([]byte, vklen+cipherLength)
	isok := C.StoreKey(C.ulong(eid), (*C.uchar)(unsafe.Pointer(&vk[0])), (C.uint)(vklen), (*C.uint8_t)(unsafe.Pointer(&cipher[0])))
	if isok != 1 {
		return nil, errors.New("store key failed")
	}
	return cipher[:], nil
}

func decrypt(eid uint64, cipher []byte) ([]byte, error) {
	vklen := len(cipher) - cipherLength
	var vk = make([]byte, vklen)
	isok := C.UnsealData(C.ulong(eid), (*C.uchar)(unsafe.Pointer(&cipher[0])), (C.uint)(vklen), (*C.uint8_t)(unsafe.Pointer(&vk[0])))
	if isok != 1 {
		return nil, errors.New("unseal data failed")
	}
	return vk[:], nil
}

func sign(eid uint64, vkCipher []byte, digest []byte, algo string) ([]byte, error) {
	var sign = make([]byte, 64)

	isok := C.Sign(C.ulong(eid), (*C.uchar)(unsafe.Pointer(&sign[0])), (*C.uchar)(unsafe.Pointer(&vkCipher[0])), (*C.uchar)(unsafe.Pointer(&digest[0])), C.CString(algo))
	if isok != 1 {
		return nil, errors.New("sign failed")
	}
	if algo == Secp256k1 {
		sign = append(sign, byte(0x00))
	}
	return sign, nil
}

func verify(eid uint64, pk, sign, digest []byte, algo string) (bool, error) {
	pk = pk[1:]
	sign = sign[:64]
	isVerify := C.Verify(C.ulong(eid), (*C.uchar)(unsafe.Pointer(&pk[0])), (*C.uchar)(unsafe.Pointer(&sign[0])), (*C.uchar)(unsafe.Pointer(&digest[0])), C.CString(algo))
	if isVerify != 1 {
		return false, errors.New("verify failed")
	}
	return true, nil
}

func generateKey(eid uint64, algo string) ([]byte, []byte, error) {
	var pk = make([]byte, 64)
	var vkCipher = make([]byte, cipherLength+32)
	isok := C.GenerateKey(C.ulong(eid), (*C.uint8_t)(unsafe.Pointer(&pk[0])), (*C.uint8_t)(unsafe.Pointer(&vkCipher[0])), C.CString(algo))
	if isok != 1 {
		return nil, nil, errors.New("generate key failed")
	}
	if algo == Secp256k1 {
		pk = append([]byte{0x04}, pk...)
	}
	return pk, vkCipher, nil
}

func getHTTPSResponse(eid uint64, method, url string, body string, head map[string]string) (Response, error) {
	//分隔发送来的请求链接如https://www.baidu.com/s?wd=ip，得到[https:  www.baidu.com s?wd=ip], s[2]既是host， www.baidu.com
	res := Response{
		Status: errStatusCode,
		Body:   "",
		Header: nil,
	}
	s := strings.Split(url, "/")
	if len(s) < 3 {
		return res, errors.New("parse url failed")
	}
	//如果得到的host后面带有端口，则解析出此端口，默认https端口为443
	index := strings.Index(s[2], ":")
	var port = "443"
	url = s[2]
	if index > 0 {
		port = s[2][index+1:]
		url = s[2][:index]
	}
	//分隔后的字符串组host后面的内容即为请求资源，因此这里将host后面的内容进行拼接，空则默认为"/"
	var source string
	if len(s) > 3 {
		for index = 3; index < len(s); index++ {
			source += "/" + s[index]
		}
	} else {
		source = "/"
	}
	//拼接请求头 请求体，请求头下面需要加一空行再加请求体
	request := method + " " + source + " HTTP/1.1\r\n"
	for k, v := range head {
		request += k + ": " + v + "\r\n"
	}
	request += "\r\n"
	if len(body) != 0 {
		request += string(body)
	}

	var status [4]byte
	header := make([]byte, bodyMaxSize)
	responseBody := make([]byte, bodyMaxSize)
	var bodyLength [4]byte
	isok := C.Get_https_reponse(C.ulong(eid), C.CString(url), C.CString(request), C.CString(port), (*C.uchar)(unsafe.Pointer(&status[0])), (*C.uchar)(unsafe.Pointer(&header[0])), (*C.uchar)(unsafe.Pointer(&responseBody[0])), (*C.uchar)(unsafe.Pointer(&bodyLength[0])))
	if isok != 0 {
		return res, errors.New("Get_https_reponse failed")
	}

	statusCode := binary.LittleEndian.Uint32(status[:])
	blen := binary.LittleEndian.Uint32(bodyLength[:])

	//解析header字符串 形成一个map[string]string结构，c++里定义的字符串是通过k和v之间用":"分开，每一条之间通过"\n"分开
	reponseHeader := make(map[string]string)
	temp := strings.Split(string(header), "\n")
	for index := 0; index < len(temp); index++ {
		s := strings.Index(temp[index], ":")
		if s == -1 {
			continue
		}
		reponseHeader[temp[index][:s]] = temp[index][s+1:]
	}
	res.Status = int(statusCode)
	res.Header = reponseHeader
	res.Body = string(responseBody)[:blen]

	return res, nil
}
