package queue

import (
	"net/http"
)

type RequestEditor func(original *http.Request) (edited *http.Request)

func (r RequestEditor) Join(other RequestEditor) RequestEditor {
	return func(original *http.Request) (edited *http.Request) {
		var e1 *http.Request = r(original)
		var e2 *http.Request = other(e1)
		return e2
	}
}
