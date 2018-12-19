/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package testdata

import (
	goo "github.com/mcc-github/blockchain/common/metrics"
)




var (
	NamedCounter = goo.CounterOpts{
		Namespace:    "namespace",
		Subsystem:    "counter",
		Name:         "name",
		Help:         "This is some help text",
		LabelNames:   []string{"label_one", "label_two"},
		StatsdFormat: "%{#fqname}.%{label_one}.%{label_two}",
	}

	NamedGauge = goo.GaugeOpts{
		Namespace:    "namespace",
		Subsystem:    "gauge",
		Name:         "name",
		Help:         "This is some help text",
		LabelNames:   []string{"label_one", "label_two"},
		StatsdFormat: "%{#fqname}.%{label_one}.%{label_two}",
	}

	NamedHistogram = goo.HistogramOpts{
		Namespace:    "namespace",
		Subsystem:    "histogram",
		Name:         "name",
		Help:         "This is some help text",
		LabelNames:   []string{"label_one", "label_two"},
		StatsdFormat: "%{#fqname}.%{label_one}.%{label_two}",
	}
)
