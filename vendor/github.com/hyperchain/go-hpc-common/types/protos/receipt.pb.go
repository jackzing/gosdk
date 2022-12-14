// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: receipt.proto

package protos

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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

type Receipt_STATUS int32

const (
	Receipt_SUCCESS Receipt_STATUS = 0
	Receipt_FAILED  Receipt_STATUS = 1
)

var Receipt_STATUS_name = map[int32]string{
	0: "SUCCESS",
	1: "FAILED",
}

var Receipt_STATUS_value = map[string]int32{
	"SUCCESS": 0,
	"FAILED":  1,
}

func (x Receipt_STATUS) String() string {
	return proto.EnumName(Receipt_STATUS_name, int32(x))
}

func (Receipt_STATUS) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ace1d6eb38fad2c8, []int{0, 0}
}

type Receipt_VmType int32

const (
	Receipt_EVM      Receipt_VmType = 0
	Receipt_JVM      Receipt_VmType = 1
	Receipt_HVM      Receipt_VmType = 2
	Receipt_BVM      Receipt_VmType = 3
	Receipt_TRANSFER Receipt_VmType = 4
	Receipt_KVSQL    Receipt_VmType = 5
	Receipt_FVM      Receipt_VmType = 6
)

var Receipt_VmType_name = map[int32]string{
	0: "EVM",
	1: "JVM",
	2: "HVM",
	3: "BVM",
	4: "TRANSFER",
	5: "KVSQL",
	6: "FVM",
}

var Receipt_VmType_value = map[string]int32{
	"EVM":      0,
	"JVM":      1,
	"HVM":      2,
	"BVM":      3,
	"TRANSFER": 4,
	"KVSQL":    5,
	"FVM":      6,
}

func (x Receipt_VmType) String() string {
	return proto.EnumName(Receipt_VmType_name, int32(x))
}

func (Receipt_VmType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ace1d6eb38fad2c8, []int{0, 1}
}

type Receipt struct {
	Version           []byte         `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Bloom             []byte         `protobuf:"bytes,2,opt,name=Bloom,proto3" json:"Bloom,omitempty"`
	CumulativeGasUsed int64          `protobuf:"varint,3,opt,name=CumulativeGasUsed,proto3" json:"CumulativeGasUsed,omitempty"`
	TxHash            []byte         `protobuf:"bytes,4,opt,name=TxHash,proto3" json:"TxHash,omitempty"`
	ContractAddress   []byte         `protobuf:"bytes,5,opt,name=ContractAddress,proto3" json:"ContractAddress,omitempty"`
	GasUsed           int64          `protobuf:"varint,6,opt,name=GasUsed,proto3" json:"GasUsed,omitempty"`
	Ret               []byte         `protobuf:"bytes,7,opt,name=Ret,proto3" json:"Ret,omitempty"`
	Logs              []byte         `protobuf:"bytes,8,opt,name=Logs,proto3" json:"Logs,omitempty"`
	Status            Receipt_STATUS `protobuf:"varint,9,opt,name=Status,proto3,enum=protos.Receipt_STATUS" json:"Status,omitempty"`
	Message           []byte         `protobuf:"bytes,10,opt,name=Message,proto3" json:"Message,omitempty"`
	VmType            Receipt_VmType `protobuf:"varint,11,opt,name=vmType,proto3,enum=protos.Receipt_VmType" json:"vmType,omitempty"`
	Oracles           []byte         `protobuf:"bytes,12,opt,name=Oracles,proto3" json:"Oracles,omitempty"`
}

func (m *Receipt) Reset()         { *m = Receipt{} }
func (m *Receipt) String() string { return proto.CompactTextString(m) }
func (*Receipt) ProtoMessage()    {}
func (*Receipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_ace1d6eb38fad2c8, []int{0}
}
func (m *Receipt) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Receipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Receipt.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Receipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Receipt.Merge(m, src)
}
func (m *Receipt) XXX_Size() int {
	return m.Size()
}
func (m *Receipt) XXX_DiscardUnknown() {
	xxx_messageInfo_Receipt.DiscardUnknown(m)
}

var xxx_messageInfo_Receipt proto.InternalMessageInfo

func (m *Receipt) GetVersion() []byte {
	if m != nil {
		return m.Version
	}
	return nil
}

func (m *Receipt) GetBloom() []byte {
	if m != nil {
		return m.Bloom
	}
	return nil
}

func (m *Receipt) GetCumulativeGasUsed() int64 {
	if m != nil {
		return m.CumulativeGasUsed
	}
	return 0
}

func (m *Receipt) GetTxHash() []byte {
	if m != nil {
		return m.TxHash
	}
	return nil
}

func (m *Receipt) GetContractAddress() []byte {
	if m != nil {
		return m.ContractAddress
	}
	return nil
}

func (m *Receipt) GetGasUsed() int64 {
	if m != nil {
		return m.GasUsed
	}
	return 0
}

func (m *Receipt) GetRet() []byte {
	if m != nil {
		return m.Ret
	}
	return nil
}

func (m *Receipt) GetLogs() []byte {
	if m != nil {
		return m.Logs
	}
	return nil
}

func (m *Receipt) GetStatus() Receipt_STATUS {
	if m != nil {
		return m.Status
	}
	return Receipt_SUCCESS
}

func (m *Receipt) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *Receipt) GetVmType() Receipt_VmType {
	if m != nil {
		return m.VmType
	}
	return Receipt_EVM
}

func (m *Receipt) GetOracles() []byte {
	if m != nil {
		return m.Oracles
	}
	return nil
}

type ReceiptWrapper struct {
	ReceiptVersion []byte `protobuf:"bytes,1,opt,name=receiptVersion,proto3" json:"receiptVersion,omitempty"`
	Receipt        []byte `protobuf:"bytes,2,opt,name=receipt,proto3" json:"receipt,omitempty"`
}

func (m *ReceiptWrapper) Reset()         { *m = ReceiptWrapper{} }
func (m *ReceiptWrapper) String() string { return proto.CompactTextString(m) }
func (*ReceiptWrapper) ProtoMessage()    {}
func (*ReceiptWrapper) Descriptor() ([]byte, []int) {
	return fileDescriptor_ace1d6eb38fad2c8, []int{1}
}
func (m *ReceiptWrapper) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReceiptWrapper) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReceiptWrapper.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReceiptWrapper) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReceiptWrapper.Merge(m, src)
}
func (m *ReceiptWrapper) XXX_Size() int {
	return m.Size()
}
func (m *ReceiptWrapper) XXX_DiscardUnknown() {
	xxx_messageInfo_ReceiptWrapper.DiscardUnknown(m)
}

var xxx_messageInfo_ReceiptWrapper proto.InternalMessageInfo

func (m *ReceiptWrapper) GetReceiptVersion() []byte {
	if m != nil {
		return m.ReceiptVersion
	}
	return nil
}

func (m *ReceiptWrapper) GetReceipt() []byte {
	if m != nil {
		return m.Receipt
	}
	return nil
}

type BlockReceipt struct {
	Number   uint64            `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	Receipts []*ReceiptWrapper `protobuf:"bytes,2,rep,name=receipts,proto3" json:"receipts,omitempty"`
}

func (m *BlockReceipt) Reset()         { *m = BlockReceipt{} }
func (m *BlockReceipt) String() string { return proto.CompactTextString(m) }
func (*BlockReceipt) ProtoMessage()    {}
func (*BlockReceipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_ace1d6eb38fad2c8, []int{2}
}
func (m *BlockReceipt) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockReceipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockReceipt.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockReceipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockReceipt.Merge(m, src)
}
func (m *BlockReceipt) XXX_Size() int {
	return m.Size()
}
func (m *BlockReceipt) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockReceipt.DiscardUnknown(m)
}

var xxx_messageInfo_BlockReceipt proto.InternalMessageInfo

func (m *BlockReceipt) GetNumber() uint64 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *BlockReceipt) GetReceipts() []*ReceiptWrapper {
	if m != nil {
		return m.Receipts
	}
	return nil
}

func init() {
	proto.RegisterEnum("protos.Receipt_STATUS", Receipt_STATUS_name, Receipt_STATUS_value)
	proto.RegisterEnum("protos.Receipt_VmType", Receipt_VmType_name, Receipt_VmType_value)
	proto.RegisterType((*Receipt)(nil), "protos.Receipt")
	proto.RegisterType((*ReceiptWrapper)(nil), "protos.ReceiptWrapper")
	proto.RegisterType((*BlockReceipt)(nil), "protos.BlockReceipt")
}

func init() { proto.RegisterFile("receipt.proto", fileDescriptor_ace1d6eb38fad2c8) }

var fileDescriptor_ace1d6eb38fad2c8 = []byte{
	// 493 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xc1, 0x8e, 0x12, 0x41,
	0x10, 0x65, 0x18, 0x18, 0xd8, 0x02, 0xb1, 0xed, 0x18, 0xd2, 0xa7, 0x09, 0x72, 0x30, 0x1c, 0x5c,
	0x30, 0xeb, 0x17, 0x00, 0x82, 0xeb, 0xca, 0xac, 0xb1, 0x07, 0xc6, 0x64, 0x6f, 0xc3, 0xd0, 0x01,
	0x22, 0x43, 0x4f, 0xba, 0x1b, 0x22, 0x7f, 0xe1, 0x67, 0x79, 0xdc, 0xa3, 0x47, 0x03, 0x3f, 0xe0,
	0x27, 0x98, 0xee, 0x69, 0x4c, 0xdc, 0xdd, 0x13, 0xf5, 0x5e, 0xbd, 0x7e, 0xaf, 0xa8, 0x29, 0x78,
	0x26, 0x58, 0xc2, 0xd6, 0x99, 0xea, 0x66, 0x82, 0x2b, 0x8e, 0x3d, 0xf3, 0x23, 0xdb, 0x7f, 0x5c,
	0xa8, 0xd0, 0xbc, 0x83, 0x09, 0x54, 0xf6, 0x4c, 0xc8, 0x35, 0xdf, 0x12, 0xa7, 0xe5, 0x74, 0xea,
	0xf4, 0x0c, 0xf1, 0x4b, 0x28, 0x0f, 0x36, 0x9c, 0xa7, 0xa4, 0x68, 0xf8, 0x1c, 0xe0, 0x37, 0xf0,
	0x62, 0xb8, 0x4b, 0x77, 0x9b, 0x58, 0xad, 0xf7, 0xec, 0x43, 0x2c, 0x67, 0x92, 0x2d, 0x88, 0xdb,
	0x72, 0x3a, 0x2e, 0x7d, 0xdc, 0xc0, 0x4d, 0xf0, 0xa6, 0xdf, 0xaf, 0x63, 0xb9, 0x22, 0x25, 0x63,
	0x62, 0x11, 0xee, 0xc0, 0xf3, 0x21, 0xdf, 0x2a, 0x11, 0x27, 0xaa, 0xbf, 0x58, 0x08, 0x26, 0x25,
	0x29, 0x1b, 0xc1, 0x43, 0x5a, 0xcf, 0x77, 0x4e, 0xf1, 0x4c, 0xca, 0x19, 0x62, 0x04, 0x2e, 0x65,
	0x8a, 0x54, 0xcc, 0x3b, 0x5d, 0x62, 0x0c, 0xa5, 0x09, 0x5f, 0x4a, 0x52, 0x35, 0x94, 0xa9, 0x71,
	0x17, 0xbc, 0x50, 0xc5, 0x6a, 0x27, 0xc9, 0x45, 0xcb, 0xe9, 0x34, 0xae, 0x9a, 0xf9, 0x2e, 0x64,
	0xd7, 0x2e, 0xa0, 0x1b, 0x4e, 0xfb, 0xd3, 0x59, 0x48, 0xad, 0x4a, 0xe7, 0x05, 0x4c, 0xca, 0x78,
	0xc9, 0x08, 0xe4, 0xfb, 0xb0, 0x50, 0x3b, 0xed, 0xd3, 0xe9, 0x21, 0x63, 0xa4, 0xf6, 0xb4, 0x53,
	0x64, 0xba, 0xd4, 0xaa, 0xb4, 0xd3, 0x67, 0x11, 0x27, 0x1b, 0x26, 0x49, 0x3d, 0x77, 0xb2, 0xb0,
	0xfd, 0x0a, 0xbc, 0x3c, 0x15, 0xd7, 0xa0, 0x12, 0xce, 0x86, 0xc3, 0x51, 0x18, 0xa2, 0x02, 0x06,
	0xf0, 0xc6, 0xfd, 0x8f, 0x93, 0xd1, 0x7b, 0xe4, 0xb4, 0x6f, 0xc1, 0xcb, 0xed, 0x70, 0x05, 0xdc,
	0x51, 0x14, 0xa0, 0x82, 0x2e, 0x6e, 0xa2, 0x00, 0x39, 0xba, 0xb8, 0x8e, 0x02, 0x54, 0xd4, 0xc5,
	0x20, 0x0a, 0x90, 0x8b, 0xeb, 0x50, 0x9d, 0xd2, 0xfe, 0x6d, 0x38, 0x1e, 0x51, 0x54, 0xc2, 0x17,
	0x50, 0xfe, 0x14, 0x85, 0x5f, 0x26, 0xa8, 0xac, 0x15, 0xe3, 0x28, 0x40, 0x5e, 0x9b, 0x42, 0xc3,
	0x8e, 0xf9, 0x55, 0xc4, 0x59, 0xc6, 0x04, 0x7e, 0x0d, 0x0d, 0x7b, 0x1d, 0xd1, 0x7f, 0xdf, 0xff,
	0x01, 0xab, 0xff, 0x86, 0x65, 0xec, 0x21, 0x9c, 0x61, 0xfb, 0x0e, 0xea, 0x83, 0x0d, 0x4f, 0xbe,
	0x9d, 0x4f, 0xa9, 0x09, 0xde, 0x76, 0x97, 0xce, 0x99, 0x30, 0x4e, 0x25, 0x6a, 0x11, 0xbe, 0x82,
	0xaa, 0x7d, 0x22, 0x49, 0xb1, 0xe5, 0x76, 0x6a, 0x8f, 0x56, 0x67, 0x67, 0xa2, 0xff, 0x74, 0x83,
	0x9b, 0x9f, 0x47, 0xdf, 0xb9, 0x3f, 0xfa, 0xce, 0xef, 0xa3, 0xef, 0xfc, 0x38, 0xf9, 0x85, 0xfb,
	0x93, 0x5f, 0xf8, 0x75, 0xf2, 0x0b, 0x77, 0x6f, 0x97, 0x6b, 0xb5, 0xda, 0xcd, 0xbb, 0x09, 0x4f,
	0x7b, 0xab, 0x43, 0xc6, 0x44, 0xb2, 0x8a, 0xd7, 0xdb, 0xde, 0x92, 0x5f, 0xae, 0xb2, 0xe4, 0x32,
	0xe1, 0x69, 0xca, 0xb7, 0x3d, 0x75, 0xc8, 0x98, 0xec, 0xe5, 0x21, 0xf3, 0xfc, 0xec, 0xdf, 0xfd,
	0x0d, 0x00, 0x00, 0xff, 0xff, 0x86, 0x8c, 0x21, 0xbb, 0x0e, 0x03, 0x00, 0x00,
}

func (m *Receipt) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Receipt) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Receipt) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Oracles) > 0 {
		i -= len(m.Oracles)
		copy(dAtA[i:], m.Oracles)
		i = encodeVarintReceipt(dAtA, i, uint64(len(m.Oracles)))
		i--
		dAtA[i] = 0x62
	}
	if m.VmType != 0 {
		i = encodeVarintReceipt(dAtA, i, uint64(m.VmType))
		i--
		dAtA[i] = 0x58
	}
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintReceipt(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x52
	}
	if m.Status != 0 {
		i = encodeVarintReceipt(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x48
	}
	if len(m.Logs) > 0 {
		i -= len(m.Logs)
		copy(dAtA[i:], m.Logs)
		i = encodeVarintReceipt(dAtA, i, uint64(len(m.Logs)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.Ret) > 0 {
		i -= len(m.Ret)
		copy(dAtA[i:], m.Ret)
		i = encodeVarintReceipt(dAtA, i, uint64(len(m.Ret)))
		i--
		dAtA[i] = 0x3a
	}
	if m.GasUsed != 0 {
		i = encodeVarintReceipt(dAtA, i, uint64(m.GasUsed))
		i--
		dAtA[i] = 0x30
	}
	if len(m.ContractAddress) > 0 {
		i -= len(m.ContractAddress)
		copy(dAtA[i:], m.ContractAddress)
		i = encodeVarintReceipt(dAtA, i, uint64(len(m.ContractAddress)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.TxHash) > 0 {
		i -= len(m.TxHash)
		copy(dAtA[i:], m.TxHash)
		i = encodeVarintReceipt(dAtA, i, uint64(len(m.TxHash)))
		i--
		dAtA[i] = 0x22
	}
	if m.CumulativeGasUsed != 0 {
		i = encodeVarintReceipt(dAtA, i, uint64(m.CumulativeGasUsed))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Bloom) > 0 {
		i -= len(m.Bloom)
		copy(dAtA[i:], m.Bloom)
		i = encodeVarintReceipt(dAtA, i, uint64(len(m.Bloom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Version) > 0 {
		i -= len(m.Version)
		copy(dAtA[i:], m.Version)
		i = encodeVarintReceipt(dAtA, i, uint64(len(m.Version)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ReceiptWrapper) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReceiptWrapper) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ReceiptWrapper) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Receipt) > 0 {
		i -= len(m.Receipt)
		copy(dAtA[i:], m.Receipt)
		i = encodeVarintReceipt(dAtA, i, uint64(len(m.Receipt)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ReceiptVersion) > 0 {
		i -= len(m.ReceiptVersion)
		copy(dAtA[i:], m.ReceiptVersion)
		i = encodeVarintReceipt(dAtA, i, uint64(len(m.ReceiptVersion)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *BlockReceipt) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockReceipt) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockReceipt) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Receipts) > 0 {
		for iNdEx := len(m.Receipts) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Receipts[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintReceipt(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Number != 0 {
		i = encodeVarintReceipt(dAtA, i, uint64(m.Number))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintReceipt(dAtA []byte, offset int, v uint64) int {
	offset -= sovReceipt(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Receipt) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovReceipt(uint64(l))
	}
	l = len(m.Bloom)
	if l > 0 {
		n += 1 + l + sovReceipt(uint64(l))
	}
	if m.CumulativeGasUsed != 0 {
		n += 1 + sovReceipt(uint64(m.CumulativeGasUsed))
	}
	l = len(m.TxHash)
	if l > 0 {
		n += 1 + l + sovReceipt(uint64(l))
	}
	l = len(m.ContractAddress)
	if l > 0 {
		n += 1 + l + sovReceipt(uint64(l))
	}
	if m.GasUsed != 0 {
		n += 1 + sovReceipt(uint64(m.GasUsed))
	}
	l = len(m.Ret)
	if l > 0 {
		n += 1 + l + sovReceipt(uint64(l))
	}
	l = len(m.Logs)
	if l > 0 {
		n += 1 + l + sovReceipt(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovReceipt(uint64(m.Status))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovReceipt(uint64(l))
	}
	if m.VmType != 0 {
		n += 1 + sovReceipt(uint64(m.VmType))
	}
	l = len(m.Oracles)
	if l > 0 {
		n += 1 + l + sovReceipt(uint64(l))
	}
	return n
}

func (m *ReceiptWrapper) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ReceiptVersion)
	if l > 0 {
		n += 1 + l + sovReceipt(uint64(l))
	}
	l = len(m.Receipt)
	if l > 0 {
		n += 1 + l + sovReceipt(uint64(l))
	}
	return n
}

func (m *BlockReceipt) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Number != 0 {
		n += 1 + sovReceipt(uint64(m.Number))
	}
	if len(m.Receipts) > 0 {
		for _, e := range m.Receipts {
			l = e.Size()
			n += 1 + l + sovReceipt(uint64(l))
		}
	}
	return n
}

func sovReceipt(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozReceipt(x uint64) (n int) {
	return sovReceipt(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Receipt) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReceipt
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
			return fmt.Errorf("proto: Receipt: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Receipt: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = append(m.Version[:0], dAtA[iNdEx:postIndex]...)
			if m.Version == nil {
				m.Version = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bloom", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bloom = append(m.Bloom[:0], dAtA[iNdEx:postIndex]...)
			if m.Bloom == nil {
				m.Bloom = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CumulativeGasUsed", wireType)
			}
			m.CumulativeGasUsed = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CumulativeGasUsed |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxHash = append(m.TxHash[:0], dAtA[iNdEx:postIndex]...)
			if m.TxHash == nil {
				m.TxHash = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractAddress = append(m.ContractAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.ContractAddress == nil {
				m.ContractAddress = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasUsed", wireType)
			}
			m.GasUsed = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasUsed |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ret", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ret = append(m.Ret[:0], dAtA[iNdEx:postIndex]...)
			if m.Ret == nil {
				m.Ret = []byte{}
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Logs", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Logs = append(m.Logs[:0], dAtA[iNdEx:postIndex]...)
			if m.Logs == nil {
				m.Logs = []byte{}
			}
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= Receipt_STATUS(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = append(m.Message[:0], dAtA[iNdEx:postIndex]...)
			if m.Message == nil {
				m.Message = []byte{}
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VmType", wireType)
			}
			m.VmType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VmType |= Receipt_VmType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Oracles", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Oracles = append(m.Oracles[:0], dAtA[iNdEx:postIndex]...)
			if m.Oracles == nil {
				m.Oracles = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipReceipt(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthReceipt
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthReceipt
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
func (m *ReceiptWrapper) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReceipt
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
			return fmt.Errorf("proto: ReceiptWrapper: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReceiptWrapper: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceiptVersion", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReceiptVersion = append(m.ReceiptVersion[:0], dAtA[iNdEx:postIndex]...)
			if m.ReceiptVersion == nil {
				m.ReceiptVersion = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receipt", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Receipt = append(m.Receipt[:0], dAtA[iNdEx:postIndex]...)
			if m.Receipt == nil {
				m.Receipt = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipReceipt(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthReceipt
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthReceipt
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
func (m *BlockReceipt) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReceipt
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
			return fmt.Errorf("proto: BlockReceipt: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockReceipt: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Number", wireType)
			}
			m.Number = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Number |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receipts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReceipt
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
				return ErrInvalidLengthReceipt
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthReceipt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Receipts = append(m.Receipts, &ReceiptWrapper{})
			if err := m.Receipts[len(m.Receipts)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipReceipt(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthReceipt
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthReceipt
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
func skipReceipt(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowReceipt
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
					return 0, ErrIntOverflowReceipt
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
					return 0, ErrIntOverflowReceipt
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
				return 0, ErrInvalidLengthReceipt
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthReceipt
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowReceipt
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
				next, err := skipReceipt(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthReceipt
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
	ErrInvalidLengthReceipt = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowReceipt   = fmt.Errorf("proto: integer overflow")
)
