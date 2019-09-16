/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"time"

	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)


type PeerClient struct {
	CommonClient
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

	override := viper.GetString("peer.tls.serverhostoverride")
	clientConfig := comm.ClientConfig{}
	clientConfig.Timeout = viper.GetDuration("peer.client.connTimeout")
	if clientConfig.Timeout == time.Duration(0) {
		clientConfig.Timeout = defaultConnTimeout
	}

	secOpts := comm.SecureOptions{
		UseTLS:            viper.GetBool("peer.tls.enabled"),
		RequireClientCert: viper.GetBool("peer.tls.clientAuthRequired"),
	}

	if secOpts.RequireClientCert {
		keyPEM, err := ioutil.ReadFile(config.GetPath("peer.tls.clientKey.file"))
		if err != nil {
			return nil, errors.WithMessage(err, "unable to load peer.tls.clientKey.file")
		}
		secOpts.Key = keyPEM
		certPEM, err := ioutil.ReadFile(config.GetPath("peer.tls.clientCert.file"))
		if err != nil {
			return nil, errors.WithMessage(err, "unable to load peer.tls.clientCert.file")
		}
		secOpts.Certificate = certPEM
	}
	clientConfig.SecOpts = secOpts

	if clientConfig.SecOpts.UseTLS {
		if tlsRootCertFile == "" {
			return nil, errors.New("tls root cert file must be set")
		}
		caPEM, res := ioutil.ReadFile(tlsRootCertFile)
		if res != nil {
			return nil, errors.WithMessagef(res, "unable to load TLS root cert file from %s", tlsRootCertFile)
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
		CommonClient: CommonClient{
			GRPCClient: gClient,
			Address:    address,
			sn:         override}}
	return pClient, nil
}


func (pc *PeerClient) Endorser() (pb.EndorserClient, error) {
	conn, err := pc.CommonClient.NewConnection(pc.Address, comm.ServerNameOverride(pc.sn))
	if err != nil {
		return nil, errors.WithMessagef(err, "endorser client failed to connect to %s", pc.Address)
	}
	return pb.NewEndorserClient(conn), nil
}


func (pc *PeerClient) Deliver() (pb.Deliver_DeliverClient, error) {
	conn, err := pc.CommonClient.NewConnection(pc.Address, comm.ServerNameOverride(pc.sn))
	if err != nil {
		return nil, errors.WithMessagef(err, "deliver client failed to connect to %s", pc.Address)
	}
	return pb.NewDeliverClient(conn).Deliver(context.TODO())
}



func (pc *PeerClient) PeerDeliver() (pb.DeliverClient, error) {
	conn, err := pc.CommonClient.NewConnection(pc.Address, comm.ServerNameOverride(pc.sn))
	if err != nil {
		return nil, errors.WithMessagef(err, "deliver client failed to connect to %s", pc.Address)
	}
	return pb.NewDeliverClient(conn), nil
}


func (pc *PeerClient) Certificate() tls.Certificate {
	return pc.CommonClient.Certificate()
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





func GetPeerDeliverClient(address, tlsRootCertFile string) (pb.DeliverClient, error) {
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
