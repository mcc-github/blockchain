/*
Copyright DTCC 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package java

import (
	"archive/tar"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"errors"

	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	ccutil "github.com/mcc-github/blockchain/core/chaincode/platforms/util"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

var logger = flogging.MustGetLogger("java/hash")







func collectChaincodeFiles(spec *pb.ChaincodeSpec, tw *tar.Writer) (string, error) {
	if spec == nil {
		return "", errors.New("Cannot collect chaincode files from nil spec")
	}

	chaincodeID := spec.ChaincodeId
	if chaincodeID == nil || chaincodeID.Path == "" {
		return "", errors.New("Cannot collect chaincode files from empty chaincode path")
	}

	codepath := chaincodeID.Path

	var err error
	if !strings.HasPrefix(codepath, "/") {
		wd := ""
		wd, err = os.Getwd()
		codepath = wd + "/" + codepath
	}

	if err != nil {
		return "", fmt.Errorf("Error getting code %s", err)
	}

	var hash []byte

	
	if spec.Input == nil || len(spec.Input.Args) == 0 {
		logger.Debugf("not using input for hash computation for %v ", chaincodeID)
	} else {
		inputbytes, err2 := proto.Marshal(spec.Input)
		if err2 != nil {
			return "", fmt.Errorf("Error marshalling constructor: %s", err)
		}
		hash = util.GenerateHashFromSignature(codepath, inputbytes)
	}

	buf, err := ioutil.ReadFile(codepath)
	if err != nil {
		logger.Errorf("Error reading %s", err)
		return "", err
	}

	
	hash = ccutil.ComputeHash(buf, hash)

	return hex.EncodeToString(hash[:]), nil
}
