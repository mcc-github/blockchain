/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	validatorstate "github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"

	"github.com/pkg/errors"
)

type StateIterator interface {
	Close() error
	Next() (*queryresult.KV, error)
}



func StateIteratorToMap(itr StateIterator) (map[string][]byte, error) {
	defer itr.Close()
	result := map[string][]byte{}
	for {
		entry, err := itr.Next()
		if err != nil {
			return nil, errors.WithMessage(err, "could not iterate over range")
		}
		if entry == nil {
			return result, nil
		}
		result[entry.Key] = entry.Value
	}
}



type ChaincodePublicLedgerShim struct {
	shim.ChaincodeStubInterface
}



func (cls *ChaincodePublicLedgerShim) GetStateRange(prefix string) (map[string][]byte, error) {
	itr, err := cls.GetStateByRange(prefix, prefix+"\x7f")
	if err != nil {
		return nil, errors.WithMessage(err, "could not get state iterator")
	}
	return StateIteratorToMap(&ChaincodeResultIteratorShim{ResultsIterator: itr})
}

type ChaincodeResultIteratorShim struct {
	ResultsIterator shim.StateQueryIteratorInterface
}

func (cris *ChaincodeResultIteratorShim) Next() (*queryresult.KV, error) {
	if !cris.ResultsIterator.HasNext() {
		return nil, nil
	}
	return cris.ResultsIterator.Next()
}

func (cris *ChaincodeResultIteratorShim) Close() error {
	return cris.ResultsIterator.Close()
}



type ChaincodePrivateLedgerShim struct {
	Stub       shim.ChaincodeStubInterface
	Collection string
}



func (cls *ChaincodePrivateLedgerShim) GetStateRange(prefix string) (map[string][]byte, error) {
	itr, err := cls.Stub.GetPrivateDataByRange(cls.Collection, prefix, prefix+"\x7f")
	if err != nil {
		return nil, errors.WithMessage(err, "could not get state iterator")
	}
	return StateIteratorToMap(&ChaincodeResultIteratorShim{ResultsIterator: itr})
}


func (cls *ChaincodePrivateLedgerShim) GetState(key string) ([]byte, error) {
	return cls.Stub.GetPrivateData(cls.Collection, key)
}


func (cls *ChaincodePrivateLedgerShim) GetStateHash(key string) ([]byte, error) {
	return cls.Stub.GetPrivateDataHash(cls.Collection, key)
}


func (cls *ChaincodePrivateLedgerShim) PutState(key string, value []byte) error {
	return cls.Stub.PutPrivateData(cls.Collection, key, value)
}


func (cls *ChaincodePrivateLedgerShim) DelState(key string) error {
	return cls.Stub.DelPrivateData(cls.Collection, key)
}

func (cls *ChaincodePrivateLedgerShim) CollectionName() string {
	return cls.Collection
}



type SimpleQueryExecutorShim struct {
	Namespace           string
	SimpleQueryExecutor ledger.SimpleQueryExecutor
}

func (sqes *SimpleQueryExecutorShim) GetState(key string) ([]byte, error) {
	return sqes.SimpleQueryExecutor.GetState(sqes.Namespace, key)
}

func (sqes *SimpleQueryExecutorShim) GetStateRange(prefix string) (map[string][]byte, error) {
	itr, err := sqes.SimpleQueryExecutor.GetStateRangeScanIterator(sqes.Namespace, prefix, prefix+"\x7f")
	if err != nil {
		return nil, errors.WithMessage(err, "could not get state iterator")
	}
	return StateIteratorToMap(&ResultsIteratorShim{ResultsIterator: itr})
}

type ResultsIteratorShim struct {
	ResultsIterator commonledger.ResultsIterator
}

func (ris *ResultsIteratorShim) Next() (*queryresult.KV, error) {
	res, err := ris.ResultsIterator.Next()
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(*queryresult.KV), err
}

func (ris *ResultsIteratorShim) Close() error {
	ris.ResultsIterator.Close()
	return nil
}

type ValidatorStateShim struct {
	ValidatorState validatorstate.State
	Namespace      string
}

func (vss *ValidatorStateShim) GetState(key string) ([]byte, error) {
	result, err := vss.ValidatorState.GetStateMultipleKeys(vss.Namespace, []string{key})
	if err != nil {
		return nil, errors.WithMessage(err, "could not get state thought validatorstate shim")
	}
	return result[0], nil
}

type PrivateQueryExecutor interface {
	GetPrivateDataHash(namespace, collection, key string) (value []byte, err error)
}

type PrivateQueryExecutorShim struct {
	Namespace  string
	Collection string
	State      PrivateQueryExecutor
}

func (pqes *PrivateQueryExecutorShim) GetStateHash(key string) ([]byte, error) {
	return pqes.State.GetPrivateDataHash(pqes.Namespace, pqes.Collection, key)
}

func (pqes *PrivateQueryExecutorShim) CollectionName() string {
	return pqes.Collection
}





type DummyQueryExecutorShim struct {
}

func (*DummyQueryExecutorShim) GetState(key string) ([]byte, error) {
	return nil, errors.New("invalid channel-less operation")
}
