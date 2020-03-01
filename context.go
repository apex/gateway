package gateway

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// Key is the type used for any items added to the request context.
type Key int

// requestContextKey is the key for the api gateway proxy `RequestContext`.
const requestContextKey Key = iota

// GetRequestContextKey is useful for creating custom claims when testing locally
// ctx := context.WithValue(r.Context(), algnhsa.GetProxyRequestContextKey(), customLocalClaims)
// r = r.Clone(ctx)y
func GetRequestContextKey() Key {
	return requestContextKey
}

// newContext returns a new Context with specific api gateway proxy values.
func newContext(ctx context.Context, e events.APIGatewayProxyRequest) context.Context {
	return context.WithValue(ctx, requestContextKey, e.RequestContext)
}

// RequestContext returns the APIGatewayProxyRequestContext value stored in ctx.
func RequestContext(ctx context.Context) (events.APIGatewayProxyRequestContext, bool) {
	c, ok := ctx.Value(requestContextKey).(events.APIGatewayProxyRequestContext)
	return c, ok
}
