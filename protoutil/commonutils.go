/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoutil

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)



func MarshalOrPanic(pb proto.Message) []byte {
	data, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	return data
}


func Marshal(pb proto.Message) ([]byte, error) {
	return proto.Marshal(pb)
}



func CreateNonceOrPanic() []byte {
	nonce, err := CreateNonce()
	if err != nil {
		panic(err)
	}
	return nonce
}


func CreateNonce() ([]byte, error) {
	nonce, err := crypto.GetRandomNonce()
	return nonce, errors.WithMessage(err, "error generating random nonce")
}



func UnmarshalPayloadOrPanic(encoded []byte) *cb.Payload {
	payload, err := UnmarshalPayload(encoded)
	if err != nil {
		panic(err)
	}
	return payload
}


func UnmarshalPayload(encoded []byte) (*cb.Payload, error) {
	payload := &cb.Payload{}
	err := proto.Unmarshal(encoded, payload)
	return payload, errors.Wrap(err, "error unmarshaling Payload")
}



func UnmarshalEnvelopeOrPanic(encoded []byte) *cb.Envelope {
	envelope, err := UnmarshalEnvelope(encoded)
	if err != nil {
		panic(err)
	}
	return envelope
}


func UnmarshalEnvelope(encoded []byte) (*cb.Envelope, error) {
	envelope := &cb.Envelope{}
	err := proto.Unmarshal(encoded, envelope)
	return envelope, errors.Wrap(err, "error unmarshaling Envelope")
}



func UnmarshalBlockOrPanic(encoded []byte) *cb.Block {
	block, err := UnmarshalBlock(encoded)
	if err != nil {
		panic(err)
	}
	return block
}


func UnmarshalBlock(encoded []byte) (*cb.Block, error) {
	block := &cb.Block{}
	err := proto.Unmarshal(encoded, block)
	return block, errors.Wrap(err, "error unmarshaling Block")
}



func UnmarshalEnvelopeOfType(envelope *cb.Envelope, headerType cb.HeaderType, message proto.Message) (*cb.ChannelHeader, error) {
	payload, err := UnmarshalPayload(envelope.Payload)
	if err != nil {
		return nil, err
	}

	if payload.Header == nil {
		return nil, errors.New("envelope must have a Header")
	}

	chdr, err := UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, err
	}

	if chdr.Type != int32(headerType) {
		return nil, errors.Errorf("invalid type %s, expected %s", cb.HeaderType(chdr.Type), headerType)
	}

	err = proto.Unmarshal(payload.Data, message)
	err = errors.Wrapf(err, "error unmarshaling message for type %s", headerType)
	return chdr, err
}



func ExtractEnvelopeOrPanic(block *cb.Block, index int) *cb.Envelope {
	envelope, err := ExtractEnvelope(block, index)
	if err != nil {
		panic(err)
	}
	return envelope
}



func ExtractEnvelope(block *cb.Block, index int) (*cb.Envelope, error) {
	if block.Data == nil {
		return nil, errors.New("block data is nil")
	}

	envelopeCount := len(block.Data.Data)
	if index < 0 || index >= envelopeCount {
		return nil, errors.New("envelope index out of bounds")
	}
	marshaledEnvelope := block.Data.Data[index]
	envelope, err := GetEnvelopeFromBlock(marshaledEnvelope)
	err = errors.WithMessagef(err, "block data does not carry an envelope at index %d", index)
	return envelope, err
}



func ExtractPayloadOrPanic(envelope *cb.Envelope) *cb.Payload {
	payload, err := ExtractPayload(envelope)
	if err != nil {
		panic(err)
	}
	return payload
}


func ExtractPayload(envelope *cb.Envelope) (*cb.Payload, error) {
	payload := &cb.Payload{}
	err := proto.Unmarshal(envelope.Payload, payload)
	err = errors.Wrap(err, "no payload in envelope")
	return payload, err
}


func MakeChannelHeader(headerType cb.HeaderType, version int32, chainID string, epoch uint64) *cb.ChannelHeader {
	return &cb.ChannelHeader{
		Type:    int32(headerType),
		Version: version,
		Timestamp: &timestamp.Timestamp{
			Seconds: time.Now().Unix(),
			Nanos:   0,
		},
		ChannelId: chainID,
		Epoch:     epoch,
	}
}


func MakeSignatureHeader(serializedCreatorCertChain []byte, nonce []byte) *cb.SignatureHeader {
	return &cb.SignatureHeader{
		Creator: serializedCreatorCertChain,
		Nonce:   nonce,
	}
}



func SetTxID(channelHeader *cb.ChannelHeader, signatureHeader *cb.SignatureHeader) error {
	txid, err := ComputeTxID(
		signatureHeader.Nonce,
		signatureHeader.Creator,
	)
	if err != nil {
		return err
	}
	channelHeader.TxId = txid
	return nil
}


func MakePayloadHeader(ch *cb.ChannelHeader, sh *cb.SignatureHeader) *cb.Header {
	return &cb.Header{
		ChannelHeader:   MarshalOrPanic(ch),
		SignatureHeader: MarshalOrPanic(sh),
	}
}


func NewSignatureHeader(id identity.Serializer) (*cb.SignatureHeader, error) {
	creator, err := id.Serialize()
	if err != nil {
		return nil, err
	}
	nonce, err := CreateNonce()
	if err != nil {
		return nil, err
	}

	return &cb.SignatureHeader{
		Creator: creator,
		Nonce:   nonce,
	}, nil
}


func NewSignatureHeaderOrPanic(id identity.Serializer) *cb.SignatureHeader {
	if id == nil {
		panic(errors.New("invalid signer. cannot be nil"))
	}

	signatureHeader, err := NewSignatureHeader(id)
	if err != nil {
		panic(fmt.Errorf("failed generating a new SignatureHeader: %s", err))
	}
	return signatureHeader
}


func SignOrPanic(signer identity.Signer, msg []byte) []byte {
	if signer == nil {
		panic(errors.New("invalid signer. cannot be nil"))
	}

	sigma, err := signer.Sign(msg)
	if err != nil {
		panic(fmt.Errorf("failed generating signature: %s", err))
	}
	return sigma
}


func UnmarshalChannelHeader(bytes []byte) (*cb.ChannelHeader, error) {
	chdr := &cb.ChannelHeader{}
	err := proto.Unmarshal(bytes, chdr)
	return chdr, errors.Wrap(err, "error unmarshaling ChannelHeader")
}



func UnmarshalChannelHeaderOrPanic(bytes []byte) *cb.ChannelHeader {
	chdr, err := UnmarshalChannelHeader(bytes)
	if err != nil {
		panic(err)
	}
	return chdr
}


func UnmarshalChaincodeID(bytes []byte) (*pb.ChaincodeID, error) {
	ccid := &pb.ChaincodeID{}
	err := proto.Unmarshal(bytes, ccid)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeID")
	}

	return ccid, nil
}



func IsConfigBlock(block *cb.Block) bool {
	envelope, err := ExtractEnvelope(block, 0)
	if err != nil {
		return false
	}

	payload, err := GetPayload(envelope)
	if err != nil {
		return false
	}

	if payload.Header == nil {
		return false
	}

	hdr, err := UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return false
	}

	return cb.HeaderType(hdr.Type) == cb.HeaderType_CONFIG || cb.HeaderType(hdr.Type) == cb.HeaderType_ORDERER_TRANSACTION
}


func ChannelHeader(env *cb.Envelope) (*cb.ChannelHeader, error) {
	envPayload, err := UnmarshalPayload(env.Payload)
	if err != nil {
		return nil, err
	}

	if envPayload.Header == nil {
		return nil, errors.New("header not set")
	}

	if envPayload.Header.ChannelHeader == nil {
		return nil, errors.New("channel header not set")
	}

	chdr, err := UnmarshalChannelHeader(envPayload.Header.ChannelHeader)
	if err != nil {
		return nil, errors.WithMessage(err, "error unmarshaling channel header")
	}

	return chdr, nil
}


func ChannelID(env *cb.Envelope) (string, error) {
	chdr, err := ChannelHeader(env)
	if err != nil {
		return "", errors.WithMessage(err, "error retrieving channel header")
	}

	return chdr.ChannelId, nil
}



func EnvelopeToConfigUpdate(configtx *cb.Envelope) (*cb.ConfigUpdateEnvelope, error) {
	configUpdateEnv := &cb.ConfigUpdateEnvelope{}
	_, err := UnmarshalEnvelopeOfType(configtx, cb.HeaderType_CONFIG_UPDATE, configUpdateEnv)
	if err != nil {
		return nil, err
	}
	return configUpdateEnv, nil
}
