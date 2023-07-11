package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	gr "github.com/redis/go-redis/v9"

	util "github.com/takanoriyanagitani/go-http-perf/util"

	queue "github.com/takanoriyanagitani/go-http-perf/queue"
	rq "github.com/takanoriyanagitani/go-http-perf/queue/external/redis/go-redis"
)

const popKey = "test-perf-request"

var client *gr.Client = gr.NewClient(&gr.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var pkr queue.PopKey[string] = rq.PopKeyNew(client)
var pre queue.Pop = pkr.ToPop(popKey)

var riu queue.RequestInfo = queue.
	RequestInfoPost.
	WithUrl("http://localhost:6280")
var rms queue.RequestMetaSource = func() queue.RequestInfo { return riu }

var prs queue.RequestSource = pre.ToRequestSource(rms)

func onRequest(ctx context.Context, req *http.Request, buf *bytes.Buffer) error {
	res, e := http.DefaultClient.Do(req)
	return util.Select(
		func() error { return e },
		func() error {
			var body io.ReadCloser = res.Body
			defer body.Close()
			buf.Reset()
			_, err := io.Copy(buf, body)
			return err
		},
		nil == e,
	)()
}

func main() {
	var buf bytes.Buffer
	var ctx context.Context = context.Background()
	for {
		req, e := prs(ctx)
		if nil != e {
			log.Printf("Unable to get a request: %v\n", e)
			time.Sleep(60 * time.Second)
			continue
		}
		e = onRequest(ctx, req, &buf)
		if nil != e {
			log.Printf("Unable to send a request: %v\n", e)
			time.Sleep(60 * time.Second)
			continue
		}
		time.Sleep(1 * time.Second)
	}
}
