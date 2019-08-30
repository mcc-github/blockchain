/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package testutil

import (
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy/mock"
)



func SampleBTLPolicy(m map[[2]string]uint64) pvtdatapolicy.BTLPolicy {
	ccInfoRetriever := &mock.CollectionInfoProvider{}
	ccInfoRetriever.CollectionInfoStub = func(ccName, collName string) (*common.StaticCollectionConfig, error) {
		btl := m[[2]string{ccName, collName}]
		return &common.StaticCollectionConfig{BlockToLive: btl}, nil
	}
	return pvtdatapolicy.ConstructBTLPolicy(ccInfoRetriever)
}
