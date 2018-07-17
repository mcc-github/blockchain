






package logrus

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TCGETS

type Termios unix.Termios
