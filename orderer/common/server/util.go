/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mcc-github/blockchain/common/ledger/blkstorage/fsblkstorage"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	"github.com/mcc-github/blockchain/common/ledger/blockledger/fileledger"
	"github.com/mcc-github/blockchain/common/ledger/blockledger/ramledger"
	"github.com/mcc-github/blockchain/common/metrics"
	config "github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/pkg/errors"
)

func createLedgerFactory(conf *config.TopLevel, metricsProvider metrics.Provider) (blockledger.Factory, string, error) {
	var lf blockledger.Factory
	var ld string
	var err error
	switch conf.General.LedgerType {
	case "file":
		ld = conf.FileLedger.Location
		if ld == "" {
			ld = createTempDir(conf.FileLedger.Prefix)
		}
		logger.Debug("Ledger dir:", ld)
		if lf, err = fileledger.New(ld, metricsProvider); err != nil {
			return nil, "", errors.WithMessage(err, "Error in opening ledger factory")
		}
		
		
		
		
		createSubDir(ld, fsblkstorage.ChainsDir)
	case "ram":
		fallthrough
	default:
		lf = ramledger.New(int(conf.RAMLedger.HistorySize))
	}
	return lf, ld, nil
}

func createTempDir(dirPrefix string) string {
	dirPath, err := ioutil.TempDir("", dirPrefix)
	if err != nil {
		logger.Panic("Error creating temp dir:", err)
	}
	return dirPath
}

func createSubDir(parentDirPath string, subDir string) (string, bool) {
	var created bool
	subDirPath := filepath.Join(parentDirPath, subDir)
	if _, err := os.Stat(subDirPath); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(subDirPath, 0755); err != nil {
				logger.Panic("Error creating sub dir:", err)
			}
			created = true
		}
	} else {
		logger.Debugf("Found %s sub-dir and using it", fsblkstorage.ChainsDir)
	}
	return subDirPath, created
}
