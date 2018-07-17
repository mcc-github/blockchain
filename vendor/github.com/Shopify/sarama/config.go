package sarama

import (
	"crypto/tls"
	"regexp"
	"time"

	"github.com/rcrowley/go-metrics"
)

const defaultClientID = "sarama"

var validID = regexp.MustCompile(`\A[A-Za-z0-9._-]+\z`)


type Config struct {
	
	
	Net struct {
		
		
		MaxOpenRequests int

		
		
		
		DialTimeout  time.Duration 
		ReadTimeout  time.Duration 
		WriteTimeout time.Duration 

		TLS struct {
			
			
			Enable bool
			
			
			Config *tls.Config
		}

		
		
		SASL struct {
			
			
			Enable bool
			
			
			
			Handshake bool
			
			User     string
			Password string
		}

		
		
		KeepAlive time.Duration
	}

	
	
	Metadata struct {
		Retry struct {
			
			
			Max int
			
			
			Backoff time.Duration
		}
		
		
		
		RefreshFrequency time.Duration

		
		
		
		
		Full bool
	}

	
	
	Producer struct {
		
		
		MaxMessageBytes int
		
		
		
		RequiredAcks RequiredAcks
		
		
		
		
		
		Timeout time.Duration
		
		
		Compression CompressionCodec
		
		
		
		Partitioner PartitionerConstructor

		
		
		
		
		
		Return struct {
			
			
			Successes bool

			
			
			Errors bool
		}

		
		
		
		
		Flush struct {
			
			
			Bytes int
			
			
			Messages int
			
			
			Frequency time.Duration
			
			
			
			MaxMessages int
		}

		Retry struct {
			
			
			Max int
			
			
			
			Backoff time.Duration
		}
	}

	
	
	
	
	
	
	
	
	
	Consumer struct {
		Retry struct {
			
			
			Backoff time.Duration
		}

		
		
		Fetch struct {
			
			
			
			
			Min int32
			
			
			
			
			
			Default int32
			
			
			
			
			
			
			Max int32
		}
		
		
		
		
		
		
		MaxWaitTime time.Duration

		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		MaxProcessingTime time.Duration

		
		
		Return struct {
			
			
			Errors bool
		}

		
		
		
		Offsets struct {
			
			CommitInterval time.Duration

			
			
			Initial int64

			
			
			
			
			
			
			Retention time.Duration
		}
	}

	
	
	
	ClientID string
	
	
	
	
	ChannelBufferSize int
	
	
	
	
	
	
	Version KafkaVersion
	
	
	
	
	
	MetricRegistry metrics.Registry
}


func NewConfig() *Config {
	c := &Config{}

	c.Net.MaxOpenRequests = 5
	c.Net.DialTimeout = 30 * time.Second
	c.Net.ReadTimeout = 30 * time.Second
	c.Net.WriteTimeout = 30 * time.Second
	c.Net.SASL.Handshake = true

	c.Metadata.Retry.Max = 3
	c.Metadata.Retry.Backoff = 250 * time.Millisecond
	c.Metadata.RefreshFrequency = 10 * time.Minute
	c.Metadata.Full = true

	c.Producer.MaxMessageBytes = 1000000
	c.Producer.RequiredAcks = WaitForLocal
	c.Producer.Timeout = 10 * time.Second
	c.Producer.Partitioner = NewHashPartitioner
	c.Producer.Retry.Max = 3
	c.Producer.Retry.Backoff = 100 * time.Millisecond
	c.Producer.Return.Errors = true

	c.Consumer.Fetch.Min = 1
	c.Consumer.Fetch.Default = 1024 * 1024
	c.Consumer.Retry.Backoff = 2 * time.Second
	c.Consumer.MaxWaitTime = 250 * time.Millisecond
	c.Consumer.MaxProcessingTime = 100 * time.Millisecond
	c.Consumer.Return.Errors = false
	c.Consumer.Offsets.CommitInterval = 1 * time.Second
	c.Consumer.Offsets.Initial = OffsetNewest

	c.ClientID = defaultClientID
	c.ChannelBufferSize = 256
	c.Version = minVersion
	c.MetricRegistry = metrics.NewRegistry()

	return c
}



func (c *Config) Validate() error {
	
	if c.Net.TLS.Enable == false && c.Net.TLS.Config != nil {
		Logger.Println("Net.TLS is disabled but a non-nil configuration was provided.")
	}
	if c.Net.SASL.Enable == false {
		if c.Net.SASL.User != "" {
			Logger.Println("Net.SASL is disabled but a non-empty username was provided.")
		}
		if c.Net.SASL.Password != "" {
			Logger.Println("Net.SASL is disabled but a non-empty password was provided.")
		}
	}
	if c.Producer.RequiredAcks > 1 {
		Logger.Println("Producer.RequiredAcks > 1 is deprecated and will raise an exception with kafka >= 0.8.2.0.")
	}
	if c.Producer.MaxMessageBytes >= int(MaxRequestSize) {
		Logger.Println("Producer.MaxMessageBytes must be smaller than MaxRequestSize; it will be ignored.")
	}
	if c.Producer.Flush.Bytes >= int(MaxRequestSize) {
		Logger.Println("Producer.Flush.Bytes must be smaller than MaxRequestSize; it will be ignored.")
	}
	if (c.Producer.Flush.Bytes > 0 || c.Producer.Flush.Messages > 0) && c.Producer.Flush.Frequency == 0 {
		Logger.Println("Producer.Flush: Bytes or Messages are set, but Frequency is not; messages may not get flushed.")
	}
	if c.Producer.Timeout%time.Millisecond != 0 {
		Logger.Println("Producer.Timeout only supports millisecond resolution; nanoseconds will be truncated.")
	}
	if c.Consumer.MaxWaitTime < 100*time.Millisecond {
		Logger.Println("Consumer.MaxWaitTime is very low, which can cause high CPU and network usage. See documentation for details.")
	}
	if c.Consumer.MaxWaitTime%time.Millisecond != 0 {
		Logger.Println("Consumer.MaxWaitTime only supports millisecond precision; nanoseconds will be truncated.")
	}
	if c.Consumer.Offsets.Retention%time.Millisecond != 0 {
		Logger.Println("Consumer.Offsets.Retention only supports millisecond precision; nanoseconds will be truncated.")
	}
	if c.ClientID == defaultClientID {
		Logger.Println("ClientID is the default of 'sarama', you should consider setting it to something application-specific.")
	}

	
	switch {
	case c.Net.MaxOpenRequests <= 0:
		return ConfigurationError("Net.MaxOpenRequests must be > 0")
	case c.Net.DialTimeout <= 0:
		return ConfigurationError("Net.DialTimeout must be > 0")
	case c.Net.ReadTimeout <= 0:
		return ConfigurationError("Net.ReadTimeout must be > 0")
	case c.Net.WriteTimeout <= 0:
		return ConfigurationError("Net.WriteTimeout must be > 0")
	case c.Net.KeepAlive < 0:
		return ConfigurationError("Net.KeepAlive must be >= 0")
	case c.Net.SASL.Enable == true && c.Net.SASL.User == "":
		return ConfigurationError("Net.SASL.User must not be empty when SASL is enabled")
	case c.Net.SASL.Enable == true && c.Net.SASL.Password == "":
		return ConfigurationError("Net.SASL.Password must not be empty when SASL is enabled")
	}

	
	switch {
	case c.Metadata.Retry.Max < 0:
		return ConfigurationError("Metadata.Retry.Max must be >= 0")
	case c.Metadata.Retry.Backoff < 0:
		return ConfigurationError("Metadata.Retry.Backoff must be >= 0")
	case c.Metadata.RefreshFrequency < 0:
		return ConfigurationError("Metadata.RefreshFrequency must be >= 0")
	}

	
	switch {
	case c.Producer.MaxMessageBytes <= 0:
		return ConfigurationError("Producer.MaxMessageBytes must be > 0")
	case c.Producer.RequiredAcks < -1:
		return ConfigurationError("Producer.RequiredAcks must be >= -1")
	case c.Producer.Timeout <= 0:
		return ConfigurationError("Producer.Timeout must be > 0")
	case c.Producer.Partitioner == nil:
		return ConfigurationError("Producer.Partitioner must not be nil")
	case c.Producer.Flush.Bytes < 0:
		return ConfigurationError("Producer.Flush.Bytes must be >= 0")
	case c.Producer.Flush.Messages < 0:
		return ConfigurationError("Producer.Flush.Messages must be >= 0")
	case c.Producer.Flush.Frequency < 0:
		return ConfigurationError("Producer.Flush.Frequency must be >= 0")
	case c.Producer.Flush.MaxMessages < 0:
		return ConfigurationError("Producer.Flush.MaxMessages must be >= 0")
	case c.Producer.Flush.MaxMessages > 0 && c.Producer.Flush.MaxMessages < c.Producer.Flush.Messages:
		return ConfigurationError("Producer.Flush.MaxMessages must be >= Producer.Flush.Messages when set")
	case c.Producer.Retry.Max < 0:
		return ConfigurationError("Producer.Retry.Max must be >= 0")
	case c.Producer.Retry.Backoff < 0:
		return ConfigurationError("Producer.Retry.Backoff must be >= 0")
	}

	if c.Producer.Compression == CompressionLZ4 && !c.Version.IsAtLeast(V0_10_0_0) {
		return ConfigurationError("lz4 compression requires Version >= V0_10_0_0")
	}

	
	switch {
	case c.Consumer.Fetch.Min <= 0:
		return ConfigurationError("Consumer.Fetch.Min must be > 0")
	case c.Consumer.Fetch.Default <= 0:
		return ConfigurationError("Consumer.Fetch.Default must be > 0")
	case c.Consumer.Fetch.Max < 0:
		return ConfigurationError("Consumer.Fetch.Max must be >= 0")
	case c.Consumer.MaxWaitTime < 1*time.Millisecond:
		return ConfigurationError("Consumer.MaxWaitTime must be >= 1ms")
	case c.Consumer.MaxProcessingTime <= 0:
		return ConfigurationError("Consumer.MaxProcessingTime must be > 0")
	case c.Consumer.Retry.Backoff < 0:
		return ConfigurationError("Consumer.Retry.Backoff must be >= 0")
	case c.Consumer.Offsets.CommitInterval <= 0:
		return ConfigurationError("Consumer.Offsets.CommitInterval must be > 0")
	case c.Consumer.Offsets.Initial != OffsetOldest && c.Consumer.Offsets.Initial != OffsetNewest:
		return ConfigurationError("Consumer.Offsets.Initial must be OffsetOldest or OffsetNewest")

	}

	
	switch {
	case c.ChannelBufferSize < 0:
		return ConfigurationError("ChannelBufferSize must be >= 0")
	case !validID.MatchString(c.ClientID):
		return ConfigurationError("ClientID is invalid")
	}

	return nil
}
