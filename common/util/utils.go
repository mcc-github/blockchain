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

package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/metadata"
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


func GetTestChainID() string {
	return testchainid
}





func GetSysCCVersion() string {
	return metadata.Version
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
