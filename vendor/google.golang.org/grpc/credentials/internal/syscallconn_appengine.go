



package internal

import (
	"net"
)


func WrapSyscallConn(rawConn, newConn net.Conn) net.Conn {
	return newConn
}
