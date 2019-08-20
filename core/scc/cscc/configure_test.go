/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cscc

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/sw"
	configtxtest "github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/common/genesis"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/chaincode"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/deliverservice"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt/ledgermgmttest"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/core/scc/cscc/mocks"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/gossip/gossip"
	gossipmetrics "github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/service"
	"github.com/mcc-github/blockchain/internal/configtxgen/configtxgentest"
	"github.com/mcc-github/blockchain/internal/configtxgen/encoder"
	genesisconfig "github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	peergossip "github.com/mcc-github/blockchain/internal/peer/gossip"
	"github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)



type aclProvider interface {
	aclmgmt.ACLProvider
}



type chaincodeStub interface {
	shim.ChaincodeStubInterface
}



type channelPolicyManagerGetter interface {
	policies.ChannelPolicyManagerGetter
}



type policyChecker interface {
	policy.PolicyChecker
}



type storeProvider interface {
	transientstore.StoreProvider
}

func TestMain(m *testing.M) {
	msptesttools.LoadMSPSetupForTesting()
	rc := m.Run()
	os.Exit(rc)

}

func TestInvokedChaincodeName(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		name, err := InvokedChaincodeName(validSignedProposal().ProposalBytes)
		assert.NoError(t, err)
		assert.Equal(t, "cscc", name)
	})

	t.Run("BadProposalBytes", func(t *testing.T) {
		_, err := InvokedChaincodeName([]byte("garbage"))
		assert.EqualError(t, err, "could not unmarshal proposal: proto: can't skip unknown wire type 7")
	})

	t.Run("BadChaincodeProposalBytes", func(t *testing.T) {
		_, err := InvokedChaincodeName(protoutil.MarshalOrPanic(&pb.Proposal{
			Payload: []byte("garbage"),
		}))
		assert.EqualError(t, err, "could not unmarshal chaincode proposal payload: proto: can't skip unknown wire type 7")
	})

	t.Run("BadChaincodeInvocationSpec", func(t *testing.T) {
		_, err := InvokedChaincodeName(protoutil.MarshalOrPanic(&pb.Proposal{
			Payload: protoutil.MarshalOrPanic(&pb.ChaincodeProposalPayload{
				Input: []byte("garbage"),
			}),
		}))
		assert.EqualError(t, err, "could not unmarshal chaincode invocation spec: proto: can't skip unknown wire type 7")
	})

	t.Run("NilChaincodeSpec", func(t *testing.T) {
		_, err := InvokedChaincodeName(protoutil.MarshalOrPanic(&pb.Proposal{
			Payload: protoutil.MarshalOrPanic(&pb.ChaincodeProposalPayload{
				Input: protoutil.MarshalOrPanic(&pb.ChaincodeInvocationSpec{}),
			}),
		}))
		assert.EqualError(t, err, "chaincode spec is nil")
	})

	t.Run("NilChaincodeID", func(t *testing.T) {
		_, err := InvokedChaincodeName(protoutil.MarshalOrPanic(&pb.Proposal{
			Payload: protoutil.MarshalOrPanic(&pb.ChaincodeProposalPayload{
				Input: protoutil.MarshalOrPanic(&pb.ChaincodeInvocationSpec{
					ChaincodeSpec: &pb.ChaincodeSpec{},
				}),
			}),
		}))
		assert.EqualError(t, err, "chaincode id is nil")
	})
}

func TestConfigerInit(t *testing.T) {
	mockACLProvider := &mocks.ACLProvider{}
	mockStub := &mocks.ChaincodeStub{}
	cscc := &PeerConfiger{
		aclProvider: mockACLProvider,
	}
	res := cscc.Init(mockStub)
	assert.Equal(t, int32(shim.OK), res.Status)
}

func TestConfigerInvokeInvalidParameters(t *testing.T) {
	mockACLProvider := &mocks.ACLProvider{}
	cscc := &PeerConfiger{
		aclProvider: mockACLProvider,
	}
	mockStub := &mocks.ChaincodeStub{}

	mockStub.GetArgsReturns(nil)
	mockStub.GetSignedProposalReturns(validSignedProposal(), nil)
	res := cscc.Invoke(mockStub)
	assert.NotEqual(
		t,
		int32(shim.OK),
		res.Status,
		"cscc invoke expected to fail having zero arguments",
	)
	assert.Equal(t, "Incorrect number of arguments, 0", res.Message)

	mockACLProvider.CheckACLReturns(errors.New("Failed authorization"))
	args := [][]byte{[]byte("GetChannels")}
	mockStub.GetArgsReturns(args)
	res = cscc.Invoke(mockStub)
	assert.NotEqual(
		t,
		int32(shim.OK),
		res.Status,
		"invoke expected to fail no signed proposal provided",
	)
	assert.Equal(t, "access denied for [GetChannels]: Failed authorization", res.Message)

	mockACLProvider.CheckACLReturns(nil)
	args = [][]byte{[]byte("fooFunction"), []byte("testChainID")}
	mockStub.GetArgsReturns(args)
	res = cscc.Invoke(mockStub)
	assert.NotEqual(
		t,
		int32(shim.OK),
		res.Status,
		"invoke invoke expected wrong function name provided",
	)
	assert.Equal(t, "Requested function fooFunction not found.", res.Message)

	mockACLProvider.CheckACLReturns(nil)
	args = [][]byte{[]byte("GetConfigBlock"), []byte("testChainID")}
	mockStub.GetArgsReturns(args)
	mockStub.GetSignedProposalReturns(&pb.SignedProposal{
		ProposalBytes: []byte("garbage"),
	}, nil)
	res = cscc.Invoke(mockStub)
	assert.NotEqual(
		t,
		int32(shim.OK),
		res.Status,
		"invoke expected to fail in ccc2cc context",
	)
	assert.Equal(
		t,
		"Could not identify the called chaincode: [could not unmarshal proposal: proto: can't skip unknown wire type 7]",
		res.Message,
	)

	mockACLProvider.CheckACLReturns(nil)
	args = [][]byte{[]byte("GetConfigBlock"), []byte("testChainID")}
	mockStub.GetArgsReturns(args)
	mockStub.GetSignedProposalReturns(&pb.SignedProposal{
		ProposalBytes: protoutil.MarshalOrPanic(&pb.Proposal{
			Payload: protoutil.MarshalOrPanic(&pb.ChaincodeProposalPayload{
				Input: protoutil.MarshalOrPanic(&pb.ChaincodeInvocationSpec{
					ChaincodeSpec: &pb.ChaincodeSpec{
						ChaincodeId: &pb.ChaincodeID{
							Name: "fake-cc2cc",
						},
					},
				}),
			}),
		}),
	}, nil)
	res = cscc.Invoke(mockStub)
	assert.NotEqual(
		t,
		int32(shim.OK),
		res.Status,
		"invoke expected to fail in ccc2cc context",
	)
	assert.Equal(
		t,
		"Cannot invoke CSCC from another chaincode, original invocation for 'fake-cc2cc'",
		res.Message,
	)

	mockACLProvider.CheckACLReturns(errors.New("Failed authorization"))
	mockStub.GetSignedProposalReturns(validSignedProposal(), nil)
	args = [][]byte{[]byte("GetConfigBlock"), []byte("testChainID")}
	mockStub.GetArgsReturns(args)
	res = cscc.Invoke(mockStub)
	assert.NotEqual(
		t,
		int32(shim.OK),
		res.Status,
		"invoke expected to fail no signed proposal provided",
	)
	assert.Equal(
		t,
		"access denied for [GetConfigBlock][testChainID]: Failed authorization",
		res.Message,
	)
}

func TestConfigerInvokeJoinChainMissingParams(t *testing.T) {
	cscc := &PeerConfiger{
		aclProvider: &mocks.ACLProvider{},
	}
	mockStub := &mocks.ChaincodeStub{}
	mockStub.GetArgsReturns([][]byte{[]byte("JoinChain")})
	res := cscc.Invoke(mockStub)
	assert.NotEqual(
		t,
		int32(shim.OK),
		res.Status,
		"cscc invoke JoinChain should have failed with invalid number of args",
	)
}

func TestConfigerInvokeJoinChainWrongParams(t *testing.T) {
	cscc := &PeerConfiger{
		aclProvider: &mocks.ACLProvider{},
	}
	mockStub := &mocks.ChaincodeStub{}
	mockStub.GetArgsReturns([][]byte{[]byte("JoinChain"), []byte("action")})
	mockStub.GetSignedProposalReturns(validSignedProposal(), nil)
	res := cscc.Invoke(mockStub)
	assert.NotEqual(
		t,
		int32(shim.OK),
		res.Status,
		"cscc invoke JoinChain should have failed with null genesis block",
	)
}

func TestConfigerInvokeJoinChainCorrectParams(t *testing.T) {
	viper.Set("chaincode.executetimeout", "3s")

	testDir, err := ioutil.TempDir("", "cscc_test")
	require.NoError(t, err, "error in creating test dir")
	defer os.Remove(testDir)

	ledgerMgr := ledgermgmt.NewLedgerMgr(ledgermgmttest.NewInitializer(testDir))
	defer ledgerMgr.Close()

	peerEndpoint := "127.0.0.1:13611"

	config := chaincode.GlobalConfig()
	config.StartupTimeout = 30 * time.Second

	grpcServer := grpc.NewServer()
	socket, err := net.Listen("tcp", peerEndpoint)
	require.NoError(t, err)

	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)
	signer := mgmt.GetLocalSigningIdentityOrPanic()

	messageCryptoService := peergossip.NewMCS(&mocks.ChannelPolicyManagerGetter{}, signer, mgmt.NewDeserializersManager(), cryptoProvider)
	secAdv := peergossip.NewSecurityAdvisor(mgmt.NewDeserializersManager())
	var defaultSecureDialOpts = func() []grpc.DialOption {
		var dialOpts []grpc.DialOption
		dialOpts = append(dialOpts, grpc.WithInsecure())
		return dialOpts
	}
	defaultDeliverClientDialOpts := []grpc.DialOption{grpc.WithBlock()}
	defaultDeliverClientDialOpts = append(
		defaultDeliverClientDialOpts,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(comm.MaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(comm.MaxSendMsgSize)))
	defaultDeliverClientDialOpts = append(
		defaultDeliverClientDialOpts,
		comm.ClientKeepaliveOptions(comm.DefaultKeepaliveOptions)...,
	)
	gossipConfig, err := gossip.GlobalConfig(peerEndpoint, nil)
	assert.NoError(t, err)

	gossipService, err := service.New(
		signer,
		gossipmetrics.NewGossipMetrics(&disabled.Provider{}),
		peerEndpoint,
		grpcServer,
		messageCryptoService,
		secAdv,
		defaultSecureDialOpts,
		nil,
		defaultDeliverClientDialOpts,
		gossipConfig,
		&service.ServiceConfig{},
		&deliverservice.DeliverServiceConfig{
			ReConnectBackoffThreshold:   deliverservice.DefaultReConnectBackoffThreshold,
			ReconnectTotalTimeThreshold: deliverservice.DefaultReConnectTotalTimeThreshold,
		},
	)
	assert.NoError(t, err)

	go grpcServer.Serve(socket)
	defer grpcServer.Stop()

	assert.NoError(t, err)
	
	mockACLProvider := &mocks.ACLProvider{}
	cscc := &PeerConfiger{
		policyChecker: &mocks.PolicyChecker{},
		aclProvider:   mockACLProvider,
		peer: &peer.Peer{
			StoreProvider:  &mocks.StoreProvider{},
			GossipService:  gossipService,
			LedgerMgr:      ledgerMgr,
			CryptoProvider: cryptoProvider,
		},
	}
	mockStub := &mocks.ChaincodeStub{}

	
	blockBytes := mockConfigBlock()
	if blockBytes == nil {
		t.Fatalf("cscc invoke JoinChain failed because invalid block")
	}
	args := [][]byte{[]byte("JoinChain"), blockBytes}
	sProp := validSignedProposal()
	sProp.Signature = sProp.ProposalBytes

	
	mockStub.GetArgsReturns([][]byte{[]byte("JoinChain"), nil})
	mockStub.GetSignedProposalReturns(sProp, nil)
	res := cscc.Invoke(mockStub)
	
	assert.Equal(t, int32(shim.ERROR), res.Status)

	
	payload, _ := proto.Marshal(&cb.Payload{})
	env, _ := proto.Marshal(&cb.Envelope{
		Payload: payload,
	})
	badBlock := &cb.Block{
		Data: &cb.BlockData{
			Data: [][]byte{env},
		},
	}
	badBlockBytes := protoutil.MarshalOrPanic(badBlock)
	mockStub.GetArgsReturns([][]byte{[]byte("JoinChain"), badBlockBytes})
	res = cscc.Invoke(mockStub)
	
	assert.Equal(t, int32(shim.ERROR), res.Status)

	
	mockStub.GetArgsReturns(args)
	mockStub.GetSignedProposalReturns(sProp, nil)
	res = cscc.Invoke(mockStub)
	assert.Equal(t, int32(shim.OK), res.Status, "invoke JoinChain failed with: %v", res.Message)

	
	sProp.Signature = nil
	mockACLProvider.CheckACLReturns(errors.New("Failed authorization"))
	mockStub.GetArgsReturns(args)
	mockStub.GetSignedProposalReturns(sProp, nil)

	res = cscc.Invoke(mockStub)
	assert.Equal(t, int32(shim.ERROR), res.Status)
	assert.Contains(t, res.Message, "access denied for [JoinChain][mytestchainid]")
	sProp.Signature = sProp.ProposalBytes

	
	
	chainID, err := protoutil.GetChainIDFromBlockBytes(blockBytes)
	if err != nil {
		t.Fatalf("cscc invoke JoinChain failed with: %v", err)
	}

	
	mockACLProvider.CheckACLReturns(errors.New("Failed authorization"))
	args = [][]byte{[]byte("GetConfigBlock"), []byte(chainID)}
	mockStub.GetArgsReturns(args)
	mockStub.GetSignedProposalReturns(sProp, nil)
	res = cscc.Invoke(mockStub)
	assert.Equal(t, int32(shim.ERROR), res.Status, "invoke GetConfigBlock should have failed: %v", res.Message)
	assert.Contains(t, res.Message, "Failed authorization")

	
	mockACLProvider.CheckACLReturns(nil)
	res = cscc.Invoke(mockStub)
	assert.Equal(t, int32(shim.OK), res.Status, "invoke GetConfigBlock failed with: %v", res.Message)

	
	mockACLProvider.CheckACLReturns(nil)
	args = [][]byte{[]byte(GetChannels)}
	mockStub.GetArgsReturns(args)
	res = cscc.Invoke(mockStub)
	if res.Status != shim.OK {
		t.FailNow()
	}

	cqr := &pb.ChannelQueryResponse{}
	err = proto.Unmarshal(res.Payload, cqr)
	if err != nil {
		t.FailNow()
	}

	
	if len(cqr.GetChannels()) != 1 {
		t.FailNow()
	}
}

func TestPeerConfiger_SubmittingOrdererGenesis(t *testing.T) {
	conf := configtxgentest.Load(genesisconfig.SampleSingleMSPSoloProfile)
	conf.Application = nil
	cg, err := encoder.NewChannelGroup(conf)
	assert.NoError(t, err)
	block := genesis.NewFactoryImpl(cg).Block("mytestchainid")
	blockBytes := protoutil.MarshalOrPanic(block)

	mockACLProvider := &mocks.ACLProvider{}
	cscc := &PeerConfiger{
		aclProvider: mockACLProvider,
	}
	mockStub := &mocks.ChaincodeStub{}
	
	args := [][]byte{[]byte("JoinChain"), []byte(blockBytes)}
	mockStub.GetArgsReturns(args)
	mockStub.GetSignedProposalReturns(validSignedProposal(), nil)
	res := cscc.Invoke(mockStub)

	assert.NotEqual(
		t,
		int32(shim.OK),
		res.Status,
		"invoke JoinChain should have failed with wrong genesis block",
	)
	assert.Contains(t, res.Message, "missing Application configuration group")
}

func mockConfigBlock() []byte {
	var blockBytes []byte = nil
	block, err := configtxtest.MakeGenesisBlock("mytestchainid")
	if err == nil {
		blockBytes = protoutil.MarshalOrPanic(block)
	}
	return blockBytes
}

func validSignedProposal() *pb.SignedProposal {
	return &pb.SignedProposal{
		ProposalBytes: protoutil.MarshalOrPanic(&pb.Proposal{
			Payload: protoutil.MarshalOrPanic(&pb.ChaincodeProposalPayload{
				Input: protoutil.MarshalOrPanic(&pb.ChaincodeInvocationSpec{
					ChaincodeSpec: &pb.ChaincodeSpec{
						ChaincodeId: &pb.ChaincodeID{
							Name: "cscc",
						},
					},
				}),
			}),
		}),
	}
}
