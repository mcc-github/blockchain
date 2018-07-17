package homedir 

import (
	"os"
)



func Key() string {
	return "USERPROFILE"
}




func Get() string {
	return os.Getenv(Key())
}



func GetShortcutString() string {
	return "%USERPROFILE%" 
}
