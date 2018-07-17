
package car

import (
	"archive/tar"
	"bytes"
)


type MetadataProvider struct {
}


func (carMetadataProv *MetadataProvider) GetMetadataAsTarEntries() ([]byte, error) {
	
	
	
	
	

	buf := bytes.NewBuffer(nil)
	tw := tar.NewWriter(buf)

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
