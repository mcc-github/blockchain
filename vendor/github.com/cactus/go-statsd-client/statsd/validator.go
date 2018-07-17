



package statsd

import (
	"fmt"
	"regexp"
)



type ValidatorFunc func(string) error

var safeName = regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)




func CheckName(stat string) error {
	if !safeName.MatchString(stat) {
		return fmt.Errorf("invalid stat name: %s", stat)
	}
	return nil
}
