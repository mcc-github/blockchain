












package prometheus

import (
	"strings"

	dto "github.com/prometheus/client_model/go"
)

const separatorByte byte = 255




type Metric interface {
	
	
	
	
	
	Desc() *Desc
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Write(*dto.Metric) error
	
	
	
	
}







type Opts struct {
	
	
	
	
	
	Namespace string
	Subsystem string
	Name      string

	
	
	
	
	Help string

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ConstLabels Labels
}








func BuildFQName(namespace, subsystem, name string) string {
	if name == "" {
		return ""
	}
	switch {
	case namespace != "" && subsystem != "":
		return strings.Join([]string{namespace, subsystem, name}, "_")
	case namespace != "":
		return strings.Join([]string{namespace, name}, "_")
	case subsystem != "":
		return strings.Join([]string{subsystem, name}, "_")
	}
	return name
}




type LabelPairSorter []*dto.LabelPair

func (s LabelPairSorter) Len() int {
	return len(s)
}

func (s LabelPairSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s LabelPairSorter) Less(i, j int) bool {
	return s[i].GetName() < s[j].GetName()
}

type hashSorter []uint64

func (s hashSorter) Len() int {
	return len(s)
}

func (s hashSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s hashSorter) Less(i, j int) bool {
	return s[i] < s[j]
}

type invalidMetric struct {
	desc *Desc
	err  error
}




func NewInvalidMetric(desc *Desc, err error) Metric {
	return &invalidMetric{desc, err}
}

func (m *invalidMetric) Desc() *Desc { return m.desc }

func (m *invalidMetric) Write(*dto.Metric) error { return m.err }
