package queue

import (
	"context"
)

type Pop func(ctx context.Context) (serialized []byte, e error)
type PopKey[K any] func(ctx context.Context, key K) (serialized []byte, e error)

func (k PopKey[K]) ToPop(key K) Pop {
	return func(ctx context.Context) (serialized []byte, e error) {
		return k(ctx, key)
	}
}
