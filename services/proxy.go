package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	//"net/url"

	"github.com/danck/tracerdemo/lib"
)

var (
	next *string
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "helloee? %s", r.Header.Get("consolidation-id"))
}

func main() {
	next = flag.String(
		"next-service",
		"",
		"<ip>:<port> of the next service endpoint",
	)
	flag.Parse()
	//handler := tracerdemo.Identifier(handler)
	handler := tracerdemo.Logger(handler)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
