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

package blockledger

import (
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
)


type Factory interface {
	
	
	GetOrCreate(chainID string) (ReadWriter, error)

	
	ChainIDs() []string

	
	Close()
}


type Iterator interface {
	
	
	Next() (*cb.Block, cb.Status)
	
	ReadyChan() <-chan struct{}
	
	Close()
}


type Reader interface {
	
	
	Iterator(startType *ab.SeekPosition) (Iterator, uint64)
	
	Height() uint64
}


type Writer interface {
	
	Append(block *cb.Block) error
}


type ReadWriter interface {
	Reader
	Writer
}