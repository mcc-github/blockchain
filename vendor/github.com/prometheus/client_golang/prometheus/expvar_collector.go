












package prometheus

import (
	"encoding/json"
	"expvar"
)

type expvarCollector struct {
	exports map[string]*Desc
}






































func NewExpvarCollector(exports map[string]*Desc) Collector {
	return &expvarCollector{
		exports: exports,
	}
}


func (e *expvarCollector) Describe(ch chan<- *Desc) {
	for _, desc := range e.exports {
		ch <- desc
	}
}


func (e *expvarCollector) Collect(ch chan<- Metric) {
	for name, desc := range e.exports {
		var m Metric
		expVar := expvar.Get(name)
		if expVar == nil {
			continue
		}
		var v interface{}
		labels := make([]string, len(desc.variableLabels))
		if err := json.Unmarshal([]byte(expVar.String()), &v); err != nil {
			ch <- NewInvalidMetric(desc, err)
			continue
		}
		var processValue func(v interface{}, i int)
		processValue = func(v interface{}, i int) {
			if i >= len(labels) {
				copiedLabels := append(make([]string, 0, len(labels)), labels...)
				switch v := v.(type) {
				case float64:
					m = MustNewConstMetric(desc, UntypedValue, v, copiedLabels...)
				case bool:
					if v {
						m = MustNewConstMetric(desc, UntypedValue, 1, copiedLabels...)
					} else {
						m = MustNewConstMetric(desc, UntypedValue, 0, copiedLabels...)
					}
				default:
					return
				}
				ch <- m
				return
			}
			vm, ok := v.(map[string]interface{})
			if !ok {
				return
			}
			for lv, val := range vm {
				labels[i] = lv
				processValue(val, i+1)
			}
		}
		processValue(v, 0)
	}
}
