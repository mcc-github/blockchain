/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package customtx

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/common"
)

var processors Processors
var once sync.Once


type Processors map[common.HeaderType]Processor


func Initialize(customTxProcessors Processors) {
	once.Do(func() {
		initialize(customTxProcessors)
	})
}

func initialize(customTxProcessors Processors) {
	processors = customTxProcessors
}


func GetProcessor(txType common.HeaderType) Processor {
	return processors[txType]
}
