/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	cb "github.com/mcc-github/blockchain-protos-go/common"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/capabilities"
	"github.com/pkg/errors"
)

const (
	
	ApplicationGroupKey = "Application"

	
	ACLsKey = "ACLs"
)


type ApplicationProtos struct {
	ACLs         *pb.ACLs
	Capabilities *cb.Capabilities
}


type ApplicationConfig struct {
	applicationOrgs map[string]ApplicationOrg
	protos          *ApplicationProtos
}


func NewApplicationConfig(appGroup *cb.ConfigGroup, mspConfig *MSPConfigHandler) (*ApplicationConfig, error) {
	ac := &ApplicationConfig{
		applicationOrgs: make(map[string]ApplicationOrg),
		protos:          &ApplicationProtos{},
	}

	if err := DeserializeProtoValuesFromGroup(appGroup, ac.protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	if !ac.Capabilities().ACLs() {
		if _, ok := appGroup.Values[ACLsKey]; ok {
			return nil, errors.New("ACLs may not be specified without the required capability")
		}
	}

	var err error
	for orgName, orgGroup := range appGroup.Groups {
		ac.applicationOrgs[orgName], err = NewApplicationOrgConfig(orgName, orgGroup, mspConfig)
		if err != nil {
			return nil, err
		}
	}

	return ac, nil
}


func (ac *ApplicationConfig) Organizations() map[string]ApplicationOrg {
	return ac.applicationOrgs
}


func (ac *ApplicationConfig) Capabilities() ApplicationCapabilities {
	return capabilities.NewApplicationProvider(ac.protos.Capabilities.Capabilities)
}


func (ac *ApplicationConfig) APIPolicyMapper() PolicyMapper {
	pm := newAPIsProvider(ac.protos.ACLs.Acls)

	return pm
}
