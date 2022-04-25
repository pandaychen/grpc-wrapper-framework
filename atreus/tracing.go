package atreus

import (
	"context"

	"grpc-wrapper-framework/atreus/tracers"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//用于客户端及服务端的tracing拦截器（jaeger）
func (c *Client) OpenTracingForClient() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		var (
			//https://pkg.go.dev/github.com/opentracing/opentracing-go#SpanContext
			parentCtx opentracing.SpanContext
		)
		//先从context中获取原始的span，可能获取不到
		parentSpan := opentracing.SpanFromContext(ctx)
		if parentSpan != nil {
			parentCtx = parentSpan.Context()
		}

		//parentSpan可能为nil
		span := c.tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
			ext.SpanKindRPCClient,
		)
		defer span.Finish()

		////从客户端context中获取metadata。md.(type) == map[string][]string（标签）
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		mdWriter := tracers.MDReaderWriter{md}
		//将span的context信息注入到carrier中
		err = c.tracer.Inject(span.Context(), opentracing.TextMap, mdWriter)
		if err != nil {
			span.LogFields(log.String("inject-err", err.Error()))
		}

		//创建一个新的context，把metadata附带上
		newCtx := metadata.NewOutgoingContext(ctx, md)
		err = invoker(newCtx, method, req, reply, cc, opts...)
		if err != nil {
			//记录错误日志
			span.LogFields(log.String("caller-err", err.Error()))
		}
		return err
	}
}
