// Code generated by protoc-gen-go. DO NOT EDIT.
// source: monitor.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Stat struct {
	CpuUsageRate            float32  `protobuf:"fixed32,1,opt,name=cpuUsageRate,proto3" json:"cpuUsageRate,omitempty"`
	MemoryUsageRate         float32  `protobuf:"fixed32,2,opt,name=memoryUsageRate,proto3" json:"memoryUsageRate,omitempty"`
	ReadNetworkIOUsageRate  float32  `protobuf:"fixed32,3,opt,name=readNetworkIOUsageRate,proto3" json:"readNetworkIOUsageRate,omitempty"`
	WriteNetworkIOUsageRate float32  `protobuf:"fixed32,4,opt,name=writeNetworkIOUsageRate,proto3" json:"writeNetworkIOUsageRate,omitempty"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_unrecognized        []byte   `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *Stat) Reset()         { *m = Stat{} }
func (m *Stat) String() string { return proto.CompactTextString(m) }
func (*Stat) ProtoMessage()    {}
func (*Stat) Descriptor() ([]byte, []int) {
	return fileDescriptor_44174b7b2a306b71, []int{0}
}

func (m *Stat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stat.Unmarshal(m, b)
}
func (m *Stat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stat.Marshal(b, m, deterministic)
}
func (m *Stat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stat.Merge(m, src)
}
func (m *Stat) XXX_Size() int {
	return xxx_messageInfo_Stat.Size(m)
}
func (m *Stat) XXX_DiscardUnknown() {
	xxx_messageInfo_Stat.DiscardUnknown(m)
}

var xxx_messageInfo_Stat proto.InternalMessageInfo

func (m *Stat) GetCpuUsageRate() float32 {
	if m != nil {
		return m.CpuUsageRate
	}
	return 0
}

func (m *Stat) GetMemoryUsageRate() float32 {
	if m != nil {
		return m.MemoryUsageRate
	}
	return 0
}

func (m *Stat) GetReadNetworkIOUsageRate() float32 {
	if m != nil {
		return m.ReadNetworkIOUsageRate
	}
	return 0
}

func (m *Stat) GetWriteNetworkIOUsageRate() float32 {
	if m != nil {
		return m.WriteNetworkIOUsageRate
	}
	return 0
}

func init() {
	proto.RegisterType((*Stat)(nil), "proto.Stat")
}

func init() { proto.RegisterFile("monitor.proto", fileDescriptor_44174b7b2a306b71) }

var fileDescriptor_44174b7b2a306b71 = []byte{
	// 140 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0xcd, 0xcf, 0xcb,
	0x2c, 0xc9, 0x2f, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0xc7, 0x18,
	0xb9, 0x58, 0x82, 0x4b, 0x12, 0x4b, 0x84, 0x94, 0xb8, 0x78, 0x92, 0x0b, 0x4a, 0x43, 0x8b, 0x13,
	0xd3, 0x53, 0x83, 0x12, 0x4b, 0x52, 0x25, 0x18, 0x15, 0x18, 0x35, 0x98, 0x82, 0x50, 0xc4, 0x84,
	0x34, 0xb8, 0xf8, 0x73, 0x53, 0x73, 0xf3, 0x8b, 0x2a, 0x11, 0xca, 0x98, 0xc0, 0xca, 0xd0, 0x85,
	0x85, 0xcc, 0xb8, 0xc4, 0x8a, 0x52, 0x13, 0x53, 0xfc, 0x52, 0x4b, 0xca, 0xf3, 0x8b, 0xb2, 0x3d,
	0xfd, 0x11, 0x1a, 0x98, 0xc1, 0x1a, 0x70, 0xc8, 0x0a, 0x59, 0x70, 0x89, 0x97, 0x17, 0x65, 0x96,
	0xa4, 0x62, 0xd1, 0xc8, 0x02, 0xd6, 0x88, 0x4b, 0x3a, 0x89, 0x0d, 0xec, 0x1f, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x8d, 0x24, 0xf5, 0xa5, 0xe7, 0x00, 0x00, 0x00,
}