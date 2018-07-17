

package grpc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"golang.org/x/net/context"
)

var (
	
	errDisabled = errors.New("proxy is disabled for the address")
	
	httpProxyFromEnvironment = http.ProxyFromEnvironment
)

func mapAddress(ctx context.Context, address string) (string, error) {
	req := &http.Request{
		URL: &url.URL{
			Scheme: "https",
			Host:   address,
		},
	}
	url, err := httpProxyFromEnvironment(req)
	if err != nil {
		return "", err
	}
	if url == nil {
		return "", errDisabled
	}
	return url.Host, nil
}






type bufConn struct {
	net.Conn
	r io.Reader
}

func (c *bufConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}

func doHTTPConnectHandshake(ctx context.Context, conn net.Conn, addr string) (_ net.Conn, err error) {
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	req := (&http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Host: addr},
		Header: map[string][]string{"User-Agent": {grpcUA}},
	})

	if err := sendHTTPRequest(ctx, req, conn); err != nil {
		return nil, fmt.Errorf("failed to write the HTTP request: %v", err)
	}

	r := bufio.NewReader(conn)
	resp, err := http.ReadResponse(r, req)
	if err != nil {
		return nil, fmt.Errorf("reading server HTTP response: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("failed to do connect handshake, status code: %s", resp.Status)
		}
		return nil, fmt.Errorf("failed to do connect handshake, response: %q", dump)
	}

	return &bufConn{Conn: conn, r: r}, nil
}




func newProxyDialer(dialer func(context.Context, string) (net.Conn, error)) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, addr string) (conn net.Conn, err error) {
		var skipHandshake bool
		newAddr, err := mapAddress(ctx, addr)
		if err != nil {
			if err != errDisabled {
				return nil, err
			}
			skipHandshake = true
			newAddr = addr
		}

		conn, err = dialer(ctx, newAddr)
		if err != nil {
			return
		}
		if !skipHandshake {
			conn, err = doHTTPConnectHandshake(ctx, conn, addr)
		}
		return
	}
}
