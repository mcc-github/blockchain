





package http2

import "net/textproto"

func traceHasWroteHeaderField(trace *clientTrace) bool { return false }

func traceWroteHeaderField(trace *clientTrace, k, v string) {}

func traceGot1xxResponseFunc(trace *clientTrace) func(int, textproto.MIMEHeader) error {
	return nil
}
