package gateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

const apiHostSuffix = ".amazonaws.com"

// NewRequest returns a new http.Request from the given Lambda event.
func NewRequest(ctx context.Context, e events.APIGatewayProxyRequest) (*http.Request, error) {
	// path
	u, err := url.Parse(e.Path)
	if err != nil {
		return nil, errors.Wrap(err, "parsing path")
	}

	// querystring
	q := u.Query()
	for k, v := range e.QueryStringParameters {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	// base64 encoded body
	body := e.Body
	if e.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, errors.Wrap(err, "decoding base64 body")
		}
		body = string(b)
	}

	// new request
	req, err := http.NewRequest(e.HTTPMethod, u.String(), strings.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	// remote addr
	req.RemoteAddr = e.RequestContext.Identity.SourceIP

	// header fields
	for k, v := range e.Headers {
		req.Header.Set(k, v)
	}

	// content-length
	if req.Header.Get("Content-Length") == "" && body != "" {
		req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}

	// custom fields
	req.Header.Set("X-Request-Id", e.RequestContext.RequestID)
	req.Header.Set("X-Stage", e.RequestContext.Stage)

	// xray support
	if traceID := ctx.Value("x-amzn-trace-id"); traceID != nil {
		req.Header.Set("X-Amzn-Trace-Id", fmt.Sprintf("%v", traceID))
	}

	// host
	req.URL.Host = req.Header.Get("Host")
	req.Host = req.URL.Host

	// add stage to path
	if strings.HasSuffix(req.Host, apiHostSuffix) {
		req.URL.Path = "/" + e.RequestContext.Stage + req.URL.Path
	}

	return req, nil
}
