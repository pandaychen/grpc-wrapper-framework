package xorm

type ctxKey string

type xormHookSpan struct{}

var (
	clientInstance     = ctxKey("_xorm_")
	xormHookSpanCtxKey = &xormHookSpan{}
)
