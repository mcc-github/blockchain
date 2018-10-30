












package prometheus

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/prometheus/common/model"
)








type Labels map[string]string



const reservedLabelPrefix = "__"

var errInconsistentCardinality = errors.New("inconsistent label cardinality")

func validateValuesInLabels(labels Labels, expectedNumberOfValues int) error {
	if len(labels) != expectedNumberOfValues {
		return errInconsistentCardinality
	}

	for name, val := range labels {
		if !utf8.ValidString(val) {
			return fmt.Errorf("label %s: value %q is not valid UTF-8", name, val)
		}
	}

	return nil
}

func validateLabelValues(vals []string, expectedNumberOfValues int) error {
	if len(vals) != expectedNumberOfValues {
		return errInconsistentCardinality
	}

	for _, val := range vals {
		if !utf8.ValidString(val) {
			return fmt.Errorf("label value %q is not valid UTF-8", val)
		}
	}

	return nil
}

func checkLabelName(l string) bool {
	return model.LabelName(l).IsValid() && !strings.HasPrefix(l, reservedLabelPrefix)
}
