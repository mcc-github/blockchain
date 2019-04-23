/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package clilogging

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/internal/peer/common"
	common2 "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name      string
	args      []string
	shouldErr bool
}

func initLoggingTest(command string) (*cobra.Command, *LoggingCmdFactory) {
	adminClient := common.GetMockAdminClient(nil)
	mockCF := &LoggingCmdFactory{
		AdminClient: adminClient,
		wrapWithEnvelope: func(msg proto.Message) *common2.Envelope {
			pl := &common2.Payload{
				Data: protoutil.MarshalOrPanic(msg),
			}
			env := &common2.Envelope{
				Payload: protoutil.MarshalOrPanic(pl),
			}
			return env
		},
	}
	var cmd *cobra.Command
	if command == "getlevel" {
		cmd = getLevelCmd(mockCF)
	} else if command == "setlevel" {
		cmd = setLevelCmd(mockCF)
	} else if command == "revertlevels" {
		cmd = revertLevelsCmd(mockCF)
	} else if command == "getlogspec" {
		cmd = getLogSpecCmd(mockCF)
	} else if command == "setlogspec" {
		cmd = setLogSpecCmd(mockCF)
	} else {
		
	}
	return cmd, mockCF
}

func runTests(t *testing.T, command string, tc []testCase) {
	cmd, _ := initLoggingTest(command)
	assert := assert.New(t)
	for i := 0; i < len(tc); i++ {
		t.Run(tc[i].name, func(t *testing.T) {
			cmd.SetArgs(tc[i].args)
			err := cmd.Execute()
			if tc[i].shouldErr {
				assert.NotNil(err)
			}
			if !tc[i].shouldErr {
				assert.Nil(err)
			}
		})
	}
}


func TestGetLevel(t *testing.T) {
	var tc []testCase
	tc = append(tc,
		testCase{"NoParameters", []string{}, true},
		testCase{"Valid", []string{"peer"}, false},
	)
	runTests(t, "getlevel", tc)
}


func TestSetLevel(t *testing.T) {
	var tc []testCase
	tc = append(tc,
		testCase{"NoParameters", []string{}, true},
		testCase{"OneParameter", []string{"peer"}, true},
		testCase{"Valid", []string{"peer", "warning"}, false},
		testCase{"InvalidLevel", []string{"peer", "invalidlevel"}, true},
	)
	runTests(t, "setlevel", tc)
}


func TestRevertLevels(t *testing.T) {
	var tc []testCase
	tc = append(tc,
		testCase{"Valid", []string{}, false},
		testCase{"ExtraParameter", []string{"peer"}, true},
	)
	runTests(t, "revertlevels", tc)
}


func TestGetLogSpec(t *testing.T) {
	var tc []testCase
	tc = append(tc,
		testCase{"Valid", []string{}, false},
		testCase{"ExtraParameter", []string{"peer"}, true},
	)
	runTests(t, "getlogspec", tc)
}


func TestSetLogSpec(t *testing.T) {
	var tc []testCase
	tc = append(tc,
		testCase{"NoParameters", []string{}, true},
		testCase{"Valid", []string{"debug"}, false},
	)
	runTests(t, "setlogspec", tc)
}
