

package archive 

import (
	"archive/tar"
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/system"
	rsystem "github.com/opencontainers/runc/libcontainer/system"
	"golang.org/x/sys/unix"
)



func fixVolumePathPrefix(srcPath string) string {
	return srcPath
}





func getWalkRoot(srcPath string, include string) string {
	return srcPath + string(filepath.Separator) + include
}




func CanonicalTarNameForPath(p string) (string, error) {
	return p, nil 
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

func getInodeFromStat(stat interface{}) (inode uint64, err error) {
	s, ok := stat.(*syscall.Stat_t)

	if ok {
		inode = s.Ino
	}

	return
}

func getFileUIDGID(stat interface{}) (idtools.IDPair, error) {
	s, ok := stat.(*syscall.Stat_t)

	if !ok {
		return idtools.IDPair{}, errors.New("cannot convert stat value to syscall.Stat_t")
	}
	return idtools.IDPair{UID: int(s.Uid), GID: int(s.Gid)}, nil
}



func handleTarTypeBlockCharFifo(hdr *tar.Header, path string) error {
	if rsystem.RunningInUserNS() {
		
		return nil
	}

	mode := uint32(hdr.Mode & 07777)
	switch hdr.Typeflag {
	case tar.TypeBlock:
		mode |= unix.S_IFBLK
	case tar.TypeChar:
		mode |= unix.S_IFCHR
	case tar.TypeFifo:
		mode |= unix.S_IFIFO
	}

	return system.Mknod(path, mode, int(system.Mkdev(hdr.Devmajor, hdr.Devminor)))
}

func handleLChmod(hdr *tar.Header, path string, hdrInfo os.FileInfo) error {
	if hdr.Typeflag == tar.TypeLink {
		if fi, err := os.Lstat(hdr.Linkname); err == nil && (fi.Mode()&os.ModeSymlink == 0) {
			if err := os.Chmod(path, hdrInfo.Mode()); err != nil {
				return err
			}
		}
	} else if hdr.Typeflag != tar.TypeSymlink {
		if err := os.Chmod(path, hdrInfo.Mode()); err != nil {
			return err
		}
	}
	return nil
}
