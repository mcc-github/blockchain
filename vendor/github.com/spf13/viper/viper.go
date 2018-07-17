


















package viper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
)

var v *Viper

func init() {
	v = New()
}

type remoteConfigFactory interface {
	Get(rp RemoteProvider) (io.Reader, error)
	Watch(rp RemoteProvider) (io.Reader, error)
}


var RemoteConfig remoteConfigFactory



type UnsupportedConfigError string


func (str UnsupportedConfigError) Error() string {
	return fmt.Sprintf("Unsupported Config Type %q", string(str))
}




type UnsupportedRemoteProviderError string


func (str UnsupportedRemoteProviderError) Error() string {
	return fmt.Sprintf("Unsupported Remote Provider Type %q", string(str))
}



type RemoteConfigError string


func (rce RemoteConfigError) Error() string {
	return fmt.Sprintf("Remote Configurations Error: %s", string(rce))
}


type ConfigFileNotFoundError struct {
	name, locations string
}


func (fnfe ConfigFileNotFoundError) Error() string {
	return fmt.Sprintf("Config File %q Not Found in %q", fnfe.name, fnfe.locations)
}



































type Viper struct {
	
	
	keyDelim string

	
	configPaths []string

	
	remoteProviders []*defaultRemoteProvider

	
	configName string
	configFile string
	configType string
	envPrefix  string

	automaticEnvApplied bool
	envKeyReplacer      *strings.Replacer

	config         map[string]interface{}
	override       map[string]interface{}
	defaults       map[string]interface{}
	kvstore        map[string]interface{}
	pflags         map[string]*pflag.Flag
	env            map[string]string
	aliases        map[string]string
	typeByDefValue bool
}


func New() *Viper {
	v := new(Viper)
	v.keyDelim = "."
	v.configName = "config"
	v.config = make(map[string]interface{})
	v.override = make(map[string]interface{})
	v.defaults = make(map[string]interface{})
	v.kvstore = make(map[string]interface{})
	v.pflags = make(map[string]*pflag.Flag)
	v.env = make(map[string]string)
	v.aliases = make(map[string]string)
	v.typeByDefValue = false

	return v
}




func Reset() {
	v = New()
	SupportedExts = []string{"json", "toml", "yaml", "yml"}
	SupportedRemoteProviders = []string{"etcd", "consul"}
}

type defaultRemoteProvider struct {
	provider      string
	endpoint      string
	path          string
	secretKeyring string
}

func (rp defaultRemoteProvider) Provider() string {
	return rp.provider
}

func (rp defaultRemoteProvider) Endpoint() string {
	return rp.endpoint
}

func (rp defaultRemoteProvider) Path() string {
	return rp.path
}

func (rp defaultRemoteProvider) SecretKeyring() string {
	return rp.secretKeyring
}





type RemoteProvider interface {
	Provider() string
	Endpoint() string
	Path() string
	SecretKeyring() string
}


var SupportedExts []string = []string{"json", "toml", "yaml", "yml", "properties", "props", "prop"}


var SupportedRemoteProviders []string = []string{"etcd", "consul"}



func SetConfigFile(in string) { v.SetConfigFile(in) }
func (v *Viper) SetConfigFile(in string) {
	if in != "" {
		v.configFile = in
	}
}




func SetEnvPrefix(in string) { v.SetEnvPrefix(in) }
func (v *Viper) SetEnvPrefix(in string) {
	if in != "" {
		v.envPrefix = in
	}
}

func (v *Viper) mergeWithEnvPrefix(in string) string {
	if v.envPrefix != "" {
		return strings.ToUpper(v.envPrefix + "_" + in)
	}

	return strings.ToUpper(in)
}








func (v *Viper) getEnv(key string) string {
	if v.envKeyReplacer != nil {
		key = v.envKeyReplacer.Replace(key)
	}
	return os.Getenv(key)
}


func ConfigFileUsed() string            { return v.ConfigFileUsed() }
func (v *Viper) ConfigFileUsed() string { return v.configFile }



func AddConfigPath(in string) { v.AddConfigPath(in) }
func (v *Viper) AddConfigPath(in string) {
	if in != "" {
		absin := absPathify(in)
		jww.INFO.Println("adding", absin, "to paths to search")
		if !stringInSlice(absin, v.configPaths) {
			v.configPaths = append(v.configPaths, absin)
		}
	}
}









func AddRemoteProvider(provider, endpoint, path string) error {
	return v.AddRemoteProvider(provider, endpoint, path)
}
func (v *Viper) AddRemoteProvider(provider, endpoint, path string) error {
	if !stringInSlice(provider, SupportedRemoteProviders) {
		return UnsupportedRemoteProviderError(provider)
	}
	if provider != "" && endpoint != "" {
		jww.INFO.Printf("adding %s:%s to remote provider list", provider, endpoint)
		rp := &defaultRemoteProvider{
			endpoint: endpoint,
			provider: provider,
			path:     path,
		}
		if !v.providerPathExists(rp) {
			v.remoteProviders = append(v.remoteProviders, rp)
		}
	}
	return nil
}











func AddSecureRemoteProvider(provider, endpoint, path, secretkeyring string) error {
	return v.AddSecureRemoteProvider(provider, endpoint, path, secretkeyring)
}

func (v *Viper) AddSecureRemoteProvider(provider, endpoint, path, secretkeyring string) error {
	if !stringInSlice(provider, SupportedRemoteProviders) {
		return UnsupportedRemoteProviderError(provider)
	}
	if provider != "" && endpoint != "" {
		jww.INFO.Printf("adding %s:%s to remote provider list", provider, endpoint)
		rp := &defaultRemoteProvider{
			endpoint:      endpoint,
			provider:      provider,
			path:          path,
			secretKeyring: secretkeyring,
		}
		if !v.providerPathExists(rp) {
			v.remoteProviders = append(v.remoteProviders, rp)
		}
	}
	return nil
}

func (v *Viper) providerPathExists(p *defaultRemoteProvider) bool {
	for _, y := range v.remoteProviders {
		if reflect.DeepEqual(y, p) {
			return true
		}
	}
	return false
}

func (v *Viper) searchMap(source map[string]interface{}, path []string) interface{} {

	if len(path) == 0 {
		return source
	}

	if next, ok := source[path[0]]; ok {
		switch next.(type) {
		case map[interface{}]interface{}:
			return v.searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			
			
			return v.searchMap(next.(map[string]interface{}), path[1:])
		default:
			return next
		}
	} else {
		return nil
	}
}















func SetTypeByDefaultValue(enable bool) { v.SetTypeByDefaultValue(enable) }
func (v *Viper) SetTypeByDefaultValue(enable bool) {
	v.typeByDefValue = enable
}








func Get(key string) interface{} { return v.Get(key) }
func (v *Viper) Get(key string) interface{} {
	path := strings.Split(key, v.keyDelim)

	lcaseKey := strings.ToLower(key)
	val := v.find(lcaseKey)

	if val == nil {
		source := v.find(path[0])
		if source == nil {
			return nil
		}

		if reflect.TypeOf(source).Kind() == reflect.Map {
			val = v.searchMap(cast.ToStringMap(source), path[1:])
		}
	}

	var valType interface{}
	if !v.typeByDefValue {
		valType = val
	} else {
		defVal, defExists := v.defaults[lcaseKey]
		if defExists {
			valType = defVal
		} else {
			valType = val
		}
	}

	switch valType.(type) {
	case bool:
		return cast.ToBool(val)
	case string:
		return cast.ToString(val)
	case int64, int32, int16, int8, int:
		return cast.ToInt(val)
	case float64, float32:
		return cast.ToFloat64(val)
	case time.Time:
		return cast.ToTime(val)
	case time.Duration:
		return cast.ToDuration(val)
	case []string:
		return cast.ToStringSlice(val)
	}
	return val
}


func GetString(key string) string { return v.GetString(key) }
func (v *Viper) GetString(key string) string {
	return cast.ToString(v.Get(key))
}


func GetBool(key string) bool { return v.GetBool(key) }
func (v *Viper) GetBool(key string) bool {
	return cast.ToBool(v.Get(key))
}


func GetInt(key string) int { return v.GetInt(key) }
func (v *Viper) GetInt(key string) int {
	return cast.ToInt(v.Get(key))
}


func GetFloat64(key string) float64 { return v.GetFloat64(key) }
func (v *Viper) GetFloat64(key string) float64 {
	return cast.ToFloat64(v.Get(key))
}


func GetTime(key string) time.Time { return v.GetTime(key) }
func (v *Viper) GetTime(key string) time.Time {
	return cast.ToTime(v.Get(key))
}


func GetDuration(key string) time.Duration { return v.GetDuration(key) }
func (v *Viper) GetDuration(key string) time.Duration {
	return cast.ToDuration(v.Get(key))
}


func GetStringSlice(key string) []string { return v.GetStringSlice(key) }
func (v *Viper) GetStringSlice(key string) []string {
	return cast.ToStringSlice(v.Get(key))
}


func GetStringMap(key string) map[string]interface{} { return v.GetStringMap(key) }
func (v *Viper) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(v.Get(key))
}


func GetStringMapString(key string) map[string]string { return v.GetStringMapString(key) }
func (v *Viper) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(v.Get(key))
}


func GetStringMapStringSlice(key string) map[string][]string { return v.GetStringMapStringSlice(key) }
func (v *Viper) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(v.Get(key))
}



func GetSizeInBytes(key string) uint { return v.GetSizeInBytes(key) }
func (v *Viper) GetSizeInBytes(key string) uint {
	sizeStr := cast.ToString(v.Get(key))
	return parseSizeInBytes(sizeStr)
}


func UnmarshalKey(key string, rawVal interface{}) error { return v.UnmarshalKey(key, rawVal) }
func (v *Viper) UnmarshalKey(key string, rawVal interface{}) error {
	return mapstructure.Decode(v.Get(key), rawVal)
}



func Unmarshal(rawVal interface{}) error { return v.Unmarshal(rawVal) }
func (v *Viper) Unmarshal(rawVal interface{}) error {
	err := mapstructure.WeakDecode(v.AllSettings(), rawVal)

	if err != nil {
		return err
	}

	v.insensitiviseMaps()

	return nil
}



func BindPFlags(flags *pflag.FlagSet) (err error) { return v.BindPFlags(flags) }
func (v *Viper) BindPFlags(flags *pflag.FlagSet) (err error) {
	flags.VisitAll(func(flag *pflag.Flag) {
		if err != nil {
			
			return
		}

		err = v.BindPFlag(flag.Name, flag)
		switch flag.Value.Type() {
		case "int", "int8", "int16", "int32", "int64":
			v.SetDefault(flag.Name, cast.ToInt(flag.Value.String()))
		case "bool":
			v.SetDefault(flag.Name, cast.ToBool(flag.Value.String()))
		default:
			v.SetDefault(flag.Name, flag.Value.String())
		}
	})
	return
}







func BindPFlag(key string, flag *pflag.Flag) (err error) { return v.BindPFlag(key, flag) }
func (v *Viper) BindPFlag(key string, flag *pflag.Flag) (err error) {
	if flag == nil {
		return fmt.Errorf("flag for %q is nil", key)
	}
	v.pflags[strings.ToLower(key)] = flag

	switch flag.Value.Type() {
	case "int", "int8", "int16", "int32", "int64":
		v.SetDefault(key, cast.ToInt(flag.Value.String()))
	case "bool":
		v.SetDefault(key, cast.ToBool(flag.Value.String()))
	default:
		v.SetDefault(key, flag.Value.String())
	}
	return nil
}





func BindEnv(input ...string) (err error) { return v.BindEnv(input...) }
func (v *Viper) BindEnv(input ...string) (err error) {
	var key, envkey string
	if len(input) == 0 {
		return fmt.Errorf("BindEnv missing key to bind to")
	}

	key = strings.ToLower(input[0])

	if len(input) == 1 {
		envkey = v.mergeWithEnvPrefix(key)
	} else {
		envkey = input[1]
	}

	v.env[key] = envkey

	return nil
}





func (v *Viper) find(key string) interface{} {
	var val interface{}
	var exists bool

	
	key = v.realKey(key)

	
	flag, exists := v.pflags[key]
	if exists {
		if flag.Changed {
			jww.TRACE.Println(key, "found in override (via pflag):", val)
			return flag.Value.String()
		}
	}

	val, exists = v.override[key]
	if exists {
		jww.TRACE.Println(key, "found in override:", val)
		return val
	}

	if v.automaticEnvApplied {
		
		
		if val = v.getEnv(v.mergeWithEnvPrefix(key)); val != "" {
			jww.TRACE.Println(key, "found in environment with val:", val)
			return val
		}
	}

	envkey, exists := v.env[key]
	if exists {
		jww.TRACE.Println(key, "registered as env var", envkey)
		if val = v.getEnv(envkey); val != "" {
			jww.TRACE.Println(envkey, "found in environment with val:", val)
			return val
		} else {
			jww.TRACE.Println(envkey, "env value unset:")
		}
	}

	val, exists = v.config[key]
	if exists {
		jww.TRACE.Println(key, "found in config:", val)
		return val
	}

	val, exists = v.kvstore[key]
	if exists {
		jww.TRACE.Println(key, "found in key/value store:", val)
		return val
	}

	val, exists = v.defaults[key]
	if exists {
		jww.TRACE.Println(key, "found in defaults:", val)
		return val
	}

	return nil
}


func IsSet(key string) bool { return v.IsSet(key) }
func (v *Viper) IsSet(key string) bool {
	t := v.Get(key)
	return t != nil
}



func AutomaticEnv() { v.AutomaticEnv() }
func (v *Viper) AutomaticEnv() {
	v.automaticEnvApplied = true
}




func SetEnvKeyReplacer(r *strings.Replacer) { v.SetEnvKeyReplacer(r) }
func (v *Viper) SetEnvKeyReplacer(r *strings.Replacer) {
	v.envKeyReplacer = r
}



func RegisterAlias(alias string, key string) { v.RegisterAlias(alias, key) }
func (v *Viper) RegisterAlias(alias string, key string) {
	v.registerAlias(alias, strings.ToLower(key))
}

func (v *Viper) registerAlias(alias string, key string) {
	alias = strings.ToLower(alias)
	if alias != key && alias != v.realKey(key) {
		_, exists := v.aliases[alias]

		if !exists {
			
			
			
			if val, ok := v.config[alias]; ok {
				delete(v.config, alias)
				v.config[key] = val
			}
			if val, ok := v.kvstore[alias]; ok {
				delete(v.kvstore, alias)
				v.kvstore[key] = val
			}
			if val, ok := v.defaults[alias]; ok {
				delete(v.defaults, alias)
				v.defaults[key] = val
			}
			if val, ok := v.override[alias]; ok {
				delete(v.override, alias)
				v.override[key] = val
			}
			v.aliases[alias] = key
		}
	} else {
		jww.WARN.Println("Creating circular reference alias", alias, key, v.realKey(key))
	}
}

func (v *Viper) realKey(key string) string {
	newkey, exists := v.aliases[key]
	if exists {
		jww.DEBUG.Println("Alias", key, "to", newkey)
		return v.realKey(newkey)
	} else {
		return key
	}
}


func InConfig(key string) bool { return v.InConfig(key) }
func (v *Viper) InConfig(key string) bool {
	
	key = v.realKey(key)

	_, exists := v.config[key]
	return exists
}



func SetDefault(key string, value interface{}) { v.SetDefault(key, value) }
func (v *Viper) SetDefault(key string, value interface{}) {
	
	key = v.realKey(strings.ToLower(key))
	v.defaults[key] = value
}




func Set(key string, value interface{}) { v.Set(key, value) }
func (v *Viper) Set(key string, value interface{}) {
	
	key = v.realKey(strings.ToLower(key))
	v.override[key] = value
}



func ReadInConfig() error { return v.ReadInConfig() }
func (v *Viper) ReadInConfig() error {
	jww.INFO.Println("Attempting to read in config file")
	if !stringInSlice(v.getConfigType(), SupportedExts) {
		return UnsupportedConfigError(v.getConfigType())
	}

	file, err := ioutil.ReadFile(v.getConfigFile())
	if err != nil {
		return err
	}

	v.config = make(map[string]interface{})

	return v.unmarshalReader(bytes.NewReader(file), v.config)
}

func ReadConfig(in io.Reader) error { return v.ReadConfig(in) }
func (v *Viper) ReadConfig(in io.Reader) error {
	v.config = make(map[string]interface{})
	return v.unmarshalReader(in, v.config)
}









func ReadRemoteConfig() error { return v.ReadRemoteConfig() }
func (v *Viper) ReadRemoteConfig() error {
	err := v.getKeyValueConfig()
	if err != nil {
		return err
	}
	return nil
}

func WatchRemoteConfig() error { return v.WatchRemoteConfig() }
func (v *Viper) WatchRemoteConfig() error {
	err := v.watchKeyValueConfig()
	if err != nil {
		return err
	}
	return nil
}



func unmarshalReader(in io.Reader, c map[string]interface{}) error {
	return v.unmarshalReader(in, c)
}

func (v *Viper) unmarshalReader(in io.Reader, c map[string]interface{}) error {
	return unmarshallConfigReader(in, c, v.getConfigType())
}

func (v *Viper) insensitiviseMaps() {
	insensitiviseMap(v.config)
	insensitiviseMap(v.defaults)
	insensitiviseMap(v.override)
	insensitiviseMap(v.kvstore)
}


func (v *Viper) getKeyValueConfig() error {
	if RemoteConfig == nil {
		return RemoteConfigError("Enable the remote features by doing a blank import of the viper/remote package: '_ github.com/spf13/viper/remote'")
	}

	for _, rp := range v.remoteProviders {
		val, err := v.getRemoteConfig(rp)
		if err != nil {
			continue
		}
		v.kvstore = val
		return nil
	}
	return RemoteConfigError("No Files Found")
}

func (v *Viper) getRemoteConfig(provider *defaultRemoteProvider) (map[string]interface{}, error) {

	reader, err := RemoteConfig.Get(provider)
	if err != nil {
		return nil, err
	}
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
}


func (v *Viper) watchKeyValueConfig() error {
	for _, rp := range v.remoteProviders {
		val, err := v.watchRemoteConfig(rp)
		if err != nil {
			continue
		}
		v.kvstore = val
		return nil
	}
	return RemoteConfigError("No Files Found")
}

func (v *Viper) watchRemoteConfig(provider *defaultRemoteProvider) (map[string]interface{}, error) {
	reader, err := RemoteConfig.Watch(provider)
	if err != nil {
		return nil, err
	}
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
}


func AllKeys() []string { return v.AllKeys() }
func (v *Viper) AllKeys() []string {
	m := map[string]struct{}{}

	for key, _ := range v.defaults {
		m[key] = struct{}{}
	}

	for key, _ := range v.config {
		m[key] = struct{}{}
	}

	for key, _ := range v.kvstore {
		m[key] = struct{}{}
	}

	for key, _ := range v.override {
		m[key] = struct{}{}
	}

	a := []string{}
	for x, _ := range m {
		a = append(a, x)
	}

	return a
}


func AllSettings() map[string]interface{} { return v.AllSettings() }
func (v *Viper) AllSettings() map[string]interface{} {
	m := map[string]interface{}{}
	for _, x := range v.AllKeys() {
		m[x] = v.Get(x)
	}

	return m
}



func SetConfigName(in string) { v.SetConfigName(in) }
func (v *Viper) SetConfigName(in string) {
	if in != "" {
		v.configName = in
	}
}



func SetConfigType(in string) { v.SetConfigType(in) }
func (v *Viper) SetConfigType(in string) {
	if in != "" {
		v.configType = in
	}
}

func (v *Viper) getConfigType() string {
	if v.configType != "" {
		return v.configType
	}

	cf := v.getConfigFile()
	ext := filepath.Ext(cf)

	if len(ext) > 1 {
		return ext[1:]
	} else {
		return ""
	}
}

func (v *Viper) getConfigFile() string {
	
	if v.configFile != "" {
		return v.configFile
	}

	cf, err := v.findConfigFile()
	if err != nil {
		return ""
	}

	v.configFile = cf
	return v.getConfigFile()
}

func (v *Viper) searchInPath(in string) (filename string) {
	jww.DEBUG.Println("Searching for config in ", in)
	for _, ext := range SupportedExts {
		jww.DEBUG.Println("Checking for", filepath.Join(in, v.configName+"."+ext))
		if b, _ := exists(filepath.Join(in, v.configName+"."+ext)); b {
			jww.DEBUG.Println("Found: ", filepath.Join(in, v.configName+"."+ext))
			return filepath.Join(in, v.configName+"."+ext)
		}
	}

	return ""
}



func (v *Viper) findConfigFile() (string, error) {

	jww.INFO.Println("Searching for config in ", v.configPaths)

	for _, cp := range v.configPaths {
		file := v.searchInPath(cp)
		if file != "" {
			return file, nil
		}
	}
	return "", ConfigFileNotFoundError{v.configName, fmt.Sprintf("%s", v.configPaths)}
}



func Debug() { v.Debug() }
func (v *Viper) Debug() {
	fmt.Println("Aliases:")
	pretty.Println(v.aliases)
	fmt.Println("Override:")
	pretty.Println(v.override)
	fmt.Println("PFlags")
	pretty.Println(v.pflags)
	fmt.Println("Env:")
	pretty.Println(v.env)
	fmt.Println("Key/Value Store:")
	pretty.Println(v.kvstore)
	fmt.Println("Config:")
	pretty.Println(v.config)
	fmt.Println("Defaults:")
	pretty.Println(v.defaults)
}
