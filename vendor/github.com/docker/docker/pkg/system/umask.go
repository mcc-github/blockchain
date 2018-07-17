

package system 

import (
	"golang.org/x/sys/unix"
)



func Umask(newmask int) (oldmask int, err error) {
	return unix.Umask(newmask), nil
}
