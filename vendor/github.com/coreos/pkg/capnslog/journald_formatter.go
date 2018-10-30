















package capnslog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coreos/go-systemd/journal"
)

func NewJournaldFormatter() (Formatter, error) {
	if !journal.Enabled() {
		return nil, errors.New("No systemd detected")
	}
	return &journaldFormatter{}, nil
}

type journaldFormatter struct{}

func (j *journaldFormatter) Format(pkg string, l LogLevel, _ int, entries ...interface{}) {
	var pri journal.Priority
	switch l {
	case CRITICAL:
		pri = journal.PriCrit
	case ERROR:
		pri = journal.PriErr
	case WARNING:
		pri = journal.PriWarning
	case NOTICE:
		pri = journal.PriNotice
	case INFO:
		pri = journal.PriInfo
	case DEBUG:
		pri = journal.PriDebug
	case TRACE:
		pri = journal.PriDebug
	default:
		panic("Unhandled loglevel")
	}
	msg := fmt.Sprint(entries...)
	tags := map[string]string{
		"PACKAGE":           pkg,
		"SYSLOG_IDENTIFIER": filepath.Base(os.Args[0]),
	}
	err := journal.Send(msg, pri, tags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (j *journaldFormatter) Flush() {}
