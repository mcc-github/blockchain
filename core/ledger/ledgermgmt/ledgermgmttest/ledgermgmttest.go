/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgermgmttest

import (
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/mcc-github/blockchain/core/ledger/mock"
)




func NewInitializer(testLedgerDir string) *ledgermgmt.Initializer {
	return &ledgermgmt.Initializer{
		Config: &ledger.Config{
			RootFSPath: testLedgerDir,
			
			StateDBConfig: &ledger.StateDBConfig{},
			HistoryDBConfig: &ledger.HistoryDBConfig{
				Enabled: false,
			},
			PrivateDataConfig: &ledger.PrivateDataConfig{
				MaxBatchSize:    5000,
				BatchesInterval: 1000,
				PurgeInterval:   100,
			},
		},
		PlatformRegistry:              platforms.NewRegistry(&golang.Platform{}),
		MetricsProvider:               &disabled.Provider{},
		DeployedChaincodeInfoProvider: &mock.DeployedChaincodeInfoProvider{},
	}
}
