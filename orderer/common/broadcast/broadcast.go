/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package broadcast

import (
	"io"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/orderer/common/msgprocessor"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/pkg/errors"
)

const pkgLogID = "orderer/common/broadcast"

var logger = flogging.MustGetLogger(pkgLogID)


type Handler interface {
	
	Handle(srv ab.AtomicBroadcast_BroadcastServer) error
}


type ChannelSupportRegistrar interface {
	
	
	
	BroadcastChannelSupport(msg *cb.Envelope) (*cb.ChannelHeader, bool, ChannelSupport, error)
}


type ChannelSupport interface {
	msgprocessor.Processor
	Consenter
}


type Consenter interface {
	
	
	Order(env *cb.Envelope, configSeq uint64) error

	
	
	Configure(config *cb.Envelope, configSeq uint64) error

	
	
	
	
	
	WaitReady() error
}

type handlerImpl struct {
	sm ChannelSupportRegistrar
}


func NewHandlerImpl(sm ChannelSupportRegistrar) Handler {
	return &handlerImpl{
		sm: sm,
	}
}


func (bh *handlerImpl) Handle(srv ab.AtomicBroadcast_BroadcastServer) error {
	addr := util.ExtractRemoteAddress(srv.Context())
	logger.Debugf("Starting new broadcast loop for %s", addr)
	for {
		msg, err := srv.Recv()
		if err == io.EOF {
			logger.Debugf("Received EOF from %s, hangup", addr)
			return nil
		}
		if err != nil {
			logger.Warningf("Error reading from %s: %s", addr, err)
			return err
		}

		chdr, isConfig, processor, err := bh.sm.BroadcastChannelSupport(msg)
		if err != nil {
			channelID := "<malformed_header>"
			if chdr != nil {
				channelID = chdr.ChannelId
			}
			logger.Warningf("[channel: %s] Could not get message processor for serving %s: %s", channelID, addr, err)
			return srv.Send(&ab.BroadcastResponse{Status: cb.Status_BAD_REQUEST, Info: err.Error()})
		}

		if err = processor.WaitReady(); err != nil {
			logger.Warningf("[channel: %s] Rejecting broadcast of message from %s with SERVICE_UNAVAILABLE: rejected by Consenter: %s", chdr.ChannelId, addr, err)
			return srv.Send(&ab.BroadcastResponse{Status: cb.Status_SERVICE_UNAVAILABLE, Info: err.Error()})
		}

		if !isConfig {
			logger.Debugf("[channel: %s] Broadcast is processing normal message from %s with txid '%s' of type %s", chdr.ChannelId, addr, chdr.TxId, cb.HeaderType_name[chdr.Type])

			configSeq, err := processor.ProcessNormalMsg(msg)
			if err != nil {
				logger.Warningf("[channel: %s] Rejecting broadcast of normal message from %s because of error: %s", chdr.ChannelId, addr, err)
				return srv.Send(&ab.BroadcastResponse{Status: ClassifyError(err), Info: err.Error()})
			}

			err = processor.Order(msg, configSeq)
			if err != nil {
				logger.Warningf("[channel: %s] Rejecting broadcast of normal message from %s with SERVICE_UNAVAILABLE: rejected by Order: %s", chdr.ChannelId, addr, err)
				return srv.Send(&ab.BroadcastResponse{Status: cb.Status_SERVICE_UNAVAILABLE, Info: err.Error()})
			}
		} else { 
			logger.Debugf("[channel: %s] Broadcast is processing config update message from %s", chdr.ChannelId, addr)

			config, configSeq, err := processor.ProcessConfigUpdateMsg(msg)
			if err != nil {
				logger.Warningf("[channel: %s] Rejecting broadcast of config message from %s because of error: %s", chdr.ChannelId, addr, err)
				return srv.Send(&ab.BroadcastResponse{Status: ClassifyError(err), Info: err.Error()})
			}

			err = processor.Configure(config, configSeq)
			if err != nil {
				logger.Warningf("[channel: %s] Rejecting broadcast of config message from %s with SERVICE_UNAVAILABLE: rejected by Configure: %s", chdr.ChannelId, addr, err)
				return srv.Send(&ab.BroadcastResponse{Status: cb.Status_SERVICE_UNAVAILABLE, Info: err.Error()})
			}
		}

		logger.Debugf("[channel: %s] Broadcast has successfully enqueued message of type %s from %s", chdr.ChannelId, cb.HeaderType_name[chdr.Type], addr)

		err = srv.Send(&ab.BroadcastResponse{Status: cb.Status_SUCCESS})
		if err != nil {
			logger.Warningf("[channel: %s] Error sending to %s: %s", chdr.ChannelId, addr, err)
			return err
		}
	}
}


func ClassifyError(err error) cb.Status {
	switch errors.Cause(err) {
	case msgprocessor.ErrChannelDoesNotExist:
		return cb.Status_NOT_FOUND
	case msgprocessor.ErrPermissionDenied:
		return cb.Status_FORBIDDEN
	default:
		return cb.Status_BAD_REQUEST
	}
}
