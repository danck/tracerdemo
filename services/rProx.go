package main

import (
	//"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	prox := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			// TODO(danck): add guid and custom logger
			r.URL.Scheme = "http"
			r.URL.Host = "localhost:9090"
			r.URL.Path = "/"
			r.Header.Set("consolidation-id", "rProxy test")
		},
	}

	log.Fatal(http.ListenAndServe(":8080", prox))
}
