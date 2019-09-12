/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
)


func ComputeSHA256(data []byte) (hash []byte) {
	hash, err := factory.GetDefault().Hash(data, &bccsp.SHA256Opts{})
	if err != nil {
		panic(fmt.Errorf("Failed computing SHA256 on [% x]", data))
	}
	return
}


func ComputeSHA3256(data []byte) (hash []byte) {
	hash, err := factory.GetDefault().Hash(data, &bccsp.SHA3_256Opts{})
	if err != nil {
		panic(fmt.Errorf("Failed computing SHA3_256 on [% x]", data))
	}
	return
}


func GenerateBytesUUID() []byte {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(fmt.Sprintf("Error generating UUID: %s", err))
	}

	
	uuid[8] = uuid[8]&^0xc0 | 0x80

	
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return uuid
}


func GenerateUUID() string {
	uuid := GenerateBytesUUID()
	return idBytesToStr(uuid)
}


func CreateUtcTimestamp() *timestamp.Timestamp {
	now := time.Now().UTC()
	secs := now.Unix()
	nanos := int32(now.UnixNano() - (secs * 1000000000))
	return &(timestamp.Timestamp{Seconds: secs, Nanos: nanos})
}

func idBytesToStr(id []byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", id[0:4], id[4:6], id[6:8], id[8:10], id[10:])
}


func ToChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

const testchainid = "testchainid"


func GetTestChannelID() string {
	return testchainid
}



func ConcatenateBytes(data ...[]byte) []byte {
	finalLength := 0
	for _, slice := range data {
		finalLength += len(slice)
	}
	result := make([]byte, finalLength)
	last := 0
	for _, slice := range data {
		for i := range slice {
			result[i+last] = slice[i]
		}
		last += len(slice)
	}
	return result
}
