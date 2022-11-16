package protos

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/crypto/sha3"
)

// InitialCheckpoint is the default value of checkpoint,
// include epoch and validator set information
func InitialCheckpoint(vset []*NodeInfo) *QuorumCheckpoint {
	return &QuorumCheckpoint{
		Checkpoint: &Checkpoint{
			Epoch:          0,
			ConsensusState: &Checkpoint_ConsensusState{},
			ExecuteState:   &Checkpoint_ExecuteState{},
			NextEpochState: &Checkpoint_NextEpochState{ValidatorSet: vset},
		},
	}
}

// Validators is a list of NodeInfo
type Validators = []*NodeInfo

// ValidateSetEquals compares two sets of validator and returns whether they are equal
func ValidateSetEquals(a, b Validators) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i].Hostname != b[i].Hostname {
			return false
		}
		if !bytes.Equal(a[i].PubKey, b[i].PubKey) {
			return false
		}
	}
	return true
}

func consensusStateEquals(a, b *Checkpoint_ConsensusState) bool {
	return a.GetRound() == b.GetRound() && a.GetId() == b.GetId()
}

func executeStateEquals(a, b *Checkpoint_ExecuteState) bool {
	return a.GetHeight() == b.GetHeight() && a.GetDigest() == b.GetDigest()
}

func epochStateEquals(a, b *Checkpoint_NextEpochState) bool {
	return ValidateSetEquals(a.GetValidatorSet(), b.GetValidatorSet())
}

// ======================== NodeInfo ========================
// String returns a formatted string for NodeInfo.
func (m *NodeInfo) String() string {
	if m == nil {
		return "NIL"
	}
	return fmt.Sprintf("{host:%s, key:%s}", m.GetHostname(), hex.EncodeToString(m.GetPubKey()))
}

// ======================= Checkpoint =======================

// Hash calculates the crypto hash of Checkpoint.
func (m *Checkpoint) Hash() []byte {
	if m == nil {
		return nil
	}
	res, jErr := proto.Marshal(m)
	if jErr != nil {
		panic(jErr)
	}
	hasher := sha3.NewLegacyKeccak256()
	//nolint
	hasher.Write(res)
	h := hasher.Sum(nil)
	return h
}

// Round returns round of checkpoint.
func (m *Checkpoint) Round() uint64 {
	return m.GetConsensusState().GetRound()
}

// ID returns id of checkpoint.
func (m *Checkpoint) ID() string {
	return m.GetConsensusState().GetId()
}

// Height returns height of checkpoint.
func (m *Checkpoint) Height() uint64 {
	return m.GetExecuteState().GetHeight()
}

// Digest returns digest of checkpoint.
func (m *Checkpoint) Digest() string {
	return m.GetExecuteState().GetDigest()
}

// SetDigest set digest to checkpoint.
func (m *Checkpoint) SetDigest(digest string) {
	if m == nil {
		return
	}
	m.ExecuteState.Digest = digest
}

// ValidatorSet returns vset of checkpoint.
func (m *Checkpoint) ValidatorSet() Validators {
	return m.GetNextEpochState().GetValidatorSet()
}

// ConsensusVersion returns consensus version of checkpoint.
func (m *Checkpoint) ConsensusVersion() string {
	return m.GetNextEpochState().GetConsensusVersion()
}

// SetValidatorSet set vset to checkpoint.
func (m *Checkpoint) SetValidatorSet(vset Validators) {
	if m == nil {
		return
	}
	if m.NextEpochState == nil {
		m.NextEpochState = &Checkpoint_NextEpochState{}
	}
	m.NextEpochState.ValidatorSet = vset
}

// Version returns version of checkpoint.
func (m *Checkpoint) Version() string {
	return m.GetNextEpochState().GetConsensusVersion()
}

// SetVersion set version to checkpoint.
func (m *Checkpoint) SetVersion(version string) {
	if m == nil {
		return
	}
	if m.NextEpochState == nil {
		m.NextEpochState = &Checkpoint_NextEpochState{}
	}
	m.NextEpochState.ConsensusVersion = version
}

// Reconfiguration returns whether this is a config checkpoint
func (m *Checkpoint) Reconfiguration() bool {
	return m.GetNextEpochState() != nil
}

// EndsEpoch returns whether this checkpoint ends the epoch
func (m *Checkpoint) EndsEpoch() bool {
	// TODO(YC): would each reconfiguration end epoch?
	return m.GetNextEpochState() != nil
}

// NextEpoch returns epoch change to checkpoint.
func (m *Checkpoint) NextEpoch() uint64 {
	epoch := m.GetEpoch()
	if m.EndsEpoch() {
		epoch++
	}
	return epoch
}

// String returns a formatted string for Checkpoint.
func (m *Checkpoint) String() string {
	if m == nil {
		return "NIL"
	}
	return fmt.Sprintf("epoch: %d, round: %d, height: %d, hash: %s, commit block id: %s, "+
		"validator set %v, consensus version %s", m.GetEpoch(), m.Round(), m.Height(), m.Digest(),
		m.ID(), m.ValidatorSet(), m.ConsensusVersion())
}

// Equals compares two checkpoint instance and returns whether they are equal
func (m *Checkpoint) Equals(n *Checkpoint) bool {
	return m.GetEpoch() == n.GetEpoch() &&
		consensusStateEquals(m.GetConsensusState(), n.GetConsensusState()) &&
		executeStateEquals(m.GetExecuteState(), n.GetExecuteState()) &&
		epochStateEquals(m.GetNextEpochState(), n.GetNextEpochState())
}

// ======================= SignedCheckpoint =======================

// Hash calculates the crypto hash of SignedCheckpoint.
func (m *SignedCheckpoint) Hash() []byte {
	return m.GetCheckpoint().Hash()
}

// Round returns round of signed checkpoint.
func (m *SignedCheckpoint) Round() uint64 {
	return m.GetCheckpoint().Round()
}

// Epoch returns epoch of signed checkpoint.
func (m *SignedCheckpoint) Epoch() uint64 {
	return m.GetCheckpoint().GetEpoch()
}

// Height returns height of signed checkpoint.
func (m *SignedCheckpoint) Height() uint64 {
	return m.GetCheckpoint().Height()
}

// Digest returns digest of signed checkpoint.
func (m *SignedCheckpoint) Digest() string {
	return m.GetCheckpoint().Digest()
}

// ValidatorSet returns validators of signed checkpoint.
func (m *SignedCheckpoint) ValidatorSet() Validators {
	return m.GetCheckpoint().ValidatorSet()
}

// String returns a formatted string for SignedCheckpoint.
func (m *SignedCheckpoint) String() string {
	if m == nil {
		return "NIL"
	}
	return fmt.Sprintf("Checkpoint %s signed by %s", m.GetCheckpoint(), m.GetAuthor())
}

// ======================= QuorumCheckpoint =======================

// Hash calculates the crypto hash of QuorumCheckpoint.
func (m *QuorumCheckpoint) Hash() []byte {
	return m.GetCheckpoint().Hash()
}

// ID returns id of quorum checkpoint.
func (m *QuorumCheckpoint) ID() string {
	return m.GetCheckpoint().ID()
}

// Digest returns digest of quorum checkpoint.
func (m *QuorumCheckpoint) Digest() string {
	return m.GetCheckpoint().Digest()
}

// ValidatorSet returns validators of quorum checkpoint.
func (m *QuorumCheckpoint) ValidatorSet() Validators {
	return m.GetCheckpoint().ValidatorSet()
}

// Version returns version of quorum checkpoint.
func (m *QuorumCheckpoint) Version() string {
	return m.GetCheckpoint().Version()
}

// Epoch returns epoch of quorum checkpoint.
func (m *QuorumCheckpoint) Epoch() uint64 {
	return m.GetCheckpoint().GetEpoch()
}

// PrevEpoch returns previous epoch change of quorum checkpoint.
func (m *QuorumCheckpoint) PrevEpoch() uint64 {
	epoch := m.Epoch()
	if epoch > 0 && m.EndsEpoch() {
		epoch--
	}
	return epoch
}

// NextEpoch returns epoch change to of quorum checkpoint.
func (m *QuorumCheckpoint) NextEpoch() uint64 {
	return m.GetCheckpoint().NextEpoch()
}

// Round returns round of quorum checkpoint.
func (m *QuorumCheckpoint) Round() uint64 {
	return m.GetCheckpoint().Round()
}

// Height returns height of quorum checkpoint.
func (m *QuorumCheckpoint) Height() uint64 {
	return m.GetCheckpoint().Height()
}

// Reconfiguration returns whether this is a config checkpoint
func (m *QuorumCheckpoint) Reconfiguration() bool {
	return m.GetCheckpoint().Reconfiguration()
}

// EndsEpoch returns whether this checkpoint ends the epoch
func (m *QuorumCheckpoint) EndsEpoch() bool {
	return m.GetCheckpoint().EndsEpoch()
}

// AddSignature adds a certified signature to QuorumCheckpoint.
func (m *QuorumCheckpoint) AddSignature(validator string, signature []byte) {
	if m == nil {
		return
	}
	m.Signatures[validator] = signature
}

// String returns a formatted string for LedgerInfoWithSignatures.
func (m *QuorumCheckpoint) String() string {
	if m == nil {
		return "NIL"
	}
	authors := make([]string, len(m.GetSignatures()))
	i := 0
	for author := range m.GetSignatures() {
		authors[i] = author
		i++
	}
	return fmt.Sprintf("CheckpointInfo: %s signed by: %+v", m.GetCheckpoint(), authors)
}

// ======================= EpochChangeProof =======================

// NewEpochChangeProof create a new proof with given checkpoints
func NewEpochChangeProof(checkpoints ...*QuorumCheckpoint) *EpochChangeProof {
	// checkpoints should end epoch
	return &EpochChangeProof{Checkpoints: checkpoints}
}

// StartEpoch returns start epoch of the proof
func (m *EpochChangeProof) StartEpoch() uint64 {
	return m.First().Epoch()
}

// NextEpoch returns the next epoch to change to of the proof
func (m *EpochChangeProof) NextEpoch() uint64 {
	return m.Last().NextEpoch()
}

// First returns the first checkpoint of the proof
func (m *EpochChangeProof) First() *QuorumCheckpoint {
	if m.IsEmpty() {
		return nil
	}
	return m.Checkpoints[0]
}

// Last returns the last checkpoint of the proof
func (m *EpochChangeProof) Last() *QuorumCheckpoint {
	if m.IsEmpty() {
		return nil
	}
	return m.Checkpoints[len(m.Checkpoints)-1]
}

// IsEmpty returns whether the proof is empty
func (m *EpochChangeProof) IsEmpty() bool {
	if m == nil {
		return true
	}
	return len(m.GetCheckpoints()) == 0
}

// String returns a formatted string for EpochChangeProof.
func (m *EpochChangeProof) String() string {
	if m == nil {
		return "NIL"
	}
	return fmt.Sprintf("EpochChangeProof: checkpoints %s, more %d, author %s", m.GetCheckpoints(), m.GetMore(), m.GetAuthor())
}

// ======================= LedgerInfo =======================

// String returns a formatted string for LedgerInfo.
func (m *LedgerInfo) String() string {
	if m == nil {
		return "NIL"
	}
	return fmt.Sprintf("[CommitBlockInfo: %s, ConsensusDataHash: %s]", m.GetCommitInfo(), m.GetConsensusDataHash())
}

// Epoch returns the epoch of LedgerInfo.
func (m *LedgerInfo) Epoch() uint64 {
	return m.GetCommitInfo().GetEpoch()
}

// Round returns the round of LedgerInfo.
func (m *LedgerInfo) Round() uint64 {
	return m.GetCommitInfo().GetRound()
}

// ConsensusBlockID returns the consensus block id of LedgerInfo.
func (m *LedgerInfo) ConsensusBlockID() string {
	return m.GetCommitInfo().GetId()
}

// TimestampNanos returns the timestamp of LedgerInfo.
func (m *LedgerInfo) TimestampNanos() int64 {
	return m.GetCommitInfo().GetTimestampNanos()
}

// Reconfiguration returns the next validator set of LedgerInfo.
func (m *LedgerInfo) Reconfiguration() bool {
	return m.GetCommitInfo().GetReconfiguration()
}

// SetConsensusDataHash sets the payload ConsensusDataHash of LedgerInfo.
func (m *LedgerInfo) SetConsensusDataHash(hash string) {
	m.ConsensusDataHash = hash
}

// Hash returns the crypto hash of LedgerInfo.
func (m *LedgerInfo) Hash() []byte {
	res, jErr := proto.Marshal(m)
	if jErr != nil {
		panic(jErr)
	}
	hasher := sha3.NewLegacyKeccak256()
	//nolint
	hasher.Write(res)
	h := hasher.Sum(nil)
	return h
}

// ======================= LedgerInfoWithSignatures =======================

// String returns a formatted string for LedgerInfoWithSignatures.
func (m *LedgerInfoWithSignatures) String() string {
	if m == nil {
		return "NIL"
	}
	authors := make([]string, len(m.Signatures))
	i := 0
	for author := range m.Signatures {
		authors[i] = author
		i++
	}

	return fmt.Sprintf("LedgerInfo: %s signed by: %+v", m.GetLedgerInfo(), authors)
}

// AddSignature adds a certified signature to LedgerInfoWithSignatures.
func (m *LedgerInfoWithSignatures) AddSignature(validator string, signature []byte) {
	m.Signatures[validator] = signature
}

// RemoveSignature removes a certified signature from LedgerInfoWithSignatures.
func (m *LedgerInfoWithSignatures) RemoveSignature(validator string) {
	delete(m.Signatures, validator)
}

// ======================= ConsensusBlockInfo =======================

// String returns a readable block info string.
func (m *ConsensusBlockInfo) String() string {
	if m == nil {
		return "NIL"
	}
	return fmt.Sprintf("[id: %s, epoch: %d, round: %d]", m.Id, m.Epoch, m.Round)
}

// HasReconfiguration returns if the block is a configuration block.
func (m *ConsensusBlockInfo) HasReconfiguration() bool {
	return m.GetReconfiguration()
}

// Equal indicates if two ConsensusBlockInfo is equal.
func (m *ConsensusBlockInfo) Equal(another *ConsensusBlockInfo) bool {
	if m.Reconfiguration != another.Reconfiguration {
		return false
	}

	return m.Epoch == another.Epoch && m.Round == another.Round && m.Id == another.Id &&
		m.Height == another.Height && m.TimestampNanos == another.TimestampNanos
}

// ======================= QuorumCert =======================

// IsConfig returns if it is a config QC.
func (qc *QuorumCert) IsConfig() bool {
	return qc.GetVoteData().GetProposed().GetReconfiguration()
}

// String returns a readable quorum cert string.
func (qc *QuorumCert) String() string {
	if qc == nil {
		return "NIL"
	}
	return fmt.Sprintf("ProposedBlockInfo: %s, ParentBlockInfo: %s, LedgerInfo: %s, Reconfiguration: %v",
		qc.CertifiedBlock(), qc.ParentBlock(), qc.LedgerInfo(), qc.IsConfig())
}

// CertifiedBlock returns the certified block info.
func (qc *QuorumCert) CertifiedBlock() *ConsensusBlockInfo {
	return qc.GetVoteData().GetProposed()
}

// CertifiedBlockID returns the certified block id.
func (qc *QuorumCert) CertifiedBlockID() string {
	return qc.CertifiedBlock().GetId()
}

// CertifiedBlockRound returns the round of certified block.
func (qc *QuorumCert) CertifiedBlockRound() uint64 {
	return qc.CertifiedBlock().GetRound()
}

// ParentBlock returns the block info of certified block's parent.
func (qc *QuorumCert) ParentBlock() *ConsensusBlockInfo {
	return qc.GetVoteData().GetParent()
}

// ParentBlockID returns the block id of certified block's parent.
func (qc *QuorumCert) ParentBlockID() string {
	return qc.ParentBlock().GetId()
}

// ParentBlockRound returns the round of certified block's parent.
func (qc *QuorumCert) ParentBlockRound() uint64 {
	return qc.ParentBlock().GetRound()
}

// CommitBlock returns the block info of committed block.
func (qc *QuorumCert) CommitBlock() *ConsensusBlockInfo {
	return qc.LedgerInfo().GetCommitInfo()
}

// CommitBlockID returns the block id of committed block.
func (qc *QuorumCert) CommitBlockID() string {
	return qc.CommitBlock().GetId()
}

// CommitBlockRound returns the round of committed block.
func (qc *QuorumCert) CommitBlockRound() uint64 {
	return qc.CommitBlock().GetRound()
}

// Epoch returns the epoch of voting proposal
func (qc *QuorumCert) Epoch() uint64 {
	return qc.CertifiedBlock().GetEpoch()
}

// EndsEpoch indicates if the QC commits reconfiguration and starts a new epoch
func (qc *QuorumCert) EndsEpoch() bool {
	return qc.CommitBlock().GetReconfiguration()
}

// Height returns the height of voting proposal
func (qc *QuorumCert) Height() uint64 {
	return qc.CertifiedBlock().GetHeight()
}

// LedgerInfo returns the signed ledger associated with this QC.
func (qc *QuorumCert) LedgerInfo() *LedgerInfo {
	return qc.GetSignedLedgerInfo().GetLedgerInfo()
}

// SignatureEntry is the entry with signature and author
type SignatureEntry struct {
	Author    string
	Signature []byte
}

// SortedSignatures is a sortable entry list
type SortedSignatures []*SignatureEntry

// Len return size of SortedSignatures
func (s SortedSignatures) Len() int {
	return len(s)
}

// Less judge the order of signature author
func (s SortedSignatures) Less(i, j int) bool {
	return s[i].Author < s[j].Author
}

// Swap signatures at i and j
func (s SortedSignatures) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Authors return author list in SortedSignatures
func (s SortedSignatures) Authors() []string {
	res := make([]string, s.Len())
	for i, entry := range s {
		res[i] = entry.Author
	}
	return res
}

// SortedSignatures return SortedSignatures in QuorumCert
func (qc *QuorumCert) SortedSignatures() SortedSignatures {
	sigs := qc.GetSignedLedgerInfo().GetSignatures()
	s := make(SortedSignatures, 0, len(sigs))
	for author, sig := range sigs {
		s = append(s, &SignatureEntry{
			Author:    author,
			Signature: sig,
		})
	}
	sort.Sort(s)
	return s
}

// ======================= VoteData =======================

// BlockID returns the block id of this vote.
func (vd *VoteData) BlockID() string {
	return vd.GetProposed().GetId()
}

// BlockEpoch returns the block epoch of this vote.
func (vd *VoteData) BlockEpoch() uint64 {
	return vd.GetProposed().GetEpoch()
}

// BlockRound returns the round of this vote.
func (vd *VoteData) BlockRound() uint64 {
	return vd.GetProposed().GetRound()
}

// BlockTimestamp returns the timestamp of this vote.
func (vd *VoteData) BlockTimestamp() int64 {
	return vd.GetProposed().GetTimestampNanos()
}

// ParentBlockID returns the parent block id of this vote.
func (vd *VoteData) ParentBlockID() string {
	return vd.GetParent().GetId()
}

// ParentBlockRound returns the parent round of this vote.
func (vd *VoteData) ParentBlockRound() uint64 {
	return vd.GetParent().GetRound()
}

// ParentBlockEpoch returns the parent epoch of this vote.
func (vd *VoteData) ParentBlockEpoch() uint64 {
	return vd.GetParent().GetEpoch()
}

// ParentBlockTimestamp returns the parent timestamp of this vote.
func (vd *VoteData) ParentBlockTimestamp() int64 {
	return vd.GetParent().GetTimestampNanos()
}

// Hash calculates the crypto hash of VoteData
func (vd *VoteData) Hash() []byte {
	res, jErr := proto.Marshal(vd)
	if jErr != nil {
		panic(jErr)
	}
	return CalculateMD5Hash(res)
}

// String returns a readable string of VoteData.
func (vd *VoteData) String() string {
	return fmt.Sprintf("VoteData: [block id: %s, epoch: %d, round: %d, "+
		"parent_block_id: %s, parent_block_round: %d]", vd.Proposed.Id, vd.Proposed.Epoch,
		vd.Proposed.Round, vd.Parent.Id, vd.Parent.Round)
}

// NewNilBlock creates nil blocks which are special: they're not carrying any real payload and are generated
// independently by different validators just to fill in the round.
// The NIL blocks are special: they're not carrying any real payload and are generated
// independently by different validators just to fill in the round with some QC.
// NOTE. HQC is not included into InnerBlock until we are ready to execute this InnerBlock.
//func NewNilBlock(round uint64, qc *QuorumCert) *InnerBlock {
//	blockData := NewNilBlockData(round, qc)
//	return &InnerBlock{
//		Id:        hex.EncodeToString(blockData.Hash()),
//		BlockData: blockData,
//		Signature: nil,
//	}
//}

//// NewProposalBlockWithSignature create a block directly.
//func NewProposalBlockWithSignature(blockData *BlockData, blockHash, signature []byte) *InnerBlock {
//	return &InnerBlock{
//		Id:        hex.EncodeToString(blockHash),
//		BlockData: blockData,
//		Signature: signature,
//	}
//}

// ======================= InnerBlock =======================

// Payload is a hash batch of transactions
type Payload = [][]byte

// Transactions is a list of transactions
type Transactions = []*Transaction

// String returns a readable block info string.
func (b *InnerBlock) String() string {
	if b == nil {
		return "NIL"
	}
	return fmt.Sprintf("[id: %s, epoch: %d, round: %d, parent_id: %s]", b.Id, b.Epoch(), b.Round(), b.ParentID())
}

// GetBlockAuthor returns block author if any.
func (b *InnerBlock) GetBlockAuthor() (string, bool) {
	return b.GetBlockData().GetBlockAuthor()
}

// Epoch returns epoch of this block.
func (b *InnerBlock) Epoch() uint64 {
	return b.GetBlockData().GetEpoch()
}

// Round returns round of this block.
func (b *InnerBlock) Round() uint64 {
	return b.GetBlockData().GetRound()
}

// ParentID returns parent id of this block.
func (b *InnerBlock) ParentID() string {
	return b.GetBlockData().GetParentId()
}

// GetBlockPayload returns block payload if any.
func (b *InnerBlock) GetBlockPayload() (Payload, bool) {
	return b.GetBlockData().GetBlockPayload()
}

// GetBlockTransactions returns block transactions if any.
func (b *InnerBlock) GetBlockTransactions() Transactions {
	return b.GetBlockData().GetTransactions()
}

// QuorumCert returns quorum cert of this block.
func (b *InnerBlock) QuorumCert() *QuorumCert {
	return b.GetBlockData().GetQuorumCert()
}

// TimestampNanos returns timestamp of this block.
func (b *InnerBlock) TimestampNanos() int64 {
	return b.GetBlockData().GetTimestampNanos()
}

// GenerateBlockInfo generates a block info from current block.
func (b *InnerBlock) GenerateBlockInfo() *ConsensusBlockInfo {
	var reconfiguration bool
	if b.IsNilBlock() {
		reconfiguration = b.QuorumCert().IsConfig()
	} else {
		reconfiguration = b.IsConfigBlock()
	}
	height := b.QuorumCert().Height()
	if !b.IsEmptyBlock() {
		height++
	}
	return &ConsensusBlockInfo{
		Epoch:           b.Epoch(),
		Round:           b.Round(),
		Id:              b.GetId(),
		Height:          height,
		TimestampNanos:  b.TimestampNanos(),
		Reconfiguration: reconfiguration,
		Proof:           b.IsProofBlock(),
	}
}

// IsGenesisBlock returns if current block is the genesis block.
func (b *InnerBlock) IsGenesisBlock() bool {
	return b.GetBlockData().IsGenesisBlock()
}

// IsNilBlock returns if current block is the nil block.
func (b *InnerBlock) IsNilBlock() bool {
	return b.GetBlockData().IsNilBlock()
}

// IsEmptyBlock returns if current block is an empty block.
func (b *InnerBlock) IsEmptyBlock() bool {
	return b.GetBlockData().IsEmpty()
}

// IsConfigBlock returns if current block is a config block.
func (b *InnerBlock) IsConfigBlock() bool {
	return b.GetBlockData().IsConfigBlock()
}

// IsProofBlock returns if current block is a proof block.
func (b *InnerBlock) IsProofBlock() bool {
	return b.GetBlockData().IsProofBlock()
}

// EndsEpoch returns if current block ends current epoch
func (b *InnerBlock) EndsEpoch() bool {
	if !b.IsConfigBlock() {
		return false
	}
	// Light block can only estimate whether ends the epoch (request complete tx for a certain judgment)
	if b.GetBlockData().IsLightBlock() {
		return len(b.GetBlockData().GetPayload()) == 1
	}
	return len(b.GetBlockTransactions()) == 1 && b.GetBlockTransactions()[0].IsConfigTx()
}

// Copy create a block instance from exist.
func (b *InnerBlock) Copy() *InnerBlock {
	if b == nil {
		return nil
	}
	return &InnerBlock{
		Id:        b.GetId(),
		BlockData: b.GetBlockData().Copy(),
		Signature: b.GetSignature(), // needn't deep-copy for contents is not overwritten
	}
}

// ======================= BlockData =======================

// GetBlockAuthor returns block author if any.
func (bd *BlockData) GetBlockAuthor() (string, bool) {
	// if current block is a nil block or genesis block, return not found directly.
	if bd.IsNilBlock() || bd.IsGenesisBlock() {
		return "", false
	}

	return bd.GetAuthor(), true
}

// GetBlockPayload returns block payload if any.
func (bd *BlockData) GetBlockPayload() (Payload, bool) {
	// if current block is a nil block or genesis block, return not found directly.
	if bd.IsEmpty() || bd.IsGenesisBlock() || !bd.IsLightBlock() {
		return nil, false
	}
	return bd.Payload, true
}

// Hash calculates the crypto hash of a BlockData.
func (bd *BlockData) Hash() []byte {
	// exclude QuorumCert and txs for calculate block hash
	qc := bd.QuorumCert
	payload := bd.Payload
	txs := bd.Transactions
	bd.QuorumCert = nil
	bd.Payload = nil
	bd.Transactions = nil

	res, jErr := proto.Marshal(bd)
	if jErr != nil {
		panic(jErr)
	}

	sortSigs := qc.SortedSignatures()

	var blockHash []byte
	if len(txs) > 0 {
		blockHash = CalculateTxRootMD5(txs, nil, res, sortSigs)
	} else {
		// use md5 root hash for payload if txs is empty
		blockHash = CalculateTxRootMD5(nil, payload, res, sortSigs)
	}
	bd.QuorumCert = qc
	bd.Payload = payload
	bd.Transactions = txs
	return blockHash
}

// IsEmpty returns if current block data is empty.
// True for IsNilBlock flag is true or payload and transactions are both empty
func (bd *BlockData) IsEmpty() bool {
	return bd.IsNilBlock() || (len(bd.GetPayload()) == 0 && len(bd.GetTransactions()) == 0)
}

// IsNilBlock returns if current block is a nil block
func (bd *BlockData) IsNilBlock() bool {
	return bd.GetBlockType() == BlockData_NIL
}

// IsGenesisBlock returns if current block is a genesis block
func (bd *BlockData) IsGenesisBlock() bool {
	return bd.GetBlockType() == BlockData_GENESIS
}

// IsLightBlock returns if current block is a light block (only includes txHash batch).
// True for transactions are empty but payload not
func (bd *BlockData) IsLightBlock() bool {
	return len(bd.Payload) > 0 && len(bd.Transactions) == 0
}

// IsConfigBlock returns if current block data is used to change configuration.
func (bd *BlockData) IsConfigBlock() bool {
	return bd.GetBlockType() == BlockData_CONFIG
}

// IsProofBlock returns if current block is a proof block
func (bd *BlockData) IsProofBlock() bool {
	return bd.GetBlockType() == BlockData_PROOF
}

// Copy create a BlockData instance from exist.
func (bd *BlockData) Copy() *BlockData {
	if bd == nil {
		return nil
	}
	return &BlockData{
		Epoch:          bd.Epoch,
		Round:          bd.Round,
		TimestampNanos: bd.TimestampNanos,
		ParentId:       bd.ParentId,
		QuorumCert:     bd.QuorumCert,   // needn't deep-copy for contents is not overwritten
		Payload:        bd.Payload,      // needn't deep-copy for contents is not overwritten
		Transactions:   bd.Transactions, // needn't deep-copy for contents is not overwritten
		Author:         bd.Author,
		BlockType:      bd.BlockType,
	}
}

// ======================= Hash Calculator ======================

// CalculateMD5Hash calculates the crypto hash of payload.
func CalculateMD5Hash(payload []byte) []byte {
	h := md5.New()
	_, _ = h.Write(payload)
	return h.Sum(nil)
}

// CalculateTxRootMD5 calculate hash with transactions.
func CalculateTxRootMD5(validTxs Transactions, payload Payload, res []byte, sigs SortedSignatures) []byte {
	h := md5.New()
	_, _ = h.Write(res)
	for _, sig := range sigs {
		_, _ = h.Write([]byte(sig.Author))
	}
	for _, tx := range validTxs {
		_, _ = h.Write(tx.TransactionHash)
	}
	for _, txHash := range payload {
		_, _ = h.Write(txHash)
	}
	return h.Sum(nil)
}
