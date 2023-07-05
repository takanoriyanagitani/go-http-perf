package ws

import (
	"context"
)

type Send func(ctx context.Context, data []byte) error

type SeedSend func(ctx context.Context) error
