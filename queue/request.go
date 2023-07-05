package queue

import (
	"context"
)

type PushRequestKey[K any] func(ctx context.Context, key K, serialized []byte) error

func (k PushRequestKey[K]) ToPushRequest(key K) PushRequest {
	return func(ctx context.Context, serialized []byte) error {
		return k(ctx, key, serialized)
	}
}

type PushRequest func(ctx context.Context, serialized []byte) error
