package grq

import (
	"context"

	gr "github.com/redis/go-redis/v9"

	queue "github.com/takanoriyanagitani/go-http-perf/queue"
)

func PushRequestKeyNew(cli *gr.Client) queue.PushRequestKey[string] {
	return func(ctx context.Context, key string, serialized []byte) error {
		var cmd *gr.IntCmd = cli.LPush(ctx, key, serialized)
		return cmd.Err()
	}
}
