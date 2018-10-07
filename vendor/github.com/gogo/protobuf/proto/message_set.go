






























package proto



import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync"
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
	return ms.find(pb) != nil
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
	return marshalMessageSet(exts, false)
}


func marshalMessageSet(exts interface{}, deterministic bool) ([]byte, error) {
	switch exts := exts.(type) {
	case *XXX_InternalExtensions:
		var u marshalInfo
		siz := u.sizeMessageSet(exts)
		b := make([]byte, 0, siz)
		return u.appendMessageSet(b, exts, deterministic)

	case map[int32]Extension:
		
		
		ie := XXX_InternalExtensions{
			p: &struct {
				mu           sync.Mutex
				extensionMap map[int32]Extension
			}{
				extensionMap: exts,
			},
		}

		var u marshalInfo
		siz := u.sizeMessageSet(&ie)
		b := make([]byte, 0, siz)
		return u.appendMessageSet(b, &ie, deterministic)

	default:
		return nil, errors.New("proto: not an extension map")
	}
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
		var mu sync.Locker
		m, mu = exts.extensionsRead()
		if m != nil {
			
			
			
			mu.Lock()
			defer mu.Unlock()
		}
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
		msd, ok := messageSetMap[id]
		if !ok {
			
			continue
		}

		if i > 0 && b.Len() > 1 {
			b.WriteByte(',')
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
