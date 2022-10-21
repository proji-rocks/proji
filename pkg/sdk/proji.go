package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/pkg/remote"
)

// ErrEmptyResponse is returned when the server responds with an empty body.
var ErrEmptyResponse = errors.New("empty response")

// defaultClient is the default HTTP client used by the SDK.
var defaultClient = &http.Client{
	Timeout: 10 * time.Second,
}

// Backend is the interface for the SDK. It's what directly communicates with a proji server.
type Backend struct {
	client    *http.Client
	serverURL string
}

// serverURLRegex is a regex that validates, that a given address starts with any protocol. If an address has no
// scheme, it will incorrectly parse the URL and use the host as the scheme. Using this regex and prefixing the address
// with "http://" will ensure that the URL is parsed correctly. The validation if the given protocol is valid at all,
// should happen afterwards, when the url was correctly parsed.
var serverURLRegex = regexp.MustCompile(`^(.+)://`)

// NewBackend takes a server URL and returns a new Backend. The server URL must be a valid URL.
func NewBackend(serverURL string) (*Backend, error) {
	if serverURL == "" {
		return nil, errors.New("server URL is empty")
	}

	// Regex check to validate that url starts with any protocol to ensure correct parsing.
	if !serverURLRegex.MatchString(serverURL) {
		serverURL = "https://" + serverURL
	}

	// Validate server URL.
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	// Normalize scheme.
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("url scheme must be http or https")
	}

	// Trim potentially redundant parts of path.
	u.Path = strings.TrimSuffix(u.Path, "/")
	u.Path = strings.TrimSuffix(u.Path, "/api/v1")

	return &Backend{
		client:    defaultClient,
		serverURL: u.String(),
	}, nil
}

func (c *Backend) newRequest(ctx context.Context, method, path, key string, body any) (*http.Request, error) {
	// Normalize path.
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = c.serverURL + path

	// Create the request.
	req, err := http.NewRequestWithContext(ctx, method, path, nil)
	if err != nil {
		return nil, err
	}

	if method != http.MethodGet {
		encodedBody := new(bytes.Buffer)
		err = json.NewEncoder(encodedBody).Encode(&body)
		if err != nil {
			return nil, errors.New("encode body")
		}

		req.Body = nopReadCloser{encodedBody}
	}

	// Set request headers.
	if key != "" {
		req.Header.Add("Authorization", "Bearer "+key)
	}

	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (c *Backend) do(req *http.Request, responseBody any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "send request")
	}
	defer func() { _ = resp.Body.Close() }()

	if !remote.IsStatusCodeOK(resp.StatusCode) {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.ContentLength == 0 {
		return ErrEmptyResponse
	}

	if req.Method != http.MethodGet {
		return nil
	}

	if err = json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return errors.Wrap(err, "read response body")
	}

	return nil
}

func (c *Backend) call(ctx context.Context, method, path, key string, requestBody, responseBody any) error {
	req, err := c.newRequest(ctx, method, path, key, requestBody)
	if err != nil {
		return errors.Wrap(err, "create request")
	}

	return c.do(req, responseBody)
}

// Call executes a request against the server. It takes a context, a method, a path, a key, a request body and a
// response body. It will return an error if the request fails.
func (c *Backend) Call(ctx context.Context, method, path, key string, requestBody, responseBody any) error {
	return c.call(ctx, method, path, key, requestBody, responseBody)
}

// Compile time check to make sure the type implements the interface.
var _ io.ReadCloser = nopReadCloser{}

// nopReadCloser is an implementation of `io.ReadCloser` that wraps an `io.Reader`. This does not alter the underlying
// `io.Reader`'s behavior. It just adds a `Close` method that does nothing. This is needed to make `http.Request`'s
// `Body` method work.
type nopReadCloser struct {
	io.Reader
}

// Close does nothing. It's here to satisfy the `io.ReadCloser` interface.
func (nopReadCloser) Close() error { return nil }
