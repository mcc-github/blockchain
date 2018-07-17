






























package ptypes




import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

const googleApis = "type.googleapis.com/"






func AnyMessageName(any *any.Any) (string, error) {
	slash := strings.LastIndex(any.TypeUrl, "/")
	if slash < 0 {
		return "", fmt.Errorf("message type url %q is invalid", any.TypeUrl)
	}
	return any.TypeUrl[slash+1:], nil
}


func MarshalAny(pb proto.Message) (*any.Any, error) {
	value, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}
	return &any.Any{TypeUrl: googleApis + proto.MessageName(pb), Value: value}, nil
}










type DynamicAny struct {
	proto.Message
}




func Empty(any *any.Any) (proto.Message, error) {
	aname, err := AnyMessageName(any)
	if err != nil {
		return nil, err
	}

	t := proto.MessageType(aname)
	if t == nil {
		return nil, fmt.Errorf("any: message type %q isn't linked in", aname)
	}
	return reflect.New(t.Elem()).Interface().(proto.Message), nil
}






func UnmarshalAny(any *any.Any, pb proto.Message) error {
	if d, ok := pb.(*DynamicAny); ok {
		if d.Message == nil {
			var err error
			d.Message, err = Empty(any)
			if err != nil {
				return err
			}
		}
		return UnmarshalAny(any, d.Message)
	}

	aname, err := AnyMessageName(any)
	if err != nil {
		return err
	}

	mname := proto.MessageName(pb)
	if aname != mname {
		return fmt.Errorf("mismatched message type: got %q want %q", aname, mname)
	}
	return proto.Unmarshal(any.Value, pb)
}


func Is(any *any.Any, pb proto.Message) bool {
	aname, err := AnyMessageName(any)
	if err != nil {
		return false
	}

	return aname == proto.MessageName(pb)
}
