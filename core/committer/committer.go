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
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
)








type Committer interface {

	
	CommitWithPvtData(blockAndPvtData *ledger.BlockAndPvtData) error

	
	
	GetPvtDataAndBlockByNum(seqNum uint64) (*ledger.BlockAndPvtData, error)

	
	
	
	GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error)

	
	LedgerHeight() (uint64, error)

	
	GetBlocks(blockSeqs []uint64) []*common.Block

	
	GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error)

	
	
	
	
	CommitPvtData(blockPvtData []*ledger.BlockPvtData) ([]*ledger.PvtdataHashMismatch, error)

	
	GetMissingPvtDataTracker() (ledger.MissingPvtDataTracker, error)

	
	Close()
}
