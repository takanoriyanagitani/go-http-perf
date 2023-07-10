package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	gr "github.com/redis/go-redis/v9"

	ew "github.com/takanoriyanagitani/go-http-perf/external/websocket"
	sw "github.com/takanoriyanagitani/go-http-perf/external/websocket/std"
	queue "github.com/takanoriyanagitani/go-http-perf/queue"
	rq "github.com/takanoriyanagitani/go-http-perf/queue/external/redis/go-redis"
)

const pushKey string = "test-perf-seed"
const popKey string = "test-perf-request"

var client *gr.Client = gr.NewClient(&gr.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var pkr queue.PopKey[string] = rq.PopKeyNew(client)
var pk queue.Pop = pkr.ToPop(popKey)

var psk queue.PopSeedKey[string, []byte] = rq.PopSeedKeyNew(client)
var ps queue.PopSeed[[]byte] = psk.ToPopSeed(pushKey)
var pss ew.PopSerializedSeed = ew.PopSerializedSeed(ps)

var prk queue.PushRequestKey[string] = rq.PushRequestKeyNew(client)
var qlk queue.QueueLengthKey[string] = rq.QueueLengthKeyNew(client)
var ql queue.QueueLength = qlk.ToQueueLength(popKey)
var pl queue.PushLimit = ql.ToPushLimit(func(length int64) (tooMany bool) {
	return 256 < length
})
var qpr queue.PushRequest = prk.
	ToPushRequest(popKey).
	WithLimit(pl)
var pr ew.PushRequest = ew.PushRequest(qpr)

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

func onConnect(con *websocket.Conn) {
	var buf []byte
	var recv ew.Recv = sw.RecvNew(con, buf)
	var send ew.Send = sw.SendNew(con)
	var ss ew.SeedSend = pss.ToSeedSender(send)
	var req2redis ew.PushReceived = pr.ToPushReceived(recv)
	var app ew.App = ew.App{
		SeedSend:     ss,
		PushReceived: req2redis,
	}
	func() {
		var ctx context.Context = context.Background()
		for {
			e := app.SendRecv(ctx)
			switch e {
			case nil:
				time.Sleep(1 * time.Second)
			case io.EOF:
				return
			case queue.ErrTooManyRequests:
				log.Printf("%v\n", e)
				time.Sleep(60 * time.Second)
			default:
				log.Printf("Err: %v\n", e)
				time.Sleep(10 * time.Second)
			}
		}
	}()
}

func main() {
	var server *http.Server = &http.Server{
		Addr:           ":7058",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.Handle("/redis2ws2redis", websocket.Handler(onConnect))
	log.Fatal(server.ListenAndServe())
}
