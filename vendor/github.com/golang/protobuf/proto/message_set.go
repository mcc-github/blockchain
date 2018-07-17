






























package proto



import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
)



var errNoMessageTypeID = errors.New("proto does not have a message type ID")












type _MessageSet_Item struct {
	TypeId  *int32 `protobuf:"varint,2,req,name=type_id"`
	Message []byte `protobuf:"bytes,3,req,name=message"`
}

type messageSet struct {
	Item             []*_MessageSet_Item `protobuf:"group,1,rep"`
	XXX_unrecognized []byte
	
}


var _ Message = (*messageSet)(nil)



type messageTypeIder interface {
	MessageTypeId() int32
}

func (ms *messageSet) find(pb Message) *_MessageSet_Item {
	mti, ok := pb.(messageTypeIder)
	if !ok {
		return nil
	}
	id := mti.MessageTypeId()
	for _, item := range ms.Item {
		if *item.TypeId == id {
			return item
		}
	}
	return nil
}

func (ms *messageSet) Has(pb Message) bool {
	if ms.find(pb) != nil {
		return true
	}
	return false
}

func (ms *messageSet) Unmarshal(pb Message) error {
	if item := ms.find(pb); item != nil {
		return Unmarshal(item.Message, pb)
	}
	if _, ok := pb.(messageTypeIder); !ok {
		return errNoMessageTypeID
	}
	return nil 
}

func (ms *messageSet) Marshal(pb Message) error {
	msg, err := Marshal(pb)
	if err != nil {
		return err
	}
	if item := ms.find(pb); item != nil {
		
		item.Message = msg
		return nil
	}

	mti, ok := pb.(messageTypeIder)
	if !ok {
		return errNoMessageTypeID
	}

	mtid := mti.MessageTypeId()
	ms.Item = append(ms.Item, &_MessageSet_Item{
		TypeId:  &mtid,
		Message: msg,
	})
	return nil
}

func (ms *messageSet) Reset()         { *ms = messageSet{} }
func (ms *messageSet) String() string { return CompactTextString(ms) }
func (*messageSet) ProtoMessage()     {}



func skipVarint(buf []byte) []byte {
	i := 0
	for ; buf[i]&0x80 != 0; i++ {
	}
	return buf[i+1:]
}



func MarshalMessageSet(exts interface{}) ([]byte, error) {
	var m map[int32]Extension
	switch exts := exts.(type) {
	case *XXX_InternalExtensions:
		if err := encodeExtensions(exts); err != nil {
			return nil, err
		}
		m, _ = exts.extensionsRead()
	case map[int32]Extension:
		if err := encodeExtensionsMap(exts); err != nil {
			return nil, err
		}
		m = exts
	default:
		return nil, errors.New("proto: not an extension map")
	}

	
	
	ids := make([]int, 0, len(m))
	for id := range m {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)

	ms := &messageSet{Item: make([]*_MessageSet_Item, 0, len(m))}
	for _, id := range ids {
		e := m[int32(id)]
		
		msg := skipVarint(skipVarint(e.enc))

		ms.Item = append(ms.Item, &_MessageSet_Item{
			TypeId:  Int32(int32(id)),
			Message: msg,
		})
	}
	return Marshal(ms)
}



func UnmarshalMessageSet(buf []byte, exts interface{}) error {
	var m map[int32]Extension
	switch exts := exts.(type) {
	case *XXX_InternalExtensions:
		m = exts.extensionsWrite()
	case map[int32]Extension:
		m = exts
	default:
		return errors.New("proto: not an extension map")
	}

	ms := new(messageSet)
	if err := Unmarshal(buf, ms); err != nil {
		return err
	}
	for _, item := range ms.Item {
		id := *item.TypeId
		msg := item.Message

		
		
		b := EncodeVarint(uint64(id)<<3 | WireBytes)
		if ext, ok := m[id]; ok {
			
			
			
			o := ext.enc[len(b):]   
			_, n := DecodeVarint(o) 
			o = o[n:]               
			msg = append(o, msg...) 
		}
		b = append(b, EncodeVarint(uint64(len(msg)))...)
		b = append(b, msg...)

		m[id] = Extension{enc: b}
	}
	return nil
}



func MarshalMessageSetJSON(exts interface{}) ([]byte, error) {
	var m map[int32]Extension
	switch exts := exts.(type) {
	case *XXX_InternalExtensions:
		m, _ = exts.extensionsRead()
	case map[int32]Extension:
		m = exts
	default:
		return nil, errors.New("proto: not an extension map")
	}
	var b bytes.Buffer
	b.WriteByte('{')

	
	ids := make([]int32, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	sort.Sort(int32Slice(ids)) 

	for i, id := range ids {
		ext := m[id]
		if i > 0 {
			b.WriteByte(',')
		}

		msd, ok := messageSetMap[id]
		if !ok {
			
			continue
		}
		fmt.Fprintf(&b, `"[%s]":`, msd.name)

		x := ext.value
		if x == nil {
			x = reflect.New(msd.t.Elem()).Interface()
			if err := Unmarshal(ext.enc, x.(Message)); err != nil {
				return nil, err
			}
		}
		d, err := json.Marshal(x)
		if err != nil {
			return nil, err
		}
		b.Write(d)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}



func UnmarshalMessageSetJSON(buf []byte, exts interface{}) error {
	
	if len(buf) == 0 || bytes.Equal(buf, []byte("{}")) {
		return nil
	}

	
	return errors.New("TODO: UnmarshalMessageSetJSON not yet implemented")
}



var messageSetMap = make(map[int32]messageSetDesc)

type messageSetDesc struct {
	t    reflect.Type 
	name string
}


func RegisterMessageSetType(m Message, fieldNum int32, name string) {
	messageSetMap[fieldNum] = messageSetDesc{
		t:    reflect.TypeOf(m),
		name: name,
	}
}
