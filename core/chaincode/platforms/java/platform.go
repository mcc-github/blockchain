/*
Copyright DTCC 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package java

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/ccmetadata"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/util"
	cutil "github.com/mcc-github/blockchain/core/container/util"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

var logger = flogging.MustGetLogger("java-platform")


type Platform struct {
}


func (javaPlatform *Platform) Name() string {
	return pb.ChaincodeSpec_JAVA.String()
}


func (javaPlatform *Platform) ValidatePath(rawPath string) error {
	path, err := url.Parse(rawPath)
	if err != nil || path == nil {
		logger.Errorf("invalid chaincode path %s %v", rawPath, err)
		return fmt.Errorf("invalid path: %s", err)
	}

	return nil
}

func (javaPlatform *Platform) ValidateCodePackage(code []byte) error {
	
	return nil
}


func (javaPlatform *Platform) GetDeploymentPayload(path string) ([]byte, error) {

	logger.Debugf("Packaging java project from path %s", path)
	var err error

	
	
	
	payload := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(payload)
	tw := tar.NewWriter(gw)

	folder := path
	if folder == "" {
		logger.Error("ChaincodeSpec's path cannot be empty")
		return nil, errors.New("ChaincodeSpec's path cannot be empty")
	}

	
	if folder[len(folder)-1] == '/' {
		folder = folder[:len(folder)-1]
	}

	if err = cutil.WriteJavaProjectToPackage(tw, folder); err != nil {

		logger.Errorf("Error writing java project to tar package %s", err)
		return nil, fmt.Errorf("Error writing Chaincode package contents: %s", err)
	}

	tw.Close()
	gw.Close()

	return payload.Bytes(), nil
}

func (javaPlatform *Platform) GenerateDockerfile() (string, error) {
	var buf []string

	buf = append(buf, "FROM "+cutil.GetDockerfileFromConfig("chaincode.java.Dockerfile"))
	buf = append(buf, "ADD binpackage.tar /root/chaincode-java/chaincode")

	dockerFileContents := strings.Join(buf, "\n")

	return dockerFileContents, nil
}

func (javaPlatform *Platform) GenerateDockerBuild(path string, code []byte, tw *tar.Writer) error {
	codepackage := bytes.NewReader(code)
	binpackage := bytes.NewBuffer(nil)
	buildOptions := util.DockerBuildOptions{
		Image:        cutil.GetDockerfileFromConfig("chaincode.java.Dockerfile"),
		Cmd:          "./build.sh",
		InputStream:  codepackage,
		OutputStream: binpackage,
	}
	logger.Debugf("Executing docker build %v, %v", buildOptions.Image, buildOptions.Cmd)
	err := util.DockerBuild(buildOptions)
	if err != nil {
		logger.Errorf("Can't build java chaincode %v", err)
		return err
	}

	resultBytes := binpackage.Bytes()
	return cutil.WriteBytesToPackage("binpackage.tar", resultBytes, tw)
}


func (javaPlatform *Platform) GetMetadataProvider(code []byte) platforms.MetadataProvider {
	return &ccmetadata.TargzMetadataProvider{Code: code}
}
