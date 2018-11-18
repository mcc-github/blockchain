

package ccmetadata

import (
	"github.com/mcc-github/blockchain/common/flogging"
)


var logger = flogging.MustGetLogger("chaincode.platform.metadata")






type MetadataProvider interface {
	GetMetadataAsTarEntries() ([]byte, error)
}
