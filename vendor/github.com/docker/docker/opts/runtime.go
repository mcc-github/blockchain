package opts 

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
)


type RuntimeOpt struct {
	name             string
	stockRuntimeName string
	values           *map[string]types.Runtime
}


func NewNamedRuntimeOpt(name string, ref *map[string]types.Runtime, stockRuntime string) *RuntimeOpt {
	if ref == nil {
		ref = &map[string]types.Runtime{}
	}
	return &RuntimeOpt{name: name, values: ref, stockRuntimeName: stockRuntime}
}


func (o *RuntimeOpt) Name() string {
	return o.name
}


func (o *RuntimeOpt) Set(val string) error {
	parts := strings.SplitN(val, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid runtime argument: %s", val)
	}

	parts[0] = strings.TrimSpace(parts[0])
	parts[1] = strings.TrimSpace(parts[1])
	if parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid runtime argument: %s", val)
	}

	parts[0] = strings.ToLower(parts[0])
	if parts[0] == o.stockRuntimeName {
		return fmt.Errorf("runtime name '%s' is reserved", o.stockRuntimeName)
	}

	if _, ok := (*o.values)[parts[0]]; ok {
		return fmt.Errorf("runtime '%s' was already defined", parts[0])
	}

	(*o.values)[parts[0]] = types.Runtime{Path: parts[1]}

	return nil
}


func (o *RuntimeOpt) String() string {
	var out []string
	for k := range *o.values {
		out = append(out, k)
	}

	return fmt.Sprintf("%v", out)
}


func (o *RuntimeOpt) GetMap() map[string]types.Runtime {
	if o.values != nil {
		return *o.values
	}

	return map[string]types.Runtime{}
}


func (o *RuntimeOpt) Type() string {
	return "runtime"
}
