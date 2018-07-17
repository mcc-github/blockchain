

package homedir 

import (
	"os"

	"github.com/opencontainers/runc/libcontainer/user"
)



func Key() string {
	return "HOME"
}




func Get() string {
	home := os.Getenv(Key())
	if home == "" {
		if u, err := user.CurrentUser(); err == nil {
			return u.Home
		}
	}
	return home
}



func GetShortcutString() string {
	return "~"
}
