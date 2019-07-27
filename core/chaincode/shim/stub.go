/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package shim

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	pb "github.com/mcc-github/blockchain/protos/peer"
)



type ChaincodeStub struct {
	TxID                       string
	ChannelId                  string
	chaincodeEvent             *pb.ChaincodeEvent
	args                       [][]byte
	handler                    *Handler
	signedProposal             *pb.SignedProposal
	proposal                   *pb.Proposal
	validationParameterMetakey string

	
	creator   []byte
	transient map[string][]byte
	binding   []byte

	decorations map[string][]byte
}



func newChaincodeStub(handler *Handler, channelId, txid string, input *pb.ChaincodeInput, signedProposal *pb.SignedProposal) (*ChaincodeStub, error) {
	stub := &ChaincodeStub{
		TxID:                       txid,
		ChannelId:                  channelId,
		args:                       input.Args,
		handler:                    handler,
		signedProposal:             signedProposal,
		decorations:                input.Decorations,
		validationParameterMetakey: pb.MetaDataKeys_VALIDATION_PARAMETER.String(),
	}

	
	
	
	if signedProposal != nil {
		var err error

		stub.proposal = &pb.Proposal{}
		err = proto.Unmarshal(signedProposal.ProposalBytes, stub.proposal)
		if err != nil {

			return nil, fmt.Errorf("failed to extract Proposal from SignedProposal: %s", err)
		}

		
		if len(stub.proposal.GetHeader()) == 0 {
			return nil, errors.New("failed to extract Proposal fields: proposal header is nil")
		}

		
		hdr := &common.Header{}
		if err := proto.Unmarshal(stub.proposal.GetHeader(), hdr); err != nil {
			return nil, fmt.Errorf("failed to extract proposal header: %s", err)
		}

		
		chdr := &common.ChannelHeader{}
		if err := proto.Unmarshal(hdr.ChannelHeader, chdr); err != nil {
			return nil, fmt.Errorf("failed to extract channel header: %s", err)
		}
		validTypes := map[common.HeaderType]bool{
			common.HeaderType_ENDORSER_TRANSACTION: true,
			common.HeaderType_CONFIG:               true,
		}
		if !validTypes[common.HeaderType(chdr.GetType())] {
			return nil, fmt.Errorf(
				"invalid channel header type. Expected %s or %s, received %s",
				common.HeaderType_ENDORSER_TRANSACTION,
				common.HeaderType_CONFIG,
				common.HeaderType(chdr.GetType()),
			)
		}

		
		shdr := &common.SignatureHeader{}
		if err := proto.Unmarshal(hdr.GetSignatureHeader(), shdr); err != nil {
			return nil, fmt.Errorf("failed to extract signature header: %s", err)
		}
		stub.creator = shdr.GetCreator()

		
		payload := &pb.ChaincodeProposalPayload{}
		if err := proto.Unmarshal(stub.proposal.GetPayload(), payload); err != nil {
			return nil, fmt.Errorf("failed to extract proposal payload: %s", err)
		}
		stub.transient = payload.GetTransientMap()

		
		epoch := make([]byte, 8)
		binary.LittleEndian.PutUint64(epoch, chdr.GetEpoch())
		digest := sha256.Sum256(append(append(shdr.GetNonce(), stub.creator...), epoch...))
		stub.binding = digest[:]

	}

	return stub, nil
}


func (s *ChaincodeStub) GetTxID() string {
	return s.TxID
}


func (s *ChaincodeStub) GetChannelID() string {
	return s.ChannelId
}

func (s *ChaincodeStub) GetDecorations() map[string][]byte {
	return s.decorations
}




func (s *ChaincodeStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response {
	
	if channel != "" {
		chaincodeName = chaincodeName + "/" + channel
	}
	return s.handler.handleInvokeChaincode(chaincodeName, args, s.ChannelId, s.TxID)
}




func (s *ChaincodeStub) GetState(key string) ([]byte, error) {
	
	collection := ""
	return s.handler.handleGetState(collection, key, s.ChannelId, s.TxID)
}


func (s *ChaincodeStub) SetStateValidationParameter(key string, ep []byte) error {
	return s.handler.handlePutStateMetadataEntry("", key, s.validationParameterMetakey, ep, s.ChannelId, s.TxID)
}


func (s *ChaincodeStub) GetStateValidationParameter(key string) ([]byte, error) {
	md, err := s.handler.handleGetStateMetadata("", key, s.ChannelId, s.TxID)
	if err != nil {
		return nil, err
	}
	if ep, ok := md[s.validationParameterMetakey]; ok {
		return ep, nil
	}
	return nil, nil
}


func (s *ChaincodeStub) PutState(key string, value []byte) error {
	if key == "" {
		return errors.New("key must not be an empty string")
	}
	
	collection := ""
	return s.handler.handlePutState(collection, key, value, s.ChannelId, s.TxID)
}

func (s *ChaincodeStub) createStateQueryIterator(response *pb.QueryResponse) *StateQueryIterator {
	return &StateQueryIterator{
		CommonIterator: &CommonIterator{
			handler:    s.handler,
			channelId:  s.ChannelId,
			txid:       s.TxID,
			response:   response,
			currentLoc: 0,
		},
	}
}


func (s *ChaincodeStub) GetQueryResult(query string) (StateQueryIteratorInterface, error) {
	
	collection := ""
	
	iterator, _, err := s.handleGetQueryResult(collection, query, nil)

	return iterator, err
}


func (s *ChaincodeStub) DelState(key string) error {
	
	collection := ""
	return s.handler.handleDelState(collection, key, s.ChannelId, s.TxID)
}




func (s *ChaincodeStub) GetPrivateData(collection string, key string) ([]byte, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must not be an empty string")
	}
	return s.handler.handleGetState(collection, key, s.ChannelId, s.TxID)
}


func (s *ChaincodeStub) GetPrivateDataHash(collection string, key string) ([]byte, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must not be an empty string")
	}
	return s.handler.handleGetPrivateDataHash(collection, key, s.ChannelId, s.TxID)
}


func (s *ChaincodeStub) PutPrivateData(collection string, key string, value []byte) error {
	if collection == "" {
		return fmt.Errorf("collection must not be an empty string")
	}
	if key == "" {
		return fmt.Errorf("key must not be an empty string")
	}
	return s.handler.handlePutState(collection, key, value, s.ChannelId, s.TxID)
}


func (s *ChaincodeStub) DelPrivateData(collection string, key string) error {
	if collection == "" {
		return fmt.Errorf("collection must not be an empty string")
	}
	return s.handler.handleDelState(collection, key, s.ChannelId, s.TxID)
}


func (s *ChaincodeStub) GetPrivateDataByRange(collection, startKey, endKey string) (StateQueryIteratorInterface, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must not be an empty string")
	}
	if startKey == "" {
		startKey = emptyKeySubstitute
	}
	if err := validateSimpleKeys(startKey, endKey); err != nil {
		return nil, err
	}
	
	iterator, _, err := s.handleGetStateByRange(collection, startKey, endKey, nil)

	return iterator, err
}

func (s *ChaincodeStub) createRangeKeysForPartialCompositeKey(objectType string, attributes []string) (string, string, error) {
	partialCompositeKey, err := s.CreateCompositeKey(objectType, attributes)
	if err != nil {
		return "", "", err
	}
	startKey := partialCompositeKey
	endKey := partialCompositeKey + string(maxUnicodeRuneValue)

	return startKey, endKey, nil
}


func (s *ChaincodeStub) GetPrivateDataByPartialCompositeKey(collection, objectType string, attributes []string) (StateQueryIteratorInterface, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must not be an empty string")
	}

	startKey, endKey, err := s.createRangeKeysForPartialCompositeKey(objectType, attributes)
	if err != nil {
		return nil, err
	}
	
	iterator, _, err := s.handleGetStateByRange(collection, startKey, endKey, nil)

	return iterator, err
}


func (s *ChaincodeStub) GetPrivateDataQueryResult(collection, query string) (StateQueryIteratorInterface, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must not be an empty string")
	}
	
	iterator, _, err := s.handleGetQueryResult(collection, query, nil)

	return iterator, err
}


func (s *ChaincodeStub) GetPrivateDataValidationParameter(collection, key string) ([]byte, error) {
	md, err := s.handler.handleGetStateMetadata(collection, key, s.ChannelId, s.TxID)
	if err != nil {
		return nil, err
	}
	if ep, ok := md[s.validationParameterMetakey]; ok {
		return ep, nil
	}
	return nil, nil
}


func (s *ChaincodeStub) SetPrivateDataValidationParameter(collection, key string, ep []byte) error {
	return s.handler.handlePutStateMetadataEntry(collection, key, s.validationParameterMetakey, ep, s.ChannelId, s.TxID)
}


type CommonIterator struct {
	handler    *Handler
	channelId  string
	txid       string
	response   *pb.QueryResponse
	currentLoc int
}


type StateQueryIterator struct {
	*CommonIterator
}


type HistoryQueryIterator struct {
	*CommonIterator
}



type queryResult interface{}

type resultType uint8

const (
	STATE_QUERY_RESULT resultType = iota + 1
	HISTORY_QUERY_RESULT
)

func createQueryResponseMetadata(metadataBytes []byte) (*pb.QueryResponseMetadata, error) {
	metadata := &pb.QueryResponseMetadata{}
	err := proto.Unmarshal(metadataBytes, metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *ChaincodeStub) handleGetStateByRange(collection, startKey, endKey string,
	metadata []byte) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {

	response, err := s.handler.handleGetStateByRange(collection, startKey, endKey, metadata, s.ChannelId, s.TxID)
	if err != nil {
		return nil, nil, err
	}

	iterator := s.createStateQueryIterator(response)
	responseMetadata, err := createQueryResponseMetadata(response.Metadata)
	if err != nil {
		return nil, nil, err
	}

	return iterator, responseMetadata, nil
}

func (s *ChaincodeStub) handleGetQueryResult(collection, query string,
	metadata []byte) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {

	response, err := s.handler.handleGetQueryResult(collection, query, metadata, s.ChannelId, s.TxID)
	if err != nil {
		return nil, nil, err
	}

	iterator := s.createStateQueryIterator(response)
	responseMetadata, err := createQueryResponseMetadata(response.Metadata)
	if err != nil {
		return nil, nil, err
	}

	return iterator, responseMetadata, nil
}


func (s *ChaincodeStub) GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error) {
	if startKey == "" {
		startKey = emptyKeySubstitute
	}
	if err := validateSimpleKeys(startKey, endKey); err != nil {
		return nil, err
	}
	collection := ""

	
	iterator, _, err := s.handleGetStateByRange(collection, startKey, endKey, nil)

	return iterator, err
}


func (s *ChaincodeStub) GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error) {
	response, err := s.handler.handleGetHistoryForKey(key, s.ChannelId, s.TxID)
	if err != nil {
		return nil, err
	}
	return &HistoryQueryIterator{CommonIterator: &CommonIterator{s.handler, s.ChannelId, s.TxID, response, 0}}, nil
}


func (s *ChaincodeStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return CreateCompositeKey(objectType, attributes)
}


func (s *ChaincodeStub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return splitCompositeKey(compositeKey)
}

func CreateCompositeKey(objectType string, attributes []string) (string, error) {
	if err := validateCompositeKeyAttribute(objectType); err != nil {
		return "", err
	}
	ck := compositeKeyNamespace + objectType + string(minUnicodeRuneValue)
	for _, att := range attributes {
		if err := validateCompositeKeyAttribute(att); err != nil {
			return "", err
		}
		ck += att + string(minUnicodeRuneValue)
	}
	return ck, nil
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

func validateCompositeKeyAttribute(str string) error {
	if !utf8.ValidString(str) {
		return fmt.Errorf("not a valid utf8 string: [%x]", str)
	}
	for index, runeValue := range str {
		if runeValue == minUnicodeRuneValue || runeValue == maxUnicodeRuneValue {
			return fmt.Errorf(`input contains unicode %#U starting at position [%d]. %#U and %#U are not allowed in the input attribute of a composite key`,
				runeValue, index, minUnicodeRuneValue, maxUnicodeRuneValue)
		}
	}
	return nil
}





func validateSimpleKeys(simpleKeys ...string) error {
	for _, key := range simpleKeys {
		if len(key) > 0 && key[0] == compositeKeyNamespace[0] {
			return fmt.Errorf(`first character of the key [%s] contains a null character which is not allowed`, key)
		}
	}
	return nil
}







func (s *ChaincodeStub) GetStateByPartialCompositeKey(objectType string, attributes []string) (StateQueryIteratorInterface, error) {
	collection := ""
	startKey, endKey, err := s.createRangeKeysForPartialCompositeKey(objectType, attributes)
	if err != nil {
		return nil, err
	}
	
	iterator, _, err := s.handleGetStateByRange(collection, startKey, endKey, nil)

	return iterator, err
}

func createQueryMetadata(pageSize int32, bookmark string) ([]byte, error) {
	
	metadata := &pb.QueryMetadata{PageSize: pageSize, Bookmark: bookmark}
	metadataBytes, err := proto.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	return metadataBytes, nil
}

func (s *ChaincodeStub) GetStateByRangeWithPagination(startKey, endKey string, pageSize int32,
	bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {

	if startKey == "" {
		startKey = emptyKeySubstitute
	}
	if err := validateSimpleKeys(startKey, endKey); err != nil {
		return nil, nil, err
	}

	collection := ""

	metadata, err := createQueryMetadata(pageSize, bookmark)
	if err != nil {
		return nil, nil, err
	}

	return s.handleGetStateByRange(collection, startKey, endKey, metadata)
}

func (s *ChaincodeStub) GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string,
	pageSize int32, bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {

	collection := ""

	metadata, err := createQueryMetadata(pageSize, bookmark)
	if err != nil {
		return nil, nil, err
	}

	startKey, endKey, err := s.createRangeKeysForPartialCompositeKey(objectType, keys)
	if err != nil {
		return nil, nil, err
	}
	return s.handleGetStateByRange(collection, startKey, endKey, metadata)
}

func (s *ChaincodeStub) GetQueryResultWithPagination(query string, pageSize int32,
	bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	
	collection := ""

	metadata, err := createQueryMetadata(pageSize, bookmark)
	if err != nil {
		return nil, nil, err
	}
	return s.handleGetQueryResult(collection, query, metadata)
}

func (iter *StateQueryIterator) Next() (*queryresult.KV, error) {
	if result, err := iter.nextResult(STATE_QUERY_RESULT); err == nil {
		return result.(*queryresult.KV), err
	} else {
		return nil, err
	}
}

func (iter *HistoryQueryIterator) Next() (*queryresult.KeyModification, error) {
	if result, err := iter.nextResult(HISTORY_QUERY_RESULT); err == nil {
		return result.(*queryresult.KeyModification), err
	} else {
		return nil, err
	}
}


func (iter *CommonIterator) HasNext() bool {
	if iter.currentLoc < len(iter.response.Results) || iter.response.HasMore {
		return true
	}
	return false
}





func (iter *CommonIterator) getResultFromBytes(queryResultBytes *pb.QueryResultBytes,
	rType resultType) (queryResult, error) {

	if rType == STATE_QUERY_RESULT {
		stateQueryResult := &queryresult.KV{}
		if err := proto.Unmarshal(queryResultBytes.ResultBytes, stateQueryResult); err != nil {
			return nil, fmt.Errorf("error unmarshaling result from bytes: %s", err)
		}
		return stateQueryResult, nil

	} else if rType == HISTORY_QUERY_RESULT {
		historyQueryResult := &queryresult.KeyModification{}
		if err := proto.Unmarshal(queryResultBytes.ResultBytes, historyQueryResult); err != nil {
			return nil, err
		}
		return historyQueryResult, nil
	}
	return nil, errors.New("wrong result type")
}

func (iter *CommonIterator) fetchNextQueryResult() error {
	if response, err := iter.handler.handleQueryStateNext(iter.response.Id, iter.channelId, iter.txid); err == nil {
		iter.currentLoc = 0
		iter.response = response
		return nil
	} else {
		return err
	}
}




func (iter *CommonIterator) nextResult(rType resultType) (queryResult, error) {
	if iter.currentLoc < len(iter.response.Results) {
		
		queryResult, err := iter.getResultFromBytes(iter.response.Results[iter.currentLoc], rType)
		if err != nil {
			return nil, err
		}
		iter.currentLoc++

		if iter.currentLoc == len(iter.response.Results) && iter.response.HasMore {
			
			if err = iter.fetchNextQueryResult(); err != nil {
				return nil, err
			}
		}

		return queryResult, err
	} else if !iter.response.HasMore {
		
		return nil, errors.New("no such key")
	}

	
	
	return nil, errors.New("invalid iterator state")
}


func (iter *CommonIterator) Close() error {
	_, err := iter.handler.handleQueryStateClose(iter.response.Id, iter.channelId, iter.txid)
	return err
}


func (s *ChaincodeStub) GetArgs() [][]byte {
	return s.args
}


func (s *ChaincodeStub) GetStringArgs() []string {
	args := s.GetArgs()
	strargs := make([]string, 0, len(args))
	for _, barg := range args {
		strargs = append(strargs, string(barg))
	}
	return strargs
}


func (s *ChaincodeStub) GetFunctionAndParameters() (function string, params []string) {
	allargs := s.GetStringArgs()
	function = ""
	params = []string{}
	if len(allargs) >= 1 {
		function = allargs[0]
		params = allargs[1:]
	}
	return
}


func (s *ChaincodeStub) GetCreator() ([]byte, error) {
	return s.creator, nil
}


func (s *ChaincodeStub) GetTransient() (map[string][]byte, error) {
	return s.transient, nil
}


func (s *ChaincodeStub) GetBinding() ([]byte, error) {
	return s.binding, nil
}


func (s *ChaincodeStub) GetSignedProposal() (*pb.SignedProposal, error) {
	return s.signedProposal, nil
}


func (s *ChaincodeStub) GetArgsSlice() ([]byte, error) {
	args := s.GetArgs()
	res := []byte{}
	for _, barg := range args {
		res = append(res, barg...)
	}
	return res, nil
}


func (s *ChaincodeStub) GetTxTimestamp() (*timestamp.Timestamp, error) {
	hdr := &common.Header{}
	if err := proto.Unmarshal(s.proposal.Header, hdr); err != nil {
		return nil, fmt.Errorf("error unmarshaling Header: %s", err)
	}

	chdr := &common.ChannelHeader{}
	if err := proto.Unmarshal(hdr.ChannelHeader, chdr); err != nil {
		return nil, fmt.Errorf("error unmarshaling ChannelHeader: %s", err)
	}

	return chdr.GetTimestamp(), nil
}




func (s *ChaincodeStub) SetEvent(name string, payload []byte) error {
	if name == "" {
		return errors.New("event name can not be empty string")
	}
	s.chaincodeEvent = &pb.ChaincodeEvent{EventName: name, Payload: payload}
	return nil
}
