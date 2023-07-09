package wstd

import (
	"context"
	"fmt"

	"golang.org/x/net/websocket"

	util "github.com/takanoriyanagitani/go-http-perf/util"

	ws "github.com/takanoriyanagitani/go-http-perf/external/websocket"
)

func SendNew(conn *websocket.Conn) ws.Send {
	var codec websocket.Codec = websocket.Message
	return func(_ context.Context, data []byte) error {
		var e error = codec.Send(conn, data)
		return util.Select(
			func() error { return fmt.Errorf("Unable to send: %w", e) },
			func() error { return nil },
			nil == e,
		)()
	}
}
