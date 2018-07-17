





package unix

















func Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error) {
	var ts *Timespec
	if timeout != nil {
		ts = &Timespec{Sec: timeout.Sec, Nsec: timeout.Usec * 1000}
	}
	return Pselect(nfd, r, w, e, ts, nil)
}












func Stat(path string, stat *Stat_t) (err error) {
	return Fstatat(AT_FDCWD, path, stat, 0)
}

func Lchown(path string, uid int, gid int) (err error) {
	return Fchownat(AT_FDCWD, path, uid, gid, AT_SYMLINK_NOFOLLOW)
}

func Lstat(path string, stat *Stat_t) (err error) {
	return Fstatat(AT_FDCWD, path, stat, AT_SYMLINK_NOFOLLOW)
}
























func setTimespec(sec, nsec int64) Timespec {
	return Timespec{Sec: sec, Nsec: nsec}
}

func setTimeval(sec, usec int64) Timeval {
	return Timeval{Sec: sec, Usec: usec}
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

func (r *PtraceRegs) PC() uint64 { return r.Pc }

func (r *PtraceRegs) SetPC(pc uint64) { r.Pc = pc }

func (iov *Iovec) SetLen(length int) {
	iov.Len = uint64(length)
}

func (msghdr *Msghdr) SetControllen(length int) {
	msghdr.Controllen = uint64(length)
}

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = uint64(length)
}

func InotifyInit() (fd int, err error) {
	return InotifyInit1(0)
}

func Dup2(oldfd int, newfd int) (err error) {
	return Dup3(oldfd, newfd, 0)
}

func Pause() (err error) {
	_, _, e1 := Syscall6(SYS_PPOLL, 0, 0, 0, 0, 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}




const (
	SYS_GETPGRP      = 1060
	SYS_UTIMES       = 1037
	SYS_FUTIMESAT    = 1066
	SYS_PAUSE        = 1061
	SYS_USTAT        = 1070
	SYS_UTIME        = 1063
	SYS_LCHOWN       = 1032
	SYS_TIME         = 1062
	SYS_EPOLL_CREATE = 1042
	SYS_EPOLL_WAIT   = 1069
)

func Poll(fds []PollFd, timeout int) (n int, err error) {
	var ts *Timespec
	if timeout >= 0 {
		ts = new(Timespec)
		*ts = NsecToTimespec(int64(timeout) * 1e6)
	}
	if len(fds) == 0 {
		return ppoll(nil, 0, ts, nil)
	}
	return ppoll(&fds[0], len(fds), ts, nil)
}
