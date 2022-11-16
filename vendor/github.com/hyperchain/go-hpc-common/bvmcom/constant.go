package bvmcom

import "time"

const (
	// ProposalPrefix the prefix of proposal
	ProposalPrefix = "_proposal"
)

const (
	// ConsensusPath the config path of consensus
	ConsensusPath = "consensus"

	// ConsensusVSetPath the config path of consensus
	ConsensusVSetPath = "consensus.vset"

	// ConsensusSetPath the config path of consensus set
	ConsensusSetPath = "consensus.set"

	// ConsensusPoolPath the config path of consensus pool
	ConsensusPoolPath = "consensus.pool"

	// FilterPath the config path of api-level tx filters' swith
	FilterPath = "filter"

	// ConsensusAlgoPath the config path of consensus algorithm
	ConsensusAlgoPath = "consensus.algo"
	// ConsensusSetSizePath the config path of rbft slice size
	ConsensusSetSizePath = "consensus.set.set_size"

	// ConsensusBatchSizePath the config path of rbft batch size
	ConsensusBatchSizePath = "consensus.pool.batch_size"

	// ConsensusPoolSizePath the config path of pool size
	ConsensusPoolSizePath = "consensus.pool.pool_size"

	// FilterEnablePath the config path of api-level tx filters' switch
	FilterEnablePath = "filter.enable"

	// FilterRulesPath the config path of api-level tx filters' rules
	FilterRulesPath = "filter.rules"

	// ProposalTimeoutPath the config path of proposal timeout
	ProposalTimeoutPath = "proposal.timeout"

	// ProposalThresholdPath the config path of proposal threshold
	ProposalThresholdPath = "proposal.threshold"

	// ProposalContractVoteEnablePath the config path of proposal contract vote enable
	ProposalContractVoteEnablePath = "proposal.contract.vote.enable"

	// ProposalContractVoteThresholdPath the config path of proposal contract vote threshold
	ProposalContractVoteThresholdPath = "proposal.contract.vote.threshold"

	// ConfigPath the config path of all config
	ConfigPath = "config"

	// GenesisAccountsKey genesis accounts key
	GenesisAccountsKey = "genesis.alloc"
	// GenesisNodesKey genesis nodes key
	GenesisNodesKey = "genesis.nodes"
	// GenesisCAModeKey genesis ca_mode key
	GenesisCAModeKey = "genesis.ca_mode"
	// GenesisRootCAsKey genesis root_ca key
	GenesisRootCAsKey = "genesis.root_ca"

	// GenesisRunningVersions the config path of genesis running versions
	GenesisRunningVersions = "genesis.running.versions.version"

	// CAModePath the key for watch caMode' change
	CAModePath = "ca.ca_mode"

	// RootCAsPath the key for watch root cas' change
	RootCAsPath = "ca.root_cas"

	// GasPrice the key of base gas price
	GasPrice = "executor.gasPrice"
)

const (
	// RolePath the key for watch role's change
	RolePath = "role"

	// VSetPath the key for watch vset's change
	VSetPath = "vset"

	// HostsPath the key for watch hosts's change
	HostsPath = "hosts"

	// CNSPath the key for watch cns's change
	CNSPath = "cns"

	// VPCertsPath the key for watch vpcerts' change
	VPCertsPath = "vpcerts"

	// CertsPath the key for watch certs' change
	CertsPath = "certs"

	// CertsFreezePath the key for freezing certs
	CertsFreezePath = "fcerts"

	// GenesisInfoKey the key for store genesis info in state db.
	GenesisInfoKey = "the_key_for_genesis_info"

	// SRSInfoPath the key for store srs info in state db.
	SRSInfoPath = "srs_info"

	// SRSListPath the srs list path in state db.
	SRSListPath = "the_key_for_srs_list"

	//CurAlgoSetPath algo set in state db
	CurAlgoSetPath = "currentAlgoSet"

	// CertListPath the cert list path in state db.
	CertListPath = "cert"

	// FreezeCertListPath the freeze cert list path in state db.
	FreezeCertListPath = "the_key_for_freeze_cert_list"

	// SupportedVersionsKey the key for storing the supportedVersions in state db
	SupportedVersionsKey = "the_key_for_supported_versions"
	// AvailableVersionsKey the key for storing the availableVersions in state db
	AvailableVersionsKey = "the_key_for_available_versions"
	// RunningVersionsKey the key for storing the runningVersions in state db
	RunningVersionsKey = "the_key_for_running_versions"
)

const (
	// ProposalTimeoutMin min value of proposal time, also the default value of proposal time.
	ProposalTimeoutMin = int64(time.Minute * 5)
)

const (
	//ChainIDKey key of the chainID
	ChainIDKey = "chainID_key"
)

var (
	// AnchorHostnamePrefix key of hostname->namespace map
	AnchorHostnamePrefix = []byte("anchor-hostname-")
	// AnchorNsPrefix key of hostname->namespace map
	AnchorNsPrefix = []byte("anchor-ns-")
	// AnchorExecResultPrefix key of anchor tx result
	AnchorExecResultPrefix = []byte("anchor-result-")
)
