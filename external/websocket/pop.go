package ws

import (
	"context"

	util "github.com/takanoriyanagitani/go-http-perf/util"
)

type PopSeed[T any] func(ctx context.Context) (seed T, e error)

type PopSerializedSeed func(ctx context.Context) (serialized []byte, e error)

func (s PopSeed[T]) ToSerialized(s2b Seed2bytes[T]) PopSerializedSeed {
	return func(ctx context.Context) (serialized []byte, e error) {
		seed, e := s(ctx)
		return util.Select(
			func() ([]byte, error) { return nil, e },
			func() ([]byte, error) { return s2b(seed) },
			nil == e,
		)()
	}
}

func (p PopSerializedSeed) ToSeedSender(send Send) SeedSend {
	return func(ctx context.Context) error {
		seed, e := p(ctx)
		return util.Select(
			func() error { return e },
			func() error { return send(ctx, seed) },
			nil == e,
		)()
	}
}
