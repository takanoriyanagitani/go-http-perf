package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	ot "go.opentelemetry.io/otel"
	ote "go.opentelemetry.io/otel/exporters/prometheus"
	otm "go.opentelemetry.io/otel/sdk/metric"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	exp, e := ote.New()
	if nil != e {
		panic(e)
	}
	var provider *otm.MeterProvider = otm.NewMeterProvider(otm.WithReader(exp))
	ot.SetMeterProvider(provider)
}

func main() {
	var server *http.Server = &http.Server{
		Addr:           ":7168",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var hdr http.Header = req.Header
		var hlen int = len(hdr)
		log.Printf("number of headers: %v\n", hlen)
		log.Printf("type: %s\n", hdr.Get("Content-Type"))
		var u *url.URL = req.URL
		var path string = u.Path
		log.Printf("path: %s\n", path)
		_, _ = w.Write([]byte(""))
	})

	log.Fatal(server.ListenAndServe())
}
