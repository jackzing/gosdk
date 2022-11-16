// Package tls partially implements TLS 1.2, as specified in RFC 5246.
package tls

// BUG(agl): The crypto/tls package only implements some countermeasures
// against Lucky13 attacks on CBC-mode encryption, and only on SHA1
// variants. See http://www.isg.rhul.ac.uk/tls/TLStiming.pdf and
// https://www.imperialviolet.org/2013/02/04/luckythirteen.html.

import (
	"bytes"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	hc "github.com/hyperchain/go-hpc-common/cryptocom"
	"github.com/hyperchain/go-hpc-msp/plugin/software"
	"github.com/meshplus/crypto"
)

// Server returns a new TLS server side connection
// using conn as the underlying transport.
// The configuration config must be non-nil and must include
// at least one certificate or else set GetCertificate.
func Server(conn net.Conn, config *Config) *Conn {
	Infof("server, get conn")
	return &Conn{conn: conn, config: config}
}

// Client returns a new TLS client side connection
// using conn as the underlying transport.
// The config cannot be nil: users must set either ServerName or
// InsecureSkipVerify in the config.
func Client(conn net.Conn, config *Config) *Conn {
	return &Conn{conn: conn, config: config, isClient: true}
}

// A listener implements a network listener (net.Listener) for TLS connections.
type listener struct {
	net.Listener
	config *Config
}

// Accept waits for and returns the next incoming TLS connection.
// The returned connection is of type *Conn.
func (l *listener) Accept() (net.Conn, error) {
	Infof("accept")
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return Server(c, l.config), nil
}

// NewListener creates a Listener which accepts connections from an inner
// Listener and wraps each connection with Server.
// The configuration config must be non-nil and must include
// at least one certificate or else set GetCertificate.
func NewListener(inner net.Listener, config *Config) net.Listener {
	l := new(listener)
	l.Listener = inner
	l.config = config
	return l
}

// Listen creates a TLS listener accepting connections on the
// given network address using net.Listen.
// The configuration config must be non-nil and must include
// at least one certificate or else set GetCertificate.
func Listen(network, laddr string, config *Config) (net.Listener, error) {
	if config == nil || (len(config.Certificates) == 0 && config.GetCertificate == nil) {
		return nil, errors.New("tls: neither Certificates nor GetCertificate set in Config")
	}
	l, err := net.Listen(network, laddr)
	if err != nil {
		return nil, err
	}
	return NewListener(l, config), nil
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "tls: DialWithDialer timed out" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

// DialWithDialer connects to the given network address using dialer.Dial and
// then initiates a TLS handshake, returning the resulting TLS connection. Any
// timeout or deadline given in the dialer apply to connection and TLS
// handshake as a whole.
//
// DialWithDialer interprets a nil configuration as equivalent to the zero
// configuration; see the documentation of Config for the defaults.
func DialWithDialer(dialer *net.Dialer, network, addr string, config *Config) (*Conn, error) {
	// We want the Timeout and Deadline values from dialer to cover the
	// whole process: TCP connection and TLS handshake. This means that we
	// also need to start our own timers now.
	timeout := dialer.Timeout

	if !dialer.Deadline.IsZero() {
		deadlineTimeout := time.Until(dialer.Deadline)
		if timeout == 0 || deadlineTimeout < timeout {
			timeout = deadlineTimeout
		}
	}

	var errChannel chan error

	if timeout != 0 {
		errChannel = make(chan error, 2)
		time.AfterFunc(timeout, func() {
			errChannel <- timeoutError{}
		})
	}

	rawConn, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	colonPos := strings.LastIndex(addr, ":")
	if colonPos == -1 {
		colonPos = len(addr)
	}
	hostname := addr[:colonPos]

	if config == nil {
		config = defaultConfig()
	}
	// If no ServerName is set, infer the ServerName
	// from the hostname we're connecting to.
	if config.ServerName == "" {
		// Make a copy to avoid polluting argument or default.
		c := config.Clone()
		c.ServerName = hostname
		config = c
	}

	conn := Client(rawConn, config)

	if timeout == 0 {
		err = conn.Handshake()
	} else {
		go func() {
			errChannel <- conn.Handshake()
		}()

		err = <-errChannel
	}

	if err != nil {
		_ = rawConn.Close()
		return nil, err
	}

	return conn, nil
}

// Dial connects to the given network address using net.Dial
// and then initiates a TLS handshake, returning the resulting
// TLS connection.
// Dial interprets a nil configuration as equivalent to
// the zero configuration; see the documentation of Config
// for the defaults.
func Dial(network, addr string, config *Config) (*Conn, error) {
	return DialWithDialer(new(net.Dialer), network, addr, config)
}

func pemDecode(multiPEM []byte) (currentBlock *pem.Block, currentPEM string, restPEM []byte) {
	multiPEM = bytes.TrimSpace(multiPEM)
	currentBlock, restPEM = pem.Decode(multiPEM)
	currentPEM = string(multiPEM[:len(multiPEM)-len(currentPEM)])
	return
}

//LoadX509KeyPairs load x509 certs, support GM double certificate
func LoadX509KeyPairs(engine hc.Engine, certFile, keyFile string) ([]Certificate, error) {
	fail := func(err error) ([]Certificate, error) { return []Certificate{}, err }
	certPEMBlock, err := ioutil.ReadFile(certFile)
	if err != nil {
		return fail(err)
	}
	keyPEMBlock, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return fail(err)
	}
	var certDerBlock, keyDerBlock *pem.Block
	var result []Certificate
	var currentPEM, currentCertPem string
	for {
		certDerBlock, currentCertPem, certPEMBlock = pemDecode(certPEMBlock)
		keyDerBlock, currentPEM, keyPEMBlock = pemDecode(keyPEMBlock)
		//判断文件类型为pem
		if certDerBlock == nil || keyDerBlock == nil {
			break
		}
		//判断pem文件类型
		if certDerBlock.Type != "CERTIFICATE" || !strings.HasSuffix(keyDerBlock.Type, "PRIVATE KEY") {
			return fail(fmt.Errorf("file is not cert or key"))
		}
		//解析并判断key类型
		x509Cert, err := software.ParseCertificate(currentCertPem)
		if err != nil {
			return fail(err)
		}

		privKey, err := software.UnmarshalPrivateKey(currentPEM)
		if err != nil {
			return fail(errors.New("tls: failed to parse private key"))
		}

		if err := checkKeyType(x509Cert.PublicKey, privKey); err != nil {
			return fail(err)
		}

		if len(result) != 0 {
			result[0].Certificate = append(result[0].Certificate, certDerBlock.Bytes)
		}
		result = append(result, Certificate{
			PrivateKey:  privKey,
			Leaf:        x509Cert,
			Certificate: [][]byte{certDerBlock.Bytes},
		})
	}
	if len(result) < 1 {
		return fail(fmt.Errorf("parse cert and key from file error"))
	}
	return result, nil
}

// LoadX509KeyPair reads and parses a public/private key pair from a pair
// of files. The files must contain PEM encoded data. The certificate file
// may contain intermediate certificates following the leaf certificate to
// form a certificate chain. On successful return, Certificate.Leaf will
// be nil because the parsed form of the certificate is not retained.
// Deprecated
func LoadX509KeyPair(certFile, keyFile string) (Certificate, error) {
	certPEMBlock, err := ioutil.ReadFile(certFile)
	if err != nil {
		return Certificate{}, err
	}
	keyPEMBlock, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return Certificate{}, err
	}
	return X509KeyPair(certPEMBlock, keyPEMBlock)
}

// X509KeyPair parses a public/private key pair from a pair of
// PEM encoded data. On successful return, Certificate.Leaf will be nil because
// the parsed form of the certificate is not retained.
func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (Certificate, error) {
	fail := func(err error) (Certificate, error) { return Certificate{}, err }

	var cert Certificate
	var skippedBlockTypes []string
	for {
		var certDERBlock *pem.Block
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		} else {
			skippedBlockTypes = append(skippedBlockTypes, certDERBlock.Type)
		}
	}

	if len(cert.Certificate) == 0 {
		if len(skippedBlockTypes) == 0 {
			return fail(errors.New("tls: failed to find any PEM data in certificate input"))
		}
		if len(skippedBlockTypes) == 1 && strings.HasSuffix(skippedBlockTypes[0], "PRIVATE KEY") {
			return fail(errors.New("tls: failed to find certificate PEM data in certificate input, but did find a private key; PEM inputs may have been switched"))
		}
		return fail(fmt.Errorf("tls: failed to find \"CERTIFICATE\" PEM block in certificate input after skipping PEM blocks of the following types: %v", skippedBlockTypes))
	}

	skippedBlockTypes = skippedBlockTypes[:0]
	var keyDERBlock *pem.Block
	var currentPEM string
	for {
		keyDERBlock, currentPEM, keyPEMBlock = pemDecode(keyPEMBlock)
		if keyDERBlock == nil {
			if len(skippedBlockTypes) == 0 {
				return fail(errors.New("tls: failed to find any PEM data in key input"))
			}
			if len(skippedBlockTypes) == 1 && skippedBlockTypes[0] == "CERTIFICATE" {
				return fail(errors.New("tls: found a certificate rather than a key in the PEM for the private key"))
			}
			return fail(fmt.Errorf("tls: failed to find PEM block with type ending in \"PRIVATE KEY\" in key input after skipping PEM blocks of the following types: %v", skippedBlockTypes))
		}
		if keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY") {
			break
		}
		skippedBlockTypes = append(skippedBlockTypes, keyDERBlock.Type)
	}

	// We don't need to parse the public key for TLS, but we so do anyway
	// to check that it looks sane and matches the private key.
	x509Cert, err := software.ParseCertificate(string(cert.Certificate[0]))
	if err != nil {
		return fail(err)
	}

	engine := software.GetSoftwareEngine()
	cert.PrivateKey, err = engine.GetSignKey(currentPEM) //PEM, cause GetSignKey need PEM
	if err != nil {
		return fail(err)
	}

	if err := checkKeyType(x509Cert.PublicKey, cert.PrivateKey); err != nil {
		return fail(err)
	}

	return cert, nil
}

func checkKeyType(pub crypto.VerifyKey, priv crypto.SignKey) error {
	if pub.GetKeyInfo() != priv.GetKeyInfo() {
		return errors.New("tls: private key type does not match public key type")
	}
	if !bytes.Equal(pub.Bytes(), priv.Bytes()) {
		return errors.New("tls: private key does not match public key")
	}
	return nil
}
