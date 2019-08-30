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

package committer

import (
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/core/ledger"
)








type Committer interface {

	
	CommitLegacy(blockAndPvtData *ledger.BlockAndPvtData, commitOpts *ledger.CommitOptions) error

	
	
	GetPvtDataAndBlockByNum(seqNum uint64) (*ledger.BlockAndPvtData, error)

	
	
	
	GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error)

	
	LedgerHeight() (uint64, error)

	
	
	DoesPvtDataInfoExistInLedger(blockNum uint64) (bool, error)

	
	GetBlocks(blockSeqs []uint64) []*common.Block

	
	GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error)

	
	
	
	
	CommitPvtDataOfOldBlocks(blockPvtData []*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error)

	
	GetMissingPvtDataTracker() (ledger.MissingPvtDataTracker, error)

	
	Close()
}
