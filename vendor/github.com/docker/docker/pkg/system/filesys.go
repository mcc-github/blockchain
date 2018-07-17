

package system 

import (
	"io/ioutil"
	"os"
	"path/filepath"
)


func MkdirAllWithACL(path string, perm os.FileMode, sddl string) error {
	return MkdirAll(path, perm, sddl)
}



func MkdirAll(path string, perm os.FileMode, sddl string) error {
	return os.MkdirAll(path, perm)
}


func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}









func CreateSequential(name string) (*os.File, error) {
	return os.Create(name)
}





func OpenSequential(name string) (*os.File, error) {
	return os.Open(name)
}






func OpenFileSequential(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}










func TempFileSequential(dir, prefix string) (f *os.File, err error) {
	return ioutil.TempFile(dir, prefix)
}
