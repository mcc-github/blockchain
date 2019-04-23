













package wal

import "github.com/prometheus/client_golang/prometheus"

var (
	walFsyncSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd",
		Subsystem: "disk",
		Name:      "wal_fsync_duration_seconds",
		Help:      "The latency distributions of fsync called by WAL.",

		
		
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
	})
)

func init() {
	prometheus.MustRegister(walFsyncSec)
}
