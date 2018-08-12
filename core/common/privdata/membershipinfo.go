/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"github.com/mcc-github/blockchain/protos/common"
)


type MembershipInfoProvider interface {
	
	AmMemberOf(collectionPolicyConfig *common.CollectionPolicyConfig) (bool, error)
}

type membershipProvider struct {
	selfSignedData common.SignedData
	cf             CollectionFilter
	channelName    string
}

func NewMembershipInfoProvider(channelName string, selfSignedData common.SignedData, filter CollectionFilter) MembershipInfoProvider {
	return &membershipProvider{channelName: channelName, selfSignedData: selfSignedData, cf: filter}
}

func (m *membershipProvider) AmMemberOf(collectionPolicyConfig *common.CollectionPolicyConfig) (bool, error) {
	filt, err := m.cf.AccessFilter(m.channelName, collectionPolicyConfig)
	if err != nil {
		return false, err
	}
	return filt(m.selfSignedData), nil
}
