/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster_test

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/common/cluster/mocks"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestRPCChangeDestination(t *testing.T) {
	t.Parallel()
	
	
	
	
	
	
	comm := &mocks.Communicator{}

	client1 := &mocks.ClusterClient{}
	client2 := &mocks.ClusterClient{}

	metrics := cluster.NewMetrics(&disabled.Provider{})

	comm.On("Remote", "mychannel", uint64(1)).Return(&cluster.RemoteContext{
		SendBuffSize: 10,
		Metrics:      metrics,
		Logger:       flogging.MustGetLogger("test"),
		Client:       client1,
		ProbeConn:    func(_ *grpc.ClientConn) error { return nil },
	}, nil)
	comm.On("Remote", "mychannel", uint64(2)).Return(&cluster.RemoteContext{
		SendBuffSize: 10,
		Metrics:      metrics,
		Logger:       flogging.MustGetLogger("test"),
		Client:       client2,
		ProbeConn:    func(_ *grpc.ClientConn) error { return nil },
	}, nil)

	streamToNode1 := &mocks.StepClient{}
	streamToNode2 := &mocks.StepClient{}
	streamToNode1.On("Context", mock.Anything).Return(context.Background())
	streamToNode2.On("Context", mock.Anything).Return(context.Background())

	client1.On("Step", mock.Anything).Return(streamToNode1, nil).Once()
	client2.On("Step", mock.Anything).Return(streamToNode2, nil).Once()

	rpc := &cluster.RPC{
		Logger:        flogging.MustGetLogger("test"),
		Timeout:       time.Hour,
		StreamsByType: cluster.NewStreamsByType(),
		Channel:       "mychannel",
		Comm:          comm,
	}

	var sent sync.WaitGroup
	sent.Add(2)

	signalSent := func(_ mock.Arguments) {
		sent.Done()
	}
	streamToNode1.On("Send", mock.Anything).Return(nil).Run(signalSent).Once()
	streamToNode2.On("Send", mock.Anything).Return(nil).Run(signalSent).Once()
	streamToNode1.On("Recv").Return(nil, io.EOF)
	streamToNode2.On("Recv").Return(nil, io.EOF)

	rpc.SendSubmit(1, &orderer.SubmitRequest{Channel: "mychannel"})
	rpc.SendSubmit(2, &orderer.SubmitRequest{Channel: "mychannel"})

	sent.Wait()
	streamToNode1.AssertNumberOfCalls(t, "Send", 1)
	streamToNode2.AssertNumberOfCalls(t, "Send", 1)
}

func TestSend(t *testing.T) {
	t.Parallel()
	submitRequest := &orderer.SubmitRequest{Channel: "mychannel"}
	submitResponse := &orderer.StepResponse{
		Payload: &orderer.StepResponse_SubmitRes{
			SubmitRes: &orderer.SubmitResponse{Status: common.Status_SUCCESS},
		},
	}

	consensusRequest := &orderer.ConsensusRequest{
		Channel: "mychannel",
	}

	submitReq := wrapSubmitReq(submitRequest)

	consensusReq := &orderer.StepRequest{
		Payload: &orderer.StepRequest_ConsensusRequest{
			ConsensusRequest: consensusRequest,
		},
	}

	submit := func(rpc *cluster.RPC) error {
		err := rpc.SendSubmit(1, submitRequest)
		return err
	}

	step := func(rpc *cluster.RPC) error {
		return rpc.SendConsensus(1, consensusRequest)
	}

	type testCase struct {
		name           string
		method         func(rpc *cluster.RPC) error
		sendReturns    error
		sendCalledWith *orderer.StepRequest
		receiveReturns []interface{}
		stepReturns    []interface{}
		remoteError    error
		expectedErr    string
	}

	l := &sync.Mutex{}
	var tst testCase

	sent := make(chan struct{})

	var sendCalls uint32

	stream := &mocks.StepClient{}
	stream.On("Context", mock.Anything).Return(context.Background())
	stream.On("Send", mock.Anything).Return(func(*orderer.StepRequest) error {
		l.Lock()
		defer l.Unlock()
		atomic.AddUint32(&sendCalls, 1)
		sent <- struct{}{}
		return tst.sendReturns
	})

	for _, tst := range []testCase{
		{
			name:           "Send and Receive submit succeed",
			method:         submit,
			sendReturns:    nil,
			stepReturns:    []interface{}{stream, nil},
			receiveReturns: []interface{}{submitResponse, nil},
			sendCalledWith: submitReq,
		},
		{
			name:           "Send step succeed",
			method:         step,
			sendReturns:    nil,
			stepReturns:    []interface{}{stream, nil},
			sendCalledWith: consensusReq,
		},
		{
			name:           "Send submit fails",
			method:         submit,
			sendReturns:    errors.New("oops"),
			stepReturns:    []interface{}{stream, nil},
			sendCalledWith: submitReq,
			expectedErr:    "stream is aborted",
		},
		{
			name:           "Send step fails",
			method:         step,
			sendReturns:    errors.New("oops"),
			stepReturns:    []interface{}{stream, nil},
			sendCalledWith: consensusReq,
			expectedErr:    "stream is aborted",
		},
		{
			name:        "Remote() fails",
			method:      submit,
			remoteError: errors.New("timed out"),
			stepReturns: []interface{}{stream, nil},
			expectedErr: "timed out",
		},
		{
			name:        "Submit fails with Send",
			method:      submit,
			stepReturns: []interface{}{nil, errors.New("deadline exceeded")},
			expectedErr: "deadline exceeded",
		},
	} {
		l.Lock()
		testCase := tst
		l.Unlock()

		t.Run(testCase.name, func(t *testing.T) {
			atomic.StoreUint32(&sendCalls, 0)
			isSend := testCase.receiveReturns == nil
			comm := &mocks.Communicator{}
			client := &mocks.ClusterClient{}
			client.On("Step", mock.Anything).Return(testCase.stepReturns...)
			rm := &cluster.RemoteContext{
				Metrics:      cluster.NewMetrics(&disabled.Provider{}),
				SendBuffSize: 1,
				Logger:       flogging.MustGetLogger("test"),
				ProbeConn:    func(_ *grpc.ClientConn) error { return nil },
				Client:       client,
			}
			defer rm.Abort()
			comm.On("Remote", "mychannel", uint64(1)).Return(rm, testCase.remoteError)

			rpc := &cluster.RPC{
				Logger:        flogging.MustGetLogger("test"),
				Timeout:       time.Hour,
				StreamsByType: cluster.NewStreamsByType(),
				Channel:       "mychannel",
				Comm:          comm,
			}

			var err error

			err = testCase.method(rpc)
			if testCase.remoteError == nil && testCase.stepReturns[1] == nil {
				<-sent
			}

			if testCase.stepReturns[1] == nil && testCase.remoteError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, testCase.expectedErr)
			}

			if testCase.remoteError == nil && testCase.expectedErr == "" && isSend {
				stream.AssertCalled(t, "Send", testCase.sendCalledWith)
				
				
				err := testCase.method(rpc)
				<-sent

				assert.NoError(t, err)
				assert.Equal(t, 2, int(atomic.LoadUint32(&sendCalls)))
				client.AssertNumberOfCalls(t, "Step", 1)
			}
		})
	}
}

func TestRPCGarbageCollection(t *testing.T) {
	
	
	
	
	

	t.Parallel()

	comm := &mocks.Communicator{}
	client := &mocks.ClusterClient{}
	stream := &mocks.StepClient{}

	remote := &cluster.RemoteContext{
		SendBuffSize: 10,
		Metrics:      cluster.NewMetrics(&disabled.Provider{}),
		Logger:       flogging.MustGetLogger("test"),
		Client:       client,
		ProbeConn:    func(_ *grpc.ClientConn) error { return nil },
	}

	var sent sync.WaitGroup

	defineMocks := func(destination uint64) {
		sent.Add(1)
		comm.On("Remote", "mychannel", destination).Return(remote, nil)
		stream.On("Context", mock.Anything).Return(context.Background())
		client.On("Step", mock.Anything).Return(stream, nil).Once()
		stream.On("Send", mock.Anything).Return(nil).Once().Run(func(_ mock.Arguments) {
			sent.Done()
		})
		stream.On("Recv").Return(nil, nil)
	}

	mapping := cluster.NewStreamsByType()

	rpc := &cluster.RPC{
		Logger:        flogging.MustGetLogger("test"),
		Timeout:       time.Hour,
		StreamsByType: mapping,
		Channel:       "mychannel",
		Comm:          comm,
	}

	defineMocks(1)

	rpc.SendSubmit(1, &orderer.SubmitRequest{Channel: "mychannel"})
	
	sent.Wait()
	
	assert.Len(t, mapping[cluster.SubmitOperation], 1)
	assert.Equal(t, uint64(1), mapping[cluster.SubmitOperation][1].ID)
	
	stream.AssertNumberOfCalls(t, "Send", 1)

	
	remote.Abort()

	
	assert.Len(t, mapping[cluster.SubmitOperation], 1)
	assert.Equal(t, uint64(1), mapping[cluster.SubmitOperation][1].ID)

	
	defineMocks(2)

	
	rpc.SendSubmit(2, &orderer.SubmitRequest{Channel: "mychannel"})
	
	assert.Len(t, mapping[cluster.SubmitOperation], 1)
	assert.Equal(t, uint64(2), mapping[cluster.SubmitOperation][2].ID)
}
