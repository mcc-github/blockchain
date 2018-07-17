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
)


type ChannelProvider struct {
	*registry
	v11 bool
}


func NewChannelProvider(capabilities map[string]*cb.Capability) *ChannelProvider {
	cp := &ChannelProvider{}
	cp.registry = newRegistry(cp, capabilities)
	_, cp.v11 = capabilities[ChannelV1_1]
	return cp
}


func (cp *ChannelProvider) Type() string {
	return channelTypeName
}


func (cp *ChannelProvider) HasCapability(capability string) bool {
	switch capability {
	
	case ChannelV1_1:
		return true
	default:
		return false
	}
}


func (cp *ChannelProvider) MSPVersion() msp.MSPVersion {
	switch {
	case cp.v11:
		return msp.MSPv1_1
	default:
		return msp.MSPv1_0
	}
}
