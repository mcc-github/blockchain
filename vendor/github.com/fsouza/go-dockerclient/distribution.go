



package docker

import (
	"encoding/json"

	"github.com/docker/docker/api/types/registry"
)


func (c *Client) InspectDistribution(name string) (*registry.DistributionInspect, error) {
	path := "/distribution/" + name + "/json"
	resp, err := c.do("GET", path, doOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var distributionInspect registry.DistributionInspect
	if err := json.NewDecoder(resp.Body).Decode(&distributionInspect); err != nil {
		return nil, err
	}
	return &distributionInspect, nil
}
