/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgerconfig

import (
	"github.com/spf13/viper"
)

const confTotalQueryLimit = "ledger.state.totalQueryLimit"
const confEnableHistoryDatabase = "ledger.history.enableHistoryDatabase"

var confCollElgProcMaxDbBatchSize = &conf{"ledger.pvtdataStore.collElgProcMaxDbBatchSize", 5000}
var confCollElgProcDbBatchesInterval = &conf{"ledger.pvtdataStore.collElgProcDbBatchesInterval", 1000}


func GetTotalQueryLimit() int {
	totalQueryLimit := viper.GetInt(confTotalQueryLimit)
	
	if !viper.IsSet(confTotalQueryLimit) {
		totalQueryLimit = 10000
	}
	return totalQueryLimit
}


func IsHistoryDBEnabled() bool {
	return viper.GetBool(confEnableHistoryDatabase)
}

type conf struct {
	Name       string
	DefaultVal int
}
