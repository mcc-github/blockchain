/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protolator

import (
	"github.com/golang/protobuf/proto"
)
































type StaticallyOpaqueFieldProto interface {
	
	StaticallyOpaqueFields() []string

	
	
	StaticallyOpaqueFieldProto(name string) (proto.Message, error)
}



type StaticallyOpaqueMapFieldProto interface {
	
	StaticallyOpaqueMapFields() []string

	
	
	StaticallyOpaqueMapFieldProto(name string, key string) (proto.Message, error)
}



type StaticallyOpaqueSliceFieldProto interface {
	
	StaticallyOpaqueSliceFields() []string

	
	
	StaticallyOpaqueSliceFieldProto(name string, index int) (proto.Message, error)
}



type VariablyOpaqueFieldProto interface {
	
	VariablyOpaqueFields() []string

	
	
	VariablyOpaqueFieldProto(name string) (proto.Message, error)
}



type VariablyOpaqueMapFieldProto interface {
	
	VariablyOpaqueMapFields() []string

	
	
	VariablyOpaqueMapFieldProto(name string, key string) (proto.Message, error)
}



type VariablyOpaqueSliceFieldProto interface {
	
	VariablyOpaqueSliceFields() []string

	
	
	VariablyOpaqueSliceFieldProto(name string, index int) (proto.Message, error)
}



type DynamicFieldProto interface {
	
	DynamicFields() []string

	
	
	DynamicFieldProto(name string, underlying proto.Message) (proto.Message, error)
}



type DynamicMapFieldProto interface {
	
	DynamicMapFields() []string

	
	
	DynamicMapFieldProto(name string, key string, underlying proto.Message) (proto.Message, error)
}



type DynamicSliceFieldProto interface {
	
	DynamicSliceFields() []string

	
	
	DynamicSliceFieldProto(name string, index int, underlying proto.Message) (proto.Message, error)
}




type DecoratedProto interface {
	
	Underlying() proto.Message
}
