



package docker

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)


type NoSuchConfig struct {
	ID  string
	Err error
}

func (err *NoSuchConfig) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return "No such config: " + err.ID
}




type CreateConfigOptions struct {
	Auth AuthConfiguration `qs:"-"`
	swarm.ConfigSpec
	Context context.Context
}





func (c *Client) CreateConfig(opts CreateConfigOptions) (*swarm.Config, error) {
	headers, err := headersWithAuth(opts.Auth)
	if err != nil {
		return nil, err
	}
	path := "/configs/create?" + queryString(opts)
	resp, err := c.do("POST", path, doOptions{
		headers:   headers,
		data:      opts.ConfigSpec,
		forceJSON: true,
		context:   opts.Context,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var config swarm.Config
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}




type RemoveConfigOptions struct {
	ID      string `qs:"-"`
	Context context.Context
}




func (c *Client) RemoveConfig(opts RemoveConfigOptions) error {
	path := "/configs/" + opts.ID
	resp, err := c.do("DELETE", path, doOptions{context: opts.Context})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchConfig{ID: opts.ID}
		}
		return err
	}
	resp.Body.Close()
	return nil
}




type UpdateConfigOptions struct {
	Auth AuthConfiguration `qs:"-"`
	swarm.ConfigSpec
	Context context.Context
	Version uint64
}






func (c *Client) UpdateConfig(id string, opts UpdateConfigOptions) error {
	headers, err := headersWithAuth(opts.Auth)
	if err != nil {
		return err
	}
	params := make(url.Values)
	params.Set("version", strconv.FormatUint(opts.Version, 10))
	resp, err := c.do("POST", "/configs/"+id+"/update?"+params.Encode(), doOptions{
		headers:   headers,
		data:      opts.ConfigSpec,
		forceJSON: true,
		context:   opts.Context,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchConfig{ID: id}
		}
		return err
	}
	defer resp.Body.Close()
	return nil
}




func (c *Client) InspectConfig(id string) (*swarm.Config, error) {
	path := "/configs/" + id
	resp, err := c.do("GET", path, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, &NoSuchConfig{ID: id}
		}
		return nil, err
	}
	defer resp.Body.Close()
	var config swarm.Config
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}




type ListConfigsOptions struct {
	Filters map[string][]string
	Context context.Context
}




func (c *Client) ListConfigs(opts ListConfigsOptions) ([]swarm.Config, error) {
	path := "/configs?" + queryString(opts)
	resp, err := c.do("GET", path, doOptions{context: opts.Context})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var configs []swarm.Config
	if err := json.NewDecoder(resp.Body).Decode(&configs); err != nil {
		return nil, err
	}
	return configs, nil
}
