/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kafka_test

import (
	"testing"

	"github.com/mcc-github/blockchain/common/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gometrics "github.com/rcrowley/go-metrics"
)


type metricsRegistry interface {
	gometrics.Registry
}


type metricsMeter interface {
	gometrics.Meter
}


type metricsHistogram interface {
	gometrics.Histogram
}


type metricsProvider interface {
	metrics.Provider
}


type metricsGauge interface {
	metrics.Gauge
}

func TestKafka(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kafka Suite")
}
