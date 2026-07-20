package kernel

type Middleware = func(Handler) Handler
