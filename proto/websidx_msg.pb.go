// Code generated by protoc-gen-go.
// source: websidx_msg.proto
// DO NOT EDIT!

package websidx_interface

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type DocumentMsg struct {
	AdType           *uint32 `protobuf:"varint,1,req,name=ad_type,json=adType" json:"ad_type,omitempty"`
	Adcontent        []byte  `protobuf:"bytes,2,req,name=adcontent" json:"adcontent,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *DocumentMsg) Reset()                    { *m = DocumentMsg{} }
func (m *DocumentMsg) String() string            { return proto.CompactTextString(m) }
func (*DocumentMsg) ProtoMessage()               {}
func (*DocumentMsg) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *DocumentMsg) GetAdType() uint32 {
	if m != nil && m.AdType != nil {
		return *m.AdType
	}
	return 0
}

func (m *DocumentMsg) GetAdcontent() []byte {
	if m != nil {
		return m.Adcontent
	}
	return nil
}

func init() {
	proto.RegisterType((*DocumentMsg)(nil), "websidx_interface.DocumentMsg")
}

func init() { proto.RegisterFile("websidx_msg.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 113 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x4f, 0x4d, 0x2a,
	0xce, 0x4c, 0xa9, 0x88, 0xcf, 0x2d, 0x4e, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x82, 0x0b,
	0x65, 0xe6, 0x95, 0xa4, 0x16, 0xa5, 0x25, 0x26, 0xa7, 0x2a, 0xb9, 0x70, 0x71, 0xbb, 0xe4, 0x27,
	0x97, 0xe6, 0xa6, 0xe6, 0x95, 0xf8, 0x16, 0xa7, 0x0b, 0x89, 0x73, 0xb1, 0x27, 0xa6, 0xc4, 0x97,
	0x54, 0x16, 0xa4, 0x4a, 0x30, 0x2a, 0x30, 0x69, 0xf0, 0x06, 0xb1, 0x25, 0xa6, 0x84, 0x00, 0x79,
	0x42, 0x32, 0x5c, 0x9c, 0x89, 0x29, 0xc9, 0xf9, 0x40, 0x7d, 0x79, 0x25, 0x12, 0x4c, 0x40, 0x29,
	0x9e, 0x20, 0x84, 0x00, 0x20, 0x00, 0x00, 0xff, 0xff, 0xe6, 0x09, 0x02, 0x8f, 0x6c, 0x00, 0x00,
	0x00,
}