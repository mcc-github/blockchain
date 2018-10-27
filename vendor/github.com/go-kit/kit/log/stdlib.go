package log

import (
	"io"
	"log"
	"regexp"
	"strings"
)







type StdlibWriter struct{}


func (w StdlibWriter) Write(p []byte) (int, error) {
	log.Print(strings.TrimSpace(string(p)))
	return len(p), nil
}




type StdlibAdapter struct {
	Logger
	timestampKey string
	fileKey      string
	messageKey   string
}


type StdlibAdapterOption func(*StdlibAdapter)


func TimestampKey(key string) StdlibAdapterOption {
	return func(a *StdlibAdapter) { a.timestampKey = key }
}


func FileKey(key string) StdlibAdapterOption {
	return func(a *StdlibAdapter) { a.fileKey = key }
}


func MessageKey(key string) StdlibAdapterOption {
	return func(a *StdlibAdapter) { a.messageKey = key }
}



func NewStdlibAdapter(logger Logger, options ...StdlibAdapterOption) io.Writer {
	a := StdlibAdapter{
		Logger:       logger,
		timestampKey: "ts",
		fileKey:      "caller",
		messageKey:   "msg",
	}
	for _, option := range options {
		option(&a)
	}
	return a
}

func (a StdlibAdapter) Write(p []byte) (int, error) {
	result := subexps(p)
	keyvals := []interface{}{}
	var timestamp string
	if date, ok := result["date"]; ok && date != "" {
		timestamp = date
	}
	if time, ok := result["time"]; ok && time != "" {
		if timestamp != "" {
			timestamp += " "
		}
		timestamp += time
	}
	if timestamp != "" {
		keyvals = append(keyvals, a.timestampKey, timestamp)
	}
	if file, ok := result["file"]; ok && file != "" {
		keyvals = append(keyvals, a.fileKey, file)
	}
	if msg, ok := result["msg"]; ok {
		keyvals = append(keyvals, a.messageKey, msg)
	}
	if err := a.Logger.Log(keyvals...); err != nil {
		return 0, err
	}
	return len(p), nil
}

const (
	logRegexpDate = `(?P<date>[0-9]{4}/[0-9]{2}/[0-9]{2})?[ ]?`
	logRegexpTime = `(?P<time>[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?)?[ ]?`
	logRegexpFile = `(?P<file>.+?:[0-9]+)?`
	logRegexpMsg  = `(: )?(?P<msg>.*)`
)

var (
	logRegexp = regexp.MustCompile(logRegexpDate + logRegexpTime + logRegexpFile + logRegexpMsg)
)

func subexps(line []byte) map[string]string {
	m := logRegexp.FindSubmatch(line)
	if len(m) < len(logRegexp.SubexpNames()) {
		return map[string]string{}
	}
	result := map[string]string{}
	for i, name := range logRegexp.SubexpNames() {
		result[name] = string(m[i])
	}
	return result
}
