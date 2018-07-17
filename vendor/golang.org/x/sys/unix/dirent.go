





package unix

import "syscall"





func ParseDirent(buf []byte, max int, names []string) (consumed int, count int, newnames []string) {
	return syscall.ParseDirent(buf, max, names)
}
