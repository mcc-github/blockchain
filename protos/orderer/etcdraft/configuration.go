/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"fmt"
	"io/ioutil"

	"github.com/mcc-github/blockchain/protos/orderer"

	"github.com/golang/protobuf/proto"
)


const TypeKey = "etcdraft"

func init() {
	orderer.ConsensusTypeMetadataMap[TypeKey] = ConsensusTypeMetadataFactory{}
}



type ConsensusTypeMetadataFactory struct{}


func (dogf ConsensusTypeMetadataFactory) NewMessage() proto.Message {
	return &Metadata{}
}



func Marshal(md *Metadata) ([]byte, error) {
	for _, c := range md.Consenters {
		
		
		clientCert, err := ioutil.ReadFile(string(c.GetClientTlsCert()))
		if err != nil {
			return nil, fmt.Errorf("cannot load client cert for consenter %s:%d: %s", c.GetHost(), c.GetPort(), err)
		}
		c.ClientTlsCert = clientCert

		serverCert, err := ioutil.ReadFile(string(c.GetServerTlsCert()))
		if err != nil {
			return nil, fmt.Errorf("cannot load server cert for consenter %s:%d: %s", c.GetHost(), c.GetPort(), err)
		}
		c.ServerTlsCert = serverCert
	}
	return proto.Marshal(md)
}
