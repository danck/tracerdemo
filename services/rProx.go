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
			r.URL.Scheme = "http"
			r.URL.Host = "localhost:9090"
			r.URL.Path = "/"
		},
	}

	log.Fatal(http.ListenAndServe(":8080", prox))
}
