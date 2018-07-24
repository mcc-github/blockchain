/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"context"

	"github.com/mcc-github/blockchain/cmd/common"
	"github.com/mcc-github/blockchain/cmd/common/comm"
	"github.com/mcc-github/blockchain/cmd/common/signer"
	"github.com/mcc-github/blockchain/discovery/client"
	. "github.com/mcc-github/blockchain/protos/discovery"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)






type ServiceResponse interface {
	
	ForChannel(string) discovery.ChannelResponse

	
	ForLocal() discovery.LocalResponse

	
	Raw() *Response
}

type response struct {
	raw *Response
	discovery.Response
}

func (r *response) Raw() *Response {
	return r.raw
}



type ClientStub struct {
}


func (stub *ClientStub) Send(server string, conf common.Config, req *discovery.Request) (ServiceResponse, error) {
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

	disc := discovery.NewClient(comm.NewDialer(server), signer.Sign, 0)

	resp, err := disc.Send(timeout, req, &AuthInfo{
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


func (stub *RawStub) Send(server string, conf common.Config, req *discovery.Request) (ServiceResponse, error) {
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

	req.Authentication = &AuthInfo{
		ClientIdentity:    signer.Creator,
		ClientTlsCertHash: comm.TLSCertHash,
	}

	payload := utils.MarshalOrPanic(req.Request)
	sig, err := signer.Sign(payload)
	if err != nil {
		return nil, err
	}

	cc, err := comm.NewDialer(server)()
	if err != nil {
		return nil, err
	}
	resp, err := NewDiscoveryClient(cc).Discover(timeout, &SignedRequest{
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
