// Package gateway provides a drop-in replacement for net/http.ListenAndServe for use in AWS Lambda & API Gateway.
package gateway

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Gateway struct {
	BasePath string
	Handler  http.Handler
}

// Serve handles incoming event from AWS Lambda by wraping them into
// http.Request which is further processed by http.Handler to reply
// as a APIGatewayProxyResponse.
func (g *Gateway) Serve(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := NewRequest(ctx, e, g.BasePath)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	w := NewResponse()
	g.Handler.ServeHTTP(w, r)

	return w.End(), nil
}

// ListenAndServe registers a listener of AWS Lambda events.
func (g *Gateway) ListenAndServe() error {
	if g.Handler == nil {
		g.Handler = http.DefaultServeMux
	}

	lambda.Start(g.Serve)

	return nil
}

// ListenAndServe is a drop-in replacement for
// http.ListenAndServe for use within AWS Lambda.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(_ string, h http.Handler) error {
	g := &Gateway{Handler: h}
	return g.ListenAndServe()
}
