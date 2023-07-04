package queue

import (
	"context"
	"net/http"

	util "github.com/takanoriyanagitani/go-http-perf/util"
)

type SerializedRequestParser func(serialized []byte) (*http.Request, error)

type RequestSender func(context.Context, *http.Request) error

type SerializedRequestSender func(ctx context.Context, request []byte) error

type Sender func(context.Context) error

func (s RequestSender) ToSerializedRequestSender(de SerializedRequestParser) SerializedRequestSender {
	return func(ctx context.Context, serialized []byte) error {
		req, e := de(serialized)
		return util.Select(
			func() error { return e },
			func() error { return s(ctx, req) },
			nil == e,
		)()
	}
}

func (s SerializedRequestSender) ToSender(pop Pop) Sender {
	return func(ctx context.Context) error {
		serialized, e := pop(ctx)
		return util.Select(
			func() error { return e },
			func() error { return s(ctx, serialized) },
			nil == e,
		)()
	}
}
