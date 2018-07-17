package kingpin

import (
	"fmt"
	"strings"
)

type cmdMixin struct {
	*flagGroup
	*argGroup
	*cmdGroup
	actionMixin
}



func (c *cmdMixin) CmdCompletion(context *ParseContext) []string {
	var options []string

	
	
	
	argsSatisfied := 0
	for _, el := range context.Elements {
		switch clause := el.Clause.(type) {
		case *ArgClause:
			if el.Value != nil && *el.Value != "" {
				argsSatisfied++
			}
		case *CmdClause:
			options = append(options, clause.completionAlts...)
		default:
		}
	}

	if argsSatisfied < len(c.argGroup.args) {
		
		options = append(options, c.argGroup.args[argsSatisfied].resolveCompletions()...)
	} else {
		
		for _, cmd := range c.cmdGroup.commandOrder {
			if !cmd.hidden {
				options = append(options, cmd.name)
			}
		}
	}

	return options
}

func (c *cmdMixin) FlagCompletion(flagName string, flagValue string) (choices []string, flagMatch bool, optionMatch bool) {
	
	
	

	options := []string{}

	for _, flag := range c.flagGroup.flagOrder {
		
		if flag.name == flagName {
			
			options = flag.resolveCompletions()
			if len(options) == 0 {
				
				return options, true, true
			}

			
			isPrefix := false
			matched := false

			for _, opt := range options {
				if flagValue == opt {
					matched = true
				} else if strings.HasPrefix(opt, flagValue) {
					isPrefix = true
				}
			}

			
			
			return options, true, !isPrefix && matched
		}

		if !flag.hidden {
			options = append(options, "--"+flag.name)
		}
	}
	
	return options, false, false

}

type cmdGroup struct {
	app          *Application
	parent       *CmdClause
	commands     map[string]*CmdClause
	commandOrder []*CmdClause
}

func (c *cmdGroup) defaultSubcommand() *CmdClause {
	for _, cmd := range c.commandOrder {
		if cmd.isDefault {
			return cmd
		}
	}
	return nil
}

func (c *cmdGroup) cmdNames() []string {
	names := make([]string, 0, len(c.commandOrder))
	for _, cmd := range c.commandOrder {
		names = append(names, cmd.name)
	}
	return names
}





func (c *cmdGroup) GetCommand(name string) *CmdClause {
	return c.commands[name]
}

func newCmdGroup(app *Application) *cmdGroup {
	return &cmdGroup{
		app:      app,
		commands: make(map[string]*CmdClause),
	}
}

func (c *cmdGroup) flattenedCommands() (out []*CmdClause) {
	for _, cmd := range c.commandOrder {
		if len(cmd.commands) == 0 {
			out = append(out, cmd)
		}
		out = append(out, cmd.flattenedCommands()...)
	}
	return
}

func (c *cmdGroup) addCommand(name, help string) *CmdClause {
	cmd := newCommand(c.app, name, help)
	c.commands[name] = cmd
	c.commandOrder = append(c.commandOrder, cmd)
	return cmd
}

func (c *cmdGroup) init() error {
	seen := map[string]bool{}
	if c.defaultSubcommand() != nil && !c.have() {
		return fmt.Errorf("default subcommand %q provided but no subcommands defined", c.defaultSubcommand().name)
	}
	defaults := []string{}
	for _, cmd := range c.commandOrder {
		if cmd.isDefault {
			defaults = append(defaults, cmd.name)
		}
		if seen[cmd.name] {
			return fmt.Errorf("duplicate command %q", cmd.name)
		}
		seen[cmd.name] = true
		for _, alias := range cmd.aliases {
			if seen[alias] {
				return fmt.Errorf("alias duplicates existing command %q", alias)
			}
			c.commands[alias] = cmd
		}
		if err := cmd.init(); err != nil {
			return err
		}
	}
	if len(defaults) > 1 {
		return fmt.Errorf("more than one default subcommand exists: %s", strings.Join(defaults, ", "))
	}
	return nil
}

func (c *cmdGroup) have() bool {
	return len(c.commands) > 0
}

type CmdClauseValidator func(*CmdClause) error



type CmdClause struct {
	cmdMixin
	app            *Application
	name           string
	aliases        []string
	help           string
	isDefault      bool
	validator      CmdClauseValidator
	hidden         bool
	completionAlts []string
}

func newCommand(app *Application, name, help string) *CmdClause {
	c := &CmdClause{
		app:  app,
		name: name,
		help: help,
	}
	c.flagGroup = newFlagGroup()
	c.argGroup = newArgGroup()
	c.cmdGroup = newCmdGroup(app)
	return c
}


func (c *CmdClause) Alias(name string) *CmdClause {
	c.aliases = append(c.aliases, name)
	return c
}


func (c *CmdClause) Validate(validator CmdClauseValidator) *CmdClause {
	c.validator = validator
	return c
}

func (c *CmdClause) FullCommand() string {
	out := []string{c.name}
	for p := c.parent; p != nil; p = p.parent {
		out = append([]string{p.name}, out...)
	}
	return strings.Join(out, " ")
}


func (c *CmdClause) Command(name, help string) *CmdClause {
	cmd := c.addCommand(name, help)
	cmd.parent = c
	return cmd
}


func (c *CmdClause) Default() *CmdClause {
	c.isDefault = true
	return c
}

func (c *CmdClause) Action(action Action) *CmdClause {
	c.addAction(action)
	return c
}

func (c *CmdClause) PreAction(action Action) *CmdClause {
	c.addPreAction(action)
	return c
}

func (c *CmdClause) init() error {
	if err := c.flagGroup.init(c.app.defaultEnvarPrefix()); err != nil {
		return err
	}
	if c.argGroup.have() && c.cmdGroup.have() {
		return fmt.Errorf("can't mix Arg()s with Command()s")
	}
	if err := c.argGroup.init(); err != nil {
		return err
	}
	if err := c.cmdGroup.init(); err != nil {
		return err
	}
	return nil
}

func (c *CmdClause) Hidden() *CmdClause {
	c.hidden = true
	return c
}
