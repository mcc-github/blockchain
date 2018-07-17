



package statsd

import "time"


type NoopClient struct{}


func (s *NoopClient) Close() error {
	return nil
}





func (s *NoopClient) Inc(stat string, value int64, rate float32) error {
	return nil
}





func (s *NoopClient) Dec(stat string, value int64, rate float32) error {
	return nil
}





func (s *NoopClient) Gauge(stat string, value int64, rate float32) error {
	return nil
}





func (s *NoopClient) GaugeDelta(stat string, value int64, rate float32) error {
	return nil
}





func (s *NoopClient) Timing(stat string, delta int64, rate float32) error {
	return nil
}





func (s *NoopClient) TimingDuration(stat string, delta time.Duration, rate float32) error {
	return nil
}





func (s *NoopClient) Set(stat string, value string, rate float32) error {
	return nil
}






func (s *NoopClient) SetInt(stat string, value int64, rate float32) error {
	return nil
}






func (s *NoopClient) Raw(stat string, value string, rate float32) error {
	return nil
}


func (s *NoopClient) SetPrefix(prefix string) {}


func (s *NoopClient) NewSubStatter(prefix string) SubStatter {
	return &NoopClient{}
}


func (s *NoopClient) SetSamplerFunc(sampler SamplerFunc) {}





func NewNoopClient(a ...interface{}) (Statter, error) {
	return &NoopClient{}, nil
}


var NewNoop = NewNoopClient
