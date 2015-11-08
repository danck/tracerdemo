package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/danck/tracerdemo/lib"
	"github.com/satori/go.uuid"
)

// handleLogger is a decorator that attaches a logger for the request lifecycle
// to a handler
func handleLogger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("consolidation-id")
		if uuid == "" {
			log.Println("Warning: no consolidation-id found")
		}
		log.Printf("Before %s", uuid)
		fn(w, r)
		log.Printf("After %s", uuid)
	}
}

// handleMarker is a decorator that attaches consolisation IDs to incoming
// requests
func handleMarker(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// generate a uuid based on timestamp and hwaddress
		//		note: if this runs in a docker container
		//		With the current implementation hwaddresses are not necessarily
		//		unique if the container's ip address isn't:
		//		https://github.com/docker/libnetwork/blob/master/netutils/utils.go#L115-L132 (08.11.2015)
		uuid := uuid.NewV1()
		r.Header.Set("consolidation-id", fmt.Sprint(uuid))
		fn(w, r)
	}
}

// handler handles the business logic for an incoming request
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "helloee? %s", r.Header.Get("consolidation-id"))
}

func main() {
	handler := tracerdemo.Identifier(handler)
	handler = tracerdemo.Logger(handler)

	http.HandleFunc("/", handler)
	log.Println(http.ListenAndServe(":8080", nil))
}
