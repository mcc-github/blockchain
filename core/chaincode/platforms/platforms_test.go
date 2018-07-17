/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package platforms

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"os"

	"archive/tar"

	"github.com/mcc-github/blockchain/common/metadata"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/stretchr/testify/assert"
)


func TestFind(t *testing.T) {
	response, err := Find(pb.ChaincodeSpec_GOLANG)
	_, ok := response.(Platform)
	if !ok {
		t.Error("Assertion error")
	}
	assert.NotNil(t, response, "Response should have been set")
	assert.Nil(t, err, "Error should have been nil")

	response, err = Find(pb.ChaincodeSpec_CAR)
	_, ok = response.(Platform)
	if !ok {
		t.Error("Assertion error")
	}
	assert.NotNil(t, response, "Response should have been set")
	assert.Nil(t, err, "Error should have been nil")

	response, err = Find(pb.ChaincodeSpec_JAVA)
	_, ok = response.(Platform)
	if !ok {
		t.Error("Assertion error")
	}
	assert.NotNil(t, response, "Response should have been set")
	assert.Nil(t, err, "Error should have been nil")

	response, err = Find(pb.ChaincodeSpec_UNDEFINED)
	_, ok = response.(Platform)
	assert.Nil(t, response, "Response should have been nil")
	assert.NotNil(t, err, "Error should have been set")
}




type FakePlatformOk struct {
	*golang.Platform
}

func (f *FakePlatformOk) GetDeploymentPayload(spec *pb.ChaincodeSpec) ([]byte, error) {
	return []byte("success"), nil
}

func FakeFindOk(chaincodeType pb.ChaincodeSpec_Type) (Platform, error) {
	platform := &FakePlatformOk{}
	return platform, nil
}

func FakeFindErr(chaincodeType pb.ChaincodeSpec_Type) (Platform, error) {
	return nil, fmt.Errorf("Unknown chaincodeType: %s", chaincodeType)
}

func TestGetDeplotmentPayload(t *testing.T) {

	old := _Find
	defer func() { _Find = old }()

	_Find = FakeFindOk

	fake := pb.ChaincodeSpec{
		Type: pb.ChaincodeSpec_GOLANG,
	}
	response, err := GetDeploymentPayload(&fake)

	fmt.Println(err)

	assert.NotNil(t, response, "Response should have been set")
	assert.Nil(t, err, "Error should have been nil")

	_Find = FakeFindErr

	response, err = GetDeploymentPayload(&fake)

	fmt.Println(err)

	assert.NotNil(t, err, "Error should have been set")
	assert.Nil(t, response, "Response should have been nil")
}




func GetPathErr(str string) string {
	return ""
}

func VGetBoolFalse(str string) bool {
	return false
}

func OSStatErr(str string) (os.FileInfo, error) {
	return nil, errors.New("error")
}

func GetPathOk(str string) string {
	return "OK"
}

func VGetBoolTrue(str string) bool {
	return true
}

func OSStatOk(str string) (os.FileInfo, error) {
	fileInfo, _ := os.Stat("./test.txt")
	return fileInfo, nil
}

func IOUtilReadFile(str string) ([]byte, error) {
	return []byte("Stub"), nil
}




func (*FakePlatformOk) GenerateDockerfile(spec *pb.ChaincodeDeploymentSpec) (string, error) {
	return "file", nil
}

type FakePlatformErr struct {
	*golang.Platform
}

func (*FakePlatformErr) GenerateDockerfile(spec *pb.ChaincodeDeploymentSpec) (string, error) {
	return "", errors.New("error")
}

func TestGenerateDockerfile(t *testing.T) {
	mockPlatform := &FakePlatformErr{}
	fakeChaincodeSpec := &pb.ChaincodeDeploymentSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type: pb.ChaincodeSpec_GOLANG,
			ChaincodeId: &pb.ChaincodeID{
				Name:    "cc",
				Version: "1",
			},
		},
	}
	response, err := generateDockerfile(mockPlatform, fakeChaincodeSpec)
	assert.NotNil(t, err, "Error should have been set")
	assert.Nil(t, response, "Response should not have been set")

	mockPlatformOk := &FakePlatformOk{}
	response, err = generateDockerfile(mockPlatformOk, fakeChaincodeSpec)
	assert.Nil(t, err, "Error should not have been set")
	assert.NotNil(t, response, "Response should not have been set")

	var buf []string
	buf = append(buf, "file")
	buf = append(buf, fmt.Sprintf("LABEL %s.chaincode.id.name=\"%s\" \\", metadata.BaseDockerLabel, "cc"))
	buf = append(buf, fmt.Sprintf("      %s.chaincode.id.version=\"%s\" \\", metadata.BaseDockerLabel, "1"))
	buf = append(buf, fmt.Sprintf("      %s.chaincode.type=\"%s\" \\", metadata.BaseDockerLabel, "GOLANG"))
	buf = append(buf, fmt.Sprintf("      %s.version=\"%s\" \\", metadata.BaseDockerLabel, metadata.Version))
	buf = append(buf, fmt.Sprintf("      %s.base.version=\"%s\"", metadata.BaseDockerLabel, metadata.BaseVersion))
	buf = append(buf, fmt.Sprintf("ENV CORE_CHAINCODE_BUILDLEVEL=%s", metadata.Version))

	contents := strings.Join(buf, "\n")
	assert.Equal(
		t,
		response,
		[]byte(contents),
		"Should return the correct values when TLS is not enabled",
	)

	response, err = generateDockerfile(mockPlatformOk, fakeChaincodeSpec)
	contents = strings.Join(buf, "\n")
	assert.Equal(
		t,
		response,
		[]byte(contents),
		"Should return the correct values when TLS is not enabled",
	)
}




func CUtilWriteBytesToPackageOk(name string, data []byte, tw *tar.Writer) error {
	return nil
}

func CUtilWriteBytesToPackageErr(name string, data []byte, tw *tar.Writer) error {
	return errors.New("error")
}

func (*FakePlatformOk) GenerateDockerBuild(cds *pb.ChaincodeDeploymentSpec, tw *tar.Writer) error {
	return nil
}

func (*FakePlatformErr) GenerateDockerBuild(cds *pb.ChaincodeDeploymentSpec, tw *tar.Writer) error {
	return errors.New("error")
}

func TestGenerateDockerBuild1(t *testing.T) {
	oldCUtilWriteBytesToPackage := _CUtilWriteBytesToPackage

	defer func() { _CUtilWriteBytesToPackage = oldCUtilWriteBytesToPackage }()

	fakeChaincodeSpec := &pb.ChaincodeDeploymentSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type: pb.ChaincodeSpec_GOLANG,
			ChaincodeId: &pb.ChaincodeID{
				Name:    "cc",
				Version: "1",
			},
		},
	}

	inputFiles := InputFiles{
		"file1": []byte("contents1"),
		"file2": []byte("contents2"),
		"file3": []byte("contents3"),
	}

	mockPlatformOk := &FakePlatformOk{}
	mockPlatformErr := &FakePlatformErr{}
	mockTw := &tar.Writer{}

	
	_CUtilWriteBytesToPackage = CUtilWriteBytesToPackageOk
	err := generateDockerBuild(mockPlatformOk, fakeChaincodeSpec, inputFiles, mockTw)
	assert.Nil(t, err, "err should not have been set")
	
	_CUtilWriteBytesToPackage = CUtilWriteBytesToPackageErr
	err = generateDockerBuild(mockPlatformOk, fakeChaincodeSpec, inputFiles, mockTw)
	assert.NotNil(t, err, "err should have been set")

	
	_CUtilWriteBytesToPackage = CUtilWriteBytesToPackageOk
	err = generateDockerBuild(mockPlatformErr, fakeChaincodeSpec, inputFiles, mockTw)
	assert.NotNil(t, err, "err should have been set")

}





func FindOk(chaincodeType pb.ChaincodeSpec_Type) (Platform, error) {
	return &FakePlatformOk{}, nil
}

func FindErr(chaincodeType pb.ChaincodeSpec_Type) (Platform, error) {
	return nil, errors.New("error")
}

func getPeerTLSCertErr() ([]byte, error) {
	return nil, errors.New("error")
}

func generateDockerfileErr(platform Platform, cds *pb.ChaincodeDeploymentSpec) ([]byte, error) {
	return nil, errors.New("error")
}

func generateDockerBuildErr(platform Platform, cds *pb.ChaincodeDeploymentSpec, inputFiles InputFiles, tw *tar.Writer) error {
	return errors.New("error")
}

func TestGenerateDockerBuild2(t *testing.T) {

	oldFind := _Find
	oldGenerateDockerfile := _generateDockerfile
	oldGenerateDockerBuild := _generateDockerBuild
	defer func() {
		_Find = oldFind
		_generateDockerfile = oldGenerateDockerfile
		_generateDockerBuild = oldGenerateDockerBuild
	}()

	_Find = FindOk
	fakeChaincodeSpec := &pb.ChaincodeDeploymentSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type: pb.ChaincodeSpec_GOLANG,
			ChaincodeId: &pb.ChaincodeID{
				Name:    "cc",
				Version: "1",
			},
		},
	}

	
	io, err := GenerateDockerBuild(fakeChaincodeSpec)
	assert.NotNil(t, io, "io should not be nil")
	assert.Nil(t, err, "error should be nil")

	
	_Find = FindErr
	io, err = GenerateDockerBuild(fakeChaincodeSpec)
	assert.Nil(t, io, "io should be nil")
	assert.NotNil(t, err, "error should not be nil")

	
	_Find = oldFind
	_generateDockerfile = generateDockerfileErr
	io, err = GenerateDockerBuild(fakeChaincodeSpec)
	assert.Nil(t, io, "io should be nil")
	assert.NotNil(t, err, "error should not be nil")

	
	_Find = oldFind
	_generateDockerfile = oldGenerateDockerfile
	_generateDockerBuild = generateDockerBuildErr
	io, err = GenerateDockerBuild(fakeChaincodeSpec)
	assert.NotNil(t, io, "io should not be nil")
	assert.Nil(t, err, "error should be nil")
}


