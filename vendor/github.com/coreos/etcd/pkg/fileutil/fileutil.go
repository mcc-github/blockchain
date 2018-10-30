














package fileutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/coreos/pkg/capnslog"
)

const (
	
	PrivateFileMode = 0600
	
	PrivateDirMode = 0700
)

var (
	plog = capnslog.NewPackageLogger("github.com/coreos/etcd", "pkg/fileutil")
)



func IsDirWriteable(dir string) error {
	f := filepath.Join(dir, ".touch")
	if err := ioutil.WriteFile(f, []byte(""), PrivateFileMode); err != nil {
		return err
	}
	return os.Remove(f)
}


func ReadDir(dirpath string) ([]string, error) {
	dir, err := os.Open(dirpath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}



func TouchDirAll(dir string) error {
	
	
	err := os.MkdirAll(dir, PrivateDirMode)
	if err != nil {
		
		
		return err
	}
	return IsDirWriteable(dir)
}



func CreateDirAll(dir string) error {
	err := TouchDirAll(dir)
	if err == nil {
		var ns []string
		ns, err = ReadDir(dir)
		if err != nil {
			return err
		}
		if len(ns) != 0 {
			err = fmt.Errorf("expected %q to be empty, got %q", dir, ns)
		}
	}
	return err
}

func Exist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}



func ZeroToEnd(f *os.File) error {
	
	off, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	lenf, lerr := f.Seek(0, io.SeekEnd)
	if lerr != nil {
		return lerr
	}
	if err = f.Truncate(off); err != nil {
		return err
	}
	
	if err = Preallocate(f, lenf, true); err != nil {
		return err
	}
	_, err = f.Seek(off, io.SeekStart)
	return err
}
