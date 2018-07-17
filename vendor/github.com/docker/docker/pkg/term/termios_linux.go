package term 

import (
	"golang.org/x/sys/unix"
)

const (
	getTermios = unix.TCGETS
	setTermios = unix.TCSETS
)


type Termios unix.Termios




func MakeRaw(fd uintptr) (*State, error) {
	termios, err := unix.IoctlGetTermios(int(fd), getTermios)
	if err != nil {
		return nil, err
	}

	var oldState State
	oldState.termios = Termios(*termios)

	termios.Iflag &^= (unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON)
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= (unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN)
	termios.Cflag &^= (unix.CSIZE | unix.PARENB)
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(int(fd), setTermios, termios); err != nil {
		return nil, err
	}
	return &oldState, nil
}
