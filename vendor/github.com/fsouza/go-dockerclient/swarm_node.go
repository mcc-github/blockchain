



package docker

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)


type NoSuchNode struct {
	ID  string
	Err error
}

func (err *NoSuchNode) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return "No such node: " + err.ID
}




type ListNodesOptions struct {
	Filters map[string][]string
	Context context.Context
}




func (c *Client) ListNodes(opts ListNodesOptions) ([]swarm.Node, error) {
	path := "/nodes?" + queryString(opts)
	resp, err := c.do("GET", path, doOptions{context: opts.Context})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var nodes []swarm.Node
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}




func (c *Client) InspectNode(id string) (*swarm.Node, error) {
	resp, err := c.do("GET", "/nodes/"+id, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, &NoSuchNode{ID: id}
		}
		return nil, err
	}
	defer resp.Body.Close()
	var node swarm.Node
	if err := json.NewDecoder(resp.Body).Decode(&node); err != nil {
		return nil, err
	}
	return &node, nil
}




type UpdateNodeOptions struct {
	swarm.NodeSpec
	Version uint64
	Context context.Context
}




func (c *Client) UpdateNode(id string, opts UpdateNodeOptions) error {
	params := make(url.Values)
	params.Set("version", strconv.FormatUint(opts.Version, 10))
	path := "/nodes/" + id + "/update?" + params.Encode()
	resp, err := c.do("POST", path, doOptions{
		context:   opts.Context,
		forceJSON: true,
		data:      opts.NodeSpec,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchNode{ID: id}
		}
		return err
	}
	resp.Body.Close()
	return nil
}




type RemoveNodeOptions struct {
	ID      string
	Force   bool
	Context context.Context
}




func (c *Client) RemoveNode(opts RemoveNodeOptions) error {
	params := make(url.Values)
	params.Set("force", strconv.FormatBool(opts.Force))
	path := "/nodes/" + opts.ID + "?" + params.Encode()
	resp, err := c.do("DELETE", path, doOptions{context: opts.Context})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchNode{ID: opts.ID}
		}
		return err
	}
	resp.Body.Close()
	return nil
}