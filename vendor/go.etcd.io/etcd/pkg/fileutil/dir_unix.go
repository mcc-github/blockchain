















package fileutil

import "os"


func OpenDir(path string) (*os.File, error) { return os.Open(path) }
