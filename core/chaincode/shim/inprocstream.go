/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shim

import (
	"fmt"

	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
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
	err = nil

	
	
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
