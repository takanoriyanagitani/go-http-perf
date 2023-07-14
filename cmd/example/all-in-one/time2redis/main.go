package main

import (
	"context"
	"log"
	"time"

	gr "github.com/redis/go-redis/v9"

	queue "github.com/takanoriyanagitani/go-http-perf/queue"
	rq "github.com/takanoriyanagitani/go-http-perf/queue/external/redis/go-redis"
)

const seedKey string = "test-seed-key"
const seedMax int64 = 255
const timeout time.Duration = 4 * time.Second
const wait time.Duration = 4 * time.Millisecond
const sleep time.Duration = 1 * time.Second

var client *gr.Client = gr.NewClient(&gr.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var qlenk queue.QueueLengthKey[string] = rq.QueueLengthKeyNew(client)
var qlen queue.QueueLength = qlenk.ToQueueLength(seedKey)
var plmt queue.PushLimit = qlen.ToPushLimit(func(l int64) (tooMany bool) {
	return seedMax < l
})

var pktime queue.PushKey[string, []byte] = rq.PushKeyNew[[]byte](client)
var praw queue.Push[[]byte] = pktime.
	ToPush(seedKey).
	WithLimit(plmt)

type SeedSender func(context.Context) error

func main() {
	var ptimeBeDefault PushTime = PushRaw(praw).ToPushTime(Time2bytesBeDefault)
	var senderBeDefault SeedSender = ptimeBeDefault.ToSeedSender(TimeSourceDefault)
	var ctx context.Context = context.Background()
	for {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		e := senderBeDefault(ctx)

		switch e {
		case queue.ErrTooMany:
			time.Sleep(sleep)
			continue
		case nil:
			time.Sleep(wait)
			continue
		default:
			log.Fatalf("Unexpected error: %v\n", e)
		}

		select {
		case <-time.After(wait):
			continue
		case <-ctx.Done():
			return
		}
	}
}
