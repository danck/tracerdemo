package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	logDump   = "http://192.168.29.130:9200/tracerdemo/event"
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
	traceLog(corrID, fmt.Sprintf("Forwaring work to %s", *next))
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
		"time":           fmt.Sprint(time.Now().Format(time.RFC3339Nano)),
		"level":          "ERROR",
		"service-id":     *serviceID,
		"correlation-id": corrID,
		"message":        fmt.Sprint(err),
	}
	json, _ := json.Marshal(line)
	//log.Panic(string(json[:]))
	req, err := http.NewRequest("POST", logDump, bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func traceLog(corrID string, message string) {
	line := map[string]string{
		"time":           fmt.Sprint(time.Now().Format(time.RFC3339Nano)),
		"level":          "INFO",
		"service-id":     *serviceID,
		"correlation-id": corrID,
		"message":        message,
	}
	json, _ := json.Marshal(line)
	//log.Println(string(json[:]))
	req, err := http.NewRequest("POST", logDump, bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func tracer(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		fn(w, r)

		line := map[string]string{
			"time":           fmt.Sprint(time.Now().Format(time.RFC3339Nano)),
			"level":          "INFO",
			"service-id":     *serviceID,
			"correlation-id": r.Header.Get("correlation-id"),
			"method":         r.Method,
			"uri":            r.RequestURI,
			"duration":       fmt.Sprint(time.Since(start)),
		}
		json, _ := json.Marshal(line)
		//log.Println(string(json[:]))
		req, err := http.NewRequest("POST", logDump, bytes.NewBuffer(json))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	})
}
