package logrus

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}


var ErrorKey = "error"





type Entry struct {
	Logger *Logger

	
	Data Fields

	
	Time time.Time

	
	
	Level Level

	
	Message string

	
	Buffer *bytes.Buffer
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		
		Data: make(Fields, 5),
	}
}



func (entry *Entry) String() (string, error) {
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}


func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField(ErrorKey, err)
}


func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}


func (entry *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time}
}


func (entry *Entry) WithTime(t time.Time) *Entry {
	return &Entry{Logger: entry.Logger, Data: entry.Data, Time: t}
}



func (entry Entry) log(level Level, msg string) {
	var buffer *bytes.Buffer

	
	
	
	
	
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}

	entry.Level = level
	entry.Message = msg

	entry.fireHooks()

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer

	entry.write()

	entry.Buffer = nil

	
	
	
	if level <= PanicLevel {
		panic(&entry)
	}
}

func (entry *Entry) fireHooks() {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	err := entry.Logger.Hooks.Fire(entry.Level, entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
	}
}

func (entry *Entry) write() {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
	} else {
		_, err = entry.Logger.Out.Write(serialized)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
	}
}

func (entry *Entry) Debug(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(DebugLevel) {
		entry.log(DebugLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Print(args ...interface{}) {
	entry.Info(args...)
}

func (entry *Entry) Info(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(InfoLevel) {
		entry.log(InfoLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warn(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(WarnLevel) {
		entry.log(WarnLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warning(args ...interface{}) {
	entry.Warn(args...)
}

func (entry *Entry) Error(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(ErrorLevel) {
		entry.log(ErrorLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(FatalLevel) {
		entry.log(FatalLevel, fmt.Sprint(args...))
	}
	Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(PanicLevel) {
		entry.log(PanicLevel, fmt.Sprint(args...))
	}
	panic(fmt.Sprint(args...))
}



func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(DebugLevel) {
		entry.Debug(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(InfoLevel) {
		entry.Info(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(WarnLevel) {
		entry.Warn(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(ErrorLevel) {
		entry.Error(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(FatalLevel) {
		entry.Fatal(fmt.Sprintf(format, args...))
	}
	Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(PanicLevel) {
		entry.Panic(fmt.Sprintf(format, args...))
	}
}



func (entry *Entry) Debugln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(DebugLevel) {
		entry.Debug(entry.sprintlnn(args...))
	}
}

func (entry *Entry) Infoln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(InfoLevel) {
		entry.Info(entry.sprintlnn(args...))
	}
}

func (entry *Entry) Println(args ...interface{}) {
	entry.Infoln(args...)
}

func (entry *Entry) Warnln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(WarnLevel) {
		entry.Warn(entry.sprintlnn(args...))
	}
}

func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warnln(args...)
}

func (entry *Entry) Errorln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(ErrorLevel) {
		entry.Error(entry.sprintlnn(args...))
	}
}

func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(FatalLevel) {
		entry.Fatal(entry.sprintlnn(args...))
	}
	Exit(1)
}

func (entry *Entry) Panicln(args ...interface{}) {
	if entry.Logger.IsLevelEnabled(PanicLevel) {
		entry.Panic(entry.sprintlnn(args...))
	}
}





func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
