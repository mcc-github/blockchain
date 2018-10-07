



package docker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/docker/docker/api/types/swarm"
)


type NoSuchSecret struct {
	ID  string
	Err error
}

func (err *NoSuchSecret) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return "No such secret: " + err.ID
}




type CreateSecretOptions struct {
	Auth AuthConfiguration `qs:"-"`
	swarm.SecretSpec
	Context context.Context
}





func (c *Client) CreateSecret(opts CreateSecretOptions) (*swarm.Secret, error) {
	headers, err := headersWithAuth(opts.Auth)
	if err != nil {
		return nil, err
	}
	path := "/secrets/create?" + queryString(opts)
	resp, err := c.do("POST", path, doOptions{
		headers:   headers,
		data:      opts.SecretSpec,
		forceJSON: true,
		context:   opts.Context,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var secret swarm.Secret
	if err := json.NewDecoder(resp.Body).Decode(&secret); err != nil {
		return nil, err
	}
	return &secret, nil
}




type RemoveSecretOptions struct {
	ID      string `qs:"-"`
	Context context.Context
}




func (c *Client) RemoveSecret(opts RemoveSecretOptions) error {
	path := "/secrets/" + opts.ID
	resp, err := c.do("DELETE", path, doOptions{context: opts.Context})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchSecret{ID: opts.ID}
		}
		return err
	}
	resp.Body.Close()
	return nil
}






type UpdateSecretOptions struct {
	Auth AuthConfiguration `qs:"-"`
	swarm.SecretSpec
	Context context.Context
	Version uint64
}




func (c *Client) UpdateSecret(id string, opts UpdateSecretOptions) error {
	headers, err := headersWithAuth(opts.Auth)
	if err != nil {
		return err
	}
	params := make(url.Values)
	params.Set("version", strconv.FormatUint(opts.Version, 10))
	resp, err := c.do("POST", "/secrets/"+id+"/update?"+params.Encode(), doOptions{
		headers:   headers,
		data:      opts.SecretSpec,
		forceJSON: true,
		context:   opts.Context,
	})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return &NoSuchSecret{ID: id}
		}
		return err
	}
	defer resp.Body.Close()
	return nil
}




func (c *Client) InspectSecret(id string) (*swarm.Secret, error) {
	path := "/secrets/" + id
	resp, err := c.do("GET", path, doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, &NoSuchSecret{ID: id}
		}
		return nil, err
	}
	defer resp.Body.Close()
	var secret swarm.Secret
	if err := json.NewDecoder(resp.Body).Decode(&secret); err != nil {
		return nil, err
	}
	return &secret, nil
}




type ListSecretsOptions struct {
	Filters map[string][]string
	Context context.Context
}




func (c *Client) ListSecrets(opts ListSecretsOptions) ([]swarm.Secret, error) {
	path := "/secrets?" + queryString(opts)
	resp, err := c.do("GET", path, doOptions{context: opts.Context})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var secrets []swarm.Secret
	if err := json.NewDecoder(resp.Body).Decode(&secrets); err != nil {
		return nil, err
	}
	return secrets, nil
}
