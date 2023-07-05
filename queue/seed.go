package queue

import (
	"context"
)

type PopSeedKey[K, T any] func(ctx context.Context, key K) (seed T, e error)

func (k PopSeedKey[K, T]) ToPopSeed(key K) PopSeed[T] {
	return func(ctx context.Context) (seed T, e error) {
		return k(ctx, key)
	}
}

type PopSeed[T any] func(ctx context.Context) (seed T, e error)
