




package unix


import "C"
import (
	"syscall"
	"unsafe"
)



func utimes(path string, times *[2]Timeval) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.utimes(C.uintptr_t(_p0), C.uintptr_t(uintptr(unsafe.Pointer(times))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func utimensat(dirfd int, path string, times *[2]Timespec, flag int) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.utimensat(C.int(dirfd), C.uintptr_t(_p0), C.uintptr_t(uintptr(unsafe.Pointer(times))), C.int(flag))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func getcwd(buf []byte) (err error) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	var _p1 int
	_p1 = len(buf)
	r0, er := C.getcwd(C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, err error) {
	r0, er := C.accept(C.int(s), C.uintptr_t(uintptr(unsafe.Pointer(rsa))), C.uintptr_t(uintptr(unsafe.Pointer(addrlen))))
	fd = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func getdirent(fd int, buf []byte) (n int, err error) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	var _p1 int
	_p1 = len(buf)
	r0, er := C.getdirent(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func wait4(pid Pid_t, status *_C_int, options int, rusage *Rusage) (wpid Pid_t, err error) {
	r0, er := C.wait4(C.int(pid), C.uintptr_t(uintptr(unsafe.Pointer(status))), C.int(options), C.uintptr_t(uintptr(unsafe.Pointer(rusage))))
	wpid = Pid_t(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func ioctl(fd int, req uint, arg uintptr) (err error) {
	r0, er := C.ioctl(C.int(fd), C.int(req), C.uintptr_t(arg))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func FcntlInt(fd uintptr, cmd int, arg int) (r int, err error) {
	r0, er := C.fcntl(C.uintptr_t(fd), C.int(cmd), C.uintptr_t(arg))
	r = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) (err error) {
	r0, er := C.fcntl(C.uintptr_t(fd), C.int(cmd), C.uintptr_t(uintptr(unsafe.Pointer(lk))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Acct(path string) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.acct(C.uintptr_t(_p0))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Chdir(path string) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.chdir(C.uintptr_t(_p0))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Chroot(path string) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.chroot(C.uintptr_t(_p0))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Close(fd int) (err error) {
	r0, er := C.close(C.int(fd))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Dup(oldfd int) (fd int, err error) {
	r0, er := C.dup(C.int(oldfd))
	fd = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Dup3(oldfd int, newfd int, flags int) (err error) {
	r0, er := C.dup3(C.int(oldfd), C.int(newfd), C.int(flags))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Exit(code int) {
	C.exit(C.int(code))
	return
}



func Faccessat(dirfd int, path string, mode uint32, flags int) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.faccessat(C.int(dirfd), C.uintptr_t(_p0), C.uint(mode), C.int(flags))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fallocate(fd int, mode uint32, off int64, len int64) (err error) {
	r0, er := C.fallocate(C.int(fd), C.uint(mode), C.longlong(off), C.longlong(len))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fchdir(fd int) (err error) {
	r0, er := C.fchdir(C.int(fd))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fchmod(fd int, mode uint32) (err error) {
	r0, er := C.fchmod(C.int(fd), C.uint(mode))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fchmodat(dirfd int, path string, mode uint32, flags int) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.fchmodat(C.int(dirfd), C.uintptr_t(_p0), C.uint(mode), C.int(flags))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fchownat(dirfd int, path string, uid int, gid int, flags int) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.fchownat(C.int(dirfd), C.uintptr_t(_p0), C.int(uid), C.int(gid), C.int(flags))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func fcntl(fd int, cmd int, arg int) (val int, err error) {
	r0, er := C.fcntl(C.uintptr_t(fd), C.int(cmd), C.uintptr_t(arg))
	val = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fdatasync(fd int) (err error) {
	r0, er := C.fdatasync(C.int(fd))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fsync(fd int) (err error) {
	r0, er := C.fsync(C.int(fd))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Getpgid(pid int) (pgid int, err error) {
	r0, er := C.getpgid(C.int(pid))
	pgid = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Getpgrp() (pid int) {
	r0, _ := C.getpgrp()
	pid = int(r0)
	return
}



func Getpid() (pid int) {
	r0, _ := C.getpid()
	pid = int(r0)
	return
}



func Getppid() (ppid int) {
	r0, _ := C.getppid()
	ppid = int(r0)
	return
}



func Getpriority(which int, who int) (prio int, err error) {
	r0, er := C.getpriority(C.int(which), C.int(who))
	prio = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Getrusage(who int, rusage *Rusage) (err error) {
	r0, er := C.getrusage(C.int(who), C.uintptr_t(uintptr(unsafe.Pointer(rusage))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Getsid(pid int) (sid int, err error) {
	r0, er := C.getsid(C.int(pid))
	sid = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Kill(pid int, sig syscall.Signal) (err error) {
	r0, er := C.kill(C.int(pid), C.int(sig))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Klogctl(typ int, buf []byte) (n int, err error) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	var _p1 int
	_p1 = len(buf)
	r0, er := C.syslog(C.int(typ), C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Mkdir(dirfd int, path string, mode uint32) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.mkdir(C.int(dirfd), C.uintptr_t(_p0), C.uint(mode))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Mkdirat(dirfd int, path string, mode uint32) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.mkdirat(C.int(dirfd), C.uintptr_t(_p0), C.uint(mode))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Mkfifo(path string, mode uint32) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.mkfifo(C.uintptr_t(_p0), C.uint(mode))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Mknod(path string, mode uint32, dev int) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.mknod(C.uintptr_t(_p0), C.uint(mode), C.int(dev))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Mknodat(dirfd int, path string, mode uint32, dev int) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.mknodat(C.int(dirfd), C.uintptr_t(_p0), C.uint(mode), C.int(dev))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Nanosleep(time *Timespec, leftover *Timespec) (err error) {
	r0, er := C.nanosleep(C.uintptr_t(uintptr(unsafe.Pointer(time))), C.uintptr_t(uintptr(unsafe.Pointer(leftover))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Open(path string, mode int, perm uint32) (fd int, err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.open64(C.uintptr_t(_p0), C.int(mode), C.uint(perm))
	fd = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Openat(dirfd int, path string, flags int, mode uint32) (fd int, err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.openat(C.int(dirfd), C.uintptr_t(_p0), C.int(flags), C.uint(mode))
	fd = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func read(fd int, p []byte) (n int, err error) {
	var _p0 *byte
	if len(p) > 0 {
		_p0 = &p[0]
	}
	var _p1 int
	_p1 = len(p)
	r0, er := C.read(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Readlink(path string, buf []byte) (n int, err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	var _p1 *byte
	if len(buf) > 0 {
		_p1 = &buf[0]
	}
	var _p2 int
	_p2 = len(buf)
	r0, er := C.readlink(C.uintptr_t(_p0), C.uintptr_t(uintptr(unsafe.Pointer(_p1))), C.size_t(_p2))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Removexattr(path string, attr string) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	_p1 := uintptr(unsafe.Pointer(C.CString(attr)))
	r0, er := C.removexattr(C.uintptr_t(_p0), C.uintptr_t(_p1))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(oldpath)))
	_p1 := uintptr(unsafe.Pointer(C.CString(newpath)))
	r0, er := C.renameat(C.int(olddirfd), C.uintptr_t(_p0), C.int(newdirfd), C.uintptr_t(_p1))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Setdomainname(p []byte) (err error) {
	var _p0 *byte
	if len(p) > 0 {
		_p0 = &p[0]
	}
	var _p1 int
	_p1 = len(p)
	r0, er := C.setdomainname(C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Sethostname(p []byte) (err error) {
	var _p0 *byte
	if len(p) > 0 {
		_p0 = &p[0]
	}
	var _p1 int
	_p1 = len(p)
	r0, er := C.sethostname(C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Setpgid(pid int, pgid int) (err error) {
	r0, er := C.setpgid(C.int(pid), C.int(pgid))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Setsid() (pid int, err error) {
	r0, er := C.setsid()
	pid = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Settimeofday(tv *Timeval) (err error) {
	r0, er := C.settimeofday(C.uintptr_t(uintptr(unsafe.Pointer(tv))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Setuid(uid int) (err error) {
	r0, er := C.setuid(C.int(uid))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Setgid(uid int) (err error) {
	r0, er := C.setgid(C.int(uid))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Setpriority(which int, who int, prio int) (err error) {
	r0, er := C.setpriority(C.int(which), C.int(who), C.int(prio))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Statx(dirfd int, path string, flags int, mask int, stat *Statx_t) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.statx(C.int(dirfd), C.uintptr_t(_p0), C.int(flags), C.int(mask), C.uintptr_t(uintptr(unsafe.Pointer(stat))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Sync() {
	C.sync()
	return
}



func Tee(rfd int, wfd int, len int, flags int) (n int64, err error) {
	r0, er := C.tee(C.int(rfd), C.int(wfd), C.int(len), C.int(flags))
	n = int64(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Times(tms *Tms) (ticks uintptr, err error) {
	r0, er := C.times(C.uintptr_t(uintptr(unsafe.Pointer(tms))))
	ticks = uintptr(r0)
	if uintptr(r0) == ^uintptr(0) && er != nil {
		err = er
	}
	return
}



func Umask(mask int) (oldmask int) {
	r0, _ := C.umask(C.int(mask))
	oldmask = int(r0)
	return
}



func Uname(buf *Utsname) (err error) {
	r0, er := C.uname(C.uintptr_t(uintptr(unsafe.Pointer(buf))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Unlink(path string) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.unlink(C.uintptr_t(_p0))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Unlinkat(dirfd int, path string, flags int) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.unlinkat(C.int(dirfd), C.uintptr_t(_p0), C.int(flags))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Unshare(flags int) (err error) {
	r0, er := C.unshare(C.int(flags))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Ustat(dev int, ubuf *Ustat_t) (err error) {
	r0, er := C.ustat(C.int(dev), C.uintptr_t(uintptr(unsafe.Pointer(ubuf))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func write(fd int, p []byte) (n int, err error) {
	var _p0 *byte
	if len(p) > 0 {
		_p0 = &p[0]
	}
	var _p1 int
	_p1 = len(p)
	r0, er := C.write(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func readlen(fd int, p *byte, np int) (n int, err error) {
	r0, er := C.read(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(p))), C.size_t(np))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func writelen(fd int, p *byte, np int) (n int, err error) {
	r0, er := C.write(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(p))), C.size_t(np))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Dup2(oldfd int, newfd int) (err error) {
	r0, er := C.dup2(C.int(oldfd), C.int(newfd))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fadvise(fd int, offset int64, length int64, advice int) (err error) {
	r0, er := C.posix_fadvise64(C.int(fd), C.longlong(offset), C.longlong(length), C.int(advice))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fchown(fd int, uid int, gid int) (err error) {
	r0, er := C.fchown(C.int(fd), C.int(uid), C.int(gid))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fstat(fd int, stat *Stat_t) (err error) {
	r0, er := C.fstat(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(stat))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fstatat(dirfd int, path string, stat *Stat_t, flags int) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.fstatat(C.int(dirfd), C.uintptr_t(_p0), C.uintptr_t(uintptr(unsafe.Pointer(stat))), C.int(flags))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Fstatfs(fd int, buf *Statfs_t) (err error) {
	r0, er := C.fstatfs(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(buf))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Ftruncate(fd int, length int64) (err error) {
	r0, er := C.ftruncate(C.int(fd), C.longlong(length))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Getegid() (egid int) {
	r0, _ := C.getegid()
	egid = int(r0)
	return
}



func Geteuid() (euid int) {
	r0, _ := C.geteuid()
	euid = int(r0)
	return
}



func Getgid() (gid int) {
	r0, _ := C.getgid()
	gid = int(r0)
	return
}



func Getuid() (uid int) {
	r0, _ := C.getuid()
	uid = int(r0)
	return
}



func Lchown(path string, uid int, gid int) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.lchown(C.uintptr_t(_p0), C.int(uid), C.int(gid))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Listen(s int, n int) (err error) {
	r0, er := C.listen(C.int(s), C.int(n))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Lstat(path string, stat *Stat_t) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.lstat(C.uintptr_t(_p0), C.uintptr_t(uintptr(unsafe.Pointer(stat))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Pause() (err error) {
	r0, er := C.pause()
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Pread(fd int, p []byte, offset int64) (n int, err error) {
	var _p0 *byte
	if len(p) > 0 {
		_p0 = &p[0]
	}
	var _p1 int
	_p1 = len(p)
	r0, er := C.pread64(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1), C.longlong(offset))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Pwrite(fd int, p []byte, offset int64) (n int, err error) {
	var _p0 *byte
	if len(p) > 0 {
		_p0 = &p[0]
	}
	var _p1 int
	_p1 = len(p)
	r0, er := C.pwrite64(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1), C.longlong(offset))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Pselect(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timespec, sigmask *Sigset_t) (n int, err error) {
	r0, er := C.pselect(C.int(nfd), C.uintptr_t(uintptr(unsafe.Pointer(r))), C.uintptr_t(uintptr(unsafe.Pointer(w))), C.uintptr_t(uintptr(unsafe.Pointer(e))), C.uintptr_t(uintptr(unsafe.Pointer(timeout))), C.uintptr_t(uintptr(unsafe.Pointer(sigmask))))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Setregid(rgid int, egid int) (err error) {
	r0, er := C.setregid(C.int(rgid), C.int(egid))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Setreuid(ruid int, euid int) (err error) {
	r0, er := C.setreuid(C.int(ruid), C.int(euid))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Shutdown(fd int, how int) (err error) {
	r0, er := C.shutdown(C.int(fd), C.int(how))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Splice(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int64, err error) {
	r0, er := C.splice(C.int(rfd), C.uintptr_t(uintptr(unsafe.Pointer(roff))), C.int(wfd), C.uintptr_t(uintptr(unsafe.Pointer(woff))), C.int(len), C.int(flags))
	n = int64(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Stat(path string, stat *Stat_t) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.stat(C.uintptr_t(_p0), C.uintptr_t(uintptr(unsafe.Pointer(stat))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Statfs(path string, buf *Statfs_t) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.statfs(C.uintptr_t(_p0), C.uintptr_t(uintptr(unsafe.Pointer(buf))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Truncate(path string, length int64) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.truncate(C.uintptr_t(_p0), C.longlong(length))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func bind(s int, addr unsafe.Pointer, addrlen _Socklen) (err error) {
	r0, er := C.bind(C.int(s), C.uintptr_t(uintptr(addr)), C.uintptr_t(uintptr(addrlen)))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func connect(s int, addr unsafe.Pointer, addrlen _Socklen) (err error) {
	r0, er := C.connect(C.int(s), C.uintptr_t(uintptr(addr)), C.uintptr_t(uintptr(addrlen)))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func getgroups(n int, list *_Gid_t) (nn int, err error) {
	r0, er := C.getgroups(C.int(n), C.uintptr_t(uintptr(unsafe.Pointer(list))))
	nn = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func setgroups(n int, list *_Gid_t) (err error) {
	r0, er := C.setgroups(C.int(n), C.uintptr_t(uintptr(unsafe.Pointer(list))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) (err error) {
	r0, er := C.getsockopt(C.int(s), C.int(level), C.int(name), C.uintptr_t(uintptr(val)), C.uintptr_t(uintptr(unsafe.Pointer(vallen))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error) {
	r0, er := C.setsockopt(C.int(s), C.int(level), C.int(name), C.uintptr_t(uintptr(val)), C.uintptr_t(vallen))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func socket(domain int, typ int, proto int) (fd int, err error) {
	r0, er := C.socket(C.int(domain), C.int(typ), C.int(proto))
	fd = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func socketpair(domain int, typ int, proto int, fd *[2]int32) (err error) {
	r0, er := C.socketpair(C.int(domain), C.int(typ), C.int(proto), C.uintptr_t(uintptr(unsafe.Pointer(fd))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func getpeername(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error) {
	r0, er := C.getpeername(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(rsa))), C.uintptr_t(uintptr(unsafe.Pointer(addrlen))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func getsockname(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error) {
	r0, er := C.getsockname(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(rsa))), C.uintptr_t(uintptr(unsafe.Pointer(addrlen))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func recvfrom(fd int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (n int, err error) {
	var _p0 *byte
	if len(p) > 0 {
		_p0 = &p[0]
	}
	var _p1 int
	_p1 = len(p)
	r0, er := C.recvfrom(C.int(fd), C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1), C.int(flags), C.uintptr_t(uintptr(unsafe.Pointer(from))), C.uintptr_t(uintptr(unsafe.Pointer(fromlen))))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func sendto(s int, buf []byte, flags int, to unsafe.Pointer, addrlen _Socklen) (err error) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	var _p1 int
	_p1 = len(buf)
	r0, er := C.sendto(C.int(s), C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1), C.int(flags), C.uintptr_t(uintptr(to)), C.uintptr_t(uintptr(addrlen)))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func recvmsg(s int, msg *Msghdr, flags int) (n int, err error) {
	r0, er := C.recvmsg(C.int(s), C.uintptr_t(uintptr(unsafe.Pointer(msg))), C.int(flags))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func sendmsg(s int, msg *Msghdr, flags int) (n int, err error) {
	r0, er := C.sendmsg(C.int(s), C.uintptr_t(uintptr(unsafe.Pointer(msg))), C.int(flags))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func munmap(addr uintptr, length uintptr) (err error) {
	r0, er := C.munmap(C.uintptr_t(addr), C.uintptr_t(length))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Madvise(b []byte, advice int) (err error) {
	var _p0 *byte
	if len(b) > 0 {
		_p0 = &b[0]
	}
	var _p1 int
	_p1 = len(b)
	r0, er := C.madvise(C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1), C.int(advice))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Mprotect(b []byte, prot int) (err error) {
	var _p0 *byte
	if len(b) > 0 {
		_p0 = &b[0]
	}
	var _p1 int
	_p1 = len(b)
	r0, er := C.mprotect(C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1), C.int(prot))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Mlock(b []byte) (err error) {
	var _p0 *byte
	if len(b) > 0 {
		_p0 = &b[0]
	}
	var _p1 int
	_p1 = len(b)
	r0, er := C.mlock(C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Mlockall(flags int) (err error) {
	r0, er := C.mlockall(C.int(flags))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Msync(b []byte, flags int) (err error) {
	var _p0 *byte
	if len(b) > 0 {
		_p0 = &b[0]
	}
	var _p1 int
	_p1 = len(b)
	r0, er := C.msync(C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1), C.int(flags))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Munlock(b []byte) (err error) {
	var _p0 *byte
	if len(b) > 0 {
		_p0 = &b[0]
	}
	var _p1 int
	_p1 = len(b)
	r0, er := C.munlock(C.uintptr_t(uintptr(unsafe.Pointer(_p0))), C.size_t(_p1))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Munlockall() (err error) {
	r0, er := C.munlockall()
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func pipe(p *[2]_C_int) (err error) {
	r0, er := C.pipe(C.uintptr_t(uintptr(unsafe.Pointer(p))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func pipe2(p *[2]_C_int, flags int) (err error) {
	r0, er := C.pipe2(C.uintptr_t(uintptr(unsafe.Pointer(p))), C.int(flags))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func poll(fds *PollFd, nfds int, timeout int) (n int, err error) {
	r0, er := C.poll(C.uintptr_t(uintptr(unsafe.Pointer(fds))), C.int(nfds), C.int(timeout))
	n = int(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func gettimeofday(tv *Timeval, tzp *Timezone) (err error) {
	r0, er := C.gettimeofday(C.uintptr_t(uintptr(unsafe.Pointer(tv))), C.uintptr_t(uintptr(unsafe.Pointer(tzp))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Time(t *Time_t) (tt Time_t, err error) {
	r0, er := C.time(C.uintptr_t(uintptr(unsafe.Pointer(t))))
	tt = Time_t(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Utime(path string, buf *Utimbuf) (err error) {
	_p0 := uintptr(unsafe.Pointer(C.CString(path)))
	r0, er := C.utime(C.uintptr_t(_p0), C.uintptr_t(uintptr(unsafe.Pointer(buf))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Getrlimit(resource int, rlim *Rlimit) (err error) {
	r0, er := C.getrlimit64(C.int(resource), C.uintptr_t(uintptr(unsafe.Pointer(rlim))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Setrlimit(resource int, rlim *Rlimit) (err error) {
	r0, er := C.setrlimit64(C.int(resource), C.uintptr_t(uintptr(unsafe.Pointer(rlim))))
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func Seek(fd int, offset int64, whence int) (off int64, err error) {
	r0, er := C.lseek64(C.int(fd), C.longlong(offset), C.int(whence))
	off = int64(r0)
	if r0 == -1 && er != nil {
		err = er
	}
	return
}



func mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error) {
	r0, er := C.mmap(C.uintptr_t(addr), C.uintptr_t(length), C.int(prot), C.int(flags), C.int(fd), C.longlong(offset))
	xaddr = uintptr(r0)
	if uintptr(r0) == ^uintptr(0) && er != nil {
		err = er
	}
	return
}
