/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"time"

	"github.com/mcc-github/blockchain/gossip/election"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/spf13/viper"
)

const (
	btlPullMarginDefault           = 10
	transientBlockRetentionDefault = 1000
)


type ServiceConfig struct {
	
	PeerTLSEnabled bool
	
	Endpoint              string
	NonBlockingCommitMode bool
	
	UseLeaderElection bool
	
	OrgLeader bool
	
	
	ElectionStartupGracePeriod time.Duration
	
	ElectionMembershipSampleInterval time.Duration
	
	
	ElectionLeaderAliveThreshold time.Duration
	
	
	ElectionLeaderElectionDuration time.Duration
	
	
	PvtDataPullRetryThreshold time.Duration
	
	
	PvtDataPushAckTimeout time.Duration
	
	
	BtlPullMargin uint64
	
	
	TransientstoreMaxBlockRetention uint64
	
	
	
	SkipPullingInvalidTransactionsDuringCommit bool
}

func GlobalConfig() *ServiceConfig {
	c := &ServiceConfig{}
	c.loadGossipConfig()
	return c
}

func (c *ServiceConfig) loadGossipConfig() {

	c.PeerTLSEnabled = viper.GetBool("peer.tls.enabled")
	c.Endpoint = viper.GetString("peer.gossip.endpoint")
	c.NonBlockingCommitMode = viper.GetBool("peer.gossip.nonBlockingCommitMode")
	c.UseLeaderElection = viper.GetBool("peer.gossip.useLeaderElection")
	c.OrgLeader = viper.GetBool("peer.gossip.orgLeader")

	c.ElectionStartupGracePeriod = util.GetDurationOrDefault("peer.gossip.election.startupGracePeriod", election.DefStartupGracePeriod)
	c.ElectionMembershipSampleInterval = util.GetDurationOrDefault("peer.gossip.election.membershipSampleInterval", election.DefMembershipSampleInterval)
	c.ElectionLeaderAliveThreshold = util.GetDurationOrDefault("peer.gossip.election.leaderAliveThreshold", election.DefLeaderAliveThreshold)
	c.ElectionLeaderElectionDuration = util.GetDurationOrDefault("peer.gossip.election.leaderElectionDuration", election.DefLeaderElectionDuration)

	c.PvtDataPushAckTimeout = viper.GetDuration("peer.gossip.pvtData.pushAckTimeout")
	c.PvtDataPullRetryThreshold = viper.GetDuration("peer.gossip.pvtData.pullRetryThreshold")
	c.SkipPullingInvalidTransactionsDuringCommit = viper.GetBool("peer.gossip.pvtData.skipPullingInvalidTransactionsDuringCommit")

	c.BtlPullMargin = btlPullMarginDefault
	if viper.IsSet("peer.gossip.pvtData.btlPullMargin") {
		btlMarginVal := viper.GetInt("peer.gossip.pvtData.btlPullMargin")
		if btlMarginVal >= 0 {
			c.BtlPullMargin = uint64(btlMarginVal)
		}
	}

	c.TransientstoreMaxBlockRetention = uint64(viper.GetInt("peer.gossip.pvtData.transientstoreMaxBlockRetention"))
	if c.TransientstoreMaxBlockRetention == 0 {
		logger.Warning("Configuration key peer.gossip.pvtData.transientstoreMaxBlockRetention isn't set, defaulting to", transientBlockRetentionDefault)
		c.TransientstoreMaxBlockRetention = transientBlockRetentionDefault
	}
}
