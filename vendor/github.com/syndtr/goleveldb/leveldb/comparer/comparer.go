







package comparer


type BasicComparer interface {
	
	
	
	
	Compare(a, b []byte) int
}



type Comparer interface {
	BasicComparer

	
	
	
	
	
	
	
	
	
	
	
	
	Name() string

	
	

	
	
	
	
	
	
	Separator(dst, a, b []byte) []byte

	
	
	
	
	
	
	Successor(dst, b []byte) []byte
}
