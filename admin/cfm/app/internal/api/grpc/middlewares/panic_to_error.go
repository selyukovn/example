package middlewares

import (
	"context"
	"google.golang.org/grpc"
	"runtime/debug"
)

// PanicToError
//
// Grpc-сервер, в отличие от net/http-сервера, не умеет сам перехватывать панику.
// Если панику не перехватить, все приложение рухнет, т.к. будет достигнута вершина стека горутины обработчика.
// Поскольку не все обработчики могут быть сломаны, имеет смысл панику перехватывать, чтоб не валить весь сервер.
func PanicToError(
	fnPanicToError func(ctx context.Context, panicValue any, debugStack []byte) error,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		rRes any,
		rErr error,
	) {
		defer func() {
			if pv := recover(); pv != nil {
				rRes = nil
				rErr = fnPanicToError(ctx, pv, debug.Stack())
			}
		}()

		return handler(ctx, req)
	}
}
