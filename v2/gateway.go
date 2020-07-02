// Package gateway provides a drop-in replacement for net/http.ListenAndServe for use in AWS Lambda & API Gateway.
package gateway

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ListenAndServe is a drop-in replacement for
// http.ListenAndServe for use within AWS Lambda.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, h http.Handler) error {
	if h == nil {
		h = http.DefaultServeMux
	}

	gw := NewGateway(h)

	lambda.StartHandler(gw)

	return nil
}

// NewGateway creates a gateway using the provided http.Handler enabling use in existing aws-lambda-go
// projects
func NewGateway(h http.Handler) *Gateway {
	return &Gateway{h: h}
}

// Gateway wrap a http handler to enable use as a lambda.Handler
type Gateway struct {
	h http.Handler
}

// Invoke Handler implementation
func (gw *Gateway) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	var evt events.APIGatewayV2HTTPRequest

	if err := json.Unmarshal(payload, &evt); err != nil {
		return []byte{}, err
	}

	r, err := NewRequest(ctx, evt)
	if err != nil {
		return []byte{}, err
	}

	w := NewResponse()
	gw.h.ServeHTTP(w, r)

	resp := w.End()

	return json.Marshal(&resp)
}
