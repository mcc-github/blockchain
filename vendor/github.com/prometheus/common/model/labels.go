












package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	
	AlertNameLabel = "alertname"

	
	
	ExportedLabelPrefix = "exported_"

	
	
	MetricNameLabel = "__name__"

	
	
	SchemeLabel = "__scheme__"

	
	
	AddressLabel = "__address__"

	
	
	MetricsPathLabel = "__metrics_path__"

	
	
	ReservedLabelPrefix = "__"

	
	
	
	MetaLabelPrefix = "__meta_"

	
	
	
	
	TmpLabelPrefix = "__tmp_"

	
	
	ParamLabelPrefix = "__param_"

	
	
	JobLabel = "job"

	
	InstanceLabel = "instance"

	
	
	BucketLabel = "le"

	
	
	QuantileLabel = "quantile"
)




var LabelNameRE = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")



type LabelName string




func (ln LabelName) IsValid() bool {
	if len(ln) == 0 {
		return false
	}
	for i, b := range ln {
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || (b >= '0' && b <= '9' && i > 0)) {
			return false
		}
	}
	return true
}


func (ln *LabelName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	if !LabelName(s).IsValid() {
		return fmt.Errorf("%q is not a valid label name", s)
	}
	*ln = LabelName(s)
	return nil
}


func (ln *LabelName) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if !LabelName(s).IsValid() {
		return fmt.Errorf("%q is not a valid label name", s)
	}
	*ln = LabelName(s)
	return nil
}


type LabelNames []LabelName

func (l LabelNames) Len() int {
	return len(l)
}

func (l LabelNames) Less(i, j int) bool {
	return l[i] < l[j]
}

func (l LabelNames) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l LabelNames) String() string {
	labelStrings := make([]string, 0, len(l))
	for _, label := range l {
		labelStrings = append(labelStrings, string(label))
	}
	return strings.Join(labelStrings, ", ")
}


type LabelValue string


func (lv LabelValue) IsValid() bool {
	return utf8.ValidString(string(lv))
}


type LabelValues []LabelValue

func (l LabelValues) Len() int {
	return len(l)
}

func (l LabelValues) Less(i, j int) bool {
	return string(l[i]) < string(l[j])
}

func (l LabelValues) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}


type LabelPair struct {
	Name  LabelName
	Value LabelValue
}



type LabelPairs []*LabelPair

func (l LabelPairs) Len() int {
	return len(l)
}

func (l LabelPairs) Less(i, j int) bool {
	switch {
	case l[i].Name > l[j].Name:
		return false
	case l[i].Name < l[j].Name:
		return true
	case l[i].Value > l[j].Value:
		return false
	case l[i].Value < l[j].Value:
		return true
	default:
		return false
	}
}

func (l LabelPairs) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
