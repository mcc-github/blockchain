
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/ginkgo/testsuite"
)

const greenColor = "\x1b[32m"
const redColor = "\x1b[91m"
const defaultStyle = "\x1b[0m"
const lightGrayColor = "\x1b[37m"

type Command struct {
	Name                      string
	AltName                   string
	FlagSet                   *flag.FlagSet
	Usage                     []string
	UsageCommand              string
	Command                   func(args []string, additionalArgs []string)
	SuppressFlagDocumentation bool
	FlagDocSubstitute         []string
}

func (c *Command) Matches(name string) bool {
	return c.Name == name || (c.AltName != "" && c.AltName == name)
}

func (c *Command) Run(args []string, additionalArgs []string) {
	c.FlagSet.Parse(args)
	c.Command(c.FlagSet.Args(), additionalArgs)
}

var DefaultCommand *Command
var Commands []*Command

func init() {
	DefaultCommand = BuildRunCommand()
	Commands = append(Commands, BuildWatchCommand())
	Commands = append(Commands, BuildBuildCommand())
	Commands = append(Commands, BuildBootstrapCommand())
	Commands = append(Commands, BuildGenerateCommand())
	Commands = append(Commands, BuildNodotCommand())
	Commands = append(Commands, BuildConvertCommand())
	Commands = append(Commands, BuildUnfocusCommand())
	Commands = append(Commands, BuildVersionCommand())
	Commands = append(Commands, BuildHelpCommand())
}

func main() {
	args := []string{}
	additionalArgs := []string{}

	foundDelimiter := false

	for _, arg := range os.Args[1:] {
		if !foundDelimiter {
			if arg == "--" {
				foundDelimiter = true
				continue
			}
		}

		if foundDelimiter {
			additionalArgs = append(additionalArgs, arg)
		} else {
			args = append(args, arg)
		}
	}

	if len(args) > 0 {
		commandToRun, found := commandMatching(args[0])
		if found {
			commandToRun.Run(args[1:], additionalArgs)
			return
		}
	}

	DefaultCommand.Run(args, additionalArgs)
}

func commandMatching(name string) (*Command, bool) {
	for _, command := range Commands {
		if command.Matches(name) {
			return command, true
		}
	}
	return nil, false
}

func usage() {
	fmt.Fprintf(os.Stderr, "Ginkgo Version %s\n\n", config.VERSION)
	usageForCommand(DefaultCommand, false)
	for _, command := range Commands {
		fmt.Fprintf(os.Stderr, "\n")
		usageForCommand(command, false)
	}
}

func usageForCommand(command *Command, longForm bool) {
	fmt.Fprintf(os.Stderr, "%s\n%s\n", command.UsageCommand, strings.Repeat("-", len(command.UsageCommand)))
	fmt.Fprintf(os.Stderr, "%s\n", strings.Join(command.Usage, "\n"))
	if command.SuppressFlagDocumentation && !longForm {
		fmt.Fprintf(os.Stderr, "%s\n", strings.Join(command.FlagDocSubstitute, "\n  "))
	} else {
		command.FlagSet.PrintDefaults()
	}
}

func complainAndQuit(complaint string) {
	fmt.Fprintf(os.Stderr, "%s\nFor usage instructions:\n\tginkgo help\n", complaint)
	os.Exit(1)
}

func findSuites(args []string, recurseForAll bool, skipPackage string, allowPrecompiled bool) ([]testsuite.TestSuite, []string) {
	suites := []testsuite.TestSuite{}

	if len(args) > 0 {
		for _, arg := range args {
			if allowPrecompiled {
				suite, err := testsuite.PrecompiledTestSuite(arg)
				if err == nil {
					suites = append(suites, suite)
					continue
				}
			}
			recurseForSuite := recurseForAll
			if strings.HasSuffix(arg, "/...") && arg != "/..." {
				arg = arg[:len(arg)-4]
				recurseForSuite = true
			}
			suites = append(suites, testsuite.SuitesInDir(arg, recurseForSuite)...)
		}
	} else {
		suites = testsuite.SuitesInDir(".", recurseForAll)
	}

	skippedPackages := []string{}
	if skipPackage != "" {
		skipFilters := strings.Split(skipPackage, ",")
		filteredSuites := []testsuite.TestSuite{}
		for _, suite := range suites {
			skip := false
			for _, skipFilter := range skipFilters {
				if strings.Contains(suite.Path, skipFilter) {
					skip = true
					break
				}
			}
			if skip {
				skippedPackages = append(skippedPackages, suite.Path)
			} else {
				filteredSuites = append(filteredSuites, suite)
			}
		}
		suites = filteredSuites
	}

	return suites, skippedPackages
}

func goFmt(path string) {
	err := exec.Command("go", "fmt", path).Run()
	if err != nil {
		complainAndQuit("Could not fmt: " + err.Error())
	}
}

func pluralizedWord(singular, plural string, count int) string {
	if count == 1 {
		return singular
	}
	return plural
}
