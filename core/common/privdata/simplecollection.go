/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	m "github.com/mcc-github/blockchain-protos-go/msp"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)



type SimpleCollection struct {
	name         string
	accessPolicy policies.Policy
	memberOrgs   []string
	conf         common.StaticCollectionConfig
}

type SimpleCollectionPersistenceConfigs struct {
	blockToLive uint64
}



func NewSimpleCollection(collectionConfig *common.StaticCollectionConfig, deserializer msp.IdentityDeserializer) (*SimpleCollection, error) {
	sc := &SimpleCollection{}
	err := sc.Setup(collectionConfig, deserializer)
	return sc, err
}


func (sc *SimpleCollection) CollectionID() string {
	return sc.name
}


func (sc *SimpleCollection) MemberOrgs() []string {
	return sc.memberOrgs
}



func (sc *SimpleCollection) RequiredPeerCount() int {
	return int(sc.conf.RequiredPeerCount)
}



func (sc *SimpleCollection) MaximumPeerCount() int {
	return int(sc.conf.MaximumPeerCount)
}



func (sc *SimpleCollection) AccessFilter() Filter {
	return func(sd protoutil.SignedData) bool {
		if err := sc.accessPolicy.Evaluate([]*protoutil.SignedData{&sd}); err != nil {
			return false
		}
		return true
	}
}



func (sc *SimpleCollection) IsMemberOnlyRead() bool {
	return sc.conf.MemberOnlyRead
}



func (sc *SimpleCollection) IsMemberOnlyWrite() bool {
	return sc.conf.MemberOnlyWrite
}



func (sc *SimpleCollection) Setup(collectionConfig *common.StaticCollectionConfig, deserializer msp.IdentityDeserializer) error {
	if collectionConfig == nil {
		return errors.New("Nil config passed to collection setup")
	}
	sc.conf = *collectionConfig
	sc.name = collectionConfig.GetName()

	
	collectionPolicyConfig := collectionConfig.GetMemberOrgsPolicy()
	if collectionPolicyConfig == nil {
		return errors.New("Collection config policy is nil")
	}
	accessPolicyEnvelope := collectionPolicyConfig.GetSignaturePolicy()
	if accessPolicyEnvelope == nil {
		return errors.New("Collection config access policy is nil")
	}

	err := sc.setupAccessPolicy(collectionPolicyConfig, deserializer)
	if err != nil {
		return err
	}

	
	for _, principal := range accessPolicyEnvelope.Identities {
		switch principal.PrincipalClassification {
		case m.MSPPrincipal_ROLE:
			
			mspRole := &m.MSPRole{}
			err := proto.Unmarshal(principal.Principal, mspRole)
			if err != nil {
				return errors.Wrap(err, "Could not unmarshal MSPRole from principal")
			}
			sc.memberOrgs = append(sc.memberOrgs, mspRole.MspIdentifier)
		case m.MSPPrincipal_IDENTITY:
			principalId, err := deserializer.DeserializeIdentity(principal.Principal)
			if err != nil {
				return errors.Wrap(err, "Invalid identity principal, not a certificate")
			}
			sc.memberOrgs = append(sc.memberOrgs, principalId.GetMSPIdentifier())
		case m.MSPPrincipal_ORGANIZATION_UNIT:
			OU := &m.OrganizationUnit{}
			err := proto.Unmarshal(principal.Principal, OU)
			if err != nil {
				return errors.Wrap(err, "Could not unmarshal OrganizationUnit from principal")
			}
			sc.memberOrgs = append(sc.memberOrgs, OU.MspIdentifier)
		default:
			return errors.New(fmt.Sprintf("Invalid principal type %d", int32(principal.PrincipalClassification)))
		}
	}

	return nil
}



func (sc *SimpleCollection) setupAccessPolicy(collectionPolicyConfig *common.CollectionPolicyConfig, deserializer msp.IdentityDeserializer) error {
	var err error
	sc.accessPolicy, err = getPolicy(collectionPolicyConfig, deserializer)
	return err
}


func (s *SimpleCollectionPersistenceConfigs) BlockToLive() uint64 {
	return s.blockToLive
}
