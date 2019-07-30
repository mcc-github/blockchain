/*
Copyright Digital Asset Holdings, LLC. All Rights Reserved.
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/mcc-github/blockchain/internal/configtxgen/configtxgentest"
	"github.com/mcc-github/blockchain/internal/configtxgen/encoder"
	genesisconfig "github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/internal/peer/chaincode/mock"
	"github.com/mcc-github/blockchain/internal/peer/common"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



type signerSerializer interface {
	identity.SignerSerializer
}



type deliver interface {
	pb.Deliver_DeliverClient
}



type peerDeliverClient interface {
	pb.DeliverClient
}

func TestCheckChaincodeCmdParamsWithNewCallingSchema(t *testing.T) {
	chaincodeCtorJSON = `{ "Args":["func", "param"] }`
	chaincodePath = "some/path"
	chaincodeName = "somename"
	require := require.New(t)
	result := checkChaincodeCmdParams(&cobra.Command{})

	require.Nil(result)
}

func TestCheckChaincodeCmdParamsWithOldCallingSchema(t *testing.T) {
	chaincodeCtorJSON = `{ "Function":"func", "Args":["param"] }`
	chaincodePath = "some/path"
	chaincodeName = "somename"
	require := require.New(t)
	result := checkChaincodeCmdParams(&cobra.Command{})

	require.Nil(result)
}

func TestCheckChaincodeCmdParamsWithoutName(t *testing.T) {
	chaincodeCtorJSON = `{ "Function":"func", "Args":["param"] }`
	chaincodePath = "some/path"
	chaincodeName = ""
	require := require.New(t)
	result := checkChaincodeCmdParams(&cobra.Command{})

	require.Error(result)
}

func TestCheckChaincodeCmdParamsWithFunctionOnly(t *testing.T) {
	chaincodeCtorJSON = `{ "Function":"func" }`
	chaincodePath = "some/path"
	chaincodeName = "somename"
	require := require.New(t)
	result := checkChaincodeCmdParams(&cobra.Command{})

	require.Error(result)
}

func TestCheckChaincodeCmdParamsEmptyCtor(t *testing.T) {
	chaincodeCtorJSON = `{}`
	chaincodePath = "some/path"
	chaincodeName = "somename"
	require := require.New(t)
	result := checkChaincodeCmdParams(&cobra.Command{})

	require.Error(result)
}

func TestCheckValidJSON(t *testing.T) {
	validJSON := `{"Args":["a","b","c"]}`
	input := &chaincodeInput{}
	if err := json.Unmarshal([]byte(validJSON), &input); err != nil {
		t.Fail()
		t.Logf("Chaincode argument error: %s", err)
		return
	}

	validJSON = `{"Function":"f", "Args":["a","b","c"]}`
	if err := json.Unmarshal([]byte(validJSON), &input); err != nil {
		t.Fail()
		t.Logf("Chaincode argument error: %s", err)
		return
	}

	validJSON = `{"Function":"f", "Args":[]}`
	if err := json.Unmarshal([]byte(validJSON), &input); err != nil {
		t.Fail()
		t.Logf("Chaincode argument error: %s", err)
		return
	}

	validJSON = `{"Function":"f"}`
	if err := json.Unmarshal([]byte(validJSON), &input); err != nil {
		t.Fail()
		t.Logf("Chaincode argument error: %s", err)
		return
	}
}

func TestCheckInvalidJSON(t *testing.T) {
	invalidJSON := `{["a","b","c"]}`
	input := &chaincodeInput{}
	if err := json.Unmarshal([]byte(invalidJSON), &input); err == nil {
		t.Fail()
		t.Logf("Bar argument error should have been caught: %s", invalidJSON)
		return
	}

	invalidJSON = `{"Function":}`
	if err := json.Unmarshal([]byte(invalidJSON), &input); err == nil {
		t.Fail()
		t.Logf("Chaincode argument error: %s", err)
		t.Logf("Bar argument error should have been caught: %s", invalidJSON)
		return
	}
}

func TestGetOrdererEndpointFromConfigTx(t *testing.T) {
	signer, err := common.GetDefaultSigner()
	assert.NoError(t, err)

	mockchain := "mockchain"
	factory.InitFactories(nil)
	config := configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile)
	pgen := encoder.New(config)
	genesisBlock := pgen.GenesisBlockForChannel(mockchain)

	mockResponse := &pb.ProposalResponse{
		Response:    &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(genesisBlock)},
		Endorsement: &pb.Endorsement{},
	}
	mockEndorserClient := common.GetMockEndorserClient(mockResponse, nil)

	ordererEndpoints, err := common.GetOrdererEndpointOfChain(mockchain, signer, mockEndorserClient)
	assert.NoError(t, err, "GetOrdererEndpointOfChain from genesis block")

	assert.Equal(t, len(ordererEndpoints), 1)
	assert.Equal(t, ordererEndpoints[0], "127.0.0.1:7050")
}

func TestGetOrdererEndpointFail(t *testing.T) {
	signer, err := common.GetDefaultSigner()
	assert.NoError(t, err)

	mockchain := "mockchain"
	factory.InitFactories(nil)

	mockResponse := &pb.ProposalResponse{
		Response:    &pb.Response{Status: 404, Payload: []byte{}},
		Endorsement: &pb.Endorsement{},
	}
	mockEndorserClient := common.GetMockEndorserClient(mockResponse, nil)

	_, err = common.GetOrdererEndpointOfChain(mockchain, signer, mockEndorserClient)
	assert.Error(t, err, "GetOrdererEndpointOfChain from invalid response")
}

const sampleCollectionConfigGood = `[
	{
		"name": "foo",
		"policy": "OR('A.member', 'B.member')",
		"requiredPeerCount": 3,
		"maxPeerCount": 483279847,
		"blockToLive":10,
		"memberOnlyRead": true,
		"memberOnlyWrite": true
	}
]`

const sampleCollectionConfigBad = `[
	{
		"name": "foo",
		"policy": "barf",
		"requiredPeerCount": 3,
		"maxPeerCount": 483279847
	}
]`

func TestCollectionParsing(t *testing.T) {
	ccp, ccpBytes, err := getCollectionConfigFromBytes([]byte(sampleCollectionConfigGood))
	assert.NoError(t, err)
	assert.NotNil(t, ccp)
	assert.NotNil(t, ccpBytes)
	conf := ccp.Config[0].GetStaticCollectionConfig()
	pol, _ := cauthdsl.FromString("OR('A.member', 'B.member')")
	assert.Equal(t, 3, int(conf.RequiredPeerCount))
	assert.Equal(t, 483279847, int(conf.MaximumPeerCount))
	assert.Equal(t, "foo", conf.Name)
	assert.True(t, proto.Equal(pol, conf.MemberOrgsPolicy.GetSignaturePolicy()))
	assert.Equal(t, 10, int(conf.BlockToLive))
	assert.Equal(t, true, conf.MemberOnlyRead)
	t.Logf("conf=%s", conf)

	ccp, ccpBytes, err = getCollectionConfigFromBytes([]byte(sampleCollectionConfigBad))
	assert.Error(t, err)
	assert.Nil(t, ccp)
	assert.Nil(t, ccpBytes)

	ccp, ccpBytes, err = getCollectionConfigFromBytes([]byte("barf"))
	assert.Error(t, err)
	assert.Nil(t, ccp)
	assert.Nil(t, ccpBytes)
}

func TestValidatePeerConnectionParams(t *testing.T) {
	defer resetFlags()
	defer viper.Reset()
	assert := assert.New(t)
	cleanup := configtest.SetDevFabricConfigPath(t)
	defer cleanup()

	
	viper.Set("peer.tls.enabled", false)

	
	resetFlags()
	peerAddresses = []string{"peer0", "peer1"}
	tlsRootCertFiles = []string{"cert0", "cert1"}
	err := validatePeerConnectionParameters("query")
	assert.Error(err)
	assert.Contains(err.Error(), "command can only be executed against one peer")

	
	
	resetFlags()
	peerAddresses = []string{"peer0"}
	err = validatePeerConnectionParameters("query")
	assert.NoError(err)
	assert.Nil(tlsRootCertFiles)

	
	
	resetFlags()
	peerAddresses = []string{"peer0"}
	tlsRootCertFiles = []string{"cert0", "cert1"}
	err = validatePeerConnectionParameters("invoke")
	assert.NoError(err)
	assert.Nil(tlsRootCertFiles)

	
	
	resetFlags()
	peerAddresses = []string{"peer0", "peer1"}
	err = validatePeerConnectionParameters("invoke")
	assert.NoError(err)
	assert.Nil(tlsRootCertFiles)

	
	viper.Set("peer.tls.enabled", true)

	
	
	resetFlags()
	peerAddresses = []string{"peer0", "peer1"}
	tlsRootCertFiles = []string{"cert0"}
	err = validatePeerConnectionParameters("invoke")
	assert.Error(err)
	assert.Contains(err.Error(), fmt.Sprintf("number of peer addresses (%d) does not match the number of TLS root cert files (%d)", len(peerAddresses), len(tlsRootCertFiles)))

	
	
	resetFlags()
	peerAddresses = []string{"peer0", "peer1"}
	tlsRootCertFiles = []string{"cert0", "cert1"}
	err = validatePeerConnectionParameters("invoke")
	assert.NoError(err)

	
	resetFlags()
	connectionProfile = "blah"
	err = validatePeerConnectionParameters("invoke")
	assert.Error(err)
	assert.Contains(err.Error(), "error reading connection profile")

	
	
	resetFlags()
	channelID = "mychannel"
	connectionProfile = "testdata/connectionprofile-uneven.yaml"
	err = validatePeerConnectionParameters("invoke")
	assert.Error(err)
	assert.Contains(err.Error(), "defined in the channel config but doesn't have associated peer config")

	
	resetFlags()
	channelID = "mychannel"
	connectionProfile = "testdata/connectionprofile.yaml"
	err = validatePeerConnectionParameters("invoke")
	assert.NoError(err)
}

func TestInitCmdFactoryFailures(t *testing.T) {
	defer resetFlags()
	assert := assert.New(t)

	
	resetFlags()
	peerAddresses = []string{"peer0", "peer1"}
	tlsRootCertFiles = []string{"cert0", "cert1"}
	cf, err := InitCmdFactory("query", true, false)
	assert.Error(err)
	assert.Contains(err.Error(), "error validating peer connection parameters: 'query' command can only be executed against one peer")
	assert.Nil(cf)

	
	resetFlags()
	peerAddresses = []string{}
	cf, err = InitCmdFactory("query", true, false)
	assert.Error(err)
	assert.Contains(err.Error(), "no endorser clients retrieved")
	assert.Nil(cf)

	
	
	resetFlags()
	peerAddresses = nil
	cf, err = InitCmdFactory("invoke", false, true)
	assert.Error(err)
	assert.Contains(err.Error(), "no ordering endpoint or endorser client supplied")
	assert.Nil(cf)
}

func TestDeliverGroupConnect(t *testing.T) {
	defer resetFlags()
	g := NewGomegaWithT(t)

	
	mockDeliverClients := []*DeliverClient{
		{
			Client:  getMockDeliverClientResponseWithTxStatusAndID(pb.TxValidationCode_VALID, "txid0"),
			Address: "peer0",
		},
		{
			Client:  getMockDeliverClientResponseWithTxStatusAndID(pb.TxValidationCode_VALID, "txid0"),
			Address: "peer1",
		},
	}
	dg := DeliverGroup{
		Clients:   mockDeliverClients,
		ChannelID: "testchannel",
		Signer:    &mock.SignerSerializer{},
		Certificate: tls.Certificate{
			Certificate: [][]byte{[]byte("test")},
		},
		TxID: "txid0",
	}
	err := dg.Connect(context.Background())
	g.Expect(err).To(BeNil())

	
	mockDC := &mock.PeerDeliverClient{}
	mockDC.DeliverFilteredReturns(nil, errors.New("icecream"))
	mockDeliverClients = []*DeliverClient{
		{
			Client:  mockDC,
			Address: "peer0",
		},
	}
	dg = DeliverGroup{
		Clients:   mockDeliverClients,
		ChannelID: "testchannel",
		Signer:    &mock.SignerSerializer{},
		Certificate: tls.Certificate{
			Certificate: [][]byte{[]byte("test")},
		},
		TxID: "txid0",
	}
	err = dg.Connect(context.Background())
	g.Expect(err.Error()).To(ContainSubstring("error connecting to deliver filtered"))
	g.Expect(err.Error()).To(ContainSubstring("icecream"))

	
	mockD := &mock.Deliver{}
	mockD.SendReturns(errors.New("blah"))
	mockDC.DeliverFilteredReturns(mockD, nil)
	mockDeliverClients = []*DeliverClient{
		{
			Client:  mockDC,
			Address: "peer0",
		},
	}
	dg = DeliverGroup{
		Clients:   mockDeliverClients,
		ChannelID: "testchannel",
		Signer:    &mock.SignerSerializer{},
		Certificate: tls.Certificate{
			Certificate: [][]byte{[]byte("test")},
		},
		TxID: "txid0",
	}
	err = dg.Connect(context.Background())
	g.Expect(err.Error()).To(ContainSubstring("error sending deliver seek info"))
	g.Expect(err.Error()).To(ContainSubstring("blah"))

	
	delayChan := make(chan struct{})
	mockDCDelay := getMockDeliverClientRegisterAfterDelay(delayChan)
	mockDeliverClients = []*DeliverClient{
		{
			Client:  mockDCDelay,
			Address: "peer0",
		},
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()
	dg = DeliverGroup{
		Clients:   mockDeliverClients,
		ChannelID: "testchannel",
		Signer:    &mock.SignerSerializer{},
		Certificate: tls.Certificate{
			Certificate: [][]byte{[]byte("test")},
		},
		TxID: "txid0",
	}
	err = dg.Connect(ctx)
	g.Expect(err.Error()).To(ContainSubstring("timed out waiting for connection to deliver on all peers"))
	close(delayChan)
}

func TestDeliverGroupWait(t *testing.T) {
	defer resetFlags()
	g := NewGomegaWithT(t)

	
	mockConn := &mock.Deliver{}
	filteredResp := &pb.DeliverResponse{
		Type: &pb.DeliverResponse_FilteredBlock{FilteredBlock: createFilteredBlock(pb.TxValidationCode_VALID, "txid0")},
	}
	mockConn.RecvReturns(filteredResp, nil)
	mockDeliverClients := []*DeliverClient{
		{
			Connection: mockConn,
			Address:    "peer0",
		},
	}
	dg := DeliverGroup{
		Clients:   mockDeliverClients,
		ChannelID: "testchannel",
		Signer:    &mock.SignerSerializer{},
		Certificate: tls.Certificate{
			Certificate: [][]byte{[]byte("test")},
		},
		TxID: "txid0",
	}
	err := dg.Wait(context.Background())
	g.Expect(err).To(BeNil())

	
	mockConn = &mock.Deliver{}
	mockConn.RecvReturns(nil, errors.New("avocado"))
	mockDeliverClients = []*DeliverClient{
		{
			Connection: mockConn,
			Address:    "peer0",
		},
	}
	dg = DeliverGroup{
		Clients:   mockDeliverClients,
		ChannelID: "testchannel",
		Signer:    &mock.SignerSerializer{},
		Certificate: tls.Certificate{
			Certificate: [][]byte{[]byte("test")},
		},
		TxID: "txid0",
	}
	err = dg.Wait(context.Background())
	g.Expect(err.Error()).To(ContainSubstring("error receiving from deliver filtered"))
	g.Expect(err.Error()).To(ContainSubstring("avocado"))

	
	mockConn = &mock.Deliver{}
	resp := &pb.DeliverResponse{
		Type: &pb.DeliverResponse_Block{},
	}
	mockConn.RecvReturns(resp, nil)
	mockDeliverClients = []*DeliverClient{
		{
			Connection: mockConn,
			Address:    "peer0",
		},
	}
	dg = DeliverGroup{
		Clients:   mockDeliverClients,
		ChannelID: "testchannel",
		Signer:    &mock.SignerSerializer{},
		Certificate: tls.Certificate{
			Certificate: [][]byte{[]byte("test")},
		},
		TxID: "txid0",
	}
	err = dg.Wait(context.Background())
	g.Expect(err.Error()).To(ContainSubstring("unexpected response type"))

	
	mockConn = &mock.Deliver{}
	mockConn.RecvReturns(nil, errors.New("barbeque"))
	mockConn2 := &mock.Deliver{}
	mockConn2.RecvReturns(nil, errors.New("tofu"))
	mockDeliverClients = []*DeliverClient{
		{
			Connection: mockConn,
			Address:    "peerBBQ",
		},
		{
			Connection: mockConn2,
			Address:    "peerTOFU",
		},
	}
	dg = DeliverGroup{
		Clients:   mockDeliverClients,
		ChannelID: "testchannel",
		Signer:    &mock.SignerSerializer{},
		Certificate: tls.Certificate{
			Certificate: [][]byte{[]byte("test")},
		},
		TxID: "txid0",
	}
	err = dg.Wait(context.Background())
	g.Expect(err.Error()).To(SatisfyAny(
		ContainSubstring("barbeque"),
		ContainSubstring("tofu")))
}

func TestChaincodeInvokeOrQuery_waitForEvent(t *testing.T) {
	defer resetFlags()

	waitForEvent = true
	mockCF, err := getMockChaincodeCmdFactory()
	assert.NoError(t, err)
	peerAddresses = []string{"peer0", "peer1"}
	channelID := "testchannel"
	txID := "txid0"

	t.Run("success - deliver clients returns event with expected txid", func(t *testing.T) {
		_, err = ChaincodeInvokeOrQuery(
			&pb.ChaincodeSpec{},
			channelID,
			txID,
			true,
			mockCF.Signer,
			mockCF.Certificate,
			mockCF.EndorserClients,
			mockCF.DeliverClients,
			mockCF.BroadcastClient,
		)
		assert.NoError(t, err)
	})

	t.Run("success - one deliver client first receives block without txid and then one with txid", func(t *testing.T) {
		filteredBlocks := []*pb.FilteredBlock{
			createFilteredBlock(pb.TxValidationCode_VALID, "theseare", "notthetxidsyouarelookingfor"),
			createFilteredBlock(pb.TxValidationCode_VALID, "txid0"),
		}
		mockDCTwoBlocks := getMockDeliverClientRespondsWithFilteredBlocks(filteredBlocks)
		mockDC := getMockDeliverClientResponseWithTxStatusAndID(pb.TxValidationCode_VALID, "txid0")
		mockDeliverClients := []pb.DeliverClient{mockDCTwoBlocks, mockDC}

		_, err = ChaincodeInvokeOrQuery(
			&pb.ChaincodeSpec{},
			channelID,
			txID,
			true,
			mockCF.Signer,
			mockCF.Certificate,
			mockCF.EndorserClients,
			mockDeliverClients,
			mockCF.BroadcastClient,
		)
		assert.NoError(t, err)
	})

	t.Run("failure - one of the deliver clients returns error", func(t *testing.T) {
		mockDCErr := getMockDeliverClientWithErr("moist")
		mockDC := getMockDeliverClientResponseWithTxStatusAndID(pb.TxValidationCode_VALID, "txid0")
		mockDeliverClients := []pb.DeliverClient{mockDCErr, mockDC}

		_, err = ChaincodeInvokeOrQuery(
			&pb.ChaincodeSpec{},
			channelID,
			txID,
			true,
			mockCF.Signer,
			mockCF.Certificate,
			mockCF.EndorserClients,
			mockDeliverClients,
			mockCF.BroadcastClient,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "moist")
	})

	t.Run("failure - transaction committed with non-success validation code", func(t *testing.T) {
		mockDC := getMockDeliverClientResponseWithTxStatusAndID(pb.TxValidationCode_VALID, "txid0")
		mockDCFail := getMockDeliverClientResponseWithTxStatusAndID(pb.TxValidationCode_ENDORSEMENT_POLICY_FAILURE, "txid0")
		mockDeliverClients := []pb.DeliverClient{mockDCFail, mockDC}

		_, err = ChaincodeInvokeOrQuery(
			&pb.ChaincodeSpec{},
			channelID,
			txID,
			true,
			mockCF.Signer,
			mockCF.Certificate,
			mockCF.EndorserClients,
			mockDeliverClients,
			mockCF.BroadcastClient,
		)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "transaction invalidated with status (ENDORSEMENT_POLICY_FAILURE)")
	})

	t.Run("failure - deliver returns response status instead of block", func(t *testing.T) {
		mockDC := &mock.PeerDeliverClient{}
		mockDF := &mock.Deliver{}
		resp := &pb.DeliverResponse{
			Type: &pb.DeliverResponse_Status{
				Status: cb.Status_FORBIDDEN,
			},
		}
		mockDF.RecvReturns(resp, nil)
		mockDC.DeliverFilteredReturns(mockDF, nil)
		mockDeliverClients := []pb.DeliverClient{mockDC}
		_, err = ChaincodeInvokeOrQuery(
			&pb.ChaincodeSpec{},
			channelID,
			txID,
			true,
			mockCF.Signer,
			mockCF.Certificate,
			mockCF.EndorserClients,
			mockDeliverClients,
			mockCF.BroadcastClient,
		)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "deliver completed with status (FORBIDDEN) before txid received")
	})

	t.Run(" failure - timeout occurs - both deliver clients don't return an event with the expected txid before timeout", func(t *testing.T) {
		delayChan := make(chan struct{})
		mockDCDelay := getMockDeliverClientRespondAfterDelay(delayChan, pb.TxValidationCode_VALID, "txid0")
		mockDeliverClients := []pb.DeliverClient{mockDCDelay, mockDCDelay}
		waitForEventTimeout = 10 * time.Millisecond

		_, err = ChaincodeInvokeOrQuery(
			&pb.ChaincodeSpec{},
			channelID,
			txID,
			true,
			mockCF.Signer,
			mockCF.Certificate,
			mockCF.EndorserClients,
			mockDeliverClients,
			mockCF.BroadcastClient,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timed out")
		close(delayChan)
	})
}

func TestProcessProposals(t *testing.T) {
	
	mockClients := []pb.EndorserClient{}
	for i := 2; i <= 5; i++ {
		response := &pb.ProposalResponse{
			Response:    &pb.Response{Status: int32(i * 100)},
			Endorsement: &pb.Endorsement{},
		}
		mockClients = append(mockClients, common.GetMockEndorserClient(response, nil))
	}
	mockErrorClient := common.GetMockEndorserClient(nil, errors.New("failed to call endorser"))
	signedProposal := &pb.SignedProposal{}
	t.Run("should process a proposal for a single peer", func(t *testing.T) {
		responses, err := processProposals([]pb.EndorserClient{mockClients[0]}, signedProposal)
		assert.NoError(t, err)
		assert.Len(t, responses, 1)
		assert.Equal(t, responses[0].Response.Status, int32(200))
	})
	t.Run("should process a proposal for multiple peers", func(t *testing.T) {
		responses, err := processProposals(mockClients, signedProposal)
		assert.NoError(t, err)
		assert.Len(t, responses, 4)
		
		statuses := []int32{}
		for _, response := range responses {
			statuses = append(statuses, response.Response.Status)
		}
		sort.Slice(statuses, func(i, j int) bool { return statuses[i] < statuses[j] })
		assert.EqualValues(t, []int32{200, 300, 400, 500}, statuses)
	})
	t.Run("should return an error from processing a proposal for a single peer", func(t *testing.T) {
		responses, err := processProposals([]pb.EndorserClient{mockErrorClient}, signedProposal)
		assert.EqualError(t, err, "failed to call endorser")
		assert.Nil(t, responses)
	})
	t.Run("should return an error from processing a proposal for a single peer within multiple peers", func(t *testing.T) {
		responses, err := processProposals([]pb.EndorserClient{mockClients[0], mockErrorClient, mockClients[1]}, signedProposal)
		assert.EqualError(t, err, "failed to call endorser")
		assert.Nil(t, responses)
	})
}
