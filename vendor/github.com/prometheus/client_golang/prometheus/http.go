












package prometheus

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/common/expfmt"
)






const (
	contentTypeHeader     = "Content-Type"
	contentLengthHeader   = "Content-Length"
	contentEncodingHeader = "Content-Encoding"
	acceptEncodingHeader  = "Accept-Encoding"
)

var bufPool sync.Pool

func getBuf() *bytes.Buffer {
	buf := bufPool.Get()
	if buf == nil {
		return &bytes.Buffer{}
	}
	return buf.(*bytes.Buffer)
}

func giveBuf(buf *bytes.Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}







func Handler() http.Handler {
	return InstrumentHandler("prometheus", UninstrumentedHandler())
}





func UninstrumentedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mfs, err := DefaultGatherer.Gather()
		if err != nil {
			http.Error(w, "An error has occurred during metrics collection:\n\n"+err.Error(), http.StatusInternalServerError)
			return
		}

		contentType := expfmt.Negotiate(req.Header)
		buf := getBuf()
		defer giveBuf(buf)
		writer, encoding := decorateWriter(req, buf)
		enc := expfmt.NewEncoder(writer, contentType)
		var lastErr error
		for _, mf := range mfs {
			if err := enc.Encode(mf); err != nil {
				lastErr = err
				http.Error(w, "An error has occurred during metrics encoding:\n\n"+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if closer, ok := writer.(io.Closer); ok {
			closer.Close()
		}
		if lastErr != nil && buf.Len() == 0 {
			http.Error(w, "No metrics encoded, last error:\n\n"+lastErr.Error(), http.StatusInternalServerError)
			return
		}
		header := w.Header()
		header.Set(contentTypeHeader, string(contentType))
		header.Set(contentLengthHeader, fmt.Sprint(buf.Len()))
		if encoding != "" {
			header.Set(contentEncodingHeader, encoding)
		}
		w.Write(buf.Bytes())
	})
}




func decorateWriter(request *http.Request, writer io.Writer) (io.Writer, string) {
	header := request.Header.Get(acceptEncodingHeader)
	parts := strings.Split(header, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "gzip" || strings.HasPrefix(part, "gzip;") {
			return gzip.NewWriter(writer), "gzip"
		}
	}
	return writer, ""
}

var instLabels = []string{"method", "code"}

type nower interface {
	Now() time.Time
}

type nowFunc func() time.Time

func (n nowFunc) Now() time.Time {
	return n()
}

var now nower = nowFunc(func() time.Time {
	return time.Now()
})




















func InstrumentHandler(handlerName string, handler http.Handler) http.HandlerFunc {
	return InstrumentHandlerFunc(handlerName, handler.ServeHTTP)
}







func InstrumentHandlerFunc(handlerName string, handlerFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return InstrumentHandlerFuncWithOpts(
		SummaryOpts{
			Subsystem:   "http",
			ConstLabels: Labels{"handler": handlerName},
			Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		handlerFunc,
	)
}





























func InstrumentHandlerWithOpts(opts SummaryOpts, handler http.Handler) http.HandlerFunc {
	return InstrumentHandlerFuncWithOpts(opts, handler.ServeHTTP)
}








func InstrumentHandlerFuncWithOpts(opts SummaryOpts, handlerFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	reqCnt := NewCounterVec(
		CounterOpts{
			Namespace:   opts.Namespace,
			Subsystem:   opts.Subsystem,
			Name:        "requests_total",
			Help:        "Total number of HTTP requests made.",
			ConstLabels: opts.ConstLabels,
		},
		instLabels,
	)
	if err := Register(reqCnt); err != nil {
		if are, ok := err.(AlreadyRegisteredError); ok {
			reqCnt = are.ExistingCollector.(*CounterVec)
		} else {
			panic(err)
		}
	}

	opts.Name = "request_duration_microseconds"
	opts.Help = "The HTTP request latencies in microseconds."
	reqDur := NewSummary(opts)
	if err := Register(reqDur); err != nil {
		if are, ok := err.(AlreadyRegisteredError); ok {
			reqDur = are.ExistingCollector.(Summary)
		} else {
			panic(err)
		}
	}

	opts.Name = "request_size_bytes"
	opts.Help = "The HTTP request sizes in bytes."
	reqSz := NewSummary(opts)
	if err := Register(reqSz); err != nil {
		if are, ok := err.(AlreadyRegisteredError); ok {
			reqSz = are.ExistingCollector.(Summary)
		} else {
			panic(err)
		}
	}

	opts.Name = "response_size_bytes"
	opts.Help = "The HTTP response sizes in bytes."
	resSz := NewSummary(opts)
	if err := Register(resSz); err != nil {
		if are, ok := err.(AlreadyRegisteredError); ok {
			resSz = are.ExistingCollector.(Summary)
		} else {
			panic(err)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		delegate := &responseWriterDelegator{ResponseWriter: w}
		out := computeApproximateRequestSize(r)

		_, cn := w.(http.CloseNotifier)
		_, fl := w.(http.Flusher)
		_, hj := w.(http.Hijacker)
		_, rf := w.(io.ReaderFrom)
		var rw http.ResponseWriter
		if cn && fl && hj && rf {
			rw = &fancyResponseWriterDelegator{delegate}
		} else {
			rw = delegate
		}
		handlerFunc(rw, r)

		elapsed := float64(time.Since(now)) / float64(time.Microsecond)

		method := sanitizeMethod(r.Method)
		code := sanitizeCode(delegate.status)
		reqCnt.WithLabelValues(method, code).Inc()
		reqDur.Observe(elapsed)
		resSz.Observe(float64(delegate.written))
		reqSz.Observe(float64(<-out))
	})
}

func computeApproximateRequestSize(r *http.Request) <-chan int {
	
	
	s := 0
	if r.URL != nil {
		s += len(r.URL.String())
	}

	out := make(chan int, 1)

	go func() {
		s += len(r.Method)
		s += len(r.Proto)
		for name, values := range r.Header {
			s += len(name)
			for _, value := range values {
				s += len(value)
			}
		}
		s += len(r.Host)

		

		if r.ContentLength != -1 {
			s += int(r.ContentLength)
		}
		out <- s
		close(out)
	}()

	return out
}

type responseWriterDelegator struct {
	http.ResponseWriter

	status      int
	written     int64
	wroteHeader bool
}

func (r *responseWriterDelegator) WriteHeader(code int) {
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}

type fancyResponseWriterDelegator struct {
	*responseWriterDelegator
}

func (f *fancyResponseWriterDelegator) CloseNotify() <-chan bool {
	return f.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (f *fancyResponseWriterDelegator) Flush() {
	f.ResponseWriter.(http.Flusher).Flush()
}

func (f *fancyResponseWriterDelegator) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return f.ResponseWriter.(http.Hijacker).Hijack()
}

func (f *fancyResponseWriterDelegator) ReadFrom(r io.Reader) (int64, error) {
	if !f.wroteHeader {
		f.WriteHeader(http.StatusOK)
	}
	n, err := f.ResponseWriter.(io.ReaderFrom).ReadFrom(r)
	f.written += n
	return n, err
}

func sanitizeMethod(m string) string {
	switch m {
	case "GET", "get":
		return "get"
	case "PUT", "put":
		return "put"
	case "HEAD", "head":
		return "head"
	case "POST", "post":
		return "post"
	case "DELETE", "delete":
		return "delete"
	case "CONNECT", "connect":
		return "connect"
	case "OPTIONS", "options":
		return "options"
	case "NOTIFY", "notify":
		return "notify"
	default:
		return strings.ToLower(m)
	}
}

func sanitizeCode(s int) string {
	switch s {
	case 100:
		return "100"
	case 101:
		return "101"

	case 200:
		return "200"
	case 201:
		return "201"
	case 202:
		return "202"
	case 203:
		return "203"
	case 204:
		return "204"
	case 205:
		return "205"
	case 206:
		return "206"

	case 300:
		return "300"
	case 301:
		return "301"
	case 302:
		return "302"
	case 304:
		return "304"
	case 305:
		return "305"
	case 307:
		return "307"

	case 400:
		return "400"
	case 401:
		return "401"
	case 402:
		return "402"
	case 403:
		return "403"
	case 404:
		return "404"
	case 405:
		return "405"
	case 406:
		return "406"
	case 407:
		return "407"
	case 408:
		return "408"
	case 409:
		return "409"
	case 410:
		return "410"
	case 411:
		return "411"
	case 412:
		return "412"
	case 413:
		return "413"
	case 414:
		return "414"
	case 415:
		return "415"
	case 416:
		return "416"
	case 417:
		return "417"
	case 418:
		return "418"

	case 500:
		return "500"
	case 501:
		return "501"
	case 502:
		return "502"
	case 503:
		return "503"
	case 504:
		return "504"
	case 505:
		return "505"

	case 428:
		return "428"
	case 429:
		return "429"
	case 431:
		return "431"
	case 511:
		return "511"

	default:
		return strconv.Itoa(s)
	}
}
