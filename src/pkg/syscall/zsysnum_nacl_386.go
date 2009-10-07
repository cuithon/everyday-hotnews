// mksysnum_nacl.sh /home/rsc/pub/nacl/native_client/src/trusted/service_runtime/include/bits/nacl_syscalls.h
// MACHINE GENERATED BY THE ABOVE COMMAND; DO NOT EDIT

package syscall

const (
	SYS_NULL		= 1;
	SYS_OPEN		= 10;
	SYS_CLOSE		= 11;
	SYS_READ		= 12;
	SYS_WRITE		= 13;
	SYS_LSEEK		= 14;
	SYS_IOCTL		= 15;
	SYS_STAT		= 16;
	SYS_FSTAT		= 17;
	SYS_CHMOD		= 18;
	SYS_SYSBRK		= 20;
	SYS_MMAP		= 21;
	SYS_MUNMAP		= 22;
	SYS_GETDENTS		= 23;
	SYS_EXIT		= 30;
	SYS_GETPID		= 31;
	SYS_SCHED_YIELD		= 32;
	SYS_SYSCONF		= 33;
	SYS_GETTIMEOFDAY	= 40;
	SYS_CLOCK		= 41;
	SYS_MULTIMEDIA_INIT	= 50;
	SYS_MULTIMEDIA_SHUTDOWN	= 51;
	SYS_VIDEO_INIT		= 52;
	SYS_VIDEO_SHUTDOWN	= 53;
	SYS_VIDEO_UPDATE	= 54;
	SYS_VIDEO_POLL_EVENT	= 55;
	SYS_AUDIO_INIT		= 56;
	SYS_AUDIO_SHUTDOWN	= 57;
	SYS_AUDIO_STREAM	= 58;
	SYS_IMC_MAKEBOUNDSOCK	= 60;
	SYS_IMC_ACCEPT		= 61;
	SYS_IMC_CONNECT		= 62;
	SYS_IMC_SENDMSG		= 63;
	SYS_IMC_RECVMSG		= 64;
	SYS_IMC_MEM_OBJ_CREATE	= 65;
	SYS_IMC_SOCKETPAIR	= 66;
	SYS_MUTEX_CREATE	= 70;
	SYS_MUTEX_LOCK		= 71;
	SYS_MUTEX_TRYLOCK	= 72;
	SYS_MUTEX_UNLOCK	= 73;
	SYS_COND_CREATE		= 74;
	SYS_COND_WAIT		= 75;
	SYS_COND_SIGNAL		= 76;
	SYS_COND_BROADCAST	= 77;
	SYS_COND_TIMED_WAIT_ABS	= 79;
	SYS_THREAD_CREATE	= 80;
	SYS_THREAD_EXIT		= 81;
	SYS_TLS_INIT		= 82;
	SYS_THREAD_NICE		= 83;
	SYS_SRPC_GET_FD		= 90;
	SYS_SEM_CREATE		= 100;
	SYS_SEM_WAIT		= 101;
	SYS_SEM_POST		= 102;
	SYS_SEM_GET_VALUE	= 103;
)
