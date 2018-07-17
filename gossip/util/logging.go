/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/op/go-logging"
)


const (
	LoggingChannelModule   = "gossip/channel"
	LoggingCommModule      = "gossip/comm"
	LoggingDiscoveryModule = "gossip/discovery"
	LoggingElectionModule  = "gossip/election"
	LoggingGossipModule    = "gossip/gossip"
	LoggingMockModule      = "gossip/comm/mock"
	LoggingPullModule      = "gossip/pull"
	LoggingServiceModule   = "gossip/service"
	LoggingStateModule     = "gossip/state"
	LoggingPrivModule      = "gossip/privdata"
)

var loggersByModules = make(map[string]*logging.Logger)
var lock = sync.Mutex{}
var testMode bool


var defaultTestSpec = "WARNING"


func GetLogger(module string, peerID string) *logging.Logger {
	if peerID != "" && testMode {
		module = module + "#" + peerID
	}

	lock.Lock()
	defer lock.Unlock()

	if lgr, ok := loggersByModules[module]; ok {
		return lgr
	}

	
	lgr := flogging.MustGetLogger(module)
	loggersByModules[module] = lgr
	return lgr
}


func SetupTestLogging() {
	testMode = true
	flogging.InitFromSpec(defaultTestSpec)
}
