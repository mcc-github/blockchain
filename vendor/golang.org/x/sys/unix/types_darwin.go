










package unix


import "C"



const (
	sizeofPtr      = C.sizeofPtr
	sizeofShort    = C.sizeof_short
	sizeofInt      = C.sizeof_int
	sizeofLong     = C.sizeof_long
	sizeofLongLong = C.sizeof_longlong
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



type Rusage C.struct_rusage

type Rlimit C.struct_rlimit

type _Gid_t C.gid_t



type Stat_t C.struct_stat64

type Statfs_t C.struct_statfs64

type Flock_t C.struct_flock

type Fstore_t C.struct_fstore

type Radvisory_t C.struct_radvisory

type Fbootstraptransfer_t C.struct_fbootstraptransfer

type Log2phys_t C.struct_log2phys

type Fsid C.struct_fsid

type Dirent C.struct_dirent



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

type Inet4Pktinfo C.struct_in_pktinfo

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
	SizeofInet4Pktinfo     = C.sizeof_struct_in_pktinfo
	SizeofInet6Pktinfo     = C.sizeof_struct_in6_pktinfo
	SizeofIPv6MTUInfo      = C.sizeof_struct_ip6_mtuinfo
	SizeofICMPv6Filter     = C.sizeof_struct_icmp6_filter
)



const (
	PTRACE_TRACEME = C.PT_TRACE_ME
	PTRACE_CONT    = C.PT_CONTINUE
	PTRACE_KILL    = C.PT_KILL
)



type Kevent_t C.struct_kevent



type FdSet C.fd_set



const (
	SizeofIfMsghdr    = C.sizeof_struct_if_msghdr
	SizeofIfData      = C.sizeof_struct_if_data
	SizeofIfaMsghdr   = C.sizeof_struct_ifa_msghdr
	SizeofIfmaMsghdr  = C.sizeof_struct_ifma_msghdr
	SizeofIfmaMsghdr2 = C.sizeof_struct_ifma_msghdr2
	SizeofRtMsghdr    = C.sizeof_struct_rt_msghdr
	SizeofRtMetrics   = C.sizeof_struct_rt_metrics
)

type IfMsghdr C.struct_if_msghdr

type IfData C.struct_if_data

type IfaMsghdr C.struct_ifa_msghdr

type IfmaMsghdr C.struct_ifma_msghdr

type IfmaMsghdr2 C.struct_ifma_msghdr2

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

type BpfHdr C.struct_bpf_hdr



type Termios C.struct_termios

type Winsize C.struct_winsize



const (
	AT_FDCWD            = C.AT_FDCWD
	AT_REMOVEDIR        = C.AT_REMOVEDIR
	AT_SYMLINK_FOLLOW   = C.AT_SYMLINK_FOLLOW
	AT_SYMLINK_NOFOLLOW = C.AT_SYMLINK_NOFOLLOW
)



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



type Utsname C.struct_utsname