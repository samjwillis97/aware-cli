package aware

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// Config is the aware connection config.
type Config struct {
	Server   string
	Token    string
	Insecure bool
	Debug    bool
}

// Client is an aware client.
type Client struct {
	transport http.RoundTripper
	insecure  bool
	server    string
	token     string
	timeout   time.Duration
	debug     bool
}

// Header is a key, value pair for request headers.
type Header map[string]string

var (
	// ErrNoResult denotes no result from the API.
	ErrNoResult = fmt.Errorf("aware: no result")
	// ErrEmptyResult denotes an empty response from the API.
	ErrEmptyResult = fmt.Errorf("aware: empty response from server")
)

// ErrUnexpectedResponse denotes a response code that was not expected.
type ErrUnexpectedResponse struct {
	Status     string
	StatusCode int
}

func (e *ErrUnexpectedResponse) Error() string {
	return e.Status
}

// NewClient creates a new aware client.
func NewClient(c Config) *Client {
	client := Client{
		server: strings.TrimSuffix(c.Server, "/"),
		token:  c.Token,
		debug:  c.Debug,
	}

	client.transport = &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: client.insecure},
		DialContext: (&net.Dialer{
			Timeout: client.timeout,
		}).DialContext,
	}

	return &client
}

func (c *Client) request(ctx context.Context, method, endpoint string, body []byte, headers Header) (*http.Response, error) {
	var (
		req *http.Request
		res *http.Response
		err error
	)

	req, err = http.NewRequest(method, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	defer func() {
		if c.debug {
			dump(req, res)
		}
	}()

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}

	res, err = c.transport.RoundTrip(req.WithContext(ctx))

	return res, err
}

func dump(req *http.Request, res *http.Response) {
	reqDump, _ := httputil.DumpRequest(req, true)
	respDump, _ := httputil.DumpResponse(res, true)

	prettyPrintDump("Request Details", reqDump)
	prettyPrintDump("Response Details", respDump)
}

func prettyPrintDump(heading string, data []byte) {
	const separatorWidth = 60

	fmt.Printf("\n\n%s", strings.ToUpper(heading))
	fmt.Printf("\n%s\n\n", strings.Repeat("-", separatorWidth))
	fmt.Print(string(data))
}

func formatUnexpectedResponse(res *http.Response) *ErrUnexpectedResponse {
	return &ErrUnexpectedResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
	}
}
