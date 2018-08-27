/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)




type RemoteCommunicator interface {
	
	
	
	Remote(channel string, id uint64) (*RemoteContext, error)
}




type SubmitClient interface {
	Send(request *orderer.SubmitRequest) error
	Recv() (*orderer.SubmitResponse, error)
	grpc.ClientStream
}





type Client interface {
	
	Submit(ctx context.Context, opts ...grpc.CallOption) (orderer.Cluster_SubmitClient, error)
	
	Step(ctx context.Context, in *orderer.StepRequest, opts ...grpc.CallOption) (*orderer.StepResponse, error)
}


type RPC struct {
	stream  orderer.Cluster_SubmitClient
	Channel string
	Comm    RemoteCommunicator
}


func (s *RPC) Step(destination uint64, msg *orderer.StepRequest) (*orderer.StepResponse, error) {
	stub, err := s.Comm.Remote(s.Channel, destination)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return stub.Step(msg)
}


func (s *RPC) SendSubmit(destination uint64, request *orderer.SubmitRequest) error {
	stream, err := s.getProposeStream(destination)
	if err != nil {
		return err
	}
	err = stream.Send(request)
	if err != nil {
		s.stream = nil
	}
	return err
}


func (s *RPC) ReceiveSubmitResponse(destination uint64) (*orderer.SubmitResponse, error) {
	stream, err := s.getProposeStream(destination)
	if err != nil {
		return nil, err
	}
	msg, err := stream.Recv()
	if err != nil {
		s.stream = nil
	}
	return msg, err
}


func (s *RPC) getProposeStream(destination uint64) (orderer.Cluster_SubmitClient, error) {
	if s.stream != nil {
		return s.stream, nil
	}
	stub, err := s.Comm.Remote(s.Channel, destination)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stream, err := stub.SubmitStream()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	s.stream = stream
	return stream, nil
}
