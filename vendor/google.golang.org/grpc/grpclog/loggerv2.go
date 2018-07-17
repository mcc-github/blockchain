

package grpclog

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)


type LoggerV2 interface {
	
	Info(args ...interface{})
	
	Infoln(args ...interface{})
	
	Infof(format string, args ...interface{})
	
	Warning(args ...interface{})
	
	Warningln(args ...interface{})
	
	Warningf(format string, args ...interface{})
	
	Error(args ...interface{})
	
	Errorln(args ...interface{})
	
	Errorf(format string, args ...interface{})
	
	
	
	Fatal(args ...interface{})
	
	
	
	Fatalln(args ...interface{})
	
	
	
	Fatalf(format string, args ...interface{})
	
	V(l int) bool
}



func SetLoggerV2(l LoggerV2) {
	logger = l
}

const (
	
	infoLog int = iota
	
	warningLog
	
	errorLog
	
	fatalLog
)


var severityName = []string{
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}


type loggerT struct {
	m []*log.Logger
	v int
}






func NewLoggerV2(infoW, warningW, errorW io.Writer) LoggerV2 {
	return NewLoggerV2WithVerbosity(infoW, warningW, errorW, 0)
}



func NewLoggerV2WithVerbosity(infoW, warningW, errorW io.Writer, v int) LoggerV2 {
	var m []*log.Logger
	m = append(m, log.New(infoW, severityName[infoLog]+": ", log.LstdFlags))
	m = append(m, log.New(io.MultiWriter(infoW, warningW), severityName[warningLog]+": ", log.LstdFlags))
	ew := io.MultiWriter(infoW, warningW, errorW) 
	m = append(m, log.New(ew, severityName[errorLog]+": ", log.LstdFlags))
	m = append(m, log.New(ew, severityName[fatalLog]+": ", log.LstdFlags))
	return &loggerT{m: m, v: v}
}



func newLoggerV2() LoggerV2 {
	errorW := ioutil.Discard
	warningW := ioutil.Discard
	infoW := ioutil.Discard

	logLevel := os.Getenv("GRPC_GO_LOG_SEVERITY_LEVEL")
	switch logLevel {
	case "", "ERROR", "error": 
		errorW = os.Stderr
	case "WARNING", "warning":
		warningW = os.Stderr
	case "INFO", "info":
		infoW = os.Stderr
	}

	var v int
	vLevel := os.Getenv("GRPC_GO_LOG_VERBOSITY_LEVEL")
	if vl, err := strconv.Atoi(vLevel); err == nil {
		v = vl
	}
	return NewLoggerV2WithVerbosity(infoW, warningW, errorW, v)
}

func (g *loggerT) Info(args ...interface{}) {
	g.m[infoLog].Print(args...)
}

func (g *loggerT) Infoln(args ...interface{}) {
	g.m[infoLog].Println(args...)
}

func (g *loggerT) Infof(format string, args ...interface{}) {
	g.m[infoLog].Printf(format, args...)
}

func (g *loggerT) Warning(args ...interface{}) {
	g.m[warningLog].Print(args...)
}

func (g *loggerT) Warningln(args ...interface{}) {
	g.m[warningLog].Println(args...)
}

func (g *loggerT) Warningf(format string, args ...interface{}) {
	g.m[warningLog].Printf(format, args...)
}

func (g *loggerT) Error(args ...interface{}) {
	g.m[errorLog].Print(args...)
}

func (g *loggerT) Errorln(args ...interface{}) {
	g.m[errorLog].Println(args...)
}

func (g *loggerT) Errorf(format string, args ...interface{}) {
	g.m[errorLog].Printf(format, args...)
}

func (g *loggerT) Fatal(args ...interface{}) {
	g.m[fatalLog].Fatal(args...)
	
}

func (g *loggerT) Fatalln(args ...interface{}) {
	g.m[fatalLog].Fatalln(args...)
	
}

func (g *loggerT) Fatalf(format string, args ...interface{}) {
	g.m[fatalLog].Fatalf(format, args...)
	
}

func (g *loggerT) V(l int) bool {
	return l <= g.v
}
