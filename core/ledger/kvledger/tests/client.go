/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)




type client struct {
	lgr            ledger.PeerLedger
	lgrID          string
	simulatedTrans []*txAndPvtdata 
	assert         *assert.Assertions
}

func newClient(lgr ledger.PeerLedger, lgrID string, t *testing.T) *client {
	return &client{lgr, lgrID, nil, assert.New(t)}
}





func (c *client) simulateDataTx(txid string, simulationLogic func(s *simulator)) *txAndPvtdata {
	if txid == "" {
		txid = util.GenerateUUID()
	}
	ledgerSimulator, err := c.lgr.NewTxSimulator(txid)
	c.assert.NoError(err)
	sim := &simulator{ledgerSimulator, txid, c.assert}
	simulationLogic(sim)
	txAndPvtdata := sim.done()
	c.simulatedTrans = append(c.simulatedTrans, txAndPvtdata)
	return txAndPvtdata
}

func (c *client) addPostOrderTx(txid string, customTxType common.HeaderType) *txAndPvtdata {
	if txid == "" {
		txid = util.GenerateUUID()
	}
	channelHeader := protoutil.MakeChannelHeader(customTxType, 0, c.lgrID, 0)
	channelHeader.TxId = txid
	paylBytes := protoutil.MarshalOrPanic(
		&common.Payload{
			Header: protoutil.MakePayloadHeader(channelHeader, &common.SignatureHeader{}),
			Data:   nil,
		},
	)
	env := &common.Envelope{
		Payload:   paylBytes,
		Signature: nil,
	}
	txAndPvtdata := &txAndPvtdata{Txid: txid, Envelope: env}
	c.simulatedTrans = append(c.simulatedTrans, txAndPvtdata)
	return txAndPvtdata
}



func (c *client) simulateDeployTx(ccName string, collConfs []*collConf) *txAndPvtdata {
	ccData := &ccprovider.ChaincodeData{Name: ccName}
	ccDataBytes, err := proto.Marshal(ccData)
	c.assert.NoError(err)

	psudoLSCCInvokeFunc := func(s *simulator) {
		s.setState("lscc", ccName, string(ccDataBytes))
		if collConfs != nil {
			protoBytes, err := convertToCollConfigProtoBytes(collConfs)
			c.assert.NoError(err)
			s.setState("lscc", privdata.BuildCollectionKVSKey(ccName), string(protoBytes))
		}
	}
	return c.simulateDataTx("", psudoLSCCInvokeFunc)
}


func (c *client) simulateUpgradeTx(ccName string, collConfs []*collConf) *txAndPvtdata {
	return c.simulateDeployTx(ccName, collConfs)
}


type simulator struct {
	ledger.TxSimulator
	txid   string
	assert *assert.Assertions
}

func (s *simulator) getState(ns, key string) string {
	val, err := s.GetState(ns, key)
	s.assert.NoError(err)
	return string(val)
}

func (s *simulator) setState(ns, key string, val string) {
	s.assert.NoError(
		s.SetState(ns, key, []byte(val)),
	)
}

func (s *simulator) setPvtdata(ns, coll, key string, val string) {
	s.assert.NoError(
		s.SetPrivateData(ns, coll, key, []byte(val)),
	)
}

func (s *simulator) done() *txAndPvtdata {
	s.Done()
	simRes, err := s.GetTxSimulationResults()
	s.assert.NoError(err)
	pubRwsetBytes, err := simRes.GetPubSimulationBytes()
	s.assert.NoError(err)
	envelope, err := constructTransaction(s.txid, pubRwsetBytes)
	s.assert.NoError(err)
	txAndPvtdata := &txAndPvtdata{Txid: s.txid, Envelope: envelope, Pvtws: simRes.PvtSimulationResults}
	return txAndPvtdata
}
