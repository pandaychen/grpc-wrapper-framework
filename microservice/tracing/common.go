package tracing

import "context"

// tracing 公共定义

// Trace key
const (
	// trace 键
	GlobalTraceID = "atreus-trace-id"
)

type ctxKey string //context uniq key

var _ctxkey ctxKey = "atreus/tracing" //key obj

// 从context.Context中获取t Trace结构（大结构，所有的span都是存储在t Trace中）
func FromContext(ctx context.Context) (t SpanTrace, ok bool) {
	t, ok = ctx.Value(_ctxkey).(SpanTrace)
	return
}

// 将Trace存储在context中
func NewContext(ctx context.Context, t SpanTrace) context.Context {
	return context.WithValue(ctx, _ctxkey, t)
}
