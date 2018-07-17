package kingpin

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

type TokenType int


const (
	TokenShort TokenType = iota
	TokenLong
	TokenArg
	TokenError
	TokenEOL
)

func (t TokenType) String() string {
	switch t {
	case TokenShort:
		return "short flag"
	case TokenLong:
		return "long flag"
	case TokenArg:
		return "argument"
	case TokenError:
		return "error"
	case TokenEOL:
		return "<EOL>"
	}
	return "?"
}

var (
	TokenEOLMarker = Token{-1, TokenEOL, ""}
)

type Token struct {
	Index int
	Type  TokenType
	Value string
}

func (t *Token) Equal(o *Token) bool {
	return t.Index == o.Index
}

func (t *Token) IsFlag() bool {
	return t.Type == TokenShort || t.Type == TokenLong
}

func (t *Token) IsEOF() bool {
	return t.Type == TokenEOL
}

func (t *Token) String() string {
	switch t.Type {
	case TokenShort:
		return "-" + t.Value
	case TokenLong:
		return "--" + t.Value
	case TokenArg:
		return t.Value
	case TokenError:
		return "error: " + t.Value
	case TokenEOL:
		return "<EOL>"
	default:
		panic("unhandled type")
	}
}


type ParseElement struct {
	
	Clause interface{}
	
	Value *string
}





type ParseContext struct {
	SelectedCommand *CmdClause
	ignoreDefault   bool
	argsOnly        bool
	peek            []*Token
	argi            int 
	args            []string
	rawArgs         []string
	flags           *flagGroup
	arguments       *argGroup
	argumenti       int 
	
	Elements []*ParseElement
}

func (p *ParseContext) nextArg() *ArgClause {
	if p.argumenti >= len(p.arguments.args) {
		return nil
	}
	arg := p.arguments.args[p.argumenti]
	if !arg.consumesRemainder() {
		p.argumenti++
	}
	return arg
}

func (p *ParseContext) next() {
	p.argi++
	p.args = p.args[1:]
}



func (p *ParseContext) HasTrailingArgs() bool {
	return len(p.args) > 0
}

func tokenize(args []string, ignoreDefault bool) *ParseContext {
	return &ParseContext{
		ignoreDefault: ignoreDefault,
		args:          args,
		rawArgs:       args,
		flags:         newFlagGroup(),
		arguments:     newArgGroup(),
	}
}

func (p *ParseContext) mergeFlags(flags *flagGroup) {
	for _, flag := range flags.flagOrder {
		if flag.shorthand != 0 {
			p.flags.short[string(flag.shorthand)] = flag
		}
		p.flags.long[flag.name] = flag
		p.flags.flagOrder = append(p.flags.flagOrder, flag)
	}
}

func (p *ParseContext) mergeArgs(args *argGroup) {
	for _, arg := range args.args {
		p.arguments.args = append(p.arguments.args, arg)
	}
}

func (p *ParseContext) EOL() bool {
	return p.Peek().Type == TokenEOL
}

func (p *ParseContext) Error() bool {
	return p.Peek().Type == TokenError
}


func (p *ParseContext) Next() *Token {
	if len(p.peek) > 0 {
		return p.pop()
	}

	
	if len(p.args) == 0 {
		return &Token{Index: p.argi, Type: TokenEOL}
	}

	arg := p.args[0]
	p.next()

	if p.argsOnly {
		return &Token{p.argi, TokenArg, arg}
	}

	
	if arg == "--" {
		p.argsOnly = true
		return p.Next()
	}

	if strings.HasPrefix(arg, "--") {
		parts := strings.SplitN(arg[2:], "=", 2)
		token := &Token{p.argi, TokenLong, parts[0]}
		if len(parts) == 2 {
			p.Push(&Token{p.argi, TokenArg, parts[1]})
		}
		return token
	}

	if strings.HasPrefix(arg, "-") {
		if len(arg) == 1 {
			return &Token{Index: p.argi, Type: TokenShort}
		}
		shortRune, size := utf8.DecodeRuneInString(arg[1:])
		short := string(shortRune)
		flag, ok := p.flags.short[short]
		
		if !ok {
		} else if fb, ok := flag.value.(boolFlag); ok && fb.IsBoolFlag() {
			
		} else {
			
			token := &Token{p.argi, TokenShort, short}
			if len(arg) > size+1 {
				p.Push(&Token{p.argi, TokenArg, arg[size+1:]})
			}
			return token
		}

		if len(arg) > size+1 {
			p.args = append([]string{"-" + arg[size+1:]}, p.args...)
		}
		return &Token{p.argi, TokenShort, short}
	} else if strings.HasPrefix(arg, "@") {
		expanded, err := ExpandArgsFromFile(arg[1:])
		if err != nil {
			return &Token{p.argi, TokenError, err.Error()}
		}
		if len(p.args) == 0 {
			p.args = expanded
		} else {
			p.args = append(expanded, p.args...)
		}
		return p.Next()
	}

	return &Token{p.argi, TokenArg, arg}
}

func (p *ParseContext) Peek() *Token {
	if len(p.peek) == 0 {
		return p.Push(p.Next())
	}
	return p.peek[len(p.peek)-1]
}

func (p *ParseContext) Push(token *Token) *Token {
	p.peek = append(p.peek, token)
	return token
}

func (p *ParseContext) pop() *Token {
	end := len(p.peek) - 1
	token := p.peek[end]
	p.peek = p.peek[0:end]
	return token
}

func (p *ParseContext) String() string {
	return p.SelectedCommand.FullCommand()
}

func (p *ParseContext) matchedFlag(flag *FlagClause, value string) {
	p.Elements = append(p.Elements, &ParseElement{Clause: flag, Value: &value})
}

func (p *ParseContext) matchedArg(arg *ArgClause, value string) {
	p.Elements = append(p.Elements, &ParseElement{Clause: arg, Value: &value})
}

func (p *ParseContext) matchedCmd(cmd *CmdClause) {
	p.Elements = append(p.Elements, &ParseElement{Clause: cmd})
	p.mergeFlags(cmd.flagGroup)
	p.mergeArgs(cmd.argGroup)
	p.SelectedCommand = cmd
}


func ExpandArgsFromFile(filename string) (out []string, err error) {
	if filename == "" {
		return nil, fmt.Errorf("expected @ file to expand arguments from")
	}
	r, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open arguments file %q: %s", filename, err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to read arguments from %q: %s", filename, err)
	}
	return
}

func parse(context *ParseContext, app *Application) (err error) {
	context.mergeFlags(app.flagGroup)
	context.mergeArgs(app.argGroup)

	cmds := app.cmdGroup
	ignoreDefault := context.ignoreDefault

loop:
	for !context.EOL() && !context.Error() {
		token := context.Peek()

		switch token.Type {
		case TokenLong, TokenShort:
			if flag, err := context.flags.parse(context); err != nil {
				if !ignoreDefault {
					if cmd := cmds.defaultSubcommand(); cmd != nil {
						cmd.completionAlts = cmds.cmdNames()
						context.matchedCmd(cmd)
						cmds = cmd.cmdGroup
						break
					}
				}
				return err
			} else if flag == HelpFlag {
				ignoreDefault = true
			}

		case TokenArg:
			if cmds.have() {
				selectedDefault := false
				cmd, ok := cmds.commands[token.String()]
				if !ok {
					if !ignoreDefault {
						if cmd = cmds.defaultSubcommand(); cmd != nil {
							cmd.completionAlts = cmds.cmdNames()
							selectedDefault = true
						}
					}
					if cmd == nil {
						return fmt.Errorf("expected command but got %q", token)
					}
				}
				if cmd == HelpCommand {
					ignoreDefault = true
				}
				cmd.completionAlts = nil
				context.matchedCmd(cmd)
				cmds = cmd.cmdGroup
				if !selectedDefault {
					context.Next()
				}
			} else if context.arguments.have() {
				if app.noInterspersed {
					
					context.argsOnly = true
				}
				arg := context.nextArg()
				if arg == nil {
					break loop
				}
				context.matchedArg(arg, token.String())
				context.Next()
			} else {
				break loop
			}

		case TokenEOL:
			break loop
		}
	}

	
	for !ignoreDefault {
		if cmd := cmds.defaultSubcommand(); cmd != nil {
			cmd.completionAlts = cmds.cmdNames()
			context.matchedCmd(cmd)
			cmds = cmd.cmdGroup
		} else {
			break
		}
	}

	if context.Error() {
		return fmt.Errorf("%s", context.Peek().Value)
	}

	if !context.EOL() {
		return fmt.Errorf("unexpected %s", context.Peek())
	}

	
	for arg := context.nextArg(); arg != nil && !arg.consumesRemainder(); arg = context.nextArg() {
		for _, defaultValue := range arg.defaultValues {
			if err := arg.value.Set(defaultValue); err != nil {
				return fmt.Errorf("invalid default value '%s' for argument '%s'", defaultValue, arg.name)
			}
		}
	}

	return
}
