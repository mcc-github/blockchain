/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deliverservice

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/deliverservice/blocksprovider"
	"github.com/mcc-github/blockchain/core/deliverservice/mocks"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const (
	goRoutineTestWaitTimeout = time.Second * 15
)

var (
	lock = sync.Mutex{}
)

type mockBlocksDelivererFactory struct {
	mockCreate func() (blocksprovider.BlocksDeliverer, error)
}

func (mock *mockBlocksDelivererFactory) Create() (blocksprovider.BlocksDeliverer, error) {
	return mock.mockCreate()
}

type mockMCS struct {
}

func (*mockMCS) Expiration(peerIdentity api.PeerIdentityType) (time.Time, error) {
	return time.Now().Add(time.Hour), nil
}

func (*mockMCS) GetPKIidOfCert(peerIdentity api.PeerIdentityType) common.PKIidType {
	return common.PKIidType("pkiID")
}

func (*mockMCS) VerifyBlock(chainID common.ChainID, seqNum uint64, signedBlock []byte) error {
	return nil
}

func (*mockMCS) Sign(msg []byte) ([]byte, error) {
	return msg, nil
}

func (*mockMCS) Verify(peerIdentity api.PeerIdentityType, signature, message []byte) error {
	return nil
}

func (*mockMCS) VerifyByChannel(chainID common.ChainID, peerIdentity api.PeerIdentityType, signature, message []byte) error {
	return nil
}

func (*mockMCS) ValidateIdentity(peerIdentity api.PeerIdentityType) error {
	return nil
}

func TestNewDeliverService(t *testing.T) {
	defer ensureNoGoroutineLeak(t)()
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64, 1)}
	factory := &struct{ mockBlocksDelivererFactory }{}

	blocksDeliverer := &mocks.MockBlocksDeliverer{}
	blocksDeliverer.MockRecv = mocks.MockRecv

	factory.mockCreate = func() (blocksprovider.BlocksDeliverer, error) {
		return blocksDeliverer, nil
	}
	abcf := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return &mocks.MockAtomicBroadcastClient{
			BD: blocksDeliverer,
		}
	}

	connFactory := func(_ string) func(string) (*grpc.ClientConn, error) {
		return func(endpoint string) (*grpc.ClientConn, error) {
			lock.Lock()
			defer lock.Unlock()
			return newConnection(), nil
		}
	}
	service, err := NewDeliverService(&Config{
		Endpoints:   []string{"a"},
		Gossip:      gossipServiceAdapter,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  abcf,
		ConnFactory: connFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.NoError(t, err)
	assert.NoError(t, service.StartDeliverForChannel("TEST_CHAINID", &mocks.MockLedgerInfo{Height: 0}, func() {}))

	
	assert.Error(t, service.StartDeliverForChannel("TEST_CHAINID", &mocks.MockLedgerInfo{Height: 0}, func() {}), "can't start delivery")
	
	assert.Error(t, service.StopDeliverForChannel("TEST_CHAINID2"), "can't stop delivery")

	
	time.Sleep(time.Second)
	assert.NoError(t, service.StopDeliverForChannel("TEST_CHAINID"))
	time.Sleep(time.Duration(10) * time.Millisecond)
	
	service.Stop()
	time.Sleep(time.Duration(500) * time.Millisecond)
	connWG.Wait()

	assertBlockDissemination(0, gossipServiceAdapter.GossipBlockDisseminations, t)
	assert.Equal(t, blocksDeliverer.RecvCount(), gossipServiceAdapter.AddPayloadCount())
	assert.Error(t, service.StartDeliverForChannel("TEST_CHAINID", &mocks.MockLedgerInfo{Height: 0}, func() {}), "Delivery service is stopping")
	assert.Error(t, service.StopDeliverForChannel("TEST_CHAINID"), "Delivery service is stopping")
}

func TestDeliverServiceRestart(t *testing.T) {
	defer ensureNoGoroutineLeak(t)()
	
	
	

	os := mocks.NewOrderer(5611, t)

	time.Sleep(time.Second)
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}

	service, err := NewDeliverService(&Config{
		Endpoints:   []string{"localhost:5611"},
		Gossip:      gossipServiceAdapter,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.NoError(t, err)

	li := &mocks.MockLedgerInfo{Height: uint64(100)}
	os.SetNextExpectedSeek(uint64(100))

	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")
	
	go os.SendBlock(uint64(100))
	assertBlockDissemination(100, gossipServiceAdapter.GossipBlockDisseminations, t)
	go os.SendBlock(uint64(101))
	assertBlockDissemination(101, gossipServiceAdapter.GossipBlockDisseminations, t)
	go os.SendBlock(uint64(102))
	assertBlockDissemination(102, gossipServiceAdapter.GossipBlockDisseminations, t)
	os.Shutdown()
	time.Sleep(time.Second * 3)
	os = mocks.NewOrderer(5611, t)
	atomic.StoreUint64(&li.Height, uint64(103))
	os.SetNextExpectedSeek(uint64(103))
	go os.SendBlock(uint64(103))
	assertBlockDissemination(103, gossipServiceAdapter.GossipBlockDisseminations, t)
	service.Stop()
	os.Shutdown()
}

func TestDeliverServiceFailover(t *testing.T) {
	defer ensureNoGoroutineLeak(t)()
	
	
	
	
	

	os1 := mocks.NewOrderer(5612, t)
	os2 := mocks.NewOrderer(5613, t)

	time.Sleep(time.Second)
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}

	service, err := NewDeliverService(&Config{
		Endpoints:   []string{"localhost:5612", "localhost:5613"},
		Gossip:      gossipServiceAdapter,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.NoError(t, err)
	li := &mocks.MockLedgerInfo{Height: uint64(100)}
	os1.SetNextExpectedSeek(uint64(100))
	os2.SetNextExpectedSeek(uint64(100))

	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")
	
	go os1.SendBlock(uint64(100))
	instance2fail := os1
	reincarnatedNodePort := 5612
	instance2failSecond := os2
	select {
	case seq := <-gossipServiceAdapter.GossipBlockDisseminations:
		assert.Equal(t, uint64(100), seq)
	case <-time.After(time.Second * 2):
		
		
		os1.Shutdown()
		time.Sleep(time.Second)
		os1 = mocks.NewOrderer(5612, t)
		instance2fail = os2
		instance2failSecond = os1
		reincarnatedNodePort = 5613
		
		
		go os2.SendBlock(uint64(100))
		assertBlockDissemination(100, gossipServiceAdapter.GossipBlockDisseminations, t)
	}

	atomic.StoreUint64(&li.Height, uint64(101))
	os1.SetNextExpectedSeek(uint64(101))
	os2.SetNextExpectedSeek(uint64(101))
	
	instance2fail.Shutdown()
	time.Sleep(time.Second)
	
	go instance2failSecond.SendBlock(uint64(101))
	assertBlockDissemination(101, gossipServiceAdapter.GossipBlockDisseminations, t)
	atomic.StoreUint64(&li.Height, uint64(102))
	
	instance2failSecond.Shutdown()
	time.Sleep(time.Second * 1)
	
	os := mocks.NewOrderer(reincarnatedNodePort, t)
	os.SetNextExpectedSeek(102)
	go os.SendBlock(uint64(102))
	assertBlockDissemination(102, gossipServiceAdapter.GossipBlockDisseminations, t)
	os.Shutdown()
	service.Stop()
}

func TestDeliverServiceUpdateEndpoints(t *testing.T) {
	
	
	
	
	
	defer ensureNoGoroutineLeak(t)()

	os1 := mocks.NewOrderer(5612, t)

	time.Sleep(time.Second)
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}

	service, err := NewDeliverService(&Config{
		Endpoints:   []string{"localhost:5612"},
		Gossip:      gossipServiceAdapter,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	defer service.Stop()

	assert.NoError(t, err)
	li := &mocks.MockLedgerInfo{Height: uint64(100)}
	os1.SetNextExpectedSeek(uint64(100))

	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")

	go os1.SendBlock(uint64(100))
	assertBlockDissemination(100, gossipServiceAdapter.GossipBlockDisseminations, t)

	os2 := mocks.NewOrderer(5613, t)
	defer os2.Shutdown()
	os2.SetNextExpectedSeek(uint64(101))

	service.UpdateEndpoints("TEST_CHAINID", []string{"localhost:5613"})
	
	
	os1.Shutdown()

	atomic.StoreUint64(&li.Height, uint64(101))
	go os2.SendBlock(uint64(101))
	assertBlockDissemination(101, gossipServiceAdapter.GossipBlockDisseminations, t)
}

func TestDeliverServiceServiceUnavailable(t *testing.T) {
	defer ensureNoGoroutineLeak(t)()
	
	
	
	
	
	
	
	

	os1 := mocks.NewOrderer(5615, t)
	os2 := mocks.NewOrderer(5616, t)

	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}

	service, err := NewDeliverService(&Config{
		Endpoints:   []string{"localhost:5615", "localhost:5616"},
		Gossip:      gossipServiceAdapter,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.NoError(t, err)
	li := &mocks.MockLedgerInfo{Height: 100}
	os1.SetNextExpectedSeek(li.Height)
	os2.SetNextExpectedSeek(li.Height)

	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")

	waitForConnectionToSomeOSN := func() (*mocks.Orderer, *mocks.Orderer) {
		for {
			if os1.ConnCount() > 0 {
				return os1, os2
			}
			if os2.ConnCount() > 0 {
				return os2, os1
			}
			time.Sleep(time.Millisecond * 100)
		}
	}

	activeInstance, backupInstance := waitForConnectionToSomeOSN()
	assert.NotNil(t, activeInstance)
	assert.NotNil(t, backupInstance)
	
	assert.Equal(t, activeInstance.ConnCount(), 1)
	
	assert.Equal(t, backupInstance.ConnCount(), 0)

	
	go activeInstance.SendBlock(li.Height)

	assertBlockDissemination(li.Height, gossipServiceAdapter.GossipBlockDisseminations, t)
	li.Height++

	
	backupInstance.SetNextExpectedSeek(li.Height)
	
	backupInstance.SendBlock(li.Height)

	
	activeInstance.Fail()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-time.After(time.Millisecond * 100):
				if backupInstance.ConnCount() > 0 {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	wg.Wait()
	assert.NoError(t, ctx.Err(), "Delivery client has not failed over to alive ordering service")
	
	assert.Equal(t, backupInstance.ConnCount(), 1)
	
	assertBlockDissemination(li.Height, gossipServiceAdapter.GossipBlockDisseminations, t)

	
	time.Sleep(time.Millisecond * 1600)

	li.Height++
	activeInstance.Resurrect()
	backupInstance.Fail()

	resurrectCtx, resCancel := context.WithTimeout(context.Background(), time.Second)
	defer resCancel()

	go func() {
		
		activeInstance.SetNextExpectedSeek(li.Height)
		
		activeInstance.SendBlock(li.Height)

	}()

	reswg := sync.WaitGroup{}
	reswg.Add(1)

	go func() {
		defer reswg.Done()
		for {
			select {
			case <-time.After(time.Millisecond * 100):
				if activeInstance.ConnCount() > 0 {
					return
				}
			case <-resurrectCtx.Done():
				return
			}
		}
	}()

	reswg.Wait()

	assert.NoError(t, resurrectCtx.Err(), "Delivery client has not failed over to alive ordering service")
	
	assert.Equal(t, activeInstance.ConnCount(), 1)
	
	assertBlockDissemination(li.Height, gossipServiceAdapter.GossipBlockDisseminations, t)

	
	os1.Shutdown()
	os2.Shutdown()
	service.Stop()
}

func TestDeliverServiceAbruptStop(t *testing.T) {
	defer ensureNoGoroutineLeak(t)()
	
	
	
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}
	service, err := NewDeliverService(&Config{
		Endpoints:   []string{"a"},
		Gossip:      gossipServiceAdapter,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.NoError(t, err)

	li := &mocks.MockLedgerInfo{Height: uint64(100)}
	service.StartDeliverForChannel("mychannel", li, func() {})
	service.StopDeliverForChannel("mychannel")
}

func TestDeliverServiceShutdown(t *testing.T) {
	defer ensureNoGoroutineLeak(t)()
	
	
	os := mocks.NewOrderer(5614, t)

	time.Sleep(time.Second)
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}

	service, err := NewDeliverService(&Config{
		Endpoints:   []string{"localhost:5614"},
		Gossip:      gossipServiceAdapter,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.NoError(t, err)

	li := &mocks.MockLedgerInfo{Height: uint64(100)}
	os.SetNextExpectedSeek(uint64(100))
	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")

	
	go os.SendBlock(uint64(100))
	assertBlockDissemination(100, gossipServiceAdapter.GossipBlockDisseminations, t)
	go os.SendBlock(uint64(101))
	assertBlockDissemination(101, gossipServiceAdapter.GossipBlockDisseminations, t)
	atomic.StoreUint64(&li.Height, uint64(102))
	os.SetNextExpectedSeek(uint64(102))
	
	service.Stop()
	go os.SendBlock(uint64(102))
	select {
	case <-gossipServiceAdapter.GossipBlockDisseminations:
		assert.Fail(t, "Disseminated a block after shutting down the delivery service")
	case <-time.After(time.Second * 2):
	}
	os.Shutdown()
	time.Sleep(time.Second)
}

func TestDeliverServiceShutdownRespawn(t *testing.T) {
	
	
	
	viper.Set("peer.deliveryclient.reconnectTotalTimeThreshold", time.Second)
	defer viper.Reset()
	defer ensureNoGoroutineLeak(t)()

	osn1 := mocks.NewOrderer(5614, t)

	time.Sleep(time.Second)
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}

	service, err := NewDeliverService(&Config{
		Endpoints:   []string{"localhost:5614", "localhost:5615"},
		Gossip:      gossipServiceAdapter,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.NoError(t, err)

	li := &mocks.MockLedgerInfo{Height: uint64(100)}
	osn1.SetNextExpectedSeek(uint64(100))
	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")

	
	go osn1.SendBlock(uint64(100))
	assertBlockDissemination(100, gossipServiceAdapter.GossipBlockDisseminations, t)
	go osn1.SendBlock(uint64(101))
	assertBlockDissemination(101, gossipServiceAdapter.GossipBlockDisseminations, t)
	atomic.StoreUint64(&li.Height, uint64(102))
	
	time.Sleep(time.Second * 2)
	
	osn2 := mocks.NewOrderer(5615, t)
	
	osn1.Shutdown()
	
	osn2.SetNextExpectedSeek(uint64(102))
	go osn2.SendBlock(uint64(102))
	
	assertBlockDissemination(102, gossipServiceAdapter.GossipBlockDisseminations, t)
	service.Stop()
	osn2.Shutdown()
}

func TestDeliverServiceDisconnectReconnect(t *testing.T) {
	
	
	
	
	
	
	
	
	viper.Set("peer.deliveryclient.reconnectTotalTimeThreshold", time.Second*2)
	defer viper.Reset()
	defer ensureNoGoroutineLeak(t)()

	osn := mocks.NewOrderer(5614, t)

	time.Sleep(time.Second)
	gossipServiceAdapter := &mocks.MockGossipServiceAdapter{GossipBlockDisseminations: make(chan uint64)}

	service, err := NewDeliverService(&Config{
		Endpoints:   []string{"localhost:5614"},
		Gossip:      gossipServiceAdapter,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.NoError(t, err)

	li := &mocks.MockLedgerInfo{Height: uint64(100)}
	osn.SetNextExpectedSeek(uint64(100))
	err = service.StartDeliverForChannel("TEST_CHAINID", li, func() {})
	assert.NoError(t, err, "can't start delivery")

	
	go osn.SendBlock(uint64(100))
	assertBlockDissemination(100, gossipServiceAdapter.GossipBlockDisseminations, t)
	go osn.SendBlock(uint64(101))
	assertBlockDissemination(101, gossipServiceAdapter.GossipBlockDisseminations, t)
	atomic.StoreUint64(&li.Height, uint64(102))

	for i := 0; i < 5; i += 1 {
		
		osn.Shutdown()
		
		assert.True(t, waitForConnectionCount(osn, 0), "deliverService can't disconnect from orderer")
		
		osn = mocks.NewOrderer(5614, t)
		osn.SetNextExpectedSeek(atomic.LoadUint64(&li.Height))
		
		assert.True(t, waitForConnectionCount(osn, 1), "deliverService can't reconnect to orderer")
	}

	
	go osn.SendBlock(uint64(102))
	
	assertBlockDissemination(102, gossipServiceAdapter.GossipBlockDisseminations, t)
	service.Stop()
	osn.Shutdown()
}

func TestDeliverServiceBadConfig(t *testing.T) {
	
	service, err := NewDeliverService(&Config{
		Endpoints:   []string{},
		Gossip:      &mocks.MockGossipServiceAdapter{},
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.Error(t, err)
	assert.Nil(t, service)

	
	service, err = NewDeliverService(&Config{
		Endpoints:   []string{"a"},
		Gossip:      nil,
		CryptoSvc:   &mockMCS{},
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.Error(t, err)
	assert.Nil(t, service)

	
	service, err = NewDeliverService(&Config{
		Endpoints:   []string{"a"},
		Gossip:      &mocks.MockGossipServiceAdapter{},
		CryptoSvc:   nil,
		ABCFactory:  DefaultABCFactory,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.Error(t, err)
	assert.Nil(t, service)

	
	service, err = NewDeliverService(&Config{
		Endpoints:   []string{"a"},
		Gossip:      &mocks.MockGossipServiceAdapter{},
		CryptoSvc:   &mockMCS{},
		ABCFactory:  nil,
		ConnFactory: DefaultConnectionFactory,
		Signer:      &mocks.SignerSerializer{},
	})
	assert.Error(t, err)
	assert.Nil(t, service)

	
	service, err = NewDeliverService(&Config{
		Endpoints:  []string{"a"},
		Gossip:     &mocks.MockGossipServiceAdapter{},
		CryptoSvc:  &mockMCS{},
		ABCFactory: DefaultABCFactory,
		Signer:     &mocks.SignerSerializer{},
	})
	assert.Error(t, err)
	assert.Nil(t, service)
}

func TestRetryPolicyOverflow(t *testing.T) {
	connFactory := func(channelID string) func(endpoint string) (*grpc.ClientConn, error) {
		return func(_ string) (*grpc.ClientConn, error) {
			return nil, errors.New("")
		}
	}
	client := (&deliverServiceImpl{conf: &Config{ConnFactory: connFactory}}).newClient("TEST", &mocks.MockLedgerInfo{Height: uint64(100)})
	assert.NotNil(t, client.shouldRetry)
	for i := 0; i < 100; i++ {
		retryTime, _ := client.shouldRetry(i, time.Second)
		assert.True(t, retryTime <= time.Hour && retryTime > 0)
	}
}

func assertBlockDissemination(expectedSeq uint64, ch chan uint64, t *testing.T) {
	select {
	case seq := <-ch:
		assert.Equal(t, expectedSeq, seq)
	case <-time.After(time.Second * 5):
		assert.FailNow(t, fmt.Sprintf("Didn't gossip a new block with seq num %d within a timely manner", expectedSeq))
		t.Fatal()
	}
}

func ensureNoGoroutineLeak(t *testing.T) func() {
	goroutineCountAtStart := runtime.NumGoroutine()
	return func() {
		start := time.Now()
		timeLimit := start.Add(goRoutineTestWaitTimeout)
		for time.Now().Before(timeLimit) {
			time.Sleep(time.Millisecond * 500)
			if goroutineCountAtStart >= runtime.NumGoroutine() {
				return
			}
		}
		assert.Fail(t, "Some goroutine(s) didn't finish: %s", getStackTrace())
	}
}

func getStackTrace() string {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	return string(buf)
}

func waitForConnectionCount(orderer *mocks.Orderer, connCount int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for {
		select {
		case <-time.After(time.Millisecond * 100):
			if orderer.ConnCount() == connCount {
				return true
			}
		case <-ctx.Done():
			return false
		}
	}
}

func TestToEndpointCriteria(t *testing.T) {
	for _, testCase := range []struct {
		description string
		input       ConnectionCriteria
		expectedOut []comm.EndpointCriteria
	}{
		{
			description: "globally defined endpoints",
			input: ConnectionCriteria{
				Organizations:    []string{"foo", "bar"},
				OrdererEndpoints: []string{"a", "b", "c"},
			},
			expectedOut: []comm.EndpointCriteria{
				{Organizations: []string{"foo", "bar"}, Endpoint: "a"},
				{Organizations: []string{"foo", "bar"}, Endpoint: "b"},
				{Organizations: []string{"foo", "bar"}, Endpoint: "c"},
			},
		},
		{
			description: "per org defined endpoints",
			input: ConnectionCriteria{
				Organizations: []string{"foo", "bar"},
				
				OrdererEndpoints: []string{"a", "b", "c"},
				OrdererEndpointsByOrg: map[string][]string{
					"foo": {"a", "b"},
					"bar": {"c"},
				},
			},
			expectedOut: []comm.EndpointCriteria{
				{Organizations: []string{"foo"}, Endpoint: "a"},
				{Organizations: []string{"foo"}, Endpoint: "b"},
				{Organizations: []string{"bar"}, Endpoint: "c"},
			},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			assert.Equal(t, testCase.expectedOut, testCase.input.toEndpointCriteria())
		})
	}
}
