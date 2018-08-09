/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabenc_test

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/common/flogging/fabenc"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestParseFormat(t *testing.T) {
	var tests = []struct {
		spec       string
		formatters []fabenc.Formatter
	}{
		{spec: "", formatters: []fabenc.Formatter{}},
		{
			spec: "%{color}",
			formatters: []fabenc.Formatter{
				fabenc.ColorFormatter{},
			},
		},
		{
			spec: "prefix %{color}",
			formatters: []fabenc.Formatter{
				fabenc.StringFormatter{Value: "prefix "},
				fabenc.ColorFormatter{},
			},
		},
		{
			spec: "%{color} suffix",
			formatters: []fabenc.Formatter{
				fabenc.ColorFormatter{},
				fabenc.StringFormatter{Value: " suffix"}},
		},
		{
			spec: "prefix %{color} suffix",
			formatters: []fabenc.Formatter{
				fabenc.StringFormatter{Value: "prefix "},
				fabenc.ColorFormatter{},
				fabenc.StringFormatter{Value: " suffix"},
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			formatters, err := fabenc.ParseFormat(tc.spec)
			assert.NoError(t, err)
			assert.Equal(t, tc.formatters, formatters)
		})
	}
}

func TestNewFormatter(t *testing.T) {
	var tests = []struct {
		verb      string
		format    string
		formatter fabenc.Formatter
		errorMsg  string
	}{
		{verb: "color", format: "", formatter: fabenc.ColorFormatter{}},
		{verb: "color", format: "bold", formatter: fabenc.ColorFormatter{Bold: true}},
		{verb: "color", format: "reset", formatter: fabenc.ColorFormatter{Reset: true}},
		{verb: "color", format: "unknown", errorMsg: "invalid color option: unknown"},
		{verb: "id", format: "", formatter: fabenc.SequenceFormatter{FormatVerb: "%d"}},
		{verb: "id", format: "04x", formatter: fabenc.SequenceFormatter{FormatVerb: "%04x"}},
		{verb: "level", format: "", formatter: fabenc.LevelFormatter{FormatVerb: "%s"}},
		{verb: "level", format: ".4s", formatter: fabenc.LevelFormatter{FormatVerb: "%.4s"}},
		{verb: "message", format: "", formatter: fabenc.MessageFormatter{FormatVerb: "%s"}},
		{verb: "message", format: "#30s", formatter: fabenc.MessageFormatter{FormatVerb: "%#30s"}},
		{verb: "module", format: "", formatter: fabenc.ModuleFormatter{FormatVerb: "%s"}},
		{verb: "module", format: "ok", formatter: fabenc.ModuleFormatter{FormatVerb: "%ok"}},
		{verb: "shortfunc", format: "", formatter: fabenc.ShortFuncFormatter{FormatVerb: "%s"}},
		{verb: "shortfunc", format: "U", formatter: fabenc.ShortFuncFormatter{FormatVerb: "%U"}},
		{verb: "time", format: "", formatter: fabenc.TimeFormatter{Layout: "2006-01-02T15:04:05.999Z07:00"}},
		{verb: "time", format: "04:05.999999Z05:00", formatter: fabenc.TimeFormatter{Layout: "04:05.999999Z05:00"}},
		{verb: "unknown", format: "", errorMsg: "unknown verb: unknown"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			f, err := fabenc.NewFormatter(tc.verb, tc.format)
			if tc.errorMsg == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.formatter, f)
			} else {
				assert.EqualError(t, err, tc.errorMsg)
			}
		})
	}
}

func TestColorFormatter(t *testing.T) {
	var tests = []struct {
		f         fabenc.ColorFormatter
		level     zapcore.Level
		formatted string
	}{
		{f: fabenc.ColorFormatter{Reset: true}, level: zapcore.DebugLevel, formatted: fabenc.ResetColor()},
		{f: fabenc.ColorFormatter{}, level: zapcore.DebugLevel, formatted: fabenc.ColorCyan.Normal()},
		{f: fabenc.ColorFormatter{Bold: true}, level: zapcore.DebugLevel, formatted: fabenc.ColorCyan.Bold()},
		{f: fabenc.ColorFormatter{}, level: zapcore.InfoLevel, formatted: fabenc.ColorBlue.Normal()},
		{f: fabenc.ColorFormatter{Bold: true}, level: zapcore.InfoLevel, formatted: fabenc.ColorBlue.Bold()},
		{f: fabenc.ColorFormatter{}, level: zapcore.WarnLevel, formatted: fabenc.ColorYellow.Normal()},
		{f: fabenc.ColorFormatter{Bold: true}, level: zapcore.WarnLevel, formatted: fabenc.ColorYellow.Bold()},
		{f: fabenc.ColorFormatter{}, level: zapcore.ErrorLevel, formatted: fabenc.ColorRed.Normal()},
		{f: fabenc.ColorFormatter{Bold: true}, level: zapcore.ErrorLevel, formatted: fabenc.ColorRed.Bold()},
		{f: fabenc.ColorFormatter{}, level: zapcore.DPanicLevel, formatted: fabenc.ColorMagenta.Normal()},
		{f: fabenc.ColorFormatter{Bold: true}, level: zapcore.DPanicLevel, formatted: fabenc.ColorMagenta.Bold()},
		{f: fabenc.ColorFormatter{}, level: zapcore.PanicLevel, formatted: fabenc.ColorMagenta.Normal()},
		{f: fabenc.ColorFormatter{Bold: true}, level: zapcore.PanicLevel, formatted: fabenc.ColorMagenta.Bold()},
		{f: fabenc.ColorFormatter{}, level: zapcore.FatalLevel, formatted: fabenc.ColorMagenta.Normal()},
		{f: fabenc.ColorFormatter{Bold: true}, level: zapcore.FatalLevel, formatted: fabenc.ColorMagenta.Bold()},
		{f: fabenc.ColorFormatter{}, level: zapcore.Level(99), formatted: fabenc.ColorNone.Normal()},
		{f: fabenc.ColorFormatter{Bold: true}, level: zapcore.Level(99), formatted: fabenc.ColorNone.Normal()},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			buf := &bytes.Buffer{}
			entry := zapcore.Entry{Level: tc.level}
			tc.f.Format(buf, entry, nil)
			assert.Equal(t, tc.formatted, buf.String())
		})
	}
}

func TestLevelFormatter(t *testing.T) {
	var tests = []struct {
		level     zapcore.Level
		formatted string
	}{
		{level: zapcore.DebugLevel, formatted: "DEBUG"},
		{level: zapcore.InfoLevel, formatted: "INFO"},
		{level: zapcore.WarnLevel, formatted: "WARN"},
		{level: zapcore.ErrorLevel, formatted: "ERROR"},
		{level: zapcore.DPanicLevel, formatted: "DPANIC"},
		{level: zapcore.PanicLevel, formatted: "PANIC"},
		{level: zapcore.FatalLevel, formatted: "FATAL"},
		{level: zapcore.Level(99), formatted: "LEVEL(99)"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			buf := &bytes.Buffer{}
			entry := zapcore.Entry{Level: tc.level}
			fabenc.LevelFormatter{FormatVerb: "%s"}.Format(buf, entry, nil)
			assert.Equal(t, tc.formatted, buf.String())
		})
	}
}

func TestMessageFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	entry := zapcore.Entry{Message: "some message text"}
	f := fabenc.MessageFormatter{FormatVerb: "%s"}
	f.Format(buf, entry, nil)
	assert.Equal(t, "some message text", buf.String())
}

func TestModuleFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	entry := zapcore.Entry{LoggerName: "logger/name"}
	f := fabenc.ModuleFormatter{FormatVerb: "%s"}
	f.Format(buf, entry, nil)
	assert.Equal(t, "logger/name", buf.String())
}

func TestSequenceFormatter(t *testing.T) {
	mutex := &sync.Mutex{}
	results := map[string]struct{}{}

	ready := &sync.WaitGroup{}
	ready.Add(100)

	finished := &sync.WaitGroup{}
	finished.Add(100)

	fabenc.SetSequence(0)
	for i := 1; i <= 100; i++ {
		go func(i int) {
			buf := &bytes.Buffer{}
			entry := zapcore.Entry{Level: zapcore.DebugLevel}
			f := fabenc.SequenceFormatter{FormatVerb: "%d"}
			ready.Done() 
			ready.Wait() 

			f.Format(buf, entry, nil) 

			mutex.Lock()
			results[buf.String()] = struct{}{}
			mutex.Unlock()

			finished.Done()
		}(i)
	}

	finished.Wait()
	for i := 1; i <= 100; i++ {
		assert.Contains(t, results, strconv.Itoa(i))
	}
}

func TestShortFuncFormatter(t *testing.T) {
	callerpc, _, _, ok := runtime.Caller(0)
	assert.True(t, ok)
	buf := &bytes.Buffer{}
	entry := zapcore.Entry{Caller: zapcore.EntryCaller{PC: callerpc}}
	fabenc.ShortFuncFormatter{FormatVerb: "%s"}.Format(buf, entry, nil)
	assert.Equal(t, "TestShortFuncFormatter", buf.String())

	buf = &bytes.Buffer{}
	entry = zapcore.Entry{Caller: zapcore.EntryCaller{PC: 0}}
	fabenc.ShortFuncFormatter{FormatVerb: "%s"}.Format(buf, entry, nil)
	assert.Equal(t, "(unknown)", buf.String())
}

func TestTimeFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	entry := zapcore.Entry{Time: time.Date(1975, time.August, 15, 12, 0, 0, 333, time.UTC)}
	f := fabenc.TimeFormatter{Layout: time.RFC3339Nano}
	f.Format(buf, entry, nil)
	assert.Equal(t, "1975-08-15T12:00:00.000000333Z", buf.String())
}
