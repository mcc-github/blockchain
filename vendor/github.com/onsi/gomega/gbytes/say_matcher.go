package gbytes

import (
	"fmt"
	"regexp"

	"github.com/onsi/gomega/format"
)


type BufferProvider interface {
	Buffer() *Buffer
}


func Say(expected string, args ...interface{}) *sayMatcher {
	formattedRegexp := expected
	if len(args) > 0 {
		formattedRegexp = fmt.Sprintf(expected, args...)
	}
	return &sayMatcher{
		re: regexp.MustCompile(formattedRegexp),
	}
}

type sayMatcher struct {
	re              *regexp.Regexp
	receivedSayings []byte
}

func (m *sayMatcher) buffer(actual interface{}) (*Buffer, bool) {
	var buffer *Buffer

	switch x := actual.(type) {
	case *Buffer:
		buffer = x
	case BufferProvider:
		buffer = x.Buffer()
	default:
		return nil, false
	}

	return buffer, true
}

func (m *sayMatcher) Match(actual interface{}) (success bool, err error) {
	buffer, ok := m.buffer(actual)
	if !ok {
		return false, fmt.Errorf("Say must be passed a *gbytes.Buffer or BufferProvider.  Got:\n%s", format.Object(actual, 1))
	}

	didSay, sayings := buffer.didSay(m.re)
	m.receivedSayings = sayings

	return didSay, nil
}

func (m *sayMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"Got stuck at:\n%s\nWaiting for:\n%s",
		format.IndentString(string(m.receivedSayings), 1),
		format.IndentString(m.re.String(), 1),
	)
}

func (m *sayMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"Saw:\n%s\nWhich matches the unexpected:\n%s",
		format.IndentString(string(m.receivedSayings), 1),
		format.IndentString(m.re.String(), 1),
	)
}

func (m *sayMatcher) MatchMayChangeInTheFuture(actual interface{}) bool {
	switch x := actual.(type) {
	case *Buffer:
		return !x.Closed()
	case BufferProvider:
		return !x.Buffer().Closed()
	default:
		return true
	}
}
