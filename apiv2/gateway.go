package apiv2

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Decorator is a wrapper function that adds functionality to the current lambda handler.
type Decorator func(handler interface{}) interface{}

// WithDecorator adds a new decorator to the lambda handler configuration.
func WithDecorator(d Decorator) Option {
	return func(opts *Options) {
		if opts.Decorators == nil {
			opts.Decorators = make([]Decorator, 0)
		}

		opts.Decorators = append(opts.Decorators, d)
	}
}

// Options represents all options that can be applied to the lambda handler.
type Options struct {
	Decorators []Decorator
}

// Apply executes options over the lambda handler.
func (o *Options) Apply(handler interface{}) interface{} {
	handler = o.applyDecorators(handler)
	return handler
}

func (o *Options) applyDecorators(handler interface{}) interface{} {
	if o.Decorators == nil {
		return handler
	}

	for _, decorator := range o.Decorators {
		handler = decorator(handler)
	}

	return handler
}

// Option is a callback that configure som handler option.
type Option func(*Options)

// ListenAndServe is a drop-in replacement for
// http.ListenAndServe for use within AWS Lambda.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(h http.Handler, opts ...Option) error {
	if h == nil {
		h = http.DefaultServeMux
	}

	optsConfig := &Options{}

	for _, opt := range opts {
		opt(optsConfig)
	}

	gw := NewGateway(h)

	lambda.Start(wrapper(gw, optsConfig))

	return nil
}

// Invoke Handler implementation
func wrapper(gw *Gateway, opts *Options) interface{} {
	handler := opts.Apply(gw.Invoke)
	return handler
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
