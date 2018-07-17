




package trace 

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/internal/timeseries"
)



var DebugUseAfterFinish = false











var AuthRequest = func(req *http.Request) (any, sensitive bool) {
	
	
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}
	switch host {
	case "localhost", "127.0.0.1", "::1":
		return true, true
	default:
		return false, false
	}
}

func init() {
	
	
	http.HandleFunc("/debug/requests", Traces)
	http.HandleFunc("/debug/events", Events)
}






func Traces(w http.ResponseWriter, req *http.Request) {
	any, sensitive := AuthRequest(req)
	if !any {
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	Render(w, req, sensitive)
}






func Events(w http.ResponseWriter, req *http.Request) {
	any, sensitive := AuthRequest(req)
	if !any {
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	RenderEvents(w, req, sensitive)
}





func Render(w io.Writer, req *http.Request, sensitive bool) {
	data := &struct {
		Families         []string
		ActiveTraceCount map[string]int
		CompletedTraces  map[string]*family

		
		Traces        traceList
		Family        string
		Bucket        int
		Expanded      bool
		Traced        bool
		Active        bool
		ShowSensitive bool 

		Histogram       template.HTML
		HistogramWindow string 

		
		
		Total int
	}{
		CompletedTraces: completedTraces,
	}

	data.ShowSensitive = sensitive
	if req != nil {
		
		
		if req.FormValue("show_sensitive") == "0" {
			data.ShowSensitive = false
		}

		if exp, err := strconv.ParseBool(req.FormValue("exp")); err == nil {
			data.Expanded = exp
		}
		if exp, err := strconv.ParseBool(req.FormValue("rtraced")); err == nil {
			data.Traced = exp
		}
	}

	completedMu.RLock()
	data.Families = make([]string, 0, len(completedTraces))
	for fam := range completedTraces {
		data.Families = append(data.Families, fam)
	}
	completedMu.RUnlock()
	sort.Strings(data.Families)

	
	
	data.ActiveTraceCount = make(map[string]int, len(data.Families))
	activeMu.RLock()
	for fam, s := range activeTraces {
		data.ActiveTraceCount[fam] = s.Len()
	}
	activeMu.RUnlock()

	var ok bool
	data.Family, data.Bucket, ok = parseArgs(req)
	switch {
	case !ok:
		
	case data.Bucket == -1:
		data.Active = true
		n := data.ActiveTraceCount[data.Family]
		data.Traces = getActiveTraces(data.Family)
		if len(data.Traces) < n {
			data.Total = n
		}
	case data.Bucket < bucketsPerFamily:
		if b := lookupBucket(data.Family, data.Bucket); b != nil {
			data.Traces = b.Copy(data.Traced)
		}
	default:
		if f := getFamily(data.Family, false); f != nil {
			var obs timeseries.Observable
			f.LatencyMu.RLock()
			switch o := data.Bucket - bucketsPerFamily; o {
			case 0:
				obs = f.Latency.Minute()
				data.HistogramWindow = "last minute"
			case 1:
				obs = f.Latency.Hour()
				data.HistogramWindow = "last hour"
			case 2:
				obs = f.Latency.Total()
				data.HistogramWindow = "all time"
			}
			f.LatencyMu.RUnlock()
			if obs != nil {
				data.Histogram = obs.(*histogram).html()
			}
		}
	}

	if data.Traces != nil {
		defer data.Traces.Free()
		sort.Sort(data.Traces)
	}

	completedMu.RLock()
	defer completedMu.RUnlock()
	if err := pageTmpl().ExecuteTemplate(w, "Page", data); err != nil {
		log.Printf("net/trace: Failed executing template: %v", err)
	}
}

func parseArgs(req *http.Request) (fam string, b int, ok bool) {
	if req == nil {
		return "", 0, false
	}
	fam, bStr := req.FormValue("fam"), req.FormValue("b")
	if fam == "" || bStr == "" {
		return "", 0, false
	}
	b, err := strconv.Atoi(bStr)
	if err != nil || b < -1 {
		return "", 0, false
	}

	return fam, b, true
}

func lookupBucket(fam string, b int) *traceBucket {
	f := getFamily(fam, false)
	if f == nil || b < 0 || b >= len(f.Buckets) {
		return nil
	}
	return f.Buckets[b]
}

type contextKeyT string

var contextKey = contextKeyT("golang.org/x/net/trace.Trace")


type Trace interface {
	
	
	
	LazyLog(x fmt.Stringer, sensitive bool)

	
	
	
	LazyPrintf(format string, a ...interface{})

	
	SetError()

	
	
	
	
	SetRecycler(f func(interface{}))

	
	
	SetTraceInfo(traceID, spanID uint64)

	
	
	
	SetMaxEvents(m int)

	
	
	Finish()
}

type lazySprintf struct {
	format string
	a      []interface{}
}

func (l *lazySprintf) String() string {
	return fmt.Sprintf(l.format, l.a...)
}


func New(family, title string) Trace {
	tr := newTrace()
	tr.ref()
	tr.Family, tr.Title = family, title
	tr.Start = time.Now()
	tr.maxEvents = maxEventsPerTrace
	tr.events = tr.eventsBuf[:0]

	activeMu.RLock()
	s := activeTraces[tr.Family]
	activeMu.RUnlock()
	if s == nil {
		activeMu.Lock()
		s = activeTraces[tr.Family] 
		if s == nil {
			s = new(traceSet)
			activeTraces[tr.Family] = s
		}
		activeMu.Unlock()
	}
	s.Add(tr)

	
	
	
	
	
	completedMu.RLock()
	if _, ok := completedTraces[tr.Family]; !ok {
		go allocFamily(tr.Family)
	}
	completedMu.RUnlock()

	return tr
}

func (tr *trace) Finish() {
	elapsed := time.Now().Sub(tr.Start)
	tr.mu.Lock()
	tr.Elapsed = elapsed
	tr.mu.Unlock()

	if DebugUseAfterFinish {
		buf := make([]byte, 4<<10) 
		n := runtime.Stack(buf, false)
		tr.finishStack = buf[:n]
	}

	activeMu.RLock()
	m := activeTraces[tr.Family]
	activeMu.RUnlock()
	m.Remove(tr)

	f := getFamily(tr.Family, true)
	tr.mu.RLock() 
	for _, b := range f.Buckets {
		if b.Cond.match(tr) {
			b.Add(tr)
		}
	}
	tr.mu.RUnlock()

	
	h := new(histogram)
	h.addMeasurement(elapsed.Nanoseconds() / 1e3)
	f.LatencyMu.Lock()
	f.Latency.Add(h)
	f.LatencyMu.Unlock()

	tr.unref() 
}

const (
	bucketsPerFamily    = 9
	tracesPerBucket     = 10
	maxActiveTraces     = 20 
	maxEventsPerTrace   = 10
	numHistogramBuckets = 38
)

var (
	
	activeMu     sync.RWMutex
	activeTraces = make(map[string]*traceSet) 

	
	completedMu     sync.RWMutex
	completedTraces = make(map[string]*family) 
)

type traceSet struct {
	mu sync.RWMutex
	m  map[*trace]bool

	
	
	
	
	
}

func (ts *traceSet) Len() int {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return len(ts.m)
}

func (ts *traceSet) Add(tr *trace) {
	ts.mu.Lock()
	if ts.m == nil {
		ts.m = make(map[*trace]bool)
	}
	ts.m[tr] = true
	ts.mu.Unlock()
}

func (ts *traceSet) Remove(tr *trace) {
	ts.mu.Lock()
	delete(ts.m, tr)
	ts.mu.Unlock()
}


func (ts *traceSet) FirstN(n int) traceList {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	if n > len(ts.m) {
		n = len(ts.m)
	}
	trl := make(traceList, 0, n)

	
	if n == len(ts.m) {
		for tr := range ts.m {
			tr.ref()
			trl = append(trl, tr)
		}
		sort.Sort(trl)
		return trl
	}

	
	
	for tr := range ts.m {
		
		
		if len(trl) < n {
			tr.ref()
			trl = append(trl, tr)
			if len(trl) == n {
				
				sort.Sort(trl)
			}
			continue
		}
		if tr.Start.After(trl[n-1].Start) {
			continue
		}

		
		tr.ref()
		i := sort.Search(n, func(i int) bool { return trl[i].Start.After(tr.Start) })
		trl[n-1].unref()
		copy(trl[i+1:], trl[i:])
		trl[i] = tr
	}

	return trl
}

func getActiveTraces(fam string) traceList {
	activeMu.RLock()
	s := activeTraces[fam]
	activeMu.RUnlock()
	if s == nil {
		return nil
	}
	return s.FirstN(maxActiveTraces)
}

func getFamily(fam string, allocNew bool) *family {
	completedMu.RLock()
	f := completedTraces[fam]
	completedMu.RUnlock()
	if f == nil && allocNew {
		f = allocFamily(fam)
	}
	return f
}

func allocFamily(fam string) *family {
	completedMu.Lock()
	defer completedMu.Unlock()
	f := completedTraces[fam]
	if f == nil {
		f = newFamily()
		completedTraces[fam] = f
	}
	return f
}


type family struct {
	
	Buckets [bucketsPerFamily]*traceBucket

	
	LatencyMu sync.RWMutex
	Latency   *timeseries.MinuteHourSeries
}

func newFamily() *family {
	return &family{
		Buckets: [bucketsPerFamily]*traceBucket{
			{Cond: minCond(0)},
			{Cond: minCond(50 * time.Millisecond)},
			{Cond: minCond(100 * time.Millisecond)},
			{Cond: minCond(200 * time.Millisecond)},
			{Cond: minCond(500 * time.Millisecond)},
			{Cond: minCond(1 * time.Second)},
			{Cond: minCond(10 * time.Second)},
			{Cond: minCond(100 * time.Second)},
			{Cond: errorCond{}},
		},
		Latency: timeseries.NewMinuteHourSeries(func() timeseries.Observable { return new(histogram) }),
	}
}



type traceBucket struct {
	Cond cond

	
	mu     sync.RWMutex
	buf    [tracesPerBucket]*trace
	start  int 
	length int 
}

func (b *traceBucket) Add(tr *trace) {
	b.mu.Lock()
	defer b.mu.Unlock()

	i := b.start + b.length
	if i >= tracesPerBucket {
		i -= tracesPerBucket
	}
	if b.length == tracesPerBucket {
		
		b.buf[i].unref()
		b.start++
		if b.start == tracesPerBucket {
			b.start = 0
		}
	}
	b.buf[i] = tr
	if b.length < tracesPerBucket {
		b.length++
	}
	tr.ref()
}






func (b *traceBucket) Copy(tracedOnly bool) traceList {
	b.mu.RLock()
	defer b.mu.RUnlock()

	trl := make(traceList, 0, b.length)
	for i, x := 0, b.start; i < b.length; i++ {
		tr := b.buf[x]
		if !tracedOnly || tr.spanID != 0 {
			tr.ref()
			trl = append(trl, tr)
		}
		x++
		if x == b.length {
			x = 0
		}
	}
	return trl
}

func (b *traceBucket) Empty() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.length == 0
}


type cond interface {
	match(t *trace) bool
	String() string
}

type minCond time.Duration

func (m minCond) match(t *trace) bool { return t.Elapsed >= time.Duration(m) }
func (m minCond) String() string      { return fmt.Sprintf("≥%gs", time.Duration(m).Seconds()) }

type errorCond struct{}

func (e errorCond) match(t *trace) bool { return t.IsError }
func (e errorCond) String() string      { return "errors" }

type traceList []*trace


func (trl traceList) Free() {
	for _, t := range trl {
		t.unref()
	}
}


func (trl traceList) Len() int           { return len(trl) }
func (trl traceList) Less(i, j int) bool { return trl[i].Start.After(trl[j].Start) }
func (trl traceList) Swap(i, j int)      { trl[i], trl[j] = trl[j], trl[i] }


type event struct {
	When       time.Time
	Elapsed    time.Duration 
	NewDay     bool          
	Recyclable bool          
	Sensitive  bool          
	What       interface{}   
}



func (e event) WhenString() string {
	if e.NewDay {
		return e.When.Format("2006/01/02 15:04:05.000000")
	}
	return e.When.Format("15:04:05.000000")
}



type discarded int

func (d *discarded) String() string {
	return fmt.Sprintf("(%d events discarded)", int(*d))
}



type trace struct {
	
	Family string

	
	Title string

	
	Start time.Time

	mu        sync.RWMutex
	events    []event 
	maxEvents int
	recycler  func(interface{})
	IsError   bool          
	Elapsed   time.Duration 
	traceID   uint64        
	spanID    uint64

	refs int32     
	disc discarded 

	finishStack []byte 

	eventsBuf [4]event 
}

func (tr *trace) reset() {
	
	tr.Family = ""
	tr.Title = ""
	tr.Start = time.Time{}

	tr.mu.Lock()
	tr.Elapsed = 0
	tr.traceID = 0
	tr.spanID = 0
	tr.IsError = false
	tr.maxEvents = 0
	tr.events = nil
	tr.recycler = nil
	tr.mu.Unlock()

	tr.refs = 0
	tr.disc = 0
	tr.finishStack = nil
	for i := range tr.eventsBuf {
		tr.eventsBuf[i] = event{}
	}
}




func (tr *trace) delta(t time.Time) (time.Duration, bool) {
	if len(tr.events) == 0 {
		return t.Sub(tr.Start), false
	}
	prev := tr.events[len(tr.events)-1].When
	return t.Sub(prev), prev.Day() != t.Day()
}

func (tr *trace) addEvent(x interface{}, recyclable, sensitive bool) {
	if DebugUseAfterFinish && tr.finishStack != nil {
		buf := make([]byte, 4<<10) 
		n := runtime.Stack(buf, false)
		log.Printf("net/trace: trace used after finish:\nFinished at:\n%s\nUsed at:\n%s", tr.finishStack, buf[:n])
	}

	

	e := event{When: time.Now(), What: x, Recyclable: recyclable, Sensitive: sensitive}
	tr.mu.Lock()
	e.Elapsed, e.NewDay = tr.delta(e.When)
	if len(tr.events) < tr.maxEvents {
		tr.events = append(tr.events, e)
	} else {
		
		di := int((tr.maxEvents - 1) / 2)
		if d, ok := tr.events[di].What.(*discarded); ok {
			(*d)++
		} else {
			
			
			tr.disc = 2
			if tr.recycler != nil && tr.events[di].Recyclable {
				go tr.recycler(tr.events[di].What)
			}
			tr.events[di].What = &tr.disc
		}
		
		
		tr.events[di].When = tr.events[di+1].When

		if tr.recycler != nil && tr.events[di+1].Recyclable {
			go tr.recycler(tr.events[di+1].What)
		}
		copy(tr.events[di+1:], tr.events[di+2:])
		tr.events[tr.maxEvents-1] = e
	}
	tr.mu.Unlock()
}

func (tr *trace) LazyLog(x fmt.Stringer, sensitive bool) {
	tr.addEvent(x, true, sensitive)
}

func (tr *trace) LazyPrintf(format string, a ...interface{}) {
	tr.addEvent(&lazySprintf{format, a}, false, false)
}

func (tr *trace) SetError() {
	tr.mu.Lock()
	tr.IsError = true
	tr.mu.Unlock()
}

func (tr *trace) SetRecycler(f func(interface{})) {
	tr.mu.Lock()
	tr.recycler = f
	tr.mu.Unlock()
}

func (tr *trace) SetTraceInfo(traceID, spanID uint64) {
	tr.mu.Lock()
	tr.traceID, tr.spanID = traceID, spanID
	tr.mu.Unlock()
}

func (tr *trace) SetMaxEvents(m int) {
	tr.mu.Lock()
	
	if len(tr.events) == 0 && m > 3 {
		tr.maxEvents = m
	}
	tr.mu.Unlock()
}

func (tr *trace) ref() {
	atomic.AddInt32(&tr.refs, 1)
}

func (tr *trace) unref() {
	if atomic.AddInt32(&tr.refs, -1) == 0 {
		tr.mu.RLock()
		if tr.recycler != nil {
			
			go func(f func(interface{}), es []event) {
				for _, e := range es {
					if e.Recyclable {
						f(e.What)
					}
				}
			}(tr.recycler, tr.events)
		}
		tr.mu.RUnlock()

		freeTrace(tr)
	}
}

func (tr *trace) When() string {
	return tr.Start.Format("2006/01/02 15:04:05.000000")
}

func (tr *trace) ElapsedTime() string {
	tr.mu.RLock()
	t := tr.Elapsed
	tr.mu.RUnlock()

	if t == 0 {
		
		t = time.Since(tr.Start)
	}
	return fmt.Sprintf("%.6f", t.Seconds())
}

func (tr *trace) Events() []event {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return tr.events
}

var traceFreeList = make(chan *trace, 1000) 


func newTrace() *trace {
	select {
	case tr := <-traceFreeList:
		return tr
	default:
		return new(trace)
	}
}



func freeTrace(tr *trace) {
	if DebugUseAfterFinish {
		return 
	}
	tr.reset()
	select {
	case traceFreeList <- tr:
	default:
	}
}

func elapsed(d time.Duration) string {
	b := []byte(fmt.Sprintf("%.6f", d.Seconds()))

	
	
	if d < time.Second {
		dot := bytes.IndexByte(b, '.')
		for i := 0; i < dot; i++ {
			b[i] = ' '
		}
		for i := dot + 1; i < len(b); i++ {
			if b[i] == '0' {
				b[i] = ' '
			} else {
				break
			}
		}
	}

	return string(b)
}

var pageTmplCache *template.Template
var pageTmplOnce sync.Once

func pageTmpl() *template.Template {
	pageTmplOnce.Do(func() {
		pageTmplCache = template.Must(template.New("Page").Funcs(template.FuncMap{
			"elapsed": elapsed,
			"add":     func(a, b int) int { return a + b },
		}).Parse(pageHTML))
	})
	return pageTmplCache
}

const pageHTML = `
{{template "Prolog" .}}
{{template "StatusTable" .}}
{{template "Epilog" .}}

{{define "Prolog"}}
<html>
	<head>
	<title>/debug/requests</title>
	<style type="text/css">
		body {
			font-family: sans-serif;
		}
		table#tr-status td.family {
			padding-right: 2em;
		}
		table#tr-status td.active {
			padding-right: 1em;
		}
		table#tr-status td.latency-first {
			padding-left: 1em;
		}
		table#tr-status td.empty {
			color: #aaa;
		}
		table#reqs {
			margin-top: 1em;
		}
		table#reqs tr.first {
			{{if $.Expanded}}font-weight: bold;{{end}}
		}
		table#reqs td {
			font-family: monospace;
		}
		table#reqs td.when {
			text-align: right;
			white-space: nowrap;
		}
		table#reqs td.elapsed {
			padding: 0 0.5em;
			text-align: right;
			white-space: pre;
			width: 10em;
		}
		address {
			font-size: smaller;
			margin-top: 5em;
		}
	</style>
	</head>
	<body>

<h1>/debug/requests</h1>
{{end}} {{}}

{{define "StatusTable"}}
<table id="tr-status">
	{{range $fam := .Families}}
	<tr>
		<td class="family">{{$fam}}</td>

		{{$n := index $.ActiveTraceCount $fam}}
		<td class="active {{if not $n}}empty{{end}}">
			{{if $n}}<a href="?fam={{$fam}}&b=-1{{if $.Expanded}}&exp=1{{end}}">{{end}}
			[{{$n}} active]
			{{if $n}}</a>{{end}}
		</td>

		{{$f := index $.CompletedTraces $fam}}
		{{range $i, $b := $f.Buckets}}
		{{$empty := $b.Empty}}
		<td {{if $empty}}class="empty"{{end}}>
		{{if not $empty}}<a href="?fam={{$fam}}&b={{$i}}{{if $.Expanded}}&exp=1{{end}}">{{end}}
		[{{.Cond}}]
		{{if not $empty}}</a>{{end}}
		</td>
		{{end}}

		{{$nb := len $f.Buckets}}
		<td class="latency-first">
		<a href="?fam={{$fam}}&b={{$nb}}">[minute]</a>
		</td>
		<td>
		<a href="?fam={{$fam}}&b={{add $nb 1}}">[hour]</a>
		</td>
		<td>
		<a href="?fam={{$fam}}&b={{add $nb 2}}">[total]</a>
		</td>

	</tr>
	{{end}}
</table>
{{end}} {{}}

{{define "Epilog"}}
{{if $.Traces}}
<hr />
<h3>Family: {{$.Family}}</h3>

{{if or $.Expanded $.Traced}}
  <a href="?fam={{$.Family}}&b={{$.Bucket}}">[Normal/Summary]</a>
{{else}}
  [Normal/Summary]
{{end}}

{{if or (not $.Expanded) $.Traced}}
  <a href="?fam={{$.Family}}&b={{$.Bucket}}&exp=1">[Normal/Expanded]</a>
{{else}}
  [Normal/Expanded]
{{end}}

{{if not $.Active}}
	{{if or $.Expanded (not $.Traced)}}
	<a href="?fam={{$.Family}}&b={{$.Bucket}}&rtraced=1">[Traced/Summary]</a>
	{{else}}
	[Traced/Summary]
	{{end}}
	{{if or (not $.Expanded) (not $.Traced)}}
	<a href="?fam={{$.Family}}&b={{$.Bucket}}&exp=1&rtraced=1">[Traced/Expanded]</a>
        {{else}}
	[Traced/Expanded]
	{{end}}
{{end}}

{{if $.Total}}
<p><em>Showing <b>{{len $.Traces}}</b> of <b>{{$.Total}}</b> traces.</em></p>
{{end}}

<table id="reqs">
	<caption>
		{{if $.Active}}Active{{else}}Completed{{end}} Requests
	</caption>
	<tr><th>When</th><th>Elapsed&nbsp;(s)</th></tr>
	{{range $tr := $.Traces}}
	<tr class="first">
		<td class="when">{{$tr.When}}</td>
		<td class="elapsed">{{$tr.ElapsedTime}}</td>
		<td>{{$tr.Title}}</td>
		{{}}
	</tr>
	{{if $.Expanded}}
	{{range $tr.Events}}
	<tr>
		<td class="when">{{.WhenString}}</td>
		<td class="elapsed">{{elapsed .Elapsed}}</td>
		<td>{{if or $.ShowSensitive (not .Sensitive)}}... {{.What}}{{else}}<em>[redacted]</em>{{end}}</td>
	</tr>
	{{end}}
	{{end}}
	{{end}}
</table>
{{end}} {{}}

{{if $.Histogram}}
<h4>Latency (&micro;s) of {{$.Family}} over {{$.HistogramWindow}}</h4>
{{$.Histogram}}
{{end}} {{}}

	</body>
</html>
{{end}} {{}}
`
