package main

import (
	"context"
	"log"
	"time"

	gr "github.com/redis/go-redis/v9"

	queue "github.com/takanoriyanagitani/go-http-perf/queue"
	rq "github.com/takanoriyanagitani/go-http-perf/queue/external/redis/go-redis"
)

const pushKey string = "test-perf-queue"
const maxQ int64 = 256
const testMax int = 512

var client *gr.Client = gr.NewClient(&gr.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var qlk queue.QueueLengthKey[string] = rq.QueueLengthKeyNew(client)
var ql queue.QueueLength = qlk.ToQueueLength(pushKey)
var pl queue.PushLimit = ql.ToPushLimit(func(length int64) (tooMany bool) {
	return maxQ < length
})

var pk queue.PushKey[string, int64] = rq.PushKeyNew[int64](client)
var push queue.Push[int64] = pk.
	ToPush(pushKey).
	WithLimit(pl)

func main() {
	var t time.Time = time.Now()
	var us int64 = t.UnixMicro()

	for i := 0; i < testMax; i++ {
		e := push(context.Background(), us)
		if nil != e {
			log.Fatalf("Unable to push: %v\n", e)
		}
	}
}
