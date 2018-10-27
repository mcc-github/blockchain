/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package metrics_test

import (
	kitmetrics "github.com/go-kit/kit/metrics"
)




type counter interface {
	kitmetrics.Counter
}


type gauge interface {
	kitmetrics.Gauge
}


type histogram interface {
	kitmetrics.Histogram
}
