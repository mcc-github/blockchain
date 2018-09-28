/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"bytes"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/orderer/consensus"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer/etcdraft"

	"code.cloudfoundry.org/clock"
	"github.com/coreos/etcd/raft"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)


type Consenter struct {
	Cert   []byte
	Logger *flogging.FabricLogger
}

func (c *Consenter) detectRaftID(m *etcdraft.Metadata) (uint64, error) {
	for i, cst := range m.Consenters {
		if bytes.Equal(c.Cert, cst.ServerTlsCert) {
			return uint64(i + 1), nil
		}
	}

	return 0, errors.Errorf("failed to detect Raft ID because no matching certificate found")
}

func (c *Consenter) HandleChain(support consensus.ConsenterSupport, metadata *common.Metadata) (consensus.Chain, error) {
	m := &etcdraft.Metadata{}
	if err := proto.Unmarshal(support.SharedConfig().ConsensusMetadata(), m); err != nil {
		return nil, errors.Errorf("failed to unmarshal consensus metadata: %s", err)
	}

	id, err := c.detectRaftID(m)
	if err != nil {
		return nil, err
	}

	peers := make([]raft.Peer, len(m.Consenters))
	for i := range peers {
		peers[i].ID = uint64(i + 1)
	}

	opts := Options{
		RaftID:  id,
		Clock:   clock.NewClock(),
		Storage: raft.NewMemoryStorage(),
		Logger:  c.Logger,

		
		TickInterval:    100 * time.Millisecond,
		ElectionTick:    10,
		HeartbeatTick:   1,
		MaxInflightMsgs: 256,
		MaxSizePerMsg:   1024 * 1024, 

		Peers: peers,
	}

	return NewChain(support, opts, nil)
}
