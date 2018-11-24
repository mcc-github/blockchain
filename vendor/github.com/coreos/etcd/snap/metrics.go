













package snap

import "github.com/prometheus/client_golang/prometheus"

var (
	
	saveDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd_debugging",
		Subsystem: "snap",
		Name:      "save_total_duration_seconds",
		Help:      "The total latency distributions of save called by snapshot.",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 14),
	})

	marshallingDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "etcd_debugging",
		Subsystem: "snap",
		Name:      "save_marshalling_duration_seconds",
		Help:      "The marshalling cost distributions of save called by snapshot.",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 14),
	})
)

func init() {
	prometheus.MustRegister(saveDurations)
	prometheus.MustRegister(marshallingDurations)
}
