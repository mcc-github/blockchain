/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/







package shimtest

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

const (
	minUnicodeRuneValue   = 0 
	compositeKeyNamespace = "\x00"
)



type MockStub struct {
	
	args [][]byte

	
	
	cc shim.Chaincode

	
	Name string

	
	State map[string][]byte

	
	Keys *list.List

	
	Invokables map[string]*MockStub

	
	
	TxID string

	TxTimestamp *timestamp.Timestamp

	
	signedProposal *pb.SignedProposal

	
	ChannelID string

	PvtState map[string]map[string][]byte

	
	EndorsementPolicies map[string]map[string][]byte

	
	ChaincodeEventsChannel chan *pb.ChaincodeEvent

	Creator []byte

	Decorations map[string][]byte
}

func (stub *MockStub) GetTxID() string {
	return stub.TxID
}

func (stub *MockStub) GetChannelID() string {
	return stub.ChannelID
}

func (stub *MockStub) GetArgs() [][]byte {
	return stub.args
}

func (stub *MockStub) GetStringArgs() []string {
	args := stub.GetArgs()
	strargs := make([]string, 0, len(args))
	for _, barg := range args {
		strargs = append(strargs, string(barg))
	}
	return strargs
}

func (stub *MockStub) GetFunctionAndParameters() (function string, params []string) {
	allargs := stub.GetStringArgs()
	function = ""
	params = []string{}
	if len(allargs) >= 1 {
		function = allargs[0]
		params = allargs[1:]
	}
	return
}




func (stub *MockStub) MockTransactionStart(txid string) {
	stub.TxID = txid
	stub.setSignedProposal(&pb.SignedProposal{})
	stub.setTxTimestamp(ptypes.TimestampNow())
}


func (stub *MockStub) MockTransactionEnd(uuid string) {
	stub.signedProposal = nil
	stub.TxID = ""
}





func (stub *MockStub) MockPeerChaincode(invokableChaincodeName string, otherStub *MockStub, channel string) {
	
	if channel != "" {
		invokableChaincodeName = invokableChaincodeName + "/" + channel
	}
	stub.Invokables[invokableChaincodeName] = otherStub
}


func (stub *MockStub) MockInit(uuid string, args [][]byte) pb.Response {
	stub.args = args
	stub.MockTransactionStart(uuid)
	res := stub.cc.Init(stub)
	stub.MockTransactionEnd(uuid)
	return res
}


func (stub *MockStub) MockInvoke(uuid string, args [][]byte) pb.Response {
	stub.args = args
	stub.MockTransactionStart(uuid)
	res := stub.cc.Invoke(stub)
	stub.MockTransactionEnd(uuid)
	return res
}

func (stub *MockStub) GetDecorations() map[string][]byte {
	return stub.Decorations
}


func (stub *MockStub) MockInvokeWithSignedProposal(uuid string, args [][]byte, sp *pb.SignedProposal) pb.Response {
	stub.args = args
	stub.MockTransactionStart(uuid)
	stub.signedProposal = sp
	res := stub.cc.Invoke(stub)
	stub.MockTransactionEnd(uuid)
	return res
}

func (stub *MockStub) GetPrivateData(collection string, key string) ([]byte, error) {
	m, in := stub.PvtState[collection]

	if !in {
		return nil, nil
	}

	return m[key], nil
}

func (stub *MockStub) GetPrivateDataHash(collection, key string) ([]byte, error) {
	return nil, errors.New("Not Implemented")
}

func (stub *MockStub) PutPrivateData(collection string, key string, value []byte) error {
	m, in := stub.PvtState[collection]
	if !in {
		stub.PvtState[collection] = make(map[string][]byte)
		m, in = stub.PvtState[collection]
	}

	m[key] = value

	return nil
}

func (stub *MockStub) DelPrivateData(collection string, key string) error {
	return errors.New("Not Implemented")
}

func (stub *MockStub) GetPrivateDataByRange(collection, startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	return nil, errors.New("Not Implemented")
}

func (stub *MockStub) GetPrivateDataByPartialCompositeKey(collection, objectType string, attributes []string) (shim.StateQueryIteratorInterface, error) {
	return nil, errors.New("Not Implemented")
}

func (stub *MockStub) GetPrivateDataQueryResult(collection, query string) (shim.StateQueryIteratorInterface, error) {
	
	
	
	return nil, errors.New("Not Implemented")
}


func (stub *MockStub) GetState(key string) ([]byte, error) {
	value := stub.State[key]
	return value, nil
}


func (stub *MockStub) PutState(key string, value []byte) error {
	if stub.TxID == "" {
		err := errors.New("cannot PutState without a transactions - call stub.MockTransactionStart()?")
		return err
	}

	
	if len(value) == 0 {
		return stub.DelState(key)
	}
	stub.State[key] = value

	
	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		elemValue := elem.Value.(string)
		comp := strings.Compare(key, elemValue)
		if comp < 0 {
			
			stub.Keys.InsertBefore(key, elem)
			break
		} else if comp == 0 {
			
			break
		} else { 
			
			if elem.Next() == nil {
				stub.Keys.PushBack(key)
				break
			}
		}
	}

	
	if stub.Keys.Len() == 0 {
		stub.Keys.PushFront(key)
	}

	return nil
}


func (stub *MockStub) DelState(key string) error {
	delete(stub.State, key)

	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		if strings.Compare(key, elem.Value.(string)) == 0 {
			stub.Keys.Remove(elem)
		}
	}

	return nil
}

func (stub *MockStub) GetStateByRange(startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	if err := validateSimpleKeys(startKey, endKey); err != nil {
		return nil, err
	}
	return NewMockStateRangeQueryIterator(stub, startKey, endKey), nil
}





func validateSimpleKeys(simpleKeys ...string) error {
	for _, key := range simpleKeys {
		if len(key) > 0 && key[0] == compositeKeyNamespace[0] {
			return fmt.Errorf(`first character of the key [%s] contains a null character which is not allowed`, key)
		}
	}
	return nil
}






func (stub *MockStub) GetQueryResult(query string) (shim.StateQueryIteratorInterface, error) {
	
	
	
	return nil, errors.New("not implemented")
}



func (stub *MockStub) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, error) {
	return nil, errors.New("not implemented")
}







func (stub *MockStub) GetStateByPartialCompositeKey(objectType string, attributes []string) (shim.StateQueryIteratorInterface, error) {
	partialCompositeKey, err := stub.CreateCompositeKey(objectType, attributes)
	if err != nil {
		return nil, err
	}
	return NewMockStateRangeQueryIterator(stub, partialCompositeKey, partialCompositeKey+string(utf8.MaxRune)), nil
}



func (stub *MockStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return shim.CreateCompositeKey(objectType, attributes)
}



func (stub *MockStub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return splitCompositeKey(compositeKey)
}

func splitCompositeKey(compositeKey string) (string, []string, error) {
	componentIndex := 1
	components := []string{}
	for i := 1; i < len(compositeKey); i++ {
		if compositeKey[i] == minUnicodeRuneValue {
			components = append(components, compositeKey[componentIndex:i])
			componentIndex = i + 1
		}
	}
	return components[0], components[1:], nil
}

func (stub *MockStub) GetStateByRangeWithPagination(startKey, endKey string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	return nil, nil, nil
}

func (stub *MockStub) GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string,
	pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	return nil, nil, nil
}

func (stub *MockStub) GetQueryResultWithPagination(query string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	return nil, nil, nil
}





func (stub *MockStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response {
	
	if channel != "" {
		chaincodeName = chaincodeName + "/" + channel
	}
	
	otherStub := stub.Invokables[chaincodeName]
	
	res := otherStub.MockInvoke(stub.TxID, args)
	return res
}

func (stub *MockStub) GetCreator() ([]byte, error) {
	return stub.Creator, nil
}


func (stub *MockStub) GetTransient() (map[string][]byte, error) {
	return nil, nil
}


func (stub *MockStub) GetBinding() ([]byte, error) {
	return nil, nil
}


func (stub *MockStub) GetSignedProposal() (*pb.SignedProposal, error) {
	return stub.signedProposal, nil
}

func (stub *MockStub) setSignedProposal(sp *pb.SignedProposal) {
	stub.signedProposal = sp
}


func (stub *MockStub) GetArgsSlice() ([]byte, error) {
	return nil, nil
}

func (stub *MockStub) setTxTimestamp(time *timestamp.Timestamp) {
	stub.TxTimestamp = time
}

func (stub *MockStub) GetTxTimestamp() (*timestamp.Timestamp, error) {
	if stub.TxTimestamp == nil {
		return nil, errors.New("TxTimestamp not set.")
	}
	return stub.TxTimestamp, nil
}

func (stub *MockStub) SetEvent(name string, payload []byte) error {
	stub.ChaincodeEventsChannel <- &pb.ChaincodeEvent{EventName: name, Payload: payload}
	return nil
}

func (stub *MockStub) SetStateValidationParameter(key string, ep []byte) error {
	return stub.SetPrivateDataValidationParameter("", key, ep)
}

func (stub *MockStub) GetStateValidationParameter(key string) ([]byte, error) {
	return stub.GetPrivateDataValidationParameter("", key)
}

func (stub *MockStub) SetPrivateDataValidationParameter(collection, key string, ep []byte) error {
	m, in := stub.EndorsementPolicies[collection]
	if !in {
		stub.EndorsementPolicies[collection] = make(map[string][]byte)
		m, in = stub.EndorsementPolicies[collection]
	}

	m[key] = ep
	return nil
}

func (stub *MockStub) GetPrivateDataValidationParameter(collection, key string) ([]byte, error) {
	m, in := stub.EndorsementPolicies[collection]

	if !in {
		return nil, nil
	}

	return m[key], nil
}


func NewMockStub(name string, cc shim.Chaincode) *MockStub {
	s := new(MockStub)
	s.Name = name
	s.cc = cc
	s.State = make(map[string][]byte)
	s.PvtState = make(map[string]map[string][]byte)
	s.EndorsementPolicies = make(map[string]map[string][]byte)
	s.Invokables = make(map[string]*MockStub)
	s.Keys = list.New()
	s.ChaincodeEventsChannel = make(chan *pb.ChaincodeEvent, 100) 
	s.Decorations = make(map[string][]byte)

	return s
}



type MockStateRangeQueryIterator struct {
	Closed   bool
	Stub     *MockStub
	StartKey string
	EndKey   string
	Current  *list.Element
}



func (iter *MockStateRangeQueryIterator) HasNext() bool {
	if iter.Closed {
		
		return false
	}

	if iter.Current == nil {
		return false
	}

	current := iter.Current
	for current != nil {
		
		if iter.StartKey == "" && iter.EndKey == "" {
			return true
		}
		comp1 := strings.Compare(current.Value.(string), iter.StartKey)
		comp2 := strings.Compare(current.Value.(string), iter.EndKey)
		if comp1 >= 0 {
			if comp2 < 0 {
				return true
			} else {
				return false

			}
		}
		current = current.Next()
	}
	return false
}


func (iter *MockStateRangeQueryIterator) Next() (*queryresult.KV, error) {
	if iter.Closed == true {
		err := errors.New("MockStateRangeQueryIterator.Next() called after Close()")
		return nil, err
	}

	if iter.HasNext() == false {
		err := errors.New("MockStateRangeQueryIterator.Next() called when it does not HaveNext()")
		return nil, err
	}

	for iter.Current != nil {
		comp1 := strings.Compare(iter.Current.Value.(string), iter.StartKey)
		comp2 := strings.Compare(iter.Current.Value.(string), iter.EndKey)
		
		
		if (comp1 >= 0 && comp2 < 0) || (iter.StartKey == "" && iter.EndKey == "") {
			key := iter.Current.Value.(string)
			value, err := iter.Stub.GetState(key)
			iter.Current = iter.Current.Next()
			return &queryresult.KV{Key: key, Value: value}, err
		}
		iter.Current = iter.Current.Next()
	}
	err := errors.New("MockStateRangeQueryIterator.Next() went past end of range")
	return nil, err
}



func (iter *MockStateRangeQueryIterator) Close() error {
	if iter.Closed == true {
		err := errors.New("MockStateRangeQueryIterator.Close() called after Close()")
		return err
	}

	iter.Closed = true
	return nil
}

func NewMockStateRangeQueryIterator(stub *MockStub, startKey string, endKey string) *MockStateRangeQueryIterator {
	iter := new(MockStateRangeQueryIterator)
	iter.Closed = false
	iter.Stub = stub
	iter.StartKey = startKey
	iter.EndKey = endKey
	iter.Current = stub.Keys.Front()
	return iter
}

func getBytes(function string, args []string) [][]byte {
	bytes := make([][]byte, 0, len(args)+1)
	bytes = append(bytes, []byte(function))
	for _, s := range args {
		bytes = append(bytes, []byte(s))
	}
	return bytes
}

func getFuncArgs(bytes [][]byte) (string, []string) {
	function := string(bytes[0])
	args := make([]string, len(bytes)-1)
	for i := 1; i < len(bytes); i++ {
		args[i-1] = string(bytes[i])
	}
	return function, args
}
