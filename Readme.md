<img src="http://tjholowaychuk.com:6000/svg/title/APEX/GATEWAY">

Package gateway provides a drop-in replacement for net/http's `ListenAndServe` for use in [AWS Lambda](https://aws.amazon.com/lambda/) & [API Gateway](https://aws.amazon.com/api-gateway/), simply swap it out for `gateway.ListenAndServe`. Extracted from [Up](https://github.com/apex/up) which provides additional middleware features and operational functionality.

There are two versions of this library, version 1.x supports AWS API Gateway 1.0 events used by the original [REST APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-rest-api.html), and 2.x which supports 2.0 events used by the [HTTP APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api.html). For more information on the options read [Choosing between HTTP APIs and REST APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-vs-rest.html) on the AWS documentation website.

# Installation

To install version 1.x for REST APIs. 

```
go get github.com/apex/gateway
```

To install version 2.x for HTTP APIs. 

```
go get github.com/apex/gateway/v2
```

# Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
)

func main() {
	http.HandleFunc("/", hello)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	// example retrieving values from the api gateway proxy request context.
	requestContext, ok := gateway.RequestContext(r.Context())
	if !ok || requestContext.Authorizer["sub"] == nil {
		fmt.Fprint(w, "Hello World from Go")
		return
	}

	userID := requestContext.Authorizer["sub"].(string)
	fmt.Fprintf(w, "Hello %s from Go", userID)
}
```

---

[![GoDoc](https://godoc.org/github.com/apex/up-go?status.svg)](https://godoc.org/github.com/apex/gateway)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)

<a href="https://apex.sh"><img src="http://tjholowaychuk.com:6000/svg/sponsor"></a>
