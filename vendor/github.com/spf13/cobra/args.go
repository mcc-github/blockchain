package cobra

import (
	"fmt"
)

type PositionalArgs func(cmd *Command, args []string) error





func legacyArgs(cmd *Command, args []string) error {
	
	if !cmd.HasSubCommands() {
		return nil
	}

	
	if !cmd.HasParent() && len(args) > 0 {
		return fmt.Errorf("unknown command %q for %q%s", args[0], cmd.CommandPath(), cmd.findSuggestions(args[0]))
	}
	return nil
}


func NoArgs(cmd *Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
	}
	return nil
}


func OnlyValidArgs(cmd *Command, args []string) error {
	if len(cmd.ValidArgs) > 0 {
		for _, v := range args {
			if !stringInSlice(v, cmd.ValidArgs) {
				return fmt.Errorf("invalid argument %q for %q%s", v, cmd.CommandPath(), cmd.findSuggestions(args[0]))
			}
		}
	}
	return nil
}


func ArbitraryArgs(cmd *Command, args []string) error {
	return nil
}


func MinimumNArgs(n int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf("requires at least %d arg(s), only received %d", n, len(args))
		}
		return nil
	}
}


func MaximumNArgs(n int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) > n {
			return fmt.Errorf("accepts at most %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}


func ExactArgs(n int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}


func RangeArgs(min int, max int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) < min || len(args) > max {
			return fmt.Errorf("accepts between %d and %d arg(s), received %d", min, max, len(args))
		}
		return nil
	}
}
