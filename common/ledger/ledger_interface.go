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

package ledger

import (
	"github.com/mcc-github/blockchain-protos-go/common"
)


type Ledger interface {
	
	GetBlockchainInfo() (*common.BlockchainInfo, error)
	
	
	GetBlockByNumber(blockNumber uint64) (*common.Block, error)
	
	
	
	GetBlocksIterator(startBlockNumber uint64) (ResultsIterator, error)
	
	Close()
}


type ResultsIterator interface {
	
	
	Next() (QueryResult, error)
	
	Close()
}


type QueryResultsIterator interface {
	ResultsIterator
	GetBookmarkAndClose() string
}


type QueryResult interface{}


type PrunePolicy interface{}
