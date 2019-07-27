/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	configtxtest "github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	mscc "github.com/mcc-github/blockchain/common/mocks/scc"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	"github.com/mcc-github/blockchain/core/deliverservice"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt/ledgermgmttest"
	"github.com/mcc-github/blockchain/core/ledger/mock"
	ledgermocks "github.com/mcc-github/blockchain/core/ledger/mock"
	"github.com/mcc-github/blockchain/core/transientstore"
	"github.com/mcc-github/blockchain/gossip/gossip"
	gossipmetrics "github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/service"
	gossipservice "github.com/mcc-github/blockchain/gossip/service"
	peergossip "github.com/mcc-github/blockchain/internal/peer/gossip"
	"github.com/mcc-github/blockchain/internal/peer/gossip/mocks"
	"github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	msptesttools.LoadMSPSetupForTesting()
	rc := m.Run()
	os.Exit(rc)
}

func NewTestPeer(t *testing.T) (*Peer, func()) {
	tempdir, err := ioutil.TempDir("", "peer-test")
	require.NoError(t, err, "failed to create temporary directory")

	
	signer := mgmt.GetLocalSigningIdentityOrPanic()
	messageCryptoService := peergossip.NewMCS(&mocks.ChannelPolicyManagerGetter{}, signer, mgmt.NewDeserializersManager())
	secAdv := peergossip.NewSecurityAdvisor(mgmt.NewDeserializersManager())
	defaultSecureDialOpts := func() []grpc.DialOption { return []grpc.DialOption{grpc.WithInsecure()} }
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
	gossipConfig, err := gossip.GlobalConfig("localhost:0", nil)

	gossipService, err := gossipservice.New(
		signer,
		gossipmetrics.NewGossipMetrics(&disabled.Provider{}),
		"localhost:0",
		grpc.NewServer(),
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
	require.NoError(t, err, "failed to create gossip service")

	ledgerMgr, err := constructLedgerMgrWithTestDefaults(filepath.Join(tempdir, "ledgersData"))
	require.NoError(t, err, "failed to create ledger manager")

	peerInstance := &Peer{
		GossipService: gossipService,
		StoreProvider: transientstore.NewStoreProvider(
			filepath.Join(tempdir, "transientstore"),
		),
		LedgerMgr: ledgerMgr,
	}

	cleanup := func() {
		ledgerMgr.Close()
		os.RemoveAll(tempdir)
	}
	return peerInstance, cleanup
}

func TestInitialize(t *testing.T) {
	peerInstance, cleanup := NewTestPeer(t)
	defer cleanup()

	peerInstance.Initialize(
		nil,
		(&mscc.MocksccProviderFactory{}).NewSystemChaincodeProvider(),
		plugin.MapBasedMapper(map[string]validation.PluginFactory{}),
		&ledgermocks.DeployedChaincodeInfoProvider{},
		nil,
		nil,
		runtime.NumCPU(),
	)
}

func TestCreateChannel(t *testing.T) {
	peerInstance, cleanup := NewTestPeer(t)
	defer cleanup()

	var initArg string
	peerInstance.Initialize(
		func(cid string) { initArg = cid },
		(&mscc.MocksccProviderFactory{}).NewSystemChaincodeProvider(),
		plugin.MapBasedMapper(map[string]validation.PluginFactory{}),
		&ledgermocks.DeployedChaincodeInfoProvider{},
		nil,
		nil,
		runtime.NumCPU(),
	)

	testChainID := fmt.Sprintf("mytestchainid-%d", rand.Int())
	block, err := configtxtest.MakeGenesisBlock(testChainID)
	if err != nil {
		fmt.Printf("Failed to create a config block, err %s\n", err)
		t.FailNow()
	}

	err = peerInstance.CreateChannel(block, nil, &mock.DeployedChaincodeInfoProvider{}, nil, nil)
	if err != nil {
		t.Fatalf("failed to create chain %s", err)
	}

	assert.Equal(t, testChainID, initArg)

	
	ledger := peerInstance.GetLedger(testChainID)
	if ledger == nil {
		t.Fatalf("failed to get correct ledger")
	}

	
	block, err = ConfigBlockFromLedger(ledger)
	assert.NoError(t, err, "Failed to get config block from ledger")
	assert.NotNil(t, block, "Config block should not be nil")
	assert.Equal(t, uint64(0), block.Header.Number, "config block should have been block 0")

	
	ledger = peerInstance.GetLedger("BogusChain")
	if ledger != nil {
		t.Fatalf("got a bogus ledger")
	}

	
	pmgr := peerInstance.GetPolicyManager(testChainID)
	if pmgr == nil {
		t.Fatal("failed to get PolicyManager")
	}

	
	pmgr = peerInstance.GetPolicyManager("BogusChain")
	if pmgr != nil {
		t.Fatal("got a bogus PolicyManager")
	}

	channels := peerInstance.GetChannelsInfo()
	if len(channels) != 1 {
		t.Fatalf("incorrect number of channels")
	}
}

func TestDeliverSupportManager(t *testing.T) {
	peerInstance, cleanup := NewTestPeer(t)
	defer cleanup()

	manager := &DeliverChainManager{Peer: peerInstance}

	chainSupport := manager.GetChain("fake")
	assert.Nil(t, chainSupport, "chain support should be nil")

	peerInstance.channels = map[string]*Channel{"testchain": {}}
	chainSupport = manager.GetChain("testchain")
	assert.NotNil(t, chainSupport, "chain support should not be nil")
}

func constructLedgerMgrWithTestDefaults(ledgersDataDir string) (*ledgermgmt.LedgerMgr, error) {
	ledgerInitializer := ledgermgmttest.NewInitializer(ledgersDataDir)
	ledgerInitializer.CustomTxProcessors = map[common.HeaderType]ledger.CustomTxProcessor{
		common.HeaderType_CONFIG: &ConfigTxProcessor{},
	}
	ledgerInitializer.Config.HistoryDBConfig = &ledger.HistoryDBConfig{
		Enabled: true,
	}
	return ledgermgmt.NewLedgerMgr(ledgerInitializer), nil
}
