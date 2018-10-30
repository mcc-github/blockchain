














package util

import (
	"bytes"
	"os"
	"syscall"
)



func SysReadFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	
	
	
	
	
	b := make([]byte, 128)
	n, err := syscall.Read(int(f.Fd()), b)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(b[:n])), nil
}
