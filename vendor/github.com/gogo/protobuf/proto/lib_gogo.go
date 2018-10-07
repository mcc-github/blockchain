



























package proto

import (
	"encoding/json"
	"strconv"
)

type Sizer interface {
	Size() int
}

type ProtoSizer interface {
	ProtoSize() int
}

func MarshalJSONEnum(m map[int32]string, value int32) ([]byte, error) {
	s, ok := m[value]
	if !ok {
		s = strconv.Itoa(int(value))
	}
	return json.Marshal(s)
}
