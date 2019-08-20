/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import (
	"fmt"
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/msp"
)


type stateEP struct {
	orgs map[string]msp.MSPRole_MSPRoleType
}



func NewStateEP(policy []byte) (KeyEndorsementPolicy, error) {
	s := &stateEP{orgs: make(map[string]msp.MSPRole_MSPRoleType)}
	if policy != nil {
		spe := &common.SignaturePolicyEnvelope{}
		if err := proto.Unmarshal(policy, spe); err != nil {
			return nil, fmt.Errorf("Error unmarshaling to SignaturePolicy: %s", err)
		}

		err := s.setMSPIDsFromSP(spe)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}


func (s *stateEP) Policy() ([]byte, error) {
	spe, err := s.policyFromMSPIDs()
	if err != nil {
		return nil, err
	}
	spBytes, err := proto.Marshal(spe)
	if err != nil {
		return nil, err
	}
	return spBytes, nil
}


func (s *stateEP) AddOrgs(role RoleType, neworgs ...string) error {
	var mspRole msp.MSPRole_MSPRoleType
	switch role {
	case RoleTypeMember:
		mspRole = msp.MSPRole_MEMBER
	case RoleTypePeer:
		mspRole = msp.MSPRole_PEER
	default:
		return &RoleTypeDoesNotExistError{RoleType: role}
	}

	
	for _, addorg := range neworgs {
		s.orgs[addorg] = mspRole
	}

	return nil
}


func (s *stateEP) DelOrgs(delorgs ...string) {
	for _, delorg := range delorgs {
		delete(s.orgs, delorg)
	}
}


func (s *stateEP) ListOrgs() []string {
	orgNames := make([]string, 0, len(s.orgs))
	for mspid := range s.orgs {
		orgNames = append(orgNames, mspid)
	}
	return orgNames
}

func (s *stateEP) setMSPIDsFromSP(sp *common.SignaturePolicyEnvelope) error {
	
	for _, identity := range sp.Identities {
		
		if identity.PrincipalClassification == msp.MSPPrincipal_ROLE {
			msprole := &msp.MSPRole{}
			err := proto.Unmarshal(identity.Principal, msprole)
			if err != nil {
				return fmt.Errorf("error unmarshaling msp principal: %s", err)
			}
			s.orgs[msprole.GetMspIdentifier()] = msprole.GetRole()
		}
	}
	return nil
}

func (s *stateEP) policyFromMSPIDs() (*common.SignaturePolicyEnvelope, error) {
	mspids := s.ListOrgs()
	sort.Strings(mspids)
	principals := make([]*msp.MSPPrincipal, len(mspids))
	sigspolicy := make([]*common.SignaturePolicy, len(mspids))
	for i, id := range mspids {
		principal, err := proto.Marshal(
			&msp.MSPRole{
				Role:          s.orgs[id],
				MspIdentifier: id,
			},
		)
		if err != nil {
			return nil, err
		}
		principals[i] = &msp.MSPPrincipal{
			PrincipalClassification: msp.MSPPrincipal_ROLE,
			Principal:               principal,
		}
		sigspolicy[i] = &common.SignaturePolicy{
			Type: &common.SignaturePolicy_SignedBy{
				SignedBy: int32(i),
			},
		}
	}

	
	p := &common.SignaturePolicyEnvelope{
		Version: 0,
		Rule: &common.SignaturePolicy{
			Type: &common.SignaturePolicy_NOutOf_{
				NOutOf: &common.SignaturePolicy_NOutOf{
					N:     int32(len(mspids)),
					Rules: sigspolicy,
				},
			},
		},
		Identities: principals,
	}
	return p, nil
}
