

package system 

import (
	"syscall"

	"golang.org/x/sys/unix"
)


func IsProcessAlive(pid int) bool {
	err := unix.Kill(pid, syscall.Signal(0))
	if err == nil || err == unix.EPERM {
		return true
	}

	return false
}


func KillProcess(pid int) {
	unix.Kill(pid, unix.SIGKILL)
}
