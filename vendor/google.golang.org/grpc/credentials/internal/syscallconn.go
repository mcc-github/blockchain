




package internal

import (
	"net"
	"syscall"
)

type sysConn = syscall.Conn










type syscallConn struct {
	net.Conn
	
	
	sysConn
}






func WrapSyscallConn(rawConn, newConn net.Conn) net.Conn {
	sysConn, ok := rawConn.(syscall.Conn)
	if !ok {
		return newConn
	}
	return &syscallConn{
		Conn:    newConn,
		sysConn: sysConn,
	}
}
