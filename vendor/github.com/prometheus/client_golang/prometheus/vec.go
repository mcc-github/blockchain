












package prometheus

import (
	"fmt"
	"sync"

	"github.com/prometheus/common/model"
)






type MetricVec struct {
	mtx      sync.RWMutex 
	children map[uint64][]metricWithLabelValues
	desc     *Desc

	newMetric   func(labelValues ...string) Metric
	hashAdd     func(h uint64, s string) uint64 
	hashAddByte func(h uint64, b byte) uint64
}



func newMetricVec(desc *Desc, newMetric func(lvs ...string) Metric) *MetricVec {
	return &MetricVec{
		children:    map[uint64][]metricWithLabelValues{},
		desc:        desc,
		newMetric:   newMetric,
		hashAdd:     hashAdd,
		hashAddByte: hashAddByte,
	}
}



type metricWithLabelValues struct {
	values []string
	metric Metric
}



func (m *MetricVec) Describe(ch chan<- *Desc) {
	ch <- m.desc
}


func (m *MetricVec) Collect(ch chan<- Metric) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	for _, metrics := range m.children {
		for _, metric := range metrics {
			ch <- metric.metric
		}
	}
}

























func (m *MetricVec) GetMetricWithLabelValues(lvs ...string) (Metric, error) {
	h, err := m.hashLabelValues(lvs)
	if err != nil {
		return nil, err
	}

	return m.getOrCreateMetricWithLabelValues(h, lvs), nil
}













func (m *MetricVec) GetMetricWith(labels Labels) (Metric, error) {
	h, err := m.hashLabels(labels)
	if err != nil {
		return nil, err
	}

	return m.getOrCreateMetricWithLabels(h, labels), nil
}




func (m *MetricVec) WithLabelValues(lvs ...string) Metric {
	metric, err := m.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return metric
}




func (m *MetricVec) With(labels Labels) Metric {
	metric, err := m.GetMetricWith(labels)
	if err != nil {
		panic(err)
	}
	return metric
}
















func (m *MetricVec) DeleteLabelValues(lvs ...string) bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	h, err := m.hashLabelValues(lvs)
	if err != nil {
		return false
	}
	return m.deleteByHashWithLabelValues(h, lvs)
}











func (m *MetricVec) Delete(labels Labels) bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	h, err := m.hashLabels(labels)
	if err != nil {
		return false
	}

	return m.deleteByHashWithLabels(h, labels)
}




func (m *MetricVec) deleteByHashWithLabelValues(h uint64, lvs []string) bool {
	metrics, ok := m.children[h]
	if !ok {
		return false
	}

	i := m.findMetricWithLabelValues(metrics, lvs)
	if i >= len(metrics) {
		return false
	}

	if len(metrics) > 1 {
		m.children[h] = append(metrics[:i], metrics[i+1:]...)
	} else {
		delete(m.children, h)
	}
	return true
}




func (m *MetricVec) deleteByHashWithLabels(h uint64, labels Labels) bool {
	metrics, ok := m.children[h]
	if !ok {
		return false
	}
	i := m.findMetricWithLabels(metrics, labels)
	if i >= len(metrics) {
		return false
	}

	if len(metrics) > 1 {
		m.children[h] = append(metrics[:i], metrics[i+1:]...)
	} else {
		delete(m.children, h)
	}
	return true
}


func (m *MetricVec) Reset() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for h := range m.children {
		delete(m.children, h)
	}
}

func (m *MetricVec) hashLabelValues(vals []string) (uint64, error) {
	if len(vals) != len(m.desc.variableLabels) {
		return 0, errInconsistentCardinality
	}
	h := hashNew()
	for _, val := range vals {
		h = m.hashAdd(h, val)
		h = m.hashAddByte(h, model.SeparatorByte)
	}
	return h, nil
}

func (m *MetricVec) hashLabels(labels Labels) (uint64, error) {
	if len(labels) != len(m.desc.variableLabels) {
		return 0, errInconsistentCardinality
	}
	h := hashNew()
	for _, label := range m.desc.variableLabels {
		val, ok := labels[label]
		if !ok {
			return 0, fmt.Errorf("label name %q missing in label map", label)
		}
		h = m.hashAdd(h, val)
		h = m.hashAddByte(h, model.SeparatorByte)
	}
	return h, nil
}





func (m *MetricVec) getOrCreateMetricWithLabelValues(hash uint64, lvs []string) Metric {
	m.mtx.RLock()
	metric, ok := m.getMetricWithLabelValues(hash, lvs)
	m.mtx.RUnlock()
	if ok {
		return metric
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	metric, ok = m.getMetricWithLabelValues(hash, lvs)
	if !ok {
		
		copiedLVs := make([]string, len(lvs))
		copy(copiedLVs, lvs)
		metric = m.newMetric(copiedLVs...)
		m.children[hash] = append(m.children[hash], metricWithLabelValues{values: copiedLVs, metric: metric})
	}
	return metric
}





func (m *MetricVec) getOrCreateMetricWithLabels(hash uint64, labels Labels) Metric {
	m.mtx.RLock()
	metric, ok := m.getMetricWithLabels(hash, labels)
	m.mtx.RUnlock()
	if ok {
		return metric
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	metric, ok = m.getMetricWithLabels(hash, labels)
	if !ok {
		lvs := m.extractLabelValues(labels)
		metric = m.newMetric(lvs...)
		m.children[hash] = append(m.children[hash], metricWithLabelValues{values: lvs, metric: metric})
	}
	return metric
}



func (m *MetricVec) getMetricWithLabelValues(h uint64, lvs []string) (Metric, bool) {
	metrics, ok := m.children[h]
	if ok {
		if i := m.findMetricWithLabelValues(metrics, lvs); i < len(metrics) {
			return metrics[i].metric, true
		}
	}
	return nil, false
}



func (m *MetricVec) getMetricWithLabels(h uint64, labels Labels) (Metric, bool) {
	metrics, ok := m.children[h]
	if ok {
		if i := m.findMetricWithLabels(metrics, labels); i < len(metrics) {
			return metrics[i].metric, true
		}
	}
	return nil, false
}



func (m *MetricVec) findMetricWithLabelValues(metrics []metricWithLabelValues, lvs []string) int {
	for i, metric := range metrics {
		if m.matchLabelValues(metric.values, lvs) {
			return i
		}
	}
	return len(metrics)
}



func (m *MetricVec) findMetricWithLabels(metrics []metricWithLabelValues, labels Labels) int {
	for i, metric := range metrics {
		if m.matchLabels(metric.values, labels) {
			return i
		}
	}
	return len(metrics)
}

func (m *MetricVec) matchLabelValues(values []string, lvs []string) bool {
	if len(values) != len(lvs) {
		return false
	}
	for i, v := range values {
		if v != lvs[i] {
			return false
		}
	}
	return true
}

func (m *MetricVec) matchLabels(values []string, labels Labels) bool {
	if len(labels) != len(values) {
		return false
	}
	for i, k := range m.desc.variableLabels {
		if values[i] != labels[k] {
			return false
		}
	}
	return true
}

func (m *MetricVec) extractLabelValues(labels Labels) []string {
	labelValues := make([]string, len(labels))
	for i, k := range m.desc.variableLabels {
		labelValues[i] = labels[k]
	}
	return labelValues
}
