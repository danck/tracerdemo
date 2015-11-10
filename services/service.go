package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	this      *string
	next      *string
	serviceID *string
	client    = &http.Client{}
)

func handler(w http.ResponseWriter, r *http.Request) {
	corrID := r.Header.Get("correlation-id")
	if *next == "" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s received %s", *this, corrID)
		return
	}
	req, err := http.NewRequest("GET", *next, nil)
	if err != nil {
		traceError(corrID, err)
	}
	req.Header.Set("correlation-id", corrID)
	req.URL.Scheme = "http"
	resp, err := client.Do(req)
	if err != nil {
		traceError(corrID, err)
	}
	defer resp.Body.Close()
	fmt.Fprintf(w, "call returned from %s", *next)
}

func main() {
	this = flag.String(
		"this-service",
		"",
		"<ip>:<port> of this service endpoint",
	)
	next = flag.String(
		"next-service",
		"",
		"<ip>:<port> of the next service endpoint",
	)
	serviceID = flag.String(
		"id",
		"",
		"Service ID",
	)

	flag.Parse()
	if *this == "" {
		log.Fatal("no endpoint given")
	}
	if *serviceID == "" {
		log.Fatal("no serviceID given")
	}

	f, err := os.OpenFile(
		fmt.Sprintf("%s.log", *serviceID),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Fatalf("unable to open logfile: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(0)

	http.HandleFunc("/", tracer(handler))
	log.Fatal(http.ListenAndServe(*this, nil))
}

func traceError(corrID string, err error) {
	line := map[string]string{
		"time":           fmt.Sprint(time.Now()),
		"level":          "ERROR",
		"service-id":     *serviceID,
		"correlation-id": corrID,
		"message":        fmt.Sprint(err),
	}
	json, _ := json.Marshal(line)
	log.Panic(string(json[:]))
}

func traceLog(corrID string, message string) {
	line := map[string]string{
		"time":           fmt.Sprint(time.Now()),
		"level":          "INFO",
		"service-id":     *serviceID,
		"correlation-id": corrID,
		"message":        message,
	}
	json, _ := json.Marshal(line)
	log.Println(string(json[:]))
}

func tracer(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		fn(w, r)

		line := map[string]string{
			"time":           fmt.Sprint(time.Now()),
			"level":          "INFO",
			"service-id":     *serviceID,
			"correlation-id": r.Header.Get("correlation-id"),
			"method":         r.Method,
			"uri":            r.RequestURI,
			"duration":       fmt.Sprint(time.Since(start)),
		}
		json, _ := json.Marshal(line)
		log.Println(string(json[:]))
	})
}
