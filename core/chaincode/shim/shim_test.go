/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package shim

import (
	"bytes"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	mockpeer "github.com/mcc-github/blockchain/common/mocks/peer"
	"github.com/mcc-github/blockchain/common/util"
	lproto "github.com/mcc-github/blockchain/protos/ledger/queryresult"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type shimTestCC struct {
}

func (t *shimTestCC) Init(stub ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var A, B string    
	var Aval, Bval int 
	var err error

	if len(args) != 4 {
		return Error("Incorrect number of arguments. Expecting 4")
	}

	
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return Error("Expecting integer value for asset holding")
	}

	
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return Error(err.Error())
	}

	return Success(nil)
}

func (t *shimTestCC) Invoke(stub ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		
		return t.invoke(stub, args)
	} else if function == "delete" {
		
		return t.delete(stub, args)
	} else if function == "query" {
		
		return t.query(stub, args)
	} else if function == "cc2cc" {
		return t.cc2cc(stub, args)
	} else if function == "rangeq" {
		return t.rangeq(stub, args)
	} else if function == "historyq" {
		return t.historyq(stub, args)
	} else if function == "richq" {
		return t.richq(stub, args)
	} else if function == "putep" {
		return t.putEP(stub)
	} else if function == "getep" {
		return t.getEP(stub)
	}

	return Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}


func (t *shimTestCC) invoke(stub ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    
	var Aval, Bval int 
	var X int          
	var err error

	if len(args) != 3 {
		return Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	
	
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return Error("Failed to get state")
	}
	if Avalbytes == nil {
		return Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X

	
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return Error(err.Error())
	}

	return Success(nil)
}


func (t *shimTestCC) delete(stub ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	
	err := stub.DelState(A)
	if err != nil {
		return Error("Failed to delete state")
	}

	return Success(nil)
}


func (t *shimTestCC) query(stub ChaincodeStubInterface, args []string) pb.Response {
	var A string 
	var err error

	if len(args) != 1 {
		return Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return Error(jsonResp)
	}

	return Success(Avalbytes)
}


func (t *shimTestCC) cc2cc(stub ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return Error("Invalid number of args for cc2cc. expecting at least 1")
	}
	return stub.InvokeChaincode(args[0], util.ToChaincodeArgs(args...), "")
}


func (t *shimTestCC) rangeq(stub ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return Error("Incorrect number of arguments. Expecting keys for range query")
	}

	A := args[0]
	B := args[0]

	
	resultsIterator, err := stub.GetStateByRange(A, B)
	if err != nil {
		return Error(err.Error())
	}
	defer resultsIterator.Close()

	
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return Error(err.Error())
		}
		
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return Success(buffer.Bytes())
}


func (t *shimTestCC) richq(stub ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return Error("Incorrect number of arguments. Expecting keys for range query")
	}

	query := args[0]

	
	resultsIterator, err := stub.GetQueryResult(query)
	if err != nil {
		return Error(err.Error())
	}
	defer resultsIterator.Close()

	
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return Error(err.Error())
		}
		
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return Success(buffer.Bytes())
}


func (t *shimTestCC) historyq(stub ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return Error(err.Error())
		}
		
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return Success(buffer.Bytes())
}

func (t *shimTestCC) putEP(stub ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	err := stub.SetStateValidationParameter(string(args[1]), args[2])
	if err != nil {
		return Error(err.Error())
	}
	return Success(nil)
}

func (t *shimTestCC) getEP(stub ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	ep, err := stub.GetStateValidationParameter(string(args[1]))
	if err != nil {
		return Error(err.Error())
	}
	return Success(ep)
}


var mockPeerCCSupport = mockpeer.NewMockPeerSupport()

func mockChaincodeStreamGetter(name string) (PeerChaincodeStream, error) {
	return mockPeerCCSupport.GetCC(name)
}

func setupcc(name string) *mockpeer.MockCCComm {
	os.Setenv("CORE_CHAINCODE_ID_NAME", name)
	send := make(chan *pb.ChaincodeMessage)
	recv := make(chan *pb.ChaincodeMessage)
	ccSide, _ := mockPeerCCSupport.AddCC(name, recv, send)
	ccSide.SetPong(true)
	return mockPeerCCSupport.GetCCMirror(name)
}


func setuperror() chan error {
	return make(chan error)
}

func processDone(t *testing.T, done chan error, expecterr bool) {
	err := <-done
	if expecterr != (err != nil) {
		if err == nil {
			t.Fatalf("Expected error but got success")
		} else {
			t.Fatalf("Expected success but got error %s", err)
		}
	}
}



func TestInvoke(t *testing.T) {
	streamGetter = mockChaincodeStreamGetter
	cc := &shimTestCC{}
	ccname := "shimTestCC"
	peerSide := setupcc(ccname)
	defer mockPeerCCSupport.RemoveCC(ccname)
	
	go Start(cc)

	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	peerDone := make(chan struct{})
	defer close(peerDone)

	
	go func() {
		respSet := &mockpeer.MockResponseSet{
			DoneFunc:  errorFunc,
			ErrorFunc: nil,
			Responses: []*mockpeer.MockResponse{
				{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_REGISTER}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_REGISTERED}},
			},
		}
		peerSide.SetResponses(respSet)
		peerSide.SetKeepAlive(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_KEEPALIVE})
		err := peerSide.Run(peerDone)
		assert.NoError(t, err, "peer side run failed")
	}()

	
	processDone(t, done, false)

	channelId := "testchannel"

	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_READY, Txid: "1", ChannelId: channelId})

	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("init"), []byte("A"), []byte("100"), []byte("B"), []byte("200")}, Decorations: nil}
	payload := marshalOrPanic(ci)
	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Txid: "2"}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "2", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Txid: "2"}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "2", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "2", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INIT, Payload: payload, Txid: "2", ChannelId: channelId})

	
	processDone(t, done, false)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Txid: "3", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: []byte("100"), Txid: "3", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Txid: "3", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: []byte("200"), Txid: "3", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Txid: "3", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "3", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Txid: "3", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "3", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "3", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("invoke"), []byte("A"), []byte("B"), []byte("10")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "3", ChannelId: channelId})

	
	processDone(t, done, false)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Txid: "3a", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: []byte("100"), Txid: "3a", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Txid: "3a", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: []byte("200"), Txid: "3a", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Txid: "3a", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Txid: "3a", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "3a", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("invoke"), []byte("A"), []byte("B"), []byte("10")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "3a", ChannelId: channelId})

	
	processDone(t, done, false)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Txid: "3b", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Txid: "3b", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "3b", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("invoke"), []byte("A"), []byte("B"), []byte("10")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "3b", ChannelId: channelId})

	
	processDone(t, done, false)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_DEL_STATE, Txid: "4", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Txid: "4", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "4", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("delete"), []byte("A")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "4", ChannelId: channelId})

	
	processDone(t, done, false)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_DEL_STATE, Txid: "4a", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "4a", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "4a", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("delete"), []byte("A")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "4a", ChannelId: channelId})

	
	processDone(t, done, false)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "5", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("badinvoke")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "5", ChannelId: channelId})

	
	processDone(t, done, false)

	

	
	rangeQueryResponse := &pb.QueryResponse{Results: []*pb.QueryResultBytes{
		{ResultBytes: marshalOrPanic(&lproto.KV{Namespace: "getputcc", Key: "A", Value: []byte("100")})},
		{ResultBytes: marshalOrPanic(&lproto.KV{Namespace: "getputcc", Key: "B", Value: []byte("200")})}},
		HasMore: true}
	rangeQPayload := marshalOrPanic(rangeQueryResponse)

	
	rangeQueryNext := &pb.QueryResponse{Results: nil, HasMore: false}

	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE_BY_RANGE, Txid: "6", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: rangeQPayload, Txid: "6", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_NEXT, Txid: "6", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: marshalOrPanic(rangeQueryNext), Txid: "6", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_CLOSE, Txid: "6", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "6", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "6", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("rangeq"), []byte("A"), []byte("B")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "6", ChannelId: channelId})

	
	processDone(t, done, false)

	

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE_BY_RANGE, Txid: "6a", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Payload: payload, Txid: "6a", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "6a", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("rangeq"), []byte("A"), []byte("B")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "6a", ChannelId: channelId})

	
	processDone(t, done, false)

	

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE_BY_RANGE, Txid: "6b", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: rangeQPayload, Txid: "6b", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_NEXT, Txid: "6b", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Txid: "6b", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_CLOSE, Txid: "6b", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "6b", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "6b", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("rangeq"), []byte("A"), []byte("B")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "6b", ChannelId: channelId})

	
	processDone(t, done, false)

	

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE_BY_RANGE, Txid: "6c", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: rangeQPayload, Txid: "6c", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_NEXT, Txid: "6c", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Txid: "6c", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_CLOSE, Txid: "6c", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Txid: "6c", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "6c", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("rangeq"), []byte("A"), []byte("B")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "6c", ChannelId: channelId})

	
	processDone(t, done, false)

	

	
	historyQueryResponse := &pb.QueryResponse{Results: []*pb.QueryResultBytes{
		{ResultBytes: marshalOrPanic(&lproto.KeyModification{TxId: "6", Value: []byte("100")})}},
		HasMore: true}
	payload = marshalOrPanic(historyQueryResponse)

	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_HISTORY_FOR_KEY, Txid: "7", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: payload, Txid: "7", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_NEXT, Txid: "7", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: marshalOrPanic(rangeQueryNext), Txid: "7", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_CLOSE, Txid: "7", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "7", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "7", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("historyq"), []byte("A")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "7", ChannelId: channelId})

	
	processDone(t, done, false)

	

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_HISTORY_FOR_KEY, Txid: "7a", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Payload: payload, Txid: "7a", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "7a", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("historyq"), []byte("A")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "7a", ChannelId: channelId})

	
	processDone(t, done, false)

	

	
	getQRResp := &pb.QueryResponse{Results: []*pb.QueryResultBytes{
		{ResultBytes: marshalOrPanic(&lproto.KV{Namespace: "getputcc", Key: "A", Value: []byte("100")})},
		{ResultBytes: marshalOrPanic(&lproto.KV{Namespace: "getputcc", Key: "B", Value: []byte("200")})}},
		HasMore: true}
	getQRRespPayload := marshalOrPanic(getQRResp)

	
	rangeQueryNext = &pb.QueryResponse{Results: nil, HasMore: false}

	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_QUERY_RESULT, Txid: "8", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: getQRRespPayload, Txid: "8", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_NEXT, Txid: "8", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: marshalOrPanic(rangeQueryNext), Txid: "8", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_CLOSE, Txid: "8", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "8", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "8", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("richq"), []byte("A")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "8", ChannelId: channelId})

	
	processDone(t, done, false)

	

	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_QUERY_RESULT, Txid: "8a", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Payload: nil, Txid: "8a", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "8a", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("richq"), []byte("A")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "8a", ChannelId: channelId})

	
	processDone(t, done, false)
}

func TestSetKeyEP(t *testing.T) {
	streamGetter = mockChaincodeStreamGetter
	cc := &shimTestCC{}
	ccname := "shimTestCC"
	peerSide := setupcc(ccname)
	defer mockPeerCCSupport.RemoveCC(ccname)
	
	go Start(cc)

	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	peerDone := make(chan struct{})
	defer close(peerDone)

	
	go func() {
		respSet := &mockpeer.MockResponseSet{
			DoneFunc:  errorFunc,
			ErrorFunc: nil,
			Responses: []*mockpeer.MockResponse{
				{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_REGISTER}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_REGISTERED}},
			},
		}
		peerSide.SetResponses(respSet)
		peerSide.SetKeepAlive(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_KEEPALIVE})
		err := peerSide.Run(peerDone)
		assert.NoError(t, err, "peer side run failed")
	}()

	
	processDone(t, done, false)

	channelID := "testchannel"

	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_READY, Txid: "1", ChannelId: channelID})

	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("init"), []byte("A"), []byte("100"), []byte("B"), []byte("200")}, Decorations: nil}
	payload := marshalOrPanic(ci)
	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Txid: "2"}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "2", ChannelId: channelID}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Txid: "2"}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "2", ChannelId: channelID}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "2", ChannelId: channelID}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INIT, Payload: payload, Txid: "2", ChannelId: channelID})

	processDone(t, done, false)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE_METADATA, Txid: "4", ChannelId: channelID}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: nil, Txid: "4", ChannelId: channelID}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "4", ChannelId: channelID}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("putep"), []byte("A"), []byte("epA")}, Decorations: nil}
	payload = marshalOrPanic(ci)

	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "4", ChannelId: channelID})

	
	processDone(t, done, false)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE_METADATA, Txid: "5", ChannelId: channelID}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: []byte("epA"), Txid: "5", ChannelId: channelID}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "5", ChannelId: channelID}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("getep"), []byte("A")}, Decorations: nil}
	payload = marshalOrPanic(ci)

	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "5", ChannelId: channelID})

	
	processDone(t, done, false)

}

func TestCC2CC(t *testing.T) {
	streamGetter = mockChaincodeStreamGetter
	cc := &shimTestCC{}
	ccname := "shimTestCC"
	peerSide := setupcc(ccname)
	defer mockPeerCCSupport.RemoveCC(ccname)
	
	go Start(cc)

	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	peerDone := make(chan struct{})
	defer close(peerDone)

	
	go func() {
		respSet := &mockpeer.MockResponseSet{
			DoneFunc:  errorFunc,
			ErrorFunc: nil,
			Responses: []*mockpeer.MockResponse{
				{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_REGISTER}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_REGISTERED}},
			},
		}
		peerSide.SetResponses(respSet)
		peerSide.SetKeepAlive(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_KEEPALIVE})
		err := peerSide.Run(peerDone)
		assert.NoError(t, err, "peer side run failed")
	}()

	
	processDone(t, done, false)

	channelId := "testchannel"

	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_READY, Txid: "1", ChannelId: channelId})

	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("init"), []byte("A"), []byte("100"), []byte("B"), []byte("200")}, Decorations: nil}
	payload := marshalOrPanic(ci)
	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Txid: "2", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "2", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Txid: "2", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Txid: "2", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "2", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INIT, Payload: payload, Txid: "2", ChannelId: channelId})

	
	processDone(t, done, false)

	
	innerResp := marshalOrPanic(&pb.Response{Status: OK, Payload: []byte("CC2CC rocks")})
	cc2ccresp := marshalOrPanic(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: innerResp})
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INVOKE_CHAINCODE, Txid: "3", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE, Payload: cc2ccresp, Txid: "3", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "3", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("cc2cc"), []byte("othercc"), []byte("arg1"), []byte("arg2")}, Decorations: nil}
	payload = marshalOrPanic(ci)
	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "3", ChannelId: channelId})

	
	processDone(t, done, false)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: errorFunc,
		Responses: []*mockpeer.MockResponse{
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INVOKE_CHAINCODE, Txid: "4", ChannelId: channelId}, RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR, Payload: cc2ccresp, Txid: "4", ChannelId: channelId}},
			{RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Txid: "4", ChannelId: channelId}, RespMsg: nil},
		},
	}
	peerSide.SetResponses(respSet)

	peerSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION, Payload: payload, Txid: "4", ChannelId: channelId})

	
	processDone(t, done, false)
}

func TestRealPeerStream(t *testing.T) {
	os.Args = []string{"chaincode", "peer.address", "127.0.0.1:12345"}
	_, err := userChaincodeStreamGetter("fake")
	assert.Error(t, err)
}
