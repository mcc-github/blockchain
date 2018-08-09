/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabenc

import (
	"fmt"
	"io"
	"regexp"
	"runtime"
	"strings"
	"sync/atomic"

	"go.uber.org/zap/zapcore"
)







var formatRegexp = regexp.MustCompile(`%{(color|id|level|message|module|shortfunc|time)(?::(.*?))?}`)



func ParseFormat(spec string) ([]Formatter, error) {
	cursor := 0
	formatters := []Formatter{}

	
	matches := formatRegexp.FindAllStringSubmatchIndex(spec, -1)
	for _, m := range matches {
		start, end := m[0], m[1]
		verbStart, verbEnd := m[2], m[3]
		formatStart, formatEnd := m[4], m[5]

		if start > cursor {
			formatters = append(formatters, StringFormatter{Value: spec[cursor:start]})
		}

		var format string
		if formatStart >= 0 {
			format = spec[formatStart:formatEnd]
		}

		formatter, err := NewFormatter(spec[verbStart:verbEnd], format)
		if err != nil {
			return nil, err
		}

		formatters = append(formatters, formatter)
		cursor = end
	}

	
	if cursor != len(spec) {
		formatters = append(formatters, StringFormatter{Value: spec[cursor:]})
	}

	return formatters, nil
}


type StringFormatter struct{ Value string }


func (s StringFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, "%s", s.Value)
}



func NewFormatter(verb, format string) (Formatter, error) {
	switch verb {
	case "color":
		return newColorFormatter(format)
	case "id":
		return newSequenceFormatter(format), nil
	case "level":
		return newLevelFormatter(format), nil
	case "message":
		return newMessageFormatter(format), nil
	case "module":
		return newModuleFormatter(format), nil
	case "shortfunc":
		return newShortFuncFormatter(format), nil
	case "time":
		return newTimeFormatter(format), nil
	default:
		return nil, fmt.Errorf("unknown verb: %s", verb)
	}
}


type ColorFormatter struct {
	Bold  bool 
	Reset bool 
}

func newColorFormatter(f string) (ColorFormatter, error) {
	switch f {
	case "bold":
		return ColorFormatter{Bold: true}, nil
	case "reset":
		return ColorFormatter{Reset: true}, nil
	case "":
		return ColorFormatter{}, nil
	default:
		return ColorFormatter{}, fmt.Errorf("invalid color option: %s", f)
	}
}


func (c ColorFormatter) LevelColor(l zapcore.Level) Color {
	switch l {
	case zapcore.DebugLevel:
		return ColorCyan
	case zapcore.InfoLevel:
		return ColorBlue
	case zapcore.WarnLevel:
		return ColorYellow
	case zapcore.ErrorLevel:
		return ColorRed
	case zapcore.DPanicLevel, zapcore.PanicLevel:
		return ColorMagenta
	case zapcore.FatalLevel:
		return ColorMagenta
	default:
		return ColorNone
	}
}


func (c ColorFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	switch {
	case c.Reset:
		fmt.Fprintf(w, ResetColor())
	case c.Bold:
		fmt.Fprintf(w, c.LevelColor(entry.Level).Bold())
	default:
		fmt.Fprintf(w, c.LevelColor(entry.Level).Normal())
	}
}


type LevelFormatter struct{ FormatVerb string }

func newLevelFormatter(f string) LevelFormatter {
	return LevelFormatter{FormatVerb: "%" + stringOrDefault(f, "s")}
}


func (l LevelFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, l.FormatVerb, entry.Level.CapitalString())
}


type MessageFormatter struct{ FormatVerb string }

func newMessageFormatter(f string) MessageFormatter {
	return MessageFormatter{FormatVerb: "%" + stringOrDefault(f, "s")}
}


func (m MessageFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, m.FormatVerb, entry.Message)
}


type ModuleFormatter struct{ FormatVerb string }

func newModuleFormatter(f string) ModuleFormatter {
	return ModuleFormatter{FormatVerb: "%" + stringOrDefault(f, "s")}
}


func (m ModuleFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, m.FormatVerb, entry.LoggerName)
}



var sequence uint64


func SetSequence(s uint64) { atomic.StoreUint64(&sequence, s) }


type SequenceFormatter struct{ FormatVerb string }

func newSequenceFormatter(f string) SequenceFormatter {
	return SequenceFormatter{FormatVerb: "%" + stringOrDefault(f, "d")}
}



func (s SequenceFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprintf(w, s.FormatVerb, atomic.AddUint64(&sequence, 1))
}


type ShortFuncFormatter struct{ FormatVerb string }

func newShortFuncFormatter(f string) ShortFuncFormatter {
	return ShortFuncFormatter{FormatVerb: "%" + stringOrDefault(f, "s")}
}



func (s ShortFuncFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	f := runtime.FuncForPC(entry.Caller.PC)
	if f == nil {
		fmt.Fprintf(w, s.FormatVerb, "(unknown)")
		return
	}

	fname := f.Name()
	funcIdx := strings.LastIndex(fname, ".")
	fmt.Fprintf(w, s.FormatVerb, fname[funcIdx+1:])
}


type TimeFormatter struct{ Layout string }

func newTimeFormatter(f string) TimeFormatter {
	return TimeFormatter{Layout: stringOrDefault(f, "2006-01-02T15:04:05.999Z07:00")}
}


func (t TimeFormatter) Format(w io.Writer, entry zapcore.Entry, fields []zapcore.Field) {
	fmt.Fprint(w, entry.Time.Format(t.Layout))
}

func stringOrDefault(str, dflt string) string {
	if str != "" {
		return str
	}
	return dflt
}
