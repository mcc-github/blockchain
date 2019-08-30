/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"context"
	"sync"
	"time"

	"github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)




type StepClient interface {
	Send(*orderer.StepRequest) error
	Recv() (*orderer.StepResponse, error)
	grpc.ClientStream
}




type ClusterClient interface {
	Step(ctx context.Context, opts ...grpc.CallOption) (orderer.Cluster_StepClient, error)
}


type RPC struct {
	consensusLock sync.Mutex
	submitLock    sync.Mutex
	Logger        *flogging.FabricLogger
	Timeout       time.Duration
	Channel       string
	Comm          Communicator
	lock          sync.RWMutex
	StreamsByType map[OperationType]map[uint64]*Stream
}



func NewStreamsByType() map[OperationType]map[uint64]*Stream {
	m := make(map[OperationType]map[uint64]*Stream)
	m[ConsensusOperation] = make(map[uint64]*Stream)
	m[SubmitOperation] = make(map[uint64]*Stream)
	return m
}



type OperationType int

const (
	ConsensusOperation OperationType = iota
	SubmitOperation
)


func (s *RPC) SendConsensus(destination uint64, msg *orderer.ConsensusRequest) error {
	if s.Logger.IsEnabledFor(zapcore.DebugLevel) {
		defer s.consensusSent(time.Now(), destination, msg)
	}

	stream, err := s.getOrCreateStream(destination, ConsensusOperation)
	if err != nil {
		return err
	}

	req := &orderer.StepRequest{
		Payload: &orderer.StepRequest_ConsensusRequest{
			ConsensusRequest: msg,
		},
	}

	s.consensusLock.Lock()
	defer s.consensusLock.Unlock()

	err = stream.Send(req)
	if err != nil {
		s.unMapStream(destination, ConsensusOperation)
	}

	return err
}


func (s *RPC) SendSubmit(destination uint64, request *orderer.SubmitRequest) error {
	if s.Logger.IsEnabledFor(zapcore.DebugLevel) {
		defer s.submitSent(time.Now(), destination, request)
	}

	stream, err := s.getOrCreateStream(destination, SubmitOperation)
	if err != nil {
		return err
	}

	req := &orderer.StepRequest{
		Payload: &orderer.StepRequest_SubmitRequest{
			SubmitRequest: request,
		},
	}

	s.submitLock.Lock()
	defer s.submitLock.Unlock()

	err = stream.Send(req)
	if err != nil {
		s.unMapStream(destination, SubmitOperation)
	}
	return err
}

func (s *RPC) submitSent(start time.Time, to uint64, msg *orderer.SubmitRequest) {
	s.Logger.Debugf("Sending msg of %d bytes to %d on channel %s took %v", submitMsgLength(msg), to, s.Channel, time.Since(start))
}

func (s *RPC) consensusSent(start time.Time, to uint64, msg *orderer.ConsensusRequest) {
	s.Logger.Debugf("Sending msg of %d bytes to %d on channel %s took %v", len(msg.Payload), to, s.Channel, time.Since(start))
}


func (s *RPC) getOrCreateStream(destination uint64, operationType OperationType) (orderer.Cluster_StepClient, error) {
	stream := s.getStream(destination, operationType)
	if stream != nil {
		return stream, nil
	}
	stub, err := s.Comm.Remote(s.Channel, destination)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stream, err = stub.NewStream(s.Timeout)
	if err != nil {
		return nil, err
	}
	s.mapStream(destination, stream, operationType)
	return stream, nil
}

func (s *RPC) getStream(destination uint64, operationType OperationType) *Stream {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.StreamsByType[operationType][destination]
}

func (s *RPC) mapStream(destination uint64, stream *Stream, operationType OperationType) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.StreamsByType[operationType][destination] = stream
	s.cleanCanceledStreams(operationType)
}

func (s *RPC) unMapStream(destination uint64, operationType OperationType) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.StreamsByType[operationType], destination)
}

func (s *RPC) cleanCanceledStreams(operationType OperationType) {
	for destination, stream := range s.StreamsByType[operationType] {
		if !stream.Canceled() {
			continue
		}
		s.Logger.Infof("Removing stream %d to %d for channel %s because it is canceled", stream.ID, destination, s.Channel)
		delete(s.StreamsByType[operationType], destination)
	}
}

func submitMsgLength(request *orderer.SubmitRequest) int {
	if request.Payload == nil {
		return 0
	}
	return len(request.Payload.Payload)
}
