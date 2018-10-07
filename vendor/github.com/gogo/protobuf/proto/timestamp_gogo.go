



























package proto

import (
	"reflect"
	"time"
)

var timeType = reflect.TypeOf((*time.Time)(nil)).Elem()

type timestamp struct {
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	Nanos   int32 `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
}

func (m *timestamp) Reset()       { *m = timestamp{} }
func (*timestamp) ProtoMessage()  {}
func (*timestamp) String() string { return "timestamp<string>" }

func init() {
	RegisterType((*timestamp)(nil), "gogo.protobuf.proto.timestamp")
}
