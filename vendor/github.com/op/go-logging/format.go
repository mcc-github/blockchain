



package logging

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)





type fmtVerb int

const (
	fmtVerbTime fmtVerb = iota
	fmtVerbLevel
	fmtVerbID
	fmtVerbPid
	fmtVerbProgram
	fmtVerbModule
	fmtVerbMessage
	fmtVerbLongfile
	fmtVerbShortfile
	fmtVerbLongpkg
	fmtVerbShortpkg
	fmtVerbLongfunc
	fmtVerbShortfunc
	fmtVerbCallpath
	fmtVerbLevelColor

	
	fmtVerbUnknown
	fmtVerbStatic
)

var fmtVerbs = []string{
	"time",
	"level",
	"id",
	"pid",
	"program",
	"module",
	"message",
	"longfile",
	"shortfile",
	"longpkg",
	"shortpkg",
	"longfunc",
	"shortfunc",
	"callpath",
	"color",
}

const rfc3339Milli = "2006-01-02T15:04:05.999Z07:00"

var defaultVerbsLayout = []string{
	rfc3339Milli,
	"s",
	"d",
	"d",
	"s",
	"s",
	"s",
	"s",
	"s",
	"s",
	"s",
	"s",
	"s",
	"0",
	"",
}

var (
	pid     = os.Getpid()
	program = filepath.Base(os.Args[0])
)

func getFmtVerbByName(name string) fmtVerb {
	for i, verb := range fmtVerbs {
		if name == verb {
			return fmtVerb(i)
		}
	}
	return fmtVerbUnknown
}


type Formatter interface {
	Format(calldepth int, r *Record, w io.Writer) error
}


var formatter struct {
	sync.RWMutex
	def Formatter
}

func getFormatter() Formatter {
	formatter.RLock()
	defer formatter.RUnlock()
	return formatter.def
}

var (
	
	DefaultFormatter = MustStringFormatter("%{message}")

	
	GlogFormatter = MustStringFormatter("%{level:.1s}%{time:0102 15:04:05.999999} %{pid} %{shortfile}] %{message}")
)





func SetFormatter(f Formatter) {
	formatter.Lock()
	defer formatter.Unlock()
	formatter.def = f
}

var formatRe = regexp.MustCompile(`%{([a-z]+)(?::(.*?[^\\]))?}`)

type part struct {
	verb   fmtVerb
	layout string
}



type stringFormatter struct {
	parts []part
}

















































func NewStringFormatter(format string) (Formatter, error) {
	var fmter = &stringFormatter{}

	
	matches := formatRe.FindAllStringSubmatchIndex(format, -1)
	if matches == nil {
		return nil, errors.New("logger: invalid log format: " + format)
	}

	
	prev := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		if start > prev {
			fmter.add(fmtVerbStatic, format[prev:start])
		}

		name := format[m[2]:m[3]]
		verb := getFmtVerbByName(name)
		if verb == fmtVerbUnknown {
			return nil, errors.New("logger: unknown variable: " + name)
		}

		
		
		layout := defaultVerbsLayout[verb]
		if m[4] != -1 {
			layout = format[m[4]:m[5]]
		}
		if verb != fmtVerbTime && verb != fmtVerbLevelColor && verb != fmtVerbCallpath {
			layout = "%" + layout
		}

		fmter.add(verb, layout)
		prev = end
	}
	end := format[prev:]
	if end != "" {
		fmter.add(fmtVerbStatic, end)
	}

	
	t, err := time.Parse(time.RFC3339, "2010-02-04T21:00:57-08:00")
	if err != nil {
		panic(err)
	}
	testFmt := "hello %s"
	r := &Record{
		ID:     12345,
		Time:   t,
		Module: "logger",
		Args:   []interface{}{"go"},
		fmt:    &testFmt,
	}
	if err := fmter.Format(0, r, &bytes.Buffer{}); err != nil {
		return nil, err
	}

	return fmter, nil
}



func MustStringFormatter(format string) Formatter {
	f, err := NewStringFormatter(format)
	if err != nil {
		panic("Failed to initialized string formatter: " + err.Error())
	}
	return f
}

func (f *stringFormatter) add(verb fmtVerb, layout string) {
	f.parts = append(f.parts, part{verb, layout})
}

func (f *stringFormatter) Format(calldepth int, r *Record, output io.Writer) error {
	for _, part := range f.parts {
		if part.verb == fmtVerbStatic {
			output.Write([]byte(part.layout))
		} else if part.verb == fmtVerbTime {
			output.Write([]byte(r.Time.Format(part.layout)))
		} else if part.verb == fmtVerbLevelColor {
			doFmtVerbLevelColor(part.layout, r.Level, output)
		} else if part.verb == fmtVerbCallpath {
			depth, err := strconv.Atoi(part.layout)
			if err != nil {
				depth = 0
			}
			output.Write([]byte(formatCallpath(calldepth+1, depth)))
		} else {
			var v interface{}
			switch part.verb {
			case fmtVerbLevel:
				v = r.Level
				break
			case fmtVerbID:
				v = r.ID
				break
			case fmtVerbPid:
				v = pid
				break
			case fmtVerbProgram:
				v = program
				break
			case fmtVerbModule:
				v = r.Module
				break
			case fmtVerbMessage:
				v = r.Message()
				break
			case fmtVerbLongfile, fmtVerbShortfile:
				_, file, line, ok := runtime.Caller(calldepth + 1)
				if !ok {
					file = "???"
					line = 0
				} else if part.verb == fmtVerbShortfile {
					file = filepath.Base(file)
				}
				v = fmt.Sprintf("%s:%d", file, line)
			case fmtVerbLongfunc, fmtVerbShortfunc,
				fmtVerbLongpkg, fmtVerbShortpkg:
				
				v = "???"
				if pc, _, _, ok := runtime.Caller(calldepth + 1); ok {
					if f := runtime.FuncForPC(pc); f != nil {
						v = formatFuncName(part.verb, f.Name())
					}
				}
			default:
				panic("unhandled format part")
			}
			fmt.Fprintf(output, part.layout, v)
		}
	}
	return nil
}






func formatFuncName(v fmtVerb, f string) string {
	i := strings.LastIndex(f, "/")
	j := strings.Index(f[i+1:], ".")
	if j < 1 {
		return "???"
	}
	pkg, fun := f[:i+j+1], f[i+j+2:]
	switch v {
	case fmtVerbLongpkg:
		return pkg
	case fmtVerbShortpkg:
		return path.Base(pkg)
	case fmtVerbLongfunc:
		return fun
	case fmtVerbShortfunc:
		i = strings.LastIndex(fun, ".")
		return fun[i+1:]
	}
	panic("unexpected func formatter")
}

func formatCallpath(calldepth int, depth int) string {
	v := ""
	callers := make([]uintptr, 64)
	n := runtime.Callers(calldepth+2, callers)
	oldPc := callers[n-1]

	start := n - 3
	if depth > 0 && start >= depth {
		start = depth - 1
		v += "~."
	}
	recursiveCall := false
	for i := start; i >= 0; i-- {
		pc := callers[i]
		if oldPc == pc {
			recursiveCall = true
			continue
		}
		oldPc = pc
		if recursiveCall {
			recursiveCall = false
			v += ".."
		}
		if i < start {
			v += "."
		}
		if f := runtime.FuncForPC(pc); f != nil {
			v += formatFuncName(fmtVerbShortfunc, f.Name())
		}
	}
	return v
}



type backendFormatter struct {
	b Backend
	f Formatter
}



func NewBackendFormatter(b Backend, f Formatter) Backend {
	return &backendFormatter{b, f}
}


func (bf *backendFormatter) Log(level Level, calldepth int, r *Record) error {
	
	r2 := *r
	r2.formatter = bf.f
	return bf.b.Log(level, calldepth+1, &r2)
}
