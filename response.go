package gateway

import (
	"bytes"
	"encoding/base64"
	"mime"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// ResponseWriter implements the http.ResponseWriter interface
// in order to support the API Gateway Lambda HTTP "protocol".
type ResponseWriter struct {
	out         events.APIGatewayProxyResponse
	buf         bytes.Buffer
	header      http.Header
	wroteHeader bool
}

// NewResponse returns a new response writer to capture http output.
func NewResponse() *ResponseWriter {
	return &ResponseWriter{}
}

// Header implementation.
func (w *ResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}

	return w.header
}

// Write implementation.
func (w *ResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	// TODO: HEAD? ignore

	return w.buf.Write(b)
}

// WriteHeader implementation.
func (w *ResponseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf8")
	}

	w.out.StatusCode = status

	h := make(map[string]string)

	for k, v := range w.Header() {
		if len(v) > 0 {
			h[k] = v[len(v)-1]
		}
	}

	w.out.Headers = h
	w.wroteHeader = true
}

// End the request.
func (w *ResponseWriter) End() events.APIGatewayProxyResponse {
	w.out.IsBase64Encoded = isBinary(w.header)

	if w.out.IsBase64Encoded {
		w.out.Body = base64.StdEncoding.EncodeToString(w.buf.Bytes())
	} else {
		w.out.Body = w.buf.String()
	}

	return w.out
}

// isBinary returns true if the response reprensents binary.
func isBinary(h http.Header) bool {
	contentType := strings.Split(h.Get("Content-Type"), ";")[0]
	if !isTextMime(contentType) {
		return true
	}

	if h.Get("Content-Encoding") == "gzip" {
		return true
	}

	return false
}

// isTextMime returns true if the content type represents textual data.
func isTextMime(kind string) bool {
	mt, _, err := mime.ParseMediaType(kind)
	if err != nil {
		return false
	}

	if strings.HasPrefix(mt, "text/") {
		return true
	}

	switch mt {
	case "svg+xml":
		return true
	case "application/json":
		return true
	case "application/xml":
		return true
	default:
		return false
	}
}
