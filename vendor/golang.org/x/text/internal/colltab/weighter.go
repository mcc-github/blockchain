



package colltab 


type Weighter interface {
	
	Start(p int, b []byte) int

	
	StartString(p int, s string) int

	
	
	
	AppendNext(buf []Elem, s []byte) (ce []Elem, n int)

	
	
	
	AppendNextString(buf []Elem, s string) (ce []Elem, n int)

	
	
	Domain() []string

	
	Top() uint32
}
