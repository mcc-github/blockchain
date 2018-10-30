













package capnslog

import (
	"errors"
	"strings"
	"sync"
)


type LogLevel int8

const (
	
	CRITICAL LogLevel = iota - 1
	
	ERROR
	
	WARNING
	
	NOTICE
	
	INFO
	
	DEBUG
	
	TRACE
)


func (l LogLevel) Char() string {
	switch l {
	case CRITICAL:
		return "C"
	case ERROR:
		return "E"
	case WARNING:
		return "W"
	case NOTICE:
		return "N"
	case INFO:
		return "I"
	case DEBUG:
		return "D"
	case TRACE:
		return "T"
	default:
		panic("Unhandled loglevel")
	}
}


func (l LogLevel) String() string {
	switch l {
	case CRITICAL:
		return "CRITICAL"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case NOTICE:
		return "NOTICE"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	default:
		panic("Unhandled loglevel")
	}
}


func (l *LogLevel) Set(s string) error {
	value, err := ParseLevel(s)
	if err != nil {
		return err
	}

	*l = value
	return nil
}


func (l *LogLevel) Type() string {
	return ""
}


func ParseLevel(s string) (LogLevel, error) {
	switch s {
	case "CRITICAL", "C":
		return CRITICAL, nil
	case "ERROR", "0", "E":
		return ERROR, nil
	case "WARNING", "1", "W":
		return WARNING, nil
	case "NOTICE", "2", "N":
		return NOTICE, nil
	case "INFO", "3", "I":
		return INFO, nil
	case "DEBUG", "4", "D":
		return DEBUG, nil
	case "TRACE", "5", "T":
		return TRACE, nil
	}
	return CRITICAL, errors.New("couldn't parse log level " + s)
}

type RepoLogger map[string]*PackageLogger

type loggerStruct struct {
	sync.Mutex
	repoMap   map[string]RepoLogger
	formatter Formatter
}


var logger = new(loggerStruct)



func SetGlobalLogLevel(l LogLevel) {
	logger.Lock()
	defer logger.Unlock()
	for _, r := range logger.repoMap {
		r.setRepoLogLevelInternal(l)
	}
}


func GetRepoLogger(repo string) (RepoLogger, error) {
	logger.Lock()
	defer logger.Unlock()
	r, ok := logger.repoMap[repo]
	if !ok {
		return nil, errors.New("no packages registered for repo " + repo)
	}
	return r, nil
}


func MustRepoLogger(repo string) RepoLogger {
	r, err := GetRepoLogger(repo)
	if err != nil {
		panic(err)
	}
	return r
}


func (r RepoLogger) SetRepoLogLevel(l LogLevel) {
	logger.Lock()
	defer logger.Unlock()
	r.setRepoLogLevelInternal(l)
}

func (r RepoLogger) setRepoLogLevelInternal(l LogLevel) {
	for _, v := range r {
		v.level = l
	}
}



func (r RepoLogger) ParseLogLevelConfig(conf string) (map[string]LogLevel, error) {
	setlist := strings.Split(conf, ",")
	out := make(map[string]LogLevel)
	for _, setstring := range setlist {
		setting := strings.Split(setstring, "=")
		if len(setting) != 2 {
			return nil, errors.New("oddly structured `pkg=level` option: " + setstring)
		}
		l, err := ParseLevel(setting[1])
		if err != nil {
			return nil, err
		}
		out[setting[0]] = l
	}
	return out, nil
}





func (r RepoLogger) SetLogLevel(m map[string]LogLevel) {
	logger.Lock()
	defer logger.Unlock()
	if l, ok := m["*"]; ok {
		r.setRepoLogLevelInternal(l)
	}
	for k, v := range m {
		l, ok := r[k]
		if !ok {
			continue
		}
		l.level = v
	}
}


func SetFormatter(f Formatter) {
	logger.Lock()
	defer logger.Unlock()
	logger.formatter = f
}



func NewPackageLogger(repo string, pkg string) (p *PackageLogger) {
	logger.Lock()
	defer logger.Unlock()
	if logger.repoMap == nil {
		logger.repoMap = make(map[string]RepoLogger)
	}
	r, rok := logger.repoMap[repo]
	if !rok {
		logger.repoMap[repo] = make(RepoLogger)
		r = logger.repoMap[repo]
	}
	p, pok := r[pkg]
	if !pok {
		r[pkg] = &PackageLogger{
			pkg:   pkg,
			level: INFO,
		}
		p = r[pkg]
	}
	return
}
