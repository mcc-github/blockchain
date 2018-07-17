

package system

import (
	"syscall"
)


func Setuid(uid int) (err error) {
	_, _, e1 := syscall.RawSyscall(syscall.SYS_SETUID32, uintptr(uid), 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}


func Setgid(gid int) (err error) {
	_, _, e1 := syscall.RawSyscall(syscall.SYS_SETGID32, uintptr(gid), 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}
