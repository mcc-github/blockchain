











package unix

import (
	"runtime"
	"syscall"
	"unsafe"
)






func Getgroups() (gids []int, err error) {
	n, err := getgroups(0, nil)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	
	if n < 0 || n > 1000 {
		return nil, EINVAL
	}

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if err != nil {
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
	
	
	
	
	var base = (*uintptr)(unsafe.Pointer(new(uint64)))
	return Getdirentries(fd, buf, base)
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



func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) {
	var status _C_int
	wpid, err = wait4(pid, &status, options, rusage)
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	return
}











func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Len = SizeofSockaddrInet4
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
}

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Len = SizeofSockaddrInet6
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
}

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) {
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) || n == 0 {
		return nil, 0, EINVAL
	}
	sa.raw.Len = byte(3 + n) 
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = int8(name[i])
	}
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
}

func (sa *SockaddrDatalink) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Index == 0 {
		return nil, 0, EINVAL
	}
	sa.raw.Len = sa.Len
	sa.raw.Family = AF_LINK
	sa.raw.Index = sa.Index
	sa.raw.Type = sa.Type
	sa.raw.Nlen = sa.Nlen
	sa.raw.Alen = sa.Alen
	sa.raw.Slen = sa.Slen
	for i := 0; i < len(sa.raw.Data); i++ {
		sa.raw.Data[i] = sa.Data[i]
	}
	return unsafe.Pointer(&sa.raw), SizeofSockaddrDatalink, nil
}

func anyToSockaddr(rsa *RawSockaddrAny) (Sockaddr, error) {
	switch rsa.Addr.Family {
	case AF_LINK:
		pp := (*RawSockaddrDatalink)(unsafe.Pointer(rsa))
		sa := new(SockaddrDatalink)
		sa.Len = pp.Len
		sa.Family = pp.Family
		sa.Index = pp.Index
		sa.Type = pp.Type
		sa.Nlen = pp.Nlen
		sa.Alen = pp.Alen
		sa.Slen = pp.Slen
		for i := 0; i < len(sa.Data); i++ {
			sa.Data[i] = pp.Data[i]
		}
		return sa, nil

	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		if pp.Len < 2 || pp.Len > SizeofSockaddrUnix {
			return nil, EINVAL
		}
		sa := new(SockaddrUnix)

		
		
		
		
		n := int(pp.Len) - 2 
		for i := 0; i < n; i++ {
			if pp.Path[i] == 0 {
				
				
				n = i
				break
			}
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
	if err != nil {
		return
	}
	if runtime.GOOS == "darwin" && len == 0 {
		
		
		
		
		Close(nfd)
		return 0, nil, ECONNABORTED
	}
	sa, err = anyToSockaddr(&rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
}

func Getsockname(fd int) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil {
		return
	}
	
	
	if runtime.GOOS == "dragonfly" && rsa.Addr.Family == AF_UNSPEC && rsa.Addr.Len == 0 {
		rsa.Addr.Family = AF_UNIX
		rsa.Addr.Len = SizeofSockaddrUnix
	}
	return anyToSockaddr(&rsa)
}



func GetsockoptByte(fd, level, opt int) (value byte, err error) {
	var n byte
	vallen := _Socklen(1)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&n), &vallen)
	return n, err
}

func GetsockoptInet4Addr(fd, level, opt int) (value [4]byte, err error) {
	vallen := _Socklen(4)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&value[0]), &vallen)
	return value, err
}

func GetsockoptIPMreq(fd, level, opt int) (*IPMreq, error) {
	var value IPMreq
	vallen := _Socklen(SizeofIPMreq)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptIPv6Mreq(fd, level, opt int) (*IPv6Mreq, error) {
	var value IPv6Mreq
	vallen := _Socklen(SizeofIPv6Mreq)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptIPv6MTUInfo(fd, level, opt int) (*IPv6MTUInfo, error) {
	var value IPv6MTUInfo
	vallen := _Socklen(SizeofIPv6MTUInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptICMPv6Filter(fd, level, opt int) (*ICMPv6Filter, error) {
	var value ICMPv6Filter
	vallen := _Socklen(SizeofICMPv6Filter)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
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





func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error) {
	var msg Msghdr
	var rsa RawSockaddrAny
	msg.Name = (*byte)(unsafe.Pointer(&rsa))
	msg.Namelen = uint32(SizeofSockaddrAny)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
		msg.SetControllen(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = recvmsg(fd, &msg, flags); err != nil {
		return
	}
	oobn = int(msg.Controllen)
	recvflags = int(msg.Flags)
	
	if rsa.Addr.Family != AF_UNSPEC {
		from, err = anyToSockaddr(&rsa)
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
		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
		msg.SetControllen(len(oob))
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



func Kevent(kq int, changes, events []Kevent_t, timeout *Timespec) (n int, err error) {
	var change, event unsafe.Pointer
	if len(changes) > 0 {
		change = unsafe.Pointer(&changes[0])
	}
	if len(events) > 0 {
		event = unsafe.Pointer(&events[0])
	}
	return kevent(kq, change, len(changes), event, len(events), timeout)
}




func sysctlmib(name string, args ...int) ([]_C_int, error) {
	
	mib, err := nametomib(name)
	if err != nil {
		return nil, err
	}

	for _, a := range args {
		mib = append(mib, _C_int(a))
	}

	return mib, nil
}

func Sysctl(name string) (string, error) {
	return SysctlArgs(name)
}

func SysctlArgs(name string, args ...int) (string, error) {
	buf, err := SysctlRaw(name, args...)
	if err != nil {
		return "", err
	}
	n := len(buf)

	
	if n > 0 && buf[n-1] == '\x00' {
		n--
	}
	return string(buf[0:n]), nil
}

func SysctlUint32(name string) (uint32, error) {
	return SysctlUint32Args(name)
}

func SysctlUint32Args(name string, args ...int) (uint32, error) {
	mib, err := sysctlmib(name, args...)
	if err != nil {
		return 0, err
	}

	n := uintptr(4)
	buf := make([]byte, 4)
	if err := sysctl(mib, &buf[0], &n, nil, 0); err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, EIO
	}
	return *(*uint32)(unsafe.Pointer(&buf[0])), nil
}

func SysctlUint64(name string, args ...int) (uint64, error) {
	mib, err := sysctlmib(name, args...)
	if err != nil {
		return 0, err
	}

	n := uintptr(8)
	buf := make([]byte, 8)
	if err := sysctl(mib, &buf[0], &n, nil, 0); err != nil {
		return 0, err
	}
	if n != 8 {
		return 0, EIO
	}
	return *(*uint64)(unsafe.Pointer(&buf[0])), nil
}

func SysctlRaw(name string, args ...int) ([]byte, error) {
	mib, err := sysctlmib(name, args...)
	if err != nil {
		return nil, err
	}

	
	n := uintptr(0)
	if err := sysctl(mib, nil, &n, nil, 0); err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	
	buf := make([]byte, n)
	if err := sysctl(mib, &buf[0], &n, nil, 0); err != nil {
		return nil, err
	}

	
	
	return buf[:n], nil
}



func Utimes(path string, tv []Timeval) error {
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
		err := utimensat(AT_FDCWD, path, nil, 0)
		if err != ENOSYS {
			return err
		}
		return utimes(path, nil)
	}
	if len(ts) != 2 {
		return EINVAL
	}
	
	err := setattrlistTimes(path, ts, 0)
	if err != ENOSYS {
		return err
	}
	err = utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
	if err != ENOSYS {
		return err
	}
	
	
	tv := [2]Timeval{
		NsecToTimeval(TimespecToNsec(ts[0])),
		NsecToTimeval(TimespecToNsec(ts[1])),
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

func UtimesNanoAt(dirfd int, path string, ts []Timespec, flags int) error {
	if ts == nil {
		return utimensat(dirfd, path, nil, flags)
	}
	if len(ts) != 2 {
		return EINVAL
	}
	err := setattrlistTimes(path, ts, flags)
	if err != ENOSYS {
		return err
	}
	return utimensat(dirfd, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), flags)
}



func Futimes(fd int, tv []Timeval) error {
	if tv == nil {
		return futimes(fd, nil)
	}
	if len(tv) != 2 {
		return EINVAL
	}
	return futimes(fd, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}





func Poll(fds []PollFd, timeout int) (n int, err error) {
	if len(fds) == 0 {
		return poll(nil, 0, timeout)
	}
	return poll(&fds[0], len(fds), timeout)
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








