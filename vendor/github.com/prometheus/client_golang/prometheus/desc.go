












package prometheus

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"

	dto "github.com/prometheus/client_model/go"
)

var (
	metricNameRE = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_:]*$`)
	labelNameRE  = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
)



const reservedLabelPrefix = "__"








type Labels map[string]string

















type Desc struct {
	
	fqName string
	
	help string
	
	
	constLabelPairs []*dto.LabelPair
	
	
	variableLabels []string
	
	
	
	id uint64
	
	
	
	dimHash uint64
	
	
	err error
}











func NewDesc(fqName, help string, variableLabels []string, constLabels Labels) *Desc {
	d := &Desc{
		fqName:         fqName,
		help:           help,
		variableLabels: variableLabels,
	}
	if help == "" {
		d.err = errors.New("empty help string")
		return d
	}
	if !metricNameRE.MatchString(fqName) {
		d.err = fmt.Errorf("%q is not a valid metric name", fqName)
		return d
	}
	
	
	labelValues := make([]string, 1, len(constLabels)+1)
	labelValues[0] = fqName
	labelNames := make([]string, 0, len(constLabels)+len(variableLabels))
	labelNameSet := map[string]struct{}{}
	
	for labelName := range constLabels {
		if !checkLabelName(labelName) {
			d.err = fmt.Errorf("%q is not a valid label name", labelName)
			return d
		}
		labelNames = append(labelNames, labelName)
		labelNameSet[labelName] = struct{}{}
	}
	sort.Strings(labelNames)
	
	for _, labelName := range labelNames {
		labelValues = append(labelValues, constLabels[labelName])
	}
	
	
	
	for _, labelName := range variableLabels {
		if !checkLabelName(labelName) {
			d.err = fmt.Errorf("%q is not a valid label name", labelName)
			return d
		}
		labelNames = append(labelNames, "$"+labelName)
		labelNameSet[labelName] = struct{}{}
	}
	if len(labelNames) != len(labelNameSet) {
		d.err = errors.New("duplicate label names")
		return d
	}
	vh := hashNew()
	for _, val := range labelValues {
		vh = hashAdd(vh, val)
		vh = hashAddByte(vh, separatorByte)
	}
	d.id = vh
	
	sort.Strings(labelNames)
	
	
	lh := hashNew()
	lh = hashAdd(lh, help)
	lh = hashAddByte(lh, separatorByte)
	for _, labelName := range labelNames {
		lh = hashAdd(lh, labelName)
		lh = hashAddByte(lh, separatorByte)
	}
	d.dimHash = lh

	d.constLabelPairs = make([]*dto.LabelPair, 0, len(constLabels))
	for n, v := range constLabels {
		d.constLabelPairs = append(d.constLabelPairs, &dto.LabelPair{
			Name:  proto.String(n),
			Value: proto.String(v),
		})
	}
	sort.Sort(LabelPairSorter(d.constLabelPairs))
	return d
}





func NewInvalidDesc(err error) *Desc {
	return &Desc{
		err: err,
	}
}

func (d *Desc) String() string {
	lpStrings := make([]string, 0, len(d.constLabelPairs))
	for _, lp := range d.constLabelPairs {
		lpStrings = append(
			lpStrings,
			fmt.Sprintf("%s=%q", lp.GetName(), lp.GetValue()),
		)
	}
	return fmt.Sprintf(
		"Desc{fqName: %q, help: %q, constLabels: {%s}, variableLabels: %v}",
		d.fqName,
		d.help,
		strings.Join(lpStrings, ","),
		d.variableLabels,
	)
}

func checkLabelName(l string) bool {
	return labelNameRE.MatchString(l) &&
		!strings.HasPrefix(l, reservedLabelPrefix)
}
