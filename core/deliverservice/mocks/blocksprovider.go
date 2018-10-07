/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	"context"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	gossip_common "github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/protos/common"
	gossip_proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/utils"
	"google.golang.org/grpc"
)




type MockGossipServiceAdapter struct {
	addPayloadCnt int32

	GossipBlockDisseminations chan uint64
}

type MockAtomicBroadcastClient struct {
	BD *MockBlocksDeliverer
}

func (mabc *MockAtomicBroadcastClient) Broadcast(ctx context.Context, opts ...grpc.CallOption) (orderer.AtomicBroadcast_BroadcastClient, error) {
	panic("Should not be used")
}
func (mabc *MockAtomicBroadcastClient) Deliver(ctx context.Context, opts ...grpc.CallOption) (orderer.AtomicBroadcast_DeliverClient, error) {
	return mabc.BD, nil
}


func (*MockGossipServiceAdapter) PeersOfChannel(gossip_common.ChainID) []discovery.NetworkMember {
	return []discovery.NetworkMember{}
}


func (mock *MockGossipServiceAdapter) AddPayload(chainID string, payload *gossip_proto.Payload) error {
	atomic.AddInt32(&mock.addPayloadCnt, 1)
	return nil
}


func (mock *MockGossipServiceAdapter) AddPayloadCount() int32 {
	return atomic.LoadInt32(&mock.addPayloadCnt)
}


func (mock *MockGossipServiceAdapter) Gossip(msg *gossip_proto.GossipMessage) {
	mock.GossipBlockDisseminations <- msg.GetDataMsg().Payload.SeqNum
}



type MockBlocksDeliverer struct {
	DisconnectCalled           chan struct{}
	DisconnectAndDisableCalled chan struct{}
	CloseCalled                chan struct{}
	Pos                        uint64
	grpc.ClientStream
	recvCnt  int32
	MockRecv func(mock *MockBlocksDeliverer) (*orderer.DeliverResponse, error)
}



func (mock *MockBlocksDeliverer) Recv() (*orderer.DeliverResponse, error) {
	atomic.AddInt32(&mock.recvCnt, 1)
	return mock.MockRecv(mock)
}


func (mock *MockBlocksDeliverer) RecvCount() int32 {
	return atomic.LoadInt32(&mock.recvCnt)
}


func MockRecv(mock *MockBlocksDeliverer) (*orderer.DeliverResponse, error) {
	pos := mock.Pos

	
	mock.Pos++
	return &orderer.DeliverResponse{
		Type: &orderer.DeliverResponse_Block{
			Block: &common.Block{
				Header: &common.BlockHeader{
					Number:       pos,
					DataHash:     []byte{},
					PreviousHash: []byte{},
				},
				Data: &common.BlockData{
					Data: [][]byte{},
				},
			}},
	}, nil
}



func (mock *MockBlocksDeliverer) Send(env *common.Envelope) error {
	payload, _ := utils.GetPayload(env)
	seekInfo := &orderer.SeekInfo{}

	proto.Unmarshal(payload.Data, seekInfo)

	
	switch t := seekInfo.Start.Type.(type) {
	case *orderer.SeekPosition_Oldest:
		mock.Pos = 0
	case *orderer.SeekPosition_Specified:
		mock.Pos = t.Specified.Number
	}
	return nil
}

func (mock *MockBlocksDeliverer) Disconnect(disableEndpoint bool) {
	if disableEndpoint {
		mock.DisconnectAndDisableCalled <- struct{}{}
	} else {
		mock.DisconnectCalled <- struct{}{}
	}
}

func (mock *MockBlocksDeliverer) Close() {
	if mock.CloseCalled == nil {
		return
	}
	mock.CloseCalled <- struct{}{}
}

func (mock *MockBlocksDeliverer) UpdateEndpoints(endpoints []string) {

}

func (mock *MockBlocksDeliverer) GetEndpoints() []string {
	return []string{} 
}



type MockLedgerInfo struct {
	Height uint64
}


func (li *MockLedgerInfo) LedgerHeight() (uint64, error) {
	return atomic.LoadUint64(&li.Height), nil
}
