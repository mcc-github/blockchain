/*
Copyright DTCC 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package java

import (
	"archive/tar"
	"fmt"
	"strings"

	"errors"

	"path/filepath"

	cutil "github.com/mcc-github/blockchain/core/container/util"
)



func writeChaincodePackage(urlLocation string, tw *tar.Writer) error {
	if urlLocation == "" {
		return errors.New("ChaincodeSpec's path/URL cannot be empty")
	}

	if strings.LastIndex(urlLocation, "/") == len(urlLocation)-1 {
		urlLocation = urlLocation[:len(urlLocation)-1]
	}

	jarname := filepath.Base(urlLocation)

	err := cutil.WriteFileToPackage(urlLocation, jarname, tw)
	if err != nil {
		return fmt.Errorf("Error writing Chaincode package contents: %s", err)
	}

	return nil
}
