/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package dockercontroller

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	coreutil "github.com/mcc-github/blockchain/core/testutil"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


func TestIntegrationPath(t *testing.T) {
	coreutil.SetupTestConfig()
	dc := NewDockerVM("", util.GenerateUUID())
	ccid := ccintf.CCID{Name: "simple"}

	err := dc.Start(ccid, nil, nil, nil, InMemBuilder{})
	require.NoError(t, err)

	
	err = dc.Stop(ccid, 0, true, true)
	require.NoError(t, err)

	err = dc.Start(ccid, nil, nil, nil, nil)
	require.NoError(t, err)

	
	_ = dc.Stop(ccid, 0, false, true)
}

func TestHostConfig(t *testing.T) {
	coreutil.SetupTestConfig()
	var hostConfig = new(docker.HostConfig)
	err := viper.UnmarshalKey("vm.docker.hostConfig", hostConfig)
	if err != nil {
		t.Fatalf("Load docker HostConfig wrong, error: %s", err.Error())
	}
	testutil.AssertNotEquals(t, hostConfig.LogConfig, nil)
	testutil.AssertEquals(t, hostConfig.LogConfig.Type, "json-file")
	testutil.AssertEquals(t, hostConfig.LogConfig.Config["max-size"], "50m")
	testutil.AssertEquals(t, hostConfig.LogConfig.Config["max-file"], "5")
}

func TestGetDockerHostConfig(t *testing.T) {
	coreutil.SetupTestConfig()
	hostConfig = nil 
	hostConfig := getDockerHostConfig()
	testutil.AssertNotNil(t, hostConfig)
	testutil.AssertEquals(t, hostConfig.NetworkMode, "host")
	testutil.AssertEquals(t, hostConfig.LogConfig.Type, "json-file")
	testutil.AssertEquals(t, hostConfig.LogConfig.Config["max-size"], "50m")
	testutil.AssertEquals(t, hostConfig.LogConfig.Config["max-file"], "5")
	testutil.AssertEquals(t, hostConfig.Memory, int64(1024*1024*1024*2))
	testutil.AssertEquals(t, hostConfig.CPUShares, int64(0))
}

func Test_Start(t *testing.T) {
	dvm := DockerVM{}
	ccid := ccintf.CCID{Name: "simple"}
	args := make([]string, 1)
	env := make([]string, 1)
	files := map[string][]byte{
		"hello": []byte("world"),
	}

	
	
	dvm.getClientFnc = getMockClient
	getClientErr = true
	err := dvm.Start(ccid, args, env, files, nil)
	testerr(t, err, false)
	getClientErr = false

	
	createErr = true
	err = dvm.Start(ccid, args, env, files, nil)
	testerr(t, err, false)
	createErr = false

	
	uploadErr = true
	err = dvm.Start(ccid, args, env, files, nil)
	testerr(t, err, false)
	uploadErr = false

	
	noSuchImgErr = true
	err = dvm.Start(ccid, args, env, files, nil)
	testerr(t, err, false)

	chaincodePath := "github.com/mcc-github/blockchain/examples/chaincode/go/example01/cmd"
	spec := &pb.ChaincodeSpec{
		Type:        pb.ChaincodeSpec_GOLANG,
		ChaincodeId: &pb.ChaincodeID{Name: "ex01", Path: chaincodePath},
		Input:       &pb.ChaincodeInput{Args: util.ToChaincodeArgs("f")},
	}
	codePackage, err := platforms.NewRegistry(&golang.Platform{}).GetDeploymentPayload(spec.CCType(), spec.Path())
	if err != nil {
		t.Fatal()
	}
	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackage}
	bldr := &mockBuilder{
		buildFunc: func() (io.Reader, error) {
			return platforms.NewRegistry(&golang.Platform{}).GenerateDockerBuild(
				cds.CCType(),
				cds.Path(),
				cds.Name(),
				cds.Version(),
				cds.Bytes(),
			)
		},
	}

	
	
	viper.Set("vm.docker.attachStdout", true)
	startErr = true
	err = dvm.Start(ccid, args, env, files, bldr)
	testerr(t, err, false)
	startErr = false

	
	err = dvm.Start(ccid, args, env, files, bldr)
	testerr(t, err, true)
	noSuchImgErr = false

	
	stopErr = true
	err = dvm.Start(ccid, args, env, files, nil)
	testerr(t, err, true)
	stopErr = false

	
	killErr = true
	err = dvm.Start(ccid, args, env, files, nil)
	testerr(t, err, true)
	killErr = false

	
	removeErr = true
	err = dvm.Start(ccid, args, env, files, nil)
	testerr(t, err, true)
	removeErr = false

	err = dvm.Start(ccid, args, env, files, nil)
	testerr(t, err, true)
}

func Test_Stop(t *testing.T) {
	dvm := DockerVM{}
	ccid := ccintf.CCID{Name: "simple"}

	
	getClientErr = true
	dvm.getClientFnc = getMockClient
	err := dvm.Stop(ccid, 10, true, true)
	testerr(t, err, false)
	getClientErr = false

	
	err = dvm.Stop(ccid, 10, true, true)
	testerr(t, err, true)
}

type testCase struct {
	name           string
	vm             *DockerVM
	ccid           ccintf.CCID
	expectedOutput string
}

func TestGetVMNameForDocker(t *testing.T) {
	tc := []testCase{
		{
			name:           "mycc",
			vm:             &DockerVM{NetworkID: "dev", PeerID: "peer0"},
			ccid:           ccintf.CCID{Name: "mycc", Version: "1.0"},
			expectedOutput: fmt.Sprintf("%s-%s", "dev-peer0-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("dev-peer0-mycc-1.0")))),
		},
		{
			name:           "mycc-nonetworkid",
			vm:             &DockerVM{PeerID: "peer1"},
			ccid:           ccintf.CCID{Name: "mycc", Version: "1.0"},
			expectedOutput: fmt.Sprintf("%s-%s", "peer1-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("peer1-mycc-1.0")))),
		},
		{
			name:           "myCC-UCids",
			vm:             &DockerVM{NetworkID: "Dev", PeerID: "Peer0"},
			ccid:           ccintf.CCID{Name: "myCC", Version: "1.0"},
			expectedOutput: fmt.Sprintf("%s-%s", "dev-peer0-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("Dev-Peer0-myCC-1.0")))),
		},
		{
			name:           "myCC-idsWithSpecialChars",
			vm:             &DockerVM{NetworkID: "Dev$dev", PeerID: "Peer*0"},
			ccid:           ccintf.CCID{Name: "myCC", Version: "1.0"},
			expectedOutput: fmt.Sprintf("%s-%s", "dev-dev-peer-0-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("Dev$dev-Peer*0-myCC-1.0")))),
		},
		{
			name:           "mycc-nopeerid",
			vm:             &DockerVM{NetworkID: "dev"},
			ccid:           ccintf.CCID{Name: "mycc", Version: "1.0"},
			expectedOutput: fmt.Sprintf("%s-%s", "dev-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("dev-mycc-1.0")))),
		},
		{
			name:           "myCC-LCids",
			vm:             &DockerVM{NetworkID: "dev", PeerID: "peer0"},
			ccid:           ccintf.CCID{Name: "myCC", Version: "1.0"},
			expectedOutput: fmt.Sprintf("%s-%s", "dev-peer0-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("dev-peer0-myCC-1.0")))),
		},
	}

	for _, test := range tc {
		name, err := test.vm.GetVMNameForDocker(test.ccid)
		assert.Nil(t, err, "Expected nil error")
		assert.Equal(t, test.expectedOutput, name, "Unexpected output for test case name: %s", test.name)
	}

}

func TestGetVMName(t *testing.T) {
	tc := []testCase{
		{
			name:           "myCC-preserveCase",
			vm:             &DockerVM{NetworkID: "Dev", PeerID: "Peer0"},
			ccid:           ccintf.CCID{Name: "myCC", Version: "1.0"},
			expectedOutput: fmt.Sprintf("%s", "Dev-Peer0-myCC-1.0"),
		},
	}

	for _, test := range tc {
		name := test.vm.GetVMName(test.ccid)
		assert.Equal(t, test.expectedOutput, name, "Unexpected output for test case name: %s", test.name)
	}

}



type InMemBuilder struct{}

func (imb InMemBuilder) Build() (io.Reader, error) {
	buf := &bytes.Buffer{}
	fmt.Fprintln(buf, "FROM busybox:latest")
	fmt.Fprintln(buf, `CMD ["tail", "-f", "/dev/null"]`)

	startTime := time.Now()
	inputbuf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(inputbuf)
	tr := tar.NewWriter(gw)
	tr.WriteHeader(&tar.Header{
		Name:       "Dockerfile",
		Size:       int64(buf.Len()),
		ModTime:    startTime,
		AccessTime: startTime,
		ChangeTime: startTime,
	})
	tr.Write(buf.Bytes())
	tr.Close()
	gw.Close()
	return inputbuf, nil
}

func testerr(t *testing.T, err error, succ bool) {
	if succ {
		assert.NoError(t, err, "Expected success but got error")
	} else {
		assert.Error(t, err, "Expected failure but succeeded")
	}
}

func getMockClient() (dockerClient, error) {
	if getClientErr {
		return nil, errors.New("Failed to get client")
	}
	return &mockClient{noSuchImgErrReturned: false}, nil
}

type mockBuilder struct {
	buildFunc func() (io.Reader, error)
}

func (m *mockBuilder) Build() (io.Reader, error) {
	return m.buildFunc()
}

type mockClient struct {
	noSuchImgErrReturned bool
}

var getClientErr, createErr, uploadErr, noSuchImgErr, buildErr, removeImgErr,
	startErr, stopErr, killErr, removeErr bool

func (c *mockClient) CreateContainer(options docker.CreateContainerOptions) (*docker.Container, error) {
	if createErr {
		return nil, errors.New("Error creating the container")
	}
	if noSuchImgErr && !c.noSuchImgErrReturned {
		c.noSuchImgErrReturned = true
		return nil, docker.ErrNoSuchImage
	}
	return &docker.Container{}, nil
}

func (c *mockClient) StartContainer(id string, cfg *docker.HostConfig) error {
	if startErr {
		return errors.New("Error starting the container")
	}
	return nil
}

func (c *mockClient) UploadToContainer(id string, opts docker.UploadToContainerOptions) error {
	if uploadErr {
		return errors.New("Error uploading archive to the container")
	}
	return nil
}

func (c *mockClient) AttachToContainer(opts docker.AttachToContainerOptions) error {
	if opts.Success != nil {
		opts.Success <- struct{}{}
	}
	return nil
}

func (c *mockClient) BuildImage(opts docker.BuildImageOptions) error {
	if buildErr {
		return errors.New("Error building image")
	}
	return nil
}

func (c *mockClient) RemoveImageExtended(id string, opts docker.RemoveImageOptions) error {
	if removeImgErr {
		return errors.New("Error removing extended image")
	}
	return nil
}

func (c *mockClient) StopContainer(id string, timeout uint) error {
	if stopErr {
		return errors.New("Error stopping container")
	}
	return nil
}

func (c *mockClient) KillContainer(opts docker.KillContainerOptions) error {
	if killErr {
		return errors.New("Error killing container")
	}
	return nil
}

func (c *mockClient) RemoveContainer(opts docker.RemoveContainerOptions) error {
	if removeErr {
		return errors.New("Error removing container")
	}
	return nil
}
