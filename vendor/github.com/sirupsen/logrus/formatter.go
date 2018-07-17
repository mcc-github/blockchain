package logrus

import "time"

const defaultTimestampFormat = time.RFC3339











type Formatter interface {
	Format(*Entry) ([]byte, error)
}













func prefixFieldClashes(data Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}
}
