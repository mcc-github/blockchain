/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package shim

import (
	pb "github.com/mcc-github/blockchain/protos/peer"
)

const (
	
	
	OK = 200

	
	ERRORTHRESHOLD = 400

	
	ERROR = 500
)

func Success(payload []byte) pb.Response {
	return pb.Response{
		Status:  OK,
		Payload: payload,
	}
}

func Error(msg string) pb.Response {
	return pb.Response{
		Status:  ERROR,
		Message: msg,
	}
}
