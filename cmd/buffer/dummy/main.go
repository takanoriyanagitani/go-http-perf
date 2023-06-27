package main

import (
	"time"

	buffer "github.com/takanoriyanagitani/go-http-perf/buffer"
)

func main() {
	var buf chan []byte = make(chan []byte, 128)

	var src buffer.SerializedRequestSource = buffer.EmptySource
	go src.StartSend(buf, 6*time.Millisecond)

	var dst buffer.SerializedRequestConsumer = buffer.NopConsumer
	dst.StartRecv(buf, 6*time.Millisecond)
}
