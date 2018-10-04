



package channelz

import (
	"syscall"
)


func GetSocketOption(socket interface{}) *SocketOptionData {
	c, ok := socket.(syscall.Conn)
	if !ok {
		return nil
	}
	data := &SocketOptionData{}
	if rawConn, err := c.SyscallConn(); err == nil {
		rawConn.Control(data.Getsockopt)
		return data
	}
	return nil
}
