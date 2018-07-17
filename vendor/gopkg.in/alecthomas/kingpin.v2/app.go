package kingpin

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	ErrCommandNotSpecified = fmt.Errorf("command not specified")
)

var (
	envarTransformRegexp = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
)

type ApplicationValidator func(*Application) error



type Application struct {
	cmdMixin
	initialized bool

	Name string
	Help string

	author         string
	version        string
	errorWriter    io.Writer 
	usageWriter    io.Writer 
	usageTemplate  string
	validator      ApplicationValidator
	terminate      func(status int) 
	noInterspersed bool             
	defaultEnvars  bool
	completion     bool

	
	HelpFlag *FlagClause
	
	HelpCommand *CmdClause
	
	VersionFlag *FlagClause
}


func New(name, help string) *Application {
	a := &Application{
		Name:          name,
		Help:          help,
		errorWriter:   os.Stderr, 
		usageWriter:   os.Stderr,
		usageTemplate: DefaultUsageTemplate,
		terminate:     os.Exit,
	}
	a.flagGroup = newFlagGroup()
	a.argGroup = newArgGroup()
	a.cmdGroup = newCmdGroup(a)
	a.HelpFlag = a.Flag("help", "Show context-sensitive help (also try --help-long and --help-man).")
	a.HelpFlag.Bool()
	a.Flag("help-long", "Generate long help.").Hidden().PreAction(a.generateLongHelp).Bool()
	a.Flag("help-man", "Generate a man page.").Hidden().PreAction(a.generateManPage).Bool()
	a.Flag("completion-bash", "Output possible completions for the given args.").Hidden().BoolVar(&a.completion)
	a.Flag("completion-script-bash", "Generate completion script for bash.").Hidden().PreAction(a.generateBashCompletionScript).Bool()
	a.Flag("completion-script-zsh", "Generate completion script for ZSH.").Hidden().PreAction(a.generateZSHCompletionScript).Bool()

	return a
}

func (a *Application) generateLongHelp(c *ParseContext) error {
	a.Writer(os.Stdout)
	if err := a.UsageForContextWithTemplate(c, 2, LongHelpTemplate); err != nil {
		return err
	}
	a.terminate(0)
	return nil
}

func (a *Application) generateManPage(c *ParseContext) error {
	a.Writer(os.Stdout)
	if err := a.UsageForContextWithTemplate(c, 2, ManPageTemplate); err != nil {
		return err
	}
	a.terminate(0)
	return nil
}

func (a *Application) generateBashCompletionScript(c *ParseContext) error {
	a.Writer(os.Stdout)
	if err := a.UsageForContextWithTemplate(c, 2, BashCompletionTemplate); err != nil {
		return err
	}
	a.terminate(0)
	return nil
}

func (a *Application) generateZSHCompletionScript(c *ParseContext) error {
	a.Writer(os.Stdout)
	if err := a.UsageForContextWithTemplate(c, 2, ZshCompletionTemplate); err != nil {
		return err
	}
	a.terminate(0)
	return nil
}






func (a *Application) DefaultEnvars() *Application {
	a.defaultEnvars = true
	return a
}



func (a *Application) Terminate(terminate func(int)) *Application {
	if terminate == nil {
		terminate = func(int) {}
	}
	a.terminate = terminate
	return a
}



func (a *Application) Writer(w io.Writer) *Application {
	a.errorWriter = w
	a.usageWriter = w
	return a
}


func (a *Application) ErrorWriter(w io.Writer) *Application {
	a.errorWriter = w
	return a
}


func (a *Application) UsageWriter(w io.Writer) *Application {
	a.usageWriter = w
	return a
}



func (a *Application) UsageTemplate(template string) *Application {
	a.usageTemplate = template
	return a
}


func (a *Application) Validate(validator ApplicationValidator) *Application {
	a.validator = validator
	return a
}



func (a *Application) ParseContext(args []string) (*ParseContext, error) {
	return a.parseContext(false, args)
}

func (a *Application) parseContext(ignoreDefault bool, args []string) (*ParseContext, error) {
	if err := a.init(); err != nil {
		return nil, err
	}
	context := tokenize(args, ignoreDefault)
	err := parse(context, a)
	return context, err
}







func (a *Application) Parse(args []string) (command string, err error) {

	context, parseErr := a.ParseContext(args)
	selected := []string{}
	var setValuesErr error

	if context == nil {
		
		
		return "", parseErr
	}

	if err = a.setDefaults(context); err != nil {
		return "", err
	}

	selected, setValuesErr = a.setValues(context)

	if err = a.applyPreActions(context, !a.completion); err != nil {
		return "", err
	}

	if a.completion {
		a.generateBashCompletion(context)
		a.terminate(0)
	} else {
		if parseErr != nil {
			return "", parseErr
		}

		a.maybeHelp(context)
		if !context.EOL() {
			return "", fmt.Errorf("unexpected argument '%s'", context.Peek())
		}

		if setValuesErr != nil {
			return "", setValuesErr
		}

		command, err = a.execute(context, selected)
		if err == ErrCommandNotSpecified {
			a.writeUsage(context, nil)
		}
	}
	return command, err
}

func (a *Application) writeUsage(context *ParseContext, err error) {
	if err != nil {
		a.Errorf("%s", err)
	}
	if err := a.UsageForContext(context); err != nil {
		panic(err)
	}
	if err != nil {
		a.terminate(1)
	} else {
		a.terminate(0)
	}
}

func (a *Application) maybeHelp(context *ParseContext) {
	for _, element := range context.Elements {
		if flag, ok := element.Clause.(*FlagClause); ok && flag == a.HelpFlag {
			
			context, _ = a.parseContext(true, context.rawArgs)
			a.writeUsage(context, nil)
		}
	}
}


func (a *Application) Version(version string) *Application {
	a.version = version
	a.VersionFlag = a.Flag("version", "Show application version.").PreAction(func(*ParseContext) error {
		fmt.Fprintln(a.usageWriter, version)
		a.terminate(0)
		return nil
	})
	a.VersionFlag.Bool()
	return a
}


func (a *Application) Author(author string) *Application {
	a.author = author
	return a
}






func (a *Application) Action(action Action) *Application {
	a.addAction(action)
	return a
}


func (a *Application) PreAction(action Action) *Application {
	a.addPreAction(action)
	return a
}


func (a *Application) Command(name, help string) *CmdClause {
	return a.addCommand(name, help)
}




func (a *Application) Interspersed(interspersed bool) *Application {
	a.noInterspersed = !interspersed
	return a
}

func (a *Application) defaultEnvarPrefix() string {
	if a.defaultEnvars {
		return a.Name
	}
	return ""
}

func (a *Application) init() error {
	if a.initialized {
		return nil
	}
	if a.cmdGroup.have() && a.argGroup.have() {
		return fmt.Errorf("can't mix top-level Arg()s with Command()s")
	}

	
	if a.cmdGroup.have() {
		var command []string
		a.HelpCommand = a.Command("help", "Show help.").PreAction(func(context *ParseContext) error {
			a.Usage(command)
			a.terminate(0)
			return nil
		})
		a.HelpCommand.Arg("command", "Show help on command.").StringsVar(&command)
		
		l := len(a.commandOrder)
		a.commandOrder = append(a.commandOrder[l-1:l], a.commandOrder[:l-1]...)
	}

	if err := a.flagGroup.init(a.defaultEnvarPrefix()); err != nil {
		return err
	}
	if err := a.cmdGroup.init(); err != nil {
		return err
	}
	if err := a.argGroup.init(); err != nil {
		return err
	}
	for _, cmd := range a.commands {
		if err := cmd.init(); err != nil {
			return err
		}
	}
	flagGroups := []*flagGroup{a.flagGroup}
	for _, cmd := range a.commandOrder {
		if err := checkDuplicateFlags(cmd, flagGroups); err != nil {
			return err
		}
	}
	a.initialized = true
	return nil
}


func checkDuplicateFlags(current *CmdClause, flagGroups []*flagGroup) error {
	
	for _, flags := range flagGroups {
		for _, flag := range current.flagOrder {
			if flag.shorthand != 0 {
				if _, ok := flags.short[string(flag.shorthand)]; ok {
					return fmt.Errorf("duplicate short flag -%c", flag.shorthand)
				}
			}
			if _, ok := flags.long[flag.name]; ok {
				return fmt.Errorf("duplicate long flag --%s", flag.name)
			}
		}
	}
	flagGroups = append(flagGroups, current.flagGroup)
	
	for _, subcmd := range current.commandOrder {
		if err := checkDuplicateFlags(subcmd, flagGroups); err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) execute(context *ParseContext, selected []string) (string, error) {
	var err error

	if err = a.validateRequired(context); err != nil {
		return "", err
	}

	if err = a.applyValidators(context); err != nil {
		return "", err
	}

	if err = a.applyActions(context); err != nil {
		return "", err
	}

	command := strings.Join(selected, " ")
	if command == "" && a.cmdGroup.have() {
		return "", ErrCommandNotSpecified
	}
	return command, err
}

func (a *Application) setDefaults(context *ParseContext) error {
	flagElements := map[string]*ParseElement{}
	for _, element := range context.Elements {
		if flag, ok := element.Clause.(*FlagClause); ok {
			if flag.name == "help" {
				return nil
			}
			flagElements[flag.name] = element
		}
	}

	argElements := map[string]*ParseElement{}
	for _, element := range context.Elements {
		if arg, ok := element.Clause.(*ArgClause); ok {
			argElements[arg.name] = element
		}
	}

	
	for _, flag := range context.flags.long {
		if flagElements[flag.name] == nil {
			if err := flag.setDefault(); err != nil {
				return err
			}
		}
	}

	for _, arg := range context.arguments.args {
		if argElements[arg.name] == nil {
			if err := arg.setDefault(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Application) validateRequired(context *ParseContext) error {
	flagElements := map[string]*ParseElement{}
	for _, element := range context.Elements {
		if flag, ok := element.Clause.(*FlagClause); ok {
			flagElements[flag.name] = element
		}
	}

	argElements := map[string]*ParseElement{}
	for _, element := range context.Elements {
		if arg, ok := element.Clause.(*ArgClause); ok {
			argElements[arg.name] = element
		}
	}

	
	for _, flag := range context.flags.long {
		if flagElements[flag.name] == nil {
			
			if flag.needsValue() {
				return fmt.Errorf("required flag --%s not provided", flag.name)
			}
		}
	}

	for _, arg := range context.arguments.args {
		if argElements[arg.name] == nil {
			if arg.needsValue() {
				return fmt.Errorf("required argument '%s' not provided", arg.name)
			}
		}
	}
	return nil
}

func (a *Application) setValues(context *ParseContext) (selected []string, err error) {
	
	var (
		lastCmd *CmdClause
		flagSet = map[string]struct{}{}
	)
	for _, element := range context.Elements {
		switch clause := element.Clause.(type) {
		case *FlagClause:
			if _, ok := flagSet[clause.name]; ok {
				if v, ok := clause.value.(repeatableFlag); !ok || !v.IsCumulative() {
					return nil, fmt.Errorf("flag '%s' cannot be repeated", clause.name)
				}
			}
			if err = clause.value.Set(*element.Value); err != nil {
				return
			}
			flagSet[clause.name] = struct{}{}

		case *ArgClause:
			if err = clause.value.Set(*element.Value); err != nil {
				return
			}

		case *CmdClause:
			if clause.validator != nil {
				if err = clause.validator(clause); err != nil {
					return
				}
			}
			selected = append(selected, clause.name)
			lastCmd = clause
		}
	}

	if lastCmd != nil && len(lastCmd.commands) > 0 {
		return nil, fmt.Errorf("must select a subcommand of '%s'", lastCmd.FullCommand())
	}

	return
}

func (a *Application) applyValidators(context *ParseContext) (err error) {
	
	for _, element := range context.Elements {
		if cmd, ok := element.Clause.(*CmdClause); ok && cmd.validator != nil {
			if err = cmd.validator(cmd); err != nil {
				return err
			}
		}
	}

	if a.validator != nil {
		err = a.validator(a)
	}
	return err
}

func (a *Application) applyPreActions(context *ParseContext, dispatch bool) error {
	if err := a.actionMixin.applyPreActions(context); err != nil {
		return err
	}
	
	if dispatch {
		for _, element := range context.Elements {
			if applier, ok := element.Clause.(actionApplier); ok {
				if err := applier.applyPreActions(context); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (a *Application) applyActions(context *ParseContext) error {
	if err := a.actionMixin.applyActions(context); err != nil {
		return err
	}
	
	for _, element := range context.Elements {
		if applier, ok := element.Clause.(actionApplier); ok {
			if err := applier.applyActions(context); err != nil {
				return err
			}
		}
	}
	return nil
}


func (a *Application) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(a.errorWriter, a.Name+": error: "+format+"\n", args...)
}


func (a *Application) Fatalf(format string, args ...interface{}) {
	a.Errorf(format, args...)
	a.terminate(1)
}



func (a *Application) FatalUsage(format string, args ...interface{}) {
	a.Errorf(format, args...)
	
	a.usageWriter = a.errorWriter
	a.Usage([]string{})
	a.terminate(1)
}



func (a *Application) FatalUsageContext(context *ParseContext, format string, args ...interface{}) {
	a.Errorf(format, args...)
	if err := a.UsageForContext(context); err != nil {
		panic(err)
	}
	a.terminate(1)
}



func (a *Application) FatalIfError(err error, format string, args ...interface{}) {
	if err != nil {
		prefix := ""
		if format != "" {
			prefix = fmt.Sprintf(format, args...) + ": "
		}
		a.Errorf(prefix+"%s", err)
		a.terminate(1)
	}
}

func (a *Application) completionOptions(context *ParseContext) []string {
	args := context.rawArgs

	var (
		currArg string
		prevArg string
		target  cmdMixin
	)

	numArgs := len(args)
	if numArgs > 1 {
		args = args[1:]
		currArg = args[len(args)-1]
	}
	if numArgs > 2 {
		prevArg = args[len(args)-2]
	}

	target = a.cmdMixin
	if context.SelectedCommand != nil {
		
		target = context.SelectedCommand.cmdMixin
	}

	if (currArg != "" && strings.HasPrefix(currArg, "--")) || strings.HasPrefix(prevArg, "--") {
		
		var (
			flagName  string 
			flagValue string 
		)

		if strings.HasPrefix(prevArg, "--") && !strings.HasPrefix(currArg, "--") {
			
			
			flagName = prevArg[2:] 
			flagValue = currArg
		} else if strings.HasPrefix(currArg, "--") {
			
			
			
			flagName = currArg[2:] 
		}

		options, flagMatched, valueMatched := target.FlagCompletion(flagName, flagValue)
		if valueMatched {
			
			return target.CmdCompletion(context)
		}

		
		if context.SelectedCommand != nil && !flagMatched {
			topOptions, topFlagMatched, topValueMatched := a.FlagCompletion(flagName, flagValue)
			if topValueMatched {
				
				return target.CmdCompletion(context)
			}

			if topFlagMatched {
				
				options = topOptions
			} else {
				
				options = append(options, topOptions...)
			}
		}
		return options
	}

	
	return target.CmdCompletion(context)
}

func (a *Application) generateBashCompletion(context *ParseContext) {
	options := a.completionOptions(context)
	fmt.Printf("%s", strings.Join(options, "\n"))
}

func envarTransform(name string) string {
	return strings.ToUpper(envarTransformRegexp.ReplaceAllString(name, "_"))
}
