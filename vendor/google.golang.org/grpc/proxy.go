

package grpc

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const proxyAuthHeaderKey = "Proxy-Authorization"

var (
	
	errDisabled = errors.New("proxy is disabled for the address")
	
	httpProxyFromEnvironment = http.ProxyFromEnvironment
)

func mapAddress(ctx context.Context, address string) (*url.URL, error) {
	req := &http.Request{
		URL: &url.URL{
			Scheme: "https",
			Host:   address,
		},
	}
	url, err := httpProxyFromEnvironment(req)
	if err != nil {
		return nil, err
	}
	if url == nil {
		return nil, errDisabled
	}
	return url, nil
}






type bufConn struct {
	net.Conn
	r io.Reader
}

func (c *bufConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func doHTTPConnectHandshake(ctx context.Context, conn net.Conn, backendAddr string, proxyURL *url.URL) (_ net.Conn, err error) {
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Host: backendAddr},
		Header: map[string][]string{"User-Agent": {grpcUA}},
	}
	if t := proxyURL.User; t != nil {
		u := t.Username()
		p, _ := t.Password()
		req.Header.Add(proxyAuthHeaderKey, "Basic "+basicAuth(u, p))
	}

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
		var newAddr string
		proxyURL, err := mapAddress(ctx, addr)
		if err != nil {
			if err != errDisabled {
				return nil, err
			}
			newAddr = addr
		} else {
			newAddr = proxyURL.Host
		}

		conn, err = dialer(ctx, newAddr)
		if err != nil {
			return
		}
		if proxyURL != nil {
			
			conn, err = doHTTPConnectHandshake(ctx, conn, addr, proxyURL)
		}
		return
	}
}

func sendHTTPRequest(ctx context.Context, req *http.Request, conn net.Conn) error {
	req = req.WithContext(ctx)
	if err := req.Write(conn); err != nil {
		return fmt.Errorf("failed to write the HTTP request: %v", err)
	}
	return nil
}
