



package docker

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)


type Env []string


func (env *Env) Get(key string) (value string) {
	return env.Map()[key]
}



func (env *Env) Exists(key string) bool {
	_, exists := env.Map()[key]
	return exists
}




func (env *Env) GetBool(key string) (value bool) {
	s := strings.ToLower(strings.Trim(env.Get(key), " \t"))
	if s == "" || s == "0" || s == "no" || s == "false" || s == "none" {
		return false
	}
	return true
}


func (env *Env) SetBool(key string, value bool) {
	if value {
		env.Set(key, "1")
	} else {
		env.Set(key, "0")
	}
}




func (env *Env) GetInt(key string) int {
	return int(env.GetInt64(key))
}


func (env *Env) SetInt(key string, value int) {
	env.Set(key, strconv.Itoa(value))
}




func (env *Env) GetInt64(key string) int64 {
	s := strings.Trim(env.Get(key), " \t")
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return val
}


func (env *Env) SetInt64(key string, value int64) {
	env.Set(key, strconv.FormatInt(value, 10))
}




func (env *Env) GetJSON(key string, iface interface{}) error {
	sval := env.Get(key)
	if sval == "" {
		return nil
	}
	return json.Unmarshal([]byte(sval), iface)
}



func (env *Env) SetJSON(key string, value interface{}) error {
	sval, err := json.Marshal(value)
	if err != nil {
		return err
	}
	env.Set(key, string(sval))
	return nil
}






func (env *Env) GetList(key string) []string {
	sval := env.Get(key)
	if sval == "" {
		return nil
	}
	var l []string
	if err := json.Unmarshal([]byte(sval), &l); err != nil {
		l = append(l, sval)
	}
	return l
}



func (env *Env) SetList(key string, value []string) error {
	return env.SetJSON(key, value)
}


func (env *Env) Set(key, value string) {
	*env = append(*env, key+"="+value)
}





func (env *Env) Decode(src io.Reader) error {
	m := make(map[string]interface{})
	if err := json.NewDecoder(src).Decode(&m); err != nil {
		return err
	}
	for k, v := range m {
		env.SetAuto(k, v)
	}
	return nil
}


func (env *Env) SetAuto(key string, value interface{}) {
	if fval, ok := value.(float64); ok {
		env.SetInt64(key, int64(fval))
	} else if sval, ok := value.(string); ok {
		env.Set(key, sval)
	} else if val, err := json.Marshal(value); err == nil {
		env.Set(key, string(val))
	} else {
		env.Set(key, fmt.Sprintf("%v", value))
	}
}


func (env *Env) Map() map[string]string {
	if len(*env) == 0 {
		return nil
	}
	m := make(map[string]string)
	for _, kv := range *env {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 1 {
			m[parts[0]] = ""
		} else {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
