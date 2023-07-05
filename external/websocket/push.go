package ws

import (
	"context"

	util "github.com/takanoriyanagitani/go-http-perf/util"
)

type PushRequest func(ctx context.Context, serialized []byte) error

func (p PushRequest) ToPushReceived(recv Recv) PushReceived {
	return func(ctx context.Context) error {
		serialized, e := recv(ctx)
		return util.Select(
			func() error { return e },
			func() error { return p(ctx, serialized) },
			nil == e,
		)()
	}
}

type PushReceived func(ctx context.Context) error
