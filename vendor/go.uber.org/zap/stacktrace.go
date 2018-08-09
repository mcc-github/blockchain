



















package zap

import (
	"runtime"
	"strings"
	"sync"

	"go.uber.org/zap/internal/bufferpool"
)

const _zapPackage = "go.uber.org/zap"

var (
	_stacktracePool = sync.Pool{
		New: func() interface{} {
			return newProgramCounters(64)
		},
	}

	
	
	_zapStacktracePrefixes       = addPrefix(_zapPackage, ".", "/")
	_zapStacktraceVendorContains = addPrefix("/vendor/", _zapStacktracePrefixes...)
)

func takeStacktrace() string {
	buffer := bufferpool.Get()
	defer buffer.Free()
	programCounters := _stacktracePool.Get().(*programCounters)
	defer _stacktracePool.Put(programCounters)

	var numFrames int
	for {
		
		
		numFrames = runtime.Callers(2, programCounters.pcs)
		if numFrames < len(programCounters.pcs) {
			break
		}
		
		
		programCounters = newProgramCounters(len(programCounters.pcs) * 2)
	}

	i := 0
	skipZapFrames := true 
	frames := runtime.CallersFrames(programCounters.pcs[:numFrames])

	
	
	
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if skipZapFrames && isZapFrame(frame.Function) {
			continue
		} else {
			skipZapFrames = false
		}

		if i != 0 {
			buffer.AppendByte('\n')
		}
		i++
		buffer.AppendString(frame.Function)
		buffer.AppendByte('\n')
		buffer.AppendByte('\t')
		buffer.AppendString(frame.File)
		buffer.AppendByte(':')
		buffer.AppendInt(int64(frame.Line))
	}

	return buffer.String()
}

func isZapFrame(function string) bool {
	for _, prefix := range _zapStacktracePrefixes {
		if strings.HasPrefix(function, prefix) {
			return true
		}
	}

	
	
	for _, contains := range _zapStacktraceVendorContains {
		if strings.Contains(function, contains) {
			return true
		}
	}

	return false
}

type programCounters struct {
	pcs []uintptr
}

func newProgramCounters(size int) *programCounters {
	return &programCounters{make([]uintptr, size)}
}

func addPrefix(prefix string, ss ...string) []string {
	withPrefix := make([]string, len(ss))
	for i, s := range ss {
		withPrefix[i] = prefix + s
	}
	return withPrefix
}
