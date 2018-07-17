

package remote

import "syscall"



func syscallDup(oldfd int, newfd int) (err error) {
	return syscall.Dup3(oldfd, newfd, 0)
}
