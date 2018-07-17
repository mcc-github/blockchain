/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import "bytes"

func init() {
	
	
	switch true {

	}
}



type PKIidType []byte




func (id PKIidType) IsNotSameFilter(that PKIidType) bool {
	return !bytes.Equal(id, that)
}




type MessageAcceptor func(interface{}) bool


type Payload struct {
	ChainID ChainID 
	Data    []byte  
	Hash    string  
	SeqNum  uint64  
}


type ChainID []byte





type MessageReplacingPolicy func(this interface{}, that interface{}) InvalidationResult



type InvalidationResult int

const (
	
	MessageNoAction InvalidationResult = iota
	
	MessageInvalidates
	
	MessageInvalidated
)
