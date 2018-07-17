





package logging

import "log/syslog"



type SyslogBackend struct {
	Writer *syslog.Writer
}




func NewSyslogBackend(prefix string) (b *SyslogBackend, err error) {
	var w *syslog.Writer
	w, err = syslog.New(syslog.LOG_CRIT, prefix)
	return &SyslogBackend{w}, err
}



func NewSyslogBackendPriority(prefix string, priority syslog.Priority) (b *SyslogBackend, err error) {
	var w *syslog.Writer
	w, err = syslog.New(priority, prefix)
	return &SyslogBackend{w}, err
}


func (b *SyslogBackend) Log(level Level, calldepth int, rec *Record) error {
	line := rec.Formatted(calldepth + 1)
	switch level {
	case CRITICAL:
		return b.Writer.Crit(line)
	case ERROR:
		return b.Writer.Err(line)
	case WARNING:
		return b.Writer.Warning(line)
	case NOTICE:
		return b.Writer.Notice(line)
	case INFO:
		return b.Writer.Info(line)
	case DEBUG:
		return b.Writer.Debug(line)
	default:
	}
	panic("unhandled log level")
}
