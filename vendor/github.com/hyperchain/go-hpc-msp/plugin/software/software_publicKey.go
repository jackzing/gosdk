package software

import (
	"crypto/rsa"
	"encoding/asn1"
	"fmt"
	"hash"
	"io"

	gm "github.com/hyperchain/go-crypto-gm"
	"github.com/hyperchain/go-crypto-standard/asym"
	"github.com/hyperchain/go-crypto-standard/ed25519"
	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/meshplus/crypto"
)

//PublicKey public Key
type PublicKey struct {
	Mode int
	//gm, ecdsa or rsa
	Key crypto.Verifier
}

//GetKeyInfo get Key information
func (p *PublicKey) GetKeyInfo() int {
	return p.Mode
}

//Verify Verify signature
func (p *PublicKey) Verify(msg []byte, hasher hash.Hash, sig []byte) bool {
	if p.Key == nil {
		return false
	}

	if p.Mode == crypto.Ed25519 {
		b, _ := p.Key.Verify(nil, sig, msg)
		return b
	}

	if len(sig) < 65 {
		return false
	}

	hasher.Reset()
	//PKCS1v15
	var hashTypeUsedInPKCS1v15 []byte
	if common.ModeIsRSAAlgo(p.Mode) {
		if len(msg) < 4 {
			return false
		}
		hashTypeUsedInPKCS1v15 = msg[len(msg)-4:]
		msg = msg[:len(msg)-4]
	}
	//批量的签名
	if p.Mode == crypto.Sm2p256v1 && sig[0] != 0x30 && sig[1] == 0x30 {
		sig = sig[1:]
	}
	//write public Key
	if _, ok := hasher.(*gm.IDHasher); ok && p.Mode == crypto.Sm2p256v1 {
		_, _ = hasher.Write(p.Bytes())
	}
	_, _ = hasher.Write(msg)
	digst := hasher.Sum(nil)
	b, err := p.Key.Verify(hashTypeUsedInPKCS1v15, sig, digst)
	return b && err == nil
}

//Encrypt encrypt
func (p *PublicKey) Encrypt(data []byte, reader io.Reader) ([]byte, error) {
	if p.Key == nil {
		return nil, errUninitialized
	}
	switch p.Mode {
	case crypto.Sm2p256v1:
		key := p.Key.(*gm.SM2PublicKey)
		return gm.Encrypt(key, data, reader)
	case crypto.Secp256k1, crypto.Secp256r1, crypto.Secp384r1, crypto.Secp521r1, crypto.Secp256k1Recover, crypto.Ed25519:
		return nil, errNotSupport
	case crypto.Rsa2048, crypto.Rsa3072, crypto.Rsa4096:
		return rsa.EncryptPKCS1v15(reader, (*rsa.PublicKey)(p.Key.(*asym.RSAPublicKey)), data)
	default:
		return nil, fmt.Errorf("this Key hasn't Init")
	}
}

//RichBytes return der
func (p *PublicKey) RichBytes() (ret []byte) {
	ret, _ = MarshalPublicKey(p)
	return
}

//Bytes return raw key
func (p *PublicKey) Bytes() (ret []byte) {
	if p.Key == nil {
		return nil
	}
	switch {
	case p.Mode == crypto.Sm2p256v1:
		k, ok := p.Key.(*gm.SM2PublicKey)
		if !ok {
			return nil
		}
		ret, _ = k.Bytes()
	case common.ModeIsECDSAAlgo(p.Mode):
		k, ok := p.Key.(*asym.ECDSAPublicKey)
		if !ok {
			return nil
		}
		ret, _ = k.Bytes()
	case common.ModeIsRSAAlgo(p.Mode):
		k, ok := p.Key.(*asym.RSAPublicKey)
		if !ok {
			return nil
		}
		ret, _ = asn1.Marshal(*k)
	case p.Mode == crypto.Ed25519:
		k, ok := p.Key.(*ed25519.EDDSAPublicKey)
		if !ok {
			return nil
		}
		ret, _ = k.Bytes()
	default:
		return nil
	}
	return ret
}
