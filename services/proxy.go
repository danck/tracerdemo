package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/danck/tracerdemo/lib"
)

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
