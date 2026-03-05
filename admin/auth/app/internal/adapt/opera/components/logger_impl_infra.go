package components

import (
	"context"
	infra_logger "example/admin/auth/internal/infra/logger"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type LoggerImplInfraLogger struct {
	l *infra_logger.Logger
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewLoggerImplInfraLogger
//
// Паникует при нулевых аргументах.
func NewLoggerImplInfraLogger(l *infra_logger.Logger) *LoggerImplInfraLogger {
	assert.Cmp[*infra_logger.Logger]().NotEq(nil).Must(l)

	return &LoggerImplInfraLogger{l: l}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// AddExtraAttrToCtx
//
// Паникует при нулевых аргументах.
func (l *LoggerImplInfraLogger) AddExtraAttrToCtx(ctx context.Context, key string, val string) context.Context {
	return l.l.AddExtraAttrToCtx(ctx, key, val)
}

// DebugFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - ctx
//   - msg
func (l *LoggerImplInfraLogger) DebugFf(ctx context.Context, msg string, msgArgs ...any) {
	l.l.CtxDebugFf(ctx, msg, msgArgs...)
}

// InfoFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - ctx
//   - msg
func (l *LoggerImplInfraLogger) InfoFf(ctx context.Context, msg string, msgArgs ...any) {
	l.l.CtxInfoFf(ctx, msg, msgArgs...)
}

// WarnFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - ctx
//   - msg
func (l *LoggerImplInfraLogger) WarnFf(ctx context.Context, msg string, msgArgs ...any) {
	l.l.CtxWarnFf(ctx, msg, msgArgs...)
}

// ErrorFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - ctx
//   - msg
func (l *LoggerImplInfraLogger) ErrorFf(ctx context.Context, msg string, msgArgs ...any) {
	l.l.CtxErrorFf(ctx, msg, msgArgs...)
}
