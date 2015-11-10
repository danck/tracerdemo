package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/satori/go.uuid"
)

var (
	serviceID *string
	hosts     = [...]string{"localhost:9090", "localhost:9091"}
	client    = &http.Client{}
	logDump   = "http://192.168.29.130:9200/tracerdemo/event"
)

func main() {
	serviceID = flag.String(
		"id",
		"",
		"Service ID",
	)

	flag.Parse()

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

	prox := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			corrId := uuid.NewV1()
			traceLog(corrId.String(), "correlation id created")
			// TODO(danck): custom logger
			r.URL.Scheme = "http"
			host := hosts[rand.Intn(len(hosts))]
			r.URL.Host = host
			r.Header.Set("correlation-id", corrId.String())
			traceLog(corrId.String(), fmt.Sprintf("dispatched to %s", host))
		},
	}

	log.Fatal(http.ListenAndServe(":8080", prox))
}

func traceLog(corrID string, message string) {
	line := map[string]string{
		"level":          "INFO",
		"time":           fmt.Sprint(time.Now().Format(time.RFC3339Nano)),
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
