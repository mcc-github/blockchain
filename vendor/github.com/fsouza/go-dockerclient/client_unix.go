





package docker

import (
	"context"
	"net"
	"net/http"
)



func (c *Client) initializeNativeClient(trFunc func() *http.Transport) {
	if c.endpointURL.Scheme != unixProtocol {
		return
	}
	sockPath := c.endpointURL.Path

	tr := trFunc()

	tr.Dial = func(network, addr string) (net.Conn, error) {
		return c.Dialer.Dial(unixProtocol, sockPath)
	}
	tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return c.Dialer.Dial(unixProtocol, sockPath)
	}
	c.HTTPClient.Transport = tr
}
