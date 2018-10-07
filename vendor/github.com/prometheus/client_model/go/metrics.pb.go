


package io_prometheus_client 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type MetricType int32

const (
	MetricType_COUNTER   MetricType = 0
	MetricType_GAUGE     MetricType = 1
	MetricType_SUMMARY   MetricType = 2
	MetricType_UNTYPED   MetricType = 3
	MetricType_HISTOGRAM MetricType = 4
)

var MetricType_name = map[int32]string{
	0: "COUNTER",
	1: "GAUGE",
	2: "SUMMARY",
	3: "UNTYPED",
	4: "HISTOGRAM",
}
var MetricType_value = map[string]int32{
	"COUNTER":   0,
	"GAUGE":     1,
	"SUMMARY":   2,
	"UNTYPED":   3,
	"HISTOGRAM": 4,
}

func (x MetricType) Enum() *MetricType {
	p := new(MetricType)
	*p = x
	return p
}
func (x MetricType) String() string {
	return proto.EnumName(MetricType_name, int32(x))
}
func (x *MetricType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MetricType_value, data, "MetricType")
	if err != nil {
		return err
	}
	*x = MetricType(value)
	return nil
}
func (MetricType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{0}
}

type LabelPair struct {
	Name                 *string  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Value                *string  `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LabelPair) Reset()         { *m = LabelPair{} }
func (m *LabelPair) String() string { return proto.CompactTextString(m) }
func (*LabelPair) ProtoMessage()    {}
func (*LabelPair) Descriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{0}
}
func (m *LabelPair) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LabelPair.Unmarshal(m, b)
}
func (m *LabelPair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LabelPair.Marshal(b, m, deterministic)
}
func (dst *LabelPair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LabelPair.Merge(dst, src)
}
func (m *LabelPair) XXX_Size() int {
	return xxx_messageInfo_LabelPair.Size(m)
}
func (m *LabelPair) XXX_DiscardUnknown() {
	xxx_messageInfo_LabelPair.DiscardUnknown(m)
}

var xxx_messageInfo_LabelPair proto.InternalMessageInfo

func (m *LabelPair) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *LabelPair) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

type Gauge struct {
	Value                *float64 `protobuf:"fixed64,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Gauge) Reset()         { *m = Gauge{} }
func (m *Gauge) String() string { return proto.CompactTextString(m) }
func (*Gauge) ProtoMessage()    {}
func (*Gauge) Descriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{1}
}
func (m *Gauge) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Gauge.Unmarshal(m, b)
}
func (m *Gauge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Gauge.Marshal(b, m, deterministic)
}
func (dst *Gauge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Gauge.Merge(dst, src)
}
func (m *Gauge) XXX_Size() int {
	return xxx_messageInfo_Gauge.Size(m)
}
func (m *Gauge) XXX_DiscardUnknown() {
	xxx_messageInfo_Gauge.DiscardUnknown(m)
}

var xxx_messageInfo_Gauge proto.InternalMessageInfo

func (m *Gauge) GetValue() float64 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

type Counter struct {
	Value                *float64 `protobuf:"fixed64,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Counter) Reset()         { *m = Counter{} }
func (m *Counter) String() string { return proto.CompactTextString(m) }
func (*Counter) ProtoMessage()    {}
func (*Counter) Descriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{2}
}
func (m *Counter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Counter.Unmarshal(m, b)
}
func (m *Counter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Counter.Marshal(b, m, deterministic)
}
func (dst *Counter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Counter.Merge(dst, src)
}
func (m *Counter) XXX_Size() int {
	return xxx_messageInfo_Counter.Size(m)
}
func (m *Counter) XXX_DiscardUnknown() {
	xxx_messageInfo_Counter.DiscardUnknown(m)
}

var xxx_messageInfo_Counter proto.InternalMessageInfo

func (m *Counter) GetValue() float64 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

type Quantile struct {
	Quantile             *float64 `protobuf:"fixed64,1,opt,name=quantile" json:"quantile,omitempty"`
	Value                *float64 `protobuf:"fixed64,2,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Quantile) Reset()         { *m = Quantile{} }
func (m *Quantile) String() string { return proto.CompactTextString(m) }
func (*Quantile) ProtoMessage()    {}
func (*Quantile) Descriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{3}
}
func (m *Quantile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Quantile.Unmarshal(m, b)
}
func (m *Quantile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Quantile.Marshal(b, m, deterministic)
}
func (dst *Quantile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Quantile.Merge(dst, src)
}
func (m *Quantile) XXX_Size() int {
	return xxx_messageInfo_Quantile.Size(m)
}
func (m *Quantile) XXX_DiscardUnknown() {
	xxx_messageInfo_Quantile.DiscardUnknown(m)
}

var xxx_messageInfo_Quantile proto.InternalMessageInfo

func (m *Quantile) GetQuantile() float64 {
	if m != nil && m.Quantile != nil {
		return *m.Quantile
	}
	return 0
}

func (m *Quantile) GetValue() float64 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

type Summary struct {
	SampleCount          *uint64     `protobuf:"varint,1,opt,name=sample_count,json=sampleCount" json:"sample_count,omitempty"`
	SampleSum            *float64    `protobuf:"fixed64,2,opt,name=sample_sum,json=sampleSum" json:"sample_sum,omitempty"`
	Quantile             []*Quantile `protobuf:"bytes,3,rep,name=quantile" json:"quantile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Summary) Reset()         { *m = Summary{} }
func (m *Summary) String() string { return proto.CompactTextString(m) }
func (*Summary) ProtoMessage()    {}
func (*Summary) Descriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{4}
}
func (m *Summary) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Summary.Unmarshal(m, b)
}
func (m *Summary) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Summary.Marshal(b, m, deterministic)
}
func (dst *Summary) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Summary.Merge(dst, src)
}
func (m *Summary) XXX_Size() int {
	return xxx_messageInfo_Summary.Size(m)
}
func (m *Summary) XXX_DiscardUnknown() {
	xxx_messageInfo_Summary.DiscardUnknown(m)
}

var xxx_messageInfo_Summary proto.InternalMessageInfo

func (m *Summary) GetSampleCount() uint64 {
	if m != nil && m.SampleCount != nil {
		return *m.SampleCount
	}
	return 0
}

func (m *Summary) GetSampleSum() float64 {
	if m != nil && m.SampleSum != nil {
		return *m.SampleSum
	}
	return 0
}

func (m *Summary) GetQuantile() []*Quantile {
	if m != nil {
		return m.Quantile
	}
	return nil
}

type Untyped struct {
	Value                *float64 `protobuf:"fixed64,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Untyped) Reset()         { *m = Untyped{} }
func (m *Untyped) String() string { return proto.CompactTextString(m) }
func (*Untyped) ProtoMessage()    {}
func (*Untyped) Descriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{5}
}
func (m *Untyped) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Untyped.Unmarshal(m, b)
}
func (m *Untyped) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Untyped.Marshal(b, m, deterministic)
}
func (dst *Untyped) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Untyped.Merge(dst, src)
}
func (m *Untyped) XXX_Size() int {
	return xxx_messageInfo_Untyped.Size(m)
}
func (m *Untyped) XXX_DiscardUnknown() {
	xxx_messageInfo_Untyped.DiscardUnknown(m)
}

var xxx_messageInfo_Untyped proto.InternalMessageInfo

func (m *Untyped) GetValue() float64 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

type Histogram struct {
	SampleCount          *uint64   `protobuf:"varint,1,opt,name=sample_count,json=sampleCount" json:"sample_count,omitempty"`
	SampleSum            *float64  `protobuf:"fixed64,2,opt,name=sample_sum,json=sampleSum" json:"sample_sum,omitempty"`
	Bucket               []*Bucket `protobuf:"bytes,3,rep,name=bucket" json:"bucket,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Histogram) Reset()         { *m = Histogram{} }
func (m *Histogram) String() string { return proto.CompactTextString(m) }
func (*Histogram) ProtoMessage()    {}
func (*Histogram) Descriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{6}
}
func (m *Histogram) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Histogram.Unmarshal(m, b)
}
func (m *Histogram) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Histogram.Marshal(b, m, deterministic)
}
func (dst *Histogram) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Histogram.Merge(dst, src)
}
func (m *Histogram) XXX_Size() int {
	return xxx_messageInfo_Histogram.Size(m)
}
func (m *Histogram) XXX_DiscardUnknown() {
	xxx_messageInfo_Histogram.DiscardUnknown(m)
}

var xxx_messageInfo_Histogram proto.InternalMessageInfo

func (m *Histogram) GetSampleCount() uint64 {
	if m != nil && m.SampleCount != nil {
		return *m.SampleCount
	}
	return 0
}

func (m *Histogram) GetSampleSum() float64 {
	if m != nil && m.SampleSum != nil {
		return *m.SampleSum
	}
	return 0
}

func (m *Histogram) GetBucket() []*Bucket {
	if m != nil {
		return m.Bucket
	}
	return nil
}

type Bucket struct {
	CumulativeCount      *uint64  `protobuf:"varint,1,opt,name=cumulative_count,json=cumulativeCount" json:"cumulative_count,omitempty"`
	UpperBound           *float64 `protobuf:"fixed64,2,opt,name=upper_bound,json=upperBound" json:"upper_bound,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bucket) Reset()         { *m = Bucket{} }
func (m *Bucket) String() string { return proto.CompactTextString(m) }
func (*Bucket) ProtoMessage()    {}
func (*Bucket) Descriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{7}
}
func (m *Bucket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bucket.Unmarshal(m, b)
}
func (m *Bucket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bucket.Marshal(b, m, deterministic)
}
func (dst *Bucket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bucket.Merge(dst, src)
}
func (m *Bucket) XXX_Size() int {
	return xxx_messageInfo_Bucket.Size(m)
}
func (m *Bucket) XXX_DiscardUnknown() {
	xxx_messageInfo_Bucket.DiscardUnknown(m)
}

var xxx_messageInfo_Bucket proto.InternalMessageInfo

func (m *Bucket) GetCumulativeCount() uint64 {
	if m != nil && m.CumulativeCount != nil {
		return *m.CumulativeCount
	}
	return 0
}

func (m *Bucket) GetUpperBound() float64 {
	if m != nil && m.UpperBound != nil {
		return *m.UpperBound
	}
	return 0
}

type Metric struct {
	Label                []*LabelPair `protobuf:"bytes,1,rep,name=label" json:"label,omitempty"`
	Gauge                *Gauge       `protobuf:"bytes,2,opt,name=gauge" json:"gauge,omitempty"`
	Counter              *Counter     `protobuf:"bytes,3,opt,name=counter" json:"counter,omitempty"`
	Summary              *Summary     `protobuf:"bytes,4,opt,name=summary" json:"summary,omitempty"`
	Untyped              *Untyped     `protobuf:"bytes,5,opt,name=untyped" json:"untyped,omitempty"`
	Histogram            *Histogram   `protobuf:"bytes,7,opt,name=histogram" json:"histogram,omitempty"`
	TimestampMs          *int64       `protobuf:"varint,6,opt,name=timestamp_ms,json=timestampMs" json:"timestamp_ms,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Metric) Reset()         { *m = Metric{} }
func (m *Metric) String() string { return proto.CompactTextString(m) }
func (*Metric) ProtoMessage()    {}
func (*Metric) Descriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{8}
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

func (m *Metric) GetLabel() []*LabelPair {
	if m != nil {
		return m.Label
	}
	return nil
}

func (m *Metric) GetGauge() *Gauge {
	if m != nil {
		return m.Gauge
	}
	return nil
}

func (m *Metric) GetCounter() *Counter {
	if m != nil {
		return m.Counter
	}
	return nil
}

func (m *Metric) GetSummary() *Summary {
	if m != nil {
		return m.Summary
	}
	return nil
}

func (m *Metric) GetUntyped() *Untyped {
	if m != nil {
		return m.Untyped
	}
	return nil
}

func (m *Metric) GetHistogram() *Histogram {
	if m != nil {
		return m.Histogram
	}
	return nil
}

func (m *Metric) GetTimestampMs() int64 {
	if m != nil && m.TimestampMs != nil {
		return *m.TimestampMs
	}
	return 0
}

type MetricFamily struct {
	Name                 *string     `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Help                 *string     `protobuf:"bytes,2,opt,name=help" json:"help,omitempty"`
	Type                 *MetricType `protobuf:"varint,3,opt,name=type,enum=io.prometheus.client.MetricType" json:"type,omitempty"`
	Metric               []*Metric   `protobuf:"bytes,4,rep,name=metric" json:"metric,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *MetricFamily) Reset()         { *m = MetricFamily{} }
func (m *MetricFamily) String() string { return proto.CompactTextString(m) }
func (*MetricFamily) ProtoMessage()    {}
func (*MetricFamily) Descriptor() ([]byte, []int) {
	return fileDescriptor_metrics_c97c9a2b9560cb8f, []int{9}
}
func (m *MetricFamily) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetricFamily.Unmarshal(m, b)
}
func (m *MetricFamily) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetricFamily.Marshal(b, m, deterministic)
}
func (dst *MetricFamily) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetricFamily.Merge(dst, src)
}
func (m *MetricFamily) XXX_Size() int {
	return xxx_messageInfo_MetricFamily.Size(m)
}
func (m *MetricFamily) XXX_DiscardUnknown() {
	xxx_messageInfo_MetricFamily.DiscardUnknown(m)
}

var xxx_messageInfo_MetricFamily proto.InternalMessageInfo

func (m *MetricFamily) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *MetricFamily) GetHelp() string {
	if m != nil && m.Help != nil {
		return *m.Help
	}
	return ""
}

func (m *MetricFamily) GetType() MetricType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return MetricType_COUNTER
}

func (m *MetricFamily) GetMetric() []*Metric {
	if m != nil {
		return m.Metric
	}
	return nil
}

func init() {
	proto.RegisterType((*LabelPair)(nil), "io.prometheus.client.LabelPair")
	proto.RegisterType((*Gauge)(nil), "io.prometheus.client.Gauge")
	proto.RegisterType((*Counter)(nil), "io.prometheus.client.Counter")
	proto.RegisterType((*Quantile)(nil), "io.prometheus.client.Quantile")
	proto.RegisterType((*Summary)(nil), "io.prometheus.client.Summary")
	proto.RegisterType((*Untyped)(nil), "io.prometheus.client.Untyped")
	proto.RegisterType((*Histogram)(nil), "io.prometheus.client.Histogram")
	proto.RegisterType((*Bucket)(nil), "io.prometheus.client.Bucket")
	proto.RegisterType((*Metric)(nil), "io.prometheus.client.Metric")
	proto.RegisterType((*MetricFamily)(nil), "io.prometheus.client.MetricFamily")
	proto.RegisterEnum("io.prometheus.client.MetricType", MetricType_name, MetricType_value)
}

func init() { proto.RegisterFile("metrics.proto", fileDescriptor_metrics_c97c9a2b9560cb8f) }

var fileDescriptor_metrics_c97c9a2b9560cb8f = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x4f, 0x4f, 0xdb, 0x4e,
	0x14, 0xfc, 0x99, 0xd8, 0x09, 0x7e, 0x86, 0x5f, 0xad, 0x15, 0x07, 0xab, 0x2d, 0x25, 0xcd, 0x89,
	0xf6, 0x10, 0x54, 0x04, 0xaa, 0x44, 0xdb, 0x03, 0x50, 0x1a, 0x2a, 0xd5, 0x40, 0x37, 0xc9, 0x81,
	0x5e, 0xac, 0x8d, 0x59, 0x25, 0x56, 0xbd, 0xb6, 0x6b, 0xef, 0x22, 0xe5, 0xdc, 0x43, 0xbf, 0x47,
	0xbf, 0x68, 0xab, 0xfd, 0xe3, 0x18, 0x24, 0xc3, 0xa9, 0xb7, 0xb7, 0xf3, 0x66, 0xde, 0x8e, 0x77,
	0xc7, 0x0b, 0x9b, 0x8c, 0xf2, 0x32, 0x89, 0xab, 0x61, 0x51, 0xe6, 0x3c, 0x47, 0x5b, 0x49, 0x2e,
	0x2b, 0x46, 0xf9, 0x82, 0x8a, 0x6a, 0x18, 0xa7, 0x09, 0xcd, 0xf8, 0xe0, 0x10, 0xdc, 0x2f, 0x64,
	0x46, 0xd3, 0x2b, 0x92, 0x94, 0x08, 0x81, 0x9d, 0x11, 0x46, 0x03, 0xab, 0x6f, 0xed, 0xba, 0x58,
	0xd5, 0x68, 0x0b, 0x9c, 0x5b, 0x92, 0x0a, 0x1a, 0xac, 0x29, 0x50, 0x2f, 0x06, 0xdb, 0xe0, 0x8c,
	0x88, 0x98, 0xdf, 0x69, 0x4b, 0x8d, 0x55, 0xb7, 0x77, 0xa0, 0x77, 0x9a, 0x8b, 0x8c, 0xd3, 0xf2,
	0x01, 0xc2, 0x7b, 0x58, 0xff, 0x2a, 0x48, 0xc6, 0x93, 0x94, 0xa2, 0xa7, 0xb0, 0xfe, 0xc3, 0xd4,
	0x86, 0xb4, 0x5a, 0xdf, 0xdf, 0x7d, 0xa5, 0xfe, 0x65, 0x41, 0x6f, 0x2c, 0x18, 0x23, 0xe5, 0x12,
	0xbd, 0x84, 0x8d, 0x8a, 0xb0, 0x22, 0xa5, 0x51, 0x2c, 0x77, 0x54, 0x13, 0x6c, 0xec, 0x69, 0x4c,
	0x99, 0x40, 0xdb, 0x00, 0x86, 0x52, 0x09, 0x66, 0x26, 0xb9, 0x1a, 0x19, 0x0b, 0x86, 0x8e, 0xee,
	0xec, 0xdf, 0xe9, 0x77, 0x76, 0xbd, 0xfd, 0x17, 0xc3, 0xb6, 0xb3, 0x1a, 0xd6, 0x8e, 0x1b, 0x7f,
	0xf2, 0x43, 0xa7, 0x19, 0x5f, 0x16, 0xf4, 0xe6, 0x81, 0x0f, 0xfd, 0x69, 0x81, 0x7b, 0x9e, 0x54,
	0x3c, 0x9f, 0x97, 0x84, 0xfd, 0x03, 0xb3, 0x07, 0xd0, 0x9d, 0x89, 0xf8, 0x3b, 0xe5, 0xc6, 0xea,
	0xf3, 0x76, 0xab, 0x27, 0x8a, 0x83, 0x0d, 0x77, 0x30, 0x81, 0xae, 0x46, 0xd0, 0x2b, 0xf0, 0x63,
	0xc1, 0x44, 0x4a, 0x78, 0x72, 0x7b, 0xdf, 0xc5, 0x93, 0x06, 0xd7, 0x4e, 0x76, 0xc0, 0x13, 0x45,
	0x41, 0xcb, 0x68, 0x96, 0x8b, 0xec, 0xc6, 0x58, 0x01, 0x05, 0x9d, 0x48, 0x64, 0xf0, 0x67, 0x0d,
	0xba, 0xa1, 0xca, 0x18, 0x3a, 0x04, 0x27, 0x95, 0x31, 0x0a, 0x2c, 0xe5, 0x6a, 0xa7, 0xdd, 0xd5,
	0x2a, 0x69, 0x58, 0xb3, 0xd1, 0x1b, 0x70, 0xe6, 0x32, 0x46, 0x6a, 0xb8, 0xb7, 0xff, 0xac, 0x5d,
	0xa6, 0x92, 0x86, 0x35, 0x13, 0xbd, 0x85, 0x5e, 0xac, 0xa3, 0x15, 0x74, 0x94, 0x68, 0xbb, 0x5d,
	0x64, 0xf2, 0x87, 0x6b, 0xb6, 0x14, 0x56, 0x3a, 0x33, 0x81, 0xfd, 0x98, 0xd0, 0x04, 0x0b, 0xd7,
	0x6c, 0x29, 0x14, 0xfa, 0x8e, 0x03, 0xe7, 0x31, 0xa1, 0x09, 0x02, 0xae, 0xd9, 0xe8, 0x03, 0xb8,
	0x8b, 0xfa, 0xea, 0x83, 0x9e, 0x92, 0x3e, 0x70, 0x30, 0xab, 0x84, 0xe0, 0x46, 0x21, 0xc3, 0xc2,
	0x13, 0x46, 0x2b, 0x4e, 0x58, 0x11, 0xb1, 0x2a, 0xe8, 0xf6, 0xad, 0xdd, 0x0e, 0xf6, 0x56, 0x58,
	0x58, 0x0d, 0x7e, 0x5b, 0xb0, 0xa1, 0x6f, 0xe0, 0x13, 0x61, 0x49, 0xba, 0x6c, 0xfd, 0x83, 0x11,
	0xd8, 0x0b, 0x9a, 0x16, 0xe6, 0x07, 0x56, 0x35, 0x3a, 0x00, 0x5b, 0x7a, 0x54, 0x47, 0xf8, 0xff,
	0x7e, 0xbf, 0xdd, 0x95, 0x9e, 0x3c, 0x59, 0x16, 0x14, 0x2b, 0xb6, 0x0c, 0x9f, 0x7e, 0x53, 0x02,
	0xfb, 0xb1, 0xf0, 0x69, 0x1d, 0x36, 0xdc, 0xd7, 0x21, 0x40, 0x33, 0x09, 0x79, 0xd0, 0x3b, 0xbd,
	0x9c, 0x5e, 0x4c, 0xce, 0xb0, 0xff, 0x1f, 0x72, 0xc1, 0x19, 0x1d, 0x4f, 0x47, 0x67, 0xbe, 0x25,
	0xf1, 0xf1, 0x34, 0x0c, 0x8f, 0xf1, 0xb5, 0xbf, 0x26, 0x17, 0xd3, 0x8b, 0xc9, 0xf5, 0xd5, 0xd9,
	0x47, 0xbf, 0x83, 0x36, 0xc1, 0x3d, 0xff, 0x3c, 0x9e, 0x5c, 0x8e, 0xf0, 0x71, 0xe8, 0xdb, 0x27,
	0x18, 0x5a, 0x5f, 0xb2, 0x6f, 0x47, 0xf3, 0x84, 0x2f, 0xc4, 0x6c, 0x18, 0xe7, 0x6c, 0xaf, 0xe9,
	0xee, 0xe9, 0x6e, 0xc4, 0xf2, 0x1b, 0x9a, 0xee, 0xcd, 0xf3, 0x77, 0x49, 0x1e, 0x35, 0xdd, 0x48,
	0x77, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0x45, 0x21, 0x7f, 0x64, 0x2b, 0x05, 0x00, 0x00,
}
