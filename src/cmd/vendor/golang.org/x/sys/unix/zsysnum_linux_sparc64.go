// go run linux/mksysnum.go -Wall -Werror -static -I/tmp/sparc64/include /tmp/sparc64/include/asm/unistd.h
// Code generated by the command above; see README.md. DO NOT EDIT.

//go:build sparc64 && linux

package unix

const (
	SYS_RESTART_SYSCALL         = 0
	SYS_EXIT                    = 1
	SYS_FORK                    = 2
	SYS_READ                    = 3
	SYS_WRITE                   = 4
	SYS_OPEN                    = 5
	SYS_CLOSE                   = 6
	SYS_WAIT4                   = 7
	SYS_CREAT                   = 8
	SYS_LINK                    = 9
	SYS_UNLINK                  = 10
	SYS_EXECV                   = 11
	SYS_CHDIR                   = 12
	SYS_CHOWN                   = 13
	SYS_MKNOD                   = 14
	SYS_CHMOD                   = 15
	SYS_LCHOWN                  = 16
	SYS_BRK                     = 17
	SYS_PERFCTR                 = 18
	SYS_LSEEK                   = 19
	SYS_GETPID                  = 20
	SYS_CAPGET                  = 21
	SYS_CAPSET                  = 22
	SYS_SETUID                  = 23
	SYS_GETUID                  = 24
	SYS_VMSPLICE                = 25
	SYS_PTRACE                  = 26
	SYS_ALARM                   = 27
	SYS_SIGALTSTACK             = 28
	SYS_PAUSE                   = 29
	SYS_UTIME                   = 30
	SYS_ACCESS                  = 33
	SYS_NICE                    = 34
	SYS_SYNC                    = 36
	SYS_KILL                    = 37
	SYS_STAT                    = 38
	SYS_SENDFILE                = 39
	SYS_LSTAT                   = 40
	SYS_DUP                     = 41
	SYS_PIPE                    = 42
	SYS_TIMES                   = 43
	SYS_UMOUNT2                 = 45
	SYS_SETGID                  = 46
	SYS_GETGID                  = 47
	SYS_SIGNAL                  = 48
	SYS_GETEUID                 = 49
	SYS_GETEGID                 = 50
	SYS_ACCT                    = 51
	SYS_MEMORY_ORDERING         = 52
	SYS_IOCTL                   = 54
	SYS_REBOOT                  = 55
	SYS_SYMLINK                 = 57
	SYS_READLINK                = 58
	SYS_EXECVE                  = 59
	SYS_UMASK                   = 60
	SYS_CHROOT                  = 61
	SYS_FSTAT                   = 62
	SYS_FSTAT64                 = 63
	SYS_GETPAGESIZE             = 64
	SYS_MSYNC                   = 65
	SYS_VFORK                   = 66
	SYS_PREAD64                 = 67
	SYS_PWRITE64                = 68
	SYS_MMAP                    = 71
	SYS_MUNMAP                  = 73
	SYS_MPROTECT                = 74
	SYS_MADVISE                 = 75
	SYS_VHANGUP                 = 76
	SYS_MINCORE                 = 78
	SYS_GETGROUPS               = 79
	SYS_SETGROUPS               = 80
	SYS_GETPGRP                 = 81
	SYS_SETITIMER               = 83
	SYS_SWAPON                  = 85
	SYS_GETITIMER               = 86
	SYS_SETHOSTNAME             = 88
	SYS_DUP2                    = 90
	SYS_FCNTL                   = 92
	SYS_SELECT                  = 93
	SYS_FSYNC                   = 95
	SYS_SETPRIORITY             = 96
	SYS_SOCKET                  = 97
	SYS_CONNECT                 = 98
	SYS_ACCEPT                  = 99
	SYS_GETPRIORITY             = 100
	SYS_RT_SIGRETURN            = 101
	SYS_RT_SIGACTION            = 102
	SYS_RT_SIGPROCMASK          = 103
	SYS_RT_SIGPENDING           = 104
	SYS_RT_SIGTIMEDWAIT         = 105
	SYS_RT_SIGQUEUEINFO         = 106
	SYS_RT_SIGSUSPEND           = 107
	SYS_SETRESUID               = 108
	SYS_GETRESUID               = 109
	SYS_SETRESGID               = 110
	SYS_GETRESGID               = 111
	SYS_RECVMSG                 = 113
	SYS_SENDMSG                 = 114
	SYS_GETTIMEOFDAY            = 116
	SYS_GETRUSAGE               = 117
	SYS_GETSOCKOPT              = 118
	SYS_GETCWD                  = 119
	SYS_READV                   = 120
	SYS_WRITEV                  = 121
	SYS_SETTIMEOFDAY            = 122
	SYS_FCHOWN                  = 123
	SYS_FCHMOD                  = 124
	SYS_RECVFROM                = 125
	SYS_SETREUID                = 126
	SYS_SETREGID                = 127
	SYS_RENAME                  = 128
	SYS_TRUNCATE                = 129
	SYS_FTRUNCATE               = 130
	SYS_FLOCK                   = 131
	SYS_LSTAT64                 = 132
	SYS_SENDTO                  = 133
	SYS_SHUTDOWN                = 134
	SYS_SOCKETPAIR              = 135
	SYS_MKDIR                   = 136
	SYS_RMDIR                   = 137
	SYS_UTIMES                  = 138
	SYS_STAT64                  = 139
	SYS_SENDFILE64              = 140
	SYS_GETPEERNAME             = 141
	SYS_FUTEX                   = 142
	SYS_GETTID                  = 143
	SYS_GETRLIMIT               = 144
	SYS_SETRLIMIT               = 145
	SYS_PIVOT_ROOT              = 146
	SYS_PRCTL                   = 147
	SYS_PCICONFIG_READ          = 148
	SYS_PCICONFIG_WRITE         = 149
	SYS_GETSOCKNAME             = 150
	SYS_INOTIFY_INIT            = 151
	SYS_INOTIFY_ADD_WATCH       = 152
	SYS_POLL                    = 153
	SYS_GETDENTS64              = 154
	SYS_INOTIFY_RM_WATCH        = 156
	SYS_STATFS                  = 157
	SYS_FSTATFS                 = 158
	SYS_UMOUNT                  = 159
	SYS_SCHED_SET_AFFINITY      = 160
	SYS_SCHED_GET_AFFINITY      = 161
	SYS_GETDOMAINNAME           = 162
	SYS_SETDOMAINNAME           = 163
	SYS_UTRAP_INSTALL           = 164
	SYS_QUOTACTL                = 165
	SYS_SET_TID_ADDRESS         = 166
	SYS_MOUNT                   = 167
	SYS_USTAT                   = 168
	SYS_SETXATTR                = 169
	SYS_LSETXATTR               = 170
	SYS_FSETXATTR               = 171
	SYS_GETXATTR                = 172
	SYS_LGETXATTR               = 173
	SYS_GETDENTS                = 174
	SYS_SETSID                  = 175
	SYS_FCHDIR                  = 176
	SYS_FGETXATTR               = 177
	SYS_LISTXATTR               = 178
	SYS_LLISTXATTR              = 179
	SYS_FLISTXATTR              = 180
	SYS_REMOVEXATTR             = 181
	SYS_LREMOVEXATTR            = 182
	SYS_SIGPENDING              = 183
	SYS_QUERY_MODULE            = 184
	SYS_SETPGID                 = 185
	SYS_FREMOVEXATTR            = 186
	SYS_TKILL                   = 187
	SYS_EXIT_GROUP              = 188
	SYS_UNAME                   = 189
	SYS_INIT_MODULE             = 190
	SYS_PERSONALITY             = 191
	SYS_REMAP_FILE_PAGES        = 192
	SYS_EPOLL_CREATE            = 193
	SYS_EPOLL_CTL               = 194
	SYS_EPOLL_WAIT              = 195
	SYS_IOPRIO_SET              = 196
	SYS_GETPPID                 = 197
	SYS_SIGACTION               = 198
	SYS_SGETMASK                = 199
	SYS_SSETMASK                = 200
	SYS_SIGSUSPEND              = 201
	SYS_OLDLSTAT                = 202
	SYS_USELIB                  = 203
	SYS_READDIR                 = 204
	SYS_READAHEAD               = 205
	SYS_SOCKETCALL              = 206
	SYS_SYSLOG                  = 207
	SYS_LOOKUP_DCOOKIE          = 208
	SYS_FADVISE64               = 209
	SYS_FADVISE64_64            = 210
	SYS_TGKILL                  = 211
	SYS_WAITPID                 = 212
	SYS_SWAPOFF                 = 213
	SYS_SYSINFO                 = 214
	SYS_IPC                     = 215
	SYS_SIGRETURN               = 216
	SYS_CLONE                   = 217
	SYS_IOPRIO_GET              = 218
	SYS_ADJTIMEX                = 219
	SYS_SIGPROCMASK             = 220
	SYS_CREATE_MODULE           = 221
	SYS_DELETE_MODULE           = 222
	SYS_GET_KERNEL_SYMS         = 223
	SYS_GETPGID                 = 224
	SYS_BDFLUSH                 = 225
	SYS_SYSFS                   = 226
	SYS_AFS_SYSCALL             = 227
	SYS_SETFSUID                = 228
	SYS_SETFSGID                = 229
	SYS__NEWSELECT              = 230
	SYS_SPLICE                  = 232
	SYS_STIME                   = 233
	SYS_STATFS64                = 234
	SYS_FSTATFS64               = 235
	SYS__LLSEEK                 = 236
	SYS_MLOCK                   = 237
	SYS_MUNLOCK                 = 238
	SYS_MLOCKALL                = 239
	SYS_MUNLOCKALL              = 240
	SYS_SCHED_SETPARAM          = 241
	SYS_SCHED_GETPARAM          = 242
	SYS_SCHED_SETSCHEDULER      = 243
	SYS_SCHED_GETSCHEDULER      = 244
	SYS_SCHED_YIELD             = 245
	SYS_SCHED_GET_PRIORITY_MAX  = 246
	SYS_SCHED_GET_PRIORITY_MIN  = 247
	SYS_SCHED_RR_GET_INTERVAL   = 248
	SYS_NANOSLEEP               = 249
	SYS_MREMAP                  = 250
	SYS__SYSCTL                 = 251
	SYS_GETSID                  = 252
	SYS_FDATASYNC               = 253
	SYS_NFSSERVCTL              = 254
	SYS_SYNC_FILE_RANGE         = 255
	SYS_CLOCK_SETTIME           = 256
	SYS_CLOCK_GETTIME           = 257
	SYS_CLOCK_GETRES            = 258
	SYS_CLOCK_NANOSLEEP         = 259
	SYS_SCHED_GETAFFINITY       = 260
	SYS_SCHED_SETAFFINITY       = 261
	SYS_TIMER_SETTIME           = 262
	SYS_TIMER_GETTIME           = 263
	SYS_TIMER_GETOVERRUN        = 264
	SYS_TIMER_DELETE            = 265
	SYS_TIMER_CREATE            = 266
	SYS_VSERVER                 = 267
	SYS_IO_SETUP                = 268
	SYS_IO_DESTROY              = 269
	SYS_IO_SUBMIT               = 270
	SYS_IO_CANCEL               = 271
	SYS_IO_GETEVENTS            = 272
	SYS_MQ_OPEN                 = 273
	SYS_MQ_UNLINK               = 274
	SYS_MQ_TIMEDSEND            = 275
	SYS_MQ_TIMEDRECEIVE         = 276
	SYS_MQ_NOTIFY               = 277
	SYS_MQ_GETSETATTR           = 278
	SYS_WAITID                  = 279
	SYS_TEE                     = 280
	SYS_ADD_KEY                 = 281
	SYS_REQUEST_KEY             = 282
	SYS_KEYCTL                  = 283
	SYS_OPENAT                  = 284
	SYS_MKDIRAT                 = 285
	SYS_MKNODAT                 = 286
	SYS_FCHOWNAT                = 287
	SYS_FUTIMESAT               = 288
	SYS_FSTATAT64               = 289
	SYS_UNLINKAT                = 290
	SYS_RENAMEAT                = 291
	SYS_LINKAT                  = 292
	SYS_SYMLINKAT               = 293
	SYS_READLINKAT              = 294
	SYS_FCHMODAT                = 295
	SYS_FACCESSAT               = 296
	SYS_PSELECT6                = 297
	SYS_PPOLL                   = 298
	SYS_UNSHARE                 = 299
	SYS_SET_ROBUST_LIST         = 300
	SYS_GET_ROBUST_LIST         = 301
	SYS_MIGRATE_PAGES           = 302
	SYS_MBIND                   = 303
	SYS_GET_MEMPOLICY           = 304
	SYS_SET_MEMPOLICY           = 305
	SYS_KEXEC_LOAD              = 306
	SYS_MOVE_PAGES              = 307
	SYS_GETCPU                  = 308
	SYS_EPOLL_PWAIT             = 309
	SYS_UTIMENSAT               = 310
	SYS_SIGNALFD                = 311
	SYS_TIMERFD_CREATE          = 312
	SYS_EVENTFD                 = 313
	SYS_FALLOCATE               = 314
	SYS_TIMERFD_SETTIME         = 315
	SYS_TIMERFD_GETTIME         = 316
	SYS_SIGNALFD4               = 317
	SYS_EVENTFD2                = 318
	SYS_EPOLL_CREATE1           = 319
	SYS_DUP3                    = 320
	SYS_PIPE2                   = 321
	SYS_INOTIFY_INIT1           = 322
	SYS_ACCEPT4                 = 323
	SYS_PREADV                  = 324
	SYS_PWRITEV                 = 325
	SYS_RT_TGSIGQUEUEINFO       = 326
	SYS_PERF_EVENT_OPEN         = 327
	SYS_RECVMMSG                = 328
	SYS_FANOTIFY_INIT           = 329
	SYS_FANOTIFY_MARK           = 330
	SYS_PRLIMIT64               = 331
	SYS_NAME_TO_HANDLE_AT       = 332
	SYS_OPEN_BY_HANDLE_AT       = 333
	SYS_CLOCK_ADJTIME           = 334
	SYS_SYNCFS                  = 335
	SYS_SENDMMSG                = 336
	SYS_SETNS                   = 337
	SYS_PROCESS_VM_READV        = 338
	SYS_PROCESS_VM_WRITEV       = 339
	SYS_KERN_FEATURES           = 340
	SYS_KCMP                    = 341
	SYS_FINIT_MODULE            = 342
	SYS_SCHED_SETATTR           = 343
	SYS_SCHED_GETATTR           = 344
	SYS_RENAMEAT2               = 345
	SYS_SECCOMP                 = 346
	SYS_GETRANDOM               = 347
	SYS_MEMFD_CREATE            = 348
	SYS_BPF                     = 349
	SYS_EXECVEAT                = 350
	SYS_MEMBARRIER              = 351
	SYS_USERFAULTFD             = 352
	SYS_BIND                    = 353
	SYS_LISTEN                  = 354
	SYS_SETSOCKOPT              = 355
	SYS_MLOCK2                  = 356
	SYS_COPY_FILE_RANGE         = 357
	SYS_PREADV2                 = 358
	SYS_PWRITEV2                = 359
	SYS_STATX                   = 360
	SYS_IO_PGETEVENTS           = 361
	SYS_PKEY_MPROTECT           = 362
	SYS_PKEY_ALLOC              = 363
	SYS_PKEY_FREE               = 364
	SYS_RSEQ                    = 365
	SYS_SEMTIMEDOP              = 392
	SYS_SEMGET                  = 393
	SYS_SEMCTL                  = 394
	SYS_SHMGET                  = 395
	SYS_SHMCTL                  = 396
	SYS_SHMAT                   = 397
	SYS_SHMDT                   = 398
	SYS_MSGGET                  = 399
	SYS_MSGSND                  = 400
	SYS_MSGRCV                  = 401
	SYS_MSGCTL                  = 402
	SYS_PIDFD_SEND_SIGNAL       = 424
	SYS_IO_URING_SETUP          = 425
	SYS_IO_URING_ENTER          = 426
	SYS_IO_URING_REGISTER       = 427
	SYS_OPEN_TREE               = 428
	SYS_MOVE_MOUNT              = 429
	SYS_FSOPEN                  = 430
	SYS_FSCONFIG                = 431
	SYS_FSMOUNT                 = 432
	SYS_FSPICK                  = 433
	SYS_PIDFD_OPEN              = 434
	SYS_CLOSE_RANGE             = 436
	SYS_OPENAT2                 = 437
	SYS_PIDFD_GETFD             = 438
	SYS_FACCESSAT2              = 439
	SYS_PROCESS_MADVISE         = 440
	SYS_EPOLL_PWAIT2            = 441
	SYS_MOUNT_SETATTR           = 442
	SYS_QUOTACTL_FD             = 443
	SYS_LANDLOCK_CREATE_RULESET = 444
	SYS_LANDLOCK_ADD_RULE       = 445
	SYS_LANDLOCK_RESTRICT_SELF  = 446
	SYS_PROCESS_MRELEASE        = 448
	SYS_FUTEX_WAITV             = 449
	SYS_SET_MEMPOLICY_HOME_NODE = 450
	SYS_CACHESTAT               = 451
)
