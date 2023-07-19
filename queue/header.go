package queue

import (
	"net/http"
)

type HeaderTransformer func(http.Header)

func (h HeaderTransformer) ToRequestEditor() RequestEditor {
	return func(req *http.Request) *http.Request {
		var hdr http.Header = req.Header
		h(hdr)
		return req
	}
}
