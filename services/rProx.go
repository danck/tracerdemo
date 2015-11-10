package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"

	"github.com/satori/go.uuid"
)

var hosts = [...]string{"localhost:9090", "localhost:9091"}

func main() {
	prox := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			corrId := uuid.NewV1()
			// TODO(danck): custom logger
			r.URL.Scheme = "http"
			r.URL.Host = hosts[rand.Intn(len(hosts))]
			r.Header.Set("correlation-id", corrId.String())
			log.Printf("Dispatched to %s with id %s", r.URL.Host, corrId)
		},
	}

	log.Fatal(http.ListenAndServe(":8080", prox))
}
