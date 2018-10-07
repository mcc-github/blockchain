



























package proto

import "reflect"

type custom interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
	Size() int
}

var customType = reflect.TypeOf((*custom)(nil)).Elem()
