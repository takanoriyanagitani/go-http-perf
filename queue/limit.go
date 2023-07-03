package queue

import (
	"context"
)

type PushLimit func(ctx context.Context) (tooMany bool, e error)
type PushLimitKey[K any] func(ctx context.Context, key K) (tooMany bool, e error)

func (l PushLimitKey[K]) ToPushLimit(key K) PushLimit {
	return func(ctx context.Context) (tooMany bool, e error) {
		return l(ctx, key)
	}
}
