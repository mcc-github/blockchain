







package unix

import "syscall"

const (
	AF_APPLETALK                  = 0x10
	AF_BYPASS                     = 0x19
	AF_CCITT                      = 0xa
	AF_CHAOS                      = 0x5
	AF_DATAKIT                    = 0x9
	AF_DECnet                     = 0xc
	AF_DLI                        = 0xd
	AF_ECMA                       = 0x8
	AF_HYLINK                     = 0xf
	AF_IMPLINK                    = 0x3
	AF_INET                       = 0x2
	AF_INET6                      = 0x18
	AF_INTF                       = 0x14
	AF_ISO                        = 0x7
	AF_LAT                        = 0xe
	AF_LINK                       = 0x12
	AF_LOCAL                      = 0x1
	AF_MAX                        = 0x1e
	AF_NDD                        = 0x17
	AF_NETWARE                    = 0x16
	AF_NS                         = 0x6
	AF_OSI                        = 0x7
	AF_PUP                        = 0x4
	AF_RIF                        = 0x15
	AF_ROUTE                      = 0x11
	AF_SNA                        = 0xb
	AF_UNIX                       = 0x1
	AF_UNSPEC                     = 0x0
	ALTWERASE                     = 0x400000
	ARPHRD_802_3                  = 0x6
	ARPHRD_802_5                  = 0x6
	ARPHRD_ETHER                  = 0x1
	ARPHRD_FDDI                   = 0x1
	B0                            = 0x0
	B110                          = 0x3
	B1200                         = 0x9
	B134                          = 0x4
	B150                          = 0x5
	B1800                         = 0xa
	B19200                        = 0xe
	B200                          = 0x6
	B2400                         = 0xb
	B300                          = 0x7
	B38400                        = 0xf
	B4800                         = 0xc
	B50                           = 0x1
	B600                          = 0x8
	B75                           = 0x2
	B9600                         = 0xd
	BRKINT                        = 0x2
	BS0                           = 0x0
	BS1                           = 0x1000
	BSDLY                         = 0x1000
	CAP_AACCT                     = 0x6
	CAP_ARM_APPLICATION           = 0x5
	CAP_BYPASS_RAC_VMM            = 0x3
	CAP_CLEAR                     = 0x0
	CAP_CREDENTIALS               = 0x7
	CAP_EFFECTIVE                 = 0x1
	CAP_EWLM_AGENT                = 0x4
	CAP_INHERITABLE               = 0x2
	CAP_MAXIMUM                   = 0x7
	CAP_NUMA_ATTACH               = 0x2
	CAP_PERMITTED                 = 0x3
	CAP_PROPAGATE                 = 0x1
	CAP_PROPOGATE                 = 0x1
	CAP_SET                       = 0x1
	CBAUD                         = 0xf
	CFLUSH                        = 0xf
	CIBAUD                        = 0xf0000
	CLOCAL                        = 0x800
	CLOCK_MONOTONIC               = 0xa
	CLOCK_PROCESS_CPUTIME_ID      = 0xb
	CLOCK_REALTIME                = 0x9
	CLOCK_THREAD_CPUTIME_ID       = 0xc
	CR0                           = 0x0
	CR1                           = 0x100
	CR2                           = 0x200
	CR3                           = 0x300
	CRDLY                         = 0x300
	CREAD                         = 0x80
	CS5                           = 0x0
	CS6                           = 0x10
	CS7                           = 0x20
	CS8                           = 0x30
	CSIOCGIFCONF                  = -0x3fef96dc
	CSIZE                         = 0x30
	CSMAP_DIR                     = "/usr/lib/nls/csmap/"
	CSTART                        = '\021'
	CSTOP                         = '\023'
	CSTOPB                        = 0x40
	CSUSP                         = 0x1a
	ECHO                          = 0x8
	ECHOCTL                       = 0x20000
	ECHOE                         = 0x10
	ECHOK                         = 0x20
	ECHOKE                        = 0x80000
	ECHONL                        = 0x40
	ECHOPRT                       = 0x40000
	ECH_ICMPID                    = 0x2
	ETHERNET_CSMACD               = 0x6
	EVENP                         = 0x80
	EXCONTINUE                    = 0x0
	EXDLOK                        = 0x3
	EXIO                          = 0x2
	EXPGIO                        = 0x0
	EXRESUME                      = 0x2
	EXRETURN                      = 0x1
	EXSIG                         = 0x4
	EXTA                          = 0xe
	EXTB                          = 0xf
	EXTRAP                        = 0x1
	EYEC_RTENTRYA                 = 0x257274656e747241
	EYEC_RTENTRYF                 = 0x257274656e747246
	E_ACC                         = 0x0
	FD_CLOEXEC                    = 0x1
	FD_SETSIZE                    = 0xfffe
	FF0                           = 0x0
	FF1                           = 0x2000
	FFDLY                         = 0x2000
	FLUSHBAND                     = 0x40
	FLUSHLOW                      = 0x8
	FLUSHO                        = 0x100000
	FLUSHR                        = 0x1
	FLUSHRW                       = 0x3
	FLUSHW                        = 0x2
	F_CLOSEM                      = 0xa
	F_DUP2FD                      = 0xe
	F_DUPFD                       = 0x0
	F_GETFD                       = 0x1
	F_GETFL                       = 0x3
	F_GETLK                       = 0xb
	F_GETLK64                     = 0xb
	F_GETOWN                      = 0x8
	F_LOCK                        = 0x1
	F_OK                          = 0x0
	F_RDLCK                       = 0x1
	F_SETFD                       = 0x2
	F_SETFL                       = 0x4
	F_SETLK                       = 0xc
	F_SETLK64                     = 0xc
	F_SETLKW                      = 0xd
	F_SETLKW64                    = 0xd
	F_SETOWN                      = 0x9
	F_TEST                        = 0x3
	F_TLOCK                       = 0x2
	F_TSTLK                       = 0xf
	F_ULOCK                       = 0x0
	F_UNLCK                       = 0x3
	F_WRLCK                       = 0x2
	HUPCL                         = 0x400
	IBSHIFT                       = 0x10
	ICANON                        = 0x2
	ICMP6_FILTER                  = 0x26
	ICMP6_SEC_SEND_DEL            = 0x46
	ICMP6_SEC_SEND_GET            = 0x47
	ICMP6_SEC_SEND_SET            = 0x44
	ICMP6_SEC_SEND_SET_CGA_ADDR   = 0x45
	ICRNL                         = 0x100
	IEXTEN                        = 0x200000
	IFA_FIRSTALIAS                = 0x2000
	IFA_ROUTE                     = 0x1
	IFF_64BIT                     = 0x4000000
	IFF_ALLCAST                   = 0x20000
	IFF_ALLMULTI                  = 0x200
	IFF_BPF                       = 0x8000000
	IFF_BRIDGE                    = 0x40000
	IFF_BROADCAST                 = 0x2
	IFF_CANTCHANGE                = 0x80c52
	IFF_CHECKSUM_OFFLOAD          = 0x10000000
	IFF_D1                        = 0x8000
	IFF_D2                        = 0x4000
	IFF_D3                        = 0x2000
	IFF_D4                        = 0x1000
	IFF_DEBUG                     = 0x4
	IFF_DEVHEALTH                 = 0x4000
	IFF_DO_HW_LOOPBACK            = 0x10000
	IFF_GROUP_ROUTING             = 0x2000000
	IFF_IFBUFMGT                  = 0x800000
	IFF_LINK0                     = 0x100000
	IFF_LINK1                     = 0x200000
	IFF_LINK2                     = 0x400000
	IFF_LOOPBACK                  = 0x8
	IFF_MULTICAST                 = 0x80000
	IFF_NOARP                     = 0x80
	IFF_NOECHO                    = 0x800
	IFF_NOTRAILERS                = 0x20
	IFF_OACTIVE                   = 0x400
	IFF_POINTOPOINT               = 0x10
	IFF_PROMISC                   = 0x100
	IFF_PSEG                      = 0x40000000
	IFF_RUNNING                   = 0x40
	IFF_SIMPLEX                   = 0x800
	IFF_SNAP                      = 0x8000
	IFF_TCP_DISABLE_CKSUM         = 0x20000000
	IFF_TCP_NOCKSUM               = 0x1000000
	IFF_UP                        = 0x1
	IFF_VIPA                      = 0x80000000
	IFNAMSIZ                      = 0x10
	IFO_FLUSH                     = 0x1
	IFT_1822                      = 0x2
	IFT_AAL5                      = 0x31
	IFT_ARCNET                    = 0x23
	IFT_ARCNETPLUS                = 0x24
	IFT_ATM                       = 0x25
	IFT_CEPT                      = 0x13
	IFT_CLUSTER                   = 0x3e
	IFT_DS3                       = 0x1e
	IFT_EON                       = 0x19
	IFT_ETHER                     = 0x6
	IFT_FCS                       = 0x3a
	IFT_FDDI                      = 0xf
	IFT_FRELAY                    = 0x20
	IFT_FRELAYDCE                 = 0x2c
	IFT_GIFTUNNEL                 = 0x3c
	IFT_HDH1822                   = 0x3
	IFT_HF                        = 0x3d
	IFT_HIPPI                     = 0x2f
	IFT_HSSI                      = 0x2e
	IFT_HY                        = 0xe
	IFT_IB                        = 0xc7
	IFT_ISDNBASIC                 = 0x14
	IFT_ISDNPRIMARY               = 0x15
	IFT_ISO88022LLC               = 0x29
	IFT_ISO88023                  = 0x7
	IFT_ISO88024                  = 0x8
	IFT_ISO88025                  = 0x9
	IFT_ISO88026                  = 0xa
	IFT_LAPB                      = 0x10
	IFT_LOCALTALK                 = 0x2a
	IFT_LOOP                      = 0x18
	IFT_MIOX25                    = 0x26
	IFT_MODEM                     = 0x30
	IFT_NSIP                      = 0x1b
	IFT_OTHER                     = 0x1
	IFT_P10                       = 0xc
	IFT_P80                       = 0xd
	IFT_PARA                      = 0x22
	IFT_PPP                       = 0x17
	IFT_PROPMUX                   = 0x36
	IFT_PROPVIRTUAL               = 0x35
	IFT_PTPSERIAL                 = 0x16
	IFT_RS232                     = 0x21
	IFT_SDLC                      = 0x11
	IFT_SIP                       = 0x1f
	IFT_SLIP                      = 0x1c
	IFT_SMDSDXI                   = 0x2b
	IFT_SMDSICIP                  = 0x34
	IFT_SN                        = 0x38
	IFT_SONET                     = 0x27
	IFT_SONETPATH                 = 0x32
	IFT_SONETVT                   = 0x33
	IFT_SP                        = 0x39
	IFT_STARLAN                   = 0xb
	IFT_T1                        = 0x12
	IFT_TUNNEL                    = 0x3b
	IFT_ULTRA                     = 0x1d
	IFT_V35                       = 0x2d
	IFT_VIPA                      = 0x37
	IFT_X25                       = 0x5
	IFT_X25DDN                    = 0x4
	IFT_X25PLE                    = 0x28
	IFT_XETHER                    = 0x1a
	IGNBRK                        = 0x1
	IGNCR                         = 0x80
	IGNPAR                        = 0x4
	IMAXBEL                       = 0x10000
	INLCR                         = 0x40
	INPCK                         = 0x10
	IN_CLASSA_HOST                = 0xffffff
	IN_CLASSA_MAX                 = 0x80
	IN_CLASSA_NET                 = 0xff000000
	IN_CLASSA_NSHIFT              = 0x18
	IN_CLASSB_HOST                = 0xffff
	IN_CLASSB_MAX                 = 0x10000
	IN_CLASSB_NET                 = 0xffff0000
	IN_CLASSB_NSHIFT              = 0x10
	IN_CLASSC_HOST                = 0xff
	IN_CLASSC_NET                 = 0xffffff00
	IN_CLASSC_NSHIFT              = 0x8
	IN_CLASSD_HOST                = 0xfffffff
	IN_CLASSD_NET                 = 0xf0000000
	IN_CLASSD_NSHIFT              = 0x1c
	IN_LOOPBACKNET                = 0x7f
	IN_USE                        = 0x1
	IPPROTO_AH                    = 0x33
	IPPROTO_BIP                   = 0x53
	IPPROTO_DSTOPTS               = 0x3c
	IPPROTO_EGP                   = 0x8
	IPPROTO_EON                   = 0x50
	IPPROTO_ESP                   = 0x32
	IPPROTO_FRAGMENT              = 0x2c
	IPPROTO_GGP                   = 0x3
	IPPROTO_GIF                   = 0x8c
	IPPROTO_GRE                   = 0x2f
	IPPROTO_HOPOPTS               = 0x0
	IPPROTO_ICMP                  = 0x1
	IPPROTO_ICMPV6                = 0x3a
	IPPROTO_IDP                   = 0x16
	IPPROTO_IGMP                  = 0x2
	IPPROTO_IP                    = 0x0
	IPPROTO_IPIP                  = 0x4
	IPPROTO_IPV6                  = 0x29
	IPPROTO_LOCAL                 = 0x3f
	IPPROTO_MAX                   = 0x100
	IPPROTO_MH                    = 0x87
	IPPROTO_NONE                  = 0x3b
	IPPROTO_PUP                   = 0xc
	IPPROTO_QOS                   = 0x2d
	IPPROTO_RAW                   = 0xff
	IPPROTO_ROUTING               = 0x2b
	IPPROTO_RSVP                  = 0x2e
	IPPROTO_SCTP                  = 0x84
	IPPROTO_TCP                   = 0x6
	IPPROTO_TP                    = 0x1d
	IPPROTO_UDP                   = 0x11
	IPV6_ADDRFORM                 = 0x16
	IPV6_ADDR_PREFERENCES         = 0x4a
	IPV6_ADD_MEMBERSHIP           = 0xc
	IPV6_AIXRAWSOCKET             = 0x39
	IPV6_CHECKSUM                 = 0x27
	IPV6_DONTFRAG                 = 0x2d
	IPV6_DROP_MEMBERSHIP          = 0xd
	IPV6_DSTOPTS                  = 0x36
	IPV6_FLOWINFO_FLOWLABEL       = 0xffffff
	IPV6_FLOWINFO_PRIFLOW         = 0xfffffff
	IPV6_FLOWINFO_PRIORITY        = 0xf000000
	IPV6_FLOWINFO_SRFLAG          = 0x10000000
	IPV6_FLOWINFO_VERSION         = 0xf0000000
	IPV6_HOPLIMIT                 = 0x28
	IPV6_HOPOPTS                  = 0x34
	IPV6_JOIN_GROUP               = 0xc
	IPV6_LEAVE_GROUP              = 0xd
	IPV6_MIPDSTOPTS               = 0x36
	IPV6_MULTICAST_HOPS           = 0xa
	IPV6_MULTICAST_IF             = 0x9
	IPV6_MULTICAST_LOOP           = 0xb
	IPV6_NEXTHOP                  = 0x30
	IPV6_NOPROBE                  = 0x1c
	IPV6_PATHMTU                  = 0x2e
	IPV6_PKTINFO                  = 0x21
	IPV6_PKTOPTIONS               = 0x24
	IPV6_PRIORITY_10              = 0xa000000
	IPV6_PRIORITY_11              = 0xb000000
	IPV6_PRIORITY_12              = 0xc000000
	IPV6_PRIORITY_13              = 0xd000000
	IPV6_PRIORITY_14              = 0xe000000
	IPV6_PRIORITY_15              = 0xf000000
	IPV6_PRIORITY_8               = 0x8000000
	IPV6_PRIORITY_9               = 0x9000000
	IPV6_PRIORITY_BULK            = 0x4000000
	IPV6_PRIORITY_CONTROL         = 0x7000000
	IPV6_PRIORITY_FILLER          = 0x1000000
	IPV6_PRIORITY_INTERACTIVE     = 0x6000000
	IPV6_PRIORITY_RESERVED1       = 0x3000000
	IPV6_PRIORITY_RESERVED2       = 0x5000000
	IPV6_PRIORITY_UNATTENDED      = 0x2000000
	IPV6_PRIORITY_UNCHARACTERIZED = 0x0
	IPV6_RECVDSTOPTS              = 0x38
	IPV6_RECVHOPLIMIT             = 0x29
	IPV6_RECVHOPOPTS              = 0x35
	IPV6_RECVHOPS                 = 0x22
	IPV6_RECVIF                   = 0x1e
	IPV6_RECVPATHMTU              = 0x2f
	IPV6_RECVPKTINFO              = 0x23
	IPV6_RECVRTHDR                = 0x33
	IPV6_RECVSRCRT                = 0x1d
	IPV6_RECVTCLASS               = 0x2a
	IPV6_RTHDR                    = 0x32
	IPV6_RTHDRDSTOPTS             = 0x37
	IPV6_RTHDR_TYPE_0             = 0x0
	IPV6_RTHDR_TYPE_2             = 0x2
	IPV6_SENDIF                   = 0x1f
	IPV6_SRFLAG_LOOSE             = 0x0
	IPV6_SRFLAG_STRICT            = 0x10000000
	IPV6_TCLASS                   = 0x2b
	IPV6_TOKEN_LENGTH             = 0x40
	IPV6_UNICAST_HOPS             = 0x4
	IPV6_USE_MIN_MTU              = 0x2c
	IPV6_V6ONLY                   = 0x25
	IPV6_VERSION                  = 0x60000000
	IP_ADDRFORM                   = 0x16
	IP_ADD_MEMBERSHIP             = 0xc
	IP_ADD_SOURCE_MEMBERSHIP      = 0x3c
	IP_BLOCK_SOURCE               = 0x3a
	IP_BROADCAST_IF               = 0x10
	IP_CACHE_LINE_SIZE            = 0x80
	IP_DEFAULT_MULTICAST_LOOP     = 0x1
	IP_DEFAULT_MULTICAST_TTL      = 0x1
	IP_DF                         = 0x4000
	IP_DHCPMODE                   = 0x11
	IP_DONTFRAG                   = 0x19
	IP_DROP_MEMBERSHIP            = 0xd
	IP_DROP_SOURCE_MEMBERSHIP     = 0x3d
	IP_FINDPMTU                   = 0x1a
	IP_HDRINCL                    = 0x2
	IP_INC_MEMBERSHIPS            = 0x14
	IP_INIT_MEMBERSHIP            = 0x14
	IP_MAXPACKET                  = 0xffff
	IP_MF                         = 0x2000
	IP_MSS                        = 0x240
	IP_MULTICAST_HOPS             = 0xa
	IP_MULTICAST_IF               = 0x9
	IP_MULTICAST_LOOP             = 0xb
	IP_MULTICAST_TTL              = 0xa
	IP_OPT                        = 0x1b
	IP_OPTIONS                    = 0x1
	IP_PMTUAGE                    = 0x1b
	IP_RECVDSTADDR                = 0x7
	IP_RECVIF                     = 0x14
	IP_RECVIFINFO                 = 0xf
	IP_RECVINTERFACE              = 0x20
	IP_RECVMACHDR                 = 0xe
	IP_RECVOPTS                   = 0x5
	IP_RECVRETOPTS                = 0x6
	IP_RECVTTL                    = 0x22
	IP_RETOPTS                    = 0x8
	IP_SOURCE_FILTER              = 0x48
	IP_TOS                        = 0x3
	IP_TTL                        = 0x4
	IP_UNBLOCK_SOURCE             = 0x3b
	IP_UNICAST_HOPS               = 0x4
	ISIG                          = 0x1
	ISTRIP                        = 0x20
	IUCLC                         = 0x800
	IXANY                         = 0x1000
	IXOFF                         = 0x400
	IXON                          = 0x200
	I_FLUSH                       = 0x20005305
	LNOFLSH                       = 0x8000
	LOCK_EX                       = 0x2
	LOCK_NB                       = 0x4
	LOCK_SH                       = 0x1
	LOCK_UN                       = 0x8
	MADV_DONTNEED                 = 0x4
	MADV_NORMAL                   = 0x0
	MADV_RANDOM                   = 0x1
	MADV_SEQUENTIAL               = 0x2
	MADV_SPACEAVAIL               = 0x5
	MADV_WILLNEED                 = 0x3
	MAP_ANON                      = 0x10
	MAP_ANONYMOUS                 = 0x10
	MAP_FILE                      = 0x0
	MAP_FIXED                     = 0x100
	MAP_PRIVATE                   = 0x2
	MAP_SHARED                    = 0x1
	MAP_TYPE                      = 0xf0
	MAP_VARIABLE                  = 0x0
	MCL_CURRENT                   = 0x100
	MCL_FUTURE                    = 0x200
	MSG_ANY                       = 0x4
	MSG_ARGEXT                    = 0x400
	MSG_BAND                      = 0x2
	MSG_COMPAT                    = 0x8000
	MSG_CTRUNC                    = 0x20
	MSG_DONTROUTE                 = 0x4
	MSG_EOR                       = 0x8
	MSG_HIPRI                     = 0x1
	MSG_MAXIOVLEN                 = 0x10
	MSG_MPEG2                     = 0x80
	MSG_NONBLOCK                  = 0x4000
	MSG_NOSIGNAL                  = 0x100
	MSG_OOB                       = 0x1
	MSG_PEEK                      = 0x2
	MSG_TRUNC                     = 0x10
	MSG_WAITALL                   = 0x40
	MSG_WAITFORONE                = 0x200
	MS_ASYNC                      = 0x10
	MS_EINTR                      = 0x80
	MS_INVALIDATE                 = 0x40
	MS_PER_SEC                    = 0x3e8
	MS_SYNC                       = 0x20
	NL0                           = 0x0
	NL1                           = 0x4000
	NL2                           = 0x8000
	NL3                           = 0xc000
	NLDLY                         = 0x4000
	NOFLSH                        = 0x80
	NOFLUSH                       = 0x80000000
	OCRNL                         = 0x8
	OFDEL                         = 0x80
	OFILL                         = 0x40
	OLCUC                         = 0x2
	ONLCR                         = 0x4
	ONLRET                        = 0x20
	ONOCR                         = 0x10
	ONOEOT                        = 0x80000
	OPOST                         = 0x1
	OXTABS                        = 0x40000
	O_ACCMODE                     = 0x23
	O_APPEND                      = 0x8
	O_CIO                         = 0x80
	O_CIOR                        = 0x800000000
	O_CLOEXEC                     = 0x800000
	O_CREAT                       = 0x100
	O_DEFER                       = 0x2000
	O_DELAY                       = 0x4000
	O_DIRECT                      = 0x8000000
	O_DIRECTORY                   = 0x80000
	O_DSYNC                       = 0x400000
	O_EFSOFF                      = 0x400000000
	O_EFSON                       = 0x200000000
	O_EXCL                        = 0x400
	O_EXEC                        = 0x20
	O_LARGEFILE                   = 0x4000000
	O_NDELAY                      = 0x8000
	O_NOCACHE                     = 0x100000
	O_NOCTTY                      = 0x800
	O_NOFOLLOW                    = 0x1000000
	O_NONBLOCK                    = 0x4
	O_NONE                        = 0x3
	O_NSHARE                      = 0x10000
	O_RAW                         = 0x100000000
	O_RDONLY                      = 0x0
	O_RDWR                        = 0x2
	O_RSHARE                      = 0x1000
	O_RSYNC                       = 0x200000
	O_SEARCH                      = 0x20
	O_SNAPSHOT                    = 0x40
	O_SYNC                        = 0x10
	O_TRUNC                       = 0x200
	O_TTY_INIT                    = 0x0
	O_WRONLY                      = 0x1
	PARENB                        = 0x100
	PAREXT                        = 0x100000
	PARMRK                        = 0x8
	PARODD                        = 0x200
	PENDIN                        = 0x20000000
	PRIO_PGRP                     = 0x1
	PRIO_PROCESS                  = 0x0
	PRIO_USER                     = 0x2
	PROT_EXEC                     = 0x4
	PROT_NONE                     = 0x0
	PROT_READ                     = 0x1
	PROT_WRITE                    = 0x2
	PR_64BIT                      = 0x20
	PR_ADDR                       = 0x2
	PR_ARGEXT                     = 0x400
	PR_ATOMIC                     = 0x1
	PR_CONNREQUIRED               = 0x4
	PR_FASTHZ                     = 0x5
	PR_INP                        = 0x40
	PR_INTRLEVEL                  = 0x8000
	PR_MLS                        = 0x100
	PR_MLS_1_LABEL                = 0x200
	PR_NOEOR                      = 0x4000
	PR_RIGHTS                     = 0x10
	PR_SLOWHZ                     = 0x2
	PR_WANTRCVD                   = 0x8
	RLIMIT_AS                     = 0x6
	RLIMIT_CORE                   = 0x4
	RLIMIT_CPU                    = 0x0
	RLIMIT_DATA                   = 0x2
	RLIMIT_FSIZE                  = 0x1
	RLIMIT_NOFILE                 = 0x7
	RLIMIT_NPROC                  = 0x9
	RLIMIT_RSS                    = 0x5
	RLIMIT_STACK                  = 0x3
	RLIM_INFINITY                 = 0x7fffffffffffffff
	RTAX_AUTHOR                   = 0x6
	RTAX_BRD                      = 0x7
	RTAX_DST                      = 0x0
	RTAX_GATEWAY                  = 0x1
	RTAX_GENMASK                  = 0x3
	RTAX_IFA                      = 0x5
	RTAX_IFP                      = 0x4
	RTAX_MAX                      = 0x8
	RTAX_NETMASK                  = 0x2
	RTA_AUTHOR                    = 0x40
	RTA_BRD                       = 0x80
	RTA_DOWNSTREAM                = 0x100
	RTA_DST                       = 0x1
	RTA_GATEWAY                   = 0x2
	RTA_GENMASK                   = 0x8
	RTA_IFA                       = 0x20
	RTA_IFP                       = 0x10
	RTA_NETMASK                   = 0x4
	RTC_IA64                      = 0x3
	RTC_POWER                     = 0x1
	RTC_POWER_PC                  = 0x2
	RTF_ACTIVE_DGD                = 0x1000000
	RTF_BCE                       = 0x80000
	RTF_BLACKHOLE                 = 0x1000
	RTF_BROADCAST                 = 0x400000
	RTF_BUL                       = 0x2000
	RTF_CLONE                     = 0x10000
	RTF_CLONED                    = 0x20000
	RTF_CLONING                   = 0x100
	RTF_DONE                      = 0x40
	RTF_DYNAMIC                   = 0x10
	RTF_FREE_IN_PROG              = 0x4000000
	RTF_GATEWAY                   = 0x2
	RTF_HOST                      = 0x4
	RTF_LLINFO                    = 0x400
	RTF_LOCAL                     = 0x200000
	RTF_MASK                      = 0x80
	RTF_MODIFIED                  = 0x20
	RTF_MULTICAST                 = 0x800000
	RTF_PERMANENT6                = 0x8000000
	RTF_PINNED                    = 0x100000
	RTF_PROTO1                    = 0x8000
	RTF_PROTO2                    = 0x4000
	RTF_PROTO3                    = 0x40000
	RTF_REJECT                    = 0x8
	RTF_SMALLMTU                  = 0x40000
	RTF_STATIC                    = 0x800
	RTF_STOPSRCH                  = 0x2000000
	RTF_UNREACHABLE               = 0x10000000
	RTF_UP                        = 0x1
	RTF_XRESOLVE                  = 0x200
	RTM_ADD                       = 0x1
	RTM_CHANGE                    = 0x3
	RTM_DELADDR                   = 0xd
	RTM_DELETE                    = 0x2
	RTM_EXPIRE                    = 0xf
	RTM_GET                       = 0x4
	RTM_GETNEXT                   = 0x11
	RTM_IFINFO                    = 0xe
	RTM_LOCK                      = 0x8
	RTM_LOSING                    = 0x5
	RTM_MISS                      = 0x7
	RTM_NEWADDR                   = 0xc
	RTM_OLDADD                    = 0x9
	RTM_OLDDEL                    = 0xa
	RTM_REDIRECT                  = 0x6
	RTM_RESOLVE                   = 0xb
	RTM_RTLOST                    = 0x10
	RTM_RTTUNIT                   = 0xf4240
	RTM_SAMEADDR                  = 0x12
	RTM_SET                       = 0x13
	RTM_VERSION                   = 0x2
	RTM_VERSION_GR                = 0x4
	RTM_VERSION_GR_COMPAT         = 0x3
	RTM_VERSION_POLICY            = 0x5
	RTM_VERSION_POLICY_EXT        = 0x6
	RTM_VERSION_POLICY_PRFN       = 0x7
	RTV_EXPIRE                    = 0x4
	RTV_HOPCOUNT                  = 0x2
	RTV_MTU                       = 0x1
	RTV_RPIPE                     = 0x8
	RTV_RTT                       = 0x40
	RTV_RTTVAR                    = 0x80
	RTV_SPIPE                     = 0x10
	RTV_SSTHRESH                  = 0x20
	RUSAGE_CHILDREN               = -0x1
	RUSAGE_SELF                   = 0x0
	RUSAGE_THREAD                 = 0x1
	SCM_RIGHTS                    = 0x1
	SHUT_RD                       = 0x0
	SHUT_RDWR                     = 0x2
	SHUT_WR                       = 0x1
	SIGMAX64                      = 0xff
	SIGQUEUE_MAX                  = 0x20
	SIOCADDIFVIPA                 = 0x20006942
	SIOCADDMTU                    = -0x7ffb9690
	SIOCADDMULTI                  = -0x7fdf96cf
	SIOCADDNETID                  = -0x7fd796a9
	SIOCADDRT                     = -0x7fc78df6
	SIOCAIFADDR                   = -0x7fbf96e6
	SIOCATMARK                    = 0x40047307
	SIOCDARP                      = -0x7fb396e0
	SIOCDELIFVIPA                 = 0x20006943
	SIOCDELMTU                    = -0x7ffb968f
	SIOCDELMULTI                  = -0x7fdf96ce
	SIOCDELPMTU                   = -0x7fd78ff6
	SIOCDELRT                     = -0x7fc78df5
	SIOCDIFADDR                   = -0x7fd796e7
	SIOCDNETOPT                   = -0x3ffe9680
	SIOCDX25XLATE                 = -0x7fd7969b
	SIOCFIFADDR                   = -0x7fdf966d
	SIOCGARP                      = -0x3fb396da
	SIOCGETMTUS                   = 0x2000696f
	SIOCGETSGCNT                  = -0x3feb8acc
	SIOCGETVIFCNT                 = -0x3feb8acd
	SIOCGHIWAT                    = 0x40047301
	SIOCGIFADDR                   = -0x3fd796df
	SIOCGIFADDRS                  = 0x2000698c
	SIOCGIFBAUDRATE               = -0x3fd79693
	SIOCGIFBRDADDR                = -0x3fd796dd
	SIOCGIFCONF                   = -0x3fef96bb
	SIOCGIFCONFGLOB               = -0x3fef9670
	SIOCGIFDSTADDR                = -0x3fd796de
	SIOCGIFFLAGS                  = -0x3fd796ef
	SIOCGIFGIDLIST                = 0x20006968
	SIOCGIFHWADDR                 = -0x3fab966b
	SIOCGIFMETRIC                 = -0x3fd796e9
	SIOCGIFMTU                    = -0x3fd796aa
	SIOCGIFNETMASK                = -0x3fd796db
	SIOCGIFOPTIONS                = -0x3fd796d6
	SIOCGISNO                     = -0x3fd79695
	SIOCGLOADF                    = -0x3ffb967e
	SIOCGLOWAT                    = 0x40047303
	SIOCGNETOPT                   = -0x3ffe96a5
	SIOCGNETOPT1                  = -0x3fdf967f
	SIOCGNMTUS                    = 0x2000696e
	SIOCGPGRP                     = 0x40047309
	SIOCGSIZIFCONF                = 0x4004696a
	SIOCGSRCFILTER                = -0x3fe796cb
	SIOCGTUNEPHASE                = -0x3ffb9676
	SIOCGX25XLATE                 = -0x3fd7969c
	SIOCIFATTACH                  = -0x7fdf9699
	SIOCIFDETACH                  = -0x7fdf969a
	SIOCIFGETPKEY                 = -0x7fdf969b
	SIOCIF_ATM_DARP               = -0x7fdf9683
	SIOCIF_ATM_DUMPARP            = -0x7fdf9685
	SIOCIF_ATM_GARP               = -0x7fdf9682
	SIOCIF_ATM_IDLE               = -0x7fdf9686
	SIOCIF_ATM_SARP               = -0x7fdf9681
	SIOCIF_ATM_SNMPARP            = -0x7fdf9687
	SIOCIF_ATM_SVC                = -0x7fdf9684
	SIOCIF_ATM_UBR                = -0x7fdf9688
	SIOCIF_DEVHEALTH              = -0x7ffb966c
	SIOCIF_IB_ARP_INCOMP          = -0x7fdf9677
	SIOCIF_IB_ARP_TIMER           = -0x7fdf9678
	SIOCIF_IB_CLEAR_PINFO         = -0x3fdf966f
	SIOCIF_IB_DEL_ARP             = -0x7fdf967f
	SIOCIF_IB_DEL_PINFO           = -0x3fdf9670
	SIOCIF_IB_DUMP_ARP            = -0x7fdf9680
	SIOCIF_IB_GET_ARP             = -0x7fdf967e
	SIOCIF_IB_GET_INFO            = -0x3f879675
	SIOCIF_IB_GET_STATS           = -0x3f879672
	SIOCIF_IB_NOTIFY_ADDR_REM     = -0x3f87966a
	SIOCIF_IB_RESET_STATS         = -0x3f879671
	SIOCIF_IB_RESIZE_CQ           = -0x7fdf9679
	SIOCIF_IB_SET_ARP             = -0x7fdf967d
	SIOCIF_IB_SET_PKEY            = -0x7fdf967c
	SIOCIF_IB_SET_PORT            = -0x7fdf967b
	SIOCIF_IB_SET_QKEY            = -0x7fdf9676
	SIOCIF_IB_SET_QSIZE           = -0x7fdf967a
	SIOCLISTIFVIPA                = 0x20006944
	SIOCSARP                      = -0x7fb396e2
	SIOCSHIWAT                    = 0xffffffff80047300
	SIOCSIFADDR                   = -0x7fd796f4
	SIOCSIFADDRORI                = -0x7fdb9673
	SIOCSIFBRDADDR                = -0x7fd796ed
	SIOCSIFDSTADDR                = -0x7fd796f2
	SIOCSIFFLAGS                  = -0x7fd796f0
	SIOCSIFGIDLIST                = 0x20006969
	SIOCSIFMETRIC                 = -0x7fd796e8
	SIOCSIFMTU                    = -0x7fd796a8
	SIOCSIFNETDUMP                = -0x7fd796e4
	SIOCSIFNETMASK                = -0x7fd796ea
	SIOCSIFOPTIONS                = -0x7fd796d7
	SIOCSIFSUBCHAN                = -0x7fd796e5
	SIOCSISNO                     = -0x7fd79694
	SIOCSLOADF                    = -0x3ffb967d
	SIOCSLOWAT                    = 0xffffffff80047302
	SIOCSNETOPT                   = -0x7ffe96a6
	SIOCSPGRP                     = 0xffffffff80047308
	SIOCSX25XLATE                 = -0x7fd7969d
	SOCK_CONN_DGRAM               = 0x6
	SOCK_DGRAM                    = 0x2
	SOCK_RAW                      = 0x3
	SOCK_RDM                      = 0x4
	SOCK_SEQPACKET                = 0x5
	SOCK_STREAM                   = 0x1
	SOL_SOCKET                    = 0xffff
	SOMAXCONN                     = 0x400
	SO_ACCEPTCONN                 = 0x2
	SO_AUDIT                      = 0x8000
	SO_BROADCAST                  = 0x20
	SO_CKSUMRECV                  = 0x800
	SO_DEBUG                      = 0x1
	SO_DONTROUTE                  = 0x10
	SO_ERROR                      = 0x1007
	SO_KEEPALIVE                  = 0x8
	SO_KERNACCEPT                 = 0x2000
	SO_LINGER                     = 0x80
	SO_NOMULTIPATH                = 0x4000
	SO_NOREUSEADDR                = 0x1000
	SO_OOBINLINE                  = 0x100
	SO_PEERID                     = 0x1009
	SO_RCVBUF                     = 0x1002
	SO_RCVLOWAT                   = 0x1004
	SO_RCVTIMEO                   = 0x1006
	SO_REUSEADDR                  = 0x4
	SO_REUSEPORT                  = 0x200
	SO_SNDBUF                     = 0x1001
	SO_SNDLOWAT                   = 0x1003
	SO_SNDTIMEO                   = 0x1005
	SO_TIMESTAMPNS                = 0x100a
	SO_TYPE                       = 0x1008
	SO_USELOOPBACK                = 0x40
	SO_USE_IFBUFS                 = 0x400
	S_BANDURG                     = 0x400
	S_EMODFMT                     = 0x3c000000
	S_ENFMT                       = 0x400
	S_ERROR                       = 0x100
	S_HANGUP                      = 0x200
	S_HIPRI                       = 0x2
	S_ICRYPTO                     = 0x80000
	S_IEXEC                       = 0x40
	S_IFBLK                       = 0x6000
	S_IFCHR                       = 0x2000
	S_IFDIR                       = 0x4000
	S_IFIFO                       = 0x1000
	S_IFJOURNAL                   = 0x10000
	S_IFLNK                       = 0xa000
	S_IFMPX                       = 0x2200
	S_IFMT                        = 0xf000
	S_IFPDIR                      = 0x4000000
	S_IFPSDIR                     = 0x8000000
	S_IFPSSDIR                    = 0xc000000
	S_IFREG                       = 0x8000
	S_IFSOCK                      = 0xc000
	S_IFSYSEA                     = 0x30000000
	S_INPUT                       = 0x1
	S_IREAD                       = 0x100
	S_IRGRP                       = 0x20
	S_IROTH                       = 0x4
	S_IRUSR                       = 0x100
	S_IRWXG                       = 0x38
	S_IRWXO                       = 0x7
	S_IRWXU                       = 0x1c0
	S_ISGID                       = 0x400
	S_ISUID                       = 0x800
	S_ISVTX                       = 0x200
	S_ITCB                        = 0x1000000
	S_ITP                         = 0x800000
	S_IWGRP                       = 0x10
	S_IWOTH                       = 0x2
	S_IWRITE                      = 0x80
	S_IWUSR                       = 0x80
	S_IXACL                       = 0x2000000
	S_IXATTR                      = 0x40000
	S_IXGRP                       = 0x8
	S_IXINTERFACE                 = 0x100000
	S_IXMOD                       = 0x40000000
	S_IXOTH                       = 0x1
	S_IXUSR                       = 0x40
	S_MSG                         = 0x8
	S_OUTPUT                      = 0x4
	S_RDBAND                      = 0x20
	S_RDNORM                      = 0x10
	S_RESERVED1                   = 0x20000
	S_RESERVED2                   = 0x200000
	S_RESERVED3                   = 0x400000
	S_RESERVED4                   = 0x80000000
	S_RESFMT1                     = 0x10000000
	S_RESFMT10                    = 0x34000000
	S_RESFMT11                    = 0x38000000
	S_RESFMT12                    = 0x3c000000
	S_RESFMT2                     = 0x14000000
	S_RESFMT3                     = 0x18000000
	S_RESFMT4                     = 0x1c000000
	S_RESFMT5                     = 0x20000000
	S_RESFMT6                     = 0x24000000
	S_RESFMT7                     = 0x28000000
	S_RESFMT8                     = 0x2c000000
	S_WRBAND                      = 0x80
	S_WRNORM                      = 0x40
	TAB0                          = 0x0
	TAB1                          = 0x400
	TAB2                          = 0x800
	TAB3                          = 0xc00
	TABDLY                        = 0xc00
	TCFLSH                        = 0x540c
	TCGETA                        = 0x5405
	TCGETS                        = 0x5401
	TCIFLUSH                      = 0x0
	TCIOFF                        = 0x2
	TCIOFLUSH                     = 0x2
	TCION                         = 0x3
	TCOFLUSH                      = 0x1
	TCOOFF                        = 0x0
	TCOON                         = 0x1
	TCP_24DAYS_WORTH_OF_SLOWTICKS = 0x3f4800
	TCP_ACLADD                    = 0x23
	TCP_ACLBIND                   = 0x26
	TCP_ACLCLEAR                  = 0x22
	TCP_ACLDEL                    = 0x24
	TCP_ACLDENY                   = 0x8
	TCP_ACLFLUSH                  = 0x21
	TCP_ACLGID                    = 0x1
	TCP_ACLLS                     = 0x25
	TCP_ACLSUBNET                 = 0x4
	TCP_ACLUID                    = 0x2
	TCP_CWND_DF                   = 0x16
	TCP_CWND_IF                   = 0x15
	TCP_DELAY_ACK_FIN             = 0x2
	TCP_DELAY_ACK_SYN             = 0x1
	TCP_FASTNAME                  = 0x101080a
	TCP_KEEPCNT                   = 0x13
	TCP_KEEPIDLE                  = 0x11
	TCP_KEEPINTVL                 = 0x12
	TCP_LSPRIV                    = 0x29
	TCP_LUID                      = 0x20
	TCP_MAXBURST                  = 0x8
	TCP_MAXDF                     = 0x64
	TCP_MAXIF                     = 0x64
	TCP_MAXSEG                    = 0x2
	TCP_MAXWIN                    = 0xffff
	TCP_MAXWINDOWSCALE            = 0xe
	TCP_MAX_SACK                  = 0x4
	TCP_MSS                       = 0x5b4
	TCP_NODELAY                   = 0x1
	TCP_NODELAYACK                = 0x14
	TCP_NOREDUCE_CWND_EXIT_FRXMT  = 0x19
	TCP_NOREDUCE_CWND_IN_FRXMT    = 0x18
	TCP_NOTENTER_SSTART           = 0x17
	TCP_OPT                       = 0x19
	TCP_RFC1323                   = 0x4
	TCP_SETPRIV                   = 0x27
	TCP_STDURG                    = 0x10
	TCP_TIMESTAMP_OPTLEN          = 0xc
	TCP_UNSETPRIV                 = 0x28
	TCSAFLUSH                     = 0x2
	TCSBRK                        = 0x5409
	TCSETA                        = 0x5406
	TCSETAF                       = 0x5408
	TCSETAW                       = 0x5407
	TCSETS                        = 0x5402
	TCSETSF                       = 0x5404
	TCSETSW                       = 0x5403
	TCXONC                        = 0x540b
	TIOC                          = 0x5400
	TIOCCBRK                      = 0x2000747a
	TIOCCDTR                      = 0x20007478
	TIOCCONS                      = 0xffffffff80047462
	TIOCEXCL                      = 0x2000740d
	TIOCFLUSH                     = 0xffffffff80047410
	TIOCGETC                      = 0x40067412
	TIOCGETD                      = 0x40047400
	TIOCGETP                      = 0x40067408
	TIOCGLTC                      = 0x40067474
	TIOCGPGRP                     = 0x40047477
	TIOCGSID                      = 0x40047448
	TIOCGSIZE                     = 0x40087468
	TIOCGWINSZ                    = 0x40087468
	TIOCHPCL                      = 0x20007402
	TIOCLBIC                      = 0xffffffff8004747e
	TIOCLBIS                      = 0xffffffff8004747f
	TIOCLGET                      = 0x4004747c
	TIOCLSET                      = 0xffffffff8004747d
	TIOCMBIC                      = 0xffffffff8004746b
	TIOCMBIS                      = 0xffffffff8004746c
	TIOCMGET                      = 0x4004746a
	TIOCMIWAIT                    = 0xffffffff80047464
	TIOCMODG                      = 0x40047403
	TIOCMODS                      = 0xffffffff80047404
	TIOCMSET                      = 0xffffffff8004746d
	TIOCM_CAR                     = 0x40
	TIOCM_CD                      = 0x40
	TIOCM_CTS                     = 0x20
	TIOCM_DSR                     = 0x100
	TIOCM_DTR                     = 0x2
	TIOCM_LE                      = 0x1
	TIOCM_RI                      = 0x80
	TIOCM_RNG                     = 0x80
	TIOCM_RTS                     = 0x4
	TIOCM_SR                      = 0x10
	TIOCM_ST                      = 0x8
	TIOCNOTTY                     = 0x20007471
	TIOCNXCL                      = 0x2000740e
	TIOCOUTQ                      = 0x40047473
	TIOCPKT                       = 0xffffffff80047470
	TIOCPKT_DATA                  = 0x0
	TIOCPKT_DOSTOP                = 0x20
	TIOCPKT_FLUSHREAD             = 0x1
	TIOCPKT_FLUSHWRITE            = 0x2
	TIOCPKT_NOSTOP                = 0x10
	TIOCPKT_START                 = 0x8
	TIOCPKT_STOP                  = 0x4
	TIOCREMOTE                    = 0xffffffff80047469
	TIOCSBRK                      = 0x2000747b
	TIOCSDTR                      = 0x20007479
	TIOCSETC                      = 0xffffffff80067411
	TIOCSETD                      = 0xffffffff80047401
	TIOCSETN                      = 0xffffffff8006740a
	TIOCSETP                      = 0xffffffff80067409
	TIOCSLTC                      = 0xffffffff80067475
	TIOCSPGRP                     = 0xffffffff80047476
	TIOCSSIZE                     = 0xffffffff80087467
	TIOCSTART                     = 0x2000746e
	TIOCSTI                       = 0xffffffff80017472
	TIOCSTOP                      = 0x2000746f
	TIOCSWINSZ                    = 0xffffffff80087467
	TIOCUCNTL                     = 0xffffffff80047466
	TOSTOP                        = 0x10000
	UTIME_NOW                     = -0x2
	UTIME_OMIT                    = -0x3
	VDISCRD                       = 0xc
	VDSUSP                        = 0xa
	VEOF                          = 0x4
	VEOL                          = 0x5
	VEOL2                         = 0x6
	VERASE                        = 0x2
	VINTR                         = 0x0
	VKILL                         = 0x3
	VLNEXT                        = 0xe
	VMIN                          = 0x4
	VQUIT                         = 0x1
	VREPRINT                      = 0xb
	VSTART                        = 0x7
	VSTOP                         = 0x8
	VSTRT                         = 0x7
	VSUSP                         = 0x9
	VT0                           = 0x0
	VT1                           = 0x8000
	VTDELAY                       = 0x2000
	VTDLY                         = 0x8000
	VTIME                         = 0x5
	VWERSE                        = 0xd
	WPARSTART                     = 0x1
	WPARSTOP                      = 0x2
	WPARTTYNAME                   = "Global"
	XCASE                         = 0x4
	XTABS                         = 0xc00
	_FDATAFLUSH                   = 0x2000000000
)


const (
	E2BIG           = syscall.Errno(0x7)
	EACCES          = syscall.Errno(0xd)
	EADDRINUSE      = syscall.Errno(0x43)
	EADDRNOTAVAIL   = syscall.Errno(0x44)
	EAFNOSUPPORT    = syscall.Errno(0x42)
	EAGAIN          = syscall.Errno(0xb)
	EALREADY        = syscall.Errno(0x38)
	EBADF           = syscall.Errno(0x9)
	EBADMSG         = syscall.Errno(0x78)
	EBUSY           = syscall.Errno(0x10)
	ECANCELED       = syscall.Errno(0x75)
	ECHILD          = syscall.Errno(0xa)
	ECHRNG          = syscall.Errno(0x25)
	ECLONEME        = syscall.Errno(0x52)
	ECONNABORTED    = syscall.Errno(0x48)
	ECONNREFUSED    = syscall.Errno(0x4f)
	ECONNRESET      = syscall.Errno(0x49)
	ECORRUPT        = syscall.Errno(0x59)
	EDEADLK         = syscall.Errno(0x2d)
	EDESTADDREQ     = syscall.Errno(0x3a)
	EDESTADDRREQ    = syscall.Errno(0x3a)
	EDIST           = syscall.Errno(0x35)
	EDOM            = syscall.Errno(0x21)
	EDQUOT          = syscall.Errno(0x58)
	EEXIST          = syscall.Errno(0x11)
	EFAULT          = syscall.Errno(0xe)
	EFBIG           = syscall.Errno(0x1b)
	EFORMAT         = syscall.Errno(0x30)
	EHOSTDOWN       = syscall.Errno(0x50)
	EHOSTUNREACH    = syscall.Errno(0x51)
	EIDRM           = syscall.Errno(0x24)
	EILSEQ          = syscall.Errno(0x74)
	EINPROGRESS     = syscall.Errno(0x37)
	EINTR           = syscall.Errno(0x4)
	EINVAL          = syscall.Errno(0x16)
	EIO             = syscall.Errno(0x5)
	EISCONN         = syscall.Errno(0x4b)
	EISDIR          = syscall.Errno(0x15)
	EL2HLT          = syscall.Errno(0x2c)
	EL2NSYNC        = syscall.Errno(0x26)
	EL3HLT          = syscall.Errno(0x27)
	EL3RST          = syscall.Errno(0x28)
	ELNRNG          = syscall.Errno(0x29)
	ELOOP           = syscall.Errno(0x55)
	EMEDIA          = syscall.Errno(0x6e)
	EMFILE          = syscall.Errno(0x18)
	EMLINK          = syscall.Errno(0x1f)
	EMSGSIZE        = syscall.Errno(0x3b)
	EMULTIHOP       = syscall.Errno(0x7d)
	ENAMETOOLONG    = syscall.Errno(0x56)
	ENETDOWN        = syscall.Errno(0x45)
	ENETRESET       = syscall.Errno(0x47)
	ENETUNREACH     = syscall.Errno(0x46)
	ENFILE          = syscall.Errno(0x17)
	ENOATTR         = syscall.Errno(0x70)
	ENOBUFS         = syscall.Errno(0x4a)
	ENOCONNECT      = syscall.Errno(0x32)
	ENOCSI          = syscall.Errno(0x2b)
	ENODATA         = syscall.Errno(0x7a)
	ENODEV          = syscall.Errno(0x13)
	ENOENT          = syscall.Errno(0x2)
	ENOEXEC         = syscall.Errno(0x8)
	ENOLCK          = syscall.Errno(0x31)
	ENOLINK         = syscall.Errno(0x7e)
	ENOMEM          = syscall.Errno(0xc)
	ENOMSG          = syscall.Errno(0x23)
	ENOPROTOOPT     = syscall.Errno(0x3d)
	ENOSPC          = syscall.Errno(0x1c)
	ENOSR           = syscall.Errno(0x76)
	ENOSTR          = syscall.Errno(0x7b)
	ENOSYS          = syscall.Errno(0x6d)
	ENOTBLK         = syscall.Errno(0xf)
	ENOTCONN        = syscall.Errno(0x4c)
	ENOTDIR         = syscall.Errno(0x14)
	ENOTEMPTY       = syscall.Errno(0x11)
	ENOTREADY       = syscall.Errno(0x2e)
	ENOTRECOVERABLE = syscall.Errno(0x5e)
	ENOTRUST        = syscall.Errno(0x72)
	ENOTSOCK        = syscall.Errno(0x39)
	ENOTSUP         = syscall.Errno(0x7c)
	ENOTTY          = syscall.Errno(0x19)
	ENXIO           = syscall.Errno(0x6)
	EOPNOTSUPP      = syscall.Errno(0x40)
	EOVERFLOW       = syscall.Errno(0x7f)
	EOWNERDEAD      = syscall.Errno(0x5f)
	EPERM           = syscall.Errno(0x1)
	EPFNOSUPPORT    = syscall.Errno(0x41)
	EPIPE           = syscall.Errno(0x20)
	EPROCLIM        = syscall.Errno(0x53)
	EPROTO          = syscall.Errno(0x79)
	EPROTONOSUPPORT = syscall.Errno(0x3e)
	EPROTOTYPE      = syscall.Errno(0x3c)
	ERANGE          = syscall.Errno(0x22)
	EREMOTE         = syscall.Errno(0x5d)
	ERESTART        = syscall.Errno(0x52)
	EROFS           = syscall.Errno(0x1e)
	ESAD            = syscall.Errno(0x71)
	ESHUTDOWN       = syscall.Errno(0x4d)
	ESOCKTNOSUPPORT = syscall.Errno(0x3f)
	ESOFT           = syscall.Errno(0x6f)
	ESPIPE          = syscall.Errno(0x1d)
	ESRCH           = syscall.Errno(0x3)
	ESTALE          = syscall.Errno(0x34)
	ESYSERROR       = syscall.Errno(0x5a)
	ETIME           = syscall.Errno(0x77)
	ETIMEDOUT       = syscall.Errno(0x4e)
	ETOOMANYREFS    = syscall.Errno(0x73)
	ETXTBSY         = syscall.Errno(0x1a)
	EUNATCH         = syscall.Errno(0x2a)
	EUSERS          = syscall.Errno(0x54)
	EWOULDBLOCK     = syscall.Errno(0xb)
	EWRPROTECT      = syscall.Errno(0x2f)
	EXDEV           = syscall.Errno(0x12)
)


const (
	SIGABRT     = syscall.Signal(0x6)
	SIGAIO      = syscall.Signal(0x17)
	SIGALRM     = syscall.Signal(0xe)
	SIGALRM1    = syscall.Signal(0x26)
	SIGBUS      = syscall.Signal(0xa)
	SIGCAPI     = syscall.Signal(0x31)
	SIGCHLD     = syscall.Signal(0x14)
	SIGCLD      = syscall.Signal(0x14)
	SIGCONT     = syscall.Signal(0x13)
	SIGCPUFAIL  = syscall.Signal(0x3b)
	SIGDANGER   = syscall.Signal(0x21)
	SIGEMT      = syscall.Signal(0x7)
	SIGFPE      = syscall.Signal(0x8)
	SIGGRANT    = syscall.Signal(0x3c)
	SIGHUP      = syscall.Signal(0x1)
	SIGILL      = syscall.Signal(0x4)
	SIGINT      = syscall.Signal(0x2)
	SIGIO       = syscall.Signal(0x17)
	SIGIOINT    = syscall.Signal(0x10)
	SIGIOT      = syscall.Signal(0x6)
	SIGKAP      = syscall.Signal(0x3c)
	SIGKILL     = syscall.Signal(0x9)
	SIGLOST     = syscall.Signal(0x6)
	SIGMAX      = syscall.Signal(0xff)
	SIGMAX32    = syscall.Signal(0x3f)
	SIGMIGRATE  = syscall.Signal(0x23)
	SIGMSG      = syscall.Signal(0x1b)
	SIGPIPE     = syscall.Signal(0xd)
	SIGPOLL     = syscall.Signal(0x17)
	SIGPRE      = syscall.Signal(0x24)
	SIGPROF     = syscall.Signal(0x20)
	SIGPTY      = syscall.Signal(0x17)
	SIGPWR      = syscall.Signal(0x1d)
	SIGQUIT     = syscall.Signal(0x3)
	SIGRECONFIG = syscall.Signal(0x3a)
	SIGRETRACT  = syscall.Signal(0x3d)
	SIGSAK      = syscall.Signal(0x3f)
	SIGSEGV     = syscall.Signal(0xb)
	SIGSOUND    = syscall.Signal(0x3e)
	SIGSTOP     = syscall.Signal(0x11)
	SIGSYS      = syscall.Signal(0xc)
	SIGSYSERROR = syscall.Signal(0x30)
	SIGTALRM    = syscall.Signal(0x26)
	SIGTERM     = syscall.Signal(0xf)
	SIGTRAP     = syscall.Signal(0x5)
	SIGTSTP     = syscall.Signal(0x12)
	SIGTTIN     = syscall.Signal(0x15)
	SIGTTOU     = syscall.Signal(0x16)
	SIGURG      = syscall.Signal(0x10)
	SIGUSR1     = syscall.Signal(0x1e)
	SIGUSR2     = syscall.Signal(0x1f)
	SIGVIRT     = syscall.Signal(0x25)
	SIGVTALRM   = syscall.Signal(0x22)
	SIGWAITING  = syscall.Signal(0x27)
	SIGWINCH    = syscall.Signal(0x1c)
	SIGXCPU     = syscall.Signal(0x18)
	SIGXFSZ     = syscall.Signal(0x19)
)


var errorList = [...]struct {
	num  syscall.Errno
	name string
	desc string
}{
	{1, "EPERM", "not owner"},
	{2, "ENOENT", "no such file or directory"},
	{3, "ESRCH", "no such process"},
	{4, "EINTR", "interrupted system call"},
	{5, "EIO", "I/O error"},
	{6, "ENXIO", "no such device or address"},
	{7, "E2BIG", "arg list too long"},
	{8, "ENOEXEC", "exec format error"},
	{9, "EBADF", "bad file number"},
	{10, "ECHILD", "no child processes"},
	{11, "EWOULDBLOCK", "resource temporarily unavailable"},
	{12, "ENOMEM", "not enough space"},
	{13, "EACCES", "permission denied"},
	{14, "EFAULT", "bad address"},
	{15, "ENOTBLK", "block device required"},
	{16, "EBUSY", "device busy"},
	{17, "ENOTEMPTY", "file exists"},
	{18, "EXDEV", "cross-device link"},
	{19, "ENODEV", "no such device"},
	{20, "ENOTDIR", "not a directory"},
	{21, "EISDIR", "is a directory"},
	{22, "EINVAL", "invalid argument"},
	{23, "ENFILE", "file table overflow"},
	{24, "EMFILE", "too many open files"},
	{25, "ENOTTY", "not a typewriter"},
	{26, "ETXTBSY", "text file busy"},
	{27, "EFBIG", "file too large"},
	{28, "ENOSPC", "no space left on device"},
	{29, "ESPIPE", "illegal seek"},
	{30, "EROFS", "read-only file system"},
	{31, "EMLINK", "too many links"},
	{32, "EPIPE", "broken pipe"},
	{33, "EDOM", "argument out of domain"},
	{34, "ERANGE", "result too large"},
	{35, "ENOMSG", "no message of desired type"},
	{36, "EIDRM", "identifier removed"},
	{37, "ECHRNG", "channel number out of range"},
	{38, "EL2NSYNC", "level 2 not synchronized"},
	{39, "EL3HLT", "level 3 halted"},
	{40, "EL3RST", "level 3 reset"},
	{41, "ELNRNG", "link number out of range"},
	{42, "EUNATCH", "protocol driver not attached"},
	{43, "ENOCSI", "no CSI structure available"},
	{44, "EL2HLT", "level 2 halted"},
	{45, "EDEADLK", "deadlock condition if locked"},
	{46, "ENOTREADY", "device not ready"},
	{47, "EWRPROTECT", "write-protected media"},
	{48, "EFORMAT", "unformatted or incompatible media"},
	{49, "ENOLCK", "no locks available"},
	{50, "ENOCONNECT", "cannot Establish Connection"},
	{52, "ESTALE", "missing file or filesystem"},
	{53, "EDIST", "requests blocked by Administrator"},
	{55, "EINPROGRESS", "operation now in progress"},
	{56, "EALREADY", "operation already in progress"},
	{57, "ENOTSOCK", "socket operation on non-socket"},
	{58, "EDESTADDREQ", "destination address required"},
	{59, "EMSGSIZE", "message too long"},
	{60, "EPROTOTYPE", "protocol wrong type for socket"},
	{61, "ENOPROTOOPT", "protocol not available"},
	{62, "EPROTONOSUPPORT", "protocol not supported"},
	{63, "ESOCKTNOSUPPORT", "socket type not supported"},
	{64, "EOPNOTSUPP", "operation not supported on socket"},
	{65, "EPFNOSUPPORT", "protocol family not supported"},
	{66, "EAFNOSUPPORT", "addr family not supported by protocol"},
	{67, "EADDRINUSE", "address already in use"},
	{68, "EADDRNOTAVAIL", "can't assign requested address"},
	{69, "ENETDOWN", "network is down"},
	{70, "ENETUNREACH", "network is unreachable"},
	{71, "ENETRESET", "network dropped connection on reset"},
	{72, "ECONNABORTED", "software caused connection abort"},
	{73, "ECONNRESET", "connection reset by peer"},
	{74, "ENOBUFS", "no buffer space available"},
	{75, "EISCONN", "socket is already connected"},
	{76, "ENOTCONN", "socket is not connected"},
	{77, "ESHUTDOWN", "can't send after socket shutdown"},
	{78, "ETIMEDOUT", "connection timed out"},
	{79, "ECONNREFUSED", "connection refused"},
	{80, "EHOSTDOWN", "host is down"},
	{81, "EHOSTUNREACH", "no route to host"},
	{82, "ERESTART", "restart the system call"},
	{83, "EPROCLIM", "too many processes"},
	{84, "EUSERS", "too many users"},
	{85, "ELOOP", "too many levels of symbolic links"},
	{86, "ENAMETOOLONG", "file name too long"},
	{88, "EDQUOT", "disk quota exceeded"},
	{89, "ECORRUPT", "invalid file system control data detected"},
	{90, "ESYSERROR", "for future use "},
	{93, "EREMOTE", "item is not local to host"},
	{94, "ENOTRECOVERABLE", "state not recoverable "},
	{95, "EOWNERDEAD", "previous owner died "},
	{109, "ENOSYS", "function not implemented"},
	{110, "EMEDIA", "media surface error"},
	{111, "ESOFT", "I/O completed, but needs relocation"},
	{112, "ENOATTR", "no attribute found"},
	{113, "ESAD", "security Authentication Denied"},
	{114, "ENOTRUST", "not a Trusted Program"},
	{115, "ETOOMANYREFS", "too many references: can't splice"},
	{116, "EILSEQ", "invalid wide character"},
	{117, "ECANCELED", "asynchronous I/O cancelled"},
	{118, "ENOSR", "out of STREAMS resources"},
	{119, "ETIME", "system call timed out"},
	{120, "EBADMSG", "next message has wrong type"},
	{121, "EPROTO", "error in protocol"},
	{122, "ENODATA", "no message on stream head read q"},
	{123, "ENOSTR", "fd not associated with a stream"},
	{124, "ENOTSUP", "unsupported attribute value"},
	{125, "EMULTIHOP", "multihop is not allowed"},
	{126, "ENOLINK", "the server link has been severed"},
	{127, "EOVERFLOW", "value too large to be stored in data type"},
}


var signalList = [...]struct {
	num  syscall.Signal
	name string
	desc string
}{
	{1, "SIGHUP", "hangup"},
	{2, "SIGINT", "interrupt"},
	{3, "SIGQUIT", "quit"},
	{4, "SIGILL", "illegal instruction"},
	{5, "SIGTRAP", "trace/BPT trap"},
	{6, "SIGIOT", "IOT/Abort trap"},
	{7, "SIGEMT", "EMT trap"},
	{8, "SIGFPE", "floating point exception"},
	{9, "SIGKILL", "killed"},
	{10, "SIGBUS", "bus error"},
	{11, "SIGSEGV", "segmentation fault"},
	{12, "SIGSYS", "bad system call"},
	{13, "SIGPIPE", "broken pipe"},
	{14, "SIGALRM", "alarm clock"},
	{15, "SIGTERM", "terminated"},
	{16, "SIGURG", "urgent I/O condition"},
	{17, "SIGSTOP", "stopped (signal)"},
	{18, "SIGTSTP", "stopped"},
	{19, "SIGCONT", "continued"},
	{20, "SIGCHLD", "child exited"},
	{21, "SIGTTIN", "stopped (tty input)"},
	{22, "SIGTTOU", "stopped (tty output)"},
	{23, "SIGIO", "I/O possible/complete"},
	{24, "SIGXCPU", "cputime limit exceeded"},
	{25, "SIGXFSZ", "filesize limit exceeded"},
	{27, "SIGMSG", "input device data"},
	{28, "SIGWINCH", "window size changes"},
	{29, "SIGPWR", "power-failure"},
	{30, "SIGUSR1", "user defined signal 1"},
	{31, "SIGUSR2", "user defined signal 2"},
	{32, "SIGPROF", "profiling timer expired"},
	{33, "SIGDANGER", "paging space low"},
	{34, "SIGVTALRM", "virtual timer expired"},
	{35, "SIGMIGRATE", "signal 35"},
	{36, "SIGPRE", "signal 36"},
	{37, "SIGVIRT", "signal 37"},
	{38, "SIGTALRM", "signal 38"},
	{39, "SIGWAITING", "signal 39"},
	{48, "SIGSYSERROR", "signal 48"},
	{49, "SIGCAPI", "signal 49"},
	{58, "SIGRECONFIG", "signal 58"},
	{59, "SIGCPUFAIL", "CPU Failure Predicted"},
	{60, "SIGGRANT", "monitor mode granted"},
	{61, "SIGRETRACT", "monitor mode retracted"},
	{62, "SIGSOUND", "sound completed"},
	{63, "SIGMAX32", "secure attention"},
	{255, "SIGMAX", "signal 255"},
}
