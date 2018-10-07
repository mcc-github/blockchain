





package archive

import (
	"archive/tar"
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"github.com/docker/docker/pkg/idtools"
	"golang.org/x/sys/unix"
)




func CanonicalTarNameForPath(p string) (string, error) {
	return p, nil 
}



func fixVolumePathPrefix(srcPath string) string {
	return srcPath
}





func getWalkRoot(srcPath string, include string) string {
	return srcPath + string(filepath.Separator) + include
}

func getInodeFromStat(stat interface{}) (inode uint64, err error) {
	s, ok := stat.(*syscall.Stat_t)

	if ok {
		inode = uint64(s.Ino)
	}

	return
}

func getFileIdentity(stat interface{}) (idtools.Identity, error) {
	s, ok := stat.(*syscall.Stat_t)

	if !ok {
		return idtools.Identity{}, errors.New("cannot convert stat value to syscall.Stat_t")
	}
	return idtools.Identity{UID: int(s.Uid), GID: int(s.Gid)}, nil
}

func chmodTarEntry(perm os.FileMode) os.FileMode {
	return perm 
}

func setHeaderForSpecialDevice(hdr *tar.Header, name string, stat interface{}) (err error) {
	s, ok := stat.(*syscall.Stat_t)

	if ok {
		
		if s.Mode&unix.S_IFBLK != 0 ||
			s.Mode&unix.S_IFCHR != 0 {
			hdr.Devmajor = int64(unix.Major(uint64(s.Rdev))) 
			hdr.Devminor = int64(unix.Minor(uint64(s.Rdev))) 
		}
	}

	return
}
