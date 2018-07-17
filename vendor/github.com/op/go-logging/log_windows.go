




package logging

import (
	"bytes"
	"io"
	"log"
	"syscall"
)

var (
	kernel32DLL                 = syscall.NewLazyDLL("kernel32.dll")
	setConsoleTextAttributeProc = kernel32DLL.NewProc("SetConsoleTextAttribute")
)






const (
	fgBlack     = 0x0000
	fgBlue      = 0x0001
	fgGreen     = 0x0002
	fgCyan      = 0x0003
	fgRed       = 0x0004
	fgMagenta   = 0x0005
	fgYellow    = 0x0006
	fgWhite     = 0x0007
	fgIntensity = 0x0008
	fgMask      = 0x000F
)

var (
	colors = []uint16{
		INFO:     fgWhite,
		CRITICAL: fgMagenta,
		ERROR:    fgRed,
		WARNING:  fgYellow,
		NOTICE:   fgGreen,
		DEBUG:    fgCyan,
	}
	boldcolors = []uint16{
		INFO:     fgWhite | fgIntensity,
		CRITICAL: fgMagenta | fgIntensity,
		ERROR:    fgRed | fgIntensity,
		WARNING:  fgYellow | fgIntensity,
		NOTICE:   fgGreen | fgIntensity,
		DEBUG:    fgCyan | fgIntensity,
	}
)

type file interface {
	Fd() uintptr
}


type LogBackend struct {
	Logger *log.Logger
	Color  bool

	
	
	f file
}


func NewLogBackend(out io.Writer, prefix string, flag int) *LogBackend {
	b := &LogBackend{Logger: log.New(out, prefix, flag)}

	
	
	if f, ok := out.(file); ok {
		b.f = f
	}

	return b
}

func (b *LogBackend) Log(level Level, calldepth int, rec *Record) error {
	if b.Color && b.f != nil {
		buf := &bytes.Buffer{}
		setConsoleTextAttribute(b.f, colors[level])
		buf.Write([]byte(rec.Formatted(calldepth + 1)))
		err := b.Logger.Output(calldepth+2, buf.String())
		setConsoleTextAttribute(b.f, fgWhite)
		return err
	}
	return b.Logger.Output(calldepth+2, rec.Formatted(calldepth+1))
}




func setConsoleTextAttribute(f file, attribute uint16) bool {
	ok, _, _ := setConsoleTextAttributeProc.Call(f.Fd(), uintptr(attribute), 0)
	return ok != 0
}

func doFmtVerbLevelColor(layout string, level Level, output io.Writer) {
	
	
}
