

package binarylog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/grpc/grpclog"
)

















func NewLoggerFromConfigString(s string) Logger {
	if s == "" {
		return nil
	}
	l := newEmptyLogger()
	methods := strings.Split(s, ",")
	for _, method := range methods {
		if err := l.fillMethodLoggerWithConfigString(method); err != nil {
			grpclog.Warningf("failed to parse binary log config: %v", err)
			return nil
		}
	}
	return l
}



func (l *logger) fillMethodLoggerWithConfigString(config string) error {
	
	if config == "" {
		return errors.New("empty string is not a valid method binary logging config")
	}

	
	if config[0] == '-' {
		s, m, suffix, err := parseMethodConfigAndSuffix(config[1:])
		if err != nil {
			return fmt.Errorf("invalid config: %q, %v", config, err)
		}
		if m == "*" {
			return fmt.Errorf("invalid config: %q, %v", config, "* not allowd in blacklist config")
		}
		if suffix != "" {
			return fmt.Errorf("invalid config: %q, %v", config, "header/message limit not allowed in blacklist config")
		}
		if err := l.setBlacklist(s + "/" + m); err != nil {
			return fmt.Errorf("invalid config: %v", err)
		}
		return nil
	}

	
	if config[0] == '*' {
		hdr, msg, err := parseHeaderMessageLengthConfig(config[1:])
		if err != nil {
			return fmt.Errorf("invalid config: %q, %v", config, err)
		}
		if err := l.setDefaultMethodLogger(&methodLoggerConfig{hdr: hdr, msg: msg}); err != nil {
			return fmt.Errorf("invalid config: %v", err)
		}
		return nil
	}

	s, m, suffix, err := parseMethodConfigAndSuffix(config)
	if err != nil {
		return fmt.Errorf("invalid config: %q, %v", config, err)
	}
	hdr, msg, err := parseHeaderMessageLengthConfig(suffix)
	if err != nil {
		return fmt.Errorf("invalid header/message length config: %q, %v", suffix, err)
	}
	if m == "*" {
		if err := l.setServiceMethodLogger(s, &methodLoggerConfig{hdr: hdr, msg: msg}); err != nil {
			return fmt.Errorf("invalid config: %v", err)
		}
	} else {
		if err := l.setMethodMethodLogger(s+"/"+m, &methodLoggerConfig{hdr: hdr, msg: msg}); err != nil {
			return fmt.Errorf("invalid config: %v", err)
		}
	}
	return nil
}

const (
	
	
	maxUInt = ^uint64(0)

	
	
	longMethodConfigRegexpStr = `^([\w./]+)/((?:\w+)|[*])(.+)?$`

	
	optionalLengthRegexpStr      = `(?::(\d+))?` 
	headerConfigRegexpStr        = `^{h` + optionalLengthRegexpStr + `}$`
	messageConfigRegexpStr       = `^{m` + optionalLengthRegexpStr + `}$`
	headerMessageConfigRegexpStr = `^{h` + optionalLengthRegexpStr + `;m` + optionalLengthRegexpStr + `}$`
)

var (
	longMethodConfigRegexp    = regexp.MustCompile(longMethodConfigRegexpStr)
	headerConfigRegexp        = regexp.MustCompile(headerConfigRegexpStr)
	messageConfigRegexp       = regexp.MustCompile(messageConfigRegexpStr)
	headerMessageConfigRegexp = regexp.MustCompile(headerMessageConfigRegexpStr)
)


func parseMethodConfigAndSuffix(c string) (service, method, suffix string, _ error) {
	
	
	
	
	match := longMethodConfigRegexp.FindStringSubmatch(c)
	if match == nil {
		return "", "", "", fmt.Errorf("%q contains invalid substring", c)
	}
	service = match[1]
	method = match[2]
	suffix = match[3]
	return
}




func parseHeaderMessageLengthConfig(c string) (hdrLenStr, msgLenStr uint64, err error) {
	if c == "" {
		return maxUInt, maxUInt, nil
	}
	
	if match := headerConfigRegexp.FindStringSubmatch(c); match != nil {
		if s := match[1]; s != "" {
			hdrLenStr, err = strconv.ParseUint(s, 10, 64)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to convert %q to uint", s)
			}
			return hdrLenStr, 0, nil
		}
		return maxUInt, 0, nil
	}

	
	if match := messageConfigRegexp.FindStringSubmatch(c); match != nil {
		if s := match[1]; s != "" {
			msgLenStr, err = strconv.ParseUint(s, 10, 64)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to convert %q to uint", s)
			}
			return 0, msgLenStr, nil
		}
		return 0, maxUInt, nil
	}

	
	if match := headerMessageConfigRegexp.FindStringSubmatch(c); match != nil {
		
		hdrLenStr = maxUInt
		msgLenStr = maxUInt
		if s := match[1]; s != "" {
			hdrLenStr, err = strconv.ParseUint(s, 10, 64)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to convert %q to uint", s)
			}
		}
		if s := match[2]; s != "" {
			msgLenStr, err = strconv.ParseUint(s, 10, 64)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to convert %q to uint", s)
			}
		}
		return hdrLenStr, msgLenStr, nil
	}
	return 0, 0, fmt.Errorf("%q contains invalid substring", c)
}
