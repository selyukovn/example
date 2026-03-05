package helpers

import (
	assert "github.com/selyukovn/go-wm-assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

// ---------------------------------------------------------------------------------------------------------------------
// Errors
// ---------------------------------------------------------------------------------------------------------------------

// ErrorInvalidArgument
//
// Паникует при нулевых аргументах.
func ErrorInvalidArgument(msg string) error {
	assert.Str().NotEmpty().Must(msg)

	return status.Error(codes.InvalidArgument, msg)
}

func ErrorUnauthenticated() error {
	return status.Error(codes.Unauthenticated, "Unauthenticated")
}

func ErrorPermissionDenied() error {
	return status.Error(codes.PermissionDenied, "PermissionDenied")
}

func ErrorNotFound() error {
	return status.Error(codes.NotFound, "NotFound")
}

// ErrorFailedPrecondition
//
// Паникует при нулевых аргументах.
func ErrorFailedPrecondition(details protoadapt.MessageV1) error {
	assert.NotNilDeepMust(details)

	st, sErr := status.New(codes.FailedPrecondition, "FailedPrecondition").WithDetails(details)
	if sErr != nil {
		return ErrorInternal()
	}
	return st.Err()
}

func ErrorInternal() error {
	return status.Error(codes.Internal, "Internal")
}

// ---------------------------------------------------------------------------------------------------------------------
