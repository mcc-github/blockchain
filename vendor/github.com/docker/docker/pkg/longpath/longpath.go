



package longpath 

import (
	"strings"
)


const Prefix = `\\?\`



func AddPrefix(path string) string {
	if !strings.HasPrefix(path, Prefix) {
		if strings.HasPrefix(path, `\\`) {
			
			path = Prefix + `UNC` + path[1:]
		} else {
			path = Prefix + path
		}
	}
	return path
}
