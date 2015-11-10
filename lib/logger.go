package tracerdemo

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Logger(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		fn(w, r)

		message := map[string]string{
			"time":           time.Now(),
			"level":          "INFO",
			"service-id":     serviceID,
			"correlation-id": r.Header.Get("correlation-id"),
			"method":         r.Method,
			"uri":            r.RequestURI,
			"duration":       time.Since(start),
		}
	})
}
