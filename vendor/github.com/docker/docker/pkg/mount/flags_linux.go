package mount 

import (
	"golang.org/x/sys/unix"
)

const (
	
	RDONLY = unix.MS_RDONLY

	
	
	NOSUID = unix.MS_NOSUID

	
	
	NODEV = unix.MS_NODEV

	
	NOEXEC = unix.MS_NOEXEC

	
	SYNCHRONOUS = unix.MS_SYNCHRONOUS

	
	
	
	DIRSYNC = unix.MS_DIRSYNC

	
	
	
	
	REMOUNT = unix.MS_REMOUNT

	
	MANDLOCK = unix.MS_MANDLOCK

	
	NOATIME = unix.MS_NOATIME

	
	NODIRATIME = unix.MS_NODIRATIME

	
	BIND = unix.MS_BIND

	
	RBIND = unix.MS_BIND | unix.MS_REC

	
	UNBINDABLE = unix.MS_UNBINDABLE

	
	RUNBINDABLE = unix.MS_UNBINDABLE | unix.MS_REC

	
	PRIVATE = unix.MS_PRIVATE

	
	RPRIVATE = unix.MS_PRIVATE | unix.MS_REC

	
	
	SLAVE = unix.MS_SLAVE

	
	RSLAVE = unix.MS_SLAVE | unix.MS_REC

	
	
	
	SHARED = unix.MS_SHARED

	
	RSHARED = unix.MS_SHARED | unix.MS_REC

	
	RELATIME = unix.MS_RELATIME

	
	
	
	STRICTATIME = unix.MS_STRICTATIME

	mntDetach = unix.MNT_DETACH
)
