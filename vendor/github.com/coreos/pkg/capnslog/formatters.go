













package capnslog

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"time"
)

type Formatter interface {
	Format(pkg string, level LogLevel, depth int, entries ...interface{})
	Flush()
}

func NewStringFormatter(w io.Writer) Formatter {
	return &StringFormatter{
		w: bufio.NewWriter(w),
	}
}

type StringFormatter struct {
	w *bufio.Writer
}

func (s *StringFormatter) Format(pkg string, l LogLevel, i int, entries ...interface{}) {
	now := time.Now().UTC()
	s.w.WriteString(now.Format(time.RFC3339))
	s.w.WriteByte(' ')
	writeEntries(s.w, pkg, l, i, entries...)
	s.Flush()
}

func writeEntries(w *bufio.Writer, pkg string, _ LogLevel, _ int, entries ...interface{}) {
	if pkg != "" {
		w.WriteString(pkg + ": ")
	}
	str := fmt.Sprint(entries...)
	endsInNL := strings.HasSuffix(str, "\n")
	w.WriteString(str)
	if !endsInNL {
		w.WriteString("\n")
	}
}

func (s *StringFormatter) Flush() {
	s.w.Flush()
}

func NewPrettyFormatter(w io.Writer, debug bool) Formatter {
	return &PrettyFormatter{
		w:     bufio.NewWriter(w),
		debug: debug,
	}
}

type PrettyFormatter struct {
	w     *bufio.Writer
	debug bool
}

func (c *PrettyFormatter) Format(pkg string, l LogLevel, depth int, entries ...interface{}) {
	now := time.Now()
	ts := now.Format("2006-01-02 15:04:05")
	c.w.WriteString(ts)
	ms := now.Nanosecond() / 1000
	c.w.WriteString(fmt.Sprintf(".%06d", ms))
	if c.debug {
		_, file, line, ok := runtime.Caller(depth) 
		if !ok {
			file = "???"
			line = 1
		} else {
			slash := strings.LastIndex(file, "/")
			if slash >= 0 {
				file = file[slash+1:]
			}
		}
		if line < 0 {
			line = 0 
		}
		c.w.WriteString(fmt.Sprintf(" [%s:%d]", file, line))
	}
	c.w.WriteString(fmt.Sprint(" ", l.Char(), " | "))
	writeEntries(c.w, pkg, l, depth, entries...)
	c.Flush()
}

func (c *PrettyFormatter) Flush() {
	c.w.Flush()
}


type LogFormatter struct {
	logger *log.Logger
	prefix string
}



func NewLogFormatter(w io.Writer, prefix string, flag int) Formatter {
	return &LogFormatter{
		logger: log.New(w, "", flag), 
		prefix: prefix,               
	}
}


func (lf *LogFormatter) Format(pkg string, _ LogLevel, _ int, entries ...interface{}) {
	str := fmt.Sprint(entries...)
	prefix := lf.prefix
	if pkg != "" {
		prefix = fmt.Sprintf("%s%s: ", prefix, pkg)
	}
	lf.logger.Output(5, fmt.Sprintf("%s%v", prefix, str)) 
}


func (lf *LogFormatter) Flush() {
	
}


type NilFormatter struct {
}



func NewNilFormatter() Formatter {
	return &NilFormatter{}
}


func (_ *NilFormatter) Format(_ string, _ LogLevel, _ int, _ ...interface{}) {
	
}


func (_ *NilFormatter) Flush() {
	
}
