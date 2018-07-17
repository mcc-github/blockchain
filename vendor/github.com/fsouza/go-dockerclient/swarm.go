



package docker

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

var (
	
	
	ErrNodeAlreadyInSwarm = errors.New("node already in a Swarm")

	
	
	ErrNodeNotInSwarm = errors.New("node is not in a Swarm")
)



type InitSwarmOptions struct {
	swarm.InitRequest
	Context context.Context
}



func (c *Client) InitSwarm(opts InitSwarmOptions) (string, error) {
	path := "/swarm/init"
	resp, err := c.do("POST", path, doOptions{
		data:      opts.InitRequest,
		forceJSON: true,
		context:   opts.Context,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && (e.Status == http.StatusNotAcceptable || e.Status == http.StatusServiceUnavailable) {
			return "", ErrNodeAlreadyInSwarm
		}
		return "", err
	}
	defer resp.Body.Close()
	var response string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	return response, nil
}



type JoinSwarmOptions struct {
	swarm.JoinRequest
	Context context.Context
}



func (c *Client) JoinSwarm(opts JoinSwarmOptions) error {
	path := "/swarm/join"
	resp, err := c.do("POST", path, doOptions{
		data:      opts.JoinRequest,
		forceJSON: true,
		context:   opts.Context,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && (e.Status == http.StatusNotAcceptable || e.Status == http.StatusServiceUnavailable) {
			return ErrNodeAlreadyInSwarm
		}
	}
	resp.Body.Close()
	return err
}



type LeaveSwarmOptions struct {
	Force   bool
	Context context.Context
}



func (c *Client) LeaveSwarm(opts LeaveSwarmOptions) error {
	params := make(url.Values)
	params.Set("force", strconv.FormatBool(opts.Force))
	path := "/swarm/leave?" + params.Encode()
	resp, err := c.do("POST", path, doOptions{
		context: opts.Context,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && (e.Status == http.StatusNotAcceptable || e.Status == http.StatusServiceUnavailable) {
			return ErrNodeNotInSwarm
		}
	}
	resp.Body.Close()
	return err
}



type UpdateSwarmOptions struct {
	Version            int
	RotateWorkerToken  bool
	RotateManagerToken bool
	Swarm              swarm.Spec
	Context            context.Context
}



func (c *Client) UpdateSwarm(opts UpdateSwarmOptions) error {
	params := make(url.Values)
	params.Set("version", strconv.Itoa(opts.Version))
	params.Set("rotateWorkerToken", strconv.FormatBool(opts.RotateWorkerToken))
	params.Set("rotateManagerToken", strconv.FormatBool(opts.RotateManagerToken))
	path := "/swarm/update?" + params.Encode()
	resp, err := c.do("POST", path, doOptions{
		data:      opts.Swarm,
		forceJSON: true,
		context:   opts.Context,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && (e.Status == http.StatusNotAcceptable || e.Status == http.StatusServiceUnavailable) {
			return ErrNodeNotInSwarm
		}
	}
	resp.Body.Close()
	return err
}



func (c *Client) InspectSwarm(ctx context.Context) (swarm.Swarm, error) {
	response := swarm.Swarm{}
	resp, err := c.do("GET", "/swarm", doOptions{
		context: ctx,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && (e.Status == http.StatusNotAcceptable || e.Status == http.StatusServiceUnavailable) {
			return response, ErrNodeNotInSwarm
		}
		return response, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
