package objx





func (m Map) Has(selector string) bool {
	if m == nil {
		return false
	}
	return !m.Get(selector).IsNil()
}


func (v *Value) IsNil() bool {
	return v == nil || v.data == nil
}
