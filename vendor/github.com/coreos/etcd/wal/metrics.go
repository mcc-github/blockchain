













package wal

import "github.com/prometheus/client_golang/prometheus"

var (
	syncDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd",
		Subsystem: "disk",
		Name:      "wal_fsync_duration_seconds",
		Help:      "The latency distributions of fsync called by wal.",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 14),
	})
)

func init() {
	prometheus.MustRegister(syncDurations)
}
