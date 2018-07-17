





package utf8internal


const (
	LoCB = 0x80 
	HiCB = 0xBF 
)


const (
	
	ASCII = as

	
	
	FirstInvalid = xx

	
	SizeMask = 7

	
	
	AcceptShift = 4

	
	
	
	
	xx = 0xF1 
	as = 0xF0 
	s1 = 0x02 
	s2 = 0x13 
	s3 = 0x03 
	s4 = 0x23 
	s5 = 0x34 
	s6 = 0x04 
	s7 = 0x44 
)


var First = [256]uint8{
	
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, 
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, 
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, 
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, 
	xx, xx, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, 
	s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, 
	s2, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s4, s3, s3, 
	s5, s6, s6, s6, s7, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, 
}



type AcceptRange struct {
	Lo uint8 
	Hi uint8 
}







var AcceptRanges = [...]AcceptRange{
	0: {LoCB, HiCB},
	1: {0xA0, HiCB},
	2: {LoCB, 0x9F},
	3: {0x90, HiCB},
	4: {LoCB, 0x8F},
}
