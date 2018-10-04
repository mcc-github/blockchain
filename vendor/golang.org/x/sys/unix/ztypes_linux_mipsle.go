




package unix

const (
	sizeofPtr      = 0x4
	sizeofShort    = 0x2
	sizeofInt      = 0x4
	sizeofLong     = 0x4
	sizeofLongLong = 0x8
	PathMax        = 0x1000
)

type (
	_C_short     int16
	_C_int       int32
	_C_long      int32
	_C_long_long int64
)

type Timespec struct {
	Sec  int32
	Nsec int32
}

type Timeval struct {
	Sec  int32
	Usec int32
}

type Timex struct {
	Modes     uint32
	Offset    int32
	Freq      int32
	Maxerror  int32
	Esterror  int32
	Status    int32
	Constant  int32
	Precision int32
	Tolerance int32
	Time      Timeval
	Tick      int32
	Ppsfreq   int32
	Jitter    int32
	Shift     int32
	Stabil    int32
	Jitcnt    int32
	Calcnt    int32
	Errcnt    int32
	Stbcnt    int32
	Tai       int32
	_         [44]byte
}

type Time_t int32

type Tms struct {
	Utime  int32
	Stime  int32
	Cutime int32
	Cstime int32
}

type Utimbuf struct {
	Actime  int32
	Modtime int32
}

type Rusage struct {
	Utime    Timeval
	Stime    Timeval
	Maxrss   int32
	Ixrss    int32
	Idrss    int32
	Isrss    int32
	Minflt   int32
	Majflt   int32
	Nswap    int32
	Inblock  int32
	Oublock  int32
	Msgsnd   int32
	Msgrcv   int32
	Nsignals int32
	Nvcsw    int32
	Nivcsw   int32
}

type Rlimit struct {
	Cur uint64
	Max uint64
}

type _Gid_t uint32

type Stat_t struct {
	Dev     uint32
	Pad1    [3]int32
	Ino     uint64
	Mode    uint32
	Nlink   uint32
	Uid     uint32
	Gid     uint32
	Rdev    uint32
	Pad2    [3]int32
	Size    int64
	Atim    Timespec
	Mtim    Timespec
	Ctim    Timespec
	Blksize int32
	Pad4    int32
	Blocks  int64
	Pad5    [14]int32
}

type StatxTimestamp struct {
	Sec  int64
	Nsec uint32
	_    int32
}

type Statx_t struct {
	Mask            uint32
	Blksize         uint32
	Attributes      uint64
	Nlink           uint32
	Uid             uint32
	Gid             uint32
	Mode            uint16
	_               [1]uint16
	Ino             uint64
	Size            uint64
	Blocks          uint64
	Attributes_mask uint64
	Atime           StatxTimestamp
	Btime           StatxTimestamp
	Ctime           StatxTimestamp
	Mtime           StatxTimestamp
	Rdev_major      uint32
	Rdev_minor      uint32
	Dev_major       uint32
	Dev_minor       uint32
	_               [14]uint64
}

type Dirent struct {
	Ino    uint64
	Off    int64
	Reclen uint16
	Type   uint8
	Name   [256]int8
	_      [5]byte
}

type Fsid struct {
	Val [2]int32
}

type Flock_t struct {
	Type   int16
	Whence int16
	_      [4]byte
	Start  int64
	Len    int64
	Pid    int32
	_      [4]byte
}

type FscryptPolicy struct {
	Version                   uint8
	Contents_encryption_mode  uint8
	Filenames_encryption_mode uint8
	Flags                     uint8
	Master_key_descriptor     [8]uint8
}

type FscryptKey struct {
	Mode uint32
	Raw  [64]uint8
	Size uint32
}

type KeyctlDHParams struct {
	Private int32
	Prime   int32
	Base    int32
}

const (
	FADV_NORMAL     = 0x0
	FADV_RANDOM     = 0x1
	FADV_SEQUENTIAL = 0x2
	FADV_WILLNEED   = 0x3
	FADV_DONTNEED   = 0x4
	FADV_NOREUSE    = 0x5
)

type RawSockaddrInet4 struct {
	Family uint16
	Port   uint16
	Addr   [4]byte 
	Zero   [8]uint8
}

type RawSockaddrInet6 struct {
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte 
	Scope_id uint32
}

type RawSockaddrUnix struct {
	Family uint16
	Path   [108]int8
}

type RawSockaddrLinklayer struct {
	Family   uint16
	Protocol uint16
	Ifindex  int32
	Hatype   uint16
	Pkttype  uint8
	Halen    uint8
	Addr     [8]uint8
}

type RawSockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
}

type RawSockaddrHCI struct {
	Family  uint16
	Dev     uint16
	Channel uint16
}

type RawSockaddrL2 struct {
	Family      uint16
	Psm         uint16
	Bdaddr      [6]uint8
	Cid         uint16
	Bdaddr_type uint8
	_           [1]byte
}

type RawSockaddrRFCOMM struct {
	Family  uint16
	Bdaddr  [6]uint8
	Channel uint8
	_       [1]byte
}

type RawSockaddrCAN struct {
	Family  uint16
	_       [2]byte
	Ifindex int32
	Addr    [8]byte
}

type RawSockaddrALG struct {
	Family uint16
	Type   [14]uint8
	Feat   uint32
	Mask   uint32
	Name   [64]uint8
}

type RawSockaddrVM struct {
	Family    uint16
	Reserved1 uint16
	Port      uint32
	Cid       uint32
	Zero      [4]uint8
}

type RawSockaddrXDP struct {
	Family         uint16
	Flags          uint16
	Ifindex        uint32
	Queue_id       uint32
	Shared_umem_fd uint32
}

type RawSockaddr struct {
	Family uint16
	Data   [14]int8
}

type RawSockaddrAny struct {
	Addr RawSockaddr
	Pad  [96]int8
}

type _Socklen uint32

type Linger struct {
	Onoff  int32
	Linger int32
}

type Iovec struct {
	Base *byte
	Len  uint32
}

type IPMreq struct {
	Multiaddr [4]byte 
	Interface [4]byte 
}

type IPMreqn struct {
	Multiaddr [4]byte 
	Address   [4]byte 
	Ifindex   int32
}

type IPv6Mreq struct {
	Multiaddr [16]byte 
	Interface uint32
}

type PacketMreq struct {
	Ifindex int32
	Type    uint16
	Alen    uint16
	Address [8]uint8
}

type Msghdr struct {
	Name       *byte
	Namelen    uint32
	Iov        *Iovec
	Iovlen     uint32
	Control    *byte
	Controllen uint32
	Flags      int32
}

type Cmsghdr struct {
	Len   uint32
	Level int32
	Type  int32
}

type Inet4Pktinfo struct {
	Ifindex  int32
	Spec_dst [4]byte 
	Addr     [4]byte 
}

type Inet6Pktinfo struct {
	Addr    [16]byte 
	Ifindex uint32
}

type IPv6MTUInfo struct {
	Addr RawSockaddrInet6
	Mtu  uint32
}

type ICMPv6Filter struct {
	Data [8]uint32
}

type Ucred struct {
	Pid int32
	Uid uint32
	Gid uint32
}

type TCPInfo struct {
	State          uint8
	Ca_state       uint8
	Retransmits    uint8
	Probes         uint8
	Backoff        uint8
	Options        uint8
	_              [2]byte
	Rto            uint32
	Ato            uint32
	Snd_mss        uint32
	Rcv_mss        uint32
	Unacked        uint32
	Sacked         uint32
	Lost           uint32
	Retrans        uint32
	Fackets        uint32
	Last_data_sent uint32
	Last_ack_sent  uint32
	Last_data_recv uint32
	Last_ack_recv  uint32
	Pmtu           uint32
	Rcv_ssthresh   uint32
	Rtt            uint32
	Rttvar         uint32
	Snd_ssthresh   uint32
	Snd_cwnd       uint32
	Advmss         uint32
	Reordering     uint32
	Rcv_rtt        uint32
	Rcv_space      uint32
	Total_retrans  uint32
}

const (
	SizeofSockaddrInet4     = 0x10
	SizeofSockaddrInet6     = 0x1c
	SizeofSockaddrAny       = 0x70
	SizeofSockaddrUnix      = 0x6e
	SizeofSockaddrLinklayer = 0x14
	SizeofSockaddrNetlink   = 0xc
	SizeofSockaddrHCI       = 0x6
	SizeofSockaddrL2        = 0xe
	SizeofSockaddrRFCOMM    = 0xa
	SizeofSockaddrCAN       = 0x10
	SizeofSockaddrALG       = 0x58
	SizeofSockaddrVM        = 0x10
	SizeofSockaddrXDP       = 0x10
	SizeofLinger            = 0x8
	SizeofIovec             = 0x8
	SizeofIPMreq            = 0x8
	SizeofIPMreqn           = 0xc
	SizeofIPv6Mreq          = 0x14
	SizeofPacketMreq        = 0x10
	SizeofMsghdr            = 0x1c
	SizeofCmsghdr           = 0xc
	SizeofInet4Pktinfo      = 0xc
	SizeofInet6Pktinfo      = 0x14
	SizeofIPv6MTUInfo       = 0x20
	SizeofICMPv6Filter      = 0x20
	SizeofUcred             = 0xc
	SizeofTCPInfo           = 0x68
)

const (
	IFA_UNSPEC           = 0x0
	IFA_ADDRESS          = 0x1
	IFA_LOCAL            = 0x2
	IFA_LABEL            = 0x3
	IFA_BROADCAST        = 0x4
	IFA_ANYCAST          = 0x5
	IFA_CACHEINFO        = 0x6
	IFA_MULTICAST        = 0x7
	IFLA_UNSPEC          = 0x0
	IFLA_ADDRESS         = 0x1
	IFLA_BROADCAST       = 0x2
	IFLA_IFNAME          = 0x3
	IFLA_INFO_KIND       = 0x1
	IFLA_MTU             = 0x4
	IFLA_LINK            = 0x5
	IFLA_QDISC           = 0x6
	IFLA_STATS           = 0x7
	IFLA_COST            = 0x8
	IFLA_PRIORITY        = 0x9
	IFLA_MASTER          = 0xa
	IFLA_WIRELESS        = 0xb
	IFLA_PROTINFO        = 0xc
	IFLA_TXQLEN          = 0xd
	IFLA_MAP             = 0xe
	IFLA_WEIGHT          = 0xf
	IFLA_OPERSTATE       = 0x10
	IFLA_LINKMODE        = 0x11
	IFLA_LINKINFO        = 0x12
	IFLA_NET_NS_PID      = 0x13
	IFLA_IFALIAS         = 0x14
	IFLA_NUM_VF          = 0x15
	IFLA_VFINFO_LIST     = 0x16
	IFLA_STATS64         = 0x17
	IFLA_VF_PORTS        = 0x18
	IFLA_PORT_SELF       = 0x19
	IFLA_AF_SPEC         = 0x1a
	IFLA_GROUP           = 0x1b
	IFLA_NET_NS_FD       = 0x1c
	IFLA_EXT_MASK        = 0x1d
	IFLA_PROMISCUITY     = 0x1e
	IFLA_NUM_TX_QUEUES   = 0x1f
	IFLA_NUM_RX_QUEUES   = 0x20
	IFLA_CARRIER         = 0x21
	IFLA_PHYS_PORT_ID    = 0x22
	IFLA_CARRIER_CHANGES = 0x23
	IFLA_PHYS_SWITCH_ID  = 0x24
	IFLA_LINK_NETNSID    = 0x25
	IFLA_PHYS_PORT_NAME  = 0x26
	IFLA_PROTO_DOWN      = 0x27
	IFLA_GSO_MAX_SEGS    = 0x28
	IFLA_GSO_MAX_SIZE    = 0x29
	IFLA_PAD             = 0x2a
	IFLA_XDP             = 0x2b
	IFLA_EVENT           = 0x2c
	IFLA_NEW_NETNSID     = 0x2d
	IFLA_IF_NETNSID      = 0x2e
	IFLA_MAX             = 0x31
	RT_SCOPE_UNIVERSE    = 0x0
	RT_SCOPE_SITE        = 0xc8
	RT_SCOPE_LINK        = 0xfd
	RT_SCOPE_HOST        = 0xfe
	RT_SCOPE_NOWHERE     = 0xff
	RT_TABLE_UNSPEC      = 0x0
	RT_TABLE_COMPAT      = 0xfc
	RT_TABLE_DEFAULT     = 0xfd
	RT_TABLE_MAIN        = 0xfe
	RT_TABLE_LOCAL       = 0xff
	RT_TABLE_MAX         = 0xffffffff
	RTA_UNSPEC           = 0x0
	RTA_DST              = 0x1
	RTA_SRC              = 0x2
	RTA_IIF              = 0x3
	RTA_OIF              = 0x4
	RTA_GATEWAY          = 0x5
	RTA_PRIORITY         = 0x6
	RTA_PREFSRC          = 0x7
	RTA_METRICS          = 0x8
	RTA_MULTIPATH        = 0x9
	RTA_FLOW             = 0xb
	RTA_CACHEINFO        = 0xc
	RTA_TABLE            = 0xf
	RTA_MARK             = 0x10
	RTA_MFC_STATS        = 0x11
	RTA_VIA              = 0x12
	RTA_NEWDST           = 0x13
	RTA_PREF             = 0x14
	RTA_ENCAP_TYPE       = 0x15
	RTA_ENCAP            = 0x16
	RTA_EXPIRES          = 0x17
	RTA_PAD              = 0x18
	RTA_UID              = 0x19
	RTA_TTL_PROPAGATE    = 0x1a
	RTA_IP_PROTO         = 0x1b
	RTA_SPORT            = 0x1c
	RTA_DPORT            = 0x1d
	RTN_UNSPEC           = 0x0
	RTN_UNICAST          = 0x1
	RTN_LOCAL            = 0x2
	RTN_BROADCAST        = 0x3
	RTN_ANYCAST          = 0x4
	RTN_MULTICAST        = 0x5
	RTN_BLACKHOLE        = 0x6
	RTN_UNREACHABLE      = 0x7
	RTN_PROHIBIT         = 0x8
	RTN_THROW            = 0x9
	RTN_NAT              = 0xa
	RTN_XRESOLVE         = 0xb
	RTNLGRP_NONE         = 0x0
	RTNLGRP_LINK         = 0x1
	RTNLGRP_NOTIFY       = 0x2
	RTNLGRP_NEIGH        = 0x3
	RTNLGRP_TC           = 0x4
	RTNLGRP_IPV4_IFADDR  = 0x5
	RTNLGRP_IPV4_MROUTE  = 0x6
	RTNLGRP_IPV4_ROUTE   = 0x7
	RTNLGRP_IPV4_RULE    = 0x8
	RTNLGRP_IPV6_IFADDR  = 0x9
	RTNLGRP_IPV6_MROUTE  = 0xa
	RTNLGRP_IPV6_ROUTE   = 0xb
	RTNLGRP_IPV6_IFINFO  = 0xc
	RTNLGRP_IPV6_PREFIX  = 0x12
	RTNLGRP_IPV6_RULE    = 0x13
	RTNLGRP_ND_USEROPT   = 0x14
	SizeofNlMsghdr       = 0x10
	SizeofNlMsgerr       = 0x14
	SizeofRtGenmsg       = 0x1
	SizeofNlAttr         = 0x4
	SizeofRtAttr         = 0x4
	SizeofIfInfomsg      = 0x10
	SizeofIfAddrmsg      = 0x8
	SizeofRtMsg          = 0xc
	SizeofRtNexthop      = 0x8
)

type NlMsghdr struct {
	Len   uint32
	Type  uint16
	Flags uint16
	Seq   uint32
	Pid   uint32
}

type NlMsgerr struct {
	Error int32
	Msg   NlMsghdr
}

type RtGenmsg struct {
	Family uint8
}

type NlAttr struct {
	Len  uint16
	Type uint16
}

type RtAttr struct {
	Len  uint16
	Type uint16
}

type IfInfomsg struct {
	Family uint8
	_      uint8
	Type   uint16
	Index  int32
	Flags  uint32
	Change uint32
}

type IfAddrmsg struct {
	Family    uint8
	Prefixlen uint8
	Flags     uint8
	Scope     uint8
	Index     uint32
}

type RtMsg struct {
	Family   uint8
	Dst_len  uint8
	Src_len  uint8
	Tos      uint8
	Table    uint8
	Protocol uint8
	Scope    uint8
	Type     uint8
	Flags    uint32
}

type RtNexthop struct {
	Len     uint16
	Flags   uint8
	Hops    uint8
	Ifindex int32
}

const (
	SizeofSockFilter = 0x8
	SizeofSockFprog  = 0x8
)

type SockFilter struct {
	Code uint16
	Jt   uint8
	Jf   uint8
	K    uint32
}

type SockFprog struct {
	Len    uint16
	_      [2]byte
	Filter *SockFilter
}

type InotifyEvent struct {
	Wd     int32
	Mask   uint32
	Cookie uint32
	Len    uint32
}

const SizeofInotifyEvent = 0x10

type PtraceRegs struct {
	Regs     [32]uint64
	Lo       uint64
	Hi       uint64
	Epc      uint64
	Badvaddr uint64
	Status   uint64
	Cause    uint64
}

type FdSet struct {
	Bits [32]int32
}

type Sysinfo_t struct {
	Uptime    int32
	Loads     [3]uint32
	Totalram  uint32
	Freeram   uint32
	Sharedram uint32
	Bufferram uint32
	Totalswap uint32
	Freeswap  uint32
	Procs     uint16
	Pad       uint16
	Totalhigh uint32
	Freehigh  uint32
	Unit      uint32
	_         [8]int8
}

type Utsname struct {
	Sysname    [65]byte
	Nodename   [65]byte
	Release    [65]byte
	Version    [65]byte
	Machine    [65]byte
	Domainname [65]byte
}

type Ustat_t struct {
	Tfree  int32
	Tinode uint32
	Fname  [6]int8
	Fpack  [6]int8
}

type EpollEvent struct {
	Events uint32
	PadFd  int32
	Fd     int32
	Pad    int32
}

const (
	AT_EMPTY_PATH   = 0x1000
	AT_FDCWD        = -0x64
	AT_NO_AUTOMOUNT = 0x800
	AT_REMOVEDIR    = 0x200

	AT_STATX_SYNC_AS_STAT = 0x0
	AT_STATX_FORCE_SYNC   = 0x2000
	AT_STATX_DONT_SYNC    = 0x4000

	AT_SYMLINK_FOLLOW   = 0x400
	AT_SYMLINK_NOFOLLOW = 0x100

	AT_EACCESS = 0x200
)

type PollFd struct {
	Fd      int32
	Events  int16
	Revents int16
}

const (
	POLLIN    = 0x1
	POLLPRI   = 0x2
	POLLOUT   = 0x4
	POLLRDHUP = 0x2000
	POLLERR   = 0x8
	POLLHUP   = 0x10
	POLLNVAL  = 0x20
)

type Sigset_t struct {
	Val [32]uint32
}

const RNDGETENTCNT = 0x40045200

const PERF_IOC_FLAG_GROUP = 0x1

type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [23]uint8
	Ispeed uint32
	Ospeed uint32
}

type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

type Taskstats struct {
	Version                   uint16
	_                         [2]byte
	Ac_exitcode               uint32
	Ac_flag                   uint8
	Ac_nice                   uint8
	_                         [6]byte
	Cpu_count                 uint64
	Cpu_delay_total           uint64
	Blkio_count               uint64
	Blkio_delay_total         uint64
	Swapin_count              uint64
	Swapin_delay_total        uint64
	Cpu_run_real_total        uint64
	Cpu_run_virtual_total     uint64
	Ac_comm                   [32]int8
	Ac_sched                  uint8
	Ac_pad                    [3]uint8
	_                         [4]byte
	Ac_uid                    uint32
	Ac_gid                    uint32
	Ac_pid                    uint32
	Ac_ppid                   uint32
	Ac_btime                  uint32
	_                         [4]byte
	Ac_etime                  uint64
	Ac_utime                  uint64
	Ac_stime                  uint64
	Ac_minflt                 uint64
	Ac_majflt                 uint64
	Coremem                   uint64
	Virtmem                   uint64
	Hiwater_rss               uint64
	Hiwater_vm                uint64
	Read_char                 uint64
	Write_char                uint64
	Read_syscalls             uint64
	Write_syscalls            uint64
	Read_bytes                uint64
	Write_bytes               uint64
	Cancelled_write_bytes     uint64
	Nvcsw                     uint64
	Nivcsw                    uint64
	Ac_utimescaled            uint64
	Ac_stimescaled            uint64
	Cpu_scaled_run_real_total uint64
	Freepages_count           uint64
	Freepages_delay_total     uint64
}

const (
	TASKSTATS_CMD_UNSPEC                  = 0x0
	TASKSTATS_CMD_GET                     = 0x1
	TASKSTATS_CMD_NEW                     = 0x2
	TASKSTATS_TYPE_UNSPEC                 = 0x0
	TASKSTATS_TYPE_PID                    = 0x1
	TASKSTATS_TYPE_TGID                   = 0x2
	TASKSTATS_TYPE_STATS                  = 0x3
	TASKSTATS_TYPE_AGGR_PID               = 0x4
	TASKSTATS_TYPE_AGGR_TGID              = 0x5
	TASKSTATS_TYPE_NULL                   = 0x6
	TASKSTATS_CMD_ATTR_UNSPEC             = 0x0
	TASKSTATS_CMD_ATTR_PID                = 0x1
	TASKSTATS_CMD_ATTR_TGID               = 0x2
	TASKSTATS_CMD_ATTR_REGISTER_CPUMASK   = 0x3
	TASKSTATS_CMD_ATTR_DEREGISTER_CPUMASK = 0x4
)

type CGroupStats struct {
	Sleeping        uint64
	Running         uint64
	Stopped         uint64
	Uninterruptible uint64
	Io_wait         uint64
}

const (
	CGROUPSTATS_CMD_UNSPEC        = 0x3
	CGROUPSTATS_CMD_GET           = 0x4
	CGROUPSTATS_CMD_NEW           = 0x5
	CGROUPSTATS_TYPE_UNSPEC       = 0x0
	CGROUPSTATS_TYPE_CGROUP_STATS = 0x1
	CGROUPSTATS_CMD_ATTR_UNSPEC   = 0x0
	CGROUPSTATS_CMD_ATTR_FD       = 0x1
)

type Genlmsghdr struct {
	Cmd      uint8
	Version  uint8
	Reserved uint16
}

const (
	CTRL_CMD_UNSPEC            = 0x0
	CTRL_CMD_NEWFAMILY         = 0x1
	CTRL_CMD_DELFAMILY         = 0x2
	CTRL_CMD_GETFAMILY         = 0x3
	CTRL_CMD_NEWOPS            = 0x4
	CTRL_CMD_DELOPS            = 0x5
	CTRL_CMD_GETOPS            = 0x6
	CTRL_CMD_NEWMCAST_GRP      = 0x7
	CTRL_CMD_DELMCAST_GRP      = 0x8
	CTRL_CMD_GETMCAST_GRP      = 0x9
	CTRL_ATTR_UNSPEC           = 0x0
	CTRL_ATTR_FAMILY_ID        = 0x1
	CTRL_ATTR_FAMILY_NAME      = 0x2
	CTRL_ATTR_VERSION          = 0x3
	CTRL_ATTR_HDRSIZE          = 0x4
	CTRL_ATTR_MAXATTR          = 0x5
	CTRL_ATTR_OPS              = 0x6
	CTRL_ATTR_MCAST_GROUPS     = 0x7
	CTRL_ATTR_OP_UNSPEC        = 0x0
	CTRL_ATTR_OP_ID            = 0x1
	CTRL_ATTR_OP_FLAGS         = 0x2
	CTRL_ATTR_MCAST_GRP_UNSPEC = 0x0
	CTRL_ATTR_MCAST_GRP_NAME   = 0x1
	CTRL_ATTR_MCAST_GRP_ID     = 0x2
)

type cpuMask uint32

const (
	_CPU_SETSIZE = 0x400
	_NCPUBITS    = 0x20
)

const (
	BDADDR_BREDR     = 0x0
	BDADDR_LE_PUBLIC = 0x1
	BDADDR_LE_RANDOM = 0x2
)

type PerfEventAttr struct {
	Type               uint32
	Size               uint32
	Config             uint64
	Sample             uint64
	Sample_type        uint64
	Read_format        uint64
	Bits               uint64
	Wakeup             uint32
	Bp_type            uint32
	Ext1               uint64
	Ext2               uint64
	Branch_sample_type uint64
	Sample_regs_user   uint64
	Sample_stack_user  uint32
	Clockid            int32
	Sample_regs_intr   uint64
	Aux_watermark      uint32
	_                  uint32
}

type PerfEventMmapPage struct {
	Version        uint32
	Compat_version uint32
	Lock           uint32
	Index          uint32
	Offset         int64
	Time_enabled   uint64
	Time_running   uint64
	Capabilities   uint64
	Pmc_width      uint16
	Time_shift     uint16
	Time_mult      uint32
	Time_offset    uint64
	Time_zero      uint64
	Size           uint32
	_              [948]uint8
	Data_head      uint64
	Data_tail      uint64
	Data_offset    uint64
	Data_size      uint64
	Aux_head       uint64
	Aux_tail       uint64
	Aux_offset     uint64
	Aux_size       uint64
}

const (
	PerfBitDisabled               uint64 = CBitFieldMaskBit0
	PerfBitInherit                       = CBitFieldMaskBit1
	PerfBitPinned                        = CBitFieldMaskBit2
	PerfBitExclusive                     = CBitFieldMaskBit3
	PerfBitExcludeUser                   = CBitFieldMaskBit4
	PerfBitExcludeKernel                 = CBitFieldMaskBit5
	PerfBitExcludeHv                     = CBitFieldMaskBit6
	PerfBitExcludeIdle                   = CBitFieldMaskBit7
	PerfBitMmap                          = CBitFieldMaskBit8
	PerfBitComm                          = CBitFieldMaskBit9
	PerfBitFreq                          = CBitFieldMaskBit10
	PerfBitInheritStat                   = CBitFieldMaskBit11
	PerfBitEnableOnExec                  = CBitFieldMaskBit12
	PerfBitTask                          = CBitFieldMaskBit13
	PerfBitWatermark                     = CBitFieldMaskBit14
	PerfBitPreciseIPBit1                 = CBitFieldMaskBit15
	PerfBitPreciseIPBit2                 = CBitFieldMaskBit16
	PerfBitMmapData                      = CBitFieldMaskBit17
	PerfBitSampleIDAll                   = CBitFieldMaskBit18
	PerfBitExcludeHost                   = CBitFieldMaskBit19
	PerfBitExcludeGuest                  = CBitFieldMaskBit20
	PerfBitExcludeCallchainKernel        = CBitFieldMaskBit21
	PerfBitExcludeCallchainUser          = CBitFieldMaskBit22
	PerfBitMmap2                         = CBitFieldMaskBit23
	PerfBitCommExec                      = CBitFieldMaskBit24
	PerfBitUseClockID                    = CBitFieldMaskBit25
	PerfBitContextSwitch                 = CBitFieldMaskBit26
)

const (
	PERF_TYPE_HARDWARE   = 0x0
	PERF_TYPE_SOFTWARE   = 0x1
	PERF_TYPE_TRACEPOINT = 0x2
	PERF_TYPE_HW_CACHE   = 0x3
	PERF_TYPE_RAW        = 0x4
	PERF_TYPE_BREAKPOINT = 0x5

	PERF_COUNT_HW_CPU_CYCLES              = 0x0
	PERF_COUNT_HW_INSTRUCTIONS            = 0x1
	PERF_COUNT_HW_CACHE_REFERENCES        = 0x2
	PERF_COUNT_HW_CACHE_MISSES            = 0x3
	PERF_COUNT_HW_BRANCH_INSTRUCTIONS     = 0x4
	PERF_COUNT_HW_BRANCH_MISSES           = 0x5
	PERF_COUNT_HW_BUS_CYCLES              = 0x6
	PERF_COUNT_HW_STALLED_CYCLES_FRONTEND = 0x7
	PERF_COUNT_HW_STALLED_CYCLES_BACKEND  = 0x8
	PERF_COUNT_HW_REF_CPU_CYCLES          = 0x9

	PERF_COUNT_HW_CACHE_L1D  = 0x0
	PERF_COUNT_HW_CACHE_L1I  = 0x1
	PERF_COUNT_HW_CACHE_LL   = 0x2
	PERF_COUNT_HW_CACHE_DTLB = 0x3
	PERF_COUNT_HW_CACHE_ITLB = 0x4
	PERF_COUNT_HW_CACHE_BPU  = 0x5
	PERF_COUNT_HW_CACHE_NODE = 0x6

	PERF_COUNT_HW_CACHE_OP_READ     = 0x0
	PERF_COUNT_HW_CACHE_OP_WRITE    = 0x1
	PERF_COUNT_HW_CACHE_OP_PREFETCH = 0x2

	PERF_COUNT_HW_CACHE_RESULT_ACCESS = 0x0
	PERF_COUNT_HW_CACHE_RESULT_MISS   = 0x1

	PERF_COUNT_SW_CPU_CLOCK        = 0x0
	PERF_COUNT_SW_TASK_CLOCK       = 0x1
	PERF_COUNT_SW_PAGE_FAULTS      = 0x2
	PERF_COUNT_SW_CONTEXT_SWITCHES = 0x3
	PERF_COUNT_SW_CPU_MIGRATIONS   = 0x4
	PERF_COUNT_SW_PAGE_FAULTS_MIN  = 0x5
	PERF_COUNT_SW_PAGE_FAULTS_MAJ  = 0x6
	PERF_COUNT_SW_ALIGNMENT_FAULTS = 0x7
	PERF_COUNT_SW_EMULATION_FAULTS = 0x8
	PERF_COUNT_SW_DUMMY            = 0x9

	PERF_SAMPLE_IP           = 0x1
	PERF_SAMPLE_TID          = 0x2
	PERF_SAMPLE_TIME         = 0x4
	PERF_SAMPLE_ADDR         = 0x8
	PERF_SAMPLE_READ         = 0x10
	PERF_SAMPLE_CALLCHAIN    = 0x20
	PERF_SAMPLE_ID           = 0x40
	PERF_SAMPLE_CPU          = 0x80
	PERF_SAMPLE_PERIOD       = 0x100
	PERF_SAMPLE_STREAM_ID    = 0x200
	PERF_SAMPLE_RAW          = 0x400
	PERF_SAMPLE_BRANCH_STACK = 0x800

	PERF_SAMPLE_BRANCH_USER       = 0x1
	PERF_SAMPLE_BRANCH_KERNEL     = 0x2
	PERF_SAMPLE_BRANCH_HV         = 0x4
	PERF_SAMPLE_BRANCH_ANY        = 0x8
	PERF_SAMPLE_BRANCH_ANY_CALL   = 0x10
	PERF_SAMPLE_BRANCH_ANY_RETURN = 0x20
	PERF_SAMPLE_BRANCH_IND_CALL   = 0x40

	PERF_FORMAT_TOTAL_TIME_ENABLED = 0x1
	PERF_FORMAT_TOTAL_TIME_RUNNING = 0x2
	PERF_FORMAT_ID                 = 0x4
	PERF_FORMAT_GROUP              = 0x8

	PERF_RECORD_MMAP       = 0x1
	PERF_RECORD_LOST       = 0x2
	PERF_RECORD_COMM       = 0x3
	PERF_RECORD_EXIT       = 0x4
	PERF_RECORD_THROTTLE   = 0x5
	PERF_RECORD_UNTHROTTLE = 0x6
	PERF_RECORD_FORK       = 0x7
	PERF_RECORD_READ       = 0x8
	PERF_RECORD_SAMPLE     = 0x9

	PERF_CONTEXT_HV     = -0x20
	PERF_CONTEXT_KERNEL = -0x80
	PERF_CONTEXT_USER   = -0x200

	PERF_CONTEXT_GUEST        = -0x800
	PERF_CONTEXT_GUEST_KERNEL = -0x880
	PERF_CONTEXT_GUEST_USER   = -0xa00

	PERF_FLAG_FD_NO_GROUP = 0x1
	PERF_FLAG_FD_OUTPUT   = 0x2
	PERF_FLAG_PID_CGROUP  = 0x4
)

const (
	CBitFieldMaskBit0  = 0x1
	CBitFieldMaskBit1  = 0x2
	CBitFieldMaskBit2  = 0x4
	CBitFieldMaskBit3  = 0x8
	CBitFieldMaskBit4  = 0x10
	CBitFieldMaskBit5  = 0x20
	CBitFieldMaskBit6  = 0x40
	CBitFieldMaskBit7  = 0x80
	CBitFieldMaskBit8  = 0x100
	CBitFieldMaskBit9  = 0x200
	CBitFieldMaskBit10 = 0x400
	CBitFieldMaskBit11 = 0x800
	CBitFieldMaskBit12 = 0x1000
	CBitFieldMaskBit13 = 0x2000
	CBitFieldMaskBit14 = 0x4000
	CBitFieldMaskBit15 = 0x8000
	CBitFieldMaskBit16 = 0x10000
	CBitFieldMaskBit17 = 0x20000
	CBitFieldMaskBit18 = 0x40000
	CBitFieldMaskBit19 = 0x80000
	CBitFieldMaskBit20 = 0x100000
	CBitFieldMaskBit21 = 0x200000
	CBitFieldMaskBit22 = 0x400000
	CBitFieldMaskBit23 = 0x800000
	CBitFieldMaskBit24 = 0x1000000
	CBitFieldMaskBit25 = 0x2000000
	CBitFieldMaskBit26 = 0x4000000
	CBitFieldMaskBit27 = 0x8000000
	CBitFieldMaskBit28 = 0x10000000
	CBitFieldMaskBit29 = 0x20000000
	CBitFieldMaskBit30 = 0x40000000
	CBitFieldMaskBit31 = 0x80000000
	CBitFieldMaskBit32 = 0x100000000
	CBitFieldMaskBit33 = 0x200000000
	CBitFieldMaskBit34 = 0x400000000
	CBitFieldMaskBit35 = 0x800000000
	CBitFieldMaskBit36 = 0x1000000000
	CBitFieldMaskBit37 = 0x2000000000
	CBitFieldMaskBit38 = 0x4000000000
	CBitFieldMaskBit39 = 0x8000000000
	CBitFieldMaskBit40 = 0x10000000000
	CBitFieldMaskBit41 = 0x20000000000
	CBitFieldMaskBit42 = 0x40000000000
	CBitFieldMaskBit43 = 0x80000000000
	CBitFieldMaskBit44 = 0x100000000000
	CBitFieldMaskBit45 = 0x200000000000
	CBitFieldMaskBit46 = 0x400000000000
	CBitFieldMaskBit47 = 0x800000000000
	CBitFieldMaskBit48 = 0x1000000000000
	CBitFieldMaskBit49 = 0x2000000000000
	CBitFieldMaskBit50 = 0x4000000000000
	CBitFieldMaskBit51 = 0x8000000000000
	CBitFieldMaskBit52 = 0x10000000000000
	CBitFieldMaskBit53 = 0x20000000000000
	CBitFieldMaskBit54 = 0x40000000000000
	CBitFieldMaskBit55 = 0x80000000000000
	CBitFieldMaskBit56 = 0x100000000000000
	CBitFieldMaskBit57 = 0x200000000000000
	CBitFieldMaskBit58 = 0x400000000000000
	CBitFieldMaskBit59 = 0x800000000000000
	CBitFieldMaskBit60 = 0x1000000000000000
	CBitFieldMaskBit61 = 0x2000000000000000
	CBitFieldMaskBit62 = 0x4000000000000000
	CBitFieldMaskBit63 = 0x8000000000000000
)

type SockaddrStorage struct {
	Family uint16
	_      [122]int8
	_      uint32
}

type TCPMD5Sig struct {
	Addr      SockaddrStorage
	Flags     uint8
	Prefixlen uint8
	Keylen    uint16
	_         uint32
	Key       [80]uint8
}

type HDDriveCmdHdr struct {
	Command uint8
	Number  uint8
	Feature uint8
	Count   uint8
}

type HDGeometry struct {
	Heads     uint8
	Sectors   uint8
	Cylinders uint16
	Start     uint32
}

type HDDriveID struct {
	Config         uint16
	Cyls           uint16
	Reserved2      uint16
	Heads          uint16
	Track_bytes    uint16
	Sector_bytes   uint16
	Sectors        uint16
	Vendor0        uint16
	Vendor1        uint16
	Vendor2        uint16
	Serial_no      [20]uint8
	Buf_type       uint16
	Buf_size       uint16
	Ecc_bytes      uint16
	Fw_rev         [8]uint8
	Model          [40]uint8
	Max_multsect   uint8
	Vendor3        uint8
	Dword_io       uint16
	Vendor4        uint8
	Capability     uint8
	Reserved50     uint16
	Vendor5        uint8
	TPIO           uint8
	Vendor6        uint8
	TDMA           uint8
	Field_valid    uint16
	Cur_cyls       uint16
	Cur_heads      uint16
	Cur_sectors    uint16
	Cur_capacity0  uint16
	Cur_capacity1  uint16
	Multsect       uint8
	Multsect_valid uint8
	Lba_capacity   uint32
	Dma_1word      uint16
	Dma_mword      uint16
	Eide_pio_modes uint16
	Eide_dma_min   uint16
	Eide_dma_time  uint16
	Eide_pio       uint16
	Eide_pio_iordy uint16
	Words69_70     [2]uint16
	Words71_74     [4]uint16
	Queue_depth    uint16
	Words76_79     [4]uint16
	Major_rev_num  uint16
	Minor_rev_num  uint16
	Command_set_1  uint16
	Command_set_2  uint16
	Cfsse          uint16
	Cfs_enable_1   uint16
	Cfs_enable_2   uint16
	Csf_default    uint16
	Dma_ultra      uint16
	Trseuc         uint16
	TrsEuc         uint16
	CurAPMvalues   uint16
	Mprc           uint16
	Hw_config      uint16
	Acoustic       uint16
	Msrqs          uint16
	Sxfert         uint16
	Sal            uint16
	Spg            uint32
	Lba_capacity_2 uint64
	Words104_125   [22]uint16
	Last_lun       uint16
	Word127        uint16
	Dlf            uint16
	Csfo           uint16
	Words130_155   [26]uint16
	Word156        uint16
	Words157_159   [3]uint16
	Cfa_power      uint16
	Words161_175   [15]uint16
	Words176_205   [30]uint16
	Words206_254   [49]uint16
	Integrity_word uint16
}

type Statfs_t struct {
	Type    int32
	Bsize   int32
	Frsize  int32
	_       [4]byte
	Blocks  uint64
	Bfree   uint64
	Files   uint64
	Ffree   uint64
	Bavail  uint64
	Fsid    Fsid
	Namelen int32
	Flags   int32
	Spare   [5]int32
	_       [4]byte
}

const (
	ST_MANDLOCK    = 0x40
	ST_NOATIME     = 0x400
	ST_NODEV       = 0x4
	ST_NODIRATIME  = 0x800
	ST_NOEXEC      = 0x8
	ST_NOSUID      = 0x2
	ST_RDONLY      = 0x1
	ST_RELATIME    = 0x1000
	ST_SYNCHRONOUS = 0x10
)

type TpacketHdr struct {
	Status  uint32
	Len     uint32
	Snaplen uint32
	Mac     uint16
	Net     uint16
	Sec     uint32
	Usec    uint32
}

type Tpacket2Hdr struct {
	Status    uint32
	Len       uint32
	Snaplen   uint32
	Mac       uint16
	Net       uint16
	Sec       uint32
	Nsec      uint32
	Vlan_tci  uint16
	Vlan_tpid uint16
	_         [4]uint8
}

type Tpacket3Hdr struct {
	Next_offset uint32
	Sec         uint32
	Nsec        uint32
	Snaplen     uint32
	Len         uint32
	Status      uint32
	Mac         uint16
	Net         uint16
	Hv1         TpacketHdrVariant1
	_           [8]uint8
}

type TpacketHdrVariant1 struct {
	Rxhash    uint32
	Vlan_tci  uint32
	Vlan_tpid uint16
	_         uint16
}

type TpacketBlockDesc struct {
	Version uint32
	To_priv uint32
	Hdr     [40]byte
}

type TpacketReq struct {
	Block_size uint32
	Block_nr   uint32
	Frame_size uint32
	Frame_nr   uint32
}

type TpacketReq3 struct {
	Block_size       uint32
	Block_nr         uint32
	Frame_size       uint32
	Frame_nr         uint32
	Retire_blk_tov   uint32
	Sizeof_priv      uint32
	Feature_req_word uint32
}

type TpacketStats struct {
	Packets uint32
	Drops   uint32
}

type TpacketStatsV3 struct {
	Packets      uint32
	Drops        uint32
	Freeze_q_cnt uint32
}

type TpacketAuxdata struct {
	Status    uint32
	Len       uint32
	Snaplen   uint32
	Mac       uint16
	Net       uint16
	Vlan_tci  uint16
	Vlan_tpid uint16
}

const (
	TPACKET_V1 = 0x0
	TPACKET_V2 = 0x1
	TPACKET_V3 = 0x2
)

const (
	SizeofTpacketHdr  = 0x18
	SizeofTpacket2Hdr = 0x20
	SizeofTpacket3Hdr = 0x30
)

const (
	NF_INET_PRE_ROUTING  = 0x0
	NF_INET_LOCAL_IN     = 0x1
	NF_INET_FORWARD      = 0x2
	NF_INET_LOCAL_OUT    = 0x3
	NF_INET_POST_ROUTING = 0x4
	NF_INET_NUMHOOKS     = 0x5
)

const (
	NF_NETDEV_INGRESS  = 0x0
	NF_NETDEV_NUMHOOKS = 0x1
)

const (
	NFPROTO_UNSPEC   = 0x0
	NFPROTO_INET     = 0x1
	NFPROTO_IPV4     = 0x2
	NFPROTO_ARP      = 0x3
	NFPROTO_NETDEV   = 0x5
	NFPROTO_BRIDGE   = 0x7
	NFPROTO_IPV6     = 0xa
	NFPROTO_DECNET   = 0xc
	NFPROTO_NUMPROTO = 0xd
)

type Nfgenmsg struct {
	Nfgen_family uint8
	Version      uint8
	Res_id       uint16
}

const (
	NFNL_BATCH_UNSPEC = 0x0
	NFNL_BATCH_GENID  = 0x1
)

const (
	NFT_REG_VERDICT                   = 0x0
	NFT_REG_1                         = 0x1
	NFT_REG_2                         = 0x2
	NFT_REG_3                         = 0x3
	NFT_REG_4                         = 0x4
	NFT_REG32_00                      = 0x8
	NFT_REG32_01                      = 0x9
	NFT_REG32_02                      = 0xa
	NFT_REG32_03                      = 0xb
	NFT_REG32_04                      = 0xc
	NFT_REG32_05                      = 0xd
	NFT_REG32_06                      = 0xe
	NFT_REG32_07                      = 0xf
	NFT_REG32_08                      = 0x10
	NFT_REG32_09                      = 0x11
	NFT_REG32_10                      = 0x12
	NFT_REG32_11                      = 0x13
	NFT_REG32_12                      = 0x14
	NFT_REG32_13                      = 0x15
	NFT_REG32_14                      = 0x16
	NFT_REG32_15                      = 0x17
	NFT_CONTINUE                      = -0x1
	NFT_BREAK                         = -0x2
	NFT_JUMP                          = -0x3
	NFT_GOTO                          = -0x4
	NFT_RETURN                        = -0x5
	NFT_MSG_NEWTABLE                  = 0x0
	NFT_MSG_GETTABLE                  = 0x1
	NFT_MSG_DELTABLE                  = 0x2
	NFT_MSG_NEWCHAIN                  = 0x3
	NFT_MSG_GETCHAIN                  = 0x4
	NFT_MSG_DELCHAIN                  = 0x5
	NFT_MSG_NEWRULE                   = 0x6
	NFT_MSG_GETRULE                   = 0x7
	NFT_MSG_DELRULE                   = 0x8
	NFT_MSG_NEWSET                    = 0x9
	NFT_MSG_GETSET                    = 0xa
	NFT_MSG_DELSET                    = 0xb
	NFT_MSG_NEWSETELEM                = 0xc
	NFT_MSG_GETSETELEM                = 0xd
	NFT_MSG_DELSETELEM                = 0xe
	NFT_MSG_NEWGEN                    = 0xf
	NFT_MSG_GETGEN                    = 0x10
	NFT_MSG_TRACE                     = 0x11
	NFT_MSG_NEWOBJ                    = 0x12
	NFT_MSG_GETOBJ                    = 0x13
	NFT_MSG_DELOBJ                    = 0x14
	NFT_MSG_GETOBJ_RESET              = 0x15
	NFT_MSG_MAX                       = 0x19
	NFTA_LIST_UNPEC                   = 0x0
	NFTA_LIST_ELEM                    = 0x1
	NFTA_HOOK_UNSPEC                  = 0x0
	NFTA_HOOK_HOOKNUM                 = 0x1
	NFTA_HOOK_PRIORITY                = 0x2
	NFTA_HOOK_DEV                     = 0x3
	NFT_TABLE_F_DORMANT               = 0x1
	NFTA_TABLE_UNSPEC                 = 0x0
	NFTA_TABLE_NAME                   = 0x1
	NFTA_TABLE_FLAGS                  = 0x2
	NFTA_TABLE_USE                    = 0x3
	NFTA_CHAIN_UNSPEC                 = 0x0
	NFTA_CHAIN_TABLE                  = 0x1
	NFTA_CHAIN_HANDLE                 = 0x2
	NFTA_CHAIN_NAME                   = 0x3
	NFTA_CHAIN_HOOK                   = 0x4
	NFTA_CHAIN_POLICY                 = 0x5
	NFTA_CHAIN_USE                    = 0x6
	NFTA_CHAIN_TYPE                   = 0x7
	NFTA_CHAIN_COUNTERS               = 0x8
	NFTA_CHAIN_PAD                    = 0x9
	NFTA_RULE_UNSPEC                  = 0x0
	NFTA_RULE_TABLE                   = 0x1
	NFTA_RULE_CHAIN                   = 0x2
	NFTA_RULE_HANDLE                  = 0x3
	NFTA_RULE_EXPRESSIONS             = 0x4
	NFTA_RULE_COMPAT                  = 0x5
	NFTA_RULE_POSITION                = 0x6
	NFTA_RULE_USERDATA                = 0x7
	NFTA_RULE_PAD                     = 0x8
	NFTA_RULE_ID                      = 0x9
	NFT_RULE_COMPAT_F_INV             = 0x2
	NFT_RULE_COMPAT_F_MASK            = 0x2
	NFTA_RULE_COMPAT_UNSPEC           = 0x0
	NFTA_RULE_COMPAT_PROTO            = 0x1
	NFTA_RULE_COMPAT_FLAGS            = 0x2
	NFT_SET_ANONYMOUS                 = 0x1
	NFT_SET_CONSTANT                  = 0x2
	NFT_SET_INTERVAL                  = 0x4
	NFT_SET_MAP                       = 0x8
	NFT_SET_TIMEOUT                   = 0x10
	NFT_SET_EVAL                      = 0x20
	NFT_SET_OBJECT                    = 0x40
	NFT_SET_POL_PERFORMANCE           = 0x0
	NFT_SET_POL_MEMORY                = 0x1
	NFTA_SET_DESC_UNSPEC              = 0x0
	NFTA_SET_DESC_SIZE                = 0x1
	NFTA_SET_UNSPEC                   = 0x0
	NFTA_SET_TABLE                    = 0x1
	NFTA_SET_NAME                     = 0x2
	NFTA_SET_FLAGS                    = 0x3
	NFTA_SET_KEY_TYPE                 = 0x4
	NFTA_SET_KEY_LEN                  = 0x5
	NFTA_SET_DATA_TYPE                = 0x6
	NFTA_SET_DATA_LEN                 = 0x7
	NFTA_SET_POLICY                   = 0x8
	NFTA_SET_DESC                     = 0x9
	NFTA_SET_ID                       = 0xa
	NFTA_SET_TIMEOUT                  = 0xb
	NFTA_SET_GC_INTERVAL              = 0xc
	NFTA_SET_USERDATA                 = 0xd
	NFTA_SET_PAD                      = 0xe
	NFTA_SET_OBJ_TYPE                 = 0xf
	NFT_SET_ELEM_INTERVAL_END         = 0x1
	NFTA_SET_ELEM_UNSPEC              = 0x0
	NFTA_SET_ELEM_KEY                 = 0x1
	NFTA_SET_ELEM_DATA                = 0x2
	NFTA_SET_ELEM_FLAGS               = 0x3
	NFTA_SET_ELEM_TIMEOUT             = 0x4
	NFTA_SET_ELEM_EXPIRATION          = 0x5
	NFTA_SET_ELEM_USERDATA            = 0x6
	NFTA_SET_ELEM_EXPR                = 0x7
	NFTA_SET_ELEM_PAD                 = 0x8
	NFTA_SET_ELEM_OBJREF              = 0x9
	NFTA_SET_ELEM_LIST_UNSPEC         = 0x0
	NFTA_SET_ELEM_LIST_TABLE          = 0x1
	NFTA_SET_ELEM_LIST_SET            = 0x2
	NFTA_SET_ELEM_LIST_ELEMENTS       = 0x3
	NFTA_SET_ELEM_LIST_SET_ID         = 0x4
	NFT_DATA_VALUE                    = 0x0
	NFT_DATA_VERDICT                  = 0xffffff00
	NFTA_DATA_UNSPEC                  = 0x0
	NFTA_DATA_VALUE                   = 0x1
	NFTA_DATA_VERDICT                 = 0x2
	NFTA_VERDICT_UNSPEC               = 0x0
	NFTA_VERDICT_CODE                 = 0x1
	NFTA_VERDICT_CHAIN                = 0x2
	NFTA_EXPR_UNSPEC                  = 0x0
	NFTA_EXPR_NAME                    = 0x1
	NFTA_EXPR_DATA                    = 0x2
	NFTA_IMMEDIATE_UNSPEC             = 0x0
	NFTA_IMMEDIATE_DREG               = 0x1
	NFTA_IMMEDIATE_DATA               = 0x2
	NFTA_BITWISE_UNSPEC               = 0x0
	NFTA_BITWISE_SREG                 = 0x1
	NFTA_BITWISE_DREG                 = 0x2
	NFTA_BITWISE_LEN                  = 0x3
	NFTA_BITWISE_MASK                 = 0x4
	NFTA_BITWISE_XOR                  = 0x5
	NFT_BYTEORDER_NTOH                = 0x0
	NFT_BYTEORDER_HTON                = 0x1
	NFTA_BYTEORDER_UNSPEC             = 0x0
	NFTA_BYTEORDER_SREG               = 0x1
	NFTA_BYTEORDER_DREG               = 0x2
	NFTA_BYTEORDER_OP                 = 0x3
	NFTA_BYTEORDER_LEN                = 0x4
	NFTA_BYTEORDER_SIZE               = 0x5
	NFT_CMP_EQ                        = 0x0
	NFT_CMP_NEQ                       = 0x1
	NFT_CMP_LT                        = 0x2
	NFT_CMP_LTE                       = 0x3
	NFT_CMP_GT                        = 0x4
	NFT_CMP_GTE                       = 0x5
	NFTA_CMP_UNSPEC                   = 0x0
	NFTA_CMP_SREG                     = 0x1
	NFTA_CMP_OP                       = 0x2
	NFTA_CMP_DATA                     = 0x3
	NFT_RANGE_EQ                      = 0x0
	NFT_RANGE_NEQ                     = 0x1
	NFTA_RANGE_UNSPEC                 = 0x0
	NFTA_RANGE_SREG                   = 0x1
	NFTA_RANGE_OP                     = 0x2
	NFTA_RANGE_FROM_DATA              = 0x3
	NFTA_RANGE_TO_DATA                = 0x4
	NFT_LOOKUP_F_INV                  = 0x1
	NFTA_LOOKUP_UNSPEC                = 0x0
	NFTA_LOOKUP_SET                   = 0x1
	NFTA_LOOKUP_SREG                  = 0x2
	NFTA_LOOKUP_DREG                  = 0x3
	NFTA_LOOKUP_SET_ID                = 0x4
	NFTA_LOOKUP_FLAGS                 = 0x5
	NFT_DYNSET_OP_ADD                 = 0x0
	NFT_DYNSET_OP_UPDATE              = 0x1
	NFT_DYNSET_F_INV                  = 0x1
	NFTA_DYNSET_UNSPEC                = 0x0
	NFTA_DYNSET_SET_NAME              = 0x1
	NFTA_DYNSET_SET_ID                = 0x2
	NFTA_DYNSET_OP                    = 0x3
	NFTA_DYNSET_SREG_KEY              = 0x4
	NFTA_DYNSET_SREG_DATA             = 0x5
	NFTA_DYNSET_TIMEOUT               = 0x6
	NFTA_DYNSET_EXPR                  = 0x7
	NFTA_DYNSET_PAD                   = 0x8
	NFTA_DYNSET_FLAGS                 = 0x9
	NFT_PAYLOAD_LL_HEADER             = 0x0
	NFT_PAYLOAD_NETWORK_HEADER        = 0x1
	NFT_PAYLOAD_TRANSPORT_HEADER      = 0x2
	NFT_PAYLOAD_CSUM_NONE             = 0x0
	NFT_PAYLOAD_CSUM_INET             = 0x1
	NFT_PAYLOAD_L4CSUM_PSEUDOHDR      = 0x1
	NFTA_PAYLOAD_UNSPEC               = 0x0
	NFTA_PAYLOAD_DREG                 = 0x1
	NFTA_PAYLOAD_BASE                 = 0x2
	NFTA_PAYLOAD_OFFSET               = 0x3
	NFTA_PAYLOAD_LEN                  = 0x4
	NFTA_PAYLOAD_SREG                 = 0x5
	NFTA_PAYLOAD_CSUM_TYPE            = 0x6
	NFTA_PAYLOAD_CSUM_OFFSET          = 0x7
	NFTA_PAYLOAD_CSUM_FLAGS           = 0x8
	NFT_EXTHDR_F_PRESENT              = 0x1
	NFT_EXTHDR_OP_IPV6                = 0x0
	NFT_EXTHDR_OP_TCPOPT              = 0x1
	NFTA_EXTHDR_UNSPEC                = 0x0
	NFTA_EXTHDR_DREG                  = 0x1
	NFTA_EXTHDR_TYPE                  = 0x2
	NFTA_EXTHDR_OFFSET                = 0x3
	NFTA_EXTHDR_LEN                   = 0x4
	NFTA_EXTHDR_FLAGS                 = 0x5
	NFTA_EXTHDR_OP                    = 0x6
	NFTA_EXTHDR_SREG                  = 0x7
	NFT_META_LEN                      = 0x0
	NFT_META_PROTOCOL                 = 0x1
	NFT_META_PRIORITY                 = 0x2
	NFT_META_MARK                     = 0x3
	NFT_META_IIF                      = 0x4
	NFT_META_OIF                      = 0x5
	NFT_META_IIFNAME                  = 0x6
	NFT_META_OIFNAME                  = 0x7
	NFT_META_IIFTYPE                  = 0x8
	NFT_META_OIFTYPE                  = 0x9
	NFT_META_SKUID                    = 0xa
	NFT_META_SKGID                    = 0xb
	NFT_META_NFTRACE                  = 0xc
	NFT_META_RTCLASSID                = 0xd
	NFT_META_SECMARK                  = 0xe
	NFT_META_NFPROTO                  = 0xf
	NFT_META_L4PROTO                  = 0x10
	NFT_META_BRI_IIFNAME              = 0x11
	NFT_META_BRI_OIFNAME              = 0x12
	NFT_META_PKTTYPE                  = 0x13
	NFT_META_CPU                      = 0x14
	NFT_META_IIFGROUP                 = 0x15
	NFT_META_OIFGROUP                 = 0x16
	NFT_META_CGROUP                   = 0x17
	NFT_META_PRANDOM                  = 0x18
	NFT_RT_CLASSID                    = 0x0
	NFT_RT_NEXTHOP4                   = 0x1
	NFT_RT_NEXTHOP6                   = 0x2
	NFT_RT_TCPMSS                     = 0x3
	NFT_HASH_JENKINS                  = 0x0
	NFT_HASH_SYM                      = 0x1
	NFTA_HASH_UNSPEC                  = 0x0
	NFTA_HASH_SREG                    = 0x1
	NFTA_HASH_DREG                    = 0x2
	NFTA_HASH_LEN                     = 0x3
	NFTA_HASH_MODULUS                 = 0x4
	NFTA_HASH_SEED                    = 0x5
	NFTA_HASH_OFFSET                  = 0x6
	NFTA_HASH_TYPE                    = 0x7
	NFTA_META_UNSPEC                  = 0x0
	NFTA_META_DREG                    = 0x1
	NFTA_META_KEY                     = 0x2
	NFTA_META_SREG                    = 0x3
	NFTA_RT_UNSPEC                    = 0x0
	NFTA_RT_DREG                      = 0x1
	NFTA_RT_KEY                       = 0x2
	NFT_CT_STATE                      = 0x0
	NFT_CT_DIRECTION                  = 0x1
	NFT_CT_STATUS                     = 0x2
	NFT_CT_MARK                       = 0x3
	NFT_CT_SECMARK                    = 0x4
	NFT_CT_EXPIRATION                 = 0x5
	NFT_CT_HELPER                     = 0x6
	NFT_CT_L3PROTOCOL                 = 0x7
	NFT_CT_SRC                        = 0x8
	NFT_CT_DST                        = 0x9
	NFT_CT_PROTOCOL                   = 0xa
	NFT_CT_PROTO_SRC                  = 0xb
	NFT_CT_PROTO_DST                  = 0xc
	NFT_CT_LABELS                     = 0xd
	NFT_CT_PKTS                       = 0xe
	NFT_CT_BYTES                      = 0xf
	NFT_CT_AVGPKT                     = 0x10
	NFT_CT_ZONE                       = 0x11
	NFT_CT_EVENTMASK                  = 0x12
	NFTA_CT_UNSPEC                    = 0x0
	NFTA_CT_DREG                      = 0x1
	NFTA_CT_KEY                       = 0x2
	NFTA_CT_DIRECTION                 = 0x3
	NFTA_CT_SREG                      = 0x4
	NFT_LIMIT_PKTS                    = 0x0
	NFT_LIMIT_PKT_BYTES               = 0x1
	NFT_LIMIT_F_INV                   = 0x1
	NFTA_LIMIT_UNSPEC                 = 0x0
	NFTA_LIMIT_RATE                   = 0x1
	NFTA_LIMIT_UNIT                   = 0x2
	NFTA_LIMIT_BURST                  = 0x3
	NFTA_LIMIT_TYPE                   = 0x4
	NFTA_LIMIT_FLAGS                  = 0x5
	NFTA_LIMIT_PAD                    = 0x6
	NFTA_COUNTER_UNSPEC               = 0x0
	NFTA_COUNTER_BYTES                = 0x1
	NFTA_COUNTER_PACKETS              = 0x2
	NFTA_COUNTER_PAD                  = 0x3
	NFTA_LOG_UNSPEC                   = 0x0
	NFTA_LOG_GROUP                    = 0x1
	NFTA_LOG_PREFIX                   = 0x2
	NFTA_LOG_SNAPLEN                  = 0x3
	NFTA_LOG_QTHRESHOLD               = 0x4
	NFTA_LOG_LEVEL                    = 0x5
	NFTA_LOG_FLAGS                    = 0x6
	NFTA_QUEUE_UNSPEC                 = 0x0
	NFTA_QUEUE_NUM                    = 0x1
	NFTA_QUEUE_TOTAL                  = 0x2
	NFTA_QUEUE_FLAGS                  = 0x3
	NFTA_QUEUE_SREG_QNUM              = 0x4
	NFT_QUOTA_F_INV                   = 0x1
	NFT_QUOTA_F_DEPLETED              = 0x2
	NFTA_QUOTA_UNSPEC                 = 0x0
	NFTA_QUOTA_BYTES                  = 0x1
	NFTA_QUOTA_FLAGS                  = 0x2
	NFTA_QUOTA_PAD                    = 0x3
	NFTA_QUOTA_CONSUMED               = 0x4
	NFT_REJECT_ICMP_UNREACH           = 0x0
	NFT_REJECT_TCP_RST                = 0x1
	NFT_REJECT_ICMPX_UNREACH          = 0x2
	NFT_REJECT_ICMPX_NO_ROUTE         = 0x0
	NFT_REJECT_ICMPX_PORT_UNREACH     = 0x1
	NFT_REJECT_ICMPX_HOST_UNREACH     = 0x2
	NFT_REJECT_ICMPX_ADMIN_PROHIBITED = 0x3
	NFTA_REJECT_UNSPEC                = 0x0
	NFTA_REJECT_TYPE                  = 0x1
	NFTA_REJECT_ICMP_CODE             = 0x2
	NFT_NAT_SNAT                      = 0x0
	NFT_NAT_DNAT                      = 0x1
	NFTA_NAT_UNSPEC                   = 0x0
	NFTA_NAT_TYPE                     = 0x1
	NFTA_NAT_FAMILY                   = 0x2
	NFTA_NAT_REG_ADDR_MIN             = 0x3
	NFTA_NAT_REG_ADDR_MAX             = 0x4
	NFTA_NAT_REG_PROTO_MIN            = 0x5
	NFTA_NAT_REG_PROTO_MAX            = 0x6
	NFTA_NAT_FLAGS                    = 0x7
	NFTA_MASQ_UNSPEC                  = 0x0
	NFTA_MASQ_FLAGS                   = 0x1
	NFTA_MASQ_REG_PROTO_MIN           = 0x2
	NFTA_MASQ_REG_PROTO_MAX           = 0x3
	NFTA_REDIR_UNSPEC                 = 0x0
	NFTA_REDIR_REG_PROTO_MIN          = 0x1
	NFTA_REDIR_REG_PROTO_MAX          = 0x2
	NFTA_REDIR_FLAGS                  = 0x3
	NFTA_DUP_UNSPEC                   = 0x0
	NFTA_DUP_SREG_ADDR                = 0x1
	NFTA_DUP_SREG_DEV                 = 0x2
	NFTA_FWD_UNSPEC                   = 0x0
	NFTA_FWD_SREG_DEV                 = 0x1
	NFTA_OBJREF_UNSPEC                = 0x0
	NFTA_OBJREF_IMM_TYPE              = 0x1
	NFTA_OBJREF_IMM_NAME              = 0x2
	NFTA_OBJREF_SET_SREG              = 0x3
	NFTA_OBJREF_SET_NAME              = 0x4
	NFTA_OBJREF_SET_ID                = 0x5
	NFTA_GEN_UNSPEC                   = 0x0
	NFTA_GEN_ID                       = 0x1
	NFTA_GEN_PROC_PID                 = 0x2
	NFTA_GEN_PROC_NAME                = 0x3
	NFTA_FIB_UNSPEC                   = 0x0
	NFTA_FIB_DREG                     = 0x1
	NFTA_FIB_RESULT                   = 0x2
	NFTA_FIB_FLAGS                    = 0x3
	NFT_FIB_RESULT_UNSPEC             = 0x0
	NFT_FIB_RESULT_OIF                = 0x1
	NFT_FIB_RESULT_OIFNAME            = 0x2
	NFT_FIB_RESULT_ADDRTYPE           = 0x3
	NFTA_FIB_F_SADDR                  = 0x1
	NFTA_FIB_F_DADDR                  = 0x2
	NFTA_FIB_F_MARK                   = 0x4
	NFTA_FIB_F_IIF                    = 0x8
	NFTA_FIB_F_OIF                    = 0x10
	NFTA_FIB_F_PRESENT                = 0x20
	NFTA_CT_HELPER_UNSPEC             = 0x0
	NFTA_CT_HELPER_NAME               = 0x1
	NFTA_CT_HELPER_L3PROTO            = 0x2
	NFTA_CT_HELPER_L4PROTO            = 0x3
	NFTA_OBJ_UNSPEC                   = 0x0
	NFTA_OBJ_TABLE                    = 0x1
	NFTA_OBJ_NAME                     = 0x2
	NFTA_OBJ_TYPE                     = 0x3
	NFTA_OBJ_DATA                     = 0x4
	NFTA_OBJ_USE                      = 0x5
	NFTA_TRACE_UNSPEC                 = 0x0
	NFTA_TRACE_TABLE                  = 0x1
	NFTA_TRACE_CHAIN                  = 0x2
	NFTA_TRACE_RULE_HANDLE            = 0x3
	NFTA_TRACE_TYPE                   = 0x4
	NFTA_TRACE_VERDICT                = 0x5
	NFTA_TRACE_ID                     = 0x6
	NFTA_TRACE_LL_HEADER              = 0x7
	NFTA_TRACE_NETWORK_HEADER         = 0x8
	NFTA_TRACE_TRANSPORT_HEADER       = 0x9
	NFTA_TRACE_IIF                    = 0xa
	NFTA_TRACE_IIFTYPE                = 0xb
	NFTA_TRACE_OIF                    = 0xc
	NFTA_TRACE_OIFTYPE                = 0xd
	NFTA_TRACE_MARK                   = 0xe
	NFTA_TRACE_NFPROTO                = 0xf
	NFTA_TRACE_POLICY                 = 0x10
	NFTA_TRACE_PAD                    = 0x11
	NFT_TRACETYPE_UNSPEC              = 0x0
	NFT_TRACETYPE_POLICY              = 0x1
	NFT_TRACETYPE_RETURN              = 0x2
	NFT_TRACETYPE_RULE                = 0x3
	NFTA_NG_UNSPEC                    = 0x0
	NFTA_NG_DREG                      = 0x1
	NFTA_NG_MODULUS                   = 0x2
	NFTA_NG_TYPE                      = 0x3
	NFTA_NG_OFFSET                    = 0x4
	NFT_NG_INCREMENTAL                = 0x0
	NFT_NG_RANDOM                     = 0x1
)

type RTCTime struct {
	Sec   int32
	Min   int32
	Hour  int32
	Mday  int32
	Mon   int32
	Year  int32
	Wday  int32
	Yday  int32
	Isdst int32
}

type RTCWkAlrm struct {
	Enabled uint8
	Pending uint8
	_       [2]byte
	Time    RTCTime
}

type RTCPLLInfo struct {
	Ctrl    int32
	Value   int32
	Max     int32
	Min     int32
	Posmult int32
	Negmult int32
	Clock   int32
}

type BlkpgIoctlArg struct {
	Op      int32
	Flags   int32
	Datalen int32
	Data    *byte
}

type BlkpgPartition struct {
	Start   int64
	Length  int64
	Pno     int32
	Devname [64]uint8
	Volname [64]uint8
	_       [4]byte
}

const (
	BLKPG                  = 0x20001269
	BLKPG_ADD_PARTITION    = 0x1
	BLKPG_DEL_PARTITION    = 0x2
	BLKPG_RESIZE_PARTITION = 0x3
)

const (
	NETNSA_NONE = 0x0
	NETNSA_NSID = 0x1
	NETNSA_PID  = 0x2
	NETNSA_FD   = 0x3
)

type XDPRingOffset struct {
	Producer uint64
	Consumer uint64
	Desc     uint64
}

type XDPMmapOffsets struct {
	Rx XDPRingOffset
	Tx XDPRingOffset
	Fr XDPRingOffset
	Cr XDPRingOffset
}

type XDPUmemReg struct {
	Addr     uint64
	Len      uint64
	Size     uint32
	Headroom uint32
}

type XDPStatistics struct {
	Rx_dropped       uint64
	Rx_invalid_descs uint64
	Tx_invalid_descs uint64
}

type XDPDesc struct {
	Addr    uint64
	Len     uint32
	Options uint32
}

const (
	NCSI_CMD_UNSPEC                 = 0x0
	NCSI_CMD_PKG_INFO               = 0x1
	NCSI_CMD_SET_INTERFACE          = 0x2
	NCSI_CMD_CLEAR_INTERFACE        = 0x3
	NCSI_ATTR_UNSPEC                = 0x0
	NCSI_ATTR_IFINDEX               = 0x1
	NCSI_ATTR_PACKAGE_LIST          = 0x2
	NCSI_ATTR_PACKAGE_ID            = 0x3
	NCSI_ATTR_CHANNEL_ID            = 0x4
	NCSI_PKG_ATTR_UNSPEC            = 0x0
	NCSI_PKG_ATTR                   = 0x1
	NCSI_PKG_ATTR_ID                = 0x2
	NCSI_PKG_ATTR_FORCED            = 0x3
	NCSI_PKG_ATTR_CHANNEL_LIST      = 0x4
	NCSI_CHANNEL_ATTR_UNSPEC        = 0x0
	NCSI_CHANNEL_ATTR               = 0x1
	NCSI_CHANNEL_ATTR_ID            = 0x2
	NCSI_CHANNEL_ATTR_VERSION_MAJOR = 0x3
	NCSI_CHANNEL_ATTR_VERSION_MINOR = 0x4
	NCSI_CHANNEL_ATTR_VERSION_STR   = 0x5
	NCSI_CHANNEL_ATTR_LINK_STATE    = 0x6
	NCSI_CHANNEL_ATTR_ACTIVE        = 0x7
	NCSI_CHANNEL_ATTR_FORCED        = 0x8
	NCSI_CHANNEL_ATTR_VLAN_LIST     = 0x9
	NCSI_CHANNEL_ATTR_VLAN_ID       = 0xa
)
