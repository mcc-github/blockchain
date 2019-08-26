/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tx

import (
	"fmt"
	"io"

	"github.com/mcc-github/blockchain/pkg/statedata"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
)









type ProcessorCreator interface {
	NewProcessor(txenv *Envelope) (processor Processor, simulatedRWSet [][]byte, err error)
}








































type Processor interface {
	Preprocess(latestChannelConfig *ChannelConfig) error
	Process(state *State, proposedWrites *statedata.ProposedWrites) error
	Done()
}


type InvalidErr struct {
	ActualErr      error
	ValidationCode peer.TxValidationCode
}

func (e *InvalidErr) msgWithoutStack() string {
	return fmt.Sprintf("ValidationCode = %s, ActualErr = %s", e.ValidationCode.String(), e.ActualErr)
}

func (e *InvalidErr) msgWithStack() string {
	return fmt.Sprintf("ValidationCode = %s, ActualErr = %+v", e.ValidationCode.String(), e.ActualErr)
}

func (e *InvalidErr) Error() string {
	return e.msgWithoutStack()
}


func (e *InvalidErr) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.msgWithStack())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msgWithoutStack())
	case 'q':
		fmt.Fprintf(s, "%q", e.msgWithoutStack())
	}
}










type ReadHinter interface {
	ReadHint(potentialWrites *statedata.WriteHint) *statedata.ReadHint
}














type Reprocessor interface {
	Reprocess(state *State, latestChannelConfig *ChannelConfig, proposedWrites *statedata.ProposedWrites)
}




type ReprocessReadHinter interface {
	ReprocessReadHint(potentialWrites *statedata.WriteHint) *statedata.ReadHint
}




type PvtdataSourceHinter interface {
	PvtdataSource() [][]byte
}


type ChannelConfig struct {
}




type State struct {
	BackingState state
}


func (s *State) GetState(ns, key string) ([]byte, error) {
	return s.BackingState.GetState(ns, key)
}


func (s *State) GetStateMetadata(ns, key string) (map[string][]byte, error) {
	return s.BackingState.GetStateMetadata(ns, key)
}



func (s *State) GetStateRangeScanIterator(ns, startKey, endKey string) (*KeyValueItr, error) {
	itr, err := s.BackingState.GetStateRangeScanIterator(ns, startKey, endKey)
	if err != nil {
		return nil, err
	}
	return &KeyValueItr{
		BackingItr: itr,
	}, nil
}



func (s *State) SetState(ns, key string, value []byte) error {
	return s.BackingState.SetState(ns, key, value)
}



func (s *State) SetStateMetadata(ns, key string, metadata map[string][]byte) error {
	return s.BackingState.SetStateMetadata(ns, key, metadata)
}


func (s *State) GetPrivateDataMetadataByHash(ns, coll string, keyHash []byte) (map[string][]byte, error) {
	return s.BackingState.GetPrivateDataMetadataByHash(ns, coll, keyHash)
}


type KeyValueItr struct {
	BackingItr keyValueItr
}


func (i *KeyValueItr) Next() (*statedata.KeyValue, error) {
	return i.BackingItr.Next()
}


func (i *KeyValueItr) Close() {
	i.BackingItr.Close()
}







type Envelope struct {
	
	SignedBytes []byte
	
	Signature []byte
	
	Data []byte
	
	HeaderBytes []byte
	
	ChannelHeaderBytes []byte
	
	SignatureHeaderBytes []byte
	
	ChannelHeader *common.ChannelHeader
	
	SignatureHeader *common.SignatureHeader
}



type state interface {
	GetState(ns, key string) ([]byte, error)
	GetStateMetadata(ns, key string) (map[string][]byte, error)
	GetStateRangeScanIterator(ns, startKey, endKey string) (keyValueItr, error)
	SetState(ns, key string, value []byte) error
	SetStateMetadata(ns, key string, metadata map[string][]byte) error
	GetPrivateDataMetadataByHash(ns, coll string, keyHash []byte) (map[string][]byte, error)
}


type keyValueItr interface {
	Next() (*statedata.KeyValue, error)
	Close()
}
