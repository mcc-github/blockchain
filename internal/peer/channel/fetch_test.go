/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/mcc-github/blockchain/internal/peer/common"
	"github.com/mcc-github/blockchain/internal/peer/common/mock"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	defer resetFlags()
	InitMSP()
	resetFlags()
	cleanup := configtest.SetDevFabricConfigPath(t)
	defer cleanup()

	mockchain := "mockchain"

	signer, err := common.GetDefaultSigner()
	if err != nil {
		t.Fatalf("Get default signer error: %v", err)
	}

	mockCF := &ChannelCmdFactory{
		BroadcastFactory: mockBroadcastClientFactory,
		Signer:           signer,
		DeliverClient:    getMockDeliverClient(mockchain),
	}

	tempDir, err := ioutil.TempDir("", "fetch-output")
	if err != nil {
		t.Fatalf("failed to create temporary directory")
	}
	defer os.RemoveAll(tempDir)

	cmd := fetchCmd(mockCF)
	AddFlags(cmd)

	
	blocksToFetch := []string{"oldest", "newest", "config", "1"}
	for _, block := range blocksToFetch {
		outputBlockPath := filepath.Join(tempDir, block+".block")
		args := []string{"-c", mockchain, block, outputBlockPath}
		cmd.SetArgs(args)

		err = cmd.Execute()
		assert.NoError(t, err, "fetch command expected to succeed")

		if _, err := os.Stat(outputBlockPath); os.IsNotExist(err) {
			
			t.Error("expected configuration block to be fetched")
			t.Fail()
		}
	}

	
	blocksToFetchBad := []string{"banana"}
	for _, block := range blocksToFetchBad {
		outputBlockPath := filepath.Join(tempDir, block+".block")
		args := []string{"-c", mockchain, block, outputBlockPath}
		cmd.SetArgs(args)

		err = cmd.Execute()
		assert.Error(t, err, "fetch command expected to fail")
		assert.Regexp(t, err.Error(), fmt.Sprintf("fetch target illegal: %s", block))

		if fileInfo, _ := os.Stat(outputBlockPath); fileInfo != nil {
			
			t.Error("expected configuration block to not be fetched")
			t.Fail()
		}
	}
}

func TestFetchArgs(t *testing.T) {
	
	cmd := fetchCmd(nil)
	AddFlags(cmd)
	err := cmd.Execute()
	assert.Error(t, err, "fetch command expected to fail")
	assert.Contains(t, err.Error(), "fetch target required")

	
	args := []string{"strawberry", "kiwi", "lemonade"}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Error(t, err, "fetch command expected to fail")
	assert.Contains(t, err.Error(), "trailing args detected")
}

func TestFetchNilCF(t *testing.T) {
	defer viper.Reset()
	defer resetFlags()

	InitMSP()
	resetFlags()
	cleanup := configtest.SetDevFabricConfigPath(t)
	defer cleanup()

	mockchain := "mockchain"
	viper.Set("peer.client.connTimeout", 10*time.Millisecond)
	cmd := fetchCmd(nil)
	AddFlags(cmd)
	args := []string{"-c", mockchain, "oldest"}
	cmd.SetArgs(args)
	err := cmd.Execute()
	assert.Error(t, err, "fetch command expected to fail")
	assert.Contains(t, err.Error(), "deliver client failed to connect to")
}

func getMockDeliverClient(channelID string) *common.DeliverClient {
	p := getMockDeliverClientWithBlock(channelID, createTestBlock())
	return p
}

func getMockDeliverClientWithBlock(channelID string, block *cb.Block) *common.DeliverClient {
	p := &common.DeliverClient{
		Service:     getMockDeliverService(block),
		ChannelID:   channelID,
		TLSCertHash: []byte("tlscerthash"),
	}
	return p
}

func getMockDeliverService(block *cb.Block) *mock.DeliverService {
	mockD := &mock.DeliverService{}
	blockResponse := &ab.DeliverResponse{
		Type: &ab.DeliverResponse_Block{
			Block: block,
		},
	}
	statusResponse := &ab.DeliverResponse{
		Type: &ab.DeliverResponse_Status{Status: cb.Status_SUCCESS},
	}
	mockD.RecvStub = func() (*ab.DeliverResponse, error) {
		
		if mockD.RecvCallCount()%2 == 1 {
			return blockResponse, nil
		}
		return statusResponse, nil
	}
	mockD.CloseSendReturns(nil)
	return mockD
}

func createTestBlock() *cb.Block {
	lc := &cb.LastConfig{Index: 0}
	lcBytes := protoutil.MarshalOrPanic(lc)
	metadata := &cb.Metadata{
		Value: lcBytes,
	}
	metadataBytes := protoutil.MarshalOrPanic(metadata)
	blockMetadata := make([][]byte, cb.BlockMetadataIndex_LAST_CONFIG+1)
	blockMetadata[cb.BlockMetadataIndex_LAST_CONFIG] = metadataBytes
	block := &cb.Block{
		Header: &cb.BlockHeader{
			Number: 0,
		},
		Metadata: &cb.BlockMetadata{
			Metadata: blockMetadata,
		},
	}

	return block
}
