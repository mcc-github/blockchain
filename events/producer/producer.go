/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package producer

import (
	"fmt"
	"io"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/comm"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

var logger = flogging.MustGetLogger("eventhub_producer")


type EventsServer struct {
	eventProcessor *eventProcessor
}


type EventsServerConfig struct {
	BufferSize       uint
	Timeout          time.Duration
	SendTimeout      time.Duration
	TimeWindow       time.Duration
	BindingInspector comm.BindingInspector
}



var globalEventsServer *EventsServer


func NewEventsServer(config *EventsServerConfig) *EventsServer {
	if globalEventsServer != nil {
		panic("Cannot create multiple event hub servers")
	}
	eventProcessor := initializeEvents(config)
	globalEventsServer := &EventsServer{eventProcessor: eventProcessor}
	return globalEventsServer
}

func (e *EventsServer) Chat(stream pb.Events_ChatServer) error {
	handler := newHandler(stream, e.eventProcessor)
	defer e.eventProcessor.cleanupHandler(handler)
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			logger.Debug("Received EOF, ending Chat")
			return nil
		}
		if err != nil {
			err := fmt.Errorf("error during Chat, stopping handler: %s", err)
			logger.Error(err.Error())
			return err
		}
		err = handler.HandleMessage(in)
		if err != nil {
			logger.Errorf("Error handling message: %s", err)
			return err
		}
	}
}
