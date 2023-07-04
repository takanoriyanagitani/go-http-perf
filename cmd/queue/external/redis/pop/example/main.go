package main

import (
	"context"
	"log"
	"net/http"

	gr "github.com/redis/go-redis/v9"

	queue "github.com/takanoriyanagitani/go-http-perf/queue"
	rq "github.com/takanoriyanagitani/go-http-perf/queue/external/redis/go-redis"
)

const popKey string = "test-perf-request"

var client *gr.Client = gr.NewClient(&gr.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var pkr queue.PopKey[string] = rq.PopKeyNew(client)
var pk queue.Pop = pkr.ToPop(popKey)

var dummyRequestSender queue.RequestSender = func(_ctx context.Context, _q *http.Request) error {
	log.Printf("trying to send...\n")
	return nil
}

var dummyRequestParser queue.SerializedRequestParser = func(_ []byte) (*http.Request, error) {
	log.Printf("trying to parse...\n")
	return nil, nil
}

var dummyBytesSender queue.SerializedRequestSender = dummyRequestSender.ToSerializedRequestSender(
	dummyRequestParser,
)

var dummySender queue.Sender = dummyBytesSender.ToSender(pk)

func main() {
	e := dummySender(context.Background())
	if nil != e {
		log.Printf("Unable to send: %v\n", e)
	}
}
