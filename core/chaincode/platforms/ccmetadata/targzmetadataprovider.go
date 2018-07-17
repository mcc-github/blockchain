
package ccmetadata

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)




const (
	ccPackageStatedbDir = "META-INF/statedb/"
)



type TargzMetadataProvider struct {
	DepSpec *pb.ChaincodeDeploymentSpec
}

func (tgzProv *TargzMetadataProvider) getCode() ([]byte, error) {
	if tgzProv.DepSpec == nil {
		logger.Errorf("nil chaincode deployment spec")
		return nil, errors.New("nil chaincode deployment spec")
	}

	if tgzProv.DepSpec.ChaincodeSpec == nil || tgzProv.DepSpec.ChaincodeSpec.ChaincodeId == nil {
		return nil, errors.New("invalid chaincode deployment spec")
	}

	if tgzProv.DepSpec.CodePackage == nil {
		return nil, errors.New(fmt.Sprintf("nil code package for %v", tgzProv.DepSpec.ChaincodeSpec.ChaincodeId))
	}

	return tgzProv.DepSpec.CodePackage, nil
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
