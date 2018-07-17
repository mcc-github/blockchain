package toml

import "strings"




type MetaData struct {
	mapping map[string]interface{}
	types   map[string]tomlType
	keys    []Key
	decoded map[string]bool
	context Key 
}








func (md *MetaData) IsDefined(key ...string) bool {
	if len(key) == 0 {
		return false
	}

	var hash map[string]interface{}
	var ok bool
	var hashOrVal interface{} = md.mapping
	for _, k := range key {
		if hash, ok = hashOrVal.(map[string]interface{}); !ok {
			return false
		}
		if hashOrVal, ok = hash[k]; !ok {
			return false
		}
	}
	return true
}





func (md *MetaData) Type(key ...string) string {
	fullkey := strings.Join(key, ".")
	if typ, ok := md.types[fullkey]; ok {
		return typ.typeString()
	}
	return ""
}



type Key []string

func (k Key) String() string {
	return strings.Join(k, ".")
}

func (k Key) maybeQuotedAll() string {
	var ss []string
	for i := range k {
		ss = append(ss, k.maybeQuoted(i))
	}
	return strings.Join(ss, ".")
}

func (k Key) maybeQuoted(i int) string {
	quote := false
	for _, c := range k[i] {
		if !isBareKeyChar(c) {
			quote = true
			break
		}
	}
	if quote {
		return "\"" + strings.Replace(k[i], "\"", "\\\"", -1) + "\""
	}
	return k[i]
}

func (k Key) add(piece string) Key {
	newKey := make(Key, len(k)+1)
	copy(newKey, k)
	newKey[len(k)] = piece
	return newKey
}








func (md *MetaData) Keys() []Key {
	return md.keys
}












func (md *MetaData) Undecoded() []Key {
	undecoded := make([]Key, 0, len(md.keys))
	for _, key := range md.keys {
		if !md.decoded[key.String()] {
			undecoded = append(undecoded, key)
		}
	}
	return undecoded
}
