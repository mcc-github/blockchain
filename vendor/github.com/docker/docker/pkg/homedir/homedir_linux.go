package homedir 

import (
	"os"

	"github.com/docker/docker/pkg/idtools"
)






func GetStatic() (string, error) {
	uid := os.Getuid()
	usr, err := idtools.LookupUID(uid)
	if err != nil {
		return "", err
	}
	return usr.Home, nil
}
