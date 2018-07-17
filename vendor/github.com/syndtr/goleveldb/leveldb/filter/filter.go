












package filter


type Buffer interface {
	
	
	Alloc(n int) []byte

	
	Write(p []byte) (n int, err error)

	
	WriteByte(c byte) error
}


type Filter interface {
	
	
	
	
	
	Name() string

	
	NewGenerator() FilterGenerator

	
	
	
	Contains(filter, key []byte) bool
}


type FilterGenerator interface {
	
	
	
	
	
	
	Add(key []byte)

	
	
	Generate(b Buffer)
}
