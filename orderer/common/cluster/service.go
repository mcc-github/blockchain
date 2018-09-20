/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"context"
	"io"

	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/op/go-logging"
	"google.golang.org/grpc"
)




type Dispatcher interface {
	DispatchSubmit(ctx context.Context, request *orderer.SubmitRequest) (*orderer.SubmitResponse, error)
	DispatchStep(ctx context.Context, request *orderer.StepRequest) (*orderer.StepResponse, error)
}





type SubmitStream interface {
	Send(response *orderer.SubmitResponse) error
	Recv() (*orderer.SubmitRequest, error)
	grpc.ServerStream
}


type Service struct {
	Dispatcher Dispatcher
	Logger     logging.Logger
}


func (s *Service) Step(ctx context.Context, request *orderer.StepRequest) (*orderer.StepResponse, error) {
	addr := util.ExtractRemoteAddress(ctx)
	s.Logger.Debugf("Connection from %s", addr)
	defer s.Logger.Debugf("Closing connection from %s", addr)
	response, err := s.Dispatcher.DispatchStep(ctx, request)
	if err != nil {
		s.Logger.Warningf("Handling of Step() from %s failed: %+v", addr, err)
	}
	return response, err
}


func (s *Service) Submit(stream orderer.Cluster_SubmitServer) error {
	addr := util.ExtractRemoteAddress(stream.Context())
	s.Logger.Debugf("Connection from %s", addr)
	defer s.Logger.Debugf("Closing connection from %s", addr)
	for {
		err := s.handleSubmit(stream, addr)
		if err == io.EOF {
			s.Logger.Debugf("%s disconnected", addr)
			return nil
		}
		if err != nil {
			return err
		}
		
	}
}

func (s *Service) handleSubmit(stream SubmitStream, addr string) error {
	request, err := stream.Recv()
	if err == io.EOF {
		return err
	}
	if err != nil {
		s.Logger.Warningf("Stream read from %s failed: %v", addr, err)
		return err
	}
	response, err := s.Dispatcher.DispatchSubmit(stream.Context(), request)
	if err != nil {
		s.Logger.Warningf("Handling of Propose() from %s failed: %+v", addr, err)
		return err
	}
	err = stream.Send(response)
	if err != nil {
		s.Logger.Warningf("Send() failed: %v", err)
	}
	return err
}
