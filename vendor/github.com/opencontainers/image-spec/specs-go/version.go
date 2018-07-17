













package specs

import "fmt"

const (
	
	VersionMajor = 1
	
	VersionMinor = 0
	
	VersionPatch = 1

	
	VersionDev = ""
)


var Version = fmt.Sprintf("%d.%d.%d%s", VersionMajor, VersionMinor, VersionPatch, VersionDev)
