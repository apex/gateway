package gateway

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// key is the type used for any items added to the request context.
type key int

// requestContextKey is the key for the api gateway proxy `Request`.
const requestContextKey key = iota

// RequestContext returns the APIGatewayV2HTTPRequestContext value stored in ctx.
func RequestContext(ctx context.Context) (events.APIGatewayV2HTTPRequestContext, bool) {
	c, ok := GetRequest(ctx)
	if ok {
		return c.RequestContext, true
	}
	return events.APIGatewayV2HTTPRequestContext{}, false
}

// GetRequest returns the APIGatewayV2HTTPRequest value stored in ctx.
func GetRequest(ctx context.Context) (events.APIGatewayV2HTTPRequest, bool) {
	c, ok := ctx.Value(requestContextKey).(events.APIGatewayV2HTTPRequest)
	return c, ok
}

// newContext returns a new Context with specific api gateway v2 values.
func newContext(ctx context.Context, e events.APIGatewayV2HTTPRequest) context.Context {
	return context.WithValue(ctx, requestContextKey, e)
}
