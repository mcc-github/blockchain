/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package dockercontroller

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/mcc-github/blockchain/common/flogging/floggingtest"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/common/metrics/metricsfakes"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	pb "github.com/mcc-github/blockchain/protos/peer"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestIntegrationPath(t *testing.T) {
	client, err := docker.NewClientFromEnv()
	assert.NoError(t, err)
	provider := Provider{
		PeerID:       "",
		NetworkID:    util.GenerateUUID(),
		BuildMetrics: NewBuildMetrics(&disabled.Provider{}),
		Client:       client,
	}
	dc := provider.NewVM()
	ccid := ccintf.CCID("simple")

	err = dc.Start(ccid, nil, nil, nil, InMemBuilder{})
	require.NoError(t, err)

	
	err = dc.Stop(ccid, 0, true, true)
	require.NoError(t, err)

	err = dc.Start(ccid, nil, nil, nil, nil)
	require.NoError(t, err)

	
	_ = dc.Stop(ccid, 0, false, true)
}

func Test_Start(t *testing.T) {
	gt := NewGomegaWithT(t)
	dvm := DockerVM{
		BuildMetrics: NewBuildMetrics(&disabled.Provider{}),
		Client:       &mockClient{},
	}
	ccid := ccintf.CCID("simple:1.0")
	args := make([]string, 1)
	env := make([]string, 1)
	files := map[string][]byte{
		"hello": []byte("world"),
	}

	
	createErr = true
	err := dvm.Start(ccid, args, env, files, nil)
	gt.Expect(err).To(HaveOccurred())
	createErr = false

	
	uploadErr = true
	err = dvm.Start(ccid, args, env, files, nil)
	gt.Expect(err).To(HaveOccurred())
	uploadErr = false

	
	noSuchImgErr = true
	buildErr = true
	err = dvm.Start(ccid, args, env, files, &mockBuilder{buildFunc: func() (io.Reader, error) { return &bytes.Buffer{}, nil }})
	gt.Expect(err).To(HaveOccurred())
	buildErr = false

	chaincodePath := "github.com/mcc-github/blockchain/core/container/dockercontroller/testdata/src/chaincodes/noop"
	spec := &pb.ChaincodeSpec{
		Type:        pb.ChaincodeSpec_GOLANG,
		ChaincodeId: &pb.ChaincodeID{Name: "ex01", Path: chaincodePath},
		Input:       &pb.ChaincodeInput{Args: util.ToChaincodeArgs("f")},
	}
	codePackage, err := platforms.NewRegistry(&golang.Platform{}).GetDeploymentPayload(spec.Type.String(), spec.ChaincodeId.Path)
	if err != nil {
		t.Fatal()
	}
	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackage}
	client, err := docker.NewClientFromEnv()
	assert.NoError(t, err)
	bldr := &mockBuilder{
		buildFunc: func() (io.Reader, error) {
			return platforms.NewRegistry(&golang.Platform{}).GenerateDockerBuild(
				cds.ChaincodeSpec.Type.String(),
				cds.ChaincodeSpec.ChaincodeId.Path,
				cds.ChaincodeSpec.ChaincodeId.Name,
				cds.ChaincodeSpec.ChaincodeId.Version,
				cds.CodePackage,
				client,
			)
		},
	}

	
	
	dvm.AttachStdOut = true
	startErr = true
	err = dvm.Start(ccid, args, env, files, bldr)
	gt.Expect(err).To(HaveOccurred())
	startErr = false

	
	err = dvm.Start(ccid, args, env, files, bldr)
	gt.Expect(err).NotTo(HaveOccurred())
	noSuchImgErr = false

	
	stopErr = true
	err = dvm.Start(ccid, args, env, files, nil)
	gt.Expect(err).NotTo(HaveOccurred())
	stopErr = false

	
	killErr = true
	err = dvm.Start(ccid, args, env, files, nil)
	gt.Expect(err).NotTo(HaveOccurred())
	killErr = false

	
	removeErr = true
	err = dvm.Start(ccid, args, env, files, nil)
	gt.Expect(err).NotTo(HaveOccurred())
	removeErr = false

	err = dvm.Start(ccid, args, env, files, nil)
	gt.Expect(err).NotTo(HaveOccurred())
}

func Test_streamOutput(t *testing.T) {
	gt := NewGomegaWithT(t)

	logger, recorder := floggingtest.NewTestLogger(t)
	containerLogger, containerRecorder := floggingtest.NewTestLogger(t)

	client := &mockClient{}
	errCh := make(chan error, 1)
	optsCh := make(chan docker.AttachToContainerOptions, 1)
	client.attachToContainerStub = func(opts docker.AttachToContainerOptions) error {
		optsCh <- opts
		return <-errCh
	}

	streamOutput(logger, client, "container-name", containerLogger)

	var opts docker.AttachToContainerOptions
	gt.Eventually(optsCh).Should(Receive(&opts))
	gt.Eventually(opts.Success).Should(BeSent(struct{}{}))
	gt.Eventually(opts.Success).Should(BeClosed())

	fmt.Fprintf(opts.OutputStream, "message-one\n")
	fmt.Fprintf(opts.OutputStream, "message-two") 
	gt.Eventually(containerRecorder).Should(gbytes.Say("message-one"))
	gt.Consistently(containerRecorder.Entries).Should(HaveLen(1))

	close(errCh)
	gt.Eventually(recorder).Should(gbytes.Say("Container container-name has closed its IO channel"))
	gt.Consistently(recorder.Entries).Should(HaveLen(1))
	gt.Consistently(containerRecorder.Entries).Should(HaveLen(1))
}

func Test_BuildMetric(t *testing.T) {
	ccid := ccintf.CCID("simple:1.0")
	client := &mockClient{}

	tests := []struct {
		desc           string
		buildErr       bool
		expectedLabels []string
	}{
		{desc: "success", buildErr: false, expectedLabels: []string{"chaincode", "simple:1.0", "success", "true"}},
		{desc: "failure", buildErr: true, expectedLabels: []string{"chaincode", "simple:1.0", "success", "false"}},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			fakeChaincodeImageBuildDuration := &metricsfakes.Histogram{}
			fakeChaincodeImageBuildDuration.WithReturns(fakeChaincodeImageBuildDuration)
			dvm := DockerVM{
				BuildMetrics: &BuildMetrics{
					ChaincodeImageBuildDuration: fakeChaincodeImageBuildDuration,
				},
				Client: client,
			}

			buildErr = tt.buildErr
			dvm.deployImage(ccid, &bytes.Buffer{})

			gt.Expect(fakeChaincodeImageBuildDuration.WithCallCount()).To(Equal(1))
			gt.Expect(fakeChaincodeImageBuildDuration.WithArgsForCall(0)).To(Equal(tt.expectedLabels))
			gt.Expect(fakeChaincodeImageBuildDuration.ObserveArgsForCall(0)).NotTo(BeZero())
			gt.Expect(fakeChaincodeImageBuildDuration.ObserveArgsForCall(0)).To(BeNumerically("<", 1.0))
		})
	}

	buildErr = false
}

func Test_Stop(t *testing.T) {
	dvm := DockerVM{Client: &mockClient{}}
	ccid := ccintf.CCID("simple")

	
	err := dvm.Stop(ccid, 10, true, true)
	assert.NoError(t, err)
}

func Test_Wait(t *testing.T) {
	dvm := DockerVM{}

	
	client := &mockClient{}
	dvm.Client = client

	client.exitCode = 99
	exitCode, err := dvm.Wait(ccintf.CCID("the-name:the-version"))
	assert.NoError(t, err)
	assert.Equal(t, 99, exitCode)
	assert.Equal(t, "the-name-the-version", client.containerID)

	
	client.waitErr = errors.New("no-wait-for-you")
	_, err = dvm.Wait(ccintf.CCID(""))
	assert.EqualError(t, err, "no-wait-for-you")
}

func Test_HealthCheck(t *testing.T) {
	dvm := DockerVM{}

	dvm.Client = &mockClient{
		pingErr: false,
	}
	err := dvm.HealthCheck(context.Background())
	assert.NoError(t, err)

	dvm.Client = &mockClient{
		pingErr: true,
	}
	err = dvm.HealthCheck(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error pinging daemon")
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
			ccid:           ccintf.CCID("mycc:1.0"),
			expectedOutput: fmt.Sprintf("%s-%s", "dev-peer0-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("dev-peer0-mycc:1.0")))),
		},
		{
			name:           "mycc-nonetworkid",
			vm:             &DockerVM{PeerID: "peer1"},
			ccid:           ccintf.CCID("mycc:1.0"),
			expectedOutput: fmt.Sprintf("%s-%s", "peer1-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("peer1-mycc:1.0")))),
		},
		{
			name:           "myCC-UCids",
			vm:             &DockerVM{NetworkID: "Dev", PeerID: "Peer0"},
			ccid:           ccintf.CCID("myCC:1.0"),
			expectedOutput: fmt.Sprintf("%s-%s", "dev-peer0-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("Dev-Peer0-myCC:1.0")))),
		},
		{
			name:           "myCC-idsWithSpecialChars",
			vm:             &DockerVM{NetworkID: "Dev$dev", PeerID: "Peer*0"},
			ccid:           ccintf.CCID("myCC:1.0"),
			expectedOutput: fmt.Sprintf("%s-%s", "dev-dev-peer-0-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("Dev$dev-Peer*0-myCC:1.0")))),
		},
		{
			name:           "mycc-nopeerid",
			vm:             &DockerVM{NetworkID: "dev"},
			ccid:           ccintf.CCID("mycc:1.0"),
			expectedOutput: fmt.Sprintf("%s-%s", "dev-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("dev-mycc:1.0")))),
		},
		{
			name:           "myCC-LCids",
			vm:             &DockerVM{NetworkID: "dev", PeerID: "peer0"},
			ccid:           ccintf.CCID("myCC:1.0"),
			expectedOutput: fmt.Sprintf("%s-%s", "dev-peer0-mycc-1.0", hex.EncodeToString(util.ComputeSHA256([]byte("dev-peer0-myCC:1.0")))),
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
			ccid:           ccintf.CCID("myCC:1.0"),
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
	pingErr              bool

	containerID string
	exitCode    int
	waitErr     error

	attachToContainerStub func(docker.AttachToContainerOptions) error
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
	if c.attachToContainerStub != nil {
		return c.attachToContainerStub(opts)
	}
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

func (c *mockClient) PingWithContext(context.Context) error {
	if c.pingErr {
		return errors.New("Error pinging daemon")
	}
	return nil
}

func (c *mockClient) WaitContainer(id string) (int, error) {
	c.containerID = id
	return c.exitCode, c.waitErr
}
