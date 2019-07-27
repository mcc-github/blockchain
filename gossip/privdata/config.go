/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"time"

	"github.com/spf13/viper"
)

const (
	reconcileSleepIntervalDefault = time.Minute
	reconcileBatchSizeDefault     = 10
)


type PrivdataConfig struct {
	
	
	ReconcileSleepInterval time.Duration
	
	ReconcileBatchSize int
	
	ReconciliationEnabled bool
}


func GlobalConfig() *PrivdataConfig {
	c := &PrivdataConfig{}
	c.loadPrivDataConfig()
	return c
}

func (c *PrivdataConfig) loadPrivDataConfig() {
	c.ReconcileSleepInterval = viper.GetDuration("peer.gossip.pvtData.reconcileSleepInterval")
	if c.ReconcileSleepInterval == 0 {
		logger.Warning("Configuration key peer.gossip.pvtData.reconcileSleepInterval isn't set, defaulting to", reconcileSleepIntervalDefault)
		c.ReconcileSleepInterval = reconcileSleepIntervalDefault
	}

	c.ReconcileBatchSize = viper.GetInt("peer.gossip.pvtData.reconcileBatchSize")
	if c.ReconcileBatchSize == 0 {
		logger.Warning("Configuration key peer.gossip.pvtData.reconcileBatchSize isn't set, defaulting to", reconcileBatchSizeDefault)
		c.ReconcileBatchSize = reconcileBatchSizeDefault
	}

	c.ReconciliationEnabled = viper.GetBool("peer.gossip.pvtData.reconciliationEnabled")

}
