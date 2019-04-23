













package snap

import "github.com/prometheus/client_golang/prometheus"

var (
	snapMarshallingSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd_debugging",
		Subsystem: "snap",
		Name:      "save_marshalling_duration_seconds",
		Help:      "The marshalling cost distributions of save called by snapshot.",

		
		
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	snapSaveSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd_debugging",
		Subsystem: "snap",
		Name:      "save_total_duration_seconds",
		Help:      "The total latency distributions of save called by snapshot.",

		
		
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	snapFsyncSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd",
		Subsystem: "snap",
		Name:      "fsync_duration_seconds",
		Help:      "The latency distributions of fsync called by snap.",

		
		
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	snapDBSaveSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd",
		Subsystem: "snap_db",
		Name:      "save_total_duration_seconds",
		Help:      "The total latency distributions of v3 snapshot save",

		
		
		Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
	})

	snapDBFsyncSec = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd",
		Subsystem: "snap_db",
		Name:      "fsync_duration_seconds",
		Help:      "The latency distributions of fsyncing .snap.db file",

		
		
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
	})
)

func init() {
	prometheus.MustRegister(snapMarshallingSec)
	prometheus.MustRegister(snapSaveSec)
	prometheus.MustRegister(snapFsyncSec)
	prometheus.MustRegister(snapDBSaveSec)
	prometheus.MustRegister(snapDBFsyncSec)
}
