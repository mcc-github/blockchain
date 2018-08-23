/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"encoding/binary"
	"math"

	"github.com/golang/protobuf/proto"
)






func EncodeReverseOrderVarUint64(number uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, math.MaxUint64-number)
	numFFBytes := 0
	for _, b := range bytes {
		if b != 0xff {
			break
		}
		numFFBytes++
	}
	size := 8 - numFFBytes
	encodedBytes := make([]byte, size+1)
	encodedBytes[0] = proto.EncodeVarint(uint64(numFFBytes))[0]
	copy(encodedBytes[1:], bytes[numFFBytes:])
	return encodedBytes
}



func DecodeReverseOrderVarUint64(bytes []byte) (uint64, int) {
	s, _ := proto.DecodeVarint(bytes)
	numFFBytes := int(s)
	decodedBytes := make([]byte, 8)
	realBytesNum := 8 - numFFBytes
	copy(decodedBytes[numFFBytes:], bytes[1:realBytesNum+1])
	numBytesConsumed := realBytesNum + 1
	for i := 0; i < numFFBytes; i++ {
		decodedBytes[i] = 0xff
	}
	return (math.MaxUint64 - binary.BigEndian.Uint64(decodedBytes)), numBytesConsumed
}
