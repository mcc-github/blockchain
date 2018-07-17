



package statsd

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var bufPool = newBufferPool()


type StatSender interface {
	Inc(string, int64, float32) error
	Dec(string, int64, float32) error
	Gauge(string, int64, float32) error
	GaugeDelta(string, int64, float32) error
	Timing(string, int64, float32) error
	TimingDuration(string, time.Duration, float32) error
	Set(string, string, float32) error
	SetInt(string, int64, float32) error
	Raw(string, string, float32) error
}


type Statter interface {
	StatSender
	NewSubStatter(string) SubStatter
	SetPrefix(string)
	Close() error
}


type SubStatter interface {
	StatSender
	SetSamplerFunc(SamplerFunc)
	NewSubStatter(string) SubStatter
}



type SamplerFunc func(float32) bool


func DefaultSampler(rate float32) bool {
	if rate < 1 {
		if rand.Float32() < rate {
			return true
		}
		return false
	}
	return true
}


type Client struct {
	
	prefix string
	
	sender Sender
	
	sampler SamplerFunc
}


func (s *Client) Close() error {
	if s == nil {
		return nil
	}

	err := s.sender.Close()
	return err
}





func (s *Client) Inc(stat string, value int64, rate float32) error {
	if !s.includeStat(rate) {
		return nil
	}

	return s.submit(stat, "", value, "|c", rate)
}





func (s *Client) Dec(stat string, value int64, rate float32) error {
	if !s.includeStat(rate) {
		return nil
	}

	return s.submit(stat, "", -value, "|c", rate)
}





func (s *Client) Gauge(stat string, value int64, rate float32) error {
	if !s.includeStat(rate) {
		return nil
	}

	return s.submit(stat, "", value, "|g", rate)
}





func (s *Client) GaugeDelta(stat string, value int64, rate float32) error {
	if !s.includeStat(rate) {
		return nil
	}

	
	
	if value >= 0 {
		return s.submit(stat, "+", value, "|g", rate)
	}
	return s.submit(stat, "", value, "|g", rate)
}





func (s *Client) Timing(stat string, delta int64, rate float32) error {
	if !s.includeStat(rate) {
		return nil
	}

	return s.submit(stat, "", delta, "|ms", rate)
}





func (s *Client) TimingDuration(stat string, delta time.Duration, rate float32) error {
	if !s.includeStat(rate) {
		return nil
	}

	ms := float64(delta) / float64(time.Millisecond)
	return s.submit(stat, "", ms, "|ms", rate)
}





func (s *Client) Set(stat string, value string, rate float32) error {
	if !s.includeStat(rate) {
		return nil
	}

	return s.submit(stat, "", value, "|s", rate)
}





func (s *Client) SetInt(stat string, value int64, rate float32) error {
	if !s.includeStat(rate) {
		return nil
	}

	return s.submit(stat, "", value, "|s", rate)
}





func (s *Client) Raw(stat string, value string, rate float32) error {
	if !s.includeStat(rate) {
		return nil
	}

	return s.submit(stat, "", value, "", rate)
}





func (s *Client) SetSamplerFunc(sampler SamplerFunc) {
	s.sampler = sampler
}


func (s *Client) submit(stat, vprefix string, value interface{}, suffix string, rate float32) error {
	data := bufPool.Get()
	defer bufPool.Put(data)

	if s.prefix != "" {
		data.WriteString(s.prefix)
		data.WriteString(".")
	}

	data.WriteString(stat)
	data.WriteString(":")

	if vprefix != "" {
		data.WriteString(vprefix)
	}

	
	
	
	b := data.Bytes()

	switch v := value.(type) {
	case string:
		b = append(b, v...)
	case int64:
		b = strconv.AppendInt(b, v, 10)
	case float64:
		b = strconv.AppendFloat(b, v, 'f', -1, 64)
	default:
		return fmt.Errorf("No matching type format")
	}

	if suffix != "" {
		b = append(b, suffix...)
	}

	if rate < 1 {
		b = append(b, "|@"...)
		b = strconv.AppendFloat(b, float64(rate), 'f', 6, 32)
	}

	_, err := s.sender.Send(b)
	return err
}


func (s *Client) includeStat(rate float32) bool {
	if s == nil {
		return false
	}

	
	
	if s.sampler != nil {
		return s.sampler(rate)
	}
	return DefaultSampler(rate)
}



func (s *Client) SetPrefix(prefix string) {
	if s == nil {
		return
	}

	s.prefix = prefix
}


func (s *Client) NewSubStatter(prefix string) SubStatter {
	var c *Client
	if s != nil {
		c = &Client{
			prefix:  joinPathComp(s.prefix, prefix),
			sender:  s.sender,
			sampler: s.sampler,
		}
	}
	return c
}







func NewClient(addr, prefix string) (Statter, error) {
	sender, err := NewSimpleSender(addr)
	if err != nil {
		return nil, err
	}

	return &Client{prefix: prefix, sender: sender}, nil
}






func NewClientWithSender(sender Sender, prefix string) (Statter, error) {
	if sender == nil {
		return nil, fmt.Errorf("Client sender may not be nil")
	}

	return &Client{prefix: prefix, sender: sender}, nil
}






func joinPathComp(prefix, suffix string) string {
	suffix = strings.TrimLeft(suffix, ".")
	if prefix != "" && suffix != "" {
		return prefix + "." + suffix
	}
	return prefix + suffix
}


var Dial = NewClient


var New = NewClient
