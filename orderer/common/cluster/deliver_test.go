/*
Copyright IBM Corp. 2018 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/common/cluster/mocks"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/roundrobin"
)


var gRPCBalancerLock = sync.Mutex{}

func init() {
	factory.InitFactories(nil)
}



type signerSerializer interface {
	identity.SignerSerializer
}

type wrappedBalancer struct {
	balancer.Balancer
	cd *countingDialer
}

func (wb *wrappedBalancer) Close() {
	defer atomic.AddUint32(&wb.cd.connectionCount, ^uint32(0))
	wb.Balancer.Close()
}

type countingDialer struct {
	name            string
	baseBuilder     balancer.Builder
	connectionCount uint32
}

func newCountingDialer() *countingDialer {
	gRPCBalancerLock.Lock()
	builder := balancer.Get(roundrobin.Name)
	gRPCBalancerLock.Unlock()

	buff := make([]byte, 16)
	rand.Read(buff)
	cb := &countingDialer{
		name:        hex.EncodeToString(buff),
		baseBuilder: builder,
	}

	gRPCBalancerLock.Lock()
	balancer.Register(cb)
	gRPCBalancerLock.Unlock()

	return cb
}

func (d *countingDialer) Build(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer {
	defer atomic.AddUint32(&d.connectionCount, 1)
	return &wrappedBalancer{Balancer: d.baseBuilder.Build(cc, opts), cd: d}
}

func (d *countingDialer) Name() string {
	return d.name
}

func (d *countingDialer) assertAllConnectionsClosed(t *testing.T) {
	timeLimit := time.Now().Add(timeout)
	for atomic.LoadUint32(&d.connectionCount) != uint32(0) && time.Now().Before(timeLimit) {
		time.Sleep(time.Millisecond)
	}
	assert.Equal(t, uint32(0), atomic.LoadUint32(&d.connectionCount))
}

func (d *countingDialer) Dial(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	gRPCBalancerLock.Lock()
	balancer := grpc.WithBalancerName(d.name)
	gRPCBalancerLock.Unlock()
	return grpc.DialContext(ctx, address, grpc.WithBlock(), grpc.WithInsecure(), balancer)
}

func noopBlockVerifierf(_ []*common.Block, _ string) error {
	return nil
}

func readSeekEnvelope(stream orderer.AtomicBroadcast_DeliverServer) (*orderer.SeekInfo, string, error) {
	env, err := stream.Recv()
	if err != nil {
		return nil, "", err
	}
	payload, err := protoutil.UnmarshalPayload(env.Payload)
	if err != nil {
		return nil, "", err
	}
	seekInfo := &orderer.SeekInfo{}
	if err = proto.Unmarshal(payload.Data, seekInfo); err != nil {
		return nil, "", err
	}
	chdr := &common.ChannelHeader{}
	if err = proto.Unmarshal(payload.Header.ChannelHeader, chdr); err != nil {
		return nil, "", err
	}
	return seekInfo, chdr.ChannelId, nil
}

type deliverServer struct {
	t *testing.T
	sync.Mutex
	err            error
	srv            *comm.GRPCServer
	seekAssertions chan func(*orderer.SeekInfo, string)
	blockResponses chan *orderer.DeliverResponse
}

func (ds *deliverServer) isFaulty() bool {
	ds.Lock()
	defer ds.Unlock()
	return ds.err != nil
}

func (*deliverServer) Broadcast(orderer.AtomicBroadcast_BroadcastServer) error {
	panic("implement me")
}

func (ds *deliverServer) Deliver(stream orderer.AtomicBroadcast_DeliverServer) error {
	ds.Lock()
	err := ds.err
	ds.Unlock()

	if err != nil {
		return err
	}
	seekInfo, channel, err := readSeekEnvelope(stream)
	if err != nil {
		panic(err)
	}
	
	seekAssert := <-ds.seekAssertions
	seekAssert(seekInfo, channel)

	if seekInfo.GetStart().GetSpecified() != nil {
		return ds.deliverBlocks(stream)
	}
	if seekInfo.GetStart().GetNewest() != nil {
		resp := <-ds.blocks()
		if resp == nil {
			return nil
		}
		return stream.Send(resp)
	}
	panic(fmt.Sprintf("expected either specified or newest seek but got %v", seekInfo.GetStart()))
}

func (ds *deliverServer) deliverBlocks(stream orderer.AtomicBroadcast_DeliverServer) error {
	for {
		blockChan := ds.blocks()
		response := <-blockChan
		
		
		
		
		if response == nil {
			return nil
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func (ds *deliverServer) blocks() chan *orderer.DeliverResponse {
	ds.Lock()
	defer ds.Unlock()
	blockChan := ds.blockResponses
	return blockChan
}

func (ds *deliverServer) setBlocks(blocks chan *orderer.DeliverResponse) {
	ds.Lock()
	defer ds.Unlock()
	ds.blockResponses = blocks
}

func (ds *deliverServer) port() int {
	_, portStr, err := net.SplitHostPort(ds.srv.Address())
	assert.NoError(ds.t, err)

	port, err := strconv.ParseInt(portStr, 10, 32)
	assert.NoError(ds.t, err)
	return int(port)
}

func (ds *deliverServer) resurrect() {
	var err error
	
	respChan := make(chan *orderer.DeliverResponse, 100)
	for resp := range ds.blocks() {
		respChan <- resp
	}
	ds.blockResponses = respChan
	ds.srv.Stop()
	
	ds.srv, err = comm.NewGRPCServer(fmt.Sprintf("127.0.0.1:%d", ds.port()), comm.ServerConfig{})
	assert.NoError(ds.t, err)
	orderer.RegisterAtomicBroadcastServer(ds.srv.Server(), ds)
	go ds.srv.Start()
}

func (ds *deliverServer) stop() {
	ds.srv.Stop()
	close(ds.blocks())
}

func (ds *deliverServer) enqueueResponse(seq uint64) {
	ds.blocks() <- &orderer.DeliverResponse{
		Type: &orderer.DeliverResponse_Block{Block: protoutil.NewBlock(seq, nil)},
	}
}

func (ds *deliverServer) addExpectProbeAssert() {
	ds.seekAssertions <- func(info *orderer.SeekInfo, _ string) {
		assert.NotNil(ds.t, info.GetStart().GetNewest())
	}
}

func (ds *deliverServer) addExpectPullAssert(seq uint64) {
	ds.seekAssertions <- func(info *orderer.SeekInfo, _ string) {
		assert.NotNil(ds.t, info.GetStart().GetSpecified())
		assert.Equal(ds.t, seq, info.GetStart().GetSpecified().Number)
	}
}

func newClusterNode(t *testing.T) *deliverServer {
	srv, err := comm.NewGRPCServer("127.0.0.1:0", comm.ServerConfig{})
	if err != nil {
		panic(err)
	}
	ds := &deliverServer{
		t:              t,
		seekAssertions: make(chan func(*orderer.SeekInfo, string), 100),
		blockResponses: make(chan *orderer.DeliverResponse, 100),
		srv:            srv,
	}
	orderer.RegisterAtomicBroadcastServer(srv.Server(), ds)
	go srv.Start()
	return ds
}

func newBlockPuller(dialer *countingDialer, orderers ...string) *cluster.BlockPuller {
	return &cluster.BlockPuller{
		Dialer:              dialer,
		Channel:             "mychannel",
		Signer:              &mocks.SignerSerializer{},
		Endpoints:           orderers,
		FetchTimeout:        time.Second,
		MaxTotalBufferBytes: 1024 * 1024, 
		RetryTimeout:        time.Millisecond * 10,
		VerifyBlockSequence: noopBlockVerifierf,
		Logger:              flogging.MustGetLogger("test"),
	}
}

func TestBlockPullerBasicHappyPath(t *testing.T) {
	
	
	osn := newClusterNode(t)
	defer osn.stop()

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn.srv.Address())

	
	osn.addExpectProbeAssert()
	
	osn.enqueueResponse(10)
	
	osn.addExpectPullAssert(5)
	
	for i := 5; i <= 10; i++ {
		osn.enqueueResponse(uint64(i))
	}

	for i := 5; i <= 10; i++ {
		assert.Equal(t, uint64(i), bp.PullBlock(uint64(i)).Header.Number)
	}
	assert.Len(t, osn.blockResponses, 0)

	bp.Close()
	dialer.assertAllConnectionsClosed(t)
}

func TestBlockPullerDuplicate(t *testing.T) {
	
	
	
	osn := newClusterNode(t)
	defer osn.stop()

	dialer := newCountingDialer()
	
	bp := newBlockPuller(dialer, osn.srv.Address(), osn.srv.Address())

	
	osn.addExpectProbeAssert()
	osn.addExpectProbeAssert()
	
	osn.enqueueResponse(3)
	osn.enqueueResponse(3)
	
	osn.addExpectPullAssert(1)
	
	for i := 1; i <= 3; i++ {
		osn.enqueueResponse(uint64(i))
	}

	for i := 1; i <= 3; i++ {
		assert.Equal(t, uint64(i), bp.PullBlock(uint64(i)).Header.Number)
	}
	assert.Len(t, osn.blockResponses, 0)

	bp.Close()
	dialer.assertAllConnectionsClosed(t)
}

func TestBlockPullerHeavyBlocks(t *testing.T) {
	
	
	
	
	

	osn := newClusterNode(t)
	defer osn.stop()
	osn.addExpectProbeAssert()
	osn.addExpectPullAssert(1)
	osn.enqueueResponse(100) 

	enqueueBlockBatch := func(start, end uint64) {
		for seq := start; seq <= end; seq++ {
			resp := &orderer.DeliverResponse{
				Type: &orderer.DeliverResponse_Block{
					Block: protoutil.NewBlock(seq, nil),
				},
			}
			data := resp.GetBlock().Data.Data
			resp.GetBlock().Data.Data = append(data, make([]byte, 1024))
			osn.blockResponses <- resp
		}
	}

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn.srv.Address())
	var gotBlockMessageCount int
	bp.Logger = bp.Logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		if strings.Contains(entry.Message, "Got block") {
			gotBlockMessageCount++
		}
		return nil
	}))
	bp.MaxTotalBufferBytes = 1024 * 10 

	
	
	
	for i := uint64(0); i < 5; i++ {
		enqueueBlockBatch(i*10+uint64(1), i*10+uint64(10))
		for seq := i*10 + uint64(1); seq <= i*10+uint64(10); seq++ {
			assert.Equal(t, seq, bp.PullBlock(seq).Header.Number)
		}
	}

	assert.Equal(t, 50, gotBlockMessageCount)
	bp.Close()
	dialer.assertAllConnectionsClosed(t)
}

func TestBlockPullerClone(t *testing.T) {
	
	
	
	
	osn1 := newClusterNode(t)
	defer osn1.stop()

	osn1.addExpectProbeAssert()
	osn1.addExpectPullAssert(1)
	
	osn1.enqueueResponse(100)
	osn1.enqueueResponse(1)
	
	
	
	osn1.blockResponses <- nil

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn1.srv.Address())
	bp.FetchTimeout = time.Millisecond * 100
	
	bp.MaxTotalBufferBytes = 1
	
	bpClone := bp.Clone()
	
	bpClone.Channel = "foo"
	
	assert.Equal(t, "mychannel", bp.Channel)

	block := bp.PullBlock(1)
	assert.Equal(t, uint64(1), block.Header.Number)

	
	
	bp.Close()
	dialer.assertAllConnectionsClosed(t)

	
	
	
	osn1.addExpectProbeAssert()
	osn1.addExpectPullAssert(2)
	osn1.enqueueResponse(100)
	osn1.enqueueResponse(2)

	block = bpClone.PullBlock(2)
	assert.Equal(t, uint64(2), block.Header.Number)

	bpClone.Close()
	dialer.assertAllConnectionsClosed(t)
}

func TestBlockPullerHeightsByEndpoints(t *testing.T) {
	
	
	
	
	osn1 := newClusterNode(t)

	osn2 := newClusterNode(t)
	defer osn2.stop()

	osn3 := newClusterNode(t)
	defer osn3.stop()

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn1.srv.Address(), osn2.srv.Address(), osn3.srv.Address())

	
	osn1.addExpectProbeAssert()
	osn2.addExpectProbeAssert()
	osn3.addExpectProbeAssert()

	
	osn1.stop()
	
	osn2.blockResponses <- &orderer.DeliverResponse{
		Type: &orderer.DeliverResponse_Status{Status: common.Status_FORBIDDEN},
	}
	
	osn3.enqueueResponse(5)

	res, err := bp.HeightsByEndpoints()
	assert.NoError(t, err)
	expected := map[string]uint64{
		osn3.srv.Address(): 6,
	}
	assert.Equal(t, expected, res)

	bp.Close()
	dialer.assertAllConnectionsClosed(t)
}

func TestBlockPullerMultipleOrderers(t *testing.T) {
	
	
	
	

	osn1 := newClusterNode(t)
	defer osn1.stop()

	osn2 := newClusterNode(t)
	defer osn2.stop()

	osn3 := newClusterNode(t)
	defer osn3.stop()

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn1.srv.Address(), osn2.srv.Address(), osn3.srv.Address())

	
	osn1.addExpectProbeAssert()
	osn2.addExpectProbeAssert()
	osn3.addExpectProbeAssert()
	
	osn1.enqueueResponse(5)
	osn2.enqueueResponse(5)
	osn3.enqueueResponse(5)

	
	osn1.addExpectPullAssert(3)
	osn2.addExpectPullAssert(3)
	osn3.addExpectPullAssert(3)

	
	for i := 3; i <= 5; i++ {
		osn1.enqueueResponse(uint64(i))
		osn2.enqueueResponse(uint64(i))
		osn3.enqueueResponse(uint64(i))
	}

	initialTotalBlockAmount := len(osn1.blockResponses) + len(osn2.blockResponses) + len(osn3.blockResponses)

	for i := 3; i <= 5; i++ {
		assert.Equal(t, uint64(i), bp.PullBlock(uint64(i)).Header.Number)
	}

	
	
	
	finalTotalBlockAmount := len(osn1.blockResponses) + len(osn2.blockResponses) + len(osn3.blockResponses)
	assert.Equal(t, initialTotalBlockAmount-6, finalTotalBlockAmount)

	bp.Close()
	dialer.assertAllConnectionsClosed(t)
}

func TestBlockPullerFailover(t *testing.T) {
	
	
	
	
	
	
	

	osn1 := newClusterNode(t)
	osn1.addExpectProbeAssert()
	osn1.addExpectPullAssert(1)
	osn1.enqueueResponse(3)
	osn1.enqueueResponse(1)

	osn2 := newClusterNode(t)
	defer osn2.stop()

	osn2.addExpectProbeAssert()
	osn2.addExpectPullAssert(1)
	
	osn2.enqueueResponse(3)
	
	osn2.enqueueResponse(1)
	osn2.enqueueResponse(2)
	osn2.enqueueResponse(3)

	
	osn2.stop()

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn1.srv.Address(), osn2.srv.Address())

	
	
	bp.FetchTimeout = time.Hour

	
	
	var pulledBlock1 sync.WaitGroup
	pulledBlock1.Add(1)
	var once sync.Once
	bp.Logger = bp.Logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		if strings.Contains(entry.Message, "Got block [1] of size") {
			once.Do(pulledBlock1.Done)
		}
		return nil
	}))

	go func() {
		
		pulledBlock1.Wait()
		
		osn1.stop()
		osn2.resurrect()
	}()

	
	assert.Equal(t, uint64(1), bp.PullBlock(uint64(1)).Header.Number)
	assert.Equal(t, uint64(2), bp.PullBlock(uint64(2)).Header.Number)
	assert.Equal(t, uint64(3), bp.PullBlock(uint64(3)).Header.Number)

	bp.Close()
	dialer.assertAllConnectionsClosed(t)
}

func TestBlockPullerNoneResponsiveOrderer(t *testing.T) {
	
	
	
	
	

	osn1 := newClusterNode(t)
	defer osn1.stop()

	osn2 := newClusterNode(t)
	defer osn2.stop()

	osn1.addExpectProbeAssert()
	osn2.addExpectProbeAssert()
	osn1.enqueueResponse(3)
	osn2.enqueueResponse(3)

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn1.srv.Address(), osn2.srv.Address())
	bp.FetchTimeout = time.Millisecond * 100

	notInUseOrdererNode := osn2
	
	
	var waitForConnection sync.WaitGroup
	waitForConnection.Add(1)
	var once sync.Once
	bp.Logger = bp.Logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		if !strings.Contains(entry.Message, "Sending request for block [1]") {
			return nil
		}
		defer once.Do(waitForConnection.Done)
		s := entry.Message[len("Sending request for block [1] to 127.0.0.1:"):]
		port, err := strconv.ParseInt(s, 10, 32)
		assert.NoError(t, err)
		
		
		if osn2.port() == int(port) {
			notInUseOrdererNode = osn1
			
			
			osn2.enqueueResponse(1)
			osn2.addExpectPullAssert(1)
		} else {
			
			osn1.enqueueResponse(1)
			osn1.addExpectPullAssert(1)
		}
		return nil
	}))

	go func() {
		waitForConnection.Wait()
		
		notInUseOrdererNode.enqueueResponse(3)
		notInUseOrdererNode.addExpectProbeAssert()
		
		notInUseOrdererNode.addExpectPullAssert(1)
		notInUseOrdererNode.enqueueResponse(1)
		notInUseOrdererNode.enqueueResponse(2)
		notInUseOrdererNode.enqueueResponse(3)
	}()

	
	assert.Equal(t, uint64(1), bp.PullBlock(uint64(1)).Header.Number)
	assert.Equal(t, uint64(2), bp.PullBlock(uint64(2)).Header.Number)
	assert.Equal(t, uint64(3), bp.PullBlock(uint64(3)).Header.Number)

	bp.Close()
	dialer.assertAllConnectionsClosed(t)
}

func TestBlockPullerNoOrdererAliveAtStartup(t *testing.T) {
	
	
	osn := newClusterNode(t)
	osn.stop()
	defer osn.stop()

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn.srv.Address())

	
	var connectionAttempt sync.WaitGroup
	connectionAttempt.Add(1)
	bp.Logger = bp.Logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		if strings.Contains(entry.Message, "Failed connecting to") {
			connectionAttempt.Done()
		}
		return nil
	}))

	go func() {
		connectionAttempt.Wait()
		osn.resurrect()
		
		osn.addExpectProbeAssert()
		
		osn.enqueueResponse(2)
		
		osn.addExpectPullAssert(1)
		
		osn.enqueueResponse(1)
		osn.enqueueResponse(2)
	}()

	assert.Equal(t, uint64(1), bp.PullBlock(1).Header.Number)
	assert.Equal(t, uint64(2), bp.PullBlock(2).Header.Number)

	bp.Close()
	dialer.assertAllConnectionsClosed(t)
}

func TestBlockPullerFailures(t *testing.T) {
	
	
	failureError := errors.New("oops, something went wrong")
	failStream := func(osn *deliverServer, _ *cluster.BlockPuller) {
		osn.Lock()
		osn.err = failureError
		osn.Unlock()
	}

	badSigErr := errors.New("bad signature")
	malformBlockSignatureAndRecreateOSNBuffer := func(osn *deliverServer, bp *cluster.BlockPuller) {
		bp.VerifyBlockSequence = func(_ []*common.Block, _ string) error {
			close(osn.blocks())
			
			defer func() {
				
				if badSigErr == nil {
					return
				}
				badSigErr = nil
				osn.setBlocks(make(chan *orderer.DeliverResponse, 100))
				osn.enqueueResponse(3)
				osn.enqueueResponse(1)
				osn.enqueueResponse(2)
				osn.enqueueResponse(3)
			}()
			return badSigErr
		}
	}

	noFailFunc := func(_ *deliverServer, _ *cluster.BlockPuller) {}

	recover := func(osn *deliverServer, bp *cluster.BlockPuller) func(entry zapcore.Entry) error {
		return func(entry zapcore.Entry) error {
			if osn.isFaulty() && strings.Contains(entry.Message, failureError.Error()) {
				osn.Lock()
				osn.err = nil
				osn.Unlock()
			}
			if strings.Contains(entry.Message, "Failed verifying") {
				bp.VerifyBlockSequence = noopBlockVerifierf
			}
			return nil
		}
	}

	failAfterConnection := func(osn *deliverServer, logTrigger string, failFunc func()) func(entry zapcore.Entry) error {
		once := &sync.Once{}
		return func(entry zapcore.Entry) error {
			if !osn.isFaulty() && strings.Contains(entry.Message, logTrigger) {
				once.Do(func() {
					failFunc()
				})
			}
			return nil
		}
	}

	for _, testCase := range []struct {
		name       string
		logTrigger string
		beforeFunc func(*deliverServer, *cluster.BlockPuller)
		failFunc   func(*deliverServer, *cluster.BlockPuller)
	}{
		{
			name:       "failure at probe",
			logTrigger: "skip this for this test case",
			beforeFunc: func(osn *deliverServer, bp *cluster.BlockPuller) {
				failStream(osn, nil)
				
				osn.addExpectProbeAssert()
				
				osn.enqueueResponse(3)
				
				osn.addExpectPullAssert(1)
			},
			failFunc: noFailFunc,
		},
		{
			name:       "failure at pull",
			logTrigger: "Sending request for block [1]",
			beforeFunc: func(osn *deliverServer, bp *cluster.BlockPuller) {
				
				osn.addExpectProbeAssert()
				
				osn.addExpectProbeAssert()
				
				osn.enqueueResponse(3)
				osn.enqueueResponse(3)
				osn.addExpectPullAssert(1)
			},
			failFunc: failStream,
		},
		{
			name:       "failure at verifying pulled block",
			logTrigger: "Sending request for block [1]",
			beforeFunc: func(osn *deliverServer, bp *cluster.BlockPuller) {
				
				osn.addExpectProbeAssert()
				osn.enqueueResponse(3)
				
				osn.addExpectPullAssert(1)
				osn.enqueueResponse(1)
				osn.enqueueResponse(2)
				osn.enqueueResponse(3)
				
				osn.addExpectProbeAssert()
				osn.addExpectPullAssert(1)
			},
			failFunc: malformBlockSignatureAndRecreateOSNBuffer,
		},
	} {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			osn := newClusterNode(t)
			defer osn.stop()

			dialer := newCountingDialer()
			bp := newBlockPuller(dialer, osn.srv.Address())

			testCase.beforeFunc(osn, bp)

			
			fail := func() {
				testCase.failFunc(osn, bp)
			}
			bp.Logger = bp.Logger.WithOptions(zap.Hooks(recover(osn, bp), failAfterConnection(osn, testCase.logTrigger, fail)))

			
			osn.enqueueResponse(1)
			osn.enqueueResponse(2)
			osn.enqueueResponse(3)

			assert.Equal(t, uint64(1), bp.PullBlock(uint64(1)).Header.Number)
			assert.Equal(t, uint64(2), bp.PullBlock(uint64(2)).Header.Number)
			assert.Equal(t, uint64(3), bp.PullBlock(uint64(3)).Header.Number)

			bp.Close()
			dialer.assertAllConnectionsClosed(t)
		})
	}
}

func TestBlockPullerBadBlocks(t *testing.T) {
	

	removeHeader := func(resp *orderer.DeliverResponse) *orderer.DeliverResponse {
		resp.GetBlock().Header = nil
		return resp
	}

	removeData := func(resp *orderer.DeliverResponse) *orderer.DeliverResponse {
		resp.GetBlock().Data = nil
		return resp
	}

	removeMetadata := func(resp *orderer.DeliverResponse) *orderer.DeliverResponse {
		resp.GetBlock().Metadata = nil
		return resp
	}

	changeType := func(resp *orderer.DeliverResponse) *orderer.DeliverResponse {
		resp.Type = nil
		return resp
	}

	statusType := func(resp *orderer.DeliverResponse) *orderer.DeliverResponse {
		resp.Type = &orderer.DeliverResponse_Status{
			Status: common.Status_INTERNAL_SERVER_ERROR,
		}
		return resp
	}

	changeSequence := func(resp *orderer.DeliverResponse) *orderer.DeliverResponse {
		resp.GetBlock().Header.Number = 3
		return resp
	}

	testcases := []struct {
		name           string
		corruptBlock   func(block *orderer.DeliverResponse) *orderer.DeliverResponse
		expectedErrMsg string
	}{
		{
			name:           "nil header",
			corruptBlock:   removeHeader,
			expectedErrMsg: "block header is nil",
		},
		{
			name:           "nil data",
			corruptBlock:   removeData,
			expectedErrMsg: "block data is nil",
		},
		{
			name:           "nil metadata",
			corruptBlock:   removeMetadata,
			expectedErrMsg: "block metadata is empty",
		},
		{
			name:           "wrong type",
			corruptBlock:   changeType,
			expectedErrMsg: "response is of type <nil>, but expected a block",
		},
		{
			name:           "bad type",
			corruptBlock:   statusType,
			expectedErrMsg: "faulty node, received: status:INTERNAL_SERVER_ERROR ",
		},
		{
			name:           "wrong number",
			corruptBlock:   changeSequence,
			expectedErrMsg: "got unexpected sequence",
		},
	}

	for _, testCase := range testcases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {

			osn := newClusterNode(t)
			defer osn.stop()

			osn.addExpectProbeAssert()
			osn.addExpectPullAssert(10)

			dialer := newCountingDialer()
			bp := newBlockPuller(dialer, osn.srv.Address())

			osn.enqueueResponse(10)
			osn.enqueueResponse(10)
			
			block := <-osn.blockResponses
			
			
			osn.blockResponses <- testCase.corruptBlock(block)
			
			
			osn.addExpectProbeAssert()
			osn.addExpectPullAssert(10)

			
			var detectedBadBlock sync.WaitGroup
			detectedBadBlock.Add(1)
			bp.Logger = bp.Logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
				if strings.Contains(entry.Message, fmt.Sprintf("Failed pulling blocks: %s", testCase.expectedErrMsg)) {
					detectedBadBlock.Done()
					
					close(osn.blocks())
					
					osn.setBlocks(make(chan *orderer.DeliverResponse, 100))
					
					osn.enqueueResponse(10)
					osn.enqueueResponse(10)
				}
				return nil
			}))

			bp.PullBlock(10)
			detectedBadBlock.Wait()

			bp.Close()
			dialer.assertAllConnectionsClosed(t)
		})
	}
}

func TestImpatientStreamFailure(t *testing.T) {
	osn := newClusterNode(t)
	dialer := newCountingDialer()
	defer dialer.assertAllConnectionsClosed(t)
	
	
	var conn *grpc.ClientConn
	var err error

	gt := gomega.NewGomegaWithT(t)
	gt.Eventually(func() (bool, error) {
		conn, err = dialer.Dial(osn.srv.Address())
		return true, err
	}).Should(gomega.BeTrue())
	newStream := cluster.NewImpatientStream(conn, time.Millisecond*100)
	defer conn.Close()
	
	osn.stop()
	
	gt.Eventually(func() (bool, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()
		conn, _ := grpc.DialContext(ctx, osn.srv.Address(), grpc.WithBlock(), grpc.WithInsecure())
		if conn != nil {
			conn.Close()
			return false, nil
		}
		return true, nil
	}).Should(gomega.BeTrue())
	stream, err := newStream()
	if err != nil {
		return
	}
	_, err = stream.Recv()
	assert.Error(t, err)
}

func TestBlockPullerMaxRetriesExhausted(t *testing.T) {
	
	
	
	
	
	
	
	

	osn := newClusterNode(t)
	defer osn.stop()

	
	osn.enqueueResponse(3)
	osn.addExpectProbeAssert()
	
	osn.addExpectPullAssert(1)
	osn.enqueueResponse(1)
	
	osn.enqueueResponse(2)
	osn.enqueueResponse(2)
	
	
	
	osn.blockResponses <- nil

	for i := 0; i < 2; i++ {
		
		osn.addExpectProbeAssert()
		
		osn.enqueueResponse(3)
		
		
		osn.addExpectPullAssert(3)
		
		osn.enqueueResponse(2)
		
		osn.blockResponses <- nil
	}

	dialer := newCountingDialer()
	bp := newBlockPuller(dialer, osn.srv.Address())

	var exhaustedRetryAttemptsLogged bool

	bp.Logger = bp.Logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Message == "Failed pulling block [3]: retry count exhausted(2)" {
			exhaustedRetryAttemptsLogged = true
		}
		return nil
	}))

	bp.MaxPullBlockRetries = 2
	
	
	bp.FetchTimeout = time.Hour
	
	
	bp.MaxTotalBufferBytes = 1

	
	assert.Equal(t, uint64(1), bp.PullBlock(uint64(1)).Header.Number)
	assert.Equal(t, uint64(2), bp.PullBlock(uint64(2)).Header.Number)
	assert.Nil(t, bp.PullBlock(uint64(3)))

	bp.Close()
	dialer.assertAllConnectionsClosed(t)
	assert.True(t, exhaustedRetryAttemptsLogged)
}
