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
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/ccmetadata"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/util"
	cutil "github.com/mcc-github/blockchain/core/container/util"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

var logger = flogging.MustGetLogger("chaincode.platform.java")


type Platform struct{}


func (p *Platform) Name() string {
	return pb.ChaincodeSpec_JAVA.String()
}


func (p *Platform) ValidatePath(rawPath string) error {
	path, err := url.Parse(rawPath)
	if err != nil || path == nil {
		logger.Errorf("invalid chaincode path %s %v", rawPath, err)
		return fmt.Errorf("invalid path: %s", err)
	}

	return nil
}

func (p *Platform) ValidateCodePackage(code []byte) error {
	
	filesToMatch := regexp.MustCompile(`^(/)?src/((src|META-INF)/.*|(build\.gradle|settings\.gradle|pom\.xml))`)
	filesToIgnore := regexp.MustCompile(`.*\.class$`)
	is := bytes.NewReader(code)
	gr, err := gzip.NewReader(is)
	if err != nil {
		return fmt.Errorf("failure opening codepackage gzip stream: %s", err)
	}
	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		
		
		
		if !filesToMatch.MatchString(header.Name) || filesToIgnore.MatchString(header.Name) {
			return fmt.Errorf("illegal file detected in payload: \"%s\"", header.Name)
		}

		
		
		
		
		
		
		
		
		
		if header.Mode&^0100666 != 0 {
			return fmt.Errorf("illegal file mode detected for file %s: %o", header.Name, header.Mode)
		}
	}
	return nil
}


func (p *Platform) GetDeploymentPayload(path string) ([]byte, error) {
	logger.Debugf("Packaging java project from path %s", path)

	if path == "" {
		logger.Error("ChaincodeSpec's path cannot be empty")
		return nil, errors.New("ChaincodeSpec's path cannot be empty")
	}

	
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	buf := &bytes.Buffer{}
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	excludedDirs := []string{"target", "build", "out"}
	excludedFileTypes := map[string]bool{".class": true}
	err := cutil.WriteFolderToTarPackage(tw, path, excludedDirs, nil, excludedFileTypes)
	if err != nil {
		logger.Errorf("Error writing java project to tar package %s", err)
		return nil, fmt.Errorf("failed to create chaincode package: %s", err)
	}

	tw.Close()
	gw.Close()

	return buf.Bytes(), nil
}

func (p *Platform) GenerateDockerfile() (string, error) {
	var buf []string

	buf = append(buf, "FROM "+util.GetDockerfileFromConfig("chaincode.java.runtime"))
	buf = append(buf, "ADD binpackage.tar /root/chaincode-java/chaincode")

	dockerFileContents := strings.Join(buf, "\n")

	return dockerFileContents, nil
}

func (p *Platform) DockerBuildOptions(path string) (util.DockerBuildOptions, error) {
	return util.DockerBuildOptions{
		Image: util.GetDockerfileFromConfig("chaincode.java.runtime"),
		Cmd:   "./build.sh",
	}, nil
}


func (p *Platform) GetMetadataAsTarEntries(code []byte) ([]byte, error) {
	metadataProvider := &ccmetadata.TargzMetadataProvider{Code: code}
	return metadataProvider.GetMetadataAsTarEntries()
}
