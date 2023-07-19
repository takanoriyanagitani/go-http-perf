package main

import (
	"context"
	"log"
	"net/http"
	"time"

	ot "go.opentelemetry.io/otel"
	ote "go.opentelemetry.io/otel/exporters/prometheus"
	otm "go.opentelemetry.io/otel/metric"
	osm "go.opentelemetry.io/otel/sdk/metric"

	oth "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var promExporter *ote.Exporter = must(ote.New())

var meterProvider *osm.MeterProvider = osm.NewMeterProvider(osm.WithReader(promExporter))

func init() {
	ot.SetMeterProvider(meterProvider)
}

func must[T any](t T, e error) T {
	if nil != e {
		panic(e)
	}
	return t
}

func main() {
	var server *http.Server = &http.Server{
		Addr:           ":7168",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.Handle("/metrics", promhttp.Handler())

	var meter otm.Meter = meterProvider.Meter("sample-meter")
	var ctr otm.Int64Counter = must(meter.Int64Counter(
		"dummy-requests",
		otm.WithDescription("Number of requests"),
	))

	var hlenCnt otm.Int64Counter = must(meter.Int64Counter(
		"dummy-header-length",
		otm.WithDescription("Number of header length"),
	))

	var original http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
		var ctx context.Context = req.Context()
		ctr.Add(ctx, 1)

		var hdr http.Header = req.Header
		var hlen int = len(hdr)
		hlenCnt.Add(ctx, int64(hlen))
		_, _ = w.Write([]byte(""))
	}
	var wh http.Handler = oth.NewHandler(original, "dummy-handler")
	http.Handle("/", wh)

	log.Fatal(server.ListenAndServe())
}
