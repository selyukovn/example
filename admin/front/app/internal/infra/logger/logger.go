package logger

import (
	"context"
	"fmt"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"io"
	"log/slog"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Logger struct {
	l *slog.Logger
	w io.Writer
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewLogger
//
// Паникует при нулевых аргументах:
//   - w
func NewLogger(w io.Writer, isDebugMode bool) *Logger {
	assert.NotNilDeepMust(w)

	var slogHandler slog.Handler

	slogHandler = slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: std.Ternary[slog.Level](isDebugMode, slog.LevelDebug, slog.LevelInfo),
	})

	slogHandler = newSlogTraceIdCtxHandler(slogHandler)

	return &Logger{
		l: slog.New(slogHandler),
		w: w,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (l *Logger) makePanicMessage(panicValue any, extraMsg string, extraMsgArgs ...any) string {
	extraInfo := ""
	if extraMsg != "" {
		extraInfo = fmt.Sprintf(extraMsg, extraMsgArgs...)
		extraInfo += ": "
	}

	switch pv := panicValue.(type) {
	case error:
		return fmt.Sprintf("ПАНИКА: %s%s", extraInfo, pv.Error())
	case string:
		return fmt.Sprintf("ПАНИКА: %s%s", extraInfo, pv)
	default:
		return fmt.Sprintf("ПАНИКА: %s%#v", extraInfo, pv)
	}
}

func (l *Logger) makeStackTraceSlogAttr(stackTrace []byte) slog.Attr {
	return slog.String("stack", string(stackTrace))
}

// General
// ---------------------------------------------------------------------------------------------------------------------

// GeneralDebugFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - msg
func (l *Logger) GeneralDebugFf(msg string, msgArgs ...any) {
	assert.Str().NotEmpty().Must(msg)

	l.l.Debug(fmt.Sprintf(msg, msgArgs...))
}

// GeneralInfoFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - msg
func (l *Logger) GeneralInfoFf(msg string, msgArgs ...any) {
	assert.Str().NotEmpty().Must(msg)

	l.l.Info(fmt.Sprintf(msg, msgArgs...))
}

// GeneralWarnFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - msg
func (l *Logger) GeneralWarnFf(msg string, msgArgs ...any) {
	assert.Str().NotEmpty().Must(msg)

	l.l.Warn(fmt.Sprintf(msg, msgArgs...))
}

// GeneralErrorFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - msg
func (l *Logger) GeneralErrorFf(msg string, msgArgs ...any) {
	assert.Str().NotEmpty().Must(msg)

	l.l.Error(fmt.Sprintf(msg, msgArgs...))
}

// GeneralPanicFf
//
// Паникует при нулевых аргументах:
//   - debugStack
//   - msg
func (l *Logger) GeneralPanicFf(panicValue any, debugStack []byte, msg string, msgArgs ...any) {
	assert.NotNilDeepMust(debugStack)
	assert.Str().NotEmpty().Must(msg)

	l.l.Error(
		l.makePanicMessage(panicValue, msg, msgArgs...),
		l.makeStackTraceSlogAttr(debugStack),
	)
}

// Ctx
// ---------------------------------------------------------------------------------------------------------------------

// AddTraceIdToCtx
//
// Паникует при нулевых аргументах.
func (l *Logger) AddTraceIdToCtx(ctx context.Context, traceId string) context.Context {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Str().NotEmpty().Must(traceId)

	return slogTraceIdCtxAdd(ctx, traceId)
}

// AddExtraAttrToCtx
//
// Паникует при нулевых аргументах.
func (l *Logger) AddExtraAttrToCtx(ctx context.Context, key string, val string) context.Context {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Str().NotEmpty().Must(key)
	assert.Str().NotEmpty().Must(val)

	return slogExtraAttrCtxAdd(ctx, key, val)
}

// GetTraceIdFromCtx
//
// Паникует при нулевых аргументах.
func (l *Logger) GetTraceIdFromCtx(ctx context.Context) (string, bool) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)

	return slogTraceIdCtxGet(ctx)
}

// CtxDebugFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - ctx
//   - msg
func (l *Logger) CtxDebugFf(ctx context.Context, msg string, msgArgs ...any) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Str().NotEmpty().Must(msg)

	l.l.DebugContext(ctx, fmt.Sprintf(msg, msgArgs...))
}

// CtxInfoFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - ctx
//   - msg
func (l *Logger) CtxInfoFf(ctx context.Context, msg string, msgArgs ...any) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Str().NotEmpty().Must(msg)

	l.l.InfoContext(ctx, fmt.Sprintf(msg, msgArgs...))
}

// CtxWarnFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - ctx
//   - msg
func (l *Logger) CtxWarnFf(ctx context.Context, msg string, msgArgs ...any) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Str().NotEmpty().Must(msg)

	l.l.WarnContext(ctx, fmt.Sprintf(msg, msgArgs...))
}

// CtxErrorFf -- см. fmt.Sprintf()
//
// Паникует при нулевых аргументах:
//   - ctx
//   - msg
func (l *Logger) CtxErrorFf(ctx context.Context, msg string, msgArgs ...any) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Str().NotEmpty().Must(msg)

	l.l.ErrorContext(ctx, fmt.Sprintf(msg, msgArgs...))
}

// CtxPanicFf
//
// Паникует при нулевых аргументах:
//   - ctx
//   - debugStack
//   - msg
func (l *Logger) CtxPanicFf(ctx context.Context, panicValue any, debugStack []byte, msg string, msgArgs ...any) {
	assert.NotNilDeepMust(ctx)
	assert.NotNilDeepMust(debugStack)
	assert.Str().NotEmpty().Must(msg)

	l.l.ErrorContext(
		ctx,
		l.makePanicMessage(panicValue, msg, msgArgs...),
		l.makeStackTraceSlogAttr(debugStack),
	)
}
