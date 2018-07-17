
package mocks

import (
	"errors"

	"github.com/Shopify/sarama"
)



type ErrorReporter interface {
	Errorf(string, ...interface{})
}



type ValueChecker func(val []byte) error

var (
	errProduceSuccess              error = nil
	errOutOfExpectations                 = errors.New("No more expectations set on mock")
	errPartitionConsumerNotStarted       = errors.New("The partition consumer was never started")
)

const AnyOffset int64 = -1000

type producerExpectation struct {
	Result        error
	CheckFunction ValueChecker
}

type consumerExpectation struct {
	Err error
	Msg *sarama.ConsumerMessage
}
