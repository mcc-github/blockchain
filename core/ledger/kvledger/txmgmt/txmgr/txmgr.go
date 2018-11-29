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

package txmgr

import (
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
)


type TxMgr interface {
	NewQueryExecutor(txid string) (ledger.QueryExecutor, error)
	NewTxSimulator(txid string) (ledger.TxSimulator, error)
	ValidateAndPrepare(blockAndPvtdata *ledger.BlockAndPvtData, doMVCCValidation bool) ([]*TxStatInfo, error)
	RemoveStaleAndCommitPvtDataOfOldBlocks(blocksPvtData map[uint64][]*ledger.TxPvtData) error
	GetLastSavepoint() (*version.Height, error)
	ShouldRecover(lastAvailableBlock uint64) (bool, uint64, error)
	CommitLostBlock(blockAndPvtdata *ledger.BlockAndPvtData) error
	Commit() error
	Rollback()
	Shutdown()
}


type TxStatInfo struct {
	ValidationCode peer.TxValidationCode
	TxType         common.HeaderType
	ChaincodeID    *peer.ChaincodeID
	NumCollections int
}


type ErrUnsupportedTransaction struct {
	Msg string
}

func (e *ErrUnsupportedTransaction) Error() string {
	return e.Msg
}




type ErrPvtdataNotAvailable struct {
	Msg string
}

func (e *ErrPvtdataNotAvailable) Error() string {
	return e.Msg
}
