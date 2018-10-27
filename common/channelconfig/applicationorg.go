/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"fmt"

	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)

const (
	
	AnchorPeersKey = "AnchorPeers"
)


type ApplicationOrgProtos struct {
	AnchorPeers *pb.AnchorPeers
}


type ApplicationOrgConfig struct {
	*OrganizationConfig
	protos *ApplicationOrgProtos
	name   string
}


func NewApplicationOrgConfig(id string, orgGroup *cb.ConfigGroup, mspConfig *MSPConfigHandler) (*ApplicationOrgConfig, error) {
	if len(orgGroup.Groups) > 0 {
		return nil, fmt.Errorf("ApplicationOrg config does not allow sub-groups")
	}

	protos := &ApplicationOrgProtos{}
	orgProtos := &OrganizationProtos{}

	if err := DeserializeProtoValuesFromGroup(orgGroup, protos, orgProtos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	aoc := &ApplicationOrgConfig{
		name:   id,
		protos: protos,
		OrganizationConfig: &OrganizationConfig{
			name:             id,
			protos:           orgProtos,
			mspConfigHandler: mspConfig,
		},
	}

	if err := aoc.Validate(); err != nil {
		return nil, err
	}

	return aoc, nil
}


func (aog *ApplicationOrgConfig) AnchorPeers() []*pb.AnchorPeer {
	return aog.protos.AnchorPeers.AnchorPeers
}

func (aoc *ApplicationOrgConfig) Validate() error {
	logger.Debugf("Anchor peers for org %s are %v", aoc.name, aoc.protos.AnchorPeers)
	return aoc.OrganizationConfig.Validate()
}
