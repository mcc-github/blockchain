/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"errors"
	"fmt"

	pb "github.com/mcc-github/blockchain-protos-go/peer"
)


type SendPanicFailure string

func (e SendPanicFailure) Error() string {
	return fmt.Sprintf("send failure %s", string(e))
}


type inProcStream struct {
	recv <-chan *pb.ChaincodeMessage
	send chan<- *pb.ChaincodeMessage
}

func newInProcStream(recv <-chan *pb.ChaincodeMessage, send chan<- *pb.ChaincodeMessage) *inProcStream {
	return &inProcStream{recv, send}
}

func (s *inProcStream) Send(msg *pb.ChaincodeMessage) (err error) {
	
	
	defer func() {
		if r := recover(); r != nil {
			err = SendPanicFailure(fmt.Sprintf("%s", r))
			return
		}
	}()
	s.send <- msg
	return
}

func (s *inProcStream) Recv() (*pb.ChaincodeMessage, error) {
	msg, ok := <-s.recv
	if !ok {
		return nil, errors.New("channel is closed")
	}
	return msg, nil
}

func (s *inProcStream) CloseSend() error {
	return nil
}
