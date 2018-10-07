

package pb

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

var tty = os.Stdin

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	
	
	
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")

	
	
	
	getConsoleMode = kernel32.NewProc("GetConsoleMode")

	
	
	
	setConsoleMode = kernel32.NewProc("SetConsoleMode")

	
	
	
	setConsoleCursorPosition = kernel32.NewProc("SetConsoleCursorPosition")
)

type (
	
	
	
	
	smallRect struct {
		Left, Top, Right, Bottom int16
	}

	
	
	
	
	
	coordinates struct {
		X, Y int16
	}

	word int16

	
	
	consoleScreenBufferInfo struct {
		dwSize              coordinates
		dwCursorPosition    coordinates
		wAttributes         word
		srWindow            smallRect
		dwMaximumWindowSize coordinates
	}
)


func terminalWidth() (width int, err error) {
	var info consoleScreenBufferInfo
	_, _, e := syscall.Syscall(procGetConsoleScreenBufferInfo.Addr(), 2, uintptr(syscall.Stdout), uintptr(unsafe.Pointer(&info)), 0)
	if e != 0 {
		return 0, error(e)
	}
	return int(info.dwSize.X) - 1, nil
}

func getCursorPos() (pos coordinates, err error) {
	var info consoleScreenBufferInfo
	_, _, e := syscall.Syscall(procGetConsoleScreenBufferInfo.Addr(), 2, uintptr(syscall.Stdout), uintptr(unsafe.Pointer(&info)), 0)
	if e != 0 {
		return info.dwCursorPosition, error(e)
	}
	return info.dwCursorPosition, nil
}

func setCursorPos(pos coordinates) error {
	_, _, e := syscall.Syscall(setConsoleCursorPosition.Addr(), 2, uintptr(syscall.Stdout), uintptr(uint32(uint16(pos.Y))<<16|uint32(uint16(pos.X))), 0)
	if e != 0 {
		return error(e)
	}
	return nil
}

var ErrPoolWasStarted = errors.New("Bar pool was started")

var echoLocked bool
var echoLockMutex sync.Mutex

var oldState word

func lockEcho() (shutdownCh chan struct{}, err error) {
	echoLockMutex.Lock()
	defer echoLockMutex.Unlock()
	if echoLocked {
		err = ErrPoolWasStarted
		return
	}
	echoLocked = true

	if _, _, e := syscall.Syscall(getConsoleMode.Addr(), 2, uintptr(syscall.Stdout), uintptr(unsafe.Pointer(&oldState)), 0); e != 0 {
		err = fmt.Errorf("Can't get terminal settings: %v", e)
		return
	}

	newState := oldState
	const ENABLE_ECHO_INPUT = 0x0004
	const ENABLE_LINE_INPUT = 0x0002
	newState = newState & (^(ENABLE_LINE_INPUT | ENABLE_ECHO_INPUT))
	if _, _, e := syscall.Syscall(setConsoleMode.Addr(), 2, uintptr(syscall.Stdout), uintptr(newState), 0); e != 0 {
		err = fmt.Errorf("Can't set terminal settings: %v", e)
		return
	}

	shutdownCh = make(chan struct{})
	return
}

func unlockEcho() (err error) {
	echoLockMutex.Lock()
	defer echoLockMutex.Unlock()
	if !echoLocked {
		return
	}
	echoLocked = false
	if _, _, e := syscall.Syscall(setConsoleMode.Addr(), 2, uintptr(syscall.Stdout), uintptr(oldState), 0); e != 0 {
		err = fmt.Errorf("Can't set terminal settings")
	}
	return
}
