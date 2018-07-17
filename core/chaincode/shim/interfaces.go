/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package shim

import (
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	pb "github.com/mcc-github/blockchain/protos/peer"
)



type Chaincode interface {
	
	
	
	Init(stub ChaincodeStubInterface) pb.Response

	
	
	
	Invoke(stub ChaincodeStubInterface) pb.Response
}



type ChaincodeStubInterface interface {
	
	
	GetArgs() [][]byte

	
	
	
	GetStringArgs() []string

	
	
	
	
	GetFunctionAndParameters() (string, []string)

	
	
	GetArgsSlice() ([]byte, error)

	
	
	
	GetTxID() string

	
	
	
	
	GetChannelID() string

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response

	
	
	
	
	
	GetState(key string) ([]byte, error)

	
	
	
	
	
	
	
	PutState(key string, value []byte) error

	
	
	
	DelState(key string) error

	
	
	
	
	
	
	
	
	
	GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)

	
	
	
	
	
	
	
	
	
	
	GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error)

	
	
	
	
	
	CreateCompositeKey(objectType string, attributes []string) (string, error)

	
	
	
	
	SplitCompositeKey(compositeKey string) (string, []string, error)

	
	
	
	
	
	
	
	
	
	
	
	GetQueryResult(query string) (StateQueryIteratorInterface, error)

	
	
	
	
	
	
	
	
	
	
	
	
	GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)

	
	
	
	
	
	GetPrivateData(collection, key string) ([]byte, error)

	
	
	
	
	
	
	
	
	
	PutPrivateData(collection string, key string, value []byte) error

	
	
	
	
	
	
	DelPrivateData(collection, key string) error

	
	
	
	
	
	
	
	
	
	GetPrivateDataByRange(collection, startKey, endKey string) (StateQueryIteratorInterface, error)

	
	
	
	
	
	
	
	
	
	
	GetPrivateDataByPartialCompositeKey(collection, objectType string, keys []string) (StateQueryIteratorInterface, error)

	
	
	
	
	
	
	
	
	
	
	
	GetPrivateDataQueryResult(collection, query string) (StateQueryIteratorInterface, error)

	
	
	
	GetCreator() ([]byte, error)

	
	
	
	
	
	
	GetTransient() (map[string][]byte, error)

	
	
	
	
	GetBinding() ([]byte, error)

	
	
	
	GetDecorations() map[string][]byte

	
	
	GetSignedProposal() (*pb.SignedProposal, error)

	
	
	
	GetTxTimestamp() (*timestamp.Timestamp, error)

	
	
	
	
	SetEvent(name string, payload []byte) error
}



type CommonIteratorInterface interface {
	
	
	HasNext() bool

	
	
	Close() error
}



type StateQueryIteratorInterface interface {
	
	CommonIteratorInterface

	
	Next() (*queryresult.KV, error)
}



type HistoryQueryIteratorInterface interface {
	
	CommonIteratorInterface

	
	Next() (*queryresult.KeyModification, error)
}





type MockQueryIteratorInterface interface {
	StateQueryIteratorInterface
}
