/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deliver_test

import (
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
)


type blockledgerReader interface {
	blockledger.Reader
}


type blockledgerIterator interface {
	blockledger.Iterator
}
