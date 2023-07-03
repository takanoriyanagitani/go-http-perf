package queue

import (
	"context"
)

type QueueLength func(ctx context.Context) (length int64, e error)
type QueueLengthKey[K any] func(ctx context.Context, key K) (length int64, e error)

func (k QueueLengthKey[K]) ToQueueLength(key K) QueueLength {
	return func(ctx context.Context) (length int64, e error) {
		return k(ctx, key)
	}
}

func (l QueueLength) ToPushLimit(checker func(length int64) (tooMany bool)) PushLimit {
	return func(ctx context.Context) (tooMany bool, e error) {
		length, e := l(ctx)
		tooMany = checker(length)
		return
	}
}
