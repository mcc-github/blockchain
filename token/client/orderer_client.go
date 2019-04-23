/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strings"

	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)




type Broadcast interface {
	Send(m *common.Envelope) error
	Recv() (*ab.BroadcastResponse, error)
	CloseSend() error
}




type OrdererClient interface {
	
	NewBroadcast(ctx context.Context, opts ...grpc.CallOption) (Broadcast, error)

	
	Certificate() *tls.Certificate
}


type ordererClient struct {
	ordererAddr        string
	serverNameOverride string
	grpcClient         *comm.GRPCClient
	conn               *grpc.ClientConn
}

func NewOrdererClient(config *ConnectionConfig) (OrdererClient, error) {
	grpcClient, err := CreateGRPCClient(config)
	if err != nil {
		err = errors.WithMessagef(err, "failed to create a GRPCClient to orderer %s", config.Address)
		return nil, err
	}
	conn, err := grpcClient.NewConnection(config.Address, config.ServerNameOverride)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to connect to orderer %s", config.Address)
	}

	return &ordererClient{
		ordererAddr:        config.Address,
		serverNameOverride: config.ServerNameOverride,
		grpcClient:         grpcClient,
		conn:               conn,
	}, nil
}


func (oc *ordererClient) NewBroadcast(ctx context.Context, opts ...grpc.CallOption) (Broadcast, error) {
	
	broadcast, err := ab.NewAtomicBroadcastClient(oc.conn).Broadcast(ctx)
	if err == nil {
		return broadcast, nil
	}

	
	oc.conn, err = oc.grpcClient.NewConnection(oc.ordererAddr, oc.serverNameOverride)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to connect to orderer %s", oc.ordererAddr)
	}

	
	broadcast, err = ab.NewAtomicBroadcastClient(oc.conn).Broadcast(ctx)
	if err != nil {
		rpcStatus, _ := status.FromError(err)
		return nil, errors.Wrapf(err, "failed to new a broadcast, rpcStatus=%+v", rpcStatus)
	}
	return broadcast, nil
}

func (oc *ordererClient) Certificate() *tls.Certificate {
	cert := oc.grpcClient.Certificate()
	return &cert
}


func BroadcastSend(broadcast Broadcast, addr string, envelope *common.Envelope) error {
	err := broadcast.Send(envelope)
	broadcast.CloseSend()

	if err != nil {
		return errors.Wrapf(err, "failed to send transaction to orderer %s", addr)
	}
	return nil
}


func BroadcastReceive(broadcast Broadcast, addr string, responses chan common.Status, errs chan error) {
	for {
		broadcastResponse, err := broadcast.Recv()
		if err == io.EOF {
			close(responses)
			return
		}

		if err != nil {
			rpcStatus, _ := status.FromError(err)
			errs <- errors.Wrapf(err, "broadcast recv error from orderer %s, rpcStatus=%+v", addr, rpcStatus)
			close(responses)
			return
		}

		if broadcastResponse.Status == common.Status_SUCCESS {
			responses <- broadcastResponse.Status
		} else {
			errs <- errors.Errorf("broadcast response error %d from orderer %s", int32(broadcastResponse.Status), addr)
		}
	}
}


func BroadcastWaitForResponse(responses chan common.Status, errs chan error) (common.Status, error) {
	var status common.Status
	allErrs := make([]error, 0)

read:
	for {
		select {
		case s, ok := <-responses:
			if !ok {
				break read
			}
			status = s
		case e := <-errs:
			allErrs = append(allErrs, e)
		}
	}

	
	for i := 0; i < len(errs); i++ {
		e := <-errs
		allErrs = append(allErrs, e)
	}
	
	close(errs)
	return status, toError(allErrs)
}


func toError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}

	errmsgs := []string{fmt.Sprint("Multiple errors occurred in order broadcast stream: ")}
	for _, err := range errs {
		errmsgs = append(errmsgs, err.Error())
	}
	return errors.New(strings.Join(errmsgs, "\n"))
}
