/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/



package shim

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/chaincode/shim/internal"
	peerpb "github.com/mcc-github/blockchain/protos/peer"
)

const (
	minUnicodeRuneValue   = 0            
	maxUnicodeRuneValue   = utf8.MaxRune 
	compositeKeyNamespace = "\x00"
	emptyKeySubstitute    = "\x01"
	connectTimeout        = 3 * time.Second
)



type peerStreamGetter func(name string) (PeerChaincodeStream, error)


var streamGetter peerStreamGetter


func userChaincodeStreamGetter(name string) (PeerChaincodeStream, error) {
	peerAddress := flag.String("peer.address", "", "peer address")
	flag.Parse()
	if *peerAddress == "" {
		return nil, errors.New("flag 'peer.address' must be set")
	}

	conf, err := internal.LoadConfig()
	if err != nil {
		return nil, err
	}

	conn, err := internal.NewClientConn(*peerAddress, conf.TLS, conf.KaOpts, connectTimeout)
	if err != nil {
		return nil, err
	}

	return internal.NewRegisterClient(conn)
}


func Start(cc Chaincode) error {
	chaincodename := os.Getenv("CORE_CHAINCODE_ID_NAME")
	if chaincodename == "" {
		return errors.New("'CORE_CHAINCODE_ID_NAME' must be set")
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



func StartInProc(chaincodename string, stream PeerChaincodeStream, cc Chaincode) error {
	return chatWithPeer(chaincodename, stream, cc)
}

func chatWithPeer(chaincodename string, stream PeerChaincodeStream, cc Chaincode) error {
	
	handler := newChaincodeHandler(stream, cc)
	defer stream.CloseSend()

	
	chaincodeID := &peerpb.ChaincodeID{Name: chaincodename}
	payload, err := proto.Marshal(chaincodeID)
	if err != nil {
		return fmt.Errorf("error marshalling chaincodeID during chaincode registration: %s", err)
	}

	
	if err = handler.serialSend(&peerpb.ChaincodeMessage{Type: peerpb.ChaincodeMessage_REGISTER, Payload: payload}); err != nil {
		return fmt.Errorf("error sending chaincode REGISTER: %s", err)

	}

	
	type recvMsg struct {
		msg *peerpb.ChaincodeMessage
		err error
	}
	msgAvail := make(chan *recvMsg, 1)
	errc := make(chan error)

	receiveMessage := func() {
		in, err := stream.Recv()
		msgAvail <- &recvMsg{in, err}
	}

	go receiveMessage()
	for {
		select {
		case rmsg := <-msgAvail:
			switch {
			case rmsg.err == io.EOF:
				err = fmt.Errorf("received EOF, ending chaincode stream: %s", rmsg.err)
				return err
			case rmsg.err != nil:
				err := fmt.Errorf("receive failed: %s", rmsg.err)
				return err
			case rmsg.msg == nil:
				err := errors.New("received nil message, ending chaincode stream")
				return err
			default:
				err := handler.handleMessage(rmsg.msg, errc)
				if err != nil {
					err = fmt.Errorf("error handling message: %s", err)
					return err
				}

				go receiveMessage()
			}

		case sendErr := <-errc:
			if sendErr != nil {
				err := fmt.Errorf("error sending: %s", sendErr)
				return err
			}
		}
	}
}
