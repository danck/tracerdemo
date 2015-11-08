package tracerdemo

import (
	"log"
	"net/http"
	"time"
)

func Logger(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		fn(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			r.Header.Get("consolidation-id"),
			time.Since(start),
		)
	})
}
