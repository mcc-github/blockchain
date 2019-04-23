/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/pkg/errors"
)

type BroadcastClient interface {
	
	Send(env *cb.Envelope) error
	Close() error
}

type broadcastClient struct {
	client ab.AtomicBroadcast_BroadcastClient
}


func GetBroadcastClient() (BroadcastClient, error) {
	oc, err := NewOrdererClientFromEnv()
	if err != nil {
		return nil, err
	}
	bc, err := oc.Broadcast()
	if err != nil {
		return nil, err
	}

	return &broadcastClient{client: bc}, nil
}

func (s *broadcastClient) getAck() error {
	msg, err := s.client.Recv()
	if err != nil {
		return err
	}
	if msg.Status != cb.Status_SUCCESS {
		return errors.Errorf("got unexpected status: %v -- %s", msg.Status, msg.Info)
	}
	return nil
}


func (s *broadcastClient) Send(env *cb.Envelope) error {
	if err := s.client.Send(env); err != nil {
		return errors.WithMessage(err, "could not send")
	}

	err := s.getAck()

	return err
}

func (s *broadcastClient) Close() error {
	return s.client.CloseSend()
}
