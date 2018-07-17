





package unix

import (
	"unsafe"
)








































func Time(t *Time_t) (tt Time_t, err error) {
	var tv Timeval
	err = Gettimeofday(&tv)
	if err != nil {
		return 0, err
	}
	if t != nil {
		*t = Time_t(tv.Sec)
	}
	return Time_t(tv.Sec), nil
}



func setTimespec(sec, nsec int64) Timespec {
	return Timespec{Sec: sec, Nsec: nsec}
}

func setTimeval(sec, usec int64) Timeval {
	return Timeval{Sec: sec, Usec: usec}
}



func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	err = pipe2(&pp, 0) 
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
}

func Pipe2(p []int, flags int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	err = pipe2(&pp, flags)
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
}

func Ioperm(from int, num int, on int) (err error) {
	return ENOSYS
}

func Iopl(level int) (err error) {
	return ENOSYS
}

func (r *PtraceRegs) PC() uint64 { return r.Psw.Addr }

func (r *PtraceRegs) SetPC(pc uint64) { r.Psw.Addr = pc }

func (iov *Iovec) SetLen(length int) {
	iov.Len = uint64(length)
}

func (msghdr *Msghdr) SetControllen(length int) {
	msghdr.Controllen = uint64(length)
}

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = uint64(length)
}



func mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error) {
	mmap_args := [6]uintptr{addr, length, uintptr(prot), uintptr(flags), uintptr(fd), uintptr(offset)}
	r0, _, e1 := Syscall(SYS_MMAP, uintptr(unsafe.Pointer(&mmap_args[0])), 0, 0)
	xaddr = uintptr(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}




const (
	
	netSocket      = 1
	netBind        = 2
	netConnect     = 3
	netListen      = 4
	netAccept      = 5
	netGetSockName = 6
	netGetPeerName = 7
	netSocketPair  = 8
	netSend        = 9
	netRecv        = 10
	netSendTo      = 11
	netRecvFrom    = 12
	netShutdown    = 13
	netSetSockOpt  = 14
	netGetSockOpt  = 15
	netSendMsg     = 16
	netRecvMsg     = 17
	netAccept4     = 18
	netRecvMMsg    = 19
	netSendMMsg    = 20
)

func accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (int, error) {
	args := [3]uintptr{uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen))}
	fd, _, err := Syscall(SYS_SOCKETCALL, netAccept, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return 0, err
	}
	return int(fd), nil
}

func accept4(s int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (int, error) {
	args := [4]uintptr{uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)), uintptr(flags)}
	fd, _, err := Syscall(SYS_SOCKETCALL, netAccept4, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return 0, err
	}
	return int(fd), nil
}

func getsockname(s int, rsa *RawSockaddrAny, addrlen *_Socklen) error {
	args := [3]uintptr{uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen))}
	_, _, err := RawSyscall(SYS_SOCKETCALL, netGetSockName, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return err
	}
	return nil
}

func getpeername(s int, rsa *RawSockaddrAny, addrlen *_Socklen) error {
	args := [3]uintptr{uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen))}
	_, _, err := RawSyscall(SYS_SOCKETCALL, netGetPeerName, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return err
	}
	return nil
}

func socketpair(domain int, typ int, flags int, fd *[2]int32) error {
	args := [4]uintptr{uintptr(domain), uintptr(typ), uintptr(flags), uintptr(unsafe.Pointer(fd))}
	_, _, err := RawSyscall(SYS_SOCKETCALL, netSocketPair, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return err
	}
	return nil
}

func bind(s int, addr unsafe.Pointer, addrlen _Socklen) error {
	args := [3]uintptr{uintptr(s), uintptr(addr), uintptr(addrlen)}
	_, _, err := Syscall(SYS_SOCKETCALL, netBind, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return err
	}
	return nil
}

func connect(s int, addr unsafe.Pointer, addrlen _Socklen) error {
	args := [3]uintptr{uintptr(s), uintptr(addr), uintptr(addrlen)}
	_, _, err := Syscall(SYS_SOCKETCALL, netConnect, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return err
	}
	return nil
}

func socket(domain int, typ int, proto int) (int, error) {
	args := [3]uintptr{uintptr(domain), uintptr(typ), uintptr(proto)}
	fd, _, err := RawSyscall(SYS_SOCKETCALL, netSocket, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return 0, err
	}
	return int(fd), nil
}

func getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) error {
	args := [5]uintptr{uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(unsafe.Pointer(vallen))}
	_, _, err := Syscall(SYS_SOCKETCALL, netGetSockOpt, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return err
	}
	return nil
}

func setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) error {
	args := [4]uintptr{uintptr(s), uintptr(level), uintptr(name), uintptr(val)}
	_, _, err := Syscall(SYS_SOCKETCALL, netSetSockOpt, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return err
	}
	return nil
}

func recvfrom(s int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (int, error) {
	var base uintptr
	if len(p) > 0 {
		base = uintptr(unsafe.Pointer(&p[0]))
	}
	args := [6]uintptr{uintptr(s), base, uintptr(len(p)), uintptr(flags), uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(fromlen))}
	n, _, err := Syscall(SYS_SOCKETCALL, netRecvFrom, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return 0, err
	}
	return int(n), nil
}

func sendto(s int, p []byte, flags int, to unsafe.Pointer, addrlen _Socklen) error {
	var base uintptr
	if len(p) > 0 {
		base = uintptr(unsafe.Pointer(&p[0]))
	}
	args := [6]uintptr{uintptr(s), base, uintptr(len(p)), uintptr(flags), uintptr(to), uintptr(addrlen)}
	_, _, err := Syscall(SYS_SOCKETCALL, netSendTo, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return err
	}
	return nil
}

func recvmsg(s int, msg *Msghdr, flags int) (int, error) {
	args := [3]uintptr{uintptr(s), uintptr(unsafe.Pointer(msg)), uintptr(flags)}
	n, _, err := Syscall(SYS_SOCKETCALL, netRecvMsg, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return 0, err
	}
	return int(n), nil
}

func sendmsg(s int, msg *Msghdr, flags int) (int, error) {
	args := [3]uintptr{uintptr(s), uintptr(unsafe.Pointer(msg)), uintptr(flags)}
	n, _, err := Syscall(SYS_SOCKETCALL, netSendMsg, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return 0, err
	}
	return int(n), nil
}

func Listen(s int, n int) error {
	args := [2]uintptr{uintptr(s), uintptr(n)}
	_, _, err := Syscall(SYS_SOCKETCALL, netListen, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return err
	}
	return nil
}

func Shutdown(s, how int) error {
	args := [2]uintptr{uintptr(s), uintptr(how)}
	_, _, err := Syscall(SYS_SOCKETCALL, netShutdown, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 {
		return err
	}
	return nil
}



func Poll(fds []PollFd, timeout int) (n int, err error) {
	if len(fds) == 0 {
		return poll(nil, 0, timeout)
	}
	return poll(&fds[0], len(fds), timeout)
}
