package sarama

import "sync"











type SyncProducer interface {

	
	
	
	SendMessage(msg *ProducerMessage) (partition int32, offset int64, err error)

	
	
	
	
	SendMessages(msgs []*ProducerMessage) error

	
	
	
	
	Close() error
}

type syncProducer struct {
	producer *asyncProducer
	wg       sync.WaitGroup
}


func NewSyncProducer(addrs []string, config *Config) (SyncProducer, error) {
	if config == nil {
		config = NewConfig()
		config.Producer.Return.Successes = true
	}

	if err := verifyProducerConfig(config); err != nil {
		return nil, err
	}

	p, err := NewAsyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}
	return newSyncProducerFromAsyncProducer(p.(*asyncProducer)), nil
}



func NewSyncProducerFromClient(client Client) (SyncProducer, error) {
	if err := verifyProducerConfig(client.Config()); err != nil {
		return nil, err
	}

	p, err := NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return newSyncProducerFromAsyncProducer(p.(*asyncProducer)), nil
}

func newSyncProducerFromAsyncProducer(p *asyncProducer) *syncProducer {
	sp := &syncProducer{producer: p}

	sp.wg.Add(2)
	go withRecover(sp.handleSuccesses)
	go withRecover(sp.handleErrors)

	return sp
}

func verifyProducerConfig(config *Config) error {
	if !config.Producer.Return.Errors {
		return ConfigurationError("Producer.Return.Errors must be true to be used in a SyncProducer")
	}
	if !config.Producer.Return.Successes {
		return ConfigurationError("Producer.Return.Successes must be true to be used in a SyncProducer")
	}
	return nil
}

func (sp *syncProducer) SendMessage(msg *ProducerMessage) (partition int32, offset int64, err error) {
	expectation := make(chan *ProducerError, 1)
	msg.expectation = expectation
	sp.producer.Input() <- msg

	if err := <-expectation; err != nil {
		return -1, -1, err.Err
	}

	return msg.Partition, msg.Offset, nil
}

func (sp *syncProducer) SendMessages(msgs []*ProducerMessage) error {
	expectations := make(chan chan *ProducerError, len(msgs))
	go func() {
		for _, msg := range msgs {
			expectation := make(chan *ProducerError, 1)
			msg.expectation = expectation
			sp.producer.Input() <- msg
			expectations <- expectation
		}
		close(expectations)
	}()

	var errors ProducerErrors
	for expectation := range expectations {
		if err := <-expectation; err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (sp *syncProducer) handleSuccesses() {
	defer sp.wg.Done()
	for msg := range sp.producer.Successes() {
		expectation := msg.expectation
		expectation <- nil
	}
}

func (sp *syncProducer) handleErrors() {
	defer sp.wg.Done()
	for err := range sp.producer.Errors() {
		expectation := err.Msg.expectation
		expectation <- err
	}
}

func (sp *syncProducer) Close() error {
	sp.producer.AsyncClose()
	sp.wg.Wait()
	return nil
}
