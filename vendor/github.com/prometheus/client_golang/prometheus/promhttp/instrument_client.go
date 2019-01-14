












package promhttp

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)




type RoundTripperFunc func(req *http.Request) (*http.Response, error)


func (rt RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}






func InstrumentRoundTripperInFlight(gauge prometheus.Gauge, next http.RoundTripper) RoundTripperFunc {
	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		gauge.Inc()
		defer gauge.Dec()
		return next.RoundTrip(r)
	})
}













func InstrumentRoundTripperCounter(counter *prometheus.CounterVec, next http.RoundTripper) RoundTripperFunc {
	code, method := checkLabels(counter)

	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		resp, err := next.RoundTrip(r)
		if err == nil {
			counter.With(labels(code, method, r.Method, resp.StatusCode)).Inc()
		}
		return resp, err
	})
}

















func InstrumentRoundTripperDuration(obs prometheus.ObserverVec, next http.RoundTripper) RoundTripperFunc {
	code, method := checkLabels(obs)

	return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		start := time.Now()
		resp, err := next.RoundTrip(r)
		if err == nil {
			obs.With(labels(code, method, r.Method, resp.StatusCode)).Observe(time.Since(start).Seconds())
		}
		return resp, err
	})
}
