/*
Copyright State Street Corp. 2018 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	pb "github.com/mcc-github/blockchain-protos-go/peer"
)


type aclsProvider struct {
	aclPolicyRefs map[string]string
}

func (ag *aclsProvider) PolicyRefForAPI(aclName string) string {
	return ag.aclPolicyRefs[aclName]
}


func newAPIsProvider(acls map[string]*pb.APIResource) *aclsProvider {
	aclPolicyRefs := make(map[string]string)

	for key, acl := range acls {
		
		
		if '/' != acl.PolicyRef[0] {
			aclPolicyRefs[key] = "/" + ChannelGroupKey + "/" + ApplicationGroupKey + "/" + acl.PolicyRef
		} else {
			aclPolicyRefs[key] = acl.PolicyRef
		}
	}

	return &aclsProvider{
		aclPolicyRefs: aclPolicyRefs,
	}
}
