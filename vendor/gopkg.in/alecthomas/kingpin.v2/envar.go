package kingpin

import (
	"os"
	"regexp"
)

var (
	envVarValuesSeparator = "\r?\n"
	envVarValuesTrimmer   = regexp.MustCompile(envVarValuesSeparator + "$")
	envVarValuesSplitter  = regexp.MustCompile(envVarValuesSeparator)
)

type envarMixin struct {
	envar   string
	noEnvar bool
}

func (e *envarMixin) HasEnvarValue() bool {
	return e.GetEnvarValue() != ""
}

func (e *envarMixin) GetEnvarValue() string {
	if e.noEnvar || e.envar == "" {
		return ""
	}
	return os.Getenv(e.envar)
}

func (e *envarMixin) GetSplitEnvarValue() []string {
	values := make([]string, 0)

	envarValue := e.GetEnvarValue()
	if envarValue == "" {
		return values
	}

	
	trimmed := envVarValuesTrimmer.ReplaceAllString(envarValue, "")
	for _, value := range envVarValuesSplitter.Split(trimmed, -1) {
		values = append(values, value)
	}

	return values
}
