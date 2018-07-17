



package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/net/context"
)



type Exec struct {
	ID string `json:"Id,omitempty" yaml:"Id,omitempty"`
}




type CreateExecOptions struct {
	AttachStdin  bool            `json:"AttachStdin,omitempty" yaml:"AttachStdin,omitempty" toml:"AttachStdin,omitempty"`
	AttachStdout bool            `json:"AttachStdout,omitempty" yaml:"AttachStdout,omitempty" toml:"AttachStdout,omitempty"`
	AttachStderr bool            `json:"AttachStderr,omitempty" yaml:"AttachStderr,omitempty" toml:"AttachStderr,omitempty"`
	Tty          bool            `json:"Tty,omitempty" yaml:"Tty,omitempty" toml:"Tty,omitempty"`
	Env          []string        `json:"Env,omitempty" yaml:"Env,omitempty" toml:"Env,omitempty"`
	Cmd          []string        `json:"Cmd,omitempty" yaml:"Cmd,omitempty" toml:"Cmd,omitempty"`
	Container    string          `json:"Container,omitempty" yaml:"Container,omitempty" toml:"Container,omitempty"`
	User         string          `json:"User,omitempty" yaml:"User,omitempty" toml:"User,omitempty"`
	Context      context.Context `json:"-"`
	Privileged   bool            `json:"Privileged,omitempty" yaml:"Privileged,omitempty" toml:"Privileged,omitempty"`
}





func (c *Client) CreateExec(opts CreateExecOptions) (*Exec, error) {
	if len(opts.Env) > 0 && c.serverAPIVersion.LessThan(apiVersion125) {
		return nil, errors.New("exec configuration Env is only supported in API#1.25 and above")
	}
	path := fmt.Sprintf("/containers/%s/exec", opts.Container)
	resp, err := c.do("POST", path, doOptions{data: opts, context: opts.Context})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, &NoSuchContainer{ID: opts.Container}
		}
		return nil, err
	}
	defer resp.Body.Close()
	var exec Exec
	if err := json.NewDecoder(resp.Body).Decode(&exec); err != nil {
		return nil, err
	}

	return &exec, nil
}




type StartExecOptions struct {
	InputStream  io.Reader `qs:"-"`
	OutputStream io.Writer `qs:"-"`
	ErrorStream  io.Writer `qs:"-"`

	Detach bool `json:"Detach,omitempty" yaml:"Detach,omitempty" toml:"Detach,omitempty"`
	Tty    bool `json:"Tty,omitempty" yaml:"Tty,omitempty" toml:"Tty,omitempty"`

	
	RawTerminal bool `qs:"-"`

	
	
	
	
	
	Success chan struct{} `json:"-"`

	Context context.Context `json:"-"`
}






func (c *Client) StartExec(id string, opts StartExecOptions) error {
	cw, err := c.StartExecNonBlocking(id, opts)
	if err != nil {
		return err
	}
	if cw != nil {
		return cw.Wait()
	}
	return nil
}






func (c *Client) StartExecNonBlocking(id string, opts StartExecOptions) (CloseWaiter, error) {
	if id == "" {
		return nil, &NoSuchExec{ID: id}
	}

	path := fmt.Sprintf("/exec/%s/start", id)

	if opts.Detach {
		resp, err := c.do("POST", path, doOptions{data: opts, context: opts.Context})
		if err != nil {
			if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
				return nil, &NoSuchExec{ID: id}
			}
			return nil, err
		}
		defer resp.Body.Close()
		return nil, nil
	}

	return c.hijack("POST", path, hijackOptions{
		success:        opts.Success,
		setRawTerminal: opts.RawTerminal,
		in:             opts.InputStream,
		stdout:         opts.OutputStream,
		stderr:         opts.ErrorStream,
		data:           opts,
	})
}






func (c *Client) ResizeExecTTY(id string, height, width int) error {
	params := make(url.Values)
	params.Set("h", strconv.Itoa(height))
	params.Set("w", strconv.Itoa(width))

	path := fmt.Sprintf("/exec/%s/resize?%s", id, params.Encode())
	resp, err := c.do("POST", path, doOptions{})
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}



type ExecProcessConfig struct {
	User       string   `json:"user,omitempty" yaml:"user,omitempty" toml:"user,omitempty"`
	Privileged bool     `json:"privileged,omitempty" yaml:"privileged,omitempty" toml:"privileged,omitempty"`
	Tty        bool     `json:"tty,omitempty" yaml:"tty,omitempty" toml:"tty,omitempty"`
	EntryPoint string   `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty" toml:"entrypoint,omitempty"`
	Arguments  []string `json:"arguments,omitempty" yaml:"arguments,omitempty" toml:"arguments,omitempty"`
}






type ExecInspect struct {
	ID            string            `json:"ID,omitempty" yaml:"ID,omitempty" toml:"ID,omitempty"`
	ExitCode      int               `json:"ExitCode,omitempty" yaml:"ExitCode,omitempty" toml:"ExitCode,omitempty"`
	Running       bool              `json:"Running,omitempty" yaml:"Running,omitempty" toml:"Running,omitempty"`
	OpenStdin     bool              `json:"OpenStdin,omitempty" yaml:"OpenStdin,omitempty" toml:"OpenStdin,omitempty"`
	OpenStderr    bool              `json:"OpenStderr,omitempty" yaml:"OpenStderr,omitempty" toml:"OpenStderr,omitempty"`
	OpenStdout    bool              `json:"OpenStdout,omitempty" yaml:"OpenStdout,omitempty" toml:"OpenStdout,omitempty"`
	ProcessConfig ExecProcessConfig `json:"ProcessConfig,omitempty" yaml:"ProcessConfig,omitempty" toml:"ProcessConfig,omitempty"`
	ContainerID   string            `json:"ContainerID,omitempty" yaml:"ContainerID,omitempty" toml:"ContainerID,omitempty"`
	DetachKeys    string            `json:"DetachKeys,omitempty" yaml:"DetachKeys,omitempty" toml:"DetachKeys,omitempty"`
	CanRemove     bool              `json:"CanRemove,omitempty" yaml:"CanRemove,omitempty" toml:"CanRemove,omitempty"`
}




func (c *Client) InspectExec(id string) (*ExecInspect, error) {
	path := fmt.Sprintf("/exec/%s/json", id)
	resp, err := c.do("GET", path, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, &NoSuchExec{ID: id}
		}
		return nil, err
	}
	defer resp.Body.Close()
	var exec ExecInspect
	if err := json.NewDecoder(resp.Body).Decode(&exec); err != nil {
		return nil, err
	}
	return &exec, nil
}


type NoSuchExec struct {
	ID string
}

func (err *NoSuchExec) Error() string {
	return "No such exec instance: " + err.ID
}
