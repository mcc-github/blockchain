package sarama

import (
	"fmt"
	"strings"

	"github.com/rcrowley/go-metrics"
)





const (
	metricsReservoirSize = 1028
	metricsAlphaFactor   = 0.015
)

func getOrRegisterHistogram(name string, r metrics.Registry) metrics.Histogram {
	return r.GetOrRegister(name, func() metrics.Histogram {
		return metrics.NewHistogram(metrics.NewExpDecaySample(metricsReservoirSize, metricsAlphaFactor))
	}).(metrics.Histogram)
}

func getMetricNameForBroker(name string, broker *Broker) string {
	
	
	return fmt.Sprintf(name+"-for-broker-%d", broker.ID())
}

func getOrRegisterBrokerMeter(name string, broker *Broker, r metrics.Registry) metrics.Meter {
	return metrics.GetOrRegisterMeter(getMetricNameForBroker(name, broker), r)
}

func getOrRegisterBrokerHistogram(name string, broker *Broker, r metrics.Registry) metrics.Histogram {
	return getOrRegisterHistogram(getMetricNameForBroker(name, broker), r)
}

func getMetricNameForTopic(name string, topic string) string {
	
	
	return fmt.Sprintf(name+"-for-topic-%s", strings.Replace(topic, ".", "_", -1))
}

func getOrRegisterTopicMeter(name string, topic string, r metrics.Registry) metrics.Meter {
	return metrics.GetOrRegisterMeter(getMetricNameForTopic(name, topic), r)
}

func getOrRegisterTopicHistogram(name string, topic string, r metrics.Registry) metrics.Histogram {
	return getOrRegisterHistogram(getMetricNameForTopic(name, topic), r)
}
