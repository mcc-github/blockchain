






package unix



















func Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error) {
	var ts *Timespec
	if timeout != nil {
		ts = &Timespec{Sec: timeout.Sec, Nsec: timeout.Usec * 1000}
	}
	return Pselect(nfd, r, w, e, ts, nil)
}


































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

type stat_t struct {
	Dev        uint32
	Pad0       [3]int32
	Ino        uint64
	Mode       uint32
	Nlink      uint32
	Uid        uint32
	Gid        uint32
	Rdev       uint32
	Pad1       [3]uint32
	Size       int64
	Atime      uint32
	Atime_nsec uint32
	Mtime      uint32
	Mtime_nsec uint32
	Ctime      uint32
	Ctime_nsec uint32
	Blksize    uint32
	Pad2       uint32
	Blocks     int64
}





func Fstat(fd int, s *Stat_t) (err error) {
	st := &stat_t{}
	err = fstat(fd, st)
	fillStat_t(s, st)
	return
}

func Lstat(path string, s *Stat_t) (err error) {
	st := &stat_t{}
	err = lstat(path, st)
	fillStat_t(s, st)
	return
}

func Stat(path string, s *Stat_t) (err error) {
	st := &stat_t{}
	err = stat(path, st)
	fillStat_t(s, st)
	return
}

func fillStat_t(s *Stat_t, st *stat_t) {
	s.Dev = st.Dev
	s.Ino = st.Ino
	s.Mode = st.Mode
	s.Nlink = st.Nlink
	s.Uid = st.Uid
	s.Gid = st.Gid
	s.Rdev = st.Rdev
	s.Size = st.Size
	s.Atim = Timespec{int64(st.Atime), int64(st.Atime_nsec)}
	s.Mtim = Timespec{int64(st.Mtime), int64(st.Mtime_nsec)}
	s.Ctim = Timespec{int64(st.Ctime), int64(st.Ctime_nsec)}
	s.Blksize = st.Blksize
	s.Blocks = st.Blocks
}

func (r *PtraceRegs) PC() uint64 { return r.Epc }

func (r *PtraceRegs) SetPC(pc uint64) { r.Epc = pc }

func (iov *Iovec) SetLen(length int) {
	iov.Len = uint64(length)
}

func (msghdr *Msghdr) SetControllen(length int) {
	msghdr.Controllen = uint64(length)
}

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = uint64(length)
}



func Poll(fds []PollFd, timeout int) (n int, err error) {
	if len(fds) == 0 {
		return poll(nil, 0, timeout)
	}
	return poll(&fds[0], len(fds), timeout)
}
