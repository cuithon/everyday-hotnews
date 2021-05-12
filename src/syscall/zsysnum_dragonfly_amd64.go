// mksysnum_dragonfly.pl
// Code generated by the command above; DO NOT EDIT.

//go:build amd64 && dragonfly
// +build amd64,dragonfly

package syscall

const (
	// SYS_NOSYS = 0;  // { int nosys(void); } syscall nosys_args int
	SYS_EXIT          = 1   // { void exit(int rval); }
	SYS_FORK          = 2   // { int fork(void); }
	SYS_READ          = 3   // { ssize_t read(int fd, void *buf, size_t nbyte); }
	SYS_WRITE         = 4   // { ssize_t write(int fd, const void *buf, size_t nbyte); }
	SYS_OPEN          = 5   // { int open(char *path, int flags, int mode); }
	SYS_CLOSE         = 6   // { int close(int fd); }
	SYS_WAIT4         = 7   // { int wait4(int pid, int *status, int options, \
	SYS_LINK          = 9   // { int link(char *path, char *link); }
	SYS_UNLINK        = 10  // { int unlink(char *path); }
	SYS_CHDIR         = 12  // { int chdir(char *path); }
	SYS_FCHDIR        = 13  // { int fchdir(int fd); }
	SYS_MKNOD         = 14  // { int mknod(char *path, int mode, int dev); }
	SYS_CHMOD         = 15  // { int chmod(char *path, int mode); }
	SYS_CHOWN         = 16  // { int chown(char *path, int uid, int gid); }
	SYS_OBREAK        = 17  // { int obreak(char *nsize); } break obreak_args int
	SYS_GETFSSTAT     = 18  // { int getfsstat(struct statfs *buf, long bufsize, \
	SYS_GETPID        = 20  // { pid_t getpid(void); }
	SYS_MOUNT         = 21  // { int mount(char *type, char *path, int flags, \
	SYS_UNMOUNT       = 22  // { int unmount(char *path, int flags); }
	SYS_SETUID        = 23  // { int setuid(uid_t uid); }
	SYS_GETUID        = 24  // { uid_t getuid(void); }
	SYS_GETEUID       = 25  // { uid_t geteuid(void); }
	SYS_PTRACE        = 26  // { int ptrace(int req, pid_t pid, caddr_t addr, \
	SYS_RECVMSG       = 27  // { int recvmsg(int s, struct msghdr *msg, int flags); }
	SYS_SENDMSG       = 28  // { int sendmsg(int s, caddr_t msg, int flags); }
	SYS_RECVFROM      = 29  // { int recvfrom(int s, caddr_t buf, size_t len, \
	SYS_ACCEPT        = 30  // { int accept(int s, caddr_t name, int *anamelen); }
	SYS_GETPEERNAME   = 31  // { int getpeername(int fdes, caddr_t asa, int *alen); }
	SYS_GETSOCKNAME   = 32  // { int getsockname(int fdes, caddr_t asa, int *alen); }
	SYS_ACCESS        = 33  // { int access(char *path, int flags); }
	SYS_CHFLAGS       = 34  // { int chflags(char *path, int flags); }
	SYS_FCHFLAGS      = 35  // { int fchflags(int fd, int flags); }
	SYS_SYNC          = 36  // { int sync(void); }
	SYS_KILL          = 37  // { int kill(int pid, int signum); }
	SYS_GETPPID       = 39  // { pid_t getppid(void); }
	SYS_DUP           = 41  // { int dup(u_int fd); }
	SYS_PIPE          = 42  // { int pipe(void); }
	SYS_GETEGID       = 43  // { gid_t getegid(void); }
	SYS_PROFIL        = 44  // { int profil(caddr_t samples, size_t size, \
	SYS_KTRACE        = 45  // { int ktrace(const char *fname, int ops, int facs, \
	SYS_GETGID        = 47  // { gid_t getgid(void); }
	SYS_GETLOGIN      = 49  // { int getlogin(char *namebuf, u_int namelen); }
	SYS_SETLOGIN      = 50  // { int setlogin(char *namebuf); }
	SYS_ACCT          = 51  // { int acct(char *path); }
	SYS_SIGALTSTACK   = 53  // { int sigaltstack(stack_t *ss, stack_t *oss); }
	SYS_IOCTL         = 54  // { int ioctl(int fd, u_long com, caddr_t data); }
	SYS_REBOOT        = 55  // { int reboot(int opt); }
	SYS_REVOKE        = 56  // { int revoke(char *path); }
	SYS_SYMLINK       = 57  // { int symlink(char *path, char *link); }
	SYS_READLINK      = 58  // { int readlink(char *path, char *buf, int count); }
	SYS_EXECVE        = 59  // { int execve(char *fname, char **argv, char **envv); }
	SYS_UMASK         = 60  // { int umask(int newmask); } umask umask_args int
	SYS_CHROOT        = 61  // { int chroot(char *path); }
	SYS_MSYNC         = 65  // { int msync(void *addr, size_t len, int flags); }
	SYS_VFORK         = 66  // { pid_t vfork(void); }
	SYS_SBRK          = 69  // { int sbrk(int incr); }
	SYS_SSTK          = 70  // { int sstk(int incr); }
	SYS_MUNMAP        = 73  // { int munmap(void *addr, size_t len); }
	SYS_MPROTECT      = 74  // { int mprotect(void *addr, size_t len, int prot); }
	SYS_MADVISE       = 75  // { int madvise(void *addr, size_t len, int behav); }
	SYS_MINCORE       = 78  // { int mincore(const void *addr, size_t len, \
	SYS_GETGROUPS     = 79  // { int getgroups(u_int gidsetsize, gid_t *gidset); }
	SYS_SETGROUPS     = 80  // { int setgroups(u_int gidsetsize, gid_t *gidset); }
	SYS_GETPGRP       = 81  // { int getpgrp(void); }
	SYS_SETPGID       = 82  // { int setpgid(int pid, int pgid); }
	SYS_SETITIMER     = 83  // { int setitimer(u_int which, struct itimerval *itv, \
	SYS_SWAPON        = 85  // { int swapon(char *name); }
	SYS_GETITIMER     = 86  // { int getitimer(u_int which, struct itimerval *itv); }
	SYS_GETDTABLESIZE = 89  // { int getdtablesize(void); }
	SYS_DUP2          = 90  // { int dup2(u_int from, u_int to); }
	SYS_FCNTL         = 92  // { int fcntl(int fd, int cmd, long arg); }
	SYS_SELECT        = 93  // { int select(int nd, fd_set *in, fd_set *ou, \
	SYS_FSYNC         = 95  // { int fsync(int fd); }
	SYS_SETPRIORITY   = 96  // { int setpriority(int which, int who, int prio); }
	SYS_SOCKET        = 97  // { int socket(int domain, int type, int protocol); }
	SYS_CONNECT       = 98  // { int connect(int s, caddr_t name, int namelen); }
	SYS_GETPRIORITY   = 100 // { int getpriority(int which, int who); }
	SYS_BIND          = 104 // { int bind(int s, caddr_t name, int namelen); }
	SYS_SETSOCKOPT    = 105 // { int setsockopt(int s, int level, int name, \
	SYS_LISTEN        = 106 // { int listen(int s, int backlog); }
	SYS_GETTIMEOFDAY  = 116 // { int gettimeofday(struct timeval *tp, \
	SYS_GETRUSAGE     = 117 // { int getrusage(int who, struct rusage *rusage); }
	SYS_GETSOCKOPT    = 118 // { int getsockopt(int s, int level, int name, \
	SYS_READV         = 120 // { int readv(int fd, struct iovec *iovp, u_int iovcnt); }
	SYS_WRITEV        = 121 // { int writev(int fd, struct iovec *iovp, \
	SYS_SETTIMEOFDAY  = 122 // { int settimeofday(struct timeval *tv, \
	SYS_FCHOWN        = 123 // { int fchown(int fd, int uid, int gid); }
	SYS_FCHMOD        = 124 // { int fchmod(int fd, int mode); }
	SYS_SETREUID      = 126 // { int setreuid(int ruid, int euid); }
	SYS_SETREGID      = 127 // { int setregid(int rgid, int egid); }
	SYS_RENAME        = 128 // { int rename(char *from, char *to); }
	SYS_FLOCK         = 131 // { int flock(int fd, int how); }
	SYS_MKFIFO        = 132 // { int mkfifo(char *path, int mode); }
	SYS_SENDTO        = 133 // { int sendto(int s, caddr_t buf, size_t len, \
	SYS_SHUTDOWN      = 134 // { int shutdown(int s, int how); }
	SYS_SOCKETPAIR    = 135 // { int socketpair(int domain, int type, int protocol, \
	SYS_MKDIR         = 136 // { int mkdir(char *path, int mode); }
	SYS_RMDIR         = 137 // { int rmdir(char *path); }
	SYS_UTIMES        = 138 // { int utimes(char *path, struct timeval *tptr); }
	SYS_ADJTIME       = 140 // { int adjtime(struct timeval *delta, \
	SYS_SETSID        = 147 // { int setsid(void); }
	SYS_QUOTACTL      = 148 // { int quotactl(char *path, int cmd, int uid, \
	SYS_STATFS        = 157 // { int statfs(char *path, struct statfs *buf); }
	SYS_FSTATFS       = 158 // { int fstatfs(int fd, struct statfs *buf); }
	SYS_GETFH         = 161 // { int getfh(char *fname, struct fhandle *fhp); }
	SYS_GETDOMAINNAME = 162 // { int getdomainname(char *domainname, int len); }
	SYS_SETDOMAINNAME = 163 // { int setdomainname(char *domainname, int len); }
	SYS_UNAME         = 164 // { int uname(struct utsname *name); }
	SYS_SYSARCH       = 165 // { int sysarch(int op, char *parms); }
	SYS_RTPRIO        = 166 // { int rtprio(int function, pid_t pid, \
	SYS_EXTPREAD      = 173 // { ssize_t extpread(int fd, void *buf, \
	SYS_EXTPWRITE     = 174 // { ssize_t extpwrite(int fd, const void *buf, \
	SYS_NTP_ADJTIME   = 176 // { int ntp_adjtime(struct timex *tp); }
	SYS_SETGID        = 181 // { int setgid(gid_t gid); }
	SYS_SETEGID       = 182 // { int setegid(gid_t egid); }
	SYS_SETEUID       = 183 // { int seteuid(uid_t euid); }
	SYS_PATHCONF      = 191 // { int pathconf(char *path, int name); }
	SYS_FPATHCONF     = 192 // { int fpathconf(int fd, int name); }
	SYS_GETRLIMIT     = 194 // { int getrlimit(u_int which, \
	SYS_SETRLIMIT     = 195 // { int setrlimit(u_int which, \
	SYS_MMAP          = 197 // { caddr_t mmap(caddr_t addr, size_t len, int prot, \
	// SYS_NOSYS = 198;  // { int nosys(void); } __syscall __syscall_args int
	SYS_LSEEK                  = 199 // { off_t lseek(int fd, int pad, off_t offset, \
	SYS_TRUNCATE               = 200 // { int truncate(char *path, int pad, off_t length); }
	SYS_FTRUNCATE              = 201 // { int ftruncate(int fd, int pad, off_t length); }
	SYS___SYSCTL               = 202 // { int __sysctl(int *name, u_int namelen, void *old, \
	SYS_MLOCK                  = 203 // { int mlock(const void *addr, size_t len); }
	SYS_MUNLOCK                = 204 // { int munlock(const void *addr, size_t len); }
	SYS_UNDELETE               = 205 // { int undelete(char *path); }
	SYS_FUTIMES                = 206 // { int futimes(int fd, struct timeval *tptr); }
	SYS_GETPGID                = 207 // { int getpgid(pid_t pid); }
	SYS_POLL                   = 209 // { int poll(struct pollfd *fds, u_int nfds, \
	SYS___SEMCTL               = 220 // { int __semctl(int semid, int semnum, int cmd, \
	SYS_SEMGET                 = 221 // { int semget(key_t key, int nsems, int semflg); }
	SYS_SEMOP                  = 222 // { int semop(int semid, struct sembuf *sops, \
	SYS_MSGCTL                 = 224 // { int msgctl(int msqid, int cmd, \
	SYS_MSGGET                 = 225 // { int msgget(key_t key, int msgflg); }
	SYS_MSGSND                 = 226 // { int msgsnd(int msqid, void *msgp, size_t msgsz, \
	SYS_MSGRCV                 = 227 // { int msgrcv(int msqid, void *msgp, size_t msgsz, \
	SYS_SHMAT                  = 228 // { caddr_t shmat(int shmid, const void *shmaddr, \
	SYS_SHMCTL                 = 229 // { int shmctl(int shmid, int cmd, \
	SYS_SHMDT                  = 230 // { int shmdt(const void *shmaddr); }
	SYS_SHMGET                 = 231 // { int shmget(key_t key, size_t size, int shmflg); }
	SYS_CLOCK_GETTIME          = 232 // { int clock_gettime(clockid_t clock_id, \
	SYS_CLOCK_SETTIME          = 233 // { int clock_settime(clockid_t clock_id, \
	SYS_CLOCK_GETRES           = 234 // { int clock_getres(clockid_t clock_id, \
	SYS_NANOSLEEP              = 240 // { int nanosleep(const struct timespec *rqtp, \
	SYS_MINHERIT               = 250 // { int minherit(void *addr, size_t len, int inherit); }
	SYS_RFORK                  = 251 // { int rfork(int flags); }
	SYS_OPENBSD_POLL           = 252 // { int openbsd_poll(struct pollfd *fds, u_int nfds, \
	SYS_ISSETUGID              = 253 // { int issetugid(void); }
	SYS_LCHOWN                 = 254 // { int lchown(char *path, int uid, int gid); }
	SYS_LCHMOD                 = 274 // { int lchmod(char *path, mode_t mode); }
	SYS_LUTIMES                = 276 // { int lutimes(char *path, struct timeval *tptr); }
	SYS_EXTPREADV              = 289 // { ssize_t extpreadv(int fd, struct iovec *iovp, \
	SYS_EXTPWRITEV             = 290 // { ssize_t extpwritev(int fd, struct iovec *iovp,\
	SYS_FHSTATFS               = 297 // { int fhstatfs(const struct fhandle *u_fhp, struct statfs *buf); }
	SYS_FHOPEN                 = 298 // { int fhopen(const struct fhandle *u_fhp, int flags); }
	SYS_MODNEXT                = 300 // { int modnext(int modid); }
	SYS_MODSTAT                = 301 // { int modstat(int modid, struct module_stat* stat); }
	SYS_MODFNEXT               = 302 // { int modfnext(int modid); }
	SYS_MODFIND                = 303 // { int modfind(const char *name); }
	SYS_KLDLOAD                = 304 // { int kldload(const char *file); }
	SYS_KLDUNLOAD              = 305 // { int kldunload(int fileid); }
	SYS_KLDFIND                = 306 // { int kldfind(const char *file); }
	SYS_KLDNEXT                = 307 // { int kldnext(int fileid); }
	SYS_KLDSTAT                = 308 // { int kldstat(int fileid, struct kld_file_stat* stat); }
	SYS_KLDFIRSTMOD            = 309 // { int kldfirstmod(int fileid); }
	SYS_GETSID                 = 310 // { int getsid(pid_t pid); }
	SYS_SETRESUID              = 311 // { int setresuid(uid_t ruid, uid_t euid, uid_t suid); }
	SYS_SETRESGID              = 312 // { int setresgid(gid_t rgid, gid_t egid, gid_t sgid); }
	SYS_AIO_RETURN             = 314 // { int aio_return(struct aiocb *aiocbp); }
	SYS_AIO_SUSPEND            = 315 // { int aio_suspend(struct aiocb * const * aiocbp, int nent, const struct timespec *timeout); }
	SYS_AIO_CANCEL             = 316 // { int aio_cancel(int fd, struct aiocb *aiocbp); }
	SYS_AIO_ERROR              = 317 // { int aio_error(struct aiocb *aiocbp); }
	SYS_AIO_READ               = 318 // { int aio_read(struct aiocb *aiocbp); }
	SYS_AIO_WRITE              = 319 // { int aio_write(struct aiocb *aiocbp); }
	SYS_LIO_LISTIO             = 320 // { int lio_listio(int mode, struct aiocb * const *acb_list, int nent, struct sigevent *sig); }
	SYS_YIELD                  = 321 // { int yield(void); }
	SYS_MLOCKALL               = 324 // { int mlockall(int how); }
	SYS_MUNLOCKALL             = 325 // { int munlockall(void); }
	SYS___GETCWD               = 326 // { int __getcwd(u_char *buf, u_int buflen); }
	SYS_SCHED_SETPARAM         = 327 // { int sched_setparam (pid_t pid, const struct sched_param *param); }
	SYS_SCHED_GETPARAM         = 328 // { int sched_getparam (pid_t pid, struct sched_param *param); }
	SYS_SCHED_SETSCHEDULER     = 329 // { int sched_setscheduler (pid_t pid, int policy, const struct sched_param *param); }
	SYS_SCHED_GETSCHEDULER     = 330 // { int sched_getscheduler (pid_t pid); }
	SYS_SCHED_YIELD            = 331 // { int sched_yield (void); }
	SYS_SCHED_GET_PRIORITY_MAX = 332 // { int sched_get_priority_max (int policy); }
	SYS_SCHED_GET_PRIORITY_MIN = 333 // { int sched_get_priority_min (int policy); }
	SYS_SCHED_RR_GET_INTERVAL  = 334 // { int sched_rr_get_interval (pid_t pid, struct timespec *interval); }
	SYS_UTRACE                 = 335 // { int utrace(const void *addr, size_t len); }
	SYS_KLDSYM                 = 337 // { int kldsym(int fileid, int cmd, void *data); }
	SYS_JAIL                   = 338 // { int jail(struct jail *jail); }
	SYS_SIGPROCMASK            = 340 // { int sigprocmask(int how, const sigset_t *set, \
	SYS_SIGSUSPEND             = 341 // { int sigsuspend(const sigset_t *sigmask); }
	SYS_SIGACTION              = 342 // { int sigaction(int sig, const struct sigaction *act, \
	SYS_SIGPENDING             = 343 // { int sigpending(sigset_t *set); }
	SYS_SIGRETURN              = 344 // { int sigreturn(ucontext_t *sigcntxp); }
	SYS_SIGTIMEDWAIT           = 345 // { int sigtimedwait(const sigset_t *set,\
	SYS_SIGWAITINFO            = 346 // { int sigwaitinfo(const sigset_t *set,\
	SYS___ACL_GET_FILE         = 347 // { int __acl_get_file(const char *path, \
	SYS___ACL_SET_FILE         = 348 // { int __acl_set_file(const char *path, \
	SYS___ACL_GET_FD           = 349 // { int __acl_get_fd(int filedes, acl_type_t type, \
	SYS___ACL_SET_FD           = 350 // { int __acl_set_fd(int filedes, acl_type_t type, \
	SYS___ACL_DELETE_FILE      = 351 // { int __acl_delete_file(const char *path, \
	SYS___ACL_DELETE_FD        = 352 // { int __acl_delete_fd(int filedes, acl_type_t type); }
	SYS___ACL_ACLCHECK_FILE    = 353 // { int __acl_aclcheck_file(const char *path, \
	SYS___ACL_ACLCHECK_FD      = 354 // { int __acl_aclcheck_fd(int filedes, acl_type_t type, \
	SYS_EXTATTRCTL             = 355 // { int extattrctl(const char *path, int cmd, \
	SYS_EXTATTR_SET_FILE       = 356 // { int extattr_set_file(const char *path, \
	SYS_EXTATTR_GET_FILE       = 357 // { int extattr_get_file(const char *path, \
	SYS_EXTATTR_DELETE_FILE    = 358 // { int extattr_delete_file(const char *path, \
	SYS_AIO_WAITCOMPLETE       = 359 // { int aio_waitcomplete(struct aiocb **aiocbp, struct timespec *timeout); }
	SYS_GETRESUID              = 360 // { int getresuid(uid_t *ruid, uid_t *euid, uid_t *suid); }
	SYS_GETRESGID              = 361 // { int getresgid(gid_t *rgid, gid_t *egid, gid_t *sgid); }
	SYS_KQUEUE                 = 362 // { int kqueue(void); }
	SYS_KEVENT                 = 363 // { int kevent(int fd, \
	SYS_SCTP_PEELOFF           = 364 // { int sctp_peeloff(int sd, caddr_t name ); }
	SYS_LCHFLAGS               = 391 // { int lchflags(char *path, int flags); }
	SYS_UUIDGEN                = 392 // { int uuidgen(struct uuid *store, int count); }
	SYS_SENDFILE               = 393 // { int sendfile(int fd, int s, off_t offset, size_t nbytes, \
	SYS_VARSYM_SET             = 450 // { int varsym_set(int level, const char *name, const char *data); }
	SYS_VARSYM_GET             = 451 // { int varsym_get(int mask, const char *wild, char *buf, int bufsize); }
	SYS_VARSYM_LIST            = 452 // { int varsym_list(int level, char *buf, int maxsize, int *marker); }
	SYS_EXEC_SYS_REGISTER      = 465 // { int exec_sys_register(void *entry); }
	SYS_EXEC_SYS_UNREGISTER    = 466 // { int exec_sys_unregister(int id); }
	SYS_SYS_CHECKPOINT         = 467 // { int sys_checkpoint(int type, int fd, pid_t pid, int retval); }
	SYS_MOUNTCTL               = 468 // { int mountctl(const char *path, int op, int fd, const void *ctl, int ctllen, void *buf, int buflen); }
	SYS_UMTX_SLEEP             = 469 // { int umtx_sleep(volatile const int *ptr, int value, int timeout); }
	SYS_UMTX_WAKEUP            = 470 // { int umtx_wakeup(volatile const int *ptr, int count); }
	SYS_JAIL_ATTACH            = 471 // { int jail_attach(int jid); }
	SYS_SET_TLS_AREA           = 472 // { int set_tls_area(int which, struct tls_info *info, size_t infosize); }
	SYS_GET_TLS_AREA           = 473 // { int get_tls_area(int which, struct tls_info *info, size_t infosize); }
	SYS_CLOSEFROM              = 474 // { int closefrom(int fd); }
	SYS_STAT                   = 475 // { int stat(const char *path, struct stat *ub); }
	SYS_FSTAT                  = 476 // { int fstat(int fd, struct stat *sb); }
	SYS_LSTAT                  = 477 // { int lstat(const char *path, struct stat *ub); }
	SYS_FHSTAT                 = 478 // { int fhstat(const struct fhandle *u_fhp, struct stat *sb); }
	SYS_GETDIRENTRIES          = 479 // { int getdirentries(int fd, char *buf, u_int count, \
	SYS_GETDENTS               = 480 // { int getdents(int fd, char *buf, size_t count); }
	SYS_USCHED_SET             = 481 // { int usched_set(pid_t pid, int cmd, void *data, \
	SYS_EXTACCEPT              = 482 // { int extaccept(int s, int flags, caddr_t name, int *anamelen); }
	SYS_EXTCONNECT             = 483 // { int extconnect(int s, int flags, caddr_t name, int namelen); }
	SYS_MCONTROL               = 485 // { int mcontrol(void *addr, size_t len, int behav, off_t value); }
	SYS_VMSPACE_CREATE         = 486 // { int vmspace_create(void *id, int type, void *data); }
	SYS_VMSPACE_DESTROY        = 487 // { int vmspace_destroy(void *id); }
	SYS_VMSPACE_CTL            = 488 // { int vmspace_ctl(void *id, int cmd, 		\
	SYS_VMSPACE_MMAP           = 489 // { int vmspace_mmap(void *id, void *addr, size_t len, \
	SYS_VMSPACE_MUNMAP         = 490 // { int vmspace_munmap(void *id, void *addr,	\
	SYS_VMSPACE_MCONTROL       = 491 // { int vmspace_mcontrol(void *id, void *addr, 	\
	SYS_VMSPACE_PREAD          = 492 // { ssize_t vmspace_pread(void *id, void *buf, \
	SYS_VMSPACE_PWRITE         = 493 // { ssize_t vmspace_pwrite(void *id, const void *buf, \
	SYS_EXTEXIT                = 494 // { void extexit(int how, int status, void *addr); }
	SYS_LWP_CREATE             = 495 // { int lwp_create(struct lwp_params *params); }
	SYS_LWP_GETTID             = 496 // { lwpid_t lwp_gettid(void); }
	SYS_LWP_KILL               = 497 // { int lwp_kill(pid_t pid, lwpid_t tid, int signum); }
	SYS_LWP_RTPRIO             = 498 // { int lwp_rtprio(int function, pid_t pid, lwpid_t tid, struct rtprio *rtp); }
	SYS_PSELECT                = 499 // { int pselect(int nd, fd_set *in, fd_set *ou, \
	SYS_STATVFS                = 500 // { int statvfs(const char *path, struct statvfs *buf); }
	SYS_FSTATVFS               = 501 // { int fstatvfs(int fd, struct statvfs *buf); }
	SYS_FHSTATVFS              = 502 // { int fhstatvfs(const struct fhandle *u_fhp, struct statvfs *buf); }
	SYS_GETVFSSTAT             = 503 // { int getvfsstat(struct statfs *buf,          \
	SYS_OPENAT                 = 504 // { int openat(int fd, char *path, int flags, int mode); }
	SYS_FSTATAT                = 505 // { int fstatat(int fd, char *path, 	\
	SYS_FCHMODAT               = 506 // { int fchmodat(int fd, char *path, int mode, \
	SYS_FCHOWNAT               = 507 // { int fchownat(int fd, char *path, int uid, int gid, \
	SYS_UNLINKAT               = 508 // { int unlinkat(int fd, char *path, int flags); }
	SYS_FACCESSAT              = 509 // { int faccessat(int fd, char *path, int amode, \
	SYS_MQ_OPEN                = 510 // { mqd_t mq_open(const char * name, int oflag, \
	SYS_MQ_CLOSE               = 511 // { int mq_close(mqd_t mqdes); }
	SYS_MQ_UNLINK              = 512 // { int mq_unlink(const char *name); }
	SYS_MQ_GETATTR             = 513 // { int mq_getattr(mqd_t mqdes, \
	SYS_MQ_SETATTR             = 514 // { int mq_setattr(mqd_t mqdes, \
	SYS_MQ_NOTIFY              = 515 // { int mq_notify(mqd_t mqdes, \
	SYS_MQ_SEND                = 516 // { int mq_send(mqd_t mqdes, const char *msg_ptr, \
	SYS_MQ_RECEIVE             = 517 // { ssize_t mq_receive(mqd_t mqdes, char *msg_ptr, \
	SYS_MQ_TIMEDSEND           = 518 // { int mq_timedsend(mqd_t mqdes, \
	SYS_MQ_TIMEDRECEIVE        = 519 // { ssize_t mq_timedreceive(mqd_t mqdes, \
	SYS_IOPRIO_SET             = 520 // { int ioprio_set(int which, int who, int prio); }
	SYS_IOPRIO_GET             = 521 // { int ioprio_get(int which, int who); }
	SYS_CHROOT_KERNEL          = 522 // { int chroot_kernel(char *path); }
	SYS_RENAMEAT               = 523 // { int renameat(int oldfd, char *old, int newfd, \
	SYS_MKDIRAT                = 524 // { int mkdirat(int fd, char *path, mode_t mode); }
	SYS_MKFIFOAT               = 525 // { int mkfifoat(int fd, char *path, mode_t mode); }
	SYS_MKNODAT                = 526 // { int mknodat(int fd, char *path, mode_t mode, \
	SYS_READLINKAT             = 527 // { int readlinkat(int fd, char *path, char *buf, \
	SYS_SYMLINKAT              = 528 // { int symlinkat(char *path1, int fd, char *path2); }
	SYS_SWAPOFF                = 529 // { int swapoff(char *name); }
	SYS_VQUOTACTL              = 530 // { int vquotactl(const char *path, \
	SYS_LINKAT                 = 531 // { int linkat(int fd1, char *path1, int fd2, \
	SYS_EACCESS                = 532 // { int eaccess(char *path, int flags); }
	SYS_LPATHCONF              = 533 // { int lpathconf(char *path, int name); }
	SYS_VMM_GUEST_CTL          = 534 // { int vmm_guest_ctl(int op, struct vmm_guest_options *options); }
	SYS_VMM_GUEST_SYNC_ADDR    = 535 // { int vmm_guest_sync_addr(long *dstaddr, long *srcaddr); }
	SYS_PIPE2                  = 538 // { int pipe2(int *fildes, int flags); }
	SYS_UTIMENSAT              = 539 // { int utimensat(int fd, const char *path, const struct timespec *ts, int flags); }
	SYS_ACCEPT4                = 541 // { int accept4(int s, caddr_t name, int *anamelen, int flags); }
	SYS_WAIT6                  = 548 // { int wait6(idtype_t idtype, id_t id, int *status, int options, struct __wrusage *wrusage, siginfo_t *info); }
)
