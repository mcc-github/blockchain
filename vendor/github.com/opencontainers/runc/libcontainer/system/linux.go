

package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)












const PR_SET_CHILD_SUBREAPER = 36

type ParentDeathSignal int

func (p ParentDeathSignal) Restore() error {
	if p == 0 {
		return nil
	}
	current, err := GetParentDeathSignal()
	if err != nil {
		return err
	}
	if p == current {
		return nil
	}
	return p.Set()
}

func (p ParentDeathSignal) Set() error {
	return SetParentDeathSignal(uintptr(p))
}

func Execv(cmd string, args []string, env []string) error {
	name, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}

	return syscall.Exec(name, args, env)
}

func Prlimit(pid, resource int, limit syscall.Rlimit) error {
	_, _, err := syscall.RawSyscall6(syscall.SYS_PRLIMIT64, uintptr(pid), uintptr(resource), uintptr(unsafe.Pointer(&limit)), uintptr(unsafe.Pointer(&limit)), 0, 0)
	if err != 0 {
		return err
	}
	return nil
}

func SetParentDeathSignal(sig uintptr) error {
	if _, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_SET_PDEATHSIG, sig, 0); err != 0 {
		return err
	}
	return nil
}

func GetParentDeathSignal() (ParentDeathSignal, error) {
	var sig int
	_, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_GET_PDEATHSIG, uintptr(unsafe.Pointer(&sig)), 0)
	if err != 0 {
		return -1, err
	}
	return ParentDeathSignal(sig), nil
}

func SetKeepCaps() error {
	if _, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_SET_KEEPCAPS, 1, 0); err != 0 {
		return err
	}

	return nil
}

func ClearKeepCaps() error {
	if _, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_SET_KEEPCAPS, 0, 0); err != 0 {
		return err
	}

	return nil
}

func Setctty() error {
	if _, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, 0, uintptr(syscall.TIOCSCTTY), 0); err != 0 {
		return err
	}
	return nil
}


func RunningInUserNS() bool {
	file, err := os.Open("/proc/self/uid_map")
	if err != nil {
		
		return false
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	l, _, err := buf.ReadLine()
	if err != nil {
		return false
	}

	line := string(l)
	var a, b, c int64
	fmt.Sscanf(line, "%d %d %d", &a, &b, &c)
	
	if a == 0 && b == 0 && c == 4294967295 {
		return false
	}
	return true
}


func SetSubreaper(i int) error {
	return Prctl(PR_SET_CHILD_SUBREAPER, uintptr(i), 0, 0, 0)
}

func Prctl(option int, arg2, arg3, arg4, arg5 uintptr) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_PRCTL, uintptr(option), arg2, arg3, arg4, arg5, 0)
	if e1 != 0 {
		err = e1
	}
	return
}