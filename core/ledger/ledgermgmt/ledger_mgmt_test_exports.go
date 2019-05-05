/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ledgermgmt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/mock"
)




func InitializeTestEnv(t *testing.T) (cleanup func()) {
	cleanup, err := InitializeTestEnvWithInitializer(nil)
	if err != nil {
		t.Fatalf("Failed to initialize test environment: %s", err)
	}
	return cleanup
}


func InitializeTestEnvWithInitializer(initializer *Initializer) (cleanup func(), err error) {
	return InitializeExistingTestEnvWithInitializer(initializer)
}




func InitializeExistingTestEnvWithInitializer(initializer *Initializer) (cleanup func(), err error) {
	if initializer == nil {
		initializer = &Initializer{}
	}
	if initializer.DeployedChaincodeInfoProvider == nil {
		initializer.DeployedChaincodeInfoProvider = &mock.DeployedChaincodeInfoProvider{}
	}
	if initializer.MetricsProvider == nil {
		initializer.MetricsProvider = &disabled.Provider{}
	}
	if initializer.PlatformRegistry == nil {
		initializer.PlatformRegistry = platforms.NewRegistry(&golang.Platform{})
	}
	if initializer.Config == nil {
		rootPath, err := ioutil.TempDir("", "ltestenv")
		if err != nil {
			return nil, err
		}
		initializer.Config = &ledger.Config{
			RootFSPath: rootPath,
			StateDB: &ledger.StateDB{
				LevelDBPath: filepath.Join(rootPath, "stateleveldb"),
			},
		}
	}
	if initializer.Config.PrivateData == nil {
		initializer.Config.PrivateData = &ledger.PrivateData{
			StorePath:       filepath.Join(initializer.Config.RootFSPath, "pvtdataStore"),
			MaxBatchSize:    5000,
			BatchesInterval: 1000,
			PurgeInterval:   100,
		}
	}
	if initializer.Config.HistoryDB == nil {
		initializer.Config.HistoryDB = &ledger.HistoryDB{
			Enabled: true,
		}
	}
	initialize(initializer)
	cleanup = func() {
		Close()
		os.RemoveAll(initializer.Config.RootFSPath)
	}
	return cleanup, nil
}
