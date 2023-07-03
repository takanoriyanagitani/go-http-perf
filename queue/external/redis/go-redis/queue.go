package grq

import (
	"context"

	gr "github.com/redis/go-redis/v9"

	queue "github.com/takanoriyanagitani/go-http-perf/queue"
)

func PushKeyNew[T any](cli *gr.Client) queue.PushKey[string, T] {
	return func(ctx context.Context, key string, seed T) error {
		var cmd *gr.IntCmd = cli.LPush(
			ctx,
			key,
			seed,
		)
		return cmd.Err()
	}
}

func QueueLengthKeyNew(cli *gr.Client) queue.QueueLengthKey[string] {
	return func(ctx context.Context, key string) (length int64, e error) {
		var cmd *gr.IntCmd = cli.LLen(ctx, key)
		return cmd.Result()
	}
}
