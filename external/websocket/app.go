package ws

import (
	"context"

	util "github.com/takanoriyanagitani/go-http-perf/util"
)

type App struct {
	SeedSend
	PushReceived
}

func (a App) SendRecv(ctx context.Context) error {
	es := a.SeedSend(ctx)
	return util.Select(
		func() error { return es },
		func() error { return a.PushReceived(ctx) },
		nil == es,
	)()
}
