package apiv2

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

// NewRequest returns a new http.Request from the given Lambda event.
func NewRequest(ctx context.Context, e events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	// path
	u, err := url.Parse(e.RawPath)
	if err != nil {
		return nil, errors.Wrap(err, "parsing path")
	}

	u.RawQuery = e.RawQueryString

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
	req, err := http.NewRequest(e.RequestContext.HTTP.Method, u.String(), strings.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	// manually set RequestURI because NewRequest is for clients and req.RequestURI is for servers
	req.RequestURI = u.RequestURI()

	// remote addr
	req.RemoteAddr = e.RequestContext.HTTP.SourceIP

	// header fields
	for k, values := range e.Headers {
		for _, v := range strings.Split(values, ",") {
			req.Header.Add(k, v)
		}
	}
	for _, c := range e.Cookies {
		req.Header.Add("Cookie", c)
	}

	// content-length
	if req.Header.Get("Content-Length") == "" && body != "" {
		req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}

	// custom fields
	req.Header.Set("X-Request-Id", e.RequestContext.RequestID)
	req.Header.Set("X-Stage", e.RequestContext.Stage)

	// custom context values
	req = req.WithContext(newContext(ctx, e))

	// xray support
	if traceID := ctx.Value("x-amzn-trace-id"); traceID != nil {
		req.Header.Set("X-Amzn-Trace-Id", fmt.Sprintf("%v", traceID))
	}

	// host
	req.URL.Host = req.Header.Get("Host")
	req.Host = req.URL.Host

	return req, nil
}
