// Code generated by protoc-gen-go. DO NOT EDIT.
// source: health.proto

package idl

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Status int32

const (
	Status_INVALID         Status = 0
	Status_NA              Status = 1
	Status_HEALTHY         Status = 2
	Status_PENDING         Status = 3
	Status_MAYBE_UNHEALTHY Status = 4
	Status_UNHEALTHY       Status = 5
	Status_DYING           Status = 6
	Status_DEAD            Status = 7
)

var Status_name = map[int32]string{
	0: "INVALID",
	1: "NA",
	2: "HEALTHY",
	3: "PENDING",
	4: "MAYBE_UNHEALTHY",
	5: "UNHEALTHY",
	6: "DYING",
	7: "DEAD",
}
var Status_value = map[string]int32{
	"INVALID":         0,
	"NA":              1,
	"HEALTHY":         2,
	"PENDING":         3,
	"MAYBE_UNHEALTHY": 4,
	"UNHEALTHY":       5,
	"DYING":           6,
	"DEAD":            7,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}
func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_health_ff0e4cf87304155b, []int{0}
}

// A value is a measurement unit of an entity's health
type Value struct {
	Status               Status   `protobuf:"varint,1,opt,name=status,proto3,enum=idl.Status" json:"status,omitempty"`
	Score                float32  `protobuf:"fixed32,2,opt,name=score,proto3" json:"score,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_health_ff0e4cf87304155b, []int{0}
}
func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (dst *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(dst, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

func (m *Value) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_INVALID
}

func (m *Value) GetScore() float32 {
	if m != nil {
		return m.Score
	}
	return 0
}

// A metric is a single aspect of an entity's health
type Metric struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value                *Value   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Metric) Reset()         { *m = Metric{} }
func (m *Metric) String() string { return proto.CompactTextString(m) }
func (*Metric) ProtoMessage()    {}
func (*Metric) Descriptor() ([]byte, []int) {
	return fileDescriptor_health_ff0e4cf87304155b, []int{1}
}
func (m *Metric) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metric.Unmarshal(m, b)
}
func (m *Metric) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metric.Marshal(b, m, deterministic)
}
func (dst *Metric) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metric.Merge(dst, src)
}
func (m *Metric) XXX_Size() int {
	return xxx_messageInfo_Metric.Size(m)
}
func (m *Metric) XXX_DiscardUnknown() {
	xxx_messageInfo_Metric.DiscardUnknown(m)
}

var xxx_messageInfo_Metric proto.InternalMessageInfo

func (m *Metric) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Metric) GetValue() *Value {
	if m != nil {
		return m.Value
	}
	return nil
}

// An observation is a collection of a metrics measuring
// an entity's health at a particular time
type Observation struct {
	Ts                   *timestamp.Timestamp `protobuf:"bytes,1,opt,name=ts,proto3" json:"ts,omitempty"`
	Metrics              map[string]*Metric   `protobuf:"bytes,2,rep,name=metrics,proto3" json:"metrics,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Observation) Reset()         { *m = Observation{} }
func (m *Observation) String() string { return proto.CompactTextString(m) }
func (*Observation) ProtoMessage()    {}
func (*Observation) Descriptor() ([]byte, []int) {
	return fileDescriptor_health_ff0e4cf87304155b, []int{2}
}
func (m *Observation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Observation.Unmarshal(m, b)
}
func (m *Observation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Observation.Marshal(b, m, deterministic)
}
func (dst *Observation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Observation.Merge(dst, src)
}
func (m *Observation) XXX_Size() int {
	return xxx_messageInfo_Observation.Size(m)
}
func (m *Observation) XXX_DiscardUnknown() {
	xxx_messageInfo_Observation.DiscardUnknown(m)
}

var xxx_messageInfo_Observation proto.InternalMessageInfo

func (m *Observation) GetTs() *timestamp.Timestamp {
	if m != nil {
		return m.Ts
	}
	return nil
}

func (m *Observation) GetMetrics() map[string]*Metric {
	if m != nil {
		return m.Metrics
	}
	return nil
}

// A report is an observation attached with the observer and the observed (subject)
type Report struct {
	Observer             string       `protobuf:"bytes,1,opt,name=observer,proto3" json:"observer,omitempty"`
	Subject              string       `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
	Observation          *Observation `protobuf:"bytes,3,opt,name=observation,proto3" json:"observation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Report) Reset()         { *m = Report{} }
func (m *Report) String() string { return proto.CompactTextString(m) }
func (*Report) ProtoMessage()    {}
func (*Report) Descriptor() ([]byte, []int) {
	return fileDescriptor_health_ff0e4cf87304155b, []int{3}
}
func (m *Report) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Report.Unmarshal(m, b)
}
func (m *Report) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Report.Marshal(b, m, deterministic)
}
func (dst *Report) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Report.Merge(dst, src)
}
func (m *Report) XXX_Size() int {
	return xxx_messageInfo_Report.Size(m)
}
func (m *Report) XXX_DiscardUnknown() {
	xxx_messageInfo_Report.DiscardUnknown(m)
}

var xxx_messageInfo_Report proto.InternalMessageInfo

func (m *Report) GetObserver() string {
	if m != nil {
		return m.Observer
	}
	return ""
}

func (m *Report) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Report) GetObservation() *Observation {
	if m != nil {
		return m.Observation
	}
	return nil
}

// A view is a continuous stream of reports made by an observer for a subject
type View struct {
	Observer             string         `protobuf:"bytes,1,opt,name=observer,proto3" json:"observer,omitempty"`
	Subject              string         `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
	Observations         []*Observation `protobuf:"bytes,3,rep,name=observations,proto3" json:"observations,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *View) Reset()         { *m = View{} }
func (m *View) String() string { return proto.CompactTextString(m) }
func (*View) ProtoMessage()    {}
func (*View) Descriptor() ([]byte, []int) {
	return fileDescriptor_health_ff0e4cf87304155b, []int{4}
}
func (m *View) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_View.Unmarshal(m, b)
}
func (m *View) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_View.Marshal(b, m, deterministic)
}
func (dst *View) XXX_Merge(src proto.Message) {
	xxx_messageInfo_View.Merge(dst, src)
}
func (m *View) XXX_Size() int {
	return xxx_messageInfo_View.Size(m)
}
func (m *View) XXX_DiscardUnknown() {
	xxx_messageInfo_View.DiscardUnknown(m)
}

var xxx_messageInfo_View proto.InternalMessageInfo

func (m *View) GetObserver() string {
	if m != nil {
		return m.Observer
	}
	return ""
}

func (m *View) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *View) GetObservations() []*Observation {
	if m != nil {
		return m.Observations
	}
	return nil
}

// A panorama is a collection of views from different observers about the same subject
type Panorama struct {
	Subject              string           `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Views                map[string]*View `protobuf:"bytes,2,rep,name=views,proto3" json:"views,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Panorama) Reset()         { *m = Panorama{} }
func (m *Panorama) String() string { return proto.CompactTextString(m) }
func (*Panorama) ProtoMessage()    {}
func (*Panorama) Descriptor() ([]byte, []int) {
	return fileDescriptor_health_ff0e4cf87304155b, []int{5}
}
func (m *Panorama) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Panorama.Unmarshal(m, b)
}
func (m *Panorama) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Panorama.Marshal(b, m, deterministic)
}
func (dst *Panorama) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Panorama.Merge(dst, src)
}
func (m *Panorama) XXX_Size() int {
	return xxx_messageInfo_Panorama.Size(m)
}
func (m *Panorama) XXX_DiscardUnknown() {
	xxx_messageInfo_Panorama.DiscardUnknown(m)
}

var xxx_messageInfo_Panorama proto.InternalMessageInfo

func (m *Panorama) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Panorama) GetViews() map[string]*View {
	if m != nil {
		return m.Views
	}
	return nil
}

// An inference is a merged result of different views about the same subject
type Inference struct {
	Subject              string       `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Observers            []string     `protobuf:"bytes,2,rep,name=observers,proto3" json:"observers,omitempty"`
	Observation          *Observation `protobuf:"bytes,3,opt,name=observation,proto3" json:"observation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Inference) Reset()         { *m = Inference{} }
func (m *Inference) String() string { return proto.CompactTextString(m) }
func (*Inference) ProtoMessage()    {}
func (*Inference) Descriptor() ([]byte, []int) {
	return fileDescriptor_health_ff0e4cf87304155b, []int{6}
}
func (m *Inference) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Inference.Unmarshal(m, b)
}
func (m *Inference) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Inference.Marshal(b, m, deterministic)
}
func (dst *Inference) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Inference.Merge(dst, src)
}
func (m *Inference) XXX_Size() int {
	return xxx_messageInfo_Inference.Size(m)
}
func (m *Inference) XXX_DiscardUnknown() {
	xxx_messageInfo_Inference.DiscardUnknown(m)
}

var xxx_messageInfo_Inference proto.InternalMessageInfo

func (m *Inference) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Inference) GetObservers() []string {
	if m != nil {
		return m.Observers
	}
	return nil
}

func (m *Inference) GetObservation() *Observation {
	if m != nil {
		return m.Observation
	}
	return nil
}

func init() {
	proto.RegisterType((*Value)(nil), "idl.Value")
	proto.RegisterType((*Metric)(nil), "idl.Metric")
	proto.RegisterType((*Observation)(nil), "idl.Observation")
	proto.RegisterMapType((map[string]*Metric)(nil), "idl.Observation.MetricsEntry")
	proto.RegisterType((*Report)(nil), "idl.Report")
	proto.RegisterType((*View)(nil), "idl.View")
	proto.RegisterType((*Panorama)(nil), "idl.Panorama")
	proto.RegisterMapType((map[string]*View)(nil), "idl.Panorama.ViewsEntry")
	proto.RegisterType((*Inference)(nil), "idl.Inference")
	proto.RegisterEnum("idl.Status", Status_name, Status_value)
}

func init() { proto.RegisterFile("health.proto", fileDescriptor_health_ff0e4cf87304155b) }

var fileDescriptor_health_ff0e4cf87304155b = []byte{
	// 521 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xc5, 0x76, 0xec, 0xd4, 0xe3, 0x00, 0xd6, 0xc2, 0xc1, 0xb2, 0x40, 0x0d, 0xe6, 0x12, 0xf5,
	0xb0, 0x95, 0x02, 0x12, 0x88, 0x03, 0x52, 0x42, 0xa2, 0x36, 0x52, 0x1b, 0x2a, 0x53, 0x22, 0xe5,
	0x84, 0x9c, 0x64, 0xda, 0xb8, 0x24, 0xde, 0x68, 0x77, 0x9d, 0xa8, 0x9f, 0xc2, 0xb7, 0xf0, 0x73,
	0xc8, 0xbb, 0x36, 0x71, 0x00, 0x71, 0xe8, 0xcd, 0x33, 0xfb, 0xf6, 0xbd, 0x37, 0xfb, 0x3c, 0xd0,
	0x5a, 0x62, 0xb2, 0x92, 0x4b, 0xba, 0xe1, 0x4c, 0x32, 0x62, 0xa5, 0x8b, 0x55, 0x78, 0x7c, 0xcb,
	0xd8, 0xed, 0x0a, 0x4f, 0x55, 0x6b, 0x96, 0xdf, 0x9c, 0xca, 0x74, 0x8d, 0x42, 0x26, 0xeb, 0x8d,
	0x46, 0x45, 0x7d, 0xb0, 0x27, 0xc9, 0x2a, 0x47, 0xf2, 0x1a, 0x1c, 0x21, 0x13, 0x99, 0x8b, 0xc0,
	0x68, 0x1b, 0x9d, 0x27, 0x5d, 0x8f, 0xa6, 0x8b, 0x15, 0xfd, 0xa2, 0x5a, 0x71, 0x79, 0x44, 0x9e,
	0x83, 0x2d, 0xe6, 0x8c, 0x63, 0x60, 0xb6, 0x8d, 0x8e, 0x19, 0xeb, 0x22, 0xfa, 0x08, 0xce, 0x25,
	0x4a, 0x9e, 0xce, 0x09, 0x81, 0x46, 0x96, 0xac, 0x51, 0x51, 0xb8, 0xb1, 0xfa, 0x26, 0x6d, 0xb0,
	0xb7, 0x85, 0x82, 0xba, 0xe3, 0x75, 0x41, 0xf1, 0x2a, 0xcd, 0x58, 0x1f, 0x44, 0x3f, 0x0d, 0xf0,
	0x3e, 0xcf, 0x04, 0xf2, 0x6d, 0x22, 0x53, 0x96, 0x91, 0x13, 0x30, 0xa5, 0xb6, 0xe1, 0x75, 0x43,
	0xaa, 0x27, 0xa0, 0xd5, 0x04, 0xf4, 0xba, 0x9a, 0x20, 0x36, 0xa5, 0x20, 0xef, 0xa0, 0xb9, 0x56,
	0xda, 0x22, 0x30, 0xdb, 0x56, 0xc7, 0xeb, 0xbe, 0x54, 0xfc, 0x35, 0x3a, 0xaa, 0xbd, 0x89, 0x61,
	0x26, 0xf9, 0x7d, 0x5c, 0xa1, 0xc3, 0x33, 0x68, 0xd5, 0x0f, 0x88, 0x0f, 0xd6, 0x77, 0xbc, 0x2f,
	0x9d, 0x17, 0x9f, 0xe4, 0xd5, 0xa1, 0x71, 0xfd, 0x20, 0xfa, 0x4e, 0xe9, 0xfc, 0x83, 0xf9, 0xde,
	0x88, 0x38, 0x38, 0x31, 0x6e, 0x18, 0x97, 0x24, 0x84, 0x23, 0xa6, 0x74, 0x91, 0x97, 0x3c, 0xbf,
	0x6b, 0x12, 0x40, 0x53, 0xe4, 0xb3, 0x3b, 0x9c, 0x4b, 0x45, 0xe7, 0xc6, 0x55, 0x49, 0xba, 0xe0,
	0xb1, 0xbd, 0xdb, 0xc0, 0x52, 0x62, 0xfe, 0x9f, 0x53, 0xc4, 0x75, 0x50, 0xc4, 0xa1, 0x31, 0x49,
	0x71, 0xf7, 0x40, 0xc5, 0xb7, 0xd0, 0xaa, 0x91, 0x89, 0xc0, 0x52, 0x0f, 0xf7, 0xb7, 0xe4, 0x01,
	0x2a, 0xfa, 0x61, 0xc0, 0xd1, 0x55, 0x92, 0x31, 0x9e, 0xac, 0x93, 0x3a, 0xb9, 0x71, 0x48, 0x4e,
	0xc1, 0xde, 0xa6, 0xb8, 0xab, 0xe2, 0x08, 0x14, 0x6b, 0x75, 0x8f, 0x16, 0xae, 0xcb, 0x24, 0x34,
	0x2c, 0xfc, 0x04, 0xb0, 0x6f, 0xfe, 0x23, 0x85, 0xe3, 0xc3, 0x14, 0x5c, 0xfd, 0xfb, 0xa4, 0xb8,
	0xab, 0x67, 0xb0, 0x03, 0x77, 0x94, 0xdd, 0x20, 0xc7, 0x6c, 0x8e, 0xff, 0xf1, 0xf6, 0x02, 0xdc,
	0xea, 0x79, 0xb4, 0x3f, 0x37, 0xde, 0x37, 0x1e, 0x12, 0xc4, 0xc9, 0x06, 0x1c, 0xbd, 0x22, 0xc4,
	0x83, 0xe6, 0x68, 0x3c, 0xe9, 0x5d, 0x8c, 0x06, 0xfe, 0x23, 0xe2, 0x80, 0x39, 0xee, 0xf9, 0x46,
	0xd1, 0x3c, 0x1f, 0xf6, 0x2e, 0xae, 0xcf, 0xa7, 0xbe, 0x59, 0x14, 0x57, 0xc3, 0xf1, 0x60, 0x34,
	0x3e, 0xf3, 0x2d, 0xf2, 0x0c, 0x9e, 0x5e, 0xf6, 0xa6, 0xfd, 0xe1, 0xb7, 0xaf, 0xe3, 0x0a, 0xd1,
	0x20, 0x8f, 0xc1, 0xdd, 0x97, 0x36, 0x71, 0xc1, 0x1e, 0x4c, 0x0b, 0xb8, 0x43, 0x8e, 0xa0, 0x31,
	0x18, 0xf6, 0x06, 0x7e, 0xb3, 0x1f, 0x42, 0x80, 0x8b, 0x9c, 0xde, 0x2d, 0x73, 0xca, 0xf8, 0x02,
	0x39, 0x5d, 0x20, 0x6e, 0xf4, 0xe2, 0xcf, 0x1c, 0xb5, 0x24, 0x6f, 0x7e, 0x05, 0x00, 0x00, 0xff,
	0xff, 0x5b, 0x37, 0x55, 0x67, 0x09, 0x04, 0x00, 0x00,
}
