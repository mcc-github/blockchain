/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster_test

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	comm_utils "github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/common/cluster/mocks"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/onsi/gomega"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

const (
	testChannel  = "test"
	testChannel2 = "test2"
	timeout      = time.Second * 10
)

var (
	
	
	
	ca = createCAOrPanic()

	lastNodeID uint64

	testSubReq = &orderer.SubmitRequest{
		Channel: "test",
	}

	testStepReq = &orderer.StepRequest{
		Channel: "test",
		Payload: []byte("test"),
	}

	testStepRes = &orderer.StepResponse{
		Payload: []byte("test"),
	}

	fooReq = &orderer.StepRequest{
		Channel: "foo",
	}

	fooRes = &orderer.StepResponse{
		Payload: []byte("foo"),
	}

	barReq = &orderer.StepRequest{
		Channel: "bar",
	}

	barRes = &orderer.StepResponse{
		Payload: []byte("bar"),
	}

	channelExtractor = &mockChannelExtractor{}
)

func nextUnusedID() uint64 {
	return atomic.AddUint64(&lastNodeID, 1)
}

func createCAOrPanic() tlsgen.CA {
	ca, err := tlsgen.NewCA()
	if err != nil {
		panic(fmt.Sprintf("failed creating CA: %+v", err))
	}
	return ca
}

type mockChannelExtractor struct{}

func (*mockChannelExtractor) TargetChannel(msg proto.Message) string {
	if stepReq, isStepReq := msg.(*orderer.StepRequest); isStepReq {
		return stepReq.Channel
	}
	if subReq, isSubReq := msg.(*orderer.SubmitRequest); isSubReq {
		return string(subReq.Channel)
	}
	return ""
}

type clusterNode struct {
	dialer       *cluster.PredicateDialer
	handler      *mocks.Handler
	nodeInfo     cluster.RemoteNode
	srv          *comm_utils.GRPCServer
	bindAddress  string
	clientConfig comm_utils.ClientConfig
	serverConfig comm_utils.ServerConfig
	c            *cluster.Comm
}

func (cn *clusterNode) Submit(stream orderer.Cluster_SubmitServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	res, err := cn.c.DispatchSubmit(stream.Context(), req)
	if err != nil {
		return err
	}
	return stream.Send(res)
}

func (cn *clusterNode) Step(ctx context.Context, msg *orderer.StepRequest) (*orderer.StepResponse, error) {
	return cn.c.DispatchStep(ctx, msg)
}

func (cn *clusterNode) resurrect() {
	gRPCServer, err := comm_utils.NewGRPCServer(cn.bindAddress, cn.serverConfig)
	if err != nil {
		panic(fmt.Errorf("failed starting gRPC server: %v", err))
	}
	cn.srv = gRPCServer
	orderer.RegisterClusterServer(gRPCServer.Server(), cn)
	go cn.srv.Start()
}

func (cn *clusterNode) stop() {
	cn.srv.Stop()
	cn.c.Shutdown()
}

func (cn *clusterNode) renewCertificates() {
	clientKeyPair, err := ca.NewClientCertKeyPair()
	if err != nil {
		panic(fmt.Errorf("failed creating client certificate %v", err))
	}
	serverKeyPair, err := ca.NewServerCertKeyPair("127.0.0.1")
	if err != nil {
		panic(fmt.Errorf("failed creating server certificate %v", err))
	}

	cn.nodeInfo.ClientTLSCert = clientKeyPair.TLSCert.Raw
	cn.nodeInfo.ServerTLSCert = serverKeyPair.TLSCert.Raw

	cn.serverConfig.SecOpts.Certificate = serverKeyPair.Cert
	cn.serverConfig.SecOpts.Key = serverKeyPair.Key

	cn.clientConfig.SecOpts.Key = clientKeyPair.Key
	cn.clientConfig.SecOpts.Certificate = clientKeyPair.Cert
	cn.dialer.SetConfig(cn.clientConfig)
}

func newTestNode(t *testing.T) *clusterNode {
	serverKeyPair, err := ca.NewServerCertKeyPair("127.0.0.1")
	assert.NoError(t, err)

	clientKeyPair, _ := ca.NewClientCertKeyPair()

	handler := &mocks.Handler{}
	clientConfig := comm_utils.ClientConfig{
		Timeout: time.Millisecond * 100,
		SecOpts: &comm_utils.SecureOptions{
			RequireClientCert: true,
			Key:               clientKeyPair.Key,
			Certificate:       clientKeyPair.Cert,
			ServerRootCAs:     [][]byte{ca.CertBytes()},
			UseTLS:            true,
			ClientRootCAs:     [][]byte{ca.CertBytes()},
		},
	}

	dialer := cluster.NewTLSPinningDialer(clientConfig)

	srvConfig := comm_utils.ServerConfig{
		SecOpts: &comm_utils.SecureOptions{
			Key:         serverKeyPair.Key,
			Certificate: serverKeyPair.Cert,
			UseTLS:      true,
		},
	}
	gRPCServer, err := comm_utils.NewGRPCServer("127.0.0.1:", srvConfig)
	assert.NoError(t, err)

	tstSrv := &clusterNode{
		dialer:       dialer,
		clientConfig: clientConfig,
		serverConfig: srvConfig,
		bindAddress:  gRPCServer.Address(),
		handler:      handler,
		nodeInfo: cluster.RemoteNode{
			Endpoint:      gRPCServer.Address(),
			ID:            nextUnusedID(),
			ServerTLSCert: serverKeyPair.TLSCert.Raw,
			ClientTLSCert: clientKeyPair.TLSCert.Raw,
		},
		srv: gRPCServer,
	}

	tstSrv.c = &cluster.Comm{
		Logger:       logging.MustGetLogger("test"),
		Chan2Members: make(cluster.MembersByChannel),
		H:            handler,
		ChanExt:      channelExtractor,
		Connections:  cluster.NewConnectionStore(dialer),
	}

	orderer.RegisterClusterServer(gRPCServer.Server(), tstSrv)
	go gRPCServer.Start()
	return tstSrv
}

func TestBasic(t *testing.T) {
	t.Parallel()
	
	

	node1 := newTestNode(t)
	node2 := newTestNode(t)

	defer node1.stop()
	defer node2.stop()

	node1.handler.On("OnStep", testChannel, node2.nodeInfo.ID, mock.Anything).Return(testStepRes, nil)
	node2.handler.On("OnStep", testChannel, node1.nodeInfo.ID, mock.Anything).Return(testStepRes, nil)

	config := []cluster.RemoteNode{node1.nodeInfo, node2.nodeInfo}
	node1.c.Configure(testChannel, config)
	node2.c.Configure(testChannel, config)

	assertBiDiCommunication(t, node1, node2, testStepReq)
}

func TestUnavailableHosts(t *testing.T) {
	t.Parallel()
	
	
	node1 := newTestNode(t)
	defer node1.stop()

	node2 := newTestNode(t)
	node2.stop()

	node1.c.Configure(testChannel, []cluster.RemoteNode{node2.nodeInfo})
	_, err := node1.c.Remote(testChannel, node2.nodeInfo.ID)
	assert.EqualError(t, err, "failed to create new connection: context deadline exceeded")
}

func TestStreamAbort(t *testing.T) {
	t.Parallel()

	
	
	
	
	
	

	node2 := newTestNode(t)
	defer node2.stop()

	invalidNodeInfo := cluster.RemoteNode{
		ID:            node2.nodeInfo.ID,
		ServerTLSCert: []byte{1, 2, 3},
		ClientTLSCert: []byte{1, 2, 3},
	}

	for _, tst := range []struct {
		testName      string
		membership    []cluster.RemoteNode
		expectedError string
	}{
		{
			testName:      "Evicted from membership",
			membership:    nil,
			expectedError: "rpc error: code = Canceled desc = context canceled",
		},
		{
			testName:      "Changed TLS certificate",
			membership:    []cluster.RemoteNode{invalidNodeInfo},
			expectedError: "rpc error: code = Canceled desc = context canceled",
		},
	} {
		t.Run(tst.testName, func(t *testing.T) {
			testStreamAbort(t, node2, tst.membership, tst.expectedError)
		})
	}
	node2.handler.AssertNumberOfCalls(t, "OnSubmit", 2)
}

func testStreamAbort(t *testing.T, node2 *clusterNode, newMembership []cluster.RemoteNode, expectedError string) {
	node1 := newTestNode(t)
	defer node1.stop()

	node1.c.Configure(testChannel, []cluster.RemoteNode{node2.nodeInfo})
	node2.c.Configure(testChannel, []cluster.RemoteNode{node1.nodeInfo})
	node1.c.Configure(testChannel2, []cluster.RemoteNode{node2.nodeInfo})
	node2.c.Configure(testChannel2, []cluster.RemoteNode{node1.nodeInfo})

	var waitForReconfigWG sync.WaitGroup
	waitForReconfigWG.Add(1)

	var streamCreated sync.WaitGroup
	streamCreated.Add(1)

	node2.handler.On("OnSubmit", testChannel, node1.nodeInfo.ID, mock.Anything).Once().Run(func(_ mock.Arguments) {
		
		streamCreated.Done()
		
		
		waitForReconfigWG.Wait()
	}).Return(nil, nil).Once()

	rm1, err := node1.c.Remote(testChannel, node2.nodeInfo.ID)
	assert.NoError(t, err)

	errorChan := make(chan error)

	go func() {
		stream, err := rm1.SubmitStream()
		assert.NoError(t, err)
		
		err = stream.Send(&orderer.SubmitRequest{
			Channel: testChannel,
		})
		assert.NoError(t, err)
		_, err = stream.Recv()
		assert.EqualError(t, err, expectedError)
		errorChan <- err
	}()

	go func() {
		
		streamCreated.Wait()
		
		node1.c.Configure(testChannel, newMembership)
		waitForReconfigWG.Done()
	}()

	<-errorChan
}

func TestDoubleReconfigure(t *testing.T) {
	t.Parallel()
	
	
	
	

	node1 := newTestNode(t)
	node2 := newTestNode(t)

	defer node1.stop()
	defer node2.stop()

	node1.c.Configure(testChannel, []cluster.RemoteNode{node2.nodeInfo})
	rm1, err := node1.c.Remote(testChannel, node2.nodeInfo.ID)
	assert.NoError(t, err)

	node1.c.Configure(testChannel, []cluster.RemoteNode{node2.nodeInfo})
	rm2, err := node1.c.Remote(testChannel, node2.nodeInfo.ID)
	assert.NoError(t, err)
	
	assert.True(t, rm1 == rm2)
}

func TestInvalidChannel(t *testing.T) {
	t.Parallel()
	
	
	

	t.Run("channel doesn't exist", func(t *testing.T) {
		t.Parallel()
		node1 := newTestNode(t)
		defer node1.stop()

		_, err := node1.c.Remote(testChannel, 0)
		assert.EqualError(t, err, "channel test doesn't exist")
	})

	t.Run("channel cannot be extracted", func(t *testing.T) {
		t.Parallel()
		node1 := newTestNode(t)
		defer node1.stop()

		node1.c.Configure(testChannel, []cluster.RemoteNode{node1.nodeInfo})
		stub, err := node1.c.Remote(testChannel, node1.nodeInfo.ID)
		assert.NoError(t, err)
		
		_, err = stub.Step(&orderer.StepRequest{})
		assert.EqualError(t, err, "rpc error: code = Unknown desc = badly formatted message, cannot extract channel")

		
		_, err = node1.c.DispatchSubmit(context.Background(), &orderer.SubmitRequest{})
		assert.EqualError(t, err, "badly formatted message, cannot extract channel")
	})
}

func TestAbortRPC(t *testing.T) {
	t.Parallel()
	
	
	
	
	

	testCases := []struct {
		name       string
		abortFunc  func(*cluster.RemoteContext)
		rpcTimeout time.Duration
	}{
		{
			name:       "Abort() called",
			rpcTimeout: 0,
			abortFunc: func(rc *cluster.RemoteContext) {
				rc.Abort()
			},
		},
		{
			name:       "RPC timeout",
			rpcTimeout: time.Millisecond * 100,
			abortFunc:  func(*cluster.RemoteContext) {},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			testAbort(t, testCase.abortFunc, testCase.rpcTimeout)
		})
	}
}

func testAbort(t *testing.T, abortFunc func(*cluster.RemoteContext), rpcTimeout time.Duration) {
	node1 := newTestNode(t)
	node1.c.RPCTimeout = rpcTimeout
	defer node1.stop()

	node2 := newTestNode(t)
	defer node2.stop()

	config := []cluster.RemoteNode{node1.nodeInfo, node2.nodeInfo}
	node1.c.Configure(testChannel, config)
	node2.c.Configure(testChannel, config)
	var onStepCalled sync.WaitGroup
	onStepCalled.Add(1)

	
	var stuckCall sync.WaitGroup
	stuckCall.Add(1)
	
	defer stuckCall.Done()

	node2.handler.On("OnStep", testChannel, node1.nodeInfo.ID, mock.Anything).Return(testStepRes, nil).Once().Run(func(_ mock.Arguments) {
		onStepCalled.Done()
		stuckCall.Wait()
	}).Once()

	rm, err := node1.c.Remote(testChannel, node2.nodeInfo.ID)
	assert.NoError(t, err)

	go func() {
		onStepCalled.Wait()
		abortFunc(rm)
	}()
	rm.Step(testStepReq)
	node2.handler.AssertNumberOfCalls(t, "OnStep", 1)
}

func TestNoTLSCertificate(t *testing.T) {
	t.Parallel()
	
	
	node1 := newTestNode(t)
	defer node1.stop()

	node1.c.Configure(testChannel, []cluster.RemoteNode{node1.nodeInfo})

	clientConfig := comm_utils.ClientConfig{
		Timeout: time.Millisecond * 100,
		SecOpts: &comm_utils.SecureOptions{
			ServerRootCAs: [][]byte{ca.CertBytes()},
			UseTLS:        true,
		},
	}
	cl, err := comm_utils.NewGRPCClient(clientConfig)
	assert.NoError(t, err)
	conn, err := cl.NewConnection(node1.srv.Address(), "")
	assert.NoError(t, err)
	echoClient := orderer.NewClusterClient(conn)
	_, err = echoClient.Step(context.Background(), testStepReq)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = no TLS certificate sent")
}

func TestReconnect(t *testing.T) {
	t.Parallel()
	
	
	
	
	
	

	node1 := newTestNode(t)
	defer node1.stop()

	node2 := newTestNode(t)
	node2.handler.On("OnStep", testChannel, node1.nodeInfo.ID, mock.Anything).Return(testStepRes, nil)
	defer node2.stop()

	config := []cluster.RemoteNode{node1.nodeInfo, node2.nodeInfo}
	node1.c.Configure(testChannel, config)
	node2.c.Configure(testChannel, config)

	
	node2.srv.Stop()
	
	
	stub, err := node1.c.Remote(testChannel, node2.nodeInfo.ID)
	assert.NoError(t, err)
	
	
	_, err = stub.Step(testStepReq)
	assert.Contains(t, err.Error(), "rpc error: code = Unavailable")
	
	node2.resurrect()
	
	
	assertEventuallyConnect(t, stub, testStepReq)
}

func TestRenewCertificates(t *testing.T) {
	t.Parallel()
	
	
	
	
	

	node1 := newTestNode(t)
	defer node1.stop()

	node2 := newTestNode(t)
	defer node2.stop()

	node1.handler.On("OnStep", testChannel, node2.nodeInfo.ID, mock.Anything).Return(testStepRes, nil)
	node2.handler.On("OnStep", testChannel, node1.nodeInfo.ID, mock.Anything).Return(testStepRes, nil)

	config := []cluster.RemoteNode{node1.nodeInfo, node2.nodeInfo}
	node1.c.Configure(testChannel, config)
	node2.c.Configure(testChannel, config)

	assertBiDiCommunication(t, node1, node2, testStepReq)

	
	node1.renewCertificates()
	node2.renewCertificates()

	
	config = []cluster.RemoteNode{node1.nodeInfo, node2.nodeInfo}
	node1.c.Configure(testChannel, config)
	node2.c.Configure(testChannel, config)

	
	
	
	_, err := node1.c.Remote(testChannel, node2.nodeInfo.ID)
	assert.Error(t, err)

	
	node1.srv.Stop()
	node1.resurrect()
	node2.srv.Stop()
	node2.resurrect()

	
	assertBiDiCommunication(t, node1, node2, testStepReq)
}

func TestMembershipReconfiguration(t *testing.T) {
	t.Parallel()
	
	
	
	
	

	node1 := newTestNode(t)
	defer node1.stop()

	node2 := newTestNode(t)
	defer node2.stop()

	node1.c.Configure(testChannel, []cluster.RemoteNode{})
	node2.c.Configure(testChannel, []cluster.RemoteNode{node1.nodeInfo})

	
	_, err := node1.c.Remote(testChannel, node2.nodeInfo.ID)
	assert.EqualError(t, err, fmt.Sprintf("node %d doesn't exist in channel test's membership", node2.nodeInfo.ID))
	

	gt := gomega.NewGomegaWithT(t)
	gt.Eventually(func() (bool, error) {
		_, err := node2.c.Remote(testChannel, node1.nodeInfo.ID)
		return true, err
	})

	stub, err := node2.c.Remote(testChannel, node1.nodeInfo.ID)
	_, err = stub.Step(testStepReq)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = certificate extracted from TLS connection isn't authorized")

	
	node1.handler.On("OnStep", testChannel, node2.nodeInfo.ID, mock.Anything).Return(testStepRes, nil)
	node2.handler.On("OnStep", testChannel, node1.nodeInfo.ID, mock.Anything).Return(testStepRes, nil)
	node1.c.Configure(testChannel, []cluster.RemoteNode{node2.nodeInfo})

	
	assertBiDiCommunication(t, node1, node2, testStepReq)
	assertBiDiCommunication(t, node2, node1, testStepReq)

	
	node2.c.Configure(testChannel, []cluster.RemoteNode{})
	
	stub, err = node1.c.Remote(testChannel, node2.nodeInfo.ID)
	
	_, err = stub.Step(testStepReq)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = certificate extracted from TLS connection isn't authorized")
}

func TestShutdown(t *testing.T) {
	t.Parallel()
	
	

	node1 := newTestNode(t)
	defer node1.stop()

	node1.c.Shutdown()

	
	_, err := node1.c.Remote(testChannel, node1.nodeInfo.ID)
	assert.EqualError(t, err, "communication has been shut down")

	node2 := newTestNode(t)
	defer node2.stop()

	node2.c.Configure(testChannel, []cluster.RemoteNode{node1.nodeInfo})
	
	node1.c.Configure(testChannel, []cluster.RemoteNode{node2.nodeInfo})
	stub, err := node2.c.Remote(testChannel, node1.nodeInfo.ID)
	assert.NoError(t, err)
	
	_, err = stub.Step(testStepReq)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = channel test doesn't exist")
}

func TestMultiChannelConfig(t *testing.T) {
	t.Parallel()
	
	
	
	
	

	node1 := newTestNode(t)
	defer node1.stop()

	node2 := newTestNode(t)
	defer node2.stop()

	node3 := newTestNode(t)
	defer node3.stop()

	node1.c.Configure("foo", []cluster.RemoteNode{node2.nodeInfo})
	node1.c.Configure("bar", []cluster.RemoteNode{node3.nodeInfo})
	node2.c.Configure("foo", []cluster.RemoteNode{node1.nodeInfo})
	node3.c.Configure("bar", []cluster.RemoteNode{node1.nodeInfo})

	t.Run("Correct channel", func(t *testing.T) {
		node1.handler.On("OnStep", "foo", node2.nodeInfo.ID, mock.Anything).Return(fooRes, nil)
		node1.handler.On("OnStep", "bar", node3.nodeInfo.ID, mock.Anything).Return(barRes, nil)

		node2toNode1, err := node2.c.Remote("foo", node1.nodeInfo.ID)
		assert.NoError(t, err)
		node3toNode1, err := node3.c.Remote("bar", node1.nodeInfo.ID)
		assert.NoError(t, err)

		res, err := node2toNode1.Step(fooReq)
		assert.NoError(t, err)
		assert.Equal(t, string(res.Payload), fooReq.Channel)

		res, err = node3toNode1.Step(barReq)
		assert.NoError(t, err)
		assert.Equal(t, string(res.Payload), barReq.Channel)
	})

	t.Run("Incorrect channel", func(t *testing.T) {
		node2toNode1, err := node2.c.Remote("foo", node1.nodeInfo.ID)
		assert.NoError(t, err)
		node3toNode1, err := node3.c.Remote("bar", node1.nodeInfo.ID)
		assert.NoError(t, err)

		_, err = node2toNode1.Step(barReq)
		assert.EqualError(t, err, "rpc error: code = Unknown desc = certificate extracted from TLS connection isn't authorized")

		_, err = node3toNode1.Step(fooReq)
		assert.EqualError(t, err, "rpc error: code = Unknown desc = certificate extracted from TLS connection isn't authorized")
	})
}

func assertBiDiCommunication(t *testing.T, node1, node2 *clusterNode, msgToSend *orderer.StepRequest) {
	for _, tst := range []struct {
		label    string
		sender   *clusterNode
		receiver *clusterNode
		target   uint64
	}{
		{label: "1->2", sender: node1, target: node2.nodeInfo.ID, receiver: node2},
		{label: "2->1", sender: node2, target: node1.nodeInfo.ID, receiver: node1},
	} {
		t.Run(tst.label, func(t *testing.T) {
			stub, err := tst.sender.c.Remote(testChannel, tst.target)
			assert.NoError(t, err)

			msg, err := stub.Step(msgToSend)
			assert.NoError(t, err)
			assert.Equal(t, msg.Payload, msgToSend.Payload)

			expectedRes := &orderer.SubmitResponse{}
			tst.receiver.handler.On("OnSubmit", testChannel, tst.sender.nodeInfo.ID, mock.Anything).Return(expectedRes, nil).Once()
			stream, err := stub.SubmitStream()
			assert.NoError(t, err)

			err = stream.Send(testSubReq)
			assert.NoError(t, err)

			res, err := stream.Recv()
			assert.NoError(t, err)

			assert.Equal(t, expectedRes, res)
		})
	}
}

func assertEventuallyConnect(t *testing.T, rpc *cluster.RemoteContext, req *orderer.StepRequest) {
	gt := gomega.NewGomegaWithT(t)
	gt.Eventually(func() (bool, error) {
		res, err := rpc.Step(req)
		if err != nil {
			return false, err
		}
		return bytes.Equal(res.Payload, req.Payload), nil
	}, timeout).Should(gomega.BeTrue())
}
