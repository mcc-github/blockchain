/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package producer

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/msp/mgmt"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

type handler struct {
	ChatStream       pb.Events_ChatServer
	interestedEvents map[string]*pb.Interest
	RemoteAddr       string
	eventProcessor   *eventProcessor
	sessionEndLock   *sync.Mutex
	sessionEndTime   time.Time
}

func newHandler(stream pb.Events_ChatServer, ep *eventProcessor) *handler {
	h := &handler{
		ChatStream:       stream,
		interestedEvents: map[string]*pb.Interest{},
		RemoteAddr:       util.ExtractRemoteAddress(stream.Context()),
		eventProcessor:   ep,
		sessionEndLock:   &sync.Mutex{},
	}
	logger.Debug("event handler created for", h.RemoteAddr)
	return h
}

func getInterestKey(interest pb.Interest) string {
	var key string
	switch interest.EventType {
	case pb.EventType_BLOCK:
		key = "/" + strconv.Itoa(int(pb.EventType_BLOCK))
	case pb.EventType_FILTEREDBLOCK:
		key = "/" + strconv.Itoa(int(pb.EventType_FILTEREDBLOCK))
	default:
		logger.Errorf("unsupported interest type: %s", interest.EventType)
	}

	return key
}

func (h *handler) register(iMsg []*pb.Interest) error {
	
	
	for _, v := range iMsg {
		if err := h.eventProcessor.registerHandler(v, h); err != nil {
			logger.Errorf("could not register %s for %s: %s", v, h.RemoteAddr, err)
			continue
		}
		h.interestedEvents[getInterestKey(*v)] = v
	}

	return nil
}

func (h *handler) deregister(iMsg []*pb.Interest) error {
	for _, v := range iMsg {
		if err := h.eventProcessor.deregisterHandler(v, h); err != nil {
			logger.Errorf("could not deregister %s for %s: %s", v, h.RemoteAddr, err)
			continue
		}
		delete(h.interestedEvents, getInterestKey(*v))
	}
	return nil
}


func (h *handler) HandleMessage(msg *pb.SignedEvent) error {
	evt, err := h.validateEventMessage(msg)
	if err != nil {
		return fmt.Errorf("event message validation failed for %s: %s", h.RemoteAddr, err)
	}

	switch evt.Event.(type) {
	case *pb.Event_Register:
		eventsObj := evt.GetRegister()
		if err := h.register(eventsObj.Events); err != nil {
			return fmt.Errorf("could not register events for %s: %s", h.RemoteAddr, err)
		}
	case *pb.Event_Unregister:
		eventsObj := evt.GetUnregister()
		if err := h.deregister(eventsObj.Events); err != nil {
			return fmt.Errorf("could not deregister events for %s: %s", h.RemoteAddr, err)
		}
	case nil:
	default:
		return fmt.Errorf("invalid type from received from %s: %T", h.RemoteAddr, evt.Event)
	}
	
	if err := h.SendMessageWithTimeout(evt, h.eventProcessor.SendTimeout); err != nil {
		return fmt.Errorf("error sending response to %s: %s", h.RemoteAddr, err)
	}

	return nil
}



func (h *handler) SendMessageWithTimeout(msg *pb.Event, timeout time.Duration) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- h.sendMessage(msg)
	}()
	t := time.NewTimer(timeout)
	select {
	case <-t.C:
		logger.Warningf("timed out sending event to %s", h.RemoteAddr)
		return fmt.Errorf("timed out sending event")
	case err := <-errChan:
		t.Stop()
		return err
	}
}


func (h *handler) sendMessage(msg *pb.Event) error {
	logger.Debug("sending event to", h.RemoteAddr)
	err := h.ChatStream.Send(msg)
	if err != nil {
		logger.Debugf("sending event failed for %s: %s", h.RemoteAddr, err)
		return fmt.Errorf("error sending message through ChatStream: %s", err)
	}
	logger.Debug("event sent successfully to", h.RemoteAddr)
	return nil
}

func (h *handler) setSessionEndTime(expiry time.Time) {
	h.sessionEndLock.Lock()
	h.sessionEndTime = expiry
	h.sessionEndLock.Unlock()
}

func (h *handler) hasSessionExpired() bool {
	h.sessionEndLock.Lock()
	defer h.sessionEndLock.Unlock()

	now := time.Now()
	if !h.sessionEndTime.IsZero() && now.After(h.sessionEndTime) {
		err := errors.Errorf("Client identity has expired for %s", h.RemoteAddr)
		logger.Warning(err.Error())
		return true
	}
	return false
}












func (h *handler) validateEventMessage(signedEvt *pb.SignedEvent) (*pb.Event, error) {
	logger.Debugf("validating for signed event %p", signedEvt)

	
	
	evt := &pb.Event{}
	err := proto.Unmarshal(signedEvt.EventBytes, evt)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling the event bytes in the SignedEvent: %s", err)
	}

	expirationTime := crypto.ExpiresAt(evt.Creator)
	if !expirationTime.IsZero() && time.Now().After(expirationTime) {
		return nil, fmt.Errorf("identity expired")
	}
	h.setSessionEndTime(expirationTime)

	if evt.GetTimestamp() != nil {
		evtTime := time.Unix(evt.GetTimestamp().Seconds, int64(evt.GetTimestamp().Nanos)).UTC()
		peerTime := time.Now()

		if math.Abs(float64(peerTime.UnixNano()-evtTime.UnixNano())) > float64(h.eventProcessor.TimeWindow.Nanoseconds()) {
			logger.Warningf("Message timestamp %s more than %s apart from current server time %s", evtTime, h.eventProcessor.TimeWindow, peerTime)
			return nil, fmt.Errorf("message timestamp out of acceptable range. must be within %s of current server time", h.eventProcessor.TimeWindow)
		}
	}

	err = h.eventProcessor.BindingInspector(h.ChatStream.Context(), evt)
	if err != nil {
		return nil, err
	}

	localMSP := mgmt.GetLocalMSP()
	principalGetter := mgmt.NewLocalMSPPrincipalGetter()

	
	principal, err := principalGetter.Get(mgmt.Members)
	if err != nil {
		return nil, fmt.Errorf("failed getting local MSP principal [member]: %s", err)
	}

	id, err := localMSP.DeserializeIdentity(evt.Creator)
	if err != nil {
		return nil, fmt.Errorf("failed deserializing event creator: %s", err)
	}

	
	err = id.SatisfiesPrincipal(principal)
	if err != nil {
		return nil, fmt.Errorf("failed verifying the creator satisfies local MSP's [member] principal: %s", err)
	}

	
	err = id.Verify(signedEvt.EventBytes, signedEvt.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed verifying the event signature: %s", err)
	}

	return evt, nil
}
