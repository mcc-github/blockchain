/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"bytes"
	"encoding/pem"
	"reflect"

	"github.com/coreos/etcd/raft"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/orderer/common/cluster"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/orderer/consensus"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer/etcdraft"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)


func EndpointconfigFromFromSupport(support consensus.ConsenterSupport) (*cluster.EndpointConfig, error) {
	lastConfigBlock, err := lastConfigBlockFromSupport(support)
	if err != nil {
		return nil, err
	}
	endpointconf, err := cluster.EndpointconfigFromConfigBlock(lastConfigBlock)
	if err != nil {
		return nil, err
	}
	return endpointconf, nil
}

func lastConfigBlockFromSupport(support consensus.ConsenterSupport) (*common.Block, error) {
	lastBlockSeq := support.Height() - 1
	lastBlock := support.Block(lastBlockSeq)
	if lastBlock == nil {
		return nil, errors.Errorf("unable to retrieve block %d", lastBlockSeq)
	}
	lastConfigBlock, err := LastConfigBlock(lastBlock, support)
	if err != nil {
		return nil, err
	}
	return lastConfigBlock, nil
}


func LastConfigBlock(block *common.Block, support consensus.ConsenterSupport) (*common.Block, error) {
	if block == nil {
		return nil, errors.New("nil block")
	}
	if support == nil {
		return nil, errors.New("nil support")
	}
	if block.Metadata == nil || len(block.Metadata.Metadata) <= int(common.BlockMetadataIndex_LAST_CONFIG) {
		return nil, errors.New("no metadata in block")
	}
	lastConfigBlockNum, err := utils.GetLastConfigIndexFromBlock(block)
	if err != nil {
		return nil, err
	}
	lastConfigBlock := support.Block(lastConfigBlockNum)
	if lastConfigBlock == nil {
		return nil, errors.Errorf("unable to retrieve last config block %d", lastConfigBlockNum)
	}
	return lastConfigBlock, nil
}


func newBlockPuller(support consensus.ConsenterSupport,
	baseDialer *cluster.PredicateDialer,
	clusterConfig localconfig.Cluster) (*cluster.BlockPuller, error) {

	verifyBlockSequence := func(blocks []*common.Block) error {
		return cluster.VerifyBlocks(blocks, support)
	}

	secureConfig, err := baseDialer.ClientConfig()
	if err != nil {
		return nil, err
	}
	secureConfig.AsyncConnect = false
	stdDialer := &cluster.StandardDialer{
		Dialer: cluster.NewTLSPinningDialer(secureConfig),
	}

	
	endpointConfig, err := EndpointconfigFromFromSupport(support)
	if err != nil {
		return nil, err
	}
	
	secureConfig.SecOpts.ServerRootCAs = endpointConfig.TLSRootCAs
	stdDialer.Dialer.SetConfig(secureConfig)

	der, _ := pem.Decode(secureConfig.SecOpts.Certificate)
	if der == nil {
		return nil, errors.Errorf("client certificate isn't in PEM format: %v",
			string(secureConfig.SecOpts.Certificate))
	}

	return &cluster.BlockPuller{
		VerifyBlockSequence: verifyBlockSequence,
		Logger:              flogging.MustGetLogger("orderer.common.cluster.puller"),
		RetryTimeout:        clusterConfig.ReplicationRetryTimeout,
		MaxTotalBufferBytes: clusterConfig.ReplicationBufferSize,
		FetchTimeout:        clusterConfig.ReplicationPullTimeout,
		Endpoints:           endpointConfig.Endpoints,
		Signer:              support,
		TLSCert:             der.Bytes,
		Channel:             support.ChainID(),
		Dialer:              stdDialer,
	}, nil
}


func RaftPeers(consenters map[uint64]*etcdraft.Consenter) []raft.Peer {
	var peers []raft.Peer

	for raftID := range consenters {
		peers = append(peers, raft.Peer{ID: raftID})
	}
	return peers
}


func ConsentersToMap(consenters []*etcdraft.Consenter) map[string]struct{} {
	set := map[string]struct{}{}
	for _, c := range consenters {
		set[string(c.ClientTlsCert)] = struct{}{}
	}
	return set
}



func MembershipByCert(consenters map[uint64]*etcdraft.Consenter) map[string]struct{} {
	set := map[string]struct{}{}
	for _, c := range consenters {
		set[string(c.ClientTlsCert)] = struct{}{}
	}
	return set
}


func ConsentersChanged(oldConsenters map[uint64]*etcdraft.Consenter, newConsenters []*etcdraft.Consenter) bool {
	if len(oldConsenters) != len(newConsenters) {
		return false
	}

	consentersSet1 := MembershipByCert(oldConsenters)
	consentersSet2 := ConsentersToMap(newConsenters)

	return reflect.DeepEqual(consentersSet1, consentersSet2)
}


type ConsenterCertificate []byte




func (conCert ConsenterCertificate) IsConsenterOfChannel(configBlock *common.Block) error {
	if configBlock == nil {
		return errors.New("nil block")
	}
	envelopeConfig, err := utils.ExtractEnvelope(configBlock, 0)
	if err != nil {
		return err
	}
	bundle, err := channelconfig.NewBundleFromEnvelope(envelopeConfig)
	if err != nil {
		return err
	}
	oc, exists := bundle.OrdererConfig()
	if !exists {
		return errors.New("no orderer config in bundle")
	}
	m := &etcdraft.Metadata{}
	if err := proto.Unmarshal(oc.ConsensusMetadata(), m); err != nil {
		return err
	}

	for _, consenter := range m.Consenters {
		if bytes.Equal(conCert, consenter.ServerTlsCert) || bytes.Equal(conCert, consenter.ClientTlsCert) {
			return nil
		}
	}
	return cluster.NotInChannelError
}
