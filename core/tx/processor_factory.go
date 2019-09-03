/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tx

import (
	"github.com/pkg/errors"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/pkg/tx"
	"github.com/mcc-github/blockchain/protoutil"
)


type ProcessorFactory struct {
	ProcessorCreators map[common.HeaderType]tx.ProcessorCreator
}






func (f *ProcessorFactory) CreateProcessor(txEnvelopeBytes []byte) (processor tx.Processor, simulatedRWSet [][]byte, err error) {
	txEnv, err := validateProtoAndConstructTxEnv(txEnvelopeBytes)
	if err != nil {
		return nil, nil, err
	}
	c, ok := f.ProcessorCreators[common.HeaderType(txEnv.ChannelHeader.Type)]
	if !ok {
		return nil, nil, &tx.InvalidErr{
			ActualErr:      errors.Errorf("invalid transaction type %d", txEnv.ChannelHeader.Type),
			ValidationCode: peer.TxValidationCode_UNKNOWN_TX_TYPE,
		}
	}
	return c.NewProcessor(txEnv)
}



func validateProtoAndConstructTxEnv(txEnvelopeBytes []byte) (*tx.Envelope, error) {
	txenv, err := protoutil.UnmarshalEnvelope(txEnvelopeBytes)
	if err != nil {
		return nil, &tx.InvalidErr{
			ActualErr:      err,
			ValidationCode: peer.TxValidationCode_INVALID_OTHER_REASON,
		}
	}

	if len(txenv.Payload) == 0 {
		return nil, &tx.InvalidErr{
			ActualErr:      errors.New("nil envelope payload"),
			ValidationCode: peer.TxValidationCode_BAD_PAYLOAD,
		}
	}

	payload, err := protoutil.UnmarshalPayload(txenv.Payload)
	if err != nil {
		return nil, &tx.InvalidErr{
			ActualErr:      err,
			ValidationCode: peer.TxValidationCode_BAD_PAYLOAD,
		}
	}

	if payload.Header == nil {
		return nil, &tx.InvalidErr{
			ActualErr:      errors.New("nil payload header"),
			ValidationCode: peer.TxValidationCode_BAD_PAYLOAD,
		}
	}

	if len(payload.Header.ChannelHeader) == 0 {
		return nil, &tx.InvalidErr{
			ActualErr:      errors.New("nil payload channel header"),
			ValidationCode: peer.TxValidationCode_BAD_PAYLOAD,
		}
	}

	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, &tx.InvalidErr{
			ActualErr:      err,
			ValidationCode: peer.TxValidationCode_BAD_PAYLOAD,
		}
	}

	if len(payload.Header.SignatureHeader) == 0 {
		return nil, &tx.InvalidErr{
			ActualErr:      errors.New("nil payload signature header"),
			ValidationCode: peer.TxValidationCode_BAD_PAYLOAD,
		}
	}

	shdr, err := protoutil.UnmarshalSignatureHeader(payload.Header.SignatureHeader)
	if err != nil {
		return nil, &tx.InvalidErr{
			ActualErr:      err,
			ValidationCode: peer.TxValidationCode_BAD_PAYLOAD,
		}
	}

	

	
	
	

	return &tx.Envelope{
			SignedBytes:          txenv.Payload,
			Signature:            txenv.Signature,
			Data:                 payload.Data,
			ChannelHeaderBytes:   payload.Header.ChannelHeader,
			SignatureHeaderBytes: payload.Header.SignatureHeader,
			ChannelHeader:        chdr,
			SignatureHeader:      shdr,
		},
		nil
}
