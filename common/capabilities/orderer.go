/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	cb "github.com/mcc-github/blockchain/protos/common"
)

const (
	ordererTypeName = "Orderer"

	
	OrdererV1_1 = "V1_1"

	
	
	
	
	
	
	OrdererV2_0 = "V2_0"
)


type OrdererProvider struct {
	*registry
	v11BugFixes   bool
	kafka2RaftMig bool
}


func NewOrdererProvider(capabilities map[string]*cb.Capability) *OrdererProvider {
	cp := &OrdererProvider{}
	cp.registry = newRegistry(cp, capabilities)
	_, cp.v11BugFixes = capabilities[OrdererV1_1]
	_, cp.kafka2RaftMig = capabilities[OrdererV2_0]
	return cp
}


func (cp *OrdererProvider) Type() string {
	return ordererTypeName
}


func (cp *OrdererProvider) HasCapability(capability string) bool {
	switch capability {
	
	case OrdererV1_1:
		return true
	case OrdererV2_0:
		return true
	default:
		return false
	}
}



func (cp *OrdererProvider) PredictableChannelTemplate() bool {
	return cp.v11BugFixes
}



func (cp *OrdererProvider) Resubmission() bool {
	return cp.v11BugFixes
}



func (cp *OrdererProvider) ExpirationCheck() bool {
	return cp.v11BugFixes
}


func (cp *OrdererProvider) Kafka2RaftMigration() bool {
	return cp.kafka2RaftMig
}
