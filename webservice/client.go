package webservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/platx/geonames/value"
)

const (
	defaultBaseURL        = "https://secure.geonames.org"
	defaultRequestTimeout = 5 * time.Second
)

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient httpDoer
	baseURL    string
	userName   string
}

type Option func(*Client)

func WithHTTPClient(httpClient httpDoer) Option {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

func WithBaseURL(baseURL string) Option {
	return func(client *Client) {
		baseURL = strings.TrimSpace(baseURL)
		if idx := strings.Index(baseURL, "?"); idx != -1 {
			baseURL = baseURL[:idx]
		}

		client.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

func NewClient(userName string, opts ...Option) *Client {
	res := &Client{
		httpClient: &http.Client{
			Transport:     http.DefaultTransport,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       defaultRequestTimeout,
		},
		baseURL:  defaultBaseURL,
		userName: userName,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func (c *Client) apiRequest(
	ctx context.Context,
	path string,
	req any,
	destination any,
) error {
	httpReq, err := c.createHTTPRequest(ctx, path, req)
	if err != nil {
		return fmt.Errorf("create http request => %w", err)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send http request => %w", err)
	}

	defer httpResp.Body.Close()

	if err := c.decodeResponse(httpResp, destination); err != nil {
		return fmt.Errorf("decode response => %w", err)
	}

	return nil
}

func (c *Client) createHTTPRequest(
	ctx context.Context,
	path string,
	req any,
) (*http.Request, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(path), nil)
	if err != nil {
		return nil, err
	}

	urlValues := url.Values{}

	NewURLEncoder(urlValues).Encode(req)

	urlValues.Set("type", "json")
	urlValues.Set("username", c.userName)
	httpReq.URL.RawQuery = urlValues.Encode()

	return httpReq, nil
}

func (c *Client) decodeResponse(
	httpRes *http.Response,
	destination any,
) error {
	if httpRes.StatusCode != http.StatusOK {
		var errResp errorResponse

		if err := c.decodeJSON(httpRes.Body, &errResp); err != nil {
			return err
		}

		return &ResponseError{code: errResp.Status.Value, message: errResp.Status.Message}
	}

	return c.decodeJSON(httpRes.Body, destination)
}

func (c *Client) url(path string) string {
	return c.baseURL + path + "?username=" + c.userName
}

func (c *Client) decodeJSON(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

type errorResponse struct {
	Status struct {
		Message string        `json:"message"`
		Value   value.ErrCode `json:"value"`
	} `json:"status"`
}
