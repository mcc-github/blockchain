/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"fmt"

	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
)


type Support interface {
	BatchSize() *ab.BatchSize
}


func NewSizeFilter(support Support) *MaxBytesRule {
	return &MaxBytesRule{support: support}
}


type MaxBytesRule struct {
	support Support
}


func (r *MaxBytesRule) Apply(message *cb.Envelope) error {
	maxBytes := r.support.BatchSize().AbsoluteMaxBytes
	if size := messageByteSize(message); size > maxBytes {
		return fmt.Errorf("message payload is %d bytes and exceeds maximum allowed %d bytes", size, maxBytes)
	}
	return nil
}

func messageByteSize(message *cb.Envelope) uint32 {
	
	
	return uint32(len(message.Payload) + len(message.Signature))
}
