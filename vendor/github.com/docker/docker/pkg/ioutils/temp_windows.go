package ioutils 

import (
	"io/ioutil"

	"github.com/docker/docker/pkg/longpath"
)


func TempDir(dir, prefix string) (string, error) {
	tempDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return "", err
	}
	return longpath.AddPrefix(tempDir), nil
}
