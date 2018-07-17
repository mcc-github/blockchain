package system 

import (
	"errors"
)

var (
	
	ErrNotSupportedPlatform = errors.New("platform and architecture is not supported")

	
	ErrNotSupportedOperatingSystem = errors.New("operating system is not supported")
)
