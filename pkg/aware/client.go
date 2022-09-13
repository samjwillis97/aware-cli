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

type Config struct {
	Server   string
	Token    string
	Insecure bool
	Debug    bool
}

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
	ErrNoResult    = fmt.Errorf("aware: no result")
	ErrEmptyResult = fmt.Errorf("aware: empty response from server")
)

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
