package connections

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type httpConnector struct {
	targetHost string

	body       []byte
	bodyDigest string
	reqMethod  string
}

// HTTPConnectorOpt is an option for a new http connector
type HTTPConnectorOpt func(*httpConnector)

// HTTPConnectorWithRandomBody sets up a randomly-generate http request body
// and calculates the body's digest to set as the "digest" header value. The
// same body will be used for all requests with the resulting http connector.
func HTTPConnectorWithRandomBody(sizeBytes int) HTTPConnectorOpt {
	return func(c *httpConnector) {
		body := make([]byte, sizeBytes)
		rand.Read(body)
		hasher := sha256.New()
		hasher.Write(body)
		c.body = body
		c.bodyDigest = fmt.Sprintf("%x", hasher.Sum(nil))
	}
}

// HTTPRequestMethod sets the request method for the resulting http connector.
func HTTPRequestMethod(method string) HTTPConnectorOpt {
	return func(c *httpConnector) {
		c.reqMethod = method
	}
}

// NewHTTPConnector returns a HTTP implementation of Connector
func NewHTTPConnector(targetHost string, opts ...HTTPConnectorOpt) Connector {
	c := &httpConnector{
		targetHost: targetHost,
		reqMethod:  "GET",
		body:       nil,
		bodyDigest: "",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Connect ...
func (c *httpConnector) Connect(ctx context.Context) error {
	client := &http.Client{}
	var body io.Reader
	body = nil
	if len(c.body) > 0 {
		body = bytes.NewReader(c.body)
	}
	targetHost := c.targetHost
	if !strings.HasPrefix(targetHost, "http://") {
		targetHost = fmt.Sprintf("http://%s", c.targetHost)
	}
	targetUrl, err := url.Parse(targetHost)
	if err != nil {
		return fmt.Errorf("parsing : '%s'", c.targetHost)
	}
	req, err := http.NewRequestWithContext(ctx, c.reqMethod, targetUrl.String(), body)
	if err != nil {
		return fmt.Errorf("creating http request object: %w", err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generating new random uuid: %w", err)
	}

	if c.bodyDigest != "" {
		req.Header.Add("digest", c.bodyDigest)
	}
	req.Header.Add("x-request-id", id.String())

	var (
		region string
		ok     bool
	)
	region, ok = ctx.Value(regionKey{}).(string)
	if !ok {
		region = "unknown"
	}
	req.Header.Add("region", region)

	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("doing http %s request: %w", c.reqMethod, err)
	}

	return nil
}

// Type ...
func (c *httpConnector) Type() string { return "http" }
