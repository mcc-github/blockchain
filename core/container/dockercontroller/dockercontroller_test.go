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
	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/mcc-github/blockchain/core/container/dockercontroller/mock"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestIntegrationPath(t *testing.T) {
	client, err := docker.NewClientFromEnv()
	assert.NoError(t, err)

	fakePlatformBuilder := &mock.PlatformBuilder{}
	fakePlatformBuilder.GenerateDockerBuildReturns(InMemBuilder{}.Build())

	provider := Provider{
		PeerID:          "",
		NetworkID:       util.GenerateUUID(),
		BuildMetrics:    NewBuildMetrics(&disabled.Provider{}),
		Client:          client,
		PlatformBuilder: fakePlatformBuilder,
	}
	dc := provider.NewVM().(*DockerVM)
	ccid := ccintf.CCID("simple")

	err = dc.Build(ccid, "type", "path", "name", "version", []byte("code-package"))
	require.NoError(t, err)

	err = dc.Start(ccid, nil, nil, nil)
	require.NoError(t, err)

	
	err = dc.Stop(ccid, 0, true, true)
	require.NoError(t, err)

	err = dc.Start(ccid, nil, nil, nil)
	require.NoError(t, err)

	
	_ = dc.Stop(ccid, 0, false, true)
}

func Test_Start(t *testing.T) {
	gt := NewGomegaWithT(t)
	dockerClient := &mock.DockerClient{}
	dvm := DockerVM{
		BuildMetrics: NewBuildMetrics(&disabled.Provider{}),
		Client:       dockerClient,
	}
	ccid := ccintf.CCID("simple:1.0")
	args := make([]string, 1)
	env := make([]string, 1)
	files := map[string][]byte{
		"hello": []byte("world"),
	}

	
	testError1 := errors.New("junk1")
	dockerClient.CreateContainerReturns(nil, testError1)
	err := dvm.Start(ccid, args, env, files)
	gt.Expect(err).To(MatchError(testError1))
	dockerClient.CreateContainerReturns(&docker.Container{}, nil)

	
	testError2 := errors.New("junk2")
	dockerClient.UploadToContainerReturns(testError2)
	err = dvm.Start(ccid, args, env, files)
	gt.Expect(err.Error()).To(ContainSubstring("junk2"))
	dockerClient.UploadToContainerReturns(nil)

	
	
	testError3 := errors.New("junk3")
	dvm.AttachStdOut = true
	dockerClient.CreateContainerReturns(nil, testError3)
	err = dvm.Start(ccid, args, env, files)
	gt.Expect(err).To(MatchError(testError3))
	dockerClient.CreateContainerReturns(&docker.Container{}, nil)

	
	err = dvm.Start(ccid, args, env, files)
	gt.Expect(err).NotTo(HaveOccurred())

	
	err = dvm.Start(ccid, args, env, files)
	gt.Expect(err).NotTo(HaveOccurred())

	
	err = dvm.Start(ccid, args, env, files)
	gt.Expect(err).NotTo(HaveOccurred())

	
	err = dvm.Start(ccid, args, env, files)
	gt.Expect(err).NotTo(HaveOccurred())

	err = dvm.Start(ccid, args, env, files)
	gt.Expect(err).NotTo(HaveOccurred())
}

func Test_streamOutput(t *testing.T) {
	gt := NewGomegaWithT(t)

	logger, recorder := floggingtest.NewTestLogger(t)
	containerLogger, containerRecorder := floggingtest.NewTestLogger(t)

	client := &mock.DockerClient{}
	errCh := make(chan error, 1)
	optsCh := make(chan docker.AttachToContainerOptions, 1)
	client.AttachToContainerStub = func(opts docker.AttachToContainerOptions) error {
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
	client := &mock.DockerClient{}

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

			if tt.buildErr {
				client.BuildImageReturns(errors.New("Error building image"))
			}
			dvm.buildImage(ccid, &bytes.Buffer{})

			gt.Expect(fakeChaincodeImageBuildDuration.WithCallCount()).To(Equal(1))
			gt.Expect(fakeChaincodeImageBuildDuration.WithArgsForCall(0)).To(Equal(tt.expectedLabels))
			gt.Expect(fakeChaincodeImageBuildDuration.ObserveArgsForCall(0)).NotTo(BeZero())
			gt.Expect(fakeChaincodeImageBuildDuration.ObserveArgsForCall(0)).To(BeNumerically("<", 1.0))
		})
	}
}

func Test_Stop(t *testing.T) {
	dvm := DockerVM{Client: &mock.DockerClient{}}
	ccid := ccintf.CCID("simple")

	
	err := dvm.Stop(ccid, 10, true, true)
	assert.NoError(t, err)
}

func Test_Wait(t *testing.T) {
	dvm := DockerVM{}

	
	client := &mock.DockerClient{}
	dvm.Client = client

	client.WaitContainerReturns(99, nil)
	exitCode, err := dvm.Wait(ccintf.CCID("the-name:the-version"))
	assert.NoError(t, err)
	assert.Equal(t, 99, exitCode)
	assert.Equal(t, "the-name-the-version", client.WaitContainerArgsForCall(0))

	
	client.WaitContainerReturns(99, errors.New("no-wait-for-you"))
	_, err = dvm.Wait(ccintf.CCID(""))
	assert.EqualError(t, err, "no-wait-for-you")
}

func TestHealthCheck(t *testing.T) {
	client := &mock.DockerClient{}
	provider := &Provider{Client: client}

	err := provider.HealthCheck(context.Background())
	assert.NoError(t, err)

	client.PingWithContextReturns(errors.New("Error pinging daemon"))
	err = provider.HealthCheck(context.Background())
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

func TestCreateNewVM(t *testing.T) {
	networkID := util.GenerateUUID()
	client := &mock.DockerClient{}

	provider := Provider{
		PeerID:       "peerID",
		NetworkID:    networkID,
		BuildMetrics: NewBuildMetrics(&disabled.Provider{}),
		Client:       client,
		NetworkMode:  "bridge",
	}
	dvm := provider.NewVM()

	expectedClient := &DockerVM{
		PeerID:       "peerID",
		NetworkID:    networkID,
		BuildMetrics: NewBuildMetrics(&disabled.Provider{}),
		Client:       client,
		NetworkMode:  "bridge",
	}
	assert.Equal(t, expectedClient, dvm)
}

func Test_buildImage(t *testing.T) {
	client := &mock.DockerClient{}
	dvm := DockerVM{
		BuildMetrics: NewBuildMetrics(&disabled.Provider{}),
		Client:       client,
		NetworkMode:  "network-mode",
	}

	err := dvm.buildImage(ccintf.CCID("simple"), &bytes.Buffer{})
	assert.NoError(t, err)
	assert.Equal(t, 1, client.BuildImageCallCount())

	opts := client.BuildImageArgsForCall(0)
	assert.Equal(t, "simple-a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1", opts.Name)
	assert.False(t, opts.Pull)
	assert.Equal(t, "network-mode", opts.NetworkMode)
	assert.Equal(t, &bytes.Buffer{}, opts.InputStream)
	assert.NotNil(t, opts.OutputStream)
}

func Test_buildImageFailure(t *testing.T) {
	client := &mock.DockerClient{}
	client.BuildImageReturns(errors.New("oh-bother-we-failed-badly"))
	dvm := DockerVM{
		BuildMetrics: NewBuildMetrics(&disabled.Provider{}),
		Client:       client,
		NetworkMode:  "network-mode",
	}

	err := dvm.buildImage(ccintf.CCID("simple"), &bytes.Buffer{})
	assert.EqualError(t, err, "oh-bother-we-failed-badly")
}

func TestBuild(t *testing.T) {
	buildMetrics := NewBuildMetrics(&disabled.Provider{})
	ccid := ccintf.CCID("chaincode-name:chaincode-version")

	t.Run("when the image does not exist", func(t *testing.T) {
		client := &mock.DockerClient{}
		client.InspectImageReturns(nil, docker.ErrNoSuchImage)

		fakePlatformBuilder := &mock.PlatformBuilder{}
		fakePlatformBuilder.GenerateDockerBuildReturns(&bytes.Buffer{}, nil)

		dvm := &DockerVM{Client: client, BuildMetrics: buildMetrics, PlatformBuilder: fakePlatformBuilder}
		err := dvm.Build(ccid, "type", "path", "name", "version", []byte("code-package"))
		assert.NoError(t, err, "should have built successfully")

		assert.Equal(t, 1, client.BuildImageCallCount())

		require.Equal(t, 1, fakePlatformBuilder.GenerateDockerBuildCallCount())
		ccType, path, name, version, codePackage := fakePlatformBuilder.GenerateDockerBuildArgsForCall(0)
		assert.Equal(t, "type", ccType)
		assert.Equal(t, "path", path)
		assert.Equal(t, "name", name)
		assert.Equal(t, "version", version)
		assert.Equal(t, []byte("code-package"), codePackage)

	})

	t.Run("when inspecting the image fails", func(t *testing.T) {
		client := &mock.DockerClient{}
		client.InspectImageReturns(nil, errors.New("inspecting-image-fails"))

		dvm := &DockerVM{Client: client, BuildMetrics: buildMetrics}
		err := dvm.Build(ccid, "type", "path", "name", "version", []byte("code-package"))
		assert.EqualError(t, err, "docker image inspection failed: inspecting-image-fails")

		assert.Equal(t, 0, client.BuildImageCallCount())
	})

	t.Run("when the image exists", func(t *testing.T) {
		client := &mock.DockerClient{}

		dvm := &DockerVM{Client: client, BuildMetrics: buildMetrics}
		err := dvm.Build(ccid, "type", "path", "name", "version", []byte("code-package"))
		assert.NoError(t, err)

		assert.Equal(t, 0, client.BuildImageCallCount())
	})

	t.Run("when the platform builder fails", func(t *testing.T) {
		client := &mock.DockerClient{}
		client.InspectImageReturns(nil, docker.ErrNoSuchImage)
		client.BuildImageReturns(errors.New("no-build-for-you"))

		fakePlatformBuilder := &mock.PlatformBuilder{}
		fakePlatformBuilder.GenerateDockerBuildReturns(nil, errors.New("fake-builder-error"))

		dvm := &DockerVM{Client: client, BuildMetrics: buildMetrics, PlatformBuilder: fakePlatformBuilder}
		err := dvm.Build(ccid, "type", "path", "name", "version", []byte("code-package"))
		assert.Equal(t, 1, client.InspectImageCallCount())
		assert.Equal(t, 1, fakePlatformBuilder.GenerateDockerBuildCallCount())
		assert.Equal(t, 0, client.BuildImageCallCount())
		assert.EqualError(t, err, "platform builder failed: fake-builder-error")
	})

	t.Run("when building the image fails", func(t *testing.T) {
		client := &mock.DockerClient{}
		client.InspectImageReturns(nil, docker.ErrNoSuchImage)
		client.BuildImageReturns(errors.New("no-build-for-you"))

		fakePlatformBuilder := &mock.PlatformBuilder{}

		dvm := &DockerVM{Client: client, BuildMetrics: buildMetrics, PlatformBuilder: fakePlatformBuilder}
		err := dvm.Build(ccid, "type", "path", "name", "version", []byte("code-package"))
		assert.Equal(t, 1, client.InspectImageCallCount())
		assert.Equal(t, 1, client.BuildImageCallCount())
		assert.EqualError(t, err, "docker image build failed: no-build-for-you")
	})
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

type mockBuilder struct {
	buildFunc func() (io.Reader, error)
}

func (m *mockBuilder) Build() (io.Reader, error) {
	return m.buildFunc()
}
