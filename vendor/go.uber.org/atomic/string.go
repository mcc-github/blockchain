



















package atomic


type String struct{ v Value }


func NewString(str string) *String {
	s := &String{}
	if str != "" {
		s.Store(str)
	}
	return s
}


func (s *String) Load() string {
	v := s.v.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}




func (s *String) Store(str string) {
	s.v.Store(str)
}
