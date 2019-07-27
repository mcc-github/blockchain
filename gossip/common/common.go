/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"bytes"
	"encoding/hex"
)



type PKIidType []byte

func (p PKIidType) String() string {
	if p == nil {
		return "<nil>"
	}
	return hex.EncodeToString(p)
}




func (id PKIidType) IsNotSameFilter(that PKIidType) bool {
	return !bytes.Equal(id, that)
}




type MessageAcceptor func(interface{}) bool


type Payload struct {
	ChannelID ChannelID 
	Data      []byte    
	Hash      string    
	SeqNum    uint64    
}


type ChannelID []byte





type MessageReplacingPolicy func(this interface{}, that interface{}) InvalidationResult



type InvalidationResult int

const (
	
	MessageNoAction InvalidationResult = iota
	
	MessageInvalidates
	
	MessageInvalidated
)
