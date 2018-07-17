/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"time"

	"github.com/mcc-github/blockchain/common/channelconfig"
	ab "github.com/mcc-github/blockchain/protos/orderer"
)


type Orderer struct {
	
	ConsensusTypeVal string
	
	BatchSizeVal *ab.BatchSize
	
	BatchTimeoutVal time.Duration
	
	KafkaBrokersVal []string
	
	MaxChannelsCountVal uint64
	
	OrganizationsVal map[string]channelconfig.Org
	
	CapabilitiesVal channelconfig.OrdererCapabilities
}


func (scm *Orderer) ConsensusType() string {
	return scm.ConsensusTypeVal
}


func (scm *Orderer) BatchSize() *ab.BatchSize {
	return scm.BatchSizeVal
}


func (scm *Orderer) BatchTimeout() time.Duration {
	return scm.BatchTimeoutVal
}


func (scm *Orderer) KafkaBrokers() []string {
	return scm.KafkaBrokersVal
}


func (scm *Orderer) MaxChannelsCount() uint64 {
	return scm.MaxChannelsCountVal
}


func (scm *Orderer) Organizations() map[string]channelconfig.Org {
	return scm.OrganizationsVal
}


func (scm *Orderer) Capabilities() channelconfig.OrdererCapabilities {
	return scm.CapabilitiesVal
}


type OrdererCapabilities struct {
	
	SupportedErr error

	
	PredictableChannelTemplateVal bool

	
	ResubmissionVal bool

	
	ExpirationVal bool
}


func (oc *OrdererCapabilities) Supported() error {
	return oc.SupportedErr
}


func (oc *OrdererCapabilities) PredictableChannelTemplate() bool {
	return oc.PredictableChannelTemplateVal
}


func (oc *OrdererCapabilities) Resubmission() bool {
	return oc.ResubmissionVal
}



func (oc *OrdererCapabilities) ExpirationCheck() bool {
	return oc.ExpirationVal
}
