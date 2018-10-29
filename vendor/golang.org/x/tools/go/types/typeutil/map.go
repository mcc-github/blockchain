





package typeutil 

import (
	"bytes"
	"fmt"
	"go/types"
	"reflect"
)











type Map struct {
	hasher Hasher             
	table  map[uint32][]entry 
	length int                
}


type entry struct {
	key   types.Type
	value interface{}
}






















func (m *Map) SetHasher(hasher Hasher) {
	m.hasher = hasher
}




func (m *Map) Delete(key types.Type) bool {
	if m != nil && m.table != nil {
		hash := m.hasher.Hash(key)
		bucket := m.table[hash]
		for i, e := range bucket {
			if e.key != nil && types.Identical(key, e.key) {
				
				
				bucket[i] = entry{}
				m.length--
				return true
			}
		}
	}
	return false
}




func (m *Map) At(key types.Type) interface{} {
	if m != nil && m.table != nil {
		for _, e := range m.table[m.hasher.Hash(key)] {
			if e.key != nil && types.Identical(key, e.key) {
				return e.value
			}
		}
	}
	return nil
}



func (m *Map) Set(key types.Type, value interface{}) (prev interface{}) {
	if m.table != nil {
		hash := m.hasher.Hash(key)
		bucket := m.table[hash]
		var hole *entry
		for i, e := range bucket {
			if e.key == nil {
				hole = &bucket[i]
			} else if types.Identical(key, e.key) {
				prev = e.value
				bucket[i].value = value
				return
			}
		}

		if hole != nil {
			*hole = entry{key, value} 
		} else {
			m.table[hash] = append(bucket, entry{key, value})
		}
	} else {
		if m.hasher.memo == nil {
			m.hasher = MakeHasher()
		}
		hash := m.hasher.Hash(key)
		m.table = map[uint32][]entry{hash: {entry{key, value}}}
	}

	m.length++
	return
}


func (m *Map) Len() int {
	if m != nil {
		return m.length
	}
	return 0
}









func (m *Map) Iterate(f func(key types.Type, value interface{})) {
	if m != nil {
		for _, bucket := range m.table {
			for _, e := range bucket {
				if e.key != nil {
					f(e.key, e.value)
				}
			}
		}
	}
}



func (m *Map) Keys() []types.Type {
	keys := make([]types.Type, 0, m.Len())
	m.Iterate(func(key types.Type, _ interface{}) {
		keys = append(keys, key)
	})
	return keys
}

func (m *Map) toString(values bool) string {
	if m == nil {
		return "{}"
	}
	var buf bytes.Buffer
	fmt.Fprint(&buf, "{")
	sep := ""
	m.Iterate(func(key types.Type, value interface{}) {
		fmt.Fprint(&buf, sep)
		sep = ", "
		fmt.Fprint(&buf, key)
		if values {
			fmt.Fprintf(&buf, ": %q", value)
		}
	})
	fmt.Fprint(&buf, "}")
	return buf.String()
}





func (m *Map) String() string {
	return m.toString(true)
}




func (m *Map) KeysString() string {
	return m.toString(false)
}










type Hasher struct {
	memo map[types.Type]uint32
}


func MakeHasher() Hasher {
	return Hasher{make(map[types.Type]uint32)}
}



func (h Hasher) Hash(t types.Type) uint32 {
	hash, ok := h.memo[t]
	if !ok {
		hash = h.hashFor(t)
		h.memo[t] = hash
	}
	return hash
}


func hashString(s string) uint32 {
	var h uint32
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}


func (h Hasher) hashFor(t types.Type) uint32 {
	
	switch t := t.(type) {
	case *types.Basic:
		return uint32(t.Kind())

	case *types.Array:
		return 9043 + 2*uint32(t.Len()) + 3*h.Hash(t.Elem())

	case *types.Slice:
		return 9049 + 2*h.Hash(t.Elem())

	case *types.Struct:
		var hash uint32 = 9059
		for i, n := 0, t.NumFields(); i < n; i++ {
			f := t.Field(i)
			if f.Anonymous() {
				hash += 8861
			}
			hash += hashString(t.Tag(i))
			hash += hashString(f.Name()) 
			hash += h.Hash(f.Type())
		}
		return hash

	case *types.Pointer:
		return 9067 + 2*h.Hash(t.Elem())

	case *types.Signature:
		var hash uint32 = 9091
		if t.Variadic() {
			hash *= 8863
		}
		return hash + 3*h.hashTuple(t.Params()) + 5*h.hashTuple(t.Results())

	case *types.Interface:
		var hash uint32 = 9103
		for i, n := 0, t.NumMethods(); i < n; i++ {
			
			
			
			m := t.Method(i)
			hash += 3*hashString(m.Name()) + 5*h.Hash(m.Type())
		}
		return hash

	case *types.Map:
		return 9109 + 2*h.Hash(t.Key()) + 3*h.Hash(t.Elem())

	case *types.Chan:
		return 9127 + 2*uint32(t.Dir()) + 3*h.Hash(t.Elem())

	case *types.Named:
		
		return uint32(reflect.ValueOf(t.Obj()).Pointer())

	case *types.Tuple:
		return h.hashTuple(t)
	}
	panic(t)
}

func (h Hasher) hashTuple(tuple *types.Tuple) uint32 {
	
	n := tuple.Len()
	var hash uint32 = 9137 + 2*uint32(n)
	for i := 0; i < n; i++ {
		hash += 3 * h.Hash(tuple.At(i).Type())
	}
	return hash
}
