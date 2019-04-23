/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	"github.com/mcc-github/blockchain/msp"
	cb "github.com/mcc-github/blockchain/protos/common"
)

const (
	channelTypeName = "Channel"

	
	ChannelV1_1 = "V1_1"

	
	ChannelV1_3 = "V1_3"

	
	ChannelV1_4_2 = "V1_4_2"

	
	ChannelV2_0 = "V2_0"
)


type ChannelProvider struct {
	*registry
	v11  bool
	v13  bool
	v142 bool
	v20  bool
}


func NewChannelProvider(capabilities map[string]*cb.Capability) *ChannelProvider {
	cp := &ChannelProvider{}
	cp.registry = newRegistry(cp, capabilities)
	_, cp.v11 = capabilities[ChannelV1_1]
	_, cp.v13 = capabilities[ChannelV1_3]
	_, cp.v142 = capabilities[ChannelV1_4_2]
	_, cp.v20 = capabilities[ChannelV2_0]
	return cp
}


func (cp *ChannelProvider) Type() string {
	return channelTypeName
}


func (cp *ChannelProvider) HasCapability(capability string) bool {
	switch capability {
	
	case ChannelV2_0:
		return true
	case ChannelV1_4_2:
		return true
	case ChannelV1_3:
		return true
	case ChannelV1_1:
		return true
	default:
		return false
	}
}


func (cp *ChannelProvider) MSPVersion() msp.MSPVersion {
	switch {
	case cp.v13 || cp.v142 || cp.v20:
		return msp.MSPv1_3
	case cp.v11:
		return msp.MSPv1_1
	default:
		return msp.MSPv1_0
	}
}


func (cp *ChannelProvider) ConsensusTypeMigration() bool {
	return cp.v142 || cp.v20
}


func (cp *ChannelProvider) OrgSpecificOrdererEndpoints() bool {
	return cp.v20
}
