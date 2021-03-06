/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	mspprotos "github.com/mcc-github/blockchain-protos-go/msp"
	ab "github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain-protos-go/orderer/etcdraft"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/common/capabilities"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)







func basicTest(t *testing.T, sv *StandardConfigValue) {
	assert.NotNil(t, sv)
	assert.NotEmpty(t, sv.Key())
	assert.NotNil(t, sv.Value())
}

func TestUtilsBasic(t *testing.T) {
	basicTest(t, ConsortiumValue("foo"))
	basicTest(t, HashingAlgorithmValue())
	basicTest(t, BlockDataHashingStructureValue())
	basicTest(t, OrdererAddressesValue([]string{"foo:1", "bar:2"}))
	basicTest(t, ConsensusTypeValue("foo", []byte("bar")))
	basicTest(t, BatchSizeValue(1, 2, 3))
	basicTest(t, BatchTimeoutValue("1s"))
	basicTest(t, ChannelRestrictionsValue(7))
	basicTest(t, KafkaBrokersValue([]string{"foo:1", "bar:2"}))
	basicTest(t, MSPValue(&mspprotos.MSPConfig{}))
	basicTest(t, CapabilitiesValue(map[string]bool{"foo": true, "bar": false}))
	basicTest(t, AnchorPeersValue([]*pb.AnchorPeer{{}, {}}))
	basicTest(t, ChannelCreationPolicyValue(&cb.Policy{}))
	basicTest(t, ACLValues(map[string]string{"foo": "fooval", "bar": "barval"}))
}


func createCfgBlockWithSupportedCapabilities(t *testing.T) *cb.Block {
	
	config := &cb.Config{
		Sequence:     0,
		ChannelGroup: protoutil.NewConfigGroup(),
	}

	
	config.ChannelGroup.Version = 0
	config.ChannelGroup.ModPolicy = AdminsPolicyKey
	config.ChannelGroup.Values[BlockDataHashingStructureKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(&cb.BlockDataHashingStructure{
			Width: defaultBlockDataHashingStructureWidth,
		}),
		ModPolicy: AdminsPolicyKey,
	}
	topCapabilities := make(map[string]bool)
	topCapabilities[capabilities.ChannelV1_1] = true
	config.ChannelGroup.Values[CapabilitiesKey] = &cb.ConfigValue{
		Value:     protoutil.MarshalOrPanic(CapabilitiesValue(topCapabilities).Value()),
		ModPolicy: AdminsPolicyKey,
	}
	config.ChannelGroup.Values[ConsortiumKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(&cb.Consortium{
			Name: "testConsortium",
		}),
		ModPolicy: AdminsPolicyKey,
	}
	config.ChannelGroup.Values[HashingAlgorithmKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(&cb.HashingAlgorithm{
			Name: defaultHashingAlgorithm,
		}),
		ModPolicy: AdminsPolicyKey,
	}
	config.ChannelGroup.Values[OrdererAddressesKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(&cb.OrdererAddresses{
			Addresses: []string{"orderer.example.com"},
		}),
		ModPolicy: AdminsPolicyKey,
	}

	
	config.ChannelGroup.Groups[ApplicationGroupKey] = protoutil.NewConfigGroup()
	config.ChannelGroup.Groups[ApplicationGroupKey].Version = 0
	config.ChannelGroup.Groups[ApplicationGroupKey].ModPolicy = AdminsPolicyKey
	config.ChannelGroup.Groups[ApplicationGroupKey].Policies[ReadersPolicyKey] = &cb.ConfigPolicy{}
	config.ChannelGroup.Groups[ApplicationGroupKey].Policies[WritersPolicyKey] = &cb.ConfigPolicy{}
	config.ChannelGroup.Groups[ApplicationGroupKey].Policies[AdminsPolicyKey] = &cb.ConfigPolicy{}
	appCapabilities := make(map[string]bool)
	appCapabilities[capabilities.ApplicationV1_1] = true
	config.ChannelGroup.Groups[ApplicationGroupKey].Values[CapabilitiesKey] = &cb.ConfigValue{
		Value:     protoutil.MarshalOrPanic(CapabilitiesValue(appCapabilities).Value()),
		ModPolicy: AdminsPolicyKey,
	}

	
	config.ChannelGroup.Groups[OrdererGroupKey] = protoutil.NewConfigGroup()
	config.ChannelGroup.Groups[OrdererGroupKey].Version = 0
	config.ChannelGroup.Groups[OrdererGroupKey].ModPolicy = AdminsPolicyKey
	config.ChannelGroup.Groups[OrdererGroupKey].Policies[ReadersPolicyKey] = &cb.ConfigPolicy{}
	config.ChannelGroup.Groups[OrdererGroupKey].Policies[WritersPolicyKey] = &cb.ConfigPolicy{}
	config.ChannelGroup.Groups[OrdererGroupKey].Policies[AdminsPolicyKey] = &cb.ConfigPolicy{}
	config.ChannelGroup.Groups[OrdererGroupKey].Values[BatchSizeKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(
			&ab.BatchSize{
				MaxMessageCount:   65535,
				AbsoluteMaxBytes:  1024000000,
				PreferredMaxBytes: 1024000000,
			}),
		ModPolicy: AdminsPolicyKey,
	}
	config.ChannelGroup.Groups[OrdererGroupKey].Values[BatchTimeoutKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(
			&ab.BatchTimeout{
				Timeout: "2s",
			}),
		ModPolicy: AdminsPolicyKey,
	}
	ordererCapabilities := make(map[string]bool)
	ordererCapabilities[capabilities.OrdererV1_1] = true
	config.ChannelGroup.Groups[OrdererGroupKey].Values[CapabilitiesKey] = &cb.ConfigValue{
		Value:     protoutil.MarshalOrPanic(CapabilitiesValue(ordererCapabilities).Value()),
		ModPolicy: AdminsPolicyKey,
	}
	config.ChannelGroup.Groups[OrdererGroupKey].Values[ConsensusTypeKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(
			&ab.ConsensusType{
				Type: "solo",
			}),
		ModPolicy: AdminsPolicyKey,
	}

	env := &cb.Envelope{
		Payload: protoutil.MarshalOrPanic(&cb.Payload{
			Header: &cb.Header{
				ChannelHeader: protoutil.MarshalOrPanic(&cb.ChannelHeader{
					ChannelId: "testChain",
					Type:      int32(cb.HeaderType_CONFIG),
				}),
			},
			Data: protoutil.MarshalOrPanic(&cb.ConfigEnvelope{
				Config: config,
			}),
		}),
	}
	configBlock := &cb.Block{
		Data: &cb.BlockData{
			Data: [][]byte{[]byte(protoutil.MarshalOrPanic(env))},
		},
	}
	return configBlock
}


func createCfgBlockWithUnsupportedCapabilities(t *testing.T) *cb.Block {
	
	config := &cb.Config{
		Sequence:     0,
		ChannelGroup: protoutil.NewConfigGroup(),
	}

	
	config.ChannelGroup.Version = 0
	config.ChannelGroup.ModPolicy = AdminsPolicyKey
	config.ChannelGroup.Values[BlockDataHashingStructureKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(&cb.BlockDataHashingStructure{
			Width: defaultBlockDataHashingStructureWidth,
		}),
		ModPolicy: AdminsPolicyKey,
	}
	topCapabilities := make(map[string]bool)
	topCapabilities["INCOMPATIBLE_CAPABILITIES"] = true
	config.ChannelGroup.Values[CapabilitiesKey] = &cb.ConfigValue{
		Value:     protoutil.MarshalOrPanic(CapabilitiesValue(topCapabilities).Value()),
		ModPolicy: AdminsPolicyKey,
	}
	config.ChannelGroup.Values[ConsortiumKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(&cb.Consortium{
			Name: "testConsortium",
		}),
		ModPolicy: AdminsPolicyKey,
	}
	config.ChannelGroup.Values[HashingAlgorithmKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(&cb.HashingAlgorithm{
			Name: defaultHashingAlgorithm,
		}),
		ModPolicy: AdminsPolicyKey,
	}
	config.ChannelGroup.Values[OrdererAddressesKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(&cb.OrdererAddresses{
			Addresses: []string{"orderer.example.com"},
		}),
		ModPolicy: AdminsPolicyKey,
	}

	
	config.ChannelGroup.Groups[ApplicationGroupKey] = protoutil.NewConfigGroup()
	config.ChannelGroup.Groups[ApplicationGroupKey].Version = 0
	config.ChannelGroup.Groups[ApplicationGroupKey].ModPolicy = AdminsPolicyKey
	config.ChannelGroup.Groups[ApplicationGroupKey].Policies[ReadersPolicyKey] = &cb.ConfigPolicy{}
	config.ChannelGroup.Groups[ApplicationGroupKey].Policies[WritersPolicyKey] = &cb.ConfigPolicy{}
	config.ChannelGroup.Groups[ApplicationGroupKey].Policies[AdminsPolicyKey] = &cb.ConfigPolicy{}
	appCapabilities := make(map[string]bool)
	appCapabilities["INCOMPATIBLE_CAPABILITIES"] = true
	config.ChannelGroup.Groups[ApplicationGroupKey].Values[CapabilitiesKey] = &cb.ConfigValue{
		Value:     protoutil.MarshalOrPanic(CapabilitiesValue(appCapabilities).Value()),
		ModPolicy: AdminsPolicyKey,
	}

	
	config.ChannelGroup.Groups[OrdererGroupKey] = protoutil.NewConfigGroup()
	config.ChannelGroup.Groups[OrdererGroupKey].Version = 0
	config.ChannelGroup.Groups[OrdererGroupKey].ModPolicy = AdminsPolicyKey
	config.ChannelGroup.Groups[OrdererGroupKey].Policies[ReadersPolicyKey] = &cb.ConfigPolicy{}
	config.ChannelGroup.Groups[OrdererGroupKey].Policies[WritersPolicyKey] = &cb.ConfigPolicy{}
	config.ChannelGroup.Groups[OrdererGroupKey].Policies[AdminsPolicyKey] = &cb.ConfigPolicy{}
	config.ChannelGroup.Groups[OrdererGroupKey].Values[BatchSizeKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(
			&ab.BatchSize{
				MaxMessageCount:   65535,
				AbsoluteMaxBytes:  1024000000,
				PreferredMaxBytes: 1024000000,
			}),
		ModPolicy: AdminsPolicyKey,
	}
	config.ChannelGroup.Groups[OrdererGroupKey].Values[BatchTimeoutKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(
			&ab.BatchTimeout{
				Timeout: "2s",
			}),
		ModPolicy: AdminsPolicyKey,
	}
	ordererCapabilities := make(map[string]bool)
	ordererCapabilities["INCOMPATIBLE_CAPABILITIES"] = true
	config.ChannelGroup.Groups[OrdererGroupKey].Values[CapabilitiesKey] = &cb.ConfigValue{
		Value:     protoutil.MarshalOrPanic(CapabilitiesValue(ordererCapabilities).Value()),
		ModPolicy: AdminsPolicyKey,
	}
	config.ChannelGroup.Groups[OrdererGroupKey].Values[ConsensusTypeKey] = &cb.ConfigValue{
		Value: protoutil.MarshalOrPanic(
			&ab.ConsensusType{
				Type: "solo",
			}),
		ModPolicy: AdminsPolicyKey,
	}

	env := &cb.Envelope{
		Payload: protoutil.MarshalOrPanic(&cb.Payload{
			Header: &cb.Header{
				ChannelHeader: protoutil.MarshalOrPanic(&cb.ChannelHeader{
					ChannelId: "testChain",
					Type:      int32(cb.HeaderType_CONFIG),
				}),
			},
			Data: protoutil.MarshalOrPanic(&cb.ConfigEnvelope{
				Config: config,
			}),
		}),
	}
	configBlock := &cb.Block{
		Data: &cb.BlockData{
			Data: [][]byte{[]byte(protoutil.MarshalOrPanic(env))},
		},
	}
	return configBlock
}

func TestValidateCapabilities(t *testing.T) {
	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)

	
	cfgBlock := createCfgBlockWithSupportedCapabilities(t)
	err = ValidateCapabilities(cfgBlock, cryptoProvider)
	assert.NoError(t, err)

	
	cfgBlock = createCfgBlockWithUnsupportedCapabilities(t)
	err = ValidateCapabilities(cfgBlock, cryptoProvider)
	assert.EqualError(t, err, "Channel capability INCOMPATIBLE_CAPABILITIES is required but not supported")
}

func TestMarshalEtcdRaftMetadata(t *testing.T) {
	md := &etcdraft.ConfigMetadata{
		Consenters: []*etcdraft.Consenter{
			{
				Host:          "node-1.example.com",
				Port:          7050,
				ClientTlsCert: []byte("testdata/tls-client-1.pem"),
				ServerTlsCert: []byte("testdata/tls-server-1.pem"),
			},
			{
				Host:          "node-2.example.com",
				Port:          7050,
				ClientTlsCert: []byte("testdata/tls-client-2.pem"),
				ServerTlsCert: []byte("testdata/tls-server-2.pem"),
			},
			{
				Host:          "node-3.example.com",
				Port:          7050,
				ClientTlsCert: []byte("testdata/tls-client-3.pem"),
				ServerTlsCert: []byte("testdata/tls-server-3.pem"),
			},
		},
	}
	packed, err := MarshalEtcdRaftMetadata(md)
	require.Nil(t, err, "marshalling should succeed")

	packed, err = MarshalEtcdRaftMetadata(md)
	require.Nil(t, err, "marshalling should succeed a second time because we did not mutate ourselves")

	unpacked := &etcdraft.ConfigMetadata{}
	require.Nil(t, proto.Unmarshal(packed, unpacked), "unmarshalling should succeed")

	var outputCerts, inputCerts [3][]byte
	for i := range unpacked.GetConsenters() {
		outputCerts[i] = []byte(unpacked.GetConsenters()[i].GetClientTlsCert())
		inputCerts[i], _ = ioutil.ReadFile(fmt.Sprintf("testdata/tls-client-%d.pem", i+1))

	}

	for i := 0; i < len(inputCerts)-1; i++ {
		require.NotEqual(t, outputCerts[i+1], outputCerts[i], "expected extracted certs to differ from each other")
	}
}
