package wstd

import (
	"context"

	"golang.org/x/net/websocket"

	ws "github.com/takanoriyanagitani/go-http-perf/external/websocket"
)

func SendNew(conn *websocket.Conn) ws.Send {
	var codec websocket.Codec = websocket.Message
	return func(_ context.Context, data []byte) error {
		return codec.Send(conn, data)
	}
}
