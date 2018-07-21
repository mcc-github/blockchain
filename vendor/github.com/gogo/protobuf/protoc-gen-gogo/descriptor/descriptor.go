



































package descriptor

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	"github.com/gogo/protobuf/proto"
)


func extractFile(gz []byte) (*FileDescriptorProto, error) {
	r, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		return nil, fmt.Errorf("failed to open gzip reader: %v", err)
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to uncompress descriptor: %v", err)
	}

	fd := new(FileDescriptorProto)
	if err := proto.Unmarshal(b, fd); err != nil {
		return nil, fmt.Errorf("malformed FileDescriptorProto: %v", err)
	}

	return fd, nil
}





type Message interface {
	proto.Message
	Descriptor() ([]byte, []int)
}



func ForMessage(msg Message) (fd *FileDescriptorProto, md *DescriptorProto) {
	gz, path := msg.Descriptor()
	fd, err := extractFile(gz)
	if err != nil {
		panic(fmt.Sprintf("invalid FileDescriptorProto for %T: %v", msg, err))
	}

	md = fd.MessageType[path[0]]
	for _, i := range path[1:] {
		md = md.NestedType[i]
	}
	return fd, md
}


func (field *FieldDescriptorProto) IsScalar() bool {
	if field.Type == nil {
		return false
	}
	switch *field.Type {
	case FieldDescriptorProto_TYPE_DOUBLE,
		FieldDescriptorProto_TYPE_FLOAT,
		FieldDescriptorProto_TYPE_INT64,
		FieldDescriptorProto_TYPE_UINT64,
		FieldDescriptorProto_TYPE_INT32,
		FieldDescriptorProto_TYPE_FIXED64,
		FieldDescriptorProto_TYPE_FIXED32,
		FieldDescriptorProto_TYPE_BOOL,
		FieldDescriptorProto_TYPE_UINT32,
		FieldDescriptorProto_TYPE_ENUM,
		FieldDescriptorProto_TYPE_SFIXED32,
		FieldDescriptorProto_TYPE_SFIXED64,
		FieldDescriptorProto_TYPE_SINT32,
		FieldDescriptorProto_TYPE_SINT64:
		return true
	default:
		return false
	}
}
