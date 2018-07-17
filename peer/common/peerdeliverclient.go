/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"context"

	ccapi "github.com/mcc-github/blockchain/peer/chaincode/api"
	pb "github.com/mcc-github/blockchain/protos/peer"
	grpc "google.golang.org/grpc"
)



type PeerDeliverClient struct {
	Client pb.DeliverClient
}


func (dc PeerDeliverClient) Deliver(ctx context.Context, opts ...grpc.CallOption) (ccapi.Deliver, error) {
	d, err := dc.Client.Deliver(ctx, opts...)
	return d, err
}


func (dc PeerDeliverClient) DeliverFiltered(ctx context.Context, opts ...grpc.CallOption) (ccapi.Deliver, error) {
	df, err := dc.Client.DeliverFiltered(ctx, opts...)
	return df, err
}
