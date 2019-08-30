/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/pkg/errors"
)




type MessageReceiver interface {
	
	Consensus(req *orderer.ConsensusRequest, sender uint64) error

	
	Submit(req *orderer.SubmitRequest, sender uint64) error
}




type ReceiverGetter interface {
	
	ReceiverByChain(channelID string) MessageReceiver
}


type Dispatcher struct {
	Logger        *flogging.FabricLogger
	ChainSelector ReceiverGetter
}


func (d *Dispatcher) OnConsensus(channel string, sender uint64, request *orderer.ConsensusRequest) error {
	receiver := d.ChainSelector.ReceiverByChain(channel)
	if receiver == nil {
		d.Logger.Warningf("An attempt to send a consensus request to a non existing channel (%s) was made by %d", channel, sender)
		return errors.Errorf("channel %s doesn't exist", channel)
	}
	return receiver.Consensus(request, sender)
}


func (d *Dispatcher) OnSubmit(channel string, sender uint64, request *orderer.SubmitRequest) error {
	receiver := d.ChainSelector.ReceiverByChain(channel)
	if receiver == nil {
		d.Logger.Warningf("An attempt to submit a transaction to a non existing channel (%s) was made by %d", channel, sender)
		return errors.Errorf("channel %s doesn't exist", channel)
	}
	return receiver.Submit(request, sender)
}
