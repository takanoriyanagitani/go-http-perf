package perf

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"google.golang.org/protobuf/proto"

	preq "github.com/takanoriyanagitani/go-http-perf/http-perf/request"
	util "github.com/takanoriyanagitani/go-http-perf/util"
)

type RequestSource interface {
	Generate() (*preq.Request, error)
}

type Bytes2Request func([]byte) (*preq.Request, error)

func (b Bytes2Request) ToTime2Request(raw Time2RequestRaw) Time2Request {
	return util.ComposeErr(
		raw,
		b,
	)
}

func Bytes2RequestNewPB(buf *preq.Request) Bytes2Request {
	return func(serialized []byte) (*preq.Request, error) {
		buf.Reset()
		e := proto.Unmarshal(serialized, buf)
		return buf, e
	}
}

type Time2Request func(time.Time) (*preq.Request, error)
type Time2RequestRaw func(time.Time) (raw []byte, e error)

type UnixtimeMicros2RequestBuilder interface {
	Build() (UnixtimeMicros2Request, error)
}

type UnixtimeMicros2Request func(micros int64) (raw []byte, e error)

func (u UnixtimeMicros2Request) ToTime2RequestRaw() Time2RequestRaw {
	return func(t time.Time) (raw []byte, e error) {
		var micros int64 = t.UnixMicro()
		return u(micros)
	}
}

func (u UnixtimeMicros2Request) ToTime2Request(b Bytes2Request) Time2Request {
	var raw Time2RequestRaw = u.ToTime2RequestRaw()
	return util.ComposeErr(
		raw,
		b,
	)
}

func (g Time2Request) Generate() (*preq.Request, error) {
	var now time.Time = time.Now()
	return g(now)
}

func RequestToStd(request *preq.Request) (*http.Request, error) {
	var method string = request.Method
	var url string = request.Url
	var content []byte = request.Body
	var body io.Reader = bytes.NewReader(content)
	req, e := http.NewRequest(method, url, body)
	return util.Select(
		func() (*http.Request, error) { return nil, e },
		func() (*http.Request, error) {
			headerToStd(req.Header, request.Header)
			return req, nil
		},
		nil == e,
	)()
}

func headerToStd(dst http.Header, src map[string]*preq.HeaderContent) {
	for key, content := range src {
		var values []string = content.Values
		for _, val := range values {
			dst.Add(key, val)
		}
	}
}
