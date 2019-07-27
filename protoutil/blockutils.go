/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"

	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/pkg/errors"
)


func NewBlock(seqNum uint64, previousHash []byte) *cb.Block {
	block := &cb.Block{}
	block.Header = &cb.BlockHeader{}
	block.Header.Number = seqNum
	block.Header.PreviousHash = previousHash
	block.Data = &cb.BlockData{}

	var metadataContents [][]byte
	for i := 0; i < len(cb.BlockMetadataIndex_name); i++ {
		metadataContents = append(metadataContents, []byte{})
	}
	block.Metadata = &cb.BlockMetadata{Metadata: metadataContents}

	return block
}

type asn1Header struct {
	Number       *big.Int
	PreviousHash []byte
	DataHash     []byte
}

func BlockHeaderBytes(b *cb.BlockHeader) []byte {
	asn1Header := asn1Header{
		PreviousHash: b.PreviousHash,
		DataHash:     b.DataHash,
		Number:       new(big.Int).SetUint64(b.Number),
	}
	result, err := asn1.Marshal(asn1Header)
	if err != nil {
		
		
		
		panic(err)
	}
	return result
}

func BlockHeaderHash(b *cb.BlockHeader) []byte {
	sum := sha256.Sum256(BlockHeaderBytes(b))
	return sum[:]
}

func BlockDataHash(b *cb.BlockData) []byte {
	sum := sha256.Sum256(bytes.Join(b.Data, nil))
	return sum[:]
}



func GetChainIDFromBlockBytes(bytes []byte) (string, error) {
	block, err := GetBlockFromBlockBytes(bytes)
	if err != nil {
		return "", err
	}

	return GetChainIDFromBlock(block)
}


func GetChainIDFromBlock(block *cb.Block) (string, error) {
	if block == nil || block.Data == nil || block.Data.Data == nil || len(block.Data.Data) == 0 {
		return "", errors.Errorf("failed to retrieve channel id - block is empty")
	}
	var err error
	envelope, err := GetEnvelopeFromBlock(block.Data.Data[0])
	if err != nil {
		return "", err
	}
	payload, err := GetPayload(envelope)
	if err != nil {
		return "", err
	}

	if payload.Header == nil {
		return "", errors.Errorf("failed to retrieve channel id - payload header is empty")
	}
	chdr, err := UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return "", err
	}

	return chdr.ChannelId, nil
}


func GetMetadataFromBlock(block *cb.Block, index cb.BlockMetadataIndex) (*cb.Metadata, error) {
	md := &cb.Metadata{}
	err := proto.Unmarshal(block.Metadata.Metadata[index], md)
	if err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling metadata from block at index [%s]", index)
	}
	return md, nil
}



func GetMetadataFromBlockOrPanic(block *cb.Block, index cb.BlockMetadataIndex) *cb.Metadata {
	md, err := GetMetadataFromBlock(block, index)
	if err != nil {
		panic(err)
	}
	return md
}




func GetConsenterMetadataFromBlock(block *cb.Block) (*cb.Metadata, error) {
	m, err := GetMetadataFromBlock(block, cb.BlockMetadataIndex_SIGNATURES)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve metadata at index: %d", cb.BlockMetadataIndex_SIGNATURES)
	}

	
	if len(m.Value) == 0 {
		return GetMetadataFromBlock(block, cb.BlockMetadataIndex_ORDERER)
	}

	obm := &cb.OrdererBlockMetadata{}
	err = proto.Unmarshal(m.Value, obm)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal orderer block metadata")
	}

	res := &cb.Metadata{}
	err = proto.Unmarshal(obm.ConsenterMetadata, res)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal consenter metadata")
	}

	return res, nil
}



func GetLastConfigIndexFromBlock(block *cb.Block) (uint64, error) {
	md, err := GetMetadataFromBlock(block, cb.BlockMetadataIndex_LAST_CONFIG)
	if err != nil {
		return 0, err
	}
	lc := &cb.LastConfig{}
	err = proto.Unmarshal(md.Value, lc)
	if err != nil {
		return 0, errors.Wrap(err, "error unmarshaling LastConfig")
	}
	return lc.Index, nil
}



func GetLastConfigIndexFromBlockOrPanic(block *cb.Block) uint64 {
	index, err := GetLastConfigIndexFromBlock(block)
	if err != nil {
		panic(err)
	}
	return index
}


func GetBlockFromBlockBytes(blockBytes []byte) (*cb.Block, error) {
	block := &cb.Block{}
	err := proto.Unmarshal(blockBytes, block)
	if err != nil {
		return block, errors.Wrap(err, "error unmarshaling block")
	}
	return block, nil
}


func CopyBlockMetadata(src *cb.Block, dst *cb.Block) {
	dst.Metadata = src.Metadata
	
	
	InitBlockMetadata(dst)
}


func InitBlockMetadata(block *cb.Block) {
	if block.Metadata == nil {
		block.Metadata = &cb.BlockMetadata{Metadata: [][]byte{{}, {}, {}}}
	} else if len(block.Metadata.Metadata) < int(cb.BlockMetadataIndex_TRANSACTIONS_FILTER+1) {
		for i := int(len(block.Metadata.Metadata)); i <= int(cb.BlockMetadataIndex_TRANSACTIONS_FILTER); i++ {
			block.Metadata.Metadata = append(block.Metadata.Metadata, []byte{})
		}
	}
}
