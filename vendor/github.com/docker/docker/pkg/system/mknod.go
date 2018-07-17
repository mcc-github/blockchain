

package system 

import (
	"golang.org/x/sys/unix"
)



func Mknod(path string, mode uint32, dev int) error {
	return unix.Mknod(path, mode, dev)
}






func Mkdev(major int64, minor int64) uint32 {
	return uint32(unix.Mkdev(uint32(major), uint32(minor)))
}
