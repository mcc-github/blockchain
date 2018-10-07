












package expfmt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	dto "github.com/prometheus/client_model/go"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/common/model"
)





type stateFn func() stateFn



type ParseError struct {
	Line int
	Msg  string
}


func (e ParseError) Error() string {
	return fmt.Sprintf("text format parsing error in line %d: %s", e.Line, e.Msg)
}



type TextParser struct {
	metricFamiliesByName map[string]*dto.MetricFamily
	buf                  *bufio.Reader 
	err                  error         
	lineCount            int           
	currentByte          byte          
	currentToken         bytes.Buffer  
	currentMF            *dto.MetricFamily
	currentMetric        *dto.Metric
	currentLabelPair     *dto.LabelPair

	
	currentLabels map[string]string 
	
	summaries       map[uint64]*dto.Metric 
	currentQuantile float64
	
	histograms    map[uint64]*dto.Metric 
	currentBucket float64
	
	
	
	currentIsSummaryCount, currentIsSummarySum     bool
	currentIsHistogramCount, currentIsHistogramSum bool
}
























func (p *TextParser) TextToMetricFamilies(in io.Reader) (map[string]*dto.MetricFamily, error) {
	p.reset(in)
	for nextState := p.startOfLine; nextState != nil; nextState = nextState() {
		
	}
	
	for k, mf := range p.metricFamiliesByName {
		if len(mf.GetMetric()) == 0 {
			delete(p.metricFamiliesByName, k)
		}
	}
	
	
	
	
	if p.err == io.EOF {
		p.parseError("unexpected end of input stream")
	}
	return p.metricFamiliesByName, p.err
}

func (p *TextParser) reset(in io.Reader) {
	p.metricFamiliesByName = map[string]*dto.MetricFamily{}
	if p.buf == nil {
		p.buf = bufio.NewReader(in)
	} else {
		p.buf.Reset(in)
	}
	p.err = nil
	p.lineCount = 0
	if p.summaries == nil || len(p.summaries) > 0 {
		p.summaries = map[uint64]*dto.Metric{}
	}
	if p.histograms == nil || len(p.histograms) > 0 {
		p.histograms = map[uint64]*dto.Metric{}
	}
	p.currentQuantile = math.NaN()
	p.currentBucket = math.NaN()
}



func (p *TextParser) startOfLine() stateFn {
	p.lineCount++
	if p.skipBlankTab(); p.err != nil {
		
		
		p.err = nil
		return nil
	}
	switch p.currentByte {
	case '#':
		return p.startComment
	case '\n':
		return p.startOfLine 
	}
	return p.readingMetricName
}



func (p *TextParser) startComment() stateFn {
	if p.skipBlankTab(); p.err != nil {
		return nil 
	}
	if p.currentByte == '\n' {
		return p.startOfLine
	}
	if p.readTokenUntilWhitespace(); p.err != nil {
		return nil 
	}
	
	
	if p.currentByte == '\n' {
		return p.startOfLine
	}
	keyword := p.currentToken.String()
	if keyword != "HELP" && keyword != "TYPE" {
		
		for p.currentByte != '\n' {
			if p.currentByte, p.err = p.buf.ReadByte(); p.err != nil {
				return nil 
			}
		}
		return p.startOfLine
	}
	
	if p.skipBlankTab(); p.err != nil {
		return nil 
	}
	if p.readTokenAsMetricName(); p.err != nil {
		return nil 
	}
	if p.currentByte == '\n' {
		
		
		return p.startOfLine
	}
	if !isBlankOrTab(p.currentByte) {
		p.parseError("invalid metric name in comment")
		return nil
	}
	p.setOrCreateCurrentMF()
	if p.skipBlankTab(); p.err != nil {
		return nil 
	}
	if p.currentByte == '\n' {
		
		
		return p.startOfLine
	}
	switch keyword {
	case "HELP":
		return p.readingHelp
	case "TYPE":
		return p.readingType
	}
	panic(fmt.Sprintf("code error: unexpected keyword %q", keyword))
}



func (p *TextParser) readingMetricName() stateFn {
	if p.readTokenAsMetricName(); p.err != nil {
		return nil
	}
	if p.currentToken.Len() == 0 {
		p.parseError("invalid metric name")
		return nil
	}
	p.setOrCreateCurrentMF()
	
	if p.currentMF.Type == nil {
		p.currentMF.Type = dto.MetricType_UNTYPED.Enum()
	}
	p.currentMetric = &dto.Metric{}
	
	
	
	
	if p.skipBlankTabIfCurrentBlankTab(); p.err != nil {
		return nil 
	}
	return p.readingLabels
}




func (p *TextParser) readingLabels() stateFn {
	
	
	
	if p.currentMF.GetType() == dto.MetricType_SUMMARY || p.currentMF.GetType() == dto.MetricType_HISTOGRAM {
		p.currentLabels = map[string]string{}
		p.currentLabels[string(model.MetricNameLabel)] = p.currentMF.GetName()
		p.currentQuantile = math.NaN()
		p.currentBucket = math.NaN()
	}
	if p.currentByte != '{' {
		return p.readingValue
	}
	return p.startLabelName
}



func (p *TextParser) startLabelName() stateFn {
	if p.skipBlankTab(); p.err != nil {
		return nil 
	}
	if p.currentByte == '}' {
		if p.skipBlankTab(); p.err != nil {
			return nil 
		}
		return p.readingValue
	}
	if p.readTokenAsLabelName(); p.err != nil {
		return nil 
	}
	if p.currentToken.Len() == 0 {
		p.parseError(fmt.Sprintf("invalid label name for metric %q", p.currentMF.GetName()))
		return nil
	}
	p.currentLabelPair = &dto.LabelPair{Name: proto.String(p.currentToken.String())}
	if p.currentLabelPair.GetName() == string(model.MetricNameLabel) {
		p.parseError(fmt.Sprintf("label name %q is reserved", model.MetricNameLabel))
		return nil
	}
	
	
	if !(p.currentMF.GetType() == dto.MetricType_SUMMARY && p.currentLabelPair.GetName() == model.QuantileLabel) &&
		!(p.currentMF.GetType() == dto.MetricType_HISTOGRAM && p.currentLabelPair.GetName() == model.BucketLabel) {
		p.currentMetric.Label = append(p.currentMetric.Label, p.currentLabelPair)
	}
	if p.skipBlankTabIfCurrentBlankTab(); p.err != nil {
		return nil 
	}
	if p.currentByte != '=' {
		p.parseError(fmt.Sprintf("expected '=' after label name, found %q", p.currentByte))
		return nil
	}
	return p.startLabelValue
}



func (p *TextParser) startLabelValue() stateFn {
	if p.skipBlankTab(); p.err != nil {
		return nil 
	}
	if p.currentByte != '"' {
		p.parseError(fmt.Sprintf("expected '\"' at start of label value, found %q", p.currentByte))
		return nil
	}
	if p.readTokenAsLabelValue(); p.err != nil {
		return nil
	}
	if !model.LabelValue(p.currentToken.String()).IsValid() {
		p.parseError(fmt.Sprintf("invalid label value %q", p.currentToken.String()))
		return nil
	}
	p.currentLabelPair.Value = proto.String(p.currentToken.String())
	
	
	
	if p.currentMF.GetType() == dto.MetricType_SUMMARY {
		if p.currentLabelPair.GetName() == model.QuantileLabel {
			if p.currentQuantile, p.err = strconv.ParseFloat(p.currentLabelPair.GetValue(), 64); p.err != nil {
				
				p.parseError(fmt.Sprintf("expected float as value for 'quantile' label, got %q", p.currentLabelPair.GetValue()))
				return nil
			}
		} else {
			p.currentLabels[p.currentLabelPair.GetName()] = p.currentLabelPair.GetValue()
		}
	}
	
	if p.currentMF.GetType() == dto.MetricType_HISTOGRAM {
		if p.currentLabelPair.GetName() == model.BucketLabel {
			if p.currentBucket, p.err = strconv.ParseFloat(p.currentLabelPair.GetValue(), 64); p.err != nil {
				
				p.parseError(fmt.Sprintf("expected float as value for 'le' label, got %q", p.currentLabelPair.GetValue()))
				return nil
			}
		} else {
			p.currentLabels[p.currentLabelPair.GetName()] = p.currentLabelPair.GetValue()
		}
	}
	if p.skipBlankTab(); p.err != nil {
		return nil 
	}
	switch p.currentByte {
	case ',':
		return p.startLabelName

	case '}':
		if p.skipBlankTab(); p.err != nil {
			return nil 
		}
		return p.readingValue
	default:
		p.parseError(fmt.Sprintf("unexpected end of label value %q", p.currentLabelPair.GetValue()))
		return nil
	}
}



func (p *TextParser) readingValue() stateFn {
	
	
	
	if p.currentMF.GetType() == dto.MetricType_SUMMARY {
		signature := model.LabelsToSignature(p.currentLabels)
		if summary := p.summaries[signature]; summary != nil {
			p.currentMetric = summary
		} else {
			p.summaries[signature] = p.currentMetric
			p.currentMF.Metric = append(p.currentMF.Metric, p.currentMetric)
		}
	} else if p.currentMF.GetType() == dto.MetricType_HISTOGRAM {
		signature := model.LabelsToSignature(p.currentLabels)
		if histogram := p.histograms[signature]; histogram != nil {
			p.currentMetric = histogram
		} else {
			p.histograms[signature] = p.currentMetric
			p.currentMF.Metric = append(p.currentMF.Metric, p.currentMetric)
		}
	} else {
		p.currentMF.Metric = append(p.currentMF.Metric, p.currentMetric)
	}
	if p.readTokenUntilWhitespace(); p.err != nil {
		return nil 
	}
	value, err := strconv.ParseFloat(p.currentToken.String(), 64)
	if err != nil {
		
		p.parseError(fmt.Sprintf("expected float as value, got %q", p.currentToken.String()))
		return nil
	}
	switch p.currentMF.GetType() {
	case dto.MetricType_COUNTER:
		p.currentMetric.Counter = &dto.Counter{Value: proto.Float64(value)}
	case dto.MetricType_GAUGE:
		p.currentMetric.Gauge = &dto.Gauge{Value: proto.Float64(value)}
	case dto.MetricType_UNTYPED:
		p.currentMetric.Untyped = &dto.Untyped{Value: proto.Float64(value)}
	case dto.MetricType_SUMMARY:
		
		if p.currentMetric.Summary == nil {
			p.currentMetric.Summary = &dto.Summary{}
		}
		switch {
		case p.currentIsSummaryCount:
			p.currentMetric.Summary.SampleCount = proto.Uint64(uint64(value))
		case p.currentIsSummarySum:
			p.currentMetric.Summary.SampleSum = proto.Float64(value)
		case !math.IsNaN(p.currentQuantile):
			p.currentMetric.Summary.Quantile = append(
				p.currentMetric.Summary.Quantile,
				&dto.Quantile{
					Quantile: proto.Float64(p.currentQuantile),
					Value:    proto.Float64(value),
				},
			)
		}
	case dto.MetricType_HISTOGRAM:
		
		if p.currentMetric.Histogram == nil {
			p.currentMetric.Histogram = &dto.Histogram{}
		}
		switch {
		case p.currentIsHistogramCount:
			p.currentMetric.Histogram.SampleCount = proto.Uint64(uint64(value))
		case p.currentIsHistogramSum:
			p.currentMetric.Histogram.SampleSum = proto.Float64(value)
		case !math.IsNaN(p.currentBucket):
			p.currentMetric.Histogram.Bucket = append(
				p.currentMetric.Histogram.Bucket,
				&dto.Bucket{
					UpperBound:      proto.Float64(p.currentBucket),
					CumulativeCount: proto.Uint64(uint64(value)),
				},
			)
		}
	default:
		p.err = fmt.Errorf("unexpected type for metric name %q", p.currentMF.GetName())
	}
	if p.currentByte == '\n' {
		return p.startOfLine
	}
	return p.startTimestamp
}



func (p *TextParser) startTimestamp() stateFn {
	if p.skipBlankTab(); p.err != nil {
		return nil 
	}
	if p.readTokenUntilWhitespace(); p.err != nil {
		return nil 
	}
	timestamp, err := strconv.ParseInt(p.currentToken.String(), 10, 64)
	if err != nil {
		
		p.parseError(fmt.Sprintf("expected integer as timestamp, got %q", p.currentToken.String()))
		return nil
	}
	p.currentMetric.TimestampMs = proto.Int64(timestamp)
	if p.readTokenUntilNewline(false); p.err != nil {
		return nil 
	}
	if p.currentToken.Len() > 0 {
		p.parseError(fmt.Sprintf("spurious string after timestamp: %q", p.currentToken.String()))
		return nil
	}
	return p.startOfLine
}



func (p *TextParser) readingHelp() stateFn {
	if p.currentMF.Help != nil {
		p.parseError(fmt.Sprintf("second HELP line for metric name %q", p.currentMF.GetName()))
		return nil
	}
	
	if p.readTokenUntilNewline(true); p.err != nil {
		return nil 
	}
	p.currentMF.Help = proto.String(p.currentToken.String())
	return p.startOfLine
}



func (p *TextParser) readingType() stateFn {
	if p.currentMF.Type != nil {
		p.parseError(fmt.Sprintf("second TYPE line for metric name %q, or TYPE reported after samples", p.currentMF.GetName()))
		return nil
	}
	
	if p.readTokenUntilNewline(false); p.err != nil {
		return nil 
	}
	metricType, ok := dto.MetricType_value[strings.ToUpper(p.currentToken.String())]
	if !ok {
		p.parseError(fmt.Sprintf("unknown metric type %q", p.currentToken.String()))
		return nil
	}
	p.currentMF.Type = dto.MetricType(metricType).Enum()
	return p.startOfLine
}



func (p *TextParser) parseError(msg string) {
	p.err = ParseError{
		Line: p.lineCount,
		Msg:  msg,
	}
}



func (p *TextParser) skipBlankTab() {
	for {
		if p.currentByte, p.err = p.buf.ReadByte(); p.err != nil || !isBlankOrTab(p.currentByte) {
			return
		}
	}
}



func (p *TextParser) skipBlankTabIfCurrentBlankTab() {
	if isBlankOrTab(p.currentByte) {
		p.skipBlankTab()
	}
}





func (p *TextParser) readTokenUntilWhitespace() {
	p.currentToken.Reset()
	for p.err == nil && !isBlankOrTab(p.currentByte) && p.currentByte != '\n' {
		p.currentToken.WriteByte(p.currentByte)
		p.currentByte, p.err = p.buf.ReadByte()
	}
}







func (p *TextParser) readTokenUntilNewline(recognizeEscapeSequence bool) {
	p.currentToken.Reset()
	escaped := false
	for p.err == nil {
		if recognizeEscapeSequence && escaped {
			switch p.currentByte {
			case '\\':
				p.currentToken.WriteByte(p.currentByte)
			case 'n':
				p.currentToken.WriteByte('\n')
			default:
				p.parseError(fmt.Sprintf("invalid escape sequence '\\%c'", p.currentByte))
				return
			}
			escaped = false
		} else {
			switch p.currentByte {
			case '\n':
				return
			case '\\':
				escaped = true
			default:
				p.currentToken.WriteByte(p.currentByte)
			}
		}
		p.currentByte, p.err = p.buf.ReadByte()
	}
}





func (p *TextParser) readTokenAsMetricName() {
	p.currentToken.Reset()
	if !isValidMetricNameStart(p.currentByte) {
		return
	}
	for {
		p.currentToken.WriteByte(p.currentByte)
		p.currentByte, p.err = p.buf.ReadByte()
		if p.err != nil || !isValidMetricNameContinuation(p.currentByte) {
			return
		}
	}
}





func (p *TextParser) readTokenAsLabelName() {
	p.currentToken.Reset()
	if !isValidLabelNameStart(p.currentByte) {
		return
	}
	for {
		p.currentToken.WriteByte(p.currentByte)
		p.currentByte, p.err = p.buf.ReadByte()
		if p.err != nil || !isValidLabelNameContinuation(p.currentByte) {
			return
		}
	}
}






func (p *TextParser) readTokenAsLabelValue() {
	p.currentToken.Reset()
	escaped := false
	for {
		if p.currentByte, p.err = p.buf.ReadByte(); p.err != nil {
			return
		}
		if escaped {
			switch p.currentByte {
			case '"', '\\':
				p.currentToken.WriteByte(p.currentByte)
			case 'n':
				p.currentToken.WriteByte('\n')
			default:
				p.parseError(fmt.Sprintf("invalid escape sequence '\\%c'", p.currentByte))
				return
			}
			escaped = false
			continue
		}
		switch p.currentByte {
		case '"':
			return
		case '\n':
			p.parseError(fmt.Sprintf("label value %q contains unescaped new-line", p.currentToken.String()))
			return
		case '\\':
			escaped = true
		default:
			p.currentToken.WriteByte(p.currentByte)
		}
	}
}

func (p *TextParser) setOrCreateCurrentMF() {
	p.currentIsSummaryCount = false
	p.currentIsSummarySum = false
	p.currentIsHistogramCount = false
	p.currentIsHistogramSum = false
	name := p.currentToken.String()
	if p.currentMF = p.metricFamiliesByName[name]; p.currentMF != nil {
		return
	}
	
	summaryName := summaryMetricName(name)
	if p.currentMF = p.metricFamiliesByName[summaryName]; p.currentMF != nil {
		if p.currentMF.GetType() == dto.MetricType_SUMMARY {
			if isCount(name) {
				p.currentIsSummaryCount = true
			}
			if isSum(name) {
				p.currentIsSummarySum = true
			}
			return
		}
	}
	histogramName := histogramMetricName(name)
	if p.currentMF = p.metricFamiliesByName[histogramName]; p.currentMF != nil {
		if p.currentMF.GetType() == dto.MetricType_HISTOGRAM {
			if isCount(name) {
				p.currentIsHistogramCount = true
			}
			if isSum(name) {
				p.currentIsHistogramSum = true
			}
			return
		}
	}
	p.currentMF = &dto.MetricFamily{Name: proto.String(name)}
	p.metricFamiliesByName[name] = p.currentMF
}

func isValidLabelNameStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isValidLabelNameContinuation(b byte) bool {
	return isValidLabelNameStart(b) || (b >= '0' && b <= '9')
}

func isValidMetricNameStart(b byte) bool {
	return isValidLabelNameStart(b) || b == ':'
}

func isValidMetricNameContinuation(b byte) bool {
	return isValidLabelNameContinuation(b) || b == ':'
}

func isBlankOrTab(b byte) bool {
	return b == ' ' || b == '\t'
}

func isCount(name string) bool {
	return len(name) > 6 && name[len(name)-6:] == "_count"
}

func isSum(name string) bool {
	return len(name) > 4 && name[len(name)-4:] == "_sum"
}

func isBucket(name string) bool {
	return len(name) > 7 && name[len(name)-7:] == "_bucket"
}

func summaryMetricName(name string) string {
	switch {
	case isCount(name):
		return name[:len(name)-6]
	case isSum(name):
		return name[:len(name)-4]
	default:
		return name
	}
}

func histogramMetricName(name string) string {
	switch {
	case isCount(name):
		return name[:len(name)-6]
	case isSum(name):
		return name[:len(name)-4]
	case isBucket(name):
		return name[:len(name)-7]
	default:
		return name
	}
}
