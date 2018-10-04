





package unix

import "runtime"




func IoctlSetWinsize(fd int, req uint, value *Winsize) error {
	
	
	err := ioctlSetWinsize(fd, req, value)
	runtime.KeepAlive(value)
	return err
}




func IoctlSetTermios(fd int, req uint, value *Termios) error {
	
	err := ioctlSetTermios(fd, req, value)
	runtime.KeepAlive(value)
	return err
}
