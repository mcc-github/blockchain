



package statsd

import "time"





















func NewBufferedClient(addr, prefix string, flushInterval time.Duration, flushBytes int) (Statter, error) {
	if flushBytes <= 0 {
		
		flushBytes = 1432
	}
	if flushInterval <= time.Duration(0) {
		flushInterval = 300 * time.Millisecond
	}
	sender, err := NewBufferedSender(addr, flushInterval, flushBytes)
	if err != nil {
		return nil, err
	}

	client := &Client{
		prefix: prefix,
		sender: sender,
	}

	return client, nil
}
