/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
)


type MembershipProvider struct {
	selfSignedData              common.SignedData
	IdentityDeserializerFactory func(chainID string) msp.IdentityDeserializer
}


func NewMembershipInfoProvider(selfSignedData common.SignedData, identityDeserializerFunc func(chainID string) msp.IdentityDeserializer) *MembershipProvider {
	return &MembershipProvider{selfSignedData: selfSignedData, IdentityDeserializerFactory: identityDeserializerFunc}
}


func (m *MembershipProvider) AmMemberOf(channelName string, collectionPolicyConfig *common.CollectionPolicyConfig) (bool, error) {
	deserializer := m.IdentityDeserializerFactory(channelName)
	accessPolicy, err := getPolicy(collectionPolicyConfig, deserializer)
	if err != nil {
		return false, err
	}
	if err := accessPolicy.Evaluate([]*common.SignedData{&m.selfSignedData}); err != nil {
		return false, nil
	}
	return true, nil
}
