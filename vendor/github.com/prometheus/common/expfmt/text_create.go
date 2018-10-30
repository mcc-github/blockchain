












package expfmt

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/common/model"

	dto "github.com/prometheus/client_model/go"
)



type enhancedWriter interface {
	io.Writer
	WriteRune(r rune) (n int, err error)
	WriteString(s string) (n int, err error)
	WriteByte(c byte) error
}

const (
	initialBufSize    = 512
	initialNumBufSize = 24
)

var (
	bufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, initialBufSize))
		},
	}
	numBufPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 0, initialNumBufSize)
			return &b
		},
	}
)










func MetricFamilyToText(out io.Writer, in *dto.MetricFamily) (written int, err error) {
	
	if len(in.Metric) == 0 {
		return 0, fmt.Errorf("MetricFamily has no metrics: %s", in)
	}
	name := in.GetName()
	if name == "" {
		return 0, fmt.Errorf("MetricFamily has no name: %s", in)
	}

	
	
	
	w, ok := out.(enhancedWriter)
	if !ok {
		b := bufPool.Get().(*bytes.Buffer)
		b.Reset()
		w = b
		defer func() {
			bWritten, bErr := out.Write(b.Bytes())
			written = bWritten
			if err == nil {
				err = bErr
			}
			bufPool.Put(b)
		}()
	}

	var n int

	
	if in.Help != nil {
		n, err = w.WriteString("# HELP ")
		written += n
		if err != nil {
			return
		}
		n, err = w.WriteString(name)
		written += n
		if err != nil {
			return
		}
		err = w.WriteByte(' ')
		written++
		if err != nil {
			return
		}
		n, err = writeEscapedString(w, *in.Help, false)
		written += n
		if err != nil {
			return
		}
		err = w.WriteByte('\n')
		written++
		if err != nil {
			return
		}
	}
	n, err = w.WriteString("# TYPE ")
	written += n
	if err != nil {
		return
	}
	n, err = w.WriteString(name)
	written += n
	if err != nil {
		return
	}
	metricType := in.GetType()
	switch metricType {
	case dto.MetricType_COUNTER:
		n, err = w.WriteString(" counter\n")
	case dto.MetricType_GAUGE:
		n, err = w.WriteString(" gauge\n")
	case dto.MetricType_SUMMARY:
		n, err = w.WriteString(" summary\n")
	case dto.MetricType_UNTYPED:
		n, err = w.WriteString(" untyped\n")
	case dto.MetricType_HISTOGRAM:
		n, err = w.WriteString(" histogram\n")
	default:
		return written, fmt.Errorf("unknown metric type %s", metricType.String())
	}
	written += n
	if err != nil {
		return
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
				w, name, "", metric, "", 0,
				metric.Counter.GetValue(),
			)
		case dto.MetricType_GAUGE:
			if metric.Gauge == nil {
				return written, fmt.Errorf(
					"expected gauge in metric %s %s", name, metric,
				)
			}
			n, err = writeSample(
				w, name, "", metric, "", 0,
				metric.Gauge.GetValue(),
			)
		case dto.MetricType_UNTYPED:
			if metric.Untyped == nil {
				return written, fmt.Errorf(
					"expected untyped in metric %s %s", name, metric,
				)
			}
			n, err = writeSample(
				w, name, "", metric, "", 0,
				metric.Untyped.GetValue(),
			)
		case dto.MetricType_SUMMARY:
			if metric.Summary == nil {
				return written, fmt.Errorf(
					"expected summary in metric %s %s", name, metric,
				)
			}
			for _, q := range metric.Summary.Quantile {
				n, err = writeSample(
					w, name, "", metric,
					model.QuantileLabel, q.GetQuantile(),
					q.GetValue(),
				)
				written += n
				if err != nil {
					return
				}
			}
			n, err = writeSample(
				w, name, "_sum", metric, "", 0,
				metric.Summary.GetSampleSum(),
			)
			written += n
			if err != nil {
				return
			}
			n, err = writeSample(
				w, name, "_count", metric, "", 0,
				float64(metric.Summary.GetSampleCount()),
			)
		case dto.MetricType_HISTOGRAM:
			if metric.Histogram == nil {
				return written, fmt.Errorf(
					"expected histogram in metric %s %s", name, metric,
				)
			}
			infSeen := false
			for _, b := range metric.Histogram.Bucket {
				n, err = writeSample(
					w, name, "_bucket", metric,
					model.BucketLabel, b.GetUpperBound(),
					float64(b.GetCumulativeCount()),
				)
				written += n
				if err != nil {
					return
				}
				if math.IsInf(b.GetUpperBound(), +1) {
					infSeen = true
				}
			}
			if !infSeen {
				n, err = writeSample(
					w, name, "_bucket", metric,
					model.BucketLabel, math.Inf(+1),
					float64(metric.Histogram.GetSampleCount()),
				)
				written += n
				if err != nil {
					return
				}
			}
			n, err = writeSample(
				w, name, "_sum", metric, "", 0,
				metric.Histogram.GetSampleSum(),
			)
			written += n
			if err != nil {
				return
			}
			n, err = writeSample(
				w, name, "_count", metric, "", 0,
				float64(metric.Histogram.GetSampleCount()),
			)
		default:
			return written, fmt.Errorf(
				"unexpected type in metric %s %s", name, metric,
			)
		}
		written += n
		if err != nil {
			return
		}
	}
	return
}






func writeSample(
	w enhancedWriter,
	name, suffix string,
	metric *dto.Metric,
	additionalLabelName string, additionalLabelValue float64,
	value float64,
) (int, error) {
	var written int
	n, err := w.WriteString(name)
	written += n
	if err != nil {
		return written, err
	}
	if suffix != "" {
		n, err = w.WriteString(suffix)
		written += n
		if err != nil {
			return written, err
		}
	}
	n, err = writeLabelPairs(
		w, metric.Label, additionalLabelName, additionalLabelValue,
	)
	written += n
	if err != nil {
		return written, err
	}
	err = w.WriteByte(' ')
	written++
	if err != nil {
		return written, err
	}
	n, err = writeFloat(w, value)
	written += n
	if err != nil {
		return written, err
	}
	if metric.TimestampMs != nil {
		err = w.WriteByte(' ')
		written++
		if err != nil {
			return written, err
		}
		n, err = writeInt(w, *metric.TimestampMs)
		written += n
		if err != nil {
			return written, err
		}
	}
	err = w.WriteByte('\n')
	written++
	if err != nil {
		return written, err
	}
	return written, nil
}








func writeLabelPairs(
	w enhancedWriter,
	in []*dto.LabelPair,
	additionalLabelName string, additionalLabelValue float64,
) (int, error) {
	if len(in) == 0 && additionalLabelName == "" {
		return 0, nil
	}
	var (
		written   int
		separator byte = '{'
	)
	for _, lp := range in {
		err := w.WriteByte(separator)
		written++
		if err != nil {
			return written, err
		}
		n, err := w.WriteString(lp.GetName())
		written += n
		if err != nil {
			return written, err
		}
		n, err = w.WriteString(`="`)
		written += n
		if err != nil {
			return written, err
		}
		n, err = writeEscapedString(w, lp.GetValue(), true)
		written += n
		if err != nil {
			return written, err
		}
		err = w.WriteByte('"')
		written++
		if err != nil {
			return written, err
		}
		separator = ','
	}
	if additionalLabelName != "" {
		err := w.WriteByte(separator)
		written++
		if err != nil {
			return written, err
		}
		n, err := w.WriteString(additionalLabelName)
		written += n
		if err != nil {
			return written, err
		}
		n, err = w.WriteString(`="`)
		written += n
		if err != nil {
			return written, err
		}
		n, err = writeFloat(w, additionalLabelValue)
		written += n
		if err != nil {
			return written, err
		}
		err = w.WriteByte('"')
		written++
		if err != nil {
			return written, err
		}
	}
	err := w.WriteByte('}')
	written++
	if err != nil {
		return written, err
	}
	return written, nil
}



var (
	escaper       = strings.NewReplacer("\\", `\\`, "\n", `\n`)
	quotedEscaper = strings.NewReplacer("\\", `\\`, "\n", `\n`, "\"", `\"`)
)

func writeEscapedString(w enhancedWriter, v string, includeDoubleQuote bool) (int, error) {
	if includeDoubleQuote {
		return quotedEscaper.WriteString(w, v)
	} else {
		return escaper.WriteString(w, v)
	}
}




func writeFloat(w enhancedWriter, f float64) (int, error) {
	switch {
	case f == 1:
		return 1, w.WriteByte('1')
	case f == 0:
		return 1, w.WriteByte('0')
	case f == -1:
		return w.WriteString("-1")
	case math.IsNaN(f):
		return w.WriteString("NaN")
	case math.IsInf(f, +1):
		return w.WriteString("+Inf")
	case math.IsInf(f, -1):
		return w.WriteString("-Inf")
	default:
		bp := numBufPool.Get().(*[]byte)
		*bp = strconv.AppendFloat((*bp)[:0], f, 'g', -1, 64)
		written, err := w.Write(*bp)
		numBufPool.Put(bp)
		return written, err
	}
}




func writeInt(w enhancedWriter, i int64) (int, error) {
	bp := numBufPool.Get().(*[]byte)
	*bp = strconv.AppendInt((*bp)[:0], i, 10)
	written, err := w.Write(*bp)
	numBufPool.Put(bp)
	return written, err
}
