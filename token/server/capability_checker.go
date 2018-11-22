/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package server

import (
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/pkg/errors"
)




type CapabilityChecker interface {
	FabToken(channelId string) (bool, error)
}


type TokenCapabilityChecker struct {
	PeerOps peer.Operations
}

func (c *TokenCapabilityChecker) FabToken(channelId string) (bool, error) {
	ac, ok := c.PeerOps.GetChannelConfig(channelId).ApplicationConfig()
	if !ok {
		return false, errors.Errorf("no application config found for channel %s", channelId)
	}
	return ac.Capabilities().FabToken(), nil
}
