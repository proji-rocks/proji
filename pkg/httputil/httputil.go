package httputil

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/pkg/logging"
)

var httpClient = &http.Client{
	// Timeout: 3 * time.Minute, // Appropriate timeout for the request.
}

func rawRequestWithClient(ctx context.Context, client *http.Client, method, url string, body []byte) (*http.Response, error) {
	logger := logging.FromContext(ctx)

	if url == "" {
		return nil, errors.New("url cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	logger.Debugf("%s %s", strings.ToUpper(method), url)

	return client.Do(req)
}

// GetWithClient returns the response for the given URL. It accepts an optional http.Client to use for the request.
func GetWithClient(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	return rawRequestWithClient(ctx, client, http.MethodGet, url, nil)
}

// Get returns the response for the given URL. It calls GetWithClient internally and uses the default HTTP client.
func Get(ctx context.Context, url string) (*http.Response, error) {
	return GetWithClient(ctx, httpClient, url)
}
