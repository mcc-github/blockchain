







package unix

import "syscall"

const (
	AF_802                        = 0x12
	AF_APPLETALK                  = 0x10
	AF_CCITT                      = 0xa
	AF_CHAOS                      = 0x5
	AF_DATAKIT                    = 0x9
	AF_DECnet                     = 0xc
	AF_DLI                        = 0xd
	AF_ECMA                       = 0x8
	AF_FILE                       = 0x1
	AF_GOSIP                      = 0x16
	AF_HYLINK                     = 0xf
	AF_IMPLINK                    = 0x3
	AF_INET                       = 0x2
	AF_INET6                      = 0x1a
	AF_INET_OFFLOAD               = 0x1e
	AF_IPX                        = 0x17
	AF_KEY                        = 0x1b
	AF_LAT                        = 0xe
	AF_LINK                       = 0x19
	AF_LOCAL                      = 0x1
	AF_MAX                        = 0x20
	AF_NBS                        = 0x7
	AF_NCA                        = 0x1c
	AF_NIT                        = 0x11
	AF_NS                         = 0x6
	AF_OSI                        = 0x13
	AF_OSINET                     = 0x15
	AF_PACKET                     = 0x20
	AF_POLICY                     = 0x1d
	AF_PUP                        = 0x4
	AF_ROUTE                      = 0x18
	AF_SNA                        = 0xb
	AF_TRILL                      = 0x1f
	AF_UNIX                       = 0x1
	AF_UNSPEC                     = 0x0
	AF_X25                        = 0x14
	ARPHRD_ARCNET                 = 0x7
	ARPHRD_ATM                    = 0x10
	ARPHRD_AX25                   = 0x3
	ARPHRD_CHAOS                  = 0x5
	ARPHRD_EETHER                 = 0x2
	ARPHRD_ETHER                  = 0x1
	ARPHRD_FC                     = 0x12
	ARPHRD_FRAME                  = 0xf
	ARPHRD_HDLC                   = 0x11
	ARPHRD_IB                     = 0x20
	ARPHRD_IEEE802                = 0x6
	ARPHRD_IPATM                  = 0x13
	ARPHRD_METRICOM               = 0x17
	ARPHRD_TUNNEL                 = 0x1f
	B0                            = 0x0
	B110                          = 0x3
	B115200                       = 0x12
	B1200                         = 0x9
	B134                          = 0x4
	B150                          = 0x5
	B153600                       = 0x13
	B1800                         = 0xa
	B19200                        = 0xe
	B200                          = 0x6
	B230400                       = 0x14
	B2400                         = 0xb
	B300                          = 0x7
	B307200                       = 0x15
	B38400                        = 0xf
	B460800                       = 0x16
	B4800                         = 0xc
	B50                           = 0x1
	B57600                        = 0x10
	B600                          = 0x8
	B75                           = 0x2
	B76800                        = 0x11
	B921600                       = 0x17
	B9600                         = 0xd
	BIOCFLUSH                     = 0x20004268
	BIOCGBLEN                     = 0x40044266
	BIOCGDLT                      = 0x4004426a
	BIOCGDLTLIST                  = -0x3fefbd89
	BIOCGDLTLIST32                = -0x3ff7bd89
	BIOCGETIF                     = 0x4020426b
	BIOCGETLIF                    = 0x4078426b
	BIOCGHDRCMPLT                 = 0x40044274
	BIOCGRTIMEOUT                 = 0x4010427b
	BIOCGRTIMEOUT32               = 0x4008427b
	BIOCGSEESENT                  = 0x40044278
	BIOCGSTATS                    = 0x4080426f
	BIOCGSTATSOLD                 = 0x4008426f
	BIOCIMMEDIATE                 = -0x7ffbbd90
	BIOCPROMISC                   = 0x20004269
	BIOCSBLEN                     = -0x3ffbbd9a
	BIOCSDLT                      = -0x7ffbbd8a
	BIOCSETF                      = -0x7fefbd99
	BIOCSETF32                    = -0x7ff7bd99
	BIOCSETIF                     = -0x7fdfbd94
	BIOCSETLIF                    = -0x7f87bd94
	BIOCSHDRCMPLT                 = -0x7ffbbd8b
	BIOCSRTIMEOUT                 = -0x7fefbd86
	BIOCSRTIMEOUT32               = -0x7ff7bd86
	BIOCSSEESENT                  = -0x7ffbbd87
	BIOCSTCPF                     = -0x7fefbd8e
	BIOCSUDPF                     = -0x7fefbd8d
	BIOCVERSION                   = 0x40044271
	BPF_A                         = 0x10
	BPF_ABS                       = 0x20
	BPF_ADD                       = 0x0
	BPF_ALIGNMENT                 = 0x4
	BPF_ALU                       = 0x4
	BPF_AND                       = 0x50
	BPF_B                         = 0x10
	BPF_DFLTBUFSIZE               = 0x100000
	BPF_DIV                       = 0x30
	BPF_H                         = 0x8
	BPF_IMM                       = 0x0
	BPF_IND                       = 0x40
	BPF_JA                        = 0x0
	BPF_JEQ                       = 0x10
	BPF_JGE                       = 0x30
	BPF_JGT                       = 0x20
	BPF_JMP                       = 0x5
	BPF_JSET                      = 0x40
	BPF_K                         = 0x0
	BPF_LD                        = 0x0
	BPF_LDX                       = 0x1
	BPF_LEN                       = 0x80
	BPF_LSH                       = 0x60
	BPF_MAJOR_VERSION             = 0x1
	BPF_MAXBUFSIZE                = 0x1000000
	BPF_MAXINSNS                  = 0x200
	BPF_MEM                       = 0x60
	BPF_MEMWORDS                  = 0x10
	BPF_MINBUFSIZE                = 0x20
	BPF_MINOR_VERSION             = 0x1
	BPF_MISC                      = 0x7
	BPF_MSH                       = 0xa0
	BPF_MUL                       = 0x20
	BPF_NEG                       = 0x80
	BPF_OR                        = 0x40
	BPF_RELEASE                   = 0x30bb6
	BPF_RET                       = 0x6
	BPF_RSH                       = 0x70
	BPF_ST                        = 0x2
	BPF_STX                       = 0x3
	BPF_SUB                       = 0x10
	BPF_TAX                       = 0x0
	BPF_TXA                       = 0x80
	BPF_W                         = 0x0
	BPF_X                         = 0x8
	BRKINT                        = 0x2
	BS0                           = 0x0
	BS1                           = 0x2000
	BSDLY                         = 0x2000
	CBAUD                         = 0xf
	CFLUSH                        = 0xf
	CIBAUD                        = 0xf0000
	CLOCAL                        = 0x800
	CLOCK_HIGHRES                 = 0x4
	CLOCK_LEVEL                   = 0xa
	CLOCK_MONOTONIC               = 0x4
	CLOCK_PROCESS_CPUTIME_ID      = 0x5
	CLOCK_PROF                    = 0x2
	CLOCK_REALTIME                = 0x3
	CLOCK_THREAD_CPUTIME_ID       = 0x2
	CLOCK_VIRTUAL                 = 0x1
	CR0                           = 0x0
	CR1                           = 0x200
	CR2                           = 0x400
	CR3                           = 0x600
	CRDLY                         = 0x600
	CREAD                         = 0x80
	CRTSCTS                       = 0x80000000
	CS5                           = 0x0
	CS6                           = 0x10
	CS7                           = 0x20
	CS8                           = 0x30
	CSIZE                         = 0x30
	CSTART                        = 0x11
	CSTATUS                       = 0x14
	CSTOP                         = 0x13
	CSTOPB                        = 0x40
	CSUSP                         = 0x1a
	CSWTCH                        = 0x1a
	DLT_AIRONET_HEADER            = 0x78
	DLT_APPLE_IP_OVER_IEEE1394    = 0x8a
	DLT_ARCNET                    = 0x7
	DLT_ARCNET_LINUX              = 0x81
	DLT_ATM_CLIP                  = 0x13
	DLT_ATM_RFC1483               = 0xb
	DLT_AURORA                    = 0x7e
	DLT_AX25                      = 0x3
	DLT_BACNET_MS_TP              = 0xa5
	DLT_CHAOS                     = 0x5
	DLT_CISCO_IOS                 = 0x76
	DLT_C_HDLC                    = 0x68
	DLT_DOCSIS                    = 0x8f
	DLT_ECONET                    = 0x73
	DLT_EN10MB                    = 0x1
	DLT_EN3MB                     = 0x2
	DLT_ENC                       = 0x6d
	DLT_ERF_ETH                   = 0xaf
	DLT_ERF_POS                   = 0xb0
	DLT_FDDI                      = 0xa
	DLT_FRELAY                    = 0x6b
	DLT_GCOM_SERIAL               = 0xad
	DLT_GCOM_T1E1                 = 0xac
	DLT_GPF_F                     = 0xab
	DLT_GPF_T                     = 0xaa
	DLT_GPRS_LLC                  = 0xa9
	DLT_HDLC                      = 0x10
	DLT_HHDLC                     = 0x79
	DLT_HIPPI                     = 0xf
	DLT_IBM_SN                    = 0x92
	DLT_IBM_SP                    = 0x91
	DLT_IEEE802                   = 0x6
	DLT_IEEE802_11                = 0x69
	DLT_IEEE802_11_RADIO          = 0x7f
	DLT_IEEE802_11_RADIO_AVS      = 0xa3
	DLT_IPNET                     = 0xe2
	DLT_IPOIB                     = 0xa2
	DLT_IP_OVER_FC                = 0x7a
	DLT_JUNIPER_ATM1              = 0x89
	DLT_JUNIPER_ATM2              = 0x87
	DLT_JUNIPER_CHDLC             = 0xb5
	DLT_JUNIPER_ES                = 0x84
	DLT_JUNIPER_ETHER             = 0xb2
	DLT_JUNIPER_FRELAY            = 0xb4
	DLT_JUNIPER_GGSN              = 0x85
	DLT_JUNIPER_MFR               = 0x86
	DLT_JUNIPER_MLFR              = 0x83
	DLT_JUNIPER_MLPPP             = 0x82
	DLT_JUNIPER_MONITOR           = 0xa4
	DLT_JUNIPER_PIC_PEER          = 0xae
	DLT_JUNIPER_PPP               = 0xb3
	DLT_JUNIPER_PPPOE             = 0xa7
	DLT_JUNIPER_PPPOE_ATM         = 0xa8
	DLT_JUNIPER_SERVICES          = 0x88
	DLT_LINUX_IRDA                = 0x90
	DLT_LINUX_LAPD                = 0xb1
	DLT_LINUX_SLL                 = 0x71
	DLT_LOOP                      = 0x6c
	DLT_LTALK                     = 0x72
	DLT_MTP2                      = 0x8c
	DLT_MTP2_WITH_PHDR            = 0x8b
	DLT_MTP3                      = 0x8d
	DLT_NULL                      = 0x0
	DLT_PCI_EXP                   = 0x7d
	DLT_PFLOG                     = 0x75
	DLT_PFSYNC                    = 0x12
	DLT_PPP                       = 0x9
	DLT_PPP_BSDOS                 = 0xe
	DLT_PPP_PPPD                  = 0xa6
	DLT_PRISM_HEADER              = 0x77
	DLT_PRONET                    = 0x4
	DLT_RAW                       = 0xc
	DLT_RAWAF_MASK                = 0x2240000
	DLT_RIO                       = 0x7c
	DLT_SCCP                      = 0x8e
	DLT_SLIP                      = 0x8
	DLT_SLIP_BSDOS                = 0xd
	DLT_SUNATM                    = 0x7b
	DLT_SYMANTEC_FIREWALL         = 0x63
	DLT_TZSP                      = 0x80
	ECHO                          = 0x8
	ECHOCTL                       = 0x200
	ECHOE                         = 0x10
	ECHOK                         = 0x20
	ECHOKE                        = 0x800
	ECHONL                        = 0x40
	ECHOPRT                       = 0x400
	EMPTY_SET                     = 0x0
	EMT_CPCOVF                    = 0x1
	EQUALITY_CHECK                = 0x0
	EXTA                          = 0xe
	EXTB                          = 0xf
	FD_CLOEXEC                    = 0x1
	FD_NFDBITS                    = 0x40
	FD_SETSIZE                    = 0x10000
	FF0                           = 0x0
	FF1                           = 0x8000
	FFDLY                         = 0x8000
	FLUSHALL                      = 0x1
	FLUSHDATA                     = 0x0
	FLUSHO                        = 0x2000
	F_ALLOCSP                     = 0xa
	F_ALLOCSP64                   = 0xa
	F_BADFD                       = 0x2e
	F_BLKSIZE                     = 0x13
	F_BLOCKS                      = 0x12
	F_CHKFL                       = 0x8
	F_COMPAT                      = 0x8
	F_DUP2FD                      = 0x9
	F_DUP2FD_CLOEXEC              = 0x24
	F_DUPFD                       = 0x0
	F_DUPFD_CLOEXEC               = 0x25
	F_FLOCK                       = 0x35
	F_FLOCK64                     = 0x35
	F_FLOCKW                      = 0x36
	F_FLOCKW64                    = 0x36
	F_FREESP                      = 0xb
	F_FREESP64                    = 0xb
	F_GETFD                       = 0x1
	F_GETFL                       = 0x3
	F_GETLK                       = 0xe
	F_GETLK64                     = 0xe
	F_GETOWN                      = 0x17
	F_GETXFL                      = 0x2d
	F_HASREMOTELOCKS              = 0x1a
	F_ISSTREAM                    = 0xd
	F_MANDDNY                     = 0x10
	F_MDACC                       = 0x20
	F_NODNY                       = 0x0
	F_NPRIV                       = 0x10
	F_OFD_GETLK                   = 0x2f
	F_OFD_GETLK64                 = 0x2f
	F_OFD_SETLK                   = 0x30
	F_OFD_SETLK64                 = 0x30
	F_OFD_SETLKW                  = 0x31
	F_OFD_SETLKW64                = 0x31
	F_PRIV                        = 0xf
	F_QUOTACTL                    = 0x11
	F_RDACC                       = 0x1
	F_RDDNY                       = 0x1
	F_RDLCK                       = 0x1
	F_REVOKE                      = 0x19
	F_RMACC                       = 0x4
	F_RMDNY                       = 0x4
	F_RWACC                       = 0x3
	F_RWDNY                       = 0x3
	F_SETFD                       = 0x2
	F_SETFL                       = 0x4
	F_SETLK                       = 0x6
	F_SETLK64                     = 0x6
	F_SETLK64_NBMAND              = 0x2a
	F_SETLKW                      = 0x7
	F_SETLKW64                    = 0x7
	F_SETLK_NBMAND                = 0x2a
	F_SETOWN                      = 0x18
	F_SHARE                       = 0x28
	F_SHARE_NBMAND                = 0x2b
	F_UNLCK                       = 0x3
	F_UNLKSYS                     = 0x4
	F_UNSHARE                     = 0x29
	F_WRACC                       = 0x2
	F_WRDNY                       = 0x2
	F_WRLCK                       = 0x2
	HUPCL                         = 0x400
	IBSHIFT                       = 0x10
	ICANON                        = 0x2
	ICRNL                         = 0x100
	IEXTEN                        = 0x8000
	IFF_ADDRCONF                  = 0x80000
	IFF_ALLMULTI                  = 0x200
	IFF_ANYCAST                   = 0x400000
	IFF_BROADCAST                 = 0x2
	IFF_CANTCHANGE                = 0x7f203003b5a
	IFF_COS_ENABLED               = 0x200000000
	IFF_DEBUG                     = 0x4
	IFF_DEPRECATED                = 0x40000
	IFF_DHCPRUNNING               = 0x4000
	IFF_DUPLICATE                 = 0x4000000000
	IFF_FAILED                    = 0x10000000
	IFF_FIXEDMTU                  = 0x1000000000
	IFF_INACTIVE                  = 0x40000000
	IFF_INTELLIGENT               = 0x400
	IFF_IPMP                      = 0x8000000000
	IFF_IPMP_CANTCHANGE           = 0x10000000
	IFF_IPMP_INVALID              = 0x1ec200080
	IFF_IPV4                      = 0x1000000
	IFF_IPV6                      = 0x2000000
	IFF_L3PROTECT                 = 0x40000000000
	IFF_LOOPBACK                  = 0x8
	IFF_MULTICAST                 = 0x800
	IFF_MULTI_BCAST               = 0x1000
	IFF_NOACCEPT                  = 0x4000000
	IFF_NOARP                     = 0x80
	IFF_NOFAILOVER                = 0x8000000
	IFF_NOLINKLOCAL               = 0x20000000000
	IFF_NOLOCAL                   = 0x20000
	IFF_NONUD                     = 0x200000
	IFF_NORTEXCH                  = 0x800000
	IFF_NOTRAILERS                = 0x20
	IFF_NOXMIT                    = 0x10000
	IFF_OFFLINE                   = 0x80000000
	IFF_POINTOPOINT               = 0x10
	IFF_PREFERRED                 = 0x400000000
	IFF_PRIVATE                   = 0x8000
	IFF_PROMISC                   = 0x100
	IFF_ROUTER                    = 0x100000
	IFF_RUNNING                   = 0x40
	IFF_STANDBY                   = 0x20000000
	IFF_TEMPORARY                 = 0x800000000
	IFF_UNNUMBERED                = 0x2000
	IFF_UP                        = 0x1
	IFF_VIRTUAL                   = 0x2000000000
	IFF_VRRP                      = 0x10000000000
	IFF_XRESOLV                   = 0x100000000
	IFNAMSIZ                      = 0x10
	IFT_1822                      = 0x2
	IFT_6TO4                      = 0xca
	IFT_AAL5                      = 0x31
	IFT_ARCNET                    = 0x23
	IFT_ARCNETPLUS                = 0x24
	IFT_ATM                       = 0x25
	IFT_CEPT                      = 0x13
	IFT_DS3                       = 0x1e
	IFT_EON                       = 0x19
	IFT_ETHER                     = 0x6
	IFT_FDDI                      = 0xf
	IFT_FRELAY                    = 0x20
	IFT_FRELAYDCE                 = 0x2c
	IFT_HDH1822                   = 0x3
	IFT_HIPPI                     = 0x2f
	IFT_HSSI                      = 0x2e
	IFT_HY                        = 0xe
	IFT_IB                        = 0xc7
	IFT_IPV4                      = 0xc8
	IFT_IPV6                      = 0xc9
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
	IFT_SONET                     = 0x27
	IFT_SONETPATH                 = 0x32
	IFT_SONETVT                   = 0x33
	IFT_STARLAN                   = 0xb
	IFT_T1                        = 0x12
	IFT_ULTRA                     = 0x1d
	IFT_V35                       = 0x2d
	IFT_X25                       = 0x5
	IFT_X25DDN                    = 0x4
	IFT_X25PLE                    = 0x28
	IFT_XETHER                    = 0x1a
	IGNBRK                        = 0x1
	IGNCR                         = 0x80
	IGNPAR                        = 0x4
	IMAXBEL                       = 0x2000
	INLCR                         = 0x40
	INPCK                         = 0x10
	IN_AUTOCONF_MASK              = 0xffff0000
	IN_AUTOCONF_NET               = 0xa9fe0000
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
	IN_CLASSE_NET                 = 0xffffffff
	IN_LOOPBACKNET                = 0x7f
	IN_PRIVATE12_MASK             = 0xfff00000
	IN_PRIVATE12_NET              = 0xac100000
	IN_PRIVATE16_MASK             = 0xffff0000
	IN_PRIVATE16_NET              = 0xc0a80000
	IN_PRIVATE8_MASK              = 0xff000000
	IN_PRIVATE8_NET               = 0xa000000
	IPPROTO_AH                    = 0x33
	IPPROTO_DSTOPTS               = 0x3c
	IPPROTO_EGP                   = 0x8
	IPPROTO_ENCAP                 = 0x4
	IPPROTO_EON                   = 0x50
	IPPROTO_ESP                   = 0x32
	IPPROTO_FRAGMENT              = 0x2c
	IPPROTO_GGP                   = 0x3
	IPPROTO_HELLO                 = 0x3f
	IPPROTO_HOPOPTS               = 0x0
	IPPROTO_ICMP                  = 0x1
	IPPROTO_ICMPV6                = 0x3a
	IPPROTO_IDP                   = 0x16
	IPPROTO_IGMP                  = 0x2
	IPPROTO_IP                    = 0x0
	IPPROTO_IPV6                  = 0x29
	IPPROTO_MAX                   = 0x100
	IPPROTO_ND                    = 0x4d
	IPPROTO_NONE                  = 0x3b
	IPPROTO_OSPF                  = 0x59
	IPPROTO_PIM                   = 0x67
	IPPROTO_PUP                   = 0xc
	IPPROTO_RAW                   = 0xff
	IPPROTO_ROUTING               = 0x2b
	IPPROTO_RSVP                  = 0x2e
	IPPROTO_SCTP                  = 0x84
	IPPROTO_TCP                   = 0x6
	IPPROTO_UDP                   = 0x11
	IPV6_ADD_MEMBERSHIP           = 0x9
	IPV6_BOUND_IF                 = 0x41
	IPV6_CHECKSUM                 = 0x18
	IPV6_DONTFRAG                 = 0x21
	IPV6_DROP_MEMBERSHIP          = 0xa
	IPV6_DSTOPTS                  = 0xf
	IPV6_FLOWINFO_FLOWLABEL       = 0xffff0f00
	IPV6_FLOWINFO_TCLASS          = 0xf00f
	IPV6_HOPLIMIT                 = 0xc
	IPV6_HOPOPTS                  = 0xe
	IPV6_JOIN_GROUP               = 0x9
	IPV6_LEAVE_GROUP              = 0xa
	IPV6_MULTICAST_HOPS           = 0x7
	IPV6_MULTICAST_IF             = 0x6
	IPV6_MULTICAST_LOOP           = 0x8
	IPV6_NEXTHOP                  = 0xd
	IPV6_PAD1_OPT                 = 0x0
	IPV6_PATHMTU                  = 0x25
	IPV6_PKTINFO                  = 0xb
	IPV6_PREFER_SRC_CGA           = 0x20
	IPV6_PREFER_SRC_CGADEFAULT    = 0x10
	IPV6_PREFER_SRC_CGAMASK       = 0x30
	IPV6_PREFER_SRC_COA           = 0x2
	IPV6_PREFER_SRC_DEFAULT       = 0x15
	IPV6_PREFER_SRC_HOME          = 0x1
	IPV6_PREFER_SRC_MASK          = 0x3f
	IPV6_PREFER_SRC_MIPDEFAULT    = 0x1
	IPV6_PREFER_SRC_MIPMASK       = 0x3
	IPV6_PREFER_SRC_NONCGA        = 0x10
	IPV6_PREFER_SRC_PUBLIC        = 0x4
	IPV6_PREFER_SRC_TMP           = 0x8
	IPV6_PREFER_SRC_TMPDEFAULT    = 0x4
	IPV6_PREFER_SRC_TMPMASK       = 0xc
	IPV6_RECVDSTOPTS              = 0x28
	IPV6_RECVHOPLIMIT             = 0x13
	IPV6_RECVHOPOPTS              = 0x14
	IPV6_RECVPATHMTU              = 0x24
	IPV6_RECVPKTINFO              = 0x12
	IPV6_RECVRTHDR                = 0x16
	IPV6_RECVRTHDRDSTOPTS         = 0x17
	IPV6_RECVTCLASS               = 0x19
	IPV6_RTHDR                    = 0x10
	IPV6_RTHDRDSTOPTS             = 0x11
	IPV6_RTHDR_TYPE_0             = 0x0
	IPV6_SEC_OPT                  = 0x22
	IPV6_SRC_PREFERENCES          = 0x23
	IPV6_TCLASS                   = 0x26
	IPV6_UNICAST_HOPS             = 0x5
	IPV6_UNSPEC_SRC               = 0x42
	IPV6_USE_MIN_MTU              = 0x20
	IPV6_V6ONLY                   = 0x27
	IP_ADD_MEMBERSHIP             = 0x13
	IP_ADD_SOURCE_MEMBERSHIP      = 0x17
	IP_BLOCK_SOURCE               = 0x15
	IP_BOUND_IF                   = 0x41
	IP_BROADCAST                  = 0x106
	IP_BROADCAST_TTL              = 0x43
	IP_DEFAULT_MULTICAST_LOOP     = 0x1
	IP_DEFAULT_MULTICAST_TTL      = 0x1
	IP_DF                         = 0x4000
	IP_DHCPINIT_IF                = 0x45
	IP_DONTFRAG                   = 0x1b
	IP_DONTROUTE                  = 0x105
	IP_DROP_MEMBERSHIP            = 0x14
	IP_DROP_SOURCE_MEMBERSHIP     = 0x18
	IP_HDRINCL                    = 0x2
	IP_MAXPACKET                  = 0xffff
	IP_MF                         = 0x2000
	IP_MSS                        = 0x240
	IP_MULTICAST_IF               = 0x10
	IP_MULTICAST_LOOP             = 0x12
	IP_MULTICAST_TTL              = 0x11
	IP_NEXTHOP                    = 0x19
	IP_OPTIONS                    = 0x1
	IP_PKTINFO                    = 0x1a
	IP_RECVDSTADDR                = 0x7
	IP_RECVIF                     = 0x9
	IP_RECVOPTS                   = 0x5
	IP_RECVPKTINFO                = 0x1a
	IP_RECVRETOPTS                = 0x6
	IP_RECVSLLA                   = 0xa
	IP_RECVTTL                    = 0xb
	IP_RETOPTS                    = 0x8
	IP_REUSEADDR                  = 0x104
	IP_SEC_OPT                    = 0x22
	IP_TOS                        = 0x3
	IP_TTL                        = 0x4
	IP_UNBLOCK_SOURCE             = 0x16
	IP_UNSPEC_SRC                 = 0x42
	ISIG                          = 0x1
	ISTRIP                        = 0x20
	IUCLC                         = 0x200
	IXANY                         = 0x800
	IXOFF                         = 0x1000
	IXON                          = 0x400
	LOCK_EX                       = 0x2
	LOCK_NB                       = 0x4
	LOCK_SH                       = 0x1
	LOCK_UN                       = 0x8
	MADV_ACCESS_DEFAULT           = 0x6
	MADV_ACCESS_LWP               = 0x7
	MADV_ACCESS_MANY              = 0x8
	MADV_DONTNEED                 = 0x4
	MADV_FREE                     = 0x5
	MADV_NORMAL                   = 0x0
	MADV_PURGE                    = 0x9
	MADV_RANDOM                   = 0x1
	MADV_SEQUENTIAL               = 0x2
	MADV_WILLNEED                 = 0x3
	MAP_32BIT                     = 0x80
	MAP_ALIGN                     = 0x200
	MAP_ANON                      = 0x100
	MAP_ANONYMOUS                 = 0x100
	MAP_FILE                      = 0x0
	MAP_FIXED                     = 0x10
	MAP_INITDATA                  = 0x800
	MAP_NORESERVE                 = 0x40
	MAP_PRIVATE                   = 0x2
	MAP_RENAME                    = 0x20
	MAP_SHARED                    = 0x1
	MAP_TEXT                      = 0x400
	MAP_TYPE                      = 0xf
	MCL_CURRENT                   = 0x1
	MCL_FUTURE                    = 0x2
	MSG_CTRUNC                    = 0x10
	MSG_DONTROUTE                 = 0x4
	MSG_DONTWAIT                  = 0x80
	MSG_DUPCTRL                   = 0x800
	MSG_EOR                       = 0x8
	MSG_MAXIOVLEN                 = 0x10
	MSG_NOTIFICATION              = 0x100
	MSG_OOB                       = 0x1
	MSG_PEEK                      = 0x2
	MSG_TRUNC                     = 0x20
	MSG_WAITALL                   = 0x40
	MSG_XPG4_2                    = 0x8000
	MS_ASYNC                      = 0x1
	MS_INVALIDATE                 = 0x2
	MS_OLDSYNC                    = 0x0
	MS_SYNC                       = 0x4
	M_FLUSH                       = 0x86
	NAME_MAX                      = 0xff
	NEWDEV                        = 0x1
	NL0                           = 0x0
	NL1                           = 0x100
	NLDLY                         = 0x100
	NOFLSH                        = 0x80
	OCRNL                         = 0x8
	OFDEL                         = 0x80
	OFILL                         = 0x40
	OLCUC                         = 0x2
	OLDDEV                        = 0x0
	ONBITSMAJOR                   = 0x7
	ONBITSMINOR                   = 0x8
	ONLCR                         = 0x4
	ONLRET                        = 0x20
	ONOCR                         = 0x10
	OPENFAIL                      = -0x1
	OPOST                         = 0x1
	O_ACCMODE                     = 0x600003
	O_APPEND                      = 0x8
	O_CLOEXEC                     = 0x800000
	O_CREAT                       = 0x100
	O_DSYNC                       = 0x40
	O_EXCL                        = 0x400
	O_EXEC                        = 0x400000
	O_LARGEFILE                   = 0x2000
	O_NDELAY                      = 0x4
	O_NOCTTY                      = 0x800
	O_NOFOLLOW                    = 0x20000
	O_NOLINKS                     = 0x40000
	O_NONBLOCK                    = 0x80
	O_RDONLY                      = 0x0
	O_RDWR                        = 0x2
	O_RSYNC                       = 0x8000
	O_SEARCH                      = 0x200000
	O_SIOCGIFCONF                 = -0x3ff796ec
	O_SIOCGLIFCONF                = -0x3fef9688
	O_SYNC                        = 0x10
	O_TRUNC                       = 0x200
	O_WRONLY                      = 0x1
	O_XATTR                       = 0x4000
	PARENB                        = 0x100
	PAREXT                        = 0x100000
	PARMRK                        = 0x8
	PARODD                        = 0x200
	PENDIN                        = 0x4000
	PRIO_PGRP                     = 0x1
	PRIO_PROCESS                  = 0x0
	PRIO_USER                     = 0x2
	PROT_EXEC                     = 0x4
	PROT_NONE                     = 0x0
	PROT_READ                     = 0x1
	PROT_WRITE                    = 0x2
	RLIMIT_AS                     = 0x6
	RLIMIT_CORE                   = 0x4
	RLIMIT_CPU                    = 0x0
	RLIMIT_DATA                   = 0x2
	RLIMIT_FSIZE                  = 0x1
	RLIMIT_NOFILE                 = 0x5
	RLIMIT_STACK                  = 0x3
	RLIM_INFINITY                 = -0x3
	RTAX_AUTHOR                   = 0x6
	RTAX_BRD                      = 0x7
	RTAX_DST                      = 0x0
	RTAX_GATEWAY                  = 0x1
	RTAX_GENMASK                  = 0x3
	RTAX_IFA                      = 0x5
	RTAX_IFP                      = 0x4
	RTAX_MAX                      = 0x9
	RTAX_NETMASK                  = 0x2
	RTAX_SRC                      = 0x8
	RTA_AUTHOR                    = 0x40
	RTA_BRD                       = 0x80
	RTA_DST                       = 0x1
	RTA_GATEWAY                   = 0x2
	RTA_GENMASK                   = 0x8
	RTA_IFA                       = 0x20
	RTA_IFP                       = 0x10
	RTA_NETMASK                   = 0x4
	RTA_NUMBITS                   = 0x9
	RTA_SRC                       = 0x100
	RTF_BLACKHOLE                 = 0x1000
	RTF_CLONING                   = 0x100
	RTF_DONE                      = 0x40
	RTF_DYNAMIC                   = 0x10
	RTF_GATEWAY                   = 0x2
	RTF_HOST                      = 0x4
	RTF_INDIRECT                  = 0x40000
	RTF_KERNEL                    = 0x80000
	RTF_LLINFO                    = 0x400
	RTF_MASK                      = 0x80
	RTF_MODIFIED                  = 0x20
	RTF_MULTIRT                   = 0x10000
	RTF_PRIVATE                   = 0x2000
	RTF_PROTO1                    = 0x8000
	RTF_PROTO2                    = 0x4000
	RTF_REJECT                    = 0x8
	RTF_SETSRC                    = 0x20000
	RTF_STATIC                    = 0x800
	RTF_UP                        = 0x1
	RTF_XRESOLVE                  = 0x200
	RTF_ZONE                      = 0x100000
	RTM_ADD                       = 0x1
	RTM_CHANGE                    = 0x3
	RTM_CHGADDR                   = 0xf
	RTM_DELADDR                   = 0xd
	RTM_DELETE                    = 0x2
	RTM_FREEADDR                  = 0x10
	RTM_GET                       = 0x4
	RTM_IFINFO                    = 0xe
	RTM_LOCK                      = 0x8
	RTM_LOSING                    = 0x5
	RTM_MISS                      = 0x7
	RTM_NEWADDR                   = 0xc
	RTM_OLDADD                    = 0x9
	RTM_OLDDEL                    = 0xa
	RTM_REDIRECT                  = 0x6
	RTM_RESOLVE                   = 0xb
	RTM_VERSION                   = 0x3
	RTV_EXPIRE                    = 0x4
	RTV_HOPCOUNT                  = 0x2
	RTV_MTU                       = 0x1
	RTV_RPIPE                     = 0x8
	RTV_RTT                       = 0x40
	RTV_RTTVAR                    = 0x80
	RTV_SPIPE                     = 0x10
	RTV_SSTHRESH                  = 0x20
	RT_AWARE                      = 0x1
	RUSAGE_CHILDREN               = -0x1
	RUSAGE_SELF                   = 0x0
	SCM_RIGHTS                    = 0x1010
	SCM_TIMESTAMP                 = 0x1013
	SCM_UCRED                     = 0x1012
	SHUT_RD                       = 0x0
	SHUT_RDWR                     = 0x2
	SHUT_WR                       = 0x1
	SIG2STR_MAX                   = 0x20
	SIOCADDMULTI                  = -0x7fdf96cf
	SIOCADDRT                     = -0x7fcf8df6
	SIOCATMARK                    = 0x40047307
	SIOCDARP                      = -0x7fdb96e0
	SIOCDELMULTI                  = -0x7fdf96ce
	SIOCDELRT                     = -0x7fcf8df5
	SIOCDXARP                     = -0x7fff9658
	SIOCGARP                      = -0x3fdb96e1
	SIOCGDSTINFO                  = -0x3fff965c
	SIOCGENADDR                   = -0x3fdf96ab
	SIOCGENPSTATS                 = -0x3fdf96c7
	SIOCGETLSGCNT                 = -0x3fef8deb
	SIOCGETNAME                   = 0x40107334
	SIOCGETPEER                   = 0x40107335
	SIOCGETPROP                   = -0x3fff8f44
	SIOCGETSGCNT                  = -0x3feb8deb
	SIOCGETSYNC                   = -0x3fdf96d3
	SIOCGETVIFCNT                 = -0x3feb8dec
	SIOCGHIWAT                    = 0x40047301
	SIOCGIFADDR                   = -0x3fdf96f3
	SIOCGIFBRDADDR                = -0x3fdf96e9
	SIOCGIFCONF                   = -0x3ff796a4
	SIOCGIFDSTADDR                = -0x3fdf96f1
	SIOCGIFFLAGS                  = -0x3fdf96ef
	SIOCGIFHWADDR                 = -0x3fdf9647
	SIOCGIFINDEX                  = -0x3fdf96a6
	SIOCGIFMEM                    = -0x3fdf96ed
	SIOCGIFMETRIC                 = -0x3fdf96e5
	SIOCGIFMTU                    = -0x3fdf96ea
	SIOCGIFMUXID                  = -0x3fdf96a8
	SIOCGIFNETMASK                = -0x3fdf96e7
	SIOCGIFNUM                    = 0x40046957
	SIOCGIP6ADDRPOLICY            = -0x3fff965e
	SIOCGIPMSFILTER               = -0x3ffb964c
	SIOCGLIFADDR                  = -0x3f87968f
	SIOCGLIFBINDING               = -0x3f879666
	SIOCGLIFBRDADDR               = -0x3f879685
	SIOCGLIFCONF                  = -0x3fef965b
	SIOCGLIFDADSTATE              = -0x3f879642
	SIOCGLIFDSTADDR               = -0x3f87968d
	SIOCGLIFFLAGS                 = -0x3f87968b
	SIOCGLIFGROUPINFO             = -0x3f4b9663
	SIOCGLIFGROUPNAME             = -0x3f879664
	SIOCGLIFHWADDR                = -0x3f879640
	SIOCGLIFINDEX                 = -0x3f87967b
	SIOCGLIFLNKINFO               = -0x3f879674
	SIOCGLIFMETRIC                = -0x3f879681
	SIOCGLIFMTU                   = -0x3f879686
	SIOCGLIFMUXID                 = -0x3f87967d
	SIOCGLIFNETMASK               = -0x3f879683
	SIOCGLIFNUM                   = -0x3ff3967e
	SIOCGLIFSRCOF                 = -0x3fef964f
	SIOCGLIFSUBNET                = -0x3f879676
	SIOCGLIFTOKEN                 = -0x3f879678
	SIOCGLIFUSESRC                = -0x3f879651
	SIOCGLIFZONE                  = -0x3f879656
	SIOCGLOWAT                    = 0x40047303
	SIOCGMSFILTER                 = -0x3ffb964e
	SIOCGPGRP                     = 0x40047309
	SIOCGSTAMP                    = -0x3fef9646
	SIOCGXARP                     = -0x3fff9659
	SIOCIFDETACH                  = -0x7fdf96c8
	SIOCILB                       = -0x3ffb9645
	SIOCLIFADDIF                  = -0x3f879691
	SIOCLIFDELND                  = -0x7f879673
	SIOCLIFGETND                  = -0x3f879672
	SIOCLIFREMOVEIF               = -0x7f879692
	SIOCLIFSETND                  = -0x7f879671
	SIOCLOWER                     = -0x7fdf96d7
	SIOCSARP                      = -0x7fdb96e2
	SIOCSCTPGOPT                  = -0x3fef9653
	SIOCSCTPPEELOFF               = -0x3ffb9652
	SIOCSCTPSOPT                  = -0x7fef9654
	SIOCSENABLESDP                = -0x3ffb9649
	SIOCSETPROP                   = -0x7ffb8f43
	SIOCSETSYNC                   = -0x7fdf96d4
	SIOCSHIWAT                    = -0x7ffb8d00
	SIOCSIFADDR                   = -0x7fdf96f4
	SIOCSIFBRDADDR                = -0x7fdf96e8
	SIOCSIFDSTADDR                = -0x7fdf96f2
	SIOCSIFFLAGS                  = -0x7fdf96f0
	SIOCSIFINDEX                  = -0x7fdf96a5
	SIOCSIFMEM                    = -0x7fdf96ee
	SIOCSIFMETRIC                 = -0x7fdf96e4
	SIOCSIFMTU                    = -0x7fdf96eb
	SIOCSIFMUXID                  = -0x7fdf96a7
	SIOCSIFNAME                   = -0x7fdf96b7
	SIOCSIFNETMASK                = -0x7fdf96e6
	SIOCSIP6ADDRPOLICY            = -0x7fff965d
	SIOCSIPMSFILTER               = -0x7ffb964b
	SIOCSLGETREQ                  = -0x3fdf96b9
	SIOCSLIFADDR                  = -0x7f879690
	SIOCSLIFBRDADDR               = -0x7f879684
	SIOCSLIFDSTADDR               = -0x7f87968e
	SIOCSLIFFLAGS                 = -0x7f87968c
	SIOCSLIFGROUPNAME             = -0x7f879665
	SIOCSLIFINDEX                 = -0x7f87967a
	SIOCSLIFLNKINFO               = -0x7f879675
	SIOCSLIFMETRIC                = -0x7f879680
	SIOCSLIFMTU                   = -0x7f879687
	SIOCSLIFMUXID                 = -0x7f87967c
	SIOCSLIFNAME                  = -0x3f87967f
	SIOCSLIFNETMASK               = -0x7f879682
	SIOCSLIFPREFIX                = -0x3f879641
	SIOCSLIFSUBNET                = -0x7f879677
	SIOCSLIFTOKEN                 = -0x7f879679
	SIOCSLIFUSESRC                = -0x7f879650
	SIOCSLIFZONE                  = -0x7f879655
	SIOCSLOWAT                    = -0x7ffb8cfe
	SIOCSLSTAT                    = -0x7fdf96b8
	SIOCSMSFILTER                 = -0x7ffb964d
	SIOCSPGRP                     = -0x7ffb8cf8
	SIOCSPROMISC                  = -0x7ffb96d0
	SIOCSQPTR                     = -0x3ffb9648
	SIOCSSDSTATS                  = -0x3fdf96d2
	SIOCSSESTATS                  = -0x3fdf96d1
	SIOCSXARP                     = -0x7fff965a
	SIOCTMYADDR                   = -0x3ff79670
	SIOCTMYSITE                   = -0x3ff7966e
	SIOCTONLINK                   = -0x3ff7966f
	SIOCUPPER                     = -0x7fdf96d8
	SIOCX25RCV                    = -0x3fdf96c4
	SIOCX25TBL                    = -0x3fdf96c3
	SIOCX25XMT                    = -0x3fdf96c5
	SIOCXPROTO                    = 0x20007337
	SOCK_CLOEXEC                  = 0x80000
	SOCK_DGRAM                    = 0x1
	SOCK_NDELAY                   = 0x200000
	SOCK_NONBLOCK                 = 0x100000
	SOCK_RAW                      = 0x4
	SOCK_RDM                      = 0x5
	SOCK_SEQPACKET                = 0x6
	SOCK_STREAM                   = 0x2
	SOCK_TYPE_MASK                = 0xffff
	SOL_FILTER                    = 0xfffc
	SOL_PACKET                    = 0xfffd
	SOL_ROUTE                     = 0xfffe
	SOL_SOCKET                    = 0xffff
	SOMAXCONN                     = 0x80
	SO_ACCEPTCONN                 = 0x2
	SO_ALL                        = 0x3f
	SO_ALLZONES                   = 0x1014
	SO_ANON_MLP                   = 0x100a
	SO_ATTACH_FILTER              = 0x40000001
	SO_BAND                       = 0x4000
	SO_BROADCAST                  = 0x20
	SO_COPYOPT                    = 0x80000
	SO_DEBUG                      = 0x1
	SO_DELIM                      = 0x8000
	SO_DETACH_FILTER              = 0x40000002
	SO_DGRAM_ERRIND               = 0x200
	SO_DOMAIN                     = 0x100c
	SO_DONTLINGER                 = -0x81
	SO_DONTROUTE                  = 0x10
	SO_ERROPT                     = 0x40000
	SO_ERROR                      = 0x1007
	SO_EXCLBIND                   = 0x1015
	SO_HIWAT                      = 0x10
	SO_ISNTTY                     = 0x800
	SO_ISTTY                      = 0x400
	SO_KEEPALIVE                  = 0x8
	SO_LINGER                     = 0x80
	SO_LOWAT                      = 0x20
	SO_MAC_EXEMPT                 = 0x100b
	SO_MAC_IMPLICIT               = 0x1016
	SO_MAXBLK                     = 0x100000
	SO_MAXPSZ                     = 0x8
	SO_MINPSZ                     = 0x4
	SO_MREADOFF                   = 0x80
	SO_MREADON                    = 0x40
	SO_NDELOFF                    = 0x200
	SO_NDELON                     = 0x100
	SO_NODELIM                    = 0x10000
	SO_OOBINLINE                  = 0x100
	SO_PROTOTYPE                  = 0x1009
	SO_RCVBUF                     = 0x1002
	SO_RCVLOWAT                   = 0x1004
	SO_RCVPSH                     = 0x100d
	SO_RCVTIMEO                   = 0x1006
	SO_READOPT                    = 0x1
	SO_RECVUCRED                  = 0x400
	SO_REUSEADDR                  = 0x4
	SO_SECATTR                    = 0x1011
	SO_SNDBUF                     = 0x1001
	SO_SNDLOWAT                   = 0x1003
	SO_SNDTIMEO                   = 0x1005
	SO_STRHOLD                    = 0x20000
	SO_TAIL                       = 0x200000
	SO_TIMESTAMP                  = 0x1013
	SO_TONSTOP                    = 0x2000
	SO_TOSTOP                     = 0x1000
	SO_TYPE                       = 0x1008
	SO_USELOOPBACK                = 0x40
	SO_VRRP                       = 0x1017
	SO_WROFF                      = 0x2
	TAB0                          = 0x0
	TAB1                          = 0x800
	TAB2                          = 0x1000
	TAB3                          = 0x1800
	TABDLY                        = 0x1800
	TCFLSH                        = 0x5407
	TCGETA                        = 0x5401
	TCGETS                        = 0x540d
	TCIFLUSH                      = 0x0
	TCIOFF                        = 0x2
	TCIOFLUSH                     = 0x2
	TCION                         = 0x3
	TCOFLUSH                      = 0x1
	TCOOFF                        = 0x0
	TCOON                         = 0x1
	TCP_ABORT_THRESHOLD           = 0x11
	TCP_ANONPRIVBIND              = 0x20
	TCP_CONN_ABORT_THRESHOLD      = 0x13
	TCP_CONN_NOTIFY_THRESHOLD     = 0x12
	TCP_CORK                      = 0x18
	TCP_EXCLBIND                  = 0x21
	TCP_INIT_CWND                 = 0x15
	TCP_KEEPALIVE                 = 0x8
	TCP_KEEPALIVE_ABORT_THRESHOLD = 0x17
	TCP_KEEPALIVE_THRESHOLD       = 0x16
	TCP_KEEPCNT                   = 0x23
	TCP_KEEPIDLE                  = 0x22
	TCP_KEEPINTVL                 = 0x24
	TCP_LINGER2                   = 0x1c
	TCP_MAXSEG                    = 0x2
	TCP_MSS                       = 0x218
	TCP_NODELAY                   = 0x1
	TCP_NOTIFY_THRESHOLD          = 0x10
	TCP_RECVDSTADDR               = 0x14
	TCP_RTO_INITIAL               = 0x19
	TCP_RTO_MAX                   = 0x1b
	TCP_RTO_MIN                   = 0x1a
	TCSAFLUSH                     = 0x5410
	TCSBRK                        = 0x5405
	TCSETA                        = 0x5402
	TCSETAF                       = 0x5404
	TCSETAW                       = 0x5403
	TCSETS                        = 0x540e
	TCSETSF                       = 0x5410
	TCSETSW                       = 0x540f
	TCXONC                        = 0x5406
	TIOC                          = 0x5400
	TIOCCBRK                      = 0x747a
	TIOCCDTR                      = 0x7478
	TIOCCILOOP                    = 0x746c
	TIOCEXCL                      = 0x740d
	TIOCFLUSH                     = 0x7410
	TIOCGETC                      = 0x7412
	TIOCGETD                      = 0x7400
	TIOCGETP                      = 0x7408
	TIOCGLTC                      = 0x7474
	TIOCGPGRP                     = 0x7414
	TIOCGPPS                      = 0x547d
	TIOCGPPSEV                    = 0x547f
	TIOCGSID                      = 0x7416
	TIOCGSOFTCAR                  = 0x5469
	TIOCGWINSZ                    = 0x5468
	TIOCHPCL                      = 0x7402
	TIOCKBOF                      = 0x5409
	TIOCKBON                      = 0x5408
	TIOCLBIC                      = 0x747e
	TIOCLBIS                      = 0x747f
	TIOCLGET                      = 0x747c
	TIOCLSET                      = 0x747d
	TIOCMBIC                      = 0x741c
	TIOCMBIS                      = 0x741b
	TIOCMGET                      = 0x741d
	TIOCMSET                      = 0x741a
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
	TIOCNOTTY                     = 0x7471
	TIOCNXCL                      = 0x740e
	TIOCOUTQ                      = 0x7473
	TIOCREMOTE                    = 0x741e
	TIOCSBRK                      = 0x747b
	TIOCSCTTY                     = 0x7484
	TIOCSDTR                      = 0x7479
	TIOCSETC                      = 0x7411
	TIOCSETD                      = 0x7401
	TIOCSETN                      = 0x740a
	TIOCSETP                      = 0x7409
	TIOCSIGNAL                    = 0x741f
	TIOCSILOOP                    = 0x746d
	TIOCSLTC                      = 0x7475
	TIOCSPGRP                     = 0x7415
	TIOCSPPS                      = 0x547e
	TIOCSSOFTCAR                  = 0x546a
	TIOCSTART                     = 0x746e
	TIOCSTI                       = 0x7417
	TIOCSTOP                      = 0x746f
	TIOCSWINSZ                    = 0x5467
	TOSTOP                        = 0x100
	VCEOF                         = 0x8
	VCEOL                         = 0x9
	VDISCARD                      = 0xd
	VDSUSP                        = 0xb
	VEOF                          = 0x4
	VEOL                          = 0x5
	VEOL2                         = 0x6
	VERASE                        = 0x2
	VERASE2                       = 0x11
	VINTR                         = 0x0
	VKILL                         = 0x3
	VLNEXT                        = 0xf
	VMIN                          = 0x4
	VQUIT                         = 0x1
	VREPRINT                      = 0xc
	VSTART                        = 0x8
	VSTATUS                       = 0x10
	VSTOP                         = 0x9
	VSUSP                         = 0xa
	VSWTCH                        = 0x7
	VT0                           = 0x0
	VT1                           = 0x4000
	VTDLY                         = 0x4000
	VTIME                         = 0x5
	VWERASE                       = 0xe
	WCONTFLG                      = 0xffff
	WCONTINUED                    = 0x8
	WCOREFLG                      = 0x80
	WEXITED                       = 0x1
	WNOHANG                       = 0x40
	WNOWAIT                       = 0x80
	WOPTMASK                      = 0xcf
	WRAP                          = 0x20000
	WSIGMASK                      = 0x7f
	WSTOPFLG                      = 0x7f
	WSTOPPED                      = 0x4
	WTRAPPED                      = 0x2
	WUNTRACED                     = 0x4
	XCASE                         = 0x4
	XTABS                         = 0x1800
)


const (
	E2BIG           = syscall.Errno(0x7)
	EACCES          = syscall.Errno(0xd)
	EADDRINUSE      = syscall.Errno(0x7d)
	EADDRNOTAVAIL   = syscall.Errno(0x7e)
	EADV            = syscall.Errno(0x44)
	EAFNOSUPPORT    = syscall.Errno(0x7c)
	EAGAIN          = syscall.Errno(0xb)
	EALREADY        = syscall.Errno(0x95)
	EBADE           = syscall.Errno(0x32)
	EBADF           = syscall.Errno(0x9)
	EBADFD          = syscall.Errno(0x51)
	EBADMSG         = syscall.Errno(0x4d)
	EBADR           = syscall.Errno(0x33)
	EBADRQC         = syscall.Errno(0x36)
	EBADSLT         = syscall.Errno(0x37)
	EBFONT          = syscall.Errno(0x39)
	EBUSY           = syscall.Errno(0x10)
	ECANCELED       = syscall.Errno(0x2f)
	ECHILD          = syscall.Errno(0xa)
	ECHRNG          = syscall.Errno(0x25)
	ECOMM           = syscall.Errno(0x46)
	ECONNABORTED    = syscall.Errno(0x82)
	ECONNREFUSED    = syscall.Errno(0x92)
	ECONNRESET      = syscall.Errno(0x83)
	EDEADLK         = syscall.Errno(0x2d)
	EDEADLOCK       = syscall.Errno(0x38)
	EDESTADDRREQ    = syscall.Errno(0x60)
	EDOM            = syscall.Errno(0x21)
	EDQUOT          = syscall.Errno(0x31)
	EEXIST          = syscall.Errno(0x11)
	EFAULT          = syscall.Errno(0xe)
	EFBIG           = syscall.Errno(0x1b)
	EHOSTDOWN       = syscall.Errno(0x93)
	EHOSTUNREACH    = syscall.Errno(0x94)
	EIDRM           = syscall.Errno(0x24)
	EILSEQ          = syscall.Errno(0x58)
	EINPROGRESS     = syscall.Errno(0x96)
	EINTR           = syscall.Errno(0x4)
	EINVAL          = syscall.Errno(0x16)
	EIO             = syscall.Errno(0x5)
	EISCONN         = syscall.Errno(0x85)
	EISDIR          = syscall.Errno(0x15)
	EL2HLT          = syscall.Errno(0x2c)
	EL2NSYNC        = syscall.Errno(0x26)
	EL3HLT          = syscall.Errno(0x27)
	EL3RST          = syscall.Errno(0x28)
	ELIBACC         = syscall.Errno(0x53)
	ELIBBAD         = syscall.Errno(0x54)
	ELIBEXEC        = syscall.Errno(0x57)
	ELIBMAX         = syscall.Errno(0x56)
	ELIBSCN         = syscall.Errno(0x55)
	ELNRNG          = syscall.Errno(0x29)
	ELOCKUNMAPPED   = syscall.Errno(0x48)
	ELOOP           = syscall.Errno(0x5a)
	EMFILE          = syscall.Errno(0x18)
	EMLINK          = syscall.Errno(0x1f)
	EMSGSIZE        = syscall.Errno(0x61)
	EMULTIHOP       = syscall.Errno(0x4a)
	ENAMETOOLONG    = syscall.Errno(0x4e)
	ENETDOWN        = syscall.Errno(0x7f)
	ENETRESET       = syscall.Errno(0x81)
	ENETUNREACH     = syscall.Errno(0x80)
	ENFILE          = syscall.Errno(0x17)
	ENOANO          = syscall.Errno(0x35)
	ENOBUFS         = syscall.Errno(0x84)
	ENOCSI          = syscall.Errno(0x2b)
	ENODATA         = syscall.Errno(0x3d)
	ENODEV          = syscall.Errno(0x13)
	ENOENT          = syscall.Errno(0x2)
	ENOEXEC         = syscall.Errno(0x8)
	ENOLCK          = syscall.Errno(0x2e)
	ENOLINK         = syscall.Errno(0x43)
	ENOMEM          = syscall.Errno(0xc)
	ENOMSG          = syscall.Errno(0x23)
	ENONET          = syscall.Errno(0x40)
	ENOPKG          = syscall.Errno(0x41)
	ENOPROTOOPT     = syscall.Errno(0x63)
	ENOSPC          = syscall.Errno(0x1c)
	ENOSR           = syscall.Errno(0x3f)
	ENOSTR          = syscall.Errno(0x3c)
	ENOSYS          = syscall.Errno(0x59)
	ENOTACTIVE      = syscall.Errno(0x49)
	ENOTBLK         = syscall.Errno(0xf)
	ENOTCONN        = syscall.Errno(0x86)
	ENOTDIR         = syscall.Errno(0x14)
	ENOTEMPTY       = syscall.Errno(0x5d)
	ENOTRECOVERABLE = syscall.Errno(0x3b)
	ENOTSOCK        = syscall.Errno(0x5f)
	ENOTSUP         = syscall.Errno(0x30)
	ENOTTY          = syscall.Errno(0x19)
	ENOTUNIQ        = syscall.Errno(0x50)
	ENXIO           = syscall.Errno(0x6)
	EOPNOTSUPP      = syscall.Errno(0x7a)
	EOVERFLOW       = syscall.Errno(0x4f)
	EOWNERDEAD      = syscall.Errno(0x3a)
	EPERM           = syscall.Errno(0x1)
	EPFNOSUPPORT    = syscall.Errno(0x7b)
	EPIPE           = syscall.Errno(0x20)
	EPROTO          = syscall.Errno(0x47)
	EPROTONOSUPPORT = syscall.Errno(0x78)
	EPROTOTYPE      = syscall.Errno(0x62)
	ERANGE          = syscall.Errno(0x22)
	EREMCHG         = syscall.Errno(0x52)
	EREMOTE         = syscall.Errno(0x42)
	ERESTART        = syscall.Errno(0x5b)
	EROFS           = syscall.Errno(0x1e)
	ESHUTDOWN       = syscall.Errno(0x8f)
	ESOCKTNOSUPPORT = syscall.Errno(0x79)
	ESPIPE          = syscall.Errno(0x1d)
	ESRCH           = syscall.Errno(0x3)
	ESRMNT          = syscall.Errno(0x45)
	ESTALE          = syscall.Errno(0x97)
	ESTRPIPE        = syscall.Errno(0x5c)
	ETIME           = syscall.Errno(0x3e)
	ETIMEDOUT       = syscall.Errno(0x91)
	ETOOMANYREFS    = syscall.Errno(0x90)
	ETXTBSY         = syscall.Errno(0x1a)
	EUNATCH         = syscall.Errno(0x2a)
	EUSERS          = syscall.Errno(0x5e)
	EWOULDBLOCK     = syscall.Errno(0xb)
	EXDEV           = syscall.Errno(0x12)
	EXFULL          = syscall.Errno(0x34)
)


const (
	SIGABRT    = syscall.Signal(0x6)
	SIGALRM    = syscall.Signal(0xe)
	SIGBUS     = syscall.Signal(0xa)
	SIGCANCEL  = syscall.Signal(0x24)
	SIGCHLD    = syscall.Signal(0x12)
	SIGCLD     = syscall.Signal(0x12)
	SIGCONT    = syscall.Signal(0x19)
	SIGEMT     = syscall.Signal(0x7)
	SIGFPE     = syscall.Signal(0x8)
	SIGFREEZE  = syscall.Signal(0x22)
	SIGHUP     = syscall.Signal(0x1)
	SIGILL     = syscall.Signal(0x4)
	SIGINFO    = syscall.Signal(0x29)
	SIGINT     = syscall.Signal(0x2)
	SIGIO      = syscall.Signal(0x16)
	SIGIOT     = syscall.Signal(0x6)
	SIGJVM1    = syscall.Signal(0x27)
	SIGJVM2    = syscall.Signal(0x28)
	SIGKILL    = syscall.Signal(0x9)
	SIGLOST    = syscall.Signal(0x25)
	SIGLWP     = syscall.Signal(0x21)
	SIGPIPE    = syscall.Signal(0xd)
	SIGPOLL    = syscall.Signal(0x16)
	SIGPROF    = syscall.Signal(0x1d)
	SIGPWR     = syscall.Signal(0x13)
	SIGQUIT    = syscall.Signal(0x3)
	SIGSEGV    = syscall.Signal(0xb)
	SIGSTOP    = syscall.Signal(0x17)
	SIGSYS     = syscall.Signal(0xc)
	SIGTERM    = syscall.Signal(0xf)
	SIGTHAW    = syscall.Signal(0x23)
	SIGTRAP    = syscall.Signal(0x5)
	SIGTSTP    = syscall.Signal(0x18)
	SIGTTIN    = syscall.Signal(0x1a)
	SIGTTOU    = syscall.Signal(0x1b)
	SIGURG     = syscall.Signal(0x15)
	SIGUSR1    = syscall.Signal(0x10)
	SIGUSR2    = syscall.Signal(0x11)
	SIGVTALRM  = syscall.Signal(0x1c)
	SIGWAITING = syscall.Signal(0x20)
	SIGWINCH   = syscall.Signal(0x14)
	SIGXCPU    = syscall.Signal(0x1e)
	SIGXFSZ    = syscall.Signal(0x1f)
	SIGXRES    = syscall.Signal(0x26)
)


var errors = [...]string{
	1:   "not owner",
	2:   "no such file or directory",
	3:   "no such process",
	4:   "interrupted system call",
	5:   "I/O error",
	6:   "no such device or address",
	7:   "arg list too long",
	8:   "exec format error",
	9:   "bad file number",
	10:  "no child processes",
	11:  "resource temporarily unavailable",
	12:  "not enough space",
	13:  "permission denied",
	14:  "bad address",
	15:  "block device required",
	16:  "device busy",
	17:  "file exists",
	18:  "cross-device link",
	19:  "no such device",
	20:  "not a directory",
	21:  "is a directory",
	22:  "invalid argument",
	23:  "file table overflow",
	24:  "too many open files",
	25:  "inappropriate ioctl for device",
	26:  "text file busy",
	27:  "file too large",
	28:  "no space left on device",
	29:  "illegal seek",
	30:  "read-only file system",
	31:  "too many links",
	32:  "broken pipe",
	33:  "argument out of domain",
	34:  "result too large",
	35:  "no message of desired type",
	36:  "identifier removed",
	37:  "channel number out of range",
	38:  "level 2 not synchronized",
	39:  "level 3 halted",
	40:  "level 3 reset",
	41:  "link number out of range",
	42:  "protocol driver not attached",
	43:  "no CSI structure available",
	44:  "level 2 halted",
	45:  "deadlock situation detected/avoided",
	46:  "no record locks available",
	47:  "operation canceled",
	48:  "operation not supported",
	49:  "disc quota exceeded",
	50:  "bad exchange descriptor",
	51:  "bad request descriptor",
	52:  "message tables full",
	53:  "anode table overflow",
	54:  "bad request code",
	55:  "invalid slot",
	56:  "file locking deadlock",
	57:  "bad font file format",
	58:  "owner of the lock died",
	59:  "lock is not recoverable",
	60:  "not a stream device",
	61:  "no data available",
	62:  "timer expired",
	63:  "out of stream resources",
	64:  "machine is not on the network",
	65:  "package not installed",
	66:  "object is remote",
	67:  "link has been severed",
	68:  "advertise error",
	69:  "srmount error",
	70:  "communication error on send",
	71:  "protocol error",
	72:  "locked lock was unmapped ",
	73:  "facility is not active",
	74:  "multihop attempted",
	77:  "not a data message",
	78:  "file name too long",
	79:  "value too large for defined data type",
	80:  "name not unique on network",
	81:  "file descriptor in bad state",
	82:  "remote address changed",
	83:  "can not access a needed shared library",
	84:  "accessing a corrupted shared library",
	85:  ".lib section in a.out corrupted",
	86:  "attempting to link in more shared libraries than system limit",
	87:  "can not exec a shared library directly",
	88:  "illegal byte sequence",
	89:  "operation not applicable",
	90:  "number of symbolic links encountered during path name traversal exceeds MAXSYMLINKS",
	91:  "error 91",
	92:  "error 92",
	93:  "directory not empty",
	94:  "too many users",
	95:  "socket operation on non-socket",
	96:  "destination address required",
	97:  "message too long",
	98:  "protocol wrong type for socket",
	99:  "option not supported by protocol",
	120: "protocol not supported",
	121: "socket type not supported",
	122: "operation not supported on transport endpoint",
	123: "protocol family not supported",
	124: "address family not supported by protocol family",
	125: "address already in use",
	126: "cannot assign requested address",
	127: "network is down",
	128: "network is unreachable",
	129: "network dropped connection because of reset",
	130: "software caused connection abort",
	131: "connection reset by peer",
	132: "no buffer space available",
	133: "transport endpoint is already connected",
	134: "transport endpoint is not connected",
	143: "cannot send after socket shutdown",
	144: "too many references: cannot splice",
	145: "connection timed out",
	146: "connection refused",
	147: "host is down",
	148: "no route to host",
	149: "operation already in progress",
	150: "operation now in progress",
	151: "stale NFS file handle",
}


var signals = [...]string{
	1:  "hangup",
	2:  "interrupt",
	3:  "quit",
	4:  "illegal Instruction",
	5:  "trace/Breakpoint Trap",
	6:  "abort",
	7:  "emulation Trap",
	8:  "arithmetic Exception",
	9:  "killed",
	10: "bus Error",
	11: "segmentation Fault",
	12: "bad System Call",
	13: "broken Pipe",
	14: "alarm Clock",
	15: "terminated",
	16: "user Signal 1",
	17: "user Signal 2",
	18: "child Status Changed",
	19: "power-Fail/Restart",
	20: "window Size Change",
	21: "urgent Socket Condition",
	22: "pollable Event",
	23: "stopped (signal)",
	24: "stopped (user)",
	25: "continued",
	26: "stopped (tty input)",
	27: "stopped (tty output)",
	28: "virtual Timer Expired",
	29: "profiling Timer Expired",
	30: "cpu Limit Exceeded",
	31: "file Size Limit Exceeded",
	32: "no runnable lwp",
	33: "inter-lwp signal",
	34: "checkpoint Freeze",
	35: "checkpoint Thaw",
	36: "thread Cancellation",
	37: "resource Lost",
	38: "resource Control Exceeded",
	39: "reserved for JVM 1",
	40: "reserved for JVM 2",
	41: "information Request",
}
