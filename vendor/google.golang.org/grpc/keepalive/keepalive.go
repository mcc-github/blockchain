



package keepalive

import (
	"time"
)







type ClientParameters struct {
	
	
	
	Time time.Duration 
	
	
	
	Timeout time.Duration 
	
	
	
	PermitWithoutStream bool 
}



type ServerParameters struct {
	
	
	
	
	MaxConnectionIdle time.Duration 
	
	
	
	
	MaxConnectionAge time.Duration 
	
	
	MaxConnectionAgeGrace time.Duration 
	
	
	
	Time time.Duration 
	
	
	
	Timeout time.Duration 
}




type EnforcementPolicy struct {
	
	
	MinTime time.Duration 
	
	
	
	PermitWithoutStream bool 
}
