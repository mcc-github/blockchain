/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deliverservice

import (
	"time"

	"github.com/mcc-github/blockchain/core/comm"
	"github.com/spf13/viper"
)

const (
	DefaultReConnectBackoffThreshold   = float64(time.Hour)
	DefaultReConnectTotalTimeThreshold = time.Second * 60 * 60
	DefaultConnectionTimeout           = time.Second * 3
)


type DeliverServiceConfig struct {
	
	PeerTLSEnabled bool
	
	ReConnectBackoffThreshold float64
	
	
	ReconnectTotalTimeThreshold time.Duration
	
	ConnectionTimeout time.Duration
	
	KeepaliveOptions comm.KeepaliveOptions
}


func GlobalConfig() *DeliverServiceConfig {
	c := &DeliverServiceConfig{}
	c.loadDeliverServiceConfig()
	return c
}

func (c *DeliverServiceConfig) loadDeliverServiceConfig() {
	c.PeerTLSEnabled = viper.GetBool("peer.tls.enabled")

	c.ReConnectBackoffThreshold = viper.GetFloat64("peer.deliveryclient.reConnectBackoffThreshold")
	if c.ReConnectBackoffThreshold == 0 {
		c.ReConnectBackoffThreshold = DefaultReConnectBackoffThreshold
	}

	c.ReconnectTotalTimeThreshold = viper.GetDuration("peer.deliveryclient.reconnectTotalTimeThreshold")
	if c.ReconnectTotalTimeThreshold == 0 {
		c.ReconnectTotalTimeThreshold = DefaultReConnectTotalTimeThreshold
	}

	c.ConnectionTimeout = viper.GetDuration("peer.deliveryclient.connTimeout")
	if c.ConnectionTimeout == 0 {
		c.ConnectionTimeout = DefaultConnectionTimeout
	}

	c.KeepaliveOptions = comm.DefaultKeepaliveOptions
	if viper.IsSet("peer.keepalive.deliveryClient.interval") {
		c.KeepaliveOptions.ClientInterval = viper.GetDuration("peer.keepalive.deliveryClient.interval")
	}
	if viper.IsSet("peer.keepalive.deliveryClient.timeout") {
		c.KeepaliveOptions.ClientTimeout = viper.GetDuration("peer.keepalive.deliveryClient.timeout")
	}

}
