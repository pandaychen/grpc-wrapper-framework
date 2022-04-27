package rdb

//go-redis/v8的hook实现

import (
	"context"

	"github.com/pkg/errors"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tlog "github.com/opentracing/opentracing-go/log"
)

var _ redis.Hook = rdbTracingHook{}

type rdbTracingHook struct{}

//注意，与xorm的hook实现略微不同，这是需要返回child_ctx
func (rdbTracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	//返回值span也是一个interface{}
	span, child_ctx := opentracing.StartSpanFromContext(ctx, combineCommand(cmd))
	ext.Component.Set(span, "redis")
	ext.DBType.Set(span, "redis")
	ext.SpanKind.Set(span, "client")
	return child_ctx, nil
}

func (rdbTracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return errors.New("not found span")
	}
	defer span.Finish()
	if cmd.Err() != nil {
		//log error
		span.LogFields(tlog.Object("err", cmd.Err().Error()))
	}
	span.LogKV("cmd", cmd.String())
	return nil
}

//pipeline
func (rdbTracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	span, child_ctx := opentracing.StartSpanFromContext(ctx, combineCommand(cmds...))
	ext.Component.Set(span, "redis")
	ext.DBType.Set(span, "redis")
	ext.SpanKind.Set(span, "client")

	return child_ctx, nil
}

func (rdbTracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.Finish()
	}
	return nil
}
