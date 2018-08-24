/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	configtxtest "github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/example"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/mcc-github/blockchain/core/ledger/testutil"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/common"
	logging "github.com/op/go-logging"
)

const (
	ledgerID = "Default"
)

var logger = logging.MustGetLogger("main")

var peerLedger ledger.PeerLedger
var marbleApp *example.MarbleApp
var committer *example.Committer
var consenter *example.Consenter

func init() {

	
	
	
	logger.Debugf("Marble Example main init()")

	
	testutil.SetupCoreYAMLConfig()

	cleanup()
	ledgermgmt.Initialize(&ledgermgmt.Initializer{
		PlatformRegistry: platforms.NewRegistry(&golang.Platform{}),
	})
	var err error
	gb, _ := configtxtest.MakeGenesisBlock(ledgerID)
	peerLedger, err = ledgermgmt.CreateLedger(gb)

	if err != nil {
		panic(fmt.Errorf("Error in NewKVLedger(): %s", err))
	}
	marbleApp = example.ConstructMarbleAppInstance(peerLedger)
	committer = example.ConstructCommitter(peerLedger)
	consenter = example.ConstructConsenter()
}

func main() {
	defer ledgermgmt.Close()
	
	
	
	

	initApp()
	transferMarble()

}

func initApp() {
	logger.Debugf("Marble Example initApp() to create a marble")
	marble := []string{"marble1", "blue", "35", "tom"}
	tx, err := marbleApp.CreateMarble(marble)
	handleError(err, true)
	rawBlock := consenter.ConstructBlock(tx)
	err = committer.Commit(rawBlock)
	handleError(err, true)
	printBlocksInfo(rawBlock)
}

func transferMarble() {
	logger.Debugf("Marble Example transferMarble()")
	tx1, err := marbleApp.TransferMarble([]string{"marble1", "jerry"})
	handleError(err, true)
	rawBlock := consenter.ConstructBlock(tx1)
	err = committer.Commit(rawBlock)
	handleError(err, true)
	printBlocksInfo(rawBlock)
}

func printBlocksInfo(block *common.Block) {
	
	txsFltr := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	numOfInvalid := 0
	
	for i := 0; i < len(block.Data.Data); i++ {
		if txsFltr.IsInvalid(i) {
			numOfInvalid++
		}
	}
	fmt.Printf("Num txs in rawBlock = [%d], num invalidTxs = [%d]\n",
		len(block.Data.Data), numOfInvalid)
}

func handleError(err error, quit bool) {
	if err != nil {
		if quit {
			panic(fmt.Errorf("Error: %s\n", err))
		} else {
			fmt.Printf("Error: %s\n", err)
		}
	}
}

func cleanup() {
	ledgerRootPath := ledgerconfig.GetRootPath()
	os.RemoveAll(ledgerRootPath)
}
