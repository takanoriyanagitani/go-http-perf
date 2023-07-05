package grq

import (
	"context"
	"errors"

	gr "github.com/redis/go-redis/v9"

	queue "github.com/takanoriyanagitani/go-http-perf/queue"
	util "github.com/takanoriyanagitani/go-http-perf/util"
)

var ErrEmpty error = errors.New("no more seeds")

func PopSeedKeyNew(cli *gr.Client) queue.PopSeedKey[string, []byte] {
	return func(ctx context.Context, key string) (seed []byte, e error) {
		var cmd *gr.StringCmd = cli.RPop(ctx, key)
		seed, e = cmd.Bytes()
		return util.Select(
			func() ([]byte, error) { return seed, nil },
			func() ([]byte, error) { return nil, e },
			errors.Is(e, gr.ErrClosed),
		)()
	}
}
