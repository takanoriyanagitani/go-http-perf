package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	h := func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("dummy serialized request"))
	}
	http.HandleFunc("/", h)
	var server *http.Server = &http.Server{
		Addr:           ":6280",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	e := server.ListenAndServe()
	if nil != e {
		log.Fatal(e)
	}
}
