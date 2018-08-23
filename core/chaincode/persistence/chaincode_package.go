/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package persistence

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)









const (
	ChaincodePackageMetadataFile = "Chaincode-Package-Metadata.json"
)


type ChaincodePackage struct {
	Metadata    *ChaincodePackageMetadata
	CodePackage []byte
}



type ChaincodePackageMetadata struct {
	Type string `json:"Type"`
	Path string `json:"Path"`
}


type ChaincodePackageParser struct{}



func (ccpp ChaincodePackageParser) Parse(source []byte) (*ChaincodePackage, error) {
	gzReader, err := gzip.NewReader(bytes.NewBuffer(source))
	if err != nil {
		return nil, errors.Wrapf(err, "error reading as gzip stream")
	}

	tarReader := tar.NewReader(gzReader)

	var codePackage []byte
	var ccPackageMetadata *ChaincodePackageMetadata
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, errors.Wrapf(err, "error inspecting next tar header")
		}

		if header.Typeflag != tar.TypeReg {
			return nil, errors.Errorf("tar entry %s is not a regular file, type %v", header.Name, header.Typeflag)
		}

		fileBytes, err := ioutil.ReadAll(tarReader)
		if err != nil {
			return nil, errors.Wrapf(err, "could not read %s from tar", header.Name)
		}

		if header.Name == ChaincodePackageMetadataFile {
			ccPackageMetadata = &ChaincodePackageMetadata{}
			err := json.Unmarshal(fileBytes, ccPackageMetadata)
			if err != nil {
				return nil, errors.Wrapf(err, "could not unmarshal %s as json", ChaincodePackageMetadataFile)
			}

			continue
		}

		if codePackage != nil {
			return nil, errors.Errorf("found too many files in archive, cannot identify which file is the code-package")
		}

		codePackage = fileBytes
	}

	if codePackage == nil {
		return nil, errors.Errorf("did not find a code package inside the package")
	}

	if ccPackageMetadata == nil {
		return nil, errors.Errorf("did not find any package metadata (missing %s)", ChaincodePackageMetadataFile)
	}

	return &ChaincodePackage{
		Metadata:    ccPackageMetadata,
		CodePackage: codePackage,
	}, nil
}
