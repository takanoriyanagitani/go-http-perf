package queue

import (
	"context"
	"errors"

	util "github.com/takanoriyanagitani/go-http-perf/util"
)

var ErrTooManyRequests error = errors.New("too many requests")

type PushRequestKey[K any] func(ctx context.Context, key K, serialized []byte) error

func (k PushRequestKey[K]) ToPushRequest(key K) PushRequest {
	return func(ctx context.Context, serialized []byte) error {
		return k(ctx, key, serialized)
	}
}

type PushRequest func(ctx context.Context, serialized []byte) error

func (r PushRequest) SkipIfTooMany(ctx context.Context, serialized []byte, tooMany bool) error {
	return util.Select(
		func() error { return r(ctx, serialized) },
		func() error { return ErrTooManyRequests },
		tooMany,
	)()
}

func (r PushRequest) WithLimit(pl PushLimit) PushRequest {
	return func(ctx context.Context, serialized []byte) error {
		tooMany, e := pl(ctx)
		return util.Select(
			func() error { return e },
			func() error { return r.SkipIfTooMany(ctx, serialized, tooMany) },
			nil == e,
		)()
	}
}
