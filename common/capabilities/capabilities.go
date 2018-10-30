/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	"github.com/mcc-github/blockchain/common/flogging"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("common.capabilities")


type provider interface {
	
	HasCapability(capability string) bool

	
	Type() string
}



type registry struct {
	provider     provider
	capabilities map[string]*cb.Capability
}

func newRegistry(p provider, capabilities map[string]*cb.Capability) *registry {
	return &registry{
		provider:     p,
		capabilities: capabilities,
	}
}


func (r *registry) Supported() error {
	for capabilityName := range r.capabilities {
		if r.provider.HasCapability(capabilityName) {
			logger.Debugf("%s capability %s is supported and is enabled", r.provider.Type(), capabilityName)
			continue
		}

		return errors.Errorf("%s capability %s is required but not supported", r.provider.Type(), capabilityName)
	}
	return nil
}
