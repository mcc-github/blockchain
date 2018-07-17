/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	cb "github.com/mcc-github/blockchain/protos/common"
)

const (
	
	ConsortiumsGroupKey = "Consortiums"
)


type ConsortiumsConfig struct {
	consortiums map[string]Consortium
}


func NewConsortiumsConfig(consortiumsGroup *cb.ConfigGroup, mspConfig *MSPConfigHandler) (*ConsortiumsConfig, error) {
	cc := &ConsortiumsConfig{
		consortiums: make(map[string]Consortium),
	}

	for consortiumName, consortiumGroup := range consortiumsGroup.Groups {
		var err error
		if cc.consortiums[consortiumName], err = NewConsortiumConfig(consortiumGroup, mspConfig); err != nil {
			return nil, err
		}
	}
	return cc, nil
}


func (cc *ConsortiumsConfig) Consortiums() map[string]Consortium {
	return cc.consortiums
}
