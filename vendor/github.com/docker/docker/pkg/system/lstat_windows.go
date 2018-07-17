package system 

import "os"



func Lstat(path string) (*StatT, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	return fromStatT(&fi)
}
