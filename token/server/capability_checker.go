/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package server

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/pkg/errors"
)




type CapabilityChecker interface {
	FabToken(channelId string) (bool, error)
}




type ChannelConfigGetter interface {
	GetChannelConfig(cid string) channelconfig.Resources
}


type TokenCapabilityChecker struct {
	ChannelConfigGetter ChannelConfigGetter
}

func (c *TokenCapabilityChecker) FabToken(channelId string) (bool, error) {
	channelConfig := c.ChannelConfigGetter.GetChannelConfig(channelId)
	if channelConfig == nil {
		
		return false, errors.Errorf("no channel config found for channel %s", channelId)
	}

	ac, ok := channelConfig.ApplicationConfig()
	if !ok {
		return false, errors.Errorf("no application config found for channel %s", channelId)
	}
	return ac.Capabilities().FabToken(), nil
}
