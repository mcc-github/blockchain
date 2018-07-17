package objx



func (m Map) Exclude(exclude []string) Map {
	excluded := make(Map)
	for k, v := range m {
		var shouldInclude = true
		for _, toExclude := range exclude {
			if k == toExclude {
				shouldInclude = false
				break
			}
		}
		if shouldInclude {
			excluded[k] = v
		}
	}
	return excluded
}


func (m Map) Copy() Map {
	copied := make(map[string]interface{})
	for k, v := range m {
		copied[k] = v
	}
	return New(copied)
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
	newMap := make(map[string]interface{})
	for k, v := range m {
		modifiedKey, modifiedVal := transformer(k, v)
		newMap[modifiedKey] = modifiedVal
	}
	return New(newMap)
}





func (m Map) TransformKeys(mapping map[string]string) Map {
	return m.Transform(func(key string, value interface{}) (string, interface{}) {
		if newKey, ok := mapping[key]; ok {
			return newKey, value
		}
		return key, value
	})
}
