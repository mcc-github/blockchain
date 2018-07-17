/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"

	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/peer/common/api"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)


type PeerClient struct {
	commonClient
}



func NewPeerClientFromEnv() (*PeerClient, error) {
	address, override, clientConfig, err := configFromEnv("peer")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to load config for PeerClient")
	}
	return newPeerClientForClientConfig(address, override, clientConfig)
}



func NewPeerClientForAddress(address, tlsRootCertFile string) (*PeerClient, error) {
	if address == "" {
		return nil, errors.New("peer address must be set")
	}

	_, override, clientConfig, err := configFromEnv("peer")
	if clientConfig.SecOpts.UseTLS {
		if tlsRootCertFile == "" {
			return nil, errors.New("tls root cert file must be set")
		}
		caPEM, res := ioutil.ReadFile(tlsRootCertFile)
		if res != nil {
			err = errors.WithMessage(res, fmt.Sprintf("unable to load TLS root cert file from %s", tlsRootCertFile))
			return nil, err
		}
		clientConfig.SecOpts.ServerRootCAs = [][]byte{caPEM}
	}
	return newPeerClientForClientConfig(address, override, clientConfig)
}

func newPeerClientForClientConfig(address, override string, clientConfig comm.ClientConfig) (*PeerClient, error) {
	gClient, err := comm.NewGRPCClient(clientConfig)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create PeerClient from config")
	}
	pClient := &PeerClient{
		commonClient: commonClient{
			GRPCClient: gClient,
			address:    address,
			sn:         override}}
	return pClient, nil
}


func (pc *PeerClient) Endorser() (pb.EndorserClient, error) {
	conn, err := pc.commonClient.NewConnection(pc.address, pc.sn)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("endorser client failed to connect to %s", pc.address))
	}
	return pb.NewEndorserClient(conn), nil
}


func (pc *PeerClient) Deliver() (pb.Deliver_DeliverClient, error) {
	conn, err := pc.commonClient.NewConnection(pc.address, pc.sn)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("deliver client failed to connect to %s", pc.address))
	}
	return pb.NewDeliverClient(conn).Deliver(context.TODO())
}



func (pc *PeerClient) PeerDeliver() (api.PeerDeliverClient, error) {
	conn, err := pc.commonClient.NewConnection(pc.address, pc.sn)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("deliver client failed to connect to %s", pc.address))
	}
	pbClient := pb.NewDeliverClient(conn)
	return &PeerDeliverClient{Client: pbClient}, nil
}


func (pc *PeerClient) Admin() (pb.AdminClient, error) {
	conn, err := pc.commonClient.NewConnection(pc.address, pc.sn)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("admin client failed to connect to %s", pc.address))
	}
	return pb.NewAdminClient(conn), nil
}


func (pc *PeerClient) Certificate() tls.Certificate {
	return pc.commonClient.Certificate()
}





func GetEndorserClient(address, tlsRootCertFile string) (pb.EndorserClient, error) {
	var peerClient *PeerClient
	var err error
	if address != "" {
		peerClient, err = NewPeerClientForAddress(address, tlsRootCertFile)
	} else {
		peerClient, err = NewPeerClientFromEnv()
	}
	if err != nil {
		return nil, err
	}
	return peerClient.Endorser()
}


func GetCertificate() (tls.Certificate, error) {
	peerClient, err := NewPeerClientFromEnv()
	if err != nil {
		return tls.Certificate{}, err
	}
	return peerClient.Certificate(), nil
}



func GetAdminClient() (pb.AdminClient, error) {
	peerClient, err := NewPeerClientFromEnv()
	if err != nil {
		return nil, err
	}
	return peerClient.Admin()
}





func GetDeliverClient(address, tlsRootCertFile string) (pb.Deliver_DeliverClient, error) {
	var peerClient *PeerClient
	var err error
	if address != "" {
		peerClient, err = NewPeerClientForAddress(address, tlsRootCertFile)
	} else {
		peerClient, err = NewPeerClientFromEnv()
	}
	if err != nil {
		return nil, err
	}
	return peerClient.Deliver()
}





func GetPeerDeliverClient(address, tlsRootCertFile string) (api.PeerDeliverClient, error) {
	var peerClient *PeerClient
	var err error
	if address != "" {
		peerClient, err = NewPeerClientForAddress(address, tlsRootCertFile)
	} else {
		peerClient, err = NewPeerClientFromEnv()
	}
	if err != nil {
		return nil, err
	}
	return peerClient.PeerDeliver()
}
