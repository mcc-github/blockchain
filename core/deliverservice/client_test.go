/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deliverclient

import (
	"crypto/sha256"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/deliverservice/blocksprovider"
	"github.com/mcc-github/blockchain/core/deliverservice/mocks"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	connWG sync.WaitGroup
)

func newConnection() *grpc.ClientConn {
	
	
	cc, _ := grpc.Dial("", grpc.WithInsecure(), grpc.WithBalancer(&balancer{}))
	return cc
}

type balancer struct{}

func (*balancer) Start(target string, config grpc.BalancerConfig) error {
	connWG.Add(1)
	return nil
}

func (*balancer) Up(addr grpc.Address) (down func(error)) {
	return func(error) {}
}

func (*balancer) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	return grpc.Address{}, func() {}, errors.New("")
}

func (*balancer) Notify() <-chan []grpc.Address {
	return nil
}

func (*balancer) Close() error {
	connWG.Done()
	return nil
}

type blocksDelivererConsumer func(blocksprovider.BlocksDeliverer) error

var blockDelivererConsumerWithRecv = func(bd blocksprovider.BlocksDeliverer) error {
	_, err := bd.Recv()
	return err
}

var blockDelivererConsumerWithSend = func(bd blocksprovider.BlocksDeliverer) error {
	return bd.Send(&common.Envelope{})
}

type abc struct {
	shouldFail bool
	grpc.ClientStream
}

func (a *abc) Send(*common.Envelope) error {
	if a.shouldFail {
		return errors.New("Failed sending")
	}
	return nil
}

func (a *abc) Recv() (*orderer.DeliverResponse, error) {
	if a.shouldFail {
		return nil, errors.New("Failed Recv")
	}
	return &orderer.DeliverResponse{}, nil
}

type abclient struct {
	shouldFail bool
	stream     *abc
}

func (ac *abclient) Broadcast(ctx context.Context, opts ...grpc.CallOption) (orderer.AtomicBroadcast_BroadcastClient, error) {
	panic("Not implemented")
}

func (ac *abclient) Deliver(ctx context.Context, opts ...grpc.CallOption) (orderer.AtomicBroadcast_DeliverClient, error) {
	if ac.stream != nil {
		return ac.stream, nil
	}
	if ac.shouldFail {
		return nil, errors.New("Failed creating ABC")
	}
	return &abc{}, nil
}

type connProducer struct {
	shouldFail      bool
	connAttempts    int
	connTime        time.Duration
	ordererEndpoint string
}

func (cp *connProducer) realConnection() (*grpc.ClientConn, string, error) {
	cc, err := grpc.Dial(cp.ordererEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, "", err
	}
	return cc, cp.ordererEndpoint, nil
}

func (cp *connProducer) NewConnection() (*grpc.ClientConn, string, error) {
	time.Sleep(cp.connTime)
	cp.connAttempts++
	if cp.ordererEndpoint != "" {
		return cp.realConnection()
	}
	if cp.shouldFail {
		return nil, "", errors.New("Failed connecting")
	}
	return newConnection(), "localhost:5611", nil
}



func (cp *connProducer) UpdateEndpoints(endpoints []string) {
	panic("Not implemented")
}

func (cp *connProducer) GetEndpoints() []string {
	panic("Not implemented")
}

func (cp *connProducer) DisableEndpoint(endpoint string) {
	panic("Not implemented")
}

func TestOrderingServiceConnFailure(t *testing.T) {
	testOrderingServiceConnFailure(t, blockDelivererConsumerWithRecv)
	testOrderingServiceConnFailure(t, blockDelivererConsumerWithSend)
	connWG.Wait()
}

func testOrderingServiceConnFailure(t *testing.T, bdc blocksDelivererConsumer) {
	
	
	cp := &connProducer{shouldFail: true}
	clFactory := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return &abclient{}
	}
	setupInvoked := 0
	setup := func(blocksprovider.BlocksDeliverer) error {
		setupInvoked++
		return nil
	}
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		
		
		
		
		if attemptNum == 1 {
			cp.shouldFail = false
		}

		return time.Duration(0), attemptNum < 2
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	defer bc.Close()
	err := bdc(bc)
	assert.NoError(t, err)
	assert.Equal(t, 2, cp.connAttempts)
	assert.Equal(t, 1, setupInvoked)
}

func TestOrderingServiceStreamFailure(t *testing.T) {
	testOrderingServiceStreamFailure(t, blockDelivererConsumerWithRecv)
	testOrderingServiceStreamFailure(t, blockDelivererConsumerWithSend)
	connWG.Wait()
}

func testOrderingServiceStreamFailure(t *testing.T, bdc blocksDelivererConsumer) {
	
	
	
	abcClient := &abclient{shouldFail: true}
	cp := &connProducer{}
	clFactory := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return abcClient
	}
	setupInvoked := 0
	setup := func(blocksprovider.BlocksDeliverer) error {
		setupInvoked++
		return nil
	}
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		
		
		if attemptNum == 1 {
			abcClient.shouldFail = false
		}
		return time.Duration(0), attemptNum < 2
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	defer bc.Close()
	err := bdc(bc)
	assert.NoError(t, err)
	assert.Equal(t, 2, cp.connAttempts)
	assert.Equal(t, 1, setupInvoked)
}

func TestOrderingServiceSetupFailure(t *testing.T) {
	testOrderingServiceSetupFailure(t, blockDelivererConsumerWithRecv)
	testOrderingServiceSetupFailure(t, blockDelivererConsumerWithSend)
	connWG.Wait()
}

func testOrderingServiceSetupFailure(t *testing.T, bdc blocksDelivererConsumer) {
	
	
	
	
	cp := &connProducer{}
	clFactory := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return &abclient{}
	}
	setupInvoked := 0
	setup := func(blocksprovider.BlocksDeliverer) error {
		setupInvoked++
		if setupInvoked == 1 {
			return errors.New("Setup failed")
		}
		return nil
	}
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return time.Duration(0), true
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	defer bc.Close()
	err := bdc(bc)
	assert.NoError(t, err)
	assert.Equal(t, 2, cp.connAttempts)
	assert.Equal(t, 2, setupInvoked)
}

func TestOrderingServiceFirstOperationFailure(t *testing.T) {
	testOrderingServiceFirstOperationFailure(t, blockDelivererConsumerWithRecv)
	testOrderingServiceFirstOperationFailure(t, blockDelivererConsumerWithSend)
	connWG.Wait()
}

func testOrderingServiceFirstOperationFailure(t *testing.T, bdc blocksDelivererConsumer) {
	
	
	
	cp := &connProducer{}
	abStream := &abc{shouldFail: true}
	abcClient := &abclient{stream: abStream}
	clFactory := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return abcClient
	}

	setupInvoked := 0
	setup := func(blocksprovider.BlocksDeliverer) error {
		
		if setupInvoked == 1 {
			abStream.shouldFail = false
		}
		setupInvoked++
		return nil
	}
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return time.Duration(0), true
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	defer bc.Close()
	err := bdc(bc)
	assert.NoError(t, err)
	assert.Equal(t, 2, setupInvoked)
	assert.Equal(t, cp.connAttempts, 2)
}

func TestOrderingServiceCrashAndRecover(t *testing.T) {
	testOrderingServiceCrashAndRecover(t, blockDelivererConsumerWithRecv)
	testOrderingServiceCrashAndRecover(t, blockDelivererConsumerWithSend)
	connWG.Wait()
}

func testOrderingServiceCrashAndRecover(t *testing.T, bdc blocksDelivererConsumer) {
	
	
	
	cp := &connProducer{}
	abStream := &abc{}
	abcClient := &abclient{stream: abStream}
	clFactory := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return abcClient
	}

	setupInvoked := 0
	setup := func(blocksprovider.BlocksDeliverer) error {
		
		if setupInvoked == 1 {
			abStream.shouldFail = false
		}
		setupInvoked++
		return nil
	}
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return time.Duration(0), true
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	defer bc.Close()
	err := bdc(bc)
	assert.NoError(t, err)
	
	abStream.shouldFail = true
	err = bdc(bc)
	assert.NoError(t, err)
	assert.Equal(t, 2, cp.connAttempts)
	assert.Equal(t, 2, setupInvoked)
}

func TestOrderingServicePermanentCrash(t *testing.T) {
	testOrderingServicePermanentCrash(t, blockDelivererConsumerWithRecv)
	testOrderingServicePermanentCrash(t, blockDelivererConsumerWithSend)
	connWG.Wait()
}

func testOrderingServicePermanentCrash(t *testing.T, bdc blocksDelivererConsumer) {
	
	
	
	cp := &connProducer{}
	abStream := &abc{}
	abcClient := &abclient{stream: abStream}
	clFactory := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return abcClient
	}

	setupInvoked := 0
	setup := func(blocksprovider.BlocksDeliverer) error {
		setupInvoked++
		return nil
	}
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return time.Duration(0), attemptNum < 10
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	defer bc.Close()
	err := bdc(bc)
	assert.NoError(t, err)
	
	abStream.shouldFail = true
	cp.shouldFail = true
	err = bdc(bc)
	assert.Error(t, err)
	assert.Equal(t, 10, cp.connAttempts)
	assert.Equal(t, 1, setupInvoked)
}

func TestLimitedConnAttempts(t *testing.T) {
	testLimitedConnAttempts(t, blockDelivererConsumerWithRecv)
	testLimitedConnAttempts(t, blockDelivererConsumerWithSend)
	connWG.Wait()
}

func testLimitedConnAttempts(t *testing.T, bdc blocksDelivererConsumer) {
	
	
	cp := &connProducer{shouldFail: true}
	clFactory := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return &abclient{}
	}
	setupInvoked := 0
	setup := func(blocksprovider.BlocksDeliverer) error {
		setupInvoked++
		return nil
	}
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return time.Duration(0), attemptNum < 10
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	defer bc.Close()
	err := bdc(bc)
	assert.Error(t, err)
	assert.Equal(t, 10, cp.connAttempts)
	assert.Equal(t, 0, setupInvoked)
}

func TestLimitedTotalConnTimeRcv(t *testing.T) {
	testLimitedTotalConnTime(t, blockDelivererConsumerWithRecv)
	connWG.Wait()
}

func TestLimitedTotalConnTimeSnd(t *testing.T) {
	testLimitedTotalConnTime(t, blockDelivererConsumerWithSend)
	connWG.Wait()
}

func testLimitedTotalConnTime(t *testing.T, bdc blocksDelivererConsumer) {
	
	
	
	
	cp := &connProducer{shouldFail: true, connTime: 1500 * time.Millisecond}
	clFactory := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return &abclient{}
	}
	setupInvoked := 0
	setup := func(blocksprovider.BlocksDeliverer) error {
		setupInvoked++
		return nil
	}
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return time.Millisecond * 500, elapsedTime.Nanoseconds() < time.Second.Nanoseconds()
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	defer bc.Close()
	err := bdc(bc)
	assert.Error(t, err)
	assert.Equal(t, 3, cp.connAttempts)
	assert.Equal(t, 0, setupInvoked)
}

func TestGreenPath(t *testing.T) {
	testGreenPath(t, blockDelivererConsumerWithRecv)
	testGreenPath(t, blockDelivererConsumerWithSend)
	connWG.Wait()
}

func testGreenPath(t *testing.T, bdc blocksDelivererConsumer) {
	
	cp := &connProducer{}
	abcClient := &abclient{}
	clFactory := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return abcClient
	}

	setupInvoked := 0
	setup := func(blocksprovider.BlocksDeliverer) error {
		setupInvoked++
		return nil
	}
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return time.Duration(0), attemptNum < 1
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	defer bc.Close()
	err := bdc(bc)
	assert.NoError(t, err)
	assert.Equal(t, 1, cp.connAttempts)
	assert.Equal(t, 1, setupInvoked)
}

func TestCloseWhileRecv(t *testing.T) {
	
	
	
	fakeOrderer := mocks.NewOrderer(5611, t)
	time.Sleep(time.Second)
	defer fakeOrderer.Shutdown()
	cp := &connProducer{ordererEndpoint: "localhost:5611"}
	clFactory := func(conn *grpc.ClientConn) orderer.AtomicBroadcastClient {
		return orderer.NewAtomicBroadcastClient(conn)
	}

	setup := func(blocksprovider.BlocksDeliverer) error {
		return nil
	}
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return 0, true
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	var flag int32
	go func() {
		for fakeOrderer.ConnCount() == 0 {
			time.Sleep(time.Second)
		}
		atomic.StoreInt32(&flag, int32(1))
		bc.Close()
		bc.Close() 
	}()
	resp, err := bc.Recv()
	
	assert.Equal(t, int32(1), atomic.LoadInt32(&flag), "Recv returned before bc.Close() was called")
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, "client is closing", err.Error())
}

func TestCloseWhileSleep(t *testing.T) {
	testCloseWhileSleep(t, blockDelivererConsumerWithRecv)
	testCloseWhileSleep(t, blockDelivererConsumerWithSend)
	connWG.Wait()
}

func testCloseWhileSleep(t *testing.T, bdc blocksDelivererConsumer) {
	
	
	
	
	cp := &connProducer{}
	abcClient := &abclient{shouldFail: true}
	clFactory := func(*grpc.ClientConn) orderer.AtomicBroadcastClient {
		return abcClient
	}

	setupInvoked := 0
	setup := func(blocksprovider.BlocksDeliverer) error {
		setupInvoked++
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(1)
	backoffStrategy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		if attemptNum == 1 {
			go func() {
				time.Sleep(time.Second)
				wg.Done()
			}()
		}
		return time.Second * 1000000, true
	}
	bc := NewBroadcastClient(cp, clFactory, setup, backoffStrategy)
	go func() {
		wg.Wait()
		bc.Close()
		bc.Close() 
	}()
	err := bdc(bc)
	assert.Error(t, err)
	assert.Equal(t, 1, cp.connAttempts)
	assert.Equal(t, 0, setupInvoked)
}

type signerMock struct {
}

func (s *signerMock) NewSignatureHeader() (*common.SignatureHeader, error) {
	return &common.SignatureHeader{}, nil
}

func (s *signerMock) Sign(message []byte) ([]byte, error) {
	hasher := sha256.New()
	hasher.Write(message)
	return hasher.Sum(nil), nil
}

func TestProductionUsage(t *testing.T) {
	defer ensureNoGoroutineLeak(t)()
	
	
	os := mocks.NewOrderer(5612, t)
	os.SetNextExpectedSeek(5)

	connFact := func(endpoint string) (*grpc.ClientConn, error) {
		return grpc.Dial(endpoint, grpc.WithInsecure(), grpc.WithBlock())
	}
	prod := comm.NewConnectionProducer(connFact, []string{"localhost:5612"})
	clFact := func(cc *grpc.ClientConn) orderer.AtomicBroadcastClient {
		return orderer.NewAtomicBroadcastClient(cc)
	}
	onConnect := func(bd blocksprovider.BlocksDeliverer) error {
		env, err := utils.CreateSignedEnvelope(common.HeaderType_CONFIG_UPDATE,
			"TEST",
			&signerMock{}, newTestSeekInfo(), 0, 0)
		assert.NoError(t, err)
		return bd.Send(env)
	}
	retryPol := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return time.Second * 3, attemptNum < 2
	}
	cl := NewBroadcastClient(prod, clFact, onConnect, retryPol)
	go os.SendBlock(5)
	resp, err := cl.Recv()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(5), resp.GetBlock().Header.Number)
	os.Shutdown()
	cl.Close()
}

func TestDisconnect(t *testing.T) {
	
	
	
	
	

	defer ensureNoGoroutineLeak(t)()
	os1 := mocks.NewOrderer(5613, t)
	os1.SetNextExpectedSeek(5)
	os2 := mocks.NewOrderer(5614, t)
	os2.SetNextExpectedSeek(5)

	defer os1.Shutdown()
	defer os2.Shutdown()

	waitForConnectionToSomeOSN := func() {
		for {
			if os1.ConnCount() > 0 || os2.ConnCount() > 0 {
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}

	connFact := func(endpoint string) (*grpc.ClientConn, error) {
		return grpc.Dial(endpoint, grpc.WithInsecure(), grpc.WithBlock())
	}
	prod := comm.NewConnectionProducer(connFact, []string{"localhost:5613", "localhost:5614"})
	clFact := func(cc *grpc.ClientConn) orderer.AtomicBroadcastClient {
		return orderer.NewAtomicBroadcastClient(cc)
	}
	onConnect := func(bd blocksprovider.BlocksDeliverer) error {
		return nil
	}
	retryPol := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return time.Millisecond * 10, attemptNum < 100
	}

	cl := NewBroadcastClient(prod, clFact, onConnect, retryPol)
	stopChan := make(chan struct{})
	go func() {
		cl.Recv()
		stopChan <- struct{}{}
	}()
	waitForConnectionToSomeOSN()
	cl.Disconnect(false)

	i := 0
	os1Connected := false
	os2Connected := false

	for (!os1Connected || !os2Connected) && i < 100 {
		if os1.ConnCount() > 0 {
			os1Connected = true
		}

		if os2.ConnCount() > 0 {
			os2Connected = true
		}

		t.Log("Attempt", i, "os1ConnCount()=", os1.ConnCount(), "os2ConnCount()=", os2.ConnCount())
		i++
		if i == 100 {
			assert.Fail(t, "Didn't switch to other instance after many attempts")
		}
		cl.Disconnect(false)
		time.Sleep(time.Millisecond * 500)
	}
	cl.Close()
	select {
	case <-stopChan:
	case <-time.After(time.Second * 20):
		assert.Fail(t, "Didn't stop within a timely manner")
	}
}

func TestDisconnectAndDisableEndpoint(t *testing.T) {
	
	
	
	
	
	
	

	defer ensureNoGoroutineLeak(t)()
	os1 := mocks.NewOrderer(5613, t)
	os1.SetNextExpectedSeek(5)
	os2 := mocks.NewOrderer(5614, t)
	os2.SetNextExpectedSeek(5)

	defer os1.Shutdown()
	defer os2.Shutdown()

	orgEndpointDisableInterval := comm.EndpointDisableInterval
	comm.EndpointDisableInterval = time.Millisecond * 1500
	defer func() { comm.EndpointDisableInterval = orgEndpointDisableInterval }()

	connFact := func(endpoint string) (*grpc.ClientConn, error) {
		return grpc.Dial(endpoint, grpc.WithInsecure(), grpc.WithBlock())
	}
	prod := comm.NewConnectionProducer(connFact, []string{"localhost:5613", "localhost:5614"})
	clFact := func(cc *grpc.ClientConn) orderer.AtomicBroadcastClient {
		return orderer.NewAtomicBroadcastClient(cc)
	}
	onConnect := func(bd blocksprovider.BlocksDeliverer) error {
		return nil
	}

	retryPol := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		return time.Millisecond * 10, attemptNum < 10
	}

	cl := NewBroadcastClient(prod, clFact, onConnect, retryPol)
	defer cl.Close()

	
	go func() {
		cl.Recv()
	}()

	assert.True(t, waitForWithTimeout(time.Millisecond*100, func() bool {
		return os1.ConnCount() == 1 || os2.ConnCount() == 1
	}), "Didn't get connection to orderer")

	connectedToOS1 := os1.ConnCount() == 1

	
	cl.Disconnect(true)

	
	assert.True(t, waitForWithTimeout(time.Millisecond*100, func() bool {
		if connectedToOS1 {
			return os1.ConnCount() == 0 && os2.ConnCount() == 1
		}
		return os2.ConnCount() == 0 && os1.ConnCount() == 1
	}), "Didn't disconnect from orderer, or reconnected to a black-listed node")

	
	cl.Disconnect(true)

	
	assert.True(t, waitForWithTimeout(time.Millisecond*100, func() bool {
		return os1.ConnCount() == 1 || os2.ConnCount() == 1
	}), "Didn't got connection to orderer")
}

func waitForWithTimeout(timeout time.Duration, f func() bool) bool {

	ctx, cancelation := context.WithTimeout(context.Background(), timeout)
	defer cancelation()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-time.After(timeout / 10):
			if f() {
				return true
			}
		}
	}
}

func newTestSeekInfo() *orderer.SeekInfo {
	return &orderer.SeekInfo{Start: &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: 5}}},
		Stop:     &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: math.MaxUint64}}},
		Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
	}
}
