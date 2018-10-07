



package docker

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types/swarm"
)


type NoSuchTask struct {
	ID  string
	Err error
}

func (err *NoSuchTask) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return "No such task: " + err.ID
}




type ListTasksOptions struct {
	Filters map[string][]string
	Context context.Context
}




func (c *Client) ListTasks(opts ListTasksOptions) ([]swarm.Task, error) {
	path := "/tasks?" + queryString(opts)
	resp, err := c.do("GET", path, doOptions{context: opts.Context})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var tasks []swarm.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}




func (c *Client) InspectTask(id string) (*swarm.Task, error) {
	resp, err := c.do("GET", "/tasks/"+id, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, &NoSuchTask{ID: id}
		}
		return nil, err
	}
	defer resp.Body.Close()
	var task swarm.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}
