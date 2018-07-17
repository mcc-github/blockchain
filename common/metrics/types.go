/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package metrics

import "io"


type Counter interface {
	
	Inc(delta int64)
}


type Gauge interface {
	
	Update(value float64)
}



type Scope interface {
	serve
	
	Counter(name string) Counter

	
	Gauge(name string) Gauge

	
	Tagged(tags map[string]string) Scope

	
	SubScope(name string) Scope
}


type serve interface {
	io.Closer
	
	Start() error
}
