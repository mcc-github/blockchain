package system 

import (
	"syscall"
	"time"
	"unsafe"
)


var maxTime time.Time

func init() {
	
	if unsafe.Sizeof(syscall.Timespec{}.Nsec) == 8 {
		
		
		maxTime = time.Unix(0, 1<<63-1)
	} else {
		
		maxTime = time.Unix(1<<31-1, 0)
	}
}
