



























package proto

import (
	"reflect"
	"time"
)

var durationType = reflect.TypeOf((*time.Duration)(nil)).Elem()

type duration struct {
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	Nanos   int32 `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
}

func (m *duration) Reset()       { *m = duration{} }
func (*duration) ProtoMessage()  {}
func (*duration) String() string { return "duration<string>" }

func init() {
	RegisterType((*duration)(nil), "gogo.protobuf.proto.duration")
}
