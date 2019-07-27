/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"net"
	"strconv"
	"time"

	"github.com/mcc-github/blockchain/gossip/comm"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/election"
	"github.com/mcc-github/blockchain/gossip/gossip/algo"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/spf13/viper"
)


type Config struct {
	
	BindPort int
	
	ID string
	
	BootstrapPeers []string
	
	PropagateIterations int
	
	PropagatePeerNum int

	
	MaxBlockCountToStore int

	
	MaxPropagationBurstSize int
	
	MaxPropagationBurstLatency time.Duration

	
	PullInterval time.Duration
	
	PullPeerNum int

	
	SkipBlockVerification bool

	
	PublishCertPeriod time.Duration
	
	PublishStateInfoInterval time.Duration
	
	RequestStateInfoInterval time.Duration

	
	TLSCerts *common.TLSCertificates

	
	InternalEndpoint string
	
	ExternalEndpoint string
	
	TimeForMembershipTracker time.Duration

	
	DigestWaitTime time.Duration
	
	RequestWaitTime time.Duration
	
	ResponseWaitTime time.Duration

	
	DialTimeout time.Duration
	
	ConnTimeout time.Duration
	
	RecvBuffSize int
	
	SendBuffSize int

	
	MsgExpirationTimeout time.Duration

	
	AliveTimeInterval time.Duration
	
	AliveExpirationTimeout time.Duration
	
	AliveExpirationCheckInterval time.Duration
	
	ReconnectInterval time.Duration
}

func GlobalConfig(endpoint string, certs *common.TLSCertificates, bootPeers ...string) (*Config, error) {
	c := &Config{}
	err := c.loadConfig(endpoint, certs, bootPeers...)
	return c, err
}

func (c *Config) loadConfig(endpoint string, certs *common.TLSCertificates, bootPeers ...string) error {
	_, p, err := net.SplitHostPort(endpoint)
	if err != nil {
		return err
	}
	port, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return err
	}

	c.BindPort = int(port)
	c.BootstrapPeers = bootPeers
	c.ID = endpoint
	c.MaxBlockCountToStore = util.GetIntOrDefault("peer.gossip.maxBlockCountToStore", 100)
	c.MaxPropagationBurstLatency = util.GetDurationOrDefault("peer.gossip.maxPropagationBurstLatency", 10*time.Millisecond)
	c.MaxPropagationBurstSize = util.GetIntOrDefault("peer.gossip.maxPropagationBurstSize", 10)
	c.PropagateIterations = util.GetIntOrDefault("peer.gossip.propagateIterations", 1)
	c.PropagatePeerNum = util.GetIntOrDefault("peer.gossip.propagatePeerNum", 3)
	c.PullInterval = util.GetDurationOrDefault("peer.gossip.pullInterval", 4*time.Second)
	c.PullPeerNum = util.GetIntOrDefault("peer.gossip.pullPeerNum", 3)
	c.InternalEndpoint = endpoint
	c.ExternalEndpoint = viper.GetString("peer.gossip.externalEndpoint")
	c.PublishCertPeriod = util.GetDurationOrDefault("peer.gossip.publishCertPeriod", 10*time.Second)
	c.RequestStateInfoInterval = util.GetDurationOrDefault("peer.gossip.requestStateInfoInterval", 4*time.Second)
	c.PublishStateInfoInterval = util.GetDurationOrDefault("peer.gossip.publishStateInfoInterval", 4*time.Second)
	c.SkipBlockVerification = viper.GetBool("peer.gossip.skipBlockVerification")
	c.TLSCerts = certs
	c.TimeForMembershipTracker = util.GetDurationOrDefault("peer.gossip.membershipTrackerInterval", 5*time.Second)
	c.DigestWaitTime = util.GetDurationOrDefault("peer.gossip.digestWaitTime", algo.DefDigestWaitTime)
	c.RequestWaitTime = util.GetDurationOrDefault("peer.gossip.requestWaitTime", algo.DefRequestWaitTime)
	c.ResponseWaitTime = util.GetDurationOrDefault("peer.gossip.responseWaitTime", algo.DefResponseWaitTime)
	c.DialTimeout = util.GetDurationOrDefault("peer.gossip.dialTimeout", comm.DefDialTimeout)
	c.ConnTimeout = util.GetDurationOrDefault("peer.gossip.connTimeout", comm.DefConnTimeout)
	c.RecvBuffSize = util.GetIntOrDefault("peer.gossip.recvBuffSize", comm.DefRecvBuffSize)
	c.SendBuffSize = util.GetIntOrDefault("peer.gossip.sendBuffSize", comm.DefSendBuffSize)
	c.MsgExpirationTimeout = util.GetDurationOrDefault("peer.gossip.election.leaderAliveThreshold", election.DefLeaderAliveThreshold) * 10
	c.AliveTimeInterval = util.GetDurationOrDefault("peer.gossip.aliveTimeInterval", discovery.DefAliveTimeInterval)
	c.AliveExpirationTimeout = util.GetDurationOrDefault("peer.gossip.aliveExpirationTimeout", 5*c.AliveTimeInterval)
	c.AliveExpirationCheckInterval = c.AliveExpirationTimeout / 10
	c.ReconnectInterval = util.GetDurationOrDefault("peer.gossip.reconnectInterval", c.AliveExpirationTimeout)

	return nil
}
