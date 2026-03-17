package logger

import (
	"context"
	"log/slog"
)

const slogTraceIdCtxKey = "logger.traceId"
const slogExtraAttrsCtxKey = "logger.extraAttrs"

type slogTraceIdCtxHandler struct {
	slog.Handler
}

func newSlogTraceIdCtxHandler(h slog.Handler) *slogTraceIdCtxHandler {
	return &slogTraceIdCtxHandler{Handler: h}
}

func slogTraceIdCtxAdd(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, slogTraceIdCtxKey, slog.String("trace_id", traceId))
}

func slogExtraAttrCtxAdd(ctx context.Context, key string, val string) context.Context {
	extraAttrs, ok := ctx.Value(slogExtraAttrsCtxKey).([]slog.Attr)
	if extraAttrs == nil || !ok {
		extraAttrs = make([]slog.Attr, 0)
	}

	extraAttrs = append(extraAttrs, slog.String(key, val))

	return context.WithValue(ctx, slogExtraAttrsCtxKey, extraAttrs)
}

func slogTraceIdCtxGet(ctx context.Context) (string, bool) {
	slogAttr := ctx.Value(slogTraceIdCtxKey)
	if slogAttr == nil {
		return "", false
	}

	v := slogAttr.(slog.Attr).Value.String()

	return v, v != ""
}

func (h *slogTraceIdCtxHandler) Handle(ctx context.Context, r slog.Record) error {
	if extraAttrs, ok := ctx.Value(slogExtraAttrsCtxKey).([]slog.Attr); ok {
		r.AddAttrs(extraAttrs...)
	}

	if traceIdAttr, ok := ctx.Value(slogTraceIdCtxKey).(slog.Attr); ok {
		r.AddAttrs(traceIdAttr)
	}

	return h.Handler.Handle(ctx, r)
}

func (h *slogTraceIdCtxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newSlogTraceIdCtxHandler(h.Handler.WithAttrs(attrs))
}

func (h *slogTraceIdCtxHandler) WithGroup(name string) slog.Handler {
	return newSlogTraceIdCtxHandler(h.Handler.WithGroup(name))
}
