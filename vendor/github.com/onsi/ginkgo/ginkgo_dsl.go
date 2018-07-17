
package ginkgo

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/internal/codelocation"
	"github.com/onsi/ginkgo/internal/failer"
	"github.com/onsi/ginkgo/internal/remote"
	"github.com/onsi/ginkgo/internal/suite"
	"github.com/onsi/ginkgo/internal/testingtproxy"
	"github.com/onsi/ginkgo/internal/writer"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/ginkgo/reporters/stenographer"
	"github.com/onsi/ginkgo/types"
)

const GINKGO_VERSION = config.VERSION
const GINKGO_PANIC = `
Your test failed.
Ginkgo panics to prevent subsequent assertions from running.
Normally Ginkgo rescues this panic so you shouldn't see it.

But, if you make an assertion in a goroutine, Ginkgo can't capture the panic.
To circumvent this, you should call

	defer GinkgoRecover()

at the top of the goroutine that caused this panic.
`
const defaultTimeout = 1

var globalSuite *suite.Suite
var globalFailer *failer.Failer

func init() {
	config.Flags(flag.CommandLine, "ginkgo", true)
	GinkgoWriter = writer.New(os.Stdout)
	globalFailer = failer.New()
	globalSuite = suite.New(globalFailer)
}





var GinkgoWriter io.Writer


type GinkgoTestingT interface {
	Fail()
}





func GinkgoRandomSeed() int64 {
	return config.GinkgoConfig.RandomSeed
}



func GinkgoParallelNode() int {
	return config.GinkgoConfig.ParallelNode
}













func GinkgoT(optionalOffset ...int) GinkgoTInterface {
	offset := 3
	if len(optionalOffset) > 0 {
		offset = optionalOffset[0]
	}
	return testingtproxy.New(GinkgoWriter, Fail, offset)
}



type GinkgoTInterface interface {
	Fail()
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Failed() bool
	Parallel()
	Skip(args ...interface{})
	Skipf(format string, args ...interface{})
	SkipNow()
	Skipped() bool
}





type Reporter reporters.Reporter



type Done chan<- interface{}









type GinkgoTestDescription struct {
	FullTestText   string
	ComponentTexts []string
	TestText       string

	IsMeasurement bool

	FileName   string
	LineNumber int

	Failed   bool
	Duration time.Duration
}


func CurrentGinkgoTestDescription() GinkgoTestDescription {
	summary, ok := globalSuite.CurrentRunningSpecSummary()
	if !ok {
		return GinkgoTestDescription{}
	}

	subjectCodeLocation := summary.ComponentCodeLocations[len(summary.ComponentCodeLocations)-1]

	return GinkgoTestDescription{
		ComponentTexts: summary.ComponentTexts[1:],
		FullTestText:   strings.Join(summary.ComponentTexts[1:], " "),
		TestText:       summary.ComponentTexts[len(summary.ComponentTexts)-1],
		IsMeasurement:  summary.IsMeasurement,
		FileName:       subjectCodeLocation.FileName,
		LineNumber:     subjectCodeLocation.LineNumber,
		Failed:         summary.HasFailureState(),
		Duration:       summary.RunTime,
	}
}











type Benchmarker interface {
	Time(name string, body func(), info ...interface{}) (elapsedTime time.Duration)
	RecordValue(name string, value float64, info ...interface{})
	RecordValueWithPrecision(name string, value float64, units string, precision int, info ...interface{})
}







func RunSpecs(t GinkgoTestingT, description string) bool {
	specReporters := []Reporter{buildDefaultReporter()}
	return RunSpecsWithCustomReporters(t, description, specReporters)
}



func RunSpecsWithDefaultAndCustomReporters(t GinkgoTestingT, description string, specReporters []Reporter) bool {
	specReporters = append(specReporters, buildDefaultReporter())
	return RunSpecsWithCustomReporters(t, description, specReporters)
}



func RunSpecsWithCustomReporters(t GinkgoTestingT, description string, specReporters []Reporter) bool {
	writer := GinkgoWriter.(*writer.Writer)
	writer.SetStream(config.DefaultReporterConfig.Verbose)
	reporters := make([]reporters.Reporter, len(specReporters))
	for i, reporter := range specReporters {
		reporters[i] = reporter
	}
	passed, hasFocusedTests := globalSuite.Run(t, description, reporters, writer, config.GinkgoConfig)
	if passed && hasFocusedTests && strings.TrimSpace(os.Getenv("GINKGO_EDITOR_INTEGRATION")) == "" {
		fmt.Println("PASS | FOCUSED")
		os.Exit(types.GINKGO_FOCUS_EXIT_CODE)
	}
	return passed
}

func buildDefaultReporter() Reporter {
	remoteReportingServer := config.GinkgoConfig.StreamHost
	if remoteReportingServer == "" {
		stenographer := stenographer.New(!config.DefaultReporterConfig.NoColor, config.GinkgoConfig.FlakeAttempts > 1)
		return reporters.NewDefaultReporter(config.DefaultReporterConfig, stenographer)
	} else {
		return remote.NewForwardingReporter(remoteReportingServer, &http.Client{}, remote.NewOutputInterceptor())
	}
}


func Skip(message string, callerSkip ...int) {
	skip := 0
	if len(callerSkip) > 0 {
		skip = callerSkip[0]
	}

	globalFailer.Skip(message, codelocation.New(skip+1))
	panic(GINKGO_PANIC)
}


func Fail(message string, callerSkip ...int) {
	skip := 0
	if len(callerSkip) > 0 {
		skip = callerSkip[0]
	}

	globalFailer.Fail(message, codelocation.New(skip+1))
	panic(GINKGO_PANIC)
}











func GinkgoRecover() {
	e := recover()
	if e != nil {
		globalFailer.Panic(codelocation.New(1), e)
	}
}







func Describe(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypeNone, codelocation.New(1))
	return true
}


func FDescribe(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypeFocused, codelocation.New(1))
	return true
}


func PDescribe(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypePending, codelocation.New(1))
	return true
}


func XDescribe(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypePending, codelocation.New(1))
	return true
}







func Context(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypeNone, codelocation.New(1))
	return true
}


func FContext(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypeFocused, codelocation.New(1))
	return true
}


func PContext(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypePending, codelocation.New(1))
	return true
}


func XContext(text string, body func()) bool {
	globalSuite.PushContainerNode(text, body, types.FlagTypePending, codelocation.New(1))
	return true
}







func When(text string, body func()) bool {
	globalSuite.PushContainerNode("when "+text, body, types.FlagTypeNone, codelocation.New(1))
	return true
}


func FWhen(text string, body func()) bool {
	globalSuite.PushContainerNode("when "+text, body, types.FlagTypeFocused, codelocation.New(1))
	return true
}


func PWhen(text string, body func()) bool {
	globalSuite.PushContainerNode("when "+text, body, types.FlagTypePending, codelocation.New(1))
	return true
}


func XWhen(text string, body func()) bool {
	globalSuite.PushContainerNode("when "+text, body, types.FlagTypePending, codelocation.New(1))
	return true
}






func It(text string, body interface{}, timeout ...float64) bool {
	globalSuite.PushItNode(text, body, types.FlagTypeNone, codelocation.New(1), parseTimeout(timeout...))
	return true
}


func FIt(text string, body interface{}, timeout ...float64) bool {
	globalSuite.PushItNode(text, body, types.FlagTypeFocused, codelocation.New(1), parseTimeout(timeout...))
	return true
}


func PIt(text string, _ ...interface{}) bool {
	globalSuite.PushItNode(text, func() {}, types.FlagTypePending, codelocation.New(1), 0)
	return true
}


func XIt(text string, _ ...interface{}) bool {
	globalSuite.PushItNode(text, func() {}, types.FlagTypePending, codelocation.New(1), 0)
	return true
}




func Specify(text string, body interface{}, timeout ...float64) bool {
	globalSuite.PushItNode(text, body, types.FlagTypeNone, codelocation.New(1), parseTimeout(timeout...))
	return true
}


func FSpecify(text string, body interface{}, timeout ...float64) bool {
	globalSuite.PushItNode(text, body, types.FlagTypeFocused, codelocation.New(1), parseTimeout(timeout...))
	return true
}


func PSpecify(text string, is ...interface{}) bool {
	globalSuite.PushItNode(text, func() {}, types.FlagTypePending, codelocation.New(1), 0)
	return true
}


func XSpecify(text string, is ...interface{}) bool {
	globalSuite.PushItNode(text, func() {}, types.FlagTypePending, codelocation.New(1), 0)
	return true
}








func By(text string, callbacks ...func()) {
	preamble := "\x1b[1mSTEP\x1b[0m"
	if config.DefaultReporterConfig.NoColor {
		preamble = "STEP"
	}
	fmt.Fprintln(GinkgoWriter, preamble+": "+text)
	if len(callbacks) == 1 {
		callbacks[0]()
	}
	if len(callbacks) > 1 {
		panic("just one callback per By, please")
	}
}






func Measure(text string, body interface{}, samples int) bool {
	globalSuite.PushMeasureNode(text, body, types.FlagTypeNone, codelocation.New(1), samples)
	return true
}


func FMeasure(text string, body interface{}, samples int) bool {
	globalSuite.PushMeasureNode(text, body, types.FlagTypeFocused, codelocation.New(1), samples)
	return true
}


func PMeasure(text string, _ ...interface{}) bool {
	globalSuite.PushMeasureNode(text, func(b Benchmarker) {}, types.FlagTypePending, codelocation.New(1), 0)
	return true
}


func XMeasure(text string, _ ...interface{}) bool {
	globalSuite.PushMeasureNode(text, func(b Benchmarker) {}, types.FlagTypePending, codelocation.New(1), 0)
	return true
}







func BeforeSuite(body interface{}, timeout ...float64) bool {
	globalSuite.SetBeforeSuiteNode(body, codelocation.New(1), parseTimeout(timeout...))
	return true
}









func AfterSuite(body interface{}, timeout ...float64) bool {
	globalSuite.SetAfterSuiteNode(body, codelocation.New(1), parseTimeout(timeout...))
	return true
}









































func SynchronizedBeforeSuite(node1Body interface{}, allNodesBody interface{}, timeout ...float64) bool {
	globalSuite.SetSynchronizedBeforeSuiteNode(
		node1Body,
		allNodesBody,
		codelocation.New(1),
		parseTimeout(timeout...),
	)
	return true
}


















func SynchronizedAfterSuite(allNodesBody interface{}, node1Body interface{}, timeout ...float64) bool {
	globalSuite.SetSynchronizedAfterSuiteNode(
		allNodesBody,
		node1Body,
		codelocation.New(1),
		parseTimeout(timeout...),
	)
	return true
}






func BeforeEach(body interface{}, timeout ...float64) bool {
	globalSuite.PushBeforeEachNode(body, codelocation.New(1), parseTimeout(timeout...))
	return true
}






func JustBeforeEach(body interface{}, timeout ...float64) bool {
	globalSuite.PushJustBeforeEachNode(body, codelocation.New(1), parseTimeout(timeout...))
	return true
}






func AfterEach(body interface{}, timeout ...float64) bool {
	globalSuite.PushAfterEachNode(body, codelocation.New(1), parseTimeout(timeout...))
	return true
}

func parseTimeout(timeout ...float64) time.Duration {
	if len(timeout) == 0 {
		return time.Duration(defaultTimeout * int64(time.Second))
	} else {
		return time.Duration(timeout[0] * float64(time.Second))
	}
}
