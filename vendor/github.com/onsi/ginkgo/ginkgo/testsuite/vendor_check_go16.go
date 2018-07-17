

package testsuite

import (
	"os"
	"path"
)



func vendorExperimentCheck(dir string) bool {
	vendorExperiment := os.Getenv("GO15VENDOREXPERIMENT")
	return vendorExperiment != "0" && path.Base(dir) == "vendor"
}
