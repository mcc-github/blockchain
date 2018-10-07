



package archive

import (
	"os"
	"path/filepath"
)




func SplitPathDirEntry(path string) (dir, base string) {
	cleanedPath := filepath.Clean(filepath.FromSlash(path))

	if specifiesCurrentDir(path) {
		cleanedPath += string(os.PathSeparator) + "."
	}

	return filepath.Dir(cleanedPath), filepath.Base(cleanedPath)
}



func specifiesCurrentDir(path string) bool {
	return filepath.Base(path) == "."
}
