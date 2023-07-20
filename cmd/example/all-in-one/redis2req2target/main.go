package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	gr "github.com/redis/go-redis/v9"

	queue "github.com/takanoriyanagitani/go-http-perf/queue"
	util "github.com/takanoriyanagitani/go-http-perf/util"

	rq "github.com/takanoriyanagitani/go-http-perf/queue/external/redis/go-redis"
)

const reqKey string = "test-request"

const timeout time.Duration = 4 * time.Second
const wait time.Duration = 10 * time.Millisecond
const sleep time.Duration = 1 * time.Second

var client *gr.Client = gr.NewClient(&gr.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var pkr queue.PopKey[string] = rq.PopKeyNew(client)
var qpr queue.Pop = pkr.ToPop(reqKey)

const target string = "http://localhost:7168/"

var reqbase queue.RequestInfo = queue.
	RequestInfoPost.
	WithUrl(target)
var rms queue.RequestMetaSource = func() queue.RequestInfo { return reqbase }

const ctyp string = "application/octet-stream"

var htType queue.HeaderTransformer = func(original http.Header) {
	original.Set("Content-Type", ctyp)
}

const path string = "/path/to/something"

var utJoin queue.UrlTransformer = func(original *url.URL) (edited *url.URL) {
	return original.JoinPath(path)
}

var rrs queue.RequestSource = rms.
	ToRequestSource(qpr).
	WithEditor(htType.ToRequestEditor()).
	WithEditor(utJoin.ToRequestEditor())

var rsender queue.RequestSender = func(_ context.Context, req *http.Request) error {
	var hcli *http.Client = http.DefaultClient
	res, e := hcli.Do(req)
	return util.Select(
		func() error { return e },
		func() error {
			var body io.ReadCloser = res.Body
			defer body.Close()
			var sink io.Writer = io.Discard
			_, e := io.Copy(sink, body)
			return e
		},
		nil == e,
	)()
}

var sender queue.Sender = rrs.ToSender(rsender)

func main() {
	var ctx context.Context = context.Background()
	var tick <-chan time.Time = time.Tick(wait)
	for next := range tick {
		_ = next
		e := sender(ctx)
		if nil != e {
			log.Fatalf("Unexpected error: %v\n", e)
		}
	}
}
