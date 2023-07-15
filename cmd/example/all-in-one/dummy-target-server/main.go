package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	var server *http.Server = &http.Server{
		Addr:           ":7168",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(""))
	})

	log.Fatal(server.ListenAndServe())
}
