/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"context"

	"github.com/mcc-github/blockchain/peer/chaincode/api"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"google.golang.org/grpc"
)




type DeliverClient interface {
	Deliver(ctx context.Context, opts ...grpc.CallOption) (DeliverService, error)
}




type DeliverService interface {
	Send(*cb.Envelope) error
	Recv() (*ab.DeliverResponse, error)
	CloseSend() error
}




type PeerDeliverClient interface {
	Deliver(ctx context.Context, opts ...grpc.CallOption) (api.Deliver, error)
	DeliverFiltered(ctx context.Context, opts ...grpc.CallOption) (api.Deliver, error)
}
