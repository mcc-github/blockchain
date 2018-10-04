





package unix

import "unsafe"



var fcntl64Syscall uintptr = SYS_FCNTL


func FcntlInt(fd uintptr, cmd, arg int) (int, error) {
	valptr, _, errno := Syscall(fcntl64Syscall, fd, uintptr(cmd), uintptr(arg))
	var err error
	if errno != 0 {
		err = errno
	}
	return int(valptr), err
}


func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error {
	_, _, errno := Syscall(fcntl64Syscall, fd, uintptr(cmd), uintptr(unsafe.Pointer(lk)))
	if errno == 0 {
		return nil
	}
	return errno
}
