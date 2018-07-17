





package main


type Class uint

const (
	L       Class = iota 
	R                    
	EN                   
	ES                   
	ET                   
	AN                   
	CS                   
	B                    
	S                    
	WS                   
	ON                   
	BN                   
	NSM                  
	AL                   
	Control              

	numClass

	LRO 
	RLO 
	LRE 
	RLE 
	PDF 
	LRI 
	RLI 
	FSI 
	PDI 

	unknownClass = ^Class(0)
)

var controlToClass = map[rune]Class{
	0x202D: LRO, 
	0x202E: RLO, 
	0x202A: LRE, 
	0x202B: RLE, 
	0x202C: PDF, 
	0x2066: LRI, 
	0x2067: RLI, 
	0x2068: FSI, 
	0x2069: PDI, 
}






const (
	openMask     = 0x10
	xorMaskShift = 5
)
