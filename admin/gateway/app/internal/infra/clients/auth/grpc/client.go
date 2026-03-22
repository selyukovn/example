package grpc

import (
	"context"
	"example/admin/gateway/internal/infra/clients/auth"
	"example/admin/gateway/internal/infra/clients/auth/grpc/pb"
	"fmt"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/netip"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ auth.ClientInterface = &ClientImplGrpc{}

type ClientImplGrpc struct {
	pbClient         pb.AuthServiceClient
	apiKey           string
	fnGetOperationId func(context.Context) string
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewClientGrpc(
	baseUrl string,
	apiKey string,
	fnGetOperationId func(context.Context) string,
) (*ClientImplGrpc, error) {
	assert.Str().NotEmpty().Must(baseUrl)
	assert.Str().NotEmpty().Must(apiKey)
	assert.NotNilDeepMust(fnGetOperationId)

	grpcClient, err := grpc.NewClient(
		baseUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // todo : TLS
	)
	if err != nil {
		return nil, err
	}

	return &ClientImplGrpc{
		pbClient:         pb.NewAuthServiceClient(grpcClient),
		apiKey:           apiKey,
		fnGetOperationId: fnGetOperationId,
	}, nil
}

func NewClientGrpcMust(baseUrl string, apiKey string, fnGetOperationId func(context.Context) string) *ClientImplGrpc {
	res, err := NewClientGrpc(baseUrl, apiKey, fnGetOperationId)
	if err != nil {
		panic(err)
	}
	return res
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (c *ClientImplGrpc) prepareCtx(ctx context.Context) context.Context {
	mData := map[string]string{
		"authorization":  "Bearer " + c.apiKey,
		"x-operation-id": c.fnGetOperationId(ctx),
	}

	return metadata.NewOutgoingContext(ctx, metadata.New(mData))
}

// SignInRequest
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - auth.ErrorValidation
//   - auth.ErrorAccountAccessDenied
//   - std.ErrorRuntime
func (c *ClientImplGrpc) SignInRequest(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	email std.Email,
) (auth.SignInRequestResult, error) {
	assert.NotNilDeepMust(ctx)
	assert.TrueMust(fromIp != netip.Addr{}, "fromIp")
	assert.Str().NotEmpty().Must(fromUserAgent)
	assert.TrueMust(email != std.EmailNil, "email")

	ctx = c.prepareCtx(ctx)
	nilRes := auth.SignInRequestResult{}

	pbRes, pbErr := c.pbClient.SignInRequest(ctx, &pb.SignInRequestRequest{
		FromIp:        fromIp.String(),
		FromUserAgent: fromUserAgent,
		Email:         email.String(),
	})

	// error
	// ----------------

	if pbErr != nil {
		st, ok := status.FromError(pbErr)
		if !ok {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "SignInRequest")
		}

		if st.Code() == codes.InvalidArgument {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "SignInRequest", "InvalidArgument")
		} else if st.Code() == codes.NotFound {
			return nilRes, std.NewErrorNotFoundFf(pbErr.Error())
		} else if st.Code() == codes.FailedPrecondition {
			for _, detail := range st.Details() {
				switch vDet := detail.(type) {
				case *pb.ErrorValidationDetail:
					return nilRes, auth.ErrorValidation{
						Field:   vDet.Field,
						Message: vDet.Message,
					}
				case *pb.ErrorAccountAccessDeniedDetail:
					return nilRes, auth.ErrorAccountAccessDenied{}
				}
			}
		} else if st.Code() == codes.Internal {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "SignInRequest", "Internal")
		}

		panic(fmt.Errorf("%T.SignInRequest не знает, как обработать статус %v: %w", c, st.Code(), pbErr))
	}

	// success
	// ----------------

	canRetryAt := time.Time{}
	if pbRes.CanRetryAt.IsValid() {
		canRetryAt = pbRes.CanRetryAt.AsTime()
	}

	expireAt := time.Time{}
	if pbRes.ExpireAt.IsValid() {
		expireAt = pbRes.ExpireAt.AsTime()
	}

	return auth.SignInRequestResult{
		SignInId:    pbRes.SignInId,
		RetriesLeft: int(pbRes.RetriesLeft),
		CanRetryAt:  canRetryAt,
		ExpireAt:    expireAt,
	}, nil
}

// SignInRequestRetry
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - auth.ErrorValidation
//   - auth.ErrorAccountAccessDenied
//   - auth.ErrorSignInFinished
//   - auth.ErrorNoAttemptsLeft
//   - auth.ErrorRequestsFrequency
//   - std.ErrorUnprocessable -- сессия уже существует
//   - std.ErrorRuntime
func (c *ClientImplGrpc) SignInRequestRetry(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	signInId string,
) (auth.SignInRequestRetryResult, error) {
	assert.NotNilDeepMust(ctx)
	assert.TrueMust(fromIp != netip.Addr{})
	assert.Str().NotEmpty().Must(fromUserAgent)
	assert.Str().NotEmpty().Must(signInId)

	ctx = c.prepareCtx(ctx)
	nilRes := auth.SignInRequestRetryResult{}

	pbRes, pbErr := c.pbClient.SignInRequestRetry(ctx, &pb.SignInRequestRetryRequest{
		FromIp:        fromIp.String(),
		FromUserAgent: fromUserAgent,
		SignInId:      signInId,
	})

	// error
	// ----------------

	if pbErr != nil {
		st, ok := status.FromError(pbErr)
		if !ok {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "SignInRequestRetry")
		}

		if st.Code() == codes.InvalidArgument {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "SignInRequestRetry", "InvalidArgument")
		} else if st.Code() == codes.NotFound {
			return nilRes, std.NewErrorNotFoundFf(pbErr.Error())
		} else if st.Code() == codes.FailedPrecondition {
			for _, detail := range st.Details() {
				switch vDet := detail.(type) {
				case *pb.ErrorValidationDetail:
					return nilRes, auth.ErrorValidation{
						Field:   vDet.Field,
						Message: vDet.Message,
					}
				case *pb.ErrorAccountAccessDeniedDetail:
					return nilRes, auth.ErrorAccountAccessDenied{}
				case *pb.ErrorSignInFinishedDetail:
					return nilRes, auth.ErrorSignInFinished{
						IsPassed:  vDet.IsPassed,
						IsFailed:  vDet.IsFailed,
						IsExpired: vDet.IsExpired,
					}
				case *pb.ErrorNoAttemptsLeftDetail:
					return nilRes, auth.ErrorNoAttemptsLeft{}
				case *pb.ErrorRequestsFrequencyDetail:
					canReqAfter := time.Time{}
					if vDet.CanReqAfter.IsValid() {
						canReqAfter = vDet.CanReqAfter.AsTime()
					}

					return nilRes, auth.ErrorRequestsFrequency{
						CanReqAfter:        canReqAfter,
						CanReqAttemptsLeft: int(vDet.CanReqAttemptsLeft),
					}
				case *pb.ErrorUnprocessableDetail:
					return nilRes, std.NewErrorUnprocessableFf(pbErr.Error())
				}
			}
		} else if st.Code() == codes.Internal {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "SignInRequestRetry", "Internal")
		}

		panic(fmt.Errorf("%T.SignInRequestRetry не знает, как обработать статус %v: %w", c, st.Code(), pbErr))
	}

	// success
	// ----------------

	canRetryAt := time.Time{}
	if pbRes.CanRetryAt.IsValid() {
		canRetryAt = pbRes.CanRetryAt.AsTime()
	}

	return auth.SignInRequestRetryResult{
		SignInId:    pbRes.SignInId,
		RetriesLeft: int(pbRes.RetriesLeft),
		CanRetryAt:  canRetryAt,
	}, nil
}

// SignInConfirm
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - auth.ErrorValidation
//   - auth.ErrorAccountAccessDenied
//   - auth.ErrorSignInFinished
//   - std.ErrorUnprocessable
//   - std.ErrorRuntime
func (c *ClientImplGrpc) SignInConfirm(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	signInId string,
	code string,
) (auth.SignInConfirmResult, error) {
	assert.NotNilDeepMust(ctx)
	assert.TrueMust(fromIp != netip.Addr{})
	assert.Str().NotEmpty().Must(fromUserAgent)
	assert.Str().NotEmpty().Must(signInId)
	assert.Str().NotEmpty().Must(code)

	ctx = c.prepareCtx(ctx)
	nilRes := auth.SignInConfirmResult{}

	pbRes, pbErr := c.pbClient.SignInConfirm(ctx, &pb.SignInConfirmRequest{
		FromIp:        fromIp.String(),
		FromUserAgent: fromUserAgent,
		SignInId:      signInId,
		Code:          code,
	})

	// error
	// ----------------

	if pbErr != nil {
		st, ok := status.FromError(pbErr)
		if !ok {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "SignInConfirm")
		}

		if st.Code() == codes.InvalidArgument {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "SignInConfirm", "InvalidArgument")
		} else if st.Code() == codes.NotFound {
			return nilRes, std.NewErrorNotFoundFf(pbErr.Error())
		} else if st.Code() == codes.FailedPrecondition {
			for _, detail := range st.Details() {
				switch vDet := detail.(type) {
				case *pb.ErrorValidationDetail:
					return nilRes, auth.ErrorValidation{
						Field:   vDet.Field,
						Message: vDet.Message,
					}
				case *pb.ErrorAccountAccessDeniedDetail:
					return nilRes, auth.ErrorAccountAccessDenied{}
				case *pb.ErrorSignInFinishedDetail:
					return nilRes, auth.ErrorSignInFinished{
						IsPassed:  vDet.IsPassed,
						IsFailed:  vDet.IsFailed,
						IsExpired: vDet.IsExpired,
					}
				case *pb.ErrorUnprocessableDetail:
					return nilRes, std.NewErrorUnprocessableFf(pbErr.Error())
				}
			}
		} else if st.Code() == codes.Internal {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "SignInConfirm", "Internal")
		}

		panic(fmt.Errorf("%T.SignInConfirm не знает, как обработать статус %v: %w", c, st.Code(), pbErr))
	}

	// success
	// ----------------

	sessExpAt := time.Time{}
	if pbRes.SessionExpireAt.IsValid() {
		sessExpAt = pbRes.SessionExpireAt.AsTime()
	}

	return auth.SignInConfirmResult{
		IsPassed:        pbRes.IsPassed,
		AttemptsLeft:    int(pbRes.AttemptsLeft),
		SessionId:       pbRes.SessionId,
		SessionExpireAt: sessExpAt,
	}, nil
}

// SignOut
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - auth.ErrorValidation
//   - auth.ErrorAccountAccessDenied
//   - std.ErrorAlreadyDone
//   - std.ErrorRuntime
func (c *ClientImplGrpc) SignOut(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	sessionId string,
) error {
	assert.NotNilDeepMust(ctx)
	assert.TrueMust(fromIp != netip.Addr{})
	assert.Str().NotEmpty().Must(fromUserAgent)
	assert.Str().NotEmpty().Must(sessionId)

	ctx = c.prepareCtx(ctx)

	_, pbErr := c.pbClient.SignOut(ctx, &pb.SignOutRequest{
		FromIp:        fromIp.String(),
		FromUserAgent: fromUserAgent,
		SessionId:     sessionId,
	})

	// error
	// ----------------

	if pbErr != nil {
		st, ok := status.FromError(pbErr)
		if !ok {
			return std.WrapErrorToRuntime(pbErr, c, "SignOut")
		}

		if st.Code() == codes.InvalidArgument {
			return std.WrapErrorToRuntime(pbErr, c, "SignOut", "InvalidArgument")
		} else if st.Code() == codes.NotFound {
			return std.NewErrorNotFoundFf(pbErr.Error())
		} else if st.Code() == codes.FailedPrecondition {
			for _, detail := range st.Details() {
				switch vDet := detail.(type) {
				case *pb.ErrorValidationDetail:
					return auth.ErrorValidation{
						Field:   vDet.Field,
						Message: vDet.Message,
					}
				case *pb.ErrorAccountAccessDeniedDetail:
					return auth.ErrorAccountAccessDenied{}
				case *pb.ErrorAlreadyDoneDetail:
					return std.NewErrorAlreadyDoneFf(pbErr.Error())
				}
			}
		} else if st.Code() == codes.Internal {
			return std.WrapErrorToRuntime(pbErr, c, "SignOut", "Internal")
		}

		panic(fmt.Errorf("%T.SignOut не знает, как обработать статус %v: %w", c, st.Code(), pbErr))
	}

	return nil
}

// CheckSession
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - auth.ErrorValidation
//   - auth.ErrorAccountAccessDenied
//   - auth.ErrorSessionClosed
//   - std.ErrorRuntime
func (c *ClientImplGrpc) CheckSession(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	sessionId string,
) (auth.CheckSessionResult, error) {
	assert.NotNilDeepMust(ctx)
	assert.TrueMust(fromIp != netip.Addr{})
	assert.Str().NotEmpty().Must(fromUserAgent)
	assert.Str().NotEmpty().Must(sessionId)

	ctx = c.prepareCtx(ctx)
	nilRes := auth.CheckSessionResult{}

	pbRes, pbErr := c.pbClient.CheckSession(ctx, &pb.CheckSessionRequest{
		FromIp:        fromIp.String(),
		FromUserAgent: fromUserAgent,
		SessionId:     sessionId,
	})

	// error
	// ----------------

	if pbErr != nil {
		st, ok := status.FromError(pbErr)
		if !ok {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "CheckSession")
		}

		if st.Code() == codes.InvalidArgument {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "CheckSession", "InvalidArgument")
		} else if st.Code() == codes.NotFound {
			return nilRes, std.NewErrorNotFoundFf(pbErr.Error())
		} else if st.Code() == codes.FailedPrecondition {
			for _, detail := range st.Details() {
				switch vDet := detail.(type) {
				case *pb.ErrorValidationDetail:
					return nilRes, auth.ErrorValidation{
						Field:   vDet.Field,
						Message: vDet.Message,
					}
				case *pb.ErrorAccountAccessDeniedDetail:
					return nilRes, auth.ErrorAccountAccessDenied{}
				case *pb.ErrorSessionClosedDetail:
					return nilRes, auth.ErrorSessionClosed{
						IsExpired: vDet.IsExpired,
						ClosedAt:  vDet.ClosedAt.AsTime(),
					}
				}
			}
		} else if st.Code() == codes.Internal {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "CheckSession", "Internal")
		}

		panic(fmt.Errorf("%T.CheckSession не знает, как обработать статус %v: %w", c, st.Code(), pbErr))
	}

	// success
	return auth.CheckSessionResult{
		AccountId:       pbRes.AccountId,
		SessionExpireAt: pbRes.SessionExpireAt.AsTime(),
	}, nil
}
