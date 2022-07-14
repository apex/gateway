package apiv2

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tj/assert"
)

func TestDecodeRequest_path(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets/luna",
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, "GET", r.Method)
	assert.Equal(t, `/pets/luna`, r.URL.Path)
	assert.Equal(t, `/pets/luna`, r.URL.String())
	assert.Equal(t, `/pets/luna`, r.RequestURI)
}

func TestDecodeRequest_method(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets/luna",
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "DELETE",
				Path:   "/pets/luna",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, "DELETE", r.Method)
}

func TestDecodeRequest_queryString(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath:        "/pets",
		RawQueryString: "fields=name%2Cspecies&order=desc",
		QueryStringParameters: map[string]string{
			"order":  "desc",
			"fields": "name,species",
		},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/pets",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `/pets?fields=name%2Cspecies&order=desc`, r.URL.String())
	assert.Equal(t, `desc`, r.URL.Query().Get("order"))
}

func TestDecodeRequest_multiValueQueryString(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath:        "/pets",
		RawQueryString: "fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc",
		QueryStringParameters: map[string]string{
			"multi_fields": strings.Join([]string{"name", "species"}, ","),
			"multi_arr[]":  strings.Join([]string{"arr1", "arr2"}, ","),
			"order":        "desc",
			"fields":       "name,species",
		},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/pets",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `/pets?fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc`, r.URL.String())
	assert.Equal(t, []string{"name", "species"}, r.URL.Query()["multi_fields"])
	assert.Equal(t, []string{"arr1", "arr2"}, r.URL.Query()["multi_arr[]"])
	assert.Equal(t, `/pets?fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc`, r.RequestURI)
}

func TestDecodeRequest_remoteAddr(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets",
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method:   "GET",
				Path:     "/pets",
				SourceIP: "1.2.3.4",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `1.2.3.4`, r.RemoteAddr)
}

func TestDecodeRequest_header(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets",
		Body:    `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Foo":        "bar",
			"Host":         "example.com",
		},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "1234",
			Stage:     "prod",
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Path:   "/pets",
				Method: "POST",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `example.com`, r.Host)
	assert.Equal(t, `prod`, r.Header.Get("X-Stage"))
	assert.Equal(t, `1234`, r.Header.Get("X-Request-Id"))
	assert.Equal(t, `18`, r.Header.Get("Content-Length"))
	assert.Equal(t, `application/json`, r.Header.Get("Content-Type"))
	assert.Equal(t, `bar`, r.Header.Get("X-Foo"))
}

func TestDecodeRequest_multiHeader(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets",
		Body:    `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"X-APEX":       strings.Join([]string{"apex1", "apex2"}, ","),
			"X-APEX-2":     strings.Join([]string{"apex-1", "apex-2"}, ","),
			"Content-Type": "application/json",
			"X-Foo":        "bar",
			"Host":         "example.com",
		},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "1234",
			Stage:     "prod",
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Path:   "/pets",
				Method: "POST",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `example.com`, r.Host)
	assert.Equal(t, `prod`, r.Header.Get("X-Stage"))
	assert.Equal(t, `1234`, r.Header.Get("X-Request-Id"))
	assert.Equal(t, `18`, r.Header.Get("Content-Length"))
	assert.Equal(t, `application/json`, r.Header.Get("Content-Type"))
	assert.Equal(t, `bar`, r.Header.Get("X-Foo"))
	assert.Equal(t, []string{"apex1", "apex2"}, r.Header["X-Apex"])
	assert.Equal(t, []string{"apex-1", "apex-2"}, r.Header["X-Apex-2"])
}

func TestDecodeRequest_cookie(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets",
		Body:    `{ "name": "Tobi" }`,
		Headers: map[string]string{},
		Cookies: []string{"TEST_COOKIE=TEST-VALUE"},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "1234",
			Stage:     "prod",
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Path:   "/pets",
				Method: "POST",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	c, err := r.Cookie("TEST_COOKIE")
	assert.NoError(t, err)

	assert.Equal(t, "TEST-VALUE", c.Value)
}

func TestDecodeRequest_body(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets",
		Body:    `{ "name": "Tobi" }`,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/pets",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{ "name": "Tobi" }`, string(b))
}

func TestDecodeRequest_bodyBinary(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath:         "/pets",
		Body:            `aGVsbG8gd29ybGQK`,
		IsBase64Encoded: true,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/pets",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, "hello world\n", string(b))
}

func TestDecodeRequest_context(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{}
	type key string

	var keyName key = "key"

	ctx := context.WithValue(context.Background(), keyName, "value")
	r, err := NewRequest(ctx, e)
	assert.NoError(t, err)
	v := r.Context().Value(keyName)
	assert.Equal(t, "value", v)
}
