/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery_test

import (
	"bytes"
	"testing"

	"github.com/mcc-github/blockchain/cmd/common"
	. "github.com/mcc-github/blockchain/discovery/client"
	"github.com/mcc-github/blockchain/discovery/cmd"
	"github.com/mcc-github/blockchain/discovery/cmd/mocks"
	discprotos "github.com/mcc-github/blockchain/protos/discovery"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEndorserCmd(t *testing.T) {
	server := "peer0"
	channel := "mychannel"
	stub := &mocks.Stub{}
	parser := &mocks.ResponseParser{}

	t.Run("no server supplied", func(t *testing.T) {
		cmd := discovery.NewEndorsersCmd(stub, parser)
		cmd.SetChannel(&channel)

		err := cmd.Execute(common.Config{})
		assert.Equal(t, err.Error(), "no server specified")
	})

	t.Run("no channel supplied", func(t *testing.T) {
		cmd := discovery.NewEndorsersCmd(stub, parser)
		cmd.SetServer(&server)

		err := cmd.Execute(common.Config{})
		assert.Equal(t, err.Error(), "no channel specified")
	})

	t.Run("Endorsement query with no chaincodes", func(t *testing.T) {
		cmd := discovery.NewEndorsersCmd(stub, parser)
		cmd.SetServer(&server)
		cmd.SetChannel(&channel)

		err := cmd.Execute(common.Config{})
		assert.Contains(t, err.Error(), "invocation chain should not be empty")
	})

	t.Run("Server return error", func(t *testing.T) {
		chaincodes := []string{"mycc"}
		cmd := discovery.NewEndorsersCmd(stub, parser)
		cmd.SetChannel(&channel)
		cmd.SetServer(&server)
		cmd.SetChaincodes(&chaincodes)
		stub.On("Send", server, mock.Anything, mock.Anything).Return(nil, errors.New("deadline exceeded")).Once()

		err := cmd.Execute(common.Config{})
		assert.Contains(t, err.Error(), "deadline exceeded")
	})

	t.Run("Endorsement query with no collections succeeds", func(t *testing.T) {
		chaincodes := []string{"mycc"}
		cmd := discovery.NewEndorsersCmd(stub, parser)
		cmd.SetChannel(&channel)
		cmd.SetServer(&server)
		cmd.SetChaincodes(&chaincodes)
		parser.On("ParseResponse", channel, mock.Anything).Return(nil).Once()
		stub.On("Send", server, mock.Anything, mock.Anything).Return(nil, nil).Once()

		err := cmd.Execute(common.Config{})
		assert.NoError(t, err)
	})

	t.Run("Endorsement query with collections succeeds", func(t *testing.T) {
		chaincodes := []string{"mycc", "yourcc"}
		collections := map[string]string{
			"mycc": "col1,col2",
		}

		stub := &mocks.Stub{}
		cmd := discovery.NewEndorsersCmd(stub, parser)
		cmd.SetChannel(&channel)
		cmd.SetServer(&server)
		cmd.SetChaincodes(&chaincodes)
		cmd.SetCollections(&collections)
		parser.On("ParseResponse", channel, mock.Anything).Return(nil).Once()
		stub.On("Send", server, mock.Anything, mock.Anything).Return(nil, nil).Once().Run(func(arg mock.Arguments) {
			
			req := arg.Get(2).(*Request)
			
			assert.Equal(t, "mycc", req.Queries[0].GetCcQuery().Interests[0].Chaincodes[0].Name)
			
			assert.Equal(t, []string{"col1", "col2"}, req.Queries[0].GetCcQuery().Interests[0].Chaincodes[0].CollectionNames)
		})

		err := cmd.Execute(common.Config{})
		assert.NoError(t, err)
		stub.AssertNumberOfCalls(t, "Send", 1)
	})

	t.Run("Endorsement query with collections that aren't mapped to any chaincode(s)", func(t *testing.T) {
		chaincodes := []string{"mycc", "yourcc"}
		collections := map[string]string{
			"mycc":  "col1,col2",
			"ourcc": "col3",
		}

		stub := &mocks.Stub{}
		cmd := discovery.NewEndorsersCmd(stub, parser)
		cmd.SetChannel(&channel)
		cmd.SetServer(&server)
		cmd.SetChaincodes(&chaincodes)
		cmd.SetCollections(&collections)
		stub.On("Send", server, mock.Anything, mock.Anything).Return(nil, nil).Once()

		err := cmd.Execute(common.Config{})
		assert.Contains(t, err.Error(), "a collection specified chaincode ourcc but it wasn't specified with a chaincode flag")
	})

	t.Run("Endorsement query with collections that aren't mapped to any chaincode(s)", func(t *testing.T) {
		chaincodes := []string{"mycc", "yourcc"}
		collections := map[string]string{
			"mycc":  "col1,col2",
			"ourcc": "col3",
		}

		stub := &mocks.Stub{}
		cmd := discovery.NewEndorsersCmd(stub, parser)
		cmd.SetChannel(&channel)
		cmd.SetServer(&server)
		cmd.SetChaincodes(&chaincodes)
		cmd.SetCollections(&collections)
		stub.On("Send", server, mock.Anything, mock.Anything).Return(nil, nil).Once()

		err := cmd.Execute(common.Config{})
		assert.Contains(t, err.Error(), "a collection specified chaincode ourcc but it wasn't specified with a chaincode flag")
	})
}

func TestParseEndorsementResponse(t *testing.T) {
	buff := &bytes.Buffer{}
	parser := &discovery.EndorserResponseParser{Writer: buff}
	res := &mocks.ServiceResponse{}

	t.Run("Server returns empty response", func(t *testing.T) {
		res.On("Raw").Return(&discprotos.Response{}).Once()
		err := parser.ParseResponse("mychannel", res)
		assert.Contains(t, err.Error(), "empty results")
	})

	t.Run("Server returns an error", func(t *testing.T) {
		res.On("Raw").Return(&discprotos.Response{
			Results: []*discprotos.QueryResult{
				{
					Result: &discprotos.QueryResult_Error{
						Error: &discprotos.Error{
							Content: "internal error",
						},
					},
				},
			},
		}).Once()
		err := parser.ParseResponse("mychannel", res)
		assert.Contains(t, err.Error(), "internal error")
	})

	t.Run("Server returns a response with the wrong type", func(t *testing.T) {
		res.On("Raw").Return(&discprotos.Response{
			Results: []*discprotos.QueryResult{
				{
					Result: &discprotos.QueryResult_Members{
						Members: &discprotos.PeerMembershipResult{PeersByOrg: map[string]*discprotos.Peers{
							"Org1MSP": {},
						}},
					},
				},
			},
		}).Once()
		err := parser.ParseResponse("mychannel", res)
		assert.Contains(t, err.Error(), "server returned response of unexpected type: *discovery.QueryResult")
	})

	t.Run("Server returns a proper response", func(t *testing.T) {
		res.On("Raw").Return(&discprotos.Response{
			Results: []*discprotos.QueryResult{
				{
					Result: endorsersResponse,
				},
			},
		}).Once()
		err := parser.ParseResponse("mychannel", res)
		assert.NoError(t, err)
	})
}

var endorsersResponse = &discprotos.QueryResult_CcQueryRes{
	CcQueryRes: &discprotos.ChaincodeQueryResult{
		Content: []*discprotos.EndorsementDescriptor{
			{
				EndorsersByGroups: map[string]*discprotos.Peers{
					"Org1MSP": {
						Peers: []*discprotos.Peer{
							{
								Identity:       []byte("identity"),
								StateInfo:      stateInfoMessage(100).Envelope,
								MembershipInfo: aliveMessage(0).Envelope,
							},
						},
					},
				},
				Layouts: []*discprotos.Layout{
					{
						QuantitiesByGroup: map[string]uint32{
							"Org1MSP": 2,
						},
					},
				},
			},
		},
	},
}
