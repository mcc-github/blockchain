/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/



package shim

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/flogging"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)


var chaincodeLogger = logging.MustGetLogger("shim")
var logOutput = os.Stderr

var key string
var cert string

const (
	minUnicodeRuneValue   = 0            
	maxUnicodeRuneValue   = utf8.MaxRune 
	compositeKeyNamespace = "\x00"
	emptyKeySubstitute    = "\x01"
)



type ChaincodeStub struct {
	TxID           string
	ChannelId      string
	chaincodeEvent *pb.ChaincodeEvent
	args           [][]byte
	handler        *Handler
	signedProposal *pb.SignedProposal
	proposal       *pb.Proposal

	
	creator   []byte
	transient map[string][]byte
	binding   []byte

	decorations map[string][]byte
}


var peerAddress string



type peerStreamGetter func(name string) (PeerChaincodeStream, error)


var streamGetter peerStreamGetter


func userChaincodeStreamGetter(name string) (PeerChaincodeStream, error) {
	flag.StringVar(&peerAddress, "peer.address", "", "peer address")
	if viper.GetBool("peer.tls.enabled") {
		keyPath := viper.GetString("tls.client.key.path")
		certPath := viper.GetString("tls.client.cert.path")

		data, err1 := ioutil.ReadFile(keyPath)
		if err1 != nil {
			err1 = errors.Wrap(err1, fmt.Sprintf("error trying to read file content %s", keyPath))
			chaincodeLogger.Errorf("%+v", err1)
			return nil, err1
		}
		key = string(data)

		data, err1 = ioutil.ReadFile(certPath)
		if err1 != nil {
			err1 = errors.Wrap(err1, fmt.Sprintf("error trying to read file content %s", certPath))
			chaincodeLogger.Errorf("%+v", err1)
			return nil, err1
		}
		cert = string(data)
	}

	flag.Parse()

	chaincodeLogger.Debugf("Peer address: %s", getPeerAddress())

	
	clientConn, err := newPeerClientConnection()
	if err != nil {
		err = errors.Wrap(err, "error trying to connect to local peer")
		chaincodeLogger.Errorf("%+v", err)
		return nil, err
	}

	chaincodeLogger.Debugf("os.Args returns: %s", os.Args)

	chaincodeSupportClient := pb.NewChaincodeSupportClient(clientConn)

	
	stream, err := chaincodeSupportClient.Register(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("error chatting with leader at address=%s", getPeerAddress()))
	}

	return stream, nil
}


func Start(cc Chaincode) error {
	
	
	SetupChaincodeLogging()

	chaincodename := viper.GetString("chaincode.id.name")
	if chaincodename == "" {
		return errors.New("error chaincode id not provided")
	}

	err := factory.InitFactories(factory.GetDefaultOpts())
	if err != nil {
		return errors.WithMessage(err, "internal error, BCCSP could not be initialized with default options")
	}

	
	if streamGetter == nil {
		streamGetter = userChaincodeStreamGetter
	}

	stream, err := streamGetter(chaincodename)
	if err != nil {
		return err
	}

	err = chatWithPeer(chaincodename, stream, cc)

	return err
}



func IsEnabledForLogLevel(logLevel string) bool {
	lvl, _ := logging.LogLevel(logLevel)
	return chaincodeLogger.IsEnabledFor(lvl)
}




func SetupChaincodeLogging() {
	viper.SetEnvPrefix("CORE")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	
	logFormat := flogging.SetFormat(viper.GetString("chaincode.logging.format"))
	flogging.InitBackend(logFormat, logOutput)

	
	chaincodeLogLevelString := viper.GetString("chaincode.logging.level")
	if chaincodeLogLevelString == "" {
		chaincodeLogger.Infof("Chaincode log level not provided; defaulting to: %s", flogging.DefaultLevel())
		flogging.InitFromSpec(flogging.DefaultLevel())
	} else {
		_, err := LogLevel(chaincodeLogLevelString)
		if err == nil {
			flogging.InitFromSpec(chaincodeLogLevelString)
		} else {
			chaincodeLogger.Warningf("Error: '%s' for chaincode log level: %s; defaulting to %s", err, chaincodeLogLevelString, flogging.DefaultLevel())
			flogging.InitFromSpec(flogging.DefaultLevel())
		}
	}

	
	
	
	
	shimLogLevelString := viper.GetString("chaincode.logging.shim")
	if shimLogLevelString != "" {
		shimLogLevel, err := LogLevel(shimLogLevelString)
		if err == nil {
			SetLoggingLevel(shimLogLevel)
		} else {
			chaincodeLogger.Warningf("Error: %s for shim log level: %s", err, shimLogLevelString)
		}
	}

	
	
	buildLevel := viper.GetString("chaincode.buildlevel")
	chaincodeLogger.Infof("Chaincode (build level: %s) starting up ...", buildLevel)
}



func StartInProc(env []string, args []string, cc Chaincode, recv <-chan *pb.ChaincodeMessage, send chan<- *pb.ChaincodeMessage) error {
	chaincodeLogger.Debugf("in proc %v", args)

	var chaincodename string
	for _, v := range env {
		if strings.Index(v, "CORE_CHAINCODE_ID_NAME=") == 0 {
			p := strings.SplitAfter(v, "CORE_CHAINCODE_ID_NAME=")
			chaincodename = p[1]
			break
		}
	}
	if chaincodename == "" {
		return errors.New("error chaincode id not provided")
	}

	stream := newInProcStream(recv, send)
	chaincodeLogger.Debugf("starting chat with peer using name=%s", chaincodename)
	err := chatWithPeer(chaincodename, stream, cc)
	return err
}

func getPeerAddress() string {
	if peerAddress != "" {
		return peerAddress
	}

	if peerAddress = viper.GetString("peer.address"); peerAddress == "" {
		chaincodeLogger.Fatalf("peer.address not configured, can't connect to peer")
	}

	return peerAddress
}

func newPeerClientConnection() (*grpc.ClientConn, error) {
	var peerAddress = getPeerAddress()
	
	kaOpts := &comm.KeepaliveOptions{
		ClientInterval: time.Duration(1) * time.Minute,
		ClientTimeout:  time.Duration(20) * time.Second,
	}
	if viper.GetBool("peer.tls.enabled") {
		return comm.NewClientConnectionWithAddress(peerAddress, true, true,
			comm.InitTLSForShim(key, cert), kaOpts)
	}
	return comm.NewClientConnectionWithAddress(peerAddress, true, false, nil, kaOpts)
}

func chatWithPeer(chaincodename string, stream PeerChaincodeStream, cc Chaincode) error {

	
	handler := newChaincodeHandler(stream, cc)

	defer stream.CloseSend()
	
	chaincodeID := &pb.ChaincodeID{Name: chaincodename}
	payload, err := proto.Marshal(chaincodeID)
	if err != nil {
		return errors.Wrap(err, "error marshalling chaincodeID during chaincode registration")
	}
	
	chaincodeLogger.Debugf("Registering.. sending %s", pb.ChaincodeMessage_REGISTER)
	if err = handler.serialSend(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_REGISTER, Payload: payload}); err != nil {
		return errors.WithMessage(err, "error sending chaincode REGISTER")
	}
	waitc := make(chan struct{})
	errc := make(chan error)
	go func() {
		defer close(waitc)
		msgAvail := make(chan *pb.ChaincodeMessage)
		var in *pb.ChaincodeMessage
		recv := true
		for {
			in = nil
			err = nil
			if recv {
				recv = false
				go func() {
					var in2 *pb.ChaincodeMessage
					in2, err = stream.Recv()
					errc <- err
					msgAvail <- in2
				}()
			}
			select {
			case sendErr := <-errc:
				
				if sendErr == nil {
					continue
				}
				
				err = errors.Wrap(sendErr, "error sending")
				return
			case in = <-msgAvail:
				if err == io.EOF {
					err = errors.Wrapf(err, "received EOF, ending chaincode stream")
					chaincodeLogger.Debugf("%+v", err)
					return
				} else if err != nil {
					chaincodeLogger.Errorf("Received error from server, ending chaincode stream: %+v", err)
					return
				} else if in == nil {
					err = errors.New("received nil message, ending chaincode stream")
					chaincodeLogger.Debugf("%+v", err)
					return
				}
				chaincodeLogger.Debugf("[%s]Received message %s from peer", shorttxid(in.Txid), in.Type)
				recv = true
			}

			err = handler.handleMessage(in, errc)
			if err != nil {
				err = errors.WithMessage(err, "error handling message")
				return
			}

			
			if in.Type == pb.ChaincodeMessage_KEEPALIVE {
				chaincodeLogger.Debug("Sending KEEPALIVE response")
				
				handler.serialSendAsync(in, nil)
			}
		}
	}()
	<-waitc
	return err
}




func (stub *ChaincodeStub) init(handler *Handler, channelId string, txid string, input *pb.ChaincodeInput, signedProposal *pb.SignedProposal) error {
	stub.TxID = txid
	stub.ChannelId = channelId
	stub.args = input.Args
	stub.handler = handler
	stub.signedProposal = signedProposal
	stub.decorations = input.Decorations

	
	
	
	if signedProposal != nil {
		var err error

		stub.proposal, err = utils.GetProposal(signedProposal.ProposalBytes)
		if err != nil {
			return errors.WithMessage(err, "failed extracting signedProposal from signed signedProposal")
		}

		
		stub.creator, stub.transient, err = utils.GetChaincodeProposalContext(stub.proposal)
		if err != nil {
			return errors.WithMessage(err, "failed extracting signedProposal fields")
		}

		stub.binding, err = utils.ComputeProposalBinding(stub.proposal)
		if err != nil {
			return errors.WithMessage(err, "failed computing binding from signedProposal")
		}
	}

	return nil
}


func (stub *ChaincodeStub) GetTxID() string {
	return stub.TxID
}


func (stub *ChaincodeStub) GetChannelID() string {
	return stub.ChannelId
}

func (stub *ChaincodeStub) GetDecorations() map[string][]byte {
	return stub.decorations
}




func (stub *ChaincodeStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response {
	
	if channel != "" {
		chaincodeName = chaincodeName + "/" + channel
	}
	return stub.handler.handleInvokeChaincode(chaincodeName, args, stub.ChannelId, stub.TxID)
}




func (stub *ChaincodeStub) GetState(key string) ([]byte, error) {
	
	collection := ""
	return stub.handler.handleGetState(collection, key, stub.ChannelId, stub.TxID)
}


func (stub *ChaincodeStub) PutState(key string, value []byte) error {
	if key == "" {
		return errors.New("key must not be an empty string")
	}
	
	collection := ""
	return stub.handler.handlePutState(collection, key, value, stub.ChannelId, stub.TxID)
}


func (stub *ChaincodeStub) GetQueryResult(query string) (StateQueryIteratorInterface, error) {
	
	collection := ""
	response, err := stub.handler.handleGetQueryResult(collection, query, stub.ChannelId, stub.TxID)
	if err != nil {
		return nil, err
	}
	return &StateQueryIterator{CommonIterator: &CommonIterator{stub.handler, stub.ChannelId, stub.TxID, response, 0}}, nil
}


func (stub *ChaincodeStub) DelState(key string) error {
	
	collection := ""
	return stub.handler.handleDelState(collection, key, stub.ChannelId, stub.TxID)
}




func (stub *ChaincodeStub) GetPrivateData(collection string, key string) ([]byte, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must not be an empty string")
	}
	return stub.handler.handleGetState(collection, key, stub.ChannelId, stub.TxID)
}


func (stub *ChaincodeStub) PutPrivateData(collection string, key string, value []byte) error {
	if collection == "" {
		return fmt.Errorf("collection must not be an empty string")
	}
	if key == "" {
		return fmt.Errorf("key must not be an empty string")
	}
	return stub.handler.handlePutState(collection, key, value, stub.ChannelId, stub.TxID)
}


func (stub *ChaincodeStub) DelPrivateData(collection string, key string) error {
	if collection == "" {
		return fmt.Errorf("collection must not be an empty string")
	}
	return stub.handler.handleDelState(collection, key, stub.ChannelId, stub.TxID)
}


func (stub *ChaincodeStub) GetPrivateDataByRange(collection, startKey, endKey string) (StateQueryIteratorInterface, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must not be an empty string")
	}
	if startKey == "" {
		startKey = emptyKeySubstitute
	}
	if err := validateSimpleKeys(startKey, endKey); err != nil {
		return nil, err
	}
	return stub.handleGetStateByRange(collection, startKey, endKey)
}


func (stub *ChaincodeStub) GetPrivateDataByPartialCompositeKey(collection, objectType string, attributes []string) (StateQueryIteratorInterface, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must not be an empty string")
	}
	if partialCompositeKey, err := stub.CreateCompositeKey(objectType, attributes); err == nil {
		return stub.handleGetStateByRange(collection, partialCompositeKey, partialCompositeKey+string(maxUnicodeRuneValue))
	} else {
		return nil, err
	}
}


func (stub *ChaincodeStub) GetPrivateDataQueryResult(collection, query string) (StateQueryIteratorInterface, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must not be an empty string")
	}
	response, err := stub.handler.handleGetQueryResult(collection, query, stub.ChannelId, stub.TxID)
	if err != nil {
		return nil, err
	}
	return &StateQueryIterator{CommonIterator: &CommonIterator{
		handler:   stub.handler,
		channelId: stub.ChannelId,
		txid:      stub.TxID,
		response:  response}}, nil
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

type resultType uint8

const (
	STATE_QUERY_RESULT resultType = iota + 1
	HISTORY_QUERY_RESULT
)

func (stub *ChaincodeStub) handleGetStateByRange(collection, startKey, endKey string) (StateQueryIteratorInterface, error) {
	response, err := stub.handler.handleGetStateByRange(collection, startKey, endKey, stub.ChannelId, stub.TxID)
	if err != nil {
		return nil, err
	}
	return &StateQueryIterator{CommonIterator: &CommonIterator{stub.handler, stub.ChannelId, stub.TxID, response, 0}}, nil
}


func (stub *ChaincodeStub) GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error) {
	if startKey == "" {
		startKey = emptyKeySubstitute
	}
	if err := validateSimpleKeys(startKey, endKey); err != nil {
		return nil, err
	}
	collection := ""
	return stub.handleGetStateByRange(collection, startKey, endKey)
}


func (stub *ChaincodeStub) GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error) {
	response, err := stub.handler.handleGetHistoryForKey(key, stub.ChannelId, stub.TxID)
	if err != nil {
		return nil, err
	}
	return &HistoryQueryIterator{CommonIterator: &CommonIterator{stub.handler, stub.ChannelId, stub.TxID, response, 0}}, nil
}


func (stub *ChaincodeStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return createCompositeKey(objectType, attributes)
}


func (stub *ChaincodeStub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return splitCompositeKey(compositeKey)
}

func createCompositeKey(objectType string, attributes []string) (string, error) {
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
		return errors.Errorf("not a valid utf8 string: [%x]", str)
	}
	for index, runeValue := range str {
		if runeValue == minUnicodeRuneValue || runeValue == maxUnicodeRuneValue {
			return errors.Errorf(`input contain unicode %#U starting at position [%d]. %#U and %#U are not allowed in the input attribute of a composite key`,
				runeValue, index, minUnicodeRuneValue, maxUnicodeRuneValue)
		}
	}
	return nil
}





func validateSimpleKeys(simpleKeys ...string) error {
	for _, key := range simpleKeys {
		if len(key) > 0 && key[0] == compositeKeyNamespace[0] {
			return errors.Errorf(`first character of the key [%s] contains a null character which is not allowed`, key)
		}
	}
	return nil
}







func (stub *ChaincodeStub) GetStateByPartialCompositeKey(objectType string, attributes []string) (StateQueryIteratorInterface, error) {
	collection := ""
	if partialCompositeKey, err := stub.CreateCompositeKey(objectType, attributes); err == nil {
		return stub.handleGetStateByRange(collection, partialCompositeKey, partialCompositeKey+string(maxUnicodeRuneValue))
	} else {
		return nil, err
	}
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
	rType resultType) (commonledger.QueryResult, error) {

	if rType == STATE_QUERY_RESULT {
		stateQueryResult := &queryresult.KV{}
		if err := proto.Unmarshal(queryResultBytes.ResultBytes, stateQueryResult); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling result from bytes")
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




func (iter *CommonIterator) nextResult(rType resultType) (commonledger.QueryResult, error) {
	if iter.currentLoc < len(iter.response.Results) {
		
		queryResult, err := iter.getResultFromBytes(iter.response.Results[iter.currentLoc], rType)
		if err != nil {
			chaincodeLogger.Errorf("Failed to decode query results: %+v", err)
			return nil, err
		}
		iter.currentLoc++

		if iter.currentLoc == len(iter.response.Results) && iter.response.HasMore {
			
			if err = iter.fetchNextQueryResult(); err != nil {
				chaincodeLogger.Errorf("Failed to fetch next results: %+v", err)
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


func (stub *ChaincodeStub) GetArgs() [][]byte {
	return stub.args
}


func (stub *ChaincodeStub) GetStringArgs() []string {
	args := stub.GetArgs()
	strargs := make([]string, 0, len(args))
	for _, barg := range args {
		strargs = append(strargs, string(barg))
	}
	return strargs
}


func (stub *ChaincodeStub) GetFunctionAndParameters() (function string, params []string) {
	allargs := stub.GetStringArgs()
	function = ""
	params = []string{}
	if len(allargs) >= 1 {
		function = allargs[0]
		params = allargs[1:]
	}
	return
}


func (stub *ChaincodeStub) GetCreator() ([]byte, error) {
	return stub.creator, nil
}


func (stub *ChaincodeStub) GetTransient() (map[string][]byte, error) {
	return stub.transient, nil
}


func (stub *ChaincodeStub) GetBinding() ([]byte, error) {
	return stub.binding, nil
}


func (stub *ChaincodeStub) GetSignedProposal() (*pb.SignedProposal, error) {
	return stub.signedProposal, nil
}


func (stub *ChaincodeStub) GetArgsSlice() ([]byte, error) {
	args := stub.GetArgs()
	res := []byte{}
	for _, barg := range args {
		res = append(res, barg...)
	}
	return res, nil
}


func (stub *ChaincodeStub) GetTxTimestamp() (*timestamp.Timestamp, error) {
	hdr, err := utils.GetHeader(stub.proposal.Header)
	if err != nil {
		return nil, err
	}
	chdr, err := utils.UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		return nil, err
	}

	return chdr.GetTimestamp(), nil
}




func (stub *ChaincodeStub) SetEvent(name string, payload []byte) error {
	if name == "" {
		return errors.New("event name can not be nil string")
	}
	stub.chaincodeEvent = &pb.ChaincodeEvent{EventName: name, Payload: payload}
	return nil
}

































type LoggingLevel logging.Level


const (
	LogDebug    = LoggingLevel(logging.DEBUG)
	LogInfo     = LoggingLevel(logging.INFO)
	LogNotice   = LoggingLevel(logging.NOTICE)
	LogWarning  = LoggingLevel(logging.WARNING)
	LogError    = LoggingLevel(logging.ERROR)
	LogCritical = LoggingLevel(logging.CRITICAL)
)

var shimLoggingLevel = LogInfo 



func SetLoggingLevel(level LoggingLevel) {
	shimLoggingLevel = level
	logging.SetLevel(logging.Level(level), "shim")
}




func LogLevel(levelString string) (LoggingLevel, error) {
	l, err := logging.LogLevel(levelString)
	level := LoggingLevel(l)
	if err != nil {
		level = LogError
	}
	return level, err
}





type ChaincodeLogger struct {
	logger *logging.Logger
}






func NewLogger(name string) *ChaincodeLogger {
	return &ChaincodeLogger{logging.MustGetLogger(name)}
}




func (c *ChaincodeLogger) SetLevel(level LoggingLevel) {
	logging.SetLevel(logging.Level(level), c.logger.Module)
}



func (c *ChaincodeLogger) IsEnabledFor(level LoggingLevel) bool {
	return c.logger.IsEnabledFor(logging.Level(level))
}



func (c *ChaincodeLogger) Debug(args ...interface{}) {
	c.logger.Debug(args...)
}



func (c *ChaincodeLogger) Info(args ...interface{}) {
	c.logger.Info(args...)
}



func (c *ChaincodeLogger) Notice(args ...interface{}) {
	c.logger.Notice(args...)
}



func (c *ChaincodeLogger) Warning(args ...interface{}) {
	c.logger.Warning(args...)
}



func (c *ChaincodeLogger) Error(args ...interface{}) {
	c.logger.Error(args...)
}


func (c *ChaincodeLogger) Critical(args ...interface{}) {
	c.logger.Critical(args...)
}



func (c *ChaincodeLogger) Debugf(format string, args ...interface{}) {
	c.logger.Debugf(format, args...)
}



func (c *ChaincodeLogger) Infof(format string, args ...interface{}) {
	c.logger.Infof(format, args...)
}



func (c *ChaincodeLogger) Noticef(format string, args ...interface{}) {
	c.logger.Noticef(format, args...)
}



func (c *ChaincodeLogger) Warningf(format string, args ...interface{}) {
	c.logger.Warningf(format, args...)
}



func (c *ChaincodeLogger) Errorf(format string, args ...interface{}) {
	c.logger.Errorf(format, args...)
}


func (c *ChaincodeLogger) Criticalf(format string, args ...interface{}) {
	c.logger.Criticalf(format, args...)
}
