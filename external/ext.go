package external

import (
	"context"

	util "github.com/takanoriyanagitani/go-http-perf/util"
)

type SerializedRequestGenerator func(context.Context) (serialized []byte, e error)

type SerializedRequestSender func(ctx context.Context, serialized []byte) error

func (s SerializedRequestSender) ToSender(gen SerializedRequestGenerator) Sender {
	return func(ctx context.Context) error {
		serialized, e := gen(ctx)
		return util.Select(
			func() error { return nil },
			func() error { return s(ctx, serialized) },
			nil == e,
		)()
	}
}

type Sender func(ctx context.Context) error
