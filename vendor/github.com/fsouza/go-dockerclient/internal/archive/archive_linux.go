



package archive

import (
	"archive/tar"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/system"
	"golang.org/x/sys/unix"
)

const (
	
	AUFSWhiteoutFormat WhiteoutFormat = iota
	
	
	OverlayWhiteoutFormat
)

func getWhiteoutConverter(format WhiteoutFormat) tarWhiteoutConverter {
	if format == OverlayWhiteoutFormat {
		return overlayWhiteoutConverter{}
	}
	return nil
}

type overlayWhiteoutConverter struct{}

func (overlayWhiteoutConverter) ConvertWrite(hdr *tar.Header, path string, fi os.FileInfo) (wo *tar.Header, err error) {
	
	if fi.Mode()&os.ModeCharDevice != 0 && hdr.Devmajor == 0 && hdr.Devminor == 0 {
		
		dir, filename := filepath.Split(hdr.Name)
		hdr.Name = filepath.Join(dir, WhiteoutPrefix+filename)
		hdr.Mode = 0600
		hdr.Typeflag = tar.TypeReg
		hdr.Size = 0
	}

	if fi.Mode()&os.ModeDir != 0 {
		
		opaque, err := system.Lgetxattr(path, "trusted.overlay.opaque")
		if err != nil {
			return nil, err
		}
		if len(opaque) == 1 && opaque[0] == 'y' {
			if hdr.Xattrs != nil {
				delete(hdr.Xattrs, "trusted.overlay.opaque")
			}

			
			
			wo = &tar.Header{
				Typeflag:   tar.TypeReg,
				Mode:       hdr.Mode & int64(os.ModePerm),
				Name:       filepath.Join(hdr.Name, WhiteoutOpaqueDir),
				Size:       0,
				Uid:        hdr.Uid,
				Uname:      hdr.Uname,
				Gid:        hdr.Gid,
				Gname:      hdr.Gname,
				AccessTime: hdr.AccessTime,
				ChangeTime: hdr.ChangeTime,
			}
		}
	}

	return
}

func (overlayWhiteoutConverter) ConvertRead(hdr *tar.Header, path string) (bool, error) {
	base := filepath.Base(path)
	dir := filepath.Dir(path)

	
	if base == WhiteoutOpaqueDir {
		err := unix.Setxattr(dir, "trusted.overlay.opaque", []byte{'y'}, 0)
		
		return false, err
	}

	
	if strings.HasPrefix(base, WhiteoutPrefix) {
		originalBase := base[len(WhiteoutPrefix):]
		originalPath := filepath.Join(dir, originalBase)

		if err := unix.Mknod(originalPath, unix.S_IFCHR, 0); err != nil {
			return false, err
		}
		if err := os.Chown(originalPath, hdr.Uid, hdr.Gid); err != nil {
			return false, err
		}

		
		return false, nil
	}

	return true, nil
}
