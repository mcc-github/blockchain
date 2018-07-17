

package user

import (
	"io"
	"os"
)


const (
	unixPasswdPath = "/etc/passwd"
	unixGroupPath  = "/etc/group"
)

func GetPasswdPath() (string, error) {
	return unixPasswdPath, nil
}

func GetPasswd() (io.ReadCloser, error) {
	return os.Open(unixPasswdPath)
}

func GetGroupPath() (string, error) {
	return unixGroupPath, nil
}

func GetGroup() (io.ReadCloser, error) {
	return os.Open(unixGroupPath)
}
