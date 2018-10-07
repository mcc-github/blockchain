package mock

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
)


type TestingT interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
}





type Call struct {
	Parent *Mock

	
	Method string

	
	Arguments Arguments

	
	
	ReturnArguments Arguments

	
	callerInfo []string

	
	
	Repeatability int

	
	totalCalls int

	
	optional bool

	
	
	WaitFor <-chan time.Time

	waitTime time.Duration

	
	
	
	RunFn func(Arguments)
}

func newCall(parent *Mock, methodName string, callerInfo []string, methodArguments ...interface{}) *Call {
	return &Call{
		Parent:          parent,
		Method:          methodName,
		Arguments:       methodArguments,
		ReturnArguments: make([]interface{}, 0),
		callerInfo:      callerInfo,
		Repeatability:   0,
		WaitFor:         nil,
		RunFn:           nil,
	}
}

func (c *Call) lock() {
	c.Parent.mutex.Lock()
}

func (c *Call) unlock() {
	c.Parent.mutex.Unlock()
}




func (c *Call) Return(returnArguments ...interface{}) *Call {
	c.lock()
	defer c.unlock()

	c.ReturnArguments = returnArguments

	return c
}




func (c *Call) Once() *Call {
	return c.Times(1)
}




func (c *Call) Twice() *Call {
	return c.Times(2)
}





func (c *Call) Times(i int) *Call {
	c.lock()
	defer c.unlock()
	c.Repeatability = i
	return c
}





func (c *Call) WaitUntil(w <-chan time.Time) *Call {
	c.lock()
	defer c.unlock()
	c.WaitFor = w
	return c
}




func (c *Call) After(d time.Duration) *Call {
	c.lock()
	defer c.unlock()
	c.waitTime = d
	return c
}









func (c *Call) Run(fn func(args Arguments)) *Call {
	c.lock()
	defer c.unlock()
	c.RunFn = fn
	return c
}



func (c *Call) Maybe() *Call {
	c.lock()
	defer c.unlock()
	c.optional = true
	return c
}







func (c *Call) On(methodName string, arguments ...interface{}) *Call {
	return c.Parent.On(methodName, arguments...)
}




type Mock struct {
	
	
	ExpectedCalls []*Call

	
	Calls []Call

	
	
	test TestingT

	
	
	testData objx.Map

	mutex sync.Mutex
}



func (m *Mock) TestData() objx.Map {

	if m.testData == nil {
		m.testData = make(objx.Map)
	}

	return m.testData
}




func (m *Mock) Test(t TestingT) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.test = t
}




func (m *Mock) fail(format string, args ...interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.test == nil {
		panic(fmt.Sprintf(format, args...))
	}
	m.test.Errorf(format, args...)
	m.test.FailNow()
}





func (m *Mock) On(methodName string, arguments ...interface{}) *Call {
	for _, arg := range arguments {
		if v := reflect.ValueOf(arg); v.Kind() == reflect.Func {
			panic(fmt.Sprintf("cannot use Func in expectations. Use mock.AnythingOfType(\"%T\")", arg))
		}
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	c := newCall(m, methodName, assert.CallerInfo(), arguments...)
	m.ExpectedCalls = append(m.ExpectedCalls, c)
	return c
}





func (m *Mock) findExpectedCall(method string, arguments ...interface{}) (int, *Call) {
	for i, call := range m.ExpectedCalls {
		if call.Method == method && call.Repeatability > -1 {

			_, diffCount := call.Arguments.Diff(arguments)
			if diffCount == 0 {
				return i, call
			}

		}
	}
	return -1, nil
}

func (m *Mock) findClosestCall(method string, arguments ...interface{}) (*Call, string) {
	var diffCount int
	var closestCall *Call
	var err string

	for _, call := range m.expectedCalls() {
		if call.Method == method {

			errInfo, tempDiffCount := call.Arguments.Diff(arguments)
			if tempDiffCount < diffCount || diffCount == 0 {
				diffCount = tempDiffCount
				closestCall = call
				err = errInfo
			}

		}
	}

	return closestCall, err
}

func callString(method string, arguments Arguments, includeArgumentValues bool) string {

	var argValsString string
	if includeArgumentValues {
		var argVals []string
		for argIndex, arg := range arguments {
			argVals = append(argVals, fmt.Sprintf("%d: %#v", argIndex, arg))
		}
		argValsString = fmt.Sprintf("\n\t\t%s", strings.Join(argVals, "\n\t\t"))
	}

	return fmt.Sprintf("%s(%s)%s", method, arguments.String(), argValsString)
}





func (m *Mock) Called(arguments ...interface{}) Arguments {
	
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Couldn't get the caller information")
	}
	functionPath := runtime.FuncForPC(pc).Name()
	
	
	
	
	re := regexp.MustCompile("\\.pN\\d+_")
	if re.MatchString(functionPath) {
		functionPath = re.Split(functionPath, -1)[0]
	}
	parts := strings.Split(functionPath, ".")
	functionName := parts[len(parts)-1]
	return m.MethodCalled(functionName, arguments...)
}





func (m *Mock) MethodCalled(methodName string, arguments ...interface{}) Arguments {
	m.mutex.Lock()
	
	found, call := m.findExpectedCall(methodName, arguments...)

	if found < 0 {
		
		
		
		
		
		

		closestCall, mismatch := m.findClosestCall(methodName, arguments...)
		m.mutex.Unlock()

		if closestCall != nil {
			m.fail("\n\nmock: Unexpected Method Call\n-----------------------------\n\n%s\n\nThe closest call I have is: \n\n%s\n\n%s\nDiff: %s",
				callString(methodName, arguments, true),
				callString(methodName, closestCall.Arguments, true),
				diffArguments(closestCall.Arguments, arguments),
				strings.TrimSpace(mismatch),
			)
		} else {
			m.fail("\nassert: mock: I don't know what to return because the method call was unexpected.\n\tEither do Mock.On(\"%s\").Return(...) first, or remove the %s() call.\n\tThis method was unexpected:\n\t\t%s\n\tat: %s", methodName, methodName, callString(methodName, arguments, true), assert.CallerInfo())
		}
	}

	if call.Repeatability == 1 {
		call.Repeatability = -1
	} else if call.Repeatability > 1 {
		call.Repeatability--
	}
	call.totalCalls++

	
	m.Calls = append(m.Calls, *newCall(m, methodName, assert.CallerInfo(), arguments...))
	m.mutex.Unlock()

	
	if call.WaitFor != nil {
		<-call.WaitFor
	} else {
		time.Sleep(call.waitTime)
	}

	m.mutex.Lock()
	runFn := call.RunFn
	m.mutex.Unlock()

	if runFn != nil {
		runFn(arguments)
	}

	m.mutex.Lock()
	returnArgs := call.ReturnArguments
	m.mutex.Unlock()

	return returnArgs
}



type assertExpectationser interface {
	AssertExpectations(TestingT) bool
}





func AssertExpectationsForObjects(t TestingT, testObjects ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	for _, obj := range testObjects {
		if m, ok := obj.(Mock); ok {
			t.Logf("Deprecated mock.AssertExpectationsForObjects(myMock.Mock) use mock.AssertExpectationsForObjects(myMock)")
			obj = &m
		}
		m := obj.(assertExpectationser)
		if !m.AssertExpectations(t) {
			t.Logf("Expectations didn't match for Mock: %+v", reflect.TypeOf(m))
			return false
		}
	}
	return true
}



func (m *Mock) AssertExpectations(t TestingT) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var somethingMissing bool
	var failedExpectations int

	
	expectedCalls := m.expectedCalls()
	for _, expectedCall := range expectedCalls {
		if !expectedCall.optional && !m.methodWasCalled(expectedCall.Method, expectedCall.Arguments) && expectedCall.totalCalls == 0 {
			somethingMissing = true
			failedExpectations++
			t.Logf("FAIL:\t%s(%s)\n\t\tat: %s", expectedCall.Method, expectedCall.Arguments.String(), expectedCall.callerInfo)
		} else {
			if expectedCall.Repeatability > 0 {
				somethingMissing = true
				failedExpectations++
				t.Logf("FAIL:\t%s(%s)\n\t\tat: %s", expectedCall.Method, expectedCall.Arguments.String(), expectedCall.callerInfo)
			} else {
				t.Logf("PASS:\t%s(%s)", expectedCall.Method, expectedCall.Arguments.String())
			}
		}
	}

	if somethingMissing {
		t.Errorf("FAIL: %d out of %d expectation(s) were met.\n\tThe code you are testing needs to make %d more call(s).\n\tat: %s", len(expectedCalls)-failedExpectations, len(expectedCalls), failedExpectations, assert.CallerInfo())
	}

	return !somethingMissing
}


func (m *Mock) AssertNumberOfCalls(t TestingT, methodName string, expectedCalls int) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var actualCalls int
	for _, call := range m.calls() {
		if call.Method == methodName {
			actualCalls++
		}
	}
	return assert.Equal(t, expectedCalls, actualCalls, fmt.Sprintf("Expected number of calls (%d) does not match the actual number of calls (%d).", expectedCalls, actualCalls))
}



func (m *Mock) AssertCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if !m.methodWasCalled(methodName, arguments) {
		var calledWithArgs []string
		for _, call := range m.calls() {
			calledWithArgs = append(calledWithArgs, fmt.Sprintf("%v", call.Arguments))
		}
		if len(calledWithArgs) == 0 {
			return assert.Fail(t, "Should have called with given arguments",
				fmt.Sprintf("Expected %q to have been called with:\n%v\nbut no actual calls happened", methodName, arguments))
		}
		return assert.Fail(t, "Should have called with given arguments",
			fmt.Sprintf("Expected %q to have been called with:\n%v\nbut actual calls were:\n        %v", methodName, arguments, strings.Join(calledWithArgs, "\n")))
	}
	return true
}



func (m *Mock) AssertNotCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.methodWasCalled(methodName, arguments) {
		return assert.Fail(t, "Should not have called with given arguments",
			fmt.Sprintf("Expected %q to not have been called with:\n%v\nbut actually it was.", methodName, arguments))
	}
	return true
}

func (m *Mock) methodWasCalled(methodName string, expected []interface{}) bool {
	for _, call := range m.calls() {
		if call.Method == methodName {

			_, differences := Arguments(expected).Diff(call.Arguments)

			if differences == 0 {
				
				return true
			}

		}
	}
	
	return false
}

func (m *Mock) expectedCalls() []*Call {
	return append([]*Call{}, m.ExpectedCalls...)
}

func (m *Mock) calls() []Call {
	return append([]Call{}, m.Calls...)
}




type Arguments []interface{}

const (
	
	
	Anything = "mock.Anything"
)



type AnythingOfTypeArgument string






func AnythingOfType(t string) AnythingOfTypeArgument {
	return AnythingOfTypeArgument(t)
}



type argumentMatcher struct {
	
	fn reflect.Value
}

func (f argumentMatcher) Matches(argument interface{}) bool {
	expectType := f.fn.Type().In(0)
	expectTypeNilSupported := false
	switch expectType.Kind() {
	case reflect.Interface, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Ptr:
		expectTypeNilSupported = true
	}

	argType := reflect.TypeOf(argument)
	var arg reflect.Value
	if argType == nil {
		arg = reflect.New(expectType).Elem()
	} else {
		arg = reflect.ValueOf(argument)
	}

	if argType == nil && !expectTypeNilSupported {
		panic(errors.New("attempting to call matcher with nil for non-nil expected type"))
	}
	if argType == nil || argType.AssignableTo(expectType) {
		result := f.fn.Call([]reflect.Value{arg})
		return result[0].Bool()
	}
	return false
}

func (f argumentMatcher) String() string {
	return fmt.Sprintf("func(%s) bool", f.fn.Type().In(0).Name())
}












func MatchedBy(fn interface{}) argumentMatcher {
	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("assert: arguments: %s is not a func", fn))
	}
	if fnType.NumIn() != 1 {
		panic(fmt.Sprintf("assert: arguments: %s does not take exactly one argument", fn))
	}
	if fnType.NumOut() != 1 || fnType.Out(0).Kind() != reflect.Bool {
		panic(fmt.Sprintf("assert: arguments: %s does not return a bool", fn))
	}

	return argumentMatcher{fn: reflect.ValueOf(fn)}
}


func (args Arguments) Get(index int) interface{} {
	if index+1 > len(args) {
		panic(fmt.Sprintf("assert: arguments: Cannot call Get(%d) because there are %d argument(s).", index, len(args)))
	}
	return args[index]
}


func (args Arguments) Is(objects ...interface{}) bool {
	for i, obj := range args {
		if obj != objects[i] {
			return false
		}
	}
	return true
}





func (args Arguments) Diff(objects []interface{}) (string, int) {
	

	var output = "\n"
	var differences int

	var maxArgCount = len(args)
	if len(objects) > maxArgCount {
		maxArgCount = len(objects)
	}

	for i := 0; i < maxArgCount; i++ {
		var actual, expected interface{}
		var actualFmt, expectedFmt string

		if len(objects) <= i {
			actual = "(Missing)"
			actualFmt = "(Missing)"
		} else {
			actual = objects[i]
			actualFmt = fmt.Sprintf("(%[1]T=%[1]v)", actual)
		}

		if len(args) <= i {
			expected = "(Missing)"
			expectedFmt = "(Missing)"
		} else {
			expected = args[i]
			expectedFmt = fmt.Sprintf("(%[1]T=%[1]v)", expected)
		}

		if matcher, ok := expected.(argumentMatcher); ok {
			if matcher.Matches(actual) {
				output = fmt.Sprintf("%s\t%d: PASS:  %s matched by %s\n", output, i, actualFmt, matcher)
			} else {
				differences++
				output = fmt.Sprintf("%s\t%d: PASS:  %s not matched by %s\n", output, i, actualFmt, matcher)
			}
		} else if reflect.TypeOf(expected) == reflect.TypeOf((*AnythingOfTypeArgument)(nil)).Elem() {

			
			if reflect.TypeOf(actual).Name() != string(expected.(AnythingOfTypeArgument)) && reflect.TypeOf(actual).String() != string(expected.(AnythingOfTypeArgument)) {
				
				differences++
				output = fmt.Sprintf("%s\t%d: FAIL:  type %s != type %s - %s\n", output, i, expected, reflect.TypeOf(actual).Name(), actualFmt)
			}

		} else {

			

			if assert.ObjectsAreEqual(expected, Anything) || assert.ObjectsAreEqual(actual, Anything) || assert.ObjectsAreEqual(actual, expected) {
				
				output = fmt.Sprintf("%s\t%d: PASS:  %s == %s\n", output, i, actualFmt, expectedFmt)
			} else {
				
				differences++
				output = fmt.Sprintf("%s\t%d: FAIL:  %s != %s\n", output, i, actualFmt, expectedFmt)
			}
		}

	}

	if differences == 0 {
		return "No differences.", differences
	}

	return output, differences

}



func (args Arguments) Assert(t TestingT, objects ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	
	diff, diffCount := args.Diff(objects)

	if diffCount == 0 {
		return true
	}

	
	t.Logf(diff)
	t.Errorf("%sArguments do not match.", assert.CallerInfo())

	return false

}






func (args Arguments) String(indexOrNil ...int) string {

	if len(indexOrNil) == 0 {
		
		var argsStr []string
		for _, arg := range args {
			argsStr = append(argsStr, fmt.Sprintf("%s", reflect.TypeOf(arg)))
		}
		return strings.Join(argsStr, ",")
	} else if len(indexOrNil) == 1 {
		
		var index = indexOrNil[0]
		var s string
		var ok bool
		if s, ok = args.Get(index).(string); !ok {
			panic(fmt.Sprintf("assert: arguments: String(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
		}
		return s
	}

	panic(fmt.Sprintf("assert: arguments: Wrong number of arguments passed to String.  Must be 0 or 1, not %d", len(indexOrNil)))

}



func (args Arguments) Int(index int) int {
	var s int
	var ok bool
	if s, ok = args.Get(index).(int); !ok {
		panic(fmt.Sprintf("assert: arguments: Int(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}



func (args Arguments) Error(index int) error {
	obj := args.Get(index)
	var s error
	var ok bool
	if obj == nil {
		return nil
	}
	if s, ok = obj.(error); !ok {
		panic(fmt.Sprintf("assert: arguments: Error(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}



func (args Arguments) Bool(index int) bool {
	var s bool
	var ok bool
	if s, ok = args.Get(index).(bool); !ok {
		panic(fmt.Sprintf("assert: arguments: Bool(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}

func typeAndKind(v interface{}) (reflect.Type, reflect.Kind) {
	t := reflect.TypeOf(v)
	k := t.Kind()

	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	return t, k
}

func diffArguments(expected Arguments, actual Arguments) string {
	if len(expected) != len(actual) {
		return fmt.Sprintf("Provided %v arguments, mocked for %v arguments", len(expected), len(actual))
	}

	for x := range expected {
		if diffString := diff(expected[x], actual[x]); diffString != "" {
			return fmt.Sprintf("Difference found in argument %v:\n\n%s", x, diffString)
		}
	}

	return ""
}



func diff(expected interface{}, actual interface{}) string {
	if expected == nil || actual == nil {
		return ""
	}

	et, ek := typeAndKind(expected)
	at, _ := typeAndKind(actual)

	if et != at {
		return ""
	}

	if ek != reflect.Struct && ek != reflect.Map && ek != reflect.Slice && ek != reflect.Array {
		return ""
	}

	e := spewConfig.Sdump(expected)
	a := spewConfig.Sdump(actual)

	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(e),
		B:        difflib.SplitLines(a),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  1,
	})

	return diff
}

var spewConfig = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
}

type tHelper interface {
	Helper()
}
