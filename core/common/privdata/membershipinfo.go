/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"github.com/mcc-github/blockchain/protos/common"
)


type MembershipProvider struct {
	selfSignedData common.SignedData
	cf             CollectionFilter
}


func NewMembershipInfoProvider(selfSignedData common.SignedData, filter CollectionFilter) *MembershipProvider {
	return &MembershipProvider{selfSignedData: selfSignedData, cf: filter}
}


func (m *MembershipProvider) AmMemberOf(channelName string, collectionPolicyConfig *common.CollectionPolicyConfig) (bool, error) {
	filt, err := m.cf.AccessFilter(channelName, collectionPolicyConfig)
	if err != nil {
		return false, err
	}
	return filt(m.selfSignedData), nil
}
