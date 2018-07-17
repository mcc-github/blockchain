package types 

import (
	"encoding/json"
	"fmt"
	"sort"
)


type PluginsListResponse []*Plugin


func (t *PluginInterfaceType) UnmarshalJSON(p []byte) error {
	versionIndex := len(p)
	prefixIndex := 0
	if len(p) < 2 || p[0] != '"' || p[len(p)-1] != '"' {
		return fmt.Errorf("%q is not a plugin interface type", p)
	}
	p = p[1 : len(p)-1]
loop:
	for i, b := range p {
		switch b {
		case '.':
			prefixIndex = i
		case '/':
			versionIndex = i
			break loop
		}
	}
	t.Prefix = string(p[:prefixIndex])
	t.Capability = string(p[prefixIndex+1 : versionIndex])
	if versionIndex < len(p) {
		t.Version = string(p[versionIndex+1:])
	}
	return nil
}


func (t *PluginInterfaceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}


func (t PluginInterfaceType) String() string {
	return fmt.Sprintf("%s.%s/%s", t.Prefix, t.Capability, t.Version)
}



type PluginPrivilege struct {
	Name        string
	Description string
	Value       []string
}


type PluginPrivileges []PluginPrivilege

func (s PluginPrivileges) Len() int {
	return len(s)
}

func (s PluginPrivileges) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s PluginPrivileges) Swap(i, j int) {
	sort.Strings(s[i].Value)
	sort.Strings(s[j].Value)
	s[i], s[j] = s[j], s[i]
}
