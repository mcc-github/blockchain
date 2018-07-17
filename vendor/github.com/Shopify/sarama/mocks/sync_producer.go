package mocks

import (
	"sync"

	"github.com/Shopify/sarama"
)





type SyncProducer struct {
	l            sync.Mutex
	t            ErrorReporter
	expectations []*producerExpectation
	lastOffset   int64
}





func NewSyncProducer(t ErrorReporter, config *sarama.Config) *SyncProducer {
	return &SyncProducer{
		t:            t,
		expectations: make([]*producerExpectation, 0),
	}
}











func (sp *SyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	sp.l.Lock()
	defer sp.l.Unlock()

	if len(sp.expectations) > 0 {
		expectation := sp.expectations[0]
		sp.expectations = sp.expectations[1:]
		if expectation.CheckFunction != nil {
			val, err := msg.Value.Encode()
			if err != nil {
				sp.t.Errorf("Input message encoding failed: %s", err.Error())
				return -1, -1, err
			}

			errCheck := expectation.CheckFunction(val)
			if errCheck != nil {
				sp.t.Errorf("Check function returned an error: %s", errCheck.Error())
				return -1, -1, errCheck
			}
		}
		if expectation.Result == errProduceSuccess {
			sp.lastOffset++
			msg.Offset = sp.lastOffset
			return 0, msg.Offset, nil
		}
		return -1, -1, expectation.Result
	}
	sp.t.Errorf("No more expectation set on this mock producer to handle the input message.")
	return -1, -1, errOutOfExpectations
}





func (sp *SyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	sp.l.Lock()
	defer sp.l.Unlock()

	if len(sp.expectations) >= len(msgs) {
		expectations := sp.expectations[0:len(msgs)]
		sp.expectations = sp.expectations[len(msgs):]

		for i, expectation := range expectations {
			if expectation.CheckFunction != nil {
				val, err := msgs[i].Value.Encode()
				if err != nil {
					sp.t.Errorf("Input message encoding failed: %s", err.Error())
					return err
				}
				errCheck := expectation.CheckFunction(val)
				if errCheck != nil {
					sp.t.Errorf("Check function returned an error: %s", errCheck.Error())
					return errCheck
				}
			}
			if expectation.Result != errProduceSuccess {
				return expectation.Result
			}
		}
		return nil
	}
	sp.t.Errorf("Insufficient expectations set on this mock producer to handle the input messages.")
	return errOutOfExpectations
}




func (sp *SyncProducer) Close() error {
	sp.l.Lock()
	defer sp.l.Unlock()

	if len(sp.expectations) > 0 {
		sp.t.Errorf("Expected to exhaust all expectations, but %d are left.", len(sp.expectations))
	}

	return nil
}









func (sp *SyncProducer) ExpectSendMessageWithCheckerFunctionAndSucceed(cf ValueChecker) {
	sp.l.Lock()
	defer sp.l.Unlock()
	sp.expectations = append(sp.expectations, &producerExpectation{Result: errProduceSuccess, CheckFunction: cf})
}





func (sp *SyncProducer) ExpectSendMessageWithCheckerFunctionAndFail(cf ValueChecker, err error) {
	sp.l.Lock()
	defer sp.l.Unlock()
	sp.expectations = append(sp.expectations, &producerExpectation{Result: err, CheckFunction: cf})
}




func (sp *SyncProducer) ExpectSendMessageAndSucceed() {
	sp.ExpectSendMessageWithCheckerFunctionAndSucceed(nil)
}




func (sp *SyncProducer) ExpectSendMessageAndFail(err error) {
	sp.ExpectSendMessageWithCheckerFunctionAndFail(nil, err)
}
