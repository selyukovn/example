package cfm

import (
	"context"
	"example/admin/auth/internal/infra/clients/cfm"
	"example/admin/auth/internal/infra/clients/cfm/grpc/pb"
	"fmt"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var ClientGrpcNil = ClientGrpc{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ cfm.ClientInterface = ClientGrpc{}

type ClientGrpc struct {
	pbClient         pb.CfmServiceClient
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
) (ClientGrpc, error) {
	assert.Str().NotEmpty().Must(baseUrl)
	assert.Str().NotEmpty().Must(apiKey)
	assert.NotNilDeepMust(fnGetOperationId)

	grpcClient, err := grpc.NewClient(
		baseUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // todo : TLS
	)
	if err != nil {
		return ClientGrpcNil, err
	}

	return ClientGrpc{
		pbClient:         pb.NewCfmServiceClient(grpcClient),
		apiKey:           apiKey,
		fnGetOperationId: fnGetOperationId,
	}, nil
}

func NewClientGrpcMust(
	baseUrl string,
	apiKey string,
	fnGetOperationId func(context.Context) string,
) ClientGrpc {
	res, err := NewClientGrpc(baseUrl, apiKey, fnGetOperationId)
	if err != nil {
		panic(err)
	}
	return res
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (c ClientGrpc) prepareCtx(ctx context.Context) context.Context {
	mData := map[string]string{
		"authorization":  "Bearer " + c.apiKey,
		"x-operation-id": c.fnGetOperationId(ctx),
	}

	return metadata.NewOutgoingContext(ctx, metadata.New(mData))
}

// CreateForEmail
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (c ClientGrpc) CreateForEmail(ctx context.Context, email std.Email) (cfm.CreateForEmailResult, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(email.IsNil())

	ctx = c.prepareCtx(ctx)
	nilRes := cfm.CreateForEmailResult{}

	pbRes, pbErr := c.pbClient.CreateForEmail(ctx, &pb.CreateForEmailRequest{
		Email: email.String(),
	})

	// error
	// ----------------

	if pbErr != nil {
		return nilRes, std.WrapErrorToRuntime(pbErr, c, "CreateForEmail")
	}

	// success
	// ----------------

	if err := assert.Str().Eq(email.String()).Check(pbRes.Email); err != nil {
		return nilRes, std.WrapErrorToRuntime(pbErr, c, "CreateForEmail", "Email")
	}

	cfmId, err := pbRes.CfmId, assert.Str().NotEmpty().Check(pbRes.CfmId)
	if err != nil {
		return nilRes, std.WrapErrorToRuntime(err, c, "CreateForEmail", "CfmId")
	}

	expireAt, err := pbRes.ExpireAt.AsTime(), assert.Bool().True().Check(pbRes.ExpireAt.IsValid())
	if err != nil {
		return nilRes, std.WrapErrorToRuntime(err, c, "CreateForEmail", "ExpireAt")
	}

	return cfm.CreateForEmailResult{
		CfmId:    cfmId,
		ExpireAt: expireAt,
	}, nil
}

// Request
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - cfm.ErrorFinished
//   - cfm.ErrorNoAttemptsLeft
//   - cfm.ErrorRequestsFrequency
//   - std.ErrorRuntime
func (c ClientGrpc) Request(ctx context.Context, cfmId string) (cfm.RequestResult, error) {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(cfmId)

	ctx = c.prepareCtx(ctx)
	nilRes := cfm.RequestResult{}

	pbRes, pbErr := c.pbClient.Request(ctx, &pb.RequestRequest{
		CfmId: cfmId,
	})

	// error
	// ----------------

	if pbErr != nil {
		st, ok := status.FromError(pbErr)
		if !ok {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "Request")
		}

		if st.Code() == codes.InvalidArgument {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "Request", "InvalidArgument")
		} else if st.Code() == codes.NotFound {
			return nilRes, std.NewErrorNotFoundFf(pbErr.Error())
		} else if st.Code() == codes.FailedPrecondition {
			for _, detail := range st.Details() {
				switch vDet := detail.(type) {
				case *pb.ErrorFinishedDetail:
					if err := assert.Str().Eq(cfmId).Check(vDet.CfmId); err != nil {
						return nilRes, std.WrapErrorToRuntime(err, c, "Request", fmt.Sprintf("%T", vDet), "CfmId")
					}

					finishedAt, err := vDet.FinishedAt.AsTime(), assert.Bool().True().Check(vDet.FinishedAt.IsValid())
					if err != nil {
						return nilRes, std.WrapErrorToRuntime(err, c, "Request", fmt.Sprintf("%T", vDet), "FinishedAt")
					}

					return nilRes, cfm.ErrorFinished{
						CfmId:       vDet.CfmId,
						FinishedAt:  finishedAt,
						IsAsPassed:  vDet.IsAsPassed,
						IsAsFailed:  vDet.IsAsFailed,
						IsAsExpired: vDet.IsAsExpired,
					}
				case *pb.ErrorNoAttemptsLeftDetail:
					if err := assert.Str().Eq(cfmId).Check(vDet.CfmId); err != nil {
						return nilRes, std.WrapErrorToRuntime(err, c, "Request", fmt.Sprintf("%T", vDet), "CfmId")
					}
					return nilRes, cfm.ErrorNoAttemptsLeft{
						CfmId: vDet.CfmId,
					}
				case *pb.ErrorRequestsFrequencyDetail:
					if err := assert.Str().Eq(cfmId).Check(vDet.CfmId); err != nil {
						return nilRes, std.WrapErrorToRuntime(err, c, "Request", fmt.Sprintf("%T", vDet), "CfmId")
					}

					canReqAfter, err := vDet.CanReqAfter.AsTime(), assert.Bool().True().Check(vDet.CanReqAfter.IsValid())
					if err != nil {
						return nilRes, std.WrapErrorToRuntime(err, c, "Request", fmt.Sprintf("%T", vDet), "CanReqAfter")
					}

					canReqAttemptsLeft := uint(vDet.CanReqAttemptsLeft)

					return nilRes, cfm.ErrorRequestsFrequency{
						CfmId:              vDet.CfmId,
						CanReqAfter:        canReqAfter,
						CanReqAttemptsLeft: canReqAttemptsLeft,
					}
				}
			}
		} else if st.Code() == codes.Internal {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "Request", "Internal")
		}

		panic(fmt.Errorf("%T.Request не знает, как обработать статус %v: %w", c, st.Code(), pbErr))
	}

	// success
	// ----------------

	if err := assert.Str().Eq(cfmId).Check(pbRes.CfmId); err != nil {
		return nilRes, std.WrapErrorToRuntime(err, c, "Request", "success", "CfmId")
	}

	canReqAfter := time.Time{}
	if pbRes.CanReqAfter.IsValid() {
		canReqAfter = pbRes.CanReqAfter.AsTime()
	}

	return cfm.RequestResult{
		CfmId:              pbRes.CfmId,
		CanReqAgain:        pbRes.CanReqAgain,
		CanReqAttemptsLeft: uint(pbRes.CanReqAttemptsLeft),
		CanReqAfter:        canReqAfter,
	}, nil
}

// Confirm
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - cfm.ErrorFinished
//   - std.ErrorUnprocessable -- если не была запрошена
//   - std.ErrorRuntime
func (c ClientGrpc) Confirm(ctx context.Context, cfmId string, code string) (cfm.ConfirmResult, error) {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(cfmId)
	assert.Str().NotEmpty().Must(code)

	ctx = c.prepareCtx(ctx)
	nilRes := cfm.ConfirmResult{}

	pbRes, pbErr := c.pbClient.Confirm(ctx, &pb.ConfirmRequest{
		CfmId: cfmId,
		Code:  code,
	})

	// error
	// ----------------

	if pbErr != nil {
		st, ok := status.FromError(pbErr)
		if !ok {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "Confirm")
		}

		if st.Code() == codes.InvalidArgument {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "Confirm", "InvalidArgument")
		} else if st.Code() == codes.NotFound {
			return nilRes, std.NewErrorNotFoundFf(pbErr.Error())
		} else if st.Code() == codes.FailedPrecondition {
			for _, detail := range st.Details() {
				switch vDet := detail.(type) {
				case *pb.ErrorFinishedDetail:
					if err := assert.Str().Eq(cfmId).Check(vDet.CfmId); err != nil {
						return nilRes, std.WrapErrorToRuntime(err, c, "Confirm", fmt.Sprintf("%T", vDet), "CfmId")
					}

					finishedAt, err := vDet.FinishedAt.AsTime(), assert.Bool().True().Check(vDet.FinishedAt.IsValid())
					if err != nil {
						return nilRes, std.WrapErrorToRuntime(err, c, "Confirm", fmt.Sprintf("%T", vDet), "FinishedAt")
					}

					return nilRes, cfm.ErrorFinished{
						CfmId:       vDet.CfmId,
						FinishedAt:  finishedAt,
						IsAsPassed:  vDet.IsAsPassed,
						IsAsFailed:  vDet.IsAsFailed,
						IsAsExpired: vDet.IsAsExpired,
					}
				case *pb.ErrorNotRequestedDetail:
					if err := assert.Str().Eq(cfmId).Check(vDet.CfmId); err != nil {
						return nilRes, std.WrapErrorToRuntime(err, c, "Confirm", fmt.Sprintf("%T", vDet), "CfmId")
					}

					return nilRes, std.NewErrorUnprocessableFf(pbErr.Error())
				}
			}
		} else if st.Code() == codes.Internal {
			return nilRes, std.WrapErrorToRuntime(pbErr, c, "Request", "Internal")
		}

		panic(fmt.Errorf("%T.Confirm не знает, как обработать статус %v: %w", c, st.Code(), pbErr))
	}

	// success
	// ----------------

	if err := assert.Str().Eq(cfmId).Check(pbRes.CfmId); err != nil {
		return nilRes, std.WrapErrorToRuntime(err, c, "Confirm", "success", "CfmId")
	}

	finishedAt := time.Time{}
	if pbRes.FinishedAt.IsValid() {
		finishedAt = pbRes.FinishedAt.AsTime()
	}

	return cfm.ConfirmResult{
		CfmId:              cfmId,
		FinishedAt:         finishedAt,
		IsFinishedAsPassed: pbRes.IsFinishedAsPassed,
		IsFinishedAsFailed: pbRes.IsFinishedAsFailed,
		FailsLeft:          uint(pbRes.FailsLeft),
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
