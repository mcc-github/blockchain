package mocks

import (
	"sync"

	"github.com/Shopify/sarama"
)







type AsyncProducer struct {
	l            sync.Mutex
	t            ErrorReporter
	expectations []*producerExpectation
	closed       chan struct{}
	input        chan *sarama.ProducerMessage
	successes    chan *sarama.ProducerMessage
	errors       chan *sarama.ProducerError
	lastOffset   int64
}





func NewAsyncProducer(t ErrorReporter, config *sarama.Config) *AsyncProducer {
	if config == nil {
		config = sarama.NewConfig()
	}
	mp := &AsyncProducer{
		t:            t,
		closed:       make(chan struct{}, 0),
		expectations: make([]*producerExpectation, 0),
		input:        make(chan *sarama.ProducerMessage, config.ChannelBufferSize),
		successes:    make(chan *sarama.ProducerMessage, config.ChannelBufferSize),
		errors:       make(chan *sarama.ProducerError, config.ChannelBufferSize),
	}

	go func() {
		defer func() {
			close(mp.successes)
			close(mp.errors)
			close(mp.closed)
		}()

		for msg := range mp.input {
			mp.l.Lock()
			if mp.expectations == nil || len(mp.expectations) == 0 {
				mp.expectations = nil
				mp.t.Errorf("No more expectation set on this mock producer to handle the input message.")
			} else {
				expectation := mp.expectations[0]
				mp.expectations = mp.expectations[1:]
				if expectation.CheckFunction != nil {
					if val, err := msg.Value.Encode(); err != nil {
						mp.t.Errorf("Input message encoding failed: %s", err.Error())
						mp.errors <- &sarama.ProducerError{Err: err, Msg: msg}
					} else {
						err = expectation.CheckFunction(val)
						if err != nil {
							mp.t.Errorf("Check function returned an error: %s", err.Error())
							mp.errors <- &sarama.ProducerError{Err: err, Msg: msg}
						}
					}
				}
				if expectation.Result == errProduceSuccess {
					mp.lastOffset++
					if config.Producer.Return.Successes {
						msg.Offset = mp.lastOffset
						mp.successes <- msg
					}
				} else {
					if config.Producer.Return.Errors {
						mp.errors <- &sarama.ProducerError{Err: expectation.Result, Msg: msg}
					}
				}
			}
			mp.l.Unlock()
		}

		mp.l.Lock()
		if len(mp.expectations) > 0 {
			mp.t.Errorf("Expected to exhaust all expectations, but %d are left.", len(mp.expectations))
		}
		mp.l.Unlock()
	}()

	return mp
}








func (mp *AsyncProducer) AsyncClose() {
	close(mp.input)
}




func (mp *AsyncProducer) Close() error {
	mp.AsyncClose()
	<-mp.closed
	return nil
}






func (mp *AsyncProducer) Input() chan<- *sarama.ProducerMessage {
	return mp.input
}


func (mp *AsyncProducer) Successes() <-chan *sarama.ProducerMessage {
	return mp.successes
}


func (mp *AsyncProducer) Errors() <-chan *sarama.ProducerError {
	return mp.errors
}










func (mp *AsyncProducer) ExpectInputWithCheckerFunctionAndSucceed(cf ValueChecker) {
	mp.l.Lock()
	defer mp.l.Unlock()
	mp.expectations = append(mp.expectations, &producerExpectation{Result: errProduceSuccess, CheckFunction: cf})
}






func (mp *AsyncProducer) ExpectInputWithCheckerFunctionAndFail(cf ValueChecker, err error) {
	mp.l.Lock()
	defer mp.l.Unlock()
	mp.expectations = append(mp.expectations, &producerExpectation{Result: err, CheckFunction: cf})
}





func (mp *AsyncProducer) ExpectInputAndSucceed() {
	mp.ExpectInputWithCheckerFunctionAndSucceed(nil)
}




func (mp *AsyncProducer) ExpectInputAndFail(err error) {
	mp.ExpectInputWithCheckerFunctionAndFail(nil, err)
}
