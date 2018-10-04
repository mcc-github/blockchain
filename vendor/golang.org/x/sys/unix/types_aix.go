











package unix


import "C"



const (
	sizeofPtr      = C.sizeofPtr
	sizeofShort    = C.sizeof_short
	sizeofInt      = C.sizeof_int
	sizeofLong     = C.sizeof_long
	sizeofLongLong = C.sizeof_longlong
	PathMax        = C.PATH_MAX
)



type (
	_C_short     C.short
	_C_int       C.int
	_C_long      C.long
	_C_long_long C.longlong
)

type off64 C.off64_t
type off C.off_t
type Mode_t C.mode_t



type Timespec C.struct_timespec

type StTimespec C.struct_st_timespec

type Timeval C.struct_timeval

type Timeval32 C.struct_timeval32

type Timex C.struct_timex

type Time_t C.time_t

type Tms C.struct_tms

type Utimbuf C.struct_utimbuf

type Timezone C.struct_timezone



type Rusage C.struct_rusage

type Rlimit C.struct_rlimit64

type Pid_t C.pid_t

type _Gid_t C.gid_t

type dev_t C.dev_t



type Stat_t C.struct_stat

type StatxTimestamp C.struct_statx_timestamp

type Statx_t C.struct_statx

type Dirent C.struct_dirent



type RawSockaddrInet4 C.struct_sockaddr_in

type RawSockaddrInet6 C.struct_sockaddr_in6

type RawSockaddrUnix C.struct_sockaddr_un

type RawSockaddr C.struct_sockaddr

type RawSockaddrAny C.struct_sockaddr_any

type _Socklen C.socklen_t

type Cmsghdr C.struct_cmsghdr

type ICMPv6Filter C.struct_icmp6_filter

type Iovec C.struct_iovec

type IPMreq C.struct_ip_mreq

type IPv6Mreq C.struct_ipv6_mreq

type IPv6MTUInfo C.struct_ip6_mtuinfo

type Linger C.struct_linger

type Msghdr C.struct_msghdr

const (
	SizeofSockaddrInet4 = C.sizeof_struct_sockaddr_in
	SizeofSockaddrInet6 = C.sizeof_struct_sockaddr_in6
	SizeofSockaddrAny   = C.sizeof_struct_sockaddr_any
	SizeofSockaddrUnix  = C.sizeof_struct_sockaddr_un
	SizeofLinger        = C.sizeof_struct_linger
	SizeofIPMreq        = C.sizeof_struct_ip_mreq
	SizeofIPv6Mreq      = C.sizeof_struct_ipv6_mreq
	SizeofIPv6MTUInfo   = C.sizeof_struct_ip6_mtuinfo
	SizeofMsghdr        = C.sizeof_struct_msghdr
	SizeofCmsghdr       = C.sizeof_struct_cmsghdr
	SizeofICMPv6Filter  = C.sizeof_struct_icmp6_filter
)



const (
	SizeofIfMsghdr = C.sizeof_struct_if_msghdr
)

type IfMsgHdr C.struct_if_msghdr



type FdSet C.fd_set

type Utsname C.struct_utsname

type Ustat_t C.struct_ustat

type Sigset_t C.sigset_t

const (
	AT_FDCWD            = C.AT_FDCWD
	AT_REMOVEDIR        = C.AT_REMOVEDIR
	AT_SYMLINK_NOFOLLOW = C.AT_SYMLINK_NOFOLLOW
)



type Termios C.struct_termios

type Termio C.struct_termio

type Winsize C.struct_winsize



type PollFd struct {
	Fd      int32
	Events  uint16
	Revents uint16
}

const (
	POLLERR    = C.POLLERR
	POLLHUP    = C.POLLHUP
	POLLIN     = C.POLLIN
	POLLNVAL   = C.POLLNVAL
	POLLOUT    = C.POLLOUT
	POLLPRI    = C.POLLPRI
	POLLRDBAND = C.POLLRDBAND
	POLLRDNORM = C.POLLRDNORM
	POLLWRBAND = C.POLLWRBAND
	POLLWRNORM = C.POLLWRNORM
)



type Flock_t C.struct_flock64



type Fsid_t C.struct_fsid_t
type Fsid64_t C.struct_fsid64_t

type Statfs_t C.struct_statfs

const RNDGETENTCNT = 0x80045200
