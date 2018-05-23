<img src="http://tjholowaychuk.com:6000/svg/title/APEX/GATEWAY">

Package gateway provides a drop-in replacement for net/http's `ListenAndServe` for use in AWS Lambda & API Gateway, simply swap it out for `gateway.ListenAndServe`. Extracted from [Up](https://github.com/apex/up) which provides additional middleware features and operational functionality.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/aws/aws-lambda-go"
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
