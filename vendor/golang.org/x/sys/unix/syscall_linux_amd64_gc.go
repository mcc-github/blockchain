






package unix

import "syscall"


func gettimeofday(tv *Timeval) (err syscall.Errno)
