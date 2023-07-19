package queue

import (
	"bytes"
	"context"
	"io"
	"net/http"

	util "github.com/takanoriyanagitani/go-http-perf/util"
)

type Pop func(ctx context.Context) (serialized []byte, e error)
type PopKey[K any] func(ctx context.Context, key K) (serialized []byte, e error)

func (p Pop) ToRequestSource(m RequestMetaSource) RequestSource {
	return m.ToRequestSource(p)
}

func (k PopKey[K]) ToPop(key K) Pop {
	return func(ctx context.Context) (serialized []byte, e error) {
		return k(ctx, key)
	}
}

type RequestInfo struct {
	method string
	url    string
}

func RequestInfoEmptyNew(trustedMethod string) RequestInfo {
	return RequestInfo{
		method: trustedMethod,
		url:    "",
	}
}

var RequestInfoGet RequestInfo = RequestInfoEmptyNew(http.MethodGet)
var RequestInfoPost RequestInfo = RequestInfoEmptyNew(http.MethodPost)
var RequestInfoPut RequestInfo = RequestInfoEmptyNew(http.MethodPut)
var RequestInfoDelete RequestInfo = RequestInfoEmptyNew(http.MethodDelete)
var RequestInfoPatch RequestInfo = RequestInfoEmptyNew(http.MethodPatch)

func (i RequestInfo) WithUrl(trustedUrl string) RequestInfo {
	return RequestInfo{
		method: i.method,
		url:    trustedUrl,
	}
}

func (i RequestInfo) ToRequest(ctx context.Context, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, i.method, i.url, body)
}

type RequestMetaSource func() RequestInfo

type RequestSource func(ctx context.Context) (*http.Request, error)

func (m RequestMetaSource) ToRequestSource(pop Pop) RequestSource {
	return func(ctx context.Context) (*http.Request, error) {
		serialized, e := pop(ctx)
		return util.Select(
			func() (*http.Request, error) { return nil, e },
			func() (*http.Request, error) {
				info := m()
				var body io.Reader = bytes.NewReader(serialized)
				return info.ToRequest(ctx, body)
			},
			nil == e,
		)()
	}
}

func (s RequestSource) WithEditor(editor RequestEditor) RequestSource {
	return func(ctx context.Context) (edited *http.Request, e error) {
		original, e := s(ctx)
		return util.Select(
			func() (*http.Request, error) { return nil, e },
			func() (*http.Request, error) {
				edited = editor(original)
				return edited, nil
			},
			nil == e,
		)()
	}
}

func (s RequestSource) ToSender(rs RequestSender) Sender {
	return func(ctx context.Context) error {
		req, e := s(ctx)
		return util.Select(
			func() error { return nil },
			func() error { return rs(ctx, req) },
			nil == e,
		)()
	}
}
