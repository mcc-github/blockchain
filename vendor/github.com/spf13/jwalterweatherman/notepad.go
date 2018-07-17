




package jwalterweatherman

import (
	"fmt"
	"io"
	"log"
)

type Threshold int

func (t Threshold) String() string {
	return prefixes[t]
}

const (
	LevelTrace Threshold = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
	LevelFatal
)

var prefixes map[Threshold]string = map[Threshold]string{
	LevelTrace:    "TRACE",
	LevelDebug:    "DEBUG",
	LevelInfo:     "INFO",
	LevelWarn:     "WARN",
	LevelError:    "ERROR",
	LevelCritical: "CRITICAL",
	LevelFatal:    "FATAL",
}


type Notepad struct {
	TRACE    *log.Logger
	DEBUG    *log.Logger
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	CRITICAL *log.Logger
	FATAL    *log.Logger

	LOG      *log.Logger
	FEEDBACK *Feedback

	loggers         [7]**log.Logger
	logHandle       io.Writer
	outHandle       io.Writer
	logThreshold    Threshold
	stdoutThreshold Threshold
	prefix          string
	flags           int

	
	logCounters [7]*logCounter
}


func NewNotepad(outThreshold Threshold, logThreshold Threshold, outHandle, logHandle io.Writer, prefix string, flags int) *Notepad {
	n := &Notepad{}

	n.loggers = [7]**log.Logger{&n.TRACE, &n.DEBUG, &n.INFO, &n.WARN, &n.ERROR, &n.CRITICAL, &n.FATAL}
	n.outHandle = outHandle
	n.logHandle = logHandle
	n.stdoutThreshold = outThreshold
	n.logThreshold = logThreshold

	if len(prefix) != 0 {
		n.prefix = "[" + prefix + "] "
	} else {
		n.prefix = ""
	}

	n.flags = flags

	n.LOG = log.New(n.logHandle,
		"LOG:   ",
		n.flags)
	n.FEEDBACK = &Feedback{out: log.New(outHandle, "", 0), log: n.LOG}

	n.init()
	return n
}


func (n *Notepad) init() {
	logAndOut := io.MultiWriter(n.outHandle, n.logHandle)

	for t, logger := range n.loggers {
		threshold := Threshold(t)
		counter := &logCounter{}
		n.logCounters[t] = counter
		prefix := n.prefix + threshold.String() + " "

		switch {
		case threshold >= n.logThreshold && threshold >= n.stdoutThreshold:
			*logger = log.New(io.MultiWriter(counter, logAndOut), prefix, n.flags)

		case threshold >= n.logThreshold:
			*logger = log.New(io.MultiWriter(counter, n.logHandle), prefix, n.flags)

		case threshold >= n.stdoutThreshold:
			*logger = log.New(io.MultiWriter(counter, n.outHandle), prefix, n.flags)

		default:
			
			
			*logger = log.New(counter, "", 0)
		}
	}
}



func (n *Notepad) SetLogThreshold(threshold Threshold) {
	n.logThreshold = threshold
	n.init()
}


func (n *Notepad) SetLogOutput(handle io.Writer) {
	n.logHandle = handle
	n.init()
}


func (n *Notepad) GetLogThreshold() Threshold {
	return n.logThreshold
}



func (n *Notepad) SetStdoutThreshold(threshold Threshold) {
	n.stdoutThreshold = threshold
	n.init()
}


func (n *Notepad) GetStdoutThreshold() Threshold {
	return n.stdoutThreshold
}



func (n *Notepad) SetPrefix(prefix string) {
	if len(prefix) != 0 {
		n.prefix = "[" + prefix + "] "
	} else {
		n.prefix = ""
	}
	n.init()
}



func (n *Notepad) SetFlags(flags int) {
	n.flags = flags
	n.init()
}



type Feedback struct {
	out *log.Logger
	log *log.Logger
}

func (fb *Feedback) Println(v ...interface{}) {
	fb.output(fmt.Sprintln(v...))
}

func (fb *Feedback) Printf(format string, v ...interface{}) {
	fb.output(fmt.Sprintf(format, v...))
}

func (fb *Feedback) Print(v ...interface{}) {
	fb.output(fmt.Sprint(v...))
}

func (fb *Feedback) output(s string) {
	if fb.out != nil {
		fb.out.Output(2, s)
	}
	if fb.log != nil {
		fb.log.Output(2, s)
	}
}
