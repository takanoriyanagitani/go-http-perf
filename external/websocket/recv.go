package ws

import (
	"context"
)

type Recv func(ctx context.Context) (data []byte, e error)
