/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/



package healthz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type AlreadyRegisteredError string

func (are AlreadyRegisteredError) Error() string {
	return fmt.Sprintf("'%s' is already registered", are)
}

const (
	
	StatusOK = "OK"
	
	StatusUnavailable = "Service Unavailable"
)







type HealthChecker interface {
	HealthCheck(context.Context) error
}


type FailedCheck struct {
	Component string `json:"component"`
	Reason    string `json:"reason"`
}


type HealthStatus struct {
	Status       string        `json:"status"`
	Time         time.Time     `json:"time"`
	FailedChecks []FailedCheck `json:"failed_checks,omitempty"`
}


func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		healthCheckers: map[string]HealthChecker{},
		now:            time.Now,
		timeout:        30 * time.Second,
	}
}




type HealthHandler struct {
	mutex          sync.RWMutex
	healthCheckers map[string]HealthChecker
	now            func() time.Time
	timeout        time.Duration
}




func (h *HealthHandler) RegisterChecker(component string, checker HealthChecker) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.healthCheckers[component]; ok {
		return AlreadyRegisteredError(component)
	}
	h.healthCheckers[component] = checker
	return nil
}


func (h *HealthHandler) DeregisterChecker(component string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.healthCheckers, component)
}



func (h *HealthHandler) SetTimeout(timeout time.Duration) {
	h.timeout = timeout
}






func (h *HealthHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	checksCtx, cancel := context.WithTimeout(req.Context(), h.timeout)
	defer cancel()

	failedCheckCh := make(chan []FailedCheck)
	go func() {
		failedCheckCh <- h.RunChecks(checksCtx)
	}()

	select {
	case failedChecks := <-failedCheckCh:
		hs := HealthStatus{
			Status: StatusOK,
			Time:   h.now(),
		}
		if len(failedChecks) > 0 {
			hs.Status = StatusUnavailable
			hs.FailedChecks = failedChecks
		}
		writeHTTPResponse(rw, hs)
	case <-checksCtx.Done():
		if checksCtx.Err() == context.DeadlineExceeded {
			rw.WriteHeader(http.StatusRequestTimeout)
		}
	}
}


func (h *HealthHandler) RunChecks(ctx context.Context) []FailedCheck {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var failedChecks []FailedCheck
	for component, checker := range h.healthCheckers {
		if err := checker.HealthCheck(ctx); err != nil {
			failedCheck := FailedCheck{
				Component: component,
				Reason:    err.Error(),
			}
			failedChecks = append(failedChecks, failedCheck)
		}
	}
	return failedChecks
}


func writeHTTPResponse(rw http.ResponseWriter, hs HealthStatus) {
	var resp []byte
	rc := http.StatusOK
	rw.Header().Set("Content-Type", "application/json")
	if len(hs.FailedChecks) > 0 {
		rc = http.StatusServiceUnavailable
	}
	resp, err := json.Marshal(hs)
	if err != nil {
		rc = http.StatusInternalServerError
	}
	rw.WriteHeader(rc)
	rw.Write(resp)
}
