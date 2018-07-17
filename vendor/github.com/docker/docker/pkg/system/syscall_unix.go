

package system 

import "golang.org/x/sys/unix"



func Unmount(dest string) error {
	return unix.Unmount(dest, 0)
}



func CommandLineToArgv(commandLine string) ([]string, error) {
	return []string{commandLine}, nil
}
