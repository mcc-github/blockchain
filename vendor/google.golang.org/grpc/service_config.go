

package grpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
)

const maxInt = int(^uint(0) >> 1)







type MethodConfig struct {
	
	
	
	WaitForReady *bool
	
	
	
	
	Timeout *time.Duration
	
	
	
	
	
	
	MaxReqSize *int
	
	
	MaxRespSize *int
	
	retryPolicy *retryPolicy
}







type ServiceConfig struct {
	
	
	LB *string

	
	
	
	
	
	
	Methods map[string]MethodConfig

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	retryThrottling *retryThrottlingPolicy
}




type retryPolicy struct {
	
	
	
	maxAttempts int

	
	
	
	
	
	
	initialBackoff    time.Duration
	maxBackoff        time.Duration
	backoffMultiplier float64

	
	
	
	
	
	
	retryableStatusCodes map[codes.Code]bool
}

type jsonRetryPolicy struct {
	MaxAttempts          int
	InitialBackoff       string
	MaxBackoff           string
	BackoffMultiplier    float64
	RetryableStatusCodes []codes.Code
}




type retryThrottlingPolicy struct {
	
	
	
	
	MaxTokens float64
	
	
	
	
	
	TokenRatio float64
}

func parseDuration(s *string) (*time.Duration, error) {
	if s == nil {
		return nil, nil
	}
	if !strings.HasSuffix(*s, "s") {
		return nil, fmt.Errorf("malformed duration %q", *s)
	}
	ss := strings.SplitN((*s)[:len(*s)-1], ".", 3)
	if len(ss) > 2 {
		return nil, fmt.Errorf("malformed duration %q", *s)
	}
	
	
	hasDigits := false
	var d time.Duration
	if len(ss[0]) > 0 {
		i, err := strconv.ParseInt(ss[0], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("malformed duration %q: %v", *s, err)
		}
		d = time.Duration(i) * time.Second
		hasDigits = true
	}
	if len(ss) == 2 && len(ss[1]) > 0 {
		if len(ss[1]) > 9 {
			return nil, fmt.Errorf("malformed duration %q", *s)
		}
		f, err := strconv.ParseInt(ss[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("malformed duration %q: %v", *s, err)
		}
		for i := 9; i > len(ss[1]); i-- {
			f *= 10
		}
		d += time.Duration(f)
		hasDigits = true
	}
	if !hasDigits {
		return nil, fmt.Errorf("malformed duration %q", *s)
	}

	return &d, nil
}

type jsonName struct {
	Service *string
	Method  *string
}

func (j jsonName) generatePath() (string, bool) {
	if j.Service == nil {
		return "", false
	}
	res := "/" + *j.Service + "/"
	if j.Method != nil {
		res += *j.Method
	}
	return res, true
}


type jsonMC struct {
	Name                    *[]jsonName
	WaitForReady            *bool
	Timeout                 *string
	MaxRequestMessageBytes  *int64
	MaxResponseMessageBytes *int64
	RetryPolicy             *jsonRetryPolicy
}


type jsonSC struct {
	LoadBalancingPolicy *string
	MethodConfig        *[]jsonMC
	RetryThrottling     *retryThrottlingPolicy
}

func parseServiceConfig(js string) (ServiceConfig, error) {
	var rsc jsonSC
	err := json.Unmarshal([]byte(js), &rsc)
	if err != nil {
		grpclog.Warningf("grpc: parseServiceConfig error unmarshaling %s due to %v", js, err)
		return ServiceConfig{}, err
	}
	sc := ServiceConfig{
		LB:              rsc.LoadBalancingPolicy,
		Methods:         make(map[string]MethodConfig),
		retryThrottling: rsc.RetryThrottling,
	}
	if rsc.MethodConfig == nil {
		return sc, nil
	}

	for _, m := range *rsc.MethodConfig {
		if m.Name == nil {
			continue
		}
		d, err := parseDuration(m.Timeout)
		if err != nil {
			grpclog.Warningf("grpc: parseServiceConfig error unmarshaling %s due to %v", js, err)
			return ServiceConfig{}, err
		}

		mc := MethodConfig{
			WaitForReady: m.WaitForReady,
			Timeout:      d,
		}
		if mc.retryPolicy, err = convertRetryPolicy(m.RetryPolicy); err != nil {
			grpclog.Warningf("grpc: parseServiceConfig error unmarshaling %s due to %v", js, err)
			return ServiceConfig{}, err
		}
		if m.MaxRequestMessageBytes != nil {
			if *m.MaxRequestMessageBytes > int64(maxInt) {
				mc.MaxReqSize = newInt(maxInt)
			} else {
				mc.MaxReqSize = newInt(int(*m.MaxRequestMessageBytes))
			}
		}
		if m.MaxResponseMessageBytes != nil {
			if *m.MaxResponseMessageBytes > int64(maxInt) {
				mc.MaxRespSize = newInt(maxInt)
			} else {
				mc.MaxRespSize = newInt(int(*m.MaxResponseMessageBytes))
			}
		}
		for _, n := range *m.Name {
			if path, valid := n.generatePath(); valid {
				sc.Methods[path] = mc
			}
		}
	}

	if sc.retryThrottling != nil {
		if sc.retryThrottling.MaxTokens <= 0 ||
			sc.retryThrottling.MaxTokens >= 1000 ||
			sc.retryThrottling.TokenRatio <= 0 {
			
			sc.retryThrottling = nil
		}
	}
	return sc, nil
}

func convertRetryPolicy(jrp *jsonRetryPolicy) (p *retryPolicy, err error) {
	if jrp == nil {
		return nil, nil
	}
	ib, err := parseDuration(&jrp.InitialBackoff)
	if err != nil {
		return nil, err
	}
	mb, err := parseDuration(&jrp.MaxBackoff)
	if err != nil {
		return nil, err
	}

	if jrp.MaxAttempts <= 1 ||
		*ib <= 0 ||
		*mb <= 0 ||
		jrp.BackoffMultiplier <= 0 ||
		len(jrp.RetryableStatusCodes) == 0 {
		grpclog.Warningf("grpc: ignoring retry policy %v due to illegal configuration", jrp)
		return nil, nil
	}

	rp := &retryPolicy{
		maxAttempts:          jrp.MaxAttempts,
		initialBackoff:       *ib,
		maxBackoff:           *mb,
		backoffMultiplier:    jrp.BackoffMultiplier,
		retryableStatusCodes: make(map[codes.Code]bool),
	}
	if rp.maxAttempts > 5 {
		
		rp.maxAttempts = 5
	}
	for _, code := range jrp.RetryableStatusCodes {
		rp.retryableStatusCodes[code] = true
	}
	return rp, nil
}

func min(a, b *int) *int {
	if *a < *b {
		return a
	}
	return b
}

func getMaxSize(mcMax, doptMax *int, defaultVal int) *int {
	if mcMax == nil && doptMax == nil {
		return &defaultVal
	}
	if mcMax != nil && doptMax != nil {
		return min(mcMax, doptMax)
	}
	if mcMax != nil {
		return mcMax
	}
	return doptMax
}

func newInt(b int) *int {
	return &b
}
