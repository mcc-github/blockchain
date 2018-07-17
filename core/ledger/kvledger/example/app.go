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

package example

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/ledger"

	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	ptestutils "github.com/mcc-github/blockchain/protos/testutils"
)


type App struct {
	name   string
	ledger ledger.PeerLedger
}


func ConstructAppInstance(ledger ledger.PeerLedger) *App {
	return &App{"PaymentApp", ledger}
}


func (app *App) Init(initialBalances map[string]int) (*common.Envelope, error) {
	var txSimulator ledger.TxSimulator
	var err error
	if txSimulator, err = app.ledger.NewTxSimulator(util.GenerateUUID()); err != nil {
		return nil, err
	}
	defer txSimulator.Done()
	for accountID, bal := range initialBalances {
		txSimulator.SetState(app.name, accountID, toBytes(bal))
	}
	var txSimulationResults *ledger.TxSimulationResults
	var pubSimBytes []byte
	if txSimulationResults, err = txSimulator.GetTxSimulationResults(); err != nil {
		return nil, err
	}
	if pubSimBytes, err = txSimulationResults.GetPubSimulationBytes(); err != nil {
		return nil, err
	}
	tx := constructTransaction(pubSimBytes)
	return tx, nil
}


func (app *App) TransferFunds(fromAccount string, toAccount string, transferAmt int) (*common.Envelope, error) {
	
	var txSimulator ledger.TxSimulator
	var err error
	if txSimulator, err = app.ledger.NewTxSimulator(util.GenerateUUID()); err != nil {
		return nil, err
	}
	defer txSimulator.Done()
	var balFromBytes []byte
	if balFromBytes, err = txSimulator.GetState(app.name, fromAccount); err != nil {
		return nil, err
	}
	balFrom := toInt(balFromBytes)
	if balFrom-transferAmt < 0 {
		return nil, fmt.Errorf("Not enough balance in account [%s]. Balance = [%d], transfer request = [%d]",
			fromAccount, balFrom, transferAmt)
	}

	var balToBytes []byte
	if balToBytes, err = txSimulator.GetState(app.name, toAccount); err != nil {
		return nil, err
	}
	balTo := toInt(balToBytes)
	txSimulator.SetState(app.name, fromAccount, toBytes(balFrom-transferAmt))
	txSimulator.SetState(app.name, toAccount, toBytes(balTo+transferAmt))
	var txSimulationResults *ledger.TxSimulationResults
	if txSimulationResults, err = txSimulator.GetTxSimulationResults(); err != nil {
		return nil, err
	}
	var pubSimBytes []byte
	if pubSimBytes, err = txSimulationResults.GetPubSimulationBytes(); err != nil {
		return nil, err
	}
	
	
	tx := constructTransaction(pubSimBytes)
	return tx, nil
}


func (app *App) QueryBalances(accounts []string) ([]int, error) {
	var queryExecutor ledger.QueryExecutor
	var err error
	if queryExecutor, err = app.ledger.NewQueryExecutor(); err != nil {
		return nil, err
	}
	defer queryExecutor.Done()
	balances := make([]int, len(accounts))
	for i := 0; i < len(accounts); i++ {
		var balBytes []byte
		if balBytes, err = queryExecutor.GetState(app.name, accounts[i]); err != nil {
			return nil, err
		}
		balances[i] = toInt(balBytes)
	}
	return balances, nil
}

func constructTransaction(simulationResults []byte) *common.Envelope {
	ccid := &pb.ChaincodeID{
		Name:    "foo",
		Version: "v1",
	}
	response := &pb.Response{Status: 200}
	txEnv, _, _ := ptestutils.ConstructSingedTxEnvWithDefaultSigner(util.GetTestChainID(), ccid, response, simulationResults, "", nil, nil)
	return txEnv
}

func toBytes(balance int) []byte {
	return proto.EncodeVarint(uint64(balance))
}

func toInt(balanceBytes []byte) int {
	v, _ := proto.DecodeVarint(balanceBytes)
	return int(v)
}
