/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"fmt"
	"io/ioutil"
	"math"

	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	mspprotos "github.com/mcc-github/blockchain-protos-go/msp"
	ab "github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain-protos-go/orderer/etcdraft"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

const (
	
	ReadersPolicyKey = "Readers"

	
	WritersPolicyKey = "Writers"

	
	AdminsPolicyKey = "Admins"

	defaultHashingAlgorithm = bccsp.SHA256

	defaultBlockDataHashingStructureWidth = math.MaxUint32
)


type ConfigValue interface {
	
	Key() string

	
	Value() proto.Message
}


type StandardConfigValue struct {
	key   string
	value proto.Message
}


func (scv *StandardConfigValue) Key() string {
	return scv.key
}


func (scv *StandardConfigValue) Value() proto.Message {
	return scv.value
}



func ConsortiumValue(name string) *StandardConfigValue {
	return &StandardConfigValue{
		key: ConsortiumKey,
		value: &cb.Consortium{
			Name: name,
		},
	}
}



func HashingAlgorithmValue() *StandardConfigValue {
	return &StandardConfigValue{
		key: HashingAlgorithmKey,
		value: &cb.HashingAlgorithm{
			Name: defaultHashingAlgorithm,
		},
	}
}



func BlockDataHashingStructureValue() *StandardConfigValue {
	return &StandardConfigValue{
		key: BlockDataHashingStructureKey,
		value: &cb.BlockDataHashingStructure{
			Width: defaultBlockDataHashingStructureWidth,
		},
	}
}



func OrdererAddressesValue(addresses []string) *StandardConfigValue {
	return &StandardConfigValue{
		key: OrdererAddressesKey,
		value: &cb.OrdererAddresses{
			Addresses: addresses,
		},
	}
}



func ConsensusTypeValue(consensusType string, consensusMetadata []byte) *StandardConfigValue {
	return &StandardConfigValue{
		key: ConsensusTypeKey,
		value: &ab.ConsensusType{
			Type:     consensusType,
			Metadata: consensusMetadata,
		},
	}
}



func BatchSizeValue(maxMessages, absoluteMaxBytes, preferredMaxBytes uint32) *StandardConfigValue {
	return &StandardConfigValue{
		key: BatchSizeKey,
		value: &ab.BatchSize{
			MaxMessageCount:   maxMessages,
			AbsoluteMaxBytes:  absoluteMaxBytes,
			PreferredMaxBytes: preferredMaxBytes,
		},
	}
}



func BatchTimeoutValue(timeout string) *StandardConfigValue {
	return &StandardConfigValue{
		key: BatchTimeoutKey,
		value: &ab.BatchTimeout{
			Timeout: timeout,
		},
	}
}



func ChannelRestrictionsValue(maxChannelCount uint64) *StandardConfigValue {
	return &StandardConfigValue{
		key: ChannelRestrictionsKey,
		value: &ab.ChannelRestrictions{
			MaxCount: maxChannelCount,
		},
	}
}



func KafkaBrokersValue(brokers []string) *StandardConfigValue {
	return &StandardConfigValue{
		key: KafkaBrokersKey,
		value: &ab.KafkaBrokers{
			Brokers: brokers,
		},
	}
}



func MSPValue(mspDef *mspprotos.MSPConfig) *StandardConfigValue {
	return &StandardConfigValue{
		key:   MSPKey,
		value: mspDef,
	}
}



func CapabilitiesValue(capabilities map[string]bool) *StandardConfigValue {
	c := &cb.Capabilities{
		Capabilities: make(map[string]*cb.Capability),
	}

	for capability, required := range capabilities {
		if !required {
			continue
		}
		c.Capabilities[capability] = &cb.Capability{}
	}

	return &StandardConfigValue{
		key:   CapabilitiesKey,
		value: c,
	}
}



func EndpointsValue(addresses []string) *StandardConfigValue {
	return &StandardConfigValue{
		key: EndpointsKey,
		value: &cb.OrdererAddresses{
			Addresses: addresses,
		},
	}
}



func AnchorPeersValue(anchorPeers []*pb.AnchorPeer) *StandardConfigValue {
	return &StandardConfigValue{
		key:   AnchorPeersKey,
		value: &pb.AnchorPeers{AnchorPeers: anchorPeers},
	}
}



func ChannelCreationPolicyValue(policy *cb.Policy) *StandardConfigValue {
	return &StandardConfigValue{
		key:   ChannelCreationPolicyKey,
		value: policy,
	}
}



func ACLValues(acls map[string]string) *StandardConfigValue {
	a := &pb.ACLs{
		Acls: make(map[string]*pb.APIResource),
	}

	for apiResource, policyRef := range acls {
		a.Acls[apiResource] = &pb.APIResource{PolicyRef: policyRef}
	}

	return &StandardConfigValue{
		key:   ACLsKey,
		value: a,
	}
}


func ValidateCapabilities(block *cb.Block, bccsp bccsp.BCCSP) error {
	envelopeConfig, err := protoutil.ExtractEnvelope(block, 0)
	if err != nil {
		return errors.Errorf("failed to %s", err)
	}

	configEnv := &cb.ConfigEnvelope{}
	_, err = protoutil.UnmarshalEnvelopeOfType(envelopeConfig, cb.HeaderType_CONFIG, configEnv)
	if err != nil {
		return errors.Errorf("malformed configuration envelope: %s", err)
	}

	if configEnv.Config == nil {
		return errors.New("nil config envelope Config")
	}

	if configEnv.Config.ChannelGroup == nil {
		return errors.New("no channel configuration was found in the config block")
	}

	if configEnv.Config.ChannelGroup.Groups == nil {
		return errors.New("no channel configuration groups are available")
	}

	_, exists := configEnv.Config.ChannelGroup.Groups[ApplicationGroupKey]
	if !exists {
		return errors.Errorf("invalid configuration block, missing %s "+
			"configuration group", ApplicationGroupKey)
	}

	cc, err := NewChannelConfig(configEnv.Config.ChannelGroup, bccsp)
	if err != nil {
		return errors.Errorf("no valid channel configuration found due to %s", err)
	}

	
	if err := cc.Capabilities().Supported(); err != nil {
		return err
	}

	
	if err := cc.ApplicationConfig().Capabilities().Supported(); err != nil {
		return err
	}

	return nil
}


func MarshalEtcdRaftMetadata(md *etcdraft.ConfigMetadata) ([]byte, error) {
	copyMd := proto.Clone(md).(*etcdraft.ConfigMetadata)
	for _, c := range copyMd.Consenters {
		
		
		clientCert, err := ioutil.ReadFile(string(c.GetClientTlsCert()))
		if err != nil {
			return nil, fmt.Errorf("cannot load client cert for consenter %s:%d: %s", c.GetHost(), c.GetPort(), err)
		}
		c.ClientTlsCert = clientCert

		serverCert, err := ioutil.ReadFile(string(c.GetServerTlsCert()))
		if err != nil {
			return nil, fmt.Errorf("cannot load server cert for consenter %s:%d: %s", c.GetHost(), c.GetPort(), err)
		}
		c.ServerTlsCert = serverCert
	}
	return proto.Marshal(copyMd)
}
