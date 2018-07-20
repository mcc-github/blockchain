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

var logger = logging.MustGetLogger("main")

const (
	ledgerID = "Default"
)

var peerLedger ledger.PeerLedger
var app *example.App
var committer *example.Committer
var consenter *example.Consenter

var accounts = []string{"account1", "account2", "account3", "account4"}

func init() {

	
	testutil.SetupCoreYAMLConfig()

	
	
	
	cleanup()
	ledgermgmt.Initialize(nil, platforms.NewRegistry(&golang.Platform{}))
	var err error

	gb, _ := configtxtest.MakeGenesisBlock(ledgerID)
	peerLedger, err = ledgermgmt.CreateLedger(gb)
	if err != nil {
		panic(fmt.Errorf("Error in NewKVLedger(): %s", err))
	}
	app = example.ConstructAppInstance(peerLedger)
	committer = example.ConstructCommitter(peerLedger)
	consenter = example.ConstructConsenter()
}

func main() {
	defer ledgermgmt.Close()

	
	
	
	

	
	initApp()

	printBalances()

	
	transferFunds()

	printBalances()

	
	
	tryInvalidTransfer()

	
	
	
	tryDoubleSpend()

	printBalances()
}

func initApp() {
	logger.Debug("Entering initApp()")
	tx, err := app.Init(map[string]int{
		accounts[0]: 100,
		accounts[1]: 100,
		accounts[2]: 100,
		accounts[3]: 100})
	handleError(err, true)
	rawBlock := consenter.ConstructBlock(tx)
	err = committer.Commit(rawBlock)
	handleError(err, true)
	printBlocksInfo(rawBlock)
	logger.Debug("Exiting initApp()")
}

func transferFunds() {
	logger.Debug("Entering transferFunds()")
	tx1, err := app.TransferFunds("account1", "account2", 50)
	handleError(err, true)
	tx2, err := app.TransferFunds("account3", "account4", 50)
	handleError(err, true)

	
	rawBlock := consenter.ConstructBlock(tx1, tx2)

	
	err = committer.Commit(rawBlock)
	handleError(err, true)
	printBlocksInfo(rawBlock)
	logger.Debug("Exiting transferFunds")
}

func tryInvalidTransfer() {
	logger.Debug("Entering tryInvalidTransfer()")
	_, err := app.TransferFunds("account1", "account2", 60)
	handleError(err, false)
	logger.Debug("Exiting tryInvalidTransfer()")
}

func tryDoubleSpend() {
	logger.Debug("Entering tryDoubleSpend()")
	tx1, err := app.TransferFunds("account1", "account2", 50)
	handleError(err, true)
	tx2, err := app.TransferFunds("account1", "account4", 50)
	handleError(err, true)
	rawBlock := consenter.ConstructBlock(tx1, tx2)
	err = committer.Commit(rawBlock)
	handleError(err, true)
	printBlocksInfo(rawBlock)
	logger.Debug("Exiting tryDoubleSpend()")
}

func printBlocksInfo(block *common.Block) {
	logger.Debug("Entering printBlocksInfo()")
	
	txsFltr := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	numOfInvalid := 0
	
	for i := 0; i < len(block.Data.Data); i++ {
		if txsFltr.IsInvalid(i) {
			numOfInvalid++
		}
	}
	fmt.Printf("Num txs in rawBlock = [%d], num invalidTxs = [%d]\n",
		len(block.Data.Data), numOfInvalid)
	logger.Debug("Exiting printBlocksInfo()")
}

func printBalances() {
	logger.Debug("Entering printBalances()")
	balances, err := app.QueryBalances(accounts)
	handleError(err, true)
	for i := 0; i < len(accounts); i++ {
		fmt.Printf("[%s] = [%d]\n", accounts[i], balances[i])
	}
	logger.Debug("Exiting printBalances()")
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
