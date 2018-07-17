/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery_test

import (
	"testing"

	"github.com/mcc-github/blockchain/discovery/cmd"
	"github.com/mcc-github/blockchain/discovery/cmd/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestAddCommands(t *testing.T) {
	app := kingpin.New("foo", "bar")
	cli := &mocks.CommandRegistrar{}
	configFunc := mock.AnythingOfType("common.CLICommand")
	cli.On("Command", discovery.PeersCommand, mock.Anything, configFunc).Return(app.Command(discovery.PeersCommand, ""))
	cli.On("Command", discovery.ConfigCommand, mock.Anything, configFunc).Return(app.Command(discovery.ConfigCommand, ""))
	cli.On("Command", discovery.EndorsersCommand, mock.Anything, configFunc).Return(app.Command(discovery.EndorsersCommand, ""))
	discovery.AddCommands(cli)
	
	for _, cmd := range []string{discovery.PeersCommand, discovery.ConfigCommand, discovery.EndorsersCommand} {
		assert.NotNil(t, app.GetCommand(cmd).GetFlag("server"))
		assert.NotNil(t, app.GetCommand(cmd).GetFlag("channel"))
	}
	
	assert.NotNil(t, app.GetCommand(discovery.EndorsersCommand).GetFlag("chaincode"))
	assert.NotNil(t, app.GetCommand(discovery.EndorsersCommand).GetFlag("collection"))
}
