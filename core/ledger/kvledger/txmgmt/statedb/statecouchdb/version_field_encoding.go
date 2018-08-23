/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package statecouchdb

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	proto "github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb/statecouchdb/msgs"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
)

func encodeVersionAndMetadata(version *version.Height, metadata []byte) (string, error) {
	msg := &msgs.VersionFieldProto{
		VersionBytes: version.ToBytes(),
		Metadata:     metadata,
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return "", err
	}
	msgBase64 := base64.StdEncoding.EncodeToString(msgBytes)
	encodedVersionField := append([]byte{byte(0)}, []byte(msgBase64)...)
	return string(encodedVersionField), nil
}

func decodeVersionAndMetadata(encodedstr string) (*version.Height, []byte, error) {
	if oldFormatEncoding(encodedstr) {
		return decodeVersionOldFormat(encodedstr), nil, nil
	}
	versionFieldBytes, err := base64.StdEncoding.DecodeString(encodedstr[1:])
	if err != nil {
		return nil, nil, err
	}
	versionFieldMsg := &msgs.VersionFieldProto{}
	if err = proto.Unmarshal(versionFieldBytes, versionFieldMsg); err != nil {
		return nil, nil, err
	}
	ver, _ := version.NewHeightFromBytes(versionFieldMsg.VersionBytes)
	return ver, versionFieldMsg.Metadata, nil
}





func encodeVersionOldFormat(version *version.Height) string {
	return fmt.Sprintf("%v:%v", version.BlockNum, version.TxNum)
}







func decodeVersionOldFormat(encodedVersion string) *version.Height {
	versionArray := strings.Split(fmt.Sprintf("%s", encodedVersion), ":")
	
	blockNum, _ := strconv.ParseUint(versionArray[0], 10, 64)
	
	txNum, _ := strconv.ParseUint(versionArray[1], 10, 64)
	return version.NewHeight(blockNum, txNum)
}

func oldFormatEncoding(encodedstr string) bool {
	return []byte(encodedstr)[0] != byte(0)
}
