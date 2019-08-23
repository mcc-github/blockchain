/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mcc-github/blockchain/common/capabilities"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/pkg/errors"
)

const (
	
	OrdererGroupKey = "Orderer"
)

const (
	
	ConsensusTypeKey = "ConsensusType"

	
	BatchSizeKey = "BatchSize"

	
	BatchTimeoutKey = "BatchTimeout"

	
	ChannelRestrictionsKey = "ChannelRestrictions"

	
	KafkaBrokersKey = "KafkaBrokers"

	
	EndpointsKey = "Endpoints"
)


type OrdererProtos struct {
	ConsensusType       *ab.ConsensusType
	BatchSize           *ab.BatchSize
	BatchTimeout        *ab.BatchTimeout
	KafkaBrokers        *ab.KafkaBrokers
	ChannelRestrictions *ab.ChannelRestrictions
	Capabilities        *cb.Capabilities
}


type OrdererConfig struct {
	protos *OrdererProtos
	orgs   map[string]OrdererOrg

	batchTimeout time.Duration
}


type OrdererOrgProtos struct {
	Endpoints *cb.OrdererAddresses
}


type OrdererOrgConfig struct {
	*OrganizationConfig
	protos *OrdererOrgProtos
	name   string
}


func (oc *OrdererOrgConfig) Endpoints() []string {
	return oc.protos.Endpoints.Addresses
}


func NewOrdererOrgConfig(orgName string, orgGroup *cb.ConfigGroup, mspConfigHandler *MSPConfigHandler, channelCapabilities ChannelCapabilities) (*OrdererOrgConfig, error) {
	if len(orgGroup.Groups) > 0 {
		return nil, fmt.Errorf("OrdererOrg config does not allow sub-groups")
	}

	if !channelCapabilities.OrgSpecificOrdererEndpoints() {
		if _, ok := orgGroup.Values[EndpointsKey]; ok {
			return nil, errors.Errorf("Orderer Org %s cannot contain endpoints value until V1_4_2+ capabilities have been enabled", orgName)
		}
	}

	protos := &OrdererOrgProtos{}
	orgProtos := &OrganizationProtos{}

	if err := DeserializeProtoValuesFromGroup(orgGroup, protos, orgProtos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	ooc := &OrdererOrgConfig{
		name:   orgName,
		protos: protos,
		OrganizationConfig: &OrganizationConfig{
			name:             orgName,
			protos:           orgProtos,
			mspConfigHandler: mspConfigHandler,
		},
	}

	if err := ooc.Validate(); err != nil {
		return nil, err
	}

	return ooc, nil
}

func (ooc *OrdererOrgConfig) Validate() error {
	return ooc.OrganizationConfig.Validate()
}


func NewOrdererConfig(ordererGroup *cb.ConfigGroup, mspConfig *MSPConfigHandler, channelCapabilities ChannelCapabilities) (*OrdererConfig, error) {
	oc := &OrdererConfig{
		protos: &OrdererProtos{},
		orgs:   make(map[string]OrdererOrg),
	}

	if err := DeserializeProtoValuesFromGroup(ordererGroup, oc.protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	if err := oc.Validate(); err != nil {
		return nil, err
	}

	for orgName, orgGroup := range ordererGroup.Groups {
		var err error
		if oc.orgs[orgName], err = NewOrdererOrgConfig(orgName, orgGroup, mspConfig, channelCapabilities); err != nil {
			return nil, err
		}
	}
	return oc, nil
}


func (oc *OrdererConfig) ConsensusType() string {
	return oc.protos.ConsensusType.Type
}


func (oc *OrdererConfig) ConsensusMetadata() []byte {
	return oc.protos.ConsensusType.Metadata
}


func (oc *OrdererConfig) ConsensusState() ab.ConsensusType_State {
	return oc.protos.ConsensusType.State
}


func (oc *OrdererConfig) BatchSize() *ab.BatchSize {
	return oc.protos.BatchSize
}


func (oc *OrdererConfig) BatchTimeout() time.Duration {
	return oc.batchTimeout
}




func (oc *OrdererConfig) KafkaBrokers() []string {
	return oc.protos.KafkaBrokers.Brokers
}


func (oc *OrdererConfig) MaxChannelsCount() uint64 {
	return oc.protos.ChannelRestrictions.MaxCount
}


func (oc *OrdererConfig) Organizations() map[string]OrdererOrg {
	return oc.orgs
}


func (oc *OrdererConfig) Capabilities() OrdererCapabilities {
	return capabilities.NewOrdererProvider(oc.protos.Capabilities.Capabilities)
}

func (oc *OrdererConfig) Validate() error {
	for _, validator := range []func() error{
		oc.validateBatchSize,
		oc.validateBatchTimeout,
		oc.validateKafkaBrokers,
	} {
		if err := validator(); err != nil {
			return err
		}
	}

	return nil
}

func (oc *OrdererConfig) validateBatchSize() error {
	if oc.protos.BatchSize.MaxMessageCount == 0 {
		return fmt.Errorf("Attempted to set the batch size max message count to an invalid value: 0")
	}
	if oc.protos.BatchSize.AbsoluteMaxBytes == 0 {
		return fmt.Errorf("Attempted to set the batch size absolute max bytes to an invalid value: 0")
	}
	if oc.protos.BatchSize.PreferredMaxBytes == 0 {
		return fmt.Errorf("Attempted to set the batch size preferred max bytes to an invalid value: 0")
	}
	if oc.protos.BatchSize.PreferredMaxBytes > oc.protos.BatchSize.AbsoluteMaxBytes {
		return fmt.Errorf("Attempted to set the batch size preferred max bytes (%v) greater than the absolute max bytes (%v).", oc.protos.BatchSize.PreferredMaxBytes, oc.protos.BatchSize.AbsoluteMaxBytes)
	}
	return nil
}

func (oc *OrdererConfig) validateBatchTimeout() error {
	var err error
	oc.batchTimeout, err = time.ParseDuration(oc.protos.BatchTimeout.Timeout)
	if err != nil {
		return fmt.Errorf("Attempted to set the batch timeout to a invalid value: %s", err)
	}
	if oc.batchTimeout <= 0 {
		return fmt.Errorf("Attempted to set the batch timeout to a non-positive value: %s", oc.batchTimeout)
	}
	return nil
}

func (oc *OrdererConfig) validateKafkaBrokers() error {
	for _, broker := range oc.protos.KafkaBrokers.Brokers {
		if !brokerEntrySeemsValid(broker) {
			return fmt.Errorf("Invalid broker entry: %s", broker)
		}
	}
	return nil
}


func brokerEntrySeemsValid(broker string) bool {
	if !strings.Contains(broker, ":") {
		return false
	}

	parts := strings.Split(broker, ":")
	if len(parts) > 2 {
		return false
	}

	host := parts[0]
	port := parts[1]

	if _, err := strconv.ParseUint(port, 10, 16); err != nil {
		return false
	}

	
	
	
	
	
	
	
	
	re, _ := regexp.Compile("^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9.-]*[a-zA-Z0-9])$")
	matched := re.FindString(host)
	return len(matched) == len(host)
}
