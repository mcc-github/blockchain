/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package operations

import (
	"sync"

	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/metrics/prometheus"
)

var (
	blockchainVersion = metrics.GaugeOpts{
		Name:         "blockchain_version",
		Help:         "The active version of Fabric.",
		LabelNames:   []string{"version"},
		StatsdFormat: "%{#fqname}.%{version}",
	}

	gaugeLock        sync.Mutex
	promVersionGauge metrics.Gauge
)

func versionGauge(provider metrics.Provider) metrics.Gauge {
	switch provider.(type) {
	case *prometheus.Provider:
		gaugeLock.Lock()
		defer gaugeLock.Unlock()
		if promVersionGauge == nil {
			promVersionGauge = provider.NewGauge(blockchainVersion)
		}
		return promVersionGauge

	default:
		return provider.NewGauge(blockchainVersion)
	}
}
