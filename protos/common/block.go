/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"encoding/asn1"
	"fmt"
	"math"

	"github.com/mcc-github/blockchain/common/util"
)


func NewBlock(seqNum uint64, previousHash []byte) *Block {
	block := &Block{}
	block.Header = &BlockHeader{}
	block.Header.Number = seqNum
	block.Header.PreviousHash = previousHash
	block.Data = &BlockData{}

	var metadataContents [][]byte
	for i := 0; i < len(BlockMetadataIndex_name); i++ {
		metadataContents = append(metadataContents, []byte{})
	}
	block.Metadata = &BlockMetadata{Metadata: metadataContents}

	return block
}

type asn1Header struct {
	Number       int64
	PreviousHash []byte
	DataHash     []byte
}


func (b *BlockHeader) Bytes() []byte {
	asn1Header := asn1Header{
		PreviousHash: b.PreviousHash,
		DataHash:     b.DataHash,
	}
	if b.Number > uint64(math.MaxInt64) {
		panic(fmt.Errorf("Golang does not currently support encoding uint64 to asn1"))
	} else {
		asn1Header.Number = int64(b.Number)
	}
	result, err := asn1.Marshal(asn1Header)
	if err != nil {
		
		
		
		panic(err)
	}
	return result
}



func (b *BlockHeader) Hash() []byte {
	return util.ComputeSHA256(b.Bytes())
}





func (b *BlockData) Bytes() []byte {
	return util.ConcatenateBytes(b.Data...)
}


func (b *BlockData) Hash() []byte {
	return util.ComputeSHA256(b.Bytes())
}
