package kingpin

import (
	"fmt"
	"strings"
)

type flagGroup struct {
	short     map[string]*FlagClause
	long      map[string]*FlagClause
	flagOrder []*FlagClause
}

func newFlagGroup() *flagGroup {
	return &flagGroup{
		short: map[string]*FlagClause{},
		long:  map[string]*FlagClause{},
	}
}





func (f *flagGroup) GetFlag(name string) *FlagClause {
	return f.long[name]
}


func (f *flagGroup) Flag(name, help string) *FlagClause {
	flag := newFlag(name, help)
	f.long[name] = flag
	f.flagOrder = append(f.flagOrder, flag)
	return flag
}

func (f *flagGroup) init(defaultEnvarPrefix string) error {
	if err := f.checkDuplicates(); err != nil {
		return err
	}
	for _, flag := range f.long {
		if defaultEnvarPrefix != "" && !flag.noEnvar && flag.envar == "" {
			flag.envar = envarTransform(defaultEnvarPrefix + "_" + flag.name)
		}
		if err := flag.init(); err != nil {
			return err
		}
		if flag.shorthand != 0 {
			f.short[string(flag.shorthand)] = flag
		}
	}
	return nil
}

func (f *flagGroup) checkDuplicates() error {
	seenShort := map[rune]bool{}
	seenLong := map[string]bool{}
	for _, flag := range f.flagOrder {
		if flag.shorthand != 0 {
			if _, ok := seenShort[flag.shorthand]; ok {
				return fmt.Errorf("duplicate short flag -%c", flag.shorthand)
			}
			seenShort[flag.shorthand] = true
		}
		if _, ok := seenLong[flag.name]; ok {
			return fmt.Errorf("duplicate long flag --%s", flag.name)
		}
		seenLong[flag.name] = true
	}
	return nil
}

func (f *flagGroup) parse(context *ParseContext) (*FlagClause, error) {
	var token *Token

loop:
	for {
		token = context.Peek()
		switch token.Type {
		case TokenEOL:
			break loop

		case TokenLong, TokenShort:
			flagToken := token
			defaultValue := ""
			var flag *FlagClause
			var ok bool
			invert := false

			name := token.Value
			if token.Type == TokenLong {
				flag, ok = f.long[name]
				if !ok {
					if strings.HasPrefix(name, "no-") {
						name = name[3:]
						invert = true
					}
					flag, ok = f.long[name]
				}
				if !ok {
					return nil, fmt.Errorf("unknown long flag '%s'", flagToken)
				}
			} else {
				flag, ok = f.short[name]
				if !ok {
					return nil, fmt.Errorf("unknown short flag '%s'", flagToken)
				}
			}

			context.Next()

			fb, ok := flag.value.(boolFlag)
			if ok && fb.IsBoolFlag() {
				if invert {
					defaultValue = "false"
				} else {
					defaultValue = "true"
				}
			} else {
				if invert {
					context.Push(token)
					return nil, fmt.Errorf("unknown long flag '%s'", flagToken)
				}
				token = context.Peek()
				if token.Type != TokenArg {
					context.Push(token)
					return nil, fmt.Errorf("expected argument for flag '%s'", flagToken)
				}
				context.Next()
				defaultValue = token.Value
			}

			context.matchedFlag(flag, defaultValue)
			return flag, nil

		default:
			break loop
		}
	}
	return nil, nil
}


type FlagClause struct {
	parserMixin
	actionMixin
	completionsMixin
	envarMixin
	name          string
	shorthand     rune
	help          string
	defaultValues []string
	placeholder   string
	hidden        bool
}

func newFlag(name, help string) *FlagClause {
	f := &FlagClause{
		name: name,
		help: help,
	}
	return f
}

func (f *FlagClause) setDefault() error {
	if f.HasEnvarValue() {
		if v, ok := f.value.(repeatableFlag); !ok || !v.IsCumulative() {
			
			return f.value.Set(f.GetEnvarValue())
		} else {
			for _, value := range f.GetSplitEnvarValue() {
				if err := f.value.Set(value); err != nil {
					return err
				}
			}
			return nil
		}
	}

	if len(f.defaultValues) > 0 {
		for _, defaultValue := range f.defaultValues {
			if err := f.value.Set(defaultValue); err != nil {
				return err
			}
		}
		return nil
	}

	return nil
}

func (f *FlagClause) needsValue() bool {
	haveDefault := len(f.defaultValues) > 0
	return f.required && !(haveDefault || f.HasEnvarValue())
}

func (f *FlagClause) init() error {
	if f.required && len(f.defaultValues) > 0 {
		return fmt.Errorf("required flag '--%s' with default value that will never be used", f.name)
	}
	if f.value == nil {
		return fmt.Errorf("no type defined for --%s (eg. .String())", f.name)
	}
	if v, ok := f.value.(repeatableFlag); (!ok || !v.IsCumulative()) && len(f.defaultValues) > 1 {
		return fmt.Errorf("invalid default for '--%s', expecting single value", f.name)
	}
	return nil
}


func (f *FlagClause) Action(action Action) *FlagClause {
	f.addAction(action)
	return f
}

func (f *FlagClause) PreAction(action Action) *FlagClause {
	f.addPreAction(action)
	return f
}


func (a *FlagClause) HintAction(action HintAction) *FlagClause {
	a.addHintAction(action)
	return a
}


func (a *FlagClause) HintOptions(options ...string) *FlagClause {
	a.addHintAction(func() []string {
		return options
	})
	return a
}

func (a *FlagClause) EnumVar(target *string, options ...string) {
	a.parserMixin.EnumVar(target, options...)
	a.addHintActionBuiltin(func() []string {
		return options
	})
}

func (a *FlagClause) Enum(options ...string) (target *string) {
	a.addHintActionBuiltin(func() []string {
		return options
	})
	return a.parserMixin.Enum(options...)
}


func (f *FlagClause) Default(values ...string) *FlagClause {
	f.defaultValues = values
	return f
}


func (f *FlagClause) OverrideDefaultFromEnvar(envar string) *FlagClause {
	return f.Envar(envar)
}




func (f *FlagClause) Envar(name string) *FlagClause {
	f.envar = name
	f.noEnvar = false
	return f
}



func (f *FlagClause) NoEnvar() *FlagClause {
	f.envar = ""
	f.noEnvar = true
	return f
}




func (f *FlagClause) PlaceHolder(placeholder string) *FlagClause {
	f.placeholder = placeholder
	return f
}


func (f *FlagClause) Hidden() *FlagClause {
	f.hidden = true
	return f
}


func (f *FlagClause) Required() *FlagClause {
	f.required = true
	return f
}


func (f *FlagClause) Short(name rune) *FlagClause {
	f.shorthand = name
	return f
}


func (f *FlagClause) Bool() (target *bool) {
	target = new(bool)
	f.SetValue(newBoolValue(target))
	return
}
