package kingpin

import (
	"os"
	"path/filepath"
)

var (
	
	CommandLine = New(filepath.Base(os.Args[0]), "")
	
	HelpFlag = CommandLine.HelpFlag
	
	HelpCommand = CommandLine.HelpCommand
	
	VersionFlag = CommandLine.VersionFlag
)


func Command(name, help string) *CmdClause {
	return CommandLine.Command(name, help)
}


func Flag(name, help string) *FlagClause {
	return CommandLine.Flag(name, help)
}


func Arg(name, help string) *ArgClause {
	return CommandLine.Arg(name, help)
}



func Parse() string {
	selected := MustParse(CommandLine.Parse(os.Args[1:]))
	if selected == "" && CommandLine.cmdGroup.have() {
		Usage()
		CommandLine.terminate(0)
	}
	return selected
}


func Errorf(format string, args ...interface{}) {
	CommandLine.Errorf(format, args...)
}


func Fatalf(format string, args ...interface{}) {
	CommandLine.Fatalf(format, args...)
}



func FatalIfError(err error, format string, args ...interface{}) {
	CommandLine.FatalIfError(err, format, args...)
}



func FatalUsage(format string, args ...interface{}) {
	CommandLine.FatalUsage(format, args...)
}



func FatalUsageContext(context *ParseContext, format string, args ...interface{}) {
	CommandLine.FatalUsageContext(context, format, args...)
}


func Usage() {
	CommandLine.Usage(os.Args[1:])
}


func UsageTemplate(template string) *Application {
	return CommandLine.UsageTemplate(template)
}


func MustParse(command string, err error) string {
	if err != nil {
		Fatalf("%s, try --help", err)
	}
	return command
}


func Version(version string) *Application {
	return CommandLine.Version(version)
}
