package types

// VSet defined to transfer vset from cm to consensus
// contains all hostname in vset,
// id of the hostname is index
type VSet []string

// Hosts defined to transfer host from cm to p2p
// contains new hosts
// key is hostname for new host, value is pub of the host
type Hosts map[string][]byte

// SignLen key length of certs table, content is signature
const SignLen = 80

// VPCerts table vpCerts in statedb
type VPCerts map[string][][]byte

// CertInfo value of table certs
type CertInfo struct {
	IsRevoked  bool
	IsReplace  bool
	IssuerHost string
	OwnerHost  string
	CertDer    []byte
}

// Certs table certs in statedb
type Certs map[[SignLen]byte]CertInfo

// InputReplace input of bvm contract func ReplaceCert
type InputReplace struct {
	HostName string
	Issuers  []string
	AddCerts [][]byte
	DelCerts [][]byte
}

// InputNVPRevoke input of bvm contract func NVPRevoke
type InputNVPRevoke struct {
	NVPHostName string
	VPHostName  string
	Cert        []byte
}

// CertItem cm item of certs
type CertItem struct {
	Revoked Certs
	Alive   Certs
}

// CertBytes all certs for the pk
type CertBytes [][]byte

// =============================================================
//					System Upgrade
// =============================================================

// SupportedVersion defines which versions the node supports, the node
// may support multiple version for one tag.
type SupportedVersion map[VersionTag][]string

// AvailableVersion defines which versions the node can upgrade to, it
// may have multiple version can be selected for one tag.
type AvailableVersion map[VersionTag][]string

// RunningVersion defines which version of the current ledger is running.
type RunningVersion map[VersionTag]string

// SupportedVersions defines which versions all node support.
// Key is hostname.
type SupportedVersions map[string]SupportedVersion

// SetSupportedVersionInput is the input struct of bvm builtin contract.
type SetSupportedVersionInput struct {
	Hostname    string
	VersionList SupportedVersion
}
