/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"fmt"
	"math/rand"
	"net"
	"testing"

	configtxtest "github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/common/localmsp"
	mscc "github.com/mcc-github/blockchain/common/mocks/scc"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/committer/txvalidator"
	"github.com/mcc-github/blockchain/core/deliverservice"
	"github.com/mcc-github/blockchain/core/deliverservice/blocksprovider"
	"github.com/mcc-github/blockchain/core/handlers/validation/api"
	ledgermocks "github.com/mcc-github/blockchain/core/ledger/mock"
	"github.com/mcc-github/blockchain/core/mocks/ccprovider"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/service"
	"github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/msp/mgmt/testtools"
	peergossip "github.com/mcc-github/blockchain/peer/gossip"
	"github.com/mcc-github/blockchain/peer/gossip/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type mockDeliveryClient struct {
}

func (ds *mockDeliveryClient) UpdateEndpoints(chainID string, endpoints []string) error {
	return nil
}



func (ds *mockDeliveryClient) StartDeliverForChannel(chainID string, ledgerInfo blocksprovider.LedgerInfo, f func()) error {
	return nil
}



func (ds *mockDeliveryClient) StopDeliverForChannel(chainID string) error {
	return nil
}


func (*mockDeliveryClient) Stop() {
}

type mockDeliveryClientFactory struct {
}

func (*mockDeliveryClientFactory) Service(g service.GossipService, endpoints []string, mcs api.MessageCryptoService) (deliverclient.DeliverService, error) {
	return &mockDeliveryClient{}, nil
}

func TestNewPeerServer(t *testing.T) {
	server, err := NewPeerServer(":4050", comm.ServerConfig{})
	assert.NoError(t, err, "NewPeerServer returned unexpected error")
	assert.Equal(t, "[::]:4050", server.Address(), "NewPeerServer returned the wrong address")
	server.Stop()

	_, err = NewPeerServer("", comm.ServerConfig{})
	assert.Error(t, err, "expected NewPeerServer to return error with missing address")
}

func TestInitChain(t *testing.T) {
	chainId := "testChain"
	chainInitializer = func(cid string) {
		assert.Equal(t, chainId, cid, "chainInitializer received unexpected cid")
	}
	InitChain(chainId)
}

func TestInitialize(t *testing.T) {
	cleanup := setupPeerFS(t)
	defer cleanup()

	Initialize(nil, &ccprovider.MockCcProviderImpl{}, (&mscc.MocksccProviderFactory{}).NewSystemChaincodeProvider(), txvalidator.MapBasedPluginMapper(map[string]validation.PluginFactory{}), nil, &ledgermocks.DeployedChaincodeInfoProvider{})
}

func TestCreateChainFromBlock(t *testing.T) {
	cleanup := setupPeerFS(t)
	defer cleanup()

	Initialize(nil, &ccprovider.MockCcProviderImpl{}, (&mscc.MocksccProviderFactory{}).NewSystemChaincodeProvider(), txvalidator.MapBasedPluginMapper(map[string]validation.PluginFactory{}), &platforms.Registry{}, &ledgermocks.DeployedChaincodeInfoProvider{})
	testChainID := fmt.Sprintf("mytestchainid-%d", rand.Int())
	block, err := configtxtest.MakeGenesisBlock(testChainID)
	if err != nil {
		fmt.Printf("Failed to create a config block, err %s\n", err)
		t.FailNow()
	}

	
	grpcServer := grpc.NewServer()
	socket, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	msptesttools.LoadMSPSetupForTesting()

	identity, _ := mgmt.GetLocalSigningIdentityOrPanic().Serialize()
	messageCryptoService := peergossip.NewMCS(&mocks.ChannelPolicyManagerGetter{}, localmsp.NewSigner(), mgmt.NewDeserializersManager())
	secAdv := peergossip.NewSecurityAdvisor(mgmt.NewDeserializersManager())
	var defaultSecureDialOpts = func() []grpc.DialOption {
		var dialOpts []grpc.DialOption
		dialOpts = append(dialOpts, grpc.WithInsecure())
		return dialOpts
	}
	err = service.InitGossipServiceCustomDeliveryFactory(
		identity, socket.Addr().String(), grpcServer, nil,
		&mockDeliveryClientFactory{},
		messageCryptoService, secAdv, defaultSecureDialOpts)

	assert.NoError(t, err)

	go grpcServer.Serve(socket)
	defer grpcServer.Stop()

	err = CreateChainFromBlock(block, nil, nil)
	if err != nil {
		t.Fatalf("failed to create chain %s", err)
	}

	
	ledger := GetLedger(testChainID)
	if ledger == nil {
		t.Fatalf("failed to get correct ledger")
	}

	
	block, err = getCurrConfigBlockFromLedger(ledger)
	assert.NoError(t, err, "Failed to get config block from ledger")
	assert.NotNil(t, block, "Config block should not be nil")
	assert.Equal(t, uint64(0), block.Header.Number, "config block should have been block 0")

	
	ledger = GetLedger("BogusChain")
	if ledger != nil {
		t.Fatalf("got a bogus ledger")
	}

	
	block = GetCurrConfigBlock(testChainID)
	if block == nil {
		t.Fatalf("failed to get correct block")
	}

	cfgSupport := configSupport{}
	chCfg := cfgSupport.GetChannelConfig(testChainID)
	assert.NotNil(t, chCfg, "failed to get channel config")

	
	block = GetCurrConfigBlock("BogusBlock")
	if block != nil {
		t.Fatalf("got a bogus block")
	}

	
	pmgr := GetPolicyManager(testChainID)
	if pmgr == nil {
		t.Fatal("failed to get PolicyManager")
	}

	
	pmgr = GetPolicyManager("BogusChain")
	if pmgr != nil {
		t.Fatal("got a bogus PolicyManager")
	}

	
	pmg := NewChannelPolicyManagerGetter()
	assert.NotNil(t, pmg, "PolicyManagerGetter should not be nil")

	pmgr, ok := pmg.Manager(testChainID)
	assert.NotNil(t, pmgr, "PolicyManager should not be nil")
	assert.Equal(t, true, ok, "expected Manage() to return true")

	SetCurrConfigBlock(block, testChainID)

	channels := GetChannelsInfo()
	if len(channels) != 1 {
		t.Fatalf("incorrect number of channels")
	}

	
	chains.Lock()
	chains.list = map[string]*chain{}
	chains.Unlock()
}

func TestGetLocalIP(t *testing.T) {
	ip := GetLocalIP()
	t.Log(ip)
}

func TestDeliverSupportManager(t *testing.T) {
	
	MockInitialize()

	manager := &DeliverChainManager{}
	chainSupport, ok := manager.GetChain("fake")
	assert.Nil(t, chainSupport, "chain support should be nil")
	assert.False(t, ok, "Should not find fake channel")

	MockCreateChain("testchain")
	chainSupport, ok = manager.GetChain("testchain")
	assert.NotNil(t, chainSupport, "chain support should not be nil")
	assert.True(t, ok, "Should find testchain channel")
}
