/*
Copyright IBM Corp. 2016-2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package common

import (
	"context"
	"crypto/tls"

	ab "github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/pkg/errors"
)



type OrdererClient struct {
	CommonClient
}



func NewOrdererClientFromEnv() (*OrdererClient, error) {
	address, override, clientConfig, err := configFromEnv("orderer")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to load config for OrdererClient")
	}
	gClient, err := comm.NewGRPCClient(clientConfig)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create OrdererClient from config")
	}
	oClient := &OrdererClient{
		CommonClient: CommonClient{
			GRPCClient: gClient,
			Address:    address,
			sn:         override}}
	return oClient, nil
}


func (oc *OrdererClient) Broadcast() (ab.AtomicBroadcast_BroadcastClient, error) {
	conn, err := oc.CommonClient.NewConnection(oc.Address, comm.ServerNameOverride(oc.sn))
	if err != nil {
		return nil, errors.WithMessagef(err, "orderer client failed to connect to %s", oc.Address)
	}
	
	return ab.NewAtomicBroadcastClient(conn).Broadcast(context.TODO())
}


func (oc *OrdererClient) Deliver() (ab.AtomicBroadcast_DeliverClient, error) {
	conn, err := oc.CommonClient.NewConnection(oc.Address, comm.ServerNameOverride(oc.sn))
	if err != nil {
		return nil, errors.WithMessagef(err, "orderer client failed to connect to %s", oc.Address)
	}
	
	return ab.NewAtomicBroadcastClient(conn).Deliver(context.TODO())

}


func (oc *OrdererClient) Certificate() tls.Certificate {
	return oc.CommonClient.Certificate()
}
