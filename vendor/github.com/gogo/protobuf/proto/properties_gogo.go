



























package proto

import (
	"reflect"
)

var sizerType = reflect.TypeOf((*Sizer)(nil)).Elem()
var protosizerType = reflect.TypeOf((*ProtoSizer)(nil)).Elem()
