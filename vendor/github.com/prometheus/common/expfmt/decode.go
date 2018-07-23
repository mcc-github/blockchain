












package expfmt

import (
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"

	dto "github.com/prometheus/client_model/go"

	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	"github.com/prometheus/common/model"
)


type Decoder interface {
	Decode(*dto.MetricFamily) error
}


type DecodeOptions struct {
	
	Timestamp model.Time
}



func ResponseFormat(h http.Header) Format {
	ct := h.Get(hdrContentType)

	mediatype, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return FmtUnknown
	}

	const textType = "text/plain"

	switch mediatype {
	case ProtoType:
		if p, ok := params["proto"]; ok && p != ProtoProtocol {
			return FmtUnknown
		}
		if e, ok := params["encoding"]; ok && e != "delimited" {
			return FmtUnknown
		}
		return FmtProtoDelim

	case textType:
		if v, ok := params["version"]; ok && v != TextVersion {
			return FmtUnknown
		}
		return FmtText
	}

	return FmtUnknown
}



func NewDecoder(r io.Reader, format Format) Decoder {
	switch format {
	case FmtProtoDelim:
		return &protoDecoder{r: r}
	}
	return &textDecoder{r: r}
}


type protoDecoder struct {
	r io.Reader
}


func (d *protoDecoder) Decode(v *dto.MetricFamily) error {
	_, err := pbutil.ReadDelimited(d.r, v)
	if err != nil {
		return err
	}
	if !model.IsValidMetricName(model.LabelValue(v.GetName())) {
		return fmt.Errorf("invalid metric name %q", v.GetName())
	}
	for _, m := range v.GetMetric() {
		if m == nil {
			continue
		}
		for _, l := range m.GetLabel() {
			if l == nil {
				continue
			}
			if !model.LabelValue(l.GetValue()).IsValid() {
				return fmt.Errorf("invalid label value %q", l.GetValue())
			}
			if !model.LabelName(l.GetName()).IsValid() {
				return fmt.Errorf("invalid label name %q", l.GetName())
			}
		}
	}
	return nil
}


type textDecoder struct {
	r    io.Reader
	p    TextParser
	fams []*dto.MetricFamily
}


func (d *textDecoder) Decode(v *dto.MetricFamily) error {
	
	if len(d.fams) == 0 {
		
		fams, err := d.p.TextToMetricFamilies(d.r)
		if err != nil {
			return err
		}
		if len(fams) == 0 {
			return io.EOF
		}
		d.fams = make([]*dto.MetricFamily, 0, len(fams))
		for _, f := range fams {
			d.fams = append(d.fams, f)
		}
	}

	*v = *d.fams[0]
	d.fams = d.fams[1:]

	return nil
}



type SampleDecoder struct {
	Dec  Decoder
	Opts *DecodeOptions

	f dto.MetricFamily
}



func (sd *SampleDecoder) Decode(s *model.Vector) error {
	err := sd.Dec.Decode(&sd.f)
	if err != nil {
		return err
	}
	*s, err = extractSamples(&sd.f, sd.Opts)
	return err
}





func ExtractSamples(o *DecodeOptions, fams ...*dto.MetricFamily) (model.Vector, error) {
	var (
		all     model.Vector
		lastErr error
	)
	for _, f := range fams {
		some, err := extractSamples(f, o)
		if err != nil {
			lastErr = err
			continue
		}
		all = append(all, some...)
	}
	return all, lastErr
}

func extractSamples(f *dto.MetricFamily, o *DecodeOptions) (model.Vector, error) {
	switch f.GetType() {
	case dto.MetricType_COUNTER:
		return extractCounter(o, f), nil
	case dto.MetricType_GAUGE:
		return extractGauge(o, f), nil
	case dto.MetricType_SUMMARY:
		return extractSummary(o, f), nil
	case dto.MetricType_UNTYPED:
		return extractUntyped(o, f), nil
	case dto.MetricType_HISTOGRAM:
		return extractHistogram(o, f), nil
	}
	return nil, fmt.Errorf("expfmt.extractSamples: unknown metric family type %v", f.GetType())
}

func extractCounter(o *DecodeOptions, f *dto.MetricFamily) model.Vector {
	samples := make(model.Vector, 0, len(f.Metric))

	for _, m := range f.Metric {
		if m.Counter == nil {
			continue
		}

		lset := make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label {
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		}
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName())

		smpl := &model.Sample{
			Metric: model.Metric(lset),
			Value:  model.SampleValue(m.Counter.GetValue()),
		}

		if m.TimestampMs != nil {
			smpl.Timestamp = model.TimeFromUnixNano(*m.TimestampMs * 1000000)
		} else {
			smpl.Timestamp = o.Timestamp
		}

		samples = append(samples, smpl)
	}

	return samples
}

func extractGauge(o *DecodeOptions, f *dto.MetricFamily) model.Vector {
	samples := make(model.Vector, 0, len(f.Metric))

	for _, m := range f.Metric {
		if m.Gauge == nil {
			continue
		}

		lset := make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label {
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		}
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName())

		smpl := &model.Sample{
			Metric: model.Metric(lset),
			Value:  model.SampleValue(m.Gauge.GetValue()),
		}

		if m.TimestampMs != nil {
			smpl.Timestamp = model.TimeFromUnixNano(*m.TimestampMs * 1000000)
		} else {
			smpl.Timestamp = o.Timestamp
		}

		samples = append(samples, smpl)
	}

	return samples
}

func extractUntyped(o *DecodeOptions, f *dto.MetricFamily) model.Vector {
	samples := make(model.Vector, 0, len(f.Metric))

	for _, m := range f.Metric {
		if m.Untyped == nil {
			continue
		}

		lset := make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label {
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		}
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName())

		smpl := &model.Sample{
			Metric: model.Metric(lset),
			Value:  model.SampleValue(m.Untyped.GetValue()),
		}

		if m.TimestampMs != nil {
			smpl.Timestamp = model.TimeFromUnixNano(*m.TimestampMs * 1000000)
		} else {
			smpl.Timestamp = o.Timestamp
		}

		samples = append(samples, smpl)
	}

	return samples
}

func extractSummary(o *DecodeOptions, f *dto.MetricFamily) model.Vector {
	samples := make(model.Vector, 0, len(f.Metric))

	for _, m := range f.Metric {
		if m.Summary == nil {
			continue
		}

		timestamp := o.Timestamp
		if m.TimestampMs != nil {
			timestamp = model.TimeFromUnixNano(*m.TimestampMs * 1000000)
		}

		for _, q := range m.Summary.Quantile {
			lset := make(model.LabelSet, len(m.Label)+2)
			for _, p := range m.Label {
				lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
			}
			
			lset[model.LabelName(model.QuantileLabel)] = model.LabelValue(fmt.Sprint(q.GetQuantile()))
			lset[model.MetricNameLabel] = model.LabelValue(f.GetName())

			samples = append(samples, &model.Sample{
				Metric:    model.Metric(lset),
				Value:     model.SampleValue(q.GetValue()),
				Timestamp: timestamp,
			})
		}

		lset := make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label {
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		}
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_sum")

		samples = append(samples, &model.Sample{
			Metric:    model.Metric(lset),
			Value:     model.SampleValue(m.Summary.GetSampleSum()),
			Timestamp: timestamp,
		})

		lset = make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label {
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		}
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_count")

		samples = append(samples, &model.Sample{
			Metric:    model.Metric(lset),
			Value:     model.SampleValue(m.Summary.GetSampleCount()),
			Timestamp: timestamp,
		})
	}

	return samples
}

func extractHistogram(o *DecodeOptions, f *dto.MetricFamily) model.Vector {
	samples := make(model.Vector, 0, len(f.Metric))

	for _, m := range f.Metric {
		if m.Histogram == nil {
			continue
		}

		timestamp := o.Timestamp
		if m.TimestampMs != nil {
			timestamp = model.TimeFromUnixNano(*m.TimestampMs * 1000000)
		}

		infSeen := false

		for _, q := range m.Histogram.Bucket {
			lset := make(model.LabelSet, len(m.Label)+2)
			for _, p := range m.Label {
				lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
			}
			lset[model.LabelName(model.BucketLabel)] = model.LabelValue(fmt.Sprint(q.GetUpperBound()))
			lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_bucket")

			if math.IsInf(q.GetUpperBound(), +1) {
				infSeen = true
			}

			samples = append(samples, &model.Sample{
				Metric:    model.Metric(lset),
				Value:     model.SampleValue(q.GetCumulativeCount()),
				Timestamp: timestamp,
			})
		}

		lset := make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label {
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		}
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_sum")

		samples = append(samples, &model.Sample{
			Metric:    model.Metric(lset),
			Value:     model.SampleValue(m.Histogram.GetSampleSum()),
			Timestamp: timestamp,
		})

		lset = make(model.LabelSet, len(m.Label)+1)
		for _, p := range m.Label {
			lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
		}
		lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_count")

		count := &model.Sample{
			Metric:    model.Metric(lset),
			Value:     model.SampleValue(m.Histogram.GetSampleCount()),
			Timestamp: timestamp,
		}
		samples = append(samples, count)

		if !infSeen {
			
			lset := make(model.LabelSet, len(m.Label)+2)
			for _, p := range m.Label {
				lset[model.LabelName(p.GetName())] = model.LabelValue(p.GetValue())
			}
			lset[model.LabelName(model.BucketLabel)] = model.LabelValue("+Inf")
			lset[model.MetricNameLabel] = model.LabelValue(f.GetName() + "_bucket")

			samples = append(samples, &model.Sample{
				Metric:    model.Metric(lset),
				Value:     count.Value,
				Timestamp: timestamp,
			})
		}
	}

	return samples
}