
package gomega

import (
	"fmt"
	"reflect"
	"time"

	"github.com/onsi/gomega/internal/assertion"
	"github.com/onsi/gomega/internal/asyncassertion"
	"github.com/onsi/gomega/internal/testingtsupport"
	"github.com/onsi/gomega/types"
)

const GOMEGA_VERSION = "1.5.0"

const nilFailHandlerPanic = `You are trying to make an assertion, but Gomega's fail handler is nil.
If you're using Ginkgo then you probably forgot to put your assertion in an It().
Alternatively, you may have forgotten to register a fail handler with RegisterFailHandler() or RegisterTestingT().
Depending on your vendoring solution you may be inadvertently importing gomega and subpackages (e.g. ghhtp, gexec,...) from different locations.
`

var globalFailWrapper *types.GomegaFailWrapper

var defaultEventuallyTimeout = time.Second
var defaultEventuallyPollingInterval = 10 * time.Millisecond
var defaultConsistentlyDuration = 100 * time.Millisecond
var defaultConsistentlyPollingInterval = 10 * time.Millisecond



func RegisterFailHandler(handler types.GomegaFailHandler) {
	RegisterFailHandlerWithT(testingtsupport.EmptyTWithHelper{}, handler)
}



func RegisterFailHandlerWithT(t types.TWithHelper, handler types.GomegaFailHandler) {
	if handler == nil {
		globalFailWrapper = nil
		return
	}

	globalFailWrapper = &types.GomegaFailWrapper{
		Fail:        handler,
		TWithHelper: t,
	}
}






















func RegisterTestingT(t types.GomegaTestingT) {
	tWithHelper, hasHelper := t.(types.TWithHelper)
	if !hasHelper {
		RegisterFailHandler(testingtsupport.BuildTestingTGomegaFailWrapper(t).Fail)
		return
	}
	RegisterFailHandlerWithT(tWithHelper, testingtsupport.BuildTestingTGomegaFailWrapper(t).Fail)
}










func InterceptGomegaFailures(f func()) []string {
	originalHandler := globalFailWrapper.Fail
	failures := []string{}
	RegisterFailHandler(func(message string, callerSkip ...int) {
		failures = append(failures, message)
	})
	f()
	RegisterFailHandler(originalHandler)
	return failures
}


















func Ω(actual interface{}, extra ...interface{}) Assertion {
	return ExpectWithOffset(0, actual, extra...)
}


















func Expect(actual interface{}, extra ...interface{}) Assertion {
	return ExpectWithOffset(0, actual, extra...)
}










func ExpectWithOffset(offset int, actual interface{}, extra ...interface{}) Assertion {
	if globalFailWrapper == nil {
		panic(nilFailHandlerPanic)
	}
	return assertion.New(actual, globalFailWrapper, offset, extra...)
}





































func Eventually(actual interface{}, intervals ...interface{}) AsyncAssertion {
	return EventuallyWithOffset(0, actual, intervals...)
}




func EventuallyWithOffset(offset int, actual interface{}, intervals ...interface{}) AsyncAssertion {
	if globalFailWrapper == nil {
		panic(nilFailHandlerPanic)
	}
	timeoutInterval := defaultEventuallyTimeout
	pollingInterval := defaultEventuallyPollingInterval
	if len(intervals) > 0 {
		timeoutInterval = toDuration(intervals[0])
	}
	if len(intervals) > 1 {
		pollingInterval = toDuration(intervals[1])
	}
	return asyncassertion.New(asyncassertion.AsyncAssertionTypeEventually, actual, globalFailWrapper, timeoutInterval, pollingInterval, offset)
}
























func Consistently(actual interface{}, intervals ...interface{}) AsyncAssertion {
	return ConsistentlyWithOffset(0, actual, intervals...)
}




func ConsistentlyWithOffset(offset int, actual interface{}, intervals ...interface{}) AsyncAssertion {
	if globalFailWrapper == nil {
		panic(nilFailHandlerPanic)
	}
	timeoutInterval := defaultConsistentlyDuration
	pollingInterval := defaultConsistentlyPollingInterval
	if len(intervals) > 0 {
		timeoutInterval = toDuration(intervals[0])
	}
	if len(intervals) > 1 {
		pollingInterval = toDuration(intervals[1])
	}
	return asyncassertion.New(asyncassertion.AsyncAssertionTypeConsistently, actual, globalFailWrapper, timeoutInterval, pollingInterval, offset)
}


func SetDefaultEventuallyTimeout(t time.Duration) {
	defaultEventuallyTimeout = t
}


func SetDefaultEventuallyPollingInterval(t time.Duration) {
	defaultEventuallyPollingInterval = t
}


func SetDefaultConsistentlyDuration(t time.Duration) {
	defaultConsistentlyDuration = t
}


func SetDefaultConsistentlyPollingInterval(t time.Duration) {
	defaultConsistentlyPollingInterval = t
}














type AsyncAssertion interface {
	Should(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool
	ShouldNot(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool
}


type GomegaAsyncAssertion = AsyncAssertion















type Assertion interface {
	Should(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool
	ShouldNot(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool

	To(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool
	ToNot(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool
	NotTo(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool
}


type GomegaAssertion = Assertion


type OmegaMatcher types.GomegaMatcher





type WithT struct {
	t types.GomegaTestingT
}


type GomegaWithT = WithT










func NewWithT(t types.GomegaTestingT) *WithT {
	return &WithT{
		t: t,
	}
}


func NewGomegaWithT(t types.GomegaTestingT) *GomegaWithT {
	return NewWithT(t)
}


func (g *WithT) Expect(actual interface{}, extra ...interface{}) Assertion {
	return assertion.New(actual, testingtsupport.BuildTestingTGomegaFailWrapper(g.t), 0, extra...)
}


func (g *WithT) Eventually(actual interface{}, intervals ...interface{}) AsyncAssertion {
	timeoutInterval := defaultEventuallyTimeout
	pollingInterval := defaultEventuallyPollingInterval
	if len(intervals) > 0 {
		timeoutInterval = toDuration(intervals[0])
	}
	if len(intervals) > 1 {
		pollingInterval = toDuration(intervals[1])
	}
	return asyncassertion.New(asyncassertion.AsyncAssertionTypeEventually, actual, testingtsupport.BuildTestingTGomegaFailWrapper(g.t), timeoutInterval, pollingInterval, 0)
}


func (g *WithT) Consistently(actual interface{}, intervals ...interface{}) AsyncAssertion {
	timeoutInterval := defaultConsistentlyDuration
	pollingInterval := defaultConsistentlyPollingInterval
	if len(intervals) > 0 {
		timeoutInterval = toDuration(intervals[0])
	}
	if len(intervals) > 1 {
		pollingInterval = toDuration(intervals[1])
	}
	return asyncassertion.New(asyncassertion.AsyncAssertionTypeConsistently, actual, testingtsupport.BuildTestingTGomegaFailWrapper(g.t), timeoutInterval, pollingInterval, 0)
}

func toDuration(input interface{}) time.Duration {
	duration, ok := input.(time.Duration)
	if ok {
		return duration
	}

	value := reflect.ValueOf(input)
	kind := reflect.TypeOf(input).Kind()

	if reflect.Int <= kind && kind <= reflect.Int64 {
		return time.Duration(value.Int()) * time.Second
	} else if reflect.Uint <= kind && kind <= reflect.Uint64 {
		return time.Duration(value.Uint()) * time.Second
	} else if reflect.Float32 <= kind && kind <= reflect.Float64 {
		return time.Duration(value.Float() * float64(time.Second))
	} else if reflect.String == kind {
		duration, err := time.ParseDuration(value.String())
		if err != nil {
			panic(fmt.Sprintf("%#v is not a valid parsable duration string.", input))
		}
		return duration
	}

	panic(fmt.Sprintf("%v is not a valid interval.  Must be time.Duration, parsable duration string or a number.", input))
}
