/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
)

var supportFactory SupportFactory


type SupportFactory interface {
	
	NewSupport() Support
}


type Support interface {
	
	
	GetApplicationConfig(cid string) (channelconfig.Application, bool)
}

type supportImpl struct {
	operations Operations
}

func (s *supportImpl) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	cc := s.operations.GetChannelConfig(cid)
	if cc == nil {
		return nil, false
	}

	return cc.ApplicationConfig()
}
