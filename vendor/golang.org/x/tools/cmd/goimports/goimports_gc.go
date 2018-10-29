





package main

import (
	"flag"
	"runtime/trace"
)

var traceProfile = flag.String("trace", "", "trace profile output")

func doTrace() func() {
	if *traceProfile != "" {
		bw, flush := bufferedFileWriter(*traceProfile)
		trace.Start(bw)
		return func() {
			flush()
			trace.Stop()
		}
	}
	return func() {}
}
