/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"encoding/hex"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/pkg/errors"
)



type PvtDataCollections []*ledger.TxPvtData


func (pvt *PvtDataCollections) Marshal() ([][]byte, error) {
	pvtDataBytes := make([][]byte, 0)
	for index, each := range *pvt {
		if each == nil {
			errMsg := fmt.Sprintf("Mallformed private data payload, rwset index %d is nil", index)
			return nil, errors.New(errMsg)
		}
		pvtBytes, err := proto.Marshal(each.WriteSet)
		if err != nil {
			errMsg := fmt.Sprintf("Could not marshal private rwset index %d, due to %s", index, err)
			return nil, errors.New(errMsg)
		}
		
		txSeqInBlock := each.SeqInBlock
		pvtDataPayload := &gossip.PvtDataPayload{TxSeqInBlock: txSeqInBlock, Payload: pvtBytes}
		payloadBytes, err := proto.Marshal(pvtDataPayload)
		if err != nil {
			errMsg := fmt.Sprintf("Could not marshal private payload with transaction index %d, due to %s", txSeqInBlock, err)
			return nil, errors.New(errMsg)
		}

		pvtDataBytes = append(pvtDataBytes, payloadBytes)
	}
	return pvtDataBytes, nil
}



func (pvt *PvtDataCollections) Unmarshal(data [][]byte) error {
	for _, each := range data {
		payload := &gossip.PvtDataPayload{}
		if err := proto.Unmarshal(each, payload); err != nil {
			return err
		}
		pvtRWSet := &rwset.TxPvtReadWriteSet{}
		if err := proto.Unmarshal(payload.Payload, pvtRWSet); err != nil {
			return err
		}
		*pvt = append(*pvt, &ledger.TxPvtData{
			SeqInBlock: payload.TxSeqInBlock,
			WriteSet:   pvtRWSet,
		})
	}

	return nil
}


func PrivateRWSets(rwsets ...PrivateRWSet) [][]byte {
	var res [][]byte
	for _, rws := range rwsets {
		res = append(res, []byte(rws))
	}
	return res
}


type PrivateRWSet []byte


func (rws PrivateRWSet) Digest() string {
	return hex.EncodeToString(util.ComputeSHA256(rws))
}



type PrivateRWSetWithConfig struct {
	RWSet            []PrivateRWSet
	CollectionConfig *common.CollectionConfig
}
