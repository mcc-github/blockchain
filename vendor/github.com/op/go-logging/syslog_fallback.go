





package logging

import (
	"fmt"
)

type Priority int

type SyslogBackend struct {
}

func NewSyslogBackend(prefix string) (b *SyslogBackend, err error) {
	return nil, fmt.Errorf("Platform does not support syslog")
}

func NewSyslogBackendPriority(prefix string, priority Priority) (b *SyslogBackend, err error) {
	return nil, fmt.Errorf("Platform does not support syslog")
}

func (b *SyslogBackend) Log(level Level, calldepth int, rec *Record) error {
	return fmt.Errorf("Platform does not support syslog")
}
