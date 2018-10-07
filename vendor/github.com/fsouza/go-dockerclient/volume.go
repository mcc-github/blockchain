



package docker

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var (
	
	ErrNoSuchVolume = errors.New("no such volume")

	
	ErrVolumeInUse = errors.New("volume in use and cannot be removed")
)




type Volume struct {
	Name       string            `json:"Name" yaml:"Name" toml:"Name"`
	Driver     string            `json:"Driver,omitempty" yaml:"Driver,omitempty" toml:"Driver,omitempty"`
	Mountpoint string            `json:"Mountpoint,omitempty" yaml:"Mountpoint,omitempty" toml:"Mountpoint,omitempty"`
	Labels     map[string]string `json:"Labels,omitempty" yaml:"Labels,omitempty" toml:"Labels,omitempty"`
	Options    map[string]string `json:"Options,omitempty" yaml:"Options,omitempty" toml:"Options,omitempty"`
	CreatedAt  time.Time         `json:"CreatedAt,omitempty" yaml:"CreatedAt,omitempty" toml:"CreatedAt,omitempty"`
}




type ListVolumesOptions struct {
	Filters map[string][]string
	Context context.Context
}




func (c *Client) ListVolumes(opts ListVolumesOptions) ([]Volume, error) {
	resp, err := c.do("GET", "/volumes?"+queryString(opts), doOptions{
		context: opts.Context,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	m := make(map[string]interface{})
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	var volumes []Volume
	volumesJSON, ok := m["Volumes"]
	if !ok {
		return volumes, nil
	}
	data, err := json.Marshal(volumesJSON)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &volumes); err != nil {
		return nil, err
	}
	return volumes, nil
}




type CreateVolumeOptions struct {
	Name       string
	Driver     string
	DriverOpts map[string]string
	Context    context.Context `json:"-"`
	Labels     map[string]string
}




func (c *Client) CreateVolume(opts CreateVolumeOptions) (*Volume, error) {
	resp, err := c.do("POST", "/volumes/create", doOptions{
		data:    opts,
		context: opts.Context,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var volume Volume
	if err := json.NewDecoder(resp.Body).Decode(&volume); err != nil {
		return nil, err
	}
	return &volume, nil
}




func (c *Client) InspectVolume(name string) (*Volume, error) {
	resp, err := c.do("GET", "/volumes/"+name, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, ErrNoSuchVolume
		}
		return nil, err
	}
	defer resp.Body.Close()
	var volume Volume
	if err := json.NewDecoder(resp.Body).Decode(&volume); err != nil {
		return nil, err
	}
	return &volume, nil
}




func (c *Client) RemoveVolume(name string) error {
	return c.RemoveVolumeWithOptions(RemoveVolumeOptions{Name: name})
}





type RemoveVolumeOptions struct {
	Context context.Context
	Name    string `qs:"-"`
	Force   bool
}





func (c *Client) RemoveVolumeWithOptions(opts RemoveVolumeOptions) error {
	path := "/volumes/" + opts.Name
	resp, err := c.do("DELETE", path+"?"+queryString(opts), doOptions{context: opts.Context})
	if err != nil {
		if e, ok := err.(*Error); ok {
			if e.Status == http.StatusNotFound {
				return ErrNoSuchVolume
			}
			if e.Status == http.StatusConflict {
				return ErrVolumeInUse
			}
		}
		return err
	}
	defer resp.Body.Close()
	return nil
}




type PruneVolumesOptions struct {
	Filters map[string][]string
	Context context.Context
}




type PruneVolumesResults struct {
	VolumesDeleted []string
	SpaceReclaimed int64
}




func (c *Client) PruneVolumes(opts PruneVolumesOptions) (*PruneVolumesResults, error) {
	path := "/volumes/prune?" + queryString(opts)
	resp, err := c.do("POST", path, doOptions{context: opts.Context})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var results PruneVolumesResults
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return &results, nil
}
