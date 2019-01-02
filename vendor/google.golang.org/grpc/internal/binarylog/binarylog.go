



package binarylog

import (
	"fmt"
	"os"

	"google.golang.org/grpc/grpclog"
)



type Logger interface {
	getMethodLogger(methodName string) *MethodLogger
}





var binLogger Logger




func SetLogger(l Logger) {
	binLogger = l
}







func GetMethodLogger(methodName string) *MethodLogger {
	if binLogger == nil {
		return nil
	}
	return binLogger.getMethodLogger(methodName)
}

func init() {
	const envStr = "GRPC_BINARY_LOG_FILTER"
	configStr := os.Getenv(envStr)
	binLogger = NewLoggerFromConfigString(configStr)
}

type methodLoggerConfig struct {
	
	hdr, msg uint64
}

type logger struct {
	all      *methodLoggerConfig
	services map[string]*methodLoggerConfig
	methods  map[string]*methodLoggerConfig

	blacklist map[string]struct{}
}



func newEmptyLogger() *logger {
	return &logger{}
}


func (l *logger) setDefaultMethodLogger(ml *methodLoggerConfig) error {
	if l.all != nil {
		return fmt.Errorf("conflicting global rules found")
	}
	l.all = ml
	return nil
}




func (l *logger) setServiceMethodLogger(service string, ml *methodLoggerConfig) error {
	if _, ok := l.services[service]; ok {
		return fmt.Errorf("conflicting rules for service %v found", service)
	}
	if l.services == nil {
		l.services = make(map[string]*methodLoggerConfig)
	}
	l.services[service] = ml
	return nil
}




func (l *logger) setMethodMethodLogger(method string, ml *methodLoggerConfig) error {
	if _, ok := l.blacklist[method]; ok {
		return fmt.Errorf("conflicting rules for method %v found", method)
	}
	if _, ok := l.methods[method]; ok {
		return fmt.Errorf("conflicting rules for method %v found", method)
	}
	if l.methods == nil {
		l.methods = make(map[string]*methodLoggerConfig)
	}
	l.methods[method] = ml
	return nil
}


func (l *logger) setBlacklist(method string) error {
	if _, ok := l.blacklist[method]; ok {
		return fmt.Errorf("conflicting rules for method %v found", method)
	}
	if _, ok := l.methods[method]; ok {
		return fmt.Errorf("conflicting rules for method %v found", method)
	}
	if l.blacklist == nil {
		l.blacklist = make(map[string]struct{})
	}
	l.blacklist[method] = struct{}{}
	return nil
}







func (l *logger) getMethodLogger(methodName string) *MethodLogger {
	s, m, err := parseMethodName(methodName)
	if err != nil {
		grpclog.Infof("binarylogging: failed to parse %q: %v", methodName, err)
		return nil
	}
	if ml, ok := l.methods[s+"/"+m]; ok {
		return newMethodLogger(ml.hdr, ml.msg)
	}
	if _, ok := l.blacklist[s+"/"+m]; ok {
		return nil
	}
	if ml, ok := l.services[s]; ok {
		return newMethodLogger(ml.hdr, ml.msg)
	}
	if l.all == nil {
		return nil
	}
	return newMethodLogger(l.all.hdr, l.all.msg)
}
