package queue

import (
	"net/http"
	"net/url"
)

type UrlTransformer func(original *url.URL) (neo *url.URL)

func (u UrlTransformer) ToRequestEditor() RequestEditor {
	return func(original *http.Request) (edited *http.Request) {
		var ou *url.URL = original.URL
		original.URL = u(ou)
		return original
	}
}
