












package expfmt

import (
	"fmt"
	"io"
	"math"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
)










func MetricFamilyToText(out io.Writer, in *dto.MetricFamily) (int, error) {
	var written int

	
	if len(in.Metric) == 0 {
		return written, fmt.Errorf("MetricFamily has no metrics: %s", in)
	}
	name := in.GetName()
	if name == "" {
		return written, fmt.Errorf("MetricFamily has no name: %s", in)
	}

	
	if in.Help != nil {
		n, err := fmt.Fprintf(
			out, "# HELP %s %s\n",
			name, escapeString(*in.Help, false),
		)
		written += n
		if err != nil {
			return written, err
		}
	}
	metricType := in.GetType()
	n, err := fmt.Fprintf(
		out, "# TYPE %s %s\n",
		name, strings.ToLower(metricType.String()),
	)
	written += n
	if err != nil {
		return written, err
	}

	
	for _, metric := range in.Metric {
		switch metricType {
		case dto.MetricType_COUNTER:
			if metric.Counter == nil {
				return written, fmt.Errorf(
					"expected counter in metric %s %s", name, metric,
				)
			}
			n, err = writeSample(
				name, metric, "", "",
				metric.Counter.GetValue(),
				out,
			)
		case dto.MetricType_GAUGE:
			if metric.Gauge == nil {
				return written, fmt.Errorf(
					"expected gauge in metric %s %s", name, metric,
				)
			}
			n, err = writeSample(
				name, metric, "", "",
				metric.Gauge.GetValue(),
				out,
			)
		case dto.MetricType_UNTYPED:
			if metric.Untyped == nil {
				return written, fmt.Errorf(
					"expected untyped in metric %s %s", name, metric,
				)
			}
			n, err = writeSample(
				name, metric, "", "",
				metric.Untyped.GetValue(),
				out,
			)
		case dto.MetricType_SUMMARY:
			if metric.Summary == nil {
				return written, fmt.Errorf(
					"expected summary in metric %s %s", name, metric,
				)
			}
			for _, q := range metric.Summary.Quantile {
				n, err = writeSample(
					name, metric,
					model.QuantileLabel, fmt.Sprint(q.GetQuantile()),
					q.GetValue(),
					out,
				)
				written += n
				if err != nil {
					return written, err
				}
			}
			n, err = writeSample(
				name+"_sum", metric, "", "",
				metric.Summary.GetSampleSum(),
				out,
			)
			if err != nil {
				return written, err
			}
			written += n
			n, err = writeSample(
				name+"_count", metric, "", "",
				float64(metric.Summary.GetSampleCount()),
				out,
			)
		case dto.MetricType_HISTOGRAM:
			if metric.Histogram == nil {
				return written, fmt.Errorf(
					"expected histogram in metric %s %s", name, metric,
				)
			}
			infSeen := false
			for _, q := range metric.Histogram.Bucket {
				n, err = writeSample(
					name+"_bucket", metric,
					model.BucketLabel, fmt.Sprint(q.GetUpperBound()),
					float64(q.GetCumulativeCount()),
					out,
				)
				written += n
				if err != nil {
					return written, err
				}
				if math.IsInf(q.GetUpperBound(), +1) {
					infSeen = true
				}
			}
			if !infSeen {
				n, err = writeSample(
					name+"_bucket", metric,
					model.BucketLabel, "+Inf",
					float64(metric.Histogram.GetSampleCount()),
					out,
				)
				if err != nil {
					return written, err
				}
				written += n
			}
			n, err = writeSample(
				name+"_sum", metric, "", "",
				metric.Histogram.GetSampleSum(),
				out,
			)
			if err != nil {
				return written, err
			}
			written += n
			n, err = writeSample(
				name+"_count", metric, "", "",
				float64(metric.Histogram.GetSampleCount()),
				out,
			)
		default:
			return written, fmt.Errorf(
				"unexpected type in metric %s %s", name, metric,
			)
		}
		written += n
		if err != nil {
			return written, err
		}
	}
	return written, nil
}





func writeSample(
	name string,
	metric *dto.Metric,
	additionalLabelName, additionalLabelValue string,
	value float64,
	out io.Writer,
) (int, error) {
	var written int
	n, err := fmt.Fprint(out, name)
	written += n
	if err != nil {
		return written, err
	}
	n, err = labelPairsToText(
		metric.Label,
		additionalLabelName, additionalLabelValue,
		out,
	)
	written += n
	if err != nil {
		return written, err
	}
	n, err = fmt.Fprintf(out, " %v", value)
	written += n
	if err != nil {
		return written, err
	}
	if metric.TimestampMs != nil {
		n, err = fmt.Fprintf(out, " %v", *metric.TimestampMs)
		written += n
		if err != nil {
			return written, err
		}
	}
	n, err = out.Write([]byte{'\n'})
	written += n
	if err != nil {
		return written, err
	}
	return written, nil
}








func labelPairsToText(
	in []*dto.LabelPair,
	additionalLabelName, additionalLabelValue string,
	out io.Writer,
) (int, error) {
	if len(in) == 0 && additionalLabelName == "" {
		return 0, nil
	}
	var written int
	separator := '{'
	for _, lp := range in {
		n, err := fmt.Fprintf(
			out, `%c%s="%s"`,
			separator, lp.GetName(), escapeString(lp.GetValue(), true),
		)
		written += n
		if err != nil {
			return written, err
		}
		separator = ','
	}
	if additionalLabelName != "" {
		n, err := fmt.Fprintf(
			out, `%c%s="%s"`,
			separator, additionalLabelName,
			escapeString(additionalLabelValue, true),
		)
		written += n
		if err != nil {
			return written, err
		}
	}
	n, err := out.Write([]byte{'}'})
	written += n
	if err != nil {
		return written, err
	}
	return written, nil
}

var (
	escape                = strings.NewReplacer("\\", `\\`, "\n", `\n`)
	escapeWithDoubleQuote = strings.NewReplacer("\\", `\\`, "\n", `\n`, "\"", `\"`)
)



func escapeString(v string, includeDoubleQuote bool) string {
	if includeDoubleQuote {
		return escapeWithDoubleQuote.Replace(v)
	}

	return escape.Replace(v)
}
