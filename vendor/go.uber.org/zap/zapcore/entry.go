



















package zapcore

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap/internal/bufferpool"
	"go.uber.org/zap/internal/exit"

	"go.uber.org/multierr"
)

var (
	_cePool = sync.Pool{New: func() interface{} {
		
		return &CheckedEntry{
			cores: make([]Core, 4),
		}
	}}
)

func getCheckedEntry() *CheckedEntry {
	ce := _cePool.Get().(*CheckedEntry)
	ce.reset()
	return ce
}

func putCheckedEntry(ce *CheckedEntry) {
	if ce == nil {
		return
	}
	_cePool.Put(ce)
}



func NewEntryCaller(pc uintptr, file string, line int, ok bool) EntryCaller {
	if !ok {
		return EntryCaller{}
	}
	return EntryCaller{
		PC:      pc,
		File:    file,
		Line:    line,
		Defined: true,
	}
}


type EntryCaller struct {
	Defined bool
	PC      uintptr
	File    string
	Line    int
}


func (ec EntryCaller) String() string {
	return ec.FullPath()
}



func (ec EntryCaller) FullPath() string {
	if !ec.Defined {
		return "undefined"
	}
	buf := bufferpool.Get()
	buf.AppendString(ec.File)
	buf.AppendByte(':')
	buf.AppendInt(int64(ec.Line))
	caller := buf.String()
	buf.Free()
	return caller
}



func (ec EntryCaller) TrimmedPath() string {
	if !ec.Defined {
		return "undefined"
	}
	
	
	
	
	
	
	
	
	
	
	
	
	
	idx := strings.LastIndexByte(ec.File, '/')
	if idx == -1 {
		return ec.FullPath()
	}
	
	idx = strings.LastIndexByte(ec.File[:idx], '/')
	if idx == -1 {
		return ec.FullPath()
	}
	buf := bufferpool.Get()
	
	buf.AppendString(ec.File[idx+1:])
	buf.AppendByte(':')
	buf.AppendInt(int64(ec.Line))
	caller := buf.String()
	buf.Free()
	return caller
}







type Entry struct {
	Level      Level
	Time       time.Time
	LoggerName string
	Message    string
	Caller     EntryCaller
	Stack      string
}



type CheckWriteAction uint8

const (
	
	
	WriteThenNoop CheckWriteAction = iota
	
	WriteThenPanic
	
	WriteThenFatal
)







type CheckedEntry struct {
	Entry
	ErrorOutput WriteSyncer
	dirty       bool 
	should      CheckWriteAction
	cores       []Core
}

func (ce *CheckedEntry) reset() {
	ce.Entry = Entry{}
	ce.ErrorOutput = nil
	ce.dirty = false
	ce.should = WriteThenNoop
	for i := range ce.cores {
		
		ce.cores[i] = nil
	}
	ce.cores = ce.cores[:0]
}




func (ce *CheckedEntry) Write(fields ...Field) {
	if ce == nil {
		return
	}

	if ce.dirty {
		if ce.ErrorOutput != nil {
			
			
			
			
			fmt.Fprintf(ce.ErrorOutput, "%v Unsafe CheckedEntry re-use near Entry %+v.\n", time.Now(), ce.Entry)
			ce.ErrorOutput.Sync()
		}
		return
	}
	ce.dirty = true

	var err error
	for i := range ce.cores {
		err = multierr.Append(err, ce.cores[i].Write(ce.Entry, fields))
	}
	if ce.ErrorOutput != nil {
		if err != nil {
			fmt.Fprintf(ce.ErrorOutput, "%v write error: %v\n", time.Now(), err)
			ce.ErrorOutput.Sync()
		}
	}

	should, msg := ce.should, ce.Message
	putCheckedEntry(ce)

	switch should {
	case WriteThenPanic:
		panic(msg)
	case WriteThenFatal:
		exit.Exit()
	}
}




func (ce *CheckedEntry) AddCore(ent Entry, core Core) *CheckedEntry {
	if ce == nil {
		ce = getCheckedEntry()
		ce.Entry = ent
	}
	ce.cores = append(ce.cores, core)
	return ce
}




func (ce *CheckedEntry) Should(ent Entry, should CheckWriteAction) *CheckedEntry {
	if ce == nil {
		ce = getCheckedEntry()
		ce.Entry = ent
	}
	ce.should = should
	return ce
}
