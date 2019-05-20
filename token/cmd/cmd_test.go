/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package token_test

import (
	"testing"

	token "github.com/mcc-github/blockchain/token/cmd"
	"github.com/mcc-github/blockchain/token/cmd/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestAddCommands(t *testing.T) {
	app := kingpin.New("foo", "bar")
	cli := &mocks.CommandRegistrar{}
	configFunc := mock.AnythingOfType("common.CLICommand")
	commands := []string{token.IssueCommand, token.TransferCommand, token.ListTokensCommad, token.RedeemCommand}
	for _, cmd := range commands {
		cli.On("Command", cmd, mock.Anything, configFunc).Return(app.Command(cmd, ""))
	}
	token.AddCommands(cli)
	
	for _, cmd := range commands {
		assert.NotNil(t, app.GetCommand(cmd).GetFlag("config"))
		assert.NotNil(t, app.GetCommand(cmd).GetFlag("channel"))
		assert.NotNil(t, app.GetCommand(cmd).GetFlag("mspPath"))
		assert.NotNil(t, app.GetCommand(cmd).GetFlag("mspId"))
	}

	
	assert.NotNil(t, app.GetCommand(token.IssueCommand).GetFlag("type"))
	assert.NotNil(t, app.GetCommand(token.IssueCommand).GetFlag("quantity"))
	assert.NotNil(t, app.GetCommand(token.IssueCommand).GetFlag("recipient"))

	
	assert.NotNil(t, app.GetCommand(token.TransferCommand).GetFlag("tokenIDs"))
	assert.NotNil(t, app.GetCommand(token.TransferCommand).GetFlag("shares"))

	
	assert.NotNil(t, app.GetCommand(token.RedeemCommand).GetFlag("tokenIDs"))
	assert.NotNil(t, app.GetCommand(token.RedeemCommand).GetFlag("quantity"))
}
