package components

import (
	"context"
)

type LoggerInterface interface {
	// AddExtraAttrToCtx
	//
	// Паникует при нулевых аргументах.
	AddExtraAttrToCtx(ctx context.Context, key string, val string) context.Context

	// DebugFf -- см. fmt.Sprintf()
	//
	// Паникует при нулевых аргументах:
	// 	- ctx
	// 	- msg
	DebugFf(ctx context.Context, msg string, msgArgs ...any)

	// InfoFf -- см. fmt.Sprintf()
	//
	// Паникует при нулевых аргументах:
	// 	- ctx
	// 	- msg
	InfoFf(ctx context.Context, msg string, msgArgs ...any)

	// WarnFf -- см. fmt.Sprintf()
	//
	// Паникует при нулевых аргументах:
	// 	- ctx
	// 	- msg
	WarnFf(ctx context.Context, msg string, msgArgs ...any)

	// ErrorFf -- см. fmt.Sprintf()
	//
	// Паникует при нулевых аргументах:
	// 	- ctx
	// 	- msg
	ErrorFf(ctx context.Context, msg string, msgArgs ...any)
}
