
package ccmetadata

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"strings"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/pkg/errors"
)




const (
	ccPackageStatedbDir = "META-INF/statedb/"
)


var logger = flogging.MustGetLogger("chaincode.platform.metadata")




type PersistenceMetadataProvider struct{}



func (t *PersistenceMetadataProvider) GetDBArtifacts(codePackage []byte) ([]byte, error) {
	return (&TargzMetadataProvider{
		Code: codePackage,
	}).GetMetadataAsTarEntries()
}



type TargzMetadataProvider struct {
	Code []byte
}

func (tgzProv *TargzMetadataProvider) getCode() ([]byte, error) {
	if tgzProv.Code == nil {
		return nil, errors.New("nil code package")
	}

	return tgzProv.Code, nil
}


func (tgzProv *TargzMetadataProvider) GetMetadataAsTarEntries() ([]byte, error) {
	code, err := tgzProv.getCode()
	if err != nil {
		return nil, err
	}

	is := bytes.NewReader(code)
	gr, err := gzip.NewReader(is)
	if err != nil {
		logger.Errorf("Failure opening codepackage gzip stream: %s", err)
		return nil, err
	}

	statedbTarBuffer := bytes.NewBuffer(nil)
	tw := tar.NewWriter(statedbTarBuffer)

	tr := tar.NewReader(gr)

	
	
	for {
		header, err := tr.Next()
		if err == io.EOF {
			
			break
		}
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(header.Name, ccPackageStatedbDir) {
			continue
		}

		if err = tw.WriteHeader(header); err != nil {
			logger.Error("Error adding header to statedb tar:", err, header.Name)
			return nil, err
		}
		if _, err := io.Copy(tw, tr); err != nil {
			logger.Error("Error copying file to statedb tar:", err, header.Name)
			return nil, err
		}
		logger.Debug("Wrote file to statedb tar:", header.Name)
	}

	if err = tw.Close(); err != nil {
		return nil, err
	}

	logger.Debug("Created metadata tar")

	return statedbTarBuffer.Bytes(), nil
}
