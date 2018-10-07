



package archive

import (
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/longpath"
)




func CanonicalTarNameForPath(p string) (string, error) {
	
	
	
	
	if strings.Contains(p, "/") {
		return "", fmt.Errorf("Windows path contains forward slash: %s", p)
	}
	return strings.Replace(p, string(os.PathSeparator), "/", -1), nil

}



func fixVolumePathPrefix(srcPath string) string {
	return longpath.AddPrefix(srcPath)
}



func getWalkRoot(srcPath string, include string) string {
	return filepath.Join(srcPath, include)
}

func getInodeFromStat(stat interface{}) (inode uint64, err error) {
	
	return
}

func getFileIdentity(stat interface{}) (idtools.Identity, error) {
	
	return idtools.Identity{}, nil
}



func chmodTarEntry(perm os.FileMode) os.FileMode {
	
	permPart := perm & os.ModePerm
	noPermPart := perm &^ os.ModePerm
	
	permPart |= 0111
	permPart &= 0755

	return noPermPart | permPart
}

func setHeaderForSpecialDevice(hdr *tar.Header, name string, stat interface{}) (err error) {
	
	return
}
