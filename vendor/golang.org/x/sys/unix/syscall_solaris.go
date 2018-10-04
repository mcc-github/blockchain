











package unix

import (
	"syscall"
	"unsafe"
)


type syscallFunc uintptr

func rawSysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)
func sysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)


type SockaddrDatalink struct {
	Family uint16
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [244]int8
	raw    RawSockaddrDatalink
}



func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	n, err := pipe(&pp)
	if n != 0 {
		return err
	}
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return nil
}

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet4, nil
}

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet6, nil
}

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) {
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = int8(name[i])
	}
	
	sl := _Socklen(2)
	if n > 0 {
		sl += _Socklen(n) + 1
	}
	if sa.raw.Path[0] == '@' {
		sa.raw.Path[0] = 0
		
		sl--
	}

	return unsafe.Pointer(&sa.raw), sl, nil
}



func Getsockname(fd int) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil {
		return
	}
	return anyToSockaddr(fd, &rsa)
}



func GetsockoptString(fd, level, opt int) (string, error) {
	buf := make([]byte, 256)
	vallen := _Socklen(len(buf))
	err := getsockopt(fd, level, opt, unsafe.Pointer(&buf[0]), &vallen)
	if err != nil {
		return "", err
	}
	return string(buf[:vallen-1]), nil
}

const ImplementsGetwd = true



func Getwd() (wd string, err error) {
	var buf [PathMax]byte
	
	_, err = Getcwd(buf[0:])
	if err != nil {
		return "", err
	}
	n := clen(buf[:])
	if n < 1 {
		return "", EINVAL
	}
	return string(buf[:n]), nil
}






func Getgroups() (gids []int, err error) {
	n, err := getgroups(0, nil)
	
	
	if n < 0 || n > 1024 {
		if err != nil {
			return nil, err
		}
		return nil, EINVAL
	} else if n == 0 {
		return nil, nil
	}

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if n == -1 {
		return nil, err
	}
	gids = make([]int, n)
	for i, v := range a[0:n] {
		gids[i] = int(v)
	}
	return
}

func Setgroups(gids []int) (err error) {
	if len(gids) == 0 {
		return setgroups(0, nil)
	}

	a := make([]_Gid_t, len(gids))
	for i, v := range gids {
		a[i] = _Gid_t(v)
	}
	return setgroups(len(a), &a[0])
}

func ReadDirent(fd int, buf []byte) (n int, err error) {
	
	
	return Getdents(fd, buf, new(uintptr))
}







type WaitStatus uint32

const (
	mask  = 0x7F
	core  = 0x80
	shift = 8

	exited  = 0
	stopped = 0x7F
)

func (w WaitStatus) Exited() bool { return w&mask == exited }

func (w WaitStatus) ExitStatus() int {
	if w&mask != exited {
		return -1
	}
	return int(w >> shift)
}

func (w WaitStatus) Signaled() bool { return w&mask != stopped && w&mask != 0 }

func (w WaitStatus) Signal() syscall.Signal {
	sig := syscall.Signal(w & mask)
	if sig == stopped || sig == 0 {
		return -1
	}
	return sig
}

func (w WaitStatus) CoreDump() bool { return w.Signaled() && w&core != 0 }

func (w WaitStatus) Stopped() bool { return w&mask == stopped && syscall.Signal(w>>shift) != SIGSTOP }

func (w WaitStatus) Continued() bool { return w&mask == stopped && syscall.Signal(w>>shift) == SIGSTOP }

func (w WaitStatus) StopSignal() syscall.Signal {
	if !w.Stopped() {
		return -1
	}
	return syscall.Signal(w>>shift) & 0xFF
}

func (w WaitStatus) TrapCause() int { return -1 }



func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (int, error) {
	var status _C_int
	rpid, err := wait4(int32(pid), &status, options, rusage)
	wpid := int(rpid)
	if wpid == -1 {
		return wpid, err
	}
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	return wpid, nil
}



func Gethostname() (name string, err error) {
	var buf [MaxHostNameLen]byte
	n, err := gethostname(buf[:])
	if n != 0 {
		return "", err
	}
	n = clen(buf[:])
	if n < 1 {
		return "", EFAULT
	}
	return string(buf[:n]), nil
}



func Utimes(path string, tv []Timeval) (err error) {
	if tv == nil {
		return utimes(path, nil)
	}
	if len(tv) != 2 {
		return EINVAL
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}



func UtimesNano(path string, ts []Timespec) error {
	if ts == nil {
		return utimensat(AT_FDCWD, path, nil, 0)
	}
	if len(ts) != 2 {
		return EINVAL
	}
	return utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
}

func UtimesNanoAt(dirfd int, path string, ts []Timespec, flags int) error {
	if ts == nil {
		return utimensat(dirfd, path, nil, flags)
	}
	if len(ts) != 2 {
		return EINVAL
	}
	return utimensat(dirfd, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), flags)
}




func FcntlInt(fd uintptr, cmd, arg int) (int, error) {
	valptr, _, errno := sysvicall6(uintptr(unsafe.Pointer(&procfcntl)), 3, uintptr(fd), uintptr(cmd), uintptr(arg), 0, 0, 0)
	var err error
	if errno != 0 {
		err = errno
	}
	return int(valptr), err
}


func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error {
	_, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procfcntl)), 3, uintptr(fd), uintptr(cmd), uintptr(unsafe.Pointer(lk)), 0, 0, 0)
	if e1 != 0 {
		return e1
	}
	return nil
}



func Futimesat(dirfd int, path string, tv []Timeval) error {
	pathp, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	if tv == nil {
		return futimesat(dirfd, pathp, nil)
	}
	if len(tv) != 2 {
		return EINVAL
	}
	return futimesat(dirfd, pathp, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}




func Futimes(fd int, tv []Timeval) error {
	if tv == nil {
		return futimesat(fd, nil, nil)
	}
	if len(tv) != 2 {
		return EINVAL
	}
	return futimesat(fd, nil, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

func anyToSockaddr(fd int, rsa *RawSockaddrAny) (Sockaddr, error) {
	switch rsa.Addr.Family {
	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)
		
		
		
		
		
		n := 0
		for n < len(pp.Path) && pp.Path[n] != 0 {
			n++
		}
		bytes := (*[10000]byte)(unsafe.Pointer(&pp.Path[0]))[0:n]
		sa.Name = string(bytes)
		return sa, nil

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i]
		}
		return sa, nil

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.ZoneId = pp.Scope_id
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i]
		}
		return sa, nil
	}
	return nil, EAFNOSUPPORT
}



func Accept(fd int) (nfd int, sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept(fd, &rsa, &len)
	if nfd == -1 {
		return
	}
	sa, err = anyToSockaddr(fd, &rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
}



func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error) {
	var msg Msghdr
	var rsa RawSockaddrAny
	msg.Name = (*byte)(unsafe.Pointer(&rsa))
	msg.Namelen = uint32(SizeofSockaddrAny)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = (*int8)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	}
	var dummy int8
	if len(oob) > 0 {
		
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		msg.Accrightslen = int32(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = recvmsg(fd, &msg, flags); n == -1 {
		return
	}
	oobn = int(msg.Accrightslen)
	
	if rsa.Addr.Family != AF_UNSPEC {
		from, err = anyToSockaddr(fd, &rsa)
	}
	return
}

func Sendmsg(fd int, p, oob []byte, to Sockaddr, flags int) (err error) {
	_, err = SendmsgN(fd, p, oob, to, flags)
	return
}



func SendmsgN(fd int, p, oob []byte, to Sockaddr, flags int) (n int, err error) {
	var ptr unsafe.Pointer
	var salen _Socklen
	if to != nil {
		ptr, salen, err = to.sockaddr()
		if err != nil {
			return 0, err
		}
	}
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(ptr))
	msg.Namelen = uint32(salen)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = (*int8)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	}
	var dummy int8
	if len(oob) > 0 {
		
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		msg.Accrightslen = int32(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = sendmsg(fd, &msg, flags); err != nil {
		return 0, err
	}
	if len(oob) > 0 && len(p) == 0 {
		n = 0
	}
	return n, nil
}



func Acct(path string) (err error) {
	if len(path) == 0 {
		
		return acct(nil)
	}

	pathp, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	return acct(pathp)
}



func Mkdev(major, minor uint32) uint64 {
	return __makedev(NEWDEV, uint(major), uint(minor))
}



func Major(dev uint64) uint32 {
	return uint32(__major(NEWDEV, dev))
}



func Minor(dev uint64) uint32 {
	return uint32(__minor(NEWDEV, dev))
}





func IoctlSetInt(fd int, req uint, value int) (err error) {
	return ioctl(fd, req, uintptr(value))
}

func ioctlSetWinsize(fd int, req uint, value *Winsize) (err error) {
	return ioctl(fd, req, uintptr(unsafe.Pointer(value)))
}

func ioctlSetTermios(fd int, req uint, value *Termios) (err error) {
	return ioctl(fd, req, uintptr(unsafe.Pointer(value)))
}

func IoctlSetTermio(fd int, req uint, value *Termio) (err error) {
	return ioctl(fd, req, uintptr(unsafe.Pointer(value)))
}

func IoctlGetInt(fd int, req uint) (int, error) {
	var value int
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return value, err
}

func IoctlGetWinsize(fd int, req uint) (*Winsize, error) {
	var value Winsize
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return &value, err
}

func IoctlGetTermios(fd int, req uint) (*Termios, error) {
	var value Termios
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return &value, err
}

func IoctlGetTermio(fd int, req uint) (*Termio, error) {
	var value Termio
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return &value, err
}



func Poll(fds []PollFd, timeout int) (n int, err error) {
	if len(fds) == 0 {
		return poll(nil, 0, timeout)
	}
	return poll(&fds[0], len(fds), timeout)
}















































































































func readlen(fd int, buf *byte, nbuf int) (n int, err error) {
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procread)), 3, uintptr(fd), uintptr(unsafe.Pointer(buf)), uintptr(nbuf), 0, 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

func writelen(fd int, buf *byte, nbuf int) (n int, err error) {
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procwrite)), 3, uintptr(fd), uintptr(unsafe.Pointer(buf)), uintptr(nbuf), 0, 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

var mapper = &mmapper{
	active: make(map[*byte][]byte),
	mmap:   mmap,
	munmap: munmap,
}

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return mapper.Mmap(fd, offset, length, prot, flags)
}

func Munmap(b []byte) (err error) {
	return mapper.Munmap(b)
}
