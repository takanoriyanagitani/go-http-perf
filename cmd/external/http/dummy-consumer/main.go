package main

import (
	"context"
	"log"

	ext "github.com/takanoriyanagitani/go-http-perf/external"
	eh "github.com/takanoriyanagitani/go-http-perf/external/http"
)

const url string = "http://localhost:6280"

var builder eh.SerializedRequestGeneratorBuilder = eh.
	DefaultBuilder.
	WithUrl(eh.UrlNewFromRawString(url))
var gen ext.SerializedRequestGenerator = builder.Build()
var snd ext.SerializedRequestSender = func(_ context.Context, serialized []byte) error {
	log.Printf("serialized request(dummy): %s\n", string(serialized))
	return nil
}
var logSender ext.Sender = snd.ToSender(gen)

func main() {
	e := logSender(context.Background())
	if nil != e {
		log.Fatalf("Unexpected error: %v\n", e)
	}
}
