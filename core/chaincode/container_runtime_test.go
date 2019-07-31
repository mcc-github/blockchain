/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode_test

import (
	"testing"

	"github.com/mcc-github/blockchain/core/chaincode"
	"github.com/mcc-github/blockchain/core/chaincode/mock"
	persistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestContainerRuntimeStart(t *testing.T) {
	fakeRouter := &mock.ContainerRouter{}

	cr := &chaincode.ContainerRuntime{
		ContainerRouter: fakeRouter,
		PeerAddress:     "peer-address",
	}

	ccci := &ccprovider.ChaincodeContainerInfo{
		Type:      "GOLANG",
		Path:      "chaincode-path",
		Name:      "chaincode-name",
		Version:   "chaincode-version",
		PackageID: "chaincode-name:chaincode-version",
	}

	err := cr.Start(ccci, []byte("code-package"))
	assert.NoError(t, err)

	assert.Equal(t, 1, fakeRouter.BuildCallCount())
	ccci, codePackage := fakeRouter.BuildArgsForCall(0)
	assert.Equal(t, &ccprovider.ChaincodeContainerInfo{
		PackageID: "chaincode-name:chaincode-version",
		Type:      "GOLANG",
		Path:      "chaincode-path",
		Name:      "chaincode-name",
		Version:   "chaincode-version",
	}, ccci)
	assert.Equal(t, []byte("code-package"), codePackage)

	assert.Equal(t, 1, fakeRouter.StartCallCount())
	ccid, peerConnection := fakeRouter.StartArgsForCall(0)
	assert.Equal(t, ccintf.CCID("chaincode-name:chaincode-version"), ccid)
	assert.Equal(t, "peer-address", peerConnection.Address)
	assert.Nil(t, peerConnection.TLSConfig)
}

func TestContainerRuntimeStartErrors(t *testing.T) {
	tests := []struct {
		chaincodeType string
		buildErr      error
		startErr      error
		errValue      string
	}{
		{pb.ChaincodeSpec_GOLANG.String(), nil, errors.New("process-failed"), "error starting container: process-failed"},
		{pb.ChaincodeSpec_GOLANG.String(), errors.New("build-failed"), nil, "error building image: build-failed"},
		{pb.ChaincodeSpec_GOLANG.String(), errors.New("build-failed"), nil, "error building image: build-failed"},
	}

	for _, tc := range tests {
		fakeRouter := &mock.ContainerRouter{}
		fakeRouter.BuildReturns(tc.buildErr)
		fakeRouter.StartReturns(tc.startErr)

		cr := &chaincode.ContainerRuntime{
			ContainerRouter: fakeRouter,
		}

		ccci := &ccprovider.ChaincodeContainerInfo{
			Type:    tc.chaincodeType,
			Name:    "chaincode-id-name",
			Version: "chaincode-version",
		}

		err := cr.Start(ccci, nil)
		assert.EqualError(t, err, tc.errValue)
	}
}

func TestContainerRuntimeStop(t *testing.T) {
	fakeRouter := &mock.ContainerRouter{}

	cr := &chaincode.ContainerRuntime{
		ContainerRouter: fakeRouter,
	}

	ccci := &ccprovider.ChaincodeContainerInfo{
		Type:      pb.ChaincodeSpec_GOLANG.String(),
		PackageID: "chaincode-id-name:chaincode-version",
	}

	err := cr.Stop(ccci)
	assert.NoError(t, err)

	assert.Equal(t, 1, fakeRouter.StopCallCount())
	ccid := fakeRouter.StopArgsForCall(0)
	assert.Equal(t, ccintf.CCID("chaincode-id-name:chaincode-version"), ccid)
}

func TestContainerRuntimeStopErrors(t *testing.T) {
	tests := []struct {
		processErr error
		errValue   string
	}{
		{errors.New("process-failed"), "error stopping container: process-failed"},
	}

	for _, tc := range tests {
		fakeRouter := &mock.ContainerRouter{}
		fakeRouter.StopReturns(tc.processErr)

		cr := &chaincode.ContainerRuntime{
			ContainerRouter: fakeRouter,
		}

		ccci := &ccprovider.ChaincodeContainerInfo{
			Type:    pb.ChaincodeSpec_GOLANG.String(),
			Name:    "chaincode-id-name",
			Version: "chaincode-version",
		}

		assert.EqualError(t, cr.Stop(ccci), tc.errValue)
	}
}

func TestContainerRuntimeWait(t *testing.T) {
	fakeRouter := &mock.ContainerRouter{}

	cr := &chaincode.ContainerRuntime{
		ContainerRouter: fakeRouter,
	}

	ccci := &ccprovider.ChaincodeContainerInfo{
		Type:      pb.ChaincodeSpec_GOLANG.String(),
		Name:      "chaincode-id-name",
		Version:   "chaincode-version",
		PackageID: persistence.PackageID("chaincode-id-name:chaincode-version"),
	}

	exitCode, err := cr.Wait(ccci)
	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, 1, fakeRouter.WaitCallCount())
	assert.Equal(t, ccintf.CCID("chaincode-id-name:chaincode-version"), fakeRouter.WaitArgsForCall(0))

	fakeRouter.WaitReturns(3, errors.New("moles-and-trolls"))
	code, err := cr.Wait(ccci)
	assert.EqualError(t, err, "moles-and-trolls")
	assert.Equal(t, code, 3)
}
