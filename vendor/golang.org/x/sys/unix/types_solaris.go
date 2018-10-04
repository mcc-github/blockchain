










package unix


import "C"



const (
	sizeofPtr      = C.sizeofPtr
	sizeofShort    = C.sizeof_short
	sizeofInt      = C.sizeof_int
	sizeofLong     = C.sizeof_long
	sizeofLongLong = C.sizeof_longlong
	PathMax        = C.PATH_MAX
	MaxHostNameLen = C.MAXHOSTNAMELEN
)



type (
	_C_short     C.short
	_C_int       C.int
	_C_long      C.long
	_C_long_long C.longlong
)



type Timespec C.struct_timespec

type Timeval C.struct_timeval

type Timeval32 C.struct_timeval32

type Tms C.struct_tms

type Utimbuf C.struct_utimbuf



type Rusage C.struct_rusage

type Rlimit C.struct_rlimit

type _Gid_t C.gid_t



type Stat_t C.struct_stat

type Flock_t C.struct_flock

type Dirent C.struct_dirent



type _Fsblkcnt_t C.fsblkcnt_t

type Statvfs_t C.struct_statvfs



type RawSockaddrInet4 C.struct_sockaddr_in

type RawSockaddrInet6 C.struct_sockaddr_in6

type RawSockaddrUnix C.struct_sockaddr_un

type RawSockaddrDatalink C.struct_sockaddr_dl

type RawSockaddr C.struct_sockaddr

type RawSockaddrAny C.struct_sockaddr_any

type _Socklen C.socklen_t

type Linger C.struct_linger

type Iovec C.struct_iovec

type IPMreq C.struct_ip_mreq

type IPv6Mreq C.struct_ipv6_mreq

type Msghdr C.struct_msghdr

type Cmsghdr C.struct_cmsghdr

type Inet6Pktinfo C.struct_in6_pktinfo

type IPv6MTUInfo C.struct_ip6_mtuinfo

type ICMPv6Filter C.struct_icmp6_filter

const (
	SizeofSockaddrInet4    = C.sizeof_struct_sockaddr_in
	SizeofSockaddrInet6    = C.sizeof_struct_sockaddr_in6
	SizeofSockaddrAny      = C.sizeof_struct_sockaddr_any
	SizeofSockaddrUnix     = C.sizeof_struct_sockaddr_un
	SizeofSockaddrDatalink = C.sizeof_struct_sockaddr_dl
	SizeofLinger           = C.sizeof_struct_linger
	SizeofIPMreq           = C.sizeof_struct_ip_mreq
	SizeofIPv6Mreq         = C.sizeof_struct_ipv6_mreq
	SizeofMsghdr           = C.sizeof_struct_msghdr
	SizeofCmsghdr          = C.sizeof_struct_cmsghdr
	SizeofInet6Pktinfo     = C.sizeof_struct_in6_pktinfo
	SizeofIPv6MTUInfo      = C.sizeof_struct_ip6_mtuinfo
	SizeofICMPv6Filter     = C.sizeof_struct_icmp6_filter
)



type FdSet C.fd_set



type Utsname C.struct_utsname

type Ustat_t C.struct_ustat

const (
	AT_FDCWD            = C.AT_FDCWD
	AT_SYMLINK_NOFOLLOW = C.AT_SYMLINK_NOFOLLOW
	AT_SYMLINK_FOLLOW   = C.AT_SYMLINK_FOLLOW
	AT_REMOVEDIR        = C.AT_REMOVEDIR
	AT_EACCESS          = C.AT_EACCESS
)



const (
	SizeofIfMsghdr  = C.sizeof_struct_if_msghdr
	SizeofIfData    = C.sizeof_struct_if_data
	SizeofIfaMsghdr = C.sizeof_struct_ifa_msghdr
	SizeofRtMsghdr  = C.sizeof_struct_rt_msghdr
	SizeofRtMetrics = C.sizeof_struct_rt_metrics
)

type IfMsghdr C.struct_if_msghdr

type IfData C.struct_if_data

type IfaMsghdr C.struct_ifa_msghdr

type RtMsghdr C.struct_rt_msghdr

type RtMetrics C.struct_rt_metrics



const (
	SizeofBpfVersion = C.sizeof_struct_bpf_version
	SizeofBpfStat    = C.sizeof_struct_bpf_stat
	SizeofBpfProgram = C.sizeof_struct_bpf_program
	SizeofBpfInsn    = C.sizeof_struct_bpf_insn
	SizeofBpfHdr     = C.sizeof_struct_bpf_hdr
)

type BpfVersion C.struct_bpf_version

type BpfStat C.struct_bpf_stat

type BpfProgram C.struct_bpf_program

type BpfInsn C.struct_bpf_insn

type BpfTimeval C.struct_bpf_timeval

type BpfHdr C.struct_bpf_hdr



type Termios C.struct_termios

type Termio C.struct_termio

type Winsize C.struct_winsize



type PollFd C.struct_pollfd

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
