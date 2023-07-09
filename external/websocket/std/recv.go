package wstd

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/net/websocket"

	util "github.com/takanoriyanagitani/go-http-perf/util"

	ws "github.com/takanoriyanagitani/go-http-perf/external/websocket"
)

func RecvNew(conn *websocket.Conn, buf []byte) ws.Recv {
	var codec websocket.Codec = websocket.Message
	return func(_ context.Context) (data []byte, e error) {
		e = codec.Receive(conn, &buf)
		return util.Select(
			func() ([]byte, error) {
				switch e {
				case io.EOF:
					return nil, e
				default:
					return nil, fmt.Errorf("Unable to recv: %w", e)
				}
			},
			func() ([]byte, error) { return buf, nil },
			nil == e,
		)()
	}
}
