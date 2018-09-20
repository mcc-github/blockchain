/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deliver

import (
	"context"
	"io"
	"math"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/comm"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("common/deliver")




type ChainManager interface {
	GetChain(chainID string) (Chain, bool)
}




type Chain interface {
	
	Sequence() uint64

	
	PolicyManager() policies.Manager

	
	Reader() blockledger.Reader

	
	Errored() <-chan struct{}
}





type PolicyChecker interface {
	CheckPolicy(envelope *cb.Envelope, channelID string) error
}



type PolicyCheckerFunc func(envelope *cb.Envelope, channelID string) error


func (pcf PolicyCheckerFunc) CheckPolicy(envelope *cb.Envelope, channelID string) error {
	return pcf(envelope, channelID)
}




type Inspector interface {
	Inspect(context.Context, proto.Message) error
}



type InspectorFunc func(context.Context, proto.Message) error


func (inspector InspectorFunc) Inspect(ctx context.Context, p proto.Message) error {
	return inspector(ctx, p)
}


type Handler struct {
	ChainManager     ChainManager
	TimeWindow       time.Duration
	BindingInspector Inspector
}




type Receiver interface {
	Recv() (*cb.Envelope, error)
}





type ResponseSender interface {
	SendStatusResponse(status cb.Status) error
	SendBlockResponse(block *cb.Block) error
}



type Server struct {
	Receiver
	PolicyChecker
	ResponseSender
}


func ExtractChannelHeaderCertHash(msg proto.Message) []byte {
	chdr, isChannelHeader := msg.(*cb.ChannelHeader)
	if !isChannelHeader || chdr == nil {
		return nil
	}
	return chdr.TlsCertHash
}


func NewHandler(cm ChainManager, timeWindow time.Duration, mutualTLS bool) *Handler {
	return &Handler{
		ChainManager:     cm,
		TimeWindow:       timeWindow,
		BindingInspector: InspectorFunc(comm.NewBindingInspector(mutualTLS, ExtractChannelHeaderCertHash)),
	}
}


func (h *Handler) Handle(ctx context.Context, srv *Server) error {
	addr := util.ExtractRemoteAddress(ctx)
	logger.Debugf("Starting new deliver loop for %s", addr)
	for {
		logger.Debugf("Attempting to read seek info message from %s", addr)
		envelope, err := srv.Recv()
		if err == io.EOF {
			logger.Debugf("Received EOF from %s, hangup", addr)
			return nil
		}

		if err != nil {
			logger.Warningf("Error reading from %s: %s", addr, err)
			return err
		}

		if err := h.deliverBlocks(ctx, srv, envelope); err != nil {
			return err
		}

		logger.Debugf("Waiting for new SeekInfo from %s", addr)
	}
}

func (h *Handler) deliverBlocks(ctx context.Context, srv *Server, envelope *cb.Envelope) error {
	addr := util.ExtractRemoteAddress(ctx)
	payload, err := utils.UnmarshalPayload(envelope.Payload)
	if err != nil {
		logger.Warningf("Received an envelope from %s with no payload: %s", addr, err)
		return srv.SendStatusResponse(cb.Status_BAD_REQUEST)
	}

	if payload.Header == nil {
		logger.Warningf("Malformed envelope received from %s with bad header", addr)
		return srv.SendStatusResponse(cb.Status_BAD_REQUEST)
	}

	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		logger.Warningf("Failed to unmarshal channel header from %s: %s", addr, err)
		return srv.SendStatusResponse(cb.Status_BAD_REQUEST)
	}

	err = h.validateChannelHeader(ctx, chdr)
	if err != nil {
		logger.Warningf("Rejecting deliver for %s due to envelope validation error: %s", addr, err)
		return srv.SendStatusResponse(cb.Status_BAD_REQUEST)
	}

	chain, ok := h.ChainManager.GetChain(chdr.ChannelId)
	if !ok {
		
		
		logger.Debugf("Rejecting deliver for %s because channel %s not found", addr, chdr.ChannelId)
		return srv.SendStatusResponse(cb.Status_NOT_FOUND)
	}

	erroredChan := chain.Errored()
	select {
	case <-erroredChan:
		logger.Warningf("[channel: %s] Rejecting deliver request for %s because of consenter error", chdr.ChannelId, addr)
		return srv.SendStatusResponse(cb.Status_SERVICE_UNAVAILABLE)
	default:

	}

	accessControl, err := NewSessionAC(chain, envelope, srv.PolicyChecker, chdr.ChannelId, crypto.ExpiresAt)
	if err != nil {
		logger.Warningf("[channel: %s] failed to create access control object due to %s", chdr.ChannelId, err)
		return srv.SendStatusResponse(cb.Status_BAD_REQUEST)
	}

	if err := accessControl.Evaluate(); err != nil {
		logger.Warningf("[channel: %s] Client authorization revoked for deliver request from %s: %s", chdr.ChannelId, addr, err)
		return srv.SendStatusResponse(cb.Status_FORBIDDEN)
	}

	seekInfo := &ab.SeekInfo{}
	if err = proto.Unmarshal(payload.Data, seekInfo); err != nil {
		logger.Warningf("[channel: %s] Received a signed deliver request from %s with malformed seekInfo payload: %s", chdr.ChannelId, addr, err)
		return srv.SendStatusResponse(cb.Status_BAD_REQUEST)
	}

	if seekInfo.Start == nil || seekInfo.Stop == nil {
		logger.Warningf("[channel: %s] Received seekInfo message from %s with missing start or stop %v, %v", chdr.ChannelId, addr, seekInfo.Start, seekInfo.Stop)
		return srv.SendStatusResponse(cb.Status_BAD_REQUEST)
	}

	logger.Debugf("[channel: %s] Received seekInfo (%p) %v from %s", chdr.ChannelId, seekInfo, seekInfo, addr)

	cursor, number := chain.Reader().Iterator(seekInfo.Start)
	defer cursor.Close()
	var stopNum uint64
	switch stop := seekInfo.Stop.Type.(type) {
	case *ab.SeekPosition_Oldest:
		stopNum = number
	case *ab.SeekPosition_Newest:
		stopNum = chain.Reader().Height() - 1
	case *ab.SeekPosition_Specified:
		stopNum = stop.Specified.Number
		if stopNum < number {
			logger.Warningf("[channel: %s] Received invalid seekInfo message from %s: start number %d greater than stop number %d", chdr.ChannelId, addr, number, stopNum)
			return srv.SendStatusResponse(cb.Status_BAD_REQUEST)
		}
	}

	for {
		if seekInfo.Behavior == ab.SeekInfo_FAIL_IF_NOT_READY {
			if number > chain.Reader().Height()-1 {
				return srv.SendStatusResponse(cb.Status_NOT_FOUND)
			}
		}

		var block *cb.Block
		var status cb.Status

		iterCh := make(chan struct{})
		go func() {
			block, status = cursor.Next()
			close(iterCh)
		}()

		select {
		case <-ctx.Done():
			logger.Debugf("Context canceled, aborting wait for next block")
			return errors.Wrapf(ctx.Err(), "context finished before block retrieved")
		case <-erroredChan:
			logger.Warningf("Aborting deliver for request because of background error")
			return srv.SendStatusResponse(cb.Status_SERVICE_UNAVAILABLE)
		case <-iterCh:
			
		}

		if status != cb.Status_SUCCESS {
			logger.Errorf("[channel: %s] Error reading from channel, cause was: %v", chdr.ChannelId, status)
			return srv.SendStatusResponse(status)
		}

		
		number++

		if err := accessControl.Evaluate(); err != nil {
			logger.Warningf("[channel: %s] Client authorization revoked for deliver request from %s: %s", chdr.ChannelId, addr, err)
			return srv.SendStatusResponse(cb.Status_FORBIDDEN)
		}

		logger.Debugf("[channel: %s] Delivering block for (%p) for %s", chdr.ChannelId, seekInfo, addr)

		if err := srv.SendBlockResponse(block); err != nil {
			logger.Warningf("[channel: %s] Error sending to %s: %s", chdr.ChannelId, addr, err)
			return err
		}

		if stopNum == block.Header.Number {
			break
		}
	}

	if err := srv.SendStatusResponse(cb.Status_SUCCESS); err != nil {
		logger.Warningf("[channel: %s] Error sending to %s: %s", chdr.ChannelId, addr, err)
		return err
	}

	logger.Debugf("[channel: %s] Done delivering to %s for (%p)", chdr.ChannelId, addr, seekInfo)

	return nil
}

func (h *Handler) validateChannelHeader(ctx context.Context, chdr *cb.ChannelHeader) error {
	if chdr.GetTimestamp() == nil {
		err := errors.New("channel header in envelope must contain timestamp")
		return err
	}

	envTime := time.Unix(chdr.GetTimestamp().Seconds, int64(chdr.GetTimestamp().Nanos)).UTC()
	serverTime := time.Now()

	if math.Abs(float64(serverTime.UnixNano()-envTime.UnixNano())) > float64(h.TimeWindow.Nanoseconds()) {
		err := errors.Errorf("envelope timestamp %s is more than %s apart from current server time %s", envTime, h.TimeWindow, serverTime)
		return err
	}

	err := h.BindingInspector.Inspect(ctx, chdr)
	if err != nil {
		return err
	}

	return nil
}
