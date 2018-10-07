





package term

import "golang.org/x/sys/unix"


func GetWinsize(fd uintptr) (*Winsize, error) {
	uws, err := unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
	ws := &Winsize{Height: uws.Row, Width: uws.Col, x: uws.Xpixel, y: uws.Ypixel}
	return ws, err
}
