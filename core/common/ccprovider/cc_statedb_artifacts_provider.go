/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccprovider

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mcc-github/blockchain/core/chaincode/platforms"
)

const (
	ccPackageStatedbDir = "META-INF/statedb/"
)


type TarFileEntry struct {
	FileHeader  *tar.Header
	FileContent []byte
}




func ExtractStatedbArtifactsForChaincode(ccname, ccversion string) (installed bool, statedbArtifactsTar []byte, err error) {
	ccpackage, err := GetChaincodeFromFS(ccname, ccversion)
	if err != nil {
		
		
		
		ccproviderLogger.Info("Error while loading installation package for ccname=%s, ccversion=%s. Err=%s", ccname, ccversion, err)
		return false, nil, nil
	}

	statedbArtifactsTar, err = ExtractStatedbArtifactsFromCCPackage(ccpackage)
	return true, statedbArtifactsTar, err
}




func ExtractStatedbArtifactsFromCCPackage(ccpackage CCPackage) (statedbArtifactsTar []byte, err error) {
	cds := ccpackage.GetDepSpec()
	pform, err := platforms.Find(cds.ChaincodeSpec.Type)
	if err != nil {
		ccproviderLogger.Infof("invalid deployment spec (bad platform type:%s)", cds.ChaincodeSpec.Type)
		return nil, fmt.Errorf("invalid deployment spec")
	}
	metaprov := pform.GetMetadataProvider(cds)
	return metaprov.GetMetadataAsTarEntries()
}









func ExtractFileEntries(tarBytes []byte, databaseType string) (map[string][]*TarFileEntry, error) {

	indexArtifacts := map[string][]*TarFileEntry{}
	tarReader := tar.NewReader(bytes.NewReader(tarBytes))
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			
			break
		}
		if err != nil {
			return nil, err
		}
		
		dir, _ := filepath.Split(hdr.Name)
		
		if strings.HasPrefix(hdr.Name, "META-INF/statedb/"+databaseType) {
			fileContent, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return nil, err
			}
			indexArtifacts[filepath.Clean(dir)] = append(indexArtifacts[filepath.Clean(dir)], &TarFileEntry{FileHeader: hdr, FileContent: fileContent})
		}
	}

	return indexArtifacts, nil
}
