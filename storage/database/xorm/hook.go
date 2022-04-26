package xorm

import (
	"context"

	"grpc-wrapper-framework/common/enums"

	"github.com/opentracing/opentracing-go"
	tlog "github.com/opentracing/opentracing-go/log"
	"xorm.io/xorm/contexts"

	"xorm.io/builder"
)

//Make 编译器知道这个是 xorm 的 Hook，防止异常
var _ contexts.Hook = &XormHook{}

type XormHook struct {
	name string
}

func NewXormHook(name string) *XormHook {
	hook := &XormHook{
		name: name,
	}
	return hook
}

// 前置钩子实现
func (h *XormHook) BeforeProcess(ctx *contexts.ContextHook) (context.Context, error) {
	span, _ := opentracing.StartSpanFromContext(ctx.Ctx, "xorm-hook")

	// 将 span 注入 c.Ctx 中
	ctx.Ctx = context.WithValue(ctx.Ctx, xormHookSpanCtxKey, span)

	return ctx.Ctx, nil
}

func (h *XormHook) AfterProcess(c *contexts.ContextHook) error {
	sp, ok := c.Ctx.Value(xormHookSpanCtxKey).(opentracing.Span)
	if !ok {
		//no span,logger?
		return nil
	}
	//结束前上报
	defer sp.Finish()

	//log details
	if c.Err != nil {
		//log error
		sp.LogFields(tlog.Object("err", c.Err))
	}

	// 使用 xorm 的 builder 将查询语句和参数结合
	sql, err := builder.ConvertToBoundSQL(c.SQL, c.Args)
	if err == nil {
		// mark sql
		sp.LogFields(tlog.String(enums.TagDBStatement, sql))
	}
	sp.LogFields(tlog.String(enums.TagDBInstance, h.name))
	sp.LogFields(tlog.Object("args", c.Args))
	sp.SetTag(enums.TagDBExecuteCosts, c.ExecuteTime)

	return nil
}
