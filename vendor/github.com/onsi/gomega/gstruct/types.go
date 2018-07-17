package gstruct


type Options int

const (
	
	IgnoreExtras Options = 1 << iota
	
	IgnoreMissing
	
	
	
	AllowDuplicates
)
