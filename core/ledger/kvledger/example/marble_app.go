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
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/mcc-github/blockchain/common/util"
	ledger "github.com/mcc-github/blockchain/core/ledger"

	"github.com/mcc-github/blockchain/protos/common"
	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("example")


type MarbleApp struct {
	name   string
	ledger ledger.PeerLedger
}


func ConstructMarbleAppInstance(ledger ledger.PeerLedger) *MarbleApp {
	return &MarbleApp{"marbles_app", ledger}
}

type Marble struct {
	Name  string `json:"asset_name"` 
	Color string `json:"color"`
	Size  int    `json:"size"`
	User  string `json:"owner"`
	Rev   string `json:"_rev"`
	Txid  string `json:"txid"`
}


func (marbleApp *MarbleApp) CreateMarble(args []string) (*common.Envelope, error) {
	
	
	logger.Debugf("Entering ----------CreateMarble()----------")
	marbleName := args[0]
	marbleJsonBytes, err := init_marble(args)
	if err != nil {
		return nil, err
	}

	var txSimulator ledger.TxSimulator
	if txSimulator, err = marbleApp.ledger.NewTxSimulator(util.GenerateUUID()); err != nil {
		return nil, err
	}
	defer txSimulator.Done()

	txSimulator.SetState(marbleApp.name, marbleName, marbleJsonBytes)

	var txSimulationResults *ledger.TxSimulationResults
	if txSimulationResults, err = txSimulator.GetTxSimulationResults(); err != nil {
		return nil, err
	}
	logger.Debugf("CreateMarble() simulation done, packaging into a transaction...")
	var pubSimBytes []byte
	if pubSimBytes, err = txSimulationResults.GetPubSimulationBytes(); err != nil {
		return nil, err
	}
	tx := constructTransaction(pubSimBytes)
	logger.Debugf("Exiting CreateMarble()")
	return tx, nil
}




func init_marble(args []string) ([]byte, error) {
	var err error

	
	
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	logger.Debugf("Entering init marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}

	size, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}

	color := strings.ToLower(args[1])
	user := strings.ToLower(args[3])

	tx := "tx000000000000001" 
	marbleJson := `{"txid": "` + tx + `",  "asset_name": "` + args[0] + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + user + `"}`
	marbleBytes := []byte(marbleJson)

	logger.Debugf("Exiting init marble")
	return marbleBytes, nil
}


func (marbleApp *MarbleApp) TransferMarble(args []string) (*common.Envelope, error) {
	
	
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	marbleName := args[0]
	marbleNewOwner := args[1]

	logger.Debugf("Entering ----------TransferMarble----------")
	var txSimulator ledger.TxSimulator
	var err error
	if txSimulator, err = marbleApp.ledger.NewTxSimulator(util.GenerateUUID()); err != nil {
		return nil, err
	}
	defer txSimulator.Done()

	marbleBytes, err := txSimulator.GetState(marbleApp.name, marbleName)
	logger.Debugf("marbleBytes is: %v", marbleBytes)
	if marbleBytes != nil {
		jsonString := string(marbleBytes[:])
		logger.Debugf("TransferMarble() Retrieved jsonString: \n   %s", jsonString)
	}

	theMarble := Marble{}
	json.Unmarshal(marbleBytes, &theMarble) 

	logger.Debugf(" theMarble after unmarshal: %v", theMarble)

	logger.Debugf(" Setting the owner to: %s", marbleNewOwner)
	theMarble.User = marbleNewOwner      
	theMarble.Txid = "tx000000000000002" 

	updatedMarbleBytes, _ := json.Marshal(theMarble)
	if updatedMarbleBytes != nil {
		updatedJsonString := string(updatedMarbleBytes[:])
		logger.Debugf("updatedJsonString:\n   %s", updatedJsonString)
	}
	err = txSimulator.SetState(marbleApp.name, marbleName, updatedMarbleBytes)
	if err != nil {
		return nil, err
	}

	var txSimulationResults *ledger.TxSimulationResults
	if txSimulationResults, err = txSimulator.GetTxSimulationResults(); err != nil {
		return nil, err
	}
	logger.Debugf("TransferMarble() simulation done, packaging into a transaction...")
	var pubSimBytes []byte
	if pubSimBytes, err = txSimulationResults.GetPubSimulationBytes(); err != nil {
		return nil, err
	}
	tx := constructTransaction(pubSimBytes)
	return tx, nil
}
