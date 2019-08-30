/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"context"

	"github.com/mcc-github/blockchain-protos-go/discovery"
	"github.com/mcc-github/blockchain/cmd/common"
	"github.com/mcc-github/blockchain/cmd/common/comm"
	"github.com/mcc-github/blockchain/cmd/common/signer"
	discoveryclient "github.com/mcc-github/blockchain/discovery/client"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)




type LocalResponse interface {
	discoveryclient.LocalResponse
}




type ChannelResponse interface {
	discoveryclient.ChannelResponse
}




type ServiceResponse interface {
	
	ForChannel(string) discoveryclient.ChannelResponse

	
	ForLocal() discoveryclient.LocalResponse

	
	Raw() *discovery.Response
}

type response struct {
	raw *discovery.Response
	discoveryclient.Response
}

func (r *response) Raw() *discovery.Response {
	return r.raw
}



type ClientStub struct {
}


func (stub *ClientStub) Send(server string, conf common.Config, req *discoveryclient.Request) (ServiceResponse, error) {
	comm, err := comm.NewClient(conf.TLSConfig)
	if err != nil {
		return nil, err
	}
	signer, err := signer.NewSigner(conf.SignerConfig)
	if err != nil {
		return nil, err
	}
	timeout, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	disc := discoveryclient.NewClient(comm.NewDialer(server), signer.Sign, 0)

	resp, err := disc.Send(timeout, req, &discovery.AuthInfo{
		ClientIdentity:    signer.Creator,
		ClientTlsCertHash: comm.TLSCertHash,
	})
	if err != nil {
		return nil, errors.Errorf("failed connecting to %s: %v", server, err)
	}
	return &response{
		Response: resp,
	}, nil
}



type RawStub struct {
}


func (stub *RawStub) Send(server string, conf common.Config, req *discoveryclient.Request) (ServiceResponse, error) {
	comm, err := comm.NewClient(conf.TLSConfig)
	if err != nil {
		return nil, err
	}
	signer, err := signer.NewSigner(conf.SignerConfig)
	if err != nil {
		return nil, err
	}
	timeout, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	req.Authentication = &discovery.AuthInfo{
		ClientIdentity:    signer.Creator,
		ClientTlsCertHash: comm.TLSCertHash,
	}

	payload := protoutil.MarshalOrPanic(req.Request)
	sig, err := signer.Sign(payload)
	if err != nil {
		return nil, err
	}

	cc, err := comm.NewDialer(server)()
	if err != nil {
		return nil, err
	}
	resp, err := discovery.NewDiscoveryClient(cc).Discover(timeout, &discovery.SignedRequest{
		Payload:   payload,
		Signature: sig,
	})

	if err != nil {
		return nil, err
	}

	return &response{
		raw: resp,
	}, nil
}
