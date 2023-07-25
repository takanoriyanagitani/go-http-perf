package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	gr "github.com/redis/go-redis/v9"

	ews "github.com/takanoriyanagitani/go-http-perf/external/websocket"
	queue "github.com/takanoriyanagitani/go-http-perf/queue"

	wss "github.com/takanoriyanagitani/go-http-perf/external/websocket/std"
	rq "github.com/takanoriyanagitani/go-http-perf/queue/external/redis/go-redis"
)

const seedKey string = "test-seed-key"
const reqKey string = "test-request"
const reqMax int64 = 255

const timeout time.Duration = 4 * time.Second
const wait time.Duration = 4 * time.Millisecond
const sleep time.Duration = 1 * time.Second

var client *gr.Client = gr.NewClient(&gr.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var pstime queue.PopSeedKey[string, []byte] = rq.PopSeedKeyNew(client)
var praw queue.PopSeed[[]byte] = pstime.ToPopSeed(seedKey)
var pss ews.PopSerializedSeed = ews.PopSerializedSeed(praw)

var qlk queue.QueueLengthKey[string] = rq.QueueLengthKeyNew(client)
var rql queue.QueueLength = qlk.ToQueueLength(reqKey)
var rpl queue.PushLimit = rql.ToPushLimit(func(qlen int64) (tooMany bool) {
	return reqMax < qlen
})

var prk queue.PushRequestKey[string] = rq.PushRequestKeyNew(client)
var qpr queue.PushRequest = prk.
	ToPushRequest(reqKey).
	WithLimit(rpl)
var wpr ews.PushRequest = ews.PushRequest(qpr)

var pubdir http.Dir = http.Dir("./public")
var pubfs http.FileSystem = pubdir
var fileh http.Handler = http.FileServer(pubfs)

func onWs(conn *websocket.Conn) {
	var send ews.Send = wss.SendNew(conn)
	var ssend ews.SeedSend = pss.ToSeedSender(send)

	var buf []byte
	var recv ews.Recv = wss.RecvNew(conn, buf)
	var pr ews.PushReceived = wpr.ToPushReceived(recv)

	var app ews.App = ews.App{
		SeedSend:     ssend,
		PushReceived: pr,
	}

	var rootCtx context.Context = context.Background()
	for {
		ctx, cancel := context.WithTimeout(rootCtx, timeout)
		err := func() error {
			defer cancel()
			return app.SendRecv(ctx)
		}()
		switch {
		case errors.Is(err, io.EOF):
			return
		case nil == err:
			time.Sleep(wait)
			continue
		case errors.Is(err, queue.ErrTooManyRequests):
			time.Sleep(sleep)
			continue
		default:
			log.Fatalf("Unexpected error: %v\n", err)
		}
	}
}

func main() {
	var server *http.Server = &http.Server{
		Addr:           "0.0.0.0:7138",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.Handle("/redis2ws2redis", websocket.Handler(onWs))
	http.Handle("/pub/", fileh)
	log.Fatal(server.ListenAndServe())
}
