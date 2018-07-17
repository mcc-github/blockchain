/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"fmt"
	"testing"

	"github.com/mcc-github/blockchain/peer/common"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestQueryCmd(t *testing.T) {
	InitMSP()

	mockCF, err := getMockChaincodeCmdFactory()
	assert.NoError(t, err, "Error getting mock chaincode command factory")
	
	channelID = ""

	
	args := []string{"-n", "example02", "-c", "{\"Args\": [\"query\",\"a\"]}"}
	cmd := newQueryCmdForTest(mockCF, args)
	err = cmd.Execute()
	assert.Error(t, err, "'peer chaincode query' command should have failed without -C flag")

	
	args = []string{"-C", "mychannel", "-n", "example02", "-c", "{\"Args\": [\"query\",\"a\"]}"}
	cmd = newQueryCmdForTest(mockCF, args)
	err = cmd.Execute()
	assert.NoError(t, err, "Run chaincode query cmd error")

	
	args = []string{"-r", "-C", "mychannel", "-n", "example02", "-c", "{\"Args\": [\"query\",\"a\"]}"}
	cmd = newQueryCmdForTest(mockCF, args)
	err = cmd.Execute()
	assert.NoError(t, err, "Run chaincode query cmd error")
	chaincodeQueryRaw = false

	
	args = []string{"-x", "-C", "mychannel", "-n", "example02", "-c", "{\"Args\": [\"query\",\"a\"]}"}
	cmd = newQueryCmdForTest(mockCF, args)
	err = cmd.Execute()
	assert.NoError(t, err, "Run chaincode query cmd error")

	
	args = []string{"-r", "-x", "-C", "mychannel", "-n", "example02", "-c", "{\"Args\": [\"query\",\"a\"]}"}
	cmd = newQueryCmdForTest(mockCF, args)
	err = cmd.Execute()
	assert.Error(t, err, "Expected error executing query command with both -r and -x options")

	
	mockCF, err = getMockChaincodeCmdFactoryWithErr()
	assert.NoError(t, err, "Error getting mock chaincode command factory")
	args = []string{"-r", "-n", "example02", "-c", "{\"Args\": [\"query\",\"a\"]}"}
	cmd = newQueryCmdForTest(mockCF, args)
	err = cmd.Execute()
	assert.Error(t, err, "Expected error executing query command")
}

func TestQueryCmdEndorsementFailure(t *testing.T) {
	InitMSP()
	args := []string{"-C", "mychannel", "-n", "example02", "-c", "{\"Args\": [\"queryinvalid\",\"a\"]}"}
	ccRespStatus := [2]int32{502, 400}
	ccRespPayload := [][]byte{[]byte("Invalid function name"), []byte("Incorrect parameters")}

	for i := 0; i < 2; i++ {
		mockCF, err := getMockChaincodeCmdFactoryEndorsementFailure(ccRespStatus[i], ccRespPayload[i])
		assert.NoError(t, err, "Error getting mock chaincode command factory")

		cmd := newQueryCmdForTest(mockCF, args)
		err = cmd.Execute()
		assert.Error(t, err)
		assert.Regexp(t, "endorsement failure during query", err.Error())
		assert.Regexp(t, fmt.Sprintf("response: status:%d payload:\"%s\"", ccRespStatus[i], ccRespPayload[i]), err.Error())
	}

	
	mockCF, err := getMockChaincodeCmdFactory()
	assert.NoError(t, err, "Error getting mock chaincode command factory")
	mockCF.EndorserClients[0] = common.GetMockEndorserClient(nil, nil)

	cmd := newQueryCmdForTest(mockCF, args)
	err = cmd.Execute()
	assert.Error(t, err)
	assert.Regexp(t, "error during query: received nil proposal response", err.Error())
}

func newQueryCmdForTest(cf *ChaincodeCmdFactory, args []string) *cobra.Command {
	cmd := queryCmd(cf)
	addFlags(cmd)
	cmd.SetArgs(args)
	return cmd
}
