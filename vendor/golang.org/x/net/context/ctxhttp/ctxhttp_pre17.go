





package ctxhttp 

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"
)

func nop() {}

var (
	testHookContextDoneBeforeHeaders = nop
	testHookDoReturned               = nop
	testHookDidBodyClose             = nop
)




func Do(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}

	
	cancel := make(chan struct{})
	req.Cancel = cancel

	type responseAndError struct {
		resp *http.Response
		err  error
	}
	result := make(chan responseAndError, 1)

	
	
	testHookDoReturned := testHookDoReturned
	testHookDidBodyClose := testHookDidBodyClose

	go func() {
		resp, err := client.Do(req)
		testHookDoReturned()
		result <- responseAndError{resp, err}
	}()

	var resp *http.Response

	select {
	case <-ctx.Done():
		testHookContextDoneBeforeHeaders()
		close(cancel)
		
		go func() {
			if r := <-result; r.resp != nil {
				testHookDidBodyClose()
				r.resp.Body.Close()
			}
		}()
		return nil, ctx.Err()
	case r := <-result:
		var err error
		resp, err = r.resp, r.err
		if err != nil {
			return resp, err
		}
	}

	c := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			close(cancel)
		case <-c:
			
		}
	}()
	resp.Body = &notifyingReader{resp.Body, c}

	return resp, nil
}


func Get(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return Do(ctx, client, req)
}


func Head(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return Do(ctx, client, req)
}


func Post(ctx context.Context, client *http.Client, url string, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return Do(ctx, client, req)
}


func PostForm(ctx context.Context, client *http.Client, url string, data url.Values) (*http.Response, error) {
	return Post(ctx, client, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}



type notifyingReader struct {
	io.ReadCloser
	notify chan<- struct{}
}

func (r *notifyingReader) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if err != nil && r.notify != nil {
		close(r.notify)
		r.notify = nil
	}
	return n, err
}

func (r *notifyingReader) Close() error {
	err := r.ReadCloser.Close()
	if r.notify != nil {
		close(r.notify)
		r.notify = nil
	}
	return err
}
