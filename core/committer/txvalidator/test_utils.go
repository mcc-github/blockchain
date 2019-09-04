/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txvalidator

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/core/ledger"
)



type ApplicationCapabilities interface {
	channelconfig.ApplicationCapabilities
}



type QueryExecutor interface {
	ledger.QueryExecutor
}
