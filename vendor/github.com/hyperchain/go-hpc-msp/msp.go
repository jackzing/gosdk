package msp

import (
	"crypto/elliptic"
	"net"

	"github.com/hyperchain/go-hpc-common/bvmcom"
	"github.com/hyperchain/go-hpc-common/cryptocom"
	"github.com/hyperchain/go-hpc-common/subscription"
	"github.com/hyperchain/go-hpc-msp/config"

	"google.golang.org/grpc/credentials"
)

// IdentityManager is used for auth.
type IdentityManager interface {
	Reloader
	GetCryptoEngine() cryptocom.Engine
	GetConfig() config.EncryptionConfigInterface
	GetEventMux() *subscription.TypeMux
	GetLogger() Logger

	/*
		不会返回nil，返回的长度可能是0
		分布式CA：根据hostname和certType匹配证书签发者（Issuer）
		非分布式CA：总是获取到default的证书，返回的长度总是1
	*/
	GetIdentities(name ...string) []Identity //get identity which respect this node (self), params are certType、hostname
	//GetCAs params will be ignored, also used in bvm
	GetCAs(_ string) []CA

	//RequestCertificate 不修改msp和系统的状态, 对默认私钥生成证书请求
	RequestCertificate(o, cn, gn string) ([]byte, error) //generate a csr,only CAF
	//GenerateIdentity 不修改msp和系统的状态
	GenerateIdentity(csr []byte) ([]byte, error) //To generate the ID,only CAF

	//NewIdentity 获取ID实例
	//第三个参数(isPersist)，根据peer的hostname而存储为hostname.cert
	//  distributed中用于新节点持久化申请到的证书（hts.PersistCert）
	//  center中用isPersist来决定是否将对方的证书保存到内存中，用于nvp。
	NewIdentity(identity, privateIndex []byte, isPersist bool) (Identity, error)

	//NewCA 获取CA实例
	NewCA(identity, privateIndex string) (CA, error)
	//RevokeIdentity local revoke identity
	RevokeIdentity(id Identity) error

	CheckAndVerify(cert Identity, msg, sign []byte) error //Verify certificates and signatures
	Sign(alias string, msg []byte) ([]byte, error)        //The signature

	Start() error
	Stop()
	CheckReplace(ch chan interface{}) (bool, error)  //Check replace cert settings
	ReplaceCerts(der []byte, issuer ...string) error //Msp revoke certs during reload
	GetVerifyTube() cryptocom.VerifyTube             //Get verify tube instance

	NewSKGenerator(version int64) SKGenerator

	SetCvpCallback(callback func(string) error)
	//根据hostname删除对应证书，这里应该传的是绑定的VP节点的hostname，这样才能指定删除哪个证书
	RemoveIdentity(name string) error

	ReloadLocal() error

	StopNs()

	CAMode() bvmcom.CAMode

	ReadCAFromLedger()
}

//Identity defines the interface
type Identity interface {
	GetName() string
	//GetCAName get ca **hostname**
	GetCAName() string
	//GetAKI get Authority Key Identifier
	GetAKI() string
	//GetPrivateKeyIndex if Identity have privateKey, return private key Index
	GetPrivateKeyIndex() (bool, string)
	//GetIdentityBytes return cert text
	//if the result is PEM, it is converted to Der
	GetIdentityBytes() []byte
	//Sign If there is no private key, an error is returned
	Sign(msg []byte) ([]byte, error)
	//SignBatch If there is no private key, an error is returned,
	//batch signature is preferred
	SignBatch(msg []byte) ([]byte, error)
	//VerifySignature verify signature
	VerifySignature(sign, msg []byte) error
	//GetAddress Get address from cert
	GetAddress() ([]byte, error)
	//GetKeyInfo return curve、publicKey、error
	GetKeyInfo() (elliptic.Curve, []byte, error)
	//CheckIdentity check identity by ca and check if revoked
	CheckIdentity() (bool, error)
	//GetVKRichBytes get VK richBytes
	GetVKRichBytes() []byte
}

//CA defines the interface
type CA interface {
	GetName() string
	//GenIdentity second param is deprecated; may return "not support", then should report an error in Distributed CA mode
	GenIdentity(cst []byte, _ bool) ([]byte, error)
	String() string
	CheckIdentity(identity []byte) (bool, error)
}

//SessionKey defines the interface
type SessionKey interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
	GetRemoteID() Identity
}

//Logger interface
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Notice(v ...interface{})
	Noticef(format string, v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Critical(v ...interface{})
	Criticalf(format string, v ...interface{})
}

//TransportLevelSecurity TLS
type TransportLevelSecurity interface {
	NewServer(caFile, certFile, keyFile string, isHTTP2 bool) (TransportLevelSecurity, error)
	NewClient(ca, serverNameOverride string, isHTTP2 bool) (TransportLevelSecurity, error)

	NewClientTLSFromFile() (credentials.TransportCredentials, error)
	NewServerTLSFromFile() (credentials.TransportCredentials, error)
	Server(conn net.Conn) (net.Conn, error)
	Client(conn net.Conn) (net.Conn, error)
	Listen(network, laddr string) (net.Listener, error)
}

// Reloader is the interface should be implement by the structures,
// which want to be informed to reload config in commit stages
// while config items or related data they watch are changed by config tx
type Reloader interface {

	// Reload will be called after writing block successfully
	// if the watched config items are changed
	Reload(newValue interface{}) error

	// NeedValue is need new value when call Reload method
	// if return true, call Reload method with newValue
	// if return false, call Reload method with nil
	NeedValue() bool
}

// SKGenerator defines the interface
type SKGenerator interface {
	//GenCertificate pack info to p2p
	GenCertificate(peerHost string, curveType string, certType string) (payLoad []byte, err error)
	//VerifyCertificate unpack and verify info from p2p
	VerifyCertificate(payLoad []byte) (err error)
	//GenShareKey generate session key
	GenShareKey() (SessionKey, error)
}
