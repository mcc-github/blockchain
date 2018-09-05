/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	cb "github.com/mcc-github/blockchain/protos/common"
)

const (
	applicationTypeName = "Application"

	
	ApplicationV1_1 = "V1_1"

	
	ApplicationV1_2 = "V1_2"

	
	ApplicationV1_3 = "V1_3"

	
	ApplicationPvtDataExperimental = "V1_1_PVTDATA_EXPERIMENTAL"

	
	ApplicationResourcesTreeExperimental = "V1_1_RESOURCETREE_EXPERIMENTAL"
)


type ApplicationProvider struct {
	*registry
	v11                    bool
	v12                    bool
	v13                    bool
	v11PvtDataExperimental bool
}


func NewApplicationProvider(capabilities map[string]*cb.Capability) *ApplicationProvider {
	ap := &ApplicationProvider{}
	ap.registry = newRegistry(ap, capabilities)
	_, ap.v11 = capabilities[ApplicationV1_1]
	_, ap.v12 = capabilities[ApplicationV1_2]
	_, ap.v13 = capabilities[ApplicationV1_3]
	_, ap.v11PvtDataExperimental = capabilities[ApplicationPvtDataExperimental]
	return ap
}


func (ap *ApplicationProvider) Type() string {
	return applicationTypeName
}


func (ap *ApplicationProvider) ACLs() bool {
	return ap.v12 || ap.v13
}



func (ap *ApplicationProvider) ForbidDuplicateTXIdInBlock() bool {
	return ap.v11 || ap.v12 || ap.v13
}




func (ap *ApplicationProvider) PrivateChannelData() bool {
	return ap.v11PvtDataExperimental || ap.v12 || ap.v13
}



func (ap ApplicationProvider) CollectionUpgrade() bool {
	return ap.v12 || ap.v13
}



func (ap *ApplicationProvider) V1_1Validation() bool {
	return ap.v11 || ap.v12 || ap.v13
}



func (ap *ApplicationProvider) V1_2Validation() bool {
	return ap.v12 || ap.v13
}



func (ap *ApplicationProvider) V1_3Validation() bool {
	return ap.v13
}




func (ap *ApplicationProvider) MetadataLifecycle() bool {
	return false
}



func (ap *ApplicationProvider) KeyLevelEndorsement() bool {
	return ap.v13
}


func (ap *ApplicationProvider) HasCapability(capability string) bool {
	switch capability {
	
	case ApplicationV1_1:
		return true
	case ApplicationV1_2:
		return true
	case ApplicationV1_3:
		return true
	case ApplicationPvtDataExperimental:
		return true
	case ApplicationResourcesTreeExperimental:
		return true
	default:
		return false
	}
}
