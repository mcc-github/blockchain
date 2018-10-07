



package docker

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)


type PluginPrivilege struct {
	Name        string   `json:"Name,omitempty" yaml:"Name,omitempty" toml:"Name,omitempty"`
	Description string   `json:"Description,omitempty" yaml:"Description,omitempty" toml:"Description,omitempty"`
	Value       []string `json:"Value,omitempty" yaml:"Value,omitempty" toml:"Value,omitempty"`
}




type InstallPluginOptions struct {
	Remote  string
	Name    string
	Plugins []PluginPrivilege `qs:"-"`

	Auth AuthConfiguration

	Context context.Context
}




func (c *Client) InstallPlugins(opts InstallPluginOptions) error {
	path := "/plugins/pull?" + queryString(opts)
	resp, err := c.do("POST", path, doOptions{
		data:    opts.Plugins,
		context: opts.Context,
	})
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}




type PluginSettings struct {
	Env     []string `json:"Env,omitempty" yaml:"Env,omitempty" toml:"Env,omitempty"`
	Args    []string `json:"Args,omitempty" yaml:"Args,omitempty" toml:"Args,omitempty"`
	Devices []string `json:"Devices,omitempty" yaml:"Devices,omitempty" toml:"Devices,omitempty"`
}




type PluginInterface struct {
	Types  []string `json:"Types,omitempty" yaml:"Types,omitempty" toml:"Types,omitempty"`
	Socket string   `json:"Socket,omitempty" yaml:"Socket,omitempty" toml:"Socket,omitempty"`
}




type PluginNetwork struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty" toml:"Type,omitempty"`
}




type PluginLinux struct {
	Capabilities    []string             `json:"Capabilities,omitempty" yaml:"Capabilities,omitempty" toml:"Capabilities,omitempty"`
	AllowAllDevices bool                 `json:"AllowAllDevices,omitempty" yaml:"AllowAllDevices,omitempty" toml:"AllowAllDevices,omitempty"`
	Devices         []PluginLinuxDevices `json:"Devices,omitempty" yaml:"Devices,omitempty" toml:"Devices,omitempty"`
}




type PluginLinuxDevices struct {
	Name        string   `json:"Name,omitempty" yaml:"Name,omitempty" toml:"Name,omitempty"`
	Description string   `json:"Documentation,omitempty" yaml:"Documentation,omitempty" toml:"Documentation,omitempty"`
	Settable    []string `json:"Settable,omitempty" yaml:"Settable,omitempty" toml:"Settable,omitempty"`
	Path        string   `json:"Path,omitempty" yaml:"Path,omitempty" toml:"Path,omitempty"`
}




type PluginEnv struct {
	Name        string   `json:"Name,omitempty" yaml:"Name,omitempty" toml:"Name,omitempty"`
	Description string   `json:"Description,omitempty" yaml:"Description,omitempty" toml:"Description,omitempty"`
	Settable    []string `json:"Settable,omitempty" yaml:"Settable,omitempty" toml:"Settable,omitempty"`
	Value       string   `json:"Value,omitempty" yaml:"Value,omitempty" toml:"Value,omitempty"`
}




type PluginArgs struct {
	Name        string   `json:"Name,omitempty" yaml:"Name,omitempty" toml:"Name,omitempty"`
	Description string   `json:"Description,omitempty" yaml:"Description,omitempty" toml:"Description,omitempty"`
	Settable    []string `json:"Settable,omitempty" yaml:"Settable,omitempty" toml:"Settable,omitempty"`
	Value       []string `json:"Value,omitempty" yaml:"Value,omitempty" toml:"Value,omitempty"`
}




type PluginUser struct {
	UID int32 `json:"UID,omitempty" yaml:"UID,omitempty" toml:"UID,omitempty"`
	GID int32 `json:"GID,omitempty" yaml:"GID,omitempty" toml:"GID,omitempty"`
}




type PluginConfig struct {
	Description     string `json:"Description,omitempty" yaml:"Description,omitempty" toml:"Description,omitempty"`
	Documentation   string
	Interface       PluginInterface `json:"Interface,omitempty" yaml:"Interface,omitempty" toml:"Interface,omitempty"`
	Entrypoint      []string        `json:"Entrypoint,omitempty" yaml:"Entrypoint,omitempty" toml:"Entrypoint,omitempty"`
	WorkDir         string          `json:"WorkDir,omitempty" yaml:"WorkDir,omitempty" toml:"WorkDir,omitempty"`
	User            PluginUser      `json:"User,omitempty" yaml:"User,omitempty" toml:"User,omitempty"`
	Network         PluginNetwork   `json:"Network,omitempty" yaml:"Network,omitempty" toml:"Network,omitempty"`
	Linux           PluginLinux     `json:"Linux,omitempty" yaml:"Linux,omitempty" toml:"Linux,omitempty"`
	PropagatedMount string          `json:"PropagatedMount,omitempty" yaml:"PropagatedMount,omitempty" toml:"PropagatedMount,omitempty"`
	Mounts          []Mount         `json:"Mounts,omitempty" yaml:"Mounts,omitempty" toml:"Mounts,omitempty"`
	Env             []PluginEnv     `json:"Env,omitempty" yaml:"Env,omitempty" toml:"Env,omitempty"`
	Args            PluginArgs      `json:"Args,omitempty" yaml:"Args,omitempty" toml:"Args,omitempty"`
}




type PluginDetail struct {
	ID       string         `json:"Id,omitempty" yaml:"Id,omitempty" toml:"Id,omitempty"`
	Name     string         `json:"Name,omitempty" yaml:"Name,omitempty" toml:"Name,omitempty"`
	Tag      string         `json:"Tag,omitempty" yaml:"Tag,omitempty" toml:"Tag,omitempty"`
	Active   bool           `json:"Enabled,omitempty" yaml:"Active,omitempty" toml:"Active,omitempty"`
	Settings PluginSettings `json:"Settings,omitempty" yaml:"Settings,omitempty" toml:"Settings,omitempty"`
	Config   PluginConfig   `json:"Config,omitempty" yaml:"Config,omitempty" toml:"Config,omitempty"`
}




func (c *Client) ListPlugins(ctx context.Context) ([]PluginDetail, error) {
	resp, err := c.do("GET", "/plugins", doOptions{
		context: ctx,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	pluginDetails := make([]PluginDetail, 0)
	if err := json.NewDecoder(resp.Body).Decode(&pluginDetails); err != nil {
		return nil, err
	}
	return pluginDetails, nil
}




type ListFilteredPluginsOptions struct {
	Filters map[string][]string
	Context context.Context
}




func (c *Client) ListFilteredPlugins(opts ListFilteredPluginsOptions) ([]PluginDetail, error) {
	path := "/plugins/json?" + queryString(opts)
	resp, err := c.do("GET", path, doOptions{
		context: opts.Context,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	pluginDetails := make([]PluginDetail, 0)
	if err := json.NewDecoder(resp.Body).Decode(&pluginDetails); err != nil {
		return nil, err
	}
	return pluginDetails, nil
}




func (c *Client) GetPluginPrivileges(name string, ctx context.Context) ([]PluginPrivilege, error) {
	resp, err := c.do("GET", "/plugins/privileges?remote="+name, doOptions{
		context: ctx,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var pluginPrivileges []PluginPrivilege
	if err := json.NewDecoder(resp.Body).Decode(&pluginPrivileges); err != nil {
		return nil, err
	}
	return pluginPrivileges, nil
}




func (c *Client) InspectPlugins(name string, ctx context.Context) (*PluginDetail, error) {
	resp, err := c.do("GET", "/plugins/"+name+"/json", doOptions{
		context: ctx,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, &NoSuchPlugin{ID: name}
		}
		return nil, err
	}
	resp.Body.Close()
	var pluginDetail PluginDetail
	if err := json.NewDecoder(resp.Body).Decode(&pluginDetail); err != nil {
		return nil, err
	}
	return &pluginDetail, nil
}




type RemovePluginOptions struct {
	
	Name string `qs:"-"`

	Force   bool `qs:"force"`
	Context context.Context
}




func (c *Client) RemovePlugin(opts RemovePluginOptions) (*PluginDetail, error) {
	path := "/plugins/" + opts.Name + "?" + queryString(opts)
	resp, err := c.do("DELETE", path, doOptions{context: opts.Context})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, &NoSuchPlugin{ID: opts.Name}
		}
		return nil, err
	}
	resp.Body.Close()
	var pluginDetail PluginDetail
	if err := json.NewDecoder(resp.Body).Decode(&pluginDetail); err != nil {
		return nil, err
	}
	return &pluginDetail, nil
}




type EnablePluginOptions struct {
	
	Name    string `qs:"-"`
	Timeout int64  `qs:"timeout"`

	Context context.Context
}




func (c *Client) EnablePlugin(opts EnablePluginOptions) error {
	path := "/plugins/" + opts.Name + "/enable?" + queryString(opts)
	resp, err := c.do("POST", path, doOptions{context: opts.Context})
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}




type DisablePluginOptions struct {
	
	Name string `qs:"-"`

	Context context.Context
}




func (c *Client) DisablePlugin(opts DisablePluginOptions) error {
	path := "/plugins/" + opts.Name + "/disable"
	resp, err := c.do("POST", path, doOptions{context: opts.Context})
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}




type CreatePluginOptions struct {
	
	Name string `qs:"name"`
	
	Path string `qs:"-"`

	Context context.Context
}




func (c *Client) CreatePlugin(opts CreatePluginOptions) (string, error) {
	path := "/plugins/create?" + queryString(opts)
	resp, err := c.do("POST", path, doOptions{
		data:    opts.Path,
		context: opts.Context})
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	containerNameBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(containerNameBytes), nil
}




type PushPluginOptions struct {
	
	Name string

	Context context.Context
}




func (c *Client) PushPlugin(opts PushPluginOptions) error {
	path := "/plugins/" + opts.Name + "/push"
	resp, err := c.do("POST", path, doOptions{context: opts.Context})
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}




type ConfigurePluginOptions struct {
	
	Name string `qs:"name"`
	Envs []string

	Context context.Context
}




func (c *Client) ConfigurePlugin(opts ConfigurePluginOptions) error {
	path := "/plugins/" + opts.Name + "/set"
	resp, err := c.do("POST", path, doOptions{
		data:    opts.Envs,
		context: opts.Context,
	})
	defer resp.Body.Close()
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchPlugin{ID: opts.Name}
		}
		return err
	}
	return nil
}


type NoSuchPlugin struct {
	ID  string
	Err error
}

func (err *NoSuchPlugin) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return "No such plugin: " + err.ID
}
