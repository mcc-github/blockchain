/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"time"

	ab "github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/common/channelconfig"
)


type Orderer struct {
	
	ConsensusTypeVal string
	
	ConsensusMetadataVal []byte
	
	ConsensusTypeStateVal ab.ConsensusType_State

	
	BatchSizeVal *ab.BatchSize
	
	BatchTimeoutVal time.Duration
	
	KafkaBrokersVal []string
	
	MaxChannelsCountVal uint64
	
	OrganizationsVal map[string]channelconfig.OrdererOrg
	
	CapabilitiesVal channelconfig.OrdererCapabilities
}


func (o *Orderer) ConsensusType() string {
	return o.ConsensusTypeVal
}


func (o *Orderer) ConsensusMetadata() []byte {
	return o.ConsensusMetadataVal
}


func (o *Orderer) ConsensusState() ab.ConsensusType_State {
	return o.ConsensusTypeStateVal
}


func (o *Orderer) BatchSize() *ab.BatchSize {
	return o.BatchSizeVal
}


func (o *Orderer) BatchTimeout() time.Duration {
	return o.BatchTimeoutVal
}


func (o *Orderer) KafkaBrokers() []string {
	return o.KafkaBrokersVal
}


func (o *Orderer) MaxChannelsCount() uint64 {
	return o.MaxChannelsCountVal
}


func (o *Orderer) Organizations() map[string]channelconfig.OrdererOrg {
	return o.OrganizationsVal
}


func (o *Orderer) Capabilities() channelconfig.OrdererCapabilities {
	return o.CapabilitiesVal
}


type OrdererCapabilities struct {
	
	SupportedErr error

	
	PredictableChannelTemplateVal bool

	
	ResubmissionVal bool

	
	ExpirationVal bool

	ConsensusTypeMigrationVal bool

	UseChannelCreationPolicyAsAdminsVal bool
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


func (oc *OrdererCapabilities) ConsensusTypeMigration() bool {
	return oc.ConsensusTypeMigrationVal
}

func (oc *OrdererCapabilities) UseChannelCreationPolicyAsAdmins() bool {
	return oc.UseChannelCreationPolicyAsAdminsVal
}
