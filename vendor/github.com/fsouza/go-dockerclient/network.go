



package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)



var ErrNetworkAlreadyExists = errors.New("network already exists")




type Network struct {
	Name       string
	ID         string `json:"Id"`
	Scope      string
	Driver     string
	IPAM       IPAMOptions
	Containers map[string]Endpoint
	Options    map[string]string
	Internal   bool
	EnableIPv6 bool `json:"EnableIPv6"`
	Labels     map[string]string
}




type Endpoint struct {
	Name        string
	ID          string `json:"EndpointID"`
	MacAddress  string
	IPv4Address string
	IPv6Address string
}




func (c *Client) ListNetworks() ([]Network, error) {
	resp, err := c.do("GET", "/networks", doOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var networks []Network
	if err := json.NewDecoder(resp.Body).Decode(&networks); err != nil {
		return nil, err
	}
	return networks, nil
}



type NetworkFilterOpts map[string]map[string]bool




func (c *Client) FilteredListNetworks(opts NetworkFilterOpts) ([]Network, error) {
	params, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	path := "/networks?filters=" + string(params)
	resp, err := c.do("GET", path, doOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var networks []Network
	if err := json.NewDecoder(resp.Body).Decode(&networks); err != nil {
		return nil, err
	}
	return networks, nil
}




func (c *Client) NetworkInfo(id string) (*Network, error) {
	path := "/networks/" + id
	resp, err := c.do("GET", path, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, &NoSuchNetwork{ID: id}
		}
		return nil, err
	}
	defer resp.Body.Close()
	var network Network
	if err := json.NewDecoder(resp.Body).Decode(&network); err != nil {
		return nil, err
	}
	return &network, nil
}





type CreateNetworkOptions struct {
	Name           string                 `json:"Name" yaml:"Name" toml:"Name"`
	Driver         string                 `json:"Driver" yaml:"Driver" toml:"Driver"`
	IPAM           *IPAMOptions           `json:"IPAM,omitempty" yaml:"IPAM" toml:"IPAM"`
	Options        map[string]interface{} `json:"Options" yaml:"Options" toml:"Options"`
	Labels         map[string]string      `json:"Labels" yaml:"Labels" toml:"Labels"`
	CheckDuplicate bool                   `json:"CheckDuplicate" yaml:"CheckDuplicate" toml:"CheckDuplicate"`
	Internal       bool                   `json:"Internal" yaml:"Internal" toml:"Internal"`
	EnableIPv6     bool                   `json:"EnableIPv6" yaml:"EnableIPv6" toml:"EnableIPv6"`
	Context        context.Context        `json:"-"`
}




type IPAMOptions struct {
	Driver  string            `json:"Driver" yaml:"Driver" toml:"Driver"`
	Config  []IPAMConfig      `json:"Config" yaml:"Config" toml:"Config"`
	Options map[string]string `json:"Options" yaml:"Options" toml:"Options"`
}




type IPAMConfig struct {
	Subnet     string            `json:",omitempty"`
	IPRange    string            `json:",omitempty"`
	Gateway    string            `json:",omitempty"`
	AuxAddress map[string]string `json:"AuxiliaryAddresses,omitempty"`
}





func (c *Client) CreateNetwork(opts CreateNetworkOptions) (*Network, error) {
	resp, err := c.do(
		"POST",
		"/networks/create",
		doOptions{
			data:    opts,
			context: opts.Context,
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type createNetworkResponse struct {
		ID string
	}
	var (
		network Network
		cnr     createNetworkResponse
	)
	if err := json.NewDecoder(resp.Body).Decode(&cnr); err != nil {
		return nil, err
	}

	network.Name = opts.Name
	network.ID = cnr.ID
	network.Driver = opts.Driver

	return &network, nil
}




func (c *Client) RemoveNetwork(id string) error {
	resp, err := c.do("DELETE", "/networks/"+id, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchNetwork{ID: id}
		}
		return err
	}
	resp.Body.Close()
	return nil
}





type NetworkConnectionOptions struct {
	Container string

	
	EndpointConfig *EndpointConfig `json:"EndpointConfig,omitempty"`

	
	Force bool

	Context context.Context `json:"-"`
}




type EndpointConfig struct {
	IPAMConfig          *EndpointIPAMConfig `json:"IPAMConfig,omitempty" yaml:"IPAMConfig,omitempty" toml:"IPAMConfig,omitempty"`
	Links               []string            `json:"Links,omitempty" yaml:"Links,omitempty" toml:"Links,omitempty"`
	Aliases             []string            `json:"Aliases,omitempty" yaml:"Aliases,omitempty" toml:"Aliases,omitempty"`
	NetworkID           string              `json:"NetworkID,omitempty" yaml:"NetworkID,omitempty" toml:"NetworkID,omitempty"`
	EndpointID          string              `json:"EndpointID,omitempty" yaml:"EndpointID,omitempty" toml:"EndpointID,omitempty"`
	Gateway             string              `json:"Gateway,omitempty" yaml:"Gateway,omitempty" toml:"Gateway,omitempty"`
	IPAddress           string              `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty" toml:"IPAddress,omitempty"`
	IPPrefixLen         int                 `json:"IPPrefixLen,omitempty" yaml:"IPPrefixLen,omitempty" toml:"IPPrefixLen,omitempty"`
	IPv6Gateway         string              `json:"IPv6Gateway,omitempty" yaml:"IPv6Gateway,omitempty" toml:"IPv6Gateway,omitempty"`
	GlobalIPv6Address   string              `json:"GlobalIPv6Address,omitempty" yaml:"GlobalIPv6Address,omitempty" toml:"GlobalIPv6Address,omitempty"`
	GlobalIPv6PrefixLen int                 `json:"GlobalIPv6PrefixLen,omitempty" yaml:"GlobalIPv6PrefixLen,omitempty" toml:"GlobalIPv6PrefixLen,omitempty"`
	MacAddress          string              `json:"MacAddress,omitempty" yaml:"MacAddress,omitempty" toml:"MacAddress,omitempty"`
}





type EndpointIPAMConfig struct {
	IPv4Address string `json:",omitempty"`
	IPv6Address string `json:",omitempty"`
}





func (c *Client) ConnectNetwork(id string, opts NetworkConnectionOptions) error {
	resp, err := c.do("POST", "/networks/"+id+"/connect", doOptions{
		data:    opts,
		context: opts.Context,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchNetworkOrContainer{NetworkID: id, ContainerID: opts.Container}
		}
		return err
	}
	resp.Body.Close()
	return nil
}





func (c *Client) DisconnectNetwork(id string, opts NetworkConnectionOptions) error {
	resp, err := c.do("POST", "/networks/"+id+"/disconnect", doOptions{data: opts})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchNetworkOrContainer{NetworkID: id, ContainerID: opts.Container}
		}
		return err
	}
	resp.Body.Close()
	return nil
}




type PruneNetworksOptions struct {
	Filters map[string][]string
	Context context.Context
}




type PruneNetworksResults struct {
	NetworksDeleted []string
}




func (c *Client) PruneNetworks(opts PruneNetworksOptions) (*PruneNetworksResults, error) {
	path := "/networks/prune?" + queryString(opts)
	resp, err := c.do("POST", path, doOptions{context: opts.Context})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var results PruneNetworksResults
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return &results, nil
}


type NoSuchNetwork struct {
	ID string
}

func (err *NoSuchNetwork) Error() string {
	return fmt.Sprintf("No such network: %s", err.ID)
}



type NoSuchNetworkOrContainer struct {
	NetworkID   string
	ContainerID string
}

func (err *NoSuchNetworkOrContainer) Error() string {
	return fmt.Sprintf("No such network (%s) or container (%s)", err.NetworkID, err.ContainerID)
}
