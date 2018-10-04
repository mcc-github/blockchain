



package credentials

import (
	"errors"
	"syscall"
)


func (c tlsConn) SyscallConn() (syscall.RawConn, error) {
	conn, ok := c.rawConn.(syscall.Conn)
	if !ok {
		return nil, errors.New("RawConn does not implement syscall.Conn")
	}
	return conn.SyscallConn()
}
