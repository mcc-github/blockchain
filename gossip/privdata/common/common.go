/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/gossip"
)






type DigKey struct {
	TxId       string
	Namespace  string
	Collection string
	BlockSeq   uint64
	SeqInBlock uint64
}

type Dig2CollectionConfig map[DigKey]*common.StaticCollectionConfig



type FetchedPvtDataContainer struct {
	AvailableElements []*gossip.PvtDataElement
	PurgedElements    []*gossip.PvtDataDigest
}
