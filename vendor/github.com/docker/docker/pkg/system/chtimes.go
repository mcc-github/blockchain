package system 

import (
	"os"
	"time"
)


func Chtimes(name string, atime time.Time, mtime time.Time) error {
	unixMinTime := time.Unix(0, 0)
	unixMaxTime := maxTime

	
	
	

	if atime.Before(unixMinTime) || atime.After(unixMaxTime) {
		atime = unixMinTime
	}

	if mtime.Before(unixMinTime) || mtime.After(unixMaxTime) {
		mtime = unixMinTime
	}

	if err := os.Chtimes(name, atime, mtime); err != nil {
		return err
	}

	
	return setCTime(name, mtime)
}
