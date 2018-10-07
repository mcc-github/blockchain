package objx



func (m Map) Exclude(exclude []string) Map {
	excluded := make(Map)
	for k, v := range m {
		if !contains(exclude, k) {
			excluded[k] = v
		}
	}
	return excluded
}


func (m Map) Copy() Map {
	copied := Map{}
	for k, v := range m {
		copied[k] = v
	}
	return copied
}





func (m Map) Merge(merge Map) Map {
	return m.Copy().MergeHere(merge)
}






func (m Map) MergeHere(merge Map) Map {
	for k, v := range merge {
		m[k] = v
	}
	return m
}




func (m Map) Transform(transformer func(key string, value interface{}) (string, interface{})) Map {
	newMap := Map{}
	for k, v := range m {
		modifiedKey, modifiedVal := transformer(k, v)
		newMap[modifiedKey] = modifiedVal
	}
	return newMap
}





func (m Map) TransformKeys(mapping map[string]string) Map {
	return m.Transform(func(key string, value interface{}) (string, interface{}) {
		if newKey, ok := mapping[key]; ok {
			return newKey, value
		}
		return key, value
	})
}


func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
