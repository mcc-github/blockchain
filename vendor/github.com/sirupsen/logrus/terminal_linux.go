






package logrus

import (
	"io"

	"golang.org/x/sys/unix"
)

const ioctlReadTermios = unix.TCGETS

type Termios unix.Termios

func initTerminal(w io.Writer) {
}
