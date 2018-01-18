package gateway

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// key is the type used for any items added to the request context.
type key int

// requestContextKey is the key for the api gateway proxy `RequestContext`.
const requestContextKey key = iota

// newContext returns a new Context with specific api gateway proxy values.
func newContext(ctx context.Context, e events.APIGatewayProxyRequest) context.Context {
	return context.WithValue(ctx, requestContextKey, e.RequestContext)
}

// RequestContext returns the APIGatewayProxyRequestContext value stored in ctx.
func RequestContext(ctx context.Context) (events.APIGatewayProxyRequestContext, bool) {
	c, ok := ctx.Value(requestContextKey).(events.APIGatewayProxyRequestContext)
	return c, ok
}
