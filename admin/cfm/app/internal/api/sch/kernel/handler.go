package kernel

import "context"

type Handler = func(ctx context.Context)
