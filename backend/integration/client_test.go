//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Client struct {
	base    string
	http    *http.Client
	tokenFn func() string
}

func NewClient(base string, tokenFn func() string) *Client {
	return &Client{base: base, http: &http.Client{}, tokenFn: tokenFn}
}

func (c *Client) Post(ctx context.Context, path string, header http.Header, body any) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, path, header, body)
}

func (c *Client) Get(ctx context.Context, path string, header http.Header) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, path, header, nil)
}

func (c *Client) Put(ctx context.Context, path string, header http.Header, body any) (*http.Response, error) {
	return c.do(ctx, http.MethodPut, path, header, body)
}

func (c *Client) Delete(ctx context.Context, path string, header http.Header) (*http.Response, error) {
	return c.do(ctx, http.MethodDelete, path, header, nil)
}

func (c *Client) do(ctx context.Context, method, path string, header http.Header, body any) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.base+path, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.tokenFn != nil {
		if tok := c.tokenFn(); tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
	}

	for k, v := range header {
		req.Header[k] = v
	}

	return c.http.Do(req)
}
