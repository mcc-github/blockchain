

package ioutils 

import "io/ioutil"


func TempDir(dir, prefix string) (string, error) {
	return ioutil.TempDir(dir, prefix)
}
