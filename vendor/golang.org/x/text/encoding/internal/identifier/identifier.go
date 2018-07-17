

























package identifier

















type Interface interface {
	
	
	
	
	
	
	
	
	ID() (mib MIB, other string)

	
	
}






type MIB uint16



const (
	
	Unofficial MIB = 10000 + iota

	
	Replacement

	
	XUserDefined

	
	MacintoshCyrillic
)
