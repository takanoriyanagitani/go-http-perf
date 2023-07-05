package wstd

import (
	"context"

	"golang.org/x/net/websocket"

	ws "github.com/takanoriyanagitani/go-http-perf/external/websocket"
)

func RecvNew(conn *websocket.Conn, buf []byte) ws.Recv {
	var codec websocket.Codec = websocket.Message
	return func(_ context.Context) (data []byte, e error) {
		return buf, codec.Receive(conn, &buf)
	}
}
