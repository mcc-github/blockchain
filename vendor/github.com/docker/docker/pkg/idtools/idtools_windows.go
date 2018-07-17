package idtools 

import (
	"os"

	"github.com/docker/docker/pkg/system"
)



func mkdirAs(path string, mode os.FileMode, ownerUID, ownerGID int, mkAll, chownExisting bool) error {
	if err := system.MkdirAll(path, mode, ""); err != nil {
		return err
	}
	return nil
}




func CanAccess(path string, pair IDPair) bool {
	return true
}
