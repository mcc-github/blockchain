












package model

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	separator = []byte{0}
	
	
	
	MetricNameRE = regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*$`)
)



type Metric LabelSet


func (m Metric) Equal(o Metric) bool {
	return LabelSet(m).Equal(LabelSet(o))
}


func (m Metric) Before(o Metric) bool {
	return LabelSet(m).Before(LabelSet(o))
}


func (m Metric) Clone() Metric {
	clone := make(Metric, len(m))
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

func (m Metric) String() string {
	metricName, hasName := m[MetricNameLabel]
	numLabels := len(m) - 1
	if !hasName {
		numLabels = len(m)
	}
	labelStrings := make([]string, 0, numLabels)
	for label, value := range m {
		if label != MetricNameLabel {
			labelStrings = append(labelStrings, fmt.Sprintf("%s=%q", label, value))
		}
	}

	switch numLabels {
	case 0:
		if hasName {
			return string(metricName)
		}
		return "{}"
	default:
		sort.Strings(labelStrings)
		return fmt.Sprintf("%s{%s}", metricName, strings.Join(labelStrings, ", "))
	}
}


func (m Metric) Fingerprint() Fingerprint {
	return LabelSet(m).Fingerprint()
}



func (m Metric) FastFingerprint() Fingerprint {
	return LabelSet(m).FastFingerprint()
}




func IsValidMetricName(n LabelValue) bool {
	if len(n) == 0 {
		return false
	}
	for i, b := range n {
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || b == ':' || (b >= '0' && b <= '9' && i > 0)) {
			return false
		}
	}
	return true
}
