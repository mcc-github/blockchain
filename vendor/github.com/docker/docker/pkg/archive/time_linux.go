package archive 

import (
	"syscall"
	"time"
)

func timeToTimespec(time time.Time) (ts syscall.Timespec) {
	if time.IsZero() {
		
		ts.Sec = 0
		ts.Nsec = ((1 << 30) - 2)
		return
	}
	return syscall.NsecToTimespec(time.UnixNano())
}
