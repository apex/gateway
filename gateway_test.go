package gateway_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/apex/gateway"
)

func Example() {
	http.HandleFunc("/", hello)
	log.Fatal(gateway.ListenAndServe("", nil))
}

func ExampleBasePath() {
	http.HandleFunc("/", hello)
	g := gateway.Gateway{
		BasePath: "v1",
		Handler:  nil,
	}
	log.Fatal(g.ListenAndServe())
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World from Go")
}
