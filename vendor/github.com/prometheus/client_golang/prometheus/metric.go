












package prometheus

import (
	"strings"
	"time"

	"github.com/golang/protobuf/proto"

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



type labelPairSorter []*dto.LabelPair

func (s labelPairSorter) Len() int {
	return len(s)
}

func (s labelPairSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s labelPairSorter) Less(i, j int) bool {
	return s[i].GetName() < s[j].GetName()
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

type timestampedMetric struct {
	Metric
	t time.Time
}

func (m timestampedMetric) Write(pb *dto.Metric) error {
	e := m.Metric.Write(pb)
	pb.TimestampMs = proto.Int64(m.t.Unix()*1000 + int64(m.t.Nanosecond()/1000000))
	return e
}














func NewMetricWithTimestamp(t time.Time, m Metric) Metric {
	return timestampedMetric{Metric: m, t: t}
}
