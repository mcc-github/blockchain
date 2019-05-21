/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package broadcast

import (
	"io"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/orderer/common/msgprocessor"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("orderer.common.broadcast")




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


type Handler struct {
	SupportRegistrar ChannelSupportRegistrar
	Metrics          *Metrics
}


func (bh *Handler) Handle(srv ab.AtomicBroadcast_BroadcastServer) error {
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

		resp := bh.ProcessMessage(msg, addr)
		err = srv.Send(resp)
		if resp.Status != cb.Status_SUCCESS {
			return err
		}

		if err != nil {
			logger.Warningf("Error sending to %s: %s", addr, err)
			return err
		}
	}

}

type MetricsTracker struct {
	ValidateStartTime time.Time
	EnqueueStartTime  time.Time
	ValidateDuration  time.Duration
	ChannelID         string
	TxType            string
	Metrics           *Metrics
}

func (mt *MetricsTracker) Record(resp *ab.BroadcastResponse) {
	labels := []string{
		"status", resp.Status.String(),
		"channel", mt.ChannelID,
		"type", mt.TxType,
	}

	if mt.ValidateDuration == 0 {
		mt.EndValidate()
	}
	mt.Metrics.ValidateDuration.With(labels...).Observe(mt.ValidateDuration.Seconds())

	if mt.EnqueueStartTime != (time.Time{}) {
		enqueueDuration := time.Since(mt.EnqueueStartTime)
		mt.Metrics.EnqueueDuration.With(labels...).Observe(enqueueDuration.Seconds())
	}

	mt.Metrics.ProcessedCount.With(labels...).Add(1)
}

func (mt *MetricsTracker) BeginValidate() {
	mt.ValidateStartTime = time.Now()
}

func (mt *MetricsTracker) EndValidate() {
	mt.ValidateDuration = time.Since(mt.ValidateStartTime)
}

func (mt *MetricsTracker) BeginEnqueue() {
	mt.EnqueueStartTime = time.Now()
}


func (bh *Handler) ProcessMessage(msg *cb.Envelope, addr string) (resp *ab.BroadcastResponse) {
	tracker := &MetricsTracker{
		ChannelID: "unknown",
		TxType:    "unknown",
		Metrics:   bh.Metrics,
	}
	defer func() {
		
		
		
		tracker.Record(resp)
	}()
	tracker.BeginValidate()

	chdr, isConfig, processor, err := bh.SupportRegistrar.BroadcastChannelSupport(msg)
	if chdr != nil {
		tracker.ChannelID = chdr.ChannelId
		tracker.TxType = cb.HeaderType(chdr.Type).String()
	}
	if err != nil {
		logger.Warningf("[channel: %s] Could not get message processor for serving %s: %s", tracker.ChannelID, addr, err)
		return &ab.BroadcastResponse{Status: cb.Status_BAD_REQUEST, Info: err.Error()}
	}

	if !isConfig {
		logger.Debugf("[channel: %s] Broadcast is processing normal message from %s with txid '%s' of type %s", chdr.ChannelId, addr, chdr.TxId, cb.HeaderType_name[chdr.Type])

		configSeq, err := processor.ProcessNormalMsg(msg)
		if err != nil {
			logger.Warningf("[channel: %s] Rejecting broadcast of normal message from %s because of error: %s", chdr.ChannelId, addr, err)
			return &ab.BroadcastResponse{Status: ClassifyError(err), Info: err.Error()}
		}
		tracker.EndValidate()

		tracker.BeginEnqueue()
		if err = processor.WaitReady(); err != nil {
			logger.Warningf("[channel: %s] Rejecting broadcast of message from %s with SERVICE_UNAVAILABLE: rejected by Consenter: %s", chdr.ChannelId, addr, err)
			return &ab.BroadcastResponse{Status: cb.Status_SERVICE_UNAVAILABLE, Info: err.Error()}
		}

		err = processor.Order(msg, configSeq)
		if err != nil {
			logger.Warningf("[channel: %s] Rejecting broadcast of normal message from %s with SERVICE_UNAVAILABLE: rejected by Order: %s", chdr.ChannelId, addr, err)
			return &ab.BroadcastResponse{Status: cb.Status_SERVICE_UNAVAILABLE, Info: err.Error()}
		}
	} else { 
		logger.Debugf("[channel: %s] Broadcast is processing config update message from %s", chdr.ChannelId, addr)

		config, configSeq, err := processor.ProcessConfigUpdateMsg(msg)
		if err != nil {
			logger.Warningf("[channel: %s] Rejecting broadcast of config message from %s because of error: %s", chdr.ChannelId, addr, err)
			return &ab.BroadcastResponse{Status: ClassifyError(err), Info: err.Error()}
		}
		tracker.EndValidate()

		tracker.BeginEnqueue()
		if err = processor.WaitReady(); err != nil {
			logger.Warningf("[channel: %s] Rejecting broadcast of message from %s with SERVICE_UNAVAILABLE: rejected by Consenter: %s", chdr.ChannelId, addr, err)
			return &ab.BroadcastResponse{Status: cb.Status_SERVICE_UNAVAILABLE, Info: err.Error()}
		}

		err = processor.Configure(config, configSeq)
		if err != nil {
			logger.Warningf("[channel: %s] Rejecting broadcast of config message from %s with SERVICE_UNAVAILABLE: rejected by Configure: %s", chdr.ChannelId, addr, err)
			return &ab.BroadcastResponse{Status: cb.Status_SERVICE_UNAVAILABLE, Info: err.Error()}
		}
	}

	logger.Debugf("[channel: %s] Broadcast has successfully enqueued message of type %s from %s", chdr.ChannelId, cb.HeaderType_name[chdr.Type], addr)

	return &ab.BroadcastResponse{Status: cb.Status_SUCCESS}
}


func ClassifyError(err error) cb.Status {
	switch errors.Cause(err) {
	case msgprocessor.ErrChannelDoesNotExist:
		return cb.Status_NOT_FOUND
	case msgprocessor.ErrPermissionDenied:
		return cb.Status_FORBIDDEN
	case msgprocessor.ErrMaintenanceMode:
		return cb.Status_SERVICE_UNAVAILABLE
	default:
		return cb.Status_BAD_REQUEST
	}
}
