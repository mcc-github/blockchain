/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/deliver"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/orderer/common/broadcast"
	localconfig "github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/orderer/common/msgprocessor"
	"github.com/mcc-github/blockchain/orderer/common/multichannel"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/pkg/errors"
)

type broadcastSupport struct {
	*multichannel.Registrar
}

func (bs broadcastSupport) BroadcastChannelSupport(msg *cb.Envelope) (*cb.ChannelHeader, bool, broadcast.ChannelSupport, error) {
	return bs.Registrar.BroadcastChannelSupport(msg)
}

type deliverSupport struct {
	*multichannel.Registrar
}

func (ds deliverSupport) GetChain(chainID string) deliver.Chain {
	chain := ds.Registrar.GetChain(chainID)
	if chain == nil {
		return nil
	}
	return chain
}

type server struct {
	bh    *broadcast.Handler
	dh    *deliver.Handler
	debug *localconfig.Debug
	*multichannel.Registrar
}

type responseSender struct {
	ab.AtomicBroadcast_DeliverServer
}

func (rs *responseSender) SendStatusResponse(status cb.Status) error {
	reply := &ab.DeliverResponse{
		Type: &ab.DeliverResponse_Status{Status: status},
	}
	return rs.Send(reply)
}

func (rs *responseSender) SendBlockResponse(block *cb.Block) error {
	response := &ab.DeliverResponse{
		Type: &ab.DeliverResponse_Block{Block: block},
	}
	return rs.Send(response)
}


func NewServer(
	r *multichannel.Registrar,
	metricsProvider metrics.Provider,
	debug *localconfig.Debug,
	timeWindow time.Duration,
	mutualTLS bool,
	expirationCheckDisabled bool,
) ab.AtomicBroadcastServer {
	s := &server{
		dh: deliver.NewHandler(deliverSupport{Registrar: r}, timeWindow, mutualTLS, deliver.NewMetrics(metricsProvider), expirationCheckDisabled),
		bh: &broadcast.Handler{
			SupportRegistrar: broadcastSupport{Registrar: r},
			Metrics:          broadcast.NewMetrics(metricsProvider),
		},
		debug:     debug,
		Registrar: r,
	}
	return s
}

type msgTracer struct {
	function string
	debug    *localconfig.Debug
}

func (mt *msgTracer) trace(traceDir string, msg *cb.Envelope, err error) {
	if err != nil {
		return
	}

	now := time.Now().UnixNano()
	path := fmt.Sprintf("%s%c%d_%p.%s", traceDir, os.PathSeparator, now, msg, mt.function)
	logger.Debugf("Writing %s request trace to %s", mt.function, path)
	go func() {
		pb, err := proto.Marshal(msg)
		if err != nil {
			logger.Debugf("Error marshaling trace msg for %s: %s", path, err)
			return
		}
		err = ioutil.WriteFile(path, pb, 0660)
		if err != nil {
			logger.Debugf("Error writing trace msg for %s: %s", path, err)
		}
	}()
}

type broadcastMsgTracer struct {
	ab.AtomicBroadcast_BroadcastServer
	msgTracer
}

func (bmt *broadcastMsgTracer) Recv() (*cb.Envelope, error) {
	msg, err := bmt.AtomicBroadcast_BroadcastServer.Recv()
	if traceDir := bmt.debug.BroadcastTraceDir; traceDir != "" {
		bmt.trace(bmt.debug.BroadcastTraceDir, msg, err)
	}
	return msg, err
}

type deliverMsgTracer struct {
	deliver.Receiver
	msgTracer
}

func (dmt *deliverMsgTracer) Recv() (*cb.Envelope, error) {
	msg, err := dmt.Receiver.Recv()
	if traceDir := dmt.debug.DeliverTraceDir; traceDir != "" {
		dmt.trace(traceDir, msg, err)
	}
	return msg, err
}


func (s *server) Broadcast(srv ab.AtomicBroadcast_BroadcastServer) error {
	logger.Debugf("Starting new Broadcast handler")
	defer func() {
		if r := recover(); r != nil {
			logger.Criticalf("Broadcast client triggered panic: %s\n%s", r, debug.Stack())
		}
		logger.Debugf("Closing Broadcast stream")
	}()
	return s.bh.Handle(&broadcastMsgTracer{
		AtomicBroadcast_BroadcastServer: srv,
		msgTracer: msgTracer{
			debug:    s.debug,
			function: "Broadcast",
		},
	})
}


func (s *server) Deliver(srv ab.AtomicBroadcast_DeliverServer) error {
	logger.Debugf("Starting new Deliver handler")
	defer func() {
		if r := recover(); r != nil {
			logger.Criticalf("Deliver client triggered panic: %s\n%s", r, debug.Stack())
		}
		logger.Debugf("Closing Deliver stream")
	}()

	policyChecker := func(env *cb.Envelope, channelID string) error {
		chain := s.GetChain(channelID)
		if chain == nil {
			return errors.Errorf("channel %s not found", channelID)
		}
		
		
		sf := msgprocessor.NewSigFilter(policies.ChannelReaders, policies.ChannelOrdererReaders, chain)
		return sf.Apply(env)
	}
	deliverServer := &deliver.Server{
		PolicyChecker: deliver.PolicyCheckerFunc(policyChecker),
		Receiver: &deliverMsgTracer{
			Receiver: srv,
			msgTracer: msgTracer{
				debug:    s.debug,
				function: "Deliver",
			},
		},
		ResponseSender: &responseSender{
			AtomicBroadcast_DeliverServer: srv,
		},
	}
	return s.dh.Handle(srv.Context(), deliverServer)
}
