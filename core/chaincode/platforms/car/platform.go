/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package car

import (
	"archive/tar"
	"io/ioutil"
	"strings"

	"bytes"
	"fmt"
	"io"

	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/util"
	cutil "github.com/mcc-github/blockchain/core/container/util"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type Platform struct {
}


func (carPlatform *Platform) Name() string {
	return pb.ChaincodeSpec_CAR.String()
}




func (carPlatform *Platform) ValidatePath(path string) error {
	return nil
}

func (carPlatform *Platform) ValidateCodePackage(codePackage []byte) error {
	
	return nil
}

func (carPlatform *Platform) GetDeploymentPayload(path string) ([]byte, error) {

	return ioutil.ReadFile(path)
}

func (carPlatform *Platform) GenerateDockerfile() (string, error) {

	var buf []string

	
	buf = append(buf, "FROM "+cutil.GetDockerfileFromConfig("chaincode.car.runtime"))
	buf = append(buf, "ADD binpackage.tar /usr/local/bin")

	dockerFileContents := strings.Join(buf, "\n")

	return dockerFileContents, nil
}

func (carPlatform *Platform) GenerateDockerBuild(path string, code []byte, tw *tar.Writer) error {

	
	codepackage, output := io.Pipe()
	go func() {
		tw := tar.NewWriter(output)

		err := cutil.WriteBytesToPackage("codepackage.car", code, tw)

		tw.Close()
		output.CloseWithError(err)
	}()

	binpackage := bytes.NewBuffer(nil)
	err := util.DockerBuild(util.DockerBuildOptions{
		Cmd:          "java -jar /usr/local/bin/chaintool buildcar /chaincode/input/codepackage.car -o /chaincode/output/chaincode",
		InputStream:  codepackage,
		OutputStream: binpackage,
	})
	if err != nil {
		return fmt.Errorf("Error building CAR: %s", err)
	}

	return cutil.WriteBytesToPackage("binpackage.tar", binpackage.Bytes(), tw)
}


func (carPlatform *Platform) GetMetadataProvider(code []byte) platforms.MetadataProvider {
	return &MetadataProvider{}
}
