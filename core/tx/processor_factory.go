/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tx

import (
	"github.com/mcc-github/blockchain/pkg/tx"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
)


type ProcessorFactory struct {
	ProcessorCreators map[common.HeaderType]tx.ProcessorCreator
}






func (f *ProcessorFactory) CreateProcessor(txEnvelopeBytes []byte) (processor tx.Processor, simulatedRWSet [][]byte, err error) {
	txEnv, err := validateProtoAndConstructTxEnv(txEnvelopeBytes)
	if err != nil {
		return nil, nil, err
	}
	c, ok := f.ProcessorCreators[common.HeaderType(txEnv.Payload.Header.ChannelHeader.Type)]
	if !ok {
		return nil, nil, &tx.InvalidErr{
			ValidationCode: peer.TxValidationCode_UNKNOWN_TX_TYPE,
		}
	}
	return c.NewProcessor(txEnv)
}



func validateProtoAndConstructTxEnv(txEnvelopeBytes []byte) (*tx.Envelope, error) {
	
	return &tx.Envelope{}, nil
}
