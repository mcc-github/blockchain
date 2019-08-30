/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/pkg/errors"
)

const (
	
	
	ChannelCreationPolicyKey = "ChannelCreationPolicy"
)


type ConsortiumProtos struct {
	ChannelCreationPolicy *cb.Policy
}


type ConsortiumConfig struct {
	protos *ConsortiumProtos
	orgs   map[string]Org
}


func NewConsortiumConfig(consortiumGroup *cb.ConfigGroup, mspConfig *MSPConfigHandler) (*ConsortiumConfig, error) {
	cc := &ConsortiumConfig{
		protos: &ConsortiumProtos{},
		orgs:   make(map[string]Org),
	}

	if err := DeserializeProtoValuesFromGroup(consortiumGroup, cc.protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	for orgName, orgGroup := range consortiumGroup.Groups {
		var err error
		if cc.orgs[orgName], err = NewOrganizationConfig(orgName, orgGroup, mspConfig); err != nil {
			return nil, err
		}
	}

	return cc, nil
}


func (cc *ConsortiumConfig) Organizations() map[string]Org {
	return cc.orgs
}



func (cc *ConsortiumConfig) ChannelCreationPolicy() *cb.Policy {
	return cc.protos.ChannelCreationPolicy
}
