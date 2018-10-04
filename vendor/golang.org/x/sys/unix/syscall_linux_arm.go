





package unix

import (
	"syscall"
	"unsafe"
)

func setTimespec(sec, nsec int64) Timespec {
	return Timespec{Sec: int32(sec), Nsec: int32(nsec)}
}

func setTimeval(sec, usec int64) Timeval {
	return Timeval{Sec: int32(sec), Usec: int32(usec)}
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



func seek(fd int, offset int64, whence int) (newoffset int64, err syscall.Errno)

func Seek(fd int, offset int64, whence int) (newoffset int64, err error) {
	newoffset, errno := seek(fd, offset, whence)
	if errno != 0 {
		return 0, errno
	}
	return newoffset, nil
}



















































func Time(t *Time_t) (Time_t, error) {
	var tv Timeval
	err := Gettimeofday(&tv)
	if err != nil {
		return 0, err
	}
	if t != nil {
		*t = Time_t(tv.Sec)
	}
	return Time_t(tv.Sec), nil
}

func Utime(path string, buf *Utimbuf) error {
	tv := []Timeval{
		{Sec: buf.Actime},
		{Sec: buf.Modtime},
	}
	return Utimes(path, tv)
}








func Fadvise(fd int, offset int64, length int64, advice int) (err error) {
	_, _, e1 := Syscall6(SYS_ARM_FADVISE64_64, uintptr(fd), uintptr(advice), uintptr(offset), uintptr(offset>>32), uintptr(length), uintptr(length>>32))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fstatfs(fd int, buf *Statfs_t) (err error) {
	_, _, e := Syscall(SYS_FSTATFS64, uintptr(fd), unsafe.Sizeof(*buf), uintptr(unsafe.Pointer(buf)))
	if e != 0 {
		err = e
	}
	return
}

func Statfs(path string, buf *Statfs_t) (err error) {
	pathp, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_STATFS64, uintptr(unsafe.Pointer(pathp)), unsafe.Sizeof(*buf), uintptr(unsafe.Pointer(buf)))
	if e != 0 {
		err = e
	}
	return
}

func mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error) {
	page := uintptr(offset / 4096)
	if offset != int64(page)*4096 {
		return 0, EINVAL
	}
	return mmap2(addr, length, prot, flags, fd, page)
}

type rlimit32 struct {
	Cur uint32
	Max uint32
}



const rlimInf32 = ^uint32(0)
const rlimInf64 = ^uint64(0)

func Getrlimit(resource int, rlim *Rlimit) (err error) {
	err = prlimit(0, resource, nil, rlim)
	if err != ENOSYS {
		return err
	}

	rl := rlimit32{}
	err = getrlimit(resource, &rl)
	if err != nil {
		return
	}

	if rl.Cur == rlimInf32 {
		rlim.Cur = rlimInf64
	} else {
		rlim.Cur = uint64(rl.Cur)
	}

	if rl.Max == rlimInf32 {
		rlim.Max = rlimInf64
	} else {
		rlim.Max = uint64(rl.Max)
	}
	return
}



func Setrlimit(resource int, rlim *Rlimit) (err error) {
	err = prlimit(0, resource, rlim, nil)
	if err != ENOSYS {
		return err
	}

	rl := rlimit32{}
	if rlim.Cur == rlimInf64 {
		rl.Cur = rlimInf32
	} else if rlim.Cur < uint64(rlimInf32) {
		rl.Cur = uint32(rlim.Cur)
	} else {
		return EINVAL
	}
	if rlim.Max == rlimInf64 {
		rl.Max = rlimInf32
	} else if rlim.Max < uint64(rlimInf32) {
		rl.Max = uint32(rlim.Max)
	} else {
		return EINVAL
	}

	return setrlimit(resource, &rl)
}

func (r *PtraceRegs) PC() uint64 { return uint64(r.Uregs[15]) }

func (r *PtraceRegs) SetPC(pc uint64) { r.Uregs[15] = uint32(pc) }

func (iov *Iovec) SetLen(length int) {
	iov.Len = uint32(length)
}

func (msghdr *Msghdr) SetControllen(length int) {
	msghdr.Controllen = uint32(length)
}

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = uint32(length)
}



func Poll(fds []PollFd, timeout int) (n int, err error) {
	if len(fds) == 0 {
		return poll(nil, 0, timeout)
	}
	return poll(&fds[0], len(fds), timeout)
}
