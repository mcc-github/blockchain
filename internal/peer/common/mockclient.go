/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	grpc "google.golang.org/grpc"
)


func GetMockEndorserClient(response *pb.ProposalResponse, err error) pb.EndorserClient {
	return &mockEndorserClient{
		response: response,
		err:      err,
	}
}

type mockEndorserClient struct {
	response *pb.ProposalResponse
	err      error
}

func (m *mockEndorserClient) ProcessProposal(ctx context.Context, in *pb.SignedProposal, opts ...grpc.CallOption) (*pb.ProposalResponse, error) {
	return m.response, m.err
}

func GetMockBroadcastClient(err error) BroadcastClient {
	return &mockBroadcastClient{err: err}
}


type mockBroadcastClient struct {
	err error
}

func (m *mockBroadcastClient) Send(env *cb.Envelope) error {
	return m.err
}

func (m *mockBroadcastClient) Close() error {
	return nil
}

func GetMockAdminClient(err error) pb.AdminClient {
	return &mockAdminClient{err: err}
}

type mockAdminClient struct {
	status *pb.ServerStatus
	err    error
}

func (m *mockAdminClient) GetStatus(ctx context.Context, in *cb.Envelope, opts ...grpc.CallOption) (*pb.ServerStatus, error) {
	return m.status, m.err
}

func (m *mockAdminClient) DumpStackTrace(ctx context.Context, in *cb.Envelope, opts ...grpc.CallOption) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (m *mockAdminClient) StartServer(ctx context.Context, in *cb.Envelope, opts ...grpc.CallOption) (*pb.ServerStatus, error) {
	m.status = &pb.ServerStatus{Status: pb.ServerStatus_STARTED}
	return m.status, m.err
}

func (m *mockAdminClient) GetModuleLogLevel(ctx context.Context, env *cb.Envelope, opts ...grpc.CallOption) (*pb.LogLevelResponse, error) {
	op := &pb.AdminOperation{}
	pl := &cb.Payload{}
	proto.Unmarshal(env.Payload, pl)
	proto.Unmarshal(pl.Data, op)
	response := &pb.LogLevelResponse{LogModule: op.GetLogReq().LogModule, LogLevel: "INFO"}
	return response, m.err
}

func (m *mockAdminClient) SetModuleLogLevel(ctx context.Context, env *cb.Envelope, opts ...grpc.CallOption) (*pb.LogLevelResponse, error) {
	op := &pb.AdminOperation{}
	pl := &cb.Payload{}
	proto.Unmarshal(env.Payload, pl)
	proto.Unmarshal(pl.Data, op)
	response := &pb.LogLevelResponse{LogModule: op.GetLogReq().LogModule, LogLevel: op.GetLogReq().LogLevel}
	return response, m.err
}

func (m *mockAdminClient) RevertLogLevels(ctx context.Context, in *cb.Envelope, opts ...grpc.CallOption) (*empty.Empty, error) {
	return &empty.Empty{}, m.err
}

func (m *mockAdminClient) GetLogSpec(ctx context.Context, in *cb.Envelope, opts ...grpc.CallOption) (*pb.LogSpecResponse, error) {
	response := &pb.LogSpecResponse{LogSpec: "info"}
	return response, m.err
}

func (m *mockAdminClient) SetLogSpec(ctx context.Context, in *cb.Envelope, opts ...grpc.CallOption) (*pb.LogSpecResponse, error) {
	response := &pb.LogSpecResponse{LogSpec: "info"}
	return response, m.err
}