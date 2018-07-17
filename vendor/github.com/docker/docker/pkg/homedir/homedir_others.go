

package homedir 

import (
	"errors"
)



func GetStatic() (string, error) {
	return "", errors.New("homedir.GetStatic() is not supported on this system")
}
