// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proposal.proto

package bvm

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type ProposalData_Status int32

const (
	ProposalData_CREATE      ProposalData_Status = 0
	ProposalData_VOTING      ProposalData_Status = 1
	ProposalData_REJECT      ProposalData_Status = 2
	ProposalData_WAITING_EXE ProposalData_Status = 3
	ProposalData_CANCEL      ProposalData_Status = 4
	ProposalData_COMPLETED   ProposalData_Status = 5
	ProposalData_TIMEOUT     ProposalData_Status = 6
)

var ProposalData_Status_name = map[int32]string{
	0: "CREATE",
	1: "VOTING",
	2: "REJECT",
	3: "WAITING_EXE",
	4: "CANCEL",
	5: "COMPLETED",
	6: "TIMEOUT",
}

var ProposalData_Status_value = map[string]int32{
	"CREATE":      0,
	"VOTING":      1,
	"REJECT":      2,
	"WAITING_EXE": 3,
	"CANCEL":      4,
	"COMPLETED":   5,
	"TIMEOUT":     6,
}

func (x ProposalData_Status) String() string {
	return proto.EnumName(ProposalData_Status_name, int32(x))
}

func (ProposalData_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c3ac5ce23bf32d05, []int{0, 0}
}

type ProposalData_PType int32

const (
	ProposalData_CONFIG     ProposalData_PType = 0
	ProposalData_PERMISSION ProposalData_PType = 1
	ProposalData_NODE       ProposalData_PType = 2
	ProposalData_CNS        ProposalData_PType = 3
	ProposalData_CONTRACT   ProposalData_PType = 4
	ProposalData_CA         ProposalData_PType = 5
)

var ProposalData_PType_name = map[int32]string{
	0: "CONFIG",
	1: "PERMISSION",
	2: "NODE",
	3: "CNS",
	4: "CONTRACT",
	5: "CA",
}

var ProposalData_PType_value = map[string]int32{
	"CONFIG":     0,
	"PERMISSION": 1,
	"NODE":       2,
	"CNS":        3,
	"CONTRACT":   4,
	"CA":         5,
}

func (x ProposalData_PType) String() string {
	return proto.EnumName(ProposalData_PType_name, int32(x))
}

func (ProposalData_PType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c3ac5ce23bf32d05, []int{0, 1}
}

type ProposalData struct {
	Id        uint64              `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Code      []byte              `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	Timestamp int64               `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Timeout   int64               `protobuf:"varint,4,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Status    ProposalData_Status `protobuf:"varint,5,opt,name=status,proto3,enum=bvm.ProposalData_Status" json:"status,omitempty"`
	Assentor  []*VoteInfo         `protobuf:"bytes,6,rep,name=assentor,proto3" json:"assentor,omitempty"`
	Objector  []*VoteInfo         `protobuf:"bytes,7,rep,name=objector,proto3" json:"objector,omitempty"`
	Threshold uint32              `protobuf:"varint,8,opt,name=threshold,proto3" json:"threshold,omitempty"`
	Score     uint32              `protobuf:"varint,9,opt,name=score,proto3" json:"score,omitempty"`
	Creator   string              `protobuf:"bytes,10,opt,name=creator,proto3" json:"creator,omitempty"`
	Version   string              `protobuf:"bytes,11,opt,name=version,proto3" json:"version,omitempty"`
	Type      ProposalData_PType  `protobuf:"varint,12,opt,name=type,proto3,enum=bvm.ProposalData_PType" json:"type,omitempty"`
	Completed []byte              `protobuf:"bytes,13,opt,name=completed,proto3" json:"completed,omitempty"`
	Cancel    []byte              `protobuf:"bytes,14,opt,name=cancel,proto3" json:"cancel,omitempty"`
}

func (m *ProposalData) Reset()         { *m = ProposalData{} }
func (m *ProposalData) String() string { return proto.CompactTextString(m) }
func (*ProposalData) ProtoMessage()    {}
func (*ProposalData) Descriptor() ([]byte, []int) {
	return fileDescriptor_c3ac5ce23bf32d05, []int{0}
}
func (m *ProposalData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProposalData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProposalData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProposalData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProposalData.Merge(m, src)
}
func (m *ProposalData) XXX_Size() int {
	return m.Size()
}
func (m *ProposalData) XXX_DiscardUnknown() {
	xxx_messageInfo_ProposalData.DiscardUnknown(m)
}

var xxx_messageInfo_ProposalData proto.InternalMessageInfo

func (m *ProposalData) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ProposalData) GetCode() []byte {
	if m != nil {
		return m.Code
	}
	return nil
}

func (m *ProposalData) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *ProposalData) GetTimeout() int64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *ProposalData) GetStatus() ProposalData_Status {
	if m != nil {
		return m.Status
	}
	return ProposalData_CREATE
}

func (m *ProposalData) GetAssentor() []*VoteInfo {
	if m != nil {
		return m.Assentor
	}
	return nil
}

func (m *ProposalData) GetObjector() []*VoteInfo {
	if m != nil {
		return m.Objector
	}
	return nil
}

func (m *ProposalData) GetThreshold() uint32 {
	if m != nil {
		return m.Threshold
	}
	return 0
}

func (m *ProposalData) GetScore() uint32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *ProposalData) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *ProposalData) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ProposalData) GetType() ProposalData_PType {
	if m != nil {
		return m.Type
	}
	return ProposalData_CONFIG
}

func (m *ProposalData) GetCompleted() []byte {
	if m != nil {
		return m.Completed
	}
	return nil
}

func (m *ProposalData) GetCancel() []byte {
	if m != nil {
		return m.Cancel
	}
	return nil
}

type VoteInfo struct {
	Addr   string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	TxHash string `protobuf:"bytes,2,opt,name=txHash,proto3" json:"txHash,omitempty"`
	Weight uint32 `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (m *VoteInfo) Reset()         { *m = VoteInfo{} }
func (m *VoteInfo) String() string { return proto.CompactTextString(m) }
func (*VoteInfo) ProtoMessage()    {}
func (*VoteInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c3ac5ce23bf32d05, []int{1}
}
func (m *VoteInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *VoteInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_VoteInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *VoteInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteInfo.Merge(m, src)
}
func (m *VoteInfo) XXX_Size() int {
	return m.Size()
}
func (m *VoteInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteInfo.DiscardUnknown(m)
}

var xxx_messageInfo_VoteInfo proto.InternalMessageInfo

func (m *VoteInfo) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *VoteInfo) GetTxHash() string {
	if m != nil {
		return m.TxHash
	}
	return ""
}

func (m *VoteInfo) GetWeight() uint32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func init() {
	proto.RegisterEnum("bvm.ProposalData_Status", ProposalData_Status_name, ProposalData_Status_value)
	proto.RegisterEnum("bvm.ProposalData_PType", ProposalData_PType_name, ProposalData_PType_value)
	proto.RegisterType((*ProposalData)(nil), "bvm.ProposalData")
	proto.RegisterType((*VoteInfo)(nil), "bvm.VoteInfo")
}

func init() { proto.RegisterFile("proposal.proto", fileDescriptor_c3ac5ce23bf32d05) }

var fileDescriptor_c3ac5ce23bf32d05 = []byte{
	// 514 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xce, 0xda, 0x8e, 0x13, 0x4f, 0x7e, 0x58, 0xad, 0x10, 0xec, 0x01, 0x59, 0x56, 0x4e, 0x46,
	0x48, 0x11, 0x2a, 0x4f, 0x10, 0xdc, 0xa5, 0x18, 0xb5, 0x76, 0xb4, 0x31, 0x85, 0x1b, 0x72, 0xec,
	0x85, 0x04, 0x25, 0x59, 0xcb, 0xde, 0x16, 0xfa, 0x16, 0x3c, 0x16, 0xc7, 0x1e, 0x39, 0xa2, 0xf6,
	0x0d, 0x78, 0x02, 0xb4, 0x1b, 0x87, 0x20, 0xc1, 0x6d, 0xbe, 0x9f, 0xb1, 0x67, 0x66, 0x3f, 0x18,
	0x57, 0xb5, 0xac, 0x64, 0x93, 0x6f, 0xa6, 0x55, 0x2d, 0x95, 0x24, 0xf6, 0xf2, 0x7a, 0x3b, 0xf9,
	0xe5, 0xc0, 0x70, 0xde, 0xf2, 0xa7, 0xb9, 0xca, 0xc9, 0x18, 0xac, 0x75, 0x49, 0x51, 0x80, 0x42,
	0x87, 0x5b, 0xeb, 0x92, 0x10, 0x70, 0x0a, 0x59, 0x0a, 0x6a, 0x05, 0x28, 0x1c, 0x72, 0x53, 0x93,
	0x27, 0xe0, 0xa9, 0xf5, 0x56, 0x34, 0x2a, 0xdf, 0x56, 0xd4, 0x0e, 0x50, 0x68, 0xf3, 0x23, 0x41,
	0x28, 0xf4, 0x34, 0x90, 0x57, 0x8a, 0x3a, 0x46, 0x3b, 0x40, 0xf2, 0x1c, 0xdc, 0x46, 0xe5, 0xea,
	0xaa, 0xa1, 0xdd, 0x00, 0x85, 0xe3, 0x13, 0x3a, 0x5d, 0x5e, 0x6f, 0xa7, 0x7f, 0xff, 0x7e, 0xba,
	0x30, 0x3a, 0x6f, 0x7d, 0xe4, 0x29, 0xf4, 0xf3, 0xa6, 0x11, 0x3b, 0x25, 0x6b, 0xea, 0x06, 0x76,
	0x38, 0x38, 0x19, 0x99, 0x9e, 0x4b, 0xa9, 0x44, 0xbc, 0xfb, 0x28, 0xf9, 0x1f, 0x59, 0x5b, 0xe5,
	0xf2, 0xb3, 0x28, 0xb4, 0xb5, 0xf7, 0x5f, 0xeb, 0x41, 0x36, 0xf3, 0xaf, 0x6a, 0xd1, 0xac, 0xe4,
	0xa6, 0xa4, 0xfd, 0x00, 0x85, 0x23, 0x7e, 0x24, 0xc8, 0x43, 0xe8, 0x36, 0x85, 0xac, 0x05, 0xf5,
	0x8c, 0xb2, 0x07, 0x7a, 0xab, 0xa2, 0x16, 0xb9, 0xfe, 0x3a, 0x04, 0x28, 0xf4, 0xf8, 0x01, 0x6a,
	0xe5, 0x5a, 0xd4, 0xcd, 0x5a, 0xee, 0xe8, 0x60, 0xaf, 0xb4, 0x90, 0x3c, 0x03, 0x47, 0xdd, 0x54,
	0x82, 0x0e, 0xcd, 0xb6, 0x8f, 0xff, 0xdd, 0x76, 0x9e, 0xdd, 0x54, 0x82, 0x1b, 0x93, 0x1e, 0xaa,
	0x90, 0xdb, 0x6a, 0x23, 0x94, 0x28, 0xe9, 0xc8, 0x5c, 0xfb, 0x48, 0x90, 0x47, 0xe0, 0x16, 0xf9,
	0xae, 0x10, 0x1b, 0x3a, 0x36, 0x52, 0x8b, 0x26, 0x02, 0xdc, 0xfd, 0xc9, 0x08, 0x80, 0x1b, 0x71,
	0x36, 0xcb, 0x18, 0xee, 0xe8, 0xfa, 0x32, 0xcd, 0xe2, 0xe4, 0x0c, 0x23, 0x5d, 0x73, 0xf6, 0x86,
	0x45, 0x19, 0xb6, 0xc8, 0x03, 0x18, 0xbc, 0x9b, 0xc5, 0x5a, 0xf8, 0xc0, 0xde, 0x33, 0x6c, 0x9b,
	0xa6, 0x59, 0x12, 0xb1, 0x73, 0xec, 0x90, 0x11, 0x78, 0x51, 0x7a, 0x31, 0x3f, 0x67, 0x19, 0x3b,
	0xc5, 0x5d, 0x32, 0x80, 0x5e, 0x16, 0x5f, 0xb0, 0xf4, 0x6d, 0x86, 0xdd, 0xc9, 0x39, 0x74, 0xcd,
	0xac, 0xa6, 0x21, 0x4d, 0x5e, 0xc5, 0x67, 0xb8, 0x43, 0xc6, 0x00, 0x73, 0xc6, 0x2f, 0xe2, 0xc5,
	0x22, 0x4e, 0x13, 0x8c, 0x48, 0x1f, 0x9c, 0x24, 0x3d, 0x65, 0xd8, 0x22, 0x3d, 0xb0, 0xa3, 0x64,
	0x81, 0x6d, 0x32, 0x84, 0x7e, 0x94, 0x26, 0x19, 0x9f, 0x45, 0x19, 0x76, 0x88, 0x0b, 0x56, 0x34,
	0xc3, 0xdd, 0x49, 0x02, 0xfd, 0xc3, 0xab, 0xe8, 0x7c, 0xe5, 0x65, 0x59, 0x9b, 0xc4, 0x79, 0xdc,
	0xd4, 0x7a, 0x59, 0xf5, 0xf5, 0x75, 0xde, 0xac, 0x4c, 0xea, 0x3c, 0xde, 0x22, 0xcd, 0x7f, 0x11,
	0xeb, 0x4f, 0x2b, 0x65, 0x42, 0x37, 0xe2, 0x2d, 0x7a, 0x49, 0xbf, 0xdf, 0xf9, 0xe8, 0xf6, 0xce,
	0x47, 0x3f, 0xef, 0x7c, 0xf4, 0xed, 0xde, 0xef, 0xdc, 0xde, 0xfb, 0x9d, 0x1f, 0xf7, 0x7e, 0x67,
	0xe9, 0x9a, 0xa8, 0xbf, 0xf8, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x10, 0x9b, 0xeb, 0x6f, 0xfc, 0x02,
	0x00, 0x00,
}

func (m *ProposalData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProposalData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintProposal(dAtA, i, uint64(m.Id))
	}
	if len(m.Code) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Code)))
		i += copy(dAtA[i:], m.Code)
	}
	if m.Timestamp != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintProposal(dAtA, i, uint64(m.Timestamp))
	}
	if m.Timeout != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintProposal(dAtA, i, uint64(m.Timeout))
	}
	if m.Status != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintProposal(dAtA, i, uint64(m.Status))
	}
	if len(m.Assentor) > 0 {
		for _, msg := range m.Assentor {
			dAtA[i] = 0x32
			i++
			i = encodeVarintProposal(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Objector) > 0 {
		for _, msg := range m.Objector {
			dAtA[i] = 0x3a
			i++
			i = encodeVarintProposal(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Threshold != 0 {
		dAtA[i] = 0x40
		i++
		i = encodeVarintProposal(dAtA, i, uint64(m.Threshold))
	}
	if m.Score != 0 {
		dAtA[i] = 0x48
		i++
		i = encodeVarintProposal(dAtA, i, uint64(m.Score))
	}
	if len(m.Creator) > 0 {
		dAtA[i] = 0x52
		i++
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Creator)))
		i += copy(dAtA[i:], m.Creator)
	}
	if len(m.Version) > 0 {
		dAtA[i] = 0x5a
		i++
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Version)))
		i += copy(dAtA[i:], m.Version)
	}
	if m.Type != 0 {
		dAtA[i] = 0x60
		i++
		i = encodeVarintProposal(dAtA, i, uint64(m.Type))
	}
	if len(m.Completed) > 0 {
		dAtA[i] = 0x6a
		i++
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Completed)))
		i += copy(dAtA[i:], m.Completed)
	}
	if len(m.Cancel) > 0 {
		dAtA[i] = 0x72
		i++
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Cancel)))
		i += copy(dAtA[i:], m.Cancel)
	}
	return i, nil
}

func (m *VoteInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *VoteInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Addr) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Addr)))
		i += copy(dAtA[i:], m.Addr)
	}
	if len(m.TxHash) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintProposal(dAtA, i, uint64(len(m.TxHash)))
		i += copy(dAtA[i:], m.TxHash)
	}
	if m.Weight != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintProposal(dAtA, i, uint64(m.Weight))
	}
	return i, nil
}

func encodeVarintProposal(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ProposalData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovProposal(uint64(m.Id))
	}
	l = len(m.Code)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if m.Timestamp != 0 {
		n += 1 + sovProposal(uint64(m.Timestamp))
	}
	if m.Timeout != 0 {
		n += 1 + sovProposal(uint64(m.Timeout))
	}
	if m.Status != 0 {
		n += 1 + sovProposal(uint64(m.Status))
	}
	if len(m.Assentor) > 0 {
		for _, e := range m.Assentor {
			l = e.Size()
			n += 1 + l + sovProposal(uint64(l))
		}
	}
	if len(m.Objector) > 0 {
		for _, e := range m.Objector {
			l = e.Size()
			n += 1 + l + sovProposal(uint64(l))
		}
	}
	if m.Threshold != 0 {
		n += 1 + sovProposal(uint64(m.Threshold))
	}
	if m.Score != 0 {
		n += 1 + sovProposal(uint64(m.Score))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovProposal(uint64(m.Type))
	}
	l = len(m.Completed)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Cancel)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	return n
}

func (m *VoteInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Addr)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.TxHash)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if m.Weight != 0 {
		n += 1 + sovProposal(uint64(m.Weight))
	}
	return n
}

func sovProposal(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozProposal(x uint64) (n int) {
	return sovProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProposalData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProposalData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProposalData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Code = append(m.Code[:0], dAtA[iNdEx:postIndex]...)
			if m.Code == nil {
				m.Code = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timeout", wireType)
			}
			m.Timeout = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timeout |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= ProposalData_Status(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Assentor", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Assentor = append(m.Assentor, &VoteInfo{})
			if err := m.Assentor[len(m.Assentor)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Objector", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Objector = append(m.Objector, &VoteInfo{})
			if err := m.Objector[len(m.Objector)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Threshold", wireType)
			}
			m.Threshold = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Threshold |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Score", wireType)
			}
			m.Score = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Score |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= ProposalData_PType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Completed", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Completed = append(m.Completed[:0], dAtA[iNdEx:postIndex]...)
			if m.Completed == nil {
				m.Completed = []byte{}
			}
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cancel", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Cancel = append(m.Cancel[:0], dAtA[iNdEx:postIndex]...)
			if m.Cancel == nil {
				m.Cancel = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthProposal
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthProposal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *VoteInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: VoteInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: VoteInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weight", wireType)
			}
			m.Weight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Weight |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthProposal
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthProposal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipProposal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProposal
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthProposal
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthProposal
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowProposal
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipProposal(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthProposal
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthProposal = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProposal   = fmt.Errorf("proto: integer overflow")
)
