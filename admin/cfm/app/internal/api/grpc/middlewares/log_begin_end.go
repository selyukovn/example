package middlewares

import (
	"context"
	"example/admin/cfm/internal/api/grpc/kernel"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LogBeginEnd() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		rRes any,
		rErr error,
	) {
		ctx = logger.AddAttrToCtx(ctx, "kernel.RequestId", kernel.RequestId(ctx))

		logger.InfoFf(ctx, "Запрос: %q", kernel.FullMethod(ctx))

		defer func() {
			code := codes.OK
			msg := "ok"
			if rErr != nil {
				st := status.Convert(rErr)
				code = st.Code()
				msg = st.Message()
			}

			logger.InfoFf(ctx, "Ответ: %d %s", code, msg)
		}()

		return handler(ctx, req)
	}
}
