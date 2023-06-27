package buffer

import (
	"errors"
	"log"
	"time"

	util "github.com/takanoriyanagitani/go-http-perf/util"
)

var ErrTooManyRequests error = errors.New("too many requests")

type SerializedRequestSource func() (serialized []byte, e error)

var EmptySource SerializedRequestSource = func() ([]byte, error) { return nil, nil }

func (s SerializedRequestSource) SendToBuffer(buf chan<- []byte) error {
	var capacity int = cap(buf)
	var length int = len(buf)
	return util.Select(
		func() error { return ErrTooManyRequests },
		func() error {
			serialized, e := s()
			return util.Select(
				func() error { return e },
				func() error {
					buf <- serialized
					return nil
				},
				nil == e,
			)()
		},
		length < capacity,
	)()
}

func (s SerializedRequestSource) StartSend(buf chan<- []byte, wait time.Duration) {
	var ct <-chan time.Time = time.Tick(wait)
	for range ct {
		e := s.SendToBuffer(buf)
		util.Select(
			func() {
				log.Printf("Unable to send: %v\n", e)
			},
			func() {},
			nil == e,
		)()
	}
}
