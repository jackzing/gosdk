package tls

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"strings"

	"github.com/hyperchain/go-hpc-common/cryptocom"
	"github.com/hyperchain/go-hpc-msp/plugin/software"

	"google.golang.org/grpc/credentials"
)

// NewClientTLSFromFile constructs TLS credentials from the input certificate file for client.
// serverNameOverride is for testing only. If set to a non empty string,
// it will override the virtual host name of authority (e.g. :authority header field) in requests.
func NewClientTLSFromFile(ca, serverNameOverride string, engine cryptocom.Engine) (credentials.TransportCredentials, error) {
	if reflect.DeepEqual(engine, nil) {
		return nil, fmt.Errorf("param engine is nil")
	}
	b, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, err
	}
	cp := software.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return &tlsCreds{cloneTLSConfig(&Config{ServerName: serverNameOverride, RootCAs: cp})}, nil
}

// NewServerTLSFromFile constructs TLS credentials from the input certificate file and key
// file for server.
func NewServerTLSFromFile(certFile, keyFile string, engine cryptocom.Engine) (credentials.TransportCredentials, error) {
	if reflect.DeepEqual(engine, nil) {
		return nil, fmt.Errorf("param engine is nil")
	}
	certs, err := LoadX509KeyPairs(software.GetSoftwareEngine(), certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tlsCreds{
		config: cloneTLSConfig(&Config{
			ClientAuth: VerifyClientCertIfGiven,
			MinVersion: VersionTLS12,
			CipherSuites: []uint16{
				TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				GMTLS_SM2_SM4_SM3,
			},
			Certificates: certs,
		}),
	}, nil
}

// tlsInfo contains the auth information for a TLS authenticated connection.
// It implements the AuthInfo interface.
type tlsInfo struct {
	State ConnectionState
}

// AuthType returns the type of tlsInfo as a string.
func (t tlsInfo) AuthType() string {
	return "tls"
}

// tlsCreds is the credentials required for authenticating a connection using TLS.
type tlsCreds struct {
	// TLS configuration
	config *Config
}

func (c tlsCreds) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
		ServerName:       c.config.ServerName,
	}
}

func (c *tlsCreds) ClientHandshake(ctx context.Context, authority string, rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	// use local cfg to avoid clobbering ServerName if using multiple endpoints
	cfg := cloneTLSConfig(c.config)
	if cfg.ServerName == "" {
		colonPos := strings.LastIndex(authority, ":")
		if colonPos == -1 {
			colonPos = len(authority)
		}
		cfg.ServerName = authority[:colonPos]
	}
	conn := Client(rawConn, cfg)
	errChannel := make(chan error, 1)
	go func() {
		errChannel <- conn.Handshake()
	}()
	select {
	case err := <-errChannel:
		if err != nil {
			return nil, nil, err
		}
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
	ttt := tlsInfo{conn.ConnectionState()}
	return conn, ttt, nil
}

func (c *tlsCreds) ServerHandshake(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	conn := Server(rawConn, c.config)
	if err := conn.Handshake(); err != nil {
		return nil, nil, err
	}
	return conn, tlsInfo{conn.ConnectionState()}, nil
}

func (c *tlsCreds) Clone() credentials.TransportCredentials {
	return &tlsCreds{cloneTLSConfig(c.config)}
}

func (c *tlsCreds) OverrideServerName(serverNameOverride string) error {
	c.config.ServerName = serverNameOverride
	return nil
}

// cloneTLSConfig returns a shallow clone of the exported
// fields of cfg, ignoring the unexported sync.Once, which
// contains a mutex and must not be copied.
//
// If cfg is nil, a new zero Config is returned.
//
func cloneTLSConfig(cfg *Config) *Config {
	if cfg == nil {
		panic("tls: config is nil")
	}

	return cfg.Clone()
}
