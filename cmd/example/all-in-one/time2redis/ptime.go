package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"time"

	queue "github.com/takanoriyanagitani/go-http-perf/queue"
	util "github.com/takanoriyanagitani/go-http-perf/util"
)

type PushTime queue.Push[time.Time]
type PushRaw queue.Push[[]byte]

type Time2us func(time.Time) (micros int64)
type Int2bytes func(int64) []byte

var Time2usDefault Time2us = func(t time.Time) (micros int64) { return t.UnixMicro() }

func Int2bytesNew(bo binary.ByteOrder, buf *bytes.Buffer) Int2bytes {
	return func(i int64) (serizlied []byte) {
		buf.Reset()
		e := binary.Write(buf, bo, i)
		if nil != e {
			panic(e)
		}
		return buf.Bytes()
	}
}

var Int2bytesBeDefault Int2bytes = Int2bytesNew(binary.BigEndian, new(bytes.Buffer))

func (t2i Time2us) ToTime2bytes(i2b Int2bytes) Time2bytes {
	return util.Compose(
		t2i,
		i2b,
	)
}

var Time2bytesBeDefault Time2bytes = Time2usDefault.ToTime2bytes(Int2bytesBeDefault)

type Time2bytes func(time.Time) []byte

func (r PushRaw) ToPushTime(t2b Time2bytes) PushTime {
	return func(ctx context.Context, tim time.Time) error {
		var serialized []byte = t2b(tim)
		return r(ctx, serialized)
	}
}

type TimeSource func() time.Time

var TimeSourceDefault TimeSource = time.Now

func (p PushTime) ToSeedSender(ts TimeSource) SeedSender {
	return func(ctx context.Context) error {
		var seed time.Time = ts()
		return p(ctx, seed)
	}
}
