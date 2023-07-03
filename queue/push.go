package queue

import (
	"context"
	"errors"

	util "github.com/takanoriyanagitani/go-http-perf/util"
)

var ErrTooMany error = errors.New("too many requests")

type Push[T any] func(ctx context.Context, seed T) error
type PushKey[K, T any] func(ctx context.Context, key K, seed T) error

func (p Push[T]) SkipIfTooMany(ctx context.Context, seed T, tooMany bool) error {
	return util.Select(
		func() error { return p(ctx, seed) },
		func() error { return ErrTooMany },
		tooMany,
	)()
}

func (p Push[T]) WithLimit(limit PushLimit) Push[T] {
	return func(ctx context.Context, seed T) error {
		tooMany, e := limit(ctx)
		return util.Select(
			func() error { return e },
			func() error { return p.SkipIfTooMany(ctx, seed, tooMany) },
			nil == e,
		)()
	}
}

func (k PushKey[K, T]) ToPush(key K) Push[T] {
	return func(ctx context.Context, seed T) error {
		return k(ctx, key, seed)
	}
}
