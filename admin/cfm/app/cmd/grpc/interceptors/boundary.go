package interceptors

import (
	"context"
	"example/admin/cfm/cmd/common/components/processing"
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/cmd/grpc/helpers"
	"github.com/google/uuid"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"strings"
)

// NewBoundary
//
// Данный перехватчик должен быть самым внешним!
func NewBoundary(ctr *container.Container) grpc.UnaryServerInterceptor {
	// Все описанные в данном перехватчике действия слишком связаны между собой,
	// чтобы выделить каждое в отдельный перехватчик.
	// Например, логирование статуса ответа не имеет смысла без operationId из обогащенного контекста,
	// а значит обогащение контекста обязано происходить до логирования статуса ответа.
	// Но статус ответа может быть изменен во внешних перехватчиках, что приведет к расхождению с уже записанными логами.
	// Кроме того, перехват паники для корректного ее логирования потребует ретрансляции и двух точек обработки,
	// что также может привести к нарушению согласованности кодов ответа.
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		rRes any,
		rErr error,
	) {
		rErrOnPanic := status.Error(codes.Internal, "Internal")

		// -------------------------------------------------------------------------------------------------------------
		// Резервный перехватчик паники.
		// -------------------------------------------------------------------------------------------------------------

		// Теоретически паника может возникнуть до основного перехватчика (например, при обогащении контекста),
		// поэтому обязательно нужен резервный перехватчик, чтобы не завалить весь сервер.

		defer func() {
			if pv := recover(); pv != nil {
				rRes = nil
				rErr = rErrOnPanic
			}
		}()

		// -------------------------------------------------------------------------------------------------------------
		// Обогащение
		// -------------------------------------------------------------------------------------------------------------

		operationId, ok := helpers.GrpcMetadataKeyFirst(ctx, "x-operation-id")
		if !ok || operationId == "" {
			return nil, helpers.ErrorInvalidArgument("x-operation-id header not found")
		}
		ctx = processing.EnrichCtx(ctx, operationId)
		ctx = logger.AddAttrToCtx(ctx, "processing.OperationId", operationId)
		/* см. */ _ = processing.OperationId

		requestId := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)
		ctx = helpers.AddGrpcInfoToCtx(ctx, requestId, info.FullMethod)
		ctx = logger.AddAttrToCtx(ctx, "helpers.RequestId", operationId)
		/* см. */ _ = helpers.RequestId

		// todo : trace... -- OpenTelemetry?

		// -------------------------------------------------------------------------------------------------------------
		// Логирование запроса
		// -------------------------------------------------------------------------------------------------------------

		logger.InfoFf(ctx, "Запрос: %q", helpers.GrpcFullMethod(ctx))

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

		// -------------------------------------------------------------------------------------------------------------
		// Основной перехватчик паники
		// -------------------------------------------------------------------------------------------------------------

		// Grpc-сервер, в отличие от net/http-сервера, не умеет сам перехватывать панику.
		// Если панику не перехватить, все приложение рухнет, т.к. будет достигнута вершина стека горутины обработчика.
		// Поскольку не все обработчики могут быть сломаны, имеет смысл панику перехватывать, чтоб не валить весь сервер.

		defer func() {
			if pv := recover(); pv != nil {
				rRes = nil
				rErr = rErrOnPanic

				logger.PanicFf(ctx, pv, debug.Stack(), "grpc.interceptors.NewBoundary")
			}
		}()

		// -------------------------------------------------------------------------------------------------------------

		return handler(ctx, req)
	}
}
