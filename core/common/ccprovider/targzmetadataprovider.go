
package ccprovider

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"strings"
)




const (
	ccPackageStatedbDir = "META-INF/statedb/"
)

type PersistenceAdapter func([]byte) ([]byte, error)

func (pa PersistenceAdapter) GetDBArtifacts(codePackage []byte) ([]byte, error) {
	return pa(codePackage)
}


func MetadataAsTarEntries(code []byte) ([]byte, error) {
	is := bytes.NewReader(code)
	gr, err := gzip.NewReader(is)
	if err != nil {
		ccproviderLogger.Errorf("Failure opening codepackage gzip stream: %s", err)
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
			ccproviderLogger.Error("Error adding header to statedb tar:", err, header.Name)
			return nil, err
		}
		if _, err := io.Copy(tw, tr); err != nil {
			ccproviderLogger.Error("Error copying file to statedb tar:", err, header.Name)
			return nil, err
		}
		ccproviderLogger.Debug("Wrote file to statedb tar:", header.Name)
	}

	if err = tw.Close(); err != nil {
		return nil, err
	}

	ccproviderLogger.Debug("Created metadata tar")

	return statedbTarBuffer.Bytes(), nil
}
