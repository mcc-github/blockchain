/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/pkg/errors"
)




type MessageReceiver interface {
	
	Step(req *orderer.StepRequest, sender uint64) error

	
	Submit(req *orderer.SubmitRequest, sender uint64) error
}




type ReceiverGetter interface {
	
	ReceiverByChain(channelID string) MessageReceiver
}


type Dispatcher struct {
	Logger        *flogging.FabricLogger
	ChainSelector ReceiverGetter
}


func (d *Dispatcher) OnStep(channel string, sender uint64, request *orderer.StepRequest) (*orderer.StepResponse, error) {
	receiver := d.ChainSelector.ReceiverByChain(channel)
	if receiver == nil {
		d.Logger.Warningf("An attempt to send a StepRequest to a non existing channel (%s) was made by %d", channel, sender)
		return nil, errors.Errorf("channel %s doesn't exist", channel)
	}
	return &orderer.StepResponse{}, receiver.Step(request, sender)
}


func (d *Dispatcher) OnSubmit(channel string, sender uint64, request *orderer.SubmitRequest) (*orderer.SubmitResponse, error) {
	receiver := d.ChainSelector.ReceiverByChain(channel)
	if receiver == nil {
		d.Logger.Warningf("An attempt to submit a transaction to a non existing channel (%s) was made by %d", channel, sender)
		return &orderer.SubmitResponse{
			Info:   fmt.Sprintf("channel %s doesn't exist", channel),
			Status: common.Status_NOT_FOUND,
		}, nil
	}
	if err := receiver.Submit(request, sender); err != nil {
		d.Logger.Errorf("Failed handling transaction on channel %s from %d: %+v", channel, sender, err)
		return &orderer.SubmitResponse{
			Info:   err.Error(),
			Status: common.Status_INTERNAL_SERVER_ERROR,
		}, nil
	}
	return &orderer.SubmitResponse{
		Status: common.Status_SUCCESS,
	}, nil
}
