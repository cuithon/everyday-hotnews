// mkerrors.sh
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

// godefs -c arm-gcc -gsyscall _const.c

// MACHINE GENERATED - DO NOT EDIT.

package syscall

// Constants
const (
	AF_APPLETALK                     = 0x5
	AF_ASH                           = 0x12
	AF_ATMPVC                        = 0x8
	AF_ATMSVC                        = 0x14
	AF_AX25                          = 0x3
	AF_BLUETOOTH                     = 0x1f
	AF_BRIDGE                        = 0x7
	AF_DECnet                        = 0xc
	AF_ECONET                        = 0x13
	AF_FILE                          = 0x1
	AF_INET                          = 0x2
	AF_INET6                         = 0xa
	AF_IPX                           = 0x4
	AF_IRDA                          = 0x17
	AF_KEY                           = 0xf
	AF_LOCAL                         = 0x1
	AF_MAX                           = 0x20
	AF_NETBEUI                       = 0xd
	AF_NETLINK                       = 0x10
	AF_NETROM                        = 0x6
	AF_PACKET                        = 0x11
	AF_PPPOX                         = 0x18
	AF_ROSE                          = 0xb
	AF_ROUTE                         = 0x10
	AF_SECURITY                      = 0xe
	AF_SNA                           = 0x16
	AF_UNIX                          = 0x1
	AF_UNSPEC                        = 0
	AF_WANPIPE                       = 0x19
	AF_X25                           = 0x9
	E2BIG                            = 0x7
	EACCES                           = 0xd
	EADDRINUSE                       = 0x62
	EADDRNOTAVAIL                    = 0x63
	EADV                             = 0x44
	EAFNOSUPPORT                     = 0x61
	EAGAIN                           = 0xb
	EALREADY                         = 0x72
	EBADE                            = 0x34
	EBADF                            = 0x9
	EBADFD                           = 0x4d
	EBADMSG                          = 0x4a
	EBADR                            = 0x35
	EBADRQC                          = 0x38
	EBADSLT                          = 0x39
	EBFONT                           = 0x3b
	EBUSY                            = 0x10
	ECANCELED                        = 0x7d
	ECHILD                           = 0xa
	ECHRNG                           = 0x2c
	ECOMM                            = 0x46
	ECONNABORTED                     = 0x67
	ECONNREFUSED                     = 0x6f
	ECONNRESET                       = 0x68
	EDEADLK                          = 0x23
	EDEADLOCK                        = 0x23
	EDESTADDRREQ                     = 0x59
	EDOM                             = 0x21
	EDOTDOT                          = 0x49
	EDQUOT                           = 0x7a
	EEXIST                           = 0x11
	EFAULT                           = 0xe
	EFBIG                            = 0x1b
	EHOSTDOWN                        = 0x70
	EHOSTUNREACH                     = 0x71
	EIDRM                            = 0x2b
	EILSEQ                           = 0x54
	EINPROGRESS                      = 0x73
	EINTR                            = 0x4
	EINVAL                           = 0x16
	EIO                              = 0x5
	EISCONN                          = 0x6a
	EISDIR                           = 0x15
	EISNAM                           = 0x78
	EKEYEXPIRED                      = 0x7f
	EKEYREJECTED                     = 0x81
	EKEYREVOKED                      = 0x80
	EL2HLT                           = 0x33
	EL2NSYNC                         = 0x2d
	EL3HLT                           = 0x2e
	EL3RST                           = 0x2f
	ELF_NGREG                        = 0x12
	ELF_PRARGSZ                      = 0x50
	ELIBACC                          = 0x4f
	ELIBBAD                          = 0x50
	ELIBEXEC                         = 0x53
	ELIBMAX                          = 0x52
	ELIBSCN                          = 0x51
	ELNRNG                           = 0x30
	ELOOP                            = 0x28
	EMEDIUMTYPE                      = 0x7c
	EMFILE                           = 0x18
	EMLINK                           = 0x1f
	EMSGSIZE                         = 0x5a
	EMULTIHOP                        = 0x48
	ENAMETOOLONG                     = 0x24
	ENAVAIL                          = 0x77
	ENETDOWN                         = 0x64
	ENETRESET                        = 0x66
	ENETUNREACH                      = 0x65
	ENFILE                           = 0x17
	ENOANO                           = 0x37
	ENOBUFS                          = 0x69
	ENOCSI                           = 0x32
	ENODATA                          = 0x3d
	ENODEV                           = 0x13
	ENOENT                           = 0x2
	ENOEXEC                          = 0x8
	ENOKEY                           = 0x7e
	ENOLCK                           = 0x25
	ENOLINK                          = 0x43
	ENOMEDIUM                        = 0x7b
	ENOMEM                           = 0xc
	ENOMSG                           = 0x2a
	ENONET                           = 0x40
	ENOPKG                           = 0x41
	ENOPROTOOPT                      = 0x5c
	ENOSPC                           = 0x1c
	ENOSR                            = 0x3f
	ENOSTR                           = 0x3c
	ENOSYS                           = 0x26
	ENOTBLK                          = 0xf
	ENOTCONN                         = 0x6b
	ENOTDIR                          = 0x14
	ENOTEMPTY                        = 0x27
	ENOTNAM                          = 0x76
	ENOTRECOVERABLE                  = 0x83
	ENOTSOCK                         = 0x58
	ENOTSUP                          = 0x5f
	ENOTTY                           = 0x19
	ENOTUNIQ                         = 0x4c
	ENXIO                            = 0x6
	EOPNOTSUPP                       = 0x5f
	EOVERFLOW                        = 0x4b
	EOWNERDEAD                       = 0x82
	EPERM                            = 0x1
	EPFNOSUPPORT                     = 0x60
	EPIPE                            = 0x20
	EPOLLERR                         = 0x8
	EPOLLET                          = -0x80000000
	EPOLLHUP                         = 0x10
	EPOLLIN                          = 0x1
	EPOLLMSG                         = 0x400
	EPOLLOUT                         = 0x4
	EPOLLPRI                         = 0x2
	EPOLLRDBAND                      = 0x80
	EPOLLRDNORM                      = 0x40
	EPOLLWRBAND                      = 0x200
	EPOLLWRNORM                      = 0x100
	EPOLL_CTL_ADD                    = 0x1
	EPOLL_CTL_DEL                    = 0x2
	EPOLL_CTL_MOD                    = 0x3
	EPROTO                           = 0x47
	EPROTONOSUPPORT                  = 0x5d
	EPROTOTYPE                       = 0x5b
	ERANGE                           = 0x22
	EREMCHG                          = 0x4e
	EREMOTE                          = 0x42
	EREMOTEIO                        = 0x79
	ERESTART                         = 0x55
	EROFS                            = 0x1e
	ESHUTDOWN                        = 0x6c
	ESOCKTNOSUPPORT                  = 0x5e
	ESPIPE                           = 0x1d
	ESRCH                            = 0x3
	ESRMNT                           = 0x45
	ESTALE                           = 0x74
	ESTRPIPE                         = 0x56
	ETIME                            = 0x3e
	ETIMEDOUT                        = 0x6e
	ETOOMANYREFS                     = 0x6d
	ETXTBSY                          = 0x1a
	EUCLEAN                          = 0x75
	EUNATCH                          = 0x31
	EUSERS                           = 0x57
	EWOULDBLOCK                      = 0xb
	EXDEV                            = 0x12
	EXFULL                           = 0x36
	EXPR_NEST_MAX                    = 0x20
	FD_CLOEXEC                       = 0x1
	FD_SETSIZE                       = 0x400
	F_DUPFD                          = 0
	F_EXLCK                          = 0x4
	F_GETFD                          = 0x1
	F_GETFL                          = 0x3
	F_GETLEASE                       = 0x401
	F_GETLK                          = 0xc
	F_GETLK64                        = 0xc
	F_GETOWN                         = 0x9
	F_GETSIG                         = 0xb
	F_LOCK                           = 0x1
	F_NOTIFY                         = 0x402
	F_OK                             = 0
	F_RDLCK                          = 0
	F_SETFD                          = 0x2
	F_SETFL                          = 0x4
	F_SETLEASE                       = 0x400
	F_SETLK                          = 0xd
	F_SETLK64                        = 0xd
	F_SETLKW                         = 0xe
	F_SETLKW64                       = 0xe
	F_SETOWN                         = 0x8
	F_SETSIG                         = 0xa
	F_SHLCK                          = 0x8
	F_TEST                           = 0x3
	F_TLOCK                          = 0x2
	F_ULOCK                          = 0
	F_UNLCK                          = 0x2
	F_WRLCK                          = 0x1
	IPPROTO_AH                       = 0x33
	IPPROTO_COMP                     = 0x6c
	IPPROTO_DSTOPTS                  = 0x3c
	IPPROTO_EGP                      = 0x8
	IPPROTO_ENCAP                    = 0x62
	IPPROTO_ESP                      = 0x32
	IPPROTO_FRAGMENT                 = 0x2c
	IPPROTO_GRE                      = 0x2f
	IPPROTO_HOPOPTS                  = 0
	IPPROTO_ICMP                     = 0x1
	IPPROTO_ICMPV6                   = 0x3a
	IPPROTO_IDP                      = 0x16
	IPPROTO_IGMP                     = 0x2
	IPPROTO_IP                       = 0
	IPPROTO_IPIP                     = 0x4
	IPPROTO_IPV6                     = 0x29
	IPPROTO_MTP                      = 0x5c
	IPPROTO_NONE                     = 0x3b
	IPPROTO_PIM                      = 0x67
	IPPROTO_PUP                      = 0xc
	IPPROTO_RAW                      = 0xff
	IPPROTO_ROUTING                  = 0x2b
	IPPROTO_RSVP                     = 0x2e
	IPPROTO_TCP                      = 0x6
	IPPROTO_TP                       = 0x1d
	IPPROTO_UDP                      = 0x11
	IPV6_ADDRFORM                    = 0x1
	IPV6_ADD_MEMBERSHIP              = 0x14
	IPV6_AUTHHDR                     = 0xa
	IPV6_CHECKSUM                    = 0x7
	IPV6_DROP_MEMBERSHIP             = 0x15
	IPV6_DSTOPTS                     = 0x4
	IPV6_HOPLIMIT                    = 0x8
	IPV6_HOPOPTS                     = 0x3
	IPV6_JOIN_GROUP                  = 0x14
	IPV6_LEAVE_GROUP                 = 0x15
	IPV6_MTU                         = 0x18
	IPV6_MTU_DISCOVER                = 0x17
	IPV6_MULTICAST_HOPS              = 0x12
	IPV6_MULTICAST_IF                = 0x11
	IPV6_MULTICAST_LOOP              = 0x13
	IPV6_NEXTHOP                     = 0x9
	IPV6_PKTINFO                     = 0x2
	IPV6_PKTOPTIONS                  = 0x6
	IPV6_PMTUDISC_DO                 = 0x2
	IPV6_PMTUDISC_DONT               = 0
	IPV6_PMTUDISC_WANT               = 0x1
	IPV6_RECVERR                     = 0x19
	IPV6_ROUTER_ALERT                = 0x16
	IPV6_RTHDR                       = 0x5
	IPV6_RTHDR_LOOSE                 = 0
	IPV6_RTHDR_STRICT                = 0x1
	IPV6_RTHDR_TYPE_0                = 0
	IPV6_RXDSTOPTS                   = 0x4
	IPV6_RXHOPOPTS                   = 0x3
	IPV6_UNICAST_HOPS                = 0x10
	IP_ADD_MEMBERSHIP                = 0x23
	IP_DEFAULT_MULTICAST_LOOP        = 0x1
	IP_DEFAULT_MULTICAST_TTL         = 0x1
	IP_DF                            = 0x4000
	IP_DROP_MEMBERSHIP               = 0x24
	IP_HDRINCL                       = 0x3
	IP_MAXPACKET                     = 0xffff
	IP_MAX_MEMBERSHIPS               = 0x14
	IP_MF                            = 0x2000
	IP_MSS                           = 0x240
	IP_MTU_DISCOVER                  = 0xa
	IP_MULTICAST_IF                  = 0x20
	IP_MULTICAST_LOOP                = 0x22
	IP_MULTICAST_TTL                 = 0x21
	IP_OFFMASK                       = 0x1fff
	IP_OPTIONS                       = 0x4
	IP_PKTINFO                       = 0x8
	IP_PKTOPTIONS                    = 0x9
	IP_PMTUDISC                      = 0xa
	IP_PMTUDISC_DO                   = 0x2
	IP_PMTUDISC_DONT                 = 0
	IP_PMTUDISC_WANT                 = 0x1
	IP_RECVERR                       = 0xb
	IP_RECVOPTS                      = 0x6
	IP_RECVRETOPTS                   = 0x7
	IP_RECVTOS                       = 0xd
	IP_RECVTTL                       = 0xc
	IP_RETOPTS                       = 0x7
	IP_RF                            = 0x8000
	IP_ROUTER_ALERT                  = 0x5
	IP_TOS                           = 0x1
	IP_TTL                           = 0x2
	NAME_MAX                         = 0xff
	O_ACCMODE                        = 0x3
	O_APPEND                         = 0x400
	O_ASYNC                          = 0x2000
	O_CREAT                          = 0x40
	O_DIRECT                         = 0x10000
	O_DIRECTORY                      = 0x4000
	O_DSYNC                          = 0x1000
	O_EXCL                           = 0x80
	O_FSYNC                          = 0x1000
	O_LARGEFILE                      = 0x20000
	O_NDELAY                         = 0x800
	O_NOCTTY                         = 0x100
	O_NOFOLLOW                       = 0x8000
	O_NONBLOCK                       = 0x800
	O_RDONLY                         = 0
	O_RDWR                           = 0x2
	O_RSYNC                          = 0x1000
	O_SYNC                           = 0x1000
	O_TRUNC                          = 0x200
	O_WRONLY                         = 0x1
	PTRACE_ATTACH                    = 0x10
	PTRACE_CONT                      = 0x7
	PTRACE_DETACH                    = 0x11
	PTRACE_EVENT_CLONE               = 0x3
	PTRACE_EVENT_EXEC                = 0x4
	PTRACE_EVENT_EXIT                = 0x6
	PTRACE_EVENT_FORK                = 0x1
	PTRACE_EVENT_VFORK               = 0x2
	PTRACE_EVENT_VFORK_DONE          = 0x5
	PTRACE_GETEVENTMSG               = 0x4201
	PTRACE_GETFPREGS                 = 0xe
	PTRACE_GETREGS                   = 0xc
	PTRACE_GETSIGINFO                = 0x4202
	PTRACE_GETWMMXREGS               = 0x12
	PTRACE_GET_THREAD_AREA           = 0x16
	PTRACE_KILL                      = 0x8
	PTRACE_OLDSETOPTIONS             = 0x15
	PTRACE_O_MASK                    = 0x7f
	PTRACE_O_TRACECLONE              = 0x8
	PTRACE_O_TRACEEXEC               = 0x10
	PTRACE_O_TRACEEXIT               = 0x40
	PTRACE_O_TRACEFORK               = 0x2
	PTRACE_O_TRACESYSGOOD            = 0x1
	PTRACE_O_TRACEVFORK              = 0x4
	PTRACE_O_TRACEVFORKDONE          = 0x20
	PTRACE_PEEKDATA                  = 0x2
	PTRACE_PEEKTEXT                  = 0x1
	PTRACE_PEEKUSR                   = 0x3
	PTRACE_POKEDATA                  = 0x5
	PTRACE_POKETEXT                  = 0x4
	PTRACE_POKEUSR                   = 0x6
	PTRACE_SETFPREGS                 = 0xf
	PTRACE_SETOPTIONS                = 0x4200
	PTRACE_SETREGS                   = 0xd
	PTRACE_SETSIGINFO                = 0x4203
	PTRACE_SETWMMXREGS               = 0x13
	PTRACE_SINGLESTEP                = 0x9
	PTRACE_SYSCALL                   = 0x18
	PTRACE_TRACEME                   = 0
	SHUT_RD                          = 0
	SHUT_RDWR                        = 0x2
	SHUT_WR                          = 0x1
	SIGABRT                          = 0x6
	SIGALRM                          = 0xe
	SIGBUS                           = 0x7
	SIGCHLD                          = 0x11
	SIGCLD                           = 0x11
	SIGCONT                          = 0x12
	SIGFPE                           = 0x8
	SIGHUP                           = 0x1
	SIGILL                           = 0x4
	SIGINT                           = 0x2
	SIGIO                            = 0x1d
	SIGIOT                           = 0x6
	SIGKILL                          = 0x9
	SIGPIPE                          = 0xd
	SIGPOLL                          = 0x1d
	SIGPROF                          = 0x1b
	SIGPWR                           = 0x1e
	SIGQUIT                          = 0x3
	SIGSEGV                          = 0xb
	SIGSTKFLT                        = 0x10
	SIGSTOP                          = 0x13
	SIGSYS                           = 0x1f
	SIGTERM                          = 0xf
	SIGTRAP                          = 0x5
	SIGTSTP                          = 0x14
	SIGTTIN                          = 0x15
	SIGTTOU                          = 0x16
	SIGUNUSED                        = 0x1f
	SIGURG                           = 0x17
	SIGUSR1                          = 0xa
	SIGUSR2                          = 0xc
	SIGVTALRM                        = 0x1a
	SIGWINCH                         = 0x1c
	SIGXCPU                          = 0x18
	SIGXFSZ                          = 0x19
	SOCK_DGRAM                       = 0x2
	SOCK_PACKET                      = 0xa
	SOCK_RAW                         = 0x3
	SOCK_RDM                         = 0x4
	SOCK_SEQPACKET                   = 0x5
	SOCK_STREAM                      = 0x1
	SOL_AAL                          = 0x109
	SOL_ATM                          = 0x108
	SOL_DECNET                       = 0x105
	SOL_ICMPV6                       = 0x3a
	SOL_IP                           = 0
	SOL_IPV6                         = 0x29
	SOL_IRDA                         = 0x10a
	SOL_PACKET                       = 0x107
	SOL_RAW                          = 0xff
	SOL_SOCKET                       = 0x1
	SOL_TCP                          = 0x6
	SOL_X25                          = 0x106
	SOMAXCONN                        = 0x80
	SO_ACCEPTCONN                    = 0x1e
	SO_ATTACH_FILTER                 = 0x1a
	SO_BINDTODEVICE                  = 0x19
	SO_BROADCAST                     = 0x6
	SO_BSDCOMPAT                     = 0xe
	SO_DEBUG                         = 0x1
	SO_DETACH_FILTER                 = 0x1b
	SO_DONTROUTE                     = 0x5
	SO_ERROR                         = 0x4
	SO_KEEPALIVE                     = 0x9
	SO_LINGER                        = 0xd
	SO_NO_CHECK                      = 0xb
	SO_OOBINLINE                     = 0xa
	SO_PASSCRED                      = 0x10
	SO_PEERCRED                      = 0x11
	SO_PEERNAME                      = 0x1c
	SO_PEERSEC                       = 0x1f
	SO_PRIORITY                      = 0xc
	SO_RCVBUF                        = 0x8
	SO_RCVLOWAT                      = 0x12
	SO_RCVTIMEO                      = 0x14
	SO_REUSEADDR                     = 0x2
	SO_SECURITY_AUTHENTICATION       = 0x16
	SO_SECURITY_ENCRYPTION_NETWORK   = 0x18
	SO_SECURITY_ENCRYPTION_TRANSPORT = 0x17
	SO_SNDBUF                        = 0x7
	SO_SNDLOWAT                      = 0x13
	SO_SNDTIMEO                      = 0x15
	SO_TIMESTAMP                     = 0x1d
	SO_TYPE                          = 0x3
	S_BLKSIZE                        = 0x200
	S_IEXEC                          = 0x40
	S_IFBLK                          = 0x6000
	S_IFCHR                          = 0x2000
	S_IFDIR                          = 0x4000
	S_IFIFO                          = 0x1000
	S_IFLNK                          = 0xa000
	S_IFMT                           = 0xf000
	S_IFREG                          = 0x8000
	S_IFSOCK                         = 0xc000
	S_IREAD                          = 0x100
	S_IRGRP                          = 0x20
	S_IROTH                          = 0x4
	S_IRUSR                          = 0x100
	S_IRWXG                          = 0x38
	S_IRWXO                          = 0x7
	S_IRWXU                          = 0x1c0
	S_ISGID                          = 0x400
	S_ISUID                          = 0x800
	S_ISVTX                          = 0x200
	S_IWGRP                          = 0x10
	S_IWOTH                          = 0x2
	S_IWRITE                         = 0x80
	S_IWUSR                          = 0x80
	S_IXGRP                          = 0x8
	S_IXOTH                          = 0x1
	S_IXUSR                          = 0x40
	TCP_CORK                         = 0x3
	TCP_DEFER_ACCEPT                 = 0x9
	TCP_INFO                         = 0xb
	TCP_KEEPCNT                      = 0x6
	TCP_KEEPIDLE                     = 0x4
	TCP_KEEPINTVL                    = 0x5
	TCP_LINGER2                      = 0x8
	TCP_MAXSEG                       = 0x2
	TCP_MAXWIN                       = 0xffff
	TCP_MAX_WINSHIFT                 = 0xe
	TCP_MSS                          = 0x200
	TCP_NODELAY                      = 0x1
	TCP_QUICKACK                     = 0xc
	TCP_SYNCNT                       = 0x7
	TCP_WINDOW_CLAMP                 = 0xa
	WALL                             = 0x40000000
	WCLONE                           = 0x80000000
	WCOREFLAG                        = 0x80
	WNOHANG                          = 0x1
	WORDSIZE                         = 0x20
	WUNTRACED                        = 0x2
)

// Types


// Error table
var errors = [...]string{
	7:   "argument list too long",
	13:  "permission denied",
	98:  "address already in use",
	99:  "cannot assign requested address",
	68:  "advertise error",
	97:  "address family not supported by protocol",
	11:  "resource temporarily unavailable",
	114: "operation already in progress",
	52:  "invalid exchange",
	9:   "bad file descriptor",
	77:  "file descriptor in bad state",
	74:  "bad message",
	53:  "invalid request descriptor",
	56:  "invalid request code",
	57:  "invalid slot",
	59:  "bad font file format",
	16:  "device or resource busy",
	125: "operation canceled",
	10:  "no child processes",
	44:  "channel number out of range",
	70:  "communication error on send",
	103: "software caused connection abort",
	111: "connection refused",
	104: "connection reset by peer",
	35:  "resource deadlock avoided",
	89:  "destination address required",
	33:  "numerical argument out of domain",
	73:  "RFS specific error",
	122: "disk quota exceeded",
	17:  "file exists",
	14:  "bad address",
	27:  "file too large",
	112: "host is down",
	113: "no route to host",
	43:  "identifier removed",
	84:  "invalid or incomplete multibyte or wide character",
	115: "operation now in progress",
	4:   "interrupted system call",
	22:  "invalid argument",
	5:   "input/output error",
	106: "transport endpoint is already connected",
	21:  "is a directory",
	120: "is a named type file",
	127: "unknown error 127",
	129: "unknown error 129",
	128: "unknown error 128",
	51:  "level 2 halted",
	45:  "level 2 not synchronized",
	46:  "level 3 halted",
	47:  "level 3 reset",
	79:  "can not access a needed shared library",
	80:  "accessing a corrupted shared library",
	83:  "cannot exec a shared library directly",
	82:  "attempting to link in too many shared libraries",
	81:  ".lib section in a.out corrupted",
	48:  "link number out of range",
	40:  "too many levels of symbolic links",
	124: "wrong medium type",
	24:  "too many open files",
	31:  "too many links",
	90:  "message too long",
	72:  "multihop attempted",
	36:  "file name too long",
	119: "no XENIX semaphores available",
	100: "network is down",
	102: "network dropped connection on reset",
	101: "network is unreachable",
	23:  "too many open files in system",
	55:  "no anode",
	105: "no buffer space available",
	50:  "no CSI structure available",
	61:  "no data available",
	19:  "no such device",
	2:   "no such file or directory",
	8:   "exec format error",
	126: "unknown error 126",
	37:  "no locks available",
	67:  "link has been severed",
	123: "no medium found",
	12:  "cannot allocate memory",
	42:  "no message of desired type",
	64:  "machine is not on the network",
	65:  "package not installed",
	92:  "protocol not available",
	28:  "no space left on device",
	63:  "out of streams resources",
	60:  "device not a stream",
	38:  "function not implemented",
	15:  "block device required",
	107: "transport endpoint is not connected",
	20:  "not a directory",
	39:  "directory not empty",
	118: "not a XENIX named type file",
	131: "unknown error 131",
	88:  "socket operation on non-socket",
	95:  "operation not supported",
	25:  "inappropriate ioctl for device",
	76:  "name not unique on network",
	6:   "no such device or address",
	75:  "value too large for defined data type",
	130: "unknown error 130",
	1:   "operation not permitted",
	96:  "protocol family not supported",
	32:  "broken pipe",
	71:  "protocol error",
	93:  "protocol not supported",
	91:  "protocol wrong type for socket",
	34:  "numerical result out of range",
	78:  "remote address changed",
	66:  "object is remote",
	121: "remote I/O error",
	85:  "interrupted system call should be restarted",
	30:  "read-only file system",
	108: "cannot send after transport endpoint shutdown",
	94:  "socket type not supported",
	29:  "illegal seek",
	3:   "no such process",
	69:  "srmount error",
	116: "stale NFS file handle",
	86:  "streams pipe error",
	62:  "timer expired",
	110: "connection timed out",
	109: "too many references: cannot splice",
	26:  "text file busy",
	117: "structure needs cleaning",
	49:  "protocol driver not attached",
	87:  "too many users",
	18:  "invalid cross-device link",
	54:  "exchange full",
}
