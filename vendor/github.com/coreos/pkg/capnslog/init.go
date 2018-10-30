















package capnslog

import (
	"io"
	"os"
	"syscall"
)








func init() {
	initHijack()

	
	SetFormatter(NewDefaultFormatter(os.Stderr))
	SetGlobalLogLevel(INFO)
}

func NewDefaultFormatter(out io.Writer) Formatter {
	if syscall.Getppid() == 1 {
		
		f, err := NewJournaldFormatter()
		if err == nil {
			return f
		}
	}
	return NewPrettyFormatter(out, false)
}
