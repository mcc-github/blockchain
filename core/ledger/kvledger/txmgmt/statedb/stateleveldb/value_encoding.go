/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package stateleveldb

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb/stateleveldb/msgs"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
)


func encodeValue(v *statedb.VersionedValue) ([]byte, error) {
	vvMsg := &msgs.VersionedValueProto{
		VersionBytes: v.Version.ToBytes(),
		Value:        v.Value,
		Metadata:     v.Metadata,
	}
	encodedValue, err := proto.Marshal(vvMsg)
	if err != nil {
		return nil, err
	}
	encodedValue = append([]byte{0}, encodedValue...)
	return encodedValue, nil
}

func decodeValue(encodedValue []byte) (*statedb.VersionedValue, error) {
	if oldFormatEncoding(encodedValue) {
		val, ver := decodeValueOldFormat(encodedValue)
		return &statedb.VersionedValue{Version: ver, Value: val, Metadata: nil}, nil
	}
	msg := &msgs.VersionedValueProto{}
	err := proto.Unmarshal(encodedValue[1:], msg)
	if err != nil {
		return nil, err
	}
	ver, _ := version.NewHeightFromBytes(msg.VersionBytes)
	val := msg.Value
	metadata := msg.Metadata
	
	if val == nil {
		val = []byte{}
	}
	return &statedb.VersionedValue{Version: ver, Value: val, Metadata: metadata}, nil
}





func encodeValueOldFormat(value []byte, version *version.Height) []byte {
	encodedValue := version.ToBytes()
	if value != nil {
		encodedValue = append(encodedValue, value...)
	}
	return encodedValue
}







func decodeValueOldFormat(encodedValue []byte) ([]byte, *version.Height) {
	height, n := version.NewHeightFromBytes(encodedValue)
	value := encodedValue[n:]
	return value, height
}

func oldFormatEncoding(encodedValue []byte) bool {
	return encodedValue[0] != byte(0)
}
