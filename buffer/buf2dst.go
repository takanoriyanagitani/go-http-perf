package buffer

import (
	"errors"
	"log"
	"time"

	util "github.com/takanoriyanagitani/go-http-perf/util"
)

var ErrBufferEmpty error = errors.New("No requests available")

type SerializedRequestConsumer func(serialized []byte) error

var NopConsumer SerializedRequestConsumer = func(_ []byte) error { return nil }

func (c SerializedRequestConsumer) RecvFromBuffer(buf <-chan []byte) error {
	var length int = len(buf)
	return util.Select(
		func() error { return ErrBufferEmpty },
		func() error {
			var req []byte = <-buf
			return c(req)
		},
		0 < length,
	)()
}

func (c SerializedRequestConsumer) StartRecv(buf <-chan []byte, wait time.Duration) {
	var ct <-chan time.Time = time.Tick(wait)
	for range ct {
		e := c.RecvFromBuffer(buf)
		util.Select(
			func() {
				log.Printf("Unable to recv: %v\n", e)
			},
			func() {},
			nil == e,
		)()
	}
}
