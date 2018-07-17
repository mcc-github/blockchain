



package docker

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)


type NoSuchService struct {
	ID  string
	Err error
}

func (err *NoSuchService) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return "No such service: " + err.ID
}




type CreateServiceOptions struct {
	Auth AuthConfiguration `qs:"-"`
	swarm.ServiceSpec
	Context context.Context
}





func (c *Client) CreateService(opts CreateServiceOptions) (*swarm.Service, error) {
	headers, err := headersWithAuth(opts.Auth)
	if err != nil {
		return nil, err
	}
	path := "/services/create?" + queryString(opts)
	resp, err := c.do("POST", path, doOptions{
		headers:   headers,
		data:      opts.ServiceSpec,
		forceJSON: true,
		context:   opts.Context,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var service swarm.Service
	if err := json.NewDecoder(resp.Body).Decode(&service); err != nil {
		return nil, err
	}
	return &service, nil
}




type RemoveServiceOptions struct {
	ID      string `qs:"-"`
	Context context.Context
}




func (c *Client) RemoveService(opts RemoveServiceOptions) error {
	path := "/services/" + opts.ID
	resp, err := c.do("DELETE", path, doOptions{context: opts.Context})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchService{ID: opts.ID}
		}
		return err
	}
	resp.Body.Close()
	return nil
}




type UpdateServiceOptions struct {
	Auth AuthConfiguration `qs:"-"`
	swarm.ServiceSpec
	Context context.Context
	Version uint64
}




func (c *Client) UpdateService(id string, opts UpdateServiceOptions) error {
	headers, err := headersWithAuth(opts.Auth)
	if err != nil {
		return err
	}
	params := make(url.Values)
	params.Set("version", strconv.FormatUint(opts.Version, 10))
	resp, err := c.do("POST", "/services/"+id+"/update?"+params.Encode(), doOptions{
		headers:   headers,
		data:      opts.ServiceSpec,
		forceJSON: true,
		context:   opts.Context,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchService{ID: id}
		}
		return err
	}
	defer resp.Body.Close()
	return nil
}




func (c *Client) InspectService(id string) (*swarm.Service, error) {
	path := "/services/" + id
	resp, err := c.do("GET", path, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, &NoSuchService{ID: id}
		}
		return nil, err
	}
	defer resp.Body.Close()
	var service swarm.Service
	if err := json.NewDecoder(resp.Body).Decode(&service); err != nil {
		return nil, err
	}
	return &service, nil
}




type ListServicesOptions struct {
	Filters map[string][]string
	Context context.Context
}




func (c *Client) ListServices(opts ListServicesOptions) ([]swarm.Service, error) {
	path := "/services?" + queryString(opts)
	resp, err := c.do("GET", path, doOptions{context: opts.Context})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var services []swarm.Service
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		return nil, err
	}
	return services, nil
}



type LogsServiceOptions struct {
	Context           context.Context
	Service           string        `qs:"-"`
	OutputStream      io.Writer     `qs:"-"`
	ErrorStream       io.Writer     `qs:"-"`
	InactivityTimeout time.Duration `qs:"-"`
	Tail              string

	
	RawTerminal bool `qs:"-"`
	Since       int64
	Follow      bool
	Stdout      bool
	Stderr      bool
	Timestamps  bool
	Details     bool
}









func (c *Client) GetServiceLogs(opts LogsServiceOptions) error {
	if opts.Service == "" {
		return &NoSuchService{ID: opts.Service}
	}
	if opts.Tail == "" {
		opts.Tail = "all"
	}
	path := "/services/" + opts.Service + "/logs?" + queryString(opts)
	return c.stream("GET", path, streamOptions{
		setRawTerminal:    opts.RawTerminal,
		stdout:            opts.OutputStream,
		stderr:            opts.ErrorStream,
		inactivityTimeout: opts.InactivityTimeout,
		context:           opts.Context,
	})
}
